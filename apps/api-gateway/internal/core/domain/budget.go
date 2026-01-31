// Package domain defines the core business entities and their behavior
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidBudgetAmount is returned when attempting to create a budget with invalid amount
	ErrInvalidBudgetAmount = errors.New("budget amount must be positive")
	// ErrInvalidBudgetPeriod is returned when attempting to create a budget with invalid period
	ErrInvalidBudgetPeriod = errors.New("budget period must be valid")
	// ErrBudgetAlreadyExists is returned when attempting to create a budget that already exists
	ErrBudgetAlreadyExists = errors.New("budget already exists for this category and period")
)

// BudgetPeriod represents the period type for budgets
type BudgetPeriod string

const (
	BudgetPeriodMonthly BudgetPeriod = "monthly"
	BudgetPeriodWeekly  BudgetPeriod = "weekly"
	BudgetPeriodYearly  BudgetPeriod = "yearly"
)

// BudgetStatus represents the current status of a budget
type BudgetStatus string

const (
	BudgetStatusOnTrack  BudgetStatus = "on_track" // Under 70% spent
	BudgetStatusWarning  BudgetStatus = "warning"  // 70-99% spent
	BudgetStatusExceeded BudgetStatus = "exceeded" // 100%+ spent
)

// Budget represents a spending limit for a category in a specific period
type Budget struct {
	ID          string       `json:"id"`
	UserID      string       `json:"user_id"`
	CategoryID  string       `json:"category_id"`
	Amount      float64      `json:"amount"`       // Budget limit
	SpentAmount float64      `json:"spent_amount"` // Current spent amount
	Period      BudgetPeriod `json:"period"`       // Monthly, weekly, yearly
	PeriodStart time.Time    `json:"period_start"` // Start of current period
	PeriodEnd   time.Time    `json:"period_end"`   // End of current period
	AlertAt     float64      `json:"alert_at"`     // Percentage to trigger alert (0.0-1.0)
	Status      BudgetStatus `json:"status"`       // Current status
	IsActive    bool         `json:"is_active"`    // Whether budget is active
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// BudgetBuilder implements the Builder pattern for creating budgets
type BudgetBuilder struct {
	budget *Budget
}

// NewBudgetBuilder creates a new BudgetBuilder instance
func NewBudgetBuilder() *BudgetBuilder {
	return &BudgetBuilder{
		budget: &Budget{
			ID:       "bud_" + uuid.New().String()[:8],
			AlertAt:  0.80, // Default alert at 80%
			IsActive: true,
			Status:   BudgetStatusOnTrack,
		},
	}
}

// SetUserID sets the user ID (Builder pattern)
func (b *BudgetBuilder) SetUserID(userID string) *BudgetBuilder {
	b.budget.UserID = userID
	return b
}

// SetCategoryID sets the category ID (Builder pattern)
func (b *BudgetBuilder) SetCategoryID(categoryID string) *BudgetBuilder {
	b.budget.CategoryID = categoryID
	return b
}

// SetAmount sets the budget amount (Builder pattern)
func (b *BudgetBuilder) SetAmount(amount float64) *BudgetBuilder {
	b.budget.Amount = amount
	return b
}

// SetPeriod sets the budget period (Builder pattern)
func (b *BudgetBuilder) SetPeriod(period BudgetPeriod) *BudgetBuilder {
	b.budget.Period = period
	b.setPeriodDates()
	return b
}

// SetAlertAt sets the alert threshold (Builder pattern)
func (b *BudgetBuilder) SetAlertAt(threshold float64) *BudgetBuilder {
	b.budget.AlertAt = threshold
	return b
}

// setPeriodDates sets the period start and end dates based on current time
func (b *BudgetBuilder) setPeriodDates() {
	now := time.Now()

	switch b.budget.Period {
	case BudgetPeriodMonthly:
		b.budget.PeriodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		b.budget.PeriodEnd = b.budget.PeriodStart.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	case BudgetPeriodWeekly:
		weekday := int(now.Weekday())
		b.budget.PeriodStart = now.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)
		b.budget.PeriodEnd = b.budget.PeriodStart.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	case BudgetPeriodYearly:
		b.budget.PeriodStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		b.budget.PeriodEnd = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}
}

// Build creates and validates the budget (Builder pattern)
func (b *BudgetBuilder) Build() (*Budget, error) {
	if err := b.budget.Validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	b.budget.CreatedAt = now
	b.budget.UpdatedAt = now

	return b.budget, nil
}

// Validate validates the budget business rules (Single Responsibility Principle)
func (b *Budget) Validate() error {
	if b.UserID == "" {
		return errors.New("user ID is required")
	}

	if b.CategoryID == "" {
		return errors.New("category ID is required")
	}

	if b.Amount <= 0 {
		return ErrInvalidBudgetAmount
	}

	if b.Period != BudgetPeriodMonthly && b.Period != BudgetPeriodWeekly && b.Period != BudgetPeriodYearly {
		return ErrInvalidBudgetPeriod
	}

	if b.AlertAt < 0 || b.AlertAt > 1 {
		return errors.New("alert threshold must be between 0 and 1")
	}

	return nil
}

// UpdateSpentAmount updates the spent amount and recalculates status (Single Responsibility)
func (b *Budget) UpdateSpentAmount(spentAmount float64) {
	b.SpentAmount = spentAmount
	b.updateStatus()
	b.UpdatedAt = time.Now()
}

// updateStatus updates the budget status based on spent percentage (Single Responsibility)
func (b *Budget) updateStatus() {
	percentage := b.GetSpentPercentage()

	if percentage >= 1.0 {
		b.Status = BudgetStatusExceeded
	} else if percentage >= b.AlertAt {
		b.Status = BudgetStatusWarning
	} else {
		b.Status = BudgetStatusOnTrack
	}
}

// GetSpentPercentage returns the percentage of budget spent (0.0-1.0+)
func (b *Budget) GetSpentPercentage() float64 {
	if b.Amount == 0 {
		return 0
	}
	return b.SpentAmount / b.Amount
}

// GetRemainingAmount returns the remaining budget amount
func (b *Budget) GetRemainingAmount() float64 {
	remaining := b.Amount - b.SpentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsAlertTriggered returns true if alert should be triggered
func (b *Budget) IsAlertTriggered() bool {
	return b.GetSpentPercentage() >= b.AlertAt && b.Status != BudgetStatusExceeded
}

// IsInCurrentPeriod checks if the budget is in the current period
func (b *Budget) IsInCurrentPeriod() bool {
	now := time.Now()
	return now.After(b.PeriodStart) && now.Before(b.PeriodEnd)
}

// ResetForNewPeriod resets the budget for a new period (Strategy pattern ready)
func (b *Budget) ResetForNewPeriod() {
	b.SpentAmount = 0
	b.Status = BudgetStatusOnTrack
	b.setPeriodDatesFromCurrent()
	b.UpdatedAt = time.Now()
}

// setPeriodDatesFromCurrent updates period dates based on current time
func (b *Budget) setPeriodDatesFromCurrent() {
	now := time.Now()

	switch b.Period {
	case BudgetPeriodMonthly:
		b.PeriodStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		b.PeriodEnd = b.PeriodStart.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	case BudgetPeriodWeekly:
		weekday := int(now.Weekday())
		b.PeriodStart = now.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)
		b.PeriodEnd = b.PeriodStart.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	case BudgetPeriodYearly:
		b.PeriodStart = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		b.PeriodEnd = time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())
	}
}

// BudgetSummary provides aggregated budget information
type BudgetSummary struct {
	TotalBudgets   int     `json:"total_budgets"`
	TotalAllocated float64 `json:"total_allocated"`
	TotalSpent     float64 `json:"total_spent"`
	OnTrackCount   int     `json:"on_track_count"`
	WarningCount   int     `json:"warning_count"`
	ExceededCount  int     `json:"exceeded_count"`
	AverageUsage   float64 `json:"average_usage"`
}

// NewBudgetID generates a new budget ID
func NewBudgetID() string {
	return "bud_" + uuid.New().String()[:8]
}
