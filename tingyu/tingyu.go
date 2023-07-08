package tingyu

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

/**
 * @Author: _niuzai
 * @Date:   2023/7/5 10:56
 * @Description:定义路由规则
 */

// HandlerFunc 定义请求处理 给用户使用的
// type HandlerFunc func(w http.ResponseWriter, r *http.Request)
type HandlerFunc func(*Context)

// Engine 实现ServeHTTP接口
type Engine struct {
	// router 存储路由规则
	//router map[string]HandlerFunc
	router *router
	*RouterGroup
	// 存储所有分组
	groups []*RouterGroup
	// html 渲染 将模版加载进内存中
	htmlTemplate *template.Template
	// 自定义模版渲染函数
	funcMap template.FuncMap
}

// RouterGroup 分组
type RouterGroup struct {
	// 前缀
	prefix string
	// 资源统一协调 这样就可以间接访问各种接口了
	engine *Engine
	// 前分组
	parent *RouterGroup
	// 中间件
	middlewares []HandlerFunc
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

// New 返回Engine对象 并对map进行初始化
func New() *Engine {
	//return &Engine{router: make(map[string]HandlerFunc)}
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Use 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// createStatisHandler 创建静态文件处理
func (group *RouterGroup) createStatisHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 找到文件的真实路径 relativePath:/assets fs:./static  absolutePath :/assets
	absolutePath := path.Join(group.prefix, relativePath)
	//http.StripPrefix函数创建一个文件服务处理器，以便正确地处理静态文件的请求
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查文件是否存在 是否有权限获取它
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static 服务静态文件
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStatisHandler(relativePath, http.Dir(root))
	// 拼接完整路径出来
	urlPattern := path.Join(relativePath, "/*filepath")
	fmt.Println(urlPattern)
	// 注册到GET的handlers 注册到路由中去
	group.GET(urlPattern, handler)
}

// Group 路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	// 新的engine中创建slice 类似套娃这种
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute 添加分组路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplate = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

// addRouter 添加路由 参数1 请求方式 参数2 请求路径 参数3 处理函数
func (e *Engine) addRouter(method string, path string, handlerFunc HandlerFunc) {
	//key := method + "-" + path
	//e.router[key] = handlerFunc
	e.router.addRoute(method, path, handlerFunc)
}

// GET Get请求方式
func (e *Engine) GET(path string, handlerFunc HandlerFunc) {
	e.addRouter("GET", path, handlerFunc)
}

// POST Post请求方式
func (e *Engine) POST(path string, handlerFunc HandlerFunc) {
	e.addRouter("POST", path, handlerFunc)
}

// Run 开启服务
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// ServeHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取请求方式和请求路径
	//key := r.Method + "-" + r.URL.Path
	// 判断key是否存在
	//handler, ok := e.router[key]
	//if ok {
	//	handler(w, r)
	//} else {
	//	w.WriteHeader(http.StatusNotFound)
	//	// 路径不存在 返回404
	//	fmt.Fprintf(w, "404 Not Found: %v\n", r.URL.Path)
	//}
	//c := newContext(w, r)
	//e.router.handler(c)
	//
	// 处理请求的时候加上中间件了
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		// 判断路径中是否有包含的前缀
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.handlers = middlewares
	c.engin = e
	e.router.handler(c)
}
