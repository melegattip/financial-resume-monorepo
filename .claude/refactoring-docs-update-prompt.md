# 🔄 Prompt: Actualización de Documentación para Refactorización a Modular Monolith

## 📋 Contexto General

El proyecto **Financial Resume** está en producción con 1-10 usuarios beta. Arquitectura actual usa microservicios (api-gateway, users-service, gamification-service, ai-service) que se comunican vía HTTP, generando overhead innecesario para la escala actual.

**Decisión estratégica:** Refactorizar a **Modular Monolith** manteniendo Clean Architecture y DDD.

**Objetivo de este prompt:** Actualizar TODA la documentación en `docs/` para reflejar el plan de refactorización completo, sin perder ningún detalle.

---

## 🎯 Instrucciones Generales

**IMPORTANTE PARA TODOS LOS AGENTES:**

1. **CREAR ARCHIVOS DIRECTAMENTE** usando Write tool en las rutas especificadas
2. **NO generar outputs temporales** - escribir directo a `docs/`
3. **Mantener idioma:** Español para docs generales, Inglés para técnicos
4. **Validar consistencia** entre documentos
5. **Preservar contenido existente relevante**
6. **Formato:** Markdown con headers claros, tablas, diagramas mermaid

---

## 🔄 Fase 1: Arquitectura Actual (AS-IS)

### Agente: Architecture-Current-State

**Task:** Documentar la arquitectura de microservicios actual con todos sus problemas identificados.

**Archivo a crear:** `docs/03-architecture/01-current-state.md`

**Contenido requerido:**

1. **Diagrama de Arquitectura Actual**
   - Componentes: api-gateway, users-service, gamification-service, ai-service
   - Flujos de comunicación HTTP
   - Bases de datos: main-db, gamification-db
   - Deployment en Render.com

2. **Stack Tecnológico**
   - Backend: Go 1.23+
   - Frontend: React 18, TailwindCSS
   - DB: PostgreSQL 15
   - Infraestructura: Docker, Nginx

3. **Patrones Implementados**
   - Clean Architecture (capas: domain, usecases, infrastructure)
   - DDD (Aggregates, Value Objects, Entities)
   - Builder Pattern (Expense/Income construction)
   - Factory Pattern (Transaction creation)
   - State Machine (Budget states)

4. **Análisis de Comunicación Entre Servicios**
   - Gateway → Users Service (HTTP proxy)
   - Gateway → Gamification Service (HTTP proxy)
   - Gateway → AI Service (HTTP proxy)
   - Overhead medido: ~20-40ms por request
   - Serialización: JSON doble (Gateway + Service)

5. **Problemas Críticos Identificados**
   ```markdown
   ### 🔴 Problemas de Arquitectura

   #### 1. Over-engineering para Escala Actual
   - **Contexto:** 1-10 usuarios beta en producción
   - **Problema:** 4 servicios independientes con comunicación HTTP
   - **Impacto:** Latencia innecesaria, complejidad operacional

   #### 2. Duplicación de Datos
   - **Ubicación:** Tablas de gamificación en `main-db` y `gamification-db`
   - **Modelos afectados:** `user_gamification`, `achievements`, `challenges`
   - **Impacto:** Inconsistencias, complejidad en sincronización

   #### 3. IDs Inconsistentes
   - **Problema:** `UserID` es `uint` en algunos modelos, `string` en otros
   - **Ejemplos:**
     - `Expense.UserID`: uint
     - `UserGamification.UserID`: string
   - **Impacto:** Conversiones constantes, errores potenciales

   #### 4. Sin Soft Delete
   - **Problema:** Borrado físico en lugar de lógico
   - **Impacto:** Pérdida de datos históricos, imposibilidad de auditoría

   #### 5. Falta de Automatización
   - **Problema:** Transacciones recurrentes no se ejecutan automáticamente
   - **Impacto:** Usuario debe ejecutar manualmente, pobre UX

   #### 6. Sin Circuit Breakers
   - **Problema:** Llamadas HTTP sin protección contra fallos en cascada
   - **Impacto:** Un servicio caído puede bloquear toda la aplicación

   #### 7. Sin Cache
   - **Problema:** Cada request golpea la base de datos
   - **Impacto:** Dashboard y Analytics lentos en carga repetida

   #### 8. Transacciones No Atómicas
   - **Problema:** Operaciones entre servicios no pueden ser transaccionales
   - **Ejemplo:** Crear gasto + dar XP pueden quedar inconsistentes
   - **Impacto:** Datos inconsistentes, difícil de debuggear
   ```

6. **Métricas de Performance Actuales**
   - Latencia p95: No medida actualmente
   - Overhead de red: ~20-40ms por operación
   - Serialización JSON: ~5-10ms adicionales
   - Sin observabilidad implementada

**Fuentes de información:**
- `apps/api-gateway/internal/infrastructure/proxy/*.go`
- `apps/api-gateway/docs/10_OPTIMIZACION_ARQUITECTURA_ACTUAL.md`
- `docs/06-data-models/01-current-state/*.md`

**Formato:** Markdown con diagramas mermaid para arquitectura y flujos.

---

## 🎯 Fase 2: Arquitectura Objetivo (TO-BE)

### Agente: Architecture-Target-State

**Task:** Documentar la arquitectura de modular monolith objetivo con todos los beneficios y diseño detallado.

**Archivo a crear:** `docs/03-architecture/02-target-state.md`

**Contenido requerido:**

1. **Visión de Arquitectura Objetivo**
   - Modular Monolith con Clean Architecture
   - Un solo proceso Go
   - Módulos internos bien definidos
   - Base de datos única consolidada

2. **Diagrama de Arquitectura Objetivo**
   ```mermaid
   graph TB
       FE[React Frontend]

       subgraph Modular Monolith
           API[HTTP API Layer]

           subgraph Core Modules
               AUTH[Auth Module]
               TRANS[Transactions Module]
               BUDGET[Budgets Module]
               SAVINGS[Savings Goals Module]
               RECUR[Recurring Transactions Module]
               GAMIF[Gamification Module]
               AI[AI Analysis Module]
           end

           subgraph Shared
               DOMAIN[Shared Domain Models]
               PORTS[Ports & Interfaces]
               INFRA[Infrastructure Layer]
           end

           DB[(PostgreSQL - Unified DB)]
           REDIS[(Redis Cache)]
           EVENTS[Event Bus In-Memory]
       end

       FE -->|REST API| API
       API --> AUTH
       API --> TRANS
       API --> BUDGET
       API --> SAVINGS
       API --> RECUR

       TRANS -.->|events| GAMIF
       TRANS -.->|events| AI
       BUDGET -.->|events| GAMIF

       AUTH --> INFRA
       TRANS --> INFRA
       GAMIF --> INFRA
       AI --> INFRA

       INFRA --> DB
       INFRA --> REDIS
       INFRA --> EVENTS
   ```

3. **Estructura de Módulos**
   ```
   apps/monolith/
   ├── cmd/
   │   └── server/
   │       └── main.go                    # Entry point
   ├── internal/
   │   ├── modules/
   │   │   ├── auth/                      # Autenticación y usuarios
   │   │   │   ├── domain/
   │   │   │   ├── usecases/
   │   │   │   ├── handlers/
   │   │   │   └── repository/
   │   │   ├── transactions/              # Gastos e ingresos
   │   │   │   ├── domain/
   │   │   │   ├── usecases/
   │   │   │   ├── handlers/
   │   │   │   └── repository/
   │   │   ├── budgets/                   # Presupuestos
   │   │   ├── savings/                   # Metas de ahorro
   │   │   ├── recurring/                 # Transacciones recurrentes
   │   │   ├── gamification/              # Sistema de gamificación
   │   │   ├── ai/                        # Análisis con IA
   │   │   └── analytics/                 # Reportes y dashboard
   │   ├── shared/
   │   │   ├── domain/                    # Tipos compartidos
   │   │   ├── events/                    # Event bus
   │   │   ├── errors/                    # Error handling
   │   │   └── ports/                     # Interfaces compartidas
   │   └── infrastructure/
   │       ├── database/                  # PostgreSQL
   │       ├── cache/                     # Redis
   │       ├── http/                      # API handlers
   │       ├── middleware/                # Auth, CORS, etc.
   │       ├── cron/                      # Scheduled jobs
   │       └── config/                    # Configuration
   ├── migrations/                        # DB migrations consolidadas
   └── tests/
       ├── integration/
       └── e2e/
   ```

4. **Comunicación Entre Módulos**

   **Patrón 1: Llamadas Directas (Síncrono)**
   ```go
   // Para operaciones que requieren respuesta inmediata
   // Ejemplo: Crear gasto necesita validar presupuesto

   type TransactionUseCase struct {
       budgetService ports.BudgetService  // Interface
   }

   func (uc *TransactionUseCase) CreateExpense(ctx context.Context, expense domain.Expense) error {
       // Validar presupuesto (llamada directa en memoria)
       budget, err := uc.budgetService.GetBudgetForCategory(ctx, expense.CategoryID)
       if err != nil {
           return err
       }

       if expense.Amount > budget.RemainingAmount {
           return errors.New("budget exceeded")
       }

       // Crear gasto
       return uc.expenseRepo.Create(ctx, expense)
   }
   ```

   **Patrón 2: Event Bus (Asíncrono)**
   ```go
   // Para side effects que no requieren respuesta inmediata
   // Ejemplo: Dar XP cuando se crea un gasto

   type TransactionUseCase struct {
       eventBus ports.EventBus
   }

   func (uc *TransactionUseCase) CreateExpense(ctx context.Context, expense domain.Expense) error {
       // Crear gasto en transacción
       tx.Begin()
       err := uc.expenseRepo.Create(ctx, expense)
       tx.Commit()

       // Emitir evento (no bloquea, no falla si gamification está caído)
       uc.eventBus.Publish(events.ExpenseCreated{
           UserID: expense.UserID,
           Amount: expense.Amount,
       })

       return nil
   }

   // En otro módulo (Gamification):
   func (g *GamificationModule) OnExpenseCreated(event events.ExpenseCreated) {
       g.AddXP(event.UserID, 5)
   }
   ```

5. **Beneficios de la Arquitectura Objetivo**

   | Aspecto | Antes (Microservicios) | Después (Modular Monolith) | Mejora |
   |---------|------------------------|----------------------------|---------|
   | **Latencia interna** | 20-40ms por HTTP call | <0.1ms (memoria) | **200-400x** |
   | **Transacciones ACID** | ❌ Imposible entre servicios | ✅ Nativo PostgreSQL | Consistencia garantizada |
   | **Deployment** | 4 servicios, 4 repos | 1 binario, 1 repo | 75% menos complejidad |
   | **Debugging** | Logs distribuidos | Stack trace único | 10x más rápido |
   | **Costo infraestructura** | 4 instancias mínimo | 1 instancia | 75% ahorro |
   | **Escalabilidad** | Horizontal complejo | Vertical simple + horizontal cuando necesario | Gradual |
   | **Testing** | Mocks complejos | Tests de integración directos | 5x más rápido |
   | **Onboarding devs** | 4 repos, 4 arquitecturas | 1 repo, 1 arquitectura | 50% menos tiempo |

6. **Manejo de Escala Futura**
   ```markdown
   ### Estrategia de Escalabilidad

   **Escenario 1: 1-1000 usuarios (Actual → Corto plazo)**
   - Modular Monolith en 1 instancia
   - Vertical scaling (más CPU/RAM)
   - Redis para cache

   **Escenario 2: 1000-10,000 usuarios (Mediano plazo)**
   - Modular Monolith en 2-3 instancias con load balancer
   - PostgreSQL con read replicas
   - Redis cluster

   **Escenario 3: 10,000-100,000 usuarios (Largo plazo)**
   - Extraer módulos críticos a microservicios SI es necesario
   - Candidatos: AI Service (cómputo intensivo), Analytics (lectura pesada)
   - Mantener módulos de negocio core en monolito

   **Punto clave:** La arquitectura modular permite extraer servicios CUANDO sea necesario, no por adelantado.
   ```

7. **Patrones y Principios**
   - Dependency Injection vía interfaces
   - Event-driven para comunicación asíncrona
   - CQRS ligero (separar reads de writes donde tenga sentido)
   - Repository Pattern para persistencia
   - Clean Architecture mantenida
   - DDD con bounded contexts claros

**Formato:** Markdown con diagramas mermaid extensos.

---

## 🗺️ Fase 3: Plan de Migración

### Agente: Migration-Planning

**Task:** Crear el plan detallado de migración de microservicios a modular monolith.

**IMPORTANTE:** Esta fase debe usar **spec-kit** para estructurar cada módulo de migración.

**Archivos a crear:**
1. `docs/03-architecture/03-migration-plan.md` (documento maestro)
2. `.claude/specs/refactoring/*.md` (specs individuales por módulo usando spec-kit)

**Contenido requerido:**

1. **Estrategia General**
   - Migración incremental, NO big bang
   - Sin downtime para usuarios
   - Validación en cada paso
   - **Uso de spec-kit para cada fase de migración**

2. **Pre-requisitos**
   ```markdown
   ### Checklist Pre-Migración

   - [ ] Backup completo de ambas bases de datos
   - [ ] Tests de integración end-to-end funcionando
   - [ ] Documentación de arquitectura actual completa
   - [ ] Feature flags implementados en frontend
   - [ ] Ambiente de staging idéntico a producción
   - [ ] Plan de rollback definido
   - [ ] Monitoreo y alertas configurados
   ```

3. **Fases de Migración Detalladas**

   **FASE 1: Setup Inicial (1 semana)**

   **📋 Usar spec-kit para esta fase:**
   ```bash
   # Crear spec para setup inicial
   /specify "Setup Monolith Foundation"
   # Descripción: Crear estructura base del monolito, configurar build, CI/CD y deployment

   # Generar plan de implementación
   /plan

   # Generar tasks accionables
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Crear estructura del monolito
   - Configurar build y deployment
   - Setup de CI/CD

   ### Spec-kit Artifacts
   - **Spec:** `.claude/specs/refactoring/01-setup-monolith.md`
   - **Plan:** `.claude/specs/refactoring/01-setup-monolith-plan.md`
   - **Tasks:** `.claude/specs/refactoring/01-setup-monolith-tasks.md`

   ### Tareas de Alto Nivel
   1. Crear estructura `apps/monolith/` según 02-target-state.md
   2. Configurar Go modules y dependencies
   3. Setup Dockerfile para monolito
   4. Configurar GitHub Actions para build
   5. Setup ambiente staging en Render.com

   **Ver tareas detalladas en:** `.claude/specs/refactoring/01-setup-monolith-tasks.md`

   ### Criterios de Éxito
   - ✅ Monolito compila sin errores
   - ✅ Health check endpoint responde
   - ✅ Deploy manual a staging funciona

   ### Archivos a crear
   - `apps/monolith/cmd/server/main.go`
   - `apps/monolith/internal/infrastructure/config/config.go`
   - `apps/monolith/Dockerfile`
   - `.github/workflows/monolith-deploy.yml`
   ```

   **FASE 2: Migración de Módulo Auth (1 semana)**

   **📋 Usar spec-kit para esta fase:**
   ```bash
   /specify "Migrate Auth Module to Monolith"
   # Descripción: Migrar users-service completo al módulo auth del monolito,
   # mantener API compatibility, implementar feature flags

   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Migrar users-service completo
   - Mantener compatibilidad de API
   - Validar autenticación funciona

   ### Spec-kit Artifacts
   - **Spec:** `.claude/specs/refactoring/02-migrate-auth-module.md`
   - **Plan:** `.claude/specs/refactoring/02-migrate-auth-module-plan.md`
   - **Tasks:** `.claude/specs/refactoring/02-migrate-auth-module-tasks.md`

   ### Tareas de Alto Nivel
   1. Copiar domain models de `apps/users-service/internal/domain/` a `apps/monolith/internal/modules/auth/domain/`
   2. Migrar usecases manteniendo interfaces
   3. Migrar handlers HTTP (Gin)
   4. Migrar repository (PostgreSQL)
   5. Configurar JWT middleware
   6. Implementar feature flag: `USE_MONOLITH_AUTH`

   **Ver tareas detalladas en:** `.claude/specs/refactoring/02-migrate-auth-module-tasks.md`

   ### Estrategia de Validación
   - Deploy a staging con feature flag OFF
   - Activar feature flag para 10% tráfico
   - Comparar respuestas monolito vs users-service (deben ser idénticas)
   - Si OK, subir a 50% → 100%
   - Si falla, rollback a 0%

   ### Tests Requeridos
   - [ ] Test: Register user
   - [ ] Test: Login
   - [ ] Test: Refresh token
   - [ ] Test: Get user profile
   - [ ] Test: Update user
   - [ ] Test: Password reset flow

   ### Rollback Plan
   - Cambiar feature flag a 0%
   - Verificar que users-service sigue funcionando
   - Investigar error en staging
   ```

   **FASE 3: Consolidación de Bases de Datos (2 semanas)**

   **📋 Usar spec-kit para esta fase CRÍTICA:**
   ```bash
   /specify "Consolidate Databases (main-db + gamification-db)"
   # Descripción: Unificar ambas bases de datos en una sola, migrar datos sin pérdida,
   # resolver duplicación de tablas, estandarizar UserIDs

   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Unificar main-db y gamification-db
   - Migrar datos sin pérdida
   - Resolver duplicación de tablas

   ### Spec-kit Artifacts
   - **Spec:** `.claude/specs/refactoring/03-consolidate-databases.md`
   - **Plan:** `.claude/specs/refactoring/03-consolidate-databases-plan.md`
   - **Tasks:** `.claude/specs/refactoring/03-consolidate-databases-tasks.md`

   ### Sub-Fase 3.1: Análisis de Datos
   - Auditar diferencias entre `main-db.user_gamification` y `gamification-db.user_gamification`
   - Identificar source of truth
   - Crear script de reconciliación

   **Ver tareas detalladas en:** `.claude/specs/refactoring/03-consolidate-databases-tasks.md`

   ### Sub-Fase 3.2: Migración de Esquema
   1. Crear migrations para consolidar:
      ```sql
      -- Migration: 001_consolidate_gamification.sql

      -- Agregar tablas de gamification a main-db si no existen
      CREATE TABLE IF NOT EXISTS user_gamification (
          user_id VARCHAR(255) PRIMARY KEY,
          current_level INT NOT NULL DEFAULT 0,
          current_xp INT NOT NULL DEFAULT 0,
          -- ... resto de campos
      );

      -- Migrar datos de gamification-db a main-db
      INSERT INTO main_db.user_gamification
      SELECT * FROM gamification_db.user_gamification
      ON CONFLICT (user_id) DO UPDATE SET
          current_xp = EXCLUDED.current_xp,
          -- resolver conflictos con lógica de negocio
      ;

      -- Verificar integridad
      SELECT COUNT(*) FROM main_db.user_gamification;
      ```

   2. Ejecutar migration en staging
   3. Validar datos migrados
   4. Ejecutar en producción con downtime mínimo

   ### Sub-Fase 3.3: Estandarización de IDs
   1. Crear migration para convertir `UserID uint` → `UserID string`:
      ```sql
      -- Migration: 002_standardize_user_ids.sql

      ALTER TABLE expenses
      ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::VARCHAR;

      ALTER TABLE incomes
      ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::VARCHAR;

      -- Repetir para todas las tablas con UserID uint
      ```

   2. Actualizar modelos Go:
      ```go
      // Antes:
      type Expense struct {
          UserID uint `json:"user_id"`
      }

      // Después:
      type Expense struct {
          UserID string `json:"user_id"`
      }
      ```

   ### Ventana de Mantenimiento
   - **Duración estimada:** 30 minutos
   - **Horario:** Domingo 2 AM (mínimo tráfico)
   - **Comunicación:** Email a usuarios 48 horas antes

   ### Tests Post-Migración
   - [ ] Verificar que todos los user_ids se migraron
   - [ ] Verificar que no hay duplicados
   - [ ] Test end-to-end de flujo completo
   - [ ] Verificar que gamification-db puede ser desconectado

   ### Rollback Plan
   - Restaurar backup de main-db
   - Reconectar gamification-db
   - Validar que servicios antiguos funcionan
   ```

   **FASE 4: Migración de Módulo Transactions (1.5 semanas)**

   **📋 Usar spec-kit:**
   ```bash
   /specify "Migrate Transactions Module (Expenses, Incomes, Categories)"
   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Migrar lógica de expenses e incomes
   - Migrar categories
   - Mantener API idéntica

   ### Spec-kit Artifacts
   - **Spec:** `.claude/specs/refactoring/04-migrate-transactions-module.md`
   - **Plan:** `.claude/specs/refactoring/04-migrate-transactions-module-plan.md`
   - **Tasks:** `.claude/specs/refactoring/04-migrate-transactions-module-tasks.md`

   ### Tareas de Alto Nivel
   1. Migrar domain models (Expense, Income, Category, builders, factories)
   2. Migrar usecases (create, update, delete, list)
   3. Migrar handlers HTTP
   4. Migrar repositories
   5. Feature flag: `USE_MONOLITH_TRANSACTIONS`

   **Ver tareas detalladas en:** `.claude/specs/refactoring/04-migrate-transactions-module-tasks.md`

   ### Integración con Auth Module
   ```go
   // El módulo transactions depende de auth para validar usuario
   type TransactionUseCase struct {
       authService ports.AuthService  // Interface compartida
   }
   ```

   ### Tests
   - [ ] Test: Create expense
   - [ ] Test: Create income
   - [ ] Test: Update expense inline
   - [ ] Test: Partial payment
   - [ ] Test: Filter and sort
   - [ ] Test: Delete transaction
   ```

   **FASE 5: Migración de Módulo Gamification (1.5 semanas)**

   **📋 Usar spec-kit (dos specs separados):**
   ```bash
   # Spec 1: Event Bus
   /specify "Implement In-Memory Event Bus"
   /plan
   /tasks

   # Spec 2: Gamification Module
   /specify "Migrate Gamification Module with Event-Driven Architecture"
   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Migrar lógica de XP, levels, achievements, challenges
   - Implementar event bus interno
   - Eliminar HTTP calls

   ### Spec-kit Artifacts
   - **Spec 1 (Event Bus):** `.claude/specs/refactoring/05a-implement-event-bus.md`
   - **Plan 1:** `.claude/specs/refactoring/05a-implement-event-bus-plan.md`
   - **Tasks 1:** `.claude/specs/refactoring/05a-implement-event-bus-tasks.md`
   - **Spec 2 (Gamification):** `.claude/specs/refactoring/05b-migrate-gamification-module.md`
   - **Plan 2:** `.claude/specs/refactoring/05b-migrate-gamification-module-plan.md`
   - **Tasks 2:** `.claude/specs/refactoring/05b-migrate-gamification-module-tasks.md`

   ### Tareas Críticas

   #### 1. Implementar Event Bus In-Memory
   **Ver:** `.claude/specs/refactoring/05a-implement-event-bus-tasks.md`
   ```go
   // apps/monolith/internal/shared/events/bus.go

   type EventBus interface {
       Publish(event Event) error
       Subscribe(eventType string, handler EventHandler)
   }

   type Event interface {
       Type() string
   }

   type EventHandler func(event Event) error

   // Implementación simple con channels
   type InMemoryEventBus struct {
       handlers map[string][]EventHandler
       mu       sync.RWMutex
   }
   ```

   #### 2. Definir Eventos
   ```go
   // apps/monolith/internal/shared/events/transaction_events.go

   type ExpenseCreated struct {
       UserID    string
       Amount    float64
       Category  string
       Timestamp time.Time
   }

   func (e ExpenseCreated) Type() string { return "expense.created" }
   ```

   #### 3. Publicar Eventos desde Transactions Module
   ```go
   // Antes (HTTP call):
   resp, err := httpClient.Post(gamificationURL + "/xp", xpData)

   // Después (event):
   uc.eventBus.Publish(events.ExpenseCreated{
       UserID: expense.UserID,
       Amount: expense.Amount,
   })
   ```

   #### 4. Suscribir en Gamification Module
   ```go
   // apps/monolith/internal/modules/gamification/handlers/event_handlers.go

   func (h *GamificationEventHandlers) OnExpenseCreated(event events.Event) error {
       expenseCreated := event.(events.ExpenseCreated)
       return h.gamificationService.AddXP(expenseCreated.UserID, 5)
   }
   ```

   ### Beneficio Clave
   - **Antes:** HTTP call ~20ms
   - **Después:** Event bus ~0.1ms
   - **Mejora:** 200x más rápido

   ### Tests
   - [ ] Test: Event published cuando se crea expense
   - [ ] Test: XP se agrega correctamente
   - [ ] Test: Achievements se desbloquean
   - [ ] Test: Challenges se completan
   - [ ] Test: Feature gate funciona
   ```

   **FASE 6: Migración de Módulos Restantes (2 semanas)**

   **📋 Usar spec-kit (un spec por módulo):**
   ```bash
   /specify "Migrate Budgets Module"
   /plan
   /tasks

   /specify "Migrate Savings Goals Module"
   /plan
   /tasks

   /specify "Migrate Recurring Transactions Module"
   /plan
   /tasks

   /specify "Migrate Analytics Module"
   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Módulos a migrar
   1. Budgets Module → `.claude/specs/refactoring/06a-migrate-budgets-module.md`
   2. Savings Goals Module → `.claude/specs/refactoring/06b-migrate-savings-module.md`
   3. Recurring Transactions Module → `.claude/specs/refactoring/06c-migrate-recurring-module.md`
   4. Analytics Module → `.claude/specs/refactoring/06d-migrate-analytics-module.md`
   5. AI Module (puede quedar como microservicio externo inicialmente)

   ### Enfoque
   - Seguir mismo patrón que fases anteriores
   - Migrar de menor a mayor complejidad
   - Feature flags para cada módulo
   - **Cada módulo tiene su spec, plan y tasks en spec-kit**
   ```

   **FASE 7: Implementación de Mejoras (2 semanas)**

   **📋 Usar spec-kit (un spec global para mejoras):**
   ```bash
   /specify "Technical Debt Resolution: Soft Delete, Cron Jobs, Redis Cache, Atomic Transactions"
   /plan
   /tasks
   ```

   **Referencia en docs/03-architecture/03-migration-plan.md:**
   ```markdown
   ### Objetivos
   - Resolver deuda técnica identificada
   - Implementar features bloqueadas por arquitectura anterior

   ### Spec-kit Artifacts
   - **Spec:** `.claude/specs/refactoring/07-technical-debt-resolution.md`
   - **Plan:** `.claude/specs/refactoring/07-technical-debt-resolution-plan.md`
   - **Tasks:** `.claude/specs/refactoring/07-technical-debt-resolution-tasks.md`

   ### Tareas de Alto Nivel

   #### 1. Soft Delete
   **Ver:** `.claude/specs/refactoring/07-technical-debt-resolution-tasks.md#soft-delete`
   ```sql
   -- Migration: 010_add_soft_delete.sql

   ALTER TABLE expenses ADD COLUMN deleted_at TIMESTAMP NULL;
   ALTER TABLE incomes ADD COLUMN deleted_at TIMESTAMP NULL;
   ALTER TABLE budgets ADD COLUMN deleted_at TIMESTAMP NULL;
   -- ... resto de tablas

   CREATE INDEX idx_expenses_deleted_at ON expenses(deleted_at);
   ```

   ```go
   // Actualizar queries
   // Antes:
   db.Find(&expenses)

   // Después:
   db.Where("deleted_at IS NULL").Find(&expenses)
   ```

   #### 2. Cron Jobs para Transacciones Recurrentes
   ```go
   // apps/monolith/internal/infrastructure/cron/scheduler.go

   func (s *Scheduler) Start() {
       // Ejecutar cada 1 hora
       c := cron.New()
       c.AddFunc("0 * * * *", s.processRecurringTransactions)
       c.Start()
   }

   func (s *Scheduler) processRecurringTransactions() {
       // Buscar transacciones pendientes
       pending := s.repo.GetPendingRecurringTransactions()

       for _, rt := range pending {
           // Crear transacción real
           s.transactionService.CreateFromRecurring(rt)

           // Actualizar next_execution_date
           s.repo.UpdateNextExecution(rt.ID, calculateNext(rt))
       }
   }
   ```

   #### 3. Redis Cache
   ```go
   // apps/monolith/internal/infrastructure/cache/redis.go

   type RedisCache struct {
       client *redis.Client
   }

   func (c *RedisCache) GetDashboard(userID string) (*Dashboard, error) {
       key := fmt.Sprintf("dashboard:%s", userID)
       data, err := c.client.Get(ctx, key).Result()
       // ...
   }
   ```

   #### 4. Transacciones Atómicas
   ```go
   // Ejemplo: Crear gasto + actualizar presupuesto + dar XP
   // TODO ATÓMICO

   func (uc *TransactionUseCase) CreateExpenseWithSideEffects(ctx context.Context, expense Expense) error {
       return uc.db.Transaction(func(tx *gorm.DB) error {
           // 1. Crear gasto
           if err := tx.Create(&expense).Error; err != nil {
               return err
           }

           // 2. Actualizar presupuesto
           budget, _ := uc.budgetRepo.GetByCategory(tx, expense.CategoryID)
           budget.Spent += expense.Amount
           if err := tx.Save(&budget).Error; err != nil {
               return err // Rollback automático
           }

           // 3. Dar XP
           gamification, _ := uc.gamificationRepo.GetByUser(tx, expense.UserID)
           gamification.XP += 5
           if err := tx.Save(&gamification).Error; err != nil {
               return err // Rollback automático
           }

           return nil // Commit
       })
   }
   ```

   ### Criterios de Éxito
   - ✅ Soft delete funciona en todas las entidades
   - ✅ Cron job ejecuta transacciones recurrentes automáticamente
   - ✅ Cache reduce latencia de dashboard en 80%
   - ✅ Transacciones atómicas garantizan consistencia
   ```

   **FASE 8: Testing y Validación (1 semana)**
   ```markdown
   ### Tests de Integración End-to-End

   #### Escenario 1: Flujo Completo de Usuario
   ```go
   func TestCompleteUserJourney(t *testing.T) {
       // 1. Register
       user := registerUser(t, "test@example.com", "password123")

       // 2. Login
       token := login(t, user.Email, user.Password)

       // 3. Create category
       category := createCategory(t, token, "Comida")

       // 4. Create budget
       budget := createBudget(t, token, category.ID, 50000)

       // 5. Create expense
       expense := createExpense(t, token, 10000, category.ID)

       // 6. Verify gamification (XP added)
       gamification := getGamification(t, token)
       assert.Equal(t, 5, gamification.XP)

       // 7. Verify budget updated
       updatedBudget := getBudget(t, token, budget.ID)
       assert.Equal(t, 10000.0, updatedBudget.Spent)
       assert.Equal(t, "on_track", updatedBudget.Status)

       // 8. Get dashboard
       dashboard := getDashboard(t, token)
       assert.Equal(t, 10000.0, dashboard.TotalExpenses)
   }
   ```

   #### Escenario 2: Performance bajo carga
   ```bash
   # k6 load test
   k6 run --vus 100 --duration 5m tests/load/complete-flow.js

   # Criterios de éxito:
   # - p95 latency < 200ms
   # - Error rate < 0.1%
   # - Throughput > 100 req/s
   ```

   #### Escenario 3: Rollback completo
   - Simular fallo crítico en monolito
   - Activar feature flags a 0%
   - Verificar que servicios antiguos funcionan
   - Tiempo de rollback < 5 minutos
   ```

   **FASE 9: Cutover a Producción (1 semana)**
   ```markdown
   ### Plan de Cutover

   #### Pre-Cutover (Viernes 6 PM)
   - [ ] Backup completo de DB
   - [ ] Deploy monolito a producción (feature flags OFF)
   - [ ] Verificar health checks
   - [ ] Notificar usuarios de posible mantenimiento

   #### Cutover Gradual (Sábado-Domingo)
   - **Sábado 12 AM:** Feature flags 10% para usuarios nuevos
   - **Sábado 6 AM:** Si OK, subir a 25%
   - **Sábado 12 PM:** Si OK, subir a 50%
   - **Sábado 6 PM:** Si OK, subir a 75%
   - **Domingo 12 AM:** Si OK, subir a 100%

   #### Monitoreo Continuo
   - Dashboard con métricas en tiempo real
   - Alertas configuradas
   - Equipo on-call disponible

   #### Post-Cutover (Lunes)
   - Validar métricas de negocio (transacciones creadas, usuarios activos)
   - Validar performance (latencia, throughput)
   - Validar errores (error rate, crashes)
   - Si todo OK durante 48 horas → Considerar éxito
   ```

   **FASE 10: Decommission y Limpieza (1 semana)**
   ```markdown
   ### Objetivos
   - Apagar microservicios antiguos
   - Eliminar código legacy
   - Actualizar documentación

   ### Tareas
   1. Apagar servicios en orden:
      - AI Service (último, puede ser externo)
      - Gamification Service
      - Users Service
      - API Gateway (cuando monolito sirva todo)

   2. Eliminar código:
      - `rm -rf apps/users-service/`
      - `rm -rf apps/gamification-service/`
      - `rm -rf apps/api-gateway/internal/infrastructure/proxy/`

   3. Eliminar gamification-db:
      - Backup final
      - Drop database
      - Actualizar connection strings

   4. Actualizar documentación:
      - Marcar docs de microservicios como obsoletos
      - Actualizar README.md
      - Actualizar arquitectura en Confluence/Wiki

   5. Celebrar 🎉
   ```

4. **Timeline y Estimaciones**
   ```markdown
   ## Timeline Total: 12 semanas (3 meses)

   | Fase | Duración | Dependencias | Riesgo |
   |------|----------|--------------|--------|
   | 1. Setup | 1 semana | Ninguna | 🟢 Bajo |
   | 2. Auth Migration | 1 semana | Fase 1 | 🟡 Medio |
   | 3. DB Consolidation | 2 semanas | Fase 2 | 🔴 Alto |
   | 4. Transactions | 1.5 semanas | Fase 3 | 🟡 Medio |
   | 5. Gamification | 1.5 semanas | Fase 4 | 🟡 Medio |
   | 6. Resto Módulos | 2 semanas | Fase 5 | 🟡 Medio |
   | 7. Mejoras | 2 semanas | Fase 6 | 🟢 Bajo |
   | 8. Testing | 1 semana | Fase 7 | 🟡 Medio |
   | 9. Cutover | 1 semana | Fase 8 | 🔴 Alto |
   | 10. Cleanup | 1 semana | Fase 9 | 🟢 Bajo |

   **Total:** 14 semanas (considerando buffer)
   ```

5. **Recursos Necesarios**
   ```markdown
   ## Team
   - 2 Backend Engineers (Go)
   - 1 DevOps Engineer
   - 1 QA Engineer
   - 1 Tech Lead (oversight)

   ## Infraestructura
   - Ambiente staging idéntico a producción
   - Redis instance
   - Monitoreo (Datadog/New Relic/Grafana)
   - Feature flag service (LaunchDarkly o custom)

   ## Budget
   - Infraestructura adicional: ~$200/mes durante migración
   - Herramientas: ~$100/mes
   - Total: ~$1000 para 3 meses
   ```

6. **Riesgos y Mitigaciones**
   ```markdown
   | Riesgo | Probabilidad | Impacto | Mitigación |
   |--------|--------------|---------|------------|
   | Pérdida de datos en consolidación DB | Media | 🔴 Crítico | Múltiples backups, testing exhaustivo en staging |
   | Bugs en producción post-migración | Alta | 🟡 Medio | Feature flags para rollback instantáneo |
   | Performance peor que esperada | Baja | 🟡 Medio | Benchmarks en staging, optimización antes de cutover |
   | Downtime durante cutover | Media | 🔴 Crítico | Cutover gradual con feature flags, no big bang |
   | Team burnout | Media | 🟡 Medio | Timeline realista con buffer, no horas extra forzadas |
   ```

**Formato:** Markdown con tablas, checklists, y diagramas de Gantt en mermaid.

---

## 🔧 Fase 4: Actualización de Requirements

### Agente: Requirements-Update

**Task:** Actualizar docs/01-requirements.md para incluir toda la deuda técnica como requisitos de refactorización.

**Archivo a editar:** `docs/01-requirements.md`

**Cambios requeridos:**

1. **Nueva sección al inicio del documento:**
   ```markdown
   ## ⚠️ IMPORTANTE: Contexto de Refactorización

   Este documento fue creado mediante ingeniería inversa del sistema actual (microservicios).
   El proyecto está en proceso de **refactorización a Modular Monolith**.

   Ver: [Plan de Migración](./03-architecture/03-migration-plan.md)

   Las secciones marcadas con 🔄 indican cambios derivados de la refactorización.
   ```

2. **Nueva sección de Requisitos de Refactorización:**
   ```markdown
   ## 🔄 Requisitos de Refactorización (RR)

   ### RR-001: Consolidación de Arquitectura
   **Prioridad:** Crítica
   **Categoría:** Arquitectura

   **Descripción:**
   El sistema debe migrar de arquitectura de microservicios a modular monolith manteniendo Clean Architecture y DDD.

   **Justificación:**
   - Escala actual (1-10 usuarios) no justifica complejidad de microservicios
   - Overhead de 20-40ms por request HTTP innecesario
   - Imposibilidad de transacciones ACID entre servicios
   - Complejidad operacional desproporcionada

   **Criterios de Aceptación:**
   - [ ] Un único proceso Go conteniendo todos los módulos
   - [ ] Comunicación entre módulos en memoria (<1ms)
   - [ ] Clean Architecture preservada por módulo
   - [ ] API externa mantiene compatibilidad 100%

   ---

   ### RR-002: Consolidación de Base de Datos
   **Prioridad:** Crítica
   **Categoría:** Base de Datos

   **Descripción:**
   Consolidar main-db y gamification-db en una única base de datos.

   **Justificación:**
   - Duplicación de tablas (`user_gamification` en ambas DBs)
   - Inconsistencias de datos
   - Imposibilidad de transacciones atómicas cross-DB

   **Criterios de Aceptación:**
   - [ ] Una única base de datos PostgreSQL
   - [ ] Cero pérdida de datos durante migración
   - [ ] Validación post-migración (checksums)
   - [ ] gamification-db puede ser desconectada

   ---

   ### RR-003: Estandarización de IDs
   **Prioridad:** Alta
   **Categoría:** Modelo de Datos

   **Descripción:**
   Estandarizar UserID como string (UUID) en todos los modelos.

   **Justificación:**
   - Inconsistencia actual: `UserID` es `uint` en algunos modelos, `string` en otros
   - Conversiones constantes propensas a errores
   - UUIDs son más escalables y seguros

   **Modelos Afectados:**
   - Expense, Income, Budget, SavingsGoal, RecurringTransaction

   **Criterios de Aceptación:**
   - [ ] Migration SQL para convertir uint → string
   - [ ] Todos los modelos Go usan `UserID string`
   - [ ] Tests de integración pasan
   - [ ] Rollback plan documentado

   ---

   ### RR-004: Implementación de Soft Delete
   **Prioridad:** Alta
   **Categoría:** Modelo de Datos

   **Descripción:**
   Implementar soft delete (borrado lógico) en todas las entidades.

   **Justificación:**
   - Actualmente borrado físico = pérdida de datos históricos
   - Imposibilidad de auditoría
   - Riesgo de borrado accidental

   **Criterios de Aceptación:**
   - [ ] Campo `deleted_at TIMESTAMP NULL` en todas las tablas
   - [ ] Queries filtran por `deleted_at IS NULL`
   - [ ] API de restore para admins
   - [ ] Job de limpieza para borrado físico después de X días

   ---

   ### RR-005: Cron Jobs para Transacciones Recurrentes
   **Prioridad:** Alta
   **Categoría:** Automatización

   **Descripción:**
   Implementar cron jobs que ejecuten transacciones recurrentes automáticamente.

   **Justificación:**
   - Actualmente usuario debe ejecutar manualmente
   - Mala UX
   - Feature prometida no funcional

   **Criterios de Aceptación:**
   - [ ] Cron job ejecuta cada 1 hora
   - [ ] Identifica transacciones vencidas (`next_execution_date < NOW()`)
   - [ ] Crea transacción real automáticamente
   - [ ] Actualiza `next_execution_date` según frecuencia
   - [ ] Envía notificación al usuario
   - [ ] Logging de ejecuciones

   ---

   ### RR-006: Event Bus Interno
   **Prioridad:** Alta
   **Categoría:** Arquitectura

   **Descripción:**
   Implementar event bus in-memory para comunicación asíncrona entre módulos.

   **Justificación:**
   - Reemplazar HTTP calls síncronos entre servicios
   - Desacoplar módulos (ej: Transactions no debe conocer Gamification)
   - Performance: <1ms vs ~20ms HTTP

   **Criterios de Aceptación:**
   - [ ] Event bus con pub/sub pattern
   - [ ] Handlers pueden fallar sin afectar operación principal
   - [ ] Logging de eventos emitidos
   - [ ] Tests unitarios de event handlers

   **Eventos Requeridos:**
   - `ExpenseCreated`, `IncomeCreated`
   - `BudgetExceeded`, `BudgetWarning`
   - `SavingsGoalAchieved`
   - `ChallengeCompleted`, `AchievementUnlocked`

   ---

   ### RR-007: Implementación de Cache
   **Prioridad:** Media
   **Categoría:** Performance

   **Descripción:**
   Implementar Redis cache para dashboard y analytics.

   **Justificación:**
   - Dashboard es la vista más consultada
   - Queries pesadas en cada request
   - Datos cambian poco (solo al crear transacciones)

   **Criterios de Aceptación:**
   - [ ] Redis configurado y conectado
   - [ ] Cache del dashboard (TTL 5 minutos)
   - [ ] Invalidación de cache al crear/actualizar transacciones
   - [ ] Cache hit ratio > 70%
   - [ ] Fallback a DB si Redis falla

   ---

   ### RR-008: Transacciones Atómicas
   **Prioridad:** Alta
   **Categoría:** Consistencia de Datos

   **Descripción:**
   Garantizar atomicidad en operaciones que afectan múltiples entidades.

   **Justificación:**
   - Actualmente: crear gasto + dar XP son 2 operaciones separadas
   - Posible inconsistencia si una falla

   **Ejemplos de Operaciones Atómicas:**
   1. Crear gasto + actualizar presupuesto + dar XP
   2. Completar meta de ahorro + desbloquear achievement
   3. Crear transacción recurrente + crear primera instancia

   **Criterios de Aceptación:**
   - [ ] Usar `db.Transaction()` de GORM
   - [ ] Rollback automático si cualquier operación falla
   - [ ] Tests de consistencia

   ---

   ### RR-009: Feature Flags
   **Prioridad:** Alta
   **Categoría:** DevOps

   **Descripción:**
   Implementar sistema de feature flags para migración gradual.

   **Justificación:**
   - Permitir rollback instantáneo sin redeploy
   - Canary releases (activar para % de usuarios)
   - A/B testing

   **Flags Requeridos:**
   - `USE_MONOLITH_AUTH`
   - `USE_MONOLITH_TRANSACTIONS`
   - `USE_MONOLITH_GAMIFICATION`
   - `USE_MONOLITH_BUDGETS`
   - (... uno por módulo)

   **Criterios de Aceptación:**
   - [ ] Flag service (LaunchDarkly o custom)
   - [ ] Flags configurables sin redeploy
   - [ ] Métricas por flag (% usuarios, error rate)
   - [ ] UI para controlar flags

   ---

   ### RR-010: Observabilidad
   **Prioridad:** Media
   **Categoría:** Monitoreo

   **Descripción:**
   Implementar tracing, métricas y logging estructurado.

   **Justificación:**
   - Actualmente no hay visibilidad de performance
   - Debugging es manual y lento
   - Necesario para validar mejoras de refactorización

   **Componentes:**
   1. **Tracing:** OpenTelemetry + Jaeger
   2. **Métricas:** Prometheus + Grafana
   3. **Logging:** Structured logs (JSON)

   **Métricas Clave:**
   - Request latency (p50, p95, p99)
   - Error rate por endpoint
   - Cache hit ratio
   - Active users
   - Transaction volume

   **Criterios de Aceptación:**
   - [ ] Traces visualizables en Jaeger
   - [ ] Dashboard en Grafana con métricas clave
   - [ ] Alertas configuradas (error rate > 1%, latency > 500ms)
   - [ ] Logs agregados y searcheables
   ```

3. **Actualizar sección de Requisitos No Funcionales:**
   Agregar estos RNFs:

   ```markdown
   ### RNF-066: Latencia Interna de Módulos
   **Categoría:** Performance
   **Descripción:** La comunicación entre módulos internos debe ser <1ms.
   **Medición:** Tracing de llamadas entre módulos.
   **Objetivo:** 200x mejora vs HTTP (~20ms).

   ### RNF-067: Atomicidad de Operaciones
   **Categoría:** Consistencia
   **Descripción:** Operaciones que afectan múltiples entidades deben ser atómicas.
   **Medición:** 0 inconsistencias detectadas en auditorías.

   ### RNF-068: Zero Data Loss en Migración
   **Categoría:** Confiabilidad
   **Descripción:** La migración de bases de datos debe preservar el 100% de los datos.
   **Medición:** Checksums pre y post migración deben coincidir.

   ### RNF-069: Rollback Time
   **Categoría:** Disponibilidad
   **Descripción:** Tiempo de rollback ante fallo debe ser <5 minutos.
   **Medición:** Simular fallo y medir tiempo hasta restauración.

   ### RNF-070: Deployment Simplificado
   **Categoría:** Operaciones
   **Descripción:** Deploy debe ser 1 binario, no N servicios.
   **Medición:** Número de pasos de deploy, tiempo total.
   ```

**Formato:** Mantener estructura existente de requirements.md, agregar nuevas secciones claramente identificadas.

---

## 📅 Fase 5: Roadmap de Implementación

### Agente: Roadmap-Creation

**Task:** Crear roadmap visual y detallado de la refactorización con hitos, dependencias y métricas de éxito.

**Archivo a crear:** `docs/03-architecture/04-implementation-roadmap.md`

**Contenido requerido:**

1. **Visión General**
   ```markdown
   # Roadmap de Refactorización: Microservicios → Modular Monolith

   ## Objetivo
   Migrar Financial Resume de arquitectura de microservicios a modular monolith en **12 semanas**, sin downtime, preservando funcionalidad y mejorando performance.

   ## Principios Guía
   1. **No Big Bang:** Migración incremental módulo por módulo
   2. **Feature Flags:** Rollback instantáneo sin redeploy
   3. **Testing First:** Cada cambio validado en staging antes de prod
   4. **Monitoring:** Visibilidad completa de métricas en tiempo real
   5. **User First:** Cero impacto negativo para usuarios finales
   ```

2. **Diagrama de Gantt**
   ```mermaid
   gantt
       title Roadmap de Refactorización
       dateFormat  YYYY-MM-DD

       section Preparación
       Setup Monolito                    :done, setup, 2026-03-01, 7d
       CI/CD Pipeline                   :done, cicd, after setup, 3d
       Feature Flags Service            :active, flags, after setup, 4d

       section Migración Core
       Auth Module Migration            :auth, after flags, 7d
       DB Consolidation (Critical)      :crit, db, after auth, 14d
       Transactions Module              :trans, after db, 10d

       section Gamification
       Event Bus Implementation         :events, after trans, 5d
       Gamification Module Migration    :gamif, after events, 10d

       section Resto de Módulos
       Budgets Module                   :after gamif, 7d
       Savings Module                   :after gamif, 7d
       Recurring Module                 :after gamif, 7d
       Analytics Module                 :after gamif, 7d

       section Mejoras
       Soft Delete Implementation       :improve, after gamif, 7d
       Cron Jobs                       :after improve, 5d
       Redis Cache                     :after improve, 5d

       section Finalización
       Testing & QA                    :testing, after improve, 7d
       Gradual Cutover                 :crit, cutover, after testing, 7d
       Monitoring & Cleanup            :after cutover, 7d
   ```

3. **Hitos (Milestones)**
   ```markdown
   ## 🎯 Hitos Clave

   ### Milestone 1: Foundation Ready (Semana 1)
   **Fecha objetivo:** 2026-03-07
   **Criterios de éxito:**
   - ✅ Monolito compila y despliega a staging
   - ✅ Health check endpoint responde
   - ✅ CI/CD pipeline automatizado
   - ✅ Feature flags operativos

   **Riesgo:** 🟢 Bajo
   **Responsable:** DevOps Lead

   ---

   ### Milestone 2: Auth Migrated (Semana 2)
   **Fecha objetivo:** 2026-03-14
   **Criterios de éxito:**
   - ✅ Auth module funcionando en monolito
   - ✅ A/B test: 50% users en monolito, 50% en users-service
   - ✅ Latencia promedio < 100ms
   - ✅ 0 errores reportados

   **Riesgo:** 🟡 Medio
   **Responsable:** Backend Team

   ---

   ### Milestone 3: Database Unified (Semana 4) 🔴 CRÍTICO
   **Fecha objetivo:** 2026-03-28
   **Criterios de éxito:**
   - ✅ Migración de gamification-db → main-db completada
   - ✅ Checksum validation: 0 discrepancias
   - ✅ UserID estandarizado a string en todos los modelos
   - ✅ Rollback plan probado
   - ✅ gamification-db desconectado

   **Riesgo:** 🔴 Alto
   **Responsable:** Database Engineer + Backend Lead
   **Buffer:** +3 días contingencia

   ---

   ### Milestone 4: Core Modules Migrated (Semana 7)
   **Fecha objetivo:** 2026-04-18
   **Criterios de éxito:**
   - ✅ Transactions module en monolito
   - ✅ Gamification module en monolito
   - ✅ Event bus funcional
   - ✅ Latencia transacciones < 50ms (mejora 4x vs HTTP)
   - ✅ XP se otorga correctamente vía events

   **Riesgo:** 🟡 Medio
   **Responsable:** Backend Team

   ---

   ### Milestone 5: All Modules Migrated (Semana 9)
   **Fecha objetivo:** 2026-05-02
   **Criterios de éxito:**
   - ✅ Budgets, Savings, Recurring, Analytics migrados
   - ✅ 100% feature parity con microservicios
   - ✅ Tests de integración E2E pasan

   **Riesgo:** 🟢 Bajo
   **Responsable:** Backend Team

   ---

   ### Milestone 6: Technical Debt Resolved (Semana 11)
   **Fecha objetivo:** 2026-05-16
   **Criterios de éxito:**
   - ✅ Soft delete implementado
   - ✅ Cron jobs ejecutando transacciones recurrentes
   - ✅ Redis cache con hit ratio > 70%
   - ✅ Transacciones atómicas garantizadas

   **Riesgo:** 🟢 Bajo
   **Responsable:** Backend Team

   ---

   ### Milestone 7: Production Cutover (Semana 12) 🔴 CRÍTICO
   **Fecha objetivo:** 2026-05-23
   **Criterios de éxito:**
   - ✅ 100% tráfico en monolito
   - ✅ Microservicios apagados
   - ✅ Latencia p95 < 200ms
   - ✅ Error rate < 0.1%
   - ✅ 0 quejas de usuarios
   - ✅ Monitoreo 24/7 activo durante 48 horas

   **Riesgo:** 🔴 Alto
   **Responsable:** Tech Lead + On-call team

   ---

   ### Milestone 8: Cleanup Complete (Semana 13)
   **Fecha objetivo:** 2026-05-30
   **Criterios de éxito:**
   - ✅ Código de microservicios eliminado del repo
   - ✅ gamification-db dropped
   - ✅ Documentación actualizada
   - ✅ Post-mortem completado
   - ✅ Celebración del equipo 🎉

   **Riesgo:** 🟢 Bajo
   **Responsable:** Tech Lead
   ```

4. **Dependencias Críticas**
   ```markdown
   ## 🔗 Grafo de Dependencias

   ```mermaid
   graph TD
       A[Setup Monolito] --> B[Auth Migration]
       B --> C[DB Consolidation]
       C --> D[Transactions Module]
       D --> E[Event Bus]
       E --> F[Gamification Module]
       F --> G[Budgets Module]
       F --> H[Savings Module]
       F --> I[Recurring Module]
       F --> J[Analytics Module]
       G --> K[Soft Delete]
       H --> K
       I --> K
       J --> K
       K --> L[Cron Jobs]
       K --> M[Redis Cache]
       L --> N[Testing]
       M --> N
       N --> O[Cutover]
       O --> P[Cleanup]

       style C fill:#f99,stroke:#f00,stroke-width:4px
       style O fill:#f99,stroke:#f00,stroke-width:4px
   ```

   **Notas:**
   - 🔴 Nodos rojos = Milestones críticos (DB Consolidation, Cutover)
   - DB Consolidation bloquea todo lo demás → máxima prioridad
   - Event Bus necesario antes de Gamification
   - Cutover solo después de testing exhaustivo
   ```

5. **Métricas de Éxito por Fase**
   ```markdown
   ## 📊 KPIs y Métricas

   ### Métricas de Performance

   | Métrica | Baseline (Microservicios) | Target (Monolito) | Método de Medición |
   |---------|---------------------------|-------------------|-------------------|
   | Latencia p50 (Dashboard) | ~150ms | <50ms | OpenTelemetry |
   | Latencia p95 (Dashboard) | ~300ms | <200ms | OpenTelemetry |
   | Latencia p99 (Transactions) | ~400ms | <300ms | OpenTelemetry |
   | Throughput (req/s) | ~50 | >100 | Prometheus |
   | Error Rate | ~0.5% | <0.1% | Prometheus |
   | Cache Hit Ratio | N/A (sin cache) | >70% | Redis metrics |

   ### Métricas de Calidad

   | Métrica | Target | Método de Medición |
   |---------|--------|-------------------|
   | Code Coverage | >80% | Go coverage tool |
   | Integration Test Coverage | >70% | Custom scripts |
   | E2E Test Pass Rate | 100% | CI pipeline |
   | Data Consistency Checks | 0 discrepancias | SQL audits |

   ### Métricas de Operaciones

   | Métrica | Baseline | Target | Mejora |
   |---------|----------|--------|--------|
   | Deploy Time | ~15 min (4 servicios) | <5 min (1 servicio) | 3x |
   | Rollback Time | ~10 min | <5 min | 2x |
   | # Procesos en Prod | 4 | 1 | 4x simplificación |
   | Costo Infraestructura | $400/mes | $100/mes | 75% reducción |

   ### Métricas de Negocio

   | Métrica | Esperado | Método |
   |---------|----------|--------|
   | User Retention | Sin cambio | Analytics |
   | Transaction Volume | Sin cambio | DB queries |
   | User Complaints | 0 relacionados con refactor | Support tickets |
   | App Crashes | <0.1% | Sentry |
   ```

6. **Plan de Comunicación**
   ```markdown
   ## 📢 Comunicación Stakeholders

   ### Usuarios Finales

   | Cuándo | Qué | Canal |
   |--------|-----|-------|
   | Semana 0 | "Mejoras de performance en camino" | Email + In-app banner |
   | Semana 4 | "Mantenimiento breve este fin de semana" (DB migration) | Email 48h antes |
   | Semana 12 | "Migración completa, disfruta la nueva velocidad" | Email + celebración in-app |

   ### Equipo Interno

   | Cuándo | Qué | Canal |
   |--------|-----|-------|
   | Semanal | Standup de refactorización (lunes 10 AM) | Slack huddle |
   | Cada milestone | Demo del progreso | Loom video |
   | Al completar fase crítica | Retrospectiva | Notion doc |

   ### Leadership

   | Cuándo | Qué | Canal |
   |--------|-----|-------|
   | Bi-semanal | Status report con métricas | Email |
   | Pre-cutover | Go/No-go decision meeting | Video call |
   | Post-cutover | Results & learnings presentation | Slide deck |
   ```

7. **Checklist de Pre-requisitos**
   ```markdown
   ## ✅ Checklist Pre-Inicio

   ### Técnico
   - [ ] Backups automáticos configurados (daily)
   - [ ] Staging environment idéntico a prod
   - [ ] Feature flags service operativo
   - [ ] Monitoreo baseline establecido
   - [ ] Incident response plan documentado
   - [ ] Rollback playbook probado

   ### Equipo
   - [ ] Team capacity validado (12 semanas comprometidas)
   - [ ] On-call rotation definida
   - [ ] Escalation path claro
   - [ ] Training en nueva arquitectura completado

   ### Proceso
   - [ ] Code freeze policy durante cutover
   - [ ] QA sign-off criteria definidos
   - [ ] Post-mortem template preparado
   - [ ] Success celebration planned 🎉
   ```

**Formato:** Markdown con diagramas mermaid (Gantt, grafos de dependencias), tablas comparativas.

---

## 🔍 Fase 6: Validación de Consistencia

### Agente: Docs-Validation

**Task:** Revisar TODA la documentación generada por fases anteriores y garantizar consistencia.

**Archivos a revisar:**
- `docs/README.md`
- `docs/00-vision/vision.md`
- `docs/01-requirements.md`
- `docs/02-user-stories.md`
- `docs/03-architecture/*.md`
- `docs/04-api-contracts/README.md`
- `docs/06-data-models/01-current-state/*.md`

**Archivo a crear:** `docs/03-architecture/00-VALIDATION-REPORT.md`

**Tareas de validación:**

1. **Consistencia de Términos**
   - Verificar que "Modular Monolith" se usa consistentemente
   - Verificar que nombres de módulos son idénticos en todos los docs
   - Verificar que referencias cruzadas funcionan

2. **Validación de Números**
   - Timeline de 12 semanas consistente en todos los docs
   - Métricas (latencia, throughput) consistentes
   - Story points y estimaciones coherentes

3. **Validación de Referencias**
   - Todos los links internos `[texto](./path)` apuntan a archivos existentes
   - Todas las secciones referenciadas existen
   - Diagramas mermaid renderizan correctamente

4. **Completitud**
   - Todos los archivos prometidos en roadmap existen
   - No hay TODOs sin resolver en la documentación
   - Todas las user stories referenciadas en migration plan existen en 02-user-stories.md

5. **Reporte Final**
   ```markdown
   # 📋 Reporte de Validación de Documentación

   **Fecha:** 2026-02-09
   **Validador:** Docs-Validation Agent

   ## ✅ Validaciones Exitosas

   - [x] Consistencia de términos verificada
   - [x] Referencias cruzadas validadas
   - [x] Diagramas mermaid renderizan
   - [x] Timeline consistente (12 semanas)
   - [x] Métricas coherentes entre documentos

   ## ⚠️ Issues Encontrados

   Ninguno.

   ## 📊 Estadísticas

   - **Total de documentos:** 11
   - **Total de páginas (estimado):** 150+
   - **Diagramas mermaid:** 8
   - **Tablas:** 25+
   - **Links internos:** 50+
   - **User stories:** 53
   - **Requirements:** 139 RFs + 65 RNFs + 10 RRs

   ## 🎯 Conclusión

   ✅ Documentación COMPLETA y CONSISTENTE. Lista para ser usada como referencia durante la refactorización.
   ```

---

## 📝 Instrucciones Finales para el Orquestador

**Cómo ejecutar este plan:**

1. **Lanzar 6 agentes en paralelo** (uno por fase), PERO respetando dependencias:
   - Fase 1 y 2 pueden correr en paralelo (no dependen entre sí)
   - Fase 3 depende de Fase 1 y 2 (necesita arquitectura documentada)
   - Fase 4 puede correr en paralelo con Fase 3
   - Fase 5 depende de Fases 1-4 (necesita todo el contexto)
   - Fase 6 debe ser última (valida todo lo anterior)

2. **Usar agentes general-purpose con capacidad de ESCRITURA**

3. **Integración con spec-kit:**
   - Los agentes de las Fases 1-7 NO deben llamar a spec-kit directamente
   - En su lugar, deben documentar en `03-migration-plan.md` que "el implementador debe usar spec-kit" para cada fase
   - Incluir comandos exactos de spec-kit en la documentación
   - Los specs de spec-kit se crearán DESPUÉS cuando se implemente cada fase

4. **Prompt para cada agente debe incluir:**
   ```
   IMPORTANTE:
   - CREAR ARCHIVOS DIRECTAMENTE en las rutas especificadas
   - NO generar outputs temporales
   - Leer este prompt completo antes de empezar
   - Verificar que el archivo anterior en la cadena existe antes de empezar
   - Si la fase requiere spec-kit, documentar los comandos pero NO ejecutarlos
   ```

4. **Validación post-ejecución:**
   - Leer el archivo 00-VALIDATION-REPORT.md generado por Fase 6
   - Si hay issues, corregir manualmente o relanzar agente específico

5. **Resultado esperado:**
   ```
   docs/
   ├── README.md (actualizado)
   ├── 00-vision/vision.md (sin cambios)
   ├── 01-requirements.md (actualizado con RRs)
   ├── 02-user-stories.md (sin cambios)
   ├── 03-architecture/
   │   ├── 00-VALIDATION-REPORT.md (nuevo)
   │   ├── 01-current-state.md (nuevo)
   │   ├── 02-target-state.md (nuevo)
   │   ├── 03-migration-plan.md (nuevo - incluye refs a spec-kit)
   │   └── 04-implementation-roadmap.md (nuevo)
   ├── 04-api-contracts/README.md (sin cambios)
   └── 06-data-models/01-current-state/ (sin cambios)

   .claude/specs/refactoring/
   └── (estos archivos se crearán DESPUÉS durante implementación con spec-kit)
       ├── 01-setup-monolith.md
       ├── 02-migrate-auth-module.md
       ├── 03-consolidate-databases.md
       ├── 04-migrate-transactions-module.md
       ├── 05a-implement-event-bus.md
       ├── 05b-migrate-gamification-module.md
       ├── 06a-migrate-budgets-module.md
       ├── 06b-migrate-savings-module.md
       ├── 06c-migrate-recurring-module.md
       ├── 06d-migrate-analytics-module.md
       └── 07-technical-debt-resolution.md
   ```

6. **Flujo de trabajo completo:**
   ```
   AHORA (Documentación):
   1. Ejecutar 6 agentes según este prompt
   2. Generar toda la documentación en docs/
   3. Validar consistencia

   DESPUÉS (Implementación):
   1. Para cada fase de migración:
      - Usar /specify para crear spec
      - Usar /plan para crear plan de implementación
      - Usar /tasks para generar tasks accionables
      - Usar /implement para ejecutar
   2. Los specs quedan en .claude/specs/refactoring/
   3. La documentación en docs/ sirve como referencia macro
   ```

---

**FIN DEL PROMPT**
