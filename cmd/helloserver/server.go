/**
 * @Author: Donlin
 * @Date: Created in 16:11 2018/6/22
 * @Version: 1.0
 * @Description: Helloworld web版
 */
package main

import (
	"net/http"
	"fmt"
)

func main() {

	http.HandleFunc("/", func(
		writer http.ResponseWriter,
		request *http.Request) {
			// 可以从 request 读取请求参数，往 writer 写响应数据
		fmt.Fprintf(writer, "<h1>Hello World! %s</h1>",
			request.FormValue("name"))

	})

	// 监听本地8888端口
	http.ListenAndServe(":8888", nil)
}
