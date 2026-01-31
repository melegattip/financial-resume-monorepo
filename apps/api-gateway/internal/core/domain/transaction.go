package domain

import "time"

// Transaction es la interfaz base para todos los tipos de transacciones
type Transaction interface {
	GetID() string
	GetUserID() string
	GetAmount() float64
	GetDescription() string
	GetCategoryID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetPercentage() float64
	CalculatePercentage(totalIncome float64)
}

// BaseTransaction contiene los campos comunes para todos los tipos de transacciones
type BaseTransaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Percentage  float64   `json:"percentage"`
}

// Implementación de los métodos de la interfaz Transaction
func (b *BaseTransaction) GetID() string           { return b.ID }
func (b *BaseTransaction) GetUserID() string       { return b.UserID }
func (b *BaseTransaction) GetAmount() float64      { return b.Amount }
func (b *BaseTransaction) GetDescription() string  { return b.Description }
func (b *BaseTransaction) GetCategoryID() string   { return b.CategoryID }
func (b *BaseTransaction) GetCreatedAt() time.Time { return b.CreatedAt }
func (b *BaseTransaction) GetUpdatedAt() time.Time { return b.UpdatedAt }
func (b *BaseTransaction) GetPercentage() float64 {
	return b.Percentage
}

func (b *BaseTransaction) CalculatePercentage(totalIncome float64) {
	if totalIncome > 0 {
		b.Percentage = (b.Amount / totalIncome) * 100
	} else {
		b.Percentage = 0
	}
}

// TransactionType representa el tipo de transacción
type TransactionType string

const (
	IncomeType  TransactionType = "income"
	ExpenseType TransactionType = "expense"
)

// TransactionFactory es la interfaz para crear diferentes tipos de transacciones
type TransactionFactory interface {
	CreateTransaction() Transaction
}
