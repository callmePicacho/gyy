package gyy

import (
	"fmt"
	"strings"
)

// 前缀树
type node struct {
	pattern  string  // 待匹配路由，例如：/p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配，part 含有 : 或 * 时为 true
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// 找到第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 注册路由
// pattern 是访问的路由全路径，例如 /p/abc
// parts 是根据'/'分隔后的路由数组，例如 [p, abc]
// height 是当前trie树的层高
func (n *node) insert(pattern string, parts []string, height int) {
	// 如果已经匹配完成，将pattern赋值给该node，表示其为完整的url
	// 即当某个节点的 pattern 不为 ""，从根节点到该节点的路径为一个路由，且完整路由存在该节点的 pattern 中
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	// 先查找树中是否已经存在节点
	child := n.matchChild(part)
	// 如果 child 为 nil，说明在 node 中未注册
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 接着去下一层节点
	// 比如注册第一个路由 /a/b/c，需要创建多层、多个节点
	child.insert(pattern, parts, height+1)
}

// 找到所有匹配成功的子节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 查找路由
func (n *node) search(parts []string, height int) *node {
	// 找到了末尾，或者遇到了通配符
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// pattern 为空字符串表示非叶子节点，匹配失败
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	// 查找所有符合条件的子节点
	children := n.matchChildren(part)

	// 遍历每个子节点，直到找到匹配项
	for _, child := range children {
		result := child.search(parts, height+1)
		// 找到了
		if result != nil {
			return result
		}
	}
	return nil
}

// 树的层级遍历，查找所有完整的 url，保存到列表中
func (n *node) Travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.Travel(list)
	}
}
