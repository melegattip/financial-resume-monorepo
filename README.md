# Financial Resume - Monorepo

Ecosistema completo de gestión financiera personal con arquitectura de microservicios.

## 🏗️ Estructura

```
├── apps/                        # Aplicaciones desplegables
│   ├── api-gateway/             # Core API + Gateway (Go)
│   ├── ai-service/              # Análisis con IA (Go)
│   ├── gamification-service/    # Sistema de gamificación (Go)
│   ├── users-service/           # Autenticación y usuarios (Go)
│   └── frontend/                # Interfaz web (React)
│
├── packages/                    # Código compartido
│   └── go-shared/               # Bibliotecas Go comunes
│
├── infrastructure/              # Configuración de infraestructura
│   └── docker/                  # Docker Compose files
│
├── scripts/                     # Scripts de desarrollo
└── docs/                        # Documentación
```

## 🚀 Inicio Rápido

### Prerrequisitos
- Go 1.23+
- Node.js 18+
- Docker & Docker Compose

### Desarrollo Local

```bash
# Clonar repositorio
git clone https://github.com/tu-usuario/financial-resume-monorepo.git
cd financial-resume-monorepo

# Levantar bases de datos
docker-compose -f infrastructure/docker/docker-compose.yml up -d postgres-main postgres-users redis

# Ejecutar todos los servicios
./scripts/dev-up.sh
```

### Comandos Útiles

```bash
# Build de todos los servicios Go
go build ./apps/...

# Tests de todos los servicios
go test ./apps/... ./packages/...

# Solo un servicio específico
cd apps/api-gateway && go run cmd/api/main.go
```

## 📦 Servicios

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| api-gateway | 8080 | API principal y orquestación |
| ai-service | 8082 | Análisis financiero con IA |
| gamification-service | 8081 | Sistema de XP y logros |
| users-service | 8083 | Autenticación JWT/2FA |
| frontend | 3000 | Interfaz React |

## 🔧 Configuración

Copiar `.env.example` a `.env` y configurar las variables requeridas.

## 📚 Documentación

- [Arquitectura](docs/architecture/)
- [API Reference](docs/api/)
- [Guía de Desarrollo](docs/development/)

## 📄 Licencia

MIT
