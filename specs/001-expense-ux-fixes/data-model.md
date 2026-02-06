# Data Model: Expense Management UX Improvements

**Feature**: 001-expense-ux-fixes
**Date**: 2026-02-05

## Database Schema Changes

### Modified Table: `expenses`

**Changes**: Add `transaction_date` column only

```sql
-- Migration: Add transaction_date column
ALTER TABLE expenses 
ADD COLUMN transaction_date TIMESTAMPTZ;

-- Backfill existing records (set transaction_date = created_at)
UPDATE expenses 
SET transaction_date = created_at 
WHERE transaction_date IS NULL;

-- Make column NOT NULL after backfill
ALTER TABLE expenses 
ALTER COLUMN transaction_date SET NOT NULL;

-- Add index for efficient filtering by month
CREATE INDEX idx_expenses_transaction_date ON expenses(transaction_date);
```

**Updated Schema**:
```sql
CREATE TABLE expenses (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount DECIMAL(10,2) NOT NULL CHECK (amount > 0),
  description TEXT,
  category VARCHAR(100) NOT NULL,
  transaction_date TIMESTAMPTZ NOT NULL,  -- NEW: actual date of expense
  due_date TIMESTAMPTZ,                    -- existing: when payment is due
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Relationships**:
- `user_id` → `users.id` (existing, unchanged)
- Referenced by: `partial_payments.expense_id` (existing table, unchanged)

---

### Existing Table: `partial_payments` (No Changes)

**Note**: This table already exists in the system. No migration needed for partial payments.

The backend should already be calculating `paid_amount` and `pending_amount` from this table.

---

## Entity Definitions

### Expense (Updated - Frontend Only)

```typescript
// Frontend TypeScript type
interface Expense {
  id: number;
  userId: number;
  amount: number;
  description?: string;
  category: string;
  transactionDate: Date;        // NEW: when expense actually occurred
  dueDate?: Date;
  paidAmount: number;           // EXISTING (calculated): sum of partial payments
  pendingAmount: number;        // EXISTING (calculated): amount - paidAmount
  createdAt: Date;
  updatedAt: Date;
}
```

```go
// Backend Go struct
type Expense struct {
    ID              int       `json:"id" db:"id"`
    UserID          int       `json:"user_id" db:"user_id"`
    Amount          float64   `json:"amount" db:"amount"`
    Description     *string   `json:"description" db:"description"`
    Category        string    `json:"category" db:"category"`
    TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`  // NEW
    DueDate         *time.Time `json:"due_date" db:"due_date"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
    
    // Calculated fields (already exist in backend, just need to be in API response)
    PaidAmount    float64 `json:"paid_amount" db:"-"`      // EXISTING
    PendingAmount float64 `json:"pending_amount" db:"-"`   // EXISTING
}
```

---

## Validation Rules

### Expense Validation
- `amount`: Must be > 0
- `category`: Must be from predefined list (validate in service layer)
- `transaction_date`: Can be any date (past, present, future)
- `description`: Optional, but if provided must not be empty string

---

## Migration Script

**File**: `apps/api-gateway/migrations/00X_add_expense_transaction_date.sql`

```sql
-- Migration UP
BEGIN;

-- Add transaction_date column to expenses
ALTER TABLE expenses 
ADD COLUMN transaction_date TIMESTAMPTZ;

-- Backfill existing records
UPDATE expenses 
SET transaction_date = created_at 
WHERE transaction_date IS NULL;

-- Make NOT NULL after backfill
ALTER TABLE expenses 
ALTER COLUMN transaction_date SET NOT NULL;

-- Add index for month filtering
CREATE INDEX idx_expenses_transaction_date ON expenses(transaction_date);

COMMIT;

-- Migration DOWN (for rollback)
BEGIN;

DROP INDEX IF EXISTS idx_expenses_transaction_date;
ALTER TABLE expenses DROP COLUMN IF EXISTS transaction_date;

COMMIT;
```

---

## Summary of Changes

**Backend (Go)**:
- Add `transaction_date` field to Expense model
- Update expense creation/update handlers to accept `transaction_date`
- Ensure list endpoint includes `paid_amount` and `pending_amount` in response (likely already does)

**Frontend (React/TypeScript)**:
- Add `transactionDate` to Expense type
- Add date picker to expense creation form
- Display partial payment badge in table row (using existing `paidAmount`/`pendingAmount`)
- Update inline editing to support `transactionDate` field

**Database**:
- Only ONE migration: Add `transaction_date` column to `expenses` table
- No changes to `partial_payments` table (already exists)
