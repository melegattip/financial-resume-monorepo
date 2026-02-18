package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("❌ DATABASE_URL no está configurada")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Error al abrir conexión: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Error al hacer ping: %v\n", err)
	}

	fmt.Println("========================================")
	fmt.Println("Conteo de Registros por Tabla")
	fmt.Println("========================================")
	fmt.Println()

	// Obtener todas las tablas
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		ORDER BY table_name
	`)
	if err != nil {
		log.Fatalf("❌ Error al obtener tablas: %v\n", err)
	}
	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("❌ Error al leer tabla: %v\n", err)
		}
		tables = append(tables, tableName)
	}

	// Contar registros en cada tabla
	fmt.Printf("%-35s %10s\n", "TABLA", "REGISTROS")
	fmt.Println("-----------------------------------------------")

	totalRecords := 0
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			fmt.Printf("%-35s %10s\n", table, "ERROR")
			continue
		}
		fmt.Printf("%-35s %10d\n", table, count)
		totalRecords += count
	}

	fmt.Println("-----------------------------------------------")
	fmt.Printf("%-35s %10d\n", "TOTAL", totalRecords)
	fmt.Println()
	fmt.Println("========================================")
}
