package main

import (
	"fmt"
	"net/http"
	"time"
	"tingyu/tingyu"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 10:56
 * @Description: 主入口
 */

type student struct {
	Name string
	Age  int8
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}
func main() {
	r := tingyu.Default()
	r.GET("/", func(c *tingyu.Context) {
		c.String(http.StatusOK, "Hello tingyu\n")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *tingyu.Context) {
		names := []string{"tingyu"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}
