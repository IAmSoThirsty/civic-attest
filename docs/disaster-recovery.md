# Disaster Recovery Plan

**Version:** 1.0
**Classification:** CONFIDENTIAL - TRUSTEES ONLY
**Last Updated:** 2026-02-22

## 1. Overview

This document outlines procedures for recovering from catastrophic failures in the Civic Attest system.

## 2. Disaster Scenarios

### 2.1 Scenario: Complete HSM Loss

**Cause:** Fire, flood, theft, or hardware destruction

**Impact:**
- Cannot create new signatures
- Existing signatures still verifiable
- Identity issuance halted

**Recovery Priority:** CRITICAL
**RTO:** 4 hours
**RPO:** 0 (no data loss)

**Recovery Procedure:**

**Step 1: Immediate Response (0-30 minutes)**
1. Confirm HSM loss
2. Alert all trustees
3. Activate backup HSM from vault
4. Transport to secure facility

**Step 2: HSM Activation (30-120 minutes)**
1. Trustee quorum assembles
2. Unseal backup HSM
3. Verify HSM integrity
4. Test signature operation
5. Verify against public key

**Step 3: Service Restoration (120-180 minutes)**
1. Configure backup HSM
2. Update system configuration
3. Test end-to-end signing
4. Restore signing service

**Step 4: Investigation (180+ minutes)**
1. Investigate cause of loss
2. Document incident
3. Public notification (if required)
4. Acquire replacement HSM

**Checklist:**
- [ ] Backup HSM retrieved
- [ ] Quorum assembled
- [ ] HSM unsealed and verified
- [ ] Service restored
- [ ] Incident documented
- [ ] Public notified (if required)

### 2.2 Scenario: Key Compromise Detected

**Cause:** Security breach, insider threat, HSM vulnerability

**Impact:**
- All signatures using compromised key untrusted
- Identity must be revoked
- Emergency rotation required

**Recovery Priority:** CRITICAL
**RTO:** 24 hours
**RPO:** N/A

**Recovery Procedure:**

**Step 1: Containment (0-60 minutes)**
1. Confirm compromise
2. Immediately disable affected HSM
3. Alert all trustees
4. Emergency quorum convened
5. Block all signing operations

**Step 2: Revocation (60-180 minutes)**
1. Create revocation record
2. Trustee quorum signs revocation
3. Append to ledger
4. Publish signed advisory
5. Notify all verifiers

**Step 3: Key Rotation (3-12 hours)**
1. Schedule emergency key ceremony
2. Generate new key pair
3. Create new identity
4. Cross-sign with old key (if possible)
5. Publish new public key

**Step 4: Trust Reestablishment (12-24 hours)**
1. Reissue trust anchors
2. Update all verifier systems
3. Communicate to stakeholders
4. Monitor for fraudulent signatures

**Step 5: Investigation (ongoing)**
1. Full forensic analysis
2. Identify attack vector
3. Implement additional controls
4. Complete incident report
5. Public disclosure (coordinated)

**Checklist:**
- [ ] Compromise confirmed
- [ ] Affected HSM disabled
- [ ] Revocation record created
- [ ] Ledger updated
- [ ] Advisory published
- [ ] New key ceremony completed
- [ ] Trust anchors reissued
- [ ] Investigation underway

### 2.3 Scenario: Ledger Corruption

**Cause:** Database corruption, ransomware, insider attack

**Impact:**
- Cannot append new entries
- Existing proofs may be invalid
- Public trust compromised

**Recovery Priority:** CRITICAL
**RTO:** 72 hours
**RPO:** 0 (no entry loss acceptable)

**Recovery Procedure:**

**Step 1: Detection and Freeze (0-30 minutes)**
1. Corruption detected (gossip protocol alert)
2. FREEZE all ledger writes immediately
3. Alert trustees and operators
4. Begin logging all events

**Step 2: Assessment (30-120 minutes)**
1. Identify corruption extent
2. Compare all mirror nodes
3. Identify last known good state
4. Document divergence points
5. Preserve corrupted state for forensics

**Step 3: Recovery (2-24 hours)**
1. Trustee quorum reviews findings
2. Approve recovery plan
3. Identify canonical ledger state
4. Restore from last valid signed tree head
5. Rebuild forward if possible
6. Verify integrity

**Step 4: Validation (24-48 hours)**
1. Verify all inclusion proofs
2. Verify all consistency proofs
3. Compare with all mirrors
4. Run full audit
5. Trustee approval to resume

**Step 5: Resumption (48-72 hours)**
1. Resume ledger writes
2. Enhanced monitoring
3. Public disclosure
4. Ongoing verification

**Checklist:**
- [ ] Writes frozen
- [ ] Corruption extent identified
- [ ] Valid state identified
- [ ] Recovery approved
- [ ] Ledger restored
- [ ] Full audit passed
- [ ] Writes resumed
- [ ] Public notified

### 2.4 Scenario: Complete Data Center Loss

**Cause:** Natural disaster, power loss, network partition

**Impact:**
- All services offline
- Cannot verify signatures
- Cannot issue new signatures

**Recovery Priority:** HIGH
**RTO:** 24 hours
**RPO:** 1 hour

**Recovery Procedure:**

**Step 1: Failover (0-60 minutes)**
1. Activate disaster recovery site
2. Restore from backups
3. Verify data integrity
4. DNS failover to DR site

**Step 2: Service Restoration (1-6 hours)**
1. Restore ledger nodes
2. Restore identity authority
3. Restore signer service
4. Verify end-to-end operation

**Step 3: Data Synchronization (6-12 hours)**
1. Sync with surviving mirrors
2. Verify ledger consistency
3. Update signed tree heads
4. Verify all services

**Step 4: Full Operation (12-24 hours)**
1. Return to full capacity
2. Monitor for issues
3. Update stakeholders
4. Plan primary site recovery

**Checklist:**
- [ ] DR site activated
- [ ] Backups restored
- [ ] Services online
- [ ] Data synchronized
- [ ] Full operation verified
- [ ] Stakeholders updated

### 2.5 Scenario: Quantum Computer Attack

**Cause:** Quantum computer breaks Ed25519

**Impact:**
- All signatures potentially forgeable
- All public keys compromised
- System trust destroyed

**Recovery Priority:** EXISTENTIAL
**RTO:** N/A (requires protocol upgrade)
**RPO:** N/A

**Recovery Procedure:**

**Step 1: Activation (immediate)**
1. Activate quantum-safe algorithm (Dilithium)
2. Emergency protocol upgrade
3. All trustees participate

**Step 2: Migration (coordinated)**
1. Generate new quantum-safe keys
2. Dual-sign all new content
3. Maintain Ed25519 for verification only
4. Gradual transition period

**Step 3: Deprecation (planned)**
1. Full migration to quantum-safe
2. Archive Ed25519 signatures
3. Update all verifiers
4. Complete transition

**Note:** This scenario has preparatory measures:
- Dilithium implementation ready (feature flag)
- Migration plan documented
- Regular quantum threat assessment
- Early warning monitoring

## 3. Recovery Resources

### 3.1 Backup HSMs

**Primary Backup:** Secure vault, same facility
**Secondary Backup:** Offsite secure facility
**Tertiary Backup:** Geographic diversity

**Access:** Requires 3 of 5 trustee quorum

### 3.2 Data Backups

**Ledger:**
- Real-time replication to 5+ mirror nodes
- Hourly snapshots
- Daily archives
- Geographic distribution

**Identity Registry:**
- Real-time replication
- Hourly backups
- Encrypted at rest

**System Configuration:**
- Version controlled
- Daily backups
- Infrastructure as code

### 3.3 Communication Channels

**Primary:** Encrypted email list
**Secondary:** Secure messaging (Signal)
**Tertiary:** Phone tree
**Emergency:** Physical assembly

### 3.4 Contact Lists

**Trustees:** [CONFIDENTIAL]
**HSM Vendor:** [CONFIDENTIAL]
**Data Center:** [CONFIDENTIAL]
**Legal:** [CONFIDENTIAL]
**PR:** [CONFIDENTIAL]

## 4. Testing and Drills

### 4.1 Drill Schedule

**Monthly:** Communication drill
**Quarterly:** HSM failover drill
**Semi-Annual:** Ledger recovery drill
**Annual:** Full disaster recovery drill

### 4.2 Drill Documentation

Each drill must document:
- Date and participants
- Scenario tested
- Results and timing
- Issues identified
- Improvements needed
- Action items

## 5. Continuous Improvement

### 5.1 Post-Incident Review

After any disaster:
1. Complete incident timeline
2. Identify root cause
3. Evaluate response
4. Document lessons learned
5. Update procedures
6. Share knowledge

### 5.2 Plan Updates

This plan reviewed:
- After each incident
- After each drill
- Quarterly by trustees
- Annually comprehensive review

## 6. Appendices

### Appendix A: Emergency Contact Cards

Distributed to all trustees - contains:
- Emergency phone tree
- Facility addresses
- Vault combinations (split)
- Initial response steps

### Appendix B: Recovery Checklists

Detailed checklists for each scenario available separately.

### Appendix C: Vendor Contacts

Complete list of vendors and support contacts.

---

**Distribution:** Trustees only
**Classification:** CONFIDENTIAL
**Review Date:** Quarterly
