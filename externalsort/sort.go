package main

import (
	"bufio"
	"fmt"
	"os"
	"readerDemo/pipeline"
)

func main() {
	p := createPipeline("small.in", 512, 4)
	writeToFile(p, "small.out")
	printFile("small.out")
}

//todo  close file
func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	chunkSize := fileSize / chunkCount
	sortResult := []<-chan int{}
	pipeline.Init()
	for i := 0; i < chunkCount; i++ {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i*chunkSize), 0)
		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResult = append(sortResult, pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortResult...)
}

func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()
	pipeline.WriteSink(writer, p)

}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)

	for v := range p {
		fmt.Println(v)
	}
}
