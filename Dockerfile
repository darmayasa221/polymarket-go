# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25

# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Download deps before copying source — maximises layer cache reuse.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

# Fully static binary: no libc, no CGO, stripped debug info.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -extldflags '-static'" \
    -trimpath \
    -o /app/bin/server \
    ./cmd/server

# ── Final stage ───────────────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /app/bin/server ./server

# SQLite database directory — mounted as a volume in compose.
RUN mkdir -p /app/db && chown appuser:appuser /app/db

USER appuser

EXPOSE 3000

# Use /health/live (no dependency checks) so the container never goes
# unhealthy due to a downstream service outage.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/health/live || exit 1

ENTRYPOINT ["./server"]
