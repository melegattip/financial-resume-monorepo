package reports

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReportRepository es un mock del repositorio de reportes
type MockReportRepository struct {
	mock.Mock
}

func (m *MockReportRepository) GetTransactions(startDate, endDate time.Time, userID string) ([]Transaction, error) {
	args := m.Called(startDate, endDate, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Transaction), args.Error(1)
}

func (m *MockReportRepository) GetCategoryNames(categoryIDs []string) (map[string]string, error) {
	args := m.Called(categoryIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

// Verificar que MockReportRepository implementa la interfaz
var _ ReportRepository = (*MockReportRepository)(nil)

func TestGenerateFinancialReport_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockReportRepository)
	service := NewGenerateFinancialReport(mockRepo)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	userID := "user123"

	transactions := []Transaction{
		{
			ID:          "income1",
			UserID:      userID,
			Amount:      1000.0,
			Description: "Salary",
			CategoryID:  "cat1",
			Type:        "income",
			CreatedAt:   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          "expense1",
			UserID:      userID,
			Amount:      100.0,
			Description: "Groceries",
			CategoryID:  "cat2",
			Type:        "expense",
			CreatedAt:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	categoryNames := map[string]string{
		"cat1": "Work",
		"cat2": "Food",
	}

	mockRepo.On("GetTransactions", startDate, endDate, userID).Return(transactions, nil)
	mockRepo.On("GetCategoryNames", []string{"cat1", "cat2"}).Return(categoryNames, nil)

	// Act
	report, err := service.Execute(startDate, endDate, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, startDate, report.StartDate)
	assert.Equal(t, endDate, report.EndDate)
	assert.Equal(t, 1000.0, report.TotalIncome)
	assert.Equal(t, 100.0, report.TotalExpenses)
	assert.Equal(t, 900.0, report.NetBalance)
	assert.Len(t, report.Transactions, 2)
	assert.Len(t, report.CategorySummary, 1) // Solo categorías de gastos

	// Verificar que el porcentaje del gasto se calculó correctamente
	expenseTransaction := report.Transactions[1]         // El gasto
	assert.Equal(t, 10.0, expenseTransaction.Percentage) // 100/1000 = 10%

	// Verificar resumen de categorías
	categorySummary := report.CategorySummary[0]
	assert.Equal(t, "cat2", categorySummary.CategoryID)
	assert.Equal(t, "Food", categorySummary.CategoryName)
	assert.Equal(t, 100.0, categorySummary.TotalAmount)
	assert.Equal(t, 100.0, categorySummary.Percentage) // 100% de los gastos

	mockRepo.AssertExpectations(t)
}

func TestGenerateFinancialReport_EmptyTransactions(t *testing.T) {
	// Arrange
	mockRepo := new(MockReportRepository)
	service := NewGenerateFinancialReport(mockRepo)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	userID := "user123"

	mockRepo.On("GetTransactions", startDate, endDate, userID).Return([]Transaction{}, nil)
	mockRepo.On("GetCategoryNames", []string{}).Return(map[string]string{}, nil)

	// Act
	report, err := service.Execute(startDate, endDate, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 0.0, report.TotalIncome)
	assert.Equal(t, 0.0, report.TotalExpenses)
	assert.Equal(t, 0.0, report.NetBalance)
	assert.Len(t, report.Transactions, 0)
	assert.Len(t, report.CategorySummary, 0)

	mockRepo.AssertExpectations(t)
}

func TestGenerateFinancialReport_OnlyExpenses(t *testing.T) {
	// Arrange
	mockRepo := new(MockReportRepository)
	service := NewGenerateFinancialReport(mockRepo)

	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	userID := "user123"

	transactions := []Transaction{
		{
			ID:          "expense1",
			UserID:      userID,
			Amount:      100.0,
			Description: "Groceries",
			CategoryID:  "cat1",
			Type:        "expense",
			CreatedAt:   time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	categoryNames := map[string]string{
		"cat1": "Food",
	}

	mockRepo.On("GetTransactions", startDate, endDate, userID).Return(transactions, nil)
	mockRepo.On("GetCategoryNames", []string{"cat1"}).Return(categoryNames, nil)

	// Act
	report, err := service.Execute(startDate, endDate, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, report)
	assert.Equal(t, 0.0, report.TotalIncome)
	assert.Equal(t, 100.0, report.TotalExpenses)
	assert.Equal(t, -100.0, report.NetBalance) // Balance negativo

	// Sin ingresos, el porcentaje del gasto debería ser 0
	expenseTransaction := report.Transactions[0]
	assert.Equal(t, 0.0, expenseTransaction.Percentage)

	mockRepo.AssertExpectations(t)
}
