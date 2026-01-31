package router

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	httpAdapters "github.com/melegattip/financial-resume-engine/internal/adapters/http"
	"github.com/melegattip/financial-resume-engine/internal/config/environment"
	"github.com/melegattip/financial-resume-engine/internal/handlers"
	budgetHandler "github.com/melegattip/financial-resume-engine/internal/handlers/budget"
	"github.com/melegattip/financial-resume-engine/internal/handlers/categories"
	configHandler "github.com/melegattip/financial-resume-engine/internal/handlers/config"
	"github.com/melegattip/financial-resume-engine/internal/handlers/dashboard"
	expensesCreateHandler "github.com/melegattip/financial-resume-engine/internal/handlers/expenses/create"
	expensesDeleteHandler "github.com/melegattip/financial-resume-engine/internal/handlers/expenses/delete"
	expensesGetHandler "github.com/melegattip/financial-resume-engine/internal/handlers/expenses/get"
	expensesListHandler "github.com/melegattip/financial-resume-engine/internal/handlers/expenses/list"
	expensesUpdateHandler "github.com/melegattip/financial-resume-engine/internal/handlers/expenses/update"
	"github.com/melegattip/financial-resume-engine/internal/handlers/insights"
	savingsGoalsHandler "github.com/melegattip/financial-resume-engine/internal/handlers/savings_goals"
	recurringTransactionsHandler "github.com/melegattip/financial-resume-engine/internal/infrastructure/http/handlers/recurring_transactions"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/http/middleware"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/proxy"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/reports"
)

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(
	incomeHandler *handlers.IncomeHandler,
	createExpenseHandler *expensesCreateHandler.Handler,
	getExpenseHandler *expensesGetHandler.Handler,
	listExpenseHandler *expensesListHandler.Handler,
	updateExpenseHandler *expensesUpdateHandler.Handler,
	deleteExpenseHandler *expensesDeleteHandler.Handler,
	categoryHandler *categories.Handler,
	reportHandler *reports.ReportHandler,
	dashboardHandler *dashboard.Handler,
	analyticsHandlers *httpAdapters.AnalyticsHandlers,
	insightsHandler *insights.Handler,
	budgetHandlerInstance *budgetHandler.Handler,
	savingsGoalsHandlerInstance *savingsGoalsHandler.Handler,
	jwtSecret string, // JWT secret del microservicio de usuarios
	recurringTransactionsHandlerInstance *recurringTransactionsHandler.Handler,
	configHandlerInstance *configHandler.Handler,
	envConfig environment.ServiceConfig, // Agregar configuración de entorno
) *gin.Engine {
	router := gin.Default()

	// === NUEVOS SERVICIOS DE SEGURIDAD ===

	// Inicializar servicio de rate limiting
	rateLimitService := services.NewInMemoryRateLimitService()

	// Configuración de rate limiting
	var rateLimitConfig middleware.RateLimitConfig

	if gin.Mode() == gin.DebugMode {
		// DESARROLLO: Rate limiting desactivado (límites muy altos)
		rateLimitConfig = middleware.RateLimitConfig{
			RequestsPerMinute:   10000, // Prácticamente ilimitado en desarrollo
			IPLimitEnabled:      false, // Desactivar límite por IP en desarrollo
			IPRequestsPerMinute: 10000,
			EndpointLimits:      map[string]int{}, // Sin límites específicos en desarrollo
			SkipPaths: []string{
				"/health",
				"/favicon.ico",
				"/robots.txt",
				"/manifest.json",
				"/swagger/",
				"/docs/",
				"/config",
			},
		}
	} else {
		// PRODUCCIÓN: Rate limiting activo
		rateLimitConfig = middleware.RateLimitConfig{
			RequestsPerMinute:   100, // 100 requests por minuto por usuario
			IPLimitEnabled:      true,
			IPRequestsPerMinute: 200, // 200 requests por minuto por IP
			EndpointLimits:      map[string]int{
				// Nota: Los endpoints de autenticación se gestionan en users-service
			},
			SkipPaths: []string{
				"/health",
				"/favicon.ico",
				"/robots.txt",
				"/manifest.json",
				"/swagger/",
				"/docs/",
				"/config",
			},
		}
	}

	// === MIDDLEWARES GLOBALES ===

	// CORS configuration
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With", "x-caller-id", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Configurar orígenes permitidos basado en el entorno
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	defaultOrigins := []string{
		// Desarrollo
		"http://localhost:3000",
		"http://localhost:8080",
		"http://localhost:8081",
		"http://localhost:8082",
		// Render.com
		"https://financial-resume-engine-frontend.onrender.com",
		"https://financial-resume-engine.onrender.com",
		"https://financial-ai-api.niloft.com",
		"https://financial-gamification-service.onrender.com",
	}

	var allowedOrigins []string
	if allowedOriginsEnv != "" {
		allowedOrigins = strings.Split(allowedOriginsEnv, ",")
		// Limpiar espacios en blanco
		for i, origin := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(origin)
		}
	} else {
		allowedOrigins = defaultOrigins
	}

	// En desarrollo, permitir todos los orígenes
	if gin.Mode() == gin.DebugMode {
		corsConfig.AllowAllOrigins = true
		corsConfig.AllowCredentials = false // No se puede usar credenciales con AllowAllOrigins
	} else {
		corsConfig.AllowOrigins = allowedOrigins
	}

	router.Use(cors.New(corsConfig))

	// NUEVO: Middlewares de seguridad y métricas
	router.Use(middleware.MetricsMiddleware(rateLimitService))
	router.Use(middleware.RateLimitMiddleware(rateLimitService, rateLimitConfig))

	// Security headers middleware
	router.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; connect-src 'self'")
		c.Next()
	})

	// === RUTAS DE SALUD Y CONFIGURACIÓN ===

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "financial-resume-engine",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "1.0.0",
		})
	})

	// === INICIALIZAR PROXIES ===

	// Inicializar proxy del servicio de usuarios
	usersServiceProxy := proxy.NewUsersServiceProxy(envConfig.UsersServiceURL)

	// Inicializar proxy de gamificación
	gamificationProxy := proxy.NewGamificationProxy(envConfig.GamificationURL)

	// === RUTAS PÚBLICAS (SIN AUTENTICACIÓN) ===

	// Config endpoint (público, no requiere autenticación)
	router.GET("/api/v1/config", configHandlerInstance.GetConfig)

	// Diagnostics endpoint (público, para debugging en producción)
	router.GET("/api/v1/diagnostics", configHandlerInstance.GetDiagnostics)

	// Grupo para autenticación (no requiere token)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", usersServiceProxy.ProxyPublicRequest)
		authGroup.POST("/login", usersServiceProxy.ProxyPublicRequest)
		authGroup.POST("/logout", usersServiceProxy.ProxyPublicRequest)
		authGroup.POST("/check-2fa", usersServiceProxy.ProxyPublicRequest)
	}

	// Endpoints públicos de gamificación (no requieren autenticación)
	publicGamificationGroup := router.Group("/api/v1/gamification")
	{
		publicGamificationGroup.GET("/action-types", gamificationProxy.ProxyPublicRequest)
		publicGamificationGroup.GET("/levels", gamificationProxy.ProxyPublicRequest)
	}

	// === RUTAS PROTEGIDAS (CON AUTENTICACIÓN JWT) ===
	// Nota: La autenticación se gestiona en el microservicio users-service

	// Middleware de autenticación JWT para todas las rutas protegidas
	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuthMiddleware(jwtSecret))

	// Rutas de autenticación que requieren token - REDIRIGIDAS AL SERVICIO DE USUARIOS
	authProtectedGroup := api.Group("/auth")
	{
		authProtectedGroup.GET("/profile", usersServiceProxy.ProxyRequest)
		authProtectedGroup.POST("/refresh", usersServiceProxy.ProxyRequest)
		authProtectedGroup.PUT("/change-password", usersServiceProxy.ProxyRequest)
	}

	// Rutas protegidas de usuarios (requieren autenticación)
	usersProtectedGroup := api.Group("/users")
	{
		// Perfil
		usersProtectedGroup.GET("/profile", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.PUT("/profile", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.POST("/logout", usersServiceProxy.ProxyRequest)

		// Preferencias
		usersProtectedGroup.GET("/preferences", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.PUT("/preferences", usersServiceProxy.ProxyRequest)

		// Notificaciones
		usersProtectedGroup.GET("/notifications/settings", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.PUT("/notifications/settings", usersServiceProxy.ProxyRequest)

		// Seguridad
		usersProtectedGroup.PUT("/security/change-password", usersServiceProxy.ProxyRequest)

		// 2FA endpoints
		usersProtectedGroup.POST("/security/2fa/setup", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.POST("/security/2fa/enable", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.POST("/security/2fa/disable", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.POST("/security/2fa/verify", usersServiceProxy.ProxyRequest)

		// Gestión de datos
		usersProtectedGroup.POST("/export", usersServiceProxy.ProxyRequest)
		usersProtectedGroup.DELETE("", usersServiceProxy.ProxyRequest)
	}

	// Rutas de categorías
	categoryRoutes := api.Group("/categories")
	{
		categoryRoutes.GET("", categoryHandler.ListCategories)
		categoryRoutes.POST("", categoryHandler.CreateCategory)
		categoryRoutes.GET("/:id", categoryHandler.GetCategory)
		categoryRoutes.PATCH("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// Rutas de gastos
	expenseRoutes := api.Group("/expenses")
	{
		expenseRoutes.GET("", listExpenseHandler.ListExpenses)
		expenseRoutes.GET("/unpaid", listExpenseHandler.ListUnpaidExpenses)
		expenseRoutes.POST("", createExpenseHandler.CreateExpense)
		expenseRoutes.GET("/:userId/:id", getExpenseHandler.GetExpense)
		expenseRoutes.PATCH("/:userId/:id", updateExpenseHandler.UpdateExpense)
		expenseRoutes.DELETE("/:userId/:id", deleteExpenseHandler.DeleteExpense)
	}

	// Rutas de ingresos
	incomeRoutes := api.Group("/incomes")
	{
		incomeRoutes.GET("", incomeHandler.ListIncomes)
		incomeRoutes.POST("", incomeHandler.CreateIncome)
		incomeRoutes.GET("/:userId/:id", incomeHandler.GetIncome)
		incomeRoutes.PATCH("/:userId/:id", incomeHandler.UpdateIncome)
		incomeRoutes.DELETE("/:userId/:id", incomeHandler.DeleteIncome)
	}

	// Rutas de dashboard
	dashboardRoutes := api.Group("/dashboard")
	{
		dashboardRoutes.GET("", dashboardHandler.GetDashboard)
	}

	// Rutas de analytics
	analyticsRoutes := api.Group("/analytics")
	{
		analyticsRoutes.GET("/expenses", analyticsHandlers.GetExpensesSummary)
		analyticsRoutes.GET("/categories", analyticsHandlers.GetCategoriesAnalytics)
		analyticsRoutes.GET("/incomes", analyticsHandlers.GetIncomesSummary)
	}

	// Rutas de insights
	insightsRoutes := api.Group("/insights")
	{
		insightsRoutes.GET("/financial-health", insightsHandler.GetFinancialHealth)
		insightsRoutes.POST("/mark-understood", insightsHandler.MarkInsightAsUnderstood)
	}

	// Rutas de IA
	aiRoutes := api.Group("/ai")
	{
		aiRoutes.GET("/insights", insightsHandler.GetAIInsights)
		aiRoutes.POST("/can-i-buy", insightsHandler.CanIBuy)
		aiRoutes.GET("/credit-improvement-plan", insightsHandler.GetCreditImprovementPlan)
	}

	// Rutas de reportes
	reportRoutes := api.Group("/reports")
	{
		reportRoutes.GET("", reportHandler.HandleGenerateReport)
	}

	// Rutas de presupuestos
	budgetHandlerInstance.RegisterRoutes(api)

	// Rutas de metas de ahorro
	savingsGoalsHandlerInstance.RegisterRoutes(api)

	// Rutas de transacciones recurrentes
	if recurringTransactionsHandlerInstance != nil {
		recurringTransactionsHandlerInstance.RegisterRoutes(api)
	}

	// Rutas de gamificación protegidas
	gamificationRoutes := api.Group("/gamification")
	{
		gamificationRoutes.GET("/profile", gamificationProxy.ProxyRequest)
		gamificationRoutes.GET("/stats", gamificationProxy.ProxyRequest)
		gamificationRoutes.GET("/achievements", gamificationProxy.ProxyRequest)
		gamificationRoutes.GET("/features", gamificationProxy.ProxyRequest)
		gamificationRoutes.GET("/features/:featureKey/access", gamificationProxy.ProxyRequest)
		gamificationRoutes.POST("/actions", gamificationProxy.ProxyRequest)

		// Challenge endpoints
		gamificationRoutes.GET("/challenges/daily", gamificationProxy.ProxyRequest)
		gamificationRoutes.GET("/challenges/weekly", gamificationProxy.ProxyRequest)
		gamificationRoutes.POST("/challenges/progress", gamificationProxy.ProxyRequest)
	}

	return router
}
