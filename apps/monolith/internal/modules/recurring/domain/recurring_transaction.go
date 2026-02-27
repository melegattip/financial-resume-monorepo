package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecurringTransaction represents a transaction that repeats automatically
type RecurringTransaction struct {
	ID             string
	UserID         string
	TenantID       string
	Amount         float64
	Description    string
	CategoryID     *string
	Type           string     // "income", "expense"
	Frequency      string     // "daily", "weekly", "monthly", "yearly"
	NextDate       time.Time
	LastExecuted   *time.Time
	IsActive       bool
	AutoCreate     bool
	NotifyBefore   int        // days before next date to notify
	EndDate        *time.Time
	ExecutionCount int
	MaxExecutions  *int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

// RecurringTransactionBuilder implements the Builder pattern
type RecurringTransactionBuilder struct {
	transaction *RecurringTransaction
}

// NewRecurringTransactionBuilder creates a new builder with sensible defaults
func NewRecurringTransactionBuilder() *RecurringTransactionBuilder {
	return &RecurringTransactionBuilder{
		transaction: &RecurringTransaction{
			IsActive:     true,
			AutoCreate:   true,
			NotifyBefore: 1,
		},
	}
}

// SetID sets the transaction ID
func (b *RecurringTransactionBuilder) SetID(id string) *RecurringTransactionBuilder {
	b.transaction.ID = id
	return b
}

// SetUserID sets the user ID
func (b *RecurringTransactionBuilder) SetUserID(userID string) *RecurringTransactionBuilder {
	b.transaction.UserID = userID
	return b
}

// SetAmount sets the amount
func (b *RecurringTransactionBuilder) SetAmount(amount float64) *RecurringTransactionBuilder {
	b.transaction.Amount = amount
	return b
}

// SetDescription sets the description
func (b *RecurringTransactionBuilder) SetDescription(description string) *RecurringTransactionBuilder {
	b.transaction.Description = description
	return b
}

// SetCategoryID sets the optional category ID
func (b *RecurringTransactionBuilder) SetCategoryID(categoryID string) *RecurringTransactionBuilder {
	if categoryID != "" {
		b.transaction.CategoryID = &categoryID
	}
	return b
}

// SetType sets the transaction type ("income" or "expense")
func (b *RecurringTransactionBuilder) SetType(transactionType string) *RecurringTransactionBuilder {
	b.transaction.Type = transactionType
	return b
}

// SetFrequency sets the frequency ("daily", "weekly", "monthly", "yearly")
func (b *RecurringTransactionBuilder) SetFrequency(frequency string) *RecurringTransactionBuilder {
	b.transaction.Frequency = frequency
	return b
}

// SetNextDate sets the next execution date
func (b *RecurringTransactionBuilder) SetNextDate(nextDate time.Time) *RecurringTransactionBuilder {
	b.transaction.NextDate = nextDate
	return b
}

// SetAutoCreate sets whether to automatically create the transaction on execution
func (b *RecurringTransactionBuilder) SetAutoCreate(autoCreate bool) *RecurringTransactionBuilder {
	b.transaction.AutoCreate = autoCreate
	return b
}

// SetNotifyBefore sets the number of days before next date to send a notification
func (b *RecurringTransactionBuilder) SetNotifyBefore(days int) *RecurringTransactionBuilder {
	b.transaction.NotifyBefore = days
	return b
}

// SetEndDate sets the optional end date
func (b *RecurringTransactionBuilder) SetEndDate(endDate *time.Time) *RecurringTransactionBuilder {
	b.transaction.EndDate = endDate
	return b
}

// SetMaxExecutions sets the optional maximum number of executions
func (b *RecurringTransactionBuilder) SetMaxExecutions(max *int) *RecurringTransactionBuilder {
	b.transaction.MaxExecutions = max
	return b
}

// Build validates and returns the built RecurringTransaction
func (b *RecurringTransactionBuilder) Build() (*RecurringTransaction, error) {
	if err := b.transaction.Validate(); err != nil {
		return nil, err
	}
	return b.transaction, nil
}

// validFrequencies is the list of accepted frequency values
var validFrequencies = []string{"daily", "weekly", "monthly", "yearly"}

// Validate validates the business rules for a recurring transaction
func (rt *RecurringTransaction) Validate() error {
	if rt.UserID == "" {
		return errors.New("user ID is required")
	}

	if rt.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if strings.TrimSpace(rt.Description) == "" {
		return errors.New("description is required")
	}

	if rt.Type != "income" && rt.Type != "expense" {
		return errors.New("type must be 'income' or 'expense'")
	}

	if !containsFrequency(validFrequencies, rt.Frequency) {
		return errors.New("frequency must be one of: daily, weekly, monthly, yearly")
	}

	if rt.NextDate.IsZero() {
		return errors.New("next date is required")
	}

	if rt.EndDate != nil && rt.EndDate.Before(rt.NextDate) {
		return errors.New("end date cannot be before next date")
	}

	if rt.MaxExecutions != nil && *rt.MaxExecutions <= 0 {
		return errors.New("max executions must be greater than 0")
	}

	if rt.NotifyBefore < 0 {
		return errors.New("notify before days cannot be negative")
	}

	return nil
}

// CalculateNextDate calculates the next execution date based on the current NextDate and Frequency
func (rt *RecurringTransaction) CalculateNextDate() time.Time {
	current := rt.NextDate

	switch rt.Frequency {
	case "daily":
		return current.AddDate(0, 0, 1)
	case "weekly":
		return current.AddDate(0, 0, 7)
	case "monthly":
		return current.AddDate(0, 1, 0)
	case "yearly":
		return current.AddDate(1, 0, 0)
	default:
		return current
	}
}

// ShouldExecute returns true if the recurring transaction is due for execution right now
func (rt *RecurringTransaction) ShouldExecute() bool {
	if !rt.IsActive {
		return false
	}

	now := time.Now()

	if now.Before(rt.NextDate) {
		return false
	}

	if rt.EndDate != nil && now.After(*rt.EndDate) {
		return false
	}

	if rt.MaxExecutions != nil && rt.ExecutionCount >= *rt.MaxExecutions {
		return false
	}

	return true
}

// Execute marks the transaction as executed, increments the counter,
// advances NextDate, and deactivates the transaction when limits are reached
func (rt *RecurringTransaction) Execute() {
	now := time.Now()
	rt.LastExecuted = &now
	rt.ExecutionCount++
	rt.NextDate = rt.CalculateNextDate()
	rt.UpdatedAt = now

	if rt.EndDate != nil && now.After(*rt.EndDate) {
		rt.IsActive = false
	}

	if rt.MaxExecutions != nil && rt.ExecutionCount >= *rt.MaxExecutions {
		rt.IsActive = false
	}
}

// Pause deactivates the recurring transaction
func (rt *RecurringTransaction) Pause() {
	rt.IsActive = false
	rt.UpdatedAt = time.Now()
}

// Resume reactivates the recurring transaction
func (rt *RecurringTransaction) Resume() {
	rt.IsActive = true
	rt.UpdatedAt = time.Now()
}

// GetDaysUntilNext returns the number of days until the next execution date.
// Returns 0 if the next date is in the past.
func (rt *RecurringTransaction) GetDaysUntilNext() int {
	now := time.Now()
	if rt.NextDate.Before(now) {
		return 0
	}

	duration := rt.NextDate.Sub(now)
	return int(duration.Hours() / 24)
}

// ShouldNotify returns true if a notification should be sent based on NotifyBefore setting
func (rt *RecurringTransaction) ShouldNotify() bool {
	if !rt.IsActive {
		return false
	}

	now := time.Now()
	notifyDate := rt.NextDate.AddDate(0, 0, -rt.NotifyBefore)

	return now.After(notifyDate) && now.Before(rt.NextDate)
}

// NewRecurringTransactionID generates a new UUID for a recurring transaction
func NewRecurringTransactionID() string {
	return uuid.New().String()
}

// containsFrequency checks whether a slice contains the given string
func containsFrequency(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
