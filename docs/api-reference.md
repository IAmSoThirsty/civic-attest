# API Reference

**Version:** 2.0
**Status:** Production
**Last Updated:** 2026-02-23

## Overview

Civic Attest provides REST and gRPC APIs for signing, verification, ledger operations, and witness coordination. This document provides complete API reference documentation.

## 1. REST API

### 1.1 Base URL

```
Production: https://api.civic-attest.example.com/v2
Staging: https://api-staging.civic-attest.example.com/v2
```

### 1.2 Authentication

**API Key Authentication:**

```http
POST /api/v2/sign HTTP/1.1
Host: api.civic-attest.example.com
Authorization: Bearer {API_KEY}
Content-Type: application/json
```

**OAuth 2.0 (for integrations):**

```http
POST /oauth/token HTTP/1.1
Host: api.civic-attest.example.com
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&
client_id={CLIENT_ID}&
client_secret={CLIENT_SECRET}
```

### 1.3 Signing Operations

#### Sign Content

**Endpoint:** `POST /api/v2/sign`

**Description:** Create a cryptographic signature for content

**Request:**

```http
POST /api/v2/sign HTTP/1.1
Authorization: Bearer {API_KEY}
Content-Type: application/json

{
  "content": "base64_encoded_content",
  "content_type": "text/plain",
  "identity_id": "mayor-springfield-v1",
  "metadata": {
    "title": "Official Statement",
    "date": "2026-02-23"
  }
}
```

**Response (200 OK):**

```json
{
  "signature_bundle": {
    "bundle_version": 2,
    "content_hash": "a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "content_hash_algorithm": "SHA-256",
    "canonical_format_version": "2.0",
    "signer_identity_id": "mayor-springfield-v1",
    "key_version": 1,
    "signatures": {
      "classical": {
        "algorithm": "Ed25519",
        "signature": "7e4c9d8a3b2c1d0e9f8a7b6c5d4e3f2a",
        "pubkey": "b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0"
      }
    },
    "timestamps": [
      {
        "tsa_id": "tsa-provider-1",
        "timestamp_token": "base64_encoded_token",
        "signed_time": "2026-02-23T12:00:00Z"
      }
    ],
    "ledger_entry_hash": "c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8",
    "merkle_inclusion_proof": {
      "leaf_index": 12345,
      "tree_size": 100000,
      "path": ["hash1", "hash2", "hash3"]
    },
    "signed_tree_head_reference": {
      "tree_size": 100000,
      "root_hash": "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9",
      "witness_quorum_met": true
    }
  },
  "operation_id": "op-2026-02-23-12345",
  "created_at": "2026-02-23T12:00:00Z",
  "processing_time_ms": 125
}
```

**Error Responses:**

```json
// 400 Bad Request - Invalid input
{
  "error": "invalid_request",
  "message": "Content exceeds maximum size of 10MB",
  "details": {
    "field": "content",
    "max_size_bytes": 10485760
  }
}

// 401 Unauthorized - Invalid or missing API key
{
  "error": "unauthorized",
  "message": "Invalid API key"
}

// 429 Too Many Requests - Rate limit exceeded
{
  "error": "rate_limit_exceeded",
  "message": "Rate limit of 1000 requests/minute exceeded",
  "retry_after": 30
}

// 503 Service Unavailable - HSM unavailable
{
  "error": "service_unavailable",
  "message": "Signing service temporarily unavailable"
}
```

#### Batch Sign

**Endpoint:** `POST /api/v2/sign/batch`

**Description:** Sign multiple content items in a single request

**Request:**

```json
{
  "items": [
    {
      "content": "base64_content_1",
      "identity_id": "mayor-v1",
      "reference_id": "doc-001"
    },
    {
      "content": "base64_content_2",
      "identity_id": "mayor-v1",
      "reference_id": "doc-002"
    }
  ],
  "batch_options": {
    "parallel_processing": true,
    "fail_fast": false
  }
}
```

**Response (200 OK):**

```json
{
  "batch_id": "batch-2026-02-23-001",
  "results": [
    {
      "reference_id": "doc-001",
      "status": "success",
      "signature_bundle": {...}
    },
    {
      "reference_id": "doc-002",
      "status": "success",
      "signature_bundle": {...}
    }
  ],
  "summary": {
    "total": 2,
    "successful": 2,
    "failed": 0
  },
  "processing_time_ms": 245
}
```

### 1.4 Verification Operations

#### Verify Signature

**Endpoint:** `POST /api/v2/verify`

**Description:** Verify a signature bundle

**Request:**

```json
{
  "content": "base64_encoded_content",
  "signature_bundle": {
    "bundle_version": 2,
    "content_hash": "a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "signatures": {...},
    ...
  },
  "verification_options": {
    "check_revocation": true,
    "verify_witness_quorum": true,
    "offline_mode": false
  }
}
```

**Response (200 OK):**

```json
{
  "verification_result": "valid",
  "checks": {
    "content_hash": {
      "status": "passed",
      "expected": "a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6",
      "actual": "a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6"
    },
    "signature": {
      "status": "passed",
      "algorithm": "Ed25519",
      "signer": "mayor-springfield-v1"
    },
    "identity_status": {
      "status": "passed",
      "identity_valid": true,
      "not_revoked": true
    },
    "timestamp": {
      "status": "passed",
      "signed_at": "2026-02-23T12:00:00Z",
      "tsa_count": 3,
      "quorum_met": true
    },
    "ledger_inclusion": {
      "status": "passed",
      "entry_index": 12345,
      "proof_valid": true
    },
    "witness_quorum": {
      "status": "passed",
      "witnesses_signed": 5,
      "quorum_required": 3,
      "quorum_met": true
    }
  },
  "signer_info": {
    "identity_id": "mayor-springfield-v1",
    "jurisdiction": "Springfield",
    "office": "Mayor",
    "valid_from": "2026-01-01T00:00:00Z",
    "valid_to": "2030-01-01T00:00:00Z"
  },
  "verified_at": "2026-02-23T13:00:00Z",
  "processing_time_ms": 45
}
```

**Invalid Signature Response:**

```json
{
  "verification_result": "invalid",
  "checks": {
    "signature": {
      "status": "failed",
      "reason": "Signature does not match content hash"
    }
  },
  "verified_at": "2026-02-23T13:00:00Z"
}
```

### 1.5 Ledger Operations

#### Get Signed Tree Head

**Endpoint:** `GET /api/v2/ledger/tree-head`

**Response (200 OK):**

```json
{
  "sth_version": 2,
  "tree_size": 100000,
  "root_hash": "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9",
  "identity_tree_root": "e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0",
  "revocation_tree_root": "f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1",
  "timestamp": "2026-02-23T12:00:00Z",
  "ledger_authority_id": "ledger-authority-v1",
  "ledger_authority_signature": "a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6",
  "witness_signatures": [
    {
      "witness_id": "witness-org-1",
      "signature": "b2c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7",
      "signed_at": "2026-02-23T12:00:01Z"
    },
    {
      "witness_id": "witness-org-2",
      "signature": "c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8",
      "signed_at": "2026-02-23T12:00:02Z"
    }
  ],
  "witness_quorum": "3-of-5",
  "quorum_met": true
}
```

#### Get Entry

**Endpoint:** `GET /api/v2/ledger/entries/{index}`

**Response (200 OK):**

```json
{
  "entry_hash": "c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8",
  "sequence_number": 12345,
  "timestamp": "2026-02-23T12:00:00Z",
  "signer_identity_id": "mayor-springfield-v1",
  "signature_hash": "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9",
  "entry_type": "signature",
  "entry_data": {
    "content_hash": "a3f2b1c4d5e6f7a8b9c0d1e2f3a4b5c6",
    "metadata": {...}
  }
}
```

#### Get Inclusion Proof

**Endpoint:** `GET /api/v2/ledger/inclusion-proof/{index}`

**Query Parameters:**
- `tree_size` (optional): Tree size for proof (defaults to current)

**Response (200 OK):**

```json
{
  "leaf_index": 12345,
  "leaf_hash": "c3d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8",
  "tree_size": 100000,
  "path": [
    "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9",
    "e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0",
    "f6a7b8c9d0e1f2a3b4c5d6e7f8a9b0c1"
  ],
  "tree_head": {
    "tree_size": 100000,
    "root_hash": "d4e5f6a7b8c9d0e1f2a3b4c5d6e7f8a9"
  }
}
```

#### Get Consistency Proof

**Endpoint:** `GET /api/v2/ledger/consistency-proof`

**Query Parameters:**
- `old_size` (required): Previous tree size
- `new_size` (required): Current tree size

**Response (200 OK):**

```json
{
  "old_tree_size": 90000,
  "new_tree_size": 100000,
  "consistency_path": [
    "hash1",
    "hash2",
    "hash3"
  ],
  "old_root_hash": "old_root_hash_value",
  "new_root_hash": "new_root_hash_value"
}
```

### 1.6 Identity Operations

#### Get Identity

**Endpoint:** `GET /api/v2/identities/{identity_id}`

**Response (200 OK):**

```json
{
  "identity_id": "mayor-springfield-v1",
  "office_id": "mayor",
  "jurisdiction": "Springfield",
  "public_key": "b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0",
  "key_version": 1,
  "key_algorithm": "Ed25519",
  "valid_from": "2026-01-01T00:00:00Z",
  "valid_to": "2030-01-01T00:00:00Z",
  "status": "active",
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

#### List Identities

**Endpoint:** `GET /api/v2/identities`

**Query Parameters:**
- `status` (optional): Filter by status (active, revoked, expired)
- `jurisdiction` (optional): Filter by jurisdiction
- `limit` (optional): Number of results (default: 100, max: 1000)
- `offset` (optional): Pagination offset

**Response (200 OK):**

```json
{
  "identities": [
    {
      "identity_id": "mayor-springfield-v1",
      "office_id": "mayor",
      "jurisdiction": "Springfield",
      "status": "active",
      ...
    }
  ],
  "pagination": {
    "total": 150,
    "limit": 100,
    "offset": 0,
    "has_more": true
  }
}
```

### 1.7 Health and Status

#### Health Check

**Endpoint:** `GET /health`

**Response (200 OK):**

```json
{
  "status": "healthy",
  "version": "2.0.0",
  "timestamp": "2026-02-23T12:00:00Z",
  "components": {
    "database": "healthy",
    "hsm": "healthy",
    "witness_network": "healthy",
    "ledger": "healthy"
  },
  "uptime_seconds": 864000
}
```

#### Readiness Check

**Endpoint:** `GET /ready`

**Response (200 OK if ready, 503 if not):**

```json
{
  "ready": true,
  "checks": {
    "database_connected": true,
    "hsm_accessible": true,
    "minimum_witnesses_available": true
  }
}
```

#### Metrics

**Endpoint:** `GET /metrics`

**Response (200 OK):**

```
# Prometheus format metrics
civic_attest_signatures_total{identity="mayor-v1"} 12345
civic_attest_verifications_total 67890
civic_attest_signature_duration_seconds{quantile="0.5"} 0.05
civic_attest_signature_duration_seconds{quantile="0.95"} 0.10
civic_attest_ledger_size 100000
civic_attest_witness_quorum_success_rate 0.99
```

## 2. Rate Limiting

### 2.1 Rate Limit Headers

All API responses include rate limit information:

```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 750
X-RateLimit-Reset: 1677158400
Retry-After: 30
```

### 2.2 Rate Limit Tiers

| Tier | Signatures/min | Verifications/min | Ledger Queries/min |
|------|----------------|-------------------|-------------------|
| Standard | 100 | 1000 | 500 |
| Professional | 500 | 5000 | 2000 |
| Enterprise | 1000 | 10000 | 5000 |

## 3. Error Codes

| HTTP Status | Error Code | Description |
|------------|------------|-------------|
| 400 | invalid_request | Malformed request |
| 401 | unauthorized | Missing or invalid authentication |
| 403 | forbidden | Insufficient permissions |
| 404 | not_found | Resource not found |
| 409 | conflict | Resource conflict |
| 422 | validation_failed | Request validation failed |
| 429 | rate_limit_exceeded | Rate limit exceeded |
| 500 | internal_error | Server error |
| 503 | service_unavailable | Service temporarily unavailable |

## 4. Webhooks

### 4.1 Webhook Events

**Event Types:**
- `signature.created` - New signature created
- `signature.verified` - Signature verification completed
- `identity.created` - New identity issued
- `identity.revoked` - Identity revoked
- `tree_head.updated` - New signed tree head
- `witness.quorum_achieved` - Witness quorum reached

**Webhook Payload:**

```json
{
  "event_id": "evt_2026-02-23-12345",
  "event_type": "signature.created",
  "timestamp": "2026-02-23T12:00:00Z",
  "data": {
    "signature_bundle": {...},
    "operation_id": "op-2026-02-23-12345"
  },
  "signature": "webhook_signature_for_verification"
}
```

### 4.2 Webhook Registration

**Endpoint:** `POST /api/v2/webhooks`

**Request:**

```json
{
  "url": "https://your-server.com/webhooks",
  "events": ["signature.created", "identity.revoked"],
  "secret": "your_webhook_secret"
}
```

## 5. SDKs and Client Libraries

### 5.1 Official SDKs

**Go:**
```go
import "github.com/civic-attest/go-sdk"

client := civicattest.NewClient("API_KEY")
result, err := client.Sign(ctx, &civicattest.SignRequest{
    Content:    []byte("message"),
    IdentityID: "mayor-v1",
})
```

**Python:**
```python
from civic_attest import Client

client = Client(api_key="API_KEY")
result = client.sign(
    content=b"message",
    identity_id="mayor-v1"
)
```

**JavaScript/TypeScript:**
```typescript
import { CivicAttestClient } from '@civic-attest/sdk';

const client = new CivicAttestClient({ apiKey: 'API_KEY' });
const result = await client.sign({
  content: Buffer.from('message'),
  identityId: 'mayor-v1'
});
```

**Java:**
```java
import com.civicattest.Client;

Client client = new Client("API_KEY");
SignResult result = client.sign(
    SignRequest.builder()
        .content("message".getBytes())
        .identityId("mayor-v1")
        .build()
);
```

## 6. Best Practices

### 6.1 API Integration

1. **Use exponential backoff** for retries
2. **Implement idempotency** using operation IDs
3. **Cache verification results** to reduce API calls
4. **Monitor rate limits** and adjust accordingly
5. **Validate webhook signatures** before processing
6. **Use batch operations** for high-volume scenarios
7. **Implement circuit breakers** for fault tolerance

### 6.2 Security

1. **Never expose API keys** in client-side code
2. **Rotate API keys** regularly (quarterly)
3. **Use TLS 1.3** for all connections
4. **Validate all inputs** before sending to API
5. **Log security events** for audit
6. **Implement request signing** for critical operations

## Appendix A: OpenAPI Specification

See `api/openapi.yaml` for complete OpenAPI 3.0 specification.

## Appendix B: Postman Collection

See `api/civic-attest.postman_collection.json` for Postman collection.

## Appendix C: Code Examples

See `examples/` directory for complete integration examples.

---

**API Version:** 2.0
**Documentation Version:** 1.0
**Last Updated:** 2026-02-23
