package tingyu

import (
	"log"
	"time"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/8 8:57
 * @Description:中间件 测试调用时间
 */

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		//time.Sleep(time.Second * 1)
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Now().Sub(t))
	}
}
