package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
)

// SecurityHandler handles HTTP requests for password management, 2FA, token
// refresh, and logout endpoints (US2 + US3).
type SecurityHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

// NewSecurityHandler creates a new SecurityHandler.
func NewSecurityHandler(authService *services.AuthService, logger zerolog.Logger) *SecurityHandler {
	return &SecurityHandler{
		authService: authService,
		logger:      logger.With().Str("component", "auth").Logger(),
	}
}

// ChangePassword handles password change for the authenticated user.
func (h *SecurityHandler) ChangePassword(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid change password request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("change password failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("password changed successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// RequestPasswordReset initiates the password reset flow by sending a reset
// link to the given email address. The response is always 200 to avoid
// revealing whether the email exists.
func (h *SecurityHandler) RequestPasswordReset(c *gin.Context) {
	var req domain.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid password reset request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.RequestPasswordReset(c.Request.Context(), req.Email); err != nil {
		// Log the error but do not expose whether the email exists.
		h.logger.Error().Err(err).Str("email", req.Email).Msg("request password reset failed")
	}

	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent"})
}

// ResetPassword completes the password reset flow using a valid reset token.
func (h *SecurityHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid reset password request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword); err != nil {
		h.logger.Error().Err(err).Msg("reset password failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	h.logger.Info().Msg("password reset successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Setup2FA generates a TOTP secret, QR code, and backup codes for the
// authenticated user.
func (h *SecurityHandler) Setup2FA(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	resp, err := h.authService.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("setup 2FA failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup 2FA"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("2FA setup completed")
	c.JSON(http.StatusOK, resp)
}

// Enable2FA activates 2FA for the authenticated user after verifying the TOTP
// code generated during setup.
func (h *SecurityHandler) Enable2FA(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.Enable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid enable 2FA request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.Enable2FA(c.Request.Context(), userID, req.Code); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("enable 2FA failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable 2FA"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("2FA enabled successfully")
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FA deactivates 2FA for the authenticated user after verifying their
// password.
func (h *SecurityHandler) Disable2FA(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.Disable2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid disable 2FA request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.Disable2FA(c.Request.Context(), userID, req.Password); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("disable 2FA failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable 2FA"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("2FA disabled successfully")
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled successfully"})
}

// Verify2FA validates a TOTP code for the authenticated user.
func (h *SecurityHandler) Verify2FA(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid verify 2FA request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authService.Verify2FA(c.Request.Context(), userID, req.Code); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("verify 2FA failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify 2FA code"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("2FA code verified")
	c.JSON(http.StatusOK, gin.H{"message": "2FA code verified"})
}

// Refresh issues a new token pair using a valid refresh token.
func (h *SecurityHandler) Refresh(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid refresh token request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error().Err(err).Msg("refresh token failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	h.logger.Info().Msg("token refreshed successfully")
	c.JSON(http.StatusOK, tokens)
}

// Logout invalidates the current session for the authenticated user.
func (h *SecurityHandler) Logout(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), userID); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("logout failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("user logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
