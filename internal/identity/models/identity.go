package models

import (
	"time"
)

// IdentityStatus represents the status of an identity
type IdentityStatus string

const (
	// StatusActive indicates the identity is active and valid
	StatusActive IdentityStatus = "active"
	// StatusRevoked indicates the identity has been revoked
	StatusRevoked IdentityStatus = "revoked"
	// StatusExpired indicates the identity has expired
	StatusExpired IdentityStatus = "expired"
	// StatusPending indicates the identity is pending activation
	StatusPending IdentityStatus = "pending"
)

// Identity represents a cryptographic identity in the system
type Identity struct {
	// OfficeID is the unique identifier for the office/role
	OfficeID string `json:"office_id" cbor:"1,keyasint"`
	// Jurisdiction is the geographic or organizational jurisdiction
	Jurisdiction string `json:"jurisdiction" cbor:"2,keyasint"`
	// PublicKey is the public key in raw bytes
	PublicKey []byte `json:"public_key" cbor:"3,keyasint"`
	// KeyVersion is the version of this key
	KeyVersion int `json:"key_version" cbor:"4,keyasint"`
	// ValidFrom is the start of the validity period
	ValidFrom time.Time `json:"valid_from" cbor:"5,keyasint"`
	// ValidTo is the end of the validity period
	ValidTo time.Time `json:"valid_to" cbor:"6,keyasint"`
	// KeyAlgorithm is the signature algorithm used
	KeyAlgorithm string `json:"key_algorithm" cbor:"7,keyasint"`
	// Status is the current status of the identity
	Status IdentityStatus `json:"status" cbor:"8,keyasint"`
	// RevocationPointer points to the revocation record if revoked
	RevocationPointer string `json:"revocation_pointer,omitempty" cbor:"9,keyasint,omitempty"`
	// IdentityID is a unique identifier for this identity version
	IdentityID string `json:"identity_id" cbor:"10,keyasint"`
}

// IsValid checks if the identity is currently valid
func (i *Identity) IsValid(at time.Time) bool {
	if i.Status != StatusActive {
		return false
	}

	if at.Before(i.ValidFrom) || at.After(i.ValidTo) {
		return false
	}

	return true
}

// KeyCeremonyRecord records the key generation ceremony
type KeyCeremonyRecord struct {
	// CeremonyID is the unique identifier for this ceremony
	CeremonyID string `json:"ceremony_id"`
	// Timestamp is when the ceremony occurred
	Timestamp time.Time `json:"timestamp"`
	// Trustees are the trustees who participated
	Trustees []string `json:"trustees"`
	// QuorumSize is the quorum size (e.g., 3 of 5)
	QuorumSize int `json:"quorum_size"`
	// TotalTrustees is the total number of trustees
	TotalTrustees int `json:"total_trustees"`
	// RecordingHash is the hash of the ceremony recording
	RecordingHash []byte `json:"recording_hash"`
	// PublicKeyHash is the hash of the generated public key
	PublicKeyHash []byte `json:"public_key_hash"`
	// LedgerEntryHash is the hash of the ledger entry
	LedgerEntryHash []byte `json:"ledger_entry_hash"`
}

// RotationRecord records a key rotation event
type RotationRecord struct {
	// RotationID is the unique identifier
	RotationID string `json:"rotation_id"`
	// OldIdentityID is the identity being rotated from
	OldIdentityID string `json:"old_identity_id"`
	// NewIdentityID is the new identity
	NewIdentityID string `json:"new_identity_id"`
	// Timestamp is when the rotation occurred
	Timestamp time.Time `json:"timestamp"`
	// Reason is the reason for rotation
	Reason string `json:"reason"`
	// CrossSignature is the signature by old key over new key
	CrossSignature []byte `json:"cross_signature"`
	// Emergency indicates if this was an emergency rotation
	Emergency bool `json:"emergency"`
}

// RevocationRecord records a key revocation
type RevocationRecord struct {
	// RevocationID is the unique identifier
	RevocationID string `json:"revocation_id"`
	// IdentityID is the identity being revoked
	IdentityID string `json:"identity_id"`
	// Timestamp is when the revocation occurred
	Timestamp time.Time `json:"timestamp"`
	// Reason is the reason for revocation
	Reason string `json:"reason"`
	// TrusteeSignatures are signatures from the quorum
	TrusteeSignatures [][]byte `json:"trustee_signatures"`
	// LedgerEntryHash is the hash of the ledger entry
	LedgerEntryHash []byte `json:"ledger_entry_hash"`
	// Irreversible marks this revocation as permanent
	Irreversible bool `json:"irreversible"`
}
