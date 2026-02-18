package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Leer las 3 bases de datos
	newDB := os.Getenv("DATABASE_URL")
	usersDB := os.Getenv("USERS_DB_URL")
	gamificationDB := os.Getenv("GAMIFICATION_DB_URL")

	if newDB == "" || usersDB == "" || gamificationDB == "" {
		log.Fatal("❌ Faltan variables de entorno")
	}

	fmt.Println("========================================")
	fmt.Println("Comparación de Tablas entre Bases de Datos")
	fmt.Println("========================================")
	fmt.Println()

	// Función para obtener tablas
	getTables := func(dbURL, dbName string) map[string]int {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("❌ Error al conectar a %s: %v\n", dbName, err)
		}
		defer db.Close()

		rows, err := db.Query(`
			SELECT table_name, 
			       (SELECT COUNT(*) FROM information_schema.columns 
			        WHERE table_schema = 'public' AND table_name = t.table_name) as col_count
			FROM information_schema.tables t
			WHERE table_schema = 'public' 
			ORDER BY table_name
		`)
		if err != nil {
			log.Fatalf("❌ Error al obtener tablas de %s: %v\n", dbName, err)
		}
		defer rows.Close()

		tables := make(map[string]int)
		for rows.Next() {
			var tableName string
			var colCount int
			if err := rows.Scan(&tableName, &colCount); err != nil {
				continue
			}
			tables[tableName] = colCount
		}
		return tables
	}

	// Obtener tablas de cada base
	fmt.Println("📊 NUEVA BASE DE DATOS (Monolith)")
	newTables := getTables(newDB, "New DB")
	for table, cols := range newTables {
		fmt.Printf("  ✅ %-40s (%d columnas)\n", table, cols)
	}
	fmt.Printf("\nTotal: %d tablas\n\n", len(newTables))

	fmt.Println("📊 USERS DB (Legacy)")
	usersTables := getTables(usersDB, "Users DB")
	for table, cols := range usersTables {
		_, existsInNew := newTables[table]
		if existsInNew {
			fmt.Printf("  ✅ %-40s (%d columnas) [MIGRADA]\n", table, cols)
		} else {
			fmt.Printf("  ⚠️  %-40s (%d columnas) [NO MIGRADA]\n", table, cols)
		}
	}
	fmt.Printf("\nTotal: %d tablas\n\n", len(usersTables))

	fmt.Println("📊 GAMIFICATION DB (Legacy)")
	gamTables := getTables(gamificationDB, "Gamification DB")
	for table, cols := range gamTables {
		_, existsInNew := newTables[table]
		if existsInNew {
			fmt.Printf("  ✅ %-40s (%d columnas) [MIGRADA]\n", table, cols)
		} else {
			fmt.Printf("  ⚠️  %-40s (%d columnas) [NO MIGRADA]\n", table, cols)
		}
	}
	fmt.Printf("\nTotal: %d tablas\n\n", len(gamTables))

	// Resumen
	fmt.Println("========================================")
	fmt.Println("RESUMEN")
	fmt.Println("========================================")
	fmt.Printf("Tablas en Users DB (legacy):        %d\n", len(usersTables))
	fmt.Printf("Tablas en Gamification DB (legacy): %d\n", len(gamTables))
	fmt.Printf("Tablas en Nueva DB (monolith):      %d\n", len(newTables))
	fmt.Println()

	// Detectar tablas no migradas
	notMigrated := []string{}
	for table := range usersTables {
		if _, exists := newTables[table]; !exists {
			notMigrated = append(notMigrated, table+" (Users DB)")
		}
	}
	for table := range gamTables {
		if _, exists := newTables[table]; !exists {
			notMigrated = append(notMigrated, table+" (Gamification DB)")
		}
	}

	if len(notMigrated) > 0 {
		fmt.Println("⚠️  TABLAS NO MIGRADAS:")
		for _, table := range notMigrated {
			fmt.Printf("   - %s\n", table)
		}
	} else {
		fmt.Println("✅ Todas las tablas fueron migradas")
	}
	fmt.Println()
}
