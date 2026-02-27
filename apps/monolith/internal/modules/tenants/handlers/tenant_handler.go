package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
)

// TenantHandler handles requests related to tenant management.
type TenantHandler struct {
	service ports.TenantService
	logger  zerolog.Logger
}

// NewTenantHandler creates a new TenantHandler.
func NewTenantHandler(service ports.TenantService, logger zerolog.Logger) *TenantHandler {
	return &TenantHandler{service: service, logger: logger}
}

// GetMyTenant returns the current user's tenant.
// GET /api/v1/tenants/me
func (h *TenantHandler) GetMyTenant(c *gin.Context) {
	userID := c.GetString("user_id")

	tenant, err := h.service.GetMyTenant(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Msg("get my tenant failed")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tenant": tenant})
}

// UpdateMyTenant updates the current tenant's name or settings.
// PUT /api/v1/tenants/me
func (h *TenantHandler) UpdateMyTenant(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	var req domain.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMyTenant(c.Request.Context(), tenantID, userID, req); err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("update tenant failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant updated"})
}

// DeleteMyTenant soft-deletes the current tenant.
// DELETE /api/v1/tenants/me
func (h *TenantHandler) DeleteMyTenant(c *gin.Context) {
	userID := c.GetString("user_id")
	tenantID := c.GetString("tenant_id")

	if err := h.service.DeleteMyTenant(c.Request.Context(), tenantID, userID); err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("delete tenant failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted"})
}

// GetMyPermissions returns the permission keys for the caller's current role.
// GET /api/v1/tenants/me/permissions
func (h *TenantHandler) GetMyPermissions(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	role := c.GetString("role")

	perms, err := h.service.GetMyPermissions(c.Request.Context(), tenantID, role)
	if err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("get permissions failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role":        role,
		"permissions": perms,
	})
}
