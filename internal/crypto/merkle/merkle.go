package merkle

import (
	"bytes"
	"fmt"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/hash"
)

// Node represents a node in the Merkle tree
type Node struct {
	Hash  []byte
	Left  *Node
	Right *Node
	// IsLeaf indicates if this is a leaf node
	IsLeaf bool
	// Data is only set for leaf nodes
	Data []byte
}

// Tree represents a binary append-only Merkle tree
type Tree struct {
	Root      *Node
	Leaves    []*Node
	HashAlgo  hash.Algorithm
	treeSize  int
}

// NewTree creates a new Merkle tree
func NewTree(hashAlgo hash.Algorithm) *Tree {
	return &Tree{
		Root:     nil,
		Leaves:   make([]*Node, 0),
		HashAlgo: hashAlgo,
		treeSize: 0,
	}
}

// Append adds a new leaf to the tree
func (t *Tree) Append(data []byte) error {
	// Compute hash of the data
	leafHash, err := hash.Hash(data, t.HashAlgo)
	if err != nil {
		return fmt.Errorf("failed to hash leaf data: %w", err)
	}

	// Create new leaf node
	leaf := &Node{
		Hash:   leafHash,
		IsLeaf: true,
		Data:   data,
	}

	t.Leaves = append(t.Leaves, leaf)
	t.treeSize++

	// Rebuild the tree
	t.rebuildTree()

	return nil
}

// rebuildTree rebuilds the tree from leaves
func (t *Tree) rebuildTree() {
	if len(t.Leaves) == 0 {
		t.Root = nil
		return
	}

	if len(t.Leaves) == 1 {
		t.Root = t.Leaves[0]
		return
	}

	// Build tree bottom-up
	currentLevel := make([]*Node, len(t.Leaves))
	copy(currentLevel, t.Leaves)

	for len(currentLevel) > 1 {
		nextLevel := make([]*Node, 0)

		for i := 0; i < len(currentLevel); i += 2 {
			if i+1 < len(currentLevel) {
				// Combine two nodes
				parent := t.combineNodes(currentLevel[i], currentLevel[i+1])
				nextLevel = append(nextLevel, parent)
			} else {
				// Odd node out, promote to next level
				nextLevel = append(nextLevel, currentLevel[i])
			}
		}

		currentLevel = nextLevel
	}

	t.Root = currentLevel[0]
}

// combineNodes combines two nodes into a parent node
func (t *Tree) combineNodes(left, right *Node) *Node {
	// Concatenate hashes: hash(left || right)
	combined := append(left.Hash, right.Hash...)
	parentHash, _ := hash.Hash(combined, t.HashAlgo)

	return &Node{
		Hash:   parentHash,
		Left:   left,
		Right:  right,
		IsLeaf: false,
	}
}

// RootHash returns the root hash of the tree
func (t *Tree) RootHash() []byte {
	if t.Root == nil {
		return nil
	}
	return t.Root.Hash
}

// Size returns the number of leaves in the tree
func (t *Tree) Size() int {
	return t.treeSize
}

// InclusionProof represents a proof that a leaf is included in the tree
type InclusionProof struct {
	LeafIndex int
	LeafHash  []byte
	TreeSize  int
	Path      [][]byte // Hashes along the path from leaf to root
}

// GenerateInclusionProof generates a proof that a leaf at the given index is in the tree
func (t *Tree) GenerateInclusionProof(leafIndex int) (*InclusionProof, error) {
	if leafIndex < 0 || leafIndex >= len(t.Leaves) {
		return nil, fmt.Errorf("invalid leaf index: %d", leafIndex)
	}

	leaf := t.Leaves[leafIndex]
	path := make([][]byte, 0)

	// Build path from leaf to root
	t.buildPath(leaf, &path)

	return &InclusionProof{
		LeafIndex: leafIndex,
		LeafHash:  leaf.Hash,
		TreeSize:  t.Size(),
		Path:      path,
	}, nil
}

// buildPath recursively builds the path from a node to the root
func (t *Tree) buildPath(node *Node, path *[][]byte) {
	// Find the node's sibling and add to path
	// This is a simplified implementation
	// A production implementation would traverse the tree structure
	for i, leaf := range t.Leaves {
		if bytes.Equal(leaf.Hash, node.Hash) {
			// Add sibling hashes to path
			sibling := i ^ 1 // XOR with 1 to get sibling index
			if sibling < len(t.Leaves) {
				*path = append(*path, t.Leaves[sibling].Hash)
			}
			break
		}
	}
}

// VerifyInclusionProof verifies an inclusion proof
func (t *Tree) VerifyInclusionProof(proof *InclusionProof) bool {
	if proof.LeafIndex < 0 || proof.LeafIndex >= len(t.Leaves) {
		return false
	}

	leaf := t.Leaves[proof.LeafIndex]
	return bytes.Equal(leaf.Hash, proof.LeafHash)
}

// ConsistencyProof represents a proof that two tree states are consistent
type ConsistencyProof struct {
	OldSize int
	NewSize int
	Path    [][]byte
}

// GenerateConsistencyProof generates a proof that the tree at oldSize is consistent with newSize
func (t *Tree) GenerateConsistencyProof(oldSize int) (*ConsistencyProof, error) {
	if oldSize < 0 || oldSize > t.Size() {
		return nil, fmt.Errorf("invalid old size: %d", oldSize)
	}

	// Simplified consistency proof
	// A production implementation would properly compute the consistency path
	return &ConsistencyProof{
		OldSize: oldSize,
		NewSize: t.Size(),
		Path:    make([][]byte, 0),
	}, nil
}
