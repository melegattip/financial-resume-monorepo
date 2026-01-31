package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpenseBuilder(t *testing.T) {
	dueDate := time.Now().Add(24 * time.Hour)
	tests := []struct {
		name        string
		id          string
		userID      string
		amount      float64
		description string
		category    string
		paid        bool
		dueDate     time.Time
	}{
		{
			name:        "Valid expense",
			id:          "123",
			userID:      "user123",
			amount:      500.75,
			description: "Grocery shopping",
			category:    "Food",
			paid:        true,
			dueDate:     dueDate,
		},
		{
			name:        "Unpaid expense",
			id:          "456",
			userID:      "user456",
			amount:      1200.00,
			description: "Rent",
			category:    "Housing",
			paid:        false,
			dueDate:     dueDate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expense := NewExpenseBuilder().
				SetID(tt.id).
				SetUserID(tt.userID).
				SetAmount(tt.amount).
				SetDescription(tt.description).
				SetCategoryID(tt.category).
				SetPaid(tt.paid).
				SetDueDate(tt.dueDate).
				Build()

			assert.Equal(t, tt.id, expense.ID)
			assert.Equal(t, tt.userID, expense.UserID)
			assert.Equal(t, tt.amount, expense.Amount)
			assert.Equal(t, tt.description, expense.Description)
			assert.Equal(t, tt.category, expense.GetCategoryID())
			assert.Equal(t, tt.paid, expense.Paid)
			assert.Equal(t, tt.dueDate, expense.DueDate)
			assert.False(t, expense.CreatedAt.IsZero())
			assert.False(t, expense.UpdatedAt.IsZero())
		})
	}
}

func TestExpenseBuilder_CustomDates(t *testing.T) {
	now := time.Now()
	dueDate := now.Add(7 * 24 * time.Hour)
	expense := NewExpenseBuilder().
		SetID("789").
		SetUserID("user789").
		SetAmount(350.25).
		SetDescription("Internet bill").
		SetCategoryID("Utilities").
		SetPaid(false).
		SetDueDate(dueDate).
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Build()

	assert.Equal(t, "789", expense.ID)
	assert.Equal(t, "user789", expense.UserID)
	assert.Equal(t, 350.25, expense.Amount)
	assert.Equal(t, "Internet bill", expense.Description)
	assert.Equal(t, "Utilities", expense.GetCategoryID())
	assert.False(t, expense.Paid)
	assert.Equal(t, dueDate, expense.DueDate)
	assert.Equal(t, now, expense.CreatedAt)
	assert.Equal(t, now, expense.UpdatedAt)
}

func TestExpense_CalculatePercentage(t *testing.T) {
	expense := &Expense{
		Amount: 50.0,
	}

	// Test con ingreso total positivo
	expense.CalculatePercentage(200.0)
	assert.Equal(t, 25.0, expense.GetPercentage(), "El porcentaje debería ser 25% (50/200 * 100)")

	// Test con ingreso total cero
	expense.CalculatePercentage(0)
	assert.Equal(t, 0.0, expense.GetPercentage(), "El porcentaje debería ser 0% cuando el ingreso total es 0")

	// Test con ingreso total negativo
	expense.CalculatePercentage(-100.0)
	assert.Equal(t, 0.0, expense.GetPercentage(), "El porcentaje debería ser 0% cuando el ingreso total es negativo")
}

func TestExpense_GetPercentage(t *testing.T) {
	expense := &Expense{
		Percentage: 30.5,
	}

	assert.Equal(t, 30.5, expense.GetPercentage(), "GetPercentage debería retornar el valor almacenado")
}

func TestExpenseFactory(t *testing.T) {
	factory := NewExpenseFactory()
	transaction := factory.CreateTransaction()

	expense, ok := transaction.(*Expense)
	assert.True(t, ok, "Expected transaction to be of type *Expense")
	assert.NotNil(t, expense)
	assert.False(t, expense.CreatedAt.IsZero())
	assert.False(t, expense.UpdatedAt.IsZero())
	assert.Equal(t, 0.0, expense.GetPercentage(), "El porcentaje inicial debería ser 0")
}
