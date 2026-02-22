.PHONY: all build test clean install lint fmt vet

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
