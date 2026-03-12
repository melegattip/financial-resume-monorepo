package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/gamification/service"
)

// GamificationHandler exposes the gamification service over HTTP.
type GamificationHandler struct {
	service *service.GamificationService
	logger  zerolog.Logger
}

// NewGamificationHandler creates a new GamificationHandler.
func NewGamificationHandler(svc *service.GamificationService, logger zerolog.Logger) *GamificationHandler {
	return &GamificationHandler{service: svc, logger: logger}
}

// ---------------------------------------------------------------------------
// Request / response types
// ---------------------------------------------------------------------------

// RecordActionRequest is the body accepted by the POST /actions endpoint.
type RecordActionRequest struct {
	ActionType  string `json:"action_type" binding:"required"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Description string `json:"description"`
}

// RecordActionResponse is the body returned after a successful action recording.
type RecordActionResponse struct {
	XPEarned     int  `json:"xp_earned"`
	TotalXP      int  `json:"total_xp"`
	CurrentLevel int  `json:"current_level"`
	LevelUp      bool `json:"level_up"`
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// GetProfile handles GET /gamification/profile
// Returns the full UserGamification aggregate for the authenticated user.
func (h *GamificationHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.service.GetUserGamification(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to get gamification profile")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get gamification profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetStats handles GET /gamification/stats
// Returns an aggregated statistics summary for the authenticated user.
func (h *GamificationHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	stats, err := h.service.GetGamificationStats(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to get gamification stats")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get gamification stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetAchievements handles GET /gamification/achievements
// Returns all achievements (completed and in-progress) for the authenticated user.
func (h *GamificationHandler) GetAchievements(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	achievements, err := h.service.GetAchievements(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to get achievements")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// GetFeatures handles GET /gamification/features
// Returns the list of application features available to the authenticated user.
func (h *GamificationHandler) GetFeatures(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unlocked_features": []string{
			"SAVINGS_GOALS",
			"BUDGETS",
			"AI_INSIGHTS",
			"dashboard",
			"expenses",
			"incomes",
			"categories",
			"analytics",
			"savings_goals",
			"budgets",
			"recurring",
			"ai_insights",
		},
		"locked_features": []interface{}{},
	})
}

// CheckFeatureAccess handles GET /gamification/features/:featureKey/access
// Returns whether a specific feature key is in trial mode.
// Since all features are currently unlocked, this always returns trial_active: false.
func (h *GamificationHandler) CheckFeatureAccess(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// All features are currently unlocked — no trial state to report.
	c.JSON(http.StatusOK, gin.H{
		"trial_active":  false,
		"trial_ends_at": nil,
	})
}

// GetDailyChallenges handles GET /gamification/challenges/daily
// Returns the list of daily challenges for the authenticated user.
// Challenges are not yet persisted — returns an empty list.
func (h *GamificationHandler) GetDailyChallenges(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, []interface{}{})
}

// GetWeeklyChallenges handles GET /gamification/challenges/weekly
// Returns the list of weekly challenges for the authenticated user.
// Challenges are not yet persisted — returns an empty list.
func (h *GamificationHandler) GetWeeklyChallenges(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, []interface{}{})
}

// ProcessChallengeProgress handles POST /gamification/challenges/progress
// Records progress on a challenge. Not yet fully implemented — returns success.
func (h *GamificationHandler) ProcessChallengeProgress(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// RecordAction handles POST /gamification/actions
// Records a user action and returns XP/level information.
func (h *GamificationHandler) RecordAction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req RecordActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.RecordAction(c.Request.Context(), userID.(string), req.ActionType, req.EntityID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID.(string)).Str("action_type", req.ActionType).Msg("failed to record action")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record action"})
		return
	}

	c.JSON(http.StatusOK, RecordActionResponse{
		XPEarned:     result.XPEarned,
		TotalXP:      result.TotalXP,
		CurrentLevel: result.CurrentLevel,
		LevelUp:      result.LevelUp,
	})
}


// GetBehaviorProfile handles GET /gamification/behavior-profile
// Returns a behavioral profile derived from the user's action history.
func (h *GamificationHandler) GetBehaviorProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	profile, err := h.service.GetBehaviorProfile(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID.(string)).Msg("failed to get behavior profile")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get behavior profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}
