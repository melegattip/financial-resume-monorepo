# Financial Gamification Service 🎮

Microservicio independiente para gestión de gamificación del Financial Resume Engine.

## Estado actual

- Arquitectura Clean con separación de `domain`, `ports`, `usecases`, `handlers`.
- DB PostgreSQL propia (`scripts/init.sql`).
- Idempotencia diaria para `view_dashboard`; `view_insight` no otorga XP.
- Suite de tests en `internal/core/usecases/*_test.go` en verde.

## 🚀 Características

- **Sistema de XP y Niveles**: 10 niveles desde "Financial Newbie" hasta "Financial Magnate"
- **Achievements**: Sistema de logros con progreso automático
- **Leaderboard**: Ranking de usuarios por XP
- **API REST**: Endpoints completos con autenticación JWT
- **Clean Architecture**: Separación clara de responsabilidades
- **PostgreSQL**: Base de datos optimizada con índices

## 📋 Endpoints

### Autenticados (requieren JWT)
- `GET /api/v1/gamification/profile` - Perfil de gamificación del usuario
- `GET /api/v1/gamification/stats` - Estadísticas detalladas
- `GET /api/v1/gamification/achievements` - Logros del usuario
- `POST /api/v1/gamification/actions` - Registrar acción y otorgar XP

### 🔒 Feature Gates (Nuevos Endpoints)
- `GET /api/v1/gamification/features` - Todas las features del usuario (desbloqueadas y bloqueadas)
- `GET /api/v1/gamification/features/{featureKey}/access` - Verificar acceso a feature específica

### Públicos
- `GET /api/v1/gamification/action-types` - Tipos de acciones disponibles (39 tipos)
- `GET /api/v1/gamification/levels` - Información de niveles (10 niveles optimizados)
- `GET /health` - Health check

### Cambios recientes
- Reglas: `view_dashboard` idempotente diario; `view_insight` XP=0.

## 🛠️ Instalación

### Prerrequisitos
- Go 1.21+
- PostgreSQL 13+

### Configuración
1. Clonar el repositorio:
```bash
git clone https://github.com/melegattip/financial-gamification-service.git
cd financial-gamification-service
```

2. Instalar dependencias:
```bash
go mod tidy
```

3. Configurar base de datos:
```bash
# Crear base de datos
createdb gamification_db

# Ejecutar migraciones
psql -d gamification_db -f scripts/init.sql
```

4. Ejecutar el servicio:
```bash
go run cmd/api/main.go
```

El servicio estará disponible en `http://localhost:8081`

## 🎯 Sistema de Puntos Rediseñado

### Acciones y XP Base (Nuevos Valores Balanceados)

#### 🏠 Acciones Básicas (Disponibles desde Nivel 1)
- **view_dashboard**: 2 XP - Ver dashboard principal
- **view_expenses/incomes/categories**: 1 XP - Ver listas básicas
- **view_analytics**: 3 XP - Ver reportes básicos

#### 💰 Transacciones (Motor Principal de Progresión)
- **create_expense/income**: 8 XP - Registrar transacciones
- **update_expense/income**: 5 XP - Actualizar transacciones
- **delete_expense/income**: 3 XP - Eliminar transacciones

#### 🏷️ Organización
- **create_category**: 10 XP - Crear categoría personalizada
- **update_category**: 5 XP - Actualizar categoría
- **assign_category**: 3 XP - Categorizar transacción

#### 🎯 Engagement y Streaks
- **daily_login**: 5 XP - Login diario
- **weekly_streak**: 25 XP - Racha de 7 días
- **monthly_streak**: 100 XP - Racha de 30 días
- **complete_profile**: 50 XP - Completar perfil

#### 🏆 Challenges
- **daily_challenge_complete**: 20 XP - Completar challenge diario
- **weekly_challenge_complete**: 75 XP - Completar challenge semanal

#### 🔓 Features Desbloqueables
- **create_savings_goal**: 15 XP - Crear meta de ahorro (Nivel 3+)
- **create_budget**: 20 XP - Crear presupuesto (Nivel 5+)
- **use_ai_analysis**: 10 XP - Usar análisis de IA (Nivel 7+)

### Multiplicadores por Entidad
- **insight**: 1.0x
- **suggestion**: 1.2x
- **pattern**: 1.1x

### Niveles (Thresholds Optimizados)
| Nivel | Nombre | XP Requerido | Features Desbloqueadas |
|-------|--------|--------------|------------------------|
| 1 | Financial Newbie | 0 | Básicas |
| 2 | Money Tracker | 75 | - |
| 3 | Smart Saver | 200 | 🔓 **Metas de Ahorro** |
| 4 | Budget Master | 400 | - |
| 5 | Financial Planner | 700 | 🔓 **Presupuestos** |
| 6 | Investment Seeker | 1,200 | - |
| 7 | Wealth Builder | 1,800 | 🔓 **IA Financiera** |
| 8 | Financial Strategist | 2,600 | - |
| 9 | Money Mentor | 3,600 | - |
| 10 | Financial Magnate | 5,500 | - |

## 🏆 Achievements Rediseñados (Sin Dependencias)

### 💰 Achievements de Transacciones (Base de Progresión)
- **🌱 Primer Paso**: Registra tu primera transacción (25 XP)
- **📝 Aprendiz Financiero**: Registra 10 transacciones (50 XP)
- **💎 Maestro de Transacciones**: Registra 100 transacciones (200 XP)

### 🏷️ Achievements de Organización
- **🎨 Creador de Categorías**: Crea 5 categorías personalizadas (75 XP)
- **📊 Expert en Organización**: Categoriza 50 transacciones (100 XP)

### 🔥 Achievements de Engagement
- **⚡ Guerrero Semanal**: Mantén una racha de 7 días (100 XP)
- **👑 Leyenda Mensual**: Mantén una racha de 30 días (500 XP)

### 📈 Achievements de Análisis
- **🔍 Explorador de Datos**: Revisa analytics 25 veces (75 XP)

### 🎯 Achievements de Features Desbloqueables (Nivel 3+)
- **💰 Pionero del Ahorro**: Crea tu primera meta de ahorro (100 XP)
- **📊 Gurú de Presupuestos**: Crea 3 presupuestos (150 XP)
- **🤖 Pionero de IA**: Usa 10 análisis de IA (200 XP)

## 🏗️ Arquitectura

```
financial-gamification-service/
├── cmd/api/                 # Punto de entrada
├── internal/
│   ├── core/
│   │   ├── domain/         # Entidades de negocio
│   │   ├── ports/          # Interfaces
│   │   └── usecases/       # Lógica de negocio
│   ├── handlers/           # Controladores HTTP
│   └── infrastructure/
│       ├── repository/     # Implementación de persistencia
│       └── http/
│           └── middleware/ # Middleware JWT
├── pkg/
│   └── db/                 # Configuración de DB
└── scripts/
    └── init.sql           # Schema de base de datos
```

## 🔧 Configuración

### Variables de Entorno
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=gamification_db
JWT_SECRET=your-secret-key
PORT=8081
```

### Docker (Opcional)
```bash
# Construir imagen
docker build -t financial-gamification-service .

# Ejecutar con Docker Compose
docker-compose up -d
```

## 📊 Base de Datos

### Tablas Principales
- **user_gamification**: Estado de gamificación del usuario
- **achievements**: Logros con progreso
- **user_actions**: Tracking de acciones para XP

### Índices Optimizados
- Búsquedas por user_id
- Ordenamiento por total_xp (leaderboard)
- Filtros por tipo de achievement
- Ordenamiento temporal de acciones

## 🧪 Testing

```bash
# Ejecutar tests
go test ./...

# Test de cobertura
go test -cover ./...

# Test con verbose
go test -v ./...
```

## 📈 Monitoreo

### Health Check
```bash
curl http://localhost:8081/health
```

### Métricas
- Tiempo de respuesta de endpoints
- Conexiones activas a DB
- XP otorgado por período
- Achievements desbloqueados

## 🚀 Despliegue

### Producción
1. Configurar variables de entorno
2. Ejecutar migraciones de DB
3. Construir binario: `go build -o gamification-service cmd/api/main.go`
4. Ejecutar con supervisor/systemd

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gamification-service cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/gamification-service .
CMD ["./gamification-service"]
```

## 🤝 Integración

### Con Financial Resume Engine
```javascript
// Frontend: Registrar acción
const recordAction = async (actionType, entityType, entityId) => {
  await fetch('http://localhost:8081/api/v1/gamification/actions', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      action_type: actionType,
      entity_type: entityType,
      entity_id: entityId,
      description: `User ${actionType} ${entityType}`
    })
  });
};
```

### Con Backend Principal
```go
// Llamar al microservicio desde el backend
func recordGamificationAction(userID, actionType string) {
    payload := map[string]string{
        "action_type": actionType,
        "entity_type": "insight",
        "entity_id": "insight-123",
    }
    
    // HTTP call to gamification service
    // ...
}
```

## 📝 Licencia

MIT License - ver archivo LICENSE para detalles.

## 👥 Contribuir

1. Fork del proyecto
2. Crear rama feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear Pull Request

---

**Desarrollado con ❤️ para el Financial Resume Engine** 