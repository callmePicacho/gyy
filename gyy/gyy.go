package gyy

import (
	"log"
	"net/http"
)

// 定义请求函数的格式
type HandlerFunc func(*Context)

// 实现 ServeHTTP 接口
type Engine struct {
	router *router
}

// 初始化 Engine
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	e.router.addRoute(method, pattern, handler)
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// 启动 http
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// 请求入口
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}
