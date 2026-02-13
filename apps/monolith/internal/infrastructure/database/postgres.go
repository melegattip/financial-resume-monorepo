package database

import (
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
)

// NewPostgresConnection creates a new GORM database connection with pooling.
func NewPostgresConnection(cfg *config.AppConfig, log zerolog.Logger) (*gorm.DB, error) {
	dsn := cfg.DSN()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Verify connectivity
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	log.Info().
		Int("max_idle", cfg.DBMaxIdleConns).
		Int("max_open", cfg.DBMaxOpenConns).
		Msg("database connected")

	return db, nil
}

// Close gracefully closes the database connection.
func Close(db *gorm.DB, log zerolog.Logger) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("failed to get sql.DB for close")
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close database")
	} else {
		log.Info().Msg("database connection closed")
	}
}
