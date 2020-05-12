// Reference: https://github.com/google/btree/blob/master/btree.go
package main

import (
	"fmt"
)

// Item represents a single object in the tree.
type Item interface {
	// Less tests whether the current item is less than the given argument.
	//
	// This must provide a strict weak ordering.
	// If !a.Less(b) && !b.Less(a), we treat this to mean a == b (i.e. we can only
	// hold one of either a or b in the tree).
	Less(than Item) bool
}

// items stores items in a node.
type items []Item

// node is an internal node in a tree.
type node struct {
	items items
	// children children
	// cow      *copyOnWriteContext
	depth int
	next  nextNodes
}

// // children stores child nodes in a node.
// type children []*node

// next stores next nodes in a node.
type nextNodes []*node

// BTree is an implementation of a B-Tree.
//
// BTree stores Item instances in an ordered structure, allowing easy insertion,
// removal, and iteration.
//
// Write operations are not safe for concurrent mutation by multiple
// goroutines, but Read operations are.
type BTree struct {
	degree int
	length int
	root   *node
	cow    *copyOnWriteContext
}
