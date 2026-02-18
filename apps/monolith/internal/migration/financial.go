package migration

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// --- Financial Data Migration (Phase 4) ---

// CopyFinancialData copies all financial tables from gamification-db to target
func CopyFinancialData(gamDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (map[string]int64, map[string]int64, error) {
	copied := make(map[string]int64)
	skipped := make(map[string]int64)

	// categories
	c, s, err := copyTableGeneric[SrcCategory](gamDB, targetDB, log, dryRun, "categories",
		`INSERT INTO categories (id, user_id, name, color, icon, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcCategory) []interface{} {
			return []interface{}{r.ID, r.UserID, r.Name, r.Color, r.Icon, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy categories: %w", err)
	}
	copied["categories"] = c
	skipped["categories"] = s

	// expenses
	c, s, err = copyTableGeneric[SrcExpense](gamDB, targetDB, log, dryRun, "expenses",
		`INSERT INTO expenses (id, user_id, category_id, amount, description, transaction_date, payment_method, notes, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcExpense) []interface{} {
			return []interface{}{r.ID, r.UserID, r.CategoryID, r.Amount, r.Description, r.TransactionDate, r.PaymentMethod, r.Notes, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy expenses: %w", err)
	}
	copied["expenses"] = c
	skipped["expenses"] = s

	// incomes
	c, s, err = copyTableGeneric[SrcIncome](gamDB, targetDB, log, dryRun, "incomes",
		`INSERT INTO incomes (id, user_id, amount, source, description, received_date, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcIncome) []interface{} {
			return []interface{}{r.ID, r.UserID, r.Amount, r.Source, r.Description, r.ReceivedDate, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy incomes: %w", err)
	}
	copied["incomes"] = c
	skipped["incomes"] = s

	// budgets
	c, s, err = copyTableGeneric[SrcBudget](gamDB, targetDB, log, dryRun, "budgets",
		`INSERT INTO budgets (id, user_id, category_id, amount, period, start_date, end_date, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcBudget) []interface{} {
			return []interface{}{r.ID, r.UserID, r.CategoryID, r.Amount, r.Period, r.StartDate, r.EndDate, r.IsActive, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy budgets: %w", err)
	}
	copied["budgets"] = c
	skipped["budgets"] = s

	// recurring_transactions
	c, s, err = copyTableGeneric[SrcRecurringTransaction](gamDB, targetDB, log, dryRun, "recurring_transactions",
		`INSERT INTO recurring_transactions (id, user_id, category_id, amount, description, transaction_type, frequency, 
		 day_of_month, day_of_week, start_date, end_date, next_execution, is_active, payment_method, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcRecurringTransaction) []interface{} {
			return []interface{}{r.ID, r.UserID, r.CategoryID, r.Amount, r.Description, r.TransactionType, r.Frequency,
				r.DayOfMonth, r.DayOfWeek, r.StartDate, r.EndDate, r.NextExecution, r.IsActive, r.PaymentMethod, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy recurring_transactions: %w", err)
	}
	copied["recurring_transactions"] = c
	skipped["recurring_transactions"] = s

	// savings_goals
	c, s, err = copyTableGeneric[SrcSavingsGoal](gamDB, targetDB, log, dryRun, "savings_goals",
		`INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount, deadline, 
		 priority, category, icon, color, is_active, achieved, achieved_at, created_at, updated_at, created_by, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcSavingsGoal) []interface{} {
			return []interface{}{r.ID, r.UserID, r.Name, r.Description, r.TargetAmount, r.CurrentAmount, r.Deadline,
				r.Priority, r.Category, r.Icon, r.Color, r.IsActive, r.Achieved, r.AchievedAt, r.CreatedAt, r.UpdatedAt, r.CreatedBy, r.UpdatedBy}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy savings_goals: %w", err)
	}
	copied["savings_goals"] = c
	skipped["savings_goals"] = s

	// savings_transactions
	c, s, err = copyTableGeneric[SrcSavingsTransaction](gamDB, targetDB, log, dryRun, "savings_transactions",
		`INSERT INTO savings_transactions (id, goal_id, amount, transaction_type, description, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcSavingsTransaction) []interface{} {
			return []interface{}{r.ID, r.GoalID, r.Amount, r.TransactionType, r.Description, r.CreatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy savings_transactions: %w", err)
	}
	copied["savings_transactions"] = c
	skipped["savings_transactions"] = s

	log.Info().Msg("financial data migration complete")
	return copied, skipped, nil
}

// Source structs for financial tables (from gamification-db)

type SrcCategory struct {
	ID        string    `gorm:"column:id"`
	UserID    string    `gorm:"column:user_id"`
	Name      string    `gorm:"column:name"`
	Color     string    `gorm:"column:color"`
	Icon      string    `gorm:"column:icon"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (SrcCategory) TableName() string { return "categories" }

type SrcExpense struct {
	ID              string    `gorm:"column:id"`
	UserID          string    `gorm:"column:user_id"`
	CategoryID      string    `gorm:"column:category_id"`
	Amount          float64   `gorm:"column:amount"`
	Description     string    `gorm:"column:description"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	PaymentMethod   string    `gorm:"column:payment_method"`
	Notes           string    `gorm:"column:notes"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (SrcExpense) TableName() string { return "expenses" }

type SrcIncome struct {
	ID           string    `gorm:"column:id"`
	UserID       string    `gorm:"column:user_id"`
	Amount       float64   `gorm:"column:amount"`
	Source       string    `gorm:"column:source"`
	Description  string    `gorm:"column:description"`
	ReceivedDate time.Time `gorm:"column:received_date"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (SrcIncome) TableName() string { return "incomes" }

type SrcBudget struct {
	ID         string    `gorm:"column:id"`
	UserID     string    `gorm:"column:user_id"`
	CategoryID string    `gorm:"column:category_id"`
	Amount     float64   `gorm:"column:amount"`
	Period     string    `gorm:"column:period"`
	StartDate  time.Time `gorm:"column:start_date"`
	EndDate    time.Time `gorm:"column:end_date"`
	IsActive   bool      `gorm:"column:is_active"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (SrcBudget) TableName() string { return "budgets" }

type SrcRecurringTransaction struct {
	ID              string     `gorm:"column:id"`
	UserID          string     `gorm:"column:user_id"`
	CategoryID      string     `gorm:"column:category_id"`
	Amount          float64    `gorm:"column:amount"`
	Description     string     `gorm:"column:description"`
	TransactionType string     `gorm:"column:transaction_type"`
	Frequency       string     `gorm:"column:frequency"`
	DayOfMonth      *int       `gorm:"column:day_of_month"`
	DayOfWeek       *string    `gorm:"column:day_of_week"`
	StartDate       time.Time  `gorm:"column:start_date"`
	EndDate         *time.Time `gorm:"column:end_date"`
	NextExecution   time.Time  `gorm:"column:next_execution"`
	IsActive        bool       `gorm:"column:is_active"`
	PaymentMethod   string     `gorm:"column:payment_method"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
}

func (SrcRecurringTransaction) TableName() string { return "recurring_transactions" }

type SrcSavingsGoal struct {
	ID            string     `gorm:"column:id"`
	UserID        string     `gorm:"column:user_id"`
	Name          string     `gorm:"column:name"`
	Description   string     `gorm:"column:description"`
	TargetAmount  float64    `gorm:"column:target_amount"`
	CurrentAmount float64    `gorm:"column:current_amount"`
	Deadline      time.Time  `gorm:"column:deadline"`
	Priority      string     `gorm:"column:priority"`
	Category      string     `gorm:"column:category"`
	Icon          string     `gorm:"column:icon"`
	Color         string     `gorm:"column:color"`
	IsActive      bool       `gorm:"column:is_active"`
	Achieved      bool       `gorm:"column:achieved"`
	AchievedAt    *time.Time `gorm:"column:achieved_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	CreatedBy     string     `gorm:"column:created_by"`
	UpdatedBy     string     `gorm:"column:updated_by"`
}

func (SrcSavingsGoal) TableName() string { return "savings_goals" }

type SrcSavingsTransaction struct {
	ID              string    `gorm:"column:id"`
	GoalID          string    `gorm:"column:goal_id"`
	Amount          float64   `gorm:"column:amount"`
	TransactionType string    `gorm:"column:transaction_type"`
	Description     string    `gorm:"column:description"`
	CreatedAt       time.Time `gorm:"column:created_at"`
}

func (SrcSavingsTransaction) TableName() string { return "savings_transactions" }
