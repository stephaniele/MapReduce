package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
	"unicode"
)

type taskState int

const maxTimeout = 10

type ChunkInfo struct {
	OffsetStart int64
	ChunkSize   int
}

type taskInfo struct {
	status      taskState
	timeoutTime time.Time
	fileName string
	chunkInfo   ChunkInfo
}

const (
	TaskIsReady      taskState = 0
	TaskInQueue      taskState = 1
	TaskIsProcessing taskState = 2
	TaskIsCompleted  taskState = 3	
	TaskHasErr       taskState = 4
)


type Master struct {
	// Your definitions here.
	inputFiles []string
	nReduce    int
	nMap       int

	tasksChan chan Task

	allDone bool

	mutex sync.Mutex

	phase      string
	taskstatus []taskInfo
}

// Your code here -- RPC handlers for the worker to call.

//
// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
//
func (m *Master) WorkerRequestHandler(args *Args, reply *Reply) error {
	if args.TaskIndex == -1 {
		//fmt.Println("[INFO][MASTER] serving request ..")
		task, ok := <-m.tasksChan //stuck
		if ok {
			reply.TodoTask = task
			// fmt.Printf("Task number %v from file %v with offset %v taken from queue \n", task.TaskIndex, task.InputFile, task.FileChunk.OffsetStart)

			reply.NReduce = m.nReduce
			reply.NMap = m.nMap
			m.taskstatus[task.TaskIndex].status = TaskIsProcessing
			m.taskstatus[task.TaskIndex].timeoutTime = time.Now().Add(maxTimeout * time.Second)
		} else {
			reply.AllTasksAreDone = true
		}
	} else {
		if args.Finished {
			//fmt.Println("[INFO][MASTER][Report] completed ", args.TaskIndex, " ", args.TaskType)
			if args.TaskType == m.phase {
				m.taskstatus[args.TaskIndex].status = TaskIsCompleted
			}
		} else {
			m.taskstatus[args.TaskIndex].status = TaskHasErr
		}
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	fmt.Println("SERVER IS UP")
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	ret := false

	// Your code here.
	finished := true
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//========= add/remove tasks to chan =========
	for i, taskInfo := range m.taskstatus {
		switch taskInfo.status {
		case TaskIsReady:
			// println(i, "calling Add Task to Queue")
			finished = false
			m.addTaskToQueue(i)
		case TaskInQueue:
			finished = false
		case TaskIsProcessing:
			finished = false
			m.checkTaskNotExpire(i)
		case TaskIsCompleted:
		case TaskHasErr:
			finished = false
			m.addTaskToQueue(i)
		default:
			fmt.Print("[ERROR][MASTER]task state error")
		}
	}

	if finished {
		if m.phase == IsMap {
			//fmt.Println("[INFO][MASTER]MAP TASKS DONE")
			m.initReduceTasks()
		} else {
			m.allDone = true
			close(m.tasksChan)
		}
	} else {
		m.allDone = false
	}
	ret = m.allDone
	return ret
}

func (m *Master) initReduceTasks() {
	m.phase = IsReduce
	m.taskstatus = make([]taskInfo, m.nReduce)
	for i := range m.taskstatus {
		m.taskstatus[i].status = TaskIsReady
	}
}

func (m *Master) addTaskToQueue(taskIndex int) {
	m.taskstatus[taskIndex].status = TaskInQueue

	task := Task{
		TaskType:  m.phase,
		TaskIndex: taskIndex,
	}
	if m.phase == IsMap {
		task.InputFile = m.taskstatus[taskIndex].fileName
		task.FileChunk = m.taskstatus[taskIndex].chunkInfo
	}

	m.tasksChan <- task

	// fmt.Printf("Task number %v from file %v with offset %v added to queue \n", task.TaskIndex, task.InputFile, task.FileChunk.OffsetStart)


}

func (m *Master) checkTaskNotExpire(taskIndex int) {
	if m.taskstatus[taskIndex].timeoutTime.Sub(time.Now()) > 0 { //expire
		m.addTaskToQueue(taskIndex)
	}
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	//initialize m
	m.inputFiles = files
	//initialize tasks with chosen chunk size
	m.taskstatus = makeTaskSlices(files,2000)
	m.nReduce = nReduce
	println("Number of map tasks: ", len(m.taskstatus), "\n")
	m.nMap = len(m.taskstatus)
	m.mutex = sync.Mutex{}
	if m.nReduce > m.nMap { //todo: # of file chunks
		m.tasksChan = make(chan Task, m.nReduce)
	} else {
		m.tasksChan = make(chan Task, m.nMap)
	}
	m.phase = IsMap

	for i := range m.taskstatus {
		m.taskstatus[i].status = TaskIsReady
	}

	m.server()
	return &m
}


func makeTaskSlices(files []string, sizeEachChunk int) []taskInfo{
	var tasks []taskInfo

	for _,file := range files {
		chunks := sliceOneFile(file, sizeEachChunk)
		for _,chunk := range chunks {
			info := taskInfo{
				status: TaskIsReady,
				fileName: file,
				chunkInfo: chunk,
			}
			tasks = append(tasks,info)
		}
	}

	return tasks
}

func sliceOneFile(file string, sizeEachChunk int) []ChunkInfo{

	var chunkInfos []ChunkInfo

	f, err := os.Open(file)
	check(err)
	var offset int64 = 0
	for {
		b := make([]byte, sizeEachChunk)
		_, err1 := f.Seek(offset, 0)
		check(err1)
		n, err2 := f.Read(b)

		// ends at punctuation or space at the end of file
		if (n == 1){
			break
		}
		if err2 != nil {
			break
		}
		end := getOffsetEnd(int64(n), b[:])

		chunkInfo := ChunkInfo{
			OffsetStart: offset,
			ChunkSize: sizeEachChunk,
		}

		chunkInfos = append(chunkInfos, chunkInfo)

		offset += end

	}

	f.Close()

	return chunkInfos

}


//offset end : exclusive
func getOffsetEnd(n int64, chunk []uint8) int64 {
	for i := n - 1; i >= 0; i-- {
		if !unicode.IsLetter(rune(chunk[:][i])) {
			//fmt.Printf(" %#U -- %d %d\n", rune(chunk[:][i]), i, n)
			return int64(i)
		}
	}
	fmt.Printf("%s\n", chunk[:n])
	return 0
}

func check(e error) {
	if e != nil {
		fmt.Println("||||||||||||||||")
		panic(e)
	}
}

