package gyy

import (
	"fmt"
	"net/http"
)

// 定义请求函数的格式
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 实现 ServeHTTP 接口
type Engine struct {
	router map[string]HandlerFunc
}

// 初始化 Engine
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由
func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 拼接 key，通过"请求方法-路由名称"映射到对应的处理函数
	key := method + "-" + pattern
	e.router[key] = handler
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

// 解析请求路径，查找 router，能找到就执行，找不到说明是 404
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, r)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND：%s\n", r.URL)
	}
}
