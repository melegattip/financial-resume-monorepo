// cmd/reconcile/main.go
// Backfills total_xp and current_level in user_gamification by computing XP
// from the actual migrated transactions (expenses, incomes, categories, etc.).
// Also diagnoses/fixes any user ID mismatches between legacy (integer) and
// UUID accounts.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// XP amounts matching domain/user_action.go XPForAction table.
const (
	xpCreateExpense  = 8
	xpCreateIncome   = 8
	xpCreateCategory = 10
	xpCreateBudget   = 20
	xpSavingsGoal    = 15
	xpDepositSavings = 8
)

// xpThresholds matches domain/user_gamification.go — cumulative XP per level.
var xpThresholds = []int{0, 75, 200, 400, 700, 1200, 1800, 2600, 3600, 5500}

func calcLevel(xp int) int {
	level := 1
	for i, threshold := range xpThresholds {
		if xp >= threshold {
			level = i + 1
		}
	}
	if level > 10 {
		level = 10
	}
	return level
}

type UserRow struct {
	ID    string `gorm:"column:id"`
	Email string `gorm:"column:email"`
}

type XPSummary struct {
	UserID           string `gorm:"column:user_id"`
	ExpenseCount     int    `gorm:"column:expense_count"`
	IncomeCount      int    `gorm:"column:income_count"`
	CategoryCount    int    `gorm:"column:category_count"`
	BudgetCount      int    `gorm:"column:budget_count"`
	SavingsGoalCount int    `gorm:"column:savings_goal_count"`
	SavingsTxCount   int    `gorm:"column:savings_tx_count"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	dryRun := len(os.Args) > 1 && os.Args[1] == "--dry-run"
	if dryRun {
		fmt.Print("[DRY RUN] — no changes will be made\n")
	}

	// 1. Compute XP per user from migrated transactions
	var summaries []XPSummary
	err = db.Raw(`
		SELECT
			u.id AS user_id,
			COALESCE((SELECT COUNT(*) FROM expenses   WHERE user_id = u.id AND deleted_at IS NULL), 0) AS expense_count,
			COALESCE((SELECT COUNT(*) FROM incomes    WHERE user_id = u.id AND deleted_at IS NULL), 0) AS income_count,
			COALESCE((SELECT COUNT(*) FROM categories WHERE user_id = u.id AND deleted_at IS NULL), 0) AS category_count,
			COALESCE((SELECT COUNT(*) FROM budgets    WHERE user_id = u.id AND deleted_at IS NULL), 0) AS budget_count,
			COALESCE((SELECT COUNT(*) FROM savings_goals WHERE user_id = u.id AND deleted_at IS NULL), 0) AS savings_goal_count,
			COALESCE((SELECT COUNT(*) FROM savings_transactions WHERE user_id = u.id), 0) AS savings_tx_count
		FROM users u
		ORDER BY u.id
	`).Scan(&summaries).Error
	if err != nil {
		log.Fatalf("compute XP summaries: %v", err)
	}

	fmt.Println("=== XP Backfill Preview ===")
	fmt.Printf("%-6s  %-5s  %-5s  %-5s  %-5s  %-5s  %-5s  %-6s  %-5s\n",
		"userID", "exp", "inc", "cat", "bud", "sav", "stx", "totalXP", "level")

	type backfillRow struct {
		userID  string
		totalXP int
		level   int
	}
	var backfills []backfillRow

	for _, s := range summaries {
		xp := s.ExpenseCount*xpCreateExpense +
			s.IncomeCount*xpCreateIncome +
			s.CategoryCount*xpCreateCategory +
			s.BudgetCount*xpCreateBudget +
			s.SavingsGoalCount*xpSavingsGoal +
			s.SavingsTxCount*xpDepositSavings

		level := calcLevel(xp)
		backfills = append(backfills, backfillRow{s.UserID, xp, level})

		fmt.Printf("%-6s  %-5d  %-5d  %-5d  %-5d  %-5d  %-5d  %-6d  %-5d\n",
			s.UserID, s.ExpenseCount, s.IncomeCount, s.CategoryCount,
			s.BudgetCount, s.SavingsGoalCount, s.SavingsTxCount, xp, level)
	}

	fmt.Println()

	if dryRun {
		fmt.Println("Run without --dry-run to apply changes.")
		return
	}

	// 2. Apply backfill
	fmt.Println("=== Applying XP Backfill ===")
	for _, b := range backfills {
		result := db.Exec(`
			UPDATE user_gamification
			SET total_xp = ?, current_level = ?, updated_at = NOW()
			WHERE user_id = ?`,
			b.totalXP, b.level, b.userID)
		if result.Error != nil {
			fmt.Printf("  ERROR user_id=%s: %v\n", b.userID, result.Error)
		} else if result.RowsAffected == 0 {
			fmt.Printf("  SKIP  user_id=%s (no user_gamification record found)\n", b.userID)
		} else {
			fmt.Printf("  OK    user_id=%-6s  xp=%-6d  level=%d\n", b.userID, b.totalXP, b.level)
		}
	}

	fmt.Println()
	fmt.Println("=== Final State ===")
	printGamificationState(db)
}

func printGamificationState(db *gorm.DB) {
	type UGRow struct {
		UserID       string `gorm:"column:user_id"`
		TotalXP      int    `gorm:"column:total_xp"`
		CurrentLevel int    `gorm:"column:current_level"`
	}
	var rows []UGRow
	db.Raw("SELECT user_id, total_xp, current_level FROM user_gamification ORDER BY total_xp DESC").Scan(&rows)
	fmt.Printf("%-6s  %-8s  %-5s\n", "userID", "total_xp", "level")
	for _, r := range rows {
		fmt.Printf("%-6s  %-8d  %-5d\n", r.UserID, r.TotalXP, r.CurrentLevel)
	}
}
