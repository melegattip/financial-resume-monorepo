// Package domain defines the core business entities and their behavior
package domain

import "time"

// Income represents a financial income transaction
type Income struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CategoryID  *string   `json:"category_id" gorm:"default:null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Percentage  float64   `json:"percentage"`
}

// GetID returns the income ID
func (i *Income) GetID() string {
	return i.ID
}

// GetUserID returns the user ID
func (i *Income) GetUserID() string {
	return i.UserID
}

// GetAmount returns the income amount
func (i *Income) GetAmount() float64 {
	return i.Amount
}

// GetDescription returns the income description
func (i *Income) GetDescription() string {
	return i.Description
}

// GetCategoryID returns the category ID
func (i *Income) GetCategoryID() string {
	if i.CategoryID == nil {
		return ""
	}
	return *i.CategoryID
}

// GetCreatedAt returns the creation timestamp
func (i *Income) GetCreatedAt() time.Time {
	return i.CreatedAt
}

// GetUpdatedAt returns the last update timestamp
func (i *Income) GetUpdatedAt() time.Time {
	return i.UpdatedAt
}

// GetPercentage retorna el porcentaje calculado del ingreso
func (i *Income) GetPercentage() float64 {
	return i.Percentage
}

// CalculatePercentage calcula el porcentaje que representa este ingreso sobre el total de ingresos
func (i *Income) CalculatePercentage(totalIncome float64) {
	// Los ingresos no necesitan calcular porcentaje
	i.Percentage = 0
}

// IncomeFactory implements TransactionFactory for creating income transactions
type IncomeFactory struct{}

// NewIncomeFactory creates a new IncomeFactory instance
func NewIncomeFactory() *IncomeFactory {
	return &IncomeFactory{}
}

// CreateTransaction creates a new income transaction with default values
func (f *IncomeFactory) CreateTransaction() Transaction {
	return &Income{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IncomeBuilder implements the builder pattern for creating Income instances
type IncomeBuilder struct {
	income *Income
}

// NewIncomeBuilder creates a new IncomeBuilder instance with default values
func NewIncomeBuilder() *IncomeBuilder {
	return &IncomeBuilder{
		income: &Income{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

// SetID sets the income ID and returns the builder
func (b *IncomeBuilder) SetID(id string) *IncomeBuilder {
	b.income.ID = id
	return b
}

// SetUserID sets the user ID and returns the builder
func (b *IncomeBuilder) SetUserID(userID string) *IncomeBuilder {
	b.income.UserID = userID
	return b
}

// SetAmount sets the income amount and returns the builder
func (b *IncomeBuilder) SetAmount(amount float64) *IncomeBuilder {
	b.income.Amount = amount
	return b
}

// SetDescription sets the income description and returns the builder
func (b *IncomeBuilder) SetDescription(description string) *IncomeBuilder {
	b.income.Description = description
	return b
}

// SetCategoryID sets the income category ID and returns the builder
func (b *IncomeBuilder) SetCategoryID(categoryID string) *IncomeBuilder {
	if categoryID == "" {
		b.income.CategoryID = nil
	} else {
		b.income.CategoryID = &categoryID
	}
	return b
}

// SetCreatedAt sets the creation timestamp and returns the builder
func (b *IncomeBuilder) SetCreatedAt(createdAt time.Time) *IncomeBuilder {
	b.income.CreatedAt = createdAt
	return b
}

// SetUpdatedAt sets the last update timestamp and returns the builder
func (b *IncomeBuilder) SetUpdatedAt(updatedAt time.Time) *IncomeBuilder {
	b.income.UpdatedAt = updatedAt
	return b
}

// Build creates and returns a new Income instance
func (b *IncomeBuilder) Build() *Income {
	return b.income
}
