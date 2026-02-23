.PHONY: all build test clean install lint fmt vet benchmark load-test help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build output directory
BIN_DIR=bin

# Binary names
SIGNER_BIN=$(BIN_DIR)/signer
VERIFIER_BIN=$(BIN_DIR)/verifier
LEDGER_NODE_BIN=$(BIN_DIR)/ledger-node
IDENTITY_AUTHORITY_BIN=$(BIN_DIR)/identity-authority
AUDITOR_BIN=$(BIN_DIR)/auditor
KEY_CEREMONY_BIN=$(BIN_DIR)/key-ceremony

all: test build

build: $(BIN_DIR)
	$(GOBUILD) -o $(SIGNER_BIN) ./cmd/signer
	$(GOBUILD) -o $(VERIFIER_BIN) ./cmd/verifier
	$(GOBUILD) -o $(LEDGER_NODE_BIN) ./cmd/ledger-node
	$(GOBUILD) -o $(IDENTITY_AUTHORITY_BIN) ./cmd/identity-authority
	$(GOBUILD) -o $(AUDITOR_BIN) ./cmd/auditor
	$(GOBUILD) -o $(KEY_CEREMONY_BIN) ./cmd/key-ceremony

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-unit:
	$(GOTEST) -v -race -short ./...

test-integration:
	$(GOTEST) -v -race -run Integration ./...

test-adversarial:
	$(GOTEST) -v -race -run Adversarial ./tests/adversarial/...

test-fuzz:
	$(GOTEST) -v -fuzz=. -fuzztime=30s ./tests/fuzz/...

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -f coverage.out

install:
	$(GOMOD) download
	$(GOMOD) tidy

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run ./...

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

# Docker targets
docker-build:
	docker build -t civic-attest:latest .

docker-run:
	docker run -it civic-attest:latest

# SBOM generation
sbom:
	@which syft > /dev/null || (echo "syft not installed" && exit 1)
	syft packages . -o json > sbom/sbom.json
	syft packages . -o spdx > sbom/sbom.spdx

# Reproducible build
reproducible-build:
	$(GOBUILD) -trimpath -ldflags="-buildid=" -o $(BIN_DIR)/signer-reproducible ./cmd/signer

# Performance benchmarking
benchmark:
	$(GOTEST) -bench=. -benchmem -benchtime=10s ./...

benchmark-baseline:
	@mkdir -p benchmarks
	$(GOTEST) -bench=. -benchmem -benchtime=10s ./... > benchmarks/baseline.txt

benchmark-compare:
	@mkdir -p benchmarks
	$(GOTEST) -bench=. -benchmem -benchtime=10s ./... > benchmarks/current.txt
	@which benchstat > /dev/null || (echo "benchstat not installed, run: go install golang.org/x/perf/cmd/benchstat@latest" && exit 1)
	benchstat benchmarks/baseline.txt benchmarks/current.txt

# Load testing
load-test:
	@which k6 > /dev/null || (echo "k6 not installed" && exit 1)
	@mkdir -p tests/load
	k6 run tests/load/signature-load.js

stress-test:
	@which k6 > /dev/null || (echo "k6 not installed" && exit 1)
	@mkdir -p tests/load
	k6 run tests/load/stress-test.js

# Profiling
profile-cpu:
	@mkdir -p profile
	$(GOTEST) -cpuprofile=profile/cpu.prof -bench=. ./...
	@echo "Opening CPU profile at http://localhost:8080"
	go tool pprof -http=:8080 profile/cpu.prof

profile-mem:
	@mkdir -p profile
	$(GOTEST) -memprofile=profile/mem.prof -bench=. ./...
	@echo "Opening memory profile at http://localhost:8080"
	go tool pprof -http=:8080 profile/mem.prof

# Operational tasks
health-check:
	@echo "Performing health check..."
	@curl -f http://localhost:8080/health || (echo "Health check failed" && exit 1)
	@curl -f http://localhost:8080/ready || (echo "Readiness check failed" && exit 1)
	@echo "Health check passed"

# Security scanning
security-scan:
	@which gosec > /dev/null || (echo "gosec not installed, run: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec -fmt=json -out=security-report.json ./...

vulnerability-check:
	@which govulncheck > /dev/null || (echo "govulncheck not installed, run: go install golang.org/x/vuln/cmd/govulncheck@latest" && exit 1)
	govulncheck ./...

# Help target
help:
	@echo "Civic Attest Makefile Targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build                 - Build all binaries"
	@echo "  reproducible-build    - Build with reproducible output"
	@echo "  docker-build          - Build Docker image"
	@echo "  docker-run            - Run Docker container"
	@echo ""
	@echo "Test targets:"
	@echo "  test                  - Run all tests with coverage"
	@echo "  test-unit             - Run unit tests only"
	@echo "  test-integration      - Run integration tests"
	@echo "  test-adversarial      - Run adversarial tests"
	@echo "  test-fuzz             - Run fuzz tests"
	@echo ""
	@echo "Performance targets:"
	@echo "  benchmark             - Run performance benchmarks"
	@echo "  benchmark-baseline    - Create baseline benchmark"
	@echo "  benchmark-compare     - Compare with baseline"
	@echo "  load-test             - Run load tests"
	@echo "  stress-test           - Run stress tests"
	@echo "  profile-cpu           - CPU profiling"
	@echo "  profile-mem           - Memory profiling"
	@echo ""
	@echo "Operational targets:"
	@echo "  health-check          - Check service health"
	@echo ""
	@echo "Security targets:"
	@echo "  security-scan         - Run security scanner"
	@echo "  vulnerability-check   - Check for vulnerabilities"
	@echo ""
	@echo "Other targets:"
	@echo "  install               - Install dependencies"
	@echo "  clean                 - Clean build artifacts"
	@echo "  lint                  - Run linter"
	@echo "  fmt                   - Format code"
	@echo "  vet                   - Run go vet"
	@echo "  sbom                  - Generate SBOM"
	@echo "  help                  - Show this help message"
