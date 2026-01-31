package insights

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
	"github.com/melegattip/financial-resume-engine/internal/infrastructure/services"
)

// CacheEntry representa una entrada en el cache
type CacheEntry struct {
	Data      *AIInsightsResponse
	Timestamp time.Time
}

// Service implementa InsightsUseCase siguiendo Clean Architecture
type Service struct {
	expenseRepo              baseRepo.ExpenseRepository
	incomeRepo               baseRepo.IncomeRepository
	categoryRepo             baseRepo.CategoryRepository
	savingsGoalRepo          ports.SavingsGoalRepository
	recurringTransactionRepo ports.RecurringTransactionRepository
	periodCalculator         usecases.PeriodCalculator
	analyticsCalculator      usecases.AnalyticsCalculator

	// ✅ NUEVO: Servicios especializados siguiendo interfaces
	aiAnalysisService       ports.AIAnalysisPort
	purchaseDecisionService ports.PurchaseDecisionPort
	creditAnalysisService   ports.CreditAnalysisPort

	// ❌ DEPRECADO: Servicio monolítico (mantenemos por compatibilidad temporal)
	aiService *services.AIService

	// Cache simple en memoria con TTL de 20 horas
	insightsCache map[string]CacheEntry
	cacheMutex    sync.RWMutex
	cacheTTL      time.Duration
}

// NewService crea una nueva instancia del servicio con dependency injection
func NewService(
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	categoryRepo baseRepo.CategoryRepository,
	savingsGoalRepo ports.SavingsGoalRepository,
	recurringTransactionRepo ports.RecurringTransactionRepository,
	periodCalculator usecases.PeriodCalculator,
	analyticsCalculator usecases.AnalyticsCalculator,
	// ✅ NUEVO: Inyección de servicios especializados
	aiAnalysisService ports.AIAnalysisPort,
	purchaseDecisionService ports.PurchaseDecisionPort,
	creditAnalysisService ports.CreditAnalysisPort,
	// ❌ DEPRECADO: Mantenemos por compatibilidad temporal
	aiService *services.AIService,
) *Service {
	return &Service{
		expenseRepo:              expenseRepo,
		incomeRepo:               incomeRepo,
		categoryRepo:             categoryRepo,
		savingsGoalRepo:          savingsGoalRepo,
		recurringTransactionRepo: recurringTransactionRepo,
		periodCalculator:         periodCalculator,
		analyticsCalculator:      analyticsCalculator,

		// ✅ NUEVO: Asignar servicios especializados
		aiAnalysisService:       aiAnalysisService,
		purchaseDecisionService: purchaseDecisionService,
		creditAnalysisService:   creditAnalysisService,

		// ❌ DEPRECADO: Mantenemos por compatibilidad temporal
		aiService: aiService,

		// Cache con TTL de 20 horas
		insightsCache: make(map[string]CacheEntry),
		cacheTTL:      20 * time.Hour,
	}
}

// GetFinancialHealth implementa el caso de uso principal
func (s *Service) GetFinancialHealth(ctx context.Context, params InsightsParams) (*FinancialHealthResponse, error) {
	// Validar parámetros
	if err := s.validateParams(params); err != nil {
		return nil, err
	}

	// Obtener y analizar datos financieros
	analyzedData, err := s.analyzeFinancialData(ctx, params)
	if err != nil {
		return nil, err
	}

	// Calcular score de salud financiera
	healthScore := s.calculateHealthScore(analyzedData)

	// Generar insights personalizados
	insights := s.generateInsights(analyzedData, healthScore)

	// Determinar nivel y mensaje
	level, message := s.getHealthLevelAndMessage(healthScore)

	return &FinancialHealthResponse{
		HealthScore:  healthScore,
		Level:        level,
		Message:      message,
		Insights:     insights,
		AnalyzedData: *analyzedData,
		GeneratedAt:  time.Now(),
	}, nil
}

// validateParams valida los parámetros de entrada
func (s *Service) validateParams(params InsightsParams) error {
	if params.UserID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}

	if params.Period.Year != nil {
		if *params.Period.Year < 1900 || *params.Period.Year > 2100 {
			return errors.NewBadRequest("Año inválido. Debe estar entre 1900 y 2100")
		}
	}

	if params.Period.Month != nil {
		if *params.Period.Month < 1 || *params.Period.Month > 12 {
			return errors.NewBadRequest("Mes inválido. Debe estar entre 1 y 12")
		}
	}

	return nil
}

// analyzeFinancialData analiza los datos financieros del usuario
func (s *Service) analyzeFinancialData(ctx context.Context, params InsightsParams) (*AnalyzedFinancialData, error) {
	// Obtener transacciones
	transactions, err := s.getTransactions(ctx, params.UserID, params.Period)
	if err != nil {
		return nil, err
	}

	// Calcular métricas básicas
	metrics := s.periodCalculator.CalculateMetrics(transactions)

	// Analizar categorías
	categoryAnalysis, err := s.analyzeCategoriesSpending(ctx, transactions, params.UserID)
	if err != nil {
		return nil, err
	}

	// Analizar patrones de gasto
	spendingPatterns := s.analyzeSpendingPatterns(transactions)

	// Analizar estabilidad de ingresos
	incomeStability := s.analyzeIncomeStability(transactions)

	// Analizar cumplimiento de presupuesto (por ahora mock)
	budgetCompliance := s.analyzeBudgetCompliance(transactions)

	// Calcular tasa de ahorro
	savingsRate := s.calculateSavingsRate(metrics.TotalIncome, metrics.TotalExpenses)

	// Construir información del período
	periodInfo := s.buildPeriodInfo(params.Period)

	return &AnalyzedFinancialData{
		Period:            periodInfo,
		TotalIncome:       metrics.TotalIncome,
		TotalExpenses:     metrics.TotalExpenses,
		Balance:           metrics.Balance,
		SavingsRate:       savingsRate,
		ExpenseCategories: categoryAnalysis,
		SpendingPatterns:  spendingPatterns,
		IncomeStability:   incomeStability,
		BudgetCompliance:  budgetCompliance,
	}, nil
}

// getTransactions obtiene todas las transacciones del usuario
func (s *Service) getTransactions(ctx context.Context, userID string, period DatePeriod) ([]usecases.Transaction, error) {
	// Obtener gastos
	expenses, err := s.expenseRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo gastos: %w", err)
	}

	// Obtener ingresos
	incomes, err := s.incomeRepo.List(userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo ingresos: %w", err)
	}

	// Convertir a transacciones genéricas
	var transactions []usecases.Transaction
	for _, expense := range expenses {
		transactions = append(transactions, NewExpenseTransactionFromDomain(expense))
	}
	for _, income := range incomes {
		transactions = append(transactions, NewIncomeTransactionFromDomain(income))
	}

	// Filtrar por período
	usecasesPeriod := usecases.DatePeriod{
		Year:  period.Year,
		Month: period.Month,
	}
	return s.periodCalculator.FilterTransactionsByPeriod(transactions, usecasesPeriod), nil
}

// calculateHealthScore calcula el score de salud financiera (0-1000)
func (s *Service) calculateHealthScore(data *AnalyzedFinancialData) int {
	var score float64

	// 1. Tasa de ahorro (40% del score)
	savingsScore := s.calculateSavingsScore(data.SavingsRate)
	score += savingsScore * 0.4

	// 2. Estabilidad de ingresos (25% del score)
	incomeScore := s.calculateIncomeStabilityScore(data.IncomeStability)
	score += incomeScore * 0.25

	// 3. Diversificación de gastos (20% del score)
	diversificationScore := s.calculateDiversificationScore(data.ExpenseCategories)
	score += diversificationScore * 0.2

	// 4. Control de gastos inusuales (15% del score)
	controlScore := s.calculateSpendingControlScore(data.SpendingPatterns)
	score += controlScore * 0.15

	// Asegurar que esté en el rango 0-1000
	finalScore := int(math.Max(0, math.Min(1000, score)))
	return finalScore
}

// getHealthLevelAndMessage determina el nivel y mensaje basado en el score
func (s *Service) getHealthLevelAndMessage(score int) (string, string) {
	if score >= 800 {
		return "Excelente", "¡Tu salud financiera es excelente! Sigue así."
	} else if score >= 600 {
		return "Bueno", "Tu salud financiera es buena, con oportunidades de mejora."
	} else if score >= 400 {
		return "Regular", "Tu salud financiera es regular. Hay varias áreas de mejora."
	} else {
		return "Mejorable", "Tu salud financiera necesita atención. ¡Pero puedes mejorar!"
	}
}

// === NUEVOS MÉTODOS DE IA ===

// GetAIInsights genera insights usando IA basados en datos financieros
func (s *Service) GetAIInsights(ctx context.Context, params InsightsParams) (*AIInsightsResponse, error) {
	// Crear clave de cache basada en userID y período (manejo correcto de nils)
	yearStr := "current"
	monthStr := "current"
	if params.Period.Year != nil {
		yearStr = fmt.Sprintf("%d", *params.Period.Year)
	}
	if params.Period.Month != nil {
		monthStr = fmt.Sprintf("%d", *params.Period.Month)
	}
	cacheKey := fmt.Sprintf("%s_%s_%s", params.UserID, yearStr, monthStr)

	// Verificar cache primero
	s.cacheMutex.RLock()
	if entry, exists := s.insightsCache[cacheKey]; exists {
		// Verificar si el cache aún es válido (TTL)
		timeSinceCache := time.Since(entry.Timestamp)
		if timeSinceCache < s.cacheTTL {
			s.cacheMutex.RUnlock()
			// Log para debugging
			fmt.Printf("✅ Cache HIT para clave '%s' - Edad: %.1f horas\n", cacheKey, timeSinceCache.Hours())
			return entry.Data, nil
		} else {
			fmt.Printf("⏰ Cache EXPIRED para clave '%s' - Edad: %.1f horas (TTL: %.1f horas)\n", cacheKey, timeSinceCache.Hours(), s.cacheTTL.Hours())
		}
	} else {
		fmt.Printf("❌ Cache MISS para clave '%s'\n", cacheKey)
	}
	s.cacheMutex.RUnlock()

	// Si no hay cache válido, generar nuevos insights
	// Obtener datos financieros analizados
	analyzedData, err := s.analyzeFinancialData(ctx, params)
	if err != nil {
		return nil, err
	}

	// ✅ NUEVO: Usar servicio especializado para análisis
	aiAnalysisData := s.convertToAnalysisData(params.UserID, analyzedData)

	// Generar insights con IA usando el servicio especializado
	insights, err := s.aiAnalysisService.GenerateInsights(ctx, aiAnalysisData)
	if err != nil {
		// Fallback al servicio monolítico si falla el especializado
		fmt.Printf("⚠️ Fallback: Error con servicio especializado, usando servicio monolítico: %v\n", err)
		aiRequest := s.buildAIRequest(params.UserID, analyzedData)
		oldInsights, fallbackErr := s.aiService.GenerateFinancialInsights(ctx, aiRequest)
		if fallbackErr != nil {
			return nil, fmt.Errorf("error generando insights con IA: %w", fallbackErr)
		}

		// Convertir formato del servicio monolítico
		aiInsights := make([]AIInsight, len(oldInsights))
		for i, insight := range oldInsights {
			aiInsights[i] = AIInsight{
				Title:       insight.Title,
				Description: insight.Description,
				Impact:      insight.Impact,
				Score:       insight.Score,
				ActionType:  insight.ActionType,
				Category:    insight.Category,
			}
		}

		source := "ai_fallback"
		response := &AIInsightsResponse{
			Insights:    aiInsights,
			GeneratedAt: time.Now(),
			Source:      source,
		}

		// Guardar en cache
		s.cacheMutex.Lock()
		s.insightsCache[cacheKey] = CacheEntry{
			Data:      response,
			Timestamp: time.Now(),
		}
		s.cacheMutex.Unlock()

		return response, nil
	}

	// Convertir a formato de respuesta
	aiInsights := make([]AIInsight, len(insights))
	for i, insight := range insights {
		aiInsights[i] = AIInsight{
			Title:       insight.Title,
			Description: insight.Description,
			Impact:      insight.Impact,
			Score:       insight.Score,
			ActionType:  insight.ActionType,
			Category:    insight.Category,
		}
	}

	source := "ai"
	if s.aiService == nil {
		source = "mock"
	}

	response := &AIInsightsResponse{
		Insights:    aiInsights,
		GeneratedAt: time.Now(),
		Source:      source,
	}

	// Guardar en cache
	s.cacheMutex.Lock()
	s.insightsCache[cacheKey] = CacheEntry{
		Data:      response,
		Timestamp: time.Now(),
	}
	s.cacheMutex.Unlock()

	// Log para debugging
	fmt.Printf("💾 Cache SAVED para clave '%s' - TTL: %.1f horas\n", cacheKey, s.cacheTTL.Hours())

	return response, nil
}

// CanIBuy analiza si el usuario puede permitirse una compra
func (s *Service) CanIBuy(ctx context.Context, params CanIBuyParams) (*CanIBuyResponse, error) {
	// Validar parámetros
	if params.UserID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}
	if params.ItemName == "" {
		return nil, errors.NewBadRequest("El nombre del artículo es requerido")
	}
	if params.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto debe ser mayor a 0")
	}

	// Obtener metas de ahorro del usuario
	savingsGoals, err := s.getSavingsGoalsForAnalysis(ctx, params.UserID, params.ItemName, params.Description)
	if err != nil {
		// Log el error pero no fallar la operación
		fmt.Printf("Warning: No se pudieron obtener metas de ahorro para usuario %s: %v\n", params.UserID, err)
		savingsGoals = []services.SavingsGoalInfo{}
	}

	// Obtener datos financieros completos incluyendo transacciones recurrentes
	totalMonthlyIncome, totalMonthlyExpenses, err := s.calculateTotalMonthlyFinancials(ctx, params.UserID, params.MonthlyIncome, params.MonthlyExpenses)
	if err != nil {
		// Log el error pero continuar con datos básicos
		fmt.Printf("Warning: Error calculando transacciones recurrentes para usuario %s: %v\n", params.UserID, err)
		totalMonthlyIncome = params.MonthlyIncome
		totalMonthlyExpenses = params.MonthlyExpenses
	}

	// ✅ NUEVO: Usar servicio especializado para decisiones de compra
	purchaseRequest := ports.PurchaseAnalysisRequest{
		UserID:       params.UserID,
		ItemName:     params.ItemName,
		Amount:       params.Amount,
		Description:  params.Description,
		PaymentTypes: params.PaymentTypes,
		IsNecessary:  params.IsNecessary,
		UserFinancialProfile: ports.UserFinancialProfile{
			CurrentBalance:  params.CurrentBalance,
			MonthlyIncome:   totalMonthlyIncome,
			MonthlyExpenses: totalMonthlyExpenses,
			SavingsGoals:    s.convertSavingsGoals(savingsGoals),
		},
	}

	// Analizar con IA usando el servicio especializado
	decision, err := s.purchaseDecisionService.CanIBuy(ctx, purchaseRequest)
	if err != nil {
		// Fallback al servicio monolítico si falla el especializado
		fmt.Printf("⚠️ Fallback: Error con servicio especializado de compra, usando servicio monolítico: %v\n", err)

		// Construir perfil financiero completo del usuario con datos corregidos
		userProfile := s.aiService.BuildUserFinancialProfile(
			params.UserID,
			params.CurrentBalance,
			totalMonthlyIncome,   // Usar datos corregidos con recurrentes
			totalMonthlyExpenses, // Usar datos corregidos con recurrentes
			savingsGoals,
		)

		// Preparar request para IA con perfil completo
		aiRequest := services.CanIBuyRequest{
			UserID:               params.UserID,
			ItemName:             params.ItemName,
			Amount:               params.Amount,
			Description:          params.Description,
			PaymentTypes:         params.PaymentTypes,
			IsNecessary:          params.IsNecessary,
			UserFinancialProfile: userProfile,
			SavingsGoal:          params.SavingsGoal, // Mantener por compatibilidad
		}

		// Analizar con IA
		response, fallbackErr := s.aiService.CanIBuyThis(ctx, aiRequest)
		if fallbackErr != nil {
			return nil, fmt.Errorf("error analizando compra con IA: %w", fallbackErr)
		}

		source := "ai_fallback"
		return &CanIBuyResponse{
			CanBuy:       response.CanBuy,
			Confidence:   response.Confidence,
			Reasoning:    response.Reasoning,
			Alternatives: response.Alternatives,
			ImpactScore:  response.ImpactScore,
			GeneratedAt:  time.Now(),
			Source:       source,
		}, nil
	}

	// ✅ NUEVO: Usar respuesta del servicio especializado
	source := "ai_specialized"
	return &CanIBuyResponse{
		CanBuy:       decision.CanBuy,
		Confidence:   decision.Confidence,
		Reasoning:    decision.Reasoning,
		Alternatives: decision.Alternatives,
		ImpactScore:  decision.ImpactScore,
		GeneratedAt:  time.Now(),
		Source:       source,
	}, nil
}

// GetCreditImprovementPlan genera un plan personalizado de mejora crediticia
func (s *Service) GetCreditImprovementPlan(ctx context.Context, params InsightsParams) (*CreditImprovementPlanResponse, error) {
	// Obtener datos financieros analizados
	analyzedData, err := s.analyzeFinancialData(ctx, params)
	if err != nil {
		return nil, err
	}

	// ✅ NUEVO: Usar servicio especializado para análisis crediticio
	aiAnalysisData := s.convertToAnalysisData(params.UserID, analyzedData)

	// Generar plan con IA usando el servicio especializado
	plan, err := s.creditAnalysisService.GenerateImprovementPlan(ctx, aiAnalysisData)
	if err != nil {
		// Fallback al servicio monolítico si falla el especializado
		fmt.Printf("⚠️ Fallback: Error con servicio especializado de crédito, usando servicio monolítico: %v\n", err)

		// Preparar request para IA
		aiRequest := s.buildAIRequest(params.UserID, analyzedData)

		// Generar plan con IA
		fallbackPlan, fallbackErr := s.aiService.GenerateCreditImprovementPlan(ctx, aiRequest)
		if fallbackErr != nil {
			return nil, fmt.Errorf("error generando plan de mejora con IA: %w", fallbackErr)
		}

		// Convertir acciones a formato de respuesta
		actions := make([]CreditAction, len(fallbackPlan.Actions))
		for i, action := range fallbackPlan.Actions {
			actions[i] = CreditAction{
				Title:       action.Title,
				Description: action.Description,
				Priority:    action.Priority,
				Timeline:    action.Timeline,
				Impact:      action.Impact,
				Difficulty:  action.Difficulty,
			}
		}

		source := "ai_fallback"
		return &CreditImprovementPlanResponse{
			CurrentScore:   fallbackPlan.CurrentScore,
			TargetScore:    fallbackPlan.TargetScore,
			TimelineMonths: fallbackPlan.TimelineMonths,
			Actions:        actions,
			KeyMetrics:     fallbackPlan.KeyMetrics,
			GeneratedAt:    time.Now(),
			Source:         source,
		}, nil
	}

	// Convertir acciones a formato de respuesta
	actions := make([]CreditAction, len(plan.Actions))
	for i, action := range plan.Actions {
		actions[i] = CreditAction{
			Title:       action.Title,
			Description: action.Description,
			Priority:    action.Priority,
			Timeline:    action.Timeline,
			Impact:      action.Impact,
			Difficulty:  action.Difficulty,
		}
	}

	source := "ai_specialized"
	return &CreditImprovementPlanResponse{
		CurrentScore:   plan.CurrentScore,
		TargetScore:    plan.TargetScore,
		TimelineMonths: plan.TimelineMonths,
		Actions:        actions,
		KeyMetrics:     plan.KeyMetrics,
		GeneratedAt:    time.Now(),
		Source:         source,
	}, nil
}

// ✅ NUEVO: Convertir datos al formato de las nuevas interfaces
func (s *Service) convertToAnalysisData(userID string, data *AnalyzedFinancialData) ports.FinancialAnalysisData {
	// Construir mapa de gastos por categoría
	expensesByCategory := make(map[string]float64)
	for _, category := range data.ExpenseCategories {
		expensesByCategory[category.CategoryName] = category.Amount
	}

	// Calcular score de salud financiera
	healthScore := s.calculateHealthScore(data)

	return ports.FinancialAnalysisData{
		UserID:             userID,
		TotalIncome:        data.TotalIncome,
		TotalExpenses:      data.TotalExpenses,
		SavingsRate:        data.SavingsRate,
		ExpensesByCategory: expensesByCategory,
		IncomeStability:    data.IncomeStability.IncomeVariation,
		FinancialScore:     healthScore,
		Period:             data.Period.Label,
	}
}

// ✅ NUEVO: Convertir metas de ahorro al formato de las nuevas interfaces
func (s *Service) convertSavingsGoals(goals []services.SavingsGoalInfo) []ports.SavingsGoalInfo {
	if len(goals) == 0 {
		return []ports.SavingsGoalInfo{}
	}

	result := make([]ports.SavingsGoalInfo, len(goals))
	for i, goal := range goals {
		result[i] = ports.SavingsGoalInfo{
			Name:          goal.Name,
			Category:      goal.Category,
			CurrentAmount: goal.CurrentAmount,
			TargetAmount:  goal.TargetAmount,
			Progress:      goal.Progress,
		}
	}
	return result
}

// buildAIRequest construye el request para el servicio de IA (DEPRECADO - para compatibilidad)
func (s *Service) buildAIRequest(userID string, data *AnalyzedFinancialData) services.FinancialAnalysisRequest {
	// Construir mapa de gastos por categoría
	expensesByCategory := make(map[string]float64)
	for _, category := range data.ExpenseCategories {
		expensesByCategory[category.CategoryName] = category.Amount
	}

	// Calcular score de salud financiera
	healthScore := s.calculateHealthScore(data)

	return services.FinancialAnalysisRequest{
		UserID:             userID,
		TotalIncome:        data.TotalIncome,
		TotalExpenses:      data.TotalExpenses,
		SavingsRate:        data.SavingsRate,
		ExpensesByCategory: expensesByCategory,
		IncomeStability:    data.IncomeStability.IncomeVariation,
		FinancialScore:     healthScore,
		Period:             data.Period.Label,
	}
}

// getSavingsGoalsForAnalysis obtiene metas de ahorro relevantes para el análisis de compra
func (s *Service) getSavingsGoalsForAnalysis(ctx context.Context, userID, itemName, description string) ([]services.SavingsGoalInfo, error) {
	// Obtener todas las metas de ahorro activas del usuario
	goals, err := s.savingsGoalRepo.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo metas de ahorro: %w", err)
	}

	var relevantGoals []services.SavingsGoalInfo
	itemNameLower := strings.ToLower(itemName)
	descriptionLower := strings.ToLower(description)

	for _, goal := range goals {
		// Solo incluir metas activas
		if goal.Status != "active" {
			continue
		}

		// Verificar si la meta es relevante para la compra
		isRelevant := false
		goalNameLower := strings.ToLower(goal.Name)
		goalCategoryLower := strings.ToLower(string(goal.Category))

		// 1. Coincidencia directa en nombre o categoría
		if strings.Contains(itemNameLower, goalCategoryLower) ||
			strings.Contains(goalNameLower, itemNameLower) ||
			strings.Contains(goalCategoryLower, itemNameLower) {
			isRelevant = true
		}

		// 2. Coincidencias específicas por categoría
		switch string(goal.Category) {
		case "car":
			if strings.Contains(itemNameLower, "auto") ||
				strings.Contains(itemNameLower, "carro") ||
				strings.Contains(itemNameLower, "vehículo") ||
				strings.Contains(itemNameLower, "vehiculo") ||
				strings.Contains(descriptionLower, "auto") ||
				strings.Contains(descriptionLower, "carro") {
				isRelevant = true
			}
		case "house":
			if strings.Contains(itemNameLower, "casa") ||
				strings.Contains(itemNameLower, "vivienda") ||
				strings.Contains(itemNameLower, "inmueble") ||
				strings.Contains(itemNameLower, "propiedad") {
				isRelevant = true
			}
		case "vacation":
			if strings.Contains(itemNameLower, "viaje") ||
				strings.Contains(itemNameLower, "vacaciones") ||
				strings.Contains(itemNameLower, "turismo") {
				isRelevant = true
			}
		case "education":
			if strings.Contains(itemNameLower, "curso") ||
				strings.Contains(itemNameLower, "educación") ||
				strings.Contains(itemNameLower, "educacion") ||
				strings.Contains(itemNameLower, "estudio") {
				isRelevant = true
			}
		}

		// 3. Análisis de descripción para palabras clave
		if !isRelevant && description != "" {
			if strings.Contains(descriptionLower, goalNameLower) ||
				strings.Contains(descriptionLower, goalCategoryLower) {
				isRelevant = true
			}
		}

		if isRelevant {
			relevantGoals = append(relevantGoals, services.SavingsGoalInfo{
				Name:          goal.Name,
				Category:      string(goal.Category),
				CurrentAmount: goal.CurrentAmount,
				TargetAmount:  goal.TargetAmount,
				Progress:      goal.Progress,
			})
		}
	}

	return relevantGoals, nil
}

// calculateTotalMonthlyFinancials calcula los ingresos y gastos mensuales totales
// incluyendo transacciones recurrentes
func (s *Service) calculateTotalMonthlyFinancials(ctx context.Context, userID string, baseMonthlyIncome, baseMonthlyExpenses float64) (float64, float64, error) {
	// Obtener proyección de transacciones recurrentes
	projection, err := s.recurringTransactionRepo.GetRecurringProjection(ctx, userID, 1) // 1 mes
	if err != nil {
		return baseMonthlyIncome, baseMonthlyExpenses, fmt.Errorf("error obteniendo proyección recurrente: %w", err)
	}

	// Sumar ingresos y gastos base + recurrentes
	totalMonthlyIncome := baseMonthlyIncome + projection.MonthlyIncome
	totalMonthlyExpenses := baseMonthlyExpenses + projection.MonthlyExpenses

	fmt.Printf("💰 Cálculo financiero completo para usuario %s:\n", userID)
	fmt.Printf("   - Ingresos base: $%.0f\n", baseMonthlyIncome)
	fmt.Printf("   - Ingresos recurrentes: $%.0f\n", projection.MonthlyIncome)
	fmt.Printf("   - TOTAL INGRESOS: $%.0f\n", totalMonthlyIncome)
	fmt.Printf("   - Gastos base: $%.0f\n", baseMonthlyExpenses)
	fmt.Printf("   - Gastos recurrentes: $%.0f\n", projection.MonthlyExpenses)
	fmt.Printf("   - TOTAL GASTOS: $%.0f\n", totalMonthlyExpenses)

	return totalMonthlyIncome, totalMonthlyExpenses, nil
}
