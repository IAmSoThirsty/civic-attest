package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/signatures"
	"github.com/IAmSoThirsty/civic-attest/internal/identity/models"
)

func main() {
	var (
		officeID     = flag.String("office", "", "Office ID")
		jurisdiction = flag.String("jurisdiction", "", "Jurisdiction")
		validYears   = flag.Int("valid-years", 1, "Validity period in years")
		output       = flag.String("output", "", "Output file for identity")
	)

	flag.Parse()

	if *officeID == "" || *jurisdiction == "" {
		flag.Usage()
		log.Fatal("office and jurisdiction are required")
	}

	fmt.Println("=== Identity Authority ===")
	fmt.Println("Generating new identity...")

	// Generate key pair (in production, this would be in HSM)
	kp, err := signatures.GenerateKeyPair(signatures.Ed25519)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	// Create identity
	now := time.Now().UTC()
	validTo := now.AddDate(*validYears, 0, 0)

	identity := &models.Identity{
		OfficeID:     *officeID,
		Jurisdiction: *jurisdiction,
		PublicKey:    kp.PublicKey,
		KeyVersion:   1,
		ValidFrom:    now,
		ValidTo:      validTo,
		KeyAlgorithm: string(signatures.Ed25519),
		Status:       models.StatusActive,
		IdentityID:   fmt.Sprintf("%s-%s-v1", *officeID, *jurisdiction),
	}

	fmt.Printf("\nIdentity created:\n")
	fmt.Printf("  ID: %s\n", identity.IdentityID)
	fmt.Printf("  Office: %s\n", identity.OfficeID)
	fmt.Printf("  Jurisdiction: %s\n", identity.Jurisdiction)
	fmt.Printf("  Public Key: %s\n", hex.EncodeToString(identity.PublicKey))
	fmt.Printf("  Valid From: %s\n", identity.ValidFrom.Format(time.RFC3339))
	fmt.Printf("  Valid To: %s\n", identity.ValidTo.Format(time.RFC3339))
	fmt.Printf("  Algorithm: %s\n", identity.KeyAlgorithm)
	fmt.Printf("  Status: %s\n", identity.Status)

	if *output != "" {
		fmt.Printf("\nIdentity data would be written to: %s\n", *output)
	}

	fmt.Println("\n⚠ Note: In production, keys would be generated in HSM")
	fmt.Println("⚠ Note: Identity issuance requires trustee quorum")
}
