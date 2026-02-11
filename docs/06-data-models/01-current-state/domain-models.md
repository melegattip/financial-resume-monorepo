# Modelos de Dominio - Estado Actual

**Última Actualización**: 2026-02-09
**Versión del Código**: master (commit 6fca155)
**Total de Modelos**: 22 modelos de dominio + 12 DTOs

---

## Resumen General

El proyecto Financial Resume implementa una arquitectura de **Clean Architecture con Domain-Driven Design (DDD)** distribuida en tres microservicios principales. Los modelos de dominio están implementados en Go y siguen patrones de diseño consistentes como **Builder Pattern**, **Factory Pattern** y **Strategy Pattern**.

### Distribución por Servicio

- **API Gateway**: 11 modelos de dominio (transacciones, presupuestos, metas, categorías)
- **Gamification Service**: 7 modelos de dominio (logros, desafíos, acciones de usuario)
- **Users Service**: 4 modelos de dominio (usuarios, preferencias, autenticación)

### Arquitectura de Dominio

- **Interfaces** para polimorfismo (`Transaction` interface)
- **Value Objects** implícitos (IDs con prefijos semánticos)
- **Aggregates** con lógica de negocio encapsulada
- **Domain Events** implícitos (logros desbloqueados, metas alcanzadas)
- **Domain Services** a través de Factories y Builders

---

## Modelos por Servicio

### API Gateway (Core Finance Domain)

#### 1. Transaction (Interface)

**Propósito**: Interfaz base para polimorfismo de transacciones financieras.

**Métodos de Contrato**:
```go
GetID() string
GetUserID() string
GetAmount() float64
GetDescription() string
GetCategoryID() string
GetCreatedAt() time.Time
GetUpdatedAt() time.Time
GetPercentage() float64
CalculatePercentage(totalIncome float64)
```

**Patrones**:
- **Interface Segregation Principle**: Contrato mínimo para todas las transacciones
- **Polimorfismo**: Permite tratar `Income` y `Expense` de forma uniforme
- **Template Method**: `CalculatePercentage` tiene comportamiento especializado por subtipo

**Implementaciones**: `Expense`, `Income`, `BaseTransaction`

---

#### 2. Expense (Aggregate Root)

**Propósito**: Representa un gasto financiero con soporte para pagos parciales y fecha de transacción.

**Campos**:
```go
ID              string    // ID único
UserID          string    // Propietario
Amount          float64   // Monto total
AmountPaid      float64   // Monto pagado
Description     string    // Descripción
CategoryID      *string   // Categoría (nullable)
Paid            bool      // Estado de pago
DueDate         time.Time // Fecha de vencimiento
TransactionDate time.Time // Fecha real de gasto
CreatedAt       time.Time
UpdatedAt       time.Time
Percentage      float64   // % sobre ingresos
```

**Métodos de Negocio**:
- `GetPendingAmount() float64`: Calcula monto pendiente
- `IsFullyPaid() bool`: Verifica si está completamente pagado
- `AddPayment(paymentAmount float64)`: Agrega un pago parcial
- `CalculatePercentage(totalIncome float64)`: Calcula % sobre ingresos

**Validaciones Implícitas**:
- `AmountPaid` no puede exceder `Amount`
- Actualiza `Paid = true` automáticamente al pago completo

**Patrones**:
- **Builder Pattern**: `ExpenseBuilder` con fluent interface
- **Factory Pattern**: `ExpenseFactory` para creación
- **Aggregate Root**: Gestiona invariantes de pago

**Problema Identificado**:
- ⚠️ **Inconsistencia**: `CategoryID` es `*string` (nullable) en Expense pero `string` en BaseTransaction

---

#### 3. Income (Aggregate Root)

**Propósito**: Representa un ingreso financiero.

**Campos**:
```go
ID          string
UserID      string
Amount      float64
Description string
CategoryID  *string   // Nullable
CreatedAt   time.Time
UpdatedAt   time.Time
Percentage  float64   // No se usa (siempre 0)
```

**Métodos de Negocio**:
- `CalculatePercentage(totalIncome float64)`: Siempre retorna 0 (no aplica para ingresos)

**Patrones**:
- **Builder Pattern**: `IncomeBuilder`
- **Factory Pattern**: `IncomeFactory`

**Problema Identificado**:
- ⚠️ **Inconsistencia semántica**: `Percentage` siempre es 0 pero está en el modelo

---

#### 4. Category (Aggregate Root)

**Propósito**: Categorización de transacciones financieras.

**Campos**:
```go
ID        string
Name      string
UserID    string
CreatedAt time.Time
UpdatedAt time.Time
```

**Validaciones**:
- `Validate() error`: Verifica que `Name` no esté vacío

**Errores de Dominio**:
```go
var ErrEmptyCategoryName = errors.New("category name cannot be empty")
```

**Patrones**:
- **Builder Pattern**: `CategoryBuilder`
- **Self-Validation**: Método `Validate()`
- **ID Generation**: Auto-genera ID con prefijo `cat_` + 8 caracteres UUID

---

#### 5. Budget (Aggregate Root)

**Propósito**: Límite de gasto por categoría en un período específico con alertas automáticas.

**Campos**:
```go
ID          string
UserID      string
CategoryID  string
Amount      float64      // Límite presupuestal
SpentAmount float64      // Gastado actual
Period      BudgetPeriod // monthly/weekly/yearly
PeriodStart time.Time
PeriodEnd   time.Time
AlertAt     float64      // 0.0-1.0 (default 0.80)
Status      BudgetStatus // on_track/warning/exceeded
IsActive    bool
CreatedAt   time.Time
UpdatedAt   time.Time
```

**Enums de Dominio**:
```go
// BudgetPeriod
const (
    BudgetPeriodMonthly = "monthly"
    BudgetPeriodWeekly  = "weekly"
    BudgetPeriodYearly  = "yearly"
)

// BudgetStatus
const (
    BudgetStatusOnTrack  = "on_track"  // <70%
    BudgetStatusWarning  = "warning"   // 70-99%
    BudgetStatusExceeded = "exceeded"  // ≥100%
)
```

**Métodos de Negocio**:
- `Validate() error`: Valida reglas de negocio
- `UpdateSpentAmount(spentAmount float64)`: Actualiza gasto y recalcula estado
- `GetSpentPercentage() float64`: Calcula % gastado
- `GetRemainingAmount() float64`: Calcula restante
- `IsAlertTriggered() bool`: Verifica si se debe alertar
- `IsInCurrentPeriod() bool`: Valida vigencia
- `ResetForNewPeriod()`: Reinicia para nuevo período

**Patrones**:
- **Builder Pattern**: `BudgetBuilder` con auto-cálculo de fechas
- **State Pattern**: `Status` cambia automáticamente según `SpentAmount`
- **Strategy Pattern Ready**: Diseñado para diferentes estrategias de período

**Lógica de Negocio Destacada**:
- **Auto-actualización de estado**: Actualiza estado automáticamente
- **Cálculo inteligente de períodos**: Calcula inicio/fin según tipo
- **Protección de invariantes**: No permite gastos negativos, AlertAt fuera de rango

---

#### 6. RecurringTransaction (Aggregate Root)

**Propósito**: Transacciones automáticas recurrentes (ingresos/gastos periódicos).

**Campos**:
```go
ID             string
UserID         string
Amount         float64
Description    string
CategoryID     *string
Type           string     // "income"/"expense"
Frequency      string     // daily/weekly/monthly/yearly
NextDate       time.Time
LastExecuted   *time.Time
IsActive       bool
AutoCreate     bool
NotifyBefore   int        // Días
EndDate        *time.Time
ExecutionCount int
MaxExecutions  *int
CreatedAt      time.Time
UpdatedAt      time.Time
```

**Métodos de Negocio**:
- `Validate() error`: Valida reglas complejas de recurrencia
- `CalculateNextDate() time.Time`: Calcula próxima ejecución según frecuencia
- `ShouldExecute() bool`: Verifica si debe ejecutarse ahora
- `ShouldNotify() bool`: Verifica si debe enviar notificación
- `Execute()`: Marca como ejecutado y calcula próxima fecha
- `Pause()` / `Resume()`: Control de estado
- `GetDaysUntilNext() int`: Días hasta próxima ejecución

**Patrones**:
- **Builder Pattern**: `RecurringTransactionBuilder`
- **State Machine**: Estados de ejecución y pausado
- **Strategy Pattern**: Diferentes frecuencias de ejecución

**Lógica de Negocio Destacada**:
- **Auto-desactivación**: Se desactiva al alcanzar EndDate o MaxExecutions
- **Notificaciones inteligentes**: Calcula según `NotifyBefore`
- **Scheduling automático**: Auto-calcula NextDate

---

#### 7. SavingsGoal (Aggregate Root)

**Propósito**: Meta de ahorro con tracking de progreso y auto-guardado.

**Campos**:
```go
ID                string
UserID            string
Name              string
Description       string
TargetAmount      float64
CurrentAmount     float64
Category          SavingsGoalCategory
Priority          SavingsGoalPriority
TargetDate        time.Time
Status            SavingsGoalStatus
MonthlyTarget     float64  // Calculado
WeeklyTarget      float64  // Calculado
DailyTarget       float64  // Calculado
Progress          float64  // 0-1
RemainingAmount   float64  // Calculado
DaysRemaining     int      // Calculado
IsAutoSave        bool
AutoSaveAmount    float64
AutoSaveFrequency string   // daily/weekly/monthly
ImageURL          string
CreatedAt         time.Time
UpdatedAt         time.Time
AchievedAt        *time.Time
```

**Enums de Dominio**:
```go
// SavingsGoalCategory
const (
    SavingsGoalCategoryVacation   = "vacation"
    SavingsGoalCategoryEmergency  = "emergency"
    SavingsGoalCategoryHouse      = "house"
    SavingsGoalCategoryCar        = "car"
    SavingsGoalCategoryEducation  = "education"
    SavingsGoalCategoryRetirement = "retirement"
    SavingsGoalCategoryInvestment = "investment"
    SavingsGoalCategoryOther      = "other"
)

// SavingsGoalStatus
const (
    SavingsGoalStatusActive    = "active"
    SavingsGoalStatusAchieved  = "achieved"
    SavingsGoalStatusPaused    = "paused"
    SavingsGoalStatusCancelled = "cancelled"
)

// SavingsGoalPriority
const (
    SavingsGoalPriorityHigh   = "high"
    SavingsGoalPriorityMedium = "medium"
    SavingsGoalPriorityLow    = "low"
)
```

**Métodos de Negocio**:
- `Validate() error`: Reglas de negocio complejas
- `AddSavings(amount float64) error`: Agrega ahorros con auto-achieve
- `WithdrawSavings(amount float64) error`: Retira ahorros
- `UpdateCalculatedFields()`: Recalcula todos los campos derivados
- `GetProgress() float64`: Progreso 0.0-1.0
- `IsOverdue() bool`: Verifica si está vencida
- `IsOnTrack() bool`: Verifica si va según plan (±10%)
- `Pause()` / `Resume()` / `Cancel()`: Gestión de estados

**Patrones**:
- **Builder Pattern**: `SavingsGoalBuilder` con auto-cálculo de targets
- **State Machine**: Active → Achieved/Paused/Cancelled
- **Computed Properties**: Campos calculados automáticamente

**Lógica de Negocio Destacada**:
- **Auto-achievement**: Marca como achieved automáticamente
- **Auto-reversión**: Revierte a active si ya no está achieved
- **Cálculo dinámico de targets**: Recalcula daily/weekly/monthly según progreso
- **Tracking de progreso**: Compara progreso real vs esperado

---

### Gamification Service

#### 8. UserGamification (Aggregate Root)

**Propósito**: Estado de gamificación y progreso del usuario.

**Campos**:
```go
ID                string
UserID            string
TotalXP           int
CurrentLevel      int
InsightsViewed    int
ActionsCompleted  int
AchievementsCount int
CurrentStreak     int
LastActivity      time.Time
CreatedAt         time.Time
UpdatedAt         time.Time
```

**Métodos de Negocio**:
- `CalculateLevel() int`: Calcula nivel basado en XP (niveles 1-10)
- `XPToNextLevel() int`: XP necesario para subir de nivel
- `ProgressToNextLevel() int`: Porcentaje 0-100 de progreso
- `GetLevelName() string`: Nombre user-friendly del nivel

**Niveles de Progresión**:
```
Nivel 1:  0 XP    - "Financial Newbie"
Nivel 2:  75 XP   - "Money Tracker"
Nivel 3:  200 XP  - "Smart Saver" 🔓 METAS DE AHORRO
Nivel 4:  400 XP  - "Budget Master"
Nivel 5:  700 XP  - "Financial Planner" 🔓 PRESUPUESTOS
Nivel 6:  1200 XP - "Investment Seeker"
Nivel 7:  1800 XP - "Wealth Builder" 🔓 IA FINANCIERA
Nivel 8:  2600 XP - "Financial Strategist"
Nivel 9:  3600 XP - "Money Mentor"
Nivel 10: 5500 XP - "Financial Magnate"
```

**Patrones**:
- **Progression System**: Unlocks de features según nivel
- **Computed Properties**: CurrentLevel se calcula, no se almacena

---

#### 9. Achievement (Aggregate Root)

**Propósito**: Logros desbloqueables con tracking de progreso.

**Campos**:
```go
ID          string
UserID      string
Type        string     // ai_partner, action_taker, etc.
Name        string     // "🤖 AI Partner"
Description string     // "100 insights de IA utilizados"
Points      int        // XP otorgados al desbloquear
Progress    int        // Actual (ej: 67)
Target      int        // Objetivo (ej: 100)
Completed   bool
UnlockedAt  *time.Time
CreatedAt   time.Time
UpdatedAt   time.Time
```

**Achievement Types**:
```go
const (
    AchievementTypeAIPartner     = "ai_partner"
    AchievementTypeActionTaker   = "action_taker"
    AchievementTypeDataExplorer  = "data_explorer"
    AchievementTypeQuickLearner  = "quick_learner"
    AchievementTypeInsightMaster = "insight_master"
    AchievementTypeStreakKeeper  = "streak_keeper"
)
```

**Métodos de Negocio**:
- `IsCompleted() bool`: Verifica si Progress ≥ Target
- `UpdateProgress(newProgress int)`: Actualiza y auto-completa si alcanza Target

**Patrones**:
- **Progress Tracking**: Sistema de progreso incremental
- **Auto-completion**: Se marca como completado automáticamente

---

#### 10. Challenge (Configuration Entity)

**Propósito**: Definición de desafíos disponibles (configuración).

**Campos**:
```go
ID                string
ChallengeKey      string  // Identificador único
Name              string
Description       string
ChallengeType     string  // daily/weekly/monthly
Icon              string
XPReward          int
RequirementType   string  // transaction_count, etc.
RequirementTarget int
RequirementData   JSONB   // Datos flexibles
Active            bool
CreatedAt         time.Time
UpdatedAt         time.Time
```

**Challenge Types**:
```go
const (
    ChallengeTypeDaily   = "daily"
    ChallengeTypeWeekly  = "weekly"
    ChallengeTypeMonthly = "monthly"
)
```

**Requirement Types**:
```go
const (
    RequirementTypeTransactionCount = "transaction_count"
    RequirementTypeCategoryVariety  = "category_variety"
    RequirementTypeViewCombo        = "view_combo"
    RequirementTypeDailyLogin       = "daily_login"
)
```

---

### Users Service

#### 11. User (Aggregate Root)

**Propósito**: Usuario del sistema con autenticación y seguridad.

**Campos**:
```go
ID                       uint
Email                    string
Password                 string  // Nunca en JSON
FirstName                string
LastName                 string
Phone                    string
Avatar                   string
IsActive                 bool
IsVerified               bool
EmailVerificationToken   string  // Nunca en JSON
PasswordResetToken       string  // Nunca en JSON
LastLogin                *time.Time
FailedLoginAttempts      int
LockedUntil              *time.Time
CreatedAt                time.Time
UpdatedAt                time.Time
```

**Patrones**:
- **Security First**: Tokens y passwords nunca se exponen en JSON
- **Account Lockout**: Sistema de intentos fallidos y bloqueo temporal
- **Email Verification**: Workflow completo de verificación

**Problema Identificado**:
- ⚠️ **Inconsistencia**: UserID es `uint` en Users Service pero `string` en API Gateway

---

## Patrones Identificados

### Patrones Consistentes (Buenas Prácticas)

#### ✅ Builder Pattern
**Uso**: En todos los agregados complejos
- `ExpenseBuilder`, `IncomeBuilder`, `CategoryBuilder`
- `BudgetBuilder`, `SavingsGoalBuilder`, `RecurringTransactionBuilder`
- Fluent Interface con return `*Builder`
- Auto-inicialización de timestamps

**Beneficio**: Construcción clara y validación en `Build()`

#### ✅ Factory Pattern
**Uso**: Creación de transacciones
- `ExpenseFactory`, `IncomeFactory`
- Implementan `TransactionFactory` interface
- Encapsulan inicialización de defaults

#### ✅ Self-Validation
**Uso**: En casi todos los agregados
- Método `Validate() error` público
- Validación de invariantes de negocio
- Errores de dominio tipados

#### ✅ Computed Properties
**Uso**: Campos derivados auto-calculados
- `Budget`: Status, SpentPercentage
- `SavingsGoal`: Progress, RemainingAmount, DailyTarget
- `UserGamification`: CurrentLevel, XPToNextLevel

#### ✅ State Machine Pattern
**Uso**: Gestión de estados
- `Budget`: OnTrack → Warning → Exceeded
- `SavingsGoal`: Active → Achieved/Paused/Cancelled
- `RecurringTransaction`: Active ↔ Paused

#### ✅ ID Generation Strategy
**Uso**: IDs semánticos con prefijos
- `exp_XXXXXXX` (Expense)
- `inc_XXXXXXX` (Income)
- `bud_XXXXXXXX` (Budget)
- `cat_XXXXXXXX` (Category)
- `goal_XXXXXXXX` (SavingsGoal)

---

### Inconsistencias Detectadas

#### ⚠️ 1. Tipo de UserID
**Problema**: Inconsistencia entre servicios
- API Gateway: `UserID string`
- Users Service: `ID uint`

**Impacto**: Requiere conversión manual en cada llamada inter-servicio

**Recomendación**: Estandarizar a `string` (UUID) en todos los servicios

---

#### ⚠️ 2. Longitud de IDs Generados
**Problema**: Diferentes longitudes en prefijos
- Expense/Income: 7 caracteres UUID
- Budget/Category/SavingsGoal: 8 caracteres UUID
- RecurringTransaction: UUID completo

**Recomendación**: Definir estándar único (ej: prefijo + 8 chars)

---

#### ⚠️ 3. CategoryID Nullable vs Non-Nullable
**Problema**: Inconsistencia en nullabilidad
- `Expense.CategoryID`: `*string` (nullable)
- `BaseTransaction.CategoryID`: `string` (non-nullable)

**Recomendación**: Estandarizar a nullable (`*string`) en todos los modelos

---

#### ⚠️ 4. Percentage en Income
**Problema**: Campo `Percentage` siempre es 0 en Income

**Recomendación**:
- Opción A: Remover `Percentage` de Income
- Opción B: Calcular % de contribution al total de ingresos

---

#### ⚠️ 5. Duplicación de Gamification Models
**Problema**: Gamification models duplicados entre API Gateway y Gamification Service

**Recomendación**: API Gateway debe usar cliente del servicio, no copiar modelos

---

#### ⚠️ 6. Domain Errors Dispersos
**Problema**: Errores en múltiples archivos
- `errors.go`, `transaction_errors.go`, `budget.go`, `savings_goal.go`

**Recomendación**: Centralizar en `domain/errors.go`

---

#### ⚠️ 7. DTOs Mezclados con Domain
**Problema**: `CreateIncomeRequest`, `UpdateIncomeRequest` en `/domain`

**Recomendación**: Mover a capa de puertos (`/ports` o `/adapters`)

---

## Recomendaciones para Backend v2

### 1. Estandarización Crítica (🔴 Alta Prioridad)

**A. Unificar Tipo de UserID**
```go
// Todos los servicios
type UserID string  // UUID format
```

**B. Estandarizar Generación de IDs**
```go
func NewEntityID(prefix string) string {
    return prefix + "_" + uuid.New().String()[:12]
}
```

**C. Centralizar Domain Errors**
```go
// domain/errors.go
var (
    ErrInvalidTransactionType = NewDomainError("INVALID_TRANSACTION_TYPE", "...")
    ErrInvalidBudgetAmount    = NewDomainError("INVALID_BUDGET_AMOUNT", "...")
)
```

---

### 2. Mejoras de Arquitectura (🟡 Media Prioridad)

**A. Implementar Domain Events**
```go
type DomainEvent interface {
    EventType() string
    AggregateID() string
    OccurredAt() time.Time
}
```

**B. Extraer Value Objects**
```go
type Money struct {
    Amount   float64
    Currency string
}

type DateRange struct {
    Start time.Time
    End   time.Time
}
```

**C. Repository Interfaces en Domain**
```go
type BudgetRepository interface {
    Save(budget *domain.Budget) error
    FindByID(id string) (*domain.Budget, error)
}
```

---

### 3. Limpieza de Código (🟢 Baja Prioridad)

**A. Eliminar Código Legacy**
- Modelo `FinancialData` deprecated

**B. Mover DTOs a Capa Correcta**
- `CreateIncomeRequest` → `ports/income_ports.go`

**C. Resolver Percentage en Income**
- Eliminar o dar sentido al campo

---

## Métricas de Calidad

### Complejidad de Modelos

| Modelo                  | Campos | Métodos | Validaciones | Score |
|-------------------------|--------|---------|--------------|-------|
| Expense                 | 11     | 7       | Builder      | ⭐⭐⭐⭐⭐ |
| Budget                  | 13     | 11      | 5 rules      | ⭐⭐⭐⭐⭐ |
| SavingsGoal             | 21     | 13      | 6 rules      | ⭐⭐⭐⭐⭐ |
| RecurringTransaction    | 15     | 11      | 8 rules      | ⭐⭐⭐⭐⭐ |
| Income                  | 8      | 3       | Builder      | ⭐⭐⭐⭐   |
| Category                | 5      | 2       | 1 rule       | ⭐⭐⭐⭐   |
| UserGamification        | 11     | 4       | None         | ⭐⭐⭐    |
| User                    | 14     | 0       | External     | ⭐⭐⭐    |

---

## Conclusiones

### Fortalezas del Diseño Actual

✅ **Patrones Consistentes**: Builder y Factory bien aplicados
✅ **Validación de Dominio**: Mayoría de modelos validan invariantes
✅ **Separación de Responsabilidades**: Clean Architecture respetada
✅ **Lógica de Negocio Encapsulada**: Comportamiento en aggregates
✅ **Testing**: Tests unitarios para modelos core

### Debilidades a Resolver

❌ **Inconsistencia de IDs**: UserID uint vs string
❌ **Falta de Domain Events**: Comunicación entre bounded contexts
❌ **Value Objects Ausentes**: Primitivos expuestos
❌ **Código Legacy**: Sin eliminar
❌ **DTOs en Domain**: Violación de arquitectura

### Score General de Calidad

**API Gateway**: ⭐⭐⭐⭐ (4/5) - Excelente con mejoras menores
**Gamification Service**: ⭐⭐⭐ (3/5) - Bueno pero falta validación
**Users Service**: ⭐⭐⭐ (3/5) - Funcional pero básico

**TOTAL**: ⭐⭐⭐⭐ (3.7/5) - **MUY BUENO con oportunidades claras de mejora**

---

**Documento generado por**: Agent: Domain Models Explorer
**Versión del código**: master (commit 6fca155)
**Total de modelos documentados**: 22 modelos de dominio + 12 DTOs
