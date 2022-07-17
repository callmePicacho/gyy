package gyy

import "log"

// 路由分组
type RouterGroup struct {
	prefix      string        // 存储前缀
	middlewares []HandlerFunc // 存储中间件
	parent      *RouterGroup  // 支持嵌套
	*router                   // 存储顶级路由
}

// 创建默认根路由分组
func newRootGroup() *RouterGroup {
	return &RouterGroup{router: newRouter()}
}

// 添加分组
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	return &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		router: g.router,
	}
}

// 添加路由
func (g *RouterGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = g.prefix + pattern
	log.Printf("Route %4s - %s", method, pattern)
	g.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}
