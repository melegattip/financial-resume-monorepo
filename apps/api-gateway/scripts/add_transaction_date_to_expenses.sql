-- =====================================================
-- MIGRATION: Add transaction_date to expenses table
-- Feature: 001-expense-ux-fixes
-- Bug: BUG-015 - Allow creating expenses with custom dates
-- =====================================================

-- Add transaction_date column (nullable initially for backfill)
ALTER TABLE expenses 
ADD COLUMN IF NOT EXISTS transaction_date TIMESTAMP;

-- Backfill existing records with created_at value
UPDATE expenses 
SET transaction_date = created_at 
WHERE transaction_date IS NULL;

-- Make column NOT NULL after backfill
ALTER TABLE expenses 
ALTER COLUMN transaction_date SET NOT NULL;

-- Add index for efficient filtering by transaction_date (month/year queries)
CREATE INDEX IF NOT EXISTS idx_expenses_transaction_date ON expenses(transaction_date);

-- Add comment for documentation
COMMENT ON COLUMN expenses.transaction_date IS 'When the expense actually occurred (can be past, present, or future)';

-- =====================================================
-- ROLLBACK (uncomment to rollback this migration)
-- =====================================================
-- DROP INDEX IF EXISTS idx_expenses_transaction_date;
-- ALTER TABLE expenses DROP COLUMN IF EXISTS transaction_date;
