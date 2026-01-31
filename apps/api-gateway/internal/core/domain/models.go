package domain

import (
	"time"
)

// BaseModel representa la estructura base para todos los modelos
type BaseModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// FinancialData representa los datos financieros
type FinancialData struct {
	BaseModel
	UserID      uint    `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}
