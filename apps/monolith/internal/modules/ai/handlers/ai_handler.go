package handlers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/ai/service"
	sharedemail "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/email"
)

// cachedInsights holds the result of a GenerateInsights call for one calendar day.
type cachedInsights struct {
	insights    []domain.AIInsight
	generatedAt time.Time
}

// AIHandler handles HTTP requests for all AI capabilities.
type AIHandler struct {
	analysisService *service.AnalysisService
	purchaseService *service.PurchaseService
	creditService   *service.CreditService
	logger          zerolog.Logger
	emailService    sharedemail.EmailService

	// Daily cache for GenerateInsights: key is "<userID>_<YYYY-MM-DD UTC>".
	mu            sync.Mutex
	insightsCache map[string]cachedInsights

	// Monthly coaching cache — key: "{userID}_{YYYY-MM}"
	monthlyMu    sync.Mutex
	monthlyCache map[string]domain.MonthlyCoachingReport
	emailSentMu  sync.Mutex
	emailSentMap map[string]bool // key: "{userID}_{YYYY-MM}" → true when email sent

	// Education cache — key: "{userID}_{YYYY-WW}"
	educationMu    sync.Mutex
	educationCache map[string]domain.EducationContent
}

// NewAIHandler creates a new AIHandler.
func NewAIHandler(
	analysis *service.AnalysisService,
	purchase *service.PurchaseService,
	credit *service.CreditService,
	logger zerolog.Logger,
	emailService sharedemail.EmailService,
) *AIHandler {
	return &AIHandler{
		analysisService: analysis,
		purchaseService: purchase,
		creditService:   credit,
		logger:          logger,
		emailService:    emailService,
		insightsCache:   make(map[string]cachedInsights),
		monthlyCache:    make(map[string]domain.MonthlyCoachingReport),
		emailSentMap:    make(map[string]bool),
		educationCache:  make(map[string]domain.EducationContent),
	}
}

// educationCacheKey returns a cache key based on user ID and the ISO week of t.
func educationCacheKey(userID string, t time.Time) string {
	year, week := t.ISOWeek()
	return fmt.Sprintf("%s_%d-W%02d", userID, year, week)
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

	// Return cached response if it was generated today (calendar day, UTC).
	today := time.Now().UTC().Format("2006-01-02")
	cacheKey := req.UserID + "_" + today

	h.mu.Lock()
	cached, hit := h.insightsCache[cacheKey]
	h.mu.Unlock()

	if hit {
		h.logger.Info().Str("user_id", req.UserID).Str("date", today).Msg("returning cached insights")
		c.JSON(http.StatusOK, gin.H{
			"success":      true,
			"data":         cached.insights,
			"generated_at": cached.generatedAt,
		})
		return
	}

	h.logger.Info().Str("user_id", req.UserID).Msg("generating financial insights")

	insights, err := h.analysisService.GenerateInsights(c.Request.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", req.UserID).Msg("failed to generate insights")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate insights"})
		return
	}

	// Store in daily cache.
	now := time.Now().UTC()
	h.mu.Lock()
	h.insightsCache[cacheKey] = cachedInsights{insights: insights, generatedAt: now}
	h.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"data":         insights,
		"generated_at": now,
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

// HandleMonthlyCoaching godoc
// POST /ai/monthly-coaching
// Accepts a MonthlyCoachingRequest body and returns a MonthlyCoachingReport.
// Results are cached per (userID, month) — the same report is never regenerated twice.
// The first generation of each month triggers an email to the user (async, best-effort).
func (h *AIHandler) HandleMonthlyCoaching(c *gin.Context) {
	var req domain.MonthlyCoachingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.FinancialData.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.FinancialData.UserID = uidStr
			}
		}
	}

	// Validate YYYY-MM format.
	if _, err := time.Parse("2006-01", req.PreviousMonth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "previous_month must be in YYYY-MM format"})
		return
	}

	userID := req.FinancialData.UserID
	cacheKey := userID + "_" + req.PreviousMonth

	h.monthlyMu.Lock()
	cached, hit := h.monthlyCache[cacheKey]
	h.monthlyMu.Unlock()

	if hit {
		h.logger.Info().Str("user_id", userID).Str("month", req.PreviousMonth).Msg("returning cached monthly coaching")
		c.JSON(http.StatusOK, gin.H{"report": cached, "cached": true})
		return
	}

	h.logger.Info().Str("user_id", userID).Str("month", req.PreviousMonth).Msg("generating monthly coaching report")

	report, err := h.analysisService.GenerateMonthlyCoaching(c.Request.Context(), req.FinancialData, req.PreviousMonth)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("failed to generate monthly coaching")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate monthly coaching report"})
		return
	}

	h.monthlyMu.Lock()
	h.monthlyCache[cacheKey] = *report
	h.monthlyMu.Unlock()

	// Send email once per (user, month), best-effort in a goroutine.
	h.emailSentMu.Lock()
	alreadySent := h.emailSentMap[cacheKey]
	if !alreadySent {
		h.emailSentMap[cacheKey] = true
	}
	h.emailSentMu.Unlock()

	if !alreadySent && h.emailService != nil {
		go func() {
			data := sharedemail.CoachingReportEmailData{
				Month:     report.Month,
				Sentiment: report.Sentiment,
				Summary:   report.Summary,
			}
			for _, w := range report.Wins {
				data.Wins = append(data.Wins, sharedemail.CoachingEmailPoint{Title: w.Title, Description: w.Description})
			}
			for _, imp := range report.Improvements {
				data.Improvements = append(data.Improvements, sharedemail.CoachingEmailPoint{Title: imp.Title, Description: imp.Description})
			}
			for _, act := range report.Actions {
				data.Actions = append(data.Actions, sharedemail.CoachingEmailAction{Title: act.Title, Detail: act.Detail})
			}
			data.BehaviorNote = report.BehaviorNote

			toEmail := req.FinancialData.UserEmail
			if toEmail == "" {
				h.logger.Warn().Str("user_id", userID).Msg("skipping coaching email: no email address in request")
				return
			}
			if err := h.emailService.SendMonthlyCoachingReport(toEmail, req.FinancialData.UserName, data); err != nil {
				h.logger.Error().Err(err).Str("user_id", userID).Str("month", report.Month).Msg("failed to send monthly coaching email")
			}
		}()
	}

	c.JSON(http.StatusOK, gin.H{"report": report, "cached": false})
}

// HandleEducationCards godoc
// POST /ai/education-cards
// Accepts an EducationRequest body and returns 3 personalized education cards.
// Results are cached for the current ISO week.
func (h *AIHandler) HandleEducationCards(c *gin.Context) {
	var req domain.EducationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format", "details": err.Error()})
		return
	}

	if req.FinancialData.UserID == "" {
		if uid, exists := c.Get("user_id"); exists {
			if uidStr, ok := uid.(string); ok {
				req.FinancialData.UserID = uidStr
			}
		}
	}

	userID := req.FinancialData.UserID
	cacheKey := educationCacheKey(userID, time.Now().UTC())

	h.educationMu.Lock()
	cached, hit := h.educationCache[cacheKey]
	h.educationMu.Unlock()

	if hit {
		h.logger.Info().Str("user_id", userID).Msg("returning cached education cards")
		c.JSON(http.StatusOK, gin.H{"cards": cached.Cards, "generated_at": cached.GeneratedAt, "cached": true})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("generating education cards")

	cards, err := h.analysisService.GenerateEducationCards(c.Request.Context(), req.FinancialData)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("failed to generate education cards")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate education cards"})
		return
	}

	now := time.Now().UTC()
	content := domain.EducationContent{Cards: cards, GeneratedAt: now}

	h.educationMu.Lock()
	h.educationCache[cacheKey] = content
	h.educationMu.Unlock()

	c.JSON(http.StatusOK, gin.H{"cards": cards, "generated_at": now, "cached": false})
}
