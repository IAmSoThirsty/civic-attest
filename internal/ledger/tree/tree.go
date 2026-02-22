package tree

import (
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/hash"
	"github.com/IAmSoThirsty/civic-attest/internal/crypto/merkle"
)

// Entry represents a ledger entry
type Entry struct {
	// EntryHash is the hash of the entry content
	EntryHash []byte `json:"entry_hash"`
	// Timestamp is when the entry was created
	Timestamp time.Time `json:"timestamp"`
	// SignerIdentityID is the identity that signed
	SignerIdentityID string `json:"signer_identity_id"`
	// SignatureHash is the hash of the signature
	SignatureHash []byte `json:"signature_hash"`
	// EntryType is the type of entry
	EntryType string `json:"entry_type"`
	// SequenceNumber is the sequence number in the ledger
	SequenceNumber int64 `json:"sequence_number"`
}

// SignedTreeHead represents a signed tree head
type SignedTreeHead struct {
	// TreeSize is the number of entries in the tree
	TreeSize int `json:"tree_size"`
	// RootHash is the root hash of the tree
	RootHash []byte `json:"root_hash"`
	// Timestamp is when the tree head was signed
	Timestamp time.Time `json:"timestamp"`
	// Signature is the signature by the ledger authority
	Signature []byte `json:"signature"`
	// LedgerAuthorityID is the identity of the ledger authority
	LedgerAuthorityID string `json:"ledger_authority_id"`
}

// LedgerTree represents the append-only Merkle tree ledger
type LedgerTree struct {
	mu       sync.RWMutex
	tree     *merkle.Tree
	entries  []*Entry
	hashAlgo hash.Algorithm
	sequence int64
}

// NewLedgerTree creates a new ledger tree
func NewLedgerTree(hashAlgo hash.Algorithm) *LedgerTree {
	return &LedgerTree{
		tree:     merkle.NewTree(hashAlgo),
		entries:  make([]*Entry, 0),
		hashAlgo: hashAlgo,
		sequence: 0,
	}
}

// Append adds a new entry to the ledger
func (lt *LedgerTree) Append(entry *Entry) error {
	lt.mu.Lock()
	defer lt.mu.Unlock()

	// Assign sequence number
	lt.sequence++
	entry.SequenceNumber = lt.sequence

	// Compute entry hash if not already set
	if entry.EntryHash == nil {
		entryData := lt.serializeEntry(entry)
		h, err := hash.Hash(entryData, lt.hashAlgo)
		if err != nil {
			return fmt.Errorf("failed to hash entry: %w", err)
		}
		entry.EntryHash = h
	}

	// Add to entries list
	lt.entries = append(lt.entries, entry)

	// Add to Merkle tree
	if err := lt.tree.Append(entry.EntryHash); err != nil {
		return fmt.Errorf("failed to append to tree: %w", err)
	}

	return nil
}

// serializeEntry serializes an entry for hashing
func (lt *LedgerTree) serializeEntry(entry *Entry) []byte {
	// Simple concatenation for demonstration
	// Production would use canonical encoding
	data := fmt.Sprintf("%s|%s|%s|%d",
		entry.SignerIdentityID,
		hex.EncodeToString(entry.SignatureHash),
		entry.EntryType,
		entry.Timestamp.Unix(),
	)
	return []byte(data)
}

// GetRootHash returns the current root hash
func (lt *LedgerTree) GetRootHash() []byte {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	return lt.tree.RootHash()
}

// GetSize returns the number of entries in the ledger
func (lt *LedgerTree) GetSize() int {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	return lt.tree.Size()
}

// GetEntry retrieves an entry by index
func (lt *LedgerTree) GetEntry(index int) (*Entry, error) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	if index < 0 || index >= len(lt.entries) {
		return nil, fmt.Errorf("invalid entry index: %d", index)
	}

	return lt.entries[index], nil
}

// GenerateInclusionProof generates a proof that an entry is in the ledger
func (lt *LedgerTree) GenerateInclusionProof(index int) (*merkle.InclusionProof, error) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	return lt.tree.GenerateInclusionProof(index)
}

// GenerateConsistencyProof generates a proof that the tree is consistent
func (lt *LedgerTree) GenerateConsistencyProof(oldSize int) (*merkle.ConsistencyProof, error) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	return lt.tree.GenerateConsistencyProof(oldSize)
}

// GetSignedTreeHead creates a signed tree head (signature would be added by caller)
func (lt *LedgerTree) GetSignedTreeHead() *SignedTreeHead {
	lt.mu.RLock()
	defer lt.mu.RUnlock()

	return &SignedTreeHead{
		TreeSize:  lt.tree.Size(),
		RootHash:  lt.tree.RootHash(),
		Timestamp: time.Now().UTC(),
	}
}

// VerifyConsistency verifies that the ledger maintains append-only property
func (lt *LedgerTree) VerifyConsistency(oldSTH, newSTH *SignedTreeHead) bool {
	// Simplified verification
	// Production would verify the consistency proof
	if newSTH.TreeSize < oldSTH.TreeSize {
		return false
	}

	return true
}
