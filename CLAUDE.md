# financial-resume-monorepo Development Guidelines

Last updated: 2026-03-18 — Estado actual: Phase 7 (multi-tenant) completo

---

## Active Technologies

### Backend (Go Monolith)
- **Go 1.24** — `go.work` workspace, single module `./apps/monolith`
- **Gin 1.10.1** — HTTP framework
- **GORM 1.31.1** + `gorm.io/driver/postgres 1.6.0` — ORM + PostgreSQL driver
- **zerolog 1.34.0** — Structured logging
- **golang-jwt/v5 5.3.1** — JWT auth (access + refresh tokens)
- **google/uuid 1.6.0** — UUID generation
- **pquerna/otp 1.5.0** — TOTP/2FA
- **x/crypto** — bcrypt password hashing
- **lib/pq 1.11.2** — PostgreSQL raw SQL support
- **gin-contrib/cors 1.7.6** — CORS middleware
- **godotenv 1.5.1** — .env loading
- **stretchr/testify 1.11.1** — Testing assertions
- **DATA-DOG/go-sqlmock 1.5.2** — SQL mocking for tests

### Frontend (React SPA)
- **React 18.2** + React DOM 18.2
- **React Router DOM 6.8.1** — Client-side routing
- **Axios 1.3.4** — HTTP client
- **Tailwind CSS 3.2.7** — Styling (dark mode via `class` strategy)
- **Recharts 2.5.0** — Charts and visualizations
- **React Hook Form 7.43.5** — Form handling
- **date-fns 2.29.3** — Date utilities
- **React Icons 4.12.0** — Icon library
- **React Hot Toast 2.4.0** — Notifications
- **Craco 7.1.0** — Create React App config override

### Infrastructure
- **PostgreSQL 15** — Single database `financial_resume`
- **Docker** — Local development via docker-compose
- **Render.com** — Production deployment
- **nginx** — Frontend static server in production

---

## Project Structure

```
financial-resume-monorepo/
├── apps/
│   ├── monolith/                    # Go backend
│   │   ├── cmd/
│   │   │   ├── server/main.go       # Entry point
│   │   │   ├── migrate/             # Data migration utility
│   │   │   └── verify-db/           # DB verification
│   │   ├── internal/
│   │   │   ├── modules/
│   │   │   │   ├── auth/            # JWT, 2FA, users, preferences
│   │   │   │   ├── transactions/    # Expenses, incomes, categories
│   │   │   │   ├── tenants/         # Multi-tenant, RBAC, audit logs
│   │   │   │   ├── gamification/    # XP, achievements, challenges
│   │   │   │   ├── budgets/         # Spending limits + alerts
│   │   │   │   ├── savings/         # Goals + transactions
│   │   │   │   ├── recurring/       # Scheduled transactions
│   │   │   │   ├── analytics/       # Dashboard + trends
│   │   │   │   └── ai/              # OpenAI integration
│   │   │   ├── shared/
│   │   │   │   ├── ports/           # EventBus interface
│   │   │   │   ├── events/          # InMemoryEventBus + DomainEvent
│   │   │   │   └── email/           # Resend/SMTP/NoOp email service
│   │   │   └── infrastructure/
│   │   │       ├── config/          # AppConfig + Load()
│   │   │       ├── middleware/      # Auth, Permission, CORS, Logging
│   │   │       ├── database/        # PostgreSQL connection + migrations
│   │   │       └── http/            # Router + Server
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── .env / .env.local
│   │
│   └── frontend/                    # React SPA
│       ├── src/
│       │   ├── App.jsx              # Routes + provider stack
│       │   ├── contexts/            # Auth, Theme, Period, Gamification, Tenant
│       │   ├── pages/               # 19 pages
│       │   ├── components/          # Layout, guards, modals, widgets
│       │   ├── services/            # API clients + service layer
│       │   ├── hooks/               # Custom hooks
│       │   └── config/              # Environment detection
│       ├── Dockerfile
│       └── package.json
│
├── scripts/
│   ├── docker-rebuild.sh            # Rebuild Docker images (preserve postgres)
│   └── dev.sh                       # Local dev without Docker
├── docs/03-architecture/            # DB schema docs
├── specs/                           # Feature specs
├── docker-compose.yml
├── go.work                          # Go workspace (./apps/monolith only)
└── render.yaml                      # Render.com deployment blueprint
```

---

## Commands

### Local Development
```bash
# Full Docker environment (recommended)
docker compose up

# Rebuild after code changes
./scripts/docker-rebuild.sh              # Rebuild backend + frontend
./scripts/docker-rebuild.sh backend      # Backend only
./scripts/docker-rebuild.sh frontend     # Frontend only
./scripts/docker-rebuild.sh --no-cache   # Force clean build

# URLs: backend http://localhost:8080 | frontend http://localhost:3000

# Dev sin Docker (postgres via docker)
./scripts/dev.sh
# Then: cd apps/monolith && go run ./cmd/server/
# Then: cd apps/frontend && npm start
```

### Backend
```bash
cd apps/monolith
go build ./...                         # Compilar
go test ./...                          # Tests
go test -race -coverprofile=coverage.out -covermode=atomic ./...  # Tests con coverage
go run ./cmd/server/                   # Ejecutar servidor
```

### Frontend
```bash
cd apps/frontend
npm run dev      # Desarrollo (craco)
npm run build    # Build producción
npm test         # Tests
```

---

## Backend Architecture

### Module Pattern
Cada módulo vive en `internal/modules/<name>/` con esta estructura:
```
<name>/
  domain/        # Structs + métodos de negocio
  repository/    # GORM models + implementación
  service/       # Lógica de negocio (en módulos complejos)
  handlers/      # Gin handlers
  ports/         # Interfaces (repos, services)
  module.go      # New() + RegisterRoutes() + RegisterSubscribers()
```

### Module Signature
```go
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig,
         eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware,
         permMW *middleware.PermissionMiddleware) *Module
```
Excepción: `ai` no recibe authMW/permMW. `savings` no recibe `cfg`.

### JWT Claims
```go
// Context keys (gin.Context):
"user_id"    string
"user_email" string
"tenant_id"  string
"role"       string
"token"      string
```

### Event Bus
```go
// Interface en internal/shared/ports/event_bus.go
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(eventType string, handler EventHandler)
}
// Implementación: InMemoryEventBus (async goroutines, panic recovery)
```

**Eventos publicados**:
- `user.registered`
- `expense.created/updated/deleted`
- `income.created/updated/deleted`
- `recurring.created/updated/deleted/executed/paused/resumed`
- `savings_goal.created/updated/deleted/achieved`
- `budget.threshold_crossed`

### Email Service
Prioridad de selección en `shared/email/service.go`:
1. Resend API key → `ResendEmailService`
2. SMTP config → `SMTPEmailService`
3. Sin credenciales → `NoOpEmailService` (logs only, dev mode)

### Middleware Stack
- `CORS(origins)` — origenes configurables
- `RequestID()` — añade request ID al contexto
- `RequestLogging(logger)` — zerolog request/response logging
- `authMW.RequireAuth()` — valida Bearer token, setea claims en contexto
- `authMW.OptionalAuth()` — igual pero no bloquea
- `permMW.Require("permission")` — valida permiso de rol en tenant

---

## Module Routes Reference

### Auth — `/api/v1`
```
POST   /auth/register
POST   /auth/login
POST   /auth/check-2fa
POST   /auth/refresh
GET    /auth/verify-email/:token
POST   /auth/request-password-reset
POST   /auth/reset-password

POST   /users/logout                    [auth]
POST   /users/switch-tenant             [auth]
GET    /users/profile                   [auth]
PUT    /users/profile                   [auth]
POST   /users/profile/avatar            [auth]
PUT    /users/change-password           [auth]
POST   /users/2fa/setup                 [auth]
POST   /users/2fa/enable                [auth]
POST   /users/2fa/disable               [auth]
POST   /users/2fa/verify                [auth]
GET    /users/preferences               [auth]
PUT    /users/preferences               [auth]
GET    /users/notifications             [auth]
PUT    /users/notifications             [auth]
GET    /users/export                    [auth]
DELETE /users/account                   [auth]
```

### Transactions — `/api/v1`
```
POST/GET       /expenses               [auth + permission]
GET/PUT/DELETE /expenses/:id           [auth + permission]
POST/GET       /incomes                [auth + permission]
GET/PUT/DELETE /incomes/:id            [auth + permission]
GET/POST       /categories             [auth]
PATCH/DELETE   /categories/:id         [auth]
```

### Tenants — `/api/v1/tenants`
```
GET    /tenants/list                                [auth]
POST   /tenants/join                                [auth]
GET    /tenants/me                                  [auth]
PUT    /tenants/me                                  [auth + manage_tenant]
DELETE /tenants/me                                  [auth + delete_tenant]
GET    /tenants/me/permissions                      [auth]
GET    /tenants/me/members                          [auth]
PUT    /tenants/me/members/:userID/role             [auth + manage_roles]
DELETE /tenants/me/members/:userID                  [auth + remove_members]
GET    /tenants/me/invitations                      [auth + invite_members]
POST   /tenants/me/invitations                      [auth + invite_members]
DELETE /tenants/me/invitations/:code                [auth + invite_members]
GET    /tenants/me/audit                            [auth + view_audit_logs]
```

### Gamification — `/api/v1/gamification`
```
GET    /gamification/profile
GET    /gamification/stats
GET    /gamification/achievements
GET    /gamification/features
GET    /gamification/features/:featureKey/access
GET    /gamification/challenges/daily
GET    /gamification/challenges/weekly
POST   /gamification/challenges/progress
POST   /gamification/actions
GET    /gamification/behavior-profile
```

### Budgets — `/api/v1/budgets`
```
POST/GET       /budgets                [auth + manage_budgets]
GET            /budgets/status
GET            /budgets/dashboard
GET/PUT/DELETE /budgets/:id
```

### Savings — `/api/v1/savings-goals`
```
POST/GET       /savings-goals          [auth + manage_savings]
GET            /savings-goals/dashboard
GET            /savings-goals/summary
GET/PUT/DELETE /savings-goals/:id
POST           /savings-goals/:id/deposit
POST           /savings-goals/:id/withdraw
POST           /savings-goals/:id/pause
POST           /savings-goals/:id/resume
POST           /savings-goals/:id/cancel
GET            /savings-goals/:id/transactions
```

### Recurring — `/api/v1/recurring-transactions`
```
POST/GET       /recurring-transactions               [auth + manage_recurring]
GET            /recurring-transactions/dashboard
GET            /recurring-transactions/due
GET            /recurring-transactions/projection
POST           /recurring-transactions/batch/process  [auth + manage_recurring]
POST           /recurring-transactions/batch/notify
GET/PUT/DELETE /recurring-transactions/:id
POST           /recurring-transactions/:id/pause
POST           /recurring-transactions/:id/resume
POST           /recurring-transactions/:id/execute
```
**Nota**: `StartScheduler(ctx)` ejecuta un goroutine que procesa recurring transactions al inicio y cada 1 hora.

### Analytics — `/api/v1`
```
GET    /dashboard                      [auth + view_data]
GET    /analytics/expenses             [auth]
GET    /analytics/incomes              [auth]
GET    /analytics/categories           [auth]
GET    /analytics/monthly              [auth]
GET    /insights/financial-health      [auth]
GET    /reports                        [auth + view_data]
```

### AI — `/api/v1/ai`
```
POST   /ai/health-analysis
POST   /ai/insights
POST   /ai/can-i-buy
POST   /ai/alternatives
POST   /ai/credit-plan
POST   /ai/credit-score
```
**Nota**: Si `OPENAI_API_KEY` no está configurado, cae en mock responses.

---

## Database Schema

### Tablas
| Tabla | Módulo | Soft Delete | Notas |
|-------|--------|-------------|-------|
| `users` | auth | `gorm.DeletedAt` | Excepción al patrón estándar |
| `user_preferences` | auth | no | 1:1 con users |
| `user_notification_settings` | auth | no | 1:1 con users |
| `user_two_fa` | auth | no | 1:1 con users |
| `categories` | transactions | `*time.Time` | indexed user_id + priority |
| `expenses` | transactions | `*time.Time` | indexed user_id + transaction_date |
| `incomes` | transactions | `*time.Time` | indexed user_id + received_date |
| `budgets` | budgets | `*time.Time` | |
| `savings_goals` | savings | `*time.Time` | |
| `savings_transactions` | savings | `*time.Time` | |
| `recurring_transactions` | recurring | `*time.Time` | |
| `user_gamification` | gamification | `*time.Time` | unique user_id |
| `achievements` | gamification | no | indexed user_id |
| `user_actions` | gamification | **NO** | audit trail inmutable |
| `tenants` | tenants | no | |
| `tenant_members` | tenants | no | |
| `role_permissions` | tenants | no | |
| `tenant_invitations` | tenants | no | |
| `audit_logs` | tenants | **NO** | inmutable |

### Soft Delete Pattern
```go
// GORM models: usar *time.Time (NO gorm.DeletedAt), excepto auth.User
DeletedAt *time.Time

// Queries: filtro explícito
WHERE deleted_at IS NULL

// Delete: actualización directa
UPDATE SET deleted_at = NOW()
```

### ID Generation
| Módulo | Formato |
|--------|---------|
| auth, transactions, recurring | `uuid.New().String()` (UUID completo) |
| budgets | `bud_` + primeros 8 chars del UUID |
| savings goals | `goal_` + primeros 8 chars del UUID |
| savings transactions | `stxn_` + primeros 8 chars del UUID |

---

## Code Conventions

### Handlers
- `userID` **siempre** desde JWT: `c.Get("user_id")`
- Validar ownership antes de mutaciones (userID check)
- Publicar eventos **después** de persistir en DB
- AutoMigrate en `New()`, nunca en `RegisterRoutes()`
- Rutas bajo `/api/v1`, protegidas con `authMW.RequireAuth()`

### Multi-tenant
- Todas las entidades tienen `tenant_id` (añadido vía migration)
- PostgreSQL RLS policies activas
- Users pueden pertenecer a múltiples tenants (switch via `/users/switch-tenant`)
- `tenant_id` y `role` en JWT claims, seteados por `RequireAuth()`

### Builder Pattern
Usado en structs complejos: `BudgetBuilder`, `SavingsGoalBuilder`, `RecurringTransactionBuilder`

### Idioma
- Documentación general → **Español**
- Código, comentarios técnicos, commits → **Inglés**

---

## Frontend Architecture

### Provider Stack (orden)
```
Router → ThemeProvider → AuthProvider → TenantProvider → PeriodProvider → GamificationProvider
```

### Contextos
| Context | Hook | Key Exports |
|---------|------|-------------|
| AuthContext | `useAuth()` | `user`, `isAuthenticated`, `login()`, `logout()`, `register()` |
| TenantContext | `useTenant()` | `currentTenant`, `myRole`, `permissions`, `hasPermission()`, `switchTenant()` |
| PeriodContext | `usePeriod()` | `selectedYear/Month`, `getFilterParams()`, `toggleBalancesVisibility()` |
| GamificationContext | `useGamification()` | `userProfile`, `achievements`, `features`, `FEATURE_GATES` |
| ThemeContext | `useTheme()` | `theme`, `isDark`, `toggleTheme()` |

### Rutas Protegidas
```
/dashboard, /expenses, /incomes, /categories, /reports
/recurring-transactions, /achievements, /settings
/budgets           → FeatureGuard: BUDGETS
/savings-goals     → FeatureGuard: SAVINGS_GOALS
/insights          → FeatureGuard: AI_INSIGHTS
/tenant-settings   → TenantSettings (role-gated)
/audit-logs        → AuditLogs (view_audit_logs permission)
```

### Services Layer
| Service | Propósito |
|---------|-----------|
| `apiClient.js` | Base axios client — auto Bearer token + X-Caller-ID |
| `authService.js` | Login, register, 2FA, profile, tokens |
| `tenantService.js` | Multi-tenant CRUD |
| `gamificationAPI.js` | Gamification backend API |
| `dataService.js` | Cache layer (5-min TTL) para dashboard/analytics |
| `configService.js` | Dynamic API URL detection |

### Entornos (auto-detect por hostname)
| Entorno | API Base URL |
|---------|-------------|
| development (localhost) | `http://localhost:8080/api/v1` |
| docker | `http://financial-resume-engine:8080/api/v1` |
| render / production | `https://financial-resume-monorepo-l71a.onrender.com/api/v1` |

### Auth Storage
```javascript
localStorage keys:
  auth_token          // JWT access token
  auth_refresh_token  // JWT refresh token
  auth_user           // User profile JSON
  auth_expires_at     // Token expiry timestamp
```

### Guards
- `<ProtectedRoute>` — require autenticación
- `<PublicOnlyRoute>` — redirige si ya autenticado
- `<FeatureGuard feature="">` — require nivel de gamificación
- `<RoleGuard permission="" role="" any>` — require permiso de tenant

---

## Deployment

### Docker Local
```yaml
# docker-compose.yml — servicios:
# postgres (5432) → financial_resume DB
# backend  (8080) → Go monolith
# frontend (3000) → React SPA
```

### Render.com (Producción)
- `render.yaml` define el blueprint de deployment
- Backend: `financial.niloft.com` (o onrender.com)
- Frontend sirve desde nginx

---

## CI/CD

**Archivo**: `.github/workflows/monolith-ci.yml`

**Triggers**: Push/PR sobre `apps/monolith/**` o `go.work`

**Jobs**:
1. **lint** — `golangci-lint` con `.golangci.yml`
2. **test** — `go test -race -coverprofile` — **falla si coverage < 80%**
3. **build** — `CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/server/`
