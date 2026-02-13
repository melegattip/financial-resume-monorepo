package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
)

// SettingsHandler handles HTTP requests for user preferences, notification
// settings, data export, and account deletion endpoints (US4 + US6).
type SettingsHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

// NewSettingsHandler creates a new SettingsHandler.
func NewSettingsHandler(authService *services.AuthService, logger zerolog.Logger) *SettingsHandler {
	return &SettingsHandler{
		authService: authService,
		logger:      logger.With().Str("component", "auth").Logger(),
	}
}

// GetPreferences returns the preferences of the authenticated user.
func (h *SettingsHandler) GetPreferences(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	prefs, err := h.authService.GetPreferences(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("get preferences failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get preferences"})
		return
	}

	c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences updates the preferences of the authenticated user.
func (h *SettingsHandler) UpdatePreferences(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var prefs domain.Preferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		h.logger.Warn().Err(err).Msg("invalid update preferences request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.UpdatePreferences(c.Request.Context(), userID, &prefs); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("update preferences failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("preferences updated successfully")
	c.JSON(http.StatusOK, prefs)
}

// GetNotifications returns the notification settings of the authenticated user.
func (h *SettingsHandler) GetNotifications(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	settings, err := h.authService.GetNotifications(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("get notifications failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get notification settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateNotifications updates the notification settings of the authenticated user.
func (h *SettingsHandler) UpdateNotifications(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var settings domain.NotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		h.logger.Warn().Err(err).Msg("invalid update notifications request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.UpdateNotifications(c.Request.Context(), userID, &settings); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("update notifications failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notification settings"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("notification settings updated successfully")
	c.JSON(http.StatusOK, settings)
}

// ExportData exports all user data for the authenticated user.
func (h *SettingsHandler) ExportData(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	data, err := h.authService.ExportData(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("export data failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to export data"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("data exported successfully")
	c.JSON(http.StatusOK, data)
}

// DeleteAccount permanently deletes the authenticated user's account after
// verifying their password.
func (h *SettingsHandler) DeleteAccount(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid delete account request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.DeleteAccount(c.Request.Context(), userID, req.Password); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("delete account failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("account deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
