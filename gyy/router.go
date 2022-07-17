package gyy

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node // 在注册路由时维护，在匹配路由时加速检测路由是否存在 and 实现动态路由
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析pattern，仅允许一个 *
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 注册路由
func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	// 解析路由
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 插入 trie 树
	// 树的第二层节点必然是某个 method 名称
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)

	r.handlers[key] = handler
}

// 匹配路由
// method 和 path 是请求中的，例如 method="get",path="/abc/ccc"
// 刚好trie树中有注册路由 "/abc/:name"，则正好 path 匹配上
func (r *router) getRoute(method, path string) (n *node, params map[string]string) {
	searchParts := parsePattern(path)
	root, ok := r.roots[method]

	if !ok {
		return
	}

	params = make(map[string]string)
	// 在前缀树中匹配路由，如果能匹配上某个路由
	if n = root.search(searchParts, 0); n != nil {
		// 这里的 n.pattern 是注册时的完整路由，即 "/abc/:name"
		parts := parsePattern(n.pattern)
		// 如果是模糊匹配，在 params 中插入 ":name" = "ccc" k-v 对
		for idx, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[idx]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
			}
		}
		return n, params
	}
	return
}

// 返回所有注册路由
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.Travel(&nodes)
	return nodes
}

// 解析请求路径，查找 router，能找到就执行，找不到说明是 404
func (r *router) handle(c *Context) {
	// 匹配路由
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 匹配到路由，赋值路由参数
		c.Params = params
		// 执行该路由对应的处理函数，使用 pattern 即注册时的路由信息作为 key 而非访问时的路由信息作为 key
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
