package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

const fileTotalNumber = 4

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: generateinputFiles needs a source file and an 'eveness' number...\n")
		os.Exit(1)
	}
	sourceFile := os.Args[1]
	v, _ := strconv.Atoi(os.Args[2])
	if v < 3 {
		println(os.Stderr, "Usage: disparity seed needs to >= 3")
		os.Exit(1)
	}
	generateFiles(sourceFile, v)

}

func generateFiles(source string, disparitySeed int) []string {
	sequence := fibonacci(disparitySeed)
	sequence = sequence[len(sequence)-fileTotalNumber:]
	sum := uint64(0)
	for _, n := range sequence {
		sum += uint64(n)
	}
	//get size
	sourceFile, err := os.Open(source)
	check(err)
	defer sourceFile.Close()
	sf, err := sourceFile.Stat()
	check(err)

	//each chunk size
	fileSize := uint64(sf.Size())
	fileNames := make([]string, fileTotalNumber)
	for i := 0; i < fileTotalNumber; i++ {
		fileNames[i] = fmt.Sprintf("input-%d.txt", i)
		chunkSize := float64(fileSize) / float64(sum) * float64(sequence[i])
		size := uint64(chunkSize)
		processFile(sourceFile, fileNames[i], size)
	}

	return fileNames
}

func fibonacci(disparitySeed int) []int64 {
	res := make([]int64, disparitySeed+1)
	res[0] = 1
	if disparitySeed == 0 {
		return res
	}
	res[1] = 1
	for i := 2; i <= disparitySeed; i++ {
		res[i] = res[i-1] + res[i-2]
		if res[i] < res[i-1] {
			return res[:i]
		}
	}
	return res
}

func processFile(source *os.File, file string, chunkSize uint64) {
	f, err := os.Create(file)
	check(err)
	defer f.Close()

	b := make([]byte, chunkSize)

	_, err = source.Seek(0, 1)

	check(err)
	n, err := source.Read(b)
	check(err)
	_, err = f.Write(b[:n])
	check(err)

	_, err = f.Write([]byte("."))
	check(err)
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
		fmt.Println(e)
		panic(e)
	}
}
