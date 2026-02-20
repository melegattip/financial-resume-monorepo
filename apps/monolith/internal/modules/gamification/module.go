package gamification

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/service"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the gamification module.
type Module struct {
	db      *gorm.DB
	handler *handlers.GamificationHandler
	svc     *service.GamificationService
	logger  zerolog.Logger
	authMW  *middleware.AuthMiddleware
}

// New creates and wires all dependencies for the gamification module.
// AutoMigrate is run here so the tables exist before any request is served.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
	if err := db.AutoMigrate(
		&repository.UserGamificationModel{},
		&repository.AchievementModel{},
		&repository.UserActionModel{},
	); err != nil {
		// Non-fatal: AutoMigrate may fail to reconcile constraint names on existing tables.
		logger.Warn().Err(err).Msg("auto-migrate warning (schema may already be up to date)")
	}

	repo := repository.NewGamificationRepository(db)
	svc := service.NewGamificationService(repo, logger)
	handler := handlers.NewGamificationHandler(svc, logger)

	return &Module{
		db:      db,
		handler: handler,
		svc:     svc,
		logger:  logger,
		authMW:  authMW,
	}
}

// RegisterRoutes adds the gamification endpoints to the provided router group.
// All routes are protected by the JWT auth middleware.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	gamification := r.Group("/gamification")
	gamification.Use(m.authMW.RequireAuth())
	{
		gamification.GET("/profile", m.handler.GetProfile)
		gamification.GET("/stats", m.handler.GetStats)
		gamification.GET("/achievements", m.handler.GetAchievements)
		gamification.GET("/features", m.handler.GetFeatures)
		gamification.GET("/features/:featureKey/access", m.handler.CheckFeatureAccess)
		gamification.GET("/challenges/daily", m.handler.GetDailyChallenges)
		gamification.GET("/challenges/weekly", m.handler.GetWeeklyChallenges)
		gamification.POST("/challenges/progress", m.handler.ProcessChallengeProgress)
		gamification.POST("/actions", m.handler.RecordAction)
	}
	m.logger.Info().Msg("gamification module routes registered")
}

// RegisterSubscribers wires the gamification module to relevant domain events.
func (m *Module) RegisterSubscribers(bus sharedports.EventBus) {
	// Initialize gamification profile when a new user registers.
	bus.Subscribe("user.registered", func(ctx context.Context, event sharedports.Event) error {
		if err := m.svc.InitializeUserGamification(ctx, event.UserID()); err != nil {
			m.logger.Warn().Err(err).Str("user_id", event.UserID()).Msg("failed to initialize gamification for new user")
		}
		return nil
	})

	// Award XP when a new expense is created.
	bus.Subscribe("expense.created", func(ctx context.Context, event sharedports.Event) error {
		if _, err := m.svc.RecordAction(ctx, event.UserID(), "create_expense", event.AggregateID()); err != nil {
			m.logger.Warn().Err(err).Str("user_id", event.UserID()).Msg("failed to record create_expense action")
		}
		return nil
	})

	// Award XP when a new income is created.
	bus.Subscribe("income.created", func(ctx context.Context, event sharedports.Event) error {
		if _, err := m.svc.RecordAction(ctx, event.UserID(), "create_income", event.AggregateID()); err != nil {
			m.logger.Warn().Err(err).Str("user_id", event.UserID()).Msg("failed to record create_income action")
		}
		return nil
	})

	m.logger.Info().Msg("gamification module subscribers registered")
}
