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
	Priority  int    `json:"priority"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateCategoryRequest is the request body for creating a category.
type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
	Priority int    `json:"priority"`
}

// UpdateCategoryRequest is the request body for updating a category.
// Accepts "new_name" (frontend convention) or "name" as the category name.
type UpdateCategoryRequest struct {
	NewName  string `json:"new_name"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
	Priority *int   `json:"priority"`
}

func toCategoryResponse(c *domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:        c.ID,
		UserID:    c.UserID,
		Name:      c.Name,
		Color:     c.Color,
		Icon:      c.Icon,
		Priority:  c.Priority,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
	}
}

// List handles GET /api/v1/categories
func (h *CategoryHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	categories, err := h.repo.FindByTenantID(c.Request.Context(), tenantID)
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

// Create handles POST /api/v1/categories
func (h *CategoryHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID := c.GetString("tenant_id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := domain.NewCategory(userID.(string), req.Name, req.Color, req.Icon, req.Priority)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	category.TenantID = tenantID

	if err := h.repo.Create(c.Request.Context(), category); err != nil {
		h.logger.Error().Err(err).Msg("failed to create category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, toCategoryResponse(category))
}

// Update handles PATCH /api/v1/categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	category, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if category.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Accept "new_name" (frontend convention) or "name"
	name := req.NewName
	if name == "" {
		name = req.Name
	}
	if name == "" {
		name = category.Name // keep existing name
	}

	color := req.Color
	if color == "" {
		color = category.Color
	}
	icon := req.Icon
	if icon == "" {
		icon = category.Icon
	}
	priority := category.Priority
	if req.Priority != nil {
		priority = *req.Priority
	}

	if err := category.Update(name, color, icon, priority); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), category); err != nil {
		h.logger.Error().Err(err).Msg("failed to update category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, toCategoryResponse(category))
}

// Delete handles DELETE /api/v1/categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	id := c.Param("id")
	category, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("failed to get category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get category"})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	if category.TenantID != tenantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().Err(err).Msg("failed to delete category")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
