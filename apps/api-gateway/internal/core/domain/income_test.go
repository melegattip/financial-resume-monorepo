package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIncomeBuilder(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		userID      string
		amount      float64
		description string
		category    string
		source      string
	}{
		{
			name:        "Valid income",
			id:          "123",
			userID:      "user123",
			amount:      5000.75,
			description: "Salary",
			category:    "Work",
			source:      "Company XYZ",
		},
		{
			name:        "Freelance income",
			id:          "456",
			userID:      "user456",
			amount:      1200.00,
			description: "Freelance work",
			category:    "Freelance",
			source:      "Client ABC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			income := NewIncomeBuilder().
				SetID(tt.id).
				SetUserID(tt.userID).
				SetAmount(tt.amount).
				SetDescription(tt.description).
				SetCategoryID(tt.category).
				Build()

			assert.Equal(t, tt.id, income.ID)
			assert.Equal(t, tt.userID, income.UserID)
			assert.Equal(t, tt.amount, income.Amount)
			assert.Equal(t, tt.description, income.Description)
			assert.Equal(t, tt.category, income.GetCategoryID())
			assert.False(t, income.CreatedAt.IsZero())
			assert.False(t, income.UpdatedAt.IsZero())
		})
	}
}

func TestIncomeBuilder_CustomDates(t *testing.T) {
	now := time.Now()
	income := NewIncomeBuilder().
		SetID("789").
		SetUserID("user789").
		SetAmount(3500.25).
		SetDescription("Bonus").
		SetCategoryID("Work").
		SetCreatedAt(now).
		SetUpdatedAt(now).
		Build()

	assert.Equal(t, "789", income.ID)
	assert.Equal(t, "user789", income.UserID)
	assert.Equal(t, 3500.25, income.Amount)
	assert.Equal(t, "Bonus", income.Description)
	assert.Equal(t, "Work", income.GetCategoryID())
	assert.Equal(t, now, income.CreatedAt)
	assert.Equal(t, now, income.UpdatedAt)
}

func TestIncome_CalculatePercentage(t *testing.T) {
	income := &Income{
		Amount: 100.0,
	}

	// Los ingresos siempre deberían tener porcentaje 0
	income.CalculatePercentage(200.0)
	assert.Equal(t, 0.0, income.GetPercentage(), "El porcentaje de un ingreso siempre debería ser 0")

	income.CalculatePercentage(0)
	assert.Equal(t, 0.0, income.GetPercentage(), "El porcentaje de un ingreso siempre debería ser 0")

	income.CalculatePercentage(-100.0)
	assert.Equal(t, 0.0, income.GetPercentage(), "El porcentaje de un ingreso siempre debería ser 0")
}

func TestIncome_GetPercentage(t *testing.T) {
	income := &Income{
		Percentage: 0.0,
	}

	assert.Equal(t, 0.0, income.GetPercentage(), "GetPercentage debería retornar siempre 0 para ingresos")
}

func TestIncomeFactory(t *testing.T) {
	factory := NewIncomeFactory()
	transaction := factory.CreateTransaction()

	income, ok := transaction.(*Income)
	assert.True(t, ok, "Expected transaction to be of type *Income")
	assert.NotNil(t, income)
	assert.False(t, income.CreatedAt.IsZero())
	assert.False(t, income.UpdatedAt.IsZero())
	assert.Equal(t, 0.0, income.GetPercentage(), "El porcentaje inicial debería ser 0")
}
