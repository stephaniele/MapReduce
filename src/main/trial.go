package main

import (
	"fmt"
	"os"
)

func main() {

	file, err := os.Open("pg-huckleberry_finn.txt")

	if err != nil {
		fmt.Println("cannot open %v", "pg-being_ernest.txt")
	}

	content := make([]byte, 2000)

	_, err2 := file.ReadAt(content,int64(1992))
	fmt.Println("CONTENT READ: ", string(content))
	fmt.Println("--------------------------------")

	if err2 != nil {
		fmt.Println("cannot open %v", "pg-being_ernest.txt")
	}

}


