package main

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
)

type node struct {
	parent *node
	value  int64

	mu       sync.Mutex
	children []*node
}

func createNode(value int64, parent *node) *node {
	n := &node{
		parent:   parent,
		value:    value,
		children: make([]*node, 0),
	}
	if parent != nil {
		parent.mu.Lock()
		if parent.children == nil {
			parent.children = make([]*node, 0)
		}
		parent.children = append(parent.children, n)
		parent.mu.Unlock()
	}
	return n
}

func calculate(n *node, x int64) (total int64) {
	if n == nil {
		return 0
	}
	total += (n.value + x)

	for _, child := range n.children {
		total += calculate(child, (n.value+x)*10)
	}
	return total
}

func getString(n *node, prefix string, isTail bool) string {
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	if isTail {
		buffer.WriteString("└── ")
	} else {
		buffer.WriteString("├── ")
	}
	buffer.WriteString(fmt.Sprintf("Node[value = %d]\n", n.value))
	sort.Slice(n.children,
		func(i, j int) bool { return n.children[i].value < n.children[j].value })
	for i, child := range n.children {
		if i == len(n.children)-1 {
			buffer.WriteString(getString(child, prefix+"  ", true))
		} else {
			buffer.WriteString(getString(child, prefix+"  ", false))
		}
	}
	return buffer.String()
}

func main() {
	root := createNode(1, nil)
	node10 := createNode(2, root)
	node11 := createNode(3, root)
	node12 := createNode(4, root)
	node20 := createNode(5, node10)
	createNode(6, node10)
	createNode(7, node10)
	createNode(8, node11)
	createNode(9, node11)
	createNode(10, node12)
	createNode(11, node12)
	createNode(12, node20)
	fmt.Println(calculate(root, 0))
	fmt.Println(getString(root, "", true))
}
