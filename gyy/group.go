package gyy

import (
	"log"
	"net/http"
	"path"
)

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

func (g *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	// 解析请求地址，映射到服务器上文件的真实地址，交给 http.FileServer 处理
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 将磁盘上的文件夹 root 映射到路由 relativePath
func (g *RouterGroup) Static(relativePath, root string) {
	// 获取处理函数
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	// 拼接路由名称
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}
