package ports

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// RecurringTransactionRepository defines the interface for recurring transaction persistence
type RecurringTransactionRepository interface {
	// CRUD operations
	Create(ctx context.Context, transaction *domain.RecurringTransaction) error
	GetByID(ctx context.Context, userID, transactionID string) (*domain.RecurringTransaction, error)
	GetByUserID(ctx context.Context, userID string, filters RecurringTransactionFilters) ([]*domain.RecurringTransaction, error)
	Update(ctx context.Context, transaction *domain.RecurringTransaction) error
	Delete(ctx context.Context, userID, transactionID string) error

	// Specialized queries
	GetPendingExecutions(ctx context.Context, beforeDate time.Time) ([]*domain.RecurringTransaction, error)
	GetPendingNotifications(ctx context.Context, beforeDate time.Time) ([]*domain.RecurringTransaction, error)
	GetByFrequency(ctx context.Context, userID, frequency string) ([]*domain.RecurringTransaction, error)
	GetByType(ctx context.Context, userID, transactionType string) ([]*domain.RecurringTransaction, error)
	GetActiveTransactions(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error)

	// Analytics
	GetTotalRecurringAmount(ctx context.Context, userID, transactionType string) (float64, error)
	GetRecurringProjection(ctx context.Context, userID string, months int) (*RecurringProjection, error)
}

// RecurringTransactionUseCase defines the business logic interface
type RecurringTransactionUseCase interface {
	// CRUD operations
	CreateRecurringTransaction(ctx context.Context, request *CreateRecurringTransactionRequest) (*RecurringTransactionResponse, error)
	GetRecurringTransaction(ctx context.Context, userID, transactionID string) (*RecurringTransactionResponse, error)
	ListRecurringTransactions(ctx context.Context, userID string, filters RecurringTransactionFilters) (*ListRecurringTransactionsResponse, error)
	UpdateRecurringTransaction(ctx context.Context, userID, transactionID string, request *UpdateRecurringTransactionRequest) (*RecurringTransactionResponse, error)
	DeleteRecurringTransaction(ctx context.Context, userID, transactionID string) error

	// Transaction control
	PauseRecurringTransaction(ctx context.Context, userID, transactionID string) error
	ResumeRecurringTransaction(ctx context.Context, userID, transactionID string) error
	ExecuteRecurringTransaction(ctx context.Context, userID, transactionID string) (*ExecutionResult, error)

	// Batch operations
	ProcessPendingTransactions(ctx context.Context) (*BatchProcessResult, error)
	SendPendingNotifications(ctx context.Context) (*NotificationResult, error)

	// Analytics
	GetRecurringTransactionsDashboard(ctx context.Context, userID string) (*RecurringDashboardResponse, error)
	GetCashFlowProjection(ctx context.Context, userID string, months int) (*CashFlowProjectionResponse, error)
}

// RecurringTransactionNotificationService defines notification interface
type RecurringTransactionNotificationService interface {
	SendUpcomingTransactionNotification(ctx context.Context, transaction *domain.RecurringTransaction) error
	SendTransactionExecutedNotification(ctx context.Context, transaction *domain.RecurringTransaction, success bool) error
	SendTransactionFailedNotification(ctx context.Context, transaction *domain.RecurringTransaction, reason string) error
}

// RecurringTransactionExecutorService defines execution interface
type RecurringTransactionExecutorService interface {
	ExecuteTransaction(ctx context.Context, transaction *domain.RecurringTransaction) error
	CreateTransactionFromRecurring(ctx context.Context, recurring *domain.RecurringTransaction) error
}

// DTOs and Request/Response types

// CreateRecurringTransactionRequest represents the request to create a recurring transaction
type CreateRecurringTransactionRequest struct {
	UserID        string  `json:"user_id" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Description   string  `json:"description" validate:"required"`
	CategoryID    string  `json:"category_id,omitempty"`
	Type          string  `json:"type" validate:"required,oneof=income expense"`
	Frequency     string  `json:"frequency" validate:"required,oneof=daily weekly monthly yearly"`
	NextDate      string  `json:"next_date" validate:"required"` // Format: YYYY-MM-DD
	AutoCreate    bool    `json:"auto_create"`
	NotifyBefore  int     `json:"notify_before" validate:"min=0"`
	EndDate       *string `json:"end_date,omitempty"` // Format: YYYY-MM-DD
	MaxExecutions *int    `json:"max_executions,omitempty" validate:"omitempty,gt=0"`
}

// UpdateRecurringTransactionRequest represents the request to update a recurring transaction
type UpdateRecurringTransactionRequest struct {
	Amount        *float64 `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Description   *string  `json:"description,omitempty"`
	CategoryID    *string  `json:"category_id,omitempty"`
	Frequency     *string  `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly yearly"`
	NextDate      *string  `json:"next_date,omitempty"` // Format: YYYY-MM-DD
	AutoCreate    *bool    `json:"auto_create,omitempty"`
	NotifyBefore  *int     `json:"notify_before,omitempty" validate:"omitempty,min=0"`
	EndDate       *string  `json:"end_date,omitempty"` // Format: YYYY-MM-DD
	MaxExecutions *int     `json:"max_executions,omitempty" validate:"omitempty,gt=0"`
}

// RecurringTransactionResponse represents a recurring transaction response
type RecurringTransactionResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	Amount           float64 `json:"amount"`
	Description      string  `json:"description"`
	CategoryID       string  `json:"category_id,omitempty"`
	CategoryName     string  `json:"category_name,omitempty"`
	Type             string  `json:"type"`
	TypeDisplay      string  `json:"type_display"`
	Frequency        string  `json:"frequency"`
	FrequencyDisplay string  `json:"frequency_display"`
	NextDate         string  `json:"next_date"`
	LastExecuted     string  `json:"last_executed,omitempty"`
	IsActive         bool    `json:"is_active"`
	AutoCreate       bool    `json:"auto_create"`
	NotifyBefore     int     `json:"notify_before"`
	EndDate          string  `json:"end_date,omitempty"`
	ExecutionCount   int     `json:"execution_count"`
	MaxExecutions    *int    `json:"max_executions,omitempty"`
	DaysUntilNext    int     `json:"days_until_next"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// ListRecurringTransactionsResponse represents the response for listing recurring transactions
type ListRecurringTransactionsResponse struct {
	Transactions []*RecurringTransactionResponse `json:"transactions"`
	Summary      *RecurringTransactionsSummary   `json:"summary"`
	Pagination   *PaginationInfo                 `json:"pagination"`
}

// RecurringTransactionsSummary contains summary information
type RecurringTransactionsSummary struct {
	TotalActive         int     `json:"total_active"`
	TotalInactive       int     `json:"total_inactive"`
	MonthlyIncomeTotal  float64 `json:"monthly_income_total"`
	MonthlyExpenseTotal float64 `json:"monthly_expense_total"`
	NetMonthlyRecurring float64 `json:"net_monthly_recurring"`
	NextExecutionDate   string  `json:"next_execution_date,omitempty"`
	PendingExecutions   int     `json:"pending_executions"`
}

// RecurringDashboardResponse represents the dashboard data
type RecurringDashboardResponse struct {
	Summary              *RecurringTransactionsSummary   `json:"summary"`
	UpcomingTransactions []*RecurringTransactionResponse `json:"upcoming_transactions"`
	RecentExecutions     []*ExecutionHistoryItem         `json:"recent_executions"`
	CategoryBreakdown    []*CategoryBreakdownItem        `json:"category_breakdown"`
	FrequencyBreakdown   []*FrequencyBreakdownItem       `json:"frequency_breakdown"`
}

// ExecutionHistoryItem represents an execution history item
type ExecutionHistoryItem struct {
	TransactionID        string  `json:"transaction_id"`
	Description          string  `json:"description"`
	Amount               float64 `json:"amount"`
	Type                 string  `json:"type"`
	ExecutedAt           string  `json:"executed_at"`
	Success              bool    `json:"success"`
	CreatedTransactionID string  `json:"created_transaction_id,omitempty"`
}

// CategoryBreakdownItem represents category breakdown
type CategoryBreakdownItem struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Count        int     `json:"count"`
	Type         string  `json:"type"`
}

// FrequencyBreakdownItem represents frequency breakdown
type FrequencyBreakdownItem struct {
	Frequency        string  `json:"frequency"`
	FrequencyDisplay string  `json:"frequency_display"`
	Count            int     `json:"count"`
	TotalAmount      float64 `json:"total_amount"`
}

// CashFlowProjectionResponse represents cash flow projection
type CashFlowProjectionResponse struct {
	ProjectionMonths   int                  `json:"projection_months"`
	MonthlyProjections []*MonthlyProjection `json:"monthly_projections"`
	Summary            *ProjectionSummary   `json:"summary"`
}

// MonthlyProjection represents a monthly projection
type MonthlyProjection struct {
	Month         string  `json:"month"`         // YYYY-MM
	MonthDisplay  string  `json:"month_display"` // "January 2024"
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	NetAmount     float64 `json:"net_amount"`
	CumulativeNet float64 `json:"cumulative_net"`
}

// ProjectionSummary contains projection summary
type ProjectionSummary struct {
	TotalProjectedIncome   float64 `json:"total_projected_income"`
	TotalProjectedExpenses float64 `json:"total_projected_expenses"`
	NetProjectedAmount     float64 `json:"net_projected_amount"`
	AverageMonthlyIncome   float64 `json:"average_monthly_income"`
	AverageMonthlyExpenses float64 `json:"average_monthly_expenses"`
}

// RecurringTransactionFilters represents filters for listing transactions
type RecurringTransactionFilters struct {
	Type       string `json:"type,omitempty"`      // "income", "expense"
	Frequency  string `json:"frequency,omitempty"` // "daily", "weekly", "monthly", "yearly"
	IsActive   *bool  `json:"is_active,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	SortBy     string `json:"sort_by,omitempty"`    // "next_date", "amount", "created_at"
	SortOrder  string `json:"sort_order,omitempty"` // "asc", "desc"
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
}

// ExecutionResult represents the result of executing a recurring transaction
type ExecutionResult struct {
	Success              bool   `json:"success"`
	CreatedTransactionID string `json:"created_transaction_id,omitempty"`
	Message              string `json:"message"`
	NextExecutionDate    string `json:"next_execution_date,omitempty"`
}

// BatchProcessResult represents the result of batch processing
type BatchProcessResult struct {
	ProcessedCount int                `json:"processed_count"`
	SuccessCount   int                `json:"success_count"`
	FailureCount   int                `json:"failure_count"`
	Results        []*ExecutionResult `json:"results"`
	Errors         []string           `json:"errors,omitempty"`
}

// NotificationResult represents the result of sending notifications
type NotificationResult struct {
	SentCount    int      `json:"sent_count"`
	FailureCount int      `json:"failure_count"`
	Errors       []string `json:"errors,omitempty"`
}

// RecurringProjection represents projection data from repository
type RecurringProjection struct {
	MonthlyIncome   float64 `json:"monthly_income"`
	MonthlyExpenses float64 `json:"monthly_expenses"`
	NetMonthly      float64 `json:"net_monthly"`
}

// PaginationInfo contains pagination information
type PaginationInfo struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrev     bool `json:"has_prev"`
}
