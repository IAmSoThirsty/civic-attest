# Key Ceremony Guide

**Version:** 1.0
**Audience:** Trustees, Technical Operators
**Last Updated:** 2026-02-22

## 1. Overview

The key ceremony is a formal, witnessed procedure for generating cryptographic keys in a Hardware Security Module (HSM). This ceremony establishes the root of trust for the Civic Attest system.

## 2. Prerequisites

### 2.1 Required Personnel

- **Trustees:** Minimum 3 of 5 trustees must be present
- **HSM Operator:** Technical operator trained on HSM procedures
- **Auditor:** Independent observer for ceremony verification
- **Recorder:** Person responsible for ceremony documentation

### 2.2 Required Equipment

- FIPS 140-2 Level 3+ certified HSM
- Air-gapped ceremony workstation
- Audio/video recording equipment
- Tamper-evident bags for materials
- Ceremony checklist printouts
- Trustee authentication tokens

### 2.3 Required Documentation

- Ceremony protocol (this document)
- HSM operational procedures
- Emergency procedures
- Trustee authorization list
- Incident report forms

## 3. Ceremony Phases

### Phase 1: Pre-Ceremony Preparation

**Duration:** 1-2 hours before ceremony

**Steps:**

1. **Facility Preparation**
   - Secure ceremony room
   - Deploy recording equipment
   - Test HSM connectivity
   - Verify air-gap isolation
   - Set up audit logging

2. **Equipment Verification**
   - Inspect HSM tamper seals
   - Verify HSM firmware version
   - Test HSM functionality
   - Document serial numbers
   - Photograph equipment state

3. **Personnel Assembly**
   - Verify trustee identities (government ID + biometric)
   - Collect signed attendance sheet
   - Confirm quorum present (3 of 5)
   - Brief all participants
   - Distribute ceremony materials

4. **Documentation**
   - Assign ceremony ID
   - Record date, time, location
   - List all participants
   - Note any exceptions or issues

**Checklist:**
- [ ] Room secured
- [ ] Recording active
- [ ] HSM verified
- [ ] Quorum present
- [ ] All participants briefed
- [ ] Documentation prepared

### Phase 2: Trustee Authorization

**Duration:** 30-45 minutes

**Steps:**

1. **Identity Verification**
   - Each trustee presents:
     - Government-issued photo ID
     - Biometric verification (fingerprint)
     - Authorization token
   - Auditor verifies and records

2. **Authorization Tokens**
   - Each trustee inserts authorization token
   - HSM validates trustee credentials
   - System records authorization
   - Minimum 3 of 5 required

3. **Quorum Establishment**
   - HSM confirms quorum
   - Display quorum status
   - Record to audit log
   - Announce quorum established

**Checklist:**
- [ ] All trustees verified
- [ ] Tokens validated
- [ ] Quorum confirmed
- [ ] Audit log updated

### Phase 3: Key Generation

**Duration:** 15-30 minutes

**Steps:**

1. **HSM Initialization**
   - Verify HSM in key generation mode
   - Confirm export disabled
   - Set key parameters:
     - Algorithm: Ed25519
     - Key usage: Digital signature
     - Exportability: NEVER
   - Display parameters for trustee review

2. **Key Generation Command**
   - Trustee #1 enters generation command
   - Trustee #2 confirms command
   - Trustee #3 authorizes execution
   - HSM generates key pair

3. **Public Key Extraction**
   - HSM extracts public key only
   - Display public key on screen
   - Trustees verify and approve
   - Record public key hash

4. **Verification**
   - Test signature operation
   - Verify signature with public key
   - Confirm private key non-exportable
   - Record test results

**Checklist:**
- [ ] HSM initialized
- [ ] Parameters verified
- [ ] Key generated
- [ ] Public key extracted
- [ ] Test signature verified
- [ ] Export disabled confirmed

### Phase 4: Identity Creation

**Duration:** 15-20 minutes

**Steps:**

1. **Identity Parameters**
   - Office ID: _______________
   - Jurisdiction: _______________
   - Valid from: _______________
   - Valid to: _______________
   - Key version: 1

2. **Identity Object Creation**
   ```json
   {
     "office_id": "...",
     "jurisdiction": "...",
     "public_key": "...",
     "key_version": 1,
     "valid_from": "...",
     "valid_to": "...",
     "key_algorithm": "Ed25519",
     "status": "active",
     "identity_id": "..."
   }
   ```

3. **Trustee Approval**
   - Display identity object
   - Each trustee reviews and approves
   - Record approval signatures
   - Sign with quorum

**Checklist:**
- [ ] Parameters confirmed
- [ ] Identity created
- [ ] Trustees approved
- [ ] Signatures recorded

### Phase 5: Ceremony Recording

**Duration:** 20-30 minutes

**Steps:**

1. **Hash Generation**
   - Compute hash of:
     - Public key
     - Identity object
     - Ceremony parameters
     - Participant list
     - Timestamp
   - Display hash: _______________

2. **Public Broadcast**
   - Prepare public announcement
   - Include ceremony hash
   - Sign with trustee quorum
   - Schedule publication

3. **Ledger Entry**
   - Create ceremony record:
     ```json
     {
       "ceremony_id": "...",
       "timestamp": "...",
       "trustees": [...],
       "quorum_size": 3,
       "total_trustees": 5,
       "recording_hash": "...",
       "public_key_hash": "...",
       "ledger_entry_hash": "..."
     }
     ```
   - Append to ledger
   - Verify inclusion proof

4. **Recording Preservation**
   - Stop audio/video recording
   - Generate recording hash
   - Store in tamper-evident media
   - Distribute copies to trustees
   - Archive master copy

**Checklist:**
- [ ] Ceremony hash computed
- [ ] Public announcement prepared
- [ ] Ledger entry appended
- [ ] Recording preserved
- [ ] Copies distributed

### Phase 6: Post-Ceremony

**Duration:** 30-45 minutes

**Steps:**

1. **Verification**
   - Verify ledger entry
   - Confirm inclusion proof
   - Test signature operation
   - Validate public key

2. **Documentation**
   - Complete ceremony report
   - Collect all signatures
   - Seal physical materials
   - Archive documentation

3. **Public Disclosure**
   - Publish ceremony hash
   - Publish public key
   - Publish identity object
   - Publish ledger entry hash

4. **HSM Securing**
   - Remove trustee tokens
   - Verify HSM sealed
   - Return HSM to vault
   - Update HSM inventory

**Checklist:**
- [ ] Verification complete
- [ ] Documentation sealed
- [ ] Public disclosure made
- [ ] HSM secured

## 4. Emergency Procedures

### 4.1 Ceremony Abort

**Triggers:**
- Security breach
- HSM malfunction
- Trustee unavailability
- Procedural violation

**Procedure:**
1. Announce abort
2. Stop all operations
3. Secure HSM
4. Remove trustee tokens
5. Document reason
6. Schedule new ceremony

### 4.2 HSM Failure

**Procedure:**
1. Stop ceremony
2. Secure failed HSM
3. Activate backup HSM
4. Restart from Phase 1
5. Document incident

### 4.3 Quorum Loss

**Procedure:**
1. Pause ceremony
2. Attempt to restore quorum
3. If unsuccessful, abort
4. Reschedule with full quorum

## 5. Security Considerations

### 5.1 Physical Security

- Ceremony room access controlled
- No electronic devices allowed (except ceremony equipment)
- Video surveillance active
- Tamper-evident seals on all materials

### 5.2 Operational Security

- Air-gapped environment
- No network connectivity during ceremony
- All USB ports disabled (except HSM)
- No removable media allowed

### 5.3 Personnel Security

- Background checks completed
- All participants under NDA
- No single person has complete access
- Dual control for all operations

## 6. Post-Ceremony Responsibilities

### 6.1 Trustees

- Maintain token security
- Attend annual renewal ceremonies
- Participate in emergency procedures
- Review audit reports

### 6.2 HSM Operator

- Maintain HSM security
- Perform regular backups
- Monitor HSM health
- Report anomalies

### 6.3 Auditor

- Verify ceremony compliance
- Review audit logs
- Investigate anomalies
- Produce audit reports

## 7. Ceremony Schedule

### 7.1 Regular Ceremonies

- **Annual Key Rotation:** Scheduled 12 months from issuance
- **New Office:** Within 30 days of office assumption
- **Jurisdiction Change:** As required by governance

### 7.2 Emergency Ceremonies

- **Key Compromise:** Within 24 hours of detection
- **Trustee Change:** Within 7 days of change
- **Security Incident:** As determined by trustees

---

**Appendix A: Ceremony Checklist**

Available in separate document: `ceremony-checklist.pdf`

**Appendix B: HSM Procedures**

Available in separate document: `hsm-procedures.pdf`

**Appendix C: Sample Forms**

- Trustee Authorization Form
- Ceremony Attendance Sheet
- Incident Report Form
- Public Disclosure Template
