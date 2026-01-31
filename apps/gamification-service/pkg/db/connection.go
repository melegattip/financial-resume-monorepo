package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// Config configuración de base de datos
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewConnection crea una nueva conexión a PostgreSQL optimizada para Supabase
func NewConnection(config Config) (*sql.DB, error) {
	var psqlInfo string

	// Detectar si es entorno de producción/Supabase
	isSupabase := strings.Contains(config.Host, "supabase") ||
		strings.Contains(config.Host, "amazonaws.com") ||
		os.Getenv("DATABASE_URL") != "" ||
		os.Getenv("ENVIRONMENT") == "production"

	// Verificar si existe DATABASE_URL (común en Render)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		log.Printf("🔗 Using DATABASE_URL for connection")

		// Agregar parámetros específicos para Supabase/pooling solo si es necesario
		if u, err := url.Parse(databaseURL); err == nil {
			query := u.Query()
			if isSupabase {
				// Configuración específica para Supabase con transaction pooling
				query.Set("prefer_simple_protocol", "true")
				query.Set("sslmode", "require")
				query.Set("binary_parameters", "yes")            // Evitar prepared statements
				query.Set("disable_prepared_statements", "true") // Deshabilitar prepared statements
				log.Printf("🔧 Optimized connection string for Supabase transaction pooling")
			} else {
				// Para desarrollo local, usar configuración básica
				query.Set("sslmode", config.SSLMode)
				log.Printf("🔧 Using local development configuration")
			}
			query.Set("connect_timeout", "30")
			query.Set("application_name", "financial-gamification-service")
			u.RawQuery = query.Encode()
			databaseURL = u.String()
		} else {
			// Si falla el parsing, agregar manualmente
			separator := "?"
			if strings.Contains(databaseURL, "?") {
				separator = "&"
			}
			if isSupabase {
				databaseURL += separator + "prefer_simple_protocol=true&sslmode=require&connect_timeout=30&application_name=financial-gamification-service"
				log.Printf("🔧 Manually added Supabase optimization parameters")
			} else {
				databaseURL += separator + fmt.Sprintf("sslmode=%s&connect_timeout=30&application_name=financial-gamification-service", config.SSLMode)
				log.Printf("🔧 Manually added local development parameters")
			}
		}

		psqlInfo = databaseURL

		// Parse DATABASE_URL para logging (sin mostrar password)
		if u, err := url.Parse(databaseURL); err == nil {
			log.Printf("🔗 Connecting to database: %s:%s/%s", u.Hostname(), u.Port(), u.Path[1:])
		}
	} else {
		// Usar variables individuales con optimizaciones condicionales
		if isSupabase {
			psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s prefer_simple_protocol=true binary_parameters=yes disable_prepared_statements=true connect_timeout=30 application_name=financial-gamification-service",
				config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
			log.Printf("🔧 Using optimized configuration for Supabase transaction pooling")
		} else {
			psqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=30 application_name=financial-gamification-service",
				config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
			log.Printf("🔧 Using standard configuration for local PostgreSQL")
		}
		log.Printf("🔗 Connecting to database: %s:%d/%s", config.Host, config.Port, config.DBName)
	}

	// Log the final connection string (sin password)
	logSafeURL := strings.ReplaceAll(psqlInfo, config.Password, "***")
	log.Printf("🔍 Final connection string: %s", logSafeURL)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configurar connection pool optimizado para Supabase con transaction pooling
	if isSupabase {
		// Configuración más conservadora para Supabase transaction pooling
		db.SetMaxOpenConns(5)                  // Muy reducido para transaction pooling
		db.SetMaxIdleConns(2)                  // Mínimo de conexiones idle
		db.SetConnMaxLifetime(2 * time.Minute) // Renovar conexiones muy frecuentemente
		db.SetConnMaxIdleTime(1 * time.Minute) // Cerrar conexiones idle muy rápido
		log.Printf("🔧 Supabase connection pool configured: MaxOpen=5, MaxIdle=2, LifeTime=2m, IdleTime=1m")
	} else {
		// Configuración estándar para desarrollo local
		db.SetMaxOpenConns(10)                 // Reducido para evitar saturar el pool
		db.SetMaxIdleConns(5)                  // Mantener menos conexiones idle
		db.SetConnMaxLifetime(5 * time.Minute) // Renovar conexiones más frecuentemente
		db.SetConnMaxIdleTime(2 * time.Minute) // Cerrar conexiones idle más rápido
		log.Printf("🔧 Local connection pool configured: MaxOpen=10, MaxIdle=5, LifeTime=5m, IdleTime=2m")
	}

	// Verificar conexión con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("🔍 Testing database connection...")
	log.Printf("🔧 Attempting database ping...")

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Printf("✅ Database ping successful!")
	log.Printf("✅ Connected to database successfully")

	return db, nil
}

// ParseDatabaseURL parsea una DATABASE_URL y devuelve un Config
func ParseDatabaseURL(databaseURL string) (Config, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return Config{}, fmt.Errorf("error parsing DATABASE_URL: %w", err)
	}

	var port int
	if u.Port() != "" {
		port, err = strconv.Atoi(u.Port())
		if err != nil {
			return Config{}, fmt.Errorf("error parsing port: %w", err)
		}
	} else {
		port = 5432 // default PostgreSQL port
	}

	password, _ := u.User.Password()

	// Determinar SSL mode
	sslMode := "disable"
	if q := u.Query(); q.Get("sslmode") != "" {
		sslMode = q.Get("sslmode")
	}

	return Config{
		Host:     u.Hostname(),
		Port:     port,
		User:     u.User.Username(),
		Password: password,
		DBName:   u.Path[1:], // Remove leading slash
		SSLMode:  sslMode,
	}, nil
}

// GetDefaultConfig retorna configuración por defecto
func GetDefaultConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBName:   "gamification_db",
		SSLMode:  "disable",
	}
}
