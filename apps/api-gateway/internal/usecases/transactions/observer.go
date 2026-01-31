package transactions

import (
	"context"
	"fmt"
	"sync"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
)

// TransactionObserver define la interfaz base para todos los observadores de transacciones
type TransactionObserver interface {
	OnTransactionCreated(ctx context.Context, transaction domain.Transaction) error
	OnTransactionUpdated(ctx context.Context, transaction domain.Transaction) error
	OnTransactionDeleted(ctx context.Context, transactionID string) error
}

// PercentageTransactionObserver extiende TransactionObserver con funcionalidad específica para el manejo de porcentajes
type PercentageTransactionObserver interface {
	TransactionObserver
	GetTotalIncome(ctx context.Context, userID string) (float64, error)
}

// TransactionSubject mantiene la lista de observadores y notifica los cambios
type TransactionSubject struct {
	observers []TransactionObserver
	mu        sync.RWMutex
}

func NewTransactionSubject() *TransactionSubject {
	return &TransactionSubject{
		observers: make([]TransactionObserver, 0),
	}
}

func (s *TransactionSubject) Attach(observer TransactionObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.observers = append(s.observers, observer)
}

func (s *TransactionSubject) Detach(observer TransactionObserver) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, o := range s.observers {
		if o == observer {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

func (s *TransactionSubject) NotifyTransactionCreated(ctx context.Context, transaction domain.Transaction) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var errs []error
	for _, observer := range s.observers {
		if err := observer.OnTransactionCreated(ctx, transaction); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores al notificar creación: %v", errs)
	}
	return nil
}

func (s *TransactionSubject) NotifyTransactionUpdated(ctx context.Context, transaction domain.Transaction) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var errs []error
	for _, observer := range s.observers {
		if err := observer.OnTransactionUpdated(ctx, transaction); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores al notificar actualización: %v", errs)
	}
	return nil
}

func (s *TransactionSubject) NotifyTransactionDeleted(ctx context.Context, transactionID string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var errs []error
	for _, observer := range s.observers {
		if err := observer.OnTransactionDeleted(ctx, transactionID); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errores al notificar eliminación: %v", errs)
	}
	return nil
}

// BalanceObserver actualiza el balance cuando se crea o actualiza una transacción
type BalanceObserver struct {
	BalanceService BalanceService
}

func (o *BalanceObserver) OnTransactionCreated(ctx context.Context, transaction domain.Transaction) error {
	return o.BalanceService.UpdateBalance(ctx, transaction.GetUserID(), transaction.GetAmount())
}

func (o *BalanceObserver) OnTransactionUpdated(ctx context.Context, transaction domain.Transaction) error {
	return o.BalanceService.UpdateBalance(ctx, transaction.GetUserID(), transaction.GetAmount())
}

func (o *BalanceObserver) OnTransactionDeleted(ctx context.Context, transactionID string) error {
	return nil
}

// NotificationObserver envía notificaciones cuando se crea una transacción
type NotificationObserver struct {
	NotificationService NotificationService
}

func (o *NotificationObserver) OnTransactionCreated(ctx context.Context, transaction domain.Transaction) error {
	return o.NotificationService.SendTransactionNotification(ctx, transaction.GetUserID(), "Nueva transacción creada")
}

func (o *NotificationObserver) OnTransactionUpdated(ctx context.Context, transaction domain.Transaction) error {
	return nil
}

func (o *NotificationObserver) OnTransactionDeleted(ctx context.Context, transactionID string) error {
	return nil
}
