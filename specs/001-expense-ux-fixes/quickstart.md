# Quickstart: Testing Expense UX Improvements Locally

**Feature**: 001-expense-ux-fixes
**Branch**: `001-expense-ux-fixes`
**Date**: 2026-02-05

## Prerequisites

- Docker Desktop running (for PostgreSQL)
- Go 1.21+ installed
- Node.js 18+ installed
- Git repository cloned

## Setup Steps

### 1. Checkout Feature Branch

```bash
git checkout 001-expense-ux-fixes
```

### 2. Run Database Migration

```bash
# Navigate to api-gateway module
cd apps/api-gateway

# Run migration (only adds transaction_date column)
go run cmd/migrate/main.go up

# Or manually:
# psql -d financial_resume -f migrations/00X_add_expense_transaction_date.sql
```

**Expected output**: `expenses` table now has `transaction_date` column.

### 3. Start Backend (api-gateway)

```bash
# From apps/api-gateway directory
go run cmd/server/main.go

# Or if using Docker:
# docker-compose up api-gateway
```

**Expected**: Server starts on `http://localhost:8000` (or configured port)

### 4. Start Frontend

```bash
# Navigate to frontend
cd apps/frontend

# Install dependencies (if needed)
npm install

# Start dev server
npm run dev
```

**Expected**: Frontend starts on `http://localhost:3000` or `http://localhost:5173` (Vite default)

## Testing the Features

### Test 1: Create Expense with Custom Date (BUG-015 Fix)

1. Navigate to expense creation form
2. **Expected**: See a date picker field for "Transaction Date"
3. Click the date picker
4. Select a date in the past (e.g., 2 months ago)
5. Fill in amount ($100), category ("Food"), description ("Old grocery bill")
6. Submit
7. **Verify**: Expense appears in the dashboard for the selected month, NOT current month

**Test with Future Date**:
1. Create expense with future date (e.g., next month)
2. **Expected**: Allowed, no error
3. **Verify**: Expense appears in future month's dashboard

**API Request Example**:
```bash
curl -X POST http://localhost:8000/api/v1/expenses \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.00,
    "description": "Old grocery bill",
    "category": "Food",
    "transaction_date": "2025-12-15T10:00:00Z"
  }'
```

### Test 2: Inline Editing (BUG-013 Fix)

1. Navigate to dashboard with expenses list
2. Click on any expense field (e.g., amount or description)
3. **Expected**: Field becomes editable inline (no modal)
4. Modify the value
5. **Expected**: Checkmark icon appears
6. Click checkmark
7. **Expected**: 
   - Value saves immediately
   - Dashboard totals update without page refresh
   - No modal or loading screen

**Test Multiple Fields**:
1. Edit amount, then description, then category (one by one)
2. **Expected**: Each edit saves independently with checkmark

**API Request Example**:
```bash
curl -X PATCH http://localhost:8000/api/v1/expenses/123 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"description": "Updated description"}'
```

### Test 3: Partial Payment Display (BUG-014 Fix)

**Note**: Partial payments system already exists. This test verifies the NEW display in table row.

#### View Expense with Partial Payment
1. Find an expense that has partial payments (create one via modal if needed)
2. **Expected - IN THE TABLE ROW** (not modal):
   - Visual badge/indicator visible
   - Shows format: "Paid: $30 of $100, $70 pending"
   - Badge only appears if expense has partial payments

#### Validate Pending Calculation
1. Note total "Pending Expenses" amount in dashboard
2. Find expense with $100 total, $30 paid
3. **Expected**: Only $70 counts toward pending expenses total

## Verification Checklist

After testing, verify all 3 bugs are fixed:

**BUG-015: Create Expense with Custom Date**
- [ ] Can create expense with past date
- [ ] Can create expense with future date
- [ ] Expense appears in correct month based on transaction_date
- [ ] Transaction date field is visible in creation form

**BUG-013: Inline Editing**
- [ ] Can click cell to edit inline (no modal)
- [ ] Checkmark appears for confirmation
- [ ] Saves without page refresh
- [ ] Dashboard totals update immediately
- [ ] Can edit multiple fields sequentially
- [ ] Can cancel edit with Escape

**BUG-014: Partial Payment Display**
- [ ] Partial payment badge appears in table row (not just modal)
- [ ] Shows "Paid: $X of $Y, $Z pending" format
- [ ] Badge only shows when paid_amount > 0
- [ ] Pending expenses calculation is correct

## Edge Cases to Test

- [ ] Future dates are accepted
- [ ] Cannot save negative amounts
- [ ] Cannot save empty required fields
- [ ] Inline editing validation shows errors inline

## Troubleshooting

### Migration Fails
- Check PostgreSQL is running
- Verify database connection string in `.env`
- Ensure no typos in migration SQL

### Inline Editing Not Working
- Check browser console for JavaScript errors
- Verify API endpoints return 200 OK
- Check network tab for failed requests

### Partial Payment Badge Not Showing
- Verify API response includes `paid_amount` and `pending_amount` fields
- Check backend is calculating these fields from `partial_payments` table
- Inspect component props in React DevTools

### Transaction Date Not Saving
- Check migration ran successfully
- Verify backend handler accepts `transaction_date` field
- Check validation isn't blocking the date

## Rollback (if needed)

```bash
# Rollback database migration
cd apps/api-gateway
go run cmd/migrate/main.go down

# Or manually:
# DROP INDEX IF EXISTS idx_expenses_transaction_date;
# ALTER TABLE expenses DROP COLUMN IF EXISTS transaction_date;

# Switch back to main branch
git checkout master
```
