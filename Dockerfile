# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN make build

# Stage 2: Runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 civic && \
    adduser -D -u 1000 -G civic civic

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/bin/* /app/bin/

# Set ownership
RUN chown -R civic:civic /app

USER civic

# Default to running the ledger node
ENTRYPOINT ["/app/bin/ledger-node"]
