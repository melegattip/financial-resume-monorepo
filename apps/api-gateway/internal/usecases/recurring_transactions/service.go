package recurring_transactions

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// Service implements the recurring transactions use case
type Service struct {
	recurringRepo   ports.RecurringTransactionRepository
	expenseRepo     baseRepo.ExpenseRepository
	incomeRepo      baseRepo.IncomeRepository
	categoryRepo    baseRepo.CategoryRepository
	notificationSvc ports.RecurringTransactionNotificationService
	executorSvc     ports.RecurringTransactionExecutorService
}

// NewService creates a new recurring transactions service
func NewService(
	recurringRepo ports.RecurringTransactionRepository,
	expenseRepo baseRepo.ExpenseRepository,
	incomeRepo baseRepo.IncomeRepository,
	categoryRepo baseRepo.CategoryRepository,
	notificationSvc ports.RecurringTransactionNotificationService,
	executorSvc ports.RecurringTransactionExecutorService,
) ports.RecurringTransactionUseCase {
	return &Service{
		recurringRepo:   recurringRepo,
		expenseRepo:     expenseRepo,
		incomeRepo:      incomeRepo,
		categoryRepo:    categoryRepo,
		notificationSvc: notificationSvc,
		executorSvc:     executorSvc,
	}
}

// CreateRecurringTransaction creates a new recurring transaction
func (s *Service) CreateRecurringTransaction(ctx context.Context, request *ports.CreateRecurringTransactionRequest) (*ports.RecurringTransactionResponse, error) {
	// Validate request
	if err := s.validateCreateRequest(request); err != nil {
		return nil, err
	}

	// Validate category if provided
	if request.CategoryID != "" {
		if _, err := s.categoryRepo.Get(request.UserID, request.CategoryID); err != nil {
			return nil, errors.NewBadRequest("La categoría especificada no existe")
		}
	}

	// Parse dates
	nextDate, err := time.Parse("2006-01-02", request.NextDate)
	if err != nil {
		return nil, errors.NewBadRequest("Formato de fecha inválido para next_date. Use YYYY-MM-DD")
	}

	var endDate *time.Time
	if request.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *request.EndDate)
		if err != nil {
			return nil, errors.NewBadRequest("Formato de fecha inválido para end_date. Use YYYY-MM-DD")
		}
		endDate = &parsed
	}

	// Build recurring transaction
	builder := domain.NewRecurringTransactionBuilder().
		SetID(domain.NewRecurringTransactionID()).
		SetUserID(request.UserID).
		SetAmount(request.Amount).
		SetDescription(request.Description).
		SetCategoryID(request.CategoryID).
		SetType(request.Type).
		SetFrequency(request.Frequency).
		SetNextDate(nextDate).
		SetAutoCreate(request.AutoCreate).
		SetNotifyBefore(request.NotifyBefore).
		SetEndDate(endDate).
		SetMaxExecutions(request.MaxExecutions)

	transaction, err := builder.Build()
	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("Error validando transacción recurrente: %v", err))
	}

	// Save to repository
	if err := s.recurringRepo.Create(ctx, transaction); err != nil {
		return nil, fmt.Errorf("error creando transacción recurrente: %w", err)
	}

	// Convert to response
	return s.toRecurringTransactionResponse(ctx, transaction), nil
}

// GetRecurringTransaction gets a recurring transaction by ID
func (s *Service) GetRecurringTransaction(ctx context.Context, userID, transactionID string) (*ports.RecurringTransactionResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}
	if transactionID == "" {
		return nil, errors.NewBadRequest("El ID de la transacción es requerido")
	}

	transaction, err := s.recurringRepo.GetByID(ctx, userID, transactionID)
	if err != nil {
		return nil, err
	}

	return s.toRecurringTransactionResponse(ctx, transaction), nil
}

// ListRecurringTransactions lists recurring transactions with filters
func (s *Service) ListRecurringTransactions(ctx context.Context, userID string, filters ports.RecurringTransactionFilters) (*ports.ListRecurringTransactionsResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	// Get transactions from repository
	transactions, err := s.recurringRepo.GetByUserID(ctx, userID, filters)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones recurrentes: %w", err)
	}

	// Convert to responses
	responses := make([]*ports.RecurringTransactionResponse, len(transactions))
	for i, transaction := range transactions {
		responses[i] = s.toRecurringTransactionResponse(ctx, transaction)
	}

	// Calculate summary
	summary := s.calculateSummary(transactions)

	// Calculate pagination
	pagination := s.calculatePagination(len(transactions), filters.Limit, filters.Offset)

	return &ports.ListRecurringTransactionsResponse{
		Transactions: responses,
		Summary:      summary,
		Pagination:   pagination,
	}, nil
}

// UpdateRecurringTransaction updates a recurring transaction
func (s *Service) UpdateRecurringTransaction(ctx context.Context, userID, transactionID string, request *ports.UpdateRecurringTransactionRequest) (*ports.RecurringTransactionResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}
	if transactionID == "" {
		return nil, errors.NewBadRequest("El ID de la transacción es requerido")
	}

	// Get existing transaction
	transaction, err := s.recurringRepo.GetByID(ctx, userID, transactionID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if err := s.applyUpdates(ctx, transaction, request); err != nil {
		return nil, err
	}

	// Validate updated transaction
	if err := transaction.Validate(); err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("Error validando transacción actualizada: %v", err))
	}

	// Save changes
	transaction.UpdatedAt = time.Now()
	if err := s.recurringRepo.Update(ctx, transaction); err != nil {
		return nil, fmt.Errorf("error actualizando transacción recurrente: %w", err)
	}

	return s.toRecurringTransactionResponse(ctx, transaction), nil
}

// DeleteRecurringTransaction deletes a recurring transaction
func (s *Service) DeleteRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	if userID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}
	if transactionID == "" {
		return errors.NewBadRequest("El ID de la transacción es requerido")
	}

	return s.recurringRepo.Delete(ctx, userID, transactionID)
}

// PauseRecurringTransaction pauses a recurring transaction
func (s *Service) PauseRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	transaction, err := s.recurringRepo.GetByID(ctx, userID, transactionID)
	if err != nil {
		return err
	}

	transaction.Pause()
	return s.recurringRepo.Update(ctx, transaction)
}

// ResumeRecurringTransaction resumes a recurring transaction
func (s *Service) ResumeRecurringTransaction(ctx context.Context, userID, transactionID string) error {
	transaction, err := s.recurringRepo.GetByID(ctx, userID, transactionID)
	if err != nil {
		return err
	}

	transaction.Resume()
	return s.recurringRepo.Update(ctx, transaction)
}

// ExecuteRecurringTransaction manually executes a recurring transaction
func (s *Service) ExecuteRecurringTransaction(ctx context.Context, userID, transactionID string) (*ports.ExecutionResult, error) {
	transaction, err := s.recurringRepo.GetByID(ctx, userID, transactionID)
	if err != nil {
		return nil, err
	}

	if !transaction.IsActive {
		return &ports.ExecutionResult{
			Success: false,
			Message: "La transacción recurrente está pausada",
		}, nil
	}

	// Execute the transaction
	var createdTransactionID string
	if transaction.AutoCreate {
		if err := s.executorSvc.CreateTransactionFromRecurring(ctx, transaction); err != nil {
			// Send failure notification
			if s.notificationSvc != nil {
				s.notificationSvc.SendTransactionFailedNotification(ctx, transaction, err.Error())
			}
			return &ports.ExecutionResult{
				Success: false,
				Message: fmt.Sprintf("Error creando transacción: %v", err),
			}, nil
		}
		createdTransactionID = fmt.Sprintf("created-from-%s", transaction.ID)
	}

	// Mark as executed
	transaction.Execute()
	if err := s.recurringRepo.Update(ctx, transaction); err != nil {
		log.Printf("Error actualizando transacción recurrente después de ejecución: %v", err)
	}

	// Send success notification
	if s.notificationSvc != nil {
		s.notificationSvc.SendTransactionExecutedNotification(ctx, transaction, true)
	}

	return &ports.ExecutionResult{
		Success:              true,
		CreatedTransactionID: createdTransactionID,
		Message:              "Transacción ejecutada exitosamente",
		NextExecutionDate:    transaction.NextDate.Format("2006-01-02"),
	}, nil
}

// ProcessPendingTransactions processes all pending recurring transactions
func (s *Service) ProcessPendingTransactions(ctx context.Context) (*ports.BatchProcessResult, error) {
	now := time.Now()
	pendingTransactions, err := s.recurringRepo.GetPendingExecutions(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones pendientes: %w", err)
	}

	result := &ports.BatchProcessResult{
		ProcessedCount: len(pendingTransactions),
		Results:        make([]*ports.ExecutionResult, 0, len(pendingTransactions)),
	}

	for _, transaction := range pendingTransactions {
		execResult, err := s.ExecuteRecurringTransaction(ctx, transaction.UserID, transaction.ID)
		if err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Error ejecutando %s: %v", transaction.ID, err))
			continue
		}

		result.Results = append(result.Results, execResult)
		if execResult.Success {
			result.SuccessCount++
		} else {
			result.FailureCount++
		}
	}

	return result, nil
}

// SendPendingNotifications sends notifications for upcoming transactions
func (s *Service) SendPendingNotifications(ctx context.Context) (*ports.NotificationResult, error) {
	if s.notificationSvc == nil {
		return &ports.NotificationResult{}, nil
	}

	now := time.Now()
	pendingNotifications, err := s.recurringRepo.GetPendingNotifications(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo notificaciones pendientes: %w", err)
	}

	result := &ports.NotificationResult{}

	for _, transaction := range pendingNotifications {
		if err := s.notificationSvc.SendUpcomingTransactionNotification(ctx, transaction); err != nil {
			result.FailureCount++
			result.Errors = append(result.Errors, fmt.Sprintf("Error enviando notificación para %s: %v", transaction.ID, err))
		} else {
			result.SentCount++
		}
	}

	return result, nil
}

// GetRecurringTransactionsDashboard gets dashboard data
func (s *Service) GetRecurringTransactionsDashboard(ctx context.Context, userID string) (*ports.RecurringDashboardResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	// Get all transactions (active and inactive) for summary calculation
	allTransactions, err := s.recurringRepo.GetByUserID(ctx, userID, ports.RecurringTransactionFilters{})
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones: %w", err)
	}

	// Get only active transactions for upcoming and breakdowns
	activeTransactions, err := s.recurringRepo.GetActiveTransactions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo transacciones activas: %w", err)
	}

	// Calculate summary using ALL transactions (active and inactive)
	summary := s.calculateSummary(allTransactions)

	// Get upcoming transactions (next 7 days) from active transactions only
	upcomingTransactions := s.getUpcomingTransactions(activeTransactions, 7)
	upcomingResponses := make([]*ports.RecurringTransactionResponse, len(upcomingTransactions))
	for i, transaction := range upcomingTransactions {
		upcomingResponses[i] = s.toRecurringTransactionResponse(ctx, transaction)
	}

	// Calculate breakdowns using active transactions only
	categoryBreakdown := s.calculateCategoryBreakdown(ctx, activeTransactions)
	frequencyBreakdown := s.calculateFrequencyBreakdown(activeTransactions)

	return &ports.RecurringDashboardResponse{
		Summary:              summary,
		UpcomingTransactions: upcomingResponses,
		RecentExecutions:     []*ports.ExecutionHistoryItem{}, // TODO: Implement execution history
		CategoryBreakdown:    categoryBreakdown,
		FrequencyBreakdown:   frequencyBreakdown,
	}, nil
}

// GetCashFlowProjection gets cash flow projection
func (s *Service) GetCashFlowProjection(ctx context.Context, userID string, months int) (*ports.CashFlowProjectionResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}
	if months <= 0 || months > 24 {
		return nil, errors.NewBadRequest("Los meses deben estar entre 1 y 24")
	}

	// Get projection data from repository
	projection, err := s.recurringRepo.GetRecurringProjection(ctx, userID, months)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo proyección: %w", err)
	}

	// Generate monthly projections
	monthlyProjections := s.generateMonthlyProjections(projection, months)

	// Calculate summary
	summary := s.calculateProjectionSummary(monthlyProjections)

	return &ports.CashFlowProjectionResponse{
		ProjectionMonths:   months,
		MonthlyProjections: monthlyProjections,
		Summary:            summary,
	}, nil
}

// Helper methods

func (s *Service) validateCreateRequest(request *ports.CreateRecurringTransactionRequest) error {
	if request.UserID == "" {
		return errors.NewBadRequest("El ID del usuario es requerido")
	}
	if request.Amount <= 0 {
		return errors.NewBadRequest("El monto debe ser mayor a 0")
	}
	if strings.TrimSpace(request.Description) == "" {
		return errors.NewBadRequest("La descripción es requerida")
	}
	if request.Type != "income" && request.Type != "expense" {
		return errors.NewBadRequest("El tipo debe ser 'income' o 'expense'")
	}
	validFrequencies := []string{"daily", "weekly", "monthly", "yearly"}
	if !contains(validFrequencies, request.Frequency) {
		return errors.NewBadRequest("La frecuencia debe ser: daily, weekly, monthly, yearly")
	}
	if request.NotifyBefore < 0 {
		return errors.NewBadRequest("Los días de notificación no pueden ser negativos")
	}
	return nil
}

func (s *Service) applyUpdates(ctx context.Context, transaction *domain.RecurringTransaction, request *ports.UpdateRecurringTransactionRequest) error {
	if request.Amount != nil {
		transaction.Amount = *request.Amount
	}
	if request.Description != nil {
		transaction.Description = *request.Description
	}
	if request.CategoryID != nil {
		if *request.CategoryID != "" {
			// Validate category exists
			if _, err := s.categoryRepo.Get(transaction.UserID, *request.CategoryID); err != nil {
				return errors.NewBadRequest("La categoría especificada no existe")
			}
			transaction.CategoryID = request.CategoryID
		} else {
			transaction.CategoryID = nil
		}
	}
	if request.Frequency != nil {
		transaction.Frequency = *request.Frequency
	}
	if request.NextDate != nil {
		nextDate, err := time.Parse("2006-01-02", *request.NextDate)
		if err != nil {
			return errors.NewBadRequest("Formato de fecha inválido para next_date")
		}
		transaction.NextDate = nextDate
	}
	if request.AutoCreate != nil {
		transaction.AutoCreate = *request.AutoCreate
	}
	if request.NotifyBefore != nil {
		transaction.NotifyBefore = *request.NotifyBefore
	}
	if request.EndDate != nil {
		if *request.EndDate != "" {
			endDate, err := time.Parse("2006-01-02", *request.EndDate)
			if err != nil {
				return errors.NewBadRequest("Formato de fecha inválido para end_date")
			}
			transaction.EndDate = &endDate
		} else {
			transaction.EndDate = nil
		}
	}
	if request.MaxExecutions != nil {
		transaction.MaxExecutions = request.MaxExecutions
	}
	return nil
}

func (s *Service) toRecurringTransactionResponse(ctx context.Context, transaction *domain.RecurringTransaction) *ports.RecurringTransactionResponse {
	response := &ports.RecurringTransactionResponse{
		ID:               transaction.ID,
		UserID:           transaction.UserID,
		Amount:           transaction.Amount,
		Description:      transaction.Description,
		Type:             transaction.Type,
		TypeDisplay:      transaction.GetTypeDisplay(),
		Frequency:        transaction.Frequency,
		FrequencyDisplay: transaction.GetFrequencyDisplay(),
		NextDate:         transaction.NextDate.Format("2006-01-02"),
		IsActive:         transaction.IsActive,
		AutoCreate:       transaction.AutoCreate,
		NotifyBefore:     transaction.NotifyBefore,
		ExecutionCount:   transaction.ExecutionCount,
		MaxExecutions:    transaction.MaxExecutions,
		DaysUntilNext:    transaction.GetDaysUntilNext(),
		CreatedAt:        transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if transaction.CategoryID != nil {
		response.CategoryID = *transaction.CategoryID
		// Try to get category name
		if category, err := s.categoryRepo.Get(transaction.UserID, *transaction.CategoryID); err == nil {
			response.CategoryName = category.Name
		}
	}

	if transaction.LastExecuted != nil {
		response.LastExecuted = transaction.LastExecuted.Format("2006-01-02T15:04:05Z07:00")
	}

	if transaction.EndDate != nil {
		response.EndDate = transaction.EndDate.Format("2006-01-02")
	}

	return response
}

func (s *Service) calculateSummary(transactions []*domain.RecurringTransaction) *ports.RecurringTransactionsSummary {
	log.Printf("🔍 Calculando resumen de transacciones")
	
	summary := &ports.RecurringTransactionsSummary{}

	var nextExecution *time.Time


	for _, transaction := range transactions {
		// Count active/inactive transactions
		if transaction.IsActive {
			summary.TotalActive++

			// Find next execution (only for active transactions)
			if nextExecution == nil || transaction.NextDate.Before(*nextExecution) {
				nextExecution = &transaction.NextDate
			}

			// Count pending executions (only for active transactions)
			if transaction.ShouldExecute() {
				summary.PendingExecutions++
			}
		} else {
			summary.TotalInactive++
		}

		// Calculate monthly amounts for ALL transactions (active and inactive)
		monthlyAmount := s.calculateMonthlyAmount(transaction)
		if transaction.Type == "income" {
			summary.MonthlyIncomeTotal += monthlyAmount
		} else {
			summary.MonthlyExpenseTotal += monthlyAmount
		}
	}

	summary.NetMonthlyRecurring = summary.MonthlyIncomeTotal - summary.MonthlyExpenseTotal

	if nextExecution != nil {
		summary.NextExecutionDate = nextExecution.Format("2006-01-02")
	}

	return summary
}

func (s *Service) calculateMonthlyAmount(transaction *domain.RecurringTransaction) float64 {
	switch transaction.Frequency {
	case "daily":
		return transaction.Amount * 30 // Approximate
	case "weekly":
		return transaction.Amount * 4.33 // 52 weeks / 12 months
	case "monthly":
		return transaction.Amount
	case "yearly":
		return transaction.Amount / 12
	default:
		return transaction.Amount
	}
}

func (s *Service) calculatePagination(totalItems, limit, offset int) *ports.PaginationInfo {
	if limit <= 0 {
		limit = 20 // Default
	}

	currentPage := (offset / limit) + 1
	totalPages := (totalItems + limit - 1) / limit

	return &ports.PaginationInfo{
		CurrentPage: currentPage,
		PageSize:    limit,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
}

func (s *Service) getUpcomingTransactions(transactions []*domain.RecurringTransaction, days int) []*domain.RecurringTransaction {
	cutoffDate := time.Now().AddDate(0, 0, days)
	var upcoming []*domain.RecurringTransaction

	for _, transaction := range transactions {
		if transaction.IsActive && transaction.NextDate.Before(cutoffDate) {
			upcoming = append(upcoming, transaction)
		}
	}

	// Sort by next date
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].NextDate.Before(upcoming[j].NextDate)
	})

	return upcoming
}

func (s *Service) calculateCategoryBreakdown(ctx context.Context, transactions []*domain.RecurringTransaction) []*ports.CategoryBreakdownItem {
	categoryMap := make(map[string]*ports.CategoryBreakdownItem)

	for _, transaction := range transactions {
		if !transaction.IsActive {
			continue
		}

		categoryID := "uncategorized"
		categoryName := "Sin categoría"

		if transaction.CategoryID != nil {
			categoryID = *transaction.CategoryID
			if category, err := s.categoryRepo.Get(transaction.UserID, categoryID); err == nil {
				categoryName = category.Name
			}
		}

		key := fmt.Sprintf("%s-%s", categoryID, transaction.Type)
		if item, exists := categoryMap[key]; exists {
			item.Amount += s.calculateMonthlyAmount(transaction)
			item.Count++
		} else {
			categoryMap[key] = &ports.CategoryBreakdownItem{
				CategoryID:   categoryID,
				CategoryName: categoryName,
				Amount:       s.calculateMonthlyAmount(transaction),
				Count:        1,
				Type:         transaction.Type,
			}
		}
	}

	// Convert to slice
	var breakdown []*ports.CategoryBreakdownItem
	for _, item := range categoryMap {
		breakdown = append(breakdown, item)
	}

	// Sort by amount descending
	sort.Slice(breakdown, func(i, j int) bool {
		return breakdown[i].Amount > breakdown[j].Amount
	})

	return breakdown
}

func (s *Service) calculateFrequencyBreakdown(transactions []*domain.RecurringTransaction) []*ports.FrequencyBreakdownItem {
	frequencyMap := make(map[string]*ports.FrequencyBreakdownItem)

	for _, transaction := range transactions {
		if !transaction.IsActive {
			continue
		}

		if item, exists := frequencyMap[transaction.Frequency]; exists {
			item.Count++
			item.TotalAmount += transaction.Amount
		} else {
			frequencyMap[transaction.Frequency] = &ports.FrequencyBreakdownItem{
				Frequency:        transaction.Frequency,
				FrequencyDisplay: transaction.GetFrequencyDisplay(),
				Count:            1,
				TotalAmount:      transaction.Amount,
			}
		}
	}

	// Convert to slice
	var breakdown []*ports.FrequencyBreakdownItem
	for _, item := range frequencyMap {
		breakdown = append(breakdown, item)
	}

	return breakdown
}

func (s *Service) generateMonthlyProjections(projection *ports.RecurringProjection, months int) []*ports.MonthlyProjection {
	projections := make([]*ports.MonthlyProjection, months)
	cumulativeNet := 0.0

	for i := 0; i < months; i++ {
		date := time.Now().AddDate(0, i, 0)
		monthlyNet := projection.MonthlyIncome - projection.MonthlyExpenses
		cumulativeNet += monthlyNet

		projections[i] = &ports.MonthlyProjection{
			Month:         date.Format("2006-01"),
			MonthDisplay:  date.Format("January 2006"),
			Income:        projection.MonthlyIncome,
			Expenses:      projection.MonthlyExpenses,
			NetAmount:     monthlyNet,
			CumulativeNet: cumulativeNet,
		}
	}

	return projections
}

func (s *Service) calculateProjectionSummary(projections []*ports.MonthlyProjection) *ports.ProjectionSummary {
	if len(projections) == 0 {
		return &ports.ProjectionSummary{}
	}

	totalIncome := 0.0
	totalExpenses := 0.0

	for _, projection := range projections {
		totalIncome += projection.Income
		totalExpenses += projection.Expenses
	}

	return &ports.ProjectionSummary{
		TotalProjectedIncome:   totalIncome,
		TotalProjectedExpenses: totalExpenses,
		NetProjectedAmount:     totalIncome - totalExpenses,
		AverageMonthlyIncome:   totalIncome / float64(len(projections)),
		AverageMonthlyExpenses: totalExpenses / float64(len(projections)),
	}
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
