package gyy

import "net/http"

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 拼接 key，通过"请求方法-路由名称"映射到对应的处理函数
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 解析请求路径，查找 router，能找到就执行，找不到说明是 404
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
