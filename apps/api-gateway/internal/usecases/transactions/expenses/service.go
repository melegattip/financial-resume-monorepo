package expenses

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// formatDate convierte time.Time en una cadena de fecha
func formatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}

// formatDateTime convierte time.Time en una cadena de fecha y hora
func formatDateTime(date time.Time) string {
	return date.Format("2006-01-02T15:04:05Z07:00")
}

// ExpenseService define las operaciones disponibles para el servicio de gastos
type ExpenseService interface {
	CreateExpense(ctx context.Context, request *CreateExpenseRequest) (*CreateExpenseResponse, error)
	GetExpense(ctx context.Context, userID string, expenseID string) (*GetExpenseResponse, error)
	ListExpenses(ctx context.Context, userID string) (*ListExpensesResponse, error)
	ListUnpaidExpenses(ctx context.Context, userID string) (*ListExpensesResponse, error)
	UpdateExpense(ctx context.Context, userID string, expenseID string, request *UpdateExpenseRequest) (*UpdateExpenseResponse, error)
	DeleteExpense(ctx context.Context, userID string, expenseID string) error
}

// ExpenseUpdater define la interfaz para actualizar gastos
type ExpenseUpdater interface {
	UpdateExpense(ctx context.Context, userID, id string, request *UpdateExpenseRequest) (*UpdateExpenseResponse, error)
}

// ExpenseServiceImpl implementa ExpenseService
type ExpenseServiceImpl struct {
	repository baseRepo.ExpenseRepository
}

func NewExpenseService(repository baseRepo.ExpenseRepository) ExpenseService {
	return &ExpenseServiceImpl{
		repository: repository,
	}
}

func (s *ExpenseServiceImpl) CreateExpense(ctx context.Context, request *CreateExpenseRequest) (*CreateExpenseResponse, error) {
	if request.Description == "" {
		return nil, errors.NewBadRequest("La descripción del gasto es requerida")
	}

	if request.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto del gasto debe ser mayor a 0")
	}

	// Parsear fecha si se proporciona
	var dueDate time.Time
	if request.DueDate != "" {
		var err error
		dueDate, err = time.Parse("2006-01-02", request.DueDate)
		if err != nil {
			return nil, errors.NewBadRequest("Formato de fecha inválido. Use YYYY-MM-DD")
		}
	}

	expense := domain.NewExpenseBuilder().
		SetID(domain.NewExpenseID()).
		SetUserID(request.UserID).
		SetAmount(request.Amount).
		SetDescription(request.Description).
		SetCategoryID(request.CategoryID).
		SetPaid(request.Paid).
		SetDueDate(dueDate).
		Build()

	if err := s.repository.Create(expense); err != nil {
		return nil, err
	}

	categoryID := ""
	if expense.CategoryID != nil {
		categoryID = *expense.CategoryID
	}

	return &CreateExpenseResponse{
		ID:            expense.ID,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		AmountPaid:    expense.AmountPaid,
		PendingAmount: expense.GetPendingAmount(),
		Description:   expense.Description,
		CategoryID:    categoryID,
		Paid:          expense.Paid,
		DueDate:       formatDate(expense.DueDate),
		CreatedAt:     formatDateTime(expense.CreatedAt),
		UpdatedAt:     formatDateTime(expense.UpdatedAt),
		Percentage:    expense.Percentage,
	}, nil
}

func (s *ExpenseServiceImpl) GetExpense(ctx context.Context, userID string, expenseID string) (*GetExpenseResponse, error) {
	expense, err := s.repository.Get(userID, expenseID)
	if err != nil {
		return nil, err
	}

	categoryID := ""
	if expense.CategoryID != nil {
		categoryID = *expense.CategoryID
	}

	return &GetExpenseResponse{
		ID:            expense.ID,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		AmountPaid:    expense.AmountPaid,
		PendingAmount: expense.GetPendingAmount(),
		Description:   expense.Description,
		CategoryID:    categoryID,
		Paid:          expense.Paid,
		DueDate:       formatDate(expense.DueDate),
		CreatedAt:     formatDateTime(expense.CreatedAt),
		UpdatedAt:     formatDateTime(expense.UpdatedAt),
		Percentage:    expense.Percentage,
	}, nil
}

func (s *ExpenseServiceImpl) ListExpenses(ctx context.Context, userID string) (*ListExpensesResponse, error) {
	expenses, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var response []GetExpenseResponse
	for _, expense := range expenses {
		categoryID := ""
		if expense.CategoryID != nil {
			categoryID = *expense.CategoryID
		}
		response = append(response, GetExpenseResponse{
			ID:            expense.ID,
			UserID:        expense.UserID,
			Amount:        expense.Amount,
			AmountPaid:    expense.AmountPaid,
			PendingAmount: expense.GetPendingAmount(),
			Description:   expense.Description,
			CategoryID:    categoryID,
			Paid:          expense.Paid,
			DueDate:       formatDate(expense.DueDate),
			CreatedAt:     formatDateTime(expense.CreatedAt),
			UpdatedAt:     formatDateTime(expense.UpdatedAt),
			Percentage:    expense.Percentage,
		})
	}

	return &ListExpensesResponse{
		Expenses: response,
	}, nil
}

func (s *ExpenseServiceImpl) ListUnpaidExpenses(ctx context.Context, userID string) (*ListExpensesResponse, error) {
	expenses, err := s.repository.ListUnpaid(userID)
	if err != nil {
		return nil, err
	}

	var response []GetExpenseResponse
	for _, expense := range expenses {
		categoryID := ""
		if expense.CategoryID != nil {
			categoryID = *expense.CategoryID
		}
		response = append(response, GetExpenseResponse{
			ID:            expense.ID,
			UserID:        expense.UserID,
			Amount:        expense.Amount,
			AmountPaid:    expense.AmountPaid,
			PendingAmount: expense.GetPendingAmount(),
			Description:   expense.Description,
			CategoryID:    categoryID,
			Paid:          expense.Paid,
			DueDate:       formatDate(expense.DueDate),
			CreatedAt:     formatDateTime(expense.CreatedAt),
			UpdatedAt:     formatDateTime(expense.UpdatedAt),
			Percentage:    expense.Percentage,
		})
	}

	return &ListExpensesResponse{
		Expenses: response,
	}, nil
}

func (s *ExpenseServiceImpl) UpdateExpense(ctx context.Context, userID string, expenseID string, request *UpdateExpenseRequest) (*UpdateExpenseResponse, error) {
	expense, err := s.repository.Get(userID, expenseID)
	if err != nil {
		return nil, err
	}

	if request.Amount != 0 {
		expense.Amount = request.Amount
	}
	if request.Description != "" {
		expense.Description = request.Description
	}
	if request.CategoryID != "" {
		if request.CategoryID == "" {
			expense.CategoryID = nil
		} else {
			expense.CategoryID = &request.CategoryID
		}
	}
	// Actualizar Paid - en el request es bool, no puntero
	expense.Paid = request.Paid

	if request.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", request.DueDate)
		if err != nil {
			return nil, errors.NewBadRequest("Formato de fecha inválido. Use YYYY-MM-DD")
		}
		expense.DueDate = dueDate
	}
	expense.UpdatedAt = time.Now()

	if err := s.repository.Update(expense); err != nil {
		return nil, err
	}

	categoryID := ""
	if expense.CategoryID != nil {
		categoryID = *expense.CategoryID
	}

	return &UpdateExpenseResponse{
		ID:            expense.ID,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		AmountPaid:    expense.AmountPaid,
		PendingAmount: expense.GetPendingAmount(),
		Description:   expense.Description,
		CategoryID:    categoryID,
		Paid:          expense.Paid,
		DueDate:       formatDate(expense.DueDate),
		CreatedAt:     formatDateTime(expense.CreatedAt),
		UpdatedAt:     formatDateTime(expense.UpdatedAt),
		Percentage:    expense.Percentage,
	}, nil
}

func (s *ExpenseServiceImpl) DeleteExpense(ctx context.Context, userID string, expenseID string) error {
	_, err := s.repository.Get(userID, expenseID)
	if err != nil {
		return err
	}

	return s.repository.Delete(userID, expenseID)
}
