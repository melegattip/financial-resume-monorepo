package list

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// Service maneja la lógica de negocio para listar ingresos
type Service struct {
	repository baseRepo.IncomeRepository
}

// IncomeLister define la interfaz para listar ingresos
type IncomeLister interface {
	ListIncomes(ctx context.Context, userID string) (*incomesDomain.ListIncomesResponse, error)
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.IncomeRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// ListIncomes lista todos los ingresos de un usuario
func (s *Service) ListIncomes(ctx context.Context, userID string) (*incomesDomain.ListIncomesResponse, error) {
	if userID == "" {
		return nil, errors.NewBadRequest("El ID del usuario es requerido")
	}

	incomes, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	incomesResponse := make([]incomesDomain.GetIncomeResponse, len(incomes))
	for i, income := range incomes {
		categoryID := ""
		if income.CategoryID != nil {
			categoryID = *income.CategoryID
		}
		incomesResponse[i] = incomesDomain.GetIncomeResponse{
			ID:          income.ID,
			UserID:      income.UserID,
			Amount:      income.Amount,
			Description: income.Description,
			CategoryID:  categoryID,
			CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &incomesDomain.ListIncomesResponse{
		Incomes: incomesResponse,
	}, nil
}
