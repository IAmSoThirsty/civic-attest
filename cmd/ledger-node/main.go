package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/hash"
	"github.com/IAmSoThirsty/civic-attest/internal/ledger/tree"
)

// LedgerNode represents a ledger node server
type LedgerNode struct {
	ledger *tree.LedgerTree
	mu     sync.RWMutex
	port   string
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	node := &LedgerNode{
		ledger: tree.NewLedgerTree(hash.SHA256),
		port:   *port,
	}

	// Setup HTTP handlers
	http.HandleFunc("/health", node.healthHandler)
	http.HandleFunc("/tree-head", node.treeHeadHandler)
	http.HandleFunc("/append", node.appendHandler)
	http.HandleFunc("/entry/", node.entryHandler)
	http.HandleFunc("/inclusion-proof/", node.inclusionProofHandler)

	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("Ledger Node starting on %s\n", addr)
	fmt.Println("Endpoints:")
	fmt.Println("  GET  /health - Health check")
	fmt.Println("  GET  /tree-head - Get signed tree head")
	fmt.Println("  POST /append - Append new entry")
	fmt.Println("  GET  /entry/{index} - Get entry by index")
	fmt.Println("  GET  /inclusion-proof/{index} - Get inclusion proof")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (ln *LedgerNode) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func (ln *LedgerNode) treeHeadHandler(w http.ResponseWriter, r *http.Request) {
	ln.mu.RLock()
	defer ln.mu.RUnlock()

	sth := ln.ledger.GetSignedTreeHead()
	fmt.Fprintf(w, "Tree Size: %d\n", sth.TreeSize)
	fmt.Fprintf(w, "Root Hash: %s\n", hex.EncodeToString(sth.RootHash))
	fmt.Fprintf(w, "Timestamp: %s\n", sth.Timestamp.Format(time.RFC3339))
}

func (ln *LedgerNode) appendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Simplified append - would validate entry in production
	ln.mu.Lock()
	defer ln.mu.Unlock()

	entry := &tree.Entry{
		SignerIdentityID: "test-identity",
		SignatureHash:    []byte("test-hash"),
		EntryType:        "signature",
		Timestamp:        time.Now().UTC(),
	}

	if err := ln.ledger.Append(entry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Entry appended: %d\n", entry.SequenceNumber)
}

func (ln *LedgerNode) entryHandler(w http.ResponseWriter, r *http.Request) {
	// Parse index from URL
	var index int
	if _, err := fmt.Sscanf(r.URL.Path, "/entry/%d", &index); err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	ln.mu.RLock()
	defer ln.mu.RUnlock()

	entry, err := ln.ledger.GetEntry(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Sequence: %d\n", entry.SequenceNumber)
	fmt.Fprintf(w, "Identity: %s\n", entry.SignerIdentityID)
	fmt.Fprintf(w, "Timestamp: %s\n", entry.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "Hash: %s\n", hex.EncodeToString(entry.EntryHash))
}

func (ln *LedgerNode) inclusionProofHandler(w http.ResponseWriter, r *http.Request) {
	var index int
	if _, err := fmt.Sscanf(r.URL.Path, "/inclusion-proof/%d", &index); err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	ln.mu.RLock()
	defer ln.mu.RUnlock()

	proof, err := ln.ledger.GenerateInclusionProof(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Leaf Index: %d\n", proof.LeafIndex)
	fmt.Fprintf(w, "Tree Size: %d\n", proof.TreeSize)
	fmt.Fprintf(w, "Leaf Hash: %s\n", hex.EncodeToString(proof.LeafHash))
}
