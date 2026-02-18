package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Leer DATABASE_URL desde variables de entorno o .env.production
	dbURL := os.Getenv("DATABASE_URL")
	
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL no está configurada. Ejecuta: set DATABASE_URL=...")
	}

	fmt.Println("========================================")
	fmt.Println("Verificación de Conexión a Supabase")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("🔍 Intentando conectar...")

	// Intentar conectar a la base de datos
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Error al abrir conexión: %v\n", err)
	}
	defer db.Close()

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Error al hacer ping a la base de datos: %v\n", err)
	}

	fmt.Println("✅ Conexión exitosa a Supabase!")
	fmt.Println()

	// Obtener información de la base de datos
	var dbName, version string
	err = db.QueryRow("SELECT current_database(), version()").Scan(&dbName, &version)
	if err != nil {
		log.Fatalf("❌ Error al consultar información: %v\n", err)
	}

	fmt.Printf("📊 Base de datos: %s\n", dbName)
	fmt.Printf("📦 Versión PostgreSQL: %s\n", version)
	fmt.Println()

	// Verificar tablas existentes
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount)
	if err != nil {
		log.Fatalf("❌ Error al contar tablas: %v\n", err)
	}

	if tableCount == 0 {
		fmt.Println("✅ Base de datos nueva (sin tablas en schema 'public')")
	} else {
		fmt.Printf("⚠️  Hay %d tablas en la base de datos\n", tableCount)
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("✅ Verificación completada exitosamente")
	fmt.Println("========================================")
}
