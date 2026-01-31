package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/financial-ai-service/internal/adapters/cache"
	"github.com/financial-ai-service/internal/adapters/http/handlers"
	"github.com/financial-ai-service/internal/adapters/openai"
	"github.com/financial-ai-service/internal/core/usecases"
	"github.com/financial-ai-service/internal/infrastructure/config"
	"github.com/financial-ai-service/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar logger
	logger.Setup()
	log.Println("🤖 Starting Financial AI Service...")

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	// Inicializar dependencias
	openaiClient := openai.NewClient(cfg.OpenAI.APIKey, cfg.OpenAI.UseMock)
	cacheClient := cache.NewRedisClient(cfg.Redis.URL)

	// Inicializar casos de uso
	analysisUseCase := usecases.NewAnalysisUseCase(openaiClient, cacheClient)
	purchaseUseCase := usecases.NewPurchaseUseCase(openaiClient, cacheClient)
	creditUseCase := usecases.NewCreditUseCase(openaiClient, cacheClient)

	// Configurar router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Configurar handlers
	handlers.SetupRoutes(router, analysisUseCase, purchaseUseCase, creditUseCase)

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("🚀 Financial AI Service listening on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Configurar graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Financial AI Service...")

	// Crear contexto con timeout para shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown del servidor
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("❌ Error during server shutdown: %v", err)
	}

	// Cerrar conexiones
	if err := cacheClient.Close(); err != nil {
		log.Printf("❌ Error closing cache client: %v", err)
	}

	log.Println("✅ Financial AI Service stopped gracefully")
}
