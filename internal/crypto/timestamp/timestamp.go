package timestamp

import (
	"encoding/asn1"
	"time"
)

// Token represents an RFC 3161 compliant timestamp token
type Token struct {
	// Version of the token format
	Version int
	// GenTime is the time at which the timestamp was generated
	GenTime time.Time
	// MessageImprint is the hash of the timestamped message
	MessageImprint []byte
	// HashAlgorithm used for the message imprint
	HashAlgorithm string
	// SerialNumber unique to this timestamp
	SerialNumber int64
	// TSA identifier
	TSA string
}

// TSAClient represents a timestamp authority client interface
type TSAClient interface {
	// Request requests a timestamp token for the given hash
	Request(messageHash []byte, hashAlgo string) (*Token, error)
}

// MockTSAClient is a mock implementation for testing
type MockTSAClient struct {
	counter int64
}

// NewMockTSAClient creates a new mock TSA client
func NewMockTSAClient() *MockTSAClient {
	return &MockTSAClient{counter: 1}
}

// Request implements TSAClient for testing
func (m *MockTSAClient) Request(messageHash []byte, hashAlgo string) (*Token, error) {
	token := &Token{
		Version:        1,
		GenTime:        time.Now().UTC(),
		MessageImprint: messageHash,
		HashAlgorithm:  hashAlgo,
		SerialNumber:   m.counter,
		TSA:            "mock-tsa",
	}
	m.counter++
	return token, nil
}

// Encode encodes the timestamp token to ASN.1 DER format
func (t *Token) Encode() ([]byte, error) {
	type tokenASN1 struct {
		Version        int
		GenTime        time.Time
		MessageImprint []byte
		HashAlgorithm  asn1.ObjectIdentifier
		SerialNumber   int64
		TSA            string
	}

	// For simplicity, using SHA-256 OID
	sha256OID := asn1.ObjectIdentifier{2, 16, 840, 1, 101, 3, 4, 2, 1}

	asn1Token := tokenASN1{
		Version:        t.Version,
		GenTime:        t.GenTime,
		MessageImprint: t.MessageImprint,
		HashAlgorithm:  sha256OID,
		SerialNumber:   t.SerialNumber,
		TSA:            t.TSA,
	}

	return asn1.Marshal(asn1Token)
}

// Verify verifies that the timestamp token is valid for the given message hash
func (t *Token) Verify(messageHash []byte, hashAlgo string) bool {
	if t.HashAlgorithm != hashAlgo {
		return false
	}

	if len(t.MessageImprint) != len(messageHash) {
		return false
	}

	for i := range t.MessageImprint {
		if t.MessageImprint[i] != messageHash[i] {
			return false
		}
	}

	return true
}
