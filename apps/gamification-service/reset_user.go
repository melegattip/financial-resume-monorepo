package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Configuración de base de datos
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "financial_resume")
	sslmode := getEnv("DB_SSLMODE", "disable")

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		host, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	userID := "2" // nivel3@test.com

	log.Printf("🔄 Reseteando usuario %s...", userID)

	// Resetear achievement category_creator
	result, err := db.Exec(`
		UPDATE achievements 
		SET progress = 3, 
		    completed = false, 
		    unlocked_at = NULL,
		    updated_at = NOW()
		WHERE user_id = $1 AND type = 'category_creator'
	`, userID)

	if err != nil {
		log.Fatalf("Error updating achievement: %v", err)
	}
	
	rows, _ := result.RowsAffected()
	log.Printf("✅ Achievement reseteado: %d filas afectadas", rows)

	// Resetear XP del usuario
	result, err = db.Exec(`
		UPDATE user_gamification 
		SET total_xp = 208,
		    achievements_count = 2,
		    updated_at = NOW()
		WHERE user_id = $1
	`, userID)

	if err != nil {
		log.Fatalf("Error updating user gamification: %v", err)
	}
	
	rows, _ = result.RowsAffected()
	log.Printf("✅ Usuario reseteado: %d filas afectadas", rows)

	// Verificar estado
	var achievementType string
	var progress, target int
	var completed bool

	err = db.QueryRow("SELECT type, progress, target, completed FROM achievements WHERE user_id = $1 AND type = 'category_creator'", userID).Scan(&achievementType, &progress, &target, &completed)
	if err != nil {
		log.Printf("Error querying achievement: %v", err)
	}

	var totalXP, achievementsCount, currentLevel int
	err = db.QueryRow("SELECT total_xp, achievements_count, current_level FROM user_gamification WHERE user_id = $1", userID).Scan(&totalXP, &achievementsCount, &currentLevel)
	if err != nil {
		log.Printf("Error querying user: %v", err)
	}

	log.Printf("🎯 ESTADO FINAL:")
	log.Printf("   Achievement: %s %d/%d (completado: %v)", achievementType, progress, target, completed)
	log.Printf("   Usuario: %d XP, %d achievements, nivel %d", totalXP, achievementsCount, currentLevel)
	log.Printf("✅ Reset completado!")
}