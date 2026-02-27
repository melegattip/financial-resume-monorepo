package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
)

// MemberHandler handles requests related to tenant membership management.
type MemberHandler struct {
	service ports.TenantService
	logger  zerolog.Logger
}

// NewMemberHandler creates a new MemberHandler.
func NewMemberHandler(service ports.TenantService, logger zerolog.Logger) *MemberHandler {
	return &MemberHandler{service: service, logger: logger}
}

// ListMembers returns all members of the current tenant.
// GET /api/v1/tenants/me/members
func (h *MemberHandler) ListMembers(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	members, err := h.service.ListMembers(c.Request.Context(), tenantID)
	if err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("list members failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list members"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": members})
}

// UpdateMemberRole changes the role of a tenant member.
// PUT /api/v1/tenants/me/members/:userID/role
func (h *MemberHandler) UpdateMemberRole(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	requesterRole := c.GetString("role")
	targetUserID := c.Param("userID")

	var req domain.UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateMemberRole(c.Request.Context(), tenantID, requesterRole, targetUserID, req.Role); err != nil {
		h.logger.Error().Err(err).
			Str("tenant_id", tenantID).
			Str("target_user_id", targetUserID).
			Msg("update member role failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

// RemoveMember removes a user from the tenant.
// DELETE /api/v1/tenants/me/members/:userID
func (h *MemberHandler) RemoveMember(c *gin.Context) {
	tenantID := c.GetString("tenant_id")
	requesterRole := c.GetString("role")
	requesterID := c.GetString("user_id")
	targetUserID := c.Param("userID")

	if err := h.service.RemoveMember(c.Request.Context(), tenantID, requesterRole, requesterID, targetUserID); err != nil {
		h.logger.Error().Err(err).
			Str("tenant_id", tenantID).
			Str("target_user_id", targetUserID).
			Msg("remove member failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}
