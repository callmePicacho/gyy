package gyy

import "log"

// 路由分组
type RouterGroup struct {
	prefix      string        // 存储前缀
	middlewares []HandlerFunc // 存储中间件
	parent      *RouterGroup  // 支持嵌套
	*router                   // 存储顶级路由
	engine      *Engine       // 存储 Engine，全部分组存储的是同一份 engine 实例
}

// 创建默认根路由分组
func newRootGroup() *RouterGroup {
	return &RouterGroup{router: newRouter()}
}

func (g *RouterGroup) addEngine(engine *Engine) {
	g.engine = engine
}

// 添加分组
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	group := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		router: g.router,
	}
	group.addEngine(g.engine)
	g.engine.groups = append(g.engine.groups, group)
	return group
}

// 在路由分组中注册中间件
func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
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
