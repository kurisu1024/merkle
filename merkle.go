package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// Hashable is a constraint for types that can be hashed
type Hashable interface {
	~[]byte | ~string
}

// Tree represents a Merkle tree stored as a flat array
// Nodes are stored in level-order (breadth-first)
// For a node at index i: left child = 2*i+1, right child = 2*i+2
type Tree[T Hashable] struct {
	Nodes      []string // Hash values stored in level-order
	LeafData   []T      // Original data for leaves
	LeafOffset int      // Index where leaves start in Nodes array
	LeafCount  int      // Number of leaves
}

// hashData converts data to a hash string
func hashData[T Hashable](data T) string {
	var bytes []byte
	switch any(data).(type) {
	case []byte:
		bytes = any(data).([]byte)
	case string:
		bytes = []byte(any(data).(string))
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

// hashNodes combines two hashes and returns the hash of the combination
func hashNodes(left, right string) string {
	combined := left + right
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// nextPowerOfTwo returns the next power of 2 >= n
func nextPowerOfTwo(n int) int {
	if n <= 0 {
		return 1
	}
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}

// NewTree creates a new Merkle tree from the given data
func NewTree[T Hashable](data []T) (*Tree[T], error) {
	if len(data) == 0 {
		return nil, errors.New("cannot create tree with empty data")
	}

	leafCount := len(data)
	// Pad to next power of 2 for a complete binary tree
	paddedLeafCount := nextPowerOfTwo(leafCount)

	// Total nodes = paddedLeafCount * 2 - 1 (complete binary tree)
	totalNodes := paddedLeafCount*2 - 1
	leafOffset := paddedLeafCount - 1

	tree := &Tree[T]{
		Nodes:      make([]string, totalNodes),
		LeafData:   make([]T, leafCount),
		LeafOffset: leafOffset,
		LeafCount:  leafCount,
	}

	copy(tree.LeafData, data)

	// Hash all leaves (including padding)
	for i := 0; i < paddedLeafCount; i++ {
		if i < leafCount {
			tree.Nodes[leafOffset+i] = hashData(data[i])
		} else {
			// Duplicate last leaf for padding
			tree.Nodes[leafOffset+i] = tree.Nodes[leafOffset+leafCount-1]
		}
	}

	// Build tree bottom-up
	for i := leafOffset - 1; i >= 0; i-- {
		left := 2*i + 1
		right := 2*i + 2
		tree.Nodes[i] = hashNodes(tree.Nodes[left], tree.Nodes[right])
	}

	return tree, nil
}

// GetRoot returns the root hash of the tree
func (t *Tree[T]) GetRoot() string {
	if len(t.Nodes) == 0 {
		return ""
	}
	return t.Nodes[0]
}

// GetProof generates a Merkle proof for the data at the given index
func (t *Tree[T]) GetProof(index int) ([]string, error) {
	if index < 0 || index >= t.LeafCount {
		return nil, errors.New("index out of range")
	}

	var proof []string
	currentIndex := t.LeafOffset + index

	for currentIndex > 0 {
		// Find sibling
		var siblingIndex int
		if currentIndex%2 == 1 {
			// Current is left child, sibling is right
			siblingIndex = currentIndex + 1
		} else {
			// Current is right child, sibling is left
			siblingIndex = currentIndex - 1
		}

		proof = append(proof, t.Nodes[siblingIndex])

		// Move to parent
		currentIndex = (currentIndex - 1) / 2
	}

	return proof, nil
}

// VerifyProof verifies a Merkle proof for the given data
func VerifyProof[T Hashable](data T, proof []string, rootHash string, index int) bool {
	hash := hashData(data)

	for _, siblingHash := range proof {
		if index%2 == 0 {
			hash = hashNodes(hash, siblingHash)
		} else {
			hash = hashNodes(siblingHash, hash)
		}
		index = index / 2
	}

	return hash == rootHash
}

// Print prints the tree structure
func (t *Tree[T]) Print() {
	if len(t.Nodes) == 0 {
		fmt.Println("Empty tree")
		return
	}
	t.printNode(0, "", true)
}

func (t *Tree[T]) printNode(index int, prefix string, isTail bool) {
	if index >= len(t.Nodes) {
		return
	}

	connector := "└── "
	if !isTail {
		connector = "├── "
	}

	fmt.Printf("%s%s%s\n", prefix, connector, t.Nodes[index][:8]+"...")

	leftChild := 2*index + 1
	rightChild := 2*index + 2

	if leftChild < len(t.Nodes) || rightChild < len(t.Nodes) {
		extension := "    "
		if !isTail {
			extension = "│   "
		}

		if rightChild < len(t.Nodes) {
			t.printNode(rightChild, prefix+extension, false)
		}
		if leftChild < len(t.Nodes) {
			t.printNode(leftChild, prefix+extension, true)
		}
	}
}
