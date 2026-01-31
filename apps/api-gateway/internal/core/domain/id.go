package domain

import (
	"github.com/google/uuid"
)

// NewID genera un nuevo ID único usando UUID v4
func NewID() string {
	return uuid.New().String()
}

// NewExpenseID genera un nuevo ID para gastos con el formato exp_ + UUID
func NewExpenseID() string {
	return "exp_" + uuid.New().String()[:7]
}

// NewIncomeID genera un nuevo ID para ingresos con el formato inc_ + UUID
func NewIncomeID() string {
	return "inc_" + uuid.New().String()[:7]
}
