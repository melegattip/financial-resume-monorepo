package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/ports"
)

// RecurringTransactionScheduler maneja la ejecución automática de transacciones recurrentes
type RecurringTransactionScheduler struct {
	useCase   ports.RecurringTransactionUseCase
	ticker    *time.Ticker
	done      chan bool
	isRunning bool
	interval  time.Duration
}

// NewRecurringTransactionScheduler crea una nueva instancia del scheduler
func NewRecurringTransactionScheduler(useCase ports.RecurringTransactionUseCase, interval time.Duration) *RecurringTransactionScheduler {
	return &RecurringTransactionScheduler{
		useCase:   useCase,
		interval:  interval,
		done:      make(chan bool),
		isRunning: false,
	}
}

// Start inicia el scheduler
func (s *RecurringTransactionScheduler) Start() {
	if s.isRunning {
		log.Println("📅 Recurring transaction scheduler is already running")
		return
	}

	log.Printf("🚀 Starting recurring transaction scheduler with interval: %v", s.interval)

	s.ticker = time.NewTicker(s.interval)
	s.isRunning = true

	go func() {
		// Ejecutar inmediatamente al iniciar
		s.processTransactions()

		// Luego ejecutar según el intervalo
		for {
			select {
			case <-s.ticker.C:
				s.processTransactions()
			case <-s.done:
				log.Println("🛑 Stopping recurring transaction scheduler")
				return
			}
		}
	}()
}

// Stop detiene el scheduler
func (s *RecurringTransactionScheduler) Stop() {
	if !s.isRunning {
		return
	}

	log.Println("🛑 Stopping recurring transaction scheduler...")
	s.ticker.Stop()
	s.done <- true
	s.isRunning = false
}

// processTransactions procesa las transacciones pendientes
func (s *RecurringTransactionScheduler) processTransactions() {
	ctx := context.Background()

	log.Println("⏰ Processing pending recurring transactions...")

	// Procesar transacciones pendientes
	result, err := s.useCase.ProcessPendingTransactions(ctx)
	if err != nil {
		log.Printf("❌ Error processing recurring transactions: %v", err)
		return
	}

	if result.ProcessedCount > 0 {
		log.Printf("✅ Processed %d recurring transactions (Success: %d, Failed: %d)",
			result.ProcessedCount, result.SuccessCount, result.FailureCount)

		if result.FailureCount > 0 {
			log.Printf("⚠️ Failed transactions errors: %v", result.Errors)
		}
	} else {
		log.Println("ℹ️ No pending recurring transactions to process")
	}

	// Enviar notificaciones pendientes
	notificationResult, err := s.useCase.SendPendingNotifications(ctx)
	if err != nil {
		log.Printf("❌ Error sending notifications: %v", err)
		return
	}

	if notificationResult.SentCount > 0 {
		log.Printf("📧 Sent %d notifications (Failed: %d)",
			notificationResult.SentCount, notificationResult.FailureCount)

		if notificationResult.FailureCount > 0 {
			log.Printf("⚠️ Failed notifications errors: %v", notificationResult.Errors)
		}
	}
}

// IsRunning retorna si el scheduler está ejecutándose
func (s *RecurringTransactionScheduler) IsRunning() bool {
	return s.isRunning
}

// GetStatus retorna el estado del scheduler
func (s *RecurringTransactionScheduler) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"running":  s.isRunning,
		"interval": s.interval.String(),
	}
}

// ProcessNow fuerza el procesamiento inmediato (para testing/admin)
func (s *RecurringTransactionScheduler) ProcessNow() {
	log.Println("🔄 Manual processing triggered")
	s.processTransactions()
}
