package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/melegattip/financial-gamification-service/internal/core/usecases"
	"github.com/melegattip/financial-gamification-service/internal/handlers"
	"github.com/melegattip/financial-gamification-service/internal/infrastructure/repository"
	"github.com/melegattip/financial-gamification-service/pkg/db"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func main() {
	// Configuración de base de datos desde variables de entorno
	dbConfig := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "financial_resume"), // ✅ MISMA BD QUE EL ENGINE PRINCIPAL
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	log.Printf("🔗 Connecting to database: %s:%d/%s", dbConfig.Host, dbConfig.Port, dbConfig.DBName)
	log.Printf("🎯 Database config: Host=%s, Port=%d, User=%s, DBName=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.DBName)

	// Conectar a base de datos
	database, err := db.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer database.Close()

	// Inicializar repositorio
	gamificationRepo := repository.NewGamificationRepository(database)

	// Inicializar casos de uso
	gamificationUseCase := usecases.NewGamificationUseCase(gamificationRepo)

	// Inicializar handlers
	gamificationHandlers := handlers.NewGamificationHandlers(gamificationUseCase)

	// Configurar router
	router := mux.NewRouter()

	// ✅ CORS middleware DEBE ir ANTES que las rutas
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Caller-ID")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// ✅ Ahora registrar las rutas DESPUÉS del middleware CORS
	gamificationHandlers.RegisterRoutes(router)

	// Health check endpoints
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			log.Printf("Error writing health check response: %v", err)
		}
	}).Methods("GET")

	router.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Ready")); err != nil {
			log.Printf("Error writing readiness check response: %v", err)
		}
	}).Methods("GET")

	// Iniciar servidor
	port := ":" + getEnv("PORT", "8081")
	log.Printf("🚀 Gamification Service starting on port %s", port)
	log.Printf("📊 Health check: http://localhost%s/health", port)
	log.Printf("🎮 API endpoints: http://localhost%s/api/v1/gamification", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
