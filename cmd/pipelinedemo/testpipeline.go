/**
 * @Author: Donlin
 * @Date: Created in 17:01 2018/6/22
 * @Version: 1.0
 * @Description: 
 */
package main

import (
	"GoMergeSort/pipeline"
	"fmt"
	"os"
	"bufio"
)

func main() {
	const fileName = "small.in"
	const n = 64

	// Create a file to store random numbers
	file, err := os.Create(fileName)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)	// Wrapper with a buffer writer
	p := pipeline.RandomSource(n)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	// Open the file
	file, err = os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Wrapper with a buffer reader
	p = pipeline.ReaderSource(bufio.NewReader(file), -1)

	// Print the top 20 lines binary context of file
	count := 0
	for v := range p{
		fmt.Println(v)
		count++
		if count > 20{
			break
		}
	}
}

func MergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemSort(
			pipeline.ArraySource(2, 4, 1, 56, 43, 2)),
		pipeline.InMemSort(
			pipeline.ArraySource(2, 4, 66, 743, 23, 42)))
	// 第一种写法
	//for {
	//	// 从channel p中不断读int，直到ok为false（没有数据）
	//	if num, ok := <- p; ok{
	//		fmt.Println(num)
	//	}else {
	//		break
	//	}
	//}
	// 第二种写法，使用 range 关键词
	for v := range p {
		fmt.Println(v)
	}
}
