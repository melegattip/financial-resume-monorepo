package transactions

import (
	"context"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// PercentageObserver se encarga de calcular y actualizar los porcentajes de los gastos
// en relación al total de ingresos
type PercentageObserver struct {
	expenseService ExpenseService
}

func NewPercentageObserver(expenseService ExpenseService) *PercentageObserver {
	return &PercentageObserver{
		expenseService: expenseService,
	}
}

func (o *PercentageObserver) OnTransactionCreated(ctx context.Context, transaction domain.Transaction) error {
	// Recalculamos porcentajes tanto para ingresos como para gastos
	return o.UpdatePercentages(ctx, transaction.GetUserID())
}

func (o *PercentageObserver) OnTransactionUpdated(ctx context.Context, transaction domain.Transaction) error {
	// Recalculamos porcentajes tanto para ingresos como para gastos
	return o.UpdatePercentages(ctx, transaction.GetUserID())
}

func (o *PercentageObserver) OnTransactionDeleted(ctx context.Context, transactionID string) error {
	// Necesitamos obtener el UserID de la transacción eliminada
	// y recalcular los porcentajes
	transaction, err := o.expenseService.GetTransaction(ctx, transactionID)
	if err != nil {
		return err
	}
	return o.UpdatePercentages(ctx, transaction.GetUserID())
}

// UpdatePercentages recalcula los porcentajes de todos los gastos
// en relación al total de ingresos
func (o *PercentageObserver) UpdatePercentages(ctx context.Context, userID string) error {
	return o.expenseService.UpdateExpensePercentages(ctx, userID)
}

// GetTotalIncome obtiene el total de ingresos de un usuario
func (o *PercentageObserver) GetTotalIncome(ctx context.Context, userID string) (float64, error) {
	return o.expenseService.GetTotalIncome(ctx, userID)
}
