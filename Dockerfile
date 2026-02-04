# Financial Resume - Unified Backend
# Multi-stage build for all Go services in a single container

FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy all source code
COPY . .

# Build API Gateway (Engine)
WORKDIR /app/apps/api-gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gateway ./cmd/api

# Build Users Service  
WORKDIR /app/apps/users-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/users-service ./cmd/api

# Build AI Service
WORKDIR /app/apps/ai-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/ai-service ./cmd/api

# Build Gamification Service
WORKDIR /app/apps/gamification-service
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gamification-service ./cmd/api

# Final stage - Alpine with all binaries
FROM alpine:3.19

RUN apk add --no-cache ca-certificates supervisor nginx

WORKDIR /app

# Copy all compiled binaries
COPY --from=builder /app/bin/api-gateway /app/bin/
COPY --from=builder /app/bin/users-service /app/bin/
COPY --from=builder /app/bin/ai-service /app/bin/
COPY --from=builder /app/bin/gamification-service /app/bin/

# Copy nginx config for reverse proxy
COPY infrastructure/docker/nginx-unified.conf /etc/nginx/nginx.conf

# Copy supervisor config
COPY infrastructure/docker/supervisord.conf /etc/supervisord.conf

# Expose single unified port
EXPOSE 8080

# Start supervisor which manages all services
CMD ["supervisord", "-c", "/etc/supervisord.conf"]
