# Operational Tuning Guide

**Version:** 2.0
**Status:** Production Operations
**Last Updated:** 2026-02-23

## Overview

This guide provides operational tuning recommendations for deploying and operating Civic Attest in production environments. These configurations balance security, performance, and resource utilization for different deployment scales.

## 1. Configuration Profiles

### 1.1 Small Deployment (< 1,000 signatures/day)

**Target Environment:**
- Small municipalities
- Local government offices
- Pilot deployments

**Resource Requirements:**
- HSM: Single FIPS 140-2 Level 3 device
- Ledger Nodes: 3 nodes (1 primary, 2 mirrors)
- Witnesses: 3 independent witnesses
- CPU: 4 cores per node
- Memory: 8 GB per node
- Storage: 100 GB SSD

**Configuration:**
```yaml
# config/small.yaml
ledger:
  max_entries_per_batch: 100
  batch_timeout_ms: 5000
  max_tree_size: 1000000
  snapshot_interval: 1000

hsm:
  rate_limit_per_minute: 100
  burst_limit: 20
  session_timeout_hours: 8

witness:
  cosign_timeout_seconds: 30
  min_witness_quorum: 2
  max_witness_count: 3

performance:
  worker_threads: 4
  io_threads: 2
  cache_size_mb: 512
  connection_pool_size: 10
```

### 1.2 Medium Deployment (1,000 - 10,000 signatures/day)

**Target Environment:**
- State/provincial agencies
- Regional governments
- Multi-office deployments

**Resource Requirements:**
- HSM: 2 HSMs (active-backup or threshold)
- Ledger Nodes: 5 nodes (1 primary, 4 mirrors)
- Witnesses: 5 independent witnesses
- CPU: 8 cores per node
- Memory: 16 GB per node
- Storage: 500 GB SSD

**Configuration:**
```yaml
# config/medium.yaml
ledger:
  max_entries_per_batch: 500
  batch_timeout_ms: 2000
  max_tree_size: 10000000
  snapshot_interval: 5000

hsm:
  rate_limit_per_minute: 500
  burst_limit: 100
  session_timeout_hours: 8
  threshold_signing_enabled: true
  threshold_config: "2-of-2"

witness:
  cosign_timeout_seconds: 20
  min_witness_quorum: 3
  max_witness_count: 5

performance:
  worker_threads: 8
  io_threads: 4
  cache_size_mb: 2048
  connection_pool_size: 50
```

### 1.3 Large Deployment (> 10,000 signatures/day)

**Target Environment:**
- National governments
- Federal agencies
- High-volume institutional use

**Resource Requirements:**
- HSM: 3 HSMs (threshold 2-of-3 configuration)
- Ledger Nodes: 7+ nodes (1 primary, 6+ mirrors)
- Witnesses: 7 independent witnesses
- CPU: 16 cores per node
- Memory: 32 GB per node
- Storage: 2 TB NVMe SSD

**Configuration:**
```yaml
# config/large.yaml
ledger:
  max_entries_per_batch: 1000
  batch_timeout_ms: 1000
  max_tree_size: 100000000
  snapshot_interval: 10000
  parallel_proof_generation: true

hsm:
  rate_limit_per_minute: 1000
  burst_limit: 200
  session_timeout_hours: 8
  threshold_signing_enabled: true
  threshold_config: "2-of-3"
  geographic_distribution: true

witness:
  cosign_timeout_seconds: 15
  min_witness_quorum: 4
  max_witness_count: 7
  parallel_verification: true

performance:
  worker_threads: 16
  io_threads: 8
  cache_size_mb: 8192
  connection_pool_size: 200
  enable_cpu_affinity: true
```

## 2. Performance Tuning

### 2.1 Ledger Node Optimization

**Merkle Tree Caching:**
```yaml
merkle_cache:
  enabled: true
  max_cache_entries: 100000
  cache_ttl_seconds: 3600
  eviction_policy: "lru"
```

**Database Tuning (for ledger storage):**
```yaml
database:
  # For PostgreSQL
  max_connections: 100
  shared_buffers: "4GB"
  effective_cache_size: "12GB"
  work_mem: "64MB"
  maintenance_work_mem: "512MB"
  checkpoint_completion_target: 0.9
  wal_buffers: "16MB"
  default_statistics_target: 100
  random_page_cost: 1.1  # For SSD
```

**Network Optimization:**
```yaml
network:
  tcp_keepalive: true
  tcp_nodelay: true
  max_concurrent_connections: 1000
  read_timeout_seconds: 30
  write_timeout_seconds: 30
  idle_timeout_seconds: 300
```

### 2.2 HSM Optimization

**Connection Pooling:**
```yaml
hsm_pool:
  min_connections: 2
  max_connections: 10
  connection_timeout_seconds: 30
  health_check_interval_seconds: 60
```

**Batch Signing (where supported):**
```yaml
batch_signing:
  enabled: true
  max_batch_size: 100
  batch_wait_ms: 50
```

### 2.3 Witness Network Optimization

**Parallel Cosigning:**
```yaml
witness_optimization:
  parallel_requests: true
  request_timeout_seconds: 15
  retry_attempts: 3
  retry_backoff_ms: 1000
  circuit_breaker_enabled: true
  circuit_breaker_threshold: 5
```

## 3. Resource Planning

### 3.1 Storage Planning

**Growth Estimates:**

| Deployment Size | Daily Entries | Monthly Growth | Annual Growth |
|----------------|---------------|----------------|---------------|
| Small | 1,000 | ~300 MB | ~3.6 GB |
| Medium | 10,000 | ~3 GB | ~36 GB |
| Large | 100,000 | ~30 GB | ~360 GB |

**Storage Allocation:**
- Ledger data: 70% of total
- Snapshots: 20% of total
- Logs: 10% of total

**Retention Policy:**
```yaml
retention:
  ledger_entries: "permanent"
  snapshots:
    hourly: "7_days"
    daily: "90_days"
    weekly: "1_year"
    monthly: "10_years"
  logs:
    application: "90_days"
    audit: "7_years"
    security: "10_years"
```

### 3.2 Bandwidth Planning

**Estimated Bandwidth (per day):**

| Component | Small | Medium | Large |
|-----------|-------|--------|-------|
| Signature Traffic | 100 MB | 1 GB | 10 GB |
| Witness Cosigning | 50 MB | 500 MB | 5 GB |
| Mirror Replication | 100 MB | 1 GB | 10 GB |
| Client Verification | 200 MB | 2 GB | 20 GB |
| **Total** | ~450 MB | ~4.5 GB | ~45 GB |

**Network Requirements:**
- Small: 10 Mbps sustained, 100 Mbps burst
- Medium: 100 Mbps sustained, 1 Gbps burst
- Large: 1 Gbps sustained, 10 Gbps burst

### 3.3 CPU Planning

**CPU Allocation by Component:**

| Component | % of Total CPU |
|-----------|---------------|
| Signature Verification | 30% |
| Merkle Tree Operations | 25% |
| Witness Communication | 20% |
| Database Operations | 15% |
| API/Networking | 10% |

**Recommended CPU:**
- Small: 4-8 cores @ 2.5+ GHz
- Medium: 8-16 cores @ 3.0+ GHz
- Large: 16-32 cores @ 3.5+ GHz

## 4. High Availability Configuration

### 4.1 Ledger Node HA

**Primary-Mirror Architecture:**
```yaml
ha_config:
  replication_mode: "synchronous"
  min_sync_replicas: 2
  failover_timeout_seconds: 30
  automatic_failover: true
  health_check_interval_seconds: 10
```

**Load Balancing:**
```yaml
load_balancer:
  algorithm: "least_connections"
  health_check:
    endpoint: "/health"
    interval_seconds: 5
    timeout_seconds: 2
    unhealthy_threshold: 3
    healthy_threshold: 2
  session_affinity: false
```

### 4.2 Geographic Distribution

**Multi-Region Deployment:**
```yaml
regions:
  - name: "primary"
    location: "us-east-1"
    role: "active"
    witnesses: 3

  - name: "secondary"
    location: "eu-west-1"
    role: "standby"
    witnesses: 2

  - name: "tertiary"
    location: "ap-south-1"
    role: "standby"
    witnesses: 2

replication:
  cross_region_enabled: true
  replication_lag_threshold_seconds: 60
  conflict_resolution: "primary_wins"
```

## 5. Monitoring Thresholds

### 5.1 Performance Metrics

**Key Metrics to Monitor:**

```yaml
alerts:
  # Latency
  signature_latency_p95_ms:
    warning: 100
    critical: 500

  verification_latency_p95_ms:
    warning: 50
    critical: 200

  ledger_append_latency_ms:
    warning: 100
    critical: 1000

  # Throughput
  signatures_per_minute:
    min_warning: 10  # for expected load
    max_warning: 900  # approaching limit

  # Resource Utilization
  cpu_usage_percent:
    warning: 70
    critical: 85

  memory_usage_percent:
    warning: 75
    critical: 90

  disk_usage_percent:
    warning: 70
    critical: 85

  # HSM
  hsm_queue_depth:
    warning: 50
    critical: 100

  hsm_error_rate_percent:
    warning: 1
    critical: 5

  # Witness
  witness_cosign_success_rate_percent:
    warning: 95
    critical: 90

  witness_timeout_rate_percent:
    warning: 5
    critical: 10
```

### 5.2 Health Checks

**Component Health Checks:**

```yaml
health_checks:
  ledger_node:
    - check: "database_connection"
      timeout_seconds: 5
      critical: true

    - check: "merkle_tree_integrity"
      interval_seconds: 300
      critical: true

    - check: "disk_space"
      threshold_percent: 80
      critical: true

  hsm:
    - check: "hsm_connection"
      timeout_seconds: 10
      critical: true

    - check: "key_accessibility"
      interval_seconds: 60
      critical: true

  witness:
    - check: "witness_connectivity"
      min_witnesses: 3
      critical: true

    - check: "consensus_status"
      interval_seconds: 30
      critical: false
```

## 6. Backup Configuration

### 6.1 Backup Strategy

**Tiered Backup Approach:**

```yaml
backups:
  # Hot backups (continuous)
  continuous:
    enabled: true
    destination: "replicas"
    min_replicas: 3
    verification: "automatic"

  # Warm backups (hourly snapshots)
  snapshots:
    hourly:
      enabled: true
      retention_hours: 168  # 7 days
      compression: true
      encryption: true

    daily:
      enabled: true
      retention_days: 90
      compression: true
      encryption: true
      offsite_copy: true

  # Cold backups (archival)
  archival:
    weekly:
      enabled: true
      retention_years: 10
      media: "tape"
      geographic_separation: true
      verification_frequency: "quarterly"
```

### 6.2 Recovery Time Objectives

**RTO/RPO Targets:**

| Failure Scenario | RTO | RPO |
|-----------------|-----|-----|
| Single node failure | 1 minute | 0 (no data loss) |
| Primary datacenter failure | 1 hour | 5 minutes |
| Region-wide failure | 4 hours | 1 hour |
| Complete system rebuild | 24 hours | 24 hours (from backup) |

## 7. Security Hardening

### 7.1 Network Security

**Firewall Rules:**
```yaml
firewall:
  inbound:
    # API endpoints
    - port: 8080
      protocol: "tcp"
      source: "trusted_networks"
      description: "REST API"

    # Witness network
    - port: 8443
      protocol: "tcp"
      source: "witness_networks"
      description: "Witness cosigning"

  outbound:
    # Timestamp authorities
    - port: 443
      protocol: "tcp"
      destination: "tsa_endpoints"
      description: "TSA queries"

  default_policy: "deny"
```

### 7.2 Access Control

**RBAC Configuration:**
```yaml
roles:
  operator:
    permissions:
      - "sign:read"
      - "ledger:append"
      - "logs:read"

  administrator:
    permissions:
      - "config:write"
      - "keys:manage"
      - "users:manage"

  auditor:
    permissions:
      - "logs:read"
      - "ledger:read"
      - "reports:generate"

  trustee:
    permissions:
      - "governance:vote"
      - "keys:ceremony"
      - "emergency:override"
```

## 8. Operational Procedures

### 8.1 Startup Sequence

1. **Pre-flight Checks:**
   - Verify HSM connectivity
   - Check disk space (>20% free)
   - Validate configuration files
   - Verify network connectivity to witnesses

2. **Start Order:**
   - Database services
   - Ledger node (read-only mode)
   - Verify ledger integrity
   - HSM connection pool
   - Enable write mode
   - API endpoints
   - Witness connectivity

3. **Post-startup Validation:**
   - Sign test message
   - Verify witness cosigning
   - Check all health endpoints
   - Monitor logs for errors

### 8.2 Shutdown Sequence

1. **Graceful Shutdown:**
   - Stop accepting new requests
   - Complete in-flight operations (timeout: 60s)
   - Flush pending ledger entries
   - Create shutdown snapshot
   - Close HSM connections
   - Stop database connections
   - Archive logs

2. **Emergency Shutdown:**
   - Immediate freeze (any 2 trustees)
   - Log emergency state
   - Close HSM (zeroize if tamper detected)
   - Create emergency snapshot
   - Notify all stakeholders

## 9. Tuning Checklist

**Pre-Production Tuning:**

- [ ] Select appropriate configuration profile (small/medium/large)
- [ ] Provision resources per requirements
- [ ] Configure high availability if required
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy
- [ ] Test failover scenarios
- [ ] Establish baseline performance metrics
- [ ] Document operational procedures
- [ ] Train operations team
- [ ] Conduct disaster recovery drill

**Regular Tuning (Monthly):**

- [ ] Review performance metrics
- [ ] Analyze resource utilization trends
- [ ] Adjust rate limits if needed
- [ ] Review and rotate logs
- [ ] Verify backup integrity
- [ ] Update capacity planning
- [ ] Review security posture

**Quarterly Tuning:**

- [ ] Full system performance audit
- [ ] Review and update configuration
- [ ] Capacity planning review
- [ ] Disaster recovery drill
- [ ] Security audit
- [ ] Update operational documentation

## 10. Performance Targets

### 10.1 Latency Targets

| Operation | P50 | P95 | P99 |
|-----------|-----|-----|-----|
| Signature (local) | <50ms | <100ms | <200ms |
| Signature (with ledger) | <100ms | <200ms | <500ms |
| Verification (offline) | <20ms | <50ms | <100ms |
| Verification (with ledger) | <50ms | <100ms | <200ms |
| Ledger append | <100ms | <200ms | <500ms |
| Witness cosigning | <1s | <2s | <5s |

### 10.2 Throughput Targets

| Component | Target (sustained) | Peak (burst) |
|-----------|-------------------|--------------|
| Signature operations | 1,000/min | 2,000/min |
| Verification operations | 10,000/min | 20,000/min |
| Ledger appends | 1,000/min | 2,000/min |

### 10.3 Availability Targets

| Service Level | Uptime | Downtime/Year |
|--------------|--------|---------------|
| Standard | 99.9% | 8.76 hours |
| High Availability | 99.95% | 4.38 hours |
| Critical | 99.99% | 52.56 minutes |

## Appendix A: Configuration Templates

See `config/` directory for complete configuration templates:
- `config/small.yaml` - Small deployment
- `config/medium.yaml` - Medium deployment
- `config/large.yaml` - Large deployment
- `config/ha.yaml` - High availability
- `config/geo-distributed.yaml` - Geographic distribution

## Appendix B: Monitoring Dashboards

See `monitoring/` directory for dashboard templates:
- `monitoring/grafana/` - Grafana dashboards
- `monitoring/prometheus/` - Prometheus scrape configs
- `monitoring/alerts/` - Alert rule definitions

---

**Document Owner:** Operations Team
**Review Frequency:** Quarterly
**Next Review:** 2026-05-23
