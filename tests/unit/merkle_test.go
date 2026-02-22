package unit

import (
	"testing"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/hash"
	"github.com/IAmSoThirsty/civic-attest/internal/crypto/merkle"
)

func TestMerkleTree(t *testing.T) {
	tree := merkle.NewTree(hash.SHA256)

	// Add some leaves
	data := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	for _, d := range data {
		err := tree.Append(d)
		if err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	// Check size
	if tree.Size() != len(data) {
		t.Errorf("Expected size %d, got %d", len(data), tree.Size())
	}

	// Check root hash exists
	root := tree.RootHash()
	if len(root) == 0 {
		t.Error("Root hash is empty")
	}
}

func TestInclusionProof(t *testing.T) {
	tree := merkle.NewTree(hash.SHA256)

	// Add leaves
	for i := 0; i < 5; i++ {
		data := []byte{byte(i)}
		err := tree.Append(data)
		if err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	// Generate proof
	proof, err := tree.GenerateInclusionProof(2)
	if err != nil {
		t.Fatalf("Failed to generate proof: %v", err)
	}

	if proof.LeafIndex != 2 {
		t.Errorf("Expected leaf index 2, got %d", proof.LeafIndex)
	}

	if proof.TreeSize != 5 {
		t.Errorf("Expected tree size 5, got %d", proof.TreeSize)
	}
}
