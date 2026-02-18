package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/database"
	apphttp "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/http/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/logging"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/analytics"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/recurring"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/events"
)

func main() {
	// Load .env file (optional — production uses real env vars)
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

	// Initialize event bus
	eventBus := events.NewInMemoryEventBus(logger)
	logger.Info().Msg("event bus initialized")

	// Setup HTTP
	healthHandler := handlers.NewHealthHandler(db)
	router := apphttp.NewRouter(logger, cfg.CORSAllowOrigins, healthHandler)

	// Register modules
	apiV1 := router.Group("/api/v1")

	authModule := auth.New(db, logger, cfg, eventBus)
	authModule.RegisterRoutes(apiV1)
	authModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("auth module registered")

	// Transactions module (requires JWT auth middleware)
	txModule := transactions.New(db, logger, cfg, eventBus)
	txModule.RegisterRoutes(apiV1)
	txModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("transactions module registered")

	// Savings module
	savingsModule := savings.New(db, logger, eventBus)
	savingsModule.RegisterRoutes(apiV1)
	savingsModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("savings module registered")

	// Gamification module
	authMW := authModule.AuthMiddleware()
	gamModule := gamification.New(db, logger, cfg, eventBus, authMW)
	gamModule.RegisterRoutes(apiV1)
	gamModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("gamification module registered")

	// Recurring transactions module
	recurringModule := recurring.New(db, logger, cfg, eventBus)
	recurringModule.RegisterRoutes(apiV1)
	recurringModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("recurring module registered")

	// Budgets module
	budgetsModule := budgets.New(db, logger, cfg, eventBus)
	budgetsModule.RegisterRoutes(apiV1)
	budgetsModule.RegisterSubscribers(eventBus)
	logger.Info().Msg("budgets module registered")

	// Analytics module
	analyticsModule := analytics.New(db, logger, cfg, eventBus, authMW)
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
