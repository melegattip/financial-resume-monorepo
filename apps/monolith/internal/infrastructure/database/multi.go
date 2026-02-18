package database

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/melegattip/financial-resume-monorepo/apps/monolith/internal/infrastructure/config"
)

// MultiConnection holds connections to all three databases used during migration.
type MultiConnection struct {
	Target       *gorm.DB // financial_resume (main-db)
	UsersDB      *gorm.DB // users_db (users-service)
	Gamification *gorm.DB // financial_gamification (gamification-db)
}

// isSupabaseConnection checks if a connection string is for Supabase (pooler)
func isSupabaseConnection(connStr string) bool {
	return strings.Contains(connStr, "pooler.supabase.com") || 
		   strings.Contains(connStr, ":6543") ||
		   strings.Contains(connStr, "supabase.com")
}

// openConnection opens a GORM connection with Supabase-compatible settings if needed
func openConnection(connStr string, log zerolog.Logger) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	var postgresCfg postgres.Config
	
	if isSupabaseConnection(connStr) {
		// Supabase pooler doesn't support prepared statements
		// Parse URL to add prefer_simple_protocol parameter
		if u, err := url.Parse(connStr); err == nil {
			query := u.Query()
			query.Set("prefer_simple_protocol", "true")
			u.RawQuery = query.Encode()
			connStr = u.String()
		}
		
		postgresCfg = postgres.Config{
			DSN:                  connStr,
			PreferSimpleProtocol: true, // Disable prepared statements for Supabase
		}
		log.Debug().Msg("using Supabase-compatible connection (prepared statements disabled)")
	} else {
		postgresCfg = postgres.Config{
			DSN: connStr,
		}
	}

	db, err := gorm.Open(postgres.New(postgresCfg), gormCfg)
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

// NewMultiConnection opens connections to the target, users, and gamification databases.
// The target database uses the standard DATABASE_URL / DB_* config.
// Source databases use USERS_DB_URL and GAMIFICATION_DB_URL.
func NewMultiConnection(cfg *config.AppConfig, log zerolog.Logger) (*MultiConnection, error) {
	// Target database (financial_resume)
	targetDB, err := openConnection(cfg.DSN(), log)
	if err != nil {
		return nil, fmt.Errorf("connect target db: %w", err)
	}
	if err := pingDB(targetDB); err != nil {
		return nil, fmt.Errorf("ping target db: %w", err)
	}
	log.Info().Msg("connected to target database (financial_resume)")

	// Users source database
	if cfg.UsersDBURL == "" {
		return nil, fmt.Errorf("USERS_DB_URL environment variable is required for migration")
	}
	usersDB, err := openConnection(cfg.UsersDBURL, log)
	if err != nil {
		return nil, fmt.Errorf("connect users db: %w", err)
	}
	if err := pingDB(usersDB); err != nil {
		return nil, fmt.Errorf("ping users db: %w", err)
	}
	log.Info().Msg("connected to source database (users_db)")

	// Gamification source database
	if cfg.GamificationDBURL == "" {
		return nil, fmt.Errorf("GAMIFICATION_DB_URL environment variable is required for migration")
	}
	gamificationDB, err := openConnection(cfg.GamificationDBURL, log)
	if err != nil {
		return nil, fmt.Errorf("connect gamification db: %w", err)
	}
	if err := pingDB(gamificationDB); err != nil {
		return nil, fmt.Errorf("ping gamification db: %w", err)
	}
	log.Info().Msg("connected to source database (financial_gamification)")

	return &MultiConnection{
		Target:       targetDB,
		UsersDB:      usersDB,
		Gamification: gamificationDB,
	}, nil
}

// Close closes all database connections.
func (mc *MultiConnection) Close(log zerolog.Logger) {
	for _, pair := range []struct {
		name string
		db   *gorm.DB
	}{
		{"target", mc.Target},
		{"users_db", mc.UsersDB},
		{"gamification_db", mc.Gamification},
	} {
		sqlDB, err := pair.db.DB()
		if err != nil {
			log.Error().Err(err).Str("db", pair.name).Msg("failed to get sql.DB for close")
			continue
		}
		if err := sqlDB.Close(); err != nil {
			log.Error().Err(err).Str("db", pair.name).Msg("failed to close database")
		} else {
			log.Info().Str("db", pair.name).Msg("database connection closed")
		}
	}
}

func pingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
