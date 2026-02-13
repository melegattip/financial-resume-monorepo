package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
)

// ProfileHandler handles HTTP requests for user profile endpoints (US4).
type ProfileHandler struct {
	authService *services.AuthService
	logger      zerolog.Logger
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(authService *services.AuthService, logger zerolog.Logger) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
		logger:      logger.With().Str("component", "auth").Logger(),
	}
}

// GetProfile returns the profile of the authenticated user.
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	resp, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("get profile failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateProfile updates the profile of the authenticated user.
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("invalid update profile request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.authService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("update profile failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	h.logger.Info().Str("user_id", userID).Msg("profile updated successfully")
	c.JSON(http.StatusOK, resp)
}

// UploadAvatar handles multipart file upload for the user's avatar.
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.logger.Warn().Err(err).Msg("failed to get avatar file from form")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Avatar file is required"})
		return
	}

	// Ensure the upload directory exists.
	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		h.logger.Error().Err(err).Str("dir", uploadDir).Msg("failed to create upload directory")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process avatar upload"})
		return
	}

	// Generate a unique filename preserving the original extension.
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		h.logger.Error().Err(err).Str("path", savePath).Msg("failed to save avatar file")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save avatar"})
		return
	}

	// Build the public URL path for the avatar.
	avatarPath := fmt.Sprintf("/api/v1/uploads/avatars/%s", filename)

	if err := h.authService.UploadAvatar(c.Request.Context(), userID, avatarPath); err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("upload avatar failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update avatar"})
		return
	}

	h.logger.Info().Str("user_id", userID).Str("avatar", avatarPath).Msg("avatar uploaded successfully")
	c.JSON(http.StatusOK, gin.H{"avatar": avatarPath})
}
