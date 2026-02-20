package savings

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/savings/repository"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the savings module
type Module struct {
	savingsHandler *handlers.SavingsHandler
	logger         zerolog.Logger
	authMW         *middleware.AuthMiddleware
}

// New creates a new savings module
func New(db *gorm.DB, logger zerolog.Logger, eventBus sharedports.EventBus, authMW *middleware.AuthMiddleware) *Module {
	savingsRepo := repository.NewSavingsRepository(db)
	savingsHandler := handlers.NewSavingsHandler(savingsRepo, logger)

	return &Module{
		savingsHandler: savingsHandler,
		logger:         logger,
		authMW:         authMW,
	}
}

// RegisterRoutes registers all HTTP routes for the savings module.
// All routes are under /savings-goals and require a valid JWT access token.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	savings := r.Group("/savings-goals")
	savings.Use(m.authMW.RequireAuth())
	{
		savings.POST("", m.savingsHandler.CreateGoal)
		savings.GET("", m.savingsHandler.ListGoals)
		savings.GET("/dashboard", m.savingsHandler.GetSummary)
		savings.GET("/summary", m.savingsHandler.GetSummary)
		savings.GET("/:id", m.savingsHandler.GetGoal)
		savings.PUT("/:id", m.savingsHandler.UpdateGoal)
		savings.DELETE("/:id", m.savingsHandler.DeleteGoal)
		savings.POST("/:id/deposit", m.savingsHandler.AddSavings)
		savings.POST("/:id/withdraw", m.savingsHandler.WithdrawSavings)
		savings.POST("/:id/pause", m.savingsHandler.PauseGoal)
		savings.POST("/:id/resume", m.savingsHandler.ResumeGoal)
		savings.POST("/:id/cancel", m.savingsHandler.CancelGoal)
		savings.GET("/:id/transactions", m.savingsHandler.ListTransactions)
	}

	m.logger.Info().Msg("savings module routes registered")
}

// RegisterSubscribers registers event subscribers for the savings module
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	// Future: subscribe to UserDeletedEvent to cascade-delete user's savings goals
	m.logger.Info().Msg("savings module subscribers registered")
}
