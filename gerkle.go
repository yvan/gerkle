package main

import "crypto/sha256"
import "fmt"
import "math"

// a node always has two children
// unless it's a leaf then it has none
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
		l,r := Node{hash:nil,left:nil,right:nil},
		Node{hash:nil,left:nil,right:nil}
		nodes = append(nodes, makeNode(d,l,r))
		// fmt.Printf("processing datum %d in makeTree to node: ", i)
		// fmt.Println(nodes[i])
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
	//tree_info := make(map[int]map[string]int)
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

// prints a vertical tree
// https://stackoverflow.com/questions/13484943/print-a-binary-tree-in-a-pretty-way
func printTree(nodes []Node) {
	// so what do we need?
	// we need a sense of how "deep" the tree is vertically
	// and we need this for every "subtree," we do not need
	// ∟ ͱ —
	//
	
	// we want to recur down into the tree and then on the way back up
	// find the "depth," but this is the same as storing info
	// about the number of nodes before hand, so why don't we do
	// that and save some compute?
	
	
}

func main () {
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}
	nodes := makeTree(data)
	fmt.Println(nodes)
}
