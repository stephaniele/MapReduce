package main

//
// start the master process, which is implemented
// in ../mr/master.go
//
// go run mrmaster.go pg*.txt
//
// Please do not change this file.
//

import (
	"fmt"
	"log"
	"os"
	"time"

	"../mr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrmaster inputfiles...\n")
		os.Exit(1)
	}

	start := time.Now()
	m := mr.MakeMaster(os.Args[1:], 10)
	for m.Done() == false {
		time.Sleep(time.Second)
	}
	elapsed := time.Since(start)
	log.Printf("%s", elapsed)
	time.Sleep(time.Second)
}
