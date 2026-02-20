package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
)

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *services.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger.With().Str("component", "auth").Logger(),
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid register request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("registration failed")

		switch {
		case strings.Contains(err.Error(), "already exists"):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "password validation"):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.Info().Str("email", req.Email).Msg("user registered successfully")
	c.JSON(http.StatusCreated, resp)
}

// Login handles user authentication.
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid login request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("login failed")

		switch {
		case errors.Is(err, services.ErrTwoFARequired):
			c.JSON(http.StatusOK, gin.H{
				"twofa_required": true,
				"message":        "2FA code required",
			})
		case strings.Contains(err.Error(), "deactivated"):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case strings.Contains(err.Error(), "locked"):
			c.JSON(http.StatusLocked, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	h.logger.Info().Str("email", req.Email).Msg("user logged in successfully")
	c.JSON(http.StatusOK, resp)
}

// Check2FA checks whether a user has 2FA enabled.
// Accepts the email either as a JSON body field or as a ?email= query parameter.
func (h *AuthHandler) Check2FA(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		var body struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			email = body.Email
		}
	}
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	resp, err := h.authService.Check2FA(c.Request.Context(), email)
	if err != nil {
		h.logger.Error().Err(err).Str("email", email).Msg("check 2FA failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// VerifyEmail handles email verification via token.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), token); err != nil {
		h.logger.Error().Err(err).Msg("email verification failed")

		switch {
		case strings.Contains(err.Error(), "expired"):
			c.JSON(http.StatusGone, gin.H{"error": "Verification token has expired"})
		case strings.Contains(err.Error(), "invalid"):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// getUserID extracts the user_id (string) from the gin context (set by auth middleware).
// Returns ("", false) if the user_id is not found or is not a string.
func getUserID(c *gin.Context) (string, bool) {
	val, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	if v, ok := val.(string); ok && v != "" {
		return v, true
	}

	return "", false
}
