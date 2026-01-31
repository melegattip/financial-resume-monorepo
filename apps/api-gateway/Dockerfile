# 🐳 Financial Resume Engine - Dockerfile Optimizado para Producción
# Multi-stage build para minimizar tamaño de imagen

# ================================
# STAGE 1: BUILD
# ================================
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git ca-certificates tzdata

# Configurar zona horaria
ENV TZ=America/Sao_Paulo

# Crear usuario no-root
RUN adduser -D -g '' appuser

# Establecer directorio de trabajo
WORKDIR /app

# Copiar go mod files primero (para cache de Docker layers)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Build optimizado para producción
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/api/main.go

# ================================
# STAGE 2: PRODUCTION
# ================================
FROM scratch

# Copiar certificados SSL
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar archivos de zona horaria
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copiar usuario
COPY --from=builder /etc/passwd /etc/passwd

# Copiar binario
COPY --from=builder /app/main /app/main

# Usar usuario no-root
USER appuser

# Puerto de la aplicación
EXPOSE 8080

# Variables de entorno por defecto
ENV GIN_MODE=release
ENV TZ=America/Sao_Paulo

# Health check - GCP Cloud Run maneja el health check automáticamente
# Removemos el HEALTHCHECK ya que scratch no tiene herramientas para hacer HTTP requests

# Comando por defecto
ENTRYPOINT ["/app/main"]
