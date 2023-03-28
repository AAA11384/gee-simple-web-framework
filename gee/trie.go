package gee

import "strings"

type node struct {
	pattern  string //只有作为路由的node此字段不为空，路径中存在:, *
	part     string //当前节点对应的目录名称
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, n := range n.children {
		if n.part == part || n.isWild {
			return n
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, n := range n.children {
		if n.part == part || n.isWild {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if height == len(parts) {
		n.pattern = pattern
		return
	}
	child := n.matchChild(parts[height])
	if child == nil {
		child = &node{
			part:   parts[height],
			isWild: parts[height][0] == '*' || parts[height][0] == ':',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		}
		return nil
	}
	for _, tempNode := range n.matchChildren(parts[height]) {
		result := tempNode.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
