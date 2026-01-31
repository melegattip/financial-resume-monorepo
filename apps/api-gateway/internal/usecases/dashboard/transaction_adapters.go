package dashboard

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// ExpenseTransactionAdapter adapta un gasto a la interfaz Transaction
type ExpenseTransactionAdapter struct {
	expense *domain.Expense
}

// NewExpenseTransaction crea un adaptador para gastos
func NewExpenseTransaction(expense *domain.Expense) usecases.Transaction {
	return &ExpenseTransactionAdapter{expense: expense}
}

func (e *ExpenseTransactionAdapter) GetID() string {
	return e.expense.ID
}

func (e *ExpenseTransactionAdapter) GetUserID() string {
	return e.expense.UserID
}

func (e *ExpenseTransactionAdapter) GetAmount() float64 {
	return e.expense.Amount
}

func (e *ExpenseTransactionAdapter) GetCreatedAt() time.Time {
	return e.expense.CreatedAt
}

func (e *ExpenseTransactionAdapter) GetType() usecases.TransactionType {
	return usecases.ExpenseTransaction
}

func (e *ExpenseTransactionAdapter) GetCategoryID() string {
	if e.expense.CategoryID == nil {
		return ""
	}
	return *e.expense.CategoryID
}

func (e *ExpenseTransactionAdapter) IsPending() bool {
	return !e.expense.Paid
}

// IncomeTransactionAdapter adapta un ingreso a la interfaz Transaction
type IncomeTransactionAdapter struct {
	income *domain.Income
}

// NewIncomeTransaction crea un adaptador para ingresos
func NewIncomeTransaction(income *domain.Income) usecases.Transaction {
	return &IncomeTransactionAdapter{income: income}
}

func (i *IncomeTransactionAdapter) GetID() string {
	return i.income.ID
}

func (i *IncomeTransactionAdapter) GetUserID() string {
	return i.income.UserID
}

func (i *IncomeTransactionAdapter) GetAmount() float64 {
	return i.income.Amount
}

func (i *IncomeTransactionAdapter) GetCreatedAt() time.Time {
	return i.income.CreatedAt
}

func (i *IncomeTransactionAdapter) GetType() usecases.TransactionType {
	return usecases.IncomeTransaction
}

func (i *IncomeTransactionAdapter) GetCategoryID() string {
	if i.income.CategoryID == nil {
		return ""
	}
	return *i.income.CategoryID
}

func (i *IncomeTransactionAdapter) IsPending() bool {
	return false // Los ingresos no tienen estado pendiente
}
