# Implementation Roadmap

**Document Version**: 1.0
**Date**: 2026-02-10
**Status**: Ready for Execution
**Dependencies**: [01-current-state.md](./01-current-state.md), [02-target-state.md](./02-target-state.md), [03-migration-plan.md](./03-migration-plan.md)

---

## Overview

This roadmap translates the migration plan into a concrete, week-by-week execution guide for transitioning Financial Resume from a Distributed Monolith (5 processes in 1 Docker container) to a Modular Monolith (1 Go binary, 1 database).

**Key facts:**
- This is a **beta environment** -- we can iterate rapidly without worrying about downtime or zero-downtime deployments
- **Total estimated duration**: ~11 weeks
- **Cost impact**: $7/month throughout (no infrastructure cost changes)
- **Each phase uses spec-kit**: `/specify` generates the spec, `/plan` produces the technical design, `/tasks` creates ordered tasks, `/implement` executes them
- **Feature flags are optional** (nice-to-have) -- since this is beta, direct cutover is acceptable
- **No monitoring/staging/canary deployments** -- we check Render.com logs directly

**What this document adds beyond `03-migration-plan.md`:**
- Visual timeline with week-by-week breakdown
- Precise pre-conditions, input artifacts, output artifacts, and validation for each phase
- ASCII dependency graph
- Risk mitigation mapped per phase
- Quick start guide to begin Phase 1 immediately

---

## Timeline Visual

```
Week     1    2    3    4    5    6    7    8    9   10   11
       +----+----+----+----+----+----+----+----+----+----+----+
Ph 1   |████|    |    |    |    |    |    |    |    |    |    |  Foundation Setup
Ph 2   |    |████|    |    |    |    |    |    |    |    |    |  Auth Module
Ph 3   |    | ░░ |████|████|    |    |    |    |    |    |    |  DB Consolidation
Ph 4   |    |    |    |    |████|░░░░|    |    |    |    |    |  Transactions Module
Ph 5   |    |    |    |    |░░░░|████|████|    |    |    |    |  Gamification + Event Bus
Ph 6   |    |    |    |    |    |    |████|████|    |    |    |  Remaining Modules
Ph 7   |    |    |    |    |    |    |    |    |████|████|████|  Tech Debt + Cleanup
       +----+----+----+----+----+----+----+----+----+----+----+

Legend:
  ████  = Primary execution window
  ░░░░  = Overlap / partial work (analysis, parallel prep)
```

**Critical path**: Phase 1 -> Phase 2 -> Phase 3 -> Phase 4 -> Phase 5 -> Phase 6 -> Phase 7

**Parallel opportunities:**
- Phase 3 DB analysis can start during Phase 2 (Week 2, marked ░░)
- Phase 5a (Event Bus) can be built in parallel with Phase 4 (Week 5, marked ░░)
- Phase 6 sub-modules (Budgets, Savings, Recurring, Analytics, AI) can be developed in parallel within Weeks 7-8

---

## Phase Details

### Phase 1: Monolith Foundation Setup (Week 1)

**Estimated effort**: 5 days

**Spec-Kit Workflow:**
1. `/specify "Setup Monolith Foundation"` -- generates `.claude/specs/refactoring/01-setup-monolith.md`
2. `/plan` -- generates plan with directory structure, dependency decisions, config approach
3. `/tasks` -- generates ordered task list for skeleton creation
4. `/implement` -- executes: create directories, write `main.go`, event bus, DB connection, Dockerfile, CI/CD

**Pre-conditions:**
- [x] Architecture docs complete (`01-current-state.md`, `02-target-state.md`, `03-migration-plan.md`)
- [ ] Full backup of both databases (`financial_resume` and `users_db`)
- [ ] Render.com deployment procedure understood

**Input artifacts (files/code to read):**
- `apps/api-gateway/cmd/api/main.go` -- understand current bootstrap pattern
- `apps/api-gateway/internal/infrastructure/database/` -- current GORM setup
- `apps/api-gateway/internal/infrastructure/router/` -- current Gin router setup
- `infrastructure/docker/supervisord.conf` -- understand current process management
- `Dockerfile` -- understand current multi-binary build
- `render.yaml` -- current Render.com configuration
- `docs/03-architecture/02-target-state.md` -- target directory layout (Appendix A)

**Output artifacts (files/code created):**
```
apps/monolith/
  cmd/server/main.go                      # Entry point with graceful shutdown
  internal/
    shared/
      events/event_bus.go                  # InMemoryEventBus implementation
      events/events.go                     # BaseEvent, Event interface
      events/event_bus_test.go             # Unit tests for event bus
      ports/repositories.go                # Shared repository interfaces
      ports/services.go                    # Shared service interfaces
      errors/errors.go                     # Shared error types
    infrastructure/
      database/postgres.go                 # Single DB connection + pooling
      config/config.go                     # Viper-based configuration
      config/features.go                   # FeatureFlags struct
      http/router.go                       # Gin router with middleware
      http/server.go                       # HTTP server setup
      middleware/auth.go                   # JWT validation middleware
      middleware/cors.go                   # CORS middleware
      middleware/logging.go                # Structured logging middleware
  Dockerfile                               # Single-binary Docker build
  go.mod                                   # Go module definition
  go.sum
.github/workflows/ci.yml                  # GitHub Actions: lint, test, build
```

**Validation:**
- [ ] `go build ./apps/monolith/...` compiles without errors
- [ ] `go test ./apps/monolith/...` passes (event bus unit tests)
- [ ] `go test -race ./apps/monolith/...` detects no race conditions
- [ ] `GET /health` returns `{"status": "ok"}` when binary runs locally
- [ ] Docker image builds successfully with simplified Dockerfile
- [ ] CI/CD pipeline runs green on push

---

### Phase 2: Auth Module Migration (Week 2)

**Estimated effort**: 5 days

**Spec-Kit Workflow:**
1. `/specify "Migrate Auth Module to Monolith"` -- generates `.claude/specs/refactoring/02-auth-module.md`
2. `/plan` -- generates plan covering User domain model copy, JWT compatibility, response parity
3. `/tasks` -- ordered task list: domain -> repository -> use cases -> handlers -> tests
4. `/implement` -- executes the migration

**Pre-conditions:**
- Phase 1 complete: monolith skeleton compiles and runs
- Database connection pool working (verified in Phase 1)
- JWT middleware functional in monolith (verified in Phase 1)

**Input artifacts (files/code to read):**
- `apps/users-service/internal/core/domain/` -- User entity, password hashing
- `apps/users-service/internal/core/usecases/` -- register, login, profile use cases
- `apps/users-service/internal/infrastructure/http/` -- Gin handlers for auth endpoints
- `apps/users-service/internal/infrastructure/database/` -- GORM user repository
- `apps/api-gateway/internal/infrastructure/proxy/users_service_proxy.go` -- proxy to understand routing
- `apps/api-gateway/internal/infrastructure/http/middleware/` -- JWT middleware to replicate
- `docs/04-api-contracts/README.md` -- API contract for auth endpoints

**Output artifacts (files/code created):**
```
apps/monolith/internal/modules/auth/
  domain/
    user.go                                # User aggregate (copied + adapted)
    password.go                            # Password value object
    email.go                               # Email value object
    events.go                              # UserRegisteredEvent, UserLoggedInEvent
  usecases/
    register_user.go                       # Registration flow
    login_user.go                          # JWT generation
    update_profile.go                      # Profile CRUD
    contracts.go                           # Use case interfaces
  handlers/
    auth_handlers.go                       # POST /auth/register, POST /auth/login
    user_handlers.go                       # GET/PUT /users/:id
    routes.go                              # Route registration
  repository/
    user_repository.go                     # GORM implementation
    contracts.go                           # Repository interfaces
  module.go                                # Module init + DI

scripts/validate_auth_parity.sh            # Response parity test script
```

**Validation:**
- [ ] `POST /api/v1/auth/register` returns structurally identical response to old `users-service`
- [ ] `POST /api/v1/auth/login` returns valid JWT accepted by frontend
- [ ] `GET /api/v1/users/:id` returns user profile
- [ ] `PUT /api/v1/users/:id` updates profile successfully
- [ ] JWT tokens issued by new module pass validation in auth middleware
- [ ] Frontend login/register flow works without any code changes
- [ ] Integration tests cover all auth endpoints
- [ ] Response parity script (`validate_auth_parity.sh`) shows zero structural differences

---

### Phase 3: Database Consolidation (Weeks 3-4)

**Estimated effort**: 10 days (this is the most critical phase -- handle with care)

**Spec-Kit Workflow:**
1. `/specify "Consolidate Databases"` -- generates `.claude/specs/refactoring/03-db-consolidation.md`
2. `/plan` -- generates plan with exact SQL migration scripts, data audit queries, validation queries
3. `/tasks` -- ordered task list: audit -> schema changes -> data copy -> dedup -> validate -> cutover
4. `/implement` -- executes migration SQL step by step with validation at each stage

**Pre-conditions:**
- Phase 2 complete: auth module functional, `users-service` proxy still available as fallback
- Full backup of both databases verified (can be restored)
- Database credentials for both `financial_resume` and `users_db` accessible

**Input artifacts (files/code to read):**
- `docs/06-data-models/01-current-state/main-db.md` -- current `financial_resume` schema
- `docs/06-data-models/01-current-state/gamification-db.md` -- current gamification tables
- `docs/03-architecture/02-target-state.md` (Database Consolidation section) -- target consolidated schema
- `docs/03-architecture/03-migration-plan.md` (Phase 3 section) -- detailed SQL migration steps
- Live database: run audit queries to discover actual data types and counts

**Output artifacts (files/code created):**
```
apps/monolith/internal/infrastructure/database/
  migrations/
    001_add_deleted_at_columns.sql          # Add soft delete to all tables
    002_create_users_table.sql              # Create users table in financial_resume
    003_create_user_preferences_table.sql   # Create user_preferences table
    004_copy_users_data.sql                 # Copy + cast user data
    005_copy_user_preferences_data.sql      # Copy preferences
    006_dedup_user_gamification.sql         # Resolve duplicate gamification rows
    007_add_composite_indexes.sql           # Performance indexes
    008_add_foreign_key_constraints.sql     # FK constraints after data is clean

scripts/
  audit_databases.sql                       # Pre-migration audit queries
  validate_migration.sql                    # Post-migration validation queries
  validate_migration.sh                     # Automated count-check script
```

**Validation:**
- [ ] Pre-migration audit complete: user counts, type inventory, duplicate detection
- [ ] `deleted_at` column present on all entity tables (expenses, incomes, budgets, savings_goals, recurring_transactions, user_gamification, categories)
- [ ] All users from `users_db` present in `financial_resume.users` (count match)
- [ ] All `user_preferences` migrated (count match)
- [ ] `user_gamification` has exactly 1 row per `user_id` (dedup verified)
- [ ] Zero orphaned foreign key references (expenses without valid user_id = 0)
- [ ] All `user_id` columns are `VARCHAR(255)` (no remaining `INT`/`uint`)
- [ ] Auth module still works correctly with unified database
- [ ] Composite indexes created and query plans verified
- [ ] `users_db` Render instance can be deleted (free tier slot freed)

---

### Phase 4: Transactions Module Migration (Weeks 5-6, first half)

**Estimated effort**: 7 days

**Spec-Kit Workflow:**
1. `/specify "Migrate Transactions Module"` -- generates `.claude/specs/refactoring/04-transactions-module.md`
2. `/plan` -- generates plan for Expense/Income/Category migration with soft delete and event publishing
3. `/tasks` -- ordered: domain models -> repositories (with soft delete) -> use cases -> handlers -> event publishing -> tests
4. `/implement` -- executes the migration

**Pre-conditions:**
- Phase 3 complete: single unified database with `deleted_at` columns
- Auth module working (used for JWT validation on all transaction endpoints)
- Event bus available from Phase 1 (even if no subscribers yet)

**Input artifacts (files/code to read):**
- `apps/api-gateway/internal/core/domain/` -- Expense, Income, Category entities, builders, factories
- `apps/api-gateway/internal/core/usecases/` -- CRUD use cases for transactions
- `apps/api-gateway/internal/infrastructure/http/handlers/` -- Gin handlers
- `apps/api-gateway/internal/infrastructure/database/` -- GORM repositories
- `apps/api-gateway/internal/core/ports/` -- repository and service interfaces
- `docs/04-api-contracts/README.md` -- API contracts for expense/income endpoints

**Output artifacts (files/code created):**
```
apps/monolith/internal/modules/transactions/
  domain/
    expense.go                             # Expense aggregate (with deleted_at)
    income.go                              # Income aggregate (with deleted_at)
    category.go                            # Category entity
    builders.go                            # ExpenseBuilder, IncomeBuilder
    factory.go                             # TransactionFactory
    events.go                              # ExpenseCreatedEvent, IncomeCreatedEvent, etc.
  usecases/
    expenses/
      create.go, update.go, delete.go, list.go, contracts.go
    incomes/
      create.go, update.go, delete.go, list.go, contracts.go
    categories/
      manage.go, contracts.go
  handlers/
    expense_handlers.go                    # POST/GET/PUT/DELETE /expenses
    income_handlers.go                     # POST/GET/PUT/DELETE /incomes
    category_handlers.go                   # GET/POST /categories
    routes.go
  repository/
    expense_repository.go                  # GORM with soft delete filtering
    income_repository.go
    category_repository.go
    contracts.go
  module.go                                # Module init + DI + route registration
```

**Validation:**
- [ ] All `/api/v1/expenses/*` endpoints return identical responses to current api-gateway
- [ ] All `/api/v1/incomes/*` endpoints return identical responses
- [ ] All `/api/v1/categories/*` endpoints return identical responses
- [ ] Soft delete works: deleted expenses return 404 but row exists in DB with `deleted_at` set
- [ ] `ExpenseCreatedEvent` and `IncomeCreatedEvent` are published to event bus (verified in logs)
- [ ] Partial payment tracking works correctly for expenses
- [ ] Transaction date vs created_at distinction preserved
- [ ] Frontend expense/income flows work without modification

---

### Phase 5: Gamification Module + Event Bus (Weeks 5-7)

**Estimated effort**: 10 days (split across two sub-phases)

This phase has two spec-kit specs because the event bus and gamification module are independent concerns.

#### Sub-Phase 5a: In-Memory Event Bus (3 days, can overlap with Phase 4)

**Spec-Kit Workflow:**
1. `/specify "Implement In-Memory Event Bus"` -- generates `.claude/specs/refactoring/05a-event-bus.md`
2. `/plan` -- generates design for thread-safe pub/sub, error handling, goroutine management
3. `/tasks` -- ordered: Event interface -> InMemoryEventBus -> subscribe/publish -> error logging -> tests
4. `/implement` -- builds and tests the event bus

**Pre-conditions:**
- Phase 1 complete: shared infrastructure directory exists
- Event bus interface defined in `shared/ports/`

**Input artifacts (files/code to read):**
- `docs/03-architecture/02-target-state.md` (Event Bus Architecture section) -- design and code examples
- `docs/03-architecture/03-migration-plan.md` (Phase 5a section) -- implementation details
- `apps/monolith/internal/shared/events/events.go` -- BaseEvent definition from Phase 1

**Output artifacts (files/code created):**
```
apps/monolith/internal/shared/events/
  event_bus.go                             # InMemoryEventBus (thread-safe, async handlers)
  events.go                                # Event interface, BaseEvent, concrete event types
  event_bus_test.go                        # Unit tests: subscribe, publish, concurrent publish, error handling

apps/monolith/internal/shared/ports/
  event_bus.go                             # EventBus interface, EventHandler type
```

**Validation:**
- [ ] `go test ./apps/monolith/internal/shared/events/...` passes
- [ ] `go test -race ./apps/monolith/internal/shared/events/...` detects no race conditions
- [ ] Event bus handles 1000 concurrent publishes without goroutine leak
- [ ] Failed handler logs error but does not crash publisher
- [ ] Subscribe/publish for multiple event types works independently

#### Sub-Phase 5b: Gamification Module Migration (7 days)

**Spec-Kit Workflow:**
1. `/specify "Migrate Gamification Module to Monolith"` -- generates `.claude/specs/refactoring/05b-gamification-module.md`
2. `/plan` -- generates plan for domain copy, event subscription wiring, XP rule verification
3. `/tasks` -- ordered: domain models -> repository -> use cases -> event handlers -> HTTP handlers -> tests
4. `/implement` -- executes the migration

**Pre-conditions:**
- Phase 4 complete: transactions module publishes `ExpenseCreatedEvent` and `IncomeCreatedEvent`
- Sub-Phase 5a complete: event bus is tested and functional
- Phase 3 complete: `user_gamification` table is deduplicated in unified DB

**Input artifacts (files/code to read):**
- `apps/gamification-service/internal/core/domain/` -- UserGamification, Achievement, Challenge entities
- `apps/gamification-service/internal/core/usecases/` -- XP calculation, achievement unlock, streak tracking
- `apps/gamification-service/internal/infrastructure/http/` -- Gin handlers
- `apps/gamification-service/internal/infrastructure/database/` -- GORM repositories
- `apps/api-gateway/internal/infrastructure/proxy/gamification_proxy.go` -- understand current proxy calls
- `docs/03-architecture/02-target-state.md` (Gamification Module section) -- XP rules, event subscriptions

**Output artifacts (files/code created):**
```
apps/monolith/internal/modules/gamification/
  domain/
    user_gamification.go                   # UserGamification aggregate
    achievement.go                         # Achievement definitions
    challenge.go                           # Challenge definitions
    xp_calculator.go                       # XP domain service
    events.go                              # LevelUpEvent, AchievementUnlockedEvent
  usecases/
    add_xp.go                              # Triggered by domain events
    unlock_achievement.go                  # Achievement check logic
    update_streak.go                       # Streak calculation
    contracts.go
  handlers/
    gamification_handlers.go               # GET /gamification/stats, achievements
    routes.go
  repository/
    gamification_repository.go             # GORM for user_gamification
    achievement_repository.go              # GORM for achievements, user_achievements
    contracts.go
  module.go                                # Module init + RegisterSubscribers()
```

**Validation:**
- [ ] Creating an expense triggers `ExpenseCreatedEvent` -> gamification adds +5 XP (end-to-end)
- [ ] XP totals match between old `gamification-service` and new module for identical action sequences
- [ ] Achievement unlock logic produces same results as old service
- [ ] Streak calculation is correct (increments on consecutive day login, resets on gap)
- [ ] `GET /api/v1/gamification/stats` returns identical response structure
- [ ] `GET /api/v1/gamification/achievements` returns identical data
- [ ] `GamificationProxy` HTTP localhost call eliminated from request path

---

### Phase 6: Remaining Modules Migration (Weeks 7-8)

**Estimated effort**: 10 days (4 sub-modules can be developed in parallel)

Each sub-module has its own spec-kit spec. They can be developed in parallel within Weeks 7-8 but must be integrated together before final validation.

#### 6a: Budgets Module (3 days)

**Spec-Kit Workflow:**
1. `/specify "Migrate Budgets Module"` -- generates `.claude/specs/refactoring/06a-budgets-module.md`
2. `/plan` -> `/tasks` -> `/implement`

**Pre-conditions:**
- Phase 4 complete: transactions module publishes `ExpenseCreatedEvent`
- Phase 5 complete: event bus and gamification subscriptions working

**Input artifacts:**
- `apps/api-gateway/internal/core/domain/budget.go` -- Budget aggregate with state machine
- `apps/api-gateway/internal/core/usecases/budgets/` -- Budget CRUD and status recalculation
- `apps/api-gateway/internal/infrastructure/http/handlers/budget_*.go`
- `apps/api-gateway/internal/infrastructure/database/budget_repository.go`

**Output artifacts:**
```
apps/monolith/internal/modules/budgets/
  domain/
    budget.go                              # Budget aggregate + state machine
    budget_period.go                       # Monthly/weekly/yearly value object
    budget_status.go                       # on_track, warning, exceeded enum
    calculator.go                          # Domain service for recalculation
    events.go                              # BudgetExceededEvent, BudgetWarningEvent
  usecases/
    create_budget.go, update_budget.go, calculate_status.go, reset_budget.go, contracts.go
  handlers/
    budget_handlers.go, routes.go
  repository/
    budget_repository.go, contracts.go
  module.go                                # RegisterSubscribers: ExpenseCreatedEvent -> recalculate
```

**Validation:**
- [ ] Budget status auto-updates when expense is created (on_track -> warning -> exceeded)
- [ ] Budget CRUD endpoints return identical responses
- [ ] `BudgetExceededEvent` published when budget crosses 100%
- [ ] Budget reset cron job can be triggered manually

#### 6b: Savings Goals Module (2 days)

**Spec-Kit Workflow:**
1. `/specify "Migrate Savings Goals Module"` -- generates `.claude/specs/refactoring/06b-savings-module.md`
2. `/plan` -> `/tasks` -> `/implement`

**Pre-conditions:**
- Phase 4 complete (income events available for optional auto-save)
- Phase 5 complete (event bus functional)

**Input artifacts:**
- `apps/api-gateway/internal/core/domain/savings_goal.go` -- SavingsGoal aggregate + state machine
- `apps/api-gateway/internal/core/usecases/savings/`
- `apps/api-gateway/internal/infrastructure/http/handlers/savings_*.go`
- `apps/api-gateway/internal/infrastructure/database/savings_repository.go`

**Output artifacts:**
```
apps/monolith/internal/modules/savings/
  domain/
    savings_goal.go, goal_status.go, progress_calculator.go, events.go
  usecases/
    create_goal.go, add_contribution.go, withdraw.go, contracts.go
  handlers/
    goal_handlers.go, routes.go
  repository/
    savings_repository.go, contracts.go
  module.go
```

**Validation:**
- [ ] Savings goal CRUD endpoints return identical responses
- [ ] Goal auto-achieves when `current_amount >= target_amount`
- [ ] `GoalAchievedEvent` published and gamification awards +50 XP
- [ ] Goal status transitions: active -> achieved, active -> paused, active -> cancelled

#### 6c: Recurring Transactions Module (2 days)

**Spec-Kit Workflow:**
1. `/specify "Migrate Recurring Transactions Module"` -- generates `.claude/specs/refactoring/06c-recurring-module.md`
2. `/plan` -> `/tasks` -> `/implement`

**Pre-conditions:**
- Phase 4 complete (transactions module available to create expenses/incomes on schedule)
- Phase 1 complete (cron scheduler infrastructure)

**Input artifacts:**
- `apps/api-gateway/internal/core/domain/recurring_transaction.go` -- RecurringTransaction aggregate
- `apps/api-gateway/internal/infrastructure/scheduler/` -- existing cron scheduler
- `apps/api-gateway/cmd/api/main.go` -- current scheduler initialization (runs every 1 hour)

**Output artifacts:**
```
apps/monolith/internal/modules/recurring/
  domain/
    recurring_transaction.go, frequency.go, scheduler.go, events.go
  usecases/
    create_recurring.go, execute_due.go, pause_resume.go, contracts.go
  handlers/
    recurring_handlers.go, routes.go
  repository/
    recurring_repository.go, contracts.go
  module.go

apps/monolith/internal/infrastructure/cron/
  scheduler.go                             # robfig/cron scheduler
  jobs.go                                  # Registered jobs: recurring execution, budget reset
```

**Validation:**
- [ ] Recurring transaction CRUD endpoints return identical responses
- [ ] Cron job executes due transactions on schedule
- [ ] No missed executions during migration window (compare execution counts before/after)
- [ ] Pause/resume works correctly
- [ ] `RecurringExecutedEvent` creates actual expense/income in transactions module

#### 6d: Analytics Module (2 days)

**Spec-Kit Workflow:**
1. `/specify "Migrate Analytics Module"` -- generates `.claude/specs/refactoring/06d-analytics-module.md`
2. `/plan` -> `/tasks` -> `/implement`

**Pre-conditions:**
- All other modules (auth, transactions, budgets, savings, gamification) are functional
- Event bus wired with all subscriptions

**Input artifacts:**
- `apps/api-gateway/internal/core/usecases/dashboard/` -- dashboard aggregation logic
- `apps/api-gateway/internal/infrastructure/http/handlers/dashboard_*.go`
- `apps/api-gateway/internal/infrastructure/http/handlers/analytics_*.go`

**Output artifacts:**
```
apps/monolith/internal/modules/analytics/
  domain/
    dashboard_data.go, report.go, trend_calculator.go
  usecases/
    get_dashboard.go, get_expense_analytics.go, contracts.go
  handlers/
    dashboard_handlers.go, analytics_handlers.go, routes.go
  repository/
    analytics_repository.go, contracts.go
  module.go                                # RegisterSubscribers for cache invalidation
```

**Validation:**
- [ ] `GET /api/v1/dashboard` returns correct data from all modules
- [ ] `GET /api/v1/expenses/analytics` returns correct aggregations
- [ ] Dashboard uses parallel goroutine queries for independent data sources
- [ ] Response structure identical to current api-gateway responses

#### 6e: AI Module (1 day -- decision point)

**Spec-Kit Workflow:**
1. `/specify "Migrate AI Module to Monolith"` -- generates `.claude/specs/refactoring/06e-ai-module.md`
2. `/plan` -> `/tasks` -> `/implement`

**Decision at start of Phase 6**: Migrate into monolith (Option A, recommended) or keep as separate service (Option B). See `03-migration-plan.md` for decision criteria. Default: Option A.

**Pre-conditions (if Option A):**
- All other modules functional
- OpenAI API key accessible from monolith config

**Input artifacts:**
- `apps/ai-service/internal/` -- AI service domain, use cases, handlers
- `apps/api-gateway/internal/infrastructure/proxy/ai_service_proxy.go` -- understand current proxy
- OpenAI integration code

**Output artifacts (if Option A):**
```
apps/monolith/internal/modules/ai/
  domain/
    health_analysis.go, recommendation.go, events.go
  usecases/
    analyze_health.go, generate_insights.go, contracts.go
  handlers/
    ai_handlers.go, routes.go
  infrastructure/
    openai_client.go                       # OpenAI SDK wrapper
  module.go
```

**Validation:**
- [ ] `POST /api/v1/ai/analyze` returns health analysis (same response structure)
- [ ] OpenAI API calls work correctly through monolith
- [ ] `context.WithTimeout` wraps AI calls (10s max)
- [ ] AI request latency is not worse than before (2000-3000ms is expected, not a regression)

---

### Phase 7: Technical Debt Resolution + Cleanup (Weeks 9-11)

**Estimated effort**: 15 days

**Spec-Kit Workflow:**
1. `/specify "Technical Debt Resolution"` -- generates `.claude/specs/refactoring/07-tech-debt.md`
2. `/plan` -- generates plan covering soft delete audit, ACID transactions, Dockerfile simplification, old code removal
3. `/tasks` -- ordered: soft delete audit -> ACID patterns -> Dockerfile -> render.yaml -> remove old services -> final validation
4. `/implement` -- executes cleanup

**Pre-conditions:**
- ALL modules (auth, transactions, budgets, savings, recurring, gamification, analytics, AI) migrated and functional
- All event bus subscriptions wired and tested
- No regressions observed for at least a few days of running

**Input artifacts:**
- All module code created in Phases 2-6
- `Dockerfile` -- current multi-binary build
- `render.yaml` -- current multi-service configuration
- `infrastructure/docker/supervisord.conf` -- to be removed
- `apps/api-gateway/internal/infrastructure/proxy/` -- all proxy files to be deleted
- `apps/users-service/`, `apps/gamification-service/`, `apps/ai-service/` -- directories to be deleted

**Output artifacts:**

Sub-Phase 7.1 -- Soft Delete Audit:
```
# Verify every repository in every module filters by deleted_at IS NULL
# Verify GORM soft delete is configured on all models
# Fix any repositories that still do hard delete
```

Sub-Phase 7.2 -- ACID Transaction Patterns:
```
# Update use cases that need cross-module consistency to use DB transactions
# Pattern: tx := db.Begin() -> operations -> tx.Commit()
# Event bus publish happens AFTER commit (not inside transaction)
```

Sub-Phase 7.3 -- Dockerfile + Deployment:
```
apps/monolith/Dockerfile                   # Simplified single-binary build
render.yaml                                # Updated: 1 web service + 1 static site + 1 DB
```

Sub-Phase 7.4 -- Old Code Removal:
```
DELETED:
  apps/api-gateway/internal/infrastructure/proxy/gamification_proxy.go
  apps/api-gateway/internal/infrastructure/proxy/users_service_proxy.go
  apps/api-gateway/internal/infrastructure/proxy/ai_service_proxy.go
  apps/users-service/                      # Entire directory
  apps/gamification-service/               # Entire directory
  apps/ai-service/                         # Entire directory (if Option A in Phase 6e)
  infrastructure/docker/supervisord.conf
  infrastructure/docker/nginx.conf         # Internal proxy config
```

Sub-Phase 7.5 -- Feature Flag Removal (if used):
```
# Remove all conditional routing (if flags.UseMonolith* { ... })
# Keep only the new module code paths
# Delete infrastructure/config/features.go
```

**Validation:**
- [ ] Single Go binary running successfully
- [ ] supervisord removed from container
- [ ] nginx (internal proxy) removed from container
- [ ] Old service directories deleted (`users-service`, `gamification-service`, `ai-service`)
- [ ] Proxy files deleted
- [ ] Docker image builds and runs correctly
- [ ] `render.yaml` reflects simplified architecture (1 web service, 1 DB)
- [ ] All API endpoints respond correctly (full E2E test pass)
- [ ] Soft delete working on all entities
- [ ] ACID transactions enforced where needed
- [ ] Cron jobs continue executing
- [ ] Frontend works without any changes

---

## Dependency Graph (ASCII)

```
                    +-------------------+
                    |   Phase 1         |
                    |   Foundation      |
                    |   (Week 1)        |
                    +--------+----------+
                             |
                             v
                    +-------------------+
                    |   Phase 2         |
                    |   Auth Module     |
                    |   (Week 2)        |
                    +--------+----------+
                             |
                             v
                    +-------------------+
                    |   Phase 3         |
                    |   DB Consolidation|
                    |   (Weeks 3-4)     |
                    +--------+----------+
                             |
                +------------+------------+
                |                         |
                v                         v
       +----------------+       +------------------+
       |  Phase 4       |       |  Phase 5a        |
       |  Transactions  |       |  Event Bus       |
       |  (Week 5)      |       |  (Week 5)        |
       +-------+--------+       +--------+---------+
               |                          |
               +--------+---------+-------+
                        |
                        v
               +------------------+
               |  Phase 5b        |
               |  Gamification    |
               |  (Week 6)        |
               +--------+---------+
                        |
                        v
   +----------+---------+--------+---------+---------+
   |          |         |        |         |         |
   v          v         v        v         v         |
+------+ +-------+ +-------+ +------+ +------+      |
|  6a  | |  6b   | |  6c   | |  6d  | |  6e  |      |
|Budget| |Savings| |Recur. | |Analyt| | AI   |      |
|3 days| |2 days | |2 days | |2 days| |1 day |      |
+--+---+ +--+----+ +--+----+ +--+---+ +--+---+      |
   |        |         |         |         |          |
   +--------+---------+---------+---------+          |
                        |                            |
                        v                            |
               +-------------------+                 |
               |   Phase 7         |<----------------+
               |   Tech Debt +     |
               |   Cleanup         |
               |   (Weeks 9-11)    |
               +-------------------+

LEGEND:
  ----->  Sequential dependency (must complete before next starts)
  Phase 4 and Phase 5a can run in PARALLEL (both depend on Phase 3)
  Phase 6a-6e can run in PARALLEL (all depend on Phase 5b)
  Phase 7 depends on ALL previous phases
```

---

## Risk Mitigation per Phase

| Phase | Risk | Probability | Impact | Mitigation |
|-------|------|-------------|--------|------------|
| **Phase 1** | Go module dependency conflicts between monolith and existing services | Low | Low | Use separate `go.mod` in `apps/monolith/`; shared dependencies pinned to same versions |
| **Phase 1** | CI/CD pipeline fails on Render.com | Low | Low | Test Docker build locally first; Render.com has straightforward Docker deployment |
| **Phase 2** | JWT tokens incompatible between old and new auth | Medium | High | Use **identical** JWT secret and signing algorithm; test token interop before cutover |
| **Phase 2** | Password hashing algorithm mismatch | Low | Critical | Copy exact bcrypt/argon2 implementation from `users-service`; verify existing passwords validate |
| **Phase 3** | Data loss during `users_db` -> `financial_resume` copy | Low | Critical | Original `users_db` is **never modified** -- it remains as natural backup; use `INSERT ... ON CONFLICT DO NOTHING`; validate row counts |
| **Phase 3** | `user_id` type cast (`INT` -> `VARCHAR`) breaks FK references | Medium | High | Run audit queries FIRST (Phase 3.1) to discover all mismatches; test cast on staging data before production |
| **Phase 3** | Gamification dedup loses the wrong record | Low | Medium | Dedup keeps row with highest `total_xp`; manual review of dedup results before applying `DELETE` |
| **Phase 4** | Soft delete breaks existing queries that expect hard delete | Medium | Medium | Audit ALL queries to include `WHERE deleted_at IS NULL`; GORM's built-in soft delete handles this automatically if `DeletedAt gorm.DeletedAt` field is used |
| **Phase 4** | Event publishing introduces latency on expense creation | Low | Low | Event bus `Publish()` is non-blocking (fires goroutines); main request returns immediately |
| **Phase 5** | Event bus goroutine leak under high concurrency | Low | Medium | Load test with 1000 concurrent events; use `go test -race`; monitor goroutine count |
| **Phase 5** | XP calculation rules diverge between old and new | Medium | Medium | Write comparison test: replay same action sequence through both old service and new module, compare final XP |
| **Phase 6** | Cron job misses executions during migration window | Low | Medium | Overlap window: both old and new scheduler ready but only one active; verify execution count |
| **Phase 6** | Dashboard aggregation returns different numbers | Medium | Medium | Snapshot current dashboard response for test users; compare field-by-field with new module |
| **Phase 6** | AI module latency affects monolith p99 | Medium | Low | Wrap OpenAI calls in `context.WithTimeout(ctx, 10*time.Second)`; AI calls already take 2-3s, this is not a regression |
| **Phase 7** | Deleting old service code breaks something still in use | Low | High | Confirm zero traffic to old code paths via logs for 1+ week before deletion; keep git history as safety net |
| **Phase 7** | Simplified Dockerfile fails on Render.com | Low | Medium | Test Docker build and run locally before deploying; Render.com supports standard Docker |

---

## Definition of Done

The migration is considered **complete** when ALL of the following are true:

### Architecture
- [ ] Single Go binary running in production (no supervisord, no nginx internal proxy)
- [ ] Single PostgreSQL database (`financial_resume`) -- `users_db` instance deleted
- [ ] All 8 modules functional: auth, transactions, budgets, savings, recurring, gamification, ai, analytics
- [ ] In-memory event bus wiring all cross-module communication
- [ ] Old service directories removed from repository (`users-service`, `gamification-service`, `ai-service`)
- [ ] Proxy files removed (`gamification_proxy.go`, `users_service_proxy.go`, `ai_service_proxy.go`)

### API Compatibility
- [ ] 100% backward compatible -- all existing API contracts preserved
- [ ] Zero frontend changes required -- React app works without modification
- [ ] Response JSON structures identical to pre-migration
- [ ] JWT authentication unchanged -- existing tokens still valid
- [ ] HTTP status codes unchanged

### Data Integrity
- [ ] Zero data loss -- all user data migrated correctly
- [ ] Soft delete implemented on all entity tables
- [ ] `user_id` standardized to `VARCHAR(255)` everywhere
- [ ] No orphaned foreign key references
- [ ] `user_gamification` deduplicated (1 row per user)

### Quality
- [ ] All integration tests passing
- [ ] All E2E tests passing
- [ ] `go test -race` detects no race conditions
- [ ] Cron jobs executing on schedule (recurring transactions, budget reset)

### Performance
- [ ] Dashboard p95 latency < 100ms (vs 150-200ms before)
- [ ] Expense creation p95 latency < 60ms (vs 90-140ms before)
- [ ] Docker image size reduced (single binary vs 4 binaries + nginx)

### Operations
- [ ] Simplified Dockerfile (single-stage build, single binary)
- [ ] Updated `render.yaml` (1 web service, 1 static site, 1 database)
- [ ] CI/CD pipeline green
- [ ] Documentation updated to reflect new architecture

---

## Quick Start Guide

Ready to begin Phase 1? Follow these steps:

### Step 1: Verify Prerequisites

```bash
# Ensure Go 1.23+ is installed
go version

# Verify you're in the monorepo root
pwd
# Expected: .../financial-resume-monorepo

# Verify the architecture docs exist
ls docs/03-architecture/
# Expected: 01-current-state.md  02-target-state.md  03-migration-plan.md  04-implementation-roadmap.md
```

### Step 2: Backup Databases

```bash
# Backup financial_resume database
pg_dump $FINANCIAL_RESUME_DB_URL > backups/financial_resume_$(date +%Y%m%d).sql

# Backup users_db database
pg_dump $USERS_DB_URL > backups/users_db_$(date +%Y%m%d).sql

# Verify backups are not empty
ls -la backups/
```

### Step 3: Start Phase 1 with Spec-Kit

```bash
# Open your terminal / Claude Code session
# Run the spec-kit specify command:
/specify "Setup Monolith Foundation"

# This generates:
#   .claude/specs/refactoring/01-setup-monolith.md
#
# Review the generated spec. Does it cover:
#   - apps/monolith/ directory structure
#   - cmd/server/main.go entry point
#   - InMemoryEventBus implementation
#   - Single PostgreSQL connection with pooling
#   - Simplified Dockerfile
#   - Health check endpoint
#   - CI/CD pipeline
#
# If satisfied, proceed:
/plan

# Review the generated plan. Then:
/tasks

# Review the task list. Finally:
/implement
```

### Step 4: Validate Phase 1

```bash
# Build the monolith
cd apps/monolith && go build ./...

# Run tests
go test ./... -v

# Run race detector
go test -race ./...

# Start the server locally
go run cmd/server/main.go

# Test health endpoint
curl http://localhost:8080/health
# Expected: {"status":"ok"}
```

### Step 5: Continue to Phase 2

Once Phase 1 validation passes:

```bash
/specify "Migrate Auth Module to Monolith"
# Review spec -> /plan -> /tasks -> /implement
```

### Repeat for Each Phase

Follow the same spec-kit workflow for each subsequent phase:

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
| Phase 6e | `/specify "Migrate AI Module to Monolith"` |
| Phase 7 | `/specify "Technical Debt Resolution"` |

**All spec artifacts are created in:** `.claude/specs/refactoring/`

---

**Document Maintenance:**
- **Next Review**: After Phase 1 completion
- **Owner**: Backend Team
- **Last Updated**: 2026-02-10

**References:**
- `docs/03-architecture/01-current-state.md` -- Current distributed monolith analysis
- `docs/03-architecture/02-target-state.md` -- Target modular monolith design
- `docs/03-architecture/03-migration-plan.md` -- Detailed migration plan with SQL scripts and code examples
