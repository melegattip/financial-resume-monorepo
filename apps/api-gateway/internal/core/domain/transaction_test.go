package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestTransactionBuilder es un builder para facilitar la creación de transacciones en pruebas
type TestTransactionBuilder struct {
	transaction *BaseTransaction
}

func NewTestTransactionBuilder() *TestTransactionBuilder {
	return &TestTransactionBuilder{
		transaction: &BaseTransaction{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func (b *TestTransactionBuilder) WithID(id string) *TestTransactionBuilder {
	b.transaction.ID = id
	return b
}

func (b *TestTransactionBuilder) WithUserID(userID string) *TestTransactionBuilder {
	b.transaction.UserID = userID
	return b
}

func (b *TestTransactionBuilder) WithAmount(amount float64) *TestTransactionBuilder {
	b.transaction.Amount = amount
	return b
}

func (b *TestTransactionBuilder) WithDescription(description string) *TestTransactionBuilder {
	b.transaction.Description = description
	return b
}

func (b *TestTransactionBuilder) WithCategoryID(categoryID string) *TestTransactionBuilder {
	b.transaction.CategoryID = categoryID
	return b
}

func (b *TestTransactionBuilder) WithCreatedAt(createdAt time.Time) *TestTransactionBuilder {
	b.transaction.CreatedAt = createdAt
	return b
}

func (b *TestTransactionBuilder) WithUpdatedAt(updatedAt time.Time) *TestTransactionBuilder {
	b.transaction.UpdatedAt = updatedAt
	return b
}

func (b *TestTransactionBuilder) Build() *BaseTransaction {
	return b.transaction
}

func TestBaseTransaction_GetterMethods(t *testing.T) {
	now := time.Now()
	transaction := NewTestTransactionBuilder().
		WithID("123").
		WithUserID("user123").
		WithAmount(100.50).
		WithDescription("Test transaction").
		WithCategoryID("Test category").
		WithCreatedAt(now).
		WithUpdatedAt(now).
		Build()

	assert.Equal(t, "123", transaction.GetID())
	assert.Equal(t, "user123", transaction.GetUserID())
	assert.Equal(t, 100.50, transaction.GetAmount())
	assert.Equal(t, "Test transaction", transaction.GetDescription())
	assert.Equal(t, "Test category", transaction.GetCategoryID())
	assert.Equal(t, now, transaction.GetCreatedAt())
	assert.Equal(t, now, transaction.GetUpdatedAt())
}

func TestTransactionType_Constants(t *testing.T) {
	assert.Equal(t, TransactionType("income"), IncomeType)
	assert.Equal(t, TransactionType("expense"), ExpenseType)
}
