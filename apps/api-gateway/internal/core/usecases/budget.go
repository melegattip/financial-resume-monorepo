package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// BudgetService implements BudgetUseCase interface (Single Responsibility Principle)
type BudgetService struct {
	budgetRepo      ports.BudgetRepository
	categoryRepo    baseRepo.CategoryRepository
	expenseRepo     baseRepo.ExpenseRepository
	notificationSvc ports.BudgetNotificationService
}

// NewBudgetService creates a new BudgetService with dependency injection
func NewBudgetService(
	budgetRepo ports.BudgetRepository,
	categoryRepo baseRepo.CategoryRepository,
	expenseRepo baseRepo.ExpenseRepository,
	notificationSvc ports.BudgetNotificationService,
) ports.BudgetUseCase {
	return &BudgetService{
		budgetRepo:      budgetRepo,
		categoryRepo:    categoryRepo,
		expenseRepo:     expenseRepo,
		notificationSvc: notificationSvc,
	}
}

// CreateBudget creates a new budget with validation (Interface Segregation Principle)
func (s *BudgetService) CreateBudget(ctx context.Context, request ports.CreateBudgetRequest) (*ports.CreateBudgetResponse, error) {
	// Validate category exists
	category, err := s.categoryRepo.Get(request.UserID, request.CategoryID)
	if err != nil {
		return nil, errors.NewBadRequest("categoría no encontrada")
	}

	// Check if budget already exists for this category and period
	existingBudget, err := s.budgetRepo.GetByCategory(ctx, request.UserID, request.CategoryID, request.Period)
	if err == nil && existingBudget != nil {
		return nil, errors.NewConflict("ya existe un presupuesto para esta categoría en este período")
	}

	// Build budget using Builder pattern
	budgetBuilder := domain.NewBudgetBuilder().
		SetUserID(request.UserID).
		SetCategoryID(request.CategoryID).
		SetAmount(request.Amount).
		SetPeriod(request.Period)

	if request.AlertAt > 0 {
		budgetBuilder.SetAlertAt(request.AlertAt)
	}

	budget, err := budgetBuilder.Build()
	if err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("error creando presupuesto: %v", err))
	}

	// Calculate initial spent amount
	spentAmount, err := s.calculateSpentAmount(ctx, budget)
	if err != nil {
		return nil, err
	}

	budget.UpdateSpentAmount(spentAmount)

	// Create budget in repository
	if err := s.budgetRepo.Create(ctx, budget); err != nil {
		return nil, fmt.Errorf("error creando presupuesto: %w", err)
	}

	// Check if alert should be triggered
	if budget.IsAlertTriggered() && s.notificationSvc != nil {
		go s.notificationSvc.NotifyBudgetAlert(ctx, budget, budget.GetSpentPercentage())
	}

	return s.buildCreateBudgetResponse(budget, category.Name), nil
}

// GetBudget retrieves a specific budget
func (s *BudgetService) GetBudget(ctx context.Context, request ports.GetBudgetRequest) (*ports.GetBudgetResponse, error) {
	budget, err := s.budgetRepo.GetByID(ctx, request.UserID, request.BudgetID)
	if err != nil {
		return nil, errors.NewResourceNotFound("presupuesto no encontrado")
	}

	// Get category name
	category, err := s.categoryRepo.Get(request.UserID, budget.CategoryID)
	if err != nil {
		return nil, err
	}

	// Refresh spent amount
	spentAmount, err := s.calculateSpentAmount(ctx, budget)
	if err != nil {
		return nil, err
	}

	budget.UpdateSpentAmount(spentAmount)

	// Update in repository if spent amount changed
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return nil, err
	}

	return s.buildGetBudgetResponse(budget, category.Name), nil
}

// ListBudgets retrieves all budgets for a user with optional filters
func (s *BudgetService) ListBudgets(ctx context.Context, request ports.ListBudgetsRequest) (*ports.ListBudgetsResponse, error) {
	var budgets []*domain.Budget
	var err error

	if request.ActiveOnly {
		budgets, err = s.budgetRepo.ListActive(ctx, request.UserID)
	} else {
		budgets, err = s.budgetRepo.List(ctx, request.UserID)
	}

	if err != nil {
		return nil, err
	}

	// Apply filters and refresh spent amounts
	var filteredBudgets []*domain.Budget
	for _, budget := range budgets {
		// Auto-reset budget if period has changed
		if !budget.IsInCurrentPeriod() {
			budget.ResetForNewPeriod()
			// Update in repository
			if err := s.budgetRepo.Update(ctx, budget); err != nil {
				continue // Skip if update fails
			}
		}

		// Apply period filter
		if request.Period != "" && budget.Period != request.Period {
			continue
		}

		// Apply category filter
		if request.CategoryID != "" && budget.CategoryID != request.CategoryID {
			continue
		}

		// Apply status filter
		if request.Status != "" && budget.Status != request.Status {
			continue
		}

		// Refresh spent amount
		spentAmount, err := s.calculateSpentAmount(ctx, budget)
		if err != nil {
			continue // Skip this budget if we can't calculate spent amount
		}

		budget.UpdateSpentAmount(spentAmount)
		filteredBudgets = append(filteredBudgets, budget)
	}

	// Build response
	var budgetResponses []ports.GetBudgetResponse
	for _, budget := range filteredBudgets {
		category, err := s.categoryRepo.Get(request.UserID, budget.CategoryID)
		if err != nil {
			continue // Skip if category not found
		}

		budgetResponses = append(budgetResponses, *s.buildGetBudgetResponse(budget, category.Name))
	}

	// Calculate summary
	summary := s.calculateBudgetSummary(filteredBudgets)

	return &ports.ListBudgetsResponse{
		Budgets: budgetResponses,
		Summary: summary,
		Count:   len(budgetResponses),
	}, nil
}

// UpdateBudget updates an existing budget
func (s *BudgetService) UpdateBudget(ctx context.Context, request ports.UpdateBudgetRequest) (*ports.UpdateBudgetResponse, error) {
	budget, err := s.budgetRepo.GetByID(ctx, request.UserID, request.BudgetID)
	if err != nil {
		return nil, errors.NewResourceNotFound("presupuesto no encontrado")
	}

	// Update fields if provided
	if request.Amount != nil {
		budget.Amount = *request.Amount
	}

	if request.AlertAt != nil {
		budget.AlertAt = *request.AlertAt
	}

	if request.IsActive != nil {
		budget.IsActive = *request.IsActive
	}

	// Validate updated budget
	if err := budget.Validate(); err != nil {
		return nil, errors.NewBadRequest(fmt.Sprintf("datos inválidos: %v", err))
	}

	// Refresh spent amount and update status
	spentAmount, err := s.calculateSpentAmount(ctx, budget)
	if err != nil {
		return nil, err
	}

	budget.UpdateSpentAmount(spentAmount)

	// Update in repository
	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return nil, err
	}

	// Get category name
	category, err := s.categoryRepo.Get(request.UserID, budget.CategoryID)
	if err != nil {
		return nil, err
	}

	return s.buildUpdateBudgetResponse(budget, category.Name), nil
}

// DeleteBudget removes a budget
func (s *BudgetService) DeleteBudget(ctx context.Context, request ports.DeleteBudgetRequest) error {
	// Verify budget exists
	_, err := s.budgetRepo.GetByID(ctx, request.UserID, request.BudgetID)
	if err != nil {
		return errors.NewResourceNotFound("presupuesto no encontrado")
	}

	return s.budgetRepo.Delete(ctx, request.UserID, request.BudgetID)
}

// GetBudgetStatus gets current status for all budgets
func (s *BudgetService) GetBudgetStatus(ctx context.Context, request ports.GetBudgetStatusRequest) (*ports.GetBudgetStatusResponse, error) {
	budgets, err := s.budgetRepo.ListActive(ctx, request.UserID)
	if err != nil {
		return nil, err
	}

	var statusItems []ports.BudgetStatusItem
	var totalSpent float64

	for _, budget := range budgets {
		// Auto-reset budget if period has changed
		if !budget.IsInCurrentPeriod() {
			budget.ResetForNewPeriod()
			// Update in repository
			if err := s.budgetRepo.Update(ctx, budget); err != nil {
				continue // Skip if update fails
			}
		}

		// Apply category filter if provided
		if request.CategoryID != "" && budget.CategoryID != request.CategoryID {
			continue
		}

		// Apply period filter if provided
		if request.Period != "" && budget.Period != request.Period {
			continue
		}

		// Refresh spent amount
		spentAmount, err := s.calculateSpentAmount(ctx, budget)
		if err != nil {
			continue
		}

		budget.UpdateSpentAmount(spentAmount)
		totalSpent += spentAmount

		// Get category name
		category, err := s.categoryRepo.Get(request.UserID, budget.CategoryID)
		if err != nil {
			continue
		}

		// Calculate days remaining
		daysRemaining := int(time.Until(budget.PeriodEnd).Hours() / 24)
		if daysRemaining < 0 {
			daysRemaining = 0
		}

		statusItems = append(statusItems, ports.BudgetStatusItem{
			ID:               budget.ID,
			CategoryID:       budget.CategoryID,
			CategoryName:     category.Name,
			Amount:           budget.Amount,
			SpentAmount:      budget.SpentAmount,
			RemainingAmount:  budget.GetRemainingAmount(),
			SpentPercentage:  budget.GetSpentPercentage(),
			Status:           budget.Status,
			IsAlertTriggered: budget.IsAlertTriggered(),
			DaysRemaining:    daysRemaining,
			Period:           budget.Period,
		})
	}

	// Calculate summary
	summary := s.calculateBudgetSummary(budgets)

	return &ports.GetBudgetStatusResponse{
		Budgets:    statusItems,
		Summary:    summary,
		TotalSpent: totalSpent,
		UpdatedAt:  time.Now(),
	}, nil
}

// RefreshBudgetAmounts recalculates spent amounts for all budgets
func (s *BudgetService) RefreshBudgetAmounts(ctx context.Context, userID string) error {
	budgets, err := s.budgetRepo.ListActive(ctx, userID)
	if err != nil {
		return err
	}

	for _, budget := range budgets {
		// Auto-reset budget if period has changed
		if !budget.IsInCurrentPeriod() {
			budget.ResetForNewPeriod()
		}

		spentAmount, err := s.calculateSpentAmount(ctx, budget)
		if err != nil {
			continue // Skip this budget if calculation fails
		}

		oldSpentAmount := budget.SpentAmount
		budget.UpdateSpentAmount(spentAmount)

		// Update in repository
		if err := s.budgetRepo.Update(ctx, budget); err != nil {
			continue // Skip if update fails
		}

		// Check for notifications
		if budget.IsAlertTriggered() && oldSpentAmount < budget.AlertAt*budget.Amount && s.notificationSvc != nil {
			go s.notificationSvc.NotifyBudgetAlert(ctx, budget, budget.GetSpentPercentage())
		}

		if budget.Status == domain.BudgetStatusExceeded && oldSpentAmount < budget.Amount && s.notificationSvc != nil {
			exceededAmount := budget.SpentAmount - budget.Amount
			go s.notificationSvc.NotifyBudgetExceeded(ctx, budget, exceededAmount)
		}
	}

	return nil
}

// Helper methods (Single Responsibility Principle)

// calculateSpentAmount calculates the spent amount for a budget based on expenses
func (s *BudgetService) calculateSpentAmount(ctx context.Context, budget *domain.Budget) (float64, error) {
	return s.budgetRepo.GetExpensesForPeriod(ctx, budget.UserID, budget.CategoryID, budget.PeriodStart, budget.PeriodEnd)
}

// calculateBudgetSummary calculates summary statistics for a list of budgets
func (s *BudgetService) calculateBudgetSummary(budgets []*domain.Budget) domain.BudgetSummary {
	summary := domain.BudgetSummary{}
	summary.TotalBudgets = len(budgets)

	if len(budgets) == 0 {
		return summary
	}

	var totalUsage float64
	for _, budget := range budgets {
		summary.TotalAllocated += budget.Amount
		summary.TotalSpent += budget.SpentAmount

		switch budget.Status {
		case domain.BudgetStatusOnTrack:
			summary.OnTrackCount++
		case domain.BudgetStatusWarning:
			summary.WarningCount++
		case domain.BudgetStatusExceeded:
			summary.ExceededCount++
		}

		totalUsage += budget.GetSpentPercentage()
	}

	summary.AverageUsage = totalUsage / float64(len(budgets))
	return summary
}

// Response builders (Factory pattern)

func (s *BudgetService) buildCreateBudgetResponse(budget *domain.Budget, categoryName string) *ports.CreateBudgetResponse {
	return &ports.CreateBudgetResponse{
		ID:           budget.ID,
		UserID:       budget.UserID,
		CategoryID:   budget.CategoryID,
		CategoryName: categoryName,
		Amount:       budget.Amount,
		SpentAmount:  budget.SpentAmount,
		Period:       budget.Period,
		PeriodStart:  budget.PeriodStart,
		PeriodEnd:    budget.PeriodEnd,
		AlertAt:      budget.AlertAt,
		Status:       budget.Status,
		IsActive:     budget.IsActive,
		CreatedAt:    budget.CreatedAt,
		UpdatedAt:    budget.UpdatedAt,
	}
}

func (s *BudgetService) buildGetBudgetResponse(budget *domain.Budget, categoryName string) *ports.GetBudgetResponse {
	daysRemaining := int(time.Until(budget.PeriodEnd).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	return &ports.GetBudgetResponse{
		ID:               budget.ID,
		UserID:           budget.UserID,
		CategoryID:       budget.CategoryID,
		CategoryName:     categoryName,
		Amount:           budget.Amount,
		SpentAmount:      budget.SpentAmount,
		RemainingAmount:  budget.GetRemainingAmount(),
		SpentPercentage:  budget.GetSpentPercentage(),
		Period:           budget.Period,
		PeriodStart:      budget.PeriodStart,
		PeriodEnd:        budget.PeriodEnd,
		AlertAt:          budget.AlertAt,
		Status:           budget.Status,
		IsActive:         budget.IsActive,
		IsAlertTriggered: budget.IsAlertTriggered(),
		DaysRemaining:    daysRemaining,
		CreatedAt:        budget.CreatedAt,
		UpdatedAt:        budget.UpdatedAt,
	}
}

func (s *BudgetService) buildUpdateBudgetResponse(budget *domain.Budget, categoryName string) *ports.UpdateBudgetResponse {
	return &ports.UpdateBudgetResponse{
		ID:              budget.ID,
		UserID:          budget.UserID,
		CategoryID:      budget.CategoryID,
		CategoryName:    categoryName,
		Amount:          budget.Amount,
		SpentAmount:     budget.SpentAmount,
		RemainingAmount: budget.GetRemainingAmount(),
		SpentPercentage: budget.GetSpentPercentage(),
		Period:          budget.Period,
		AlertAt:         budget.AlertAt,
		Status:          budget.Status,
		IsActive:        budget.IsActive,
		UpdatedAt:       budget.UpdatedAt,
	}
}
