package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/melegattip/financial-resume-engine/docs"
	"github.com/melegattip/financial-resume-engine/internal/adapters/http"
	"github.com/melegattip/financial-resume-engine/internal/config"
	"github.com/melegattip/financial-resume-engine/internal/config/environment"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
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
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/calculators"
	recurringTransactionsHandler "github.com/melegattip/financial-resume-engine/internal/infrastructure/http/handlers/recurring_transactions"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/proxy"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/repository"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/router"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/scheduler"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
	"github.com/melegattip/financial-resume-engine/internal/usecases/analytics"
	categoriesUsecases "github.com/melegattip/financial-resume-engine/internal/usecases/categories"
	dashboardUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/dashboard"
	insightsUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/insights"
	"github.com/melegattip/financial-resume-engine/internal/usecases/recurring_transactions"
	"github.com/melegattip/financial-resume-engine/internal/usecases/reports"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
	expensesCreateUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/create"
	expensesDeleteUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/delete"
	expensesGetUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/get"
	expensesListUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/list"
	expensesUpdateUsecase "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses/update"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Financial Resume Engine API
// @version 1.0
// @description API para gestión de finanzas personales con análisis inteligente

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter JWT token with 'Bearer ' prefix (e.g., 'Bearer your_token')
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	} else {
		log.Println("✅ .env file loaded successfully")
	}

	// Setup environment configuration with automatic detection
	environment.SetUp()

	// JWT secret para validar tokens del microservicio de usuarios
	jwtSecret := getEnv("JWT_SECRET", "financial_resume_secret_key_2024")

	db := config.InitDB()

	// Nota: userRepo eliminado - la gestión de usuarios se hace en users-service
	incomeRepo := repository.NewIncomeRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	recurringTransactionRepo := repository.NewRecurringTransactionRepository(db)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("❌ Failed to get SQL DB from GORM:", err)
	}
	savingsGoalRepo := repository.NewSavingsGoalRepository(sqlDB)

	periodCalculator := calculators.NewPeriodCalculator()
	analyticsCalculator := calculators.NewAnalyticsCalculator()

	categoryService := services.NewCategoryService(categoryRepo)

	// Usar configuración centralizada de environment
	envConfig := environment.GetConfig()

	gamificationServiceURL := envConfig.GamificationURL
	gamificationHelper := services.NewGamificationHelper(gamificationServiceURL)

	// Configurar URL del servicio de usuarios usando configuración centralizada
	usersServiceURL := envConfig.UsersServiceURL
	log.Printf("👤 [Main] Users Service URL configurada: %s", usersServiceURL)

	// Nota: authService eliminado - se gestiona en users-service

	incomeService := incomes.NewIncomeService(incomeRepo)
	expenseService := transactions.NewExpenseService(expenseRepo, incomeRepo)

	percentageObserver := transactions.NewPercentageObserver(expenseService)

	dashboardService := dashboardUsecase.NewService(
		expenseRepo,
		incomeRepo,
		periodCalculator,
		analyticsCalculator,
	)

	aiServiceURL := envConfig.AIServiceURL
	aiServiceProxy := proxy.NewAIServiceProxy(aiServiceURL)

	legacyAIService := services.NewAIService()

	insightsService := insightsUsecase.NewService(
		expenseRepo,
		incomeRepo,
		categoryRepo,
		savingsGoalRepo,
		recurringTransactionRepo,
		periodCalculator,
		analyticsCalculator,
		aiServiceProxy,
		aiServiceProxy,
		aiServiceProxy,
		legacyAIService,
	)

	expensesAnalyticsService := analytics.NewExpensesAnalyticsService(
		expenseRepo,
		incomeRepo,
		categoryService,
		periodCalculator,
		analyticsCalculator,
	)

	categoriesAnalyticsService := analytics.NewCategoriesAnalyticsService(
		expenseRepo,
		incomeRepo,
		categoryService,
		periodCalculator,
		analyticsCalculator,
	)

	incomesAnalyticsService := analytics.NewIncomesAnalyticsService(
		incomeRepo,
		categoryService,
		periodCalculator,
		analyticsCalculator,
	)

	budgetService := usecases.NewBudgetService(budgetRepo, categoryRepo, expenseRepo, nil) // nil for notification service for now
	savingsGoalService := usecases.NewSavingsGoalUseCase(savingsGoalRepo, nil)             // nil for notification service for now

	recurringExecutorService := services.NewRecurringTransactionExecutorService(expenseRepo, incomeRepo, gamificationHelper)
	recurringNotificationService := services.NewRecurringTransactionNotificationService()

	recurringTransactionService := recurring_transactions.NewService(
		recurringTransactionRepo,
		expenseRepo,
		incomeRepo,
		categoryRepo,
		recurringNotificationService,
		recurringExecutorService,
	)

	// Nota: authHandlerInstance eliminado - se gestiona en users-service

	incomeHandler := handlers.NewIncomeHandler(incomeService, gamificationHelper)

	createExpenseHandler := expensesCreateHandler.NewHandler(expensesCreateUsecase.NewService(expenseRepo, categoryRepo, percentageObserver), gamificationHelper)
	getExpenseHandler := expensesGetHandler.NewHandler(expensesGetUsecase.NewService(expenseRepo), gamificationHelper)
	listExpenseHandler := expensesListHandler.NewHandler(expensesListUsecase.NewService(expenseRepo), gamificationHelper)
	updateExpenseHandler := expensesUpdateHandler.NewHandler(expensesUpdateUsecase.NewService(expenseRepo, percentageObserver), gamificationHelper)
	deleteExpenseHandler := expensesDeleteHandler.NewHandler(expensesDeleteUsecase.NewService(expenseRepo, percentageObserver), gamificationHelper)

	categoryUseCase := categoriesUsecases.NewService(categoryRepo)
	categoryHandler := categories.NewHandler(categoryUseCase, gamificationHelper)

	reportHandler := reports.NewReportHandler(db)

	dashboardHandler := dashboard.NewHandler(dashboardService, gamificationHelper)

	insightsHandler := insights.NewHandler(insightsService, gamificationHelper)

	analyticsHandlers := http.NewAnalyticsHandlers(
		expensesAnalyticsService,
		categoriesAnalyticsService,
		incomesAnalyticsService,
	)

	budgetHandlerInstance := budgetHandler.NewHandler(budgetService)
	savingsGoalsHandlerInstance := savingsGoalsHandler.NewHandler(savingsGoalService)

	recurringTransactionsHandlerInstance := recurringTransactionsHandler.NewHandler(recurringTransactionService)

	configHandlerInstance := configHandler.NewHandler(db)

	router := router.SetupRouter(
		incomeHandler,
		createExpenseHandler,
		getExpenseHandler,
		listExpenseHandler,
		updateExpenseHandler,
		deleteExpenseHandler,
		categoryHandler,
		reportHandler,
		dashboardHandler,
		analyticsHandlers,
		insightsHandler,
		budgetHandlerInstance,
		savingsGoalsHandlerInstance,
		jwtSecret, // Usar JWT secret en lugar del servicio
		recurringTransactionsHandlerInstance,
		configHandlerInstance,
		envConfig, // Pasar configuración de entorno
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	schedulerIntervalHours := getEnv("RECURRING_SCHEDULER_INTERVAL_HOURS", "1")
	var schedulerInterval time.Duration
	if hours := schedulerIntervalHours; hours != "" {
		if h, err := time.ParseDuration(hours + "h"); err == nil {
			schedulerInterval = h
		} else {
			schedulerInterval = 1 * time.Hour
		}
	} else {
		schedulerInterval = 1 * time.Hour
	}

	recurringScheduler := scheduler.NewRecurringTransactionScheduler(recurringTransactionService, schedulerInterval)
	recurringScheduler.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	port := getEnv("PORT", "8080")
	log.Printf("🚀 Financial Resume Engine starting on port %s", port)
	log.Printf("📊 Swagger documentation available at: http://localhost:%s/swagger/index.html", port)
	log.Printf("🔐 JWT Authentication delegated to users-service")
	log.Printf("⏰ Recurring transactions scheduler started (Interval: %v)", schedulerInterval)

	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatal("❌ Failed to start server:", err)
		}
	}()

	<-c
	log.Println("🛑 Shutting down gracefully...")

	recurringScheduler.Stop()

	log.Println("✅ Shutdown complete")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
