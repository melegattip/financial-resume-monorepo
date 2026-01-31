package reports

import (
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"gorm.io/gorm"
)

// ReportRepository define la interfaz para el repositorio de reportes
type ReportRepository interface {
	GetTransactions(startDate, endDate time.Time, userID string) ([]Transaction, error)
	GetCategoryNames(categoryIDs []string) (map[string]string, error)
}

type ReportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{db: db}
}

func (r *ReportRepositoryImpl) GetTransactions(startDate, endDate time.Time, userID string) ([]Transaction, error) {
	var transactions []Transaction

	// Obtener gastos
	var expenses []domain.Expense
	if err := r.db.Where("created_at BETWEEN ? AND ? AND user_id = ?", startDate, endDate, userID).Find(&expenses).Error; err != nil {
		return nil, err
	}

	// Convertir gastos a transacciones
	for _, expense := range expenses {
		categoryID := ""
		if expense.CategoryID != nil {
			categoryID = *expense.CategoryID
		}
		transactions = append(transactions, Transaction{
			ID:            expense.ID,
			UserID:        expense.UserID,
			Amount:        expense.Amount,
			AmountPaid:    expense.AmountPaid,
			PendingAmount: expense.GetPendingAmount(),
			Description:   expense.Description,
			CategoryID:    categoryID,
			CreatedAt:     expense.CreatedAt,
			UpdatedAt:     expense.UpdatedAt,
			Percentage:    expense.Percentage,
			Type:          "expense",
			Paid:          expense.Paid,
			DueDate:       expense.DueDate,
		})
	}

	// Obtener ingresos
	var incomes []domain.Income
	if err := r.db.Where("created_at BETWEEN ? AND ? AND user_id = ?", startDate, endDate, userID).Find(&incomes).Error; err != nil {
		return nil, err
	}

	// Convertir ingresos a transacciones
	for _, income := range incomes {
		categoryID := ""
		if income.CategoryID != nil {
			categoryID = *income.CategoryID
		}
		transactions = append(transactions, Transaction{
			ID:          income.ID,
			UserID:      income.UserID,
			Amount:      income.Amount,
			Description: income.Description,
			CategoryID:  categoryID,
			CreatedAt:   income.CreatedAt,
			UpdatedAt:   income.UpdatedAt,
			Percentage:  income.Percentage,
			Type:        "income",
		})
	}

	return transactions, nil
}

func (r *ReportRepositoryImpl) GetCategoryNames(categoryIDs []string) (map[string]string, error) {
	var categories []domain.Category
	if err := r.db.Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
		return nil, err
	}

	categoryMap := make(map[string]string)
	for _, category := range categories {
		categoryMap[category.ID] = category.Name
	}

	return categoryMap, nil
}
