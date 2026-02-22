package signatures

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Algorithm represents a signature algorithm
type Algorithm string

const (
	// Ed25519 is the primary signature algorithm
	Ed25519 Algorithm = "Ed25519"
	// Ed448 is the secondary profile (reserved for future implementation)
	Ed448 Algorithm = "Ed448"
)

// KeyPair represents a public/private key pair
type KeyPair struct {
	PublicKey  []byte
	PrivateKey []byte
	Algorithm  Algorithm
}

// GenerateKeyPair generates a new key pair for the specified algorithm
func GenerateKeyPair(algo Algorithm) (*KeyPair, error) {
	switch algo {
	case Ed25519:
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Ed25519 key: %w", err)
		}
		return &KeyPair{
			PublicKey:  pub,
			PrivateKey: priv,
			Algorithm:  Ed25519,
		}, nil
	case Ed448:
		return nil, fmt.Errorf("Ed448 not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported signature algorithm: %s", algo)
	}
}

// Sign creates a signature for the given message
func Sign(privateKey []byte, message []byte, algo Algorithm) ([]byte, error) {
	switch algo {
	case Ed25519:
		if len(privateKey) != ed25519.PrivateKeySize {
			return nil, fmt.Errorf("invalid Ed25519 private key size")
		}
		signature := ed25519.Sign(privateKey, message)
		return signature, nil
	case Ed448:
		return nil, fmt.Errorf("Ed448 not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported signature algorithm: %s", algo)
	}
}

// Verify verifies a signature against a message and public key
func Verify(publicKey []byte, message []byte, signature []byte, algo Algorithm) (bool, error) {
	switch algo {
	case Ed25519:
		if len(publicKey) != ed25519.PublicKeySize {
			return false, fmt.Errorf("invalid Ed25519 public key size")
		}
		if len(signature) != ed25519.SignatureSize {
			return false, fmt.Errorf("invalid Ed25519 signature size")
		}
		return ed25519.Verify(publicKey, message, signature), nil
	case Ed448:
		return false, fmt.Errorf("Ed448 not yet implemented")
	default:
		return false, fmt.Errorf("unsupported signature algorithm: %s", algo)
	}
}

// PublicKeyString returns hex-encoded public key
func PublicKeyString(publicKey []byte) string {
	return hex.EncodeToString(publicKey)
}

// SignatureString returns hex-encoded signature
func SignatureString(signature []byte) string {
	return hex.EncodeToString(signature)
}
