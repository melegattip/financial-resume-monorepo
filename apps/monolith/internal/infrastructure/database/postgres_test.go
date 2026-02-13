package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openTestGorm(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Allow any queries during GORM init (SELECT VERSION(), etc.)
	mock.MatchExpectationsInOrder(false)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           "sqlmock",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		// If GORM init fails with mock, skip the test
		t.Skipf("GORM could not init with sqlmock: %v", err)
	}

	t.Cleanup(func() { mockDB.Close() })
	return gormDB, mock
}

func TestClose_Success(t *testing.T) {
	gormDB, mock := openTestGorm(t)

	mock.ExpectClose()

	logger := zerolog.Nop()
	Close(gormDB, logger)

	assert.NoError(t, mock.ExpectationsWereMet())
}
