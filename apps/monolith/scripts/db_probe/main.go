package main

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

type DBTarget struct {
	Name string
	URL  string
}

func decodeURL(rawURL string) string {
	decoded, err := url.QueryUnescape(rawURL)
	if err != nil {
		return rawURL
	}
	return decoded
}

func probeDB(target DBTarget) {
	fmt.Printf("\n========================================\n")
	fmt.Printf("DATABASE: %s\n", target.Name)
	fmt.Printf("========================================\n")

	connStr := decodeURL(target.URL)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("ERROR opening connection: %v\n", err)
		return
	}
	defer db.Close()

	// Test ping
	if err := db.Ping(); err != nil {
		fmt.Printf("ERROR pinging database: %v\n", err)
		return
	}
	fmt.Printf("Connection: OK\n\n")

	// List tables
	fmt.Println("--- Tables in public schema ---")
	rows, err := db.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`)
	if err != nil {
		fmt.Printf("ERROR listing tables: %v\n", err)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			fmt.Printf("ERROR scanning row: %v\n", err)
			continue
		}
		tables = append(tables, tableName)
		fmt.Printf("  - %s\n", tableName)
	}

	if len(tables) == 0 {
		fmt.Println("  (no tables found)")
		return
	}

	// Row counts
	fmt.Println("\n--- Row counts ---")
	for _, table := range tables {
		var count int64
		query := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"`, table)
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			fmt.Printf("  %-35s  ERROR: %v\n", table, err)
		} else {
			fmt.Printf("  %-35s  %d rows\n", table, count)
		}
	}
}

func main() {
	targets := []DBTarget{
		{
			Name: "supabase-monolith (njgddjhjqzhhruklxrzg) - NEW TARGET",
			URL:  "postgresql://postgres.njgddjhjqzhhruklxrzg:%24c%408MPpN%409sSwX2@aws-1-us-east-1.pooler.supabase.com:6543/postgres?sslmode=require",
		},
		{
			Name: "supabase-users (akngrdpnwboujagnziqb) - LEGACY SOURCE",
			URL:  "postgresql://postgres.akngrdpnwboujagnziqb:%24c%408MPpN%409sSwX2@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require",
		},
		{
			Name: "supabase-gamification (gtzkqlbkqgnaittehfey) - LEGACY SOURCE",
			URL:  "postgresql://postgres.gtzkqlbkqgnaittehfey:%24c%408MPpN%409sSwX2@aws-0-sa-east-1.pooler.supabase.com:6543/postgres?sslmode=require",
		},
	}

	fmt.Println("Supabase Database Probe")
	fmt.Println("=======================")

	for _, target := range targets {
		probeDB(target)
	}

	fmt.Println("\n========================================")
	fmt.Println("Probe complete.")
	fmt.Println("========================================")
}
