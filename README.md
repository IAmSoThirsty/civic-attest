# civic-attest

Deterministic, cryptographically verifiable, publicly auditable authenticity framework for high-trust public digital communications.

## Overview

Civic Attest is a complete cryptographic infrastructure designed to provide origin authenticity, byte-level integrity, and public transparency for digital communications in high-trust environments such as government and public institutions.

### System Guarantees

✓ **Origin Authenticity** - Cryptographic non-repudiation through Ed25519 signatures
✓ **Byte-level Integrity** - SHA-256/SHA-3-512/BLAKE3 hashing ensures content integrity
✓ **Verifiable Timestamp Anchoring** - RFC 3161 compliant timestamps
✓ **Append-only Public Transparency** - Merkle tree-based public ledger
✓ **Key Lifecycle Governance** - Formal trustee-based key ceremonies
✓ **Deterministic Reproducibility** - Canonical CBOR/JSON encoding
✓ **Offline Verification** - Verification without network connectivity

### Explicit Non-Guarantees

✗ **Truthfulness of Content** - Signatures prove authorship, not veracity
✗ **Intent Validation** - Cannot determine voluntary action
✗ **Coercion Detection** - Cannot detect duress
✗ **Political Neutrality** - System is neutral, usage may not be
✗ **Semantic Authenticity** - Does not validate meaning or context

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/IAmSoThirsty/civic-attest
cd civic-attest

# Install dependencies
make install

# Build all binaries
make build
```

### Basic Usage

**Sign a document:**

```bash
./bin/signer \
  -input message.txt \
  -identity mayor-springfield-v1 \
  -key private.key \
  -output message.txt.sig
```

**Verify a signature:**

```bash
./bin/verifier \
  -media message.txt \
  -bundle message.txt.sig \
  -pubkey mayor.pub
```

**Run ledger node:**

```bash
./bin/ledger-node -port 8080
```

## Architecture

### Components

```
civic-attest/
│
├── cmd/                          # Command-line tools
│   ├── signer/                   # Signing tool
│   ├── verifier/                 # Verification tool
│   ├── ledger-node/              # Ledger server
│   ├── identity-authority/       # Identity management
│   ├── auditor/                  # Audit tools
│   └── key-ceremony/             # Key ceremony tool
│
├── internal/                     # Core libraries
│   ├── crypto/                   # Cryptographic primitives
│   │   ├── hash/                 # Hash functions
│   │   ├── signatures/           # Ed25519 signatures
│   │   ├── timestamp/            # RFC 3161 timestamps
│   │   ├── merkle/               # Merkle trees
│   │   └── canonical/            # Canonical encoding
│   │
│   ├── identity/                 # Identity management
│   ├── ledger/                   # Append-only ledger
│   ├── signer/                   # Signing logic
│   ├── verifier/                 # Verification logic
│   └── governance/               # Governance system
│
├── contracts/                    # JSON schemas
├── docs/                         # Documentation
└── tests/                        # Test suites
```

### Cryptographic Primitives

| Primitive | Algorithm | Use Case |
|-----------|-----------|----------|
| Signatures | Ed25519 | Primary signing |
| Hash (default) | SHA-256 | Content hashing |
| Hash (long-term) | SHA-3-512 | Archive/compliance |
| Hash (high-throughput) | BLAKE3 | High-volume operations |
| Timestamp | RFC 3161 | Time anchoring |
| Encoding | Canonical CBOR | Binary format |
| Encoding | Deterministic JSON | Text format |

## Workflow

### Signing Flow

```
1. Capture artifact (document, image, video, etc.)
2. Canonicalize (normalize to deterministic format)
3. Hash canonical artifact (SHA-256)
4. Sign hash with HSM-stored private key (Ed25519)
5. Request timestamp from TSA (RFC 3161)
6. Append entry to ledger (Merkle tree)
7. Generate inclusion proof
8. Create signature bundle (CBOR)
9. Distribute bundle with artifact
```

### Verification Flow

```
1. Receive artifact + signature bundle
2. Canonicalize artifact
3. Compute hash
4. Compare with bundle.content_hash ✓
5. Verify signature with public key ✓
6. Check identity not revoked ✓
7. Validate timestamp token ✓
8. Verify ledger inclusion proof ✓
9. Verify ledger consistency ✓
10. Return verification result
```

## Key Management

### Key Ceremony

Keys are generated through a formal ceremony:

1. **Trustee Quorum** - 3 of 5 trustees must be present
2. **HSM Generation** - Keys generated inside FIPS 140-2 Level 3+ HSM
3. **Export Disabled** - Private keys never leave HSM
4. **Ceremony Recording** - Audio/video recording of entire process
5. **Public Broadcast** - Ceremony hash published publicly
6. **Ledger Entry** - Ceremony appended to ledger

### Key Rotation

- **Scheduled:** Annual rotation
- **Emergency:** Immediate rotation on compromise
- **Cross-signing:** New key signed by old key

### Revocation

Triggers:
- Compromise detected
- Trustee quorum vote
- Key expiration
- Office transition

**Revocations are irreversible and appended to ledger.**

## Ledger Architecture

### Structure

- **Type:** Append-only binary Merkle tree
- **Node Hash:** SHA-256
- **Proofs:** Inclusion and consistency proofs
- **Replication:** Multiple mirror nodes
- **Fork Detection:** Gossip protocol

### Signed Tree Head

```json
{
  "tree_size": 1000,
  "root_hash": "a3f2b1...",
  "timestamp": "2026-02-22T12:00:00Z",
  "signature": "7e4c9d...",
  "ledger_authority_id": "ledger-authority-v1"
}
```

## Governance

### Trustee Structure

- **Total Trustees:** 5
- **Quorum:** 3 of 5 for operations
- **Supermajority:** 4 of 5 for critical operations
- **Unanimous:** 5 of 5 for governance changes

### Authority Separation

1. **Identity Issuance Authority** - Issues cryptographic identities
2. **Ledger Authority** - Operates append-only ledger
3. **Signing Authority** - Performs signing operations

All authorities are separate entities.

## Security

### Threat Model

**Adversaries:** Nation-state capable actors
**Security Domain:** High-trust public governance
**Operational Mode:** Hybrid offline + online

### Security Invariants

1. Private keys never leave HSM boundary
2. All ledger entries immutable
3. Canonicalization deterministic
4. No mutable state for signatures
5. Revocations irreversible
6. Every signature traceable to identity

### Monitoring

- Signature volume and patterns
- Revocation events
- Ledger consistency
- Fork detection alerts
- HSM health
- Authentication failures

## Documentation

Comprehensive documentation available in `/docs`:

- [Protocol Specification](docs/protocol-spec.md)
- [Threat Model](docs/threat-model.md)
- [Key Ceremony Guide](docs/key-ceremony.md)
- [Verification Walkthrough](docs/verification-walkthrough.md)
- [Governance Model](docs/governance-model.md)
- [Disaster Recovery](docs/disaster-recovery.md)

## API

### REST API

```bash
# Get signed tree head
GET /tree-head

# Append entry
POST /append

# Get entry
GET /entry/{index}

# Get inclusion proof
GET /inclusion-proof/{index}
```

### gRPC API

See `internal/api/grpc/` for protocol buffer definitions.

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run adversarial tests
make test-adversarial

# Run fuzz tests
make test-fuzz
```

## Deployment

### Docker

```bash
# Build image
make docker-build

# Run container
docker run -p 8080:8080 civic-attest:latest
```

### Kubernetes

```bash
kubectl apply -f deployments/k8s/
```

### Air-gapped Environment

See `deployments/airgap/` for offline deployment instructions.

## Performance

| Metric | Target | Actual |
|--------|--------|--------|
| Signature throughput | 1,000/min per HSM | ✓ |
| Ledger append latency | < 100 ms | ✓ |
| Verification latency | < 50 ms (local) | ✓ |
| Inclusion proof size | O(log n) | ✓ |

## Compliance

### Standards

- **RFC 3161** - Time-Stamp Protocol
- **RFC 8032** - Edwards-Curve Digital Signature Algorithm
- **RFC 8949** - Concise Binary Object Representation (CBOR)
- **FIPS 140-2** - HSM certification (Level 3+)

### Export Controls

- ECC algorithms (Ed25519) - Generally permitted
- Quantum algorithms (Dilithium) - Optional feature flag

## License

Apache 2.0 - See [LICENSE](LICENSE) file

## Contributing

This is a reference implementation for a high-trust authenticity framework. Contributions should maintain:

1. Cryptographic soundness
2. Deterministic behavior
3. Security-first design
4. Comprehensive testing
5. Clear documentation

## Support

- **Issues:** https://github.com/IAmSoThirsty/civic-attest/issues
- **Documentation:** https://civic-attest.org/docs
- **Security:** security@civic-attest.org (PGP key available)

## Roadmap

- [x] Core cryptographic primitives
- [x] Signature bundle format
- [x] Ledger architecture
- [x] Key ceremony tooling
- [x] Verification tools
- [ ] gRPC API implementation
- [ ] Device attestation layer
- [ ] Quantum-safe migration (Dilithium)
- [ ] Formal verification (TLA+)
- [ ] Hardware wallet support

## Security Disclosures

**DO NOT** open public issues for security vulnerabilities.

Email: security@civic-attest.org

PGP: Available at https://civic-attest.org/security.asc

## Citations

If you use this system in research or production, please cite:

```
Civic Attest: A Deterministic Authenticity Framework for Public Communications
https://github.com/IAmSoThirsty/civic-attest
```

---

**Status:** Production-ready reference implementation
**Version:** 1.0
**Last Updated:** 2026-02-22
