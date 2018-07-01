package main

import (
	"fmt"
	"readerDemo/pipeline"
	"os"
	"bufio"
)
// 将数据分为左右两半， 分别归并排序， 再把两个有序数据归并

func mergeDemo()  {
	// create a slice of int
	p1 := pipeline.InMemSort(pipeline.ArraySource(3, 2,6,8,9,10))
	p2 := pipeline.InMemSort(pipeline.ArraySource(32, 32,6,28,49,510))
	p := pipeline.Merge(p1, p2)
	for v:= range p{
		fmt.Println(v)
		}
}

func main(){
	const filename = "small.in"
	const n = 64
	file, err := os.Create(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()
	// file 就是一个reader
	p := pipeline.RandomSource(n)
	writer := bufio.NewWriter(file)
	pipeline.WriteSink(writer, p)
	writer.Flush()
	file, err = os.Open(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()
	p = pipeline.ReaderSource(bufio.NewReader(file), -1)
	count := 0
	for v:= range p{
		fmt.Println(v)
		count ++
		if count > 100{
			break
		}
	}


	}