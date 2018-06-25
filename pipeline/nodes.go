/**
 * @Author: Donlin
 * @Date: Created in 16:55 2018/6/22
 * @Version: 1.0
 * @Description: 归并排序节点工具
 */
package pipeline

import (
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"time"
	"fmt"
)

var startTime time.Time

func Init(){
	startTime = time.Now()
}

/**
	Read data from slice a(...int), "...int" is a slice actually
 */
func ArraySource(a ...int) <-chan int{ 	// "<-chan" for the goroutine which receive it
 	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

/**
	Use build-in sort in memory, pull data in memory and sort it, push the sorted result into a new channel(int) out
 */
func InMemSort(in <-chan int) <-chan int{
	out := make(chan int, 1024)
	go func() {
		// Create a empty slice or we can use this way: "a:= make([]int, 0)"
		var a = []int{}
		// Read into memory
		for v := range in {
			a = append(a, v)
		}
		fmt.Println("Read done:", time.Now().Sub(startTime))

		// Sort
		sort.Ints(a)
		fmt.Println("InMemSort done:", time.Now().Sub(startTime))

		// Output
		for _, v := range a{
			out <- v
		}
		close(out)
	}()
	return out
}

/**
	Merge two sorted arrays into one sorted array
 */
func Merge(in1, in2 <-chan int) <-chan int{
	out := make(chan int, 1024)
	go func() {
		// Read two integers from two channels
		int1, ok1 := <- in1
		int2, ok2 := <- in2
		// Just one channel exists num
		for ok1 || ok2 {
			if !ok2 || ( ok1 && int1 < int2) {
				out <- int1
				int1, ok1 = <- in1
			}else {
				out <- int2
				int2, ok2 = <- in2
			}
		}
		close(out)
		fmt.Println("Merge done:", time.Now().Sub(startTime))
	}()
	return out
}

/**
	We can read from a instance which implements io.Reader interface, push data into a channel(int) out
 */
func ReaderSource(
	reader io.Reader, chunkSize int) <-chan int{
		// chunkSize == -1 means we can read all the time
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		bytesRead := 0
		for  {
			// Read() return the length of "buffer"--n and an error
			n, err := reader.Read(buffer)
			bytesRead += n
			if n > 0{
				// Uint64() encapsulates []byte--buffer to a uint--v
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil ||
				(chunkSize != -1 && bytesRead >= chunkSize){
				break
			}
		}
		close(out)
	}()
	return out
}

/**
	Receive from a channel and write in a io.Writer
 */
func WriterSink(writer io.Writer, in <-chan int)  {
	buffer := make([]byte, 8)
	for v := range in {
		// Uint64() inverse process
		binary.BigEndian.PutUint64(buffer,uint64(v))
		writer.Write(buffer)
	}
}

/**
	Generate "count" random numbers
 */
func RandomSource(count int) <-chan int{
	out := make(chan int)
	go func() {
		for i := 0; i < count ; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

/**
	Merge N arrays into one array
 */
func MergeN(inputs ...<-chan int) <-chan int {
	if len(inputs) == 1 {
		return inputs[0]
	}

	m := len(inputs) / 2
	// Merge inputs[0..m) and inputs [m..end)
	return Merge(
		MergeN(inputs[:m]...),
		MergeN(inputs[m:]...))
}