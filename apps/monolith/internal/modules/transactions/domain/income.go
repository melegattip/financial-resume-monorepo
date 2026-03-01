package domain

import (
	"errors"
	"time"
)

// Income represents a financial income transaction
type Income struct {
	ID           string
	UserID       string
	TenantID     string
	CategoryID   string
	Amount       float64
	Source       string
	Description  string
	ReceivedDate time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// NewIncome creates a new income with validation
func NewIncome(userID string, amount float64, source, description string, receivedDate time.Time) (*Income, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}
	if source == "" {
		return nil, errors.New("source is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}

	now := time.Now().UTC()
	return &Income{
		UserID:       userID,
		Amount:       amount,
		Source:       source,
		Description:  description,
		ReceivedDate: receivedDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates mutable fields
func (i *Income) Update(amount float64, source, description string, receivedDate time.Time) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if source == "" {
		return errors.New("source is required")
	}
	if description == "" {
		return errors.New("description is required")
	}

	i.Amount = amount
	i.Source = source
	i.Description = description
	i.ReceivedDate = receivedDate
	i.UpdatedAt = time.Now().UTC()
	return nil
}

// SoftDelete marks the income as deleted
func (i *Income) SoftDelete() {
	now := time.Now().UTC()
	i.DeletedAt = &now
}

// IsDeleted checks if the income has been soft-deleted
func (i *Income) IsDeleted() bool {
	return i.DeletedAt != nil
}
