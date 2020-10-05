package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

const (
	IsMap    string = "map"
	IsReduce string = "reduce"
)

const (
	stateIdle       string = "idle"
	stateProcessing string = "processing"
	stateCompleted  string = "done"
)

type Task struct {
	taskType   string
	inputFile  string
	taskIndex  int //to generate intermediate file name
	state      string
	outputFile string //[]string for reduce file

	//time:how long been working on this task

}

type Args struct {
	taskIndex int
	finished  bool
}

type Reply struct {
	todoTask        Task
	allTasksAreDone bool
	nReduce         int
	nMap            int
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
