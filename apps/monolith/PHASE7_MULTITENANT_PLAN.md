# Phase 7 — Multi-Tenant Implementation Plan

**Fecha:** 2026-02-25
**Decisiones de diseño:**
- Grupo flexible (familia/empresa/equipo)
- Visibilidad: todos los miembros ven todo dentro del tenant
- Aislamiento: app-level enforcement (tenant_id en queries) + PostgreSQL RLS como capa defensiva
- Onboarding: código de invitación
- Roles: owner, admin, member, viewer
- Permisos: RBAC dinámico en DB (per-tenant customizable)
- Gamificación: personal (sin cambios)
- Un user → un solo tenant activo
- Migración: cada user existente → owner de tenant personal automático

---

## Alcance de cambios

| Categoría | Cantidad |
|---|---|
| Archivos nuevos (Go) | ~28 |
| Archivos modificados (Go) | ~32 |
| Archivos nuevos (Frontend) | ~4 |
| Archivos modificados (Frontend) | ~8 |
| Migraciones SQL | ~12 |

---

## Fase A — Schema de Base de Datos

**Objetivo:** Crear las nuevas tablas, migrar datos existentes, sin romper nada.

### A1. Nuevas tablas (SQL migration)

```sql
-- tenants
CREATE TABLE tenants (
  id          VARCHAR(50) PRIMARY KEY,        -- "tnt_" + 8-char UUID
  name        VARCHAR(255) NOT NULL,
  slug        VARCHAR(100) UNIQUE NOT NULL,   -- URL-friendly identifier
  owner_id    VARCHAR(255) NOT NULL REFERENCES users(id),
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  plan        VARCHAR(20) NOT NULL DEFAULT 'free',
  settings    JSONB,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at  TIMESTAMPTZ
);

-- tenant_members
CREATE TABLE tenant_members (
  id          VARCHAR(50) PRIMARY KEY,
  tenant_id   VARCHAR(50) NOT NULL REFERENCES tenants(id),
  user_id     VARCHAR(255) NOT NULL REFERENCES users(id),
  role        VARCHAR(20) NOT NULL DEFAULT 'member',  -- owner|admin|member|viewer
  invited_by  VARCHAR(255) REFERENCES users(id),
  joined_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, user_id)
);

-- permissions (catalog)
CREATE TABLE permissions (
  key         VARCHAR(100) PRIMARY KEY,
  description TEXT,
  category    VARCHAR(50)  -- 'data'|'member_management'|'settings'|'admin'
);

-- role_permissions (per-tenant customizable)
CREATE TABLE role_permissions (
  id             VARCHAR(50) PRIMARY KEY,
  tenant_id      VARCHAR(50) NOT NULL REFERENCES tenants(id),
  role           VARCHAR(20) NOT NULL,
  permission_key VARCHAR(100) NOT NULL REFERENCES permissions(key),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(tenant_id, role, permission_key)
);
CREATE INDEX idx_role_permissions_lookup ON role_permissions(tenant_id, role);

-- tenant_invitations
CREATE TABLE tenant_invitations (
  id          VARCHAR(50) PRIMARY KEY,
  tenant_id   VARCHAR(50) NOT NULL REFERENCES tenants(id),
  code        VARCHAR(20) UNIQUE NOT NULL,
  role        VARCHAR(20) NOT NULL DEFAULT 'member',
  created_by  VARCHAR(255) NOT NULL REFERENCES users(id),
  expires_at  TIMESTAMPTZ,
  max_uses    INT NOT NULL DEFAULT 10,
  used_count  INT NOT NULL DEFAULT 0,
  is_active   BOOLEAN NOT NULL DEFAULT TRUE,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- audit_logs (sin soft delete — inmutable)
CREATE TABLE audit_logs (
  id           VARCHAR(50) PRIMARY KEY,
  tenant_id    VARCHAR(50) NOT NULL,
  user_id      VARCHAR(255) NOT NULL,
  action       VARCHAR(50) NOT NULL,   -- CREATE|UPDATE|DELETE|LOGIN|INVITE|ROLE_CHANGE|SETTING_CHANGE
  entity_type  VARCHAR(50),            -- 'expense'|'income'|'budget'|'member'|...
  entity_id    VARCHAR(255),
  old_values   JSONB,
  new_values   JSONB,
  ip_address   VARCHAR(45),
  user_agent   VARCHAR(500),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_tenant ON audit_logs(tenant_id, created_at DESC);
CREATE INDEX idx_audit_logs_user ON audit_logs(tenant_id, user_id, created_at DESC);
```

### A2. Seed de permisos y roles por defecto

```sql
-- permissions catalog
INSERT INTO permissions (key, description, category) VALUES
  ('view_data',               'Ver todos los datos del tenant',           'data'),
  ('create_transaction',      'Crear gastos/ingresos',                    'data'),
  ('edit_any_transaction',    'Editar cualquier transacción',             'data'),
  ('delete_any_transaction',  'Eliminar cualquier transacción',           'data'),
  ('manage_budgets',          'Gestionar presupuestos',                   'data'),
  ('manage_savings',          'Gestionar metas de ahorro',                'data'),
  ('manage_recurring',        'Gestionar transacciones recurrentes',      'data'),
  ('invite_members',          'Generar códigos de invitación',            'member_management'),
  ('manage_roles',            'Cambiar roles de miembros',               'member_management'),
  ('remove_members',          'Remover miembros del tenant',             'member_management'),
  ('view_audit_logs',         'Ver logs de auditoría',                   'admin'),
  ('manage_tenant',           'Editar configuración del tenant',         'settings'),
  ('delete_tenant',           'Eliminar el tenant',                      'admin'),
  ('transfer_ownership',      'Transferir ownership del tenant',         'admin');

-- Default permissions seeded per tenant on creation (via function or trigger)
-- owner gets ALL permissions
-- admin gets all except delete_tenant and transfer_ownership
-- member gets: view_data, create_transaction, edit_any_transaction, delete_any_transaction,
--              manage_budgets, manage_savings, manage_recurring
-- viewer gets: view_data only
```

### A3. Agregar tenant_id a tablas existentes

Tablas que reciben `tenant_id` (sin `user_preferences`, `user_notification_settings`, `user_two_fa` — permanecen user-personal):

```sql
-- Para cada una de las siguientes tablas:
ALTER TABLE expenses              ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE incomes               ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE categories            ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE budgets               ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE savings_goals         ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE savings_transactions  ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE recurring_transactions ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE user_gamification     ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE achievements          ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);
ALTER TABLE user_actions          ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(50);

-- Índices de tenant_id
CREATE INDEX IF NOT EXISTS idx_expenses_tenant             ON expenses(tenant_id, transaction_date DESC);
CREATE INDEX IF NOT EXISTS idx_incomes_tenant              ON incomes(tenant_id, received_date DESC);
CREATE INDEX IF NOT EXISTS idx_categories_tenant           ON categories(tenant_id);
CREATE INDEX IF NOT EXISTS idx_budgets_tenant              ON budgets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_savings_goals_tenant        ON savings_goals(tenant_id);
CREATE INDEX IF NOT EXISTS idx_savings_transactions_tenant ON savings_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_recurring_tenant            ON recurring_transactions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_gamification_tenant         ON user_gamification(tenant_id);
CREATE INDEX IF NOT EXISTS idx_achievements_tenant         ON achievements(tenant_id);
CREATE INDEX IF NOT EXISTS idx_user_actions_tenant         ON user_actions(tenant_id);
```

### A4. Migración de datos existentes

```sql
-- 1. Crear tenant personal para cada usuario
INSERT INTO tenants (id, name, slug, owner_id, plan, created_at, updated_at)
SELECT
  'tnt_' || LEFT(REPLACE(id::text, '-', ''), 8),
  first_name || ' ' || last_name || ' (personal)',
  LEFT(LOWER(REPLACE(email, '@', '_at_')), 80),
  id,
  'free',
  NOW(), NOW()
FROM users
WHERE deleted_at IS NULL;

-- 2. Registrar a cada usuario como owner de su tenant personal
INSERT INTO tenant_members (id, tenant_id, user_id, role, joined_at, created_at)
SELECT
  'tmb_' || LEFT(REPLACE(gen_random_uuid()::text, '-', ''), 8),
  'tnt_' || LEFT(REPLACE(u.id::text, '-', ''), 8),
  u.id,
  'owner',
  NOW(), NOW()
FROM users u
WHERE u.deleted_at IS NULL;

-- 3. Poblar tenant_id en todas las tablas (cada registro → tenant del user que lo creó)
UPDATE expenses e SET tenant_id = 'tnt_' || LEFT(REPLACE(e.user_id::text, '-', ''), 8)
  WHERE e.tenant_id IS NULL;
-- ... mismo para incomes, categories, budgets, savings_goals, savings_transactions,
-- recurring_transactions, user_gamification, achievements, user_actions

-- 4. Seed role_permissions default para cada tenant existente
-- (via stored procedure o script Go)
```

### A5. Constraints NOT NULL (después de poblar)

```sql
ALTER TABLE expenses               ALTER COLUMN tenant_id SET NOT NULL;
ALTER TABLE incomes                ALTER COLUMN tenant_id SET NOT NULL;
-- ... mismo para todas las tablas anteriores
```

---

## Fase B — Auth Layer

**Objetivo:** Extender JWT con tenant_id y role. Crear el middleware de permisos.

### B1. `internal/modules/auth/domain/requests.go`
**Cambio:** Agregar `TenantID` y `Role` a Claims struct.

```go
type Claims struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    TenantID  string `json:"tenant_id"`   // NUEVO
    Role      string `json:"role"`        // NUEVO: owner|admin|member|viewer
    TokenType string `json:"token_type"`
    jwt.RegisteredClaims
}
```

### B2. `internal/modules/auth/ports/services.go`
**Cambio:** Actualizar `JWTService` interface + agregar 2 nuevas interfaces.

```go
// JWTService — actualizar GenerateTokens
type JWTService interface {
    GenerateTokens(userID, email, tenantID, role string) (*domain.TokenPair, error)
    // ... resto sin cambios
}

// TenantCreator — nueva interface (implementada por tenants module)
type TenantCreator interface {
    CreatePersonalTenant(ctx context.Context, userID, email string) (tenantID string, err error)
}

// TenantMemberFinder — nueva interface (para Login)
type TenantMemberFinder interface {
    FindTenantByUserID(ctx context.Context, userID string) (tenantID, role string, err error)
}
```

### B3. `internal/modules/auth/services/jwt_service.go`
**Cambio:** `GenerateTokens` recibe `tenantID, role string` y los pasa a `generateToken`.

```go
func (j *jwtService) GenerateTokens(userID, email, tenantID, role string) (*domain.TokenPair, error) {
    // Pasa tenantID y role al generateToken
}

func (j *jwtService) generateToken(userID, email, tenantID, role, tokenType string, expiry time.Duration) (string, time.Time, error) {
    claims := domain.Claims{
        UserID:    userID,
        Email:     email,
        TenantID:  tenantID,  // NUEVO
        Role:      role,      // NUEVO
        TokenType: tokenType,
        // ...
    }
    // ...
}
```

### B4. `internal/modules/auth/services/auth_service.go`
**Cambios:**
1. Constructor recibe `TenantCreator` y `TenantMemberFinder` interfaces
2. `Register()` — después de crear user, llama `tenantCreator.CreatePersonalTenant()` y usa el tenantID en GenerateTokens
3. `Login()` — después de validar credenciales, llama `tenantMemberFinder.FindTenantByUserID()` y usa tenantID+role en GenerateTokens
4. `RefreshToken()` — preserva tenant_id y role del token anterior (no re-query)

```go
// NewAuthService — nueva firma
func NewAuthService(
    userRepo ports.UserRepository,
    prefsRepo ports.PreferencesRepository,
    notifRepo ports.NotificationSettingsRepository,
    twoFARepo ports.TwoFARepository,
    jwtSvc ports.JWTService,
    pwSvc ports.PasswordService,
    twoFASvc ports.TwoFAService,
    tenantCreator ports.TenantCreator,        // NUEVO
    tenantFinder ports.TenantMemberFinder,    // NUEVO
    eventBus ports.EventBus,
    logger zerolog.Logger,
    maxLoginAttempts int,
    lockoutDuration time.Duration,
) *AuthService
```

### B5. `internal/modules/auth/module.go`
**Cambio:** Constructor `New()` recibe `TenantCreator` y `TenantMemberFinder`.

```go
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig,
         eventBus ports.EventBus,
         tenantCreator authports.TenantCreator,
         tenantFinder authports.TenantMemberFinder) *Module
```

### B6. `internal/infrastructure/middleware/auth.go`
**Cambio:** `RequireAuth()` extrae y setea `tenant_id` y `role` además de `user_id`.

```go
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...validate token...
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("tenant_id", claims.TenantID)  // NUEVO
        c.Set("role", claims.Role)           // NUEVO
        c.Set("token", token)
        c.Next()
    }
}
```

### B7. `internal/infrastructure/middleware/permission.go` (NUEVO)

```go
package middleware

import "net/http"
import "github.com/gin-gonic/gin"

type PermissionChecker interface {
    HasPermission(ctx context.Context, tenantID, role, permission string) (bool, error)
}

type PermissionMiddleware struct {
    checker PermissionChecker
    logger  zerolog.Logger
}

func NewPermissionMiddleware(checker PermissionChecker, logger zerolog.Logger) *PermissionMiddleware

func (m *PermissionMiddleware) Require(permission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tenantID := c.GetString("tenant_id")
        role := c.GetString("role")
        hasPermission, _ := m.checker.HasPermission(c.Request.Context(), tenantID, role, permission)
        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## Fase C — Módulo Tenants (nuevo)

**Path:** `internal/modules/tenants/`

### Estructura de archivos

```
internal/modules/tenants/
  domain/
    tenant.go              — Tenant struct + builder
    tenant_member.go       — TenantMember, Role constants
    invitation.go          — TenantInvitation + GenerateCode()
    permission.go          — Permission, DefaultPermissions map
    audit_log.go           — AuditLog, AuditAction constants
  ports/
    ports.go               — interfaces: TenantRepository, MemberRepository,
                             InvitationRepository, PermissionRepository,
                             AuditRepository, TenantService
  repository/
    gorm_repository.go     — implementa TODOS los ports del módulo
  services/
    tenant_service.go      — lógica de negocio (CreateTenant, SeedPermissions, ...)
  handlers/
    tenant_handler.go      — CRUD del tenant
    member_handler.go      — gestión de miembros
    invitation_handler.go  — códigos de invitación
    audit_handler.go       — consulta de audit logs
  module.go                — New(), RegisterRoutes(), RegisterSubscribers(),
                             PermissionMiddleware(), TenantCreator(), TenantMemberFinder()
```

### C1. Domain structs

**`domain/tenant.go`**
```go
type Tenant struct {
    ID        string       // "tnt_" + 8 chars
    Name      string
    Slug      string       // unique, URL-friendly
    OwnerID   string
    IsActive  bool
    Plan      string       // "free" | "premium"
    Settings  interface{}  // JSON
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

**`domain/tenant_member.go`**
```go
const (
    RoleOwner  = "owner"
    RoleAdmin  = "admin"
    RoleMember = "member"
    RoleViewer = "viewer"
)

type TenantMember struct {
    ID        string
    TenantID  string
    UserID    string
    Role      string    // RoleOwner|RoleAdmin|RoleMember|RoleViewer
    InvitedBy *string
    JoinedAt  time.Time
    CreatedAt time.Time
}
```

**`domain/invitation.go`**
```go
type TenantInvitation struct {
    ID        string
    TenantID  string
    Code      string    // 8-char alphanumeric unique
    Role      string
    CreatedBy string
    ExpiresAt *time.Time
    MaxUses   int
    UsedCount int
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
}
func GenerateInviteCode() string  // e.g., "HOGAR-2026-XYZ" or 8-char random
```

**`domain/audit_log.go`**
```go
const (
    AuditCreate       = "CREATE"
    AuditUpdate       = "UPDATE"
    AuditDelete       = "DELETE"
    AuditLogin        = "LOGIN"
    AuditInvite       = "INVITE"
    AuditRoleChange   = "ROLE_CHANGE"
    AuditSettingChange = "SETTING_CHANGE"
)

type AuditLog struct {
    ID         string
    TenantID   string
    UserID     string
    Action     string
    EntityType string     // "expense"|"income"|"budget"|"member"|...
    EntityID   string
    OldValues  interface{} // JSON
    NewValues  interface{} // JSON
    IPAddress  string
    UserAgent  string
    CreatedAt  time.Time
}
```

### C2. Ports interfaces

**`ports/ports.go`**
```go
type TenantRepository interface {
    Create(ctx, tenant) error
    FindByID(ctx, id) (*Tenant, error)
    FindByOwnerID(ctx, ownerID) (*Tenant, error)
    Update(ctx, tenant) error
    Delete(ctx, id) error
}

type MemberRepository interface {
    AddMember(ctx, member) error
    FindByUserID(ctx, userID) (*TenantMember, error)    // para Login
    FindByTenantID(ctx, tenantID) ([]*TenantMember, error)
    UpdateRole(ctx, tenantID, userID, role) error
    Remove(ctx, tenantID, userID) error
}

type InvitationRepository interface {
    Create(ctx, inv) error
    FindByCode(ctx, code) (*TenantInvitation, error)
    FindByTenantID(ctx, tenantID) ([]*TenantInvitation, error)
    IncrementUsed(ctx, id) error
    Deactivate(ctx, id) error
}

type PermissionRepository interface {
    HasPermission(ctx, tenantID, role, permKey) (bool, error)
    GetRolePermissions(ctx, tenantID, role) ([]string, error)
    GetAllRolePermissions(ctx, tenantID) (map[string][]string, error)
    SetRolePermissions(ctx, tenantID, role string, permissions []string) error
    SeedDefaultPermissions(ctx, tenantID) error
}

type AuditRepository interface {
    Create(ctx, log) error
    FindByTenantID(ctx, tenantID string, filters AuditFilters) ([]*AuditLog, error)
}

type TenantService interface {
    CreatePersonalTenant(ctx, userID, email string) (tenantID string, err error)
    FindTenantByUserID(ctx, userID string) (tenantID, role string, err error)
    // ... más métodos de negocio
}
```

### C3. Rutas del módulo

**Registradas en `module.go`:**

```go
// Todas bajo /api/v1/tenants y protegidas con RequireAuth()
tenants := router.Group("/tenants")
tenants.Use(authMW.RequireAuth())
{
    tenants.GET("/current",                    h.tenant.GetCurrent)
    tenants.PUT("/current",                    permMW.Require("manage_tenant"), h.tenant.Update)
    tenants.DELETE("/current",                 permMW.Require("delete_tenant"), h.tenant.Delete)
    tenants.POST("/transfer-ownership",        permMW.Require("transfer_ownership"), h.tenant.TransferOwnership)

    // Miembros
    tenants.GET("/members",                    h.member.List)
    tenants.DELETE("/members/:userID",         permMW.Require("remove_members"), h.member.Remove)
    tenants.PATCH("/members/:userID/role",     permMW.Require("manage_roles"), h.member.ChangeRole)

    // Roles y permisos (RBAC dinámico)
    tenants.GET("/roles",                      h.member.GetRoles)
    tenants.PUT("/roles/:role/permissions",    permMW.Require("manage_roles"), h.member.SetRolePermissions)

    // Invitaciones
    tenants.POST("/invitations",               permMW.Require("invite_members"), h.invitation.Create)
    tenants.GET("/invitations",                permMW.Require("invite_members"), h.invitation.List)
    tenants.DELETE("/invitations/:id",         permMW.Require("invite_members"), h.invitation.Revoke)

    // Unirse (público — no requiere tenant en JWT, solo auth)
    tenants.POST("/join",                      h.invitation.Join)

    // Audit logs
    tenants.GET("/audit-logs",                 permMW.Require("view_audit_logs"), h.audit.List)
}
```

### C4. `module.go` expone para otros módulos

```go
func (m *Module) PermissionMiddleware() *middleware.PermissionMiddleware
func (m *Module) TenantCreator() authports.TenantCreator    // usado por auth module
func (m *Module) TenantMemberFinder() authports.TenantMemberFinder  // usado por auth module
```

---

## Fase D — Agregar tenant_id a todos los módulos existentes

**Patrón para CADA módulo:**

### D.Pattern — cambios en domain struct
```go
type Expense struct {
    ID       string
    TenantID string  // NUEVO — siempre presente
    UserID   string  // mantener — indica el creador
    // ... resto igual
}
```

### D.Pattern — cambios en GORM model
```go
type ExpenseModel struct {
    ID       string `gorm:"primaryKey"`
    TenantID string `gorm:"column:tenant_id;not null;index"` // NUEVO
    UserID   string `gorm:"column:user_id;not null;index"`
    // ... resto igual
}
```

### D.Pattern — cambios en repository (queries)
```go
// List — cambia de user_id a tenant_id
func (r *ExpenseRepo) FindByTenantID(ctx, tenantID string, limit, offset int) ([]*Expense, error) {
    WHERE("tenant_id = ? AND deleted_at IS NULL", tenantID)
}

// FindByID — usa tenant_id para verificar pertenencia
func (r *ExpenseRepo) FindByID(ctx, id, tenantID string) (*Expense, error) {
    WHERE("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID)
}
```

### D.Pattern — cambios en handlers
```go
func (h *ExpenseHandler) Create(c *gin.Context) {
    userID   := c.GetString("user_id")
    tenantID := c.GetString("tenant_id")  // NUEVO — del JWT

    // Crear con tenant_id
    expense, err := domain.NewExpense(userID, tenantID, ...)

    // Audit
    h.audit.Log(ctx, AuditEntry{TenantID: tenantID, UserID: userID,
                                Action: "CREATE", EntityType: "expense", ...})
}

func (h *ExpenseHandler) List(c *gin.Context) {
    tenantID := c.GetString("tenant_id")
    // Lista POR tenant (todos los miembros ven todos los gastos)
    expenses, err := h.repo.FindByTenantID(ctx, tenantID, limit, offset)
}

func (h *ExpenseHandler) GetByID(c *gin.Context) {
    tenantID := c.GetString("tenant_id")
    expense, err := h.repo.FindByID(ctx, id, tenantID)  // tenant check embebido
    // Ya no verifica user_id — cualquier miembro puede ver
}

func (h *ExpenseHandler) Delete(c *gin.Context) {
    // Permission check ya hecho por middleware (delete_any_transaction)
    tenantID := c.GetString("tenant_id")
    expense, err := h.repo.FindByID(ctx, id, tenantID)
    // No verifica user_id — el permiso del rol ya lo garantizó el middleware
}
```

### D.Pattern — routes con permisos
```go
expenses := router.Group("/expenses")
expenses.Use(authMW.RequireAuth())
{
    expenses.GET("",     permMW.Require("view_data"), h.expense.List)
    expenses.GET("/:id", permMW.Require("view_data"), h.expense.GetByID)
    expenses.POST("",    permMW.Require("create_transaction"), h.expense.Create)
    expenses.PUT("/:id", permMW.Require("edit_any_transaction"), h.expense.Update)
    expenses.DELETE("/:id", permMW.Require("delete_any_transaction"), h.expense.Delete)
}
```

### Módulos afectados y archivos a modificar

| Módulo | Domain | GORM Model | Repository | Handler | Routes |
|---|---|---|---|---|---|
| **transactions** | expense.go, income.go, category.go | expense_repo.go, income_repo.go, category_repo.go | gorm_repository.go | expense_handler.go, income_handler.go, category_handler.go | module.go |
| **budgets** | budget.go | gorm_repository.go | gorm_repository.go | budget_handler.go | module.go |
| **savings** | savings_goal.go, savings_transaction.go | gorm_repository.go | gorm_repository.go | savings_handler.go | module.go |
| **recurring** | recurring_transaction.go | gorm_repository.go | gorm_repository.go | recurring_handler.go | module.go |
| **gamification** | (user_gamification, achievement, user_action) | gorm_repository.go | gorm_repository.go | handler.go | module.go |
| **analytics** | — | — | gorm_repository.go | handler.go | module.go |
| **ai** | — | — | — | handler.go | module.go |

**Nota sobre gamification:** user_gamification y achievements siguen siendo por user_id (XP personal), pero se agrega tenant_id para contexto. user_actions también recibe tenant_id.

**Nota sobre auth:** user_preferences, user_notification_settings, user_two_fa NO reciben tenant_id — son datos estrictamente personales del usuario.

---

## Fase E — Audit Trail

**Objetivo:** Registrar acciones importantes async via event bus.

### E1. `internal/shared/ports/audit.go` (nuevo)
```go
type AuditEntry struct {
    TenantID   string
    UserID     string
    Action     string
    EntityType string
    EntityID   string
    OldValues  interface{}
    NewValues  interface{}
    IPAddress  string
    UserAgent  string
}

type AuditLogger interface {
    Log(ctx context.Context, entry AuditEntry)  // non-blocking, async
}
```

### E2. `internal/modules/tenants/services/audit_service.go`
Implementación de `AuditLogger`:
- Recibe entries
- Publica evento `audit.action` al event bus (async)
- El subscriber escribe en `audit_logs`

### E3. Inyección en módulos
Los módulos que modifican datos reciben `AuditLogger` en su constructor:
```go
// transactions/module.go New() — agregar auditLogger
func New(db, logger, cfg, eventBus, authMW, permMW, auditLogger)
```

Handlers llaman:
```go
h.audit.Log(c.Request.Context(), AuditEntry{
    TenantID:   tenantID,
    UserID:     userID,
    Action:     "CREATE",
    EntityType: "expense",
    EntityID:   expense.ID,
    NewValues:  expense,
    IPAddress:  c.ClientIP(),
    UserAgent:  c.GetHeader("User-Agent"),
})
```

---

## Fase F — PostgreSQL RLS (defensa en profundidad)

**Objetivo:** Agregar policies a nivel DB como segunda capa de seguridad.

### F1. Habilitar RLS en tablas tenant-scoped
```sql
ALTER TABLE expenses ENABLE ROW LEVEL SECURITY;
CREATE POLICY expenses_tenant_isolation ON expenses
  USING (tenant_id = current_setting('app.current_tenant', TRUE));
-- ... mismo para todas las tablas tenant-scoped
```

### F2. GORM middleware para setear session variable
```go
// internal/infrastructure/database/tenant_scope.go
func WithTenantScope(db *gorm.DB, tenantID string) *gorm.DB {
    db.Exec("SET LOCAL app.current_tenant = ?", tenantID)
    return db
}
```

**Nota:** Requiere usar transacciones para que `SET LOCAL` sea efectivo. Es una capa adicional de seguridad — el app layer ya garantiza el aislamiento.

---

## Fase G — Frontend

### G1. Nuevas páginas

**`src/pages/TenantSettings.jsx`**
- Info del tenant (nombre, slug)
- Lista de miembros con roles
- Generar/revocar códigos de invitación
- Editar permisos por rol (RBAC dinámico)

**`src/pages/AuditLogs.jsx`**
- Tabla de audit logs con filtros (acción, usuario, entidad, fecha)
- Solo visible para admin+

### G2. Archivos modificados

| Archivo | Cambio |
|---|---|
| `src/contexts/AuthContext.js` | Almacenar `tenant_id` y `role` del JWT |
| `src/contexts/TenantContext.js` | **NUEVO** — info del tenant, members |
| `src/services/authService.js` | Leer `tenant_id`/`role` del token response |
| `src/services/tenantService.js` | **NUEVO** — API calls para tenant management |
| `src/components/Layout/Sidebar.jsx` | Agregar "Configuración del Tenant" (admin+) |
| `src/components/RoleGuard.jsx` | **NUEVO** — similar a FeatureGuard, chequea role |
| `src/pages/Dashboard.jsx` | Mostrar members widget si admin+ |

### G3. Role-based UI pattern
```jsx
// Ocultar botones de crear/editar para Viewers
<RoleGuard roles={['owner', 'admin', 'member']}>
  <button onClick={onCreateExpense}>+ Nuevo Gasto</button>
</RoleGuard>

// Solo admins ven audit logs en el sidebar
<RoleGuard roles={['owner', 'admin']}>
  <SidebarLink to="/audit-logs">Auditoría</SidebarLink>
</RoleGuard>
```

---

## Fase H — main.go (wiring final)

```go
// Orden de inicialización en main.go

// 1. Tenants module (se inicializa ANTES que auth — provee TenantCreator y TenantMemberFinder)
tenantsModule := tenants.New(db, logger, cfg, eventBus)
tenantCreator := tenantsModule.TenantCreator()
tenantFinder  := tenantsModule.TenantMemberFinder()

// 2. Auth module — recibe las interfaces de tenants
authModule := auth.New(db, logger, cfg, eventBus, tenantCreator, tenantFinder)
authModule.RegisterRoutes(apiV1)
authModule.RegisterSubscribers(eventBus)
authMW := authModule.AuthMiddleware()

// 3. Permission middleware (del tenants module)
permMW := tenantsModule.PermissionMiddleware()

// 4. Audit logger
auditLogger := tenantsModule.AuditLogger()

// 5. Resto de módulos (ahora reciben permMW y auditLogger además de authMW)
txModule := transactions.New(db, logger, cfg, eventBus, authMW, permMW, auditLogger)
// ... etc

// 6. Registrar rutas de tenants (al final, cuando authMW está disponible)
tenantsModule.SetAuthMiddleware(authMW)
tenantsModule.RegisterRoutes(apiV1)
tenantsModule.RegisterSubscribers(eventBus)
```

---

## Orden de implementación recomendado

```
1. [A1-A2] Schema SQL: nuevas tablas + seed de permisos
2. [A3-A5] Migración de datos: tenant_id a tablas existentes
3. [C1-C2] Tenants module: domain structs + ports (sin handlers aún)
4. [C3]    Tenants repository: implementar MemberRepository (para Login)
5. [B1-B4] Auth layer: Claims + JWTService + AuthService (Register + Login con tenant)
6. [B5-B7] Middleware: RequireAuth extendido + PermissionMiddleware
7. [B5]    Auth module.go: actualizar constructor
8. [C4]    Tenants module: completar service + handlers + module.go
9. [D]     Transactions module: tenant_id + permisos en rutas
10. [D]    Budgets module: idem
11. [D]    Savings module: idem
12. [D]    Recurring module: idem
13. [D]    Gamification module: idem
14. [D]    Analytics + AI modules: idem
15. [E]    AuditLogger: service + inyección en handlers
16. [H]    main.go: rewiring completo
17. Build + tests: go build ./... + prueba manual
18. [F]    PostgreSQL RLS: policies (opcional post-validate)
19. [G]    Frontend: TenantSettings + AuditLogs + RoleGuard + context
```

---

## Resumen de archivos nuevos (Go)

```
internal/modules/tenants/
  domain/tenant.go
  domain/tenant_member.go
  domain/invitation.go
  domain/permission.go
  domain/audit_log.go
  ports/ports.go
  repository/gorm_repository.go
  services/tenant_service.go
  services/audit_service.go
  handlers/tenant_handler.go
  handlers/member_handler.go
  handlers/invitation_handler.go
  handlers/audit_handler.go
  module.go

internal/infrastructure/middleware/permission.go
internal/shared/ports/audit.go
```

## Resumen de archivos modificados (Go)

```
internal/modules/auth/domain/requests.go      — Claims: +TenantID, +Role
internal/modules/auth/ports/services.go       — JWTService + 2 nuevas interfaces
internal/modules/auth/services/jwt_service.go — GenerateTokens con tenantID+role
internal/modules/auth/services/auth_service.go — Register+Login con tenant
internal/modules/auth/module.go               — constructor extendido
internal/infrastructure/middleware/auth.go    — RequireAuth: +tenant_id, +role en context
cmd/server/main.go                            — wiring completo

internal/modules/transactions/domain/expense.go       — +TenantID
internal/modules/transactions/domain/income.go        — +TenantID
internal/modules/transactions/repository/*.go          — queries por tenant_id
internal/modules/transactions/handlers/expense_handler.go — tenant context
internal/modules/transactions/handlers/income_handler.go  — tenant context
internal/modules/transactions/handlers/category_handler.go — tenant context
internal/modules/transactions/module.go               — +permMW, +auditLogger

internal/modules/budgets/domain/budget.go     — +TenantID
internal/modules/budgets/repository/*.go      — queries por tenant_id
internal/modules/budgets/handlers/*.go        — tenant context
internal/modules/budgets/module.go            — +permMW, +auditLogger

internal/modules/savings/domain/savings_goal.go         — +TenantID
internal/modules/savings/domain/savings_transaction.go  — +TenantID
internal/modules/savings/repository/*.go                — queries por tenant_id
internal/modules/savings/handlers/*.go                  — tenant context
internal/modules/savings/module.go                      — +permMW, +auditLogger

internal/modules/recurring/domain/recurring_transaction.go — +TenantID
internal/modules/recurring/repository/*.go              — queries por tenant_id
internal/modules/recurring/handlers/*.go                — tenant context
internal/modules/recurring/module.go                    — +permMW, +auditLogger

internal/modules/gamification/repository/*.go  — queries +tenant_id
internal/modules/gamification/handlers/*.go    — tenant context
internal/modules/gamification/module.go        — +permMW
```
