/**
 * @Author: Donlin
 * @Date: Created in 22:00 2018/6/22
 * @Version: 1.0
 * @Description: 
 */
package pipeline

import (
	"net"
	"bufio"
)

/**
	Writer into network
 */
func NetworkSink(addr string, in <-chan int){
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		writer := bufio.NewWriter(conn)
		defer writer.Flush()

		WriterSink(writer,in)
	}()

}

/**
	Read from network
 */
func NetworkSource(addr string) <-chan int{
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		r := ReaderSource(
			bufio.NewReader(conn), -1)
		for v := range r {
			out <- v
		}
		close(out)
	}()
	return out
}