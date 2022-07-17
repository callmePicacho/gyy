package gyy

import (
	"net/http"
)

// 定义请求函数的格式
type HandlerFunc func(*Context)

// 实现 ServeHTTP 接口
type Engine struct {
	*RouterGroup
}

// 初始化 Engine
func New() *Engine {
	group := newRootGroup()
	return &Engine{group}
}

// 启动 http
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// 请求入口
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.handle(c)
}
