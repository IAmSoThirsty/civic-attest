# Verification Walkthrough

**Version:** 1.0
**Audience:** Verifiers, System Integrators
**Last Updated:** 2026-02-22

## 1. Introduction

This guide walks through the complete process of verifying a signature bundle in the Civic Attest system. It covers both online and offline verification modes.

## 2. Quick Start

### 2.1 Installation

```bash
# Build from source
git clone https://github.com/IAmSoThirsty/civic-attest
cd civic-attest
make build
```

### 2.2 Basic Verification

```bash
./bin/verifier \
  -media message.txt \
  -bundle message.sig \
  -pubkey pubkey.hex
```

### 2.3 Expected Output

```
=== Civic Attest Verifier ===

✓ Hash verification: PASSED
✓ Signature verification: PASSED
✓ Timestamp token: PRESENT
✓ Ledger inclusion: VERIFIED

=== VERIFICATION SUCCESSFUL ===
Signer Identity: mayor-springfield-v1
Key Version: 1
Bundle Version: 1.0
```

## 3. Detailed Verification Steps

### Step 1: Obtain Required Inputs

You need three components:

1. **Media File** - The signed content (document, image, video, etc.)
2. **Signature Bundle** - The `.sig` file containing cryptographic proofs
3. **Public Key** - The signer's public key

**Example:**
```bash
# Download from official source (replace with actual government URLs)
curl -O https://example.gov/announcement.txt
curl -O https://example.gov/announcement.txt.sig
curl -O https://example.gov/keys/mayor-v1.pub
```

### Step 2: Canonicalize Media

The verifier automatically canonicalizes the media file to match the signing process.

**What happens internally:**

```go
// Read media
media, _ := os.ReadFile("announcement.txt")

// Canonicalize using CBOR
canonical, _ := canonical.Encode(media, canonical.CBOR)

// Compute hash
hash, _ := hash.Hash(canonical, hash.SHA256)
```

**Canonical Formats Supported:**
- CBOR (default)
- Deterministic JSON

### Step 3: Verify Content Hash

The verifier compares the computed hash with the hash in the signature bundle.

**Check performed:**
```go
if !bytes.Equal(computedHash, bundle.ContentHash) {
    return errors.New("content hash mismatch")
}
```

**What this proves:** The media has not been altered since signing.

**Failure modes:**
- ❌ File modified
- ❌ Wrong canonical format
- ❌ Corrupted bundle

### Step 4: Verify Digital Signature

The verifier checks the cryptographic signature against the public key.

**Verification:**
```go
valid := ed25519.Verify(
    publicKey,
    bundle.ContentHash,  // Message signed
    bundle.Signature,    // Signature
)
```

**What this proves:** The signature was created by the holder of the private key.

**Failure modes:**
- ❌ Wrong public key
- ❌ Signature tampered
- ❌ Bundle corrupted

### Step 5: Verify Identity

The verifier checks that the identity is valid and not revoked.

**Checks performed:**
1. Identity exists in system
2. Status is "active"
3. Current time within validity period
4. No revocation record

**Example identity check:**
```go
identity := getIdentity(bundle.SignerIdentityID)

if identity.Status != "active" {
    return errors.New("identity not active")
}

if time.Now().After(identity.ValidTo) {
    return errors.New("identity expired")
}

if hasRevocation(identity.IdentityID) {
    return errors.New("identity revoked")
}
```

**What this proves:** The key was authorized at signing time.

### Step 6: Verify Timestamp

The verifier validates the RFC 3161 timestamp token.

**Checks:**
1. Token signature valid
2. Token references correct content hash
3. TSA is trusted
4. Timestamp within acceptable range

**What this proves:** The signature was created at a specific time.

### Step 7: Verify Ledger Inclusion

The verifier checks the Merkle inclusion proof.

**Online mode:**
```bash
./bin/verifier \
  -media announcement.txt \
  -bundle announcement.txt.sig \
  -pubkey mayor-v1.pub
  # Connects to ledger node
```

**Offline mode:**
```bash
./bin/verifier \
  -media announcement.txt \
  -bundle announcement.txt.sig \
  -pubkey mayor-v1.pub \
  -offline
  # Uses cached signed tree head
```

**What happens:**
```go
proof := bundle.MerkleInclusionProof

// Verify proof against ledger root
valid := verifyInclusionProof(
    proof.LeafHash,
    proof.LeafIndex,
    proof.TreeSize,
    proof.Path,
    ledgerRootHash,
)
```

**What this proves:** The signature was permanently recorded in the public ledger.

### Step 8: Verify Ledger Consistency

The verifier checks that the ledger hasn't forked.

**Procedure:**
1. Retrieve current signed tree head (STH)
2. Retrieve STH at bundle creation time
3. Verify consistency proof
4. Check for gossip alerts

**What this proves:** The ledger maintains append-only property.

## 4. Verification Modes

### 4.1 Standard Online Verification

**Features:**
- Full ledger connectivity
- Real-time revocation checking
- Latest signed tree head
- Gossip protocol participation

**Command:**
```bash
./bin/verifier \
  -media file.txt \
  -bundle file.txt.sig \
  -pubkey signer.pub
```

**Use when:**
- Network connectivity available
- Maximum assurance required
- Real-time verification needed

### 4.2 Offline Verification

**Features:**
- No network required
- Uses cached data
- Portable verification
- Privacy preserving

**Command:**
```bash
./bin/verifier \
  -media file.txt \
  -bundle file.txt.sig \
  -pubkey signer.pub \
  -offline
```

**Requires:**
- Recent signed tree head cached
- Public key cached
- Revocation list cached (if checking revocations)

**Use when:**
- No network connectivity
- Air-gapped environment
- Privacy concerns

### 4.3 Audit Verification

**Features:**
- Full ledger replay
- All historical checks
- Comprehensive validation
- Detailed reporting

**Command:**
```bash
./bin/verifier \
  -media file.txt \
  -bundle file.txt.sig \
  -pubkey signer.pub \
  -audit
```

**Use when:**
- Forensic analysis
- Legal proceedings
- Incident investigation
- Compliance audits

## 5. Verification Result Interpretation

### 5.1 Successful Verification

```
=== VERIFICATION SUCCESSFUL ===
Signer Identity: mayor-springfield-v1
Key Version: 1
Bundle Version: 1.0
```

**Meaning:** All checks passed. The signature is cryptographically valid.

**What you can conclude:**
✓ Content is authentic (from stated signer)
✓ Content is unmodified (byte-level integrity)
✓ Signature was created at stated time
✓ Signature is publicly recorded
✓ Signer's key was valid at signing time

**What you CANNOT conclude:**
✗ Content is truthful (only that it's from stated source)
✗ Signer acted voluntarily (no coercion detection)
✗ Content is legally binding (depends on jurisdiction)

### 5.2 Failed Verification

```
=== VERIFICATION FAILED ===
Error: Content hash mismatch
```

**Possible causes:**
- File has been modified
- Wrong file provided
- Bundle corruption
- Using wrong canonical format

**Action:** DO NOT trust the content. Obtain fresh copy from authoritative source.

### 5.3 Warnings

```
=== VERIFICATION SUCCESSFUL ===
⚠ Timestamp token: MISSING
```

**Meaning:** Verification passed but with caveats.

**Common warnings:**
- Missing timestamp (still valid, but no time proof)
- Offline mode (didn't check latest revocations)
- Old ledger snapshot (may be outdated)

## 6. Troubleshooting

### 6.1 Hash Mismatch

**Error:**
```
❌ Hash verification: FAILED
   Expected: a3f2b1...
   Computed: b7e4c9...
```

**Solutions:**
1. Verify you have the correct file
2. Check file hasn't been modified
3. Ensure same canonical format as signing
4. Try re-downloading bundle

### 6.2 Invalid Signature

**Error:**
```
❌ Signature verification: FAILED
```

**Solutions:**
1. Verify public key is correct
2. Check bundle integrity
3. Ensure bundle matches media file
4. Verify key algorithm matches

### 6.3 Revoked Key

**Error:**
```
❌ Identity verification: FAILED
Error: Identity has been revoked
```

**Meaning:** The signing key was revoked. Signature may have been valid when created, but key is no longer trusted.

**Action:** Check revocation date vs. signature date. Contact issuing authority.

### 6.4 Ledger Connection Failed

**Error:**
```
❌ Ledger inclusion: FAILED
Error: Connection refused
```

**Solutions:**
1. Check network connectivity
2. Verify ledger node URL
3. Use offline mode with cached data
4. Try alternative ledger mirror

## 7. Integration Examples

### 7.1 Command Line

```bash
#!/bin/bash
# Verify script

if ./bin/verifier -media "$1" -bundle "$2" -pubkey "$3"; then
    echo "✓ Verified"
    exit 0
else
    echo "✗ Verification failed"
    exit 1
fi
```

### 7.2 Python Integration

```python
import subprocess

def verify_signature(media_file, bundle_file, pubkey_file):
    result = subprocess.run([
        './bin/verifier',
        '-media', media_file,
        '-bundle', bundle_file,
        '-pubkey', pubkey_file
    ], capture_output=True)

    return result.returncode == 0
```

### 7.3 Go Integration

```go
import "os/exec"

func VerifySignature(media, bundle, pubkey string) (bool, error) {
    cmd := exec.Command("./bin/verifier",
        "-media", media,
        "-bundle", bundle,
        "-pubkey", pubkey)

    err := cmd.Run()
    return err == nil, err
}
```

## 8. Advanced Topics

### 8.1 Batch Verification

Verify multiple files efficiently:

```bash
for sig in *.sig; do
    media="${sig%.sig}"
    if ./bin/verifier -media "$media" -bundle "$sig" -pubkey mayor.pub; then
        echo "✓ $media"
    else
        echo "✗ $media"
    fi
done
```

### 8.2 Automated Monitoring

Set up continuous verification:

```bash
# Crontab entry
*/5 * * * * /usr/local/bin/verify-watch.sh
```

### 8.3 Custom Verification Policies

Implement organization-specific checks:

```go
func CustomVerify(bundle *Bundle) error {
    // Standard verification
    if err := StandardVerify(bundle); err != nil {
        return err
    }

    // Custom: Require timestamp within 1 hour
    if time.Since(bundle.Timestamp) > time.Hour {
        return errors.New("signature too old")
    }

    // Custom: Require specific jurisdiction
    if bundle.Identity.Jurisdiction != "Springfield" {
        return errors.New("wrong jurisdiction")
    }

    return nil
}
```

## 9. Best Practices

### 9.1 For Verifiers

✓ **DO:**
- Verify immediately upon receipt
- Use online mode when possible
- Check revocation status
- Archive verified bundles
- Maintain cached public keys
- Report verification failures

✗ **DON'T:**
- Skip verification steps
- Trust expired signatures
- Ignore warnings
- Use untrusted public keys
- Verify without hash check

### 9.2 For System Integrators

✓ **DO:**
- Implement automated verification
- Log all verification attempts
- Handle errors gracefully
- Support offline mode
- Cache public keys securely
- Monitor verification rates

✗ **DON'T:**
- Disable security checks
- Ignore verification failures
- Hard-code public keys
- Skip revocation checks
- Trust external verification claims

---

**Support:** For questions or issues, see https://github.com/IAmSoThirsty/civic-attest/issues
