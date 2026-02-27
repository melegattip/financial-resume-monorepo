package database

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// openGormMock creates a GORM DB backed by go-sqlmock.
// Note: postgres.Config{Conn: rawDB} does NOT run SELECT VERSION() at Open,
// so no VERSION expectation is needed.
func openGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: rawDB}), &gorm.Config{})
	require.NoError(t, err)
	t.Cleanup(func() { rawDB.Close() })
	return gormDB, mock
}

// ─── MigrateRLSPolicies ───────────────────────────────────────────────────────

func TestMigrateRLSPolicies_ExecutesForAllTables(t *testing.T) {
	gormDB, mock := openGormMock(t)
	logger := zerolog.Nop()

	// For each tenant-scoped table we expect two Exec calls:
	//   1. ALTER TABLE ... ENABLE ROW LEVEL SECURITY
	//   2. DO $$ BEGIN ... CREATE POLICY ... END $$
	for _, table := range tenantScopedTables {
		mock.ExpectExec("ALTER TABLE " + table + " ENABLE ROW LEVEL SECURITY").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DO").
			WillReturnResult(sqlmock.NewResult(0, 0))
	}

	MigrateRLSPolicies(gormDB, logger)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMigrateRLSPolicies_ContinuesOnAlterError(t *testing.T) {
	gormDB, mock := openGormMock(t)
	logger := zerolog.Nop()

	// First table: ALTER TABLE fails — should log warning and continue
	mock.ExpectExec("ALTER TABLE " + tenantScopedTables[0] + " ENABLE ROW LEVEL SECURITY").
		WillReturnError(errors.New("syntax error"))
	mock.ExpectExec("DO").WillReturnResult(sqlmock.NewResult(0, 0))

	// Remaining tables succeed
	for _, table := range tenantScopedTables[1:] {
		mock.ExpectExec("ALTER TABLE " + table + " ENABLE ROW LEVEL SECURITY").
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec("DO").WillReturnResult(sqlmock.NewResult(0, 0))
	}

	// Should not panic
	assert.NotPanics(t, func() { MigrateRLSPolicies(gormDB, logger) })
}

func TestMigrateRLSPolicies_TablesCovered(t *testing.T) {
	// Ensure the expected tables are in the list — catches accidental removals.
	required := []string{
		"expenses", "incomes", "categories",
		"budgets", "savings_goals", "savings_transactions",
		"recurring_transactions", "audit_logs",
	}
	for _, want := range required {
		found := false
		for _, got := range tenantScopedTables {
			if got == want {
				found = true
				break
			}
		}
		assert.True(t, found, "table %q should be in tenantScopedTables", want)
	}
}

// ─── WithTenant ───────────────────────────────────────────────────────────────

func TestWithTenant_SetsLocalVariableAndCommits(t *testing.T) {
	gormDB, mock := openGormMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("SET LOCAL app.current_tenant").
		WithArgs("tnt_abc").
		WillReturnResult(sqlmock.NewResult(0, 0))
	// Simulate the user's fn querying a table
	mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("e1"))
	mock.ExpectCommit()

	called := false
	err := WithTenant(gormDB, "tnt_abc", func(tx *gorm.DB) error {
		called = true
		// Execute any query to consume the SELECT expectation
		var result struct{ ID string }
		return tx.Raw("SELECT id FROM expenses LIMIT 1").Scan(&result).Error
	})

	require.NoError(t, err)
	assert.True(t, called, "user function should be invoked")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTenant_RollsBackOnFnError(t *testing.T) {
	gormDB, mock := openGormMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("SET LOCAL app.current_tenant").
		WithArgs("tnt_xyz").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	fnErr := errors.New("something went wrong")
	err := WithTenant(gormDB, "tnt_xyz", func(_ *gorm.DB) error {
		return fnErr
	})

	assert.ErrorIs(t, err, fnErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTenant_RollsBackOnSetLocalError(t *testing.T) {
	gormDB, mock := openGormMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("SET LOCAL app.current_tenant").
		WithArgs("tnt_bad").
		WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	err := WithTenant(gormDB, "tnt_bad", func(_ *gorm.DB) error {
		return nil
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rls: set tenant context")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWithTenant_EmptyTenantID(t *testing.T) {
	gormDB, mock := openGormMock(t)

	// WithTenant should still set the variable (empty string) and execute fn
	mock.ExpectBegin()
	mock.ExpectExec("SET LOCAL app.current_tenant").
		WithArgs("").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	err := WithTenant(gormDB, "", func(_ *gorm.DB) error { return nil })
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
