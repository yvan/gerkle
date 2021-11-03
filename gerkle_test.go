package main

import (
	"testing"
	"crypto/sha256"
)

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

func TestNodeIsLeaf(t *testing.T) {
	data := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}	
	nodes := makeTree(data)
	test_hash := sha256.Sum256([]byte("some_utxo3"))
	check_node_root := nodes[0]
	check_node := nodes[0].right.right
	if (string(test_hash[:]) != string(check_node.hash)) {
		t.Errorf("right most node does not match")
	} else {
		if !nodeIsLeaf(*check_node) {
			t.Errorf("the node that should be a leaf is not a leaf, node: %v", *check_node)
		}
		if nodeIsLeaf(check_node_root) {
			t.Errorf("root node of multi node tree is beingclassified as leaf, should not be, node %v", check_node_root)
		}
	}
}

func TestVerifyTree(t *testing.T) {
	data_one := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}
	data_two := [][]byte{[]byte("some_utxo0_diff"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}
	data_three := [][]byte{[]byte("some_utxo0"),
		[]byte("some_utxo1"),
		[]byte("some_utxo2"),
		[]byte("some_utxo3")}	
	nodes_one := makeTree(data_one)
	nodes_two := makeTree(data_two)
	nodes_three := makeTree(data_three)

	// these are the tests for shallow verification at the root node
	
	// test that verification for two different trees is false
	verified, problems := verifyTree(nodes_two[0],nodes_one[0],false)
	if verified {
		t.Errorf("root verification returned true but should return false, root one: %v, root two: %v", nodes_two[0], nodes_one[0])
	}	
	// test that the number of problem leaf nodes is 1 if verification is false
	if len(problems) != 0 {
		t.Errorf("root verification should return 0 values, but returned %d", len(problems))
	}
	// test that verification for two identical trees is true
	verified, _ = verifyTree(nodes_three[0],nodes_one[0],false)
	if !verified {
		t.Errorf("root verification returned false but should return true, root one: %v, root two: %v", nodes_three[0], nodes_one[0])
	}

	// these are the tests for deep verification

	// deep verification of identical trees
	verified, problems = verifyTree(nodes_three[0],nodes_one[0],true)
	if !verified {
		t.Errorf("two identical trees should be verified but are not, tree one root: %v, tre two root: %v", nodes_three[0], nodes_one[0])
	}

	if len(problems) > 0 {
		t.Errorf("problems generated from two identical trees should be 0, but are %d, the problems: %v", len(problems), problems)
	}

	// deep verification of differing trees
	verified, problems = verifyTree(nodes_two[0],nodes_one[0],true)

	// different tree should not verify
	if verified {
		t.Errorf("two different trees should not verify, but they did, tree root one: %v, tree root two: %v", nodes_two[0], nodes_one[0] )
	}

	if len(problems) != 1 {
		t.Errorf("two trees with one node differing should return one problem node, but got %d problem nodes, the problems: %v", len(problems), problems)
	}
}
