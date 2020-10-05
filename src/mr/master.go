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
)

type taskState int

const maxTimeout = 10

type taskInfo struct {
	status      taskState
	timeoutTime time.Time
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
	fmt.Print("[INFO][MASTER] start to serve worker's request ...")
	if args.taskIndex == -1 {
		task, ok := <-m.tasksChan
		if ok {
			reply.todoTask = task
			reply.nReduce = m.nReduce
			reply.nMap = m.nMap

			m.taskstatus[task.taskIndex].status = TaskIsProcessing
			m.taskstatus[task.taskIndex].timeoutTime = time.Now().Add(maxTimeout * time.Second)
		} else {
			reply.allTasksAreDone = true
		}
	} else {
		if args.finished {
			m.taskstatus[args.taskIndex].status = TaskIsCompleted
		} else {
			m.taskstatus[args.taskIndex].status = TaskHasErr
		}
	}
	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
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
			//init reduce tasks
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
	m.allDone = false
	m.taskstatus = make([]taskInfo, m.nReduce)
	for i := range m.taskstatus {
		m.taskstatus[i].status = TaskIsReady
	}
}

func (m *Master) addTaskToQueue(taskIndex int) {
	m.taskstatus[taskIndex].status = TaskInQueue
	task := Task{
		taskType:  m.phase,
		inputFile: "",
		taskIndex: taskIndex,
		//time:how long been working on this task
	}
	if m.phase == IsMap {
		task.inputFile = m.inputFiles[taskIndex]
	}
	m.tasksChan <- task
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
	m.nReduce = nReduce
	m.nMap = len(files)
	m.mutex = sync.Mutex{}
	if m.nReduce > m.nMap {
		m.tasksChan = make(chan Task, m.nReduce)
	} else {
		m.tasksChan = make(chan Task, m.nMap)
	}
	m.phase = IsMap

	//initialize tasks
	m.taskstatus = make([]taskInfo, m.nMap)
	for i := range m.taskstatus {
		m.taskstatus[i].status = TaskIsReady
	}

	m.server()
	return &m
}
