package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"time"
	"fmt"
)
var startztime time.Time
func Init(){
	startztime = time.Now()
}

func ArraySource(a ...int) <-chan int {
	out  := make(chan int)
	go func(){
		for _, v := range a{
			out <- v
		}
		close(out)

	}()
	return out
}

func InMemSort(in <-chan int) <-chan int{
	out := make(chan int)
	go func() {
		// read into memory
		a := []int{}
		for v:= range in {
			a = append(a, v)
		}
		fmt.Println("Read done: ", time.Now().Sub(startztime))
		//sort
		sort.Ints(a)
		fmt.Println("Sort done: ", time.Now().Sub(startztime))

		//output
		for _,v := range a{
			out <- v
		}
		close(out)
	}()
	return  out
}

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
	}()
	return out
}

func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func(){
		byteSize := 0
		buffer := make([]byte, 8)
		for {
			n, err := reader.Read(buffer)
			byteSize = byteSize + n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && byteSize >= chunkSize) {
				break
			}
		}
		close(out)
	}()

	return out
}

func WriteSink(write io.Writer, in <-chan int){
	for v:= range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		write.Write(buffer)
	}
}

func RandomSource(count  int) <- chan int {
		out := make(chan int)
		go func() {
			for i :=0; i < count; i++{
				out <- rand.Int()
			}
			close(out)
		}()

	return out
}

func MergeN(inputs ...<-chan int) <-chan int {

	if len(inputs) == 1{
		return inputs[0]
	}

	m := len(inputs) /2

	return Merge(
		MergeN(inputs[:m]...),
		MergeN(inputs[m:]...),
	)

}