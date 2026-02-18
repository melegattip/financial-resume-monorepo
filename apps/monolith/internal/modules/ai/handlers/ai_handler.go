package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/service"
)

// AIHandler handles HTTP requests for all AI capabilities.
type AIHandler struct {
	analysisService *service.AnalysisService
	purchaseService *service.PurchaseService
	creditService   *service.CreditService
	logger          zerolog.Logger
}

// NewAIHandler creates a new AIHandler.
func NewAIHandler(
	analysis *service.AnalysisService,
	purchase *service.PurchaseService,
	credit *service.CreditService,
	logger zerolog.Logger,
) *AIHandler {
	return &AIHandler{
		analysisService: analysis,
		purchaseService: purchase,
		creditService:   credit,
		logger:          logger,
	}
}

// AnalyzeFinancialHealth godoc
// POST /ai/health-analysis
// Accepts a FinancialAnalysisData body and returns a HealthAnalysis.
func (h *AIHandler) AnalyzeFinancialHealth(c *gin.Context) {
	var req domain.FinancialAnalysisData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	// Inject user_id from JWT context when not provided in the body.
	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().Str("user_id", req.UserID).Msg("analyzing financial health")

	result, err := h.analysisService.AnalyzeFinancialHealth(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to analyze financial health")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze financial health"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GenerateInsights godoc
// POST /ai/insights
// Accepts a FinancialAnalysisData body and returns a list of AIInsights.
func (h *AIHandler) GenerateInsights(c *gin.Context) {
	var req domain.FinancialAnalysisData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().Str("user_id", req.UserID).Msg("generating financial insights")

	insights, err := h.analysisService.GenerateInsights(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to generate insights")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate insights"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    insights,
	})
}

// CanIBuy godoc
// POST /ai/can-i-buy
// Accepts a PurchaseAnalysisRequest body and returns a PurchaseDecision.
func (h *AIHandler) CanIBuy(c *gin.Context) {
	var req domain.PurchaseAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().
		Str("user_id", req.UserID).
		Str("item", req.ItemName).
		Float64("amount", req.Amount).
		Msg("analyzing purchase decision")

	decision, err := h.purchaseService.CanIBuy(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to analyze purchase decision")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to analyze purchase decision"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    decision,
	})
}

// SuggestAlternatives godoc
// POST /ai/alternatives
// Accepts a PurchaseAnalysisRequest body and returns a list of Alternative suggestions.
func (h *AIHandler) SuggestAlternatives(c *gin.Context) {
	var req domain.PurchaseAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().Str("user_id", req.UserID).Str("item", req.ItemName).Msg("generating purchase alternatives")

	alternatives, err := h.purchaseService.SuggestAlternatives(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to generate alternatives")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate purchase alternatives"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    alternatives,
	})
}

// GenerateCreditPlan godoc
// POST /ai/credit-plan
// Accepts a FinancialAnalysisData body and returns a CreditPlan.
func (h *AIHandler) GenerateCreditPlan(c *gin.Context) {
	var req domain.FinancialAnalysisData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().Str("user_id", req.UserID).Msg("generating credit improvement plan")

	plan, err := h.creditService.GenerateCreditPlan(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to generate credit plan")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate credit improvement plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plan,
	})
}

// CalculateCreditScore godoc
// POST /ai/credit-score
// Accepts a FinancialAnalysisData body and returns a CreditScoreResponse.
func (h *AIHandler) CalculateCreditScore(c *gin.Context) {
	var req domain.FinancialAnalysisData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.UserID = uidStr
			}
		}
	}

	h.logger.Info().Str("user_id", req.UserID).Msg("calculating credit score")

	score, err := h.creditService.CalculateCreditScore(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to calculate credit score")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate credit score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": domain.CreditScoreResponse{
			Score:        score,
			UserID:       req.UserID,
			CalculatedAt: time.Now(),
		},
	})
}
