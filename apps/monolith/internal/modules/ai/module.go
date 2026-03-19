package ai

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/service"
	sharedemail "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/email"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the AI module and all its dependencies.
type Module struct {
	aiHandler *handlers.AIHandler
	logger    zerolog.Logger
}

// New creates and wires the AI module.
// db and cfg are accepted for interface consistency with other modules, but are not used
// because the AI module is stateless (no DB access needed).
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus sharedports.EventBus, emailService sharedemail.EmailService) *Module {
	openaiKey := os.Getenv("OPENAI_API_KEY")

	openaiClient := service.NewOpenAIClient(openaiKey)
	analysisService := service.NewAnalysisService(openaiClient)
	purchaseService := service.NewPurchaseService(openaiClient)
	creditService := service.NewCreditService(openaiClient)

	aiHandler := handlers.NewAIHandler(analysisService, purchaseService, creditService, logger, emailService)

	if openaiKey == "" {
		logger.Warn().Msg("OPENAI_API_KEY not set — ai module running in mock mode")
	}

	return &Module{
		aiHandler: aiHandler,
		logger:    logger,
	}
}

// RegisterRoutes registers all HTTP routes for the AI module under the provided router group.
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	ai := r.Group("/ai")
	{
		// Financial health analysis
		ai.POST("/health-analysis", m.aiHandler.AnalyzeFinancialHealth)
		ai.POST("/insights", m.aiHandler.GenerateInsights)

		// Purchase decisions
		ai.POST("/can-i-buy", m.aiHandler.CanIBuy)
		ai.POST("/alternatives", m.aiHandler.SuggestAlternatives)

		// Credit analysis
		ai.POST("/credit-plan", m.aiHandler.GenerateCreditPlan)
		ai.POST("/credit-score", m.aiHandler.CalculateCreditScore)

		// Monthly coaching + education
		ai.POST("/monthly-coaching", m.aiHandler.HandleMonthlyCoaching)
		ai.POST("/education-cards", m.aiHandler.HandleEducationCards)
	}

	m.logger.Info().Msg("ai module routes registered")
}

// RegisterSubscribers registers event subscribers for the AI module.
// Currently the AI module does not subscribe to any domain events.
func (m *Module) RegisterSubscribers(eventBus sharedports.EventBus) {
	m.logger.Info().Msg("ai module subscribers registered")
}
