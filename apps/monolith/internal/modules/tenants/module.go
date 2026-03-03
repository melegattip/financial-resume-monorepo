package tenants

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/handlers"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/services"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module encapsulates the tenants module.
type Module struct {
	db                *gorm.DB
	repo              *repository.GormRepository
	tenantHandler     *handlers.TenantHandler
	memberHandler     *handlers.MemberHandler
	invitationHandler *handlers.InvitationHandler
	auditHandler      *handlers.AuditHandler
	logger            zerolog.Logger
	authMW            *middleware.AuthMiddleware
	permMW            *middleware.PermissionMiddleware
}

// New creates a new tenants Module, wiring all dependencies.
func New(
	db *gorm.DB,
	logger zerolog.Logger,
	_ *config.AppConfig,
	_ sharedports.EventBus,
	authMW *middleware.AuthMiddleware,
	permMW *middleware.PermissionMiddleware,
) *Module {
	repo := repository.NewGormRepository(db)
	svc := services.NewTenantService(repo, logger)

	return &Module{
		db:                db,
		repo:              repo,
		tenantHandler:     handlers.NewTenantHandler(svc, logger),
		memberHandler:     handlers.NewMemberHandler(svc, logger),
		invitationHandler: handlers.NewInvitationHandler(svc, logger),
		auditHandler:      handlers.NewAuditHandler(svc, logger),
		logger:            logger,
		authMW:            authMW,
		permMW:            permMW,
	}
}

// RegisterRoutes registers all HTTP routes for the tenants module.
//
// Routes:
//
//	GET    /tenants/me                         → get current tenant
//	PUT    /tenants/me                         → update tenant (manage_tenant)
//	DELETE /tenants/me                         → delete tenant (delete_tenant)
//	GET    /tenants/me/permissions             → list caller's permissions
//	GET    /tenants/me/members                 → list members
//	PUT    /tenants/me/members/:userID/role    → change role (manage_roles)
//	DELETE /tenants/me/members/:userID         → remove member (remove_members)
//	GET    /tenants/me/invitations             → list invitations (invite_members)
//	POST   /tenants/me/invitations             → create invitation (invite_members)
//	DELETE /tenants/me/invitations/:code       → revoke invitation (invite_members)
//	POST   /tenants/join                       → join via code (auth only)
//	GET    /tenants/me/audit                   → audit logs (view_audit_logs)
func (m *Module) RegisterRoutes(r *gin.RouterGroup) {
	tenants := r.Group("/tenants")
	tenants.Use(m.authMW.RequireAuth())

	// List all tenants the user belongs to (for tenant switcher)
	tenants.GET("/list", m.tenantHandler.ListMyTenants)

	// Join via invitation code — no extra permission required
	tenants.POST("/join", m.invitationHandler.Join)

	me := tenants.Group("/me")
	{
		me.GET("", m.tenantHandler.GetMyTenant)
		me.PUT("", m.permMW.Require("manage_tenant"), m.tenantHandler.UpdateMyTenant)
		me.DELETE("", m.permMW.Require("delete_tenant"), m.tenantHandler.DeleteMyTenant)
		me.GET("/permissions", m.tenantHandler.GetMyPermissions)

		// Members
		members := me.Group("/members")
		{
			members.GET("", m.memberHandler.ListMembers)
			members.PUT("/:userID/role", m.permMW.Require("manage_roles"), m.memberHandler.UpdateMemberRole)
			members.DELETE("/:userID", m.permMW.Require("remove_members"), m.memberHandler.RemoveMember)
		}

		// Invitations
		invitations := me.Group("/invitations")
		invitations.Use(m.permMW.Require("invite_members"))
		{
			invitations.GET("", m.invitationHandler.List)
			invitations.POST("", m.invitationHandler.Create)
			invitations.DELETE("/:code", m.invitationHandler.Revoke)
		}

		// Audit logs
		me.GET("/audit", m.permMW.Require("view_audit_logs"), m.auditHandler.List)
	}

	m.logger.Info().Msg("tenants module routes registered")
}

// RegisterSubscribers wires the tenants module to all auditable domain events.
// The AuditLogger writes an immutable audit_log entry for every event.
func (m *Module) RegisterSubscribers(bus sharedports.EventBus) {
	auditLogger := services.NewAuditLogger(m.db, m.repo, m.logger)

	auditable := []string{
		"expense.created", "expense.updated", "expense.deleted",
		"income.created", "income.updated", "income.deleted",
		"recurring.created", "recurring.updated", "recurring.deleted",
		"recurring.executed", "recurring.paused", "recurring.resumed",
		"budget.created", "budget.updated", "budget.deleted",
		"savings_goal.created", "savings_goal.updated", "savings_goal.deleted",
		"user.registered",
	}

	for _, eventType := range auditable {
		et := eventType // capture loop variable (pre-Go-1.22 safety)
		bus.Subscribe(et, func(ctx context.Context, event sharedports.Event) error {
			return auditLogger.Handle(ctx, event)
		})
	}

	m.logger.Info().Msg("tenants module subscribers registered")
}
