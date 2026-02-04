# Financial Resume - Monorepo

Backend unificado + Frontend para gestión financiera personal.

## 🏗️ Arquitectura

```
                    ┌─────────────────────────────────┐
                    │     Unified Backend (Docker)     │
                    │           Port 8080              │
                    │  ┌─────────────────────────────┐ │
                    │  │         Nginx Proxy          │ │
                    │  └──────────┬──────────────────┘ │
                    │             │                    │
     ┌──────────────┼─────────────┼────────────────────┼──────────────┐
     │              │             │                    │              │
     ▼              ▼             ▼                    ▼              │
┌─────────┐   ┌──────────┐  ┌───────────┐    ┌─────────────────┐      │
│ Users   │   │ Engine   │  │ AI        │    │ Gamification    │      │
│:8083    │   │:8081     │  │:8082      │    │:8084            │      │
└─────────┘   └──────────┘  └───────────┘    └─────────────────┘      │
                    └─────────────────────────────────────────────────┘
```

## 📂 Estructura

```
├── apps/
│   ├── api-gateway/        # Engine principal (Go)
│   ├── ai-service/         # Análisis con IA (Go)
│   ├── gamification-service/# XP y logros (Go)
│   ├── users-service/      # Auth JWT/2FA (Go)
│   └── frontend/           # React
├── infrastructure/docker/
│   ├── nginx-unified.conf  # Reverse proxy config
│   └── supervisord.conf    # Process manager
├── Dockerfile              # Multi-stage build
└── render.yaml             # Render.com blueprint
```

## 🚀 Deploy a Render

1. Push el repositorio a GitHub
2. En Render: **New** → **Blueprint** → Conectar repo
3. Render detectará `render.yaml` automáticamente
4. Configurar variables de entorno (DBs, JWT_SECRET, etc.)

**Resultado**: 2 servicios en Render
- `financial-resume-backend` (~$7/mes)
- `financial-resume-frontend` (~$7/mes)

## 🔧 Desarrollo Local

```bash
# Levantar DBs
docker-compose -f infrastructure/docker/docker-compose.yml up -d postgres-main postgres-users redis

# Opción 1: Correr servicios individuales
cd apps/api-gateway && go run cmd/api/main.go
cd apps/users-service && go run cmd/api/main.go
# etc...

# Opción 2: Build Docker completo
docker build -t financial-resume-backend .
docker run -p 8080:8080 financial-resume-backend
```

## 📡 API Routing

| Path | Servicio |
|------|----------|
| `/api/v1/users/*` | users-service |
| `/api/v1/auth/*` | users-service |
| `/api/v1/ai/*` | ai-service |
| `/api/v1/gamification/*` | gamification-service |
| `/api/v1/*` (resto) | api-gateway (engine) |
| `/swagger/*` | api-gateway |
