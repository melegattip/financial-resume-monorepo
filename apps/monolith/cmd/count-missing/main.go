package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	gamificationDB := os.Getenv("GAMIFICATION_DB_URL")
	if gamificationDB == "" {
		log.Fatal("❌ GAMIFICATION_DB_URL no está configurada")
	}

	db, err := sql.Open("postgres", gamificationDB)
	if err != nil {
		log.Fatalf("❌ Error: %v\n", err)
	}
	defer db.Close()

	fmt.Println("========================================")
	fmt.Println("Datos Financieros en Gamification DB")
	fmt.Println("========================================")
	fmt.Println()

	tables := []string{
		"expenses",
		"incomes",
		"budgets",
		"categories",
		"recurring_transactions",
		"recurring_transaction_executions",
		"savings_goals",
		"savings_transactions",
		"transaction_models",
	}

	fmt.Printf("%-40s %12s\n", "TABLA", "REGISTROS")
	fmt.Println("--------------------------------------------------------")

	total := 0
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			fmt.Printf("%-40s %12s\n", table, "ERROR")
			continue
		}
		fmt.Printf("%-40s %12d\n", table, count)
		total += count
	}

	fmt.Println("--------------------------------------------------------")
	fmt.Printf("%-40s %12d\n", "TOTAL REGISTROS NO MIGRADOS", total)
	fmt.Println()
	fmt.Println("========================================")
}
