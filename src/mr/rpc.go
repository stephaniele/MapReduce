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
	TaskType  string
	InputFile string
	TaskIndex int //to generate intermediate file name
	fileChunk FileChunk
}

type Args struct {
	TaskIndex int
	Finished  bool
	TaskType  string
}

type Reply struct {
	TodoTask        Task
	AllTasksAreDone bool
	NReduce         int
	NMap            int
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
