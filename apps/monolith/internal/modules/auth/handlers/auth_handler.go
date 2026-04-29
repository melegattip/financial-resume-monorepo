package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
)

// authSvc is the subset of AuthService methods used by AuthHandler.
// Using an interface decouples the handler from the concrete service and
// makes handler-level unit tests possible without a real database.
type authSvc interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	Check2FA(ctx context.Context, email string) (*domain.Check2FAResponse, error)
	SwitchTenant(ctx context.Context, userID, tenantID string) (*domain.TokenPair, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, email string) error
}

// AuthHandler handles HTTP requests for authentication endpoints.
type AuthHandler struct {
	authService authSvc
	logger      zerolog.Logger
}

// NewAuthHandler creates a new AuthHandler.
// authService is accepted as the interface so tests can inject a mock.
func NewAuthHandler(authService authSvc, logger zerolog.Logger) *AuthHandler {
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

		msg := err.Error()
		switch {
		case msg == "2FA_REQUIRED":
			c.JSON(http.StatusOK, gin.H{
				"twofa_required": true,
				"message":        "2FA code required",
			})
		case msg == "EMAIL_NOT_VERIFIED":
			c.JSON(http.StatusForbidden, gin.H{"error": "EMAIL_NOT_VERIFIED"})
		case strings.Contains(msg, "deactivated"):
			c.JSON(http.StatusForbidden, gin.H{"error": msg})
		case strings.Contains(msg, "locked"):
			c.JSON(http.StatusLocked, gin.H{"error": msg})
		case msg == "invalid email or password":
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
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

// SwitchTenant issues new tokens scoped to a different tenant the user belongs to.
// POST /api/v1/users/switch-tenant
func (h *AuthHandler) SwitchTenant(c *gin.Context) {
	userID := c.GetString("user_id")

	var body struct {
		TenantID string `json:"tenant_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id is required"})
		return
	}

	tokens, err := h.authService.SwitchTenant(c.Request.Context(), userID, body.TenantID)
	if err != nil {
		h.logger.Warn().Err(err).Str("user_id", userID).Str("tenant_id", body.TenantID).Msg("switch tenant failed")
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
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

		if strings.Contains(err.Error(), "token is expired") {
			c.JSON(http.StatusGone, gin.H{"error": "Verification token has expired"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification token"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ResendVerification sends a new email verification link for an unverified account.
// Always returns 200 to prevent email enumeration.
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	if err := h.authService.ResendVerificationEmail(c.Request.Context(), req.Email); err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("resend verification failed")
	}

	c.JSON(http.StatusOK, gin.H{"message": "If that email is registered and unverified, a new verification link has been sent"})
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
