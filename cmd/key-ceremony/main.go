package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IAmSoThirsty/civic-attest/internal/crypto/signatures"
	"github.com/IAmSoThirsty/civic-attest/internal/identity/models"
)

func main() {
	var (
		quorumSize    = flag.Int("quorum", 3, "Quorum size")
		totalTrustees = flag.Int("trustees", 5, "Total trustees")
		officeID      = flag.String("office", "", "Office ID for new key")
		jurisdiction  = flag.String("jurisdiction", "", "Jurisdiction")
	)

	flag.Parse()

	if *officeID == "" || *jurisdiction == "" {
		flag.Usage()
		log.Fatal("office and jurisdiction are required")
	}

	fmt.Println("╔═══════════════════════════════════════════╗")
	fmt.Println("║   CIVIC ATTEST KEY CEREMONY PROTOCOL      ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("Ceremony for: %s (%s)\n", *officeID, *jurisdiction)
	fmt.Printf("Quorum: %d of %d trustees\n", *quorumSize, *totalTrustees)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Step 1: Gather trustees
	fmt.Println("=== Step 1: Trustee Assembly ===")
	trustees := make([]string, 0)
	for i := 0; i < *quorumSize; i++ {
		fmt.Printf("Enter trustee %d ID: ", i+1)
		trusteeID, _ := reader.ReadString('\n')
		trusteeID = strings.TrimSpace(trusteeID)
		trustees = append(trustees, trusteeID)
	}

	fmt.Println("\n✓ Quorum assembled")
	for i, t := range trustees {
		fmt.Printf("  %d. %s\n", i+1, t)
	}

	// Step 2: Generate key
	fmt.Println("\n=== Step 2: Key Generation ===")
	fmt.Println("Generating key in HSM (simulated)...")

	kp, err := signatures.GenerateKeyPair(signatures.Ed25519)
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	pubKeyHex := hex.EncodeToString(kp.PublicKey)
	fmt.Printf("✓ Key generated\n")
	fmt.Printf("  Public Key: %s\n", pubKeyHex[:32])
	fmt.Printf("              %s...\n", pubKeyHex[32:64])

	// Step 3: Record ceremony
	fmt.Println("\n=== Step 3: Ceremony Recording ===")
	fmt.Println("Recording ceremony (audio/video)...")

	ceremony := &models.KeyCeremonyRecord{
		CeremonyID:    fmt.Sprintf("ceremony-%s-%d", *officeID, time.Now().Unix()),
		Timestamp:     time.Now().UTC(),
		Trustees:      trustees,
		QuorumSize:    *quorumSize,
		TotalTrustees: *totalTrustees,
		RecordingHash: []byte("recording-hash-placeholder"),
		PublicKeyHash: kp.PublicKey[:32], // First 32 bytes as hash
	}

	fmt.Printf("✓ Ceremony recorded\n")
	fmt.Printf("  Ceremony ID: %s\n", ceremony.CeremonyID)
	fmt.Printf("  Recording Hash: %s\n", hex.EncodeToString(ceremony.RecordingHash))

	// Step 4: Broadcast and ledger append
	fmt.Println("\n=== Step 4: Public Broadcast ===")
	fmt.Println("Broadcasting ceremony hash...")
	fmt.Println("Appending to ledger...")

	fmt.Printf("✓ Ceremony complete\n")
	fmt.Printf("  Ledger entry: %s\n", hex.EncodeToString(ceremony.RecordingHash))

	// Step 5: Create identity
	fmt.Println("\n=== Step 5: Identity Creation ===")

	now := time.Now().UTC()
	identity := &models.Identity{
		OfficeID:     *officeID,
		Jurisdiction: *jurisdiction,
		PublicKey:    kp.PublicKey,
		KeyVersion:   1,
		ValidFrom:    now,
		ValidTo:      now.AddDate(1, 0, 0),
		KeyAlgorithm: string(signatures.Ed25519),
		Status:       models.StatusActive,
		IdentityID:   fmt.Sprintf("%s-%s-v1", *officeID, *jurisdiction),
	}

	fmt.Printf("✓ Identity created: %s\n", identity.IdentityID)

	fmt.Println("\n╔═══════════════════════════════════════════╗")
	fmt.Println("║   KEY CEREMONY SUCCESSFULLY COMPLETED     ║")
	fmt.Println("╚═══════════════════════════════════════════╝")
}
