/**
 * @Author: Donlin
 * @Date: Created in 16:25 2018/6/22
 * @Version: 1.0
 * @Description: helloworld 并发执行版(deadlock 版)
 */
package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)
	for i := 1; i < 50 ; i++ {
		// go starts a goroutine
		go printHelloWorld(i, ch)
	}

	// 通过通道进行通信，从另一个 goroutine 接收 string
	for {
		msg := <- ch
		fmt.Println(msg)
	}
}

func printHelloWorld(i int, ch chan string)  {
	ch <- fmt.Sprintf("Hello world from goroutine %v!\n", i)
}
