# Main Database Schema - Current State (Production)

**Database**: `financial_resume` (PostgreSQL Main)
**Location**: Production (postgres-main, port 5432)
**Last Updated**: 2026-02-09
**Status**: ⚠️ PRODUCTION - Any changes require migration plan

## Overview

This document describes the **current production schema** for the main financial database. This schema handles:
- Financial transactions (income, expenses)
- Categories
- Budgets
- Savings goals
- Recurring transactions
- Gamification data (user_gamification, achievements, user_actions)

⚠️ **Note**: Gamification tables are duplicated between main-db and gamification-db (see [gamification-db.md](./gamification-db.md))

---

## Entity Relationship Diagram

```
┌─────────────┐
│  categories │────┐
└─────────────┘    │
                   │
       ┌───────────┼───────────────────────────┐
       │           │                           │
       ▼           ▼                           ▼
┌──────────┐  ┌──────────┐           ┌─────────────────────┐
│ incomes  │  │ expenses │           │ recurring_transactions│
└──────────┘  └──────────┘           └─────────────────────┘
                   │                           │
                   │                           ▼
                   │                 ┌──────────────────────────┐
                   │                 │ recurring_transaction_   │
                   │                 │    executions            │
                   │                 └──────────────────────────┘
                   │
                   ▼                 ┌──────────────────────────┐
           ┌──────────────┐          │ recurring_transaction_   │
           │   budgets    │          │    notifications         │
           └──────────────┘          └──────────────────────────┘
                   │
                   ▼
         ┌──────────────────┐
         │ budget_notifications│
         └──────────────────┘

┌──────────────────┐
│ savings_goals    │
└──────────────────┘
         │
         ▼
┌─────────────────────┐
│ savings_transactions│
└─────────────────────┘

┌────────────────────┐
│ user_gamification  │
└────────────────────┘
         │
         ├────────────────────┐
         ▼                    ▼
┌───────────────┐   ┌────────────────┐
│ achievements  │   │ user_actions   │
└───────────────┘   └────────────────┘
```

---

## Table Schemas

### 1. Categories

**Purpose**: User-defined categories for organizing transactions

```sql
CREATE TABLE categories (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    icon VARCHAR(10) DEFAULT '📂',
    color VARCHAR(7) DEFAULT '#3B82F6',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);
```

**Indexes**:
- `idx_incomes_category_id` on `categories(id)`
- `idx_expenses_category_id` on `categories(id)`

**Business Rules**:
- Each user can have multiple categories
- Category names must be unique per user
- Default icon and color provided if not specified

---

### 2. Incomes

**Purpose**: Track all income transactions

```sql
CREATE TABLE incomes (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    category_id VARCHAR(50),
    source VARCHAR(50),
    percentage DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

**Indexes**:
- `idx_incomes_user_id` on `incomes(user_id)`
- `idx_incomes_category_id` on `incomes(category_id)`

**Business Rules**:
- Amount must be positive (enforced at application layer)
- `percentage` calculated as `(amount / total_income) * 100` for reports
- Soft delete not implemented (physical delete only)

---

### 3. Expenses

**Purpose**: Track all expense transactions with payment tracking

```sql
CREATE TABLE expenses (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    amount_paid DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    description TEXT,
    category_id VARCHAR(50),
    due_date DATE,
    paid BOOLEAN DEFAULT FALSE,
    percentage DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

**Indexes**:
- `idx_expenses_user_id` on `expenses(user_id)`
- `idx_expenses_category_id` on `expenses(category_id)`
- `idx_expenses_due_date` on `expenses(due_date)`
- `idx_expenses_paid` on `expenses(paid)`
- `idx_expenses_amount_paid` on `expenses(amount_paid)`
- `idx_expenses_paid_status` on `expenses(paid, amount_paid)`

**Business Rules**:
- `amount_paid` <= `amount` (enforced at application layer)
- `paid` = TRUE when `amount_paid >= amount`
- `due_date` can be in the future (for budgeting purposes)
- `percentage` calculated as `(amount / total_income) * 100` for reports

**Known Issues**:
- No partial payment history (only current `amount_paid`)
- No tracking of payment method

---

### 4. Budgets

**Purpose**: Set spending limits per category for a time period

```sql
CREATE TABLE budgets (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    category_id VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    spent_amount DECIMAL(15,2) DEFAULT 0.00 CHECK (spent_amount >= 0),
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'weekly', 'yearly')),
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    alert_at DECIMAL(3,2) DEFAULT 0.80 CHECK (alert_at >= 0 AND alert_at <= 1),
    status VARCHAR(20) DEFAULT 'on_track' CHECK (status IN ('on_track', 'warning', 'exceeded')),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_budgets_user_id` on `budgets(user_id)`
- `idx_budgets_category_id` on `budgets(category_id)`
- `idx_budgets_period` on `budgets(period)`
- `idx_budgets_status` on `budgets(status)`
- `idx_budgets_active` on `budgets(is_active)`
- `idx_budgets_period_dates` on `budgets(period_start, period_end)`
- `idx_budgets_user_active` on `budgets(user_id, is_active)`

**Business Rules**:
- Status calculation:
  - `on_track`: `spent_amount / amount < alert_at` (default <80%)
  - `warning`: `spent_amount / amount >= alert_at AND < 1.0` (80-99%)
  - `exceeded`: `spent_amount / amount >= 1.0` (100%+)
- Automatic period rollover not implemented (requires cron job)
- `spent_amount` updated manually via application logic

---

### 5. Budget Notifications

**Purpose**: Track budget alert notifications sent to users

```sql
CREATE TABLE budget_notifications (
    id VARCHAR(50) PRIMARY KEY,
    budget_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    notification_type VARCHAR(30) NOT NULL CHECK (notification_type IN ('alert', 'exceeded', 'reset')),
    message TEXT NOT NULL,
    threshold_percentage DECIMAL(5,2),
    exceeded_amount DECIMAL(15,2),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_budget_notifications_user_id` on `budget_notifications(user_id)`
- `idx_budget_notifications_budget_id` on `budget_notifications(budget_id)`
- `idx_budget_notifications_type` on `budget_notifications(notification_type)`
- `idx_budget_notifications_read` on `budget_notifications(is_read)`
- `idx_budget_notifications_created` on `budget_notifications(created_at)`

---

### 6. Recurring Transactions

**Purpose**: Manage recurring income/expense templates

```sql
CREATE TABLE recurring_transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    description VARCHAR(255) NOT NULL,
    category_id VARCHAR(36),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    frequency VARCHAR(10) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    next_date DATE NOT NULL,
    last_executed TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    auto_create BOOLEAN DEFAULT TRUE,
    notify_before INTEGER DEFAULT 1 CHECK (notify_before >= 0),
    end_date DATE,
    execution_count INTEGER DEFAULT 0 CHECK (execution_count >= 0),
    max_executions INTEGER CHECK (max_executions > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_recurring_user_id` on `recurring_transactions(user_id)`
- `idx_recurring_next_date` on `recurring_transactions(next_date)`
- `idx_recurring_active` on `recurring_transactions(is_active)`
- `idx_recurring_type` on `recurring_transactions(type)`
- `idx_recurring_frequency` on `recurring_transactions(frequency)`
- `idx_recurring_category` on `recurring_transactions(category_id)`
- `idx_recurring_user_active` on `recurring_transactions(user_id, is_active)`
- `idx_recurring_pending_execution` on `recurring_transactions(is_active, next_date, auto_create)`

**Business Rules**:
- `auto_create = TRUE`: Automatically create transaction on `next_date`
- `notify_before`: Days before execution to send notification (default 1 day)
- `end_date`: Optional end date for finite recurring transactions
- `max_executions`: Optional limit on number of executions
- Deactivated when `execution_count >= max_executions` OR `next_date > end_date`

---

### 7. Recurring Transaction Executions

**Purpose**: Audit trail for executed recurring transactions

```sql
CREATE TABLE recurring_transaction_executions (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN NOT NULL,
    created_transaction_id VARCHAR(36),
    error_message TEXT,
    execution_method VARCHAR(20) DEFAULT 'automatic'
);
```

**Indexes**:
- `idx_execution_recurring_id` on `recurring_transaction_executions(recurring_transaction_id)`
- `idx_execution_user_id` on `recurring_transaction_executions(user_id)`
- `idx_execution_date` on `recurring_transaction_executions(executed_at)`
- `idx_execution_success` on `recurring_transaction_executions(success)`
- `idx_execution_created_transaction` on `recurring_transaction_executions(created_transaction_id)`

**Known Issues**:
- No foreign key to `created_transaction_id` (incomes/expenses tables)
- `execution_method` not validated by CHECK constraint

---

### 8. Recurring Transaction Notifications

**Purpose**: Track notifications for upcoming recurring transactions

```sql
CREATE TABLE recurring_transaction_notifications (
    id VARCHAR(36) PRIMARY KEY,
    recurring_transaction_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    notification_type VARCHAR(20) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivery_status VARCHAR(20) DEFAULT 'pending',
    delivery_method VARCHAR(20) DEFAULT 'system',
    error_message TEXT
);
```

**Indexes**:
- `idx_notification_recurring_id` on `recurring_transaction_notifications(recurring_transaction_id)`
- `idx_notification_user_id` on `recurring_transaction_notifications(user_id)`
- `idx_notification_type` on `recurring_transaction_notifications(notification_type)`
- `idx_notification_status` on `recurring_transaction_notifications(delivery_status)`
- `idx_notification_date` on `recurring_transaction_notifications(sent_at)`

---

### 9. Savings Goals

**Purpose**: Track user savings goals with progress and targets

```sql
CREATE TABLE savings_goals (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
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
    progress DECIMAL(5,4) DEFAULT 0.0000 CHECK (progress >= 0 AND progress <= 1),
    remaining_amount DECIMAL(15,2) DEFAULT 0.00,
    days_remaining INTEGER DEFAULT 0,
    is_auto_save BOOLEAN DEFAULT FALSE,
    auto_save_amount DECIMAL(15,2) DEFAULT 0.00,
    auto_save_frequency VARCHAR(10) CHECK (auto_save_frequency IN ('daily', 'weekly', 'monthly')),
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    achieved_at TIMESTAMP
);
```

**Indexes**:
- `idx_savings_goals_user_id` on `savings_goals(user_id)`
- `idx_savings_goals_status` on `savings_goals(status)`
- `idx_savings_goals_category` on `savings_goals(category)`
- `idx_savings_goals_priority` on `savings_goals(priority)`
- `idx_savings_goals_target_date` on `savings_goals(target_date)`
- `idx_savings_goals_user_status` on `savings_goals(user_id, status)`
- `idx_savings_goals_user_category` on `savings_goals(user_id, category)`
- `idx_savings_goals_auto_save` on `savings_goals(is_auto_save, auto_save_frequency)`

**Calculated Fields** (computed in application layer):
- `progress = current_amount / target_amount`
- `remaining_amount = target_amount - current_amount`
- `days_remaining = (target_date - NOW())`
- `daily_target = remaining_amount / days_remaining`
- `weekly_target = daily_target * 7`
- `monthly_target = remaining_amount / (days_remaining / 30)`

**Business Rules**:
- Status automatically changed to `achieved` when `current_amount >= target_amount`
- `achieved_at` timestamp set when goal is achieved
- Auto-save requires cron job to execute (not implemented in DB)

---

### 10. Savings Transactions

**Purpose**: Track deposits and withdrawals for savings goals

```sql
CREATE TABLE savings_transactions (
    id VARCHAR(50) PRIMARY KEY,
    goal_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    amount DECIMAL(15,2) NOT NULL CHECK (amount != 0),
    type VARCHAR(10) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (goal_id) REFERENCES savings_goals(id) ON DELETE CASCADE
);
```

**Indexes**:
- `idx_savings_transactions_goal_id` on `savings_transactions(goal_id)`
- `idx_savings_transactions_user_id` on `savings_transactions(user_id)`
- `idx_savings_transactions_type` on `savings_transactions(type)`
- `idx_savings_transactions_created` on `savings_transactions(created_at)`
- `idx_savings_transactions_goal_created` on `savings_transactions(goal_id, created_at)`

**Business Rules**:
- `amount > 0` for deposits, `amount < 0` for withdrawals (NOT enforced in DB, should be `amount > 0` always)
- CASCADE delete when parent savings goal is deleted
- No validation that withdrawal amount <= current savings (application layer)

---

### 11. User Gamification

**Purpose**: Track user XP, level, and engagement metrics

```sql
CREATE TABLE user_gamification (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    total_xp INTEGER NOT NULL DEFAULT 0,
    current_level INTEGER NOT NULL DEFAULT 0,
    insights_viewed INTEGER NOT NULL DEFAULT 0,
    actions_completed INTEGER NOT NULL DEFAULT 0,
    achievements_count INTEGER NOT NULL DEFAULT 0,
    current_streak INTEGER NOT NULL DEFAULT 0,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_user_gamification_user_id` on `user_gamification(user_id)`
- `idx_user_gamification_total_xp` on `user_gamification(total_xp DESC)`
- `idx_user_gamification_level` on `user_gamification(current_level DESC)`
- `idx_user_gamification_last_activity` on `user_gamification(last_activity DESC)`

**⚠️ Data Duplication Issue**:
This table exists in **BOTH** main-db and gamification-db. Potential sync issues.

---

### 12. Achievements

**Purpose**: Track user achievements/badges

```sql
CREATE TABLE achievements (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,
    progress INTEGER NOT NULL DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    unlocked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);
```

**Indexes**:
- `idx_achievements_user_id` on `achievements(user_id)`
- `idx_achievements_type` on `achievements(type)`
- `idx_achievements_completed` on `achievements(completed)`
- `idx_achievements_user_type` on `achievements(user_id, type)`

**⚠️ Data Duplication Issue**:
This table exists in **BOTH** main-db and gamification-db.

---

### 13. User Actions

**Purpose**: Audit trail for user actions that earn XP

```sql
CREATE TABLE user_actions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(100) NOT NULL,
    xp_earned INTEGER NOT NULL DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);
```

**Indexes**:
- `idx_user_actions_user_id` on `user_actions(user_id)`
- `idx_user_actions_action_type` on `user_actions(action_type)`
- `idx_user_actions_entity_type` on `user_actions(entity_type)`
- `idx_user_actions_created_at` on `user_actions(created_at DESC)`
- `idx_user_actions_user_date` on `user_actions(user_id, created_at DESC)`

**⚠️ Data Duplication Issue**:
This table exists in **BOTH** main-db and gamification-db.

---

## Triggers

### Auto-update `updated_at`

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language plpgsql;
```

**Applied to**:
- `user_gamification`
- `achievements`
- `budgets`
- `recurring_transactions`
- `savings_goals`

---

## Known Issues & Technical Debt

### 1. Gamification Data Duplication
**Problem**: `user_gamification`, `achievements`, `user_actions` exist in both main-db and gamification-db
**Impact**: Potential data inconsistency, sync issues
**Recommendation**: Consolidate to one DB (gamification-db or eliminate microservice)

### 2. No Foreign Key to Users
**Problem**: `user_id` is VARCHAR with no FK constraint (users are in separate DB)
**Impact**: Orphaned data if users are deleted, referential integrity not enforced
**Recommendation**: Either consolidate users table OR implement application-layer cleanup

### 3. Soft Delete Not Implemented
**Problem**: All deletes are physical (data loss risk)
**Impact**: No audit trail for deleted data, cannot restore accidentally deleted items
**Recommendation**: Add `deleted_at TIMESTAMP` column to key tables

### 4. No Payment History for Expenses
**Problem**: Only `amount_paid` field, no breakdown of partial payments
**Impact**: Cannot track payment history over time
**Recommendation**: Add `expense_payments` table

### 5. Recurring Transactions Not Automated
**Problem**: No cron job or background worker to execute `auto_create = TRUE` transactions
**Impact**: Manual intervention required
**Recommendation**: Implement background worker (Go cron or external scheduler)

### 6. Budget Period Rollover Not Automated
**Problem**: Budgets don't auto-reset at end of period
**Impact**: Manual reset required or stale data
**Recommendation**: Background job to reset budgets

### 7. Inconsistent ID Formats
**Problem**: Some tables use `VARCHAR(36)` (UUID), others `VARCHAR(50)` (custom prefixes like `goal_`, `bud_`)
**Impact**: Inconsistent API contracts, confusion
**Recommendation**: Standardize to UUID v4 or ULID

### 8. No Multi-tenancy Safeguards
**Problem**: No DB-level isolation between users (relies on application WHERE clauses)
**Impact**: Risk of data leakage via bugs
**Recommendation**: Add Row-Level Security (RLS) in PostgreSQL

### 9. No Audit Logging
**Problem**: No tracking of who changed what and when
**Impact**: Cannot debug data issues or comply with audit requirements
**Recommendation**: Add audit log table or use PostgreSQL event triggers

---

## Migration Considerations

When migrating to Backend v2:

### Must Address
1. ✅ **Consolidate gamification tables** (eliminate duplication)
2. ✅ **Add soft delete** (`deleted_at`) to critical tables
3. ✅ **Add audit logging** for compliance
4. ✅ **Standardize ID generation** (UUIDs or ULIDs)

### Should Address
5. ⚠️ **Add expense payment history** table
6. ⚠️ **Implement background workers** for recurring transactions and budget resets
7. ⚠️ **Add Row-Level Security** for multi-tenancy

### Could Address (Nice-to-have)
8. 💡 **Add database-level validation** (more CHECK constraints)
9. 💡 **Add materialized views** for expensive queries (reports)
10. 💡 **Partition large tables** if data grows (>1M rows)

---

## Conclusion

The current schema is **functional but has technical debt** that should be addressed in Backend v2:
- **Data duplication** between main-db and gamification-db
- **No soft deletes** or audit logging
- **Manual processes** for recurring transactions and budget resets
- **Inconsistent ID formats**

With only **1-10 users in production**, this is the **perfect time** to consolidate, clean up, and re-architect before scaling.

---

**Next Steps**:
1. Document gamification-db schema (see [gamification-db.md](./gamification-db.md))
2. Document domain models (Go structs) (see [domain-models.md](./domain-models.md))
3. Define target schema for Backend v2 (see [../02-target-state/](../02-target-state/))
4. Plan migration strategy (see [../03-migrations/strategy.md](../03-migrations/strategy.md))
