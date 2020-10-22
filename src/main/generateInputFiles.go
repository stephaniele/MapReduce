package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: generateinputFiles needs a source file and an 'eveness' number...\n")
		os.Exit(1)
	}
	sourceFile := os.Args[1]
	v, _ := strconv.Atoi(os.Args[2])
	generateFiles(sourceFile, v)

}

func generateFiles(source string, disparitySeed int) []string {
	sequence, sum := fibonacci(disparitySeed)
	//get size
	sourceFile, err := os.Open(source)
	check(err)
	defer sourceFile.Close()
	sf, err := sourceFile.Stat()
	check(err)

	//each chunk size
	fileNames := make([]string, disparitySeed+1)
	for i := 0; i <= disparitySeed; i++ {
		fileNames[i] = fmt.Sprintf("input-%d.txt", i)
		chunkSize := (sf.Size() * sequence[i]) / sum
		processFile(sourceFile, fileNames[i], chunkSize)
	}

	return fileNames
}

func fibonacci(disparitySeed int) ([]int64, int64) {
	res := make([]int64, disparitySeed+1)
	sum := int64(1)
	res[0] = 1
	if disparitySeed == 0 {
		return res, sum
	}
	res[1] = 1
	sum++
	for i := 2; i <= disparitySeed; i++ {
		res[i] = res[i-1] + res[i-2]
		sum += int64(res[i])
	}
	return res, sum
}

func processFile(source *os.File, file string, chunkSize int64) {
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
