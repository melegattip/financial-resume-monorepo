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

	// budgets — monolith uses period_start/period_end (not start_date/end_date)
	c, s, err = copyTableGeneric[SrcBudget](gamDB, targetDB, log, dryRun, "budgets",
		`INSERT INTO budgets (id, user_id, category_id, amount, spent_amount, period, period_start, period_end, alert_at, status, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, 0, ?, ?, ?, 0.8, 'on_track', ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcBudget) []interface{} {
			return []interface{}{r.ID, r.UserID, r.CategoryID, r.Amount, r.Period, r.StartDate, r.EndDate, r.IsActive, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy budgets: %w", err)
	}
	copied["budgets"] = c
	skipped["budgets"] = s

	// recurring_transactions — monolith uses type/next_date (not transaction_type/next_execution)
	c, s, err = copyTableGeneric[SrcRecurringTransaction](gamDB, targetDB, log, dryRun, "recurring_transactions",
		`INSERT INTO recurring_transactions (id, user_id, category_id, amount, description, type, frequency,
		 end_date, next_date, is_active, auto_create, notify_before, execution_count, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, false, 0, 0, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcRecurringTransaction) []interface{} {
			return []interface{}{r.ID, r.UserID, r.CategoryID, r.Amount, r.Description, r.TransactionType, r.Frequency,
				r.EndDate, r.NextExecution, r.IsActive, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy recurring_transactions: %w", err)
	}
	copied["recurring_transactions"] = c
	skipped["recurring_transactions"] = s

	// savings_goals — monolith uses target_date (not deadline); status derived from achieved/is_active
	c, s, err = copyTableGeneric[SrcSavingsGoal](gamDB, targetDB, log, dryRun, "savings_goals",
		`INSERT INTO savings_goals (id, user_id, name, description, target_amount, current_amount,
		 category, priority, target_date, status, monthly_target, weekly_target, daily_target,
		 progress, remaining_amount, days_remaining, is_auto_save, auto_save_amount, auto_save_frequency,
		 image_url, achieved_at, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, 0, 0, ?, ?, 0, false, 0, '', '', ?, ?, ?)
		 ON CONFLICT (id) DO NOTHING`,
		func(r SrcSavingsGoal) []interface{} {
			status := "active"
			if r.Achieved {
				status = "completed"
			} else if !r.IsActive {
				status = "paused"
			}
			progress := 0.0
			if r.TargetAmount > 0 {
				progress = (r.CurrentAmount / r.TargetAmount) * 100
			}
			remaining := r.TargetAmount - r.CurrentAmount
			if remaining < 0 {
				remaining = 0
			}
			return []interface{}{r.ID, r.UserID, r.Name, r.Description, r.TargetAmount, r.CurrentAmount,
				r.Category, r.Priority, r.Deadline, status, progress, remaining, r.AchievedAt, r.CreatedAt, r.UpdatedAt}
		})
	if err != nil {
		return copied, skipped, fmt.Errorf("copy savings_goals: %w", err)
	}
	copied["savings_goals"] = c
	skipped["savings_goals"] = s

	// savings_transactions — monolith uses type (not transaction_type) and requires user_id (via JOIN)
	stCopied, stSkipped, err := copySavingsTransactions(gamDB, targetDB, log, dryRun)
	if err != nil {
		return copied, skipped, fmt.Errorf("copy savings_transactions: %w", err)
	}
	copied["savings_transactions"] = stCopied
	skipped["savings_transactions"] = stSkipped

	log.Info().Msg("financial data migration complete")
	return copied, skipped, nil
}

// copySavingsTransactions copies savings_transactions joining with savings_goals
// to obtain the user_id, which the monolith schema requires.
func copySavingsTransactions(sourceDB, targetDB *gorm.DB, log zerolog.Logger, dryRun bool) (int64, int64, error) {
	type srcRow struct {
		ID              string    `gorm:"column:id"`
		GoalID          string    `gorm:"column:goal_id"`
		UserID          string    `gorm:"column:user_id"`
		Amount          float64   `gorm:"column:amount"`
		TransactionType string    `gorm:"column:transaction_type"`
		Description     string    `gorm:"column:description"`
		CreatedAt       time.Time `gorm:"column:created_at"`
	}

	var rows []srcRow
	err := sourceDB.Raw(`
		SELECT st.id, st.goal_id, sg.user_id, st.amount,
		       COALESCE(st.type, 'deposit') AS transaction_type,
		       st.description, st.created_at
		FROM savings_transactions st
		JOIN savings_goals sg ON sg.id = st.goal_id
	`).Scan(&rows).Error
	if err != nil {
		return 0, 0, fmt.Errorf("read savings_transactions with user_id: %w", err)
	}

	log.Info().Str("table", "savings_transactions").Int("source_count", len(rows)).Msg("read from gamification-db")

	if dryRun {
		log.Info().Str("table", "savings_transactions").Int("count", len(rows)).Msg("[DRY RUN] would copy")
		return int64(len(rows)), 0, nil
	}

	var copied, skipped int64
	for _, r := range rows {
		result := targetDB.Exec(`
			INSERT INTO savings_transactions (id, goal_id, user_id, amount, type, description, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.GoalID, r.UserID, r.Amount, r.TransactionType, r.Description, r.CreatedAt)
		if result.Error != nil {
			return copied, skipped, fmt.Errorf("insert savings_transaction %s: %w", r.ID, result.Error)
		}
		if result.RowsAffected == 0 {
			skipped++
		} else {
			copied++
		}
	}

	log.Info().Str("table", "savings_transactions").Int64("copied", copied).Int64("skipped", skipped).Msg("copy complete")
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
