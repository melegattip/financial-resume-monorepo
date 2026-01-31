// Package config provides configuration utilities for the application
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/melegattip/financial-resume-engine/internal/core/domain"
	"github.com/melegattip/financial-resume-engine/internal/usecases/transactions"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB initializes and returns a database connection using environment variables
// It also performs necessary migrations and schema updates
func InitDB() *gorm.DB {
	// Get environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "financial_resume")
	dbPort := getEnvOrDefault("DB_PORT", "5432")

	// Build DSN with SSL mode configuration for production
	sslMode := getEnvOrDefault("DB_SSLMODE", "disable")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)

	// Pool mode configuration for connection pooling
	if poolMode := os.Getenv("DB_POOL_MODE"); poolMode != "" {
		log.Printf("🔧 [Database] Using pool mode: %s", poolMode)
	}

	// Log connection attempt for debugging
	log.Printf("🔗 [Database] Connecting to: %s:%s/%s (SSL: %s)", dbHost, dbPort, dbName, sslMode)

	// Configure GORM
	gormConfig := &gorm.Config{
		// Use silent logger in production to avoid spam
		Logger: logger.Default.LogMode(logger.Silent),
	}

	// Configure postgres driver with prepared statement options
	var postgresConfig postgres.Config
	if sslMode == "require" {
		// For cloud providers like Render/Supabase, disable prepared statements
		postgresConfig = postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // Disable prepared statements
		}
		log.Println("🔧 [Database] Disabled prepared statements for cloud deployment")
	} else {
		postgresConfig = postgres.Config{
			DSN: dsn,
		}
	}

	// Try to connect to the database
	db, err := gorm.Open(postgres.New(postgresConfig), gormConfig)
	if err != nil {
		log.Printf("❌ [Database] Connection failed with DSN: host=%s port=%s dbname=%s sslmode=%s", dbHost, dbPort, dbName, sslMode)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ [Database] Failed to get SQL DB instance: %v", err)
	}

	// Connection pool settings based on environment
	maxConnections := getEnvAsInt("DB_MAX_CONNECTIONS", 25)
	maxIdleConnections := maxConnections / 5 // 20% of max connections
	if maxIdleConnections < 2 {
		maxIdleConnections = 2
	}

	sqlDB.SetMaxOpenConns(maxConnections)
	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)

	log.Printf("✅ [Database] Connected successfully (Pool: %d max, %d idle)", maxConnections, maxIdleConnections)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ [Database] Ping failed: %v", err)
	}

	// Commented: Automatic migrations are problematic
	// This should be done manually or through migration scripts
	// if err := runMigrations(db); err != nil {
	// 	log.Printf("⚠️ [Database] Migration warnings: %v", err)
	// }

	return db
}

// runMigrations runs database migrations (commented for safety)
func runMigrations(db *gorm.DB) error {
	log.Println("🔄 [Database] Running migrations...")

	// Auto-migrate all domain models
	err := db.AutoMigrate(
		&transactions.TransactionModel{},
		&domain.Category{},
		&domain.Budget{},
		&domain.RecurringTransaction{},
		&domain.SavingsGoal{},
	)

	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ [Database] Migrations completed")
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt returns the value of an environment variable as integer or a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
