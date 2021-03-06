package main

import (
	"math"
	"crypto/sha256"
)

// make a simple stack type
type stack []Node

func Push(stack []Node, n Node) ([]Node) {
	return append(stack,n)
}

func Pop(stack []Node) (Node, []Node) {
	return stack[0],stack[1:]
}

// a node -- has two children with hash values
// a leaf -- has two children with nil hash values
type Node struct {
	hash []byte
	left *Node
	right *Node
}

// make a node for the merkle tree
func makeNode(data []byte, left Node, right Node) (n Node) {
	var hash [32]byte
	// if this is a leaf
	if len(data) != 0 {
		// create a hash from the data
		hash = sha256.Sum256(data)
	// if this is a node (no raw data)
	} else {
		// create a hash from left and right node hashes
		hash = sha256.Sum256(append(left.hash, right.hash...))
	}
	// assign the children that were passed in
	n = Node{
		hash: hash[:],
		left: &left,
		right: &right,
	}
	return
}

// returns a list of leaf nodes at the bottom of the tree
func makeLeaves(data [][]byte)(nodes []Node){
	//take the nodes by 1 and made nodes
	for _, d := range data {
		l,r := Node{hash:nil,left:nil,right:nil},Node{hash:nil,left:nil,right:nil}
		nodes = append(nodes, makeNode(d,l,r))
	}

	return
}


// returns a list of nodes with only the root node
func makeParents(nodes []Node) (new_nodes []Node){
	//take the nodes in pairs and make new nodes
	for i := 0; i < len(nodes)-1; i=i+2 {
		n1 := nodes[i]
		n2 := nodes[i+1]		
		new_nodes = append(new_nodes, makeNode(nil, n1, n2))
	}
	return
}

func makeTree(data [][]byte) (root []Node){
	nodes := makeLeaves(data)
	for len(nodes) > 1 {
	 	nodes = makeParents(nodes)
	}
	return nodes
}

func getTreeLevels(root Node)(int){
	if (root.hash != nil) {
		return 1+getTreeLevels(*root.left)
	}
	return 0
}

func getNumNodes(tree Node) (numNodes int) {
	numLevels := getTreeLevels(tree)
	numNodes = int(math.Pow(2, float64(numLevels)) - 1)
	return
}

func getNumNodesLevel(tree Node, level int) (numNodes int) {
	numLevels := getTreeLevels(tree)
	if (level > numLevels) {
		return 0
	} else {
		return int(math.Pow(2,float64(level)) - math.Pow(2,float64(level-1)))
	}
}

func nodeIsLeaf(node Node) (bool) {
	return (node.left.hash == nil) && (node.right.hash == nil)
}

// verifies a tree by checking the hashes inside it
// can do a fast check, or tell you which leaves
// are responsible.
func verifyTree(new_tree Node, old_tree Node, deep bool) (bool, []Node) {
	// check that the trees have the same number of nodes
	numNew := getNumNodes(new_tree)
	numOld := getNumNodes(old_tree)
	if numNew != numOld {
		return false, nil
	}
	// if the trees have the same number of nodes
	// then check each node against its analog
	if deep {		
		var problem_nodes []Node
		var stack_old stack
		var stack_new stack
		stack_old = Push(stack_old, old_tree)
		stack_new = Push(stack_new, new_tree)
		var node_old Node
		var node_new Node

		// while there's nodes in the stack
		for (len(stack_old) > 0) && (len(stack_new) > 0) {
			// get two nodes
			node_old, stack_old = Pop(stack_old)
			node_new, stack_new = Pop(stack_new)
			// if we are not a nil node (nil nodes are not
			// part of the tree
			if (node_new.left != nil) && (node_new.right != nil) {
				
				// push it's children to the stack
				stack_old = Push(stack_old, *node_old.left)
				stack_old = Push(stack_old, *node_old.right)
				stack_new = Push(stack_new, *node_new.left)
				stack_new = Push(stack_new, *node_new.right)
				
				// examine the nodes
				if string(node_old.hash) != string(node_new.hash) {
					problem_nodes = append(problem_nodes,node_new)
				}
			}
		}

		// if we have problem nodes
		// get only the leaves
		if len(problem_nodes) > 0 {			
			var problem_leaves []Node	
			for _,n := range problem_nodes {
				if nodeIsLeaf(n){
					problem_leaves = append(problem_leaves,n)
				}
			}			
			return false, problem_leaves
		} else {
			return true, nil
		}
	} else {
		return string(new_tree.hash) == string(old_tree.hash), nil
	}	
}

func main () {

}

