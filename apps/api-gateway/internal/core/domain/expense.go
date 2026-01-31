// Package domain defines the core business entities and their behavior
package domain

import "time"

// Expense represents a financial expense transaction
type Expense struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`      // Monto total original
	AmountPaid  float64   `json:"amount_paid"` // Monto ya pagado
	Description string    `json:"description"`
	CategoryID  *string   `json:"category_id" gorm:"default:null"`
	Paid        bool      `json:"paid"`
	DueDate     time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Percentage  float64   `json:"percentage"`
}

// GetID returns the expense ID
func (e *Expense) GetID() string {
	return e.ID
}

// GetUserID returns the user ID
func (e *Expense) GetUserID() string {
	return e.UserID
}

// GetAmount returns the expense amount
func (e *Expense) GetAmount() float64 {
	return e.Amount
}

// GetDescription returns the expense description
func (e *Expense) GetDescription() string {
	return e.Description
}

// GetCategoryID returns the category ID
func (e *Expense) GetCategoryID() string {
	if e.CategoryID == nil {
		return ""
	}
	return *e.CategoryID
}

// GetCreatedAt returns the creation timestamp
func (e *Expense) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// GetUpdatedAt returns the last update timestamp
func (e *Expense) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

// GetPercentage retorna el porcentaje calculado del gasto
func (e *Expense) GetPercentage() float64 {
	return e.Percentage
}

// CalculatePercentage calcula el porcentaje que representa este gasto sobre el total de ingresos
func (e *Expense) CalculatePercentage(totalIncome float64) {
	if totalIncome <= 0 {
		e.Percentage = 0
		return
	}
	e.Percentage = (e.Amount / totalIncome) * 100
}

// GetPendingAmount retorna el monto pendiente de pago
func (e *Expense) GetPendingAmount() float64 {
	pending := e.Amount - e.AmountPaid
	if pending < 0 {
		return 0
	}
	return pending
}

// IsFullyPaid verifica si el gasto está completamente pagado
func (e *Expense) IsFullyPaid() bool {
	return e.AmountPaid >= e.Amount
}

// AddPayment agrega un pago al gasto y actualiza el estado
func (e *Expense) AddPayment(paymentAmount float64) {
	e.AmountPaid += paymentAmount
	if e.AmountPaid >= e.Amount {
		e.AmountPaid = e.Amount // No permitir sobrepagos
		e.Paid = true
	}
}

// ExpenseFactory implements TransactionFactory for creating expense transactions
type ExpenseFactory struct{}

// NewExpenseFactory creates a new ExpenseFactory instance
func NewExpenseFactory() *ExpenseFactory {
	return &ExpenseFactory{}
}

// CreateTransaction creates a new expense transaction with default values
func (f *ExpenseFactory) CreateTransaction() Transaction {
	return &Expense{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ExpenseBuilder implements the builder pattern for creating Expense instances
type ExpenseBuilder struct {
	expense *Expense
}

// NewExpenseBuilder creates a new ExpenseBuilder instance with default values
func NewExpenseBuilder() *ExpenseBuilder {
	return &ExpenseBuilder{
		expense: &Expense{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// SetID sets the expense ID and returns the builder
func (b *ExpenseBuilder) SetID(id string) *ExpenseBuilder {
	b.expense.ID = id
	return b
}

// SetUserID sets the user ID and returns the builder
func (b *ExpenseBuilder) SetUserID(userID string) *ExpenseBuilder {
	b.expense.UserID = userID
	return b
}

// SetAmount sets the expense amount and returns the builder
func (b *ExpenseBuilder) SetAmount(amount float64) *ExpenseBuilder {
	b.expense.Amount = amount
	return b
}

// SetAmountPaid sets the amount paid and returns the builder
func (b *ExpenseBuilder) SetAmountPaid(amountPaid float64) *ExpenseBuilder {
	b.expense.AmountPaid = amountPaid
	return b
}

// SetPercentage sets the expense percentage and returns the builder
func (b *ExpenseBuilder) SetPercentage(percentage float64) *ExpenseBuilder {
	b.expense.Percentage = percentage
	return b
}

// SetDescription sets the expense description and returns the builder
func (b *ExpenseBuilder) SetDescription(description string) *ExpenseBuilder {
	b.expense.Description = description
	return b
}

// SetCategoryID sets the expense category ID and returns the builder
func (b *ExpenseBuilder) SetCategoryID(categoryID string) *ExpenseBuilder {
	if categoryID == "" {
		b.expense.CategoryID = nil
	} else {
		b.expense.CategoryID = &categoryID
	}
	return b
}

// SetPaid sets the paid status and returns the builder
func (b *ExpenseBuilder) SetPaid(paid bool) *ExpenseBuilder {
	b.expense.Paid = paid
	return b
}

// SetDueDate sets the due date and returns the builder
func (b *ExpenseBuilder) SetDueDate(dueDate time.Time) *ExpenseBuilder {
	b.expense.DueDate = dueDate
	return b
}

// SetCreatedAt sets the creation timestamp and returns the builder
func (b *ExpenseBuilder) SetCreatedAt(createdAt time.Time) *ExpenseBuilder {
	b.expense.CreatedAt = createdAt
	return b
}

// SetUpdatedAt sets the last update timestamp and returns the builder
func (b *ExpenseBuilder) SetUpdatedAt(updatedAt time.Time) *ExpenseBuilder {
	b.expense.UpdatedAt = updatedAt
	return b
}

// Build creates and returns a new Expense instance
func (b *ExpenseBuilder) Build() *Expense {
	return b.expense
}
