# Gamification Database Schema - Current State (Production)

**Database**: `financial_gamification` (PostgreSQL Gamification)
**Location**: Production (separate instance/port)
**Last Updated**: 2026-02-09
**Status**: ⚠️ PRODUCTION - Part of gamification-service microservice

## Overview

This document describes the **current production schema** for the gamification service database. This schema handles:
- User XP and leveling system
- Achievements/Badges
- Daily/Weekly challenges
- User action tracking (audit trail for XP)
- Challenge progress tracking

⚠️ **CRITICAL ISSUE**: Tables `user_gamification`, `achievements`, and `user_actions` are **DUPLICATED** in main-db. This creates data consistency risks.

---

## Architecture Context

```
┌─────────────────────────┐
│  main-db (PostgreSQL)   │
│  ┌──────────────────┐   │      ┌─────────────────────────────┐
│  │user_gamification │───┼──────│  gamification-db (PostgreSQL)│
│  │achievements      │   │      │  ┌──────────────────┐        │
│  │user_actions      │   │ SYNC?│  │user_gamification │        │
│  └──────────────────┘   │      │  │achievements      │        │
└─────────────────────────┘      │  │user_actions      │        │
                                  │  │challenges        │        │
                                  │  │user_challenges   │        │
                                  │  │challenge_progress│        │
                                  │  └──────────────────┘        │
                                  └─────────────────────────────┘
```

**Problem**: No clear sync mechanism between duplicated tables.

---

## Entity Relationship Diagram

```
┌────────────────────┐
│ user_gamification  │ (Primary user stats)
└────────────────────┘
         │
         │ (FK: user_id)
         │
         ├─────────────────────┬──────────────────────────┐
         │                     │                          │
         ▼                     ▼                          ▼
┌───────────────┐   ┌────────────────┐     ┌───────────────────────┐
│ achievements  │   │ user_actions   │     │  user_challenges      │
└───────────────┘   └────────────────┘     └───────────────────────┘
                                                    │
                                                    │ (FK: challenge_id)
                                                    ▼
                                           ┌──────────────┐
                                           │  challenges  │
                                           └──────────────┘

┌───────────────────────────┐
│ challenge_progress_tracking│ (Daily/weekly aggregates)
└───────────────────────────┘
```

---

## Table Schemas

### 1. User Gamification

**Purpose**: Core Financial Health Score stats per user

```sql
CREATE TABLE user_gamification (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    financial_health_score INTEGER DEFAULT 1,
    engagement_component INTEGER DEFAULT 0,
    health_component INTEGER DEFAULT 0,
    current_level INTEGER DEFAULT 1,
    insights_viewed INTEGER DEFAULT 0,
    actions_completed INTEGER DEFAULT 0,
    achievements_count INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_score_calculation TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_user_gamification_user_id` on `user_gamification(user_id)`
- `idx_user_gamification_score` on `user_gamification(financial_health_score DESC)`
- `idx_user_gamification_level` on `user_gamification(current_level)`
- `idx_user_gamification_last_calc` on `user_gamification(last_score_calculation)`

**Business Logic**:
- Score → Level mapping (application layer):
  - Level 1 Aprendiz: 1-400 Score 🌱
  - Level 2 Manager: 401-500 Score 💼 (desbloquea Presupuestos)
  - Level 3 Guru: 501-600 Score 🧠 (desbloquea Metas de Ahorro)
  - Level 4 Master: 601-900 Score 👑 (desbloquea IA Financiera)
  - Level 5 Magnate: 901-1000 Score 💎
- `financial_health_score` = `engagement_component` + `health_component` (1-1000)
- `engagement_component`: 0-400 pts (Racha, Adopción, Insights, Challenges)
- `health_component`: 0-600 pts (Savings Rate, Budget Adherence, Goals Progress, Expense Management)
- `current_streak`: Consecutive days with activity (contributes to engagement_component)
- `last_activity`: Updated on any user action
- `last_score_calculation`: Timestamp of last full score recalculation

**⚠️ Issue**: **DUPLICATED** in main-db. Which is the source of truth?

---

### 2. Achievements

**Purpose**: Track user achievements/badges

```sql
CREATE TABLE achievements (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    points INTEGER DEFAULT 0,
    progress INTEGER DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    unlocked_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);
```

**Indexes**:
- `idx_achievements_user_id` on `achievements(user_id)`
- `idx_achievements_type` on `achievements(type)`
- `idx_achievements_completed` on `achievements(completed)`

**Achievement Types** (examples from code):
- `first_transaction`: Log first expense/income
- `budget_master`: Set first budget
- `saver`: Complete first savings goal
- `streak_7`: Maintain 7-day streak
- `xp_milestone_100`: Reach 100 XP

**Business Logic**:
- `progress / target` = completion percentage
- `completed = TRUE` when `progress >= target`
- `unlocked_at` timestamp set on completion

**⚠️ Issue**: **DUPLICATED** in main-db.

---

### 3. User Actions

**Purpose**: Audit trail for all actions that contribute to Score

```sql
CREATE TABLE user_actions (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100) NOT NULL,
    entity_id VARCHAR(255),
    score_contribution INTEGER DEFAULT 0,
    component_affected VARCHAR(50), -- 'engagement' or 'health'
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE
);
```

**Indexes**:
- `idx_user_actions_user_id` on `user_actions(user_id)`
- `idx_user_actions_type` on `user_actions(action_type)`
- `idx_user_actions_created_at` on `user_actions(created_at DESC)`
- `idx_user_actions_component` on `user_actions(component_affected)`

**Action Types** (Score contributions):
- `create_expense`: Contributes to Engagement → Adopción de features
- `create_income`: Contributes to Engagement → Adopción de features
- `create_budget`: Contributes to Engagement → Adopción de features
- `create_savings_goal`: Contributes to Engagement → Adopción de features
- `view_dashboard`: Contributes to Engagement → Consumo de insights (once per day)
- `view_analytics`: Contributes to Engagement → Consumo de insights
- `complete_challenge`: Contributes to Engagement → Challenges completados
- `achieve_savings_milestone`: Contributes to Health → Goals Progress
- `stay_within_budget`: Contributes to Health → Budget Adherence

**Business Logic**:
- Each action appends a record (append-only log)
- `score_contribution` indicates points added to specific component
- `component_affected` indicates which component ('engagement' or 'health')
- Score is recalculated periodically based on aggregated actions and financial data

**⚠️ Issue**: **DUPLICATED** in main-db.

---

### 4. Challenges

**Purpose**: Define available challenges (daily/weekly/monthly)

```sql
CREATE TABLE challenges (
    id VARCHAR(255) PRIMARY KEY,
    challenge_key VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    challenge_type VARCHAR(50) NOT NULL, -- 'daily', 'weekly', 'monthly'
    icon VARCHAR(20) DEFAULT '🎯',
    xp_reward INTEGER DEFAULT 0,
    requirement_type VARCHAR(100) NOT NULL,
    requirement_target INTEGER NOT NULL,
    requirement_data JSONB,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes**:
- `idx_challenges_type` on `challenges(challenge_type)`
- `idx_challenges_active` on `challenges(active)`
- `idx_challenges_key` on `challenges(challenge_key)`

**Challenge Types**:
- `daily`: Resets every day at midnight
- `weekly`: Resets every Monday
- `monthly`: Resets on 1st of month

**Requirement Types**:
- `transaction_count`: Log N transactions
- `category_variety`: Use N different categories
- `view_combo`: View dashboard + analytics
- `daily_login`: Login once per day
- `daily_login_count`: Login N days in period

**Requirement Data** (JSONB examples):
```json
{
  "actions": ["create_expense", "create_income"]
}
```

**Default Challenges** (seeded):
```sql
-- Daily
'transaction_master': 'Registra 3 transacciones hoy' → 20 XP
'category_organizer': 'Usa 3 categorías diferentes' → 15 XP
'analytics_explorer': 'Revisa dashboard y analytics' → 10 XP
'streak_keeper': 'Mantén tu racha diaria' → 5 XP

-- Weekly
'transaction_champion': 'Registra 15 transacciones esta semana' → 75 XP
'category_master': 'Usa al menos 5 categorías diferentes' → 50 XP
'engagement_hero': 'Inicia sesión 5 días esta semana' → 60 XP
```

---

### 5. User Challenges

**Purpose**: Track user progress on challenges

```sql
CREATE TABLE user_challenges (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL,
    progress INTEGER DEFAULT 0,
    target INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE,
    FOREIGN KEY (challenge_id) REFERENCES challenges(id) ON DELETE CASCADE,
    UNIQUE(user_id, challenge_id, challenge_date)
);
```

**Indexes**:
- `idx_user_challenges_user_date` on `user_challenges(user_id, challenge_date)`
- `idx_user_challenges_completed` on `user_challenges(completed)`
- `idx_user_challenges_user_id` on `user_challenges(user_id)`

**Business Logic**:
- `challenge_date`: For daily challenges, this is the specific day (e.g., 2026-02-09)
- `challenge_date`: For weekly challenges, this is the Monday of that week
- `challenge_date`: For monthly challenges, this is the 1st day of month
- `UNIQUE` constraint prevents duplicate tracking for same challenge on same date
- `completed = TRUE` when `progress >= target`
- `completed_at` timestamp set on completion

**Example**:
```
user_id='user123', challenge_id='transaction_master', challenge_date='2026-02-09'
progress=2, target=3, completed=FALSE
```

---

### 6. Challenge Progress Tracking

**Purpose**: Aggregate daily/weekly action counts for challenge evaluation

```sql
CREATE TABLE challenge_progress_tracking (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    challenge_date DATE NOT NULL,
    action_type VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    count INTEGER DEFAULT 1,
    unique_entities JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user_gamification(user_id) ON DELETE CASCADE,
    UNIQUE(user_id, challenge_date, action_type, entity_type)
);
```

**Indexes**:
- `idx_challenge_tracking_user_date` on `challenge_progress_tracking(user_id, challenge_date)`
- `idx_challenge_tracking_action` on `challenge_progress_tracking(action_type)`

**Purpose**:
Efficiently track aggregated metrics for challenge requirements without scanning all `user_actions`.

**Example**:
```
user_id='user123', challenge_date='2026-02-09', action_type='create_expense'
count=5, unique_entities='{"categories": ["food", "transport", "entertainment"]}'
```

This enables:
- `transaction_count` challenges: Check `count` field
- `category_variety` challenges: Check `unique_entities.categories` array length

**Business Logic**:
- Updated on each user action
- UPSERT (INSERT ... ON CONFLICT UPDATE)
- `UNIQUE` constraint enables efficient increments

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

CREATE TRIGGER update_user_gamification_updated_at
    BEFORE UPDATE ON user_gamification
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_achievements_updated_at
    BEFORE UPDATE ON achievements
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_challenges_updated_at
    BEFORE UPDATE ON challenges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_challenges_updated_at
    BEFORE UPDATE ON user_challenges
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## Known Issues & Technical Debt

### 1. Data Duplication with main-db (CRITICAL)
**Problem**: `user_gamification`, `achievements`, `user_actions` exist in **BOTH** databases
**Impact**:
- Which is the source of truth?
- Sync issues if both are written to
- Data inconsistency risk

**Current Behavior** (from code analysis):
- gamification-service **reads/writes** to `gamification-db`
- api-gateway **proxies** requests to gamification-service (HTTP calls)
- ⚠️ **Unclear if main-db tables are still being used**

**Recommendation**:
- **Option A**: Delete gamification tables from main-db, use only gamification-db
- **Option B**: Eliminate gamification-service microservice, consolidate to main-db (monolith approach)
- **Option C**: Implement event-driven sync (publish events on write, subscribe and replicate)

**For monolith modular approach**: Choose **Option B** (consolidate to main-db)

---

### 2. No Foreign Key to Users
**Problem**: `user_id` is VARCHAR with no FK constraint (users are in separate `users_db`)
**Impact**: Orphaned gamification data if users are deleted
**Recommendation**: Application-layer cleanup OR consolidate users table

---

### 3. Challenge Reset Logic Not in DB
**Problem**: No automatic reset of daily/weekly challenges at period boundaries
**Impact**: Requires cron job or background worker
**Current Implementation**: Assumed to be in gamification-service application code
**Recommendation**: Verify cron job exists, document in architecture

---

### 4. Score Calculation is Eventually Consistent
**Problem**: `financial_health_score` in `user_gamification` is calculated from multiple sources (user_actions, financial data)
**Impact**: Score may be stale if recalculation is not triggered properly
**Recommendation**:
- Use database transactions to ensure atomicity of action logging
- Implement periodic score recalculation (e.g., daily at midnight)
- Cache score but mark as "needs_recalc" flag when relevant data changes

---

### 5. No Soft Delete
**Problem**: Deletes are physical (data loss)
**Impact**: Cannot restore accidentally deleted challenges or achievements
**Recommendation**: Add `deleted_at TIMESTAMP` to tables

---

### 6. JSONB Fields Not Validated
**Problem**: `requirement_data` in `challenges` and `unique_entities` in `challenge_progress_tracking` have no schema validation
**Impact**: Invalid JSON can be inserted, breaking application logic
**Recommendation**: Use PostgreSQL CHECK constraints with `jsonb_typeof()` or JSON Schema validation

---

### 7. No Audit Logging
**Problem**: No tracking of who changed challenge definitions or modified XP
**Impact**: Cannot debug admin actions or XP manipulation
**Recommendation**: Add audit log table

---

### 8. Challenge Completion Idempotency
**Problem**: Completing the same challenge multiple times in one day could inflate Score
**Impact**: Score inflation
**Current Mitigation**: `UNIQUE(user_id, challenge_id, challenge_date)` constraint prevents duplicate tracking
**Status**: ✅ Handled correctly

---

## Migration Considerations for Backend v2

### Must Address
1. ✅ **Consolidate gamification data to one DB** (eliminate duplication with main-db)
2. ✅ **Verify/implement challenge reset cron job**
3. ✅ **Ensure Score calculations are transactional** (atomic updates)
4. ✅ **Implement Financial Health Score calculation algorithm** (combining engagement + health components)

### Should Address
4. ⚠️ **Add soft delete** (`deleted_at`)
5. ⚠️ **Add JSONB schema validation**
6. ⚠️ **Add audit logging** for challenge/achievement changes

### Could Address (Nice-to-have)
7. 💡 **Add database views** for leaderboards (top XP, top streaks)
8. 💡 **Add materialized views** for expensive aggregations
9. 💡 **Add partitioning** for `user_actions` if it grows large

---

## Data Consolidation Strategy (for Monolith Modular)

When consolidating gamification into main-db:

### Step 1: Verify Current Usage
- [ ] Confirm gamification-service is the only writer to `gamification-db`
- [ ] Confirm main-db gamification tables are **NOT** being written to
- [ ] Identify if any reads happen from main-db gamification tables

### Step 2: Plan Migration
```sql
-- Option A: If main-db tables are unused, just drop them
DROP TABLE IF EXISTS user_actions CASCADE;
DROP TABLE IF EXISTS achievements CASCADE;
DROP TABLE IF EXISTS user_gamification CASCADE;

-- Option B: If main-db has data, migrate it first
-- (Export from main-db, import to gamification-db, verify, then drop main-db tables)
```

### Step 3: Consolidate Service
- Merge gamification-service logic into api-gateway as a **domain module**
- Eliminate HTTP proxy calls (internal function calls instead)
- Keep clean architecture (domain → usecases → repository)

### Step 4: Single Database
- Migrate `challenges`, `user_challenges`, `challenge_progress_tracking` to main-db
- Update connection strings
- Deploy backend v2

---

## Conclusion

The gamification-db schema is **well-designed** for its purpose and has been **updated to support the unified Financial Health Score**:
- ✅ Efficient challenge tracking with `challenge_progress_tracking`
- ✅ Proper unique constraints to prevent duplicate Score contributions
- ✅ JSONB for flexible requirement definitions
- ✅ New fields for tracking Score components (engagement_component, health_component)
- ✅ Support for both engagement-based and financial health-based scoring

**Schema Changes Required**:
1. Rename `total_xp` to `financial_health_score`
2. Add `engagement_component` (0-400)
3. Add `health_component` (0-600)
4. Add `last_score_calculation` timestamp
5. Rename `xp_earned` to `score_contribution` in user_actions
6. Add `component_affected` field in user_actions

**However**, the **duplication with main-db** is still a **critical architectural issue** that must be resolved:
- For **monolith modular** approach: Consolidate to main-db, eliminate microservice
- For **microservices** approach: Ensure proper event-driven sync OR remove duplicated tables

With only **1-10 users**, this is the **perfect time** to consolidate and implement the unified Financial Health Score system.

---

**Next Steps**:
1. Document domain models (Go structs) (see [domain-models.md](./domain-models.md))
2. Define target schema for Backend v2 (see [../02-target-state/](../02-target-state/))
3. Plan consolidation migration (see [../03-migrations/strategy.md](../03-migrations/strategy.md))
