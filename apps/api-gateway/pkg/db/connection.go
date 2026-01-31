package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	// Configuración para Supabase con transaction pooling
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require prefer_simple_protocol=true binary_parameters=yes disable_prepared_statements=true",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Configurar connection pool para Supabase
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	return db, db.Ping()
}
