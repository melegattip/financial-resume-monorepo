package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/database"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	apphttp "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/logging"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants"
	authdomain "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets"
	budgetsrepo "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring"
	recurringrepo "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring/repository"
	tenantsrepo "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings"
	savingsrepo "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions"
	txrepo "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/events"
)

func main() {
	// Load .env.local first (local dev overrides), then fall back to .env.
	// godotenv.Load does NOT override vars already set in the environment,
	// so loading .env.local first means its values take precedence over .env.
	_ = godotenv.Load(".env.local")
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup structured logger
	logger := logging.Setup(cfg.LogLevel, cfg.Environment)

	logger.Info().
		Str("port", cfg.ServerPort).
		Str("env", cfg.Environment).
		Msg("starting monolith")

	// Initialize database
	db, err := database.NewPostgresConnection(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer database.Close(db, logger)

	// Run schema migrations for all modules
	runMigrations(db, logger)

	// Initialize event bus
	eventBus := events.NewInMemoryEventBus(logger)
	logger.Info().Msg("event bus initialized")

	// Tenants repository — initialized before auth module to provide TenantCreator and TenantMemberFinder
	tenantsRepo := tenantsrepo.NewGormRepository(db)

	// Setup HTTP
	healthHandler := handlers.NewHealthHandler(db)
	router := apphttp.NewRouter(logger, cfg.CORSAllowOrigins, healthHandler)

	// Register modules
	apiV1 := router.Group("/api/v1")

	authModule := auth.New(db, logger, cfg, eventBus, tenantsRepo, tenantsRepo)
	authModule.RegisterRoutes(apiV1)
	authModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("auth module registered")

	// Auth middleware — created once and shared across all protected modules.
	authMW := authModule.AuthMiddleware()

	// Permission middleware — used by tenants module and future modules (Phase D)
	permMW := middleware.NewPermissionMiddleware(tenantsRepo, logger)

	// Tenants module
	tenantsModule := tenants.New(db, logger, cfg, eventBus, authMW, permMW)
	tenantsModule.RegisterRoutes(apiV1)
	tenantsModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("tenants module registered")

	// Transactions module
	txModule := transactions.New(db, logger, cfg, eventBus, authMW, permMW)
	txModule.RegisterRoutes(apiV1)
	txModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("transactions module registered")

	// Savings module
	savingsModule := savings.New(db, logger, eventBus, authMW, permMW)
	savingsModule.RegisterRoutes(apiV1)
	savingsModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("savings module registered")

	// Gamification module
	gamModule := gamification.New(db, logger, cfg, eventBus, authMW, permMW)
	gamModule.RegisterRoutes(apiV1)
	gamModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("gamification module registered")

	// Recurring transactions module
	recurringModule := recurring.New(db, logger, cfg, eventBus, authMW, permMW)
	recurringModule.RegisterRoutes(apiV1)
	recurringModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("recurring module registered")

	// Budgets module
	budgetsModule := budgets.New(db, logger, cfg, eventBus, authMW, permMW)
	budgetsModule.RegisterRoutes(apiV1)
	budgetsModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("budgets module registered")

	// Analytics module
	analyticsModule := analytics.New(db, logger, cfg, eventBus, authMW, permMW)
	analyticsModule.RegisterRoutes(apiV1)
	logger.Info().Msg("analytics module registered")

	// AI module
	aiModule := ai.New(db, logger, cfg, eventBus)
	aiModule.RegisterRoutes(apiV1)
	aiModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("ai module registered")

	server := apphttp.NewServer(cfg.ServerPort, router, logger)

	// Start server in background
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
	server.Shutdown()
	logger.Info().Msg("monolith stopped")
}

// migrateUsersColumns adds missing columns to the existing users table using
// ADD COLUMN IF NOT EXISTS (idempotent). This is needed because the table was
// created by the old microservice and AutoMigrate cannot reconcile the
// constraint name mismatch (uni_users_email), causing it to roll back.
func migrateUsersColumns(db *gorm.DB, logger zerolog.Logger) {
	cols := []struct{ col, def string }{
		{"phone", "VARCHAR(20)"},
		{"avatar", "VARCHAR(500)"},
		{"is_active", "BOOLEAN NOT NULL DEFAULT TRUE"},
		{"is_verified", "BOOLEAN NOT NULL DEFAULT FALSE"},
		{"email_verification_token", "VARCHAR(500)"},
		{"email_verification_expires", "TIMESTAMPTZ"},
		{"password_reset_token", "VARCHAR(500)"},
		{"password_reset_expires", "TIMESTAMPTZ"},
		{"last_login", "TIMESTAMPTZ"},
		{"failed_login_attempts", "INTEGER NOT NULL DEFAULT 0"},
		{"locked_until", "TIMESTAMPTZ"},
		{"deleted_at", "TIMESTAMPTZ"},
	}
	for _, c := range cols {
		sql := "ALTER TABLE users ADD COLUMN IF NOT EXISTS " + c.col + " " + c.def
		if err := db.Exec(sql).Error; err != nil {
			logger.Warn().Err(err).Str("column", c.col).Msg("users column migration warning")
		}
	}
	logger.Info().Msg("users columns ensured")
}

// runMigrations creates missing tables for all modules.
// Each model is migrated individually so a failure on one table
// (e.g. a constraint name mismatch on an existing table) does not
// prevent the remaining tables from being created.
func runMigrations(db *gorm.DB, logger zerolog.Logger) {
	// Ensure all columns exist on the legacy users table before module AutoMigrate runs
	migrateUsersColumns(db, logger)

	models := []struct {
		name  string
		model any
	}{
		// Auth — users already exists; migrate the companion tables separately
		{"user_two_fa", &authdomain.TwoFA{}},
		{"user_notification_settings", &authdomain.NotificationSettings{}},
		// Transactions
		{"categories", &txrepo.CategoryModel{}},
		{"expenses", &txrepo.ExpenseModel{}},
		{"incomes", &txrepo.IncomeModel{}},
		// Budgets
		{"budgets", &budgetsrepo.BudgetModel{}},
		// Savings
		{"savings_goals", &savingsrepo.SavingsGoalModel{}},
		{"savings_transactions", &savingsrepo.SavingsTransactionModel{}},
		// Recurring transactions
		{"recurring_transactions", &recurringrepo.RecurringTransactionModel{}},
		// Tenants (multi-tenant)
		{"tenants",             &tenantsrepo.TenantModel{}},
		{"tenant_members",      &tenantsrepo.TenantMemberModel{}},
		{"permissions",         &tenantsrepo.PermissionModel{}},
		{"role_permissions",    &tenantsrepo.RolePermissionModel{}},
		{"tenant_invitations",  &tenantsrepo.TenantInvitationModel{}},
		{"audit_logs",          &tenantsrepo.AuditLogModel{}},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m.model); err != nil {
			logger.Warn().Err(err).Str("table", m.name).Msg("auto-migrate warning")
		} else {
			logger.Info().Str("table", m.name).Msg("table ready")
		}
	}

	// Add tenant_id column to existing tables (idempotent)
	database.MigrateTenantColumns(db, logger)

	// Add category_id column to incomes table (idempotent)
	database.MigrateIncomeCategoryColumn(db, logger)

	// Add description column to audit_logs table (idempotent)
	database.MigrateAuditLogDescriptionColumn(db, logger)

	// Create personal tenants for existing users and backfill tenant_id (idempotent)
	database.MigrateExistingDataToTenants(db, logger)

	// Enable PostgreSQL Row-Level Security policies on all tenant-scoped tables (idempotent)
	database.MigrateRLSPolicies(db, logger)
}
