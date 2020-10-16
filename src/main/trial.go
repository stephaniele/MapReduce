package main

import (
	"fmt"
	"os"
	"unicode"
)

func main() {
	files := []string{"tt", "pg-being_ernest.txt", "pg-dorian_gray.txt"}
	var sizeEachChunk int64 = 50

	processFile(files[0], sizeEachChunk)

}

func processFile(file string, chunkSize int64) {
	f, err := os.Open(file)
	check(err)
	var offset int64 = 0
	for {
		b := make([]byte, chunkSize)
		_, err1 := f.Seek(offset, 0)
		check(err1)
		n, err2 := f.Read(b)

		// ends at punctuation or space at the end of file
		if (n == 1){
			break
		}
		if err2 != nil {
			check(err2)
		}
		end := getOffsetEnd(int64(n), b[:])
		fmt.Printf("%v - %v |%s|%s|%s| end: %v \n", offset, offset+end-1, b[:end], b[:n], b[:], end)
		offset += end

	}
	f.Close()
}

//offset end : exclusive
func getOffsetEnd(n int64, chunk []uint8) int64 {
	for i := n - 1; i >= 0; i-- {
		if !unicode.IsLetter(rune(chunk[:][i])) {
			fmt.Printf(" NOT LETTER: %#U -- %d %d\n", rune(chunk[:][i]), i, n)
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
