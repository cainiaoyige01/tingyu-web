package main

import (
	"fmt"
	"log"
	"net/http"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 9:48
 * @Description: 看一下标准库的http请求
 */
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Path)
		w.Write([]byte("ping"))
	})
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
