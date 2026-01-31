package get

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// Service maneja la lógica de negocio para obtener ingresos
type Service struct {
	repository baseRepo.IncomeRepository
}

// IncomeGetter define la interfaz para obtener ingresos
type IncomeGetter interface {
	GetIncome(ctx context.Context, userID, id string) (*incomesDomain.GetIncomeResponse, error)
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.IncomeRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// GetIncome obtiene un ingreso por su ID
func (s *Service) GetIncome(ctx context.Context, userID, id string) (*incomesDomain.GetIncomeResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	if id == "" {
		return nil, errors.NewBadRequest("El ID del ingreso es requerido")
	}

	income, err := s.repository.Get(userID, id)
	if err != nil {
		return nil, err
	}

	categoryID := ""
	if income.CategoryID != nil {
		categoryID = *income.CategoryID
	}

	return &incomesDomain.GetIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
