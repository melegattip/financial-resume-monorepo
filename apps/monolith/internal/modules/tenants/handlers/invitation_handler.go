package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
)

// InvitationHandler handles requests related to tenant invitations.
type InvitationHandler struct {
	service ports.TenantService
	logger  zerolog.Logger
}

// NewInvitationHandler creates a new InvitationHandler.
func NewInvitationHandler(service ports.TenantService, logger zerolog.Logger) *InvitationHandler {
	return &InvitationHandler{service: service, logger: logger}
}

// Create generates a new invitation code for the current tenant.
// POST /api/v1/tenants/me/invitations
func (h *InvitationHandler) Create(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	userID := c.GetString("user_id")

	var req domain.CreateInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	inv, err := h.service.CreateInvitation(c.Request.Context(), tenantID, userID, req)
	if err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("create invitation failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invitation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invitation": inv})
}

// List returns all active invitations for the current tenant.
// GET /api/v1/tenants/me/invitations
func (h *InvitationHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	invitations, err := h.service.ListInvitations(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("list invitations failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invitations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invitations": invitations})
}

// Revoke deactivates an invitation by its code.
// DELETE /api/v1/tenants/me/invitations/:code
func (h *InvitationHandler) Revoke(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	code := c.Param("code")

	if err := h.service.RevokeInvitation(c.Request.Context(), tenantID, code); err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Str("code", code).Msg("revoke invitation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "invitation revoked"})
}

// Join adds the calling user to a tenant using an invitation code.
// POST /api/v1/tenants/join
func (h *InvitationHandler) Join(c *gin.Context) {
	userID := c.GetString("user_id")

	var req domain.JoinTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.service.JoinTenant(c.Request.Context(), userID, req.Code)
	if err != nil {
		h.logger.Error().Err(err).Str("user_id", userID).Str("code", req.Code).Msg("join tenant failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "joined tenant successfully; please log in again to activate the new context",
		"tenant":  tenant,
	})
}
