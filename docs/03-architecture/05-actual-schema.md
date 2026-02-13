# Actual Database Schema (As-Built)

**Document Version**: 1.0
**Date**: 2026-02-13
**Status**: Reflects implemented code after Phase 1 + Phase 2 + Phase 3
**Source of truth**: GORM models in `apps/monolith/` + legacy `init.sql` + migration tool (`cmd/migrate`)

---

## Overview

Single consolidated database: **`financial_resume`** (PostgreSQL 15).

- **21 tables** total
- **All IDs**: `VARCHAR(255)` (string UUIDs)
- **All entity tables**: `deleted_at TIMESTAMP NULL` (GORM soft delete)
- **Auth tables** (4): managed by GORM AutoMigrate
- **Gamification tables** (6): managed by GORM AutoMigrate via migration tool
- **Financial tables** (11): created from legacy `init.sql`, modified by migration tool

---

## Entity-Relationship Diagram

```mermaid
erDiagram
    %% ============ AUTH (migrated from users_db) ============
    users {
        varchar_255 id PK
        varchar_255 email UK "NOT NULL"
        varchar_255 password_hash "NOT NULL"
        varchar_100 first_name "NOT NULL"
        varchar_100 last_name "NOT NULL"
        varchar_20 phone
        varchar_500 avatar
        boolean is_active "DEFAULT true"
        boolean is_verified "DEFAULT false"
        varchar_500 email_verification_token
        timestamp email_verification_expires
        varchar_500 password_reset_token
        timestamp password_reset_expires
        timestamp last_login
        int failed_login_attempts "DEFAULT 0"
        timestamp locked_until
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    user_preferences {
        varchar_255 user_id PK_FK
        varchar_10 currency "DEFAULT USD"
        varchar_10 language "DEFAULT en"
        varchar_10 theme "DEFAULT light"
        varchar_20 date_format "DEFAULT YYYY-MM-DD"
        varchar_50 timezone "DEFAULT UTC"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    user_notification_settings {
        varchar_255 user_id PK_FK
        boolean email_notifications "DEFAULT true"
        boolean push_notifications "DEFAULT true"
        boolean weekly_reports "DEFAULT true"
        boolean expense_alerts "DEFAULT true"
        boolean budget_alerts "DEFAULT true"
        boolean achievement_notifications "DEFAULT true"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    user_two_fa {
        varchar_255 user_id PK_FK
        varchar_255 secret "NOT NULL"
        boolean enabled "DEFAULT false"
        text_array backup_codes
        varchar_50 last_used_code
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ TRANSACTIONS ============
    categories {
        varchar_50 id PK
        varchar_255 user_id FK
        varchar_50 name "NOT NULL"
        varchar_10 icon "DEFAULT folder"
        varchar_7 color "DEFAULT #3B82F6"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    expenses {
        varchar_50 id PK
        varchar_255 user_id FK
        decimal amount "NOT NULL"
        decimal amount_paid "DEFAULT 0"
        text description
        varchar_50 category_id FK
        date due_date
        boolean paid "DEFAULT false"
        decimal percentage
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    incomes {
        varchar_36 id PK
        varchar_255 user_id FK
        decimal amount "NOT NULL"
        text description
        varchar_50 category_id FK
        varchar_50 source
        decimal percentage
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ BUDGETS ============
    budgets {
        varchar_50 id PK
        varchar_255 user_id FK
        varchar_50 category_id FK "NOT NULL"
        decimal amount "NOT NULL, GT 0"
        decimal spent_amount "DEFAULT 0"
        varchar_20 period "monthly|weekly|yearly"
        timestamp period_start "NOT NULL"
        timestamp period_end "NOT NULL"
        decimal alert_at "DEFAULT 0.80"
        varchar_20 status "on_track|warning|exceeded"
        boolean is_active "DEFAULT true"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    budget_notifications {
        varchar_50 id PK
        varchar_50 budget_id FK
        varchar_255 user_id FK
        varchar_30 notification_type "alert|exceeded|reset"
        text message "NOT NULL"
        decimal threshold_percentage
        decimal exceeded_amount
        boolean is_read "DEFAULT false"
        timestamp created_at
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ SAVINGS ============
    savings_goals {
        varchar_50 id PK
        varchar_255 user_id FK
        varchar_255 name "NOT NULL"
        text description
        decimal target_amount "NOT NULL, GT 0"
        decimal current_amount "DEFAULT 0"
        varchar_20 category "vacation|emergency|house|..."
        varchar_10 priority "high|medium|low"
        date target_date "NOT NULL"
        varchar_15 status "active|achieved|paused|cancelled"
        decimal monthly_target
        decimal weekly_target
        decimal daily_target
        decimal progress "0 to 1"
        decimal remaining_amount
        int days_remaining
        boolean is_auto_save "DEFAULT false"
        decimal auto_save_amount
        varchar_10 auto_save_frequency
        varchar_500 image_url
        timestamp created_at
        timestamp updated_at
        timestamp achieved_at
        timestamp deleted_at "SOFT DELETE"
    }

    savings_transactions {
        varchar_50 id PK
        varchar_50 goal_id FK
        varchar_255 user_id FK
        decimal amount "NOT NULL"
        varchar_10 type "deposit|withdrawal"
        varchar_255 description
        timestamp created_at
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ RECURRING ============
    recurring_transactions {
        varchar_36 id PK
        varchar_255 user_id FK
        decimal amount "NOT NULL, GT 0"
        varchar_255 description "NOT NULL"
        varchar_36 category_id FK
        varchar_10 type "income|expense"
        varchar_10 frequency "daily|weekly|monthly|yearly"
        date next_date "NOT NULL"
        timestamp last_executed
        boolean is_active "DEFAULT true"
        boolean auto_create "DEFAULT true"
        int notify_before "DEFAULT 1"
        date end_date
        int execution_count "DEFAULT 0"
        int max_executions
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    recurring_transaction_executions {
        varchar_36 id PK
        varchar_36 recurring_transaction_id FK
        varchar_255 user_id FK
        decimal amount "NOT NULL"
        varchar_255 description "NOT NULL"
        varchar_10 type "NOT NULL"
        timestamp executed_at
        boolean success "NOT NULL"
        varchar_36 created_transaction_id
        text error_message
        varchar_20 execution_method "DEFAULT automatic"
        timestamp deleted_at "SOFT DELETE"
    }

    recurring_transaction_notifications {
        varchar_36 id PK
        varchar_36 recurring_transaction_id FK
        varchar_255 user_id FK
        varchar_20 notification_type "NOT NULL"
        varchar_255 title "NOT NULL"
        text message "NOT NULL"
        timestamp sent_at
        varchar_20 delivery_status "DEFAULT pending"
        varchar_20 delivery_method "DEFAULT system"
        text error_message
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ GAMIFICATION (consolidated from gamification_db) ============
    user_gamification {
        varchar_255 id PK
        varchar_255 user_id UK_FK
        int financial_health_score "DEFAULT 1"
        int engagement_component "DEFAULT 0"
        int health_component "DEFAULT 0"
        int current_level "DEFAULT 1"
        int insights_viewed "DEFAULT 0"
        int actions_completed "DEFAULT 0"
        int achievements_count "DEFAULT 0"
        int current_streak "DEFAULT 0"
        timestamp last_activity
        timestamp last_score_calculation
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    achievements {
        varchar_255 id PK
        varchar_255 user_id FK
        varchar_100 type "NOT NULL"
        varchar_255 name "NOT NULL"
        text description
        int points "DEFAULT 0"
        int progress "DEFAULT 0"
        int target "NOT NULL"
        boolean completed "DEFAULT false"
        timestamp unlocked_at
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    user_actions {
        varchar_255 id PK
        varchar_255 user_id FK
        varchar_100 action_type "NOT NULL"
        varchar_100 entity_type "NOT NULL"
        varchar_255 entity_id
        int score_contribution "DEFAULT 0"
        varchar_50 component_affected
        text description
        timestamp created_at
        timestamp deleted_at "SOFT DELETE"
    }

    challenges {
        varchar_255 id PK
        varchar_100 challenge_key UK "NOT NULL"
        varchar_255 name "NOT NULL"
        text description "NOT NULL"
        varchar_50 challenge_type "NOT NULL"
        varchar_20 icon "DEFAULT target"
        int xp_reward "DEFAULT 0"
        varchar_100 requirement_type "NOT NULL"
        int requirement_target "NOT NULL"
        jsonb requirement_data
        boolean active "DEFAULT true"
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    user_challenges {
        varchar_255 id PK
        varchar_255 user_id FK
        varchar_255 challenge_id FK
        date challenge_date "NOT NULL"
        int progress "DEFAULT 0"
        int target "NOT NULL"
        boolean completed "DEFAULT false"
        timestamp completed_at
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at "SOFT DELETE"
    }

    challenge_progress_tracking {
        varchar_255 id PK
        varchar_255 user_id FK
        date challenge_date "NOT NULL"
        varchar_100 action_type "NOT NULL"
        varchar_100 entity_type
        int count "DEFAULT 1"
        jsonb unique_entities
        timestamp created_at
        timestamp deleted_at "SOFT DELETE"
    }

    %% ============ RELATIONSHIPS ============
    users ||--o| user_preferences : "1:1"
    users ||--o| user_notification_settings : "1:1"
    users ||--o| user_two_fa : "1:0..1"
    users ||--o| user_gamification : "1:0..1"
    users ||--o{ categories : "owns"
    users ||--o{ expenses : "owns"
    users ||--o{ incomes : "owns"
    users ||--o{ budgets : "owns"
    users ||--o{ savings_goals : "owns"
    users ||--o{ recurring_transactions : "owns"
    users ||--o{ achievements : "earns"
    users ||--o{ user_actions : "performs"
    users ||--o{ user_challenges : "participates"
    users ||--o{ challenge_progress_tracking : "tracks"

    categories ||--o{ expenses : "categorizes"
    categories ||--o{ incomes : "categorizes"
    categories ||--o{ budgets : "limits"

    budgets ||--o{ budget_notifications : "notifies"
    savings_goals ||--o{ savings_transactions : "tracks"
    recurring_transactions ||--o{ recurring_transaction_executions : "logs"
    recurring_transactions ||--o{ recurring_transaction_notifications : "notifies"
    challenges ||--o{ user_challenges : "assigned"
```

---

## Divergences from Target State (`02-target-state.md`)

The original target-state document describes an **idealized schema**. The actual implementation
diverges in several areas due to preserving backward compatibility with the legacy frontend
and existing data.

### Auth Tables

| Aspect | Target State | Actual |
|--------|-------------|--------|
| User name field | `name VARCHAR(255)` | `first_name VARCHAR(100)` + `last_name VARCHAR(100)` |
| User fields | Basic (id, email, password, name) | Extended: 2FA, lockout, email verification, avatar, phone |
| Auth tables | 2 (users, user_preferences) | 4 (+ user_notification_settings, user_two_fa) |
| Preferences | currency, language, timezone, gamification_enabled | currency, language, theme, date_format, timezone |

### Gamification Tables

| Aspect | Target State | Actual |
|--------|-------------|--------|
| XP tracking | `total_xp`, `longest_streak` | `financial_health_score`, `engagement_component`, `health_component` |
| Achievements model | Catalog-style: `condition_type`, `condition_value`, separate `user_achievements` junction | Per-user: `user_id`, `progress/target`, `completed`, no junction table |
| Extra tables | 3 (user_gamification, achievements, user_achievements) | 6 (+ user_actions, challenges, user_challenges, challenge_progress_tracking) |

### Financial Tables

| Aspect | Target State | Actual |
|--------|-------------|--------|
| Expenses | `payment_status` enum, `transaction_date`, `is_recurring`, `recurring_transaction_id` | `paid` boolean, `amount_paid`, `due_date`, `percentage` (no transaction_date) |
| Incomes | `source`, `transaction_date`, `is_recurring`, `recurring_transaction_id` | `source`, `category_id`, `percentage` (no transaction_date) |
| Categories | `type` (expense/income), `is_default` | `user_id` (per-user), `icon`, `color`, UNIQUE(user_id, name) |
| Budgets | `limit_amount`, `current_amount`, `start_date`, `end_date` | `amount`, `spent_amount`, `period_start`, `period_end`, `alert_at` |
| Recurring | `start_date`, `next_execution`, `last_execution` | `next_date`, `last_executed`, `auto_create`, `notify_before`, `execution_count`, `max_executions` |
| Extra tables | None | `budget_notifications`, `savings_transactions`, `recurring_transaction_executions`, `recurring_transaction_notifications` |

### Summary

The target-state should be updated in future phases to reflect these realities, or the
implementation should be adjusted to match the target. The key structural goals are met:

- **Single database**: `financial_resume`
- **All user_id**: `VARCHAR(255)` (string UUIDs)
- **Soft delete**: `deleted_at` on all 21 tables
- **Single connection pool**: 10 idle / 25 max open

---

## SQL Schema

### Auth Tables (GORM AutoMigrate)

```sql
CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    avatar VARCHAR(500),
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    email_verification_token VARCHAR(500),
    email_verification_expires TIMESTAMP,
    password_reset_token VARCHAR(500),
    password_reset_expires TIMESTAMP,
    last_login TIMESTAMP,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE user_preferences (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    language VARCHAR(10) NOT NULL DEFAULT 'en',
    theme VARCHAR(10) NOT NULL DEFAULT 'light',
    date_format VARCHAR(20) NOT NULL DEFAULT 'YYYY-MM-DD',
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE user_notification_settings (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    email_notifications BOOLEAN NOT NULL DEFAULT true,
    push_notifications BOOLEAN NOT NULL DEFAULT true,
    weekly_reports BOOLEAN NOT NULL DEFAULT true,
    expense_alerts BOOLEAN NOT NULL DEFAULT true,
    budget_alerts BOOLEAN NOT NULL DEFAULT true,
    achievement_notifications BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE user_two_fa (
    user_id VARCHAR(255) PRIMARY KEY REFERENCES users(id),
    secret VARCHAR(255) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    backup_codes TEXT[],
    last_used_code VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### Financial Tables (from legacy init.sql + migration tool)

```sql
CREATE TABLE categories (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(10) DEFAULT '📂',
    color VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE(user_id, name)
);

CREATE TABLE expenses (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    amount_paid DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    description TEXT,
    category_id VARCHAR(50) REFERENCES categories(id),
    due_date DATE,
    paid BOOLEAN DEFAULT FALSE,
    percentage DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE incomes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    category_id VARCHAR(50) REFERENCES categories(id),
    source VARCHAR(50),
    percentage DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE budgets (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    category_id VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    spent_amount DECIMAL(15,2) DEFAULT 0.00 CHECK (spent_amount >= 0),
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'weekly', 'yearly')),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    alert_at DECIMAL(3,2) DEFAULT 0.80,
    status VARCHAR(20) DEFAULT 'on_track' CHECK (status IN ('on_track', 'warning', 'exceeded')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE budget_notifications (
    id VARCHAR(50) PRIMARY KEY,
    budget_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    notification_type VARCHAR(30) NOT NULL CHECK (notification_type IN ('alert', 'exceeded', 'reset')),
    message TEXT NOT NULL,
    threshold_percentage DECIMAL(5,2),
    exceeded_amount DECIMAL(15,2),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE savings_goals (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    target_amount DECIMAL(15,2) NOT NULL CHECK (target_amount > 0),
    current_amount DECIMAL(15,2) DEFAULT 0.00 CHECK (current_amount >= 0),
    category VARCHAR(20) NOT NULL CHECK (category IN ('vacation', 'emergency', 'house', 'car', 'education', 'retirement', 'investment', 'other')),
    priority VARCHAR(10) NOT NULL DEFAULT 'medium' CHECK (priority IN ('high', 'medium', 'low')),
    target_date DATE NOT NULL,
    status VARCHAR(15) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'achieved', 'paused', 'cancelled')),
    monthly_target DECIMAL(15,2) DEFAULT 0.00,
    weekly_target DECIMAL(15,2) DEFAULT 0.00,
    daily_target DECIMAL(15,2) DEFAULT 0.00,
    progress DECIMAL(5,4) DEFAULT 0.0000,
    remaining_amount DECIMAL(15,2) DEFAULT 0.00,
    days_remaining INTEGER DEFAULT 0,
    is_auto_save BOOLEAN DEFAULT FALSE,
    auto_save_amount DECIMAL(15,2) DEFAULT 0.00,
    auto_save_frequency VARCHAR(10) CHECK (auto_save_frequency IN ('daily', 'weekly', 'monthly')),
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    achieved_at TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE savings_transactions (
    id VARCHAR(50) PRIMARY KEY,
    goal_id VARCHAR(50) NOT NULL REFERENCES savings_goals(id) ON DELETE CASCADE,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount != 0),
    type VARCHAR(10) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE recurring_transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    category_id VARCHAR(36),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    next_date DATE NOT NULL,
    last_executed TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    auto_create BOOLEAN DEFAULT TRUE,
    notify_before INTEGER DEFAULT 1,
    end_date DATE,
    execution_count INTEGER DEFAULT 0,
    max_executions INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE recurring_transaction_executions (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN NOT NULL,
    created_transaction_id VARCHAR(36),
    error_message TEXT,
    execution_method VARCHAR(20) DEFAULT 'automatic',
    deleted_at TIMESTAMP NULL
);

CREATE TABLE recurring_transaction_notifications (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    notification_type VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivery_status VARCHAR(20) DEFAULT 'pending',
    delivery_method VARCHAR(20) DEFAULT 'system',
    error_message TEXT,
    deleted_at TIMESTAMP NULL
);
```

### Gamification Tables (GORM AutoMigrate via migration tool)

```sql
CREATE TABLE user_gamification (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    financial_health_score INT DEFAULT 1,
    engagement_component INT DEFAULT 0,
    health_component INT DEFAULT 0,
    current_level INT DEFAULT 1,
    insights_viewed INT DEFAULT 0,
    actions_completed INT DEFAULT 0,
    achievements_count INT DEFAULT 0,
    current_streak INT DEFAULT 0,
    last_activity TIMESTAMP,
    last_score_calculation TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE achievements (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    points INT DEFAULT 0,
    progress INT DEFAULT 0,
    target INT NOT NULL,
    completed BOOLEAN DEFAULT false,
    unlocked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE user_actions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255),
    score_contribution INT DEFAULT 0,
    component_affected VARCHAR(50),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE challenges (
    id VARCHAR(255) PRIMARY KEY,
    challenge_key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    challenge_type VARCHAR(50) NOT NULL,
    icon VARCHAR(20) DEFAULT '🎯',
    xp_reward INT DEFAULT 0,
    requirement_type VARCHAR(100) NOT NULL,
    requirement_target INT NOT NULL,
    requirement_data JSONB,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE user_challenges (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL,
    progress INT DEFAULT 0,
    target INT NOT NULL,
    completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE challenge_progress_tracking (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    count INT DEFAULT 1,
    unique_entities JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### Partial Indexes (created by migration tool)

```sql
CREATE INDEX idx_expenses_user_date_active ON expenses(user_id, transaction_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_expenses_category_active ON expenses(category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_budgets_user_active_soft ON budgets(user_id, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_recurring_next_exec_active ON recurring_transactions(next_date, is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_gamification_user_active ON user_gamification(user_id) WHERE deleted_at IS NULL;
```

### Triggers (from legacy init.sql)

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;

CREATE TRIGGER update_user_gamification_updated_at BEFORE UPDATE ON user_gamification FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_achievements_updated_at BEFORE UPDATE ON achievements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_recurring_transactions_updated_at BEFORE UPDATE ON recurring_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_savings_goals_updated_at BEFORE UPDATE ON savings_goals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## References

- **Target state (vision)**: [02-target-state.md](./02-target-state.md) -- idealized schema the roadmap targets
- **Legacy schemas**: `apps/api-gateway/scripts/init.sql`, `apps/users-service/scripts/user_init.sql`, `apps/gamification-service/scripts/init.sql`
- **GORM models (auth)**: `apps/monolith/internal/modules/auth/domain/`
- **GORM models (migration)**: `apps/monolith/internal/migration/models.go`
- **Migration tool**: `apps/monolith/cmd/migrate/`
