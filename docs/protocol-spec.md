# Civic Attest Protocol Specification

**Version:** 1.0
**Status:** Draft
**Last Updated:** 2026-02-22

## Abstract

Civic Attest is a deterministic, cryptographically verifiable, publicly auditable authenticity framework designed for high-trust public digital communications. This document defines the complete protocol specification.

## 1. System Guarantees

### 1.1 Explicit Guarantees

The system provides the following cryptographic guarantees:

1. **Origin Authenticity** - Cryptographic non-repudiation through digital signatures
2. **Byte-level Integrity** - Detection of any modification to signed content
3. **Verifiable Timestamp Anchoring** - RFC 3161 compliant timestamps
4. **Append-only Public Transparency** - Immutable public ledger
5. **Key Lifecycle Governance** - Formal key management with trustee oversight
6. **Deterministic Reproducibility** - Canonical encoding ensures identical results
7. **Offline Verification Capability** - Verification possible without network access

### 1.2 Explicit Non-Guarantees

The system explicitly does NOT guarantee:

1. **Truthfulness of Content** - Digital signatures prove authorship, not veracity
2. **Intent Validation** - Cannot determine if signer acted voluntarily
3. **Coercion Detection** - Cannot detect if signer was under duress
4. **Political Neutrality** - System is neutral but usage may not be
5. **Semantic Authenticity** - Does not validate meaning or context

## 2. Cryptographic Primitives

### 2.1 Signature Algorithms

- **Primary:** Ed25519
- **Secondary:** Ed448 (reserved)
- **Quantum Migration:** Dilithium (feature flag, reserved)

### 2.2 Hash Functions

| Use Case | Algorithm | Output Size |
|----------|-----------|-------------|
| Default | SHA-256 | 256 bits |
| Long-term | SHA-3-512 | 512 bits |
| High-throughput | BLAKE3 | 256 bits |

### 2.3 Canonical Encoding

All cryptographic outputs use:
- **Primary:** Canonical CBOR (RFC 8949)
- **Alternative:** Deterministic JSON (RFC 8785)

**Invariant:** No base64 encoding without canonical ordering.

### 2.4 Timestamp Authority

- **Standard:** RFC 3161 compliant
- **Protocol:** HTTP/HTTPS TSA protocol
- **Format:** ASN.1 DER encoded TimeStampToken

## 3. Identity Model

### 3.1 Identity Object

```json
{
  "office_id": "string",
  "jurisdiction": "string",
  "public_key": "bytes",
  "key_version": "integer",
  "valid_from": "timestamp",
  "valid_to": "timestamp",
  "key_algorithm": "string",
  "status": "enum",
  "revocation_pointer": "string?",
  "identity_id": "string"
}
```

### 3.2 Key Generation

1. Generated inside HSM
2. Export disabled
3. Public key extracted
4. Certificate self-signed by trustee quorum

### 3.3 Key Ceremony

Required participants: 3-of-5 trustee quorum

**Procedure:**
1. Trustees assemble
2. Key generated in HSM
3. Ceremony recorded (audio/video)
4. Public hash broadcast
5. Ledger entry appended

### 3.4 Rotation Policy

- **Scheduled:** Annual rotation
- **Emergency:** Immediate rotation on compromise
- **Cross-signing:** Successor key signed by predecessor

### 3.5 Revocation Triggers

- Compromise detected
- Trustee quorum vote
- Expiration reached
- Office transition

## 4. Signature Bundle Format

### 4.1 Bundle Structure

```cbor
{
  1: content_hash,
  2: content_hash_algorithm,
  3: canonical_format_version,
  4: signer_identity_id,
  5: key_version,
  6: signature,
  7: timestamp_token,
  8: ledger_entry_hash,
  9: merkle_inclusion_proof,
  10: bundle_version
}
```

### 4.2 Invariants

1. `content_hash` computed on canonical byte stream only
2. `signature` must reference exact hash
3. `ledger_entry_hash` must match append record
4. `inclusion_proof` must verify to ledger root

## 5. Ledger Architecture

### 5.1 Structure

Type: Append-only Binary Merkle Tree

**Leaf Node:**
```
{
  entry_hash,
  timestamp,
  signer_identity_id,
  signature_hash
}
```

**Internal Node:**
```
{
  left_hash,
  right_hash,
  parent_hash = H(left_hash || right_hash)
}
```

**Root:**
```
{
  tree_size,
  root_hash,
  signed_tree_head
}
```

### 5.2 Signed Tree Head

- Signed by ledger authority key
- Published publicly
- Mirrored by third parties
- Enables fork detection via gossip

### 5.3 Proofs

**Inclusion Proof:** Binary hash path from leaf to root

**Consistency Proof:** Demonstrates append-only property between tree states

## 6. Signing Flow

### 6.1 Deterministic Signing Procedure

```
1. Capture master artifact
2. Canonicalize (normalize metadata, freeze codec)
3. Hash canonical artifact
4. Send hash to signer
5. Signer requests HSM sign
6. Signer requests timestamp authority
7. Append entry to ledger
8. Generate inclusion proof
9. Emit signature bundle
```

### 6.2 Atomicity

If ledger append fails, signature is invalidated.

## 7. Verification Flow

### 7.1 Standard Verification

**Input:** `media_file + signature_bundle`

**Procedure:**
```
1. Canonicalize media
2. Compute hash
3. Compare to bundle.content_hash
4. Verify signature using identity public key
5. Validate key not revoked
6. Validate timestamp token
7. Verify inclusion proof
8. Verify ledger root consistency
9. Return verification result
```

### 7.2 Offline Mode

Skip ledger live validation, use cached signed tree head.

### 7.3 Audit Mode

Full ledger replay validation.

## 8. Security Invariants

1. Private keys never leave hardware boundary
2. All historical ledger entries immutable
3. Canonicalization deterministic
4. No mutable state for signature artifacts
5. Revocations are irreversible
6. Every signature traceable to identity version

## 9. Performance Constraints

| Metric | Target |
|--------|--------|
| Signature throughput | 1,000/min per HSM |
| Ledger append latency | < 100 ms |
| Verification latency | < 50 ms (local) |
| Ledger growth | O(n) |
| Inclusion proof size | O(log n) |

## 10. Threat Model

### 10.1 Security Domain

High-trust public governance

### 10.2 Adversary Model

Nation-state capable adversaries

### 10.3 Operational Mode

Hybrid offline + online

## 11. Governance

### 11.1 Trustee Structure

- 5 trustees total
- 3-of-5 quorum for operations

### 11.2 Authority Separation

- Identity issuance authority
- Ledger authority
- Signer authority

All separate entities.

### 11.3 Emergency Override

- Supermajority required
- Publicly logged

## 12. Compliance

### 12.1 Export Controls

- ECC algorithms permitted
- Quantum algorithms optional (feature flag)

### 12.2 License

Apache 2.0 or AGPL (deployment dependent)

---

**References:**
- RFC 3161: Time-Stamp Protocol (TSP)
- RFC 8032: Edwards-Curve Digital Signature Algorithm (EdDSA)
- RFC 8949: Concise Binary Object Representation (CBOR)
