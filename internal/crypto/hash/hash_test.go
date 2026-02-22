package hash

import (
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	testCases := []struct {
		name     string
		data     []byte
		algo     Algorithm
		expected string
	}{
		{
			name:     "SHA256 empty",
			data:     []byte{},
			algo:     SHA256,
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "SHA256 hello",
			data:     []byte("hello"),
			algo:     SHA256,
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Hash(tc.data, tc.algo)
			if err != nil {
				t.Fatalf("Hash failed: %v", err)
			}

			resultHex := hex.EncodeToString(result)
			if resultHex != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, resultHex)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	data := []byte("test data")
	hash, err := Hash(data, SHA256)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	valid, err := Verify(data, hash, SHA256)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if !valid {
		t.Error("Verification should pass")
	}

	// Test with wrong data
	wrongData := []byte("wrong data")
	valid, err = Verify(wrongData, hash, SHA256)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if valid {
		t.Error("Verification should fail with wrong data")
	}
}
