package gyy

import (
	"net/http"
	"strings"
)

// 定义请求函数的格式
type HandlerFunc func(*Context)

// 实现 ServeHTTP 接口
type Engine struct {
	*RouterGroup                // 存储根分组
	groups       []*RouterGroup // 存储全部分组
}

// 初始化 Engine
func New() *Engine {
	group := newRootGroup()
	engine := &Engine{RouterGroup: group, groups: []*RouterGroup{group}}
	group.addEngine(engine)
	return engine
}

// 启动 http
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

// 请求入口
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	// 通过相同前缀，找到当前路由需要应用的中间件
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			c.handlers = append(c.handlers, group.middlewares...)
		}
	}
	e.handle(c)
}
