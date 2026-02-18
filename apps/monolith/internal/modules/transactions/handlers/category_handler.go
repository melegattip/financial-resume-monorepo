package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/transactions/ports"
)

type CategoryHandler struct {
	repo   ports.CategoryRepository
	logger zerolog.Logger
}

func NewCategoryHandler(repo ports.CategoryRepository, logger zerolog.Logger) *CategoryHandler {
	return &CategoryHandler{
		repo:   repo,
		logger: logger,
	}
}

// CategoryResponse is the response format for a category
type CategoryResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func toCategoryResponse(c *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Color:     c.Color,
		Icon:      c.Icon,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

// ListCategories handles GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	categories, err := h.repo.FindByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to list categories")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		response[i] = toCategoryResponse(cat)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  response,
		"total": len(response),
	})
}
