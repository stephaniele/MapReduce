package mr

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"sort"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {

	// 1. ask for task
	for {
		reply, ok := callMasterForTask(-1, false)
		if !ok || reply.AllTasksAreDone {
			break
		}
		task := reply.TodoTask
		var err error
		switch task.TaskType {
		case IsMap:
			err = doMapTask(mapf, task.InputFile, reply.NReduce, task.TaskIndex)
		case IsReduce:
			err = doReduceTask(reducef, reply.NMap, task.TaskIndex)
		default:
			fmt.Printf("[DEBUG][WORKER] unknown task type %v \n", task.TaskType)
		}

		if err != nil {
			callMasterForTask(task.TaskIndex, false)
		} else {
			callMasterForTask(task.TaskIndex, true)
		}
	}
}

func doReduceTask(reducef func(string, []string) string, nMap int, taskIndex int) error {
	fmt.Println("[INFO][WORKER] start doing reduce task...")
	kvMap := make(map[string][]string)
	for i := 0; i < nMap; i++ {
		fileName := getIntermediateFileName(i, taskIndex)
		ofile, err := os.Open(fileName)
		if err != nil {
			log.Fatalf("cannot open %v", fileName)
			return err
		}
		decoder := json.NewDecoder(ofile)
		for decoder.More() {
			var kv KeyValue
			err := decoder.Decode(&kv)
			if err != nil {
				log.Fatalf("error decoding %v", err)
				return err
			}
			_, ok := kvMap[kv.Key]
			if !ok {
				kvMap[kv.Key] = make([]string, 0)
			}
			kvMap[kv.Key] = append(kvMap[kv.Key], kv.Value)
		}
		ofile.Close()
	}
	//========= sort key =============
	var keys []string
	for k := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	//========= write frequency to output =============
	fileName := fmt.Sprintf("mr-out-%d", taskIndex)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("cannot open %v", fileName)
		return err
	}
	for k, v := range kvMap {
		reduceOutput := reducef(k, v)
		fmt.Fprintf(file, "%v %v\n", k, reduceOutput)
	}
	file.Close()
	return nil
}

func doMapTask(mapf func(string, string) []KeyValue, filename string, nReduce int, taskIndex int) error {
	fmt.Println("[INFO][WORKER] start doing map task...")
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("cannot open %v", filename)
		return err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("cannot read %v", filename)
		return err
	}
	file.Close()
	kva := mapf(filename, string(content))

	//======== Create intermediate =============
	keyToReduceTasks := make([][]KeyValue, nReduce)
	for _, kv := range kva {
		reduceIndex := ihash(kv.Key) % nReduce
		keyToReduceTasks[reduceIndex] = append(keyToReduceTasks[reduceIndex], kv)
	}
	for reduceIndex, file := range keyToReduceTasks {
		fileName := getIntermediateFileName(taskIndex, reduceIndex)
		ofile, _ := os.Create(fileName)
		encoder := json.NewEncoder(ofile)
		for _, kv := range file {
			if err := encoder.Encode(&kv); err != nil {
				log.Fatalf("cannot write %v", err)
				return err
			}
		}
		ofile.Close()
	}
	return nil
}

func getIntermediateFileName(mapIndex int, reduceIndex int) string {
	return fmt.Sprintf("mr-%d-%d", mapIndex, reduceIndex)
}

func callMasterForTask(taskIndex int, isDone bool) (reply Reply, ok bool) {

	args := Args{
		TaskIndex: taskIndex,
		Finished:  isDone,
	}
	reply = Reply{}

	// send the RPC request, wait for the reply.
	ok = call("Master.WorkerRequestHandler", &args, &reply)

	fmt.Printf("[INFO][WORKER Reply]task.todoTask %v, is %v \n", reply.TodoTask.InputFile, ok)
	return
}

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
