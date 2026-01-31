package ports

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// BudgetRepository defines the interface for budget data access (Repository pattern)
type BudgetRepository interface {
	// Create creates a new budget
	Create(ctx context.Context, budget *domain.Budget) error

	// GetByID retrieves a budget by ID
	GetByID(ctx context.Context, userID, budgetID string) (*domain.Budget, error)

	// GetByCategory retrieves a budget by category and period
	GetByCategory(ctx context.Context, userID, categoryID string, period domain.BudgetPeriod) (*domain.Budget, error)

	// List retrieves all budgets for a user
	List(ctx context.Context, userID string) ([]*domain.Budget, error)

	// ListActive retrieves all active budgets for a user
	ListActive(ctx context.Context, userID string) ([]*domain.Budget, error)

	// Update updates an existing budget
	Update(ctx context.Context, budget *domain.Budget) error

	// Delete removes a budget
	Delete(ctx context.Context, userID, budgetID string) error

	// GetExpensesForPeriod gets total expenses for a category in a period (for calculating spent amount)
	GetExpensesForPeriod(ctx context.Context, userID, categoryID string, startDate, endDate time.Time) (float64, error)
}

// BudgetUseCase defines the business logic interface for budgets (Clean Architecture)
type BudgetUseCase interface {
	// CreateBudget creates a new budget with validation
	CreateBudget(ctx context.Context, request CreateBudgetRequest) (*CreateBudgetResponse, error)

	// GetBudget retrieves a specific budget
	GetBudget(ctx context.Context, request GetBudgetRequest) (*GetBudgetResponse, error)

	// ListBudgets retrieves all budgets for a user
	ListBudgets(ctx context.Context, request ListBudgetsRequest) (*ListBudgetsResponse, error)

	// UpdateBudget updates an existing budget
	UpdateBudget(ctx context.Context, request UpdateBudgetRequest) (*UpdateBudgetResponse, error)

	// DeleteBudget removes a budget
	DeleteBudget(ctx context.Context, request DeleteBudgetRequest) error

	// GetBudgetStatus gets current status and spent amounts for all budgets
	GetBudgetStatus(ctx context.Context, request GetBudgetStatusRequest) (*GetBudgetStatusResponse, error)

	// RefreshBudgetAmounts recalculates spent amounts for all budgets
	RefreshBudgetAmounts(ctx context.Context, userID string) error
}

// BudgetNotificationService defines interface for budget notifications (Single Responsibility)
type BudgetNotificationService interface {
	// NotifyBudgetAlert sends an alert when budget threshold is reached
	NotifyBudgetAlert(ctx context.Context, budget *domain.Budget, spentPercentage float64) error

	// NotifyBudgetExceeded sends notification when budget is exceeded
	NotifyBudgetExceeded(ctx context.Context, budget *domain.Budget, exceededAmount float64) error

	// NotifyBudgetReset sends notification when budget period resets
	NotifyBudgetReset(ctx context.Context, budget *domain.Budget) error
}

// Request/Response DTOs (Data Transfer Objects)

// CreateBudgetRequest represents a request to create a budget
type CreateBudgetRequest struct {
	UserID     string              `json:"user_id" validate:"required"`
	CategoryID string              `json:"category_id" validate:"required"`
	Amount     float64             `json:"amount" validate:"required,gt=0"`
	Period     domain.BudgetPeriod `json:"period" validate:"required,oneof=monthly weekly yearly"`
	AlertAt    float64             `json:"alert_at,omitempty" validate:"gte=0,lte=1"`
}

// CreateBudgetResponse represents the response after creating a budget
type CreateBudgetResponse struct {
	ID           string              `json:"id"`
	UserID       string              `json:"user_id"`
	CategoryID   string              `json:"category_id"`
	CategoryName string              `json:"category_name"`
	Amount       float64             `json:"amount"`
	SpentAmount  float64             `json:"spent_amount"`
	Period       domain.BudgetPeriod `json:"period"`
	PeriodStart  time.Time           `json:"period_start"`
	PeriodEnd    time.Time           `json:"period_end"`
	AlertAt      float64             `json:"alert_at"`
	Status       domain.BudgetStatus `json:"status"`
	IsActive     bool                `json:"is_active"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

// GetBudgetRequest represents a request to get a specific budget
type GetBudgetRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	BudgetID string `json:"budget_id" validate:"required"`
}

// GetBudgetResponse represents the response for getting a budget
type GetBudgetResponse struct {
	ID               string              `json:"id"`
	UserID           string              `json:"user_id"`
	CategoryID       string              `json:"category_id"`
	CategoryName     string              `json:"category_name"`
	Amount           float64             `json:"amount"`
	SpentAmount      float64             `json:"spent_amount"`
	RemainingAmount  float64             `json:"remaining_amount"`
	SpentPercentage  float64             `json:"spent_percentage"`
	Period           domain.BudgetPeriod `json:"period"`
	PeriodStart      time.Time           `json:"period_start"`
	PeriodEnd        time.Time           `json:"period_end"`
	AlertAt          float64             `json:"alert_at"`
	Status           domain.BudgetStatus `json:"status"`
	IsActive         bool                `json:"is_active"`
	IsAlertTriggered bool                `json:"is_alert_triggered"`
	DaysRemaining    int                 `json:"days_remaining"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}

// ListBudgetsRequest represents a request to list budgets
type ListBudgetsRequest struct {
	UserID     string              `json:"user_id" validate:"required"`
	Period     domain.BudgetPeriod `json:"period,omitempty"`
	CategoryID string              `json:"category_id,omitempty"`
	Status     domain.BudgetStatus `json:"status,omitempty"`
	ActiveOnly bool                `json:"active_only,omitempty"`
}

// ListBudgetsResponse represents the response for listing budgets
type ListBudgetsResponse struct {
	Budgets []GetBudgetResponse  `json:"budgets"`
	Summary domain.BudgetSummary `json:"summary"`
	Count   int                  `json:"count"`
}

// UpdateBudgetRequest represents a request to update a budget
type UpdateBudgetRequest struct {
	UserID   string   `json:"user_id" validate:"required"`
	BudgetID string   `json:"budget_id" validate:"required"`
	Amount   *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	AlertAt  *float64 `json:"alert_at,omitempty" validate:"omitempty,gte=0,lte=1"`
	IsActive *bool    `json:"is_active,omitempty"`
}

// UpdateBudgetResponse represents the response after updating a budget
type UpdateBudgetResponse struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	CategoryID      string              `json:"category_id"`
	CategoryName    string              `json:"category_name"`
	Amount          float64             `json:"amount"`
	SpentAmount     float64             `json:"spent_amount"`
	RemainingAmount float64             `json:"remaining_amount"`
	SpentPercentage float64             `json:"spent_percentage"`
	Period          domain.BudgetPeriod `json:"period"`
	AlertAt         float64             `json:"alert_at"`
	Status          domain.BudgetStatus `json:"status"`
	IsActive        bool                `json:"is_active"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

// DeleteBudgetRequest represents a request to delete a budget
type DeleteBudgetRequest struct {
	UserID   string `json:"user_id" validate:"required"`
	BudgetID string `json:"budget_id" validate:"required"`
}

// GetBudgetStatusRequest represents a request to get budget status
type GetBudgetStatusRequest struct {
	UserID     string              `json:"user_id" validate:"required"`
	CategoryID string              `json:"category_id,omitempty"`
	Period     domain.BudgetPeriod `json:"period,omitempty"`
}

// BudgetStatusItem represents status for a single budget
type BudgetStatusItem struct {
	ID               string              `json:"id"`
	CategoryID       string              `json:"category_id"`
	CategoryName     string              `json:"category_name"`
	Amount           float64             `json:"amount"`
	SpentAmount      float64             `json:"spent_amount"`
	RemainingAmount  float64             `json:"remaining_amount"`
	SpentPercentage  float64             `json:"spent_percentage"`
	Status           domain.BudgetStatus `json:"status"`
	IsAlertTriggered bool                `json:"is_alert_triggered"`
	DaysRemaining    int                 `json:"days_remaining"`
	Period           domain.BudgetPeriod `json:"period"`
}

// GetBudgetStatusResponse represents the response for budget status
type GetBudgetStatusResponse struct {
	Budgets    []BudgetStatusItem   `json:"budgets"`
	Summary    domain.BudgetSummary `json:"summary"`
	TotalSpent float64              `json:"total_spent"`
	UpdatedAt  time.Time            `json:"updated_at"`
}
