# --- Stage 1: build the binary ---
FROM golang:1.23.1-alpine AS builder

RUN apk add --no-cache git
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o ha2trmnl ./cmd/ha2trmnl

# --- Stage 2: runtime container ---
FROM alpine:3.19

RUN apk add --no-cache bash curl ca-certificates busybox-suid shadow \
  && useradd -u 1000 ha2trmnl || true

COPY --from=builder /src/ha2trmnl /usr/local/bin/ha2trmnl
COPY docker/entrypoint.sh /entrypoint.sh

RUN chmod +x /usr/local/bin/ha2trmnl /entrypoint.sh

ENV CRON_SCHEDULE="*/5 * * * *"
CMD ["/entrypoint.sh"]