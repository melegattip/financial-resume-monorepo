package domain

import (
	"errors"
	"time"
)

// Expense represents a financial expense transaction
type Expense struct {
	ID              string
	UserID          string
	CategoryID      string
	Amount          float64
	Description     string
	TransactionDate time.Time
	PaymentMethod   string
	Notes           string
	Paid            bool
	AmountPaid      float64
	PendingAmount   float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

// NewExpense creates a new expense with validation
func NewExpense(userID, categoryID string, amount float64, description string, transactionDate time.Time, paymentMethod string) (*Expense, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if categoryID == "" {
		return nil, errors.New("category_id is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}

	now := time.Now().UTC()
	return &Expense{
		UserID:          userID,
		CategoryID:      categoryID,
		Amount:          amount,
		Description:     description,
		TransactionDate: transactionDate,
		PaymentMethod:   paymentMethod,
		Paid:            false,
		AmountPaid:      0,
		PendingAmount:   amount,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Update updates mutable fields
func (e *Expense) Update(categoryID string, amount float64, description string, transactionDate time.Time, paymentMethod, notes string) error {
	if categoryID == "" {
		return errors.New("category_id is required")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if description == "" {
		return errors.New("description is required")
	}

	e.CategoryID = categoryID
	e.Amount = amount
	e.Description = description
	e.TransactionDate = transactionDate
	e.PaymentMethod = paymentMethod
	e.Notes = notes
	e.UpdatedAt = time.Now().UTC()
	return nil
}

// ApplyPayment updates the paid state. amountPaid is the new cumulative total paid.
// Pass paid=false to mark the expense as unpaid (resets amountPaid to 0).
func (e *Expense) ApplyPayment(paid bool, amountPaid float64) {
	e.Paid = paid
	if !paid {
		e.AmountPaid = 0
		e.PendingAmount = e.Amount
	} else {
		e.AmountPaid = amountPaid
		pending := e.Amount - amountPaid
		if pending < 0 {
			pending = 0
		}
		e.PendingAmount = pending
	}
	e.UpdatedAt = time.Now().UTC()
}

// SoftDelete marks the expense as deleted
func (e *Expense) SoftDelete() {
	now := time.Now().UTC()
	e.DeletedAt = &now
}

// IsDeleted checks if the expense has been soft-deleted
func (e *Expense) IsDeleted() bool {
	return e.DeletedAt != nil
}
