# HSM Operational Constraints Specification

**Version:** 2.0
**Status:** Enhanced Security Requirements
**Last Updated:** 2026-02-23
**Classification:** CONFIDENTIAL - TRUSTEES AND OPERATORS ONLY

## Abstract

This document specifies enhanced operational constraints and security controls for Hardware Security Modules (HSMs) used in the Civic Attest system. These requirements address nation-state level threats and prevent single points of failure.

## 1. HSM Requirements

### 1.1 Certification

**Minimum:** FIPS 140-2 Level 3
**Recommended:** FIPS 140-3 Level 3 or higher
**Additional:** Common Criteria EAL4+

**Required Features:**
- Tamper detection and response
- Physical access controls
- Role-based access control
- Audit logging
- Key export disabled
- Zeroization on tamper detection

### 1.2 Approved HSM Models

**Primary Options:**
- Thales Luna Network HSM 7
- nCipher nShield Connect+
- Utimaco SecurityServer Se Gen2
- AWS CloudHSM (FIPS 140-2 Level 3)

**Prohibited:**
- Software-only HSMs
- FIPS 140-2 Level 1 or 2
- End-of-life models
- Uncertified devices

## 2. Rate Limiting Controls

### 2.1 Signature Rate Limits

**Enforced at Hardware Level:**

**Standard Operations:**
- Maximum: 1,000 signatures per minute per HSM
- Burst window: 10 seconds
- Burst maximum: 100 signatures
- Cooldown after burst: 30 seconds

**Rate Limit Configuration:**
```yaml
hsm_rate_limits:
  operations_per_minute: 1000
  burst_window_seconds: 10
  burst_maximum: 100
  cooldown_seconds: 30
  action_on_exceed: REJECT
  alert_threshold_percent: 80
  log_all_rejections: true
```

**Enforcement:**
- Implemented inside HSM firmware
- Cannot be bypassed by operator
- Rate limit counters stored in tamper-resistant memory
- Alerts sent when threshold reached

**Per-Identity Rate Limits:**
```yaml
identity_rate_limits:
  standard_identity:
    signatures_per_hour: 100
    signatures_per_day: 1000

  verified_organization:
    signatures_per_hour: 1000
    signatures_per_day: 10000

  high_volume_entity:
    signatures_per_hour: 10000
    signatures_per_day: 100000
    requires_approval: true
```

### 2.2 Authentication Rate Limits

**Login Attempts:**
- Maximum: 5 failed attempts per account per hour
- Lockout duration: 1 hour (first offense)
- Lockout duration: 24 hours (subsequent offenses)
- Permanent lockout: After 10 lockouts in 30 days

**Session Limits:**
- Maximum concurrent sessions per operator: 1
- Session timeout: 8 hours
- Idle timeout: 30 minutes
- Re-authentication required after timeout

## 3. Anti-Rollback Protection

### 3.1 Firmware Version Enforcement

**Monotonic Version Counter:**
- Hardware-enforced monotonic counter
- Version number cannot decrease
- Rollback attempts logged and rejected
- Counter stored in tamper-resistant memory

**Firmware Update Process:**
1. Vendor releases new firmware with version N+1
2. Vendor signs firmware with vendor key
3. Trustee quorum reviews firmware (security audit required)
4. Minimum 3-of-5 trustees co-sign firmware
5. HSM verifies dual signature (vendor + trustee quorum)
6. HSM checks: new_version > current_version
7. HSM updates firmware
8. HSM increments monotonic counter to new_version
9. Rollback permanently impossible

**Firmware Signing Schema:**
```json
{
  "firmware_version": "2.4.1",
  "firmware_hash": "a3f2b1c4...",
  "vendor_signature": "7e4c9d...",
  "trustee_signatures": [
    {
      "trustee_id": "trustee-1",
      "signature": "8d3e2f...",
      "signed_at": "2026-02-23T12:00:00Z"
    },
    {
      "trustee_id": "trustee-2",
      "signature": "9e4f3a...",
      "signed_at": "2026-02-23T12:05:00Z"
    },
    {
      "trustee_id": "trustee-3",
      "signature": "1a2b3c...",
      "signed_at": "2026-02-23T12:10:00Z"
    }
  ],
  "quorum": "3-of-5",
  "previous_version": "2.4.0",
  "release_date": "2026-02-20"
}
```

### 3.2 Configuration Rollback Protection

**Configuration Version:**
- All configuration changes versioned
- Configuration version monotonically increasing
- Rollback to previous configuration logged and requires justification
- Automatic rollback only for emergency (with 3-of-5 approval)

## 4. Dual-Control Activation

### 4.1 M-of-N Authentication

**Requirement:** Minimum 2-of-3 operator authentication

**Authentication Factors (all required per operator):**
1. **Physical Presence:** Proximity token or smartcard
2. **Biometric:** Fingerprint or iris scan
3. **Knowledge:** PIN (minimum 12 characters)
4. **Hardware Token:** FIDO2 security key

**Activation Procedure:**
1. Operator 1 presents all 4 factors
2. HSM logs authentication attempt
3. HSM waits for Operator 2 (within 5-minute window)
4. Operator 2 presents all 4 factors (must be different operator)
5. HSM verifies both authentications
6. HSM verifies operators have distinct roles
7. HSM unlocks for operational session
8. Session valid for 8 hours (idle timeout: 30 minutes)
9. Re-authentication required after session expiry

**Operator Separation:**
- Operator 1 and Operator 2 must be different individuals
- No operator can authenticate twice
- Minimum separation: Different physical locations
- Recommended: Different organizations

### 4.2 Role-Based Access Control

**Roles:**

**1. HSM Administrator:**
- Firmware updates (with trustee approval)
- Configuration changes
- Audit log access
- Cannot perform signing operations

**2. Signing Operator:**
- Signature operations only
- Cannot modify configuration
- Cannot access admin functions

**3. Audit Officer:**
- Read-only access to logs
- Cannot perform operations
- Cannot modify configuration

**4. Emergency Operator:**
- Limited emergency functions
- Requires 3-of-5 trustee approval
- All actions heavily logged
- Time-limited authorization

**Role Separation:**
- No operator can hold multiple roles
- Role assignment requires 3-of-5 trustee approval
- Role changes logged to transparency ledger
- Annual role review mandatory

## 5. Sealed Audit Logging

### 5.1 Comprehensive Logging

**All Events Logged:**
- Every signature operation (identity, hash, result)
- Every authentication attempt (success and failure)
- Every configuration change
- Every error condition
- Every rate limit enforcement
- Every alarm/alert triggered
- Clock synchronization events
- Tamper detection events

**Log Entry Format:**
```json
{
  "timestamp": "2026-02-23T12:00:00.123456Z",
  "sequence_number": 1234567,
  "hsm_id": "hsm-primary-1",
  "event_type": "SIGNATURE_OPERATION",
  "operator_id": "operator-1",
  "operator_role": "SIGNING_OPERATOR",
  "identity_id": "mayor-springfield-v1",
  "content_hash": "a3f2b1c4d5e6...",
  "signature_hash": "7e8f9a0b1c2d...",
  "operation_status": "SUCCESS",
  "operation_duration_ms": 45,
  "rate_limit_remaining": 985,
  "log_signature": "ed4c3b2a1098...",
  "previous_log_hash": "9f8e7d6c5b4a..."
}
```

### 5.2 Log Integrity

**Cryptographic Sealing:**
- Each log entry signed by HSM
- Chained hash linking: hash(current_entry + previous_hash)
- Tampering detectable (breaks chain)
- Log signature key separate from signing key
- Log signature key never exported

**Chain Structure:**
```
Entry N: {
  data: {...},
  previous_hash: hash(Entry N-1),
  signature: sign(hash(data || previous_hash))
}
```

**Validation:**
```python
def validate_log_chain(entries):
    for i in range(1, len(entries)):
        current = entries[i]
        previous = entries[i-1]

        # Verify previous hash
        computed_hash = hash(previous)
        if current.previous_hash != computed_hash:
            return False, f"Chain broken at entry {i}"

        # Verify signature
        data_to_sign = current.data + current.previous_hash
        if not verify_signature(data_to_sign, current.signature, hsm_log_pubkey):
            return False, f"Invalid signature at entry {i}"

    return True, "Chain valid"
```

### 5.3 Log Export and Archival

**Real-Time Export:**
- Stream to separate log server via encrypted channel
- TLS 1.3 with mutual authentication
- Log server cannot modify HSM logs
- Connection loss triggers alarm

**Hourly Export:**
- Export to append-only storage
- Encrypted with separate key
- Multiple geographic copies
- Automated integrity verification

**Daily Export:**
- Export to offline archive storage
- Physical media (write-once)
- Offsite storage in vault
- Annual verification of archive integrity

**Weekly Export:**
- Export to geographically distant location
- Different jurisdiction
- Different physical security domain

**Export Configuration:**
```yaml
log_export:
  realtime:
    enabled: true
    destination: syslog-server.internal:6514
    protocol: TLS_1_3
    mutual_auth: true
    connection_timeout_alarm: 60s

  hourly:
    enabled: true
    destination: s3://civic-attest-logs-primary/
    encryption: AES-256-GCM
    key_rotation: daily
    integrity_check: true

  daily:
    enabled: true
    destination: /mnt/archive/hsm-logs/
    media: WORM
    offsite_copy: true
    vault_storage: true

  weekly:
    enabled: true
    destination: geo-distant-archive
    jurisdiction: different
    verification_required: true
```

## 6. Threshold Key Signing (Optional Advanced Deployment)

### 6.1 Overview

**Purpose:** Eliminate single HSM compromise risk

**Architecture:** 2-of-3 threshold signature scheme

**Benefits:**
- No single HSM can sign alone
- Geographic distribution possible
- Byzantine fault tolerance
- Enhanced resilience

### 6.2 Threshold Signature Scheme

**Algorithm:** ECDSA threshold signatures (compatible with Ed25519)

**Key Generation:**
1. 3 HSMs participate in distributed key generation (DKG)
2. Each HSM generates key share
3. No HSM knows complete private key
4. Public key derived from shares
5. Minimum 2 HSMs required to sign

**Signing Protocol:**
1. Content hash distributed to all HSMs
2. Each participating HSM (minimum 2) computes signature share
3. Signature shares combined to form complete signature
4. Combined signature verifiable with public key
5. Individual shares cannot produce valid signature

**Configuration:**
```yaml
threshold_signing:
  enabled: true
  threshold: 2
  total_shares: 3
  hsm_locations:
    - hsm-1: datacenter-us-east
    - hsm-2: datacenter-eu-west
    - hsm-3: datacenter-ap-south
  key_share_storage: tamper_resistant_memory
  signature_timeout: 30s
  quorum_requirement: any_2_of_3
```

### 6.3 Geographic Distribution

**Deployment:**
- HSM 1: Primary datacenter (jurisdiction A)
- HSM 2: Secondary datacenter (jurisdiction B)
- HSM 3: Tertiary datacenter (jurisdiction C)

**Benefits:**
- Survives physical destruction of 1 HSM
- Survives compromise of 1 HSM
- Survives jurisdiction-level attack
- Network partition tolerant (with degradation)

## 7. Physical Security

### 7.1 Physical Access Controls

**HSM Housing:**
- Locked cage or vault
- 24/7 video surveillance (90-day retention)
- Access log (badge reader)
- Two-person rule for access
- Alarm on unauthorized access

**Environmental Controls:**
- Temperature monitoring (alert on deviation)
- Humidity control
- Fire suppression system
- UPS backup power
- Seismic protection (if applicable)

### 7.2 Tamper Detection

**Tamper Sensors:**
- Intrusion detection
- Temperature sensors (detect heating attacks)
- Voltage sensors (detect power analysis)
- Light sensors (detect case opening)
- Vibration sensors

**Tamper Response:**
1. Immediate zeroization of all keys
2. Alarm to security operations center
3. Log event (if possible before zeroization)
4. Lock HSM (requires trustee unlock)
5. Incident investigation initiated

## 8. Operational Procedures

### 8.1 Daily Checks

**Automated:**
- Health check every 5 minutes
- Log export verification hourly
- Rate limit counter check
- Clock synchronization check

**Manual (daily):**
- Review overnight logs
- Verify no alarms
- Check rate limit usage
- Verify backup HSM ready

### 8.2 Weekly Procedures

**Manual:**
- Review weekly log summary
- Verify log archive integrity
- Test backup HSM activation
- Review operator access logs
- Update operational documentation

### 8.3 Monthly Procedures

**Manual:**
- Full audit of HSM configuration
- Review all operator roles
- Test emergency procedures
- Verify physical security
- Update risk assessment

### 8.4 Quarterly Procedures

**Manual:**
- Independent security audit
- Penetration testing (authorized)
- Firmware update review
- Trustee review of operations
- Disaster recovery drill

## 9. Emergency Procedures

### 9.1 HSM Compromise Suspected

**Immediate Actions:**
1. Activate backup HSM
2. Freeze signing operations on suspect HSM
3. Alert all trustees
4. Initiate forensic investigation
5. Revoke keys if compromise confirmed

### 9.2 HSM Failure

**Immediate Actions:**
1. Activate backup HSM (within 4 hours)
2. Trustee quorum unlock backup
3. Verify backup integrity
4. Resume operations
5. Schedule primary HSM repair/replacement

### 9.3 Tamper Detection

**Immediate Actions:**
1. HSM automatically zeroizes keys
2. Security operations center notified
3. Physical investigation of HSM location
4. Trustee quorum convened
5. Incident response plan activated

## 10. Compliance and Audit

### 10.1 Audit Requirements

**Quarterly:**
- Independent security audit
- Review all logs
- Verify procedures followed
- Test emergency procedures

**Annual:**
- Full FIPS compliance audit
- Penetration testing
- Physical security audit
- Trustee comprehensive review

### 10.2 Metrics

**Track and Report:**
- Signature operation count (per day/week/month)
- Rate limit enforcements
- Authentication failures
- Configuration changes
- Alarms triggered
- HSM uptime/availability
- Log export success rate

## 11. Configuration Management

### 11.1 Change Control

**All Changes Require:**
1. Change proposal documentation
2. Security impact assessment
3. 3-of-5 trustee approval
4. Testing in non-production environment
5. Scheduled maintenance window
6. Rollback plan documented
7. Post-change verification

### 11.2 Configuration Backup

**Backup Schedule:**
- Immediate: After every configuration change
- Daily: Automated full configuration backup
- Weekly: Offsite configuration backup

**Backup Contents:**
- HSM configuration (no private keys)
- Operator roles and permissions
- Rate limit configurations
- Network settings
- Audit settings

## 12. Appendices

### Appendix A: HSM Procurement Checklist

See separate document: `hsm-procurement-checklist.md`

### Appendix B: Operator Training Requirements

See separate document: `hsm-operator-training.md`

### Appendix C: Emergency Contact List

CONFIDENTIAL - Maintained separately by trustees

---

**Document Classification:** CONFIDENTIAL - TRUSTEES AND OPERATORS ONLY
**Review Frequency:** Quarterly
**Next Review:** 2026-05-23
**Owner:** Chief Security Officer + Trustee Board
