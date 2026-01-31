package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// BaseModelBuilder es un builder para facilitar la creación de modelos base en pruebas
type BaseModelBuilder struct {
	model *BaseModel
}

func NewBaseModelBuilder() *BaseModelBuilder {
	return &BaseModelBuilder{
		model: &BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *BaseModelBuilder) WithID(id uint) *BaseModelBuilder {
	b.model.ID = id
	return b
}

func (b *BaseModelBuilder) WithCreatedAt(createdAt time.Time) *BaseModelBuilder {
	b.model.CreatedAt = createdAt
	return b
}

func (b *BaseModelBuilder) WithUpdatedAt(updatedAt time.Time) *BaseModelBuilder {
	b.model.UpdatedAt = updatedAt
	return b
}

func (b *BaseModelBuilder) Build() *BaseModel {
	return b.model
}

// FinancialDataBuilder es un builder para facilitar la creación de datos financieros en pruebas
type FinancialDataBuilder struct {
	data *FinancialData
}

func NewFinancialDataBuilder() *FinancialDataBuilder {
	return &FinancialDataBuilder{
		data: &FinancialData{
			BaseModel: BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}

func (b *FinancialDataBuilder) WithID(id uint) *FinancialDataBuilder {
	b.data.ID = id
	return b
}

func (b *FinancialDataBuilder) WithUserID(userID uint) *FinancialDataBuilder {
	b.data.UserID = userID
	return b
}

func (b *FinancialDataBuilder) WithAmount(amount float64) *FinancialDataBuilder {
	b.data.Amount = amount
	return b
}

func (b *FinancialDataBuilder) WithDescription(description string) *FinancialDataBuilder {
	b.data.Description = description
	return b
}

func (b *FinancialDataBuilder) WithCategory(category string) *FinancialDataBuilder {
	b.data.Category = category
	return b
}

func (b *FinancialDataBuilder) WithDate(date string) *FinancialDataBuilder {
	b.data.Date = date
	return b
}

func (b *FinancialDataBuilder) Build() *FinancialData {
	return b.data
}

func TestBaseModelBuilder(t *testing.T) {
	now := time.Now()
	model := NewBaseModelBuilder().
		WithID(1).
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()

	assert.Equal(t, uint(1), model.ID)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
}

func TestFinancialDataBuilder(t *testing.T) {
	data := NewFinancialDataBuilder().
		WithID(1).
		WithUserID(1).
		WithAmount(1000.50).
		WithDescription("Test transaction").
		WithCategory("Test category").
		WithDate("2024-03-25").
		Build()

	assert.Equal(t, uint(1), data.ID)
	assert.Equal(t, uint(1), data.UserID)
	assert.Equal(t, 1000.50, data.Amount)
	assert.Equal(t, "Test transaction", data.Description)
	assert.Equal(t, "Test category", data.Category)
	assert.Equal(t, "2024-03-25", data.Date)
	assert.False(t, data.CreatedAt.IsZero())
	assert.False(t, data.UpdatedAt.IsZero())
}
