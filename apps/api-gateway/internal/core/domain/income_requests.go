package domain

import "time"

// CreateIncomeRequest representa la solicitud para crear un ingreso
type CreateIncomeRequest struct {
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CategoryID  string    `json:"category_id"`
	UserID      string    `json:"user_id"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

// UpdateIncomeRequest representa la solicitud para actualizar un ingreso
type UpdateIncomeRequest struct {
	Amount      float64   `json:"amount,omitempty"`
	Description string    `json:"description,omitempty"`
	CategoryID  string    `json:"category_id,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
	Received    bool      `json:"received,omitempty"`
}

// DeleteIncomeRequest representa la solicitud para eliminar un ingreso
type DeleteIncomeRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}
