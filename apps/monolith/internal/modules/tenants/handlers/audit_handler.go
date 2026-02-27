package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
)

// AuditHandler handles requests for tenant audit logs.
type AuditHandler struct {
	service ports.TenantService
	logger  zerolog.Logger
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(service ports.TenantService, logger zerolog.Logger) *AuditHandler {
	return &AuditHandler{service: service, logger: logger}
}

// List returns paginated audit logs for the current tenant.
// GET /api/v1/tenants/me/audit?limit=50&offset=0
func (h *AuditHandler) List(c *gin.Context) {
	tenantID := c.GetString("tenant_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.service.ListAuditLogs(c.Request.Context(), tenantID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("tenant_id", tenantID).Msg("list audit logs failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load audit logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"limit":  limit,
		"offset": offset,
	})
}
