package create

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// ServiceInterface define la interfaz para el servicio de creación de ingresos
type ServiceInterface interface {
	CreateIncome(ctx context.Context, request *incomesDomain.CreateIncomeRequest) (*incomesDomain.CreateIncomeResponse, error)
}

// Service maneja la lógica de negocio para la creación de ingresos
type Service struct {
	repository         baseRepo.IncomeRepository
	categoryRepository baseRepo.CategoryRepository
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.IncomeRepository, categoryRepository baseRepo.CategoryRepository) *Service {
	return &Service{
		repository:         repository,
		categoryRepository: categoryRepository,
	}
}

// CreateIncome crea un nuevo ingreso
func (s *Service) CreateIncome(ctx context.Context, request *incomesDomain.CreateIncomeRequest) (*incomesDomain.CreateIncomeResponse, error) {
	if request.Description == "" {
		return nil, errors.NewBadRequest("La descripción del ingreso es requerida")
	}

	if request.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto del ingreso debe ser mayor a 0")
	}

	if request.Source == "" {
		return nil, errors.NewBadRequest("La fuente del ingreso es requerida")
	}

	// Validar que la categoría existe si se proporciona
	if request.CategoryID != "" {
		_, err := s.categoryRepository.Get(request.UserID, request.CategoryID)
		if err != nil {
			return nil, errors.NewBadRequest("La categoría especificada no existe")
		}
	}

	income := domain.NewIncomeBuilder().
		SetID(domain.NewID()).
		SetUserID(request.UserID).
		SetAmount(request.Amount).
		SetDescription(request.Description).
		SetCategoryID(request.CategoryID).
		Build()

	if err := s.repository.Create(income); err != nil {
		return nil, err
	}

	categoryID := ""
	if income.CategoryID != nil {
		categoryID = *income.CategoryID
	}

	return &incomesDomain.CreateIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
