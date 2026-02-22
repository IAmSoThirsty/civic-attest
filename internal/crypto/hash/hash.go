package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/zeebo/blake3"
	"golang.org/x/crypto/sha3"
)

// Algorithm represents a cryptographic hash algorithm
type Algorithm string

const (
	// SHA256 is the default hash algorithm
	SHA256 Algorithm = "SHA-256"
	// SHA3_512 is the long-term hash profile
	SHA3_512 Algorithm = "SHA-3-512"
	// BLAKE3 is the high-throughput mode
	BLAKE3 Algorithm = "BLAKE3"
)

// Hash computes a cryptographic hash of the input data using the specified algorithm
func Hash(data []byte, algo Algorithm) ([]byte, error) {
	switch algo {
	case SHA256:
		h := sha256.Sum256(data)
		return h[:], nil
	case SHA3_512:
		h := sha3.Sum512(data)
		return h[:], nil
	case BLAKE3:
		h := blake3.Sum256(data)
		return h[:], nil
	default:
		return nil, fmt.Errorf("unsupported hash algorithm: %s", algo)
	}
}

// HashString returns the hex-encoded hash
func HashString(data []byte, algo Algorithm) (string, error) {
	h, err := Hash(data, algo)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h), nil
}

// Verify checks if the hash matches the expected value
func Verify(data []byte, expectedHash []byte, algo Algorithm) (bool, error) {
	computed, err := Hash(data, algo)
	if err != nil {
		return false, err
	}

	if len(computed) != len(expectedHash) {
		return false, nil
	}

	for i := range computed {
		if computed[i] != expectedHash[i] {
			return false, nil
		}
	}

	return true, nil
}
