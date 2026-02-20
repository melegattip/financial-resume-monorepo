//go:build ignore

package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func tryConnect(label, connStr string) {
	fmt.Printf("\nTrying: %s\n", label)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("  Open error: %v\n", err)
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Printf("  Ping error: %v\n", err)
		return
	}
	fmt.Printf("  SUCCESS\n")
	var count int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
	fmt.Printf("  Public tables count: %d\n", count)
}

func main() {
	// password decoded from %24c%408MPpN%409sSwX2: $c@8MPpN@9sSwX2
	// The @ in the password is the issue — need to use key=value DSN format

	tryConnect("key=value DSN port 6543",
		"host=aws-1-us-east-1.pooler.supabase.com port=6543 user=postgres.njgddjhjqzhhruklxrzg password=$c@8MPpN@9sSwX2 dbname=postgres sslmode=require")

	tryConnect("key=value DSN port 5432",
		"host=aws-1-us-east-1.pooler.supabase.com port=5432 user=postgres.njgddjhjqzhhruklxrzg password=$c@8MPpN@9sSwX2 dbname=postgres sslmode=require")
}
