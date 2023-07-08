package tingyu

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/8 16:12
 * @Description: 错误恢复
 */

// trance 打印堆栈跟踪信息
func trance(message string) string {
	// 定义uintptr数组 用于存储调用栈信息
	var pcs [32]uintptr
	//跳过前3个调用者 即当前函数、trace函数和调用trace函数的函数
	n := runtime.Callers(3, pcs[:])
	// 用于构建字符串
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	// 遍历栈中有效的信息 对每一个调用点进行处理
	for _, pc := range pcs[:n] {
		// 获取包含调用点的函数对象
		fn := runtime.FuncForPC(pc)
		// 获取文件名和行号
		file, line := fn.FileLine(pc)
		// 将文件名和行号进行格式化为字符串 添加到str
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			err := recover()
			if err != nil {
				message := fmt.Sprintf("%v", err)
				log.Printf("%s\n\n", trance(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
