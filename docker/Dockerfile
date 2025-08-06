# --- Stage 1: build the binary ---
FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git
WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 go build -o ha2trmnl ./cmd/ha2trmnl

# --- Stage 2: runtime container ---
FROM alpine:3.19

# Install runtime deps
RUN apk add --no-cache bash curl ca-certificates busybox-suid shadow tzdata \
  && useradd -u 1000 ha2trmnl || true

# Copy binary from builder
COPY --from=builder /src/ha2trmnl /usr/local/bin/ha2trmnl

# Ensure it's executable
RUN chmod +x /usr/local/bin/ha2trmnl

# Use non-root by default
USER ha2trmnl

# Set default mode (can override via `docker run ... push`)
ENTRYPOINT ["/usr/local/bin/ha2trmnl"]
CMD ["serve", "/config/config.yaml"]