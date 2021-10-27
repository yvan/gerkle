package main

import "crypto/sha256"
import "testing"
//import "fmt"

// a test to see that 1- the function returns
// a node with a hash inside it and two null
// links
func TestLeaf(t *testing.T){
	data := []byte("some_utxo")
	var left Node
	var right Node
	var hash [32]byte
	left = Node{hash:nil,left:nil,right:nil,}
	right = Node{hash:nil,left:nil,right:nil,}
	hash = sha256.Sum256(data)
	n_test := Node{
		hash: hash[:],
		left: &left,
		right: &right,
		
	}		
	n := makeNode(data,left,right)
	n1 := n
	n2 := n_test
	// first the lenght of bytes
	if (len(n1.hash) != len(n2.hash)) {
		t.Errorf("hashes are a different length, %d and %d", len(n1.hash), len(n2.hash))
	}
	// the bytes must match
	for i := 0; i < len(n1.hash); i++ {
		if (n1.hash[i] != n2.hash[i]) {
			t.Errorf("hash bytes do not match for byte %d, found %x, expected %x ",i,n1.hash[i],n2.hash[i])
		}
	}
}

func checkTree(n Node, t *testing.T) {
	if (n.hash == nil) {
		return
	}
	node_hash := n.hash
	_temp := sha256.Sum256(append(n.left.hash,n.right.hash...))
	check_hash := _temp[:]
	for i := range node_hash {
		// if not a leaf node and the hashes don't match
		if (!(len(n.left.hash) == 0 && len(n.right.hash) == 0) && (node_hash[i] != check_hash[i])) {
			t.Errorf("node does not properly combine hashes of left and right. expected %x, got %x from left hash %x and right hash %x",node_hash, check_hash, n.left.hash, n.right.hash)
		}
	}
	checkTree(*n.left, t)
	checkTree(*n.right, t)
}


func TestTree (t *testing.T) {
	// make a tree using test data
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}
	nodes := makeTree(data)

	if (len(nodes) != 1) {
		t.Errorf("tree does not have a single root node")
	}
	
	// traverse the tree
	checkTree(nodes[0], t)	
}

func TestGetTreeLevels(t *testing.T){
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}	
	nodes := makeTree(data)
	correctTreeLevels := 3
	testTreeLevels := getTreeLevels(nodes[0])
	if (correctTreeLevels != testTreeLevels) {
		t.Errorf("getTreeLevels returns %x levels, should be %x", testTreeLevels, correctTreeLevels)
	}
}

func TestGetNumNodes(t *testing.T){
	var testTreeNodes int = 7
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}	
	nodes := makeTree(data)
	n := getNumNodes(nodes[0])
	if (n != testTreeNodes) {
		t.Errorf("tree should have %d nodes, but got %d", testTreeNodes, n)
	}
	
}

func TestGetNumNodesLevel(t *testing.T) {
	var correctNodes = 4
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}	
	nodes := makeTree(data)
	n := getNumNodesLevel(nodes[0], 3)
	if (n != correctNodes) {
		t.Errorf("tree level 3 should have %d nodes, but got %d", correctNodes, n)
	}
}
