# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go workspace files
COPY go.work go.work
COPY apps/monolith/go.mod apps/monolith/go.sum apps/monolith/

# Download dependencies
WORKDIR /app/apps/monolith
RUN go mod download

# Copy source code
WORKDIR /app
COPY apps/monolith/ apps/monolith/

# Build single binary
WORKDIR /app/apps/monolith
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd/server/

# Stage 2: Runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/server .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
