package ports

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// SavingsGoalRepository defines the interface for savings goal data access (Repository pattern)
type SavingsGoalRepository interface {
	// Create creates a new savings goal
	Create(ctx context.Context, goal *domain.SavingsGoal) error

	// GetByID retrieves a savings goal by ID
	GetByID(ctx context.Context, userID, goalID string) (*domain.SavingsGoal, error)

	// List retrieves all savings goals for a user
	List(ctx context.Context, userID string) ([]*domain.SavingsGoal, error)

	// ListByStatus retrieves savings goals by status
	ListByStatus(ctx context.Context, userID string, status domain.SavingsGoalStatus) ([]*domain.SavingsGoal, error)

	// ListByCategory retrieves savings goals by category
	ListByCategory(ctx context.Context, userID string, category domain.SavingsGoalCategory) ([]*domain.SavingsGoal, error)

	// Update updates an existing savings goal
	Update(ctx context.Context, goal *domain.SavingsGoal) error

	// Delete removes a savings goal
	Delete(ctx context.Context, userID, goalID string) error

	// CreateTransaction creates a savings transaction
	CreateTransaction(ctx context.Context, transaction *domain.SavingsTransaction) error

	// GetTransactionsByGoal retrieves transactions for a specific goal
	GetTransactionsByGoal(ctx context.Context, userID, goalID string) ([]*domain.SavingsTransaction, error)

	// GetTransactionsByUser retrieves all transactions for a user
	GetTransactionsByUser(ctx context.Context, userID string) ([]*domain.SavingsTransaction, error)
}

// SavingsGoalUseCase defines the business logic interface for savings goals (Clean Architecture)
type SavingsGoalUseCase interface {
	// CreateGoal creates a new savings goal with validation
	CreateGoal(ctx context.Context, request CreateSavingsGoalRequest) (*CreateSavingsGoalResponse, error)

	// GetGoal retrieves a specific savings goal
	GetGoal(ctx context.Context, request GetSavingsGoalRequest) (*GetSavingsGoalResponse, error)

	// ListGoals retrieves all savings goals for a user
	ListGoals(ctx context.Context, request ListSavingsGoalsRequest) (*ListSavingsGoalsResponse, error)

	// UpdateGoal updates an existing savings goal
	UpdateGoal(ctx context.Context, request UpdateSavingsGoalRequest) (*UpdateSavingsGoalResponse, error)

	// DeleteGoal removes a savings goal
	DeleteGoal(ctx context.Context, request DeleteSavingsGoalRequest) error

	// AddSavings adds money to a savings goal
	AddSavings(ctx context.Context, request AddSavingsRequest) (*AddSavingsResponse, error)

	// WithdrawSavings removes money from a savings goal
	WithdrawSavings(ctx context.Context, request WithdrawSavingsRequest) (*WithdrawSavingsResponse, error)

	// PauseGoal pauses a savings goal
	PauseGoal(ctx context.Context, request PauseGoalRequest) error

	// ResumeGoal resumes a paused savings goal
	ResumeGoal(ctx context.Context, request ResumeGoalRequest) error

	// CancelGoal cancels a savings goal
	CancelGoal(ctx context.Context, request CancelGoalRequest) error

	// GetGoalSummary gets summary statistics for all goals
	GetGoalSummary(ctx context.Context, request GetGoalSummaryRequest) (*GetGoalSummaryResponse, error)

	// GetGoalTransactions gets transaction history for a goal
	GetGoalTransactions(ctx context.Context, request GetGoalTransactionsRequest) (*GetGoalTransactionsResponse, error)
}

// SavingsGoalNotificationService defines interface for goal notifications (Single Responsibility)
type SavingsGoalNotificationService interface {
	// NotifyGoalAchieved sends notification when goal is achieved
	NotifyGoalAchieved(ctx context.Context, goal *domain.SavingsGoal) error

	// NotifyGoalOverdue sends notification when goal is overdue
	NotifyGoalOverdue(ctx context.Context, goal *domain.SavingsGoal) error

	// NotifyMilestoneReached sends notification when milestone is reached (25%, 50%, 75%)
	NotifyMilestoneReached(ctx context.Context, goal *domain.SavingsGoal, milestone float64) error

	// NotifyAutoSaveExecuted sends notification when auto-save is executed
	NotifyAutoSaveExecuted(ctx context.Context, goal *domain.SavingsGoal, amount float64) error
}

// Request/Response DTOs (Data Transfer Objects)

// CreateSavingsGoalRequest represents a request to create a savings goal
type CreateSavingsGoalRequest struct {
	UserID            string                     `json:"user_id" validate:"required"`
	Name              string                     `json:"name" validate:"required,min=1,max=100"`
	Description       string                     `json:"description,omitempty" validate:"max=500"`
	TargetAmount      float64                    `json:"target_amount" validate:"required,gt=0"`
	Category          domain.SavingsGoalCategory `json:"category" validate:"required"`
	Priority          domain.SavingsGoalPriority `json:"priority,omitempty"`
	TargetDate        time.Time                  `json:"target_date,omitempty"`
	IsAutoSave        bool                       `json:"is_auto_save,omitempty"`
	AutoSaveAmount    float64                    `json:"auto_save_amount,omitempty" validate:"omitempty,gt=0"`
	AutoSaveFrequency string                     `json:"auto_save_frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly"`
	ImageURL          string                     `json:"image_url,omitempty" validate:"omitempty,url"`
}

// CreateSavingsGoalResponse represents the response after creating a savings goal
type CreateSavingsGoalResponse struct {
	ID                string                     `json:"id"`
	UserID            string                     `json:"user_id"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	TargetAmount      float64                    `json:"target_amount"`
	CurrentAmount     float64                    `json:"current_amount"`
	Category          domain.SavingsGoalCategory `json:"category"`
	Priority          domain.SavingsGoalPriority `json:"priority"`
	TargetDate        time.Time                  `json:"target_date"`
	Status            domain.SavingsGoalStatus   `json:"status"`
	MonthlyTarget     float64                    `json:"monthly_target"`
	WeeklyTarget      float64                    `json:"weekly_target"`
	DailyTarget       float64                    `json:"daily_target"`
	Progress          float64                    `json:"progress"`
	DaysRemaining     int                        `json:"days_remaining"`
	IsAutoSave        bool                       `json:"is_auto_save"`
	AutoSaveAmount    float64                    `json:"auto_save_amount"`
	AutoSaveFrequency string                     `json:"auto_save_frequency"`
	ImageURL          string                     `json:"image_url"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

// GetSavingsGoalRequest represents a request to get a specific savings goal
type GetSavingsGoalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// GetSavingsGoalResponse represents the response for getting a savings goal
type GetSavingsGoalResponse struct {
	ID                string                     `json:"id"`
	UserID            string                     `json:"user_id"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	TargetAmount      float64                    `json:"target_amount"`
	CurrentAmount     float64                    `json:"current_amount"`
	RemainingAmount   float64                    `json:"remaining_amount"`
	Category          domain.SavingsGoalCategory `json:"category"`
	Priority          domain.SavingsGoalPriority `json:"priority"`
	TargetDate        time.Time                  `json:"target_date"`
	Status            domain.SavingsGoalStatus   `json:"status"`
	Progress          float64                    `json:"progress"`
	MonthlyTarget     float64                    `json:"monthly_target"`
	WeeklyTarget      float64                    `json:"weekly_target"`
	DailyTarget       float64                    `json:"daily_target"`
	DaysRemaining     int                        `json:"days_remaining"`
	IsOverdue         bool                       `json:"is_overdue"`
	IsOnTrack         bool                       `json:"is_on_track"`
	IsAutoSave        bool                       `json:"is_auto_save"`
	AutoSaveAmount    float64                    `json:"auto_save_amount"`
	AutoSaveFrequency string                     `json:"auto_save_frequency"`
	ImageURL          string                     `json:"image_url"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	AchievedAt        *time.Time                 `json:"achieved_at,omitempty"`
}

// ListSavingsGoalsRequest represents a request to list savings goals
type ListSavingsGoalsRequest struct {
	UserID   string                     `json:"user_id" validate:"required"`
	Status   domain.SavingsGoalStatus   `json:"status,omitempty"`
	Category domain.SavingsGoalCategory `json:"category,omitempty"`
	Priority domain.SavingsGoalPriority `json:"priority,omitempty"`
}

// ListSavingsGoalsResponse represents the response for listing savings goals
type ListSavingsGoalsResponse struct {
	Goals   []GetSavingsGoalResponse  `json:"goals"`
	Summary domain.SavingsGoalSummary `json:"summary"`
	Count   int                       `json:"count"`
}

// UpdateSavingsGoalRequest represents a request to update a savings goal
type UpdateSavingsGoalRequest struct {
	UserID            string                      `json:"user_id" validate:"required"`
	GoalID            string                      `json:"goal_id" validate:"required"`
	Name              *string                     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description       *string                     `json:"description,omitempty" validate:"omitempty,max=500"`
	TargetAmount      *float64                    `json:"target_amount,omitempty" validate:"omitempty,gt=0"`
	Category          *domain.SavingsGoalCategory `json:"category,omitempty"`
	Priority          *domain.SavingsGoalPriority `json:"priority,omitempty"`
	TargetDate        *time.Time                  `json:"target_date,omitempty"`
	IsAutoSave        *bool                       `json:"is_auto_save,omitempty"`
	AutoSaveAmount    *float64                    `json:"auto_save_amount,omitempty" validate:"omitempty,gt=0"`
	AutoSaveFrequency *string                     `json:"auto_save_frequency,omitempty" validate:"omitempty,oneof=daily weekly monthly"`
	ImageURL          *string                     `json:"image_url,omitempty" validate:"omitempty,url"`
}

// UpdateSavingsGoalResponse represents the response after updating a savings goal
type UpdateSavingsGoalResponse struct {
	ID                string                     `json:"id"`
	UserID            string                     `json:"user_id"`
	Name              string                     `json:"name"`
	Description       string                     `json:"description"`
	TargetAmount      float64                    `json:"target_amount"`
	CurrentAmount     float64                    `json:"current_amount"`
	RemainingAmount   float64                    `json:"remaining_amount"`
	Progress          float64                    `json:"progress"`
	Category          domain.SavingsGoalCategory `json:"category"`
	Priority          domain.SavingsGoalPriority `json:"priority"`
	TargetDate        time.Time                  `json:"target_date"`
	Status            domain.SavingsGoalStatus   `json:"status"`
	MonthlyTarget     float64                    `json:"monthly_target"`
	WeeklyTarget      float64                    `json:"weekly_target"`
	DailyTarget       float64                    `json:"daily_target"`
	IsAutoSave        bool                       `json:"is_auto_save"`
	AutoSaveAmount    float64                    `json:"auto_save_amount"`
	AutoSaveFrequency string                     `json:"auto_save_frequency"`
	UpdatedAt         time.Time                  `json:"updated_at"`
}

// DeleteSavingsGoalRequest represents a request to delete a savings goal
type DeleteSavingsGoalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// AddSavingsRequest represents a request to add money to a savings goal
type AddSavingsRequest struct {
	UserID      string  `json:"user_id" validate:"required"`
	GoalID      string  `json:"goal_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description,omitempty" validate:"max=200"`
}

// AddSavingsResponse represents the response after adding savings
type AddSavingsResponse struct {
	TransactionID    string                   `json:"transaction_id"`
	GoalID           string                   `json:"goal_id"`
	Amount           float64                  `json:"amount"`
	NewCurrentAmount float64                  `json:"new_current_amount"`
	NewProgress      float64                  `json:"new_progress"`
	RemainingAmount  float64                  `json:"remaining_amount"`
	IsAchieved       bool                     `json:"is_achieved"`
	Status           domain.SavingsGoalStatus `json:"status"`
	CreatedAt        time.Time                `json:"created_at"`
}

// WithdrawSavingsRequest represents a request to withdraw money from a savings goal
type WithdrawSavingsRequest struct {
	UserID      string  `json:"user_id" validate:"required"`
	GoalID      string  `json:"goal_id" validate:"required"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description,omitempty" validate:"max=200"`
}

// WithdrawSavingsResponse represents the response after withdrawing savings
type WithdrawSavingsResponse struct {
	TransactionID    string                   `json:"transaction_id"`
	GoalID           string                   `json:"goal_id"`
	Amount           float64                  `json:"amount"`
	NewCurrentAmount float64                  `json:"new_current_amount"`
	NewProgress      float64                  `json:"new_progress"`
	RemainingAmount  float64                  `json:"remaining_amount"`
	Status           domain.SavingsGoalStatus `json:"status"`
	CreatedAt        time.Time                `json:"created_at"`
}

// PauseGoalRequest represents a request to pause a savings goal
type PauseGoalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// ResumeGoalRequest represents a request to resume a savings goal
type ResumeGoalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// CancelGoalRequest represents a request to cancel a savings goal
type CancelGoalRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
}

// GetGoalSummaryRequest represents a request to get goal summary
type GetGoalSummaryRequest struct {
	UserID   string                     `json:"user_id" validate:"required"`
	Category domain.SavingsGoalCategory `json:"category,omitempty"`
}

// GetGoalSummaryResponse represents the response for goal summary
type GetGoalSummaryResponse struct {
	Summary            domain.SavingsGoalSummary `json:"summary"`
	MonthlyTargetTotal float64                   `json:"monthly_target_total"`
	WeeklyTargetTotal  float64                   `json:"weekly_target_total"`
	DailyTargetTotal   float64                   `json:"daily_target_total"`
	NextMilestones     []GoalMilestone           `json:"next_milestones"`
	OverdueGoals       []string                  `json:"overdue_goals"`
	UpdatedAt          time.Time                 `json:"updated_at"`
}

// GoalMilestone represents a milestone for a goal
type GoalMilestone struct {
	GoalID        string    `json:"goal_id"`
	GoalName      string    `json:"goal_name"`
	Milestone     float64   `json:"milestone"` // 0.25, 0.50, 0.75, 1.0
	AmountNeeded  float64   `json:"amount_needed"`
	EstimatedDate time.Time `json:"estimated_date"`
}

// GetGoalTransactionsRequest represents a request to get goal transactions
type GetGoalTransactionsRequest struct {
	UserID string `json:"user_id" validate:"required"`
	GoalID string `json:"goal_id" validate:"required"`
	Limit  int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset,omitempty" validate:"omitempty,min=0"`
}

// GetGoalTransactionsResponse represents the response for goal transactions
type GetGoalTransactionsResponse struct {
	Transactions []domain.SavingsTransaction `json:"transactions"`
	Total        int                         `json:"total"`
	Count        int                         `json:"count"`
}
