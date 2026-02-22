package bundle

import (
	"time"
)

// SignatureBundle represents the complete signature bundle format
type SignatureBundle struct {
	// ContentHash is the hash of the canonical content
	ContentHash []byte `json:"content_hash" cbor:"1,keyasint"`
	// ContentHashAlgorithm is the algorithm used for content hash
	ContentHashAlgorithm string `json:"content_hash_algorithm" cbor:"2,keyasint"`
	// CanonicalFormatVersion is the version of the canonicalization format
	CanonicalFormatVersion string `json:"canonical_format_version" cbor:"3,keyasint"`
	// SignerIdentityID is the identity that created the signature
	SignerIdentityID string `json:"signer_identity_id" cbor:"4,keyasint"`
	// KeyVersion is the version of the key used
	KeyVersion int `json:"key_version" cbor:"5,keyasint"`
	// Signature is the cryptographic signature
	Signature []byte `json:"signature" cbor:"6,keyasint"`
	// TimestampToken is the RFC 3161 timestamp token
	TimestampToken []byte `json:"timestamp_token" cbor:"7,keyasint"`
	// LedgerEntryHash is the hash of the ledger entry
	LedgerEntryHash []byte `json:"ledger_entry_hash" cbor:"8,keyasint"`
	// MerkleInclusionProof is the proof of inclusion in the ledger
	MerkleInclusionProof *InclusionProof `json:"merkle_inclusion_proof" cbor:"9,keyasint"`
	// BundleVersion is the version of the bundle format
	BundleVersion string `json:"bundle_version" cbor:"10,keyasint"`
}

// InclusionProof represents a Merkle inclusion proof
type InclusionProof struct {
	// LeafIndex is the index of the leaf in the tree
	LeafIndex int `json:"leaf_index" cbor:"1,keyasint"`
	// LeafHash is the hash of the leaf
	LeafHash []byte `json:"leaf_hash" cbor:"2,keyasint"`
	// TreeSize is the size of the tree at proof time
	TreeSize int `json:"tree_size" cbor:"3,keyasint"`
	// Path is the hash path from leaf to root
	Path [][]byte `json:"path" cbor:"4,keyasint"`
}

// Metadata contains non-cryptographic metadata about the signature
type Metadata struct {
	// CreatedAt is when the bundle was created
	CreatedAt time.Time `json:"created_at"`
	// SignerOffice is the office of the signer
	SignerOffice string `json:"signer_office"`
	// Jurisdiction is the jurisdiction
	Jurisdiction string `json:"jurisdiction"`
	// ContentType describes the type of content signed
	ContentType string `json:"content_type,omitempty"`
	// ContentDescription is a human-readable description
	ContentDescription string `json:"content_description,omitempty"`
}

// VerificationResult represents the result of bundle verification
type VerificationResult struct {
	// Valid indicates if the bundle is valid
	Valid bool `json:"valid"`
	// Timestamp is when verification occurred
	Timestamp time.Time `json:"timestamp"`
	// Checks contains detailed check results
	Checks map[string]bool `json:"checks"`
	// Errors contains any errors encountered
	Errors []string `json:"errors,omitempty"`
	// Warnings contains any warnings
	Warnings []string `json:"warnings,omitempty"`
	// IdentityInfo contains info about the signer identity
	IdentityInfo *IdentityInfo `json:"identity_info,omitempty"`
}

// IdentityInfo contains information about the signer's identity
type IdentityInfo struct {
	// IdentityID is the identity identifier
	IdentityID string `json:"identity_id"`
	// Office is the office/role
	Office string `json:"office"`
	// Jurisdiction is the jurisdiction
	Jurisdiction string `json:"jurisdiction"`
	// KeyVersion is the key version used
	KeyVersion int `json:"key_version"`
	// ValidFrom is the start of validity
	ValidFrom time.Time `json:"valid_from"`
	// ValidTo is the end of validity
	ValidTo time.Time `json:"valid_to"`
	// Status is the current status
	Status string `json:"status"`
}

// DeviceAttestation represents optional device attestation layer
type DeviceAttestation struct {
	// DeviceCertChain is the device certificate chain
	DeviceCertChain [][]byte `json:"device_cert_chain" cbor:"1,keyasint"`
	// FirmwareHash is the hash of the firmware
	FirmwareHash []byte `json:"firmware_hash" cbor:"2,keyasint"`
	// CaptureTimestamp is when the capture occurred
	CaptureTimestamp time.Time `json:"capture_timestamp" cbor:"3,keyasint"`
	// SensorSignature is the signature from the sensor
	SensorSignature []byte `json:"sensor_signature" cbor:"4,keyasint"`
}
