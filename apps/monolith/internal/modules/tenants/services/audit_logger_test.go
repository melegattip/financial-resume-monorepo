package services

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/modules/tenants/domain"
	sharedports "github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/shared/ports"
)

// ─── Test doubles ────────────────────────────────────────────────────────────

type mockAuditRepo struct {
	calls []domain.AuditLog
	err   error
}

func (m *mockAuditRepo) CreateAuditLog(_ context.Context, log domain.AuditLog) error {
	m.calls = append(m.calls, log)
	return m.err
}

func (m *mockAuditRepo) ListAuditLogs(_ context.Context, _ string, _, _ int) ([]domain.AuditLog, error) {
	return nil, nil
}

type mockEvent struct {
	eventType   string
	aggregateID string
	userID      string
}

func (e mockEvent) EventType() string   { return e.eventType }
func (e mockEvent) AggregateID() string { return e.aggregateID }
func (e mockEvent) UserID() string      { return e.userID }
func (e mockEvent) OccurredAt() string  { return "" }

var _ sharedports.Event = mockEvent{}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func openGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	// postgres.Config{Conn: rawDB} reuses the existing *sql.DB without
	// running any initialisation queries, so no VERSION() expectation needed.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: rawDB}), &gorm.Config{})
	require.NoError(t, err)

	t.Cleanup(func() { rawDB.Close() })
	return gormDB, mock
}

// ─── resolveTable ─────────────────────────────────────────────────────────────

func TestResolveTable(t *testing.T) {
	cases := []struct {
		eventType string
		want      string
	}{
		{"expense.created", "expenses"},
		{"expense.updated", "expenses"},
		{"expense.deleted", "expenses"},
		{"income.created", "incomes"},
		{"income.updated", "incomes"},
		{"income.deleted", "incomes"},
		{"recurring.created", "recurring_transactions"},
		{"recurring.updated", "recurring_transactions"},
		{"recurring.deleted", "recurring_transactions"},
		{"recurring.executed", "recurring_transactions"},
		{"recurring.paused", "recurring_transactions"},
		{"recurring.resumed", "recurring_transactions"},
		// user.registered resolves via tenant_members, not a table prefix
		{"user.registered", ""},
		{"budget.created", ""},
		{"unknown.event", ""},
	}

	for _, tc := range cases {
		t.Run(tc.eventType, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveTable(tc.eventType))
		})
	}
}

// ─── resolveEntityType ────────────────────────────────────────────────────────

func TestResolveEntityType(t *testing.T) {
	cases := []struct {
		eventType string
		want      string
	}{
		{"expense.created", "expense"},
		{"expense.updated", "expense"},
		{"expense.deleted", "expense"},
		{"income.created", "income"},
		{"income.updated", "income"},
		{"income.deleted", "income"},
		{"recurring.created", "recurring_transaction"},
		{"recurring.executed", "recurring_transaction"},
		{"recurring.paused", "recurring_transaction"},
		{"recurring.resumed", "recurring_transaction"},
		{"user.registered", "user"},
		{"budget.created", ""},
		{"unknown.event", ""},
	}

	for _, tc := range cases {
		t.Run(tc.eventType, func(t *testing.T) {
			assert.Equal(t, tc.want, resolveEntityType(tc.eventType))
		})
	}
}

// ─── Handle — happy paths ─────────────────────────────────────────────────────

func TestAuditLogger_Handle_Expense(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM expenses").
		WithArgs("exp_123").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_abc"))

	err := al.Handle(context.Background(), mockEvent{
		eventType:   "expense.created",
		aggregateID: "exp_123",
		userID:      "user_xyz",
	})

	require.NoError(t, err)
	require.Len(t, repo.calls, 1)

	entry := repo.calls[0]
	assert.Equal(t, "tnt_abc", entry.TenantID)
	assert.Equal(t, "user_xyz", entry.UserID)
	assert.Equal(t, "expense.created", entry.Action)
	assert.Equal(t, "expense", entry.EntityType)
	assert.Equal(t, "exp_123", entry.EntityID)
	assert.True(t, len(entry.ID) > 0, "ID should be set")
	assert.False(t, entry.CreatedAt.IsZero(), "CreatedAt should be set")
}

func TestAuditLogger_Handle_Income(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM incomes").
		WithArgs("inc_456").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_abc"))

	err := al.Handle(context.Background(), mockEvent{
		eventType:   "income.deleted",
		aggregateID: "inc_456",
		userID:      "u1",
	})

	require.NoError(t, err)
	require.Len(t, repo.calls, 1)
	assert.Equal(t, "income", repo.calls[0].EntityType)
	assert.Equal(t, "income.deleted", repo.calls[0].Action)
	assert.Equal(t, "tnt_abc", repo.calls[0].TenantID)
}

func TestAuditLogger_Handle_Recurring(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM recurring_transactions").
		WithArgs("rec_789").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_def"))

	err := al.Handle(context.Background(), mockEvent{
		eventType:   "recurring.executed",
		aggregateID: "rec_789",
		userID:      "u2",
	})

	require.NoError(t, err)
	require.Len(t, repo.calls, 1)
	assert.Equal(t, "recurring_transaction", repo.calls[0].EntityType)
	assert.Equal(t, "tnt_def", repo.calls[0].TenantID)
}

func TestAuditLogger_Handle_UserRegistered(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM tenant_members").
		WithArgs("user_new").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_personal"))

	err := al.Handle(context.Background(), mockEvent{
		eventType:   "user.registered",
		aggregateID: "user_new",
		userID:      "user_new",
	})

	require.NoError(t, err)
	require.Len(t, repo.calls, 1)
	assert.Equal(t, "user", repo.calls[0].EntityType)
	assert.Equal(t, "tnt_personal", repo.calls[0].TenantID)
	assert.Equal(t, "user.registered", repo.calls[0].Action)
}

func TestAuditLogger_Handle_AuditID_HasPrefix(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM expenses").
		WithArgs("e1").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_x"))

	_ = al.Handle(context.Background(), mockEvent{
		eventType: "expense.created", aggregateID: "e1", userID: "u1",
	})

	require.Len(t, repo.calls, 1)
	assert.True(t, len(repo.calls[0].ID) > 4, "ID should not be empty")
	assert.Contains(t, repo.calls[0].ID, "aud_", "ID should have aud_ prefix")
}

// ─── Handle — skip conditions ─────────────────────────────────────────────────

func TestAuditLogger_Handle_TenantNotFound_Skips(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	// Return empty result set — entity exists but tenant_id is null or entity not found
	mock.ExpectQuery("SELECT tenant_id FROM expenses").
		WithArgs("exp_gone").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}))

	err := al.Handle(context.Background(), mockEvent{
		eventType:   "expense.created",
		aggregateID: "exp_gone",
		userID:      "u1",
	})

	require.NoError(t, err)
	assert.Empty(t, repo.calls, "should not write audit log when tenant_id is not found")
}

func TestAuditLogger_Handle_UnknownEvent_Skips(t *testing.T) {
	gormDB, _ := openGormMock(t)
	repo := &mockAuditRepo{}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	// No DB expectation — unknown events should not hit the DB
	err := al.Handle(context.Background(), mockEvent{
		eventType:   "budget.created",
		aggregateID: "bud_1",
		userID:      "u1",
	})

	require.NoError(t, err)
	assert.Empty(t, repo.calls, "unknown event types should be silently skipped")
}

func TestAuditLogger_Handle_RepoError_NonFatal(t *testing.T) {
	gormDB, mock := openGormMock(t)
	repo := &mockAuditRepo{err: assert.AnError}
	al := NewAuditLogger(gormDB, repo, zerolog.Nop())

	mock.ExpectQuery("SELECT tenant_id FROM expenses").
		WithArgs("e1").
		WillReturnRows(sqlmock.NewRows([]string{"tenant_id"}).AddRow("tnt_x"))

	// Repo error must NOT propagate — audit is best-effort
	err := al.Handle(context.Background(), mockEvent{
		eventType:   "expense.created",
		aggregateID: "e1",
		userID:      "u1",
	})

	require.NoError(t, err, "repo error should be swallowed, not propagated")
}
