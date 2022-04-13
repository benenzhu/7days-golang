package gee

import (
	"fmt"
	"strings"
)

// 这是一个数的节点
type node struct {
	pattern  string  // 待匹配的路由
	part     string  // 路由中的一部分
	children []*node // 子节点 [doc tutorial intro]
	isWild   bool    // 是否精确匹配.
	// insert(pattern string, parts []string, height int) // 递归插入路由
	// search(parts []string, height int) *node 	// 递归寻找路由
	// travel(list *([]*node)) 		// 对树进行遍历, 为了获取所有的路由
	// matchChild(part string) *node		// 第一个匹配成功的节点
	//										, 用于插入
	// matchChildren(part string) []*node 	// 所有匹配成功的节点, 用于查找
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height { // 如果到了最后一层了. 直接赋值返回.
		n.pattern = pattern
		return
	}

	part := parts[height]       // 这一层的
	child := n.matchChild(part) // 第一个match到的节点.
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'} // 如果为空 新创建一个 child 节点
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") { // 查找是否有当前这个点.
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1) // 在里面找到第一个对吧. 那为什么还需要返回所有的孩子呢?
		if result != nil {
			return result
		}
	}

	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}

func (n *node) matchChild(part string) *node {
	// 返回第一个找到的.
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0) // 在所有的里面, 如果完全匹配或者是 wild 的话就全部返回
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
