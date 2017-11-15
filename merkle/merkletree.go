package merkle

import "crypto/sha256"

//Tree defined by the root
type Tree struct {
	Root *Node
}

//Node of dual SLL
type Node struct {
	Left  *Node
	Right *Node
	Data  []byte
}

func NewNode(left, right *Node, data []byte) (n *Node) {
	n = &Node{}

	if IsLeaf(left, right) {
		hash := sha256.Sum256(data)
		n.Data = hash[:]
	} else {
		prevHash := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHash)
		n.Data = hash[:]
	}
	n.Left = left
	n.Right = right
	return n
}

func IsLeaf(left, right *Node) bool {
	return left == nil && right == nil
}

func roundUpData(data [][]byte) [][]byte {
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	return data
}

func buildLeaves(data [][]byte) []*Node {
	var nodes []*Node

	for _, datum := range data {
		node := NewNode(nil, nil, datum)
		nodes = append(nodes, node)
	}
	return nodes
}

func NewTree(data [][]byte) (tree *Tree) {
	data = roundUpData(data)
	nodes := buildLeaves(data)
	for i := 0; i < len(data)/2; i++ {
		var newLevel []*Node

		for j := 0; j < len(nodes); j += 2 {
			node := NewNode(nodes[j], nodes[j+1], nil)
			newLevel = append(newLevel, node)
		}
		nodes = newLevel
	}

	tree = &Tree{nodes[0]}
	return tree
}
