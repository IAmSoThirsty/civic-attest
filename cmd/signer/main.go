package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/canonical"
	"github.com/IAmSoThirsty/civic-attest/internal/crypto/hash"
	"github.com/IAmSoThirsty/civic-attest/internal/crypto/signatures"
	"github.com/IAmSoThirsty/civic-attest/internal/crypto/timestamp"
	"github.com/IAmSoThirsty/civic-attest/internal/ledger/tree"
	"github.com/IAmSoThirsty/civic-attest/internal/signer/bundle"
)

func main() {
	var (
		inputFile    = flag.String("input", "", "Input file to sign")
		identityID   = flag.String("identity", "", "Signer identity ID")
		keyFile      = flag.String("key", "", "Private key file (PEM format)")
		outputFile   = flag.String("output", "", "Output signature bundle file")
		canonFormat  = flag.String("canon", "CBOR", "Canonical format (CBOR or JSON)")
	)

	flag.Parse()

	if *inputFile == "" || *identityID == "" || *keyFile == "" || *outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Step 1: Read master artifact
	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	// Step 2: Canonicalize
	var format canonical.Format
	switch *canonFormat {
	case "CBOR":
		format = canonical.CBOR
	case "JSON":
		format = canonical.JSON
	default:
		log.Fatalf("Unsupported canonical format: %s", *canonFormat)
	}

	canonicalData, err := canonical.Encode(content, format)
	if err != nil {
		log.Fatalf("Failed to canonicalize content: %v", err)
	}

	// Step 3: Hash canonical artifact
	contentHash, err := hash.Hash(canonicalData, hash.SHA256)
	if err != nil {
		log.Fatalf("Failed to hash content: %v", err)
	}

	fmt.Printf("Content hash: %s\n", hex.EncodeToString(contentHash))

	// Step 4-5: Sign with HSM (simulated with file key for demo)
	// In production, this would interface with actual HSM
	privateKey, err := loadPrivateKey(*keyFile)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	signature, err := signatures.Sign(privateKey, contentHash, signatures.Ed25519)
	if err != nil {
		log.Fatalf("Failed to sign: %v", err)
	}

	fmt.Printf("Signature: %s\n", hex.EncodeToString(signature))

	// Step 6: Request timestamp
	tsaClient := timestamp.NewMockTSAClient()
	tsToken, err := tsaClient.Request(contentHash, string(hash.SHA256))
	if err != nil {
		log.Fatalf("Failed to get timestamp: %v", err)
	}

	tsData, err := tsToken.Encode()
	if err != nil {
		log.Fatalf("Failed to encode timestamp: %v", err)
	}

	// Step 7: Append to ledger (simulated)
	ledger := tree.NewLedgerTree(hash.SHA256)
	entry := &tree.Entry{
		SignerIdentityID: *identityID,
		SignatureHash:    contentHash,
		EntryType:        "signature",
		Timestamp:        time.Now().UTC(),
	}

	if err := ledger.Append(entry); err != nil {
		log.Fatalf("Failed to append to ledger: %v", err)
	}

	// Step 8: Generate inclusion proof
	proof, err := ledger.GenerateInclusionProof(0)
	if err != nil {
		log.Fatalf("Failed to generate inclusion proof: %v", err)
	}

	// Step 9: Create signature bundle
	bundleData := &bundle.SignatureBundle{
		ContentHash:            contentHash,
		ContentHashAlgorithm:   string(hash.SHA256),
		CanonicalFormatVersion: "1.0",
		SignerIdentityID:       *identityID,
		KeyVersion:             1,
		Signature:              signature,
		TimestampToken:         tsData,
		LedgerEntryHash:        entry.EntryHash,
		MerkleInclusionProof: &bundle.InclusionProof{
			LeafIndex: proof.LeafIndex,
			LeafHash:  proof.LeafHash,
			TreeSize:  proof.TreeSize,
			Path:      proof.Path,
		},
		BundleVersion: "1.0",
	}

	// Encode bundle
	bundleBytes, err := canonical.Encode(bundleData, canonical.CBOR)
	if err != nil {
		log.Fatalf("Failed to encode bundle: %v", err)
	}

	// Write bundle
	if err := os.WriteFile(*outputFile, bundleBytes, 0644); err != nil {
		log.Fatalf("Failed to write bundle: %v", err)
	}

	fmt.Printf("Signature bundle written to: %s\n", *outputFile)
	fmt.Printf("Ledger entry hash: %s\n", hex.EncodeToString(entry.EntryHash))
}

func loadPrivateKey(filename string) ([]byte, error) {
	// For demo purposes, generate a new key if file doesn't exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		kp, err := signatures.GenerateKeyPair(signatures.Ed25519)
		if err != nil {
			return nil, err
		}
		return kp.PrivateKey, nil
	}

	return os.ReadFile(filename)
}
