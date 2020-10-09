package main

import (
	"fmt"
	"os"
	"unicode"
)

func main() {
	files := []string{"tt", "pg-being_ernest.txt", "pg-dorian_gray.txt"}
	var sizeEachChunk int64 = 15

	//tasks := make([]taskInfo, 100)
	//for _, file := range files {

	// f, err := os.Open(files[0])
	// check(err)

	// o2, err := f.Seek(6, 0)
	// check(err)

	// b2 := make([]byte, sizeEachChunk)
	// n2, err := f.Read(b2)
	// check(err)

	// index := getOffsetEnd(n2, b2[:n2])
	// fmt.Printf("%s bytes @ %d: %s ", b2[:], o2, b2[:index])

	// f.Close()
	//}

	// f, err := os.Open(files[0])
	// check(err)

	// o2, err := f.Seek(6, 0)
	// check(err)

	// b2 := make([]byte, sizeEachChunk)

	// n2, err := f.Read(b2)
	// check(err)
	// fmt.Printf("%s bytes @ %d: %d \n", b2[:n2], o2, n2)

	// o2, err = f.Seek(2, 0)
	// n2, err = f.Read(b2)
	// check(err)
	// fmt.Printf("%s bytes @ %d: %d \n", b2[:n2], o2, n2)
	processFile(files[0], sizeEachChunk)

}

func processFile(file string, chunkSize int64) {
	f, err := os.Open(file)
	check(err)
	b := make([]byte, chunkSize)
	var offset int64 = 0
	for {
		_, err1 := f.Seek(offset, 0)
		check(err1)
		n, err2 := f.Read(b)
		if err2 != nil {
			break
		}
		end := getOffsetEnd(chunkSize, b[:])
		fmt.Printf("%v - %v |%s|%s|%s| end: %v \n", offset, offset+end-1, b[:end], b[:n], b[:], end)
		offset += end
	}
	f.Close()
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
