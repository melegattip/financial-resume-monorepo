package repository

import (
	"context"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/core/errors"
	"github.com/melegattip/financial-resume-engine/internal/core/ports"
	"gorm.io/gorm"
)

// RecurringTransactionRepository implements the recurring transaction repository
type RecurringTransactionRepository struct {
	db *gorm.DB
}

// NewRecurringTransactionRepository creates a new recurring transaction repository
func NewRecurringTransactionRepository(db *gorm.DB) ports.RecurringTransactionRepository {
	return &RecurringTransactionRepository{
		db: db,
	}
}

// Create creates a new recurring transaction
func (r *RecurringTransactionRepository) Create(ctx context.Context, transaction *domain.RecurringTransaction) error {
	if err := r.db.WithContext(ctx).Create(transaction).Error; err != nil {
		return r.handleDBError(err, "error creando transacción recurrente")
	}
	return nil
}

// GetByID gets a recurring transaction by ID and user ID
func (r *RecurringTransactionRepository) GetByID(ctx context.Context, userID, transactionID string) (*domain.RecurringTransaction, error) {
	var transaction domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", transactionID, userID).
		First(&transaction).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewResourceNotFound("transacción recurrente no encontrada")
		}
		return nil, r.handleDBError(err, "error obteniendo transacción recurrente")
	}

	return &transaction, nil
}

// GetByUserID gets recurring transactions for a user with filters
func (r *RecurringTransactionRepository) GetByUserID(ctx context.Context, userID string, filters ports.RecurringTransactionFilters) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)

	// Apply filters
	query = r.applyFilters(query, filters)

	// Apply sorting
	query = r.applySorting(query, filters)

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return nil, r.handleDBError(err, "error obteniendo transacciones recurrentes")
	}

	// Convert to pointers
	result := make([]*domain.RecurringTransaction, len(transactions))
	for i := range transactions {
		result[i] = &transactions[i]
	}

	return result, nil
}

// Update updates a recurring transaction
func (r *RecurringTransactionRepository) Update(ctx context.Context, transaction *domain.RecurringTransaction) error {
	if err := r.db.WithContext(ctx).Save(transaction).Error; err != nil {
		return r.handleDBError(err, "error actualizando transacción recurrente")
	}
	return nil
}

// Delete deletes a recurring transaction
func (r *RecurringTransactionRepository) Delete(ctx context.Context, userID, transactionID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", transactionID, userID).
		Delete(&domain.RecurringTransaction{})

	if result.Error != nil {
		return r.handleDBError(result.Error, "error eliminando transacción recurrente")
	}

	if result.RowsAffected == 0 {
		return errors.NewResourceNotFound("transacción recurrente no encontrada")
	}

	return nil
}

// GetPendingExecutions gets transactions that should be executed before the given date
func (r *RecurringTransactionRepository) GetPendingExecutions(ctx context.Context, beforeDate time.Time) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("is_active = ? AND next_date <= ? AND auto_create = ?", true, beforeDate, true).
		Find(&transactions).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo ejecuciones pendientes")
	}

	// Filter by business rules
	var result []*domain.RecurringTransaction
	for i := range transactions {
		if transactions[i].ShouldExecute() {
			result = append(result, &transactions[i])
		}
	}

	return result, nil
}

// GetPendingNotifications gets transactions that should send notifications before the given date
func (r *RecurringTransactionRepository) GetPendingNotifications(ctx context.Context, beforeDate time.Time) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&transactions).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo notificaciones pendientes")
	}

	// Filter by business rules
	var result []*domain.RecurringTransaction
	for i := range transactions {
		if transactions[i].ShouldNotify() {
			result = append(result, &transactions[i])
		}
	}

	return result, nil
}

// GetByFrequency gets recurring transactions by frequency
func (r *RecurringTransactionRepository) GetByFrequency(ctx context.Context, userID, frequency string) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND frequency = ?", userID, frequency).
		Find(&transactions).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo transacciones por frecuencia")
	}

	// Convert to pointers
	result := make([]*domain.RecurringTransaction, len(transactions))
	for i := range transactions {
		result[i] = &transactions[i]
	}

	return result, nil
}

// GetByType gets recurring transactions by type
func (r *RecurringTransactionRepository) GetByType(ctx context.Context, userID, transactionType string) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, transactionType).
		Find(&transactions).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo transacciones por tipo")
	}

	// Convert to pointers
	result := make([]*domain.RecurringTransaction, len(transactions))
	for i := range transactions {
		result[i] = &transactions[i]
	}

	return result, nil
}

// GetActiveTransactions gets all active recurring transactions for a user
func (r *RecurringTransactionRepository) GetActiveTransactions(ctx context.Context, userID string) ([]*domain.RecurringTransaction, error) {
	var transactions []domain.RecurringTransaction

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("next_date ASC").
		Find(&transactions).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo transacciones activas")
	}

	// Convert to pointers
	result := make([]*domain.RecurringTransaction, len(transactions))
	for i := range transactions {
		result[i] = &transactions[i]
	}

	return result, nil
}

// GetTotalRecurringAmount gets the total recurring amount for a user and type
func (r *RecurringTransactionRepository) GetTotalRecurringAmount(ctx context.Context, userID, transactionType string) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Model(&domain.RecurringTransaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND type = ? AND is_active = ?", userID, transactionType, true).
		Scan(&total).Error

	if err != nil {
		return 0, r.handleDBError(err, "error calculando total recurrente")
	}

	return total, nil
}

// GetRecurringProjection gets projection data for cash flow analysis
func (r *RecurringTransactionRepository) GetRecurringProjection(ctx context.Context, userID string, months int) (*ports.RecurringProjection, error) {
	// Get monthly totals by type
	type ProjectionData struct {
		Type   string  `gorm:"column:type"`
		Amount float64 `gorm:"column:monthly_amount"`
	}

	var projectionData []ProjectionData

	// This query calculates the monthly equivalent for each transaction type
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			type,
			SUM(
				CASE 
					WHEN frequency = 'daily' THEN amount * 30
					WHEN frequency = 'weekly' THEN amount * 4.33
					WHEN frequency = 'monthly' THEN amount
					WHEN frequency = 'yearly' THEN amount / 12
					ELSE amount
				END
			) as monthly_amount
		FROM recurring_transactions 
		WHERE user_id = ? AND is_active = true
		GROUP BY type
	`, userID).Scan(&projectionData).Error

	if err != nil {
		return nil, r.handleDBError(err, "error obteniendo proyección recurrente")
	}

	projection := &ports.RecurringProjection{}

	for _, data := range projectionData {
		if data.Type == "income" {
			projection.MonthlyIncome = data.Amount
		} else if data.Type == "expense" {
			projection.MonthlyExpenses = data.Amount
		}
	}

	projection.NetMonthly = projection.MonthlyIncome - projection.MonthlyExpenses

	return projection, nil
}

// Helper methods

func (r *RecurringTransactionRepository) applyFilters(query *gorm.DB, filters ports.RecurringTransactionFilters) *gorm.DB {
	if filters.Type != "" {
		query = query.Where("type = ?", filters.Type)
	}

	if filters.Frequency != "" {
		query = query.Where("frequency = ?", filters.Frequency)
	}

	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	return query
}

func (r *RecurringTransactionRepository) applySorting(query *gorm.DB, filters ports.RecurringTransactionFilters) *gorm.DB {
	sortBy := filters.SortBy
	if sortBy == "" {
		sortBy = "next_date" // Default sort
	}

	sortOrder := filters.SortOrder
	if sortOrder == "" {
		sortOrder = "asc" // Default order
	}

	// Validate sort fields
	validSortFields := map[string]bool{
		"next_date":   true,
		"amount":      true,
		"created_at":  true,
		"updated_at":  true,
		"description": true,
		"frequency":   true,
		"type":        true,
	}

	if !validSortFields[sortBy] {
		sortBy = "next_date"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	return query.Order(sortBy + " " + sortOrder)
}

func (r *RecurringTransactionRepository) handleDBError(err error, context string) error {
	if err == nil {
		return nil
	}

	// Handle specific GORM errors
	switch err {
	case gorm.ErrRecordNotFound:
		return errors.NewResourceNotFound("recurso no encontrado")
	case gorm.ErrInvalidTransaction:
		return errors.NewBadRequest("transacción inválida")
	case gorm.ErrNotImplemented:
		return errors.NewInternalServerError("operación no implementada")
	default:
		// For other database errors, wrap with context
		return errors.NewInternalServerError(context + ": " + err.Error())
	}
}
