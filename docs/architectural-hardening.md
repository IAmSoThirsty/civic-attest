# Architectural Hardening Specification

**Version:** 2.0
**Status:** Draft - Critical Security Enhancements
**Last Updated:** 2026-02-23
**Classification:** PUBLIC

## Executive Summary

This document specifies critical architectural hardening enhancements to Civic Attest in response to adversarial security review. These enhancements are designed to defend against nation-state adversaries and ensure the system can serve as trusted public cryptographic infrastructure for democratic institutions.

**Threat Model:** Nation-state capable adversaries
**Security Domain:** Public cryptographic institution
**Operational Classification:** Federated multi-operator transparency log

## 1. Transparency Log Hardening

### 1.1 Witness Cosigning System

**Current State:** Single-authority signed tree heads (vulnerable to equivocation)

**Required Enhancement:** Multi-witness cosigning architecture

#### 1.1.1 Witness Set Specification

**Minimum Witnesses:** 3 independent organizations
**Recommended:** 5+ geographically distributed witnesses

**Witness Requirements:**
- Independent legal entity
- Separate infrastructure
- Different geographic jurisdiction
- Independent operational control
- Public accountability commitment

**Witness Responsibilities:**
1. Monitor ledger append operations
2. Verify tree head integrity
3. Sign valid tree heads independently
4. Refuse to sign on inconsistency detection
5. Publish signed tree heads publicly
6. Participate in gossip protocol

#### 1.1.2 Signed Tree Head Format (Enhanced)

```json
{
  "tree_size": 1000,
  "root_hash": "a3f2b1...",
  "timestamp": "2026-02-23T12:00:00Z",
  "ledger_authority_signature": "7e4c9d...",
  "ledger_authority_id": "ledger-authority-v1",
  "witness_signatures": [
    {
      "witness_id": "witness-org-1",
      "witness_pubkey_hash": "f3a2b1...",
      "signature": "8d3e2f...",
      "signed_at": "2026-02-23T12:00:01Z"
    },
    {
      "witness_id": "witness-org-2",
      "witness_pubkey_hash": "a1b2c3...",
      "signature": "9e4f3a...",
      "signed_at": "2026-02-23T12:00:02Z"
    },
    {
      "witness_id": "witness-org-3",
      "witness_pubkey_hash": "d4e5f6...",
      "signature": "1a2b3c...",
      "signed_at": "2026-02-23T12:00:03Z"
    }
  ],
  "witness_quorum": "3-of-5",
  "sth_version": 2
}
```

#### 1.1.3 Witness Cosigning Protocol

**Step 1: Tree Head Generation**
1. Ledger authority computes new tree head
2. Signs tree head with ledger authority key
3. Publishes to witness network

**Step 2: Witness Verification**
1. Each witness independently:
   - Verifies tree head cryptographic integrity
   - Verifies consistency with previous tree head
   - Verifies append-only property
   - Checks for fork attempts
2. If valid, witness signs tree head
3. Publishes witness signature

**Step 3: Quorum Achievement**
1. Collect witness signatures
2. Verify quorum threshold (≥3 of 5)
3. Bundle signatures into complete STH
4. Distribute to verifiers

**Step 4: Fork Detection**
1. Compare tree heads across witnesses
2. Detect conflicting signatures
3. Alert on split-view detection
4. Initiate emergency freeze if divergence detected

#### 1.1.4 Gossip Verification Protocol

**Client Verifiers MUST:**
1. Request tree heads from multiple witnesses
2. Compare root hashes for same tree size
3. Verify witness signatures
4. Detect and report inconsistencies
5. Refuse to accept signatures without witness quorum

**Gossip Network:**
- Decentralized peer-to-peer monitoring
- Automated split-view detection
- Public alert mechanisms
- Community-driven verification

### 1.2 Multi-Operator Federation Model

**Architecture:** Federated multi-operator transparency log

**Federation Structure:**

**Primary Operator:** Core ledger authority
- Maintains canonical ledger
- Coordinates append operations
- Signs tree heads first

**Mirror Operators (Minimum 5):**
- Independent infrastructure
- Real-time replication
- Independent tree head signing
- Geographic distribution
- Serve as witnesses

**Federation Requirements:**
1. No single operator controls majority
2. Minimum 3 operators in different legal jurisdictions
3. Public disclosure of all operators
4. Transparent governance of federation
5. Documented removal/addition procedures

### 1.3 Split-View Detection Guarantees

**Detection Mechanisms:**

**1. Client-side Detection:**
- Verifiers query multiple witnesses
- Compare tree heads
- Alert on mismatch

**2. Witness Cross-verification:**
- Witnesses monitor each other
- Automated divergence detection
- Public alert publication

**3. Community Monitoring:**
- Public monitoring tools
- Open-source verification clients
- Bounty program for fork detection

**4. Cryptographic Proofs:**
- Non-equivocation guarantees via witness cosigning
- Append-only consistency proofs required
- Public audit trail of all tree heads

## 2. Identity Layer Hardening

### 2.1 Identity State Merkleization

**Current Gap:** Identity verification may require online lookup

**Required Enhancement:** Complete identity state committed to ledger root

#### 2.1.1 Identity Merkle Tree

**Structure:**
```
Identity Tree Root
├── Active Identities Subtree
│   ├── Identity 1 (hash of identity record)
│   ├── Identity 2
│   └── ...
├── Revoked Identities Subtree
│   ├── Revoked 1 (hash of revocation record)
│   └── ...
└── Metadata Subtree
    ├── Tree version
    └── Last update timestamp
```

**Identity Tree Root:** SHA-256 hash of complete identity state

**Commitment:** Identity tree root embedded in ledger signed tree head

#### 2.1.2 Enhanced Signed Tree Head with Identity Commitment

```json
{
  "tree_size": 1000,
  "root_hash": "a3f2b1...",
  "identity_tree_root": "b4c3d2...",
  "revocation_tree_root": "c5d4e3...",
  "timestamp": "2026-02-23T12:00:00Z",
  "ledger_authority_signature": "7e4c9d...",
  "witness_signatures": [...],
  "commitment_version": 2
}
```

### 2.2 Offline Identity Verification

**Capability:** Verify identity status without network access

**Implementation:**

**1. Identity Inclusion Proof:**
```json
{
  "identity_id": "mayor-springfield-v1",
  "identity_record": {...},
  "merkle_proof": ["hash1", "hash2", ...],
  "tree_root": "b4c3d2...",
  "proof_type": "inclusion"
}
```

**2. Revocation Exclusion Proof:**
```json
{
  "identity_id": "mayor-springfield-v1",
  "merkle_proof": ["hash1", "hash2", ...],
  "revocation_tree_root": "c5d4e3...",
  "proof_type": "non-revocation"
}
```

**3. Bundled Verification:**
- Signature bundles include identity inclusion proof
- Include non-revocation proof
- Verifier validates against committed roots
- No online lookup required

### 2.3 Revocation Distribution Models

**Multiple Distribution Channels:**

#### 2.3.1 CRL Snapshot Embedding

**Certificate Revocation List snapshots:**
- Published as ledger entries
- Merkleized and committed to tree root
- Included in signature bundles for offline verification
- Updated hourly (configurable)

#### 2.3.2 OCSP-like Real-time Query

**Online Status Protocol (OSP):**
- Real-time revocation queries
- Signed responses from identity authority
- Response includes STH reference
- Cached responses valid for 1 hour
- Fallback to embedded CRL if unavailable

#### 2.3.3 Bloom Filter Distribution

**Privacy-preserving revocation checking:**
- Compact Bloom filter of revoked identities
- Published daily
- False positive rate: 0.1%
- Download size: ~100KB for 10,000 identities
- Offline verification friendly

## 3. Time Anchoring Hardening

### 3.1 Multi-TSA Architecture

**Current:** RFC 3161 TSA (single point of trust)

**Enhancement:** Multiple independent timestamp authorities

**Minimum TSAs:** 3 independent providers
**Timestamp Bundle:** Includes timestamps from all TSAs

```json
{
  "content_hash": "a3f2b1...",
  "timestamps": [
    {
      "tsa_id": "tsa-provider-1",
      "timestamp_token": "...",
      "signed_time": "2026-02-23T12:00:00Z",
      "tsa_signature": "..."
    },
    {
      "tsa_id": "tsa-provider-2",
      "timestamp_token": "...",
      "signed_time": "2026-02-23T12:00:01Z",
      "tsa_signature": "..."
    }
  ],
  "timestamp_quorum": "2-of-3"
}
```

### 3.2 Blockchain Time Anchoring

**Public Blockchain Anchoring:** Tree head hash published to multiple public blockchains

**Anchoring Targets:**

**1. Bitcoin OP_RETURN:**
```
OP_RETURN <civic-attest-v2> <tree_head_hash>
```
- Published every 1000 entries or 1 hour
- Immutable timestamp proof
- Publicly verifiable

**2. Ethereum Calldata:**
```solidity
contract CivicAttestAnchor {
    event TreeHeadAnchored(
        uint256 indexed treeSize,
        bytes32 treeRoot,
        bytes32 identityRoot,
        uint256 timestamp
    );

    function anchorTreeHead(
        uint256 treeSize,
        bytes32 treeRoot,
        bytes32 identityRoot
    ) external onlyAuthorized {
        emit TreeHeadAnchored(treeSize, treeRoot, identityRoot, block.timestamp);
    }
}
```

**3. Additional Blockchains:**
- Polygon (low-cost anchoring)
- Avalanche (fast finality)
- Solana (high throughput)

**4. Public Newspaper Hash Publication:**
- Weekly hash publication in major newspapers
- QR code with tree head hash
- Physical archive for long-term verification
- Public verifiability without technology

**Anchoring Frequency:**
- Bitcoin: Every 1000 entries or 1 hour
- Ethereum: Every 500 entries or 30 minutes
- Other chains: Every 100 entries or 10 minutes
- Newspapers: Weekly

**Benefits:**
1. Decentralized timestamp proof
2. Survives individual TSA compromise
3. Publicly verifiable by anyone
4. Long-term archival properties
5. Multi-jurisdiction redundancy

## 4. HSM Operational Hardening

### 4.1 Enhanced HSM Controls

**Current:** FIPS 140-2 Level 3+ HSM with export disabled

**Enhancements Required:**

#### 4.1.1 Rate Limiting Inside HSM

**Signature Rate Limits:**
- Maximum: 1,000 signatures per minute per HSM
- Burst limit: 100 signatures in 10 seconds
- Enforced at hardware level
- Alerts on approaching limits

**Implementation:**
```
HSM Configuration:
  rate_limit:
    max_per_minute: 1000
    burst_window: 10s
    burst_max: 100
    action_on_exceed: REJECT
    alert_threshold: 80%
```

#### 4.1.2 Anti-Rollback Firmware Enforcement

**Requirements:**
- Firmware version monotonically increasing
- Rollback attempts detected and rejected
- Firmware updates signed by vendor + trustee quorum
- Version attestation in every signature operation

**Firmware Update Procedure:**
1. Vendor releases new firmware
2. Vendor signs firmware
3. Trustee quorum reviews firmware
4. Trustee quorum co-signs firmware
5. HSM verifies dual signature
6. HSM updates and locks version counter
7. Rollback impossible

#### 4.1.3 Dual-Control Activation

**M-of-N Control:**
- Minimum 2-of-3 operator authentication required
- Physical presence required (proximity tokens)
- Biometric + PIN + hardware token
- No single operator can activate HSM

**Activation Procedure:**
1. Operator 1 presents credentials
2. HSM logs authentication
3. Operator 2 presents credentials (within 5 minutes)
4. HSM verifies both authentications
5. HSM unlocks for operation
6. Session timeout: 8 hours
7. Re-authentication required

#### 4.1.4 Sealed Audit Logs

**Comprehensive Logging:**
- Every signature operation
- Every authentication attempt
- Every configuration change
- Every error condition
- Cryptographically sealed

**Export Schedule:**
- Real-time stream to separate log server
- Hourly export to append-only storage
- Daily export to offline archive
- Weekly backup to geographic distant location

**Log Format:**
```json
{
  "timestamp": "2026-02-23T12:00:00.123Z",
  "hsm_id": "hsm-primary-1",
  "event_type": "SIGNATURE_OPERATION",
  "operator_id": "operator-1",
  "identity_id": "mayor-springfield-v1",
  "operation_hash": "a3f2b1...",
  "sequence_number": 12345,
  "log_signature": "7e4c9d..."
}
```

**Log Integrity:**
- Each log entry signed by HSM
- Chained hash linking
- Tampering detectable
- Exported logs cannot be modified

#### 4.1.5 Threshold Key Signing (Optional)

**Advanced Deployment Option:**
- Split signing key across 2-of-3 HSMs
- Threshold signature scheme (TSS)
- No single HSM can sign alone
- Byzantine fault tolerance

**Benefits:**
- Eliminates single HSM compromise risk
- Geographic distribution possible
- Enhanced resilience

**Implementation:**
- ECDSA threshold signatures (Ed25519 compatible)
- Secure multi-party computation (MPC)
- Each HSM holds key share
- Cooperative signing protocol

## 5. Governance Hardening

### 5.1 Publicly Auditable Vote Publication

**Current:** Governance decisions documented

**Enhancement:** All trustee votes published to transparency ledger

**Vote Record Format:**
```json
{
  "proposal_id": "PROP-2026-001",
  "proposal_type": "KEY_ROTATION",
  "proposal_hash": "a3f2b1...",
  "votes": [
    {
      "trustee_id": "trustee-1",
      "vote": "APPROVE",
      "signature": "7e4c9d...",
      "voted_at": "2026-02-23T12:00:00Z"
    },
    {
      "trustee_id": "trustee-2",
      "vote": "APPROVE",
      "signature": "8d3e2f...",
      "voted_at": "2026-02-23T12:00:05Z"
    }
  ],
  "quorum_met": true,
  "decision": "APPROVED",
  "executed_at": "2026-02-23T18:00:00Z"
}
```

**Publication:**
- Every vote appended to transparency ledger
- Cryptographically signed by each trustee
- Publicly verifiable
- Cannot be altered retroactively

### 5.2 Delayed Execution Window

**Critical Operations Delay:**

**Tier 1: Critical Changes (72-hour delay)**
- Key rotation (non-emergency)
- Governance amendments
- Trustee changes
- Major policy changes

**Tier 2: Moderate Changes (24-hour delay)**
- Configuration changes
- Operational parameter adjustments
- Non-critical updates

**Tier 3: Emergency Operations (No delay)**
- Key compromise response
- Active attack mitigation
- System integrity threats

**Delayed Execution Process:**
1. Proposal submitted and voted
2. If approved, enters delay period
3. Public announcement with countdown
4. Community review period
5. Execution only after delay expires
6. Emergency override requires supermajority (4-of-5)

**Benefits:**
- Prevents hasty decisions under coercion
- Allows community oversight
- Enables intervention if needed
- Transparent governance timeline

### 5.3 Emergency Freeze Mechanism

**Freeze Authority:** Any 2 trustees can initiate freeze

**Freeze Triggers:**
- Suspected key compromise
- Ledger inconsistency detected
- Coordinated attack detected
- Trustee coercion suspected

**Freeze Effects:**
1. All signing operations halted
2. Ledger writes paused
3. Emergency quorum convened (4 hours)
4. Investigation initiated
5. Resume requires 3-of-5 vote

**Freeze Duration:**
- Initial: 24 hours
- Extended: 3-of-5 vote required
- Maximum: 7 days without full quorum review

### 5.4 Trustee Rotation Process

**Regular Rotation:**
- Staggered terms (every 12 months)
- Overlap period (3 months)
- Knowledge transfer protocol
- Public attestation

**Emergency Rotation:**
- Triggered by: compromise, unavailability, coercion
- Requires 4-of-5 vote
- Expedited onboarding
- Public disclosure

**Rotation Ceremony:**
1. Outgoing trustee public statement
2. Incoming trustee public commitment
3. Key handoff ceremony (recorded)
4. Updated public trustee registry
5. Ledger entry with attestations
6. Community announcement

### 5.5 Governance Transparency Ledger

**All Governance Actions Logged:**
- Trustee appointments/removals
- Policy changes
- Emergency actions
- Vote records
- Freeze events
- Rotation ceremonies

**Separate Governance Ledger Stream:**
- Dedicated ledger entries for governance
- Cryptographically linked to main ledger
- Public API for governance queries
- Real-time monitoring dashboard

## 6. Deterministic Encoding Hardening

### 6.1 Formal Canonicalization Specification

**Current Gap:** Canonical CBOR/JSON mentioned but not formally specified

**Required:** Complete formal specification

#### 6.1.1 Canonical CBOR Rules

**RFC 8949 Section 4.2 Compliance (Deterministic Encoding):**

**1. Integer Encoding:**
- Use shortest possible encoding
- Prefer integer types over bignum
- No leading zeros

**2. Floating Point Normalization:**
- PROHIBITED in cryptographic contexts
- If required: IEEE 754 canonical NaN
- Prefer rational representation
- -0.0 normalized to +0.0

**3. Map Key Sorting:**
- Sort by byte-wise lexicographic order
- After canonical encoding of keys
- Deterministic iteration order

**4. No Indefinite Length:**
- All arrays, maps, strings must use definite length
- Streaming not permitted for cryptographic content

**5. Disallowed Types:**
- No undefined values
- No tags unless explicitly specified
- No duplicate map keys
- No NaN, Infinity in cryptographic contexts

**6. String Encoding:**
- UTF-8 normalization: NFC (Canonical Decomposition followed by Canonical Composition)
- No unnormalized Unicode
- No overlong UTF-8 sequences

**7. Whitespace:**
- N/A for CBOR (binary format)

#### 6.1.2 Deterministic JSON Rules

**RFC 8785 (JSON Canonicalization Scheme - JCS) Compliance:**

**1. Unicode Normalization:**
- NFC (Canonical Decomposition + Canonical Composition)
- Apply before encoding

**2. Whitespace:**
- Remove all unnecessary whitespace
- No spaces around ':' or ','
- No newlines or indentation

**3. String Escaping:**
- Minimal escaping
- Only control characters U+0000 through U+001F
- Forward slash not escaped
- Use lowercase hex in escapes (\u00XX)

**4. Number Representation:**
- No leading zeros (except "0")
- No trailing zeros after decimal point
- No leading '+' sign
- Use lowercase 'e' for exponent
- No floating point: use string representation for precision

**5. Map Key Sorting:**
- Sort by UTF-16 code unit order
- After escaping

**6. Boolean/Null:**
- Lowercase: true, false, null

**7. Disallowed:**
- No NaN, Infinity (use string representation)
- No comments
- No duplicate keys

#### 6.1.3 Canonical Encoding Verification

**Test Vectors Required:**
- Minimum 100 test vectors per encoding
- Edge cases documented
- Cross-implementation validation
- Fuzzing for non-determinism

**Reference Implementations:**
- Canonical CBOR: cbor-deterministic (reference)
- Canonical JSON: jcs (reference)
- Cross-language implementations required

**Validation:**
```
For any input A:
  canonical(A) == canonical(canonical(A))
  hash(canonical(A)) == hash(canonical(canonical(A)))
```

### 6.2 Encoding Version Field

**All Cryptographic Structures Include:**
```json
{
  "canonical_encoding_version": "2.0",
  "encoding_type": "CBOR_DETERMINISTIC",
  "unicode_normalization": "NFC",
  ...
}
```

**Version Evolution:**
- Backwards compatibility required
- Verifiers support all versions
- Signers use latest version
- Migration path documented

## 7. Quantum Migration Strategy

### 7.1 Dual-Signature Architecture

**Current:** Optional Dilithium flag

**Enhancement:** Production-ready dual-signature system

#### 7.1.1 Hybrid Signature Bundle

```json
{
  "content_hash": "a3f2b1...",
  "signatures": {
    "classical": {
      "algorithm": "Ed25519",
      "signature": "7e4c9d...",
      "pubkey": "a3f2b1...",
      "key_version": 1
    },
    "post_quantum": {
      "algorithm": "Dilithium3",
      "signature": "9f5e3a...",
      "pubkey": "b4c3d2...",
      "key_version": 1,
      "optional": true
    }
  },
  "signature_policy": "REQUIRE_CLASSICAL_AND_OPTIONAL_PQ"
}
```

**Signature Policies:**

**Phase 1: Classical Only (Current)**
- Require: Ed25519
- Optional: None

**Phase 2: Hybrid Optional (2026-2027)**
- Require: Ed25519
- Optional: Dilithium3
- Signers encouraged to dual-sign

**Phase 3: Hybrid Required (2028+)**
- Require: Ed25519 AND Dilithium3
- Both signatures must validate

**Phase 4: Post-Quantum Only (When Quantum Threat Realized)**
- Require: Dilithium3
- Optional: Ed25519 (for historical compatibility)

#### 7.1.2 Ledger Algorithm Agility

**Ledger Entry Format Enhancement:**
```json
{
  "entry_hash": "a3f2b1...",
  "signature_algorithm_version": 2,
  "supported_algorithms": ["Ed25519", "Dilithium3"],
  "primary_algorithm": "Ed25519",
  "timestamp": "2026-02-23T12:00:00Z",
  ...
}
```

**Algorithm Version Registry:**
```json
{
  "version": 1,
  "algorithms": {
    "signature": "Ed25519",
    "hash": "SHA-256"
  },
  "valid_from": "2024-01-01",
  "valid_to": "2030-01-01"
},
{
  "version": 2,
  "algorithms": {
    "signature": "Ed25519+Dilithium3",
    "hash": "SHA-256"
  },
  "valid_from": "2026-01-01",
  "valid_to": null
}
```

#### 7.1.3 Migration Timeline

**2026 Q1:** Dilithium implementation complete
**2026 Q2:** Dual-key generation support
**2026 Q3:** Optional PQ signatures in production
**2027 Q1:** Dual-signature encouraged
**2027 Q4:** Dual-signature required
**2028+:** Monitor quantum threat landscape
**Quantum Threat Detected:** Emergency transition to PQ-only

#### 7.1.4 Historic Validation Preservation

**Critical Requirement:** Migration must not break historic signature validation

**Solution:**
1. Verifiers support all algorithm versions
2. Verification includes algorithm version check
3. Historic signatures remain valid
4. Algorithm version committed to ledger
5. Time-based algorithm policies

**Verification Logic:**
```
function verify(signature, content, timestamp):
  algorithm_version = get_algorithm_version(timestamp)

  if algorithm_version == 1:
    return verify_ed25519(signature.classical, content)

  if algorithm_version == 2:
    ed25519_valid = verify_ed25519(signature.classical, content)
    if signature.post_quantum exists:
      pq_valid = verify_dilithium(signature.post_quantum, content)
      return ed25519_valid AND pq_valid
    return ed25519_valid

  if algorithm_version == 3:
    return verify_dilithium(signature.post_quantum, content)
```

## 8. Denial of Service Protection

### 8.1 Ledger Append Rate Limiting

**Protection Mechanisms:**

#### 8.1.1 Per-Identity Rate Limits

```
Rate Limit Tiers:
  Standard Identity: 100 signatures/hour
  Verified Organization: 1,000 signatures/hour
  High-Volume Entity: 10,000 signatures/hour
  Emergency Override: Trustee approval required
```

**Implementation:**
- Token bucket algorithm
- Per-identity quota
- Refill rate: hourly
- Burst allowance: 2x quota for 5 minutes

#### 8.1.2 Proof-of-Work Gating

**For Anonymous/Unregistered Appends:**
- Hashcash-style PoW required
- Difficulty adjusted based on ledger load
- Target: 1 second computation on commodity hardware
- Prevents automated spam

**PoW Verification:**
```json
{
  "ledger_entry": {...},
  "proof_of_work": {
    "nonce": 12345678,
    "difficulty": 20,
    "hash": "000000a3f2b1..."
  }
}
```

#### 8.1.3 Economic Deposit Mechanism

**Optional Stake-Based Admission:**
- Small deposit required for append (e.g., $0.01 USD equivalent)
- Refunded after 24 hours if entry valid
- Forfeited if entry spam/malicious
- Deposit accumulation funds operations

**Benefits:**
- Economic disincentive for spam
- Self-funding mechanism
- Sybil attack mitigation

#### 8.1.4 Resource Bounding

**Per-Entry Limits:**
- Maximum entry size: 10 KB
- Maximum attachment size: 1 MB
- Maximum batch size: 100 entries
- Maximum processing time: 5 seconds

**Global Limits:**
- Maximum append rate: 10,000 entries/minute
- Maximum ledger growth: 1 GB/day
- Backpressure when approaching limits

### 8.2 Network-Level Protection

**DDoS Mitigation:**
- Rate limiting at load balancer
- Geo-distribution with Anycast
- CDN for static content
- Connection limits per IP
- Challenge-response for suspicious traffic

**Application-Level:**
- API rate limiting (100 req/min per IP)
- Request size limits
- Timeout enforcement
- Resource quotas

## 9. Disaster Recovery Enhancements

### 9.1 Ledger Rebuild from Snapshots

**Snapshot Strategy:**
- Hourly: Last 24 hours
- Daily: Last 30 days
- Weekly: Last 12 weeks
- Monthly: Permanent archive

**Snapshot Format:**
```json
{
  "snapshot_version": 2,
  "tree_size": 1000000,
  "root_hash": "a3f2b1...",
  "identity_tree_root": "b4c3d2...",
  "entries": [...],
  "witness_signatures": [...],
  "created_at": "2026-02-23T12:00:00Z",
  "snapshot_signature": "7e4c9d..."
}
```

**Rebuild Procedure:**
1. Identify last valid snapshot
2. Load snapshot into new ledger instance
3. Verify cryptographic integrity
4. Replay entries from snapshot point
5. Cross-verify with mirrors
6. Achieve witness quorum on rebuilt state
7. Resume operations

### 9.2 Cross-Verify Mirror Roots

**Continuous Cross-Verification:**
- Every hour: Compare roots across all mirrors
- Alert on divergence
- Automatic investigation trigger
- Fork detection and resolution

**Mirror Synchronization Protocol:**
1. Primary publishes new tree head
2. Mirrors replicate and verify
3. Mirrors compute independent root
4. Compare roots across network
5. Consensus on canonical state
6. Witness signatures collected

### 9.3 Trustee Emergency Rotation

**Emergency Rotation Triggers:**
- Trustee compromise suspected
- Trustee unavailability (>48 hours)
- Trustee coercion reported
- Security incident

**Rapid Rotation Procedure:**
1. Emergency quorum (2 trustees minimum)
2. Temporary trustee appointment
3. Expedited background check (parallel)
4. Emergency key ceremony (within 24 hours)
5. Limited powers until full verification
6. Full powers after background check complete

### 9.4 Public Incident Disclosure Template

**Required Disclosures:**

```markdown
# Security Incident Disclosure

**Incident ID:** INC-2026-001
**Severity:** [CRITICAL/HIGH/MEDIUM/LOW]
**Disclosure Date:** 2026-02-23
**Incident Date:** 2026-02-20

## Summary
Brief description of incident

## Impact
What was affected

## Timeline
- Detection
- Containment
- Resolution

## Root Cause
Technical details

## Remediation
Actions taken

## User Action Required
What users/verifiers should do

## Contact
security@civic-attest-incident-response.org
```

**Disclosure Timeline:**
- Detection: Immediate internal notification
- Containment: Within 4 hours
- Initial Public Disclosure: Within 24 hours
- Detailed Report: Within 7 days
- Post-Incident Review: Within 30 days

## 10. Federation Model and Operational Classification

### 10.1 Operational Classification

**Answer to Critical Question:**

**Civic Attest is a FEDERATED MULTI-OPERATOR TRANSPARENCY LOG (Option B)**

**Characteristics:**

**Not Single-Operator (Option A):**
- Multiple independent operators
- No single point of trust
- Distributed governance

**Not Publicly Permissionless (Option C):**
- Identity verification required
- Quorum governance
- Controlled admission
- Not blockchain-style permissionless

**Not Simple Consortium (Option D):**
- Public transparency
- Open verification
- Community oversight
- Federation can evolve

### 10.2 Federation Structure

**Tier 1: Core Federation (Minimum 5 Operators)**
- Primary ledger authority
- 4+ mirror authorities
- Each with independent infrastructure
- Geographic distribution required
- Different legal jurisdictions

**Tier 2: Witness Network (Minimum 3, Recommended 7)**
- Independent observers
- Sign tree heads
- Participate in gossip
- Public accountability

**Tier 3: Community Monitors (Unlimited)**
- Anyone can run monitoring node
- Verify consistency
- Report anomalies
- Open-source verification tools

### 10.3 Federation Governance

**Admission Criteria:**
- Legal entity with public accountability
- Technical capability (infrastructure)
- Security certification
- Geographic diversity
- Governance agreement signature
- Public disclosure commitment

**Removal Criteria:**
- Persistent operational failures
- Security violations
- Governance violations
- Vote by 4-of-5 current operators

**Federation Evolution:**
- Add operators via supermajority vote
- Remove via supermajority vote
- Maximum 10 core operators
- Unlimited witnesses
- Public registry of all participants

### 10.4 Survivability Classification

**Byzantine Fault Tolerance:** System survives f failures where n ≥ 3f + 1

**With 5 Operators:**
- Survives 1 Byzantine failure
- Requires 3 honest operators

**With 7 Operators:**
- Survives 2 Byzantine failures
- Requires 5 honest operators

**Survivability Guarantees:**
1. **Single operator compromise:** System continues
2. **Coordinated attack on 2 operators:** System continues (if 7+ operators)
3. **Nation-state targeting single jurisdiction:** System survives (geographic distribution)
4. **Network partition:** System detects and freezes until resolved
5. **Quantum attack:** System can migrate without service interruption

## 11. Implementation Priorities

### 11.1 Critical Path (Months 1-3)

**Priority 1: Witness Cosigning**
- Implement witness protocol
- Deploy 3 initial witnesses
- Update verifier clients

**Priority 2: Identity Merkleization**
- Implement identity tree
- Update ledger format
- Deploy offline verification

**Priority 3: Governance Transparency**
- Implement vote logging
- Delayed execution
- Emergency freeze

### 11.2 High Priority (Months 4-6)

**Priority 4: HSM Hardening**
- Rate limiting
- Dual control
- Audit log export

**Priority 5: Multi-TSA**
- Integrate 3 TSA providers
- Implement quorum logic

**Priority 6: DoS Protection**
- Rate limiting
- PoW gating
- Resource bounds

### 11.3 Important (Months 7-12)

**Priority 7: Blockchain Anchoring**
- Bitcoin OP_RETURN
- Ethereum contract
- Additional chains

**Priority 8: Quantum Migration**
- Dilithium implementation
- Dual-signature support
- Migration tooling

**Priority 9: Formal Canonicalization**
- Complete specification
- Test vectors
- Reference implementations

### 11.4 Ongoing

**Security Audits:** Quarterly
**Penetration Testing:** Semi-annual
**Formal Verification:** Continuous
**Community Engagement:** Continuous

## 12. Conclusion

These architectural hardenings transform Civic Attest from a well-designed system into a genuinely resilient public cryptographic institution capable of:

1. **Surviving nation-state attacks** through federation and witness cosigning
2. **Preventing equivocation** via multi-party verification
3. **Operating offline** through complete state commitment
4. **Migrating to quantum-safe** without breaking historic validation
5. **Resisting capture** through distributed governance
6. **Recovering from disaster** through redundancy and documented procedures
7. **Maintaining public trust** through radical transparency

**This is infrastructure for democratic institutions.**

The difference between implementation quality now determines whether this becomes:
- **Implemented correctly:** Digital notarization infrastructure for democracies
- **Implemented poorly:** A false sense of security

**The stakes demand excellence.**

---

**Appendix A: Threat Model Updates**

See updated threat-model.md

**Appendix B: Implementation Roadmap**

See implementation-roadmap.md (to be created)

**Appendix C: Security Audit Requirements**

See security-audit-requirements.md (to be created)

**Appendix D: Formal Verification Specifications**

See formal-verification.md (to be created)
