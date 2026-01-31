package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecurringTransaction representa una transacción que se repite automáticamente
type RecurringTransaction struct {
	ID             string     `json:"id" gorm:"primaryKey"`
	UserID         string     `json:"user_id" gorm:"not null;index"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Description    string     `json:"description" gorm:"not null"`
	CategoryID     *string    `json:"category_id" gorm:"index"`
	Type           string     `json:"type" gorm:"not null"`      // "income", "expense"
	Frequency      string     `json:"frequency" gorm:"not null"` // "daily", "weekly", "monthly", "yearly"
	NextDate       time.Time  `json:"next_date" gorm:"not null;index"`
	LastExecuted   *time.Time `json:"last_executed"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	AutoCreate     bool       `json:"auto_create" gorm:"default:true"`
	NotifyBefore   int        `json:"notify_before" gorm:"default:1"` // Days before to notify
	EndDate        *time.Time `json:"end_date"`
	ExecutionCount int        `json:"execution_count" gorm:"default:0"`
	MaxExecutions  *int       `json:"max_executions"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// RecurringTransactionBuilder implements Builder pattern
type RecurringTransactionBuilder struct {
	transaction *RecurringTransaction
}

// NewRecurringTransactionBuilder creates a new builder
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

// SetCategoryID sets the category ID
func (b *RecurringTransactionBuilder) SetCategoryID(categoryID string) *RecurringTransactionBuilder {
	if categoryID != "" {
		b.transaction.CategoryID = &categoryID
	}
	return b
}

// SetType sets the transaction type
func (b *RecurringTransactionBuilder) SetType(transactionType string) *RecurringTransactionBuilder {
	b.transaction.Type = transactionType
	return b
}

// SetFrequency sets the frequency
func (b *RecurringTransactionBuilder) SetFrequency(frequency string) *RecurringTransactionBuilder {
	b.transaction.Frequency = frequency
	return b
}

// SetNextDate sets the next execution date
func (b *RecurringTransactionBuilder) SetNextDate(nextDate time.Time) *RecurringTransactionBuilder {
	b.transaction.NextDate = nextDate
	return b
}

// SetAutoCreate sets auto creation flag
func (b *RecurringTransactionBuilder) SetAutoCreate(autoCreate bool) *RecurringTransactionBuilder {
	b.transaction.AutoCreate = autoCreate
	return b
}

// SetNotifyBefore sets notification days before
func (b *RecurringTransactionBuilder) SetNotifyBefore(days int) *RecurringTransactionBuilder {
	b.transaction.NotifyBefore = days
	return b
}

// SetEndDate sets the end date
func (b *RecurringTransactionBuilder) SetEndDate(endDate *time.Time) *RecurringTransactionBuilder {
	b.transaction.EndDate = endDate
	return b
}

// SetMaxExecutions sets maximum executions
func (b *RecurringTransactionBuilder) SetMaxExecutions(max *int) *RecurringTransactionBuilder {
	b.transaction.MaxExecutions = max
	return b
}

// Build creates and validates the recurring transaction
func (b *RecurringTransactionBuilder) Build() (*RecurringTransaction, error) {
	if err := b.transaction.Validate(); err != nil {
		return nil, err
	}
	return b.transaction, nil
}

// Validate validates the recurring transaction business rules
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

	validFrequencies := []string{"daily", "weekly", "monthly", "yearly"}
	if !contains(validFrequencies, rt.Frequency) {
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

// CalculateNextDate calculates the next execution date based on frequency
func (rt *RecurringTransaction) CalculateNextDate() time.Time {
	currentNext := rt.NextDate

	switch rt.Frequency {
	case "daily":
		return currentNext.AddDate(0, 0, 1)
	case "weekly":
		return currentNext.AddDate(0, 0, 7)
	case "monthly":
		return currentNext.AddDate(0, 1, 0)
	case "yearly":
		return currentNext.AddDate(1, 0, 0)
	default:
		return currentNext
	}
}

// ShouldExecute checks if the transaction should be executed now
func (rt *RecurringTransaction) ShouldExecute() bool {
	now := time.Now()

	// Not active
	if !rt.IsActive {
		return false
	}

	// Next date not reached
	if now.Before(rt.NextDate) {
		return false
	}

	// End date passed
	if rt.EndDate != nil && now.After(*rt.EndDate) {
		return false
	}

	// Max executions reached
	if rt.MaxExecutions != nil && rt.ExecutionCount >= *rt.MaxExecutions {
		return false
	}

	return true
}

// ShouldNotify checks if a notification should be sent
func (rt *RecurringTransaction) ShouldNotify() bool {
	if !rt.IsActive {
		return false
	}

	now := time.Now()
	notifyDate := rt.NextDate.AddDate(0, 0, -rt.NotifyBefore)

	return now.After(notifyDate) && now.Before(rt.NextDate)
}

// Execute marks the transaction as executed and calculates next date
func (rt *RecurringTransaction) Execute() {
	now := time.Now()
	rt.LastExecuted = &now
	rt.ExecutionCount++
	rt.NextDate = rt.CalculateNextDate()
	rt.UpdatedAt = now

	// Check if should be deactivated
	if rt.EndDate != nil && now.After(*rt.EndDate) {
		rt.IsActive = false
	}

	if rt.MaxExecutions != nil && rt.ExecutionCount >= *rt.MaxExecutions {
		rt.IsActive = false
	}
}

// Pause pauses the recurring transaction
func (rt *RecurringTransaction) Pause() {
	rt.IsActive = false
	rt.UpdatedAt = time.Now()
}

// Resume resumes the recurring transaction
func (rt *RecurringTransaction) Resume() {
	rt.IsActive = true
	rt.UpdatedAt = time.Now()
}

// GetDaysUntilNext returns days until next execution
func (rt *RecurringTransaction) GetDaysUntilNext() int {
	now := time.Now()
	if rt.NextDate.Before(now) {
		return 0
	}

	duration := rt.NextDate.Sub(now)
	return int(duration.Hours() / 24)
}

// GetFrequencyDisplay returns user-friendly frequency text
func (rt *RecurringTransaction) GetFrequencyDisplay() string {
	switch rt.Frequency {
	case "daily":
		return "Diario"
	case "weekly":
		return "Semanal"
	case "monthly":
		return "Mensual"
	case "yearly":
		return "Anual"
	default:
		return rt.Frequency
	}
}

// GetTypeDisplay returns user-friendly type text
func (rt *RecurringTransaction) GetTypeDisplay() string {
	switch rt.Type {
	case "income":
		return "Ingreso"
	case "expense":
		return "Gasto"
	default:
		return rt.Type
	}
}

// NewRecurringTransactionID generates a new UUID for recurring transaction
func NewRecurringTransactionID() string {
	return uuid.New().String()
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
