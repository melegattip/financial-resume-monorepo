package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/middleware"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/handlers"
	authports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/ports"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/repository"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/auth/services"
	budgetsdomain "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/budgets/domain"
	sharedemail "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/email"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// Module implements the ports.Module interface for authentication.
type Module struct {
	db              *gorm.DB
	logger          zerolog.Logger
	cfg             *config.AppConfig
	eventBus        ports.EventBus
	authSvc         *services.AuthService
	authHandler     *handlers.AuthHandler
	securityHandler *handlers.SecurityHandler
	profileHandler  *handlers.ProfileHandler
	settingsHandler *handlers.SettingsHandler
	authMiddleware  *middleware.AuthMiddleware
}

// New creates and initializes the auth module with all dependencies.
func New(db *gorm.DB, logger zerolog.Logger, cfg *config.AppConfig, eventBus ports.EventBus,
	tenantCreator authports.TenantCreator, tenantFinder authports.TenantMemberFinder,
	tenantCleaner authports.TenantAccountCleaner) *Module {
	repo := repository.New(db)

	jwtSvc := services.NewJWTService(cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry, cfg.JWT.Issuer)
	pwSvc := services.NewPasswordService(cfg.Security.PasswordMinLength)
	twoFASvc := services.NewTwoFAService(cfg.JWT.Issuer)
	emailSvc := sharedemail.NewServiceWithResend(
		sharedemail.SMTPConfig{
			Host:     cfg.Email.Host,
			Port:     cfg.Email.Port,
			User:     cfg.Email.User,
			Password: cfg.Email.Password,
			From:     cfg.Email.From,
		},
		sharedemail.ResendConfig{
			APIKey: cfg.Email.ResendAPIKey,
			From:   cfg.Email.From,
		},
		logger,
	)

	authSvc := services.NewAuthService(
		repo, repo, repo, repo,
		jwtSvc, pwSvc, twoFASvc,
		tenantCreator, tenantFinder, tenantCleaner,
		emailSvc, cfg.AppURL,
		eventBus, logger,
		cfg.Security.MaxLoginAttempts,
		cfg.Security.LockoutDuration,
	)

	authMW := middleware.NewAuthMiddleware(jwtSvc)

	return &Module{
		db:              db,
		logger:          logger,
		cfg:             cfg,
		eventBus:        eventBus,
		authSvc:         authSvc,
		authHandler:     handlers.NewAuthHandler(authSvc, logger),
		securityHandler: handlers.NewSecurityHandler(authSvc, logger),
		profileHandler:  handlers.NewProfileHandler(authSvc, logger),
		settingsHandler: handlers.NewSettingsHandler(authSvc, logger),
		authMiddleware:  authMW,
	}
}

// Name returns the module identifier.
func (m *Module) Name() string { return "auth" }

// AuthMiddleware returns the shared auth middleware for other modules to use.
func (m *Module) AuthMiddleware() *middleware.AuthMiddleware {
	return m.authMiddleware
}

// RegisterRoutes adds auth endpoints to the router group.
// The router group should be the /api/v1 group.
func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	if err := m.db.AutoMigrate(
		&domain.User{},
		&domain.Preferences{},
		&domain.NotificationSettings{},
		&domain.TwoFA{},
	); err != nil {
		// Non-fatal: AutoMigrate may fail to reconcile constraint names on existing tables
		// (e.g. "uni_users_email" vs the name used by a previous migration).
		// The schema is correct as long as the tables and columns exist.
		m.logger.Warn().Err(err).Msg("auto-migrate warning (schema may already be up to date)")
	}

	// --- Public auth routes: /api/v1/auth/* ---
	auth := router.Group("/auth")
	{
		auth.POST("/register", m.authHandler.Register)
		auth.POST("/login", m.authHandler.Login)
		auth.POST("/check-2fa", m.authHandler.Check2FA)
		auth.POST("/refresh", m.securityHandler.Refresh)
		auth.GET("/verify-email/:token", m.authHandler.VerifyEmail)
		auth.POST("/request-password-reset", m.securityHandler.RequestPasswordReset)
		auth.POST("/reset-password", m.securityHandler.ResetPassword)
	}

	// --- Protected user routes: /api/v1/users/* ---
	users := router.Group("/users")
	users.Use(m.authMiddleware.RequireAuth())
	{
		// Session
		users.POST("/logout", m.securityHandler.Logout)
		users.POST("/switch-tenant", m.authHandler.SwitchTenant)

		// Profile
		users.GET("/profile", m.profileHandler.GetProfile)
		users.PUT("/profile", m.profileHandler.UpdateProfile)
		users.POST("/profile/avatar", m.profileHandler.UploadAvatar)

		// Password
		users.PUT("/change-password", m.securityHandler.ChangePassword)

		// 2FA
		twoFA := users.Group("/2fa")
		{
			twoFA.POST("/setup", m.securityHandler.Setup2FA)
			twoFA.POST("/enable", m.securityHandler.Enable2FA)
			twoFA.POST("/disable", m.securityHandler.Disable2FA)
			twoFA.POST("/verify", m.securityHandler.Verify2FA)
		}

		// Preferences
		users.GET("/preferences", m.settingsHandler.GetPreferences)
		users.PUT("/preferences", m.settingsHandler.UpdatePreferences)

		// Notifications
		users.GET("/notifications", m.settingsHandler.GetNotifications)
		users.PUT("/notifications", m.settingsHandler.UpdateNotifications)

		// Data management
		users.GET("/export", m.settingsHandler.ExportData)
		users.DELETE("/account", m.settingsHandler.DeleteAccount)
	}

	m.logger.Info().Str("component", "auth").Msg("auth routes registered")
}

// RegisterSubscribers registers event bus subscriptions for the auth module.
func (m *Module) RegisterSubscribers(bus ports.EventBus) {
	bus.Subscribe("budget.threshold_crossed", func(ctx context.Context, event ports.Event) error {
		ev, ok := event.(budgetsdomain.BudgetThresholdCrossedEvent)
		if !ok {
			return nil
		}
		if err := m.authSvc.SendBudgetAlertNotification(ctx, ev.User, ev.CategoryID, ev.Period, ev.NewStatus, ev.SpentAmount, ev.BudgetLimit); err != nil {
			m.logger.Warn().Err(err).Str("user_id", ev.User).Msg("auth subscriber: failed to send budget alert")
		}
		return nil
	})

	m.logger.Info().Msg("auth module subscribers registered")
}
