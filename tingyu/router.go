package tingyu

import (
	"net/http"
	"strings"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 21:37
 * @Description:路由 提供了动态的支持
 */

type router struct {
	handlers map[string]HandlerFunc
	// 存储动态路由
	roots map[string]*node
}

// newRouter creates a new router
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

// parsePattern 解析pattern 也就是把路由/p/:niuzai/doc 查分[p,:niuzai,doc]
func parsePattern(pattern string) []string {
	// 分割
	split := strings.Split(pattern, "/")
	// 定义切片
	parts := make([]string, 0)
	// 拼接内容 还要判断是否含有"*"
	for _, v := range split {
		if v != "" {
			parts = append(parts, v)
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRoute 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//log.Printf("Route %4s - %s", method, pattern)
	//key := method + "-" + pattern
	//r.handlers[key] = handler

	parts := parsePattern(pattern)
	key := method + "-" + pattern
	// 判断map是否存在路由了 相同的路由也会GET与POST的
	_, ok := r.roots[method]
	if !ok {
		// 不存在
		r.roots[method] = &node{}
	}
	// 添加路由
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler

}

// handler 给用户使用
func (r *router) handler(c *Context) {
	//key := c.Method + "-" + c.Path
	//handler, ok := r.handlers[key]
	//if ok {
	//	handler(c)
	//} else {
	//	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	//}

	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {

		// 是n不是c
		key := c.Method + "-" + n.pattern
		c.Params = params
		// 用户调用的web框架函数做最后的处理
		c.handlers = append(c.handlers, r.handlers[key])

	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}

// getRoute 匹配路由
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchPaths := parsePattern(pattern)
	// params 出现的意义是让我们可以context中获取到:niuzai 参数的值
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchPaths, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchPaths[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchPaths[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
