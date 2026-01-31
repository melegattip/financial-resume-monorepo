package incomes

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	baseRepo "github.com/melegattip/financial-resume-engine/internal/core/repository"
)

// IncomeRepository define las operaciones para el repositorio de ingresos
type IncomeRepository interface {
	Create(income *domain.Income) error
	Get(userID string, id string) (*domain.Income, error)
	List(userID string) ([]*domain.Income, error)
	Update(income *domain.Income) error
	Delete(userID string, id string) error
	ListUnreceived(userID string) ([]*domain.Income, error)
}

// IncomeService define las operaciones disponibles para el servicio de ingresos
type IncomeService interface {
	CreateIncome(ctx context.Context, request *CreateIncomeRequest) (*CreateIncomeResponse, error)
	GetIncome(ctx context.Context, userID string, incomeID string) (*GetIncomeResponse, error)
	ListIncomes(ctx context.Context, userID string) (*ListIncomesResponse, error)
	UpdateIncome(ctx context.Context, userID string, incomeID string, request *UpdateIncomeRequest) (*UpdateIncomeResponse, error)
	DeleteIncome(ctx context.Context, userID string, incomeID string) error
}

// IncomeUpdater define la interfaz para actualizar ingresos
type IncomeUpdater interface {
	UpdateIncome(ctx context.Context, userID, id string, request *UpdateIncomeRequest) (*UpdateIncomeResponse, error)
}

// IncomeServiceImpl implementa IncomeService
type IncomeServiceImpl struct {
	repository baseRepo.IncomeRepository
}

// NewIncomeService crea una nueva instancia del servicio de ingresos
func NewIncomeService(repository baseRepo.IncomeRepository) IncomeService {
	return &IncomeServiceImpl{
		repository: repository,
	}
}

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

// toGetIncomeResponse convierte un Income en GetIncomeResponse
func toGetIncomeResponse(income *domain.Income) GetIncomeResponse {
	return GetIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
	}
}

// CreateIncome crea un nuevo ingreso
func (s *IncomeServiceImpl) CreateIncome(ctx context.Context, request *CreateIncomeRequest) (*CreateIncomeResponse, error) {
	if request.Description == "" {
		return nil, errors.NewBadRequest("La descripción del ingreso es requerida")
	}

	if request.Amount <= 0 {
		return nil, errors.NewBadRequest("El monto del ingreso debe ser mayor a 0")
	}

	income := domain.NewIncomeBuilder().
		SetID(domain.NewIncomeID()).
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

	return &CreateIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetIncome obtiene un ingreso por su ID
func (s *IncomeServiceImpl) GetIncome(ctx context.Context, userID string, incomeID string) (*GetIncomeResponse, error) {
	income, err := s.repository.Get(userID, incomeID)
	if err != nil {
		return nil, err
	}

	categoryID := ""
	if income.CategoryID != nil {
		categoryID = *income.CategoryID
	}

	return &GetIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListIncomes obtiene todos los ingresos de un usuario
func (s *IncomeServiceImpl) ListIncomes(ctx context.Context, userID string) (*ListIncomesResponse, error) {
	incomes, err := s.repository.List(userID)
	if err != nil {
		return nil, err
	}

	var response []GetIncomeResponse
	for _, income := range incomes {
		categoryID := ""
		if income.CategoryID != nil {
			categoryID = *income.CategoryID
		}
		response = append(response, GetIncomeResponse{
			ID:          income.ID,
			UserID:      income.UserID,
			Amount:      income.Amount,
			Description: income.Description,
			CategoryID:  categoryID,
			CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &ListIncomesResponse{
		Incomes: response,
	}, nil
}

// UpdateIncome actualiza un ingreso
func (s *IncomeServiceImpl) UpdateIncome(ctx context.Context, userID string, incomeID string, request *UpdateIncomeRequest) (*UpdateIncomeResponse, error) {
	income, err := s.repository.Get(userID, incomeID)
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
		if request.CategoryID == "" {
			income.CategoryID = nil
		} else {
			income.CategoryID = &request.CategoryID
		}
	}
	income.UpdatedAt = time.Now()

	if err := s.repository.Update(income); err != nil {
		return nil, err
	}

	categoryID := ""
	if income.CategoryID != nil {
		categoryID = *income.CategoryID
	}

	return &UpdateIncomeResponse{
		ID:          income.ID,
		UserID:      income.UserID,
		Amount:      income.Amount,
		Description: income.Description,
		CategoryID:  categoryID,
		CreatedAt:   income.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   income.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteIncome elimina un ingreso
func (s *IncomeServiceImpl) DeleteIncome(ctx context.Context, userID string, incomeID string) error {
	_, err := s.repository.Get(userID, incomeID)
	if err != nil {
		return err
	}

	return s.repository.Delete(userID, incomeID)
}
