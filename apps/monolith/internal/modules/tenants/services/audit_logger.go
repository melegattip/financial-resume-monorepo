package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/ports"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// AuditLogger subscribes to domain events and writes immutable audit log entries
// to the tenant-scoped audit_logs table.
//
// Because the base Event interface does not expose a TenantID, this service
// queries the affected entity's table by AggregateID to resolve the tenant.
type AuditLogger struct {
	db     *gorm.DB
	repo   ports.AuditRepository
	logger zerolog.Logger
}

// NewAuditLogger creates a new AuditLogger.
func NewAuditLogger(db *gorm.DB, repo ports.AuditRepository, logger zerolog.Logger) *AuditLogger {
	return &AuditLogger{db: db, repo: repo, logger: logger}
}

// Handle processes a domain event and writes an audit log entry.
// Errors are logged as warnings and never propagated to the caller — audit
// logging is best-effort and must not disrupt the primary business flow.
func (al *AuditLogger) Handle(ctx context.Context, event sharedports.Event) error {
	tenantID, err := al.resolveTenantID(ctx, event)
	if err != nil || tenantID == "" {
		al.logger.Warn().
			Err(err).
			Str("event_type", event.EventType()).
			Str("aggregate_id", event.AggregateID()).
			Msg("audit: could not resolve tenant_id, skipping")
		return nil
	}

	entityType := resolveEntityType(event.EventType())
	description := al.buildDescription(ctx, event.EventType(), event.AggregateID())

	entry := domain.AuditLog{
		ID:          "aud_" + strings.ReplaceAll(uuid.New().String(), "-", "")[:12],
		TenantID:    tenantID,
		UserID:      event.UserID(),
		Action:      event.EventType(),
		Description: description,
		EntityType:  entityType,
		EntityID:    event.AggregateID(),
		CreatedAt:   time.Now().UTC(),
	}

	if err := al.repo.CreateAuditLog(ctx, entry); err != nil {
		al.logger.Warn().
			Err(err).
			Str("event_type", event.EventType()).
			Str("tenant_id", tenantID).
			Msg("audit: failed to persist audit log entry")
	}

	return nil
}

// resolveTenantID looks up the tenant_id for the entity referenced by the event.
// For expense/income/recurring events it queries the entity table directly.
// For user.registered it queries tenant_members to find the user's personal tenant.
func (al *AuditLogger) resolveTenantID(ctx context.Context, event sharedports.Event) (string, error) {
	eventType := event.EventType()

	if eventType == "user.registered" {
		var tenantID string
		if err := al.db.WithContext(ctx).
			Raw("SELECT tenant_id FROM tenant_members WHERE user_id = ? AND role = 'owner' LIMIT 1", event.UserID()).
			Scan(&tenantID).Error; err != nil {
			return "", err
		}
		return tenantID, nil
	}

	table := resolveTable(eventType)
	if table == "" {
		return "", nil
	}

	var tenantID string
	//nolint:gosec // table is derived from a controlled switch, not user input
	query := "SELECT tenant_id FROM " + table + " WHERE id = ? LIMIT 1"
	if err := al.db.WithContext(ctx).Raw(query, event.AggregateID()).Scan(&tenantID).Error; err != nil {
		return "", err
	}
	return tenantID, nil
}

// buildDescription fetches entity details from the DB and returns a concise
// human-readable summary of what happened. Best-effort: returns empty string on failure.
func (al *AuditLogger) buildDescription(ctx context.Context, eventType, aggregateID string) string {
	switch {
	case strings.HasPrefix(eventType, "expense."):
		type row struct {
			Description string
			Amount      float64
		}
		var r row
		if err := al.db.WithContext(ctx).
			Raw("SELECT description, amount FROM expenses WHERE id = ? LIMIT 1", aggregateID).
			Scan(&r).Error; err != nil || r.Description == "" {
			return ""
		}
		return fmt.Sprintf("%s · $%.2f", r.Description, r.Amount)

	case strings.HasPrefix(eventType, "income."):
		type row struct {
			Source string
			Amount float64
		}
		var r row
		if err := al.db.WithContext(ctx).
			Raw("SELECT source, amount FROM incomes WHERE id = ? LIMIT 1", aggregateID).
			Scan(&r).Error; err != nil || r.Source == "" {
			return ""
		}
		return fmt.Sprintf("%s · $%.2f", r.Source, r.Amount)

	case strings.HasPrefix(eventType, "recurring."):
		type row struct {
			Description string
			Amount      float64
			Type        string
		}
		var r row
		if err := al.db.WithContext(ctx).
			Raw("SELECT description, amount, type FROM recurring_transactions WHERE id = ? LIMIT 1", aggregateID).
			Scan(&r).Error; err != nil || r.Description == "" {
			return ""
		}
		return fmt.Sprintf("%s · $%.2f", r.Description, r.Amount)

	case strings.HasPrefix(eventType, "budget."):
		type row struct {
			Amount float64
			Period string
		}
		var r row
		if err := al.db.WithContext(ctx).
			Raw("SELECT amount, period FROM budgets WHERE id = ? LIMIT 1", aggregateID).
			Scan(&r).Error; err != nil || r.Amount == 0 {
			return ""
		}
		return fmt.Sprintf("$%.2f (%s)", r.Amount, r.Period)

	case strings.HasPrefix(eventType, "savings_goal."):
		type row struct {
			Name         string
			TargetAmount float64
		}
		var r row
		if err := al.db.WithContext(ctx).
			Raw("SELECT name, target_amount FROM savings_goals WHERE id = ? LIMIT 1", aggregateID).
			Scan(&r).Error; err != nil || r.Name == "" {
			return ""
		}
		return fmt.Sprintf("%s · $%.2f", r.Name, r.TargetAmount)

	case eventType == "user.registered":
		return "Usuario se unió al espacio"
	}
	return ""
}

// resolveTable returns the database table name for a given event type prefix.
func resolveTable(eventType string) string {
	switch {
	case strings.HasPrefix(eventType, "expense."):
		return "expenses"
	case strings.HasPrefix(eventType, "income."):
		return "incomes"
	case strings.HasPrefix(eventType, "recurring."):
		return "recurring_transactions"
	case strings.HasPrefix(eventType, "budget."):
		return "budgets"
	case strings.HasPrefix(eventType, "savings_goal."):
		return "savings_goals"
	}
	return ""
}

// resolveEntityType maps an event type to a human-readable entity type string.
func resolveEntityType(eventType string) string {
	switch {
	case strings.HasPrefix(eventType, "expense."):
		return "expense"
	case strings.HasPrefix(eventType, "income."):
		return "income"
	case strings.HasPrefix(eventType, "recurring."):
		return "recurring_transaction"
	case strings.HasPrefix(eventType, "budget."):
		return "budget"
	case strings.HasPrefix(eventType, "savings_goal."):
		return "savings_goal"
	case eventType == "user.registered":
		return "user"
	}
	return ""
}
