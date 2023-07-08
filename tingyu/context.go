package tingyu

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 21:00
 * @Description: 设计context
 *Context目前只包含了http.ResponseWriter和*http.Request，
 *另外提供了对 Method 和 Path 这两个常用属性的直接访问。
 *提供了访问Query和PostForm参数的方法。
 *提供了快速构造String/Data/JSON/HTML响应的方法
 */

// H 用于map数据返回 看过gin框架就会很熟悉！构建json数据显得更加简洁
type H map[string]interface{}

// Context 封装请求与响应信息
type Context struct {
	// http来源
	Writer http.ResponseWriter
	Req    *http.Request
	// Request 信息
	Path   string
	Method string
	// 增加一个map 用于获取动态的路由的值
	Params map[string]string
	// 响应状态码
	StatusCode int
	// 中间件
	handlers []HandlerFunc
	index    int
	engin    *Engine
}

// newContext 返回Context指针
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		// 状态码有用户来设计
		//StatusCode: 0,
		// 用于判断slice中有没有中间件的
		index: -1,
	}
}

func (c *Context) Next() {
	c.index++
	// 遍历调用HnadlerFunc
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// Param 获取动态路由中的值
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// PostForm 表单根据key获取value值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取路径参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Status 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 放回string类型数据给客户端
func (c *Context) String(code int, format string, values ...interface{}) {
	// 设置头信息
	c.SetHeader("Content-type", "text/plain")
	// 设置状态码
	c.Status(code)
	// 返回数据给客户端
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 返回json格式的数据
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-type", "application/json")
	c.Status(code)
	// 把Writer中数据设置为json格式返回
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// Data 返回字节数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// HTML 返回HTML页面
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-type", "text/html; charset=utf-8")
	c.Status(code)
	err := c.engin.htmlTemplate.ExecuteTemplate(c.Writer, name, data)
	if err != nil {
		fmt.Println(500, err.Error())
	}

}
