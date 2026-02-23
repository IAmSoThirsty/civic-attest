# Performance Scaling Guide

**Version:** 2.0
**Status:** Production Operations
**Last Updated:** 2026-02-23

## Overview

This guide provides strategies for scaling Civic Attest to meet increasing performance demands while maintaining security guarantees and operational stability.

## 1. Scaling Dimensions

### 1.1 Vertical Scaling (Scale-Up)

**When to Scale Up:**
- CPU utilization consistently >70%
- Memory pressure causes swapping
- Single-threaded bottlenecks
- I/O wait times increasing

**Hardware Upgrade Path:**

| Current | Target | Performance Gain |
|---------|--------|------------------|
| 4 cores → 8 cores | 2x throughput | ~1.8x |
| 8 GB → 16 GB RAM | Larger caches | ~1.5x |
| SATA SSD → NVMe | Lower latency | ~3x IOPS |
| 1 Gbps → 10 Gbps NIC | Higher throughput | ~10x |

**Limitations:**
- Diminishing returns above 16 cores (for single instance)
- Memory scaling limited by workload characteristics
- Cost increases non-linearly

### 1.2 Horizontal Scaling (Scale-Out)

**When to Scale Out:**
- Vertical scaling limits reached
- Geographic distribution required
- High availability requirements
- Load distribution needed

**Scaling Strategy:**

```yaml
horizontal_scaling:
  # Read scaling
  read_replicas:
    min: 3
    max: 10
    auto_scaling:
      metric: "cpu_utilization"
      target: 60
      scale_up_threshold: 75
      scale_down_threshold: 40

  # Witness scaling
  witnesses:
    min: 3
    max: 20
    distribution: "geographic"
    auto_scaling: false  # Manual approval required

  # API endpoint scaling
  api_servers:
    min: 2
    max: 20
    auto_scaling:
      metric: "requests_per_second"
      target: 1000
      scale_up_threshold: 1500
      scale_down_threshold: 500
```

### 1.3 Functional Scaling (Component Separation)

**Microservices Architecture:**

```
┌─────────────────┐
│  Load Balancer  │
└────────┬────────┘
         │
    ┌────┴────┬────────────┬──────────┐
    │         │            │          │
┌───▼───┐ ┌──▼──┐  ┌──────▼─────┐ ┌─▼────────┐
│ Sign  │ │Verify│  │ Ledger API │ │ Witness  │
│ Service│ │Service│  │  Service   │ │ Service  │
└───┬───┘ └──┬──┘  └──────┬─────┘ └─┬────────┘
    │        │            │          │
    └────────┴────────────┴──────────┘
                   │
            ┌──────▼──────┐
            │Ledger Storage│
            └─────────────┘
```

**Benefits:**
- Independent scaling per component
- Technology optimization per service
- Failure isolation
- Easier testing and deployment

## 2. Bottleneck Analysis

### 2.1 Identifying Bottlenecks

**Diagnostic Tools:**

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace analysis
go test -trace=trace.out -bench=.
go tool trace trace.out

# Real-time monitoring
pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

**Common Bottlenecks:**

| Component | Symptom | Solution |
|-----------|---------|----------|
| HSM | High queue depth | Add HSM capacity, enable batching |
| Merkle Tree | CPU-bound operations | Add compute, optimize algorithms |
| Database | I/O wait times | Add IOPS, optimize queries |
| Network | High latency | Increase bandwidth, use CDN |
| Witness | Timeout rate increasing | Add witnesses, optimize protocol |

### 2.2 Performance Profiling

**Benchmark Suite:**

```go
// internal/benchmark/signature_bench_test.go
func BenchmarkSignature(b *testing.B) {
    // Setup
    signer := setupSigner()
    data := generateTestData(1024) // 1KB payload

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := signer.Sign(data)
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkLedgerAppend(b *testing.B) {
    ledger := setupLedger()
    entries := generateTestEntries(b.N)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        err := ledger.Append(entries[i])
        if err != nil {
            b.Fatal(err)
        }
    }
}

func BenchmarkWitnessCosigning(b *testing.B) {
    witnesses := setupWitnesses(5)
    treeHead := generateTreeHead()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := collectWitnessSignatures(witnesses, treeHead)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

**Performance Regression Testing:**

```yaml
# .github/workflows/performance.yml
name: Performance Regression

on: [pull_request]

jobs:
  benchmark:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.21'

      - name: Run benchmarks
        run: make benchmark

      - name: Compare with baseline
        run: |
          benchstat baseline.txt current.txt
          # Fail if >10% regression
```

## 3. Optimization Strategies

### 3.1 Cryptographic Optimizations

**Ed25519 Batch Verification:**

```go
// Verify multiple signatures in batch (8x faster than individual)
func BatchVerifySignatures(messages [][]byte, signatures [][]byte, publicKeys [][]byte) (bool, error) {
    // Use edwards25519 batch verification
    // Amortizes expensive elliptic curve operations
    return ed25519.BatchVerify(publicKeys, messages, signatures)
}
```

**Hardware Acceleration:**

```yaml
crypto_acceleration:
  # Use AES-NI for symmetric crypto
  aes_ni_enabled: true

  # Use Intel SHA extensions
  sha_extensions: true

  # Use AVX2 for Merkle tree hashing
  avx2_enabled: true
```

**Signature Caching:**

```go
// Cache verified signatures to avoid re-verification
type SignatureCache struct {
    cache *lru.Cache // LRU cache with TTL
}

func (sc *SignatureCache) Verify(sig, msg, pubkey []byte) (bool, error) {
    cacheKey := hash(sig + msg + pubkey)

    if result, ok := sc.cache.Get(cacheKey); ok {
        return result.(bool), nil
    }

    valid := ed25519.Verify(pubkey, msg, sig)
    sc.cache.Add(cacheKey, valid)
    return valid, nil
}
```

### 3.2 Database Optimizations

**Ledger Storage Schema:**

```sql
-- Optimized ledger table with partitioning
CREATE TABLE ledger_entries (
    sequence_number BIGINT NOT NULL,
    entry_hash BYTEA NOT NULL,
    entry_data JSONB NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (sequence_number, timestamp)
) PARTITION BY RANGE (timestamp);

-- Create monthly partitions
CREATE TABLE ledger_entries_2026_02 PARTITION OF ledger_entries
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

-- Index on frequently queried fields
CREATE INDEX idx_entry_hash ON ledger_entries USING HASH (entry_hash);
CREATE INDEX idx_timestamp ON ledger_entries (timestamp DESC);
```

**Query Optimization:**

```go
// Use prepared statements
stmt, err := db.Prepare("SELECT entry_data FROM ledger_entries WHERE sequence_number = $1")
defer stmt.Close()

// Use connection pooling
db.SetMaxOpenConns(100)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(time.Hour)

// Batch inserts
tx, _ := db.Begin()
for _, entry := range entries {
    tx.Exec("INSERT INTO ledger_entries VALUES ($1, $2, $3)", ...)
}
tx.Commit()
```

### 3.3 Merkle Tree Optimizations

**Incremental Tree Updates:**

```go
type IncrementalMerkleTree struct {
    nodes map[int][]byte  // Sparse storage
    size  uint64
}

// Only recompute affected branches
func (t *IncrementalMerkleTree) Append(leaf []byte) []byte {
    index := t.size
    hash := hashLeaf(leaf)

    // Update only O(log n) nodes
    for level := 0; index > 0; level++ {
        if index%2 == 1 {
            // Right node - compute parent with left sibling
            leftSibling := t.nodes[level*maxNodes + index-1]
            parentHash := hashNodes(leftSibling, hash)
            t.nodes[(level+1)*maxNodes + index/2] = parentHash
            hash = parentHash
        }
        index /= 2
    }

    t.size++
    return hash // new root
}
```

**Parallel Proof Generation:**

```go
func GenerateInclusionProofParallel(tree *MerkleTree, index uint64) ([][]byte, error) {
    // Generate multiple proofs concurrently
    var wg sync.WaitGroup
    proofChan := make(chan []byte, int(math.Log2(float64(tree.Size()))))

    // Parallelize sibling hash collection
    for level := 0; index > 0; level++ {
        wg.Add(1)
        go func(lvl, idx uint64) {
            defer wg.Done()
            sibling := tree.GetSibling(lvl, idx)
            proofChan <- sibling
        }(level, index)
        index /= 2
    }

    go func() {
        wg.Wait()
        close(proofChan)
    }()

    proof := [][]byte{}
    for p := range proofChan {
        proof = append(proof, p)
    }
    return proof, nil
}
```

### 3.4 Network Optimizations

**HTTP/2 with Multiplexing:**

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  120 * time.Second,
    // Enable HTTP/2
    TLSConfig: &tls.Config{
        NextProtos: []string{"h2", "http/1.1"},
    },
}
```

**Connection Pooling:**

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
    // Enable keep-alive
    DisableKeepAlives: false,
}

client := &http.Client{
    Transport: transport,
    Timeout:   30 * time.Second,
}
```

**Compression:**

```go
// Use gzip compression for API responses
func gzipMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        gz := gzip.NewWriter(w)
        defer gz.Close()

        w.Header().Set("Content-Encoding", "gzip")
        gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
        next.ServeHTTP(gzw, r)
    })
}
```

### 3.5 Caching Strategies

**Multi-Level Cache:**

```go
type CacheHierarchy struct {
    l1 *LocalCache     // In-memory (100ms)
    l2 *RedisCache     // Distributed (10ms)
    l3 *DatabaseCache  // Persistent (100ms)
}

func (ch *CacheHierarchy) Get(key string) (interface{}, error) {
    // Try L1 (fastest)
    if val, ok := ch.l1.Get(key); ok {
        return val, nil
    }

    // Try L2
    if val, err := ch.l2.Get(key); err == nil {
        ch.l1.Set(key, val) // Promote to L1
        return val, nil
    }

    // Fall back to L3
    val, err := ch.l3.Get(key)
    if err == nil {
        ch.l2.Set(key, val) // Promote to L2
        ch.l1.Set(key, val) // Promote to L1
    }
    return val, err
}
```

**Cache Invalidation:**

```yaml
cache_policy:
  # Time-based expiration
  signature_bundles_ttl: 3600  # 1 hour
  tree_heads_ttl: 300           # 5 minutes
  identity_records_ttl: 86400   # 24 hours

  # Event-based invalidation
  invalidate_on:
    - "key_rotation"
    - "identity_revocation"
    - "tree_head_update"
```

## 4. Load Testing

### 4.1 Load Test Scenarios

**Baseline Load Test:**

```bash
# Using hey (HTTP load generator)
hey -n 10000 -c 100 -m POST \
    -H "Content-Type: application/json" \
    -d '{"data":"test"}' \
    http://localhost:8080/api/sign

# Using k6 (modern load testing)
k6 run --vus 100 --duration 30s load-test.js
```

**Load Test Script (k6):**

```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '2m', target: 100 },  // Ramp up to 100 users
        { duration: '5m', target: 100 },  // Stay at 100 users
        { duration: '2m', target: 200 },  // Ramp up to 200 users
        { duration: '5m', target: 200 },  // Stay at 200 users
        { duration: '2m', target: 0 },    // Ramp down to 0 users
    ],
    thresholds: {
        'http_req_duration': ['p(95)<200'],  // 95% of requests < 200ms
        'http_req_failed': ['rate<0.01'],    // Error rate < 1%
    },
};

export default function() {
    let payload = JSON.stringify({
        content: 'test message',
        identity: 'test-identity-1',
    });

    let params = {
        headers: { 'Content-Type': 'application/json' },
    };

    let res = http.post('http://localhost:8080/api/sign', payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
    });

    sleep(1);
}
```

### 4.2 Stress Testing

**Breaking Point Test:**

```javascript
// stress-test.js
export let options = {
    stages: [
        { duration: '2m', target: 100 },
        { duration: '5m', target: 200 },
        { duration: '5m', target: 400 },
        { duration: '5m', target: 800 },
        { duration: '5m', target: 1600 }, // Push to breaking point
        { duration: '2m', target: 0 },
    ],
};

// Same test function as above
```

**Chaos Testing:**

```yaml
# chaos-test.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: ledger-node-failure
spec:
  action: pod-kill
  mode: one
  selector:
    labelSelectors:
      app: ledger-node
  scheduler:
    cron: '@every 10m'
```

### 4.3 Soak Testing

**Long-Duration Test:**

```javascript
// soak-test.js
export let options = {
    stages: [
        { duration: '5m', target: 50 },   // Ramp up
        { duration: '24h', target: 50 },  // Sustain for 24 hours
        { duration: '5m', target: 0 },    // Ramp down
    ],
};

// Monitor for memory leaks, resource exhaustion
```

## 5. Scalability Benchmarks

### 5.1 Baseline Benchmarks

**Single Node Performance:**

| Operation | Throughput | Latency (p95) |
|-----------|-----------|---------------|
| Sign (in-memory) | 10,000/s | 5ms |
| Sign (with ledger) | 1,000/s | 100ms |
| Verify (offline) | 50,000/s | 1ms |
| Verify (with ledger) | 10,000/s | 50ms |
| Ledger append | 1,000/s | 100ms |

**Scaling Benchmarks:**

| Nodes | Signatures/s | Verifications/s | Latency (p95) |
|-------|-------------|-----------------|---------------|
| 1 | 1,000 | 10,000 | 100ms |
| 3 | 2,500 | 25,000 | 120ms |
| 5 | 4,000 | 40,000 | 150ms |
| 10 | 7,000 | 70,000 | 200ms |

### 5.2 Witness Scaling

**Cosigning Performance:**

| Witnesses | Cosign Time (p95) | Throughput Impact |
|-----------|-------------------|-------------------|
| 3 | 1s | Baseline |
| 5 | 1.5s | -10% |
| 7 | 2s | -20% |
| 10 | 3s | -30% |

**Optimization:** Parallel witness requests reduce impact to <5%

## 6. Cost Optimization

### 6.1 Resource Efficiency

**Right-Sizing:**

```yaml
cost_optimization:
  # Use auto-scaling to match demand
  auto_scaling:
    enabled: true
    min_instances: 2
    max_instances: 10
    target_cpu: 60

  # Use spot instances for non-critical components
  spot_instances:
    verification_workers: true
    batch_processors: true
    read_replicas: true

  # Storage tiering
  storage:
    hot: "nvme-ssd"     # Last 30 days
    warm: "sata-ssd"    # 30-90 days
    cold: "hdd"         # >90 days
    archive: "tape"     # >1 year
```

### 6.2 Cost Monitoring

**Cost Metrics:**

| Resource | Cost/Month (Medium) | Optimization |
|----------|-------------------|--------------|
| Compute | $2,000 | Auto-scaling, spot instances |
| Storage | $500 | Compression, tiering |
| Network | $300 | CDN, compression |
| HSM | $5,000 | Shared infrastructure |
| **Total** | **$7,800** | |

## 7. Scaling Checklist

**Pre-Scaling:**

- [ ] Identify bottleneck (CPU, memory, I/O, network)
- [ ] Establish current baseline performance
- [ ] Review monitoring dashboards
- [ ] Check resource utilization trends
- [ ] Review application logs for errors
- [ ] Run load tests to confirm capacity limits

**Scaling Execution:**

- [ ] Select scaling strategy (vertical, horizontal, functional)
- [ ] Update capacity plan
- [ ] Provision new resources
- [ ] Update configuration
- [ ] Deploy changes (canary or blue-green)
- [ ] Run validation tests
- [ ] Monitor for regressions

**Post-Scaling:**

- [ ] Verify performance improvements
- [ ] Update documentation
- [ ] Adjust monitoring thresholds
- [ ] Review cost impact
- [ ] Update capacity planning
- [ ] Document lessons learned

## Appendix A: Benchmark Results

See `benchmarks/` directory for detailed results:
- `benchmarks/baseline/` - Single-node benchmarks
- `benchmarks/scaling/` - Multi-node scaling tests
- `benchmarks/stress/` - Stress test results
- `benchmarks/soak/` - Long-duration test results

## Appendix B: Profiling Tools

**Recommended Tools:**
- `pprof` - Go profiling
- `perf` - Linux performance analysis
- `flamegraph` - Visualization
- `k6` - Load testing
- `hey` - HTTP benchmarking
- `sysbench` - System benchmarking

---

**Document Owner:** Performance Engineering Team
**Review Frequency:** Quarterly
**Next Review:** 2026-05-23
