package engine

import (
	"fmt"
	"strings"
)

var emptyMiddleWaring = "Empty string can't be contained in the pattern's middle,it has been ignored"

//词典树结构体
type trieNode struct {
	//当前节点对应的字符串模板
	pattern string
	//当前节点当前值
	part string
	//当前节点子节点集合
	children []*trieNode
	//是否使用模糊匹配，若当前节点符合规定的字符模板同样可以匹配成功
	isWild bool
}

func rootNode() *trieNode {
	return new(trieNode)
}

func (n *trieNode) AddNode(pattern string) {
	n.insert(patternToParts(pattern), 0)
}
func (n *trieNode) insert(parts []string, height int) {
	if len(parts) == height {
		n.pattern = "/" + strings.Join(parts, "/")
		return
	}
	matchNode := n.matchNode(parts[height])
	if matchNode == nil {
		matchNode = &trieNode{
			part:   parts[height],
			isWild: parts[height][0] == '$' || parts[height][0] == '*',
		}
		n.children = append(n.children, matchNode)
	}
	matchNode.insert(parts, height+1)

}
func (n *trieNode) search(pattern string) *trieNode {
	return n.searchRecursiveAble(patternToParts(pattern), 0)
}
func (n *trieNode) searchRecursiveAble(parts []string, height int) *trieNode {
	if len(parts) == height {
		if n.pattern != "" {
			return n
		}
	}
	matches := n.matchNodes(parts[height])
	for _, item := range matches {
		if res := item.searchRecursiveAble(parts, height+1); res != nil {
			return res
		}
	}
	return nil
}
func (n *trieNode) matchNode(part string) *trieNode {
	for _, node := range n.children {
		if node.part == part {
			return node
		}
	}
	return nil
}
func (n *trieNode) matchNodes(part string) []*trieNode {
	matches := make([]*trieNode, 0)
	for _, node := range n.children {
		if node.part == part || node.isWild {
			matches = append(matches, node)
		}
	}
	return matches
}
func patternToParts(pattern string) []string {
	splitRes := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for index, item := range splitRes {
		if item != "" {
			parts = append(parts, item)
		} else {
			if index != 0 {
				fmt.Printf("%c[1;37;33m%s%c[0m\n\n", 0x1B, emptyMiddleWaring, 0x1B)
			}
		}
	}
	return parts
}
