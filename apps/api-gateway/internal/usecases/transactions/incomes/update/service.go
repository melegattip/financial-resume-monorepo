package update

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
	incomesDomain "github.com/melegattip/financial-resume-engine/internal/usecases/transactions/incomes"
)

// Service maneja la lógica de negocio para actualizar ingresos
type Service struct {
	repository         baseRepo.IncomeRepository
	categoryRepository baseRepo.CategoryRepository
}

// IncomeUpdater define la interfaz para actualizar ingresos
type IncomeUpdater interface {
	UpdateIncome(ctx context.Context, userID string, id string, request *incomesDomain.UpdateIncomeRequest) (*incomesDomain.UpdateIncomeResponse, error)
}

// NewService crea una nueva instancia del servicio
func NewService(repository baseRepo.IncomeRepository, categoryRepository baseRepo.CategoryRepository) *Service {
	return &Service{
		repository:         repository,
		categoryRepository: categoryRepository,
	}
}

// UpdateIncome actualiza un ingreso existente
func (s *Service) UpdateIncome(ctx context.Context, userID string, id string, request *incomesDomain.UpdateIncomeRequest) (*incomesDomain.UpdateIncomeResponse, error) {
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

	if request.Amount != 0 {
		income.Amount = request.Amount
	}
	if request.Description != "" {
		income.Description = request.Description
	}
	if request.CategoryID != "" {
		// Validar que la categoría existe
		_, err := s.categoryRepository.Get(userID, request.CategoryID)
		if err != nil {
			return nil, errors.NewBadRequest("La categoría especificada no existe")
		}
		income.CategoryID = &request.CategoryID
	}
	income.UpdatedAt = time.Now()

	if err := s.repository.Update(income); err != nil {
		return nil, err
	}

	categoryID := ""
	if income.CategoryID != nil {
		categoryID = *income.CategoryID
	}

	return &incomesDomain.UpdateIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
