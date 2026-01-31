package transactions

import "time"

// TransactionModel implementa la interfaz Transaction
type TransactionModel struct {
	ID          string
	UserID      string
	Amount      float64
	Description string
	CategoryID  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *TransactionModel) GetID() string {
	return t.ID
}

func (t *TransactionModel) GetUserID() string {
	return t.UserID
}

func (t *TransactionModel) GetAmount() float64 {
	return t.Amount
}

func (t *TransactionModel) GetDescription() string {
	return t.Description
}

func (t *TransactionModel) GetCategoryID() string {
	return t.CategoryID
}

func (t *TransactionModel) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *TransactionModel) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
