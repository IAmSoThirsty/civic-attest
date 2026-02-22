# Threat Model

**Version:** 1.0
**Last Updated:** 2026-02-22

## 1. Threat Actors

### 1.1 Nation-State Adversaries

**Capabilities:**
- Advanced persistent threats (APTs)
- Supply chain compromise attempts
- Zero-day exploits
- Physical attacks on infrastructure
- Social engineering at scale
- Cryptanalysis resources

**Motivations:**
- Undermine public trust
- Forge official communications
- Deny authentic communications
- Disrupt critical operations

### 1.2 Insider Threats

**Capabilities:**
- Access to systems and keys
- Knowledge of procedures
- Physical access
- Social trust

**Motivations:**
- Financial gain
- Coercion
- Ideological
- Negligence

### 1.3 Criminal Organizations

**Capabilities:**
- Sophisticated technical skills
- Financial resources
- Social engineering
- Physical threats

**Motivations:**
- Financial fraud
- Extortion
- Disinformation campaigns

## 2. Assets

### 2.1 Critical Assets

1. **Private Signing Keys** - Most critical asset
   - Protection: HSM with export disabled
   - Access: Controlled by trustees
   - Monitoring: All operations logged

2. **Ledger Integrity** - Essential for trust
   - Protection: Append-only structure
   - Verification: Gossip protocol
   - Replication: Multiple mirror nodes

3. **Trustee Credentials** - Key to governance
   - Protection: Multi-factor authentication
   - Access: Physical + cryptographic
   - Monitoring: All actions audited

4. **Timestamp Authority** - Temporal anchoring
   - Protection: RFC 3161 compliant
   - Backup: Multiple TSA providers
   - Monitoring: Timestamp verification

### 2.2 Supporting Assets

- Identity records
- Revocation records
- Audit logs
- Ceremony recordings
- System configurations

## 3. Threat Scenarios

### 3.1 Key Compromise

**Scenario:** Attacker gains access to private signing key

**Attack Vectors:**
- HSM vulnerability exploitation
- Side-channel attacks
- Social engineering of trustees
- Physical theft of HSM
- Supply chain compromise

**Mitigations:**
- FIPS 140-2 Level 3+ HSM
- Export disabled
- Tamper detection
- Multi-party authorization
- Regular key rotation
- Emergency revocation procedures

**Detection:**
- Unusual signing patterns
- Unauthorized signatures discovered
- HSM tamper alerts
- Trustee reports

**Response:**
1. Immediate key revocation
2. Ledger entry appended
3. Public advisory issued
4. New key ceremony
5. Investigation initiated

### 3.2 Ledger Fork Attack

**Scenario:** Attacker creates divergent ledger view

**Attack Vectors:**
- Compromise ledger node
- Network partition exploitation
- DNS/BGP hijacking
- Malicious mirror node

**Mitigations:**
- Gossip protocol for fork detection
- Multiple independent mirrors
- Signed tree heads
- Consistency proofs
- Certificate pinning

**Detection:**
- Gossip protocol divergence alerts
- Multiple signed tree heads for same size
- Mirror node disagreement
- User reports of inconsistencies

**Response:**
1. Freeze ledger writes
2. Compare all mirror nodes
3. Identify correct tree head
4. Revoke compromised node authority
5. Restore from valid state
6. Public disclosure

### 3.3 Timestamp Manipulation

**Scenario:** Attacker backdates or forward-dates signatures

**Attack Vectors:**
- Compromise TSA
- Replay old timestamp tokens
- Time synchronization attacks
- TSA key compromise

**Mitigations:**
- RFC 3161 compliant TSA
- Multiple independent TSAs
- Timestamp verification against ledger
- NTP security
- Token freshness checks

**Detection:**
- Timestamp inconsistent with ledger
- Duplicate timestamp serial numbers
- TSA signature verification failure
- Chronological anomalies

**Response:**
1. Verify with alternate TSAs
2. Investigate TSA compromise
3. Revoke affected signatures if necessary
4. Switch to backup TSA
5. Public notification

### 3.4 Canonicalization Attack

**Scenario:** Attacker finds canonicalization collision

**Attack Vectors:**
- Exploit canonical encoding weakness
- Find encoding that produces same hash
- Unicode normalization issues
- Metadata manipulation

**Mitigations:**
- Strict canonical encoding (CBOR/JSON)
- No custom encodings allowed
- Reference implementation mandatory
- Extensive testing
- Formal verification

**Detection:**
- Different artifacts with same hash
- Canonicalization non-determinism
- Hash collision detected

**Response:**
1. Immediate protocol freeze
2. Security advisory
3. Analyze collision
4. Update canonical spec if needed
5. Re-sign affected content

### 3.5 Revocation DoS

**Scenario:** Attacker forces mass revocations

**Attack Vectors:**
- Compromise detection false positives
- Social engineering of trustees
- Automated revocation trigger exploitation
- Trustee credential theft

**Mitigations:**
- Quorum requirement for revocation
- Rate limiting on revocations
- Multi-factor trustee auth
- Revocation cooldown periods
- Emergency override requires supermajority

**Detection:**
- Unusual revocation rate
- Multiple revocations in short period
- Trustee voting anomalies

**Response:**
1. Pause automated revocations
2. Verify trustee identities
3. Investigate trigger cause
4. Restore improperly revoked keys if needed
5. Adjust revocation policies

### 3.6 Side-Channel Attacks

**Scenario:** Extract keys via side channels

**Attack Vectors:**
- Timing attacks on HSM
- Power analysis
- EM radiation analysis
- Acoustic cryptanalysis
- Cache timing

**Mitigations:**
- HSM with side-channel protections
- Constant-time operations
- Physical security
- EM shielding
- Regular security audits

**Detection:**
- HSM tamper detection
- Physical security monitoring
- Unusual access patterns

**Response:**
1. Investigate HSM integrity
2. Rotate keys if compromise suspected
3. Enhanced physical security
4. HSM replacement if necessary

### 3.7 Social Engineering

**Scenario:** Trick trustees or operators

**Attack Vectors:**
- Phishing for credentials
- Pretexting for information
- Impersonation
- Baiting with malware
- Quid pro quo

**Mitigations:**
- Security awareness training
- Multi-factor authentication
- Out-of-band verification
- Strict procedures
- Dual control requirements

**Detection:**
- Unusual requests
- Procedural violations
- Failed authentication attempts
- User reports

**Response:**
1. Verify incident
2. Revoke compromised credentials
3. Additional training
4. Procedure review
5. Incident analysis

## 4. Security Properties

### 4.1 Maintained Properties

Under all attack scenarios, system maintains:

1. **Append-only ledger** - Historical entries immutable
2. **Non-repudiation** - Valid signatures traceable to identity
3. **Transparency** - All operations publicly auditable
4. **Revocation irreversibility** - Cannot un-revoke

### 4.2 Degraded Operations

System can operate in degraded mode:

1. **Offline verification** - Without live ledger access
2. **Read-only ledger** - During investigation
3. **Backup TSA** - If primary TSA unavailable
4. **Emergency governance** - Supermajority override

## 5. Security Boundaries

### 5.1 Trust Boundaries

1. **HSM boundary** - Private keys never cross
2. **Trustee quorum** - No single point of trust
3. **Ledger authority** - Separate from signing
4. **Network boundary** - Assume hostile network

### 5.2 Failure Modes

**Safe Failures:**
- Fail closed (reject on doubt)
- Ledger freeze on inconsistency
- Signature rejection on any error

**Unsafe Failures:**
- ❌ Accept invalid signature
- ❌ Allow ledger modification
- ❌ Bypass revocation check

## 6. Monitoring and Detection

### 6.1 Security Monitoring

Continuous monitoring for:
- Unusual signature patterns
- Revocation anomalies
- Ledger inconsistencies
- Failed authentication
- HSM alerts
- Network anomalies

### 6.2 Audit Logging

Immutable logs for:
- All signature operations
- All revocations
- All trustee actions
- All ledger operations
- All governance decisions

## 7. Incident Response

### 7.1 Response Tiers

**Tier 1: Minor**
- Failed authentication
- Network issues
- User errors

**Tier 2: Moderate**
- Suspected compromise attempt
- Unusual patterns
- Component failure

**Tier 3: Critical**
- Key compromise
- Ledger corruption
- Active attack

### 7.2 Response Procedures

1. Detection and triage
2. Containment
3. Investigation
4. Remediation
5. Recovery
6. Post-incident review
7. Public disclosure (if warranted)

---

**Conclusion:** This threat model assumes competent, well-resourced adversaries and designs defenses accordingly. Regular review and updates required as threat landscape evolves.
