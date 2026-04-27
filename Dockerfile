# ============================================================
# Multi-stage Dockerfile for backend-gmao monorepo
# Build any service using: --build-arg SERVICE_NAME=<name>
# Example: docker build --build-arg SERVICE_NAME=analytics-service .
# ============================================================

# --- Stage 1: Build ---
FROM golang:1.26-alpine AS builder

ARG SERVICE_NAME

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go module files first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy entire source tree
COPY . .

# Build the specific service binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/service \
    ./apps/${SERVICE_NAME}/cmd/api

# --- Stage 2: Runtime ---
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8080

ENTRYPOINT ["./service"]
