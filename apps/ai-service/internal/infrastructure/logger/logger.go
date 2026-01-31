package logger

import (
	"log"
	"os"
)

// Setup configura el logger del servicio
func Setup() {
	// Configurar formato de log con timestamp
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Usar stdout para logs
	log.SetOutput(os.Stdout)

	log.Println("📋 Logger configured successfully")
}
