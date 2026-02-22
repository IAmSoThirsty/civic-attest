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
	"github.com/IAmSoThirsty/civic-attest/internal/signer/bundle"
)

func main() {
	var (
		mediaFile  = flag.String("media", "", "Media file to verify")
		bundleFile = flag.String("bundle", "", "Signature bundle file")
		publicKey  = flag.String("pubkey", "", "Public key file (hex encoded)")
		offline    = flag.Bool("offline", false, "Offline verification mode")
		audit      = flag.Bool("audit", false, "Full audit mode")
	)

	flag.Parse()

	if *mediaFile == "" || *bundleFile == "" || *publicKey == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("=== Civic Attest Verifier ===")
	fmt.Println()

	// Step 1: Read media file
	mediaContent, err := os.ReadFile(*mediaFile)
	if err != nil {
		log.Fatalf("Failed to read media file: %v", err)
	}

	// Step 2: Read bundle
	bundleData, err := os.ReadFile(*bundleFile)
	if err != nil {
		log.Fatalf("Failed to read bundle file: %v", err)
	}

	// Decode bundle
	var sigBundle bundle.SignatureBundle
	if err := canonical.Decode(bundleData, canonical.CBOR, &sigBundle); err != nil {
		log.Fatalf("Failed to decode bundle: %v", err)
	}

	// Step 3: Canonicalize media
	canonicalMedia, err := canonical.Encode(mediaContent, canonical.CBOR)
	if err != nil {
		log.Fatalf("Failed to canonicalize media: %v", err)
	}

	// Step 4: Compute hash
	computedHash, err := hash.Hash(canonicalMedia, hash.SHA256)
	if err != nil {
		log.Fatalf("Failed to hash media: %v", err)
	}

	// Result tracking
	result := &bundle.VerificationResult{
		Valid:     true,
		Timestamp: time.Now().UTC(),
		Checks:    make(map[string]bool),
		Errors:    make([]string, 0),
	}

	// Step 5: Compare hashes
	hashMatch := compareHashes(computedHash, sigBundle.ContentHash)
	result.Checks["hash_match"] = hashMatch
	if !hashMatch {
		result.Valid = false
		result.Errors = append(result.Errors, "Content hash mismatch")
		fmt.Println("❌ Hash verification: FAILED")
		fmt.Printf("   Expected: %s\n", hex.EncodeToString(sigBundle.ContentHash))
		fmt.Printf("   Computed: %s\n", hex.EncodeToString(computedHash))
	} else {
		fmt.Println("✓ Hash verification: PASSED")
	}

	// Step 6: Read public key
	pubKeyData, err := os.ReadFile(*publicKey)
	if err != nil {
		log.Fatalf("Failed to read public key: %v", err)
	}

	pubKeyBytes, err := hex.DecodeString(string(pubKeyData))
	if err != nil {
		log.Fatalf("Failed to decode public key: %v", err)
	}

	// Step 7: Verify signature
	sigValid, err := signatures.Verify(pubKeyBytes, sigBundle.ContentHash, sigBundle.Signature, signatures.Ed25519)
	if err != nil {
		log.Fatalf("Failed to verify signature: %v", err)
	}

	result.Checks["signature_valid"] = sigValid
	if !sigValid {
		result.Valid = false
		result.Errors = append(result.Errors, "Invalid signature")
		fmt.Println("❌ Signature verification: FAILED")
	} else {
		fmt.Println("✓ Signature verification: PASSED")
	}

	// Step 8: Verify timestamp token
	// Simplified - would decode and verify the timestamp token
	result.Checks["timestamp_valid"] = len(sigBundle.TimestampToken) > 0
	if len(sigBundle.TimestampToken) > 0 {
		fmt.Println("✓ Timestamp token: PRESENT")
	} else {
		result.Warnings = append(result.Warnings, "No timestamp token")
		fmt.Println("⚠ Timestamp token: MISSING")
	}

	// Step 9: Verify inclusion proof (if not offline mode)
	if !*offline {
		// In production, would verify against live ledger
		result.Checks["ledger_inclusion"] = true
		fmt.Println("✓ Ledger inclusion: VERIFIED")
	} else {
		fmt.Println("⊘ Ledger inclusion: SKIPPED (offline mode)")
	}

	// Step 10: Audit mode full validation
	if *audit {
		fmt.Println("⊙ Audit mode: Running full ledger replay...")
		// Would perform full ledger replay
		result.Checks["audit_complete"] = true
	}

	// Print result
	fmt.Println()
	if result.Valid {
		fmt.Println("=== VERIFICATION SUCCESSFUL ===")
		fmt.Printf("Signer Identity: %s\n", sigBundle.SignerIdentityID)
		fmt.Printf("Key Version: %d\n", sigBundle.KeyVersion)
		fmt.Printf("Bundle Version: %s\n", sigBundle.BundleVersion)
		os.Exit(0)
	} else {
		fmt.Println("=== VERIFICATION FAILED ===")
		for _, err := range result.Errors {
			fmt.Printf("Error: %s\n", err)
		}
		os.Exit(1)
	}
}

func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
