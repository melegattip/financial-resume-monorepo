package insights

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/usecases"
)

// ExpenseTransactionWrapper adapta domain.Expense a usecases.Transaction
type ExpenseTransactionWrapper struct {
	expense *domain.Expense
}

// NewExpenseTransactionFromDomain crea un wrapper de Transaction para un Expense
func NewExpenseTransactionFromDomain(expense *domain.Expense) usecases.Transaction {
	return &ExpenseTransactionWrapper{
		expense: expense,
	}
}

// GetID implementa Transaction.GetID()
func (w *ExpenseTransactionWrapper) GetID() string {
	return w.expense.GetID()
}

// GetUserID implementa Transaction.GetUserID()
func (w *ExpenseTransactionWrapper) GetUserID() string {
	return w.expense.GetUserID()
}

// GetAmount implementa Transaction.GetAmount()
func (w *ExpenseTransactionWrapper) GetAmount() float64 {
	return w.expense.GetAmount()
}

// GetCreatedAt implementa Transaction.GetCreatedAt()
func (w *ExpenseTransactionWrapper) GetCreatedAt() time.Time {
	return w.expense.GetCreatedAt()
}

// GetType implementa Transaction.GetType()
func (w *ExpenseTransactionWrapper) GetType() usecases.TransactionType {
	return usecases.ExpenseTransaction
}

// GetCategoryID implementa Transaction.GetCategoryID()
func (w *ExpenseTransactionWrapper) GetCategoryID() string {
	return w.expense.GetCategoryID()
}

// IsPending implementa Transaction.IsPending()
func (w *ExpenseTransactionWrapper) IsPending() bool {
	return !w.expense.Paid
}

// IncomeTransactionWrapper adapta domain.Income a usecases.Transaction
type IncomeTransactionWrapper struct {
	income *domain.Income
}

// NewIncomeTransactionFromDomain crea un wrapper de Transaction para un Income
func NewIncomeTransactionFromDomain(income *domain.Income) usecases.Transaction {
	return &IncomeTransactionWrapper{
		income: income,
	}
}

// GetID implementa Transaction.GetID()
func (w *IncomeTransactionWrapper) GetID() string {
	return w.income.GetID()
}

// GetUserID implementa Transaction.GetUserID()
func (w *IncomeTransactionWrapper) GetUserID() string {
	return w.income.GetUserID()
}

// GetAmount implementa Transaction.GetAmount()
func (w *IncomeTransactionWrapper) GetAmount() float64 {
	return w.income.GetAmount()
}

// GetCreatedAt implementa Transaction.GetCreatedAt()
func (w *IncomeTransactionWrapper) GetCreatedAt() time.Time {
	return w.income.GetCreatedAt()
}

// GetType implementa Transaction.GetType()
func (w *IncomeTransactionWrapper) GetType() usecases.TransactionType {
	return usecases.IncomeTransaction
}

// GetCategoryID implementa Transaction.GetCategoryID()
func (w *IncomeTransactionWrapper) GetCategoryID() string {
	return w.income.GetCategoryID()
}

// IsPending implementa Transaction.IsPending()
func (w *IncomeTransactionWrapper) IsPending() bool {
	return false // Los ingresos normalmente no tienen estado pendiente
}
