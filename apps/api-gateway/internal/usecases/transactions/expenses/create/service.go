package create

import (
	"context"
	"log"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
	expensesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/expenses"
)

// ServiceInterface define la interfaz para el servicio de creación de gastos
type ServiceInterface interface {
	CreateExpense(ctx context.Context, request *expensesDomain.CreateExpenseRequest) (*expensesDomain.CreateExpenseResponse, error)
}

// Service maneja la lógica de negocio para la creación de gastos
type Service struct {
	repository         baseRepo.ExpenseRepository
	categoryRepository baseRepo.CategoryRepository
	percentageObserver transactions.PercentageTransactionObserver
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.ExpenseRepository, categoryRepository baseRepo.CategoryRepository, percentageObserver transactions.PercentageTransactionObserver) *Service {
	return &Service{
		repository:         repository,
		categoryRepository: categoryRepository,
		percentageObserver: percentageObserver,
	}
}

// CreateExpense crea un nuevo gasto
func (s *Service) CreateExpense(ctx context.Context, request *expensesDomain.CreateExpenseRequest) (*expensesDomain.CreateExpenseResponse, error) {
	if request.Description == "" {
		return nil, errors.NewBadRequest("El nombre del gasto es requerido")
	}

	if request.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto del gasto debe ser mayor a 0")
	}

	// Validar que la categoría existe si se proporciona
	if request.CategoryID != "" {
		_, err := s.categoryRepository.Get(request.UserID, request.CategoryID)
		if err != nil {
			return nil, errors.NewBadRequest("La categoría especificada no existe")
		}
	}

	var dueDate time.Time
	var err error
	if request.DueDate != "" {
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

	// Calcular el porcentaje antes de crear el gasto
	totalIncome, err := s.percentageObserver.GetTotalIncome(ctx, request.UserID)
	if err != nil {
		return nil, err
	}
	expense.CalculatePercentage(totalIncome)

	// Crear el gasto con el porcentaje calculado
	if err := s.repository.Create(expense); err != nil {
		return nil, err
	}

	// Notificar al PercentageObserver para que actualice los porcentajes de los demás gastos
	if err := s.percentageObserver.OnTransactionCreated(ctx, expense); err != nil {
		// Log del error pero no lo retornamos para no fallar la creación del gasto
		log.Printf("Error al recalcular porcentajes: %v", err)
	}

	return &expensesDomain.CreateExpenseResponse{
		ID:            expense.ID,
		UserID:        expense.UserID,
		Amount:        expense.Amount,
		AmountPaid:    expense.AmountPaid,
		PendingAmount: expense.GetPendingAmount(),
		Description:   expense.Description,
		CategoryID: func() string {
			if expense.CategoryID != nil {
				return *expense.CategoryID
			}
			return ""
		}(),
		Paid:       expense.Paid,
		DueDate:    expense.DueDate.Format("2006-01-02"),
		CreatedAt:  expense.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  expense.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Percentage: expense.Percentage,
	}, nil
}
