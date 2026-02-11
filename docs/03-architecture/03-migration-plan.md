# Migration Plan: Distributed Monolith → Modular Monolith

**Document Version**: 1.0
**Date**: 2026-02-10
**Status**: Approved - Ready for Implementation
**Total Duration**: ~11 weeks
**Current State**: See [01-current-state.md](./01-current-state.md)
**Target State**: See [02-target-state.md](./02-target-state.md)

---

## Table of Contents

1. [Overview](#overview)
2. [Pre-Migration Checklist](#pre-migration-checklist)
3. [Migration Strategy](#migration-strategy)
4. [Migration Phases](#migration-phases)
   - [Phase 1: Monolith Foundation Setup (1 week)](#phase-1-monolith-foundation-setup-1-week)
   - [Phase 2: Auth Module Migration (1 week)](#phase-2-auth-module-migration-1-week)
   - [Phase 3: Database Consolidation (2 weeks)](#phase-3-database-consolidation-2-weeks--critical)
   - [Phase 4: Transactions Module Migration (1.5 weeks)](#phase-4-transactions-module-migration-15-weeks)
   - [Phase 5: Gamification Module Migration (1.5 weeks)](#phase-5-gamification-module-migration-15-weeks)
   - [Phase 6: Remaining Modules Migration (2 weeks)](#phase-6-remaining-modules-migration-2-weeks)
   - [Phase 7: Technical Debt Resolution (2 weeks)](#phase-7-technical-debt-resolution-2-weeks)
5. [Dependency Map](#dependency-map)
6. [Risk Assessment](#risk-assessment)
7. [Success Metrics](#success-metrics)
8. [Rollback Strategy](#rollback-strategy)
9. [Implementation Notes](#implementation-notes)

---

## Overview

### Strategy: Incremental Migration (NO Big Bang)

This migration transitions Financial Resume from a **Distributed Monolith** (5 processes managed by supervisord in 1 Docker container, communicating via HTTP localhost) to a **Modular Monolith** (1 Go binary, 1 database, in-memory function calls).

**Core Principles:**
- **Incremental**: One module at a time, never a full rewrite in a single step
- **Rollback available**: Every phase can be reverted independently before cutting over traffic
- **Beta environment**: This repo coexists with the production project — we can break things freely, the only sensitive asset is the database data

### Why Incremental (Not Big Bang)?

```
Incremental approach benefits:
  - Each phase is independently verifiable
  - Rollback is always one step back, never starting from scratch
  - Confidence increases with each successful phase
  - Easier to debug issues in smaller changesets
```

> **Note**: Since this is a beta environment that coexists with production, we have
> full freedom to experiment. Downtime is acceptable. The only critical asset is the
> database data (which will be migrated to a new unified DB anyway).

### Context: Why This Migration Is Simpler Than Typical Microservices

Since the current architecture is already a **Distributed Monolith** (not real microservices), significant advantages reduce migration risk:

- All 4 services live in the **same git repository** — no cross-repo coordination
- Everything deploys as **1 Docker container** on Render.com — no infrastructure topology changes
- Services are **co-located in the same container** — "migration" is eliminating unnecessary HTTP calls, not moving network topology
- **Same $7/month cost** throughout — no budget risk during migration
- The elimination targets: supervisord, nginx (internal proxy), HTTP localhost calls, JWT enrichment proxies

---

## Pre-Migration Checklist

Complete before starting Phase 1.

### Database Safety
- [ ] Full backup of both databases (`financial_resume` and `users_db`)
- [ ] Verify backup can be restored

### Documentation
- [ ] Architecture docs complete (current-state + target-state: ✅ done)
- [ ] API contracts documented (✅ done)

### Development
- [ ] Feature flag strategy defined (see below — optional but recommended)
- [ ] Deployment procedure for Render.com understood

> **Note**: No staging environment needed — this IS the beta environment.
> No monitoring setup needed — we can check logs directly on Render.

---

## Migration Strategy

### Approach: Strangler Fig Pattern

The migration uses the **Strangler Fig Pattern**: the new modular monolith grows alongside the existing services. Traffic is progressively routed from old services to new modules. Old code is only removed once the new module handles 100% of traffic with validated correctness.

```
Phase N: Old service handles 100% traffic
         ↓ (implement new module)
Phase N + validation: Old service still runs, new module tested in parallel
         ↓ (feature flag flipped)
Phase N + cutover: New module handles 100% traffic
         ↓ (old service verified idle)
Phase N + cleanup: Old service process removed from supervisord.conf
```

**Applied to this project:**
- `users-service` → `auth` module (Phase 2)
- `api-gateway` transactions logic → `transactions` module (Phase 4)
- `gamification-service` → `gamification` module (Phase 5)
- `ai-service` → `ai` module (Phase 6, or kept as separate service)

### Feature Flags

Feature flags are **optional but recommended** as a good engineering practice. Since this is a beta environment, we can also migrate directly without flags if preferred. If used, flags are environment variables toggled in Render.com dashboard.

```go
// In infrastructure/config/features.go
type FeatureFlags struct {
    UseMonolithAuth         bool   // Phase 2: Route auth requests to new module
    UseMonolithTransactions bool   // Phase 4: Route transaction requests to new module
    UseEventBus             bool   // Phase 5: Enable in-memory event bus
    UseMonolithGamification bool   // Phase 5: Route gamification requests to new module
    UseMonolithBudgets      bool   // Phase 6: Route budget requests to new module
    UseMonolithSavings      bool   // Phase 6: Route savings requests to new module
    UseMonolithAI           bool   // Phase 6: Route AI requests to new module
}

func LoadFeatureFlags() FeatureFlags {
    return FeatureFlags{
        UseMonolithAuth:         getEnvBool("USE_MONOLITH_AUTH", false),
        UseMonolithTransactions: getEnvBool("USE_MONOLITH_TRANSACTIONS", false),
        UseEventBus:             getEnvBool("USE_EVENT_BUS", false),
        UseMonolithGamification: getEnvBool("USE_MONOLITH_GAMIFICATION", false),
        UseMonolithBudgets:      getEnvBool("USE_MONOLITH_BUDGETS", false),
        UseMonolithSavings:      getEnvBool("USE_MONOLITH_SAVINGS", false),
        UseMonolithAI:           getEnvBool("USE_MONOLITH_AI", false),
    }
}
```

**Feature flag routing pattern:**
```go
// In HTTP router (during migration period)
func (r *Router) RegisterAuthRoutes(oldProxy *UsersServiceProxy, newModule *auth.Module, flags FeatureFlags) {
    r.engine.POST("/api/v1/auth/register", func(c *gin.Context) {
        if flags.UseMonolithAuth {
            newModule.RegisterHandler(c)  // New module
        } else {
            oldProxy.ProxyRequest(c)      // Old users-service
        }
    })
}
```

**Flag lifecycle:**
1. Flag defaults to `false` (old service used)
2. Flip to `true` — verify functionality works
3. If no regressions: remove flag and old code path in next release

---

## Migration Phases

### Phase 1: Monolith Foundation Setup (1 week)

**Objectives:**
- Create the new monolith application skeleton
- Set up shared infrastructure (database connection, event bus, config, logging)
- Establish CI/CD pipeline
- Produce a deployable binary (health check only — no business logic yet)

**High-Level Tasks:**

| Task | Description |
|------|-------------|
| Create monolith app structure | `apps/monolith/` directory with Go module |
| Implement shared event bus | `InMemoryEventBus` with subscribe/publish |
| Unified database connection | Single `*gorm.DB` with connection pooling |
| Structured logging | Replace `fmt.Println` with `zerolog` |
| Configuration management | `viper`-based config (env vars + files) |
| Health check endpoint | `GET /health` returns 200 OK |
| Dockerfile (simplified) | Single binary, no supervisord, no nginx |
| CI/CD pipeline | GitHub Actions: lint → test → build → deploy |
| Feature flags infrastructure | `FeatureFlags` struct loaded from env |

**Files to Create:**
```
apps/monolith/
├── cmd/server/main.go               # Entry point
├── internal/
│   ├── shared/
│   │   ├── events/event_bus.go      # InMemoryEventBus
│   │   ├── events/events.go         # Base event types
│   │   ├── ports/repositories.go    # Shared repository interfaces
│   │   ├── ports/services.go        # Shared service interfaces
│   │   └── errors/errors.go         # Shared error types
│   └── infrastructure/
│       ├── database/postgres.go     # Single DB connection + pooling
│       ├── config/config.go         # Viper configuration
│       ├── config/features.go       # Feature flags
│       ├── http/router.go           # Gin router setup
│       ├── http/server.go           # HTTP server with graceful shutdown
│       └── middleware/
│           ├── auth.go              # JWT middleware
│           ├── cors.go
│           └── logging.go
├── Dockerfile                       # Simplified: single binary
├── go.mod
└── go.sum
```

**Spec-Kit Command (when implementing):**
```
/specify "Setup Monolith Foundation"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 01-setup-monolith.md             # Feature spec
└── 01-setup-monolith-tasks.md       # Implementation tasks
```

**Success Criteria:**
- [ ] `go build ./...` compiles without errors
- [ ] `GET /health` returns `{"status": "ok"}` in staging
- [ ] CI/CD pipeline runs and deploys to staging
- [ ] Unit tests for event bus pass
- [ ] Database connection pool configured (MaxIdleConns: 10, MaxOpenConns: 100)

---

### Phase 2: Auth Module Migration (1 week)

**Objectives:**
- Migrate `users-service` (port 8083) into the `auth` module
- Eliminate `UsersServiceProxy` from `api-gateway`
- Standardize `user_id` type to `VARCHAR(255)` in all tables

**Feature Flag:** `USE_MONOLITH_AUTH`

**High-Level Tasks:**

| Task | Description |
|------|-------------|
| Copy User domain model | From `users-service/internal/core/domain/` |
| Implement auth use cases | Register, login, JWT issuance, profile CRUD |
| Implement auth handlers | POST /auth/register, POST /auth/login, GET/PUT /users/:id |
| Create user repository | GORM implementation for `users` + `user_preferences` tables |
| JWT middleware for monolith | Validate tokens — same secret, same format |
| Feature flag routing | Toggle between `UsersServiceProxy` and new `auth.Module` |
| Integration tests | Verify response parity against old users-service |

**Validation: Response Parity Testing**

Before flipping `USE_MONOLITH_AUTH=true` in production, responses from both old and new implementations must be identical:

```bash
# Script: scripts/validate_auth_parity.sh
# Run both old and new in staging, compare responses

OLD_RESPONSE=$(curl -s -X POST http://staging-old/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}')

NEW_RESPONSE=$(curl -s -X POST http://staging-new/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}')

# Compare (ignoring timestamps and token values, check structure)
diff <(echo "$OLD_RESPONSE" | jq 'keys') <(echo "$NEW_RESPONSE" | jq 'keys')
```

**Spec-Kit Command (when implementing):**
```
/specify "Migrate Auth Module to Monolith"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 02-auth-module.md
└── 02-auth-module-tasks.md
```

**Success Criteria:**
- [ ] All auth endpoints return structurally identical responses to `users-service`
- [ ] JWT tokens issued by new module are validated correctly
- [ ] Frontend login/register flow works without modification
- [ ] `USE_MONOLITH_AUTH=true` works correctly (if using feature flags)
- [ ] `UsersServiceProxy` still available as fallback (not yet deleted)

---

### Phase 3: Database Consolidation (2 weeks) ⚠️ CRITICAL

**Objectives:**
- Merge `users_db` into `financial_resume` (single database)
- Resolve duplicate `user_gamification` table data
- Standardize all `user_id` columns to `VARCHAR(255)`
- Shut down `users_db` Render instance

**This phase requires care with data.** Even though this is a beta environment, user data should be preserved. The new unified DB will be created fresh, so the original databases remain untouched as a natural backup.

**Sub-Phase 3.1: Analysis & Preparation (3 days)**

Before writing any migration SQL:

```sql
-- Audit: How many users exist in each DB?
SELECT COUNT(*) FROM users_db.public.users;                        -- users-service users
SELECT COUNT(DISTINCT user_id) FROM financial_resume.public.expenses; -- api-gateway users

-- Audit: Are there user_id values in financial_resume that don't exist in users_db?
SELECT DISTINCT e.user_id
FROM financial_resume.public.expenses e
WHERE e.user_id NOT IN (SELECT CAST(id AS VARCHAR) FROM users_db.public.users);

-- Audit: Gamification data duplication
SELECT user_id, COUNT(*) as records, MAX(total_xp) as max_xp, MIN(total_xp) as min_xp
FROM financial_resume.public.user_gamification
GROUP BY user_id
HAVING COUNT(*) > 1;

-- Audit: user_id types in use
SELECT column_name, data_type, table_name
FROM information_schema.columns
WHERE column_name = 'user_id'
  AND table_schema = 'public';
```

**Sub-Phase 3.2: Schema Migration — `financial_resume` DB (3 days)**

Apply changes to the existing `financial_resume` database:

```sql
-- Step 1: Add deleted_at to all tables (non-destructive, all NULLs)
ALTER TABLE expenses          ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE incomes           ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE budgets           ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE savings_goals     ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE recurring_transactions ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE user_gamification ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE categories        ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Step 2: Create users table IN financial_resume (copy from users_db)
-- NOTE: Run this from a migration script with access to both DBs
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,    -- Changed from SERIAL INT to VARCHAR(255)
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Step 3: Create user_preferences IN financial_resume
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    currency VARCHAR(10) DEFAULT 'USD',
    language VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50) DEFAULT 'UTC',
    gamification_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Step 4: Add composite indexes
CREATE INDEX IF NOT EXISTS idx_expenses_user_date ON expenses(user_id, transaction_date)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_expenses_category ON expenses(category_id)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_budgets_user_active ON budgets(user_id, is_active)
    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_recurring_next_exec ON recurring_transactions(next_execution, is_active)
    WHERE deleted_at IS NULL;
```

**Sub-Phase 3.3: Data Migration (3 days)**

```sql
-- Step 1: Copy users from users_db to financial_resume
-- Cast SERIAL int ID to VARCHAR(255)
INSERT INTO financial_resume.public.users (id, email, password_hash, name, created_at, updated_at)
SELECT
    CAST(id AS VARCHAR(255)) AS id,
    email,
    password_hash,
    name,
    created_at,
    updated_at
FROM users_db.public.users
ON CONFLICT (id) DO NOTHING;

-- Step 2: Copy user_preferences
INSERT INTO financial_resume.public.user_preferences (user_id, gamification_enabled, created_at, updated_at)
SELECT
    CAST(user_id AS VARCHAR(255)) AS user_id,
    gamification_enabled,
    created_at,
    updated_at
FROM users_db.public.user_preferences
ON CONFLICT (user_id) DO NOTHING;

-- Step 3: Resolve duplicate user_gamification rows
-- Keep the row with the highest total_xp per user
-- First, create a temp table with deduplicated data
CREATE TEMP TABLE user_gamification_dedup AS
SELECT DISTINCT ON (user_id) *
FROM user_gamification
ORDER BY user_id, total_xp DESC;

-- Verify dedup looks correct before applying
SELECT COUNT(*) FROM user_gamification;        -- before
SELECT COUNT(*) FROM user_gamification_dedup;  -- after dedup

-- Apply dedup (only if counts look right)
DELETE FROM user_gamification;
INSERT INTO user_gamification SELECT * FROM user_gamification_dedup;
DROP TABLE user_gamification_dedup;

-- Step 4: Verify referential integrity
-- No expense should have a user_id that doesn't exist in the new users table
SELECT COUNT(*) FROM expenses e
WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = e.user_id);
-- Expected: 0 rows
```

**Sub-Phase 3.4: Validation (2 days)**

```bash
# Count validation: totals must match before and after
USERS_BEFORE=$(psql $USERS_DB_URL -c "SELECT COUNT(*) FROM users" -t)
USERS_AFTER=$(psql $MAIN_DB_URL -c "SELECT COUNT(*) FROM users" -t)
echo "Users: before=$USERS_BEFORE, after=$USERS_AFTER"  # Must match

EXPENSES_COUNT=$(psql $MAIN_DB_URL -c "SELECT COUNT(*) FROM expenses" -t)
echo "Expenses: $EXPENSES_COUNT"  # Must match pre-migration count

# Foreign key consistency
ORPHANED=$(psql $MAIN_DB_URL -c "
    SELECT COUNT(*) FROM expenses e
    WHERE NOT EXISTS (SELECT 1 FROM users u WHERE u.id = e.user_id)
" -t)
echo "Orphaned expenses: $ORPHANED"  # Must be 0
```

**Rollback Plan for Phase 3:**

Phase 3 is the only phase with a complex rollback. The key is that `users_db` is NOT shut down until Sub-Phase 3.4 validation passes completely.

```
Rollback trigger: Any data count mismatch OR orphaned FK records

Rollback steps:
  1. Do NOT shut down users_db (it was never touched destructively)
  2. Set USE_MONOLITH_AUTH=false → auth requests route back to users-service proxy
  3. Drop the newly created `users` and `user_preferences` tables from financial_resume
  4. Revert schema changes (delete deleted_at columns) — or leave them (they are NULL and harmless)
  5. Root-cause analysis before retry

Recovery time objective: <2 hours (no data loss possible if users_db is untouched)
```

**Spec-Kit Command (when implementing):**
```
/specify "Consolidate Databases"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 03-db-consolidation.md
└── 03-db-consolidation-tasks.md
```

**Success Criteria:**
- [ ] All users present in `financial_resume.users` table
- [ ] Zero orphaned foreign key references
- [ ] `user_gamification` has exactly 1 row per `user_id`
- [ ] All `user_id` columns are `VARCHAR(255)` (no remaining `INT`/`uint` columns)
- [ ] `deleted_at` column present on all entity tables
- [ ] `users_db` Render instance deleted (cost reduction: free tier slot freed)
- [ ] Auth module still functional after DB consolidation

---

### Phase 4: Transactions Module Migration (1.5 weeks)

**Objectives:**
- Migrate `Expenses`, `Incomes`, and `Categories` from `api-gateway` internal code into the `transactions` module
- Remove HTTP proxy overhead for transaction-related operations
- Publish `ExpenseCreatedEvent` and `IncomeCreatedEvent` to the event bus

**Feature Flag:** `USE_MONOLITH_TRANSACTIONS`

**High-Level Tasks:**

| Task | Description |
|------|-------------|
| Copy domain models | `Expense`, `Income`, `Category` aggregates from `api-gateway/internal/` |
| Migrate use cases | CreateExpense, UpdateExpense, DeleteExpense, ListExpenses (and incomes) |
| Migrate HTTP handlers | All `/api/v1/expenses/*` and `/api/v1/incomes/*` routes |
| Implement soft delete | Add `deleted_at` filtering to all repository queries |
| Publish domain events | `ExpenseCreatedEvent`, `IncomeCreatedEvent` on all create operations |
| Wire feature flag routing | Toggle between api-gateway handlers and new module |
| Integration tests | Response parity validation for all transaction endpoints |

**Key Code Change — Soft Delete:**
```go
// BEFORE (hard delete — current production)
func (r *ExpenseRepository) Delete(id string) error {
    return r.db.Delete(&domain.Expense{}, "id = ?", id).Error
}

// AFTER (soft delete — monolith module)
func (r *ExpenseRepository) Delete(id string) error {
    return r.db.Model(&domain.Expense{}).
        Where("id = ?", id).
        Update("deleted_at", time.Now().UTC()).Error
}

func (r *ExpenseRepository) FindAll(userID string) ([]domain.Expense, error) {
    var expenses []domain.Expense
    err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).
        Find(&expenses).Error
    return expenses, err
}
```

**Spec-Kit Command (when implementing):**
```
/specify "Migrate Transactions Module"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 04-transactions-module.md
└── 04-transactions-module-tasks.md
```

**Success Criteria:**
- [ ] All `/api/v1/expenses/*` and `/api/v1/incomes/*` endpoints return identical responses
- [ ] Soft delete: deleted expenses return 404 but remain in DB with `deleted_at` set
- [ ] `ExpenseCreatedEvent` and `IncomeCreatedEvent` are published to event bus on creation
- [ ] `USE_MONOLITH_TRANSACTIONS=true` works correctly (if using feature flags)
- [ ] Existing api-gateway expense handlers still available as fallback

---

### Phase 5: Gamification Module Migration (1.5 weeks)

**Objectives:**
- Replace `GamificationProxy` HTTP calls with in-memory event subscriptions
- Migrate `gamification-service` (port 8084) into the `gamification` module
- Remove the only remaining cross-service HTTP call from the critical path

This phase has **two distinct spec-kit specs** because the event bus and the gamification module are independent concerns:

**Feature Flags:** `USE_EVENT_BUS` (independent), `USE_MONOLITH_GAMIFICATION`

**Sub-Phase 5.1: Implement In-Memory Event Bus**

The event bus is infrastructure shared by all modules. It must be tested in isolation before any module depends on it.

```go
// shared/events/event_bus.go
type InMemoryEventBus struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func (bus *InMemoryEventBus) Publish(ctx context.Context, event Event) error {
    bus.mu.RLock()
    handlers := bus.handlers[event.Type()]
    bus.mu.RUnlock()

    for _, handler := range handlers {
        go func(h EventHandler) {
            if err := h(ctx, event); err != nil {
                log.Error().Err(err).
                    Str("event_type", event.Type()).
                    Msg("Event handler failed")
            }
        }(handler)
    }
    return nil
}
```

**Spec-Kit Command for event bus:**
```
/specify "Implement In-Memory Event Bus"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 05a-event-bus.md
└── 05a-event-bus-tasks.md
```

**Sub-Phase 5.2: Gamification Module**

| Task | Description |
|------|-------------|
| Copy gamification domain | `UserGamification`, `Achievement`, `Challenge` from `gamification-service` |
| Subscribe to domain events | `ExpenseCreatedEvent` (+5 XP), `GoalAchievedEvent` (+50 XP), `UserLoggedInEvent` (streak) |
| Migrate gamification handlers | GET /gamification/stats, GET /gamification/achievements |
| Feature flag routing | Toggle between `GamificationProxy` and new module |
| Consolidate XP rules | Verify XP rules match between old service and new module |
| Integration tests | XP calculation parity, achievement unlock parity |

**Event Subscriptions to Implement:**
```go
// gamification/module.go
func (m *GamificationModule) RegisterSubscribers(bus ports.EventBus) {
    bus.Subscribe("expense.created",  m.OnExpenseCreated)   // +5 XP
    bus.Subscribe("income.created",   m.OnIncomeCreated)    // +5 XP
    bus.Subscribe("budget.created",   m.OnBudgetCreated)    // +10 XP
    bus.Subscribe("goal.achieved",    m.OnGoalAchieved)     // +50 XP
    bus.Subscribe("user.logged_in",   m.OnUserLoggedIn)     // streak update
    bus.Subscribe("challenge.completed", m.OnChallengeCompleted) // +100 XP
}
```

**Spec-Kit Command for gamification module:**
```
/specify "Migrate Gamification Module to Monolith"
```

**Spec Artifacts Location:**
```
.claude/specs/refactoring/
├── 05b-gamification-module.md
└── 05b-gamification-module-tasks.md
```

**Success Criteria:**
- [ ] Event bus unit tests pass (subscribe, publish, async execution)
- [ ] XP is awarded correctly on expense creation (no regression)
- [ ] Achievement unlock logic matches old `gamification-service` behavior
- [ ] Streak calculation is correct after `UserLoggedInEvent`
- [ ] `GamificationProxy` removed from critical path (HTTP localhost call eliminated)
- [ ] `USE_MONOLITH_GAMIFICATION=true` works correctly (if using feature flags)

---

### Phase 6: Remaining Modules Migration (2 weeks)

**Objectives:**
- Migrate `Budgets`, `SavingsGoals`, `RecurringTransactions`, and `Analytics` modules
- Wire all event subscriptions between modules
- Optionally migrate or keep `ai-service` (see notes)

**Feature Flags:** `USE_MONOLITH_BUDGETS`, `USE_MONOLITH_SAVINGS`, `USE_MONOLITH_AI`

**Modules in this phase:**

#### Budgets Module

| Task | Description |
|------|-------------|
| Copy Budget domain | `Budget` aggregate, state machine, `BudgetPeriod` value object |
| Subscribe to transaction events | Recalculate budget status on `ExpenseCreatedEvent` |
| Publish budget events | `BudgetExceededEvent`, `BudgetWarningEvent` |
| Cron job: budget reset | Monthly reset triggered by cron scheduler |
| Integration tests | Budget state machine (on_track → warning → exceeded) |

**Spec-Kit Command:**
```
/specify "Migrate Budgets Module"
```

**Spec Artifacts:** `.claude/specs/refactoring/06a-budgets-module.md`

---

#### Savings Goals Module

| Task | Description |
|------|-------------|
| Copy SavingsGoal domain | `SavingsGoal` aggregate, `GoalStatus` state machine |
| Subscribe to income events | Optional auto-save from `IncomeCreatedEvent` |
| Publish savings events | `GoalAchievedEvent` (auto-detected when target reached) |
| Integration tests | Goal progress tracking, auto-achievement |

**Spec-Kit Command:**
```
/specify "Migrate Savings Goals Module"
```

**Spec Artifacts:** `.claude/specs/refactoring/06b-savings-module.md`

---

#### Recurring Transactions Module

**Note:** Recurring transactions cron is already running in production (initialized in `apps/api-gateway/cmd/api/main.go`, runs every 1 hour). This module migration primarily means moving the scheduler initialization into the monolith's `main.go`.

| Task | Description |
|------|-------------|
| Copy RecurringTransaction domain | Keep existing scheduler logic intact |
| Move cron initialization | From api-gateway's `main.go` to monolith's `main.go` |
| Use `robfig/cron` scheduler | More control over scheduling granularity |
| Verify cron job continuity | No missed executions during migration window |

**Spec-Kit Command:**
```
/specify "Migrate Recurring Transactions Module"
```

**Spec Artifacts:** `.claude/specs/refactoring/06c-recurring-module.md`

---

#### Analytics Module

| Task | Description |
|------|-------------|
| Migrate dashboard endpoint | GET /api/v1/dashboard → reads from all modules via interfaces |
| Migrate analytics endpoints | GET /api/v1/expenses/analytics and related |
| Parallel DB queries | Use goroutines for independent queries (expenses, budgets, goals) |
| Cache invalidation hooks | Subscribe to events that should invalidate dashboard cache |

**Spec-Kit Command:**
```
/specify "Migrate Analytics Module"
```

**Spec Artifacts:** `.claude/specs/refactoring/06d-analytics-module.md`

---

#### AI Module (Decision Point)

The AI module calls OpenAI API (2000-3000ms latency). Two options:

**Option A: Migrate into monolith** (recommended for current scale)
- Simple: no separate process, no HTTP proxy
- Acceptable: OpenAI calls are already slow — monolith doesn't make them slower
- The existing `ai-service` proxy adds 25-45ms overhead that can be eliminated

**Option B: Keep as separate service** (defer to Stage 4 scalability)
- Valid if AI calls start blocking other requests
- Requires keeping `AIServiceProxy` and a separate process

**Default recommendation**: Option A (migrate into monolith). Revisit if AI becomes a bottleneck.

**Spec-Kit Command (if Option A):**
```
/specify "Migrate AI Module to Monolith"
```

**Spec Artifacts:** `.claude/specs/refactoring/06e-ai-module.md`

**Overall Phase 6 Success Criteria:**
- [ ] Budget status updates automatically when expense is created (via event)
- [ ] Savings goal auto-achieves when target is reached
- [ ] Recurring transactions continue executing on schedule (no missed executions)
- [ ] Dashboard returns correct data aggregated from all modules
- [ ] All feature flags set to `true` and working correctly (if using feature flags)

---

### Phase 7: Technical Debt Resolution (2 weeks)

**Objectives:**
- Implement remaining technical debt items that are now tractable in the monolith
- Decommission all old services, supervisord, nginx
- Simplify Dockerfile to single-binary

**Spec-Kit Command:**
```
/specify "Technical Debt Resolution"
```

**Spec Artifacts:** `.claude/specs/refactoring/07-tech-debt.md`

#### 7.1 Soft Delete (if not completed in Phase 4)

Ensure `deleted_at` is respected in all queries across all modules:

```go
// All repositories: filter out soft-deleted records
func (r *Repository) FindAll(userID string) ([]T, error) {
    var records []T
    err := r.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&records).Error
    return records, err
}
```

#### 7.2 Redis Cache (optional — add when >50 concurrent users)

```go
// analytics/infrastructure/cache.go
type DashboardCache struct {
    redis  *redis.Client
    ttl    time.Duration
}

func (c *DashboardCache) Get(userID string) (*DashboardData, error) {
    val, err := c.redis.Get(context.Background(), "dashboard:"+userID).Result()
    if err == redis.Nil {
        return nil, nil // cache miss
    }
    // ...
}
```

#### 7.3 Atomic Transactions (ACID)

Replace the current non-atomic pattern with proper DB transactions:

```go
// BEFORE (non-atomic — current production):
func (uc *CreateExpenseUseCase) Execute(params) error {
    expense := uc.expenseRepo.Create(expense) // commits
    uc.gamificationProxy.TrackAction(...)     // separate commit, can fail
}

// AFTER (ACID — monolith):
func (uc *CreateExpenseUseCase) Execute(params) error {
    tx := uc.db.Begin()
    defer tx.Rollback()

    expense := uc.expenseRepo.CreateTx(tx, expense)
    // Gamification handled via event bus after commit — no atomicity issue
    tx.Commit()
    uc.eventBus.Publish(ctx, ExpenseCreatedEvent{...}) // async, non-blocking
}
```

#### 7.4 Final Cleanup: Decommission Old Architecture

This is the last step. Only execute after ALL modules are migrated and verified.

```bash
# Checklist for final cleanup:
# [ ] All feature flags set to true in production
# [ ] Zero errors for 1 week with all flags enabled
# [ ] Old proxy code paths have zero traffic (confirmed via logs)

# Steps:
# 1. Remove feature flag conditional routing (keep only new module code paths)
# 2. Delete proxy files:
#    - apps/api-gateway/internal/infrastructure/proxy/gamification_proxy.go
#    - apps/api-gateway/internal/infrastructure/proxy/users_service_proxy.go
#    - apps/api-gateway/internal/infrastructure/proxy/ai_service_proxy.go
# 3. Update Dockerfile: remove multi-binary build, keep only monolith
# 4. Remove supervisord.conf
# 5. Remove nginx config (internal proxy)
# 6. Update render.yaml: remove old service entries
# 7. Delete apps/users-service/ directory
# 8. Delete apps/gamification-service/ directory
# 9. Delete apps/ai-service/ directory (if migrated to monolith)
```

**Updated Dockerfile (target state):**
```dockerfile
# Simple single-binary Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY apps/monolith/ .
RUN go build -o server ./cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**Updated render.yaml (target state):**
```yaml
services:
  - type: web
    name: financial-resume-backend
    runtime: docker
    plan: starter           # $7/month — unchanged
    healthCheckPath: /health
    autoDeploy: true
    # Single binary: no supervisord, no nginx, no multi-process

  - type: static
    name: financial-resume-frontend
    # Unchanged

databases:
  - name: financial-resume-db       # Only 1 database now
    databaseName: financial_resume
    plan: free
```

**Phase 7 Success Criteria:**
- [ ] Single Go binary running in production
- [ ] supervisord removed from container
- [ ] nginx (internal) removed from container
- [ ] Old service directories deleted from repo
- [ ] Docker image size reduced (estimate: 40-60% smaller)
- [ ] Deployment time reduced to <3 minutes
- [ ] All technical debt items addressed (soft delete, ACID transactions)

---

## Dependency Map

Phases that **must be sequential** (output of one is input of next):

```
Phase 1 (Foundation)
    └── Phase 2 (Auth)
            └── Phase 3 (DB Consolidation)      ← Requires auth module to be working
                    └── Phase 4 (Transactions)  ← Requires single DB
                            └── Phase 5 (Gamification)  ← Requires transaction events
                                    └── Phase 6 (Remaining)  ← Requires event bus + transactions
                                            └── Phase 7 (Cleanup)  ← Requires ALL modules migrated
```

Phases that **can partially overlap** (with caution):

```
Phase 2 (Auth) + Phase 3 (DB Analysis sub-phase)
    → Analysis and preparation for DB consolidation can start while Auth is in staging validation

Phase 5a (Event Bus) + Phase 4 (Transactions)
    → Event bus implementation can proceed in parallel with transactions module
    → But transactions cannot publish events until event bus is complete

Phase 6 modules (Budgets, Savings, Recurring, Analytics)
    → All four can be developed in parallel by different sessions
    → But must be wired together before the final integration test
```

**Critical Path (minimum timeline):**
```
Week 1:  Phase 1 (Foundation)
Week 2:  Phase 2 (Auth) + Phase 3 analysis starts
Week 3-4: Phase 3 (DB Consolidation — 2 weeks, critical path)
Week 5:  Phase 4 (Transactions) + Phase 5a (Event Bus)
Week 6:  Phase 5b (Gamification)
Week 7-8: Phase 6 (Remaining Modules)
Week 9-10: Phase 7 (Tech Debt + Cleanup)
Week 11: Final validation, monitoring, rollback period
```

---

## Risk Assessment

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| **DB migration corrupts data** (Phase 3) | Low | Critical | Full backup before each sub-phase; rollback plan keeps `users_db` alive; validation queries before any destructive step |
| **user_id type inconsistency breaks FK constraints** (Phase 3) | Medium | High | Audit queries in Phase 3.1 identify all type mismatches; explicit CAST in migration SQL; validation count checks |
| **Gamification XP rules diverge between old service and new module** (Phase 5) | Medium | Medium | Response parity testing: run both old and new in staging, compare XP totals for same sequence of events |
| **Cron job misses executions during migration window** (Phase 6) | Low | Medium | RecurringTransaction scheduler logs each execution; verify execution count before and after migration; overlap window where both old and new schedulers are ready |
| **Feature flag routing bug sends some requests to old and some to new** | Low | High | Route toggle is at the handler level (per endpoint), not at the network level; each flag controls a complete, isolated code path |
| **Event bus goroutine leak under load** (Phase 5) | Low | Medium | Load test event bus with 1000 concurrent events in staging; use `go test -race` to detect races; monitor goroutine count in production |
| **AI module latency affects overall monolith p99** (Phase 6) | Medium | Low | AI calls are already slow (2-3s); wrap in `context.WithTimeout(ctx, 10s)`; consider background job pattern for AI analysis |

---

## Success Metrics

The migration is considered successful when **all** of the following KPIs are met:

### Latency (p95)

| Endpoint | Current | Target | Measurement |
|----------|---------|--------|-------------|
| `GET /api/v1/dashboard` | 150-200ms | <100ms | Render.com metrics |
| `POST /api/v1/expenses` | 90-140ms | <60ms | Render.com metrics |
| `GET /api/v1/expenses/analytics` | 270-380ms | <150ms | Render.com metrics |
| `POST /api/v1/auth/login` | ~50ms | <30ms | Render.com metrics |

### Reliability

| Metric | Current | Target |
|--------|---------|--------|
| Error rate (5xx) | 1.25% | <0.1% |
| Data consistency | Eventually consistent | 100% ACID |
| Cron execution rate | ~100% (already working) | 100% |
| Orphaned data on delete | Possible (hard delete) | Zero (soft delete) |

### Operational

| Metric | Current | Target |
|--------|---------|--------|
| Running processes | 5 (supervisord + nginx + 4 Go) | 1 Go binary |
| Docker image size | ~250MB (4 binaries + nginx) | <80MB |
| Deploy time | 3-5 minutes | <2 minutes |
| Databases | 2 | 1 |
| Monthly infrastructure cost | $7 | $7 (same) |

### Developer Experience

| Metric | Current | Target |
|--------|---------|--------|
| Local dev setup time | Run 4 services + 2 DBs | Run 1 service + 1 DB |
| Test execution time | Slow (mock HTTP calls) | Fast (direct function calls) |
| Stack trace on error | 4 disconnected log streams | Single unified stack trace |

---

## Rollback Strategy

### General Principle

**At no point during the migration should rollback require more than 2 steps.**

The feature flag pattern ensures rollback is always: (1) set flag back to `false`, (2) deploy or toggle env var.

### Per-Phase Rollback

| Phase | Rollback Trigger | Rollback Action | Recovery Time |
|-------|-----------------|-----------------|---------------|
| Phase 1 (Foundation) | Monolith binary fails to start | Delete new `apps/monolith/` directory; no production impact | <1 hour |
| Phase 2 (Auth) | `USE_MONOLITH_AUTH=true` causes auth errors | Set `USE_MONOLITH_AUTH=false` in Render.com env vars | <5 minutes |
| Phase 3 (DB Consolidation) | Count mismatch or FK violations | `users_db` untouched; Auth flag already false-able | <2 hours |
| Phase 4 (Transactions) | `USE_MONOLITH_TRANSACTIONS=true` causes errors | Set `USE_MONOLITH_TRANSACTIONS=false` | <5 minutes |
| Phase 5 (Gamification) | XP not awarded or event errors | Set `USE_MONOLITH_GAMIFICATION=false` and `USE_EVENT_BUS=false` | <5 minutes |
| Phase 6 (Remaining) | Any module flag causes errors | Set relevant module flag to `false` | <5 minutes |
| Phase 7 (Cleanup) | Only attempted after 1+ week stable | Restore from git history (old services still in git) | <4 hours |

### Post-Migration Safety Net

Keep old service directories in git history (do not force-push or rebase). If a critical regression is discovered after cleanup, the old architecture can be restored from git.

---

## Implementation Notes

### On Spec-Kit Usage

**Agents must create spec-kit specs DURING implementation, not before.**

The `/specify` commands documented in each phase are meant to be run at the START of that phase's implementation session. Do not batch-create all specs upfront — each spec should be grounded in the current state of the codebase at implementation time.

**Exact spec-kit commands per phase:**

| Phase | Command |
|-------|---------|
| Phase 1 | `/specify "Setup Monolith Foundation"` |
| Phase 2 | `/specify "Migrate Auth Module to Monolith"` |
| Phase 3 | `/specify "Consolidate Databases"` |
| Phase 4 | `/specify "Migrate Transactions Module"` |
| Phase 5a | `/specify "Implement In-Memory Event Bus"` |
| Phase 5b | `/specify "Migrate Gamification Module to Monolith"` |
| Phase 6a | `/specify "Migrate Budgets Module"` |
| Phase 6b | `/specify "Migrate Savings Goals Module"` |
| Phase 6c | `/specify "Migrate Recurring Transactions Module"` |
| Phase 6d | `/specify "Migrate Analytics Module"` |
| Phase 6e | `/specify "Migrate AI Module to Monolith"` (if Option A chosen) |
| Phase 7 | `/specify "Technical Debt Resolution"` |

**All spec artifacts will be created in:** `.claude/specs/refactoring/`

### On the Beta Environment

This repository is a **beta environment** that coexists with the production project. This means:
- **Downtime is acceptable** — no real users depend on 100% uptime here
- **We can break things freely** — the production project is separate
- **The only sensitive asset is the database data** — but since we create a new unified DB, the original databases remain as natural backups
- **Feature flags are a nice-to-have**, not a critical safety requirement

### On the AI Module (Phase 6 Decision)

The AI module decision (Option A: migrate to monolith vs Option B: keep as service) should be deferred until Phase 6 begins. Evaluate at that time:
- Is the current `ai-service` causing issues?
- Is the 25-45ms proxy overhead meaningful given 2000-3000ms OpenAI calls?
- Has OpenAI latency become a bottleneck for other operations?

If the answer to any of these is "no," choose Option A (migrate) for simplicity.

### On the Gamification Event Bus (Phase 5)

The current `GamificationProxy` makes a synchronous HTTP call that blocks the `POST /expenses` response until gamification processing completes. The event bus makes this asynchronous (non-blocking). This changes observable behavior:

```
Current (synchronous):
  POST /expenses → expense saved → gamification updated → response returned
  (XP is always updated before the client receives 201)

New (asynchronous via event bus):
  POST /expenses → expense saved → response returned → gamification updated (async)
  (XP may not be updated the instant the client receives 201 — lag of <100ms)
```

This is an acceptable tradeoff. The current behavior is only coincidentally "consistent" because the HTTP call happens to complete before the response. The new behavior is correct: the expense is created successfully, and XP is an eventual side effect.

If strict consistency is required for XP (e.g., the client immediately reads gamification stats after creating an expense), the gamification update can be done synchronously before the event bus publish, inside the same DB transaction.

---

**Document Maintenance:**
- **Next Review**: After Phase 1 completion
- **Owner**: Backend Team
- **Last Updated**: 2026-02-10

**References:**
- `docs/03-architecture/01-current-state.md` — Current distributed monolith analysis
- `docs/03-architecture/02-target-state.md` — Target modular monolith design
- `docs/06-data-models/01-current-state/domain-models.md` — Go domain model inventory
- `docs/06-data-models/01-current-state/main-db.md` — Current DB schema
