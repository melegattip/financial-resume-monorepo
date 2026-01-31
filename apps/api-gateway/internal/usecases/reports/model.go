package reports

import (
	"time"
)

type Transaction struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Amount        float64   `json:"amount"`                   // Monto total original
	AmountPaid    float64   `json:"amount_paid,omitempty"`    // Monto ya pagado (solo gastos)
	PendingAmount float64   `json:"pending_amount,omitempty"` // Monto pendiente (solo gastos)
	Description   string    `json:"description"`
	CategoryID    string    `json:"category_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Percentage    float64   `json:"percentage"`
	Type          string    `json:"type"`               // "income" o "expense"
	Paid          bool      `json:"paid,omitempty"`     // Solo para gastos
	DueDate       time.Time `json:"due_date,omitempty"` // Solo para gastos
	Received      bool      `json:"received,omitempty"` // Solo para ingresos
}

// FinancialReport representa el reporte financiero
type FinancialReport struct {
	StartDate       time.Time         `json:"start_date"`
	EndDate         time.Time         `json:"end_date"`
	Transactions    []Transaction     `json:"transactions"`
	TotalIncome     float64           `json:"total_income"`
	TotalExpenses   float64           `json:"total_expenses"`
	NetBalance      float64           `json:"net_balance"`
	CategorySummary []CategorySummary `json:"category_summary"`
}

// CategorySummary representa el resumen de egresos por categoría
type CategorySummary struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Percentage   float64 `json:"percentage"`
}

// GenerateReportRequest representa los parámetros para generar un reporte
type GenerateReportRequest struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
}

// Builder pattern
type FinancialReportBuilder struct {
	report *FinancialReport
}

func NewFinancialReportBuilder() *FinancialReportBuilder {
	return &FinancialReportBuilder{
		report: &FinancialReport{},
	}
}

func (b *FinancialReportBuilder) WithStartDate(date time.Time) *FinancialReportBuilder {
	b.report.StartDate = date
	return b
}

func (b *FinancialReportBuilder) WithEndDate(date time.Time) *FinancialReportBuilder {
	b.report.EndDate = date
	return b
}

func (b *FinancialReportBuilder) WithTransactions(transactions []Transaction) *FinancialReportBuilder {
	b.report.Transactions = transactions
	var totalIncome, totalExpenses float64
	for _, t := range transactions {
		if t.Type == "income" {
			totalIncome += t.Amount
		} else if t.Type == "expense" {
			totalExpenses += t.Amount
		}
	}
	b.report.TotalIncome = totalIncome
	b.report.TotalExpenses = totalExpenses
	b.report.NetBalance = totalIncome - totalExpenses
	return b
}

func (b *FinancialReportBuilder) WithCategorySummary(summary []CategorySummary) *FinancialReportBuilder {
	b.report.CategorySummary = summary
	return b
}

func (b *FinancialReportBuilder) Build() *FinancialReport {
	return b.report
}
