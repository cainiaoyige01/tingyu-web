package tingyu

import "strings"

/**
 * @Author: _niuzai
 * @Date:   2023/7/6 10:49
 * @Description:设计路由树 目的实现动态路由匹配
 */

// node 定义数节点
type node struct {
	// 待匹配路由 如/p/:niuzai/doc
	pattern string
	// 路由中的一部份 如doc、:niuzai
	part string
	// 子节点  如[]
	children []*node
	// 是否精确匹配？ 如果包含":"或者"*" 就为true
	isWild bool
}

// insert 注册路由
func (n *node) insert(pattern string, parts []string, height int) {
	// 考虑"/"以及"/p/:niuzai"情况 防止索引越界
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 获取part
	part := parts[height]
	// 进行路由路径匹配
	child := n.matchChild(part)
	// 如果返回的child是nil 也就是不存在的！我们就创建一个child节点
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 拼接到n的children中去
		n.children = append(n.children, child)
	}
	// 进行递归注册所有路径
	child.insert(pattern, parts, height+1)
}

// search 查找路由
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	// 匹配子节点出来
	children := n.matchChildren(part)
	// 肯能会存在多个节点 比如 /p /:niuzai 这种
	for _, child := range children {
		// 进行递归查询
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

// matchChild 精确匹配每一个节点
func (n *node) matchChild(part string) *node {
	// 获取n的子节点进行遍历查询
	for _, child := range n.children {
		// 判断节点中part是否相等 或者节点是否精确匹配的？
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 所有匹配成功的节点 用于查询
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
