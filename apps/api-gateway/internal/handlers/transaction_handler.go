package handlers

import (
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
)

type TransactionHandler struct {
	subject *transactions.TransactionSubject
}

func NewTransactionHandler(
	expenseService transactions.ExpenseService,
	balanceService transactions.BalanceService,
	notificationService transactions.NotificationService,
) *TransactionHandler {
	handler := &TransactionHandler{
		subject: transactions.NewTransactionSubject(),
	}

	// Registrar observadores
	handler.subject.Attach(&transactions.BalanceObserver{
		BalanceService: balanceService,
	})
	handler.subject.Attach(&transactions.NotificationObserver{
		NotificationService: notificationService,
	})
	handler.subject.Attach(transactions.NewPercentageObserver(expenseService))

	return handler
}
