package signatures

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	if len(kp.PublicKey) == 0 {
		t.Error("Public key is empty")
	}

	if len(kp.PrivateKey) == 0 {
		t.Error("Private key is empty")
	}

	if kp.Algorithm != Ed25519 {
		t.Errorf("Expected algorithm %s, got %s", Ed25519, kp.Algorithm)
	}
}

func TestSignAndVerify(t *testing.T) {
	// Generate key pair
	kp, err := GenerateKeyPair(Ed25519)
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	// Message to sign
	message := []byte("test message")

	// Sign
	signature, err := Sign(kp.PrivateKey, message, Ed25519)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	// Verify
	valid, err := Verify(kp.PublicKey, message, signature, Ed25519)
	if err != nil {
		t.Fatalf("Failed to verify: %v", err)
	}

	if !valid {
		t.Error("Signature verification failed")
	}

	// Test with wrong message
	wrongMessage := []byte("wrong message")
	valid, err = Verify(kp.PublicKey, wrongMessage, signature, Ed25519)
	if err != nil {
		t.Fatalf("Failed to verify: %v", err)
	}

	if valid {
		t.Error("Signature should not verify with wrong message")
	}
}
