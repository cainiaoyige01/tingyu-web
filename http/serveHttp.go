package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 10:22
 * @Description: 实现ServeHTTP的例子
 */
func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":8080", engine))
}

// Engine 定义一个结构体
type Engine struct {
}

// ServeHTTP 实现接口
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取http请求路径
	path := r.URL.Path
	// 多个路由 使用switch
	switch path {
	case "/":
		fmt.Fprintf(w, "url path:%s", path)
	case "/hello":
		fmt.Fprintf(w, "url path:%s----context:%s", path, "hello world")
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", path)
	}
}
