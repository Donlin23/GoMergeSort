/**
 * @Author: Donlin
 * @Date: Created in 20:40 2018/6/22
 * @Version: 1.0
 * @Description: 
 */
package main

import (
	"os"
	"GoMergeSort/pipeline"
	"bufio"
	"fmt"
	"strconv"
)

func main() {
	p := createNetworkPipeline(
		"lager.in", 800000000, 4)
	writeToFile(p, "lager.out" )
	printFile("lager.out")
}

/**
	Print the sorted file
 */
func printFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)

	count := 0
	for v := range p {
		fmt.Println(v)
		count ++
		if count >= 100{
			break
		}
	}
}

/**
	Write the sorted result into File
 */
func writeToFile(
	p <-chan int, fileName string) {

	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

/**
	Reading file in "chunkCount" chunks
 */
func createPipeline(fileName string,
	fileSize, chunkCount int) <-chan int{

	chunkSize := fileSize / chunkCount 	//chunkCount means we will divides into "chunkCount" chunks
	pipeline.Init() // Init a startTime to log
	sortResults := []<-chan int{}
	for i := 0; i < chunkCount; i++{
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)

		source := pipeline.ReaderSource(
			bufio.NewReader(file), chunkSize)

		sortResults = append(
			sortResults, pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortResults...)

}
/**
	Reading from network in "chunkCount" chunks
 */
func createNetworkPipeline(fileName string,
	fileSize, chunkCount int) <-chan int{

	chunkSize := fileSize / chunkCount 	//chunkCount means we will divides into "chunkCount" chunks
	pipeline.Init() // Init a startTime to log

	sortAddr := []string{}
	for i := 0; i < chunkCount; i++{
		file, err := os.Open(fileName)
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)

		source := pipeline.ReaderSource(
			bufio.NewReader(file), chunkSize)

		addr := ":" + strconv.Itoa(7000 + i)
		pipeline.NetworkSink(
			addr, pipeline.InMemSort(source))
		sortAddr = append(sortAddr, addr)
	}

	sortResults := []<-chan int{}
	for _, addr := range sortAddr  {
		sortResults = append(
			sortResults, pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortResults...)

}
