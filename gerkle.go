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
func printTree(node Node) {
	
	// elbow := "∟"
	//tee := "ͱ"
	//dash := "—"
	vert := "|"
	numNodes := getNumNodes(node)
	to_print := make([]string, numNodes)

	temp_node := node
	// deal with root
	to_print = append(to_print,string(node.hash))	

	count := 0
	for temp_node.right.hash != nil {
		temp_str := ""
		for i := 0; i < count; i++ {
			temp_str += vert+"   "
		}
		fmt.Println(vert)
		// go right
		to_print = append(to_print,temp_str+string(temp_node.right.hash))
		temp_node = *temp_node.right
		count += 1
	}

	for i,l := range to_print {
		i = i
		fmt.Printf("%x\n",l)
	}
	
	// go left
	//node.left.hash

	//recur

	
	// first print the root
	//fmt.Printf("%x.\n",string(node.hash[:2]))
	// print a tee
	//fmt.Println(tee+dash+dash+dash)
	// 
	
	
	// so what do we need?
	// we need a sense of how deep the tree is vertically
	// and we need this for every subtree
	
	
}

func main () {
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}
	nodes := makeTree(data)
	printTree(nodes[0])
}
