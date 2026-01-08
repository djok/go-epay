# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies (git, ca-certificates, gcc for CGO/SQLite)
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -a -o goepay ./cmd/goepay

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install ca-certificates for HTTPS calls to billing APIs
RUN apk --no-cache add ca-certificates wget

# Copy binary from builder
COPY --from=builder /app/goepay .

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Create non-root user for security
RUN adduser -D -g '' appuser && chown -R appuser:appuser /app/data
USER appuser

# Expose port
EXPOSE 8080

# Health check disabled in Dockerfile - use docker-compose healthcheck instead
# This allows configurable PORT via environment variable

# Run the binary
ENTRYPOINT ["./goepay"]
