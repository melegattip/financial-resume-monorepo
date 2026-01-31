package transactions

import (
	"testing"

	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
	"github.com/stretchr/testify/assert"
)

func TestToTransaction(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected Transaction
	}{
		{
			name: "CreateIncomeResponse",
			input: &incomes.CreateIncomeResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
		},
		{
			name: "GetIncomeResponse",
			input: &incomes.GetIncomeResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
		},
		{
			name: "UpdateIncomeResponse",
			input: &incomes.UpdateIncomeResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test income",
			},
		},
		{
			name: "CreateExpenseResponse",
			input: &expenses.CreateExpenseResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
		},
		{
			name: "GetExpenseResponse",
			input: &expenses.GetExpenseResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
		},
		{
			name: "UpdateExpenseResponse",
			input: &expenses.UpdateExpenseResponse{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
			expected: &TransactionModel{
				ID:          "1",
				UserID:      "user1",
				Amount:      100.0,
				Description: "Test expense",
				CategoryID:  "food",
			},
		},
		{
			name:     "Invalid type",
			input:    "invalid",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTransaction(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToTransactionSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
	}{
		{
			name: "Income slice",
			input: []incomes.GetIncomeResponse{
				{
					ID:          "1",
					UserID:      "user1",
					Amount:      100.0,
					Description: "Test income 1",
				},
				{
					ID:          "2",
					UserID:      "user1",
					Amount:      200.0,
					Description: "Test income 2",
				},
			},
			expected: 2,
		},
		{
			name: "Expense slice",
			input: []expenses.GetExpenseResponse{
				{
					ID:          "1",
					UserID:      "user1",
					Amount:      100.0,
					Description: "Test expense 1",
					CategoryID:  "food",
				},
				{
					ID:          "2",
					UserID:      "user1",
					Amount:      200.0,
					Description: "Test expense 2",
					CategoryID:  "transport",
				},
			},
			expected: 2,
		},
		{
			name:     "Invalid type",
			input:    []string{"invalid"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toTransactionSlice(tt.input)
			assert.Len(t, result, tt.expected)
		})
	}
}
