package transactions

import (
	"context"
)

// BalanceService define las operaciones relacionadas con el balance
type BalanceService interface {
	UpdateBalance(ctx context.Context, userID string, amount float64) error
}

// NotificationService define las operaciones de notificación
type NotificationService interface {
	SendTransactionNotification(ctx context.Context, userID string, message string) error
}
