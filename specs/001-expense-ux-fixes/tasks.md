# Tasks: Expense Management UX Improvements

**Feature**: 001-expense-ux-fixes
**Branch**: `001-expense-ux-fixes`
**Input**: Design documents from `/specs/001-expense-ux-fixes/`
**Prerequisites**: ✅ spec.md, plan.md, research.md, data-model.md, contracts/expenses-api.yaml

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)

---

## Phase 1: Database Migration (Shared Infrastructure)

**Purpose**: Add transaction_date column to expenses table

- [ ] **T001** [US1] Create migration file `apps/api-gateway/migrations/00X_add_expense_transaction_date.sql`
  - Add transaction_date column (TIMESTAMPTZ, nullable initially)
  - Backfill existing records: `SET transaction_date = created_at`
  - Make column NOT NULL after backfill
  - Add index `idx_expenses_transaction_date`
  - Include rollback (DOWN) migration

- [ ] **T002** [US1] Run migration and verify
  - Execute migration on local database
  - Verify column exists: `\d expenses` in psql
  - Verify index created
  - Verify existing records have transaction_date populated

---

## Phase 2: User Story 1 - Custom Transaction Dates (P1)

**Goal**: Allow creating expenses with past/future dates

### Backend (Go)

- [ ] **T003** [P] [US1] Update Expense model in `apps/api-gateway/internal/models/expense.go`
  - Add `TransactionDate time.Time` field with json tag `transaction_date`
  - Add db tag `transaction_date`

- [ ] **T004** [P] [US1] Update expense creation handler in `apps/api-gateway/internal/handlers/expenses.go`
  - Accept `transaction_date` from request body
  - Validate: must be valid date (allow past and future)
  - Default to current time if not provided (backward compatibility)
  - Save to database with transaction_date value

- [ ] **T005** [P] [US1] Update expense list handler to filter by transaction_date
  - Modify query to filter by `transaction_date` instead of `created_at` for month/year filters
  - Ensure month filter uses transaction_date month
  - Return expenses sorted by transaction_date DESC

- [ ] **T006** [US1] Update expense validator in `apps/api-gateway/internal/validators/expense_validator.go`
  - Add validation for transaction_date (must be valid date, allow any date)
  - Add test cases for validation

### Frontend (React/TypeScript)

- [ ] **T007** [P] [US1] Update Expense type in `apps/frontend/src/types/expense.ts`
  - Add `transactionDate: Date` field
  - Update ExpenseFormData type to include transactionDate

- [ ] **T008** [P] [US1] Update ExpenseForm component in `apps/frontend/src/components/forms/ExpenseForm.tsx`
  - Add date picker field for transaction_date
  - Default to today's date
  - Allow selecting any past or future date
  - Validate date is selected before form submission

- [ ] **T009** [P] [US1] Update expenseService in `apps/frontend/src/services/expenseService.ts`
  - Include transaction_date in POST request body
  - Ensure date is serialized correctly (ISO 8601)

- [ ] **T010** [US1] Update dashboard month filter logic
  - Filter expenses by transactionDate instead of createdAt
  - Verify expenses appear in correct month

### Testing (US1)

- [ ] **T011** [US1] Manual test: Create expense with past date
  - Create expense dated 2 months ago
  - Verify appears in that month's dashboard, not current month

- [ ] **T012** [US1] Manual test: Create expense with future date
  - Create expense dated next month
  - Verify appears in future month's dashboard

---

## Phase 3: User Story 2 - Inline Excel-Style Editing (P1)

**Goal**: Enable editing expense fields directly in table without modal

### Frontend (React/TypeScript)

- [ ] **T013** [P] [US2] Create InlineEditCell component in `apps/frontend/src/components/expenses/InlineEditCell.tsx`
  - Accept props: initialValue, onSave, fieldType (text, number, date, select)
  - Double-click or click to enter edit mode
  - Show checkmark icon when in edit mode
  - ESC to cancel, checkmark click to save
  - Show loading state while saving

- [ ] **T014** [P] [US2] Create useInlineEdit hook in `apps/frontend/src/hooks/useInlineEdit.ts`
  - Manage edit state (isEditing, isSaving, value)
  - Implement optimistic update pattern
  - Handle save with rollback on error
  - Return { value, isEditing, save, cancel, setValue }

- [ ] **T015** [US2] Update ExpenseTable component in `apps/frontend/src/components/expenses/ExpenseTable.tsx`
  - Replace read-only cells with InlineEditCell for: amount, description, category, transaction_date, due_date
  - Pass expense.id and field name to InlineEditCell
  - Connect onSave to expenseService.update() with PATCH request

- [ ] **T016** [US2] Update expenseService in `apps/frontend/src/services/expenseService.ts`
  - Add `updateField(id, field, value)` method
  - Use PATCH endpoint: `/api/v1/expenses/:id`
  - Return updated expense with all calculated fields

- [ ] **T017** [US2] Implement dashboard auto-update after inline edit
  - After successful save, update expense in local state
  - Recalculate dashboard totals (pending, paid, balance) client-side
  - No page refresh required

### Backend (Go)

- [ ] **T018** [US2] Update PATCH handler in `apps/api-gateway/internal/handlers/expenses.go`
  - Accept partial updates (any subset of fields)
  - Validate each field that's provided
  - Return full expense object with calculated fields after update

### Testing (US2)

- [ ] **T019** [US2] Manual test: Inline edit amount
  - Click amount cell, modify value, click checkmark
  - Verify saves without modal, dashboard totals update instantly

- [ ] **T020** [US2] Manual test: Inline edit multiple fields sequentially
  - Edit amount, then description, then category
  - Verify each saves independently

- [ ] **T021** [US2] Manual test: Cancel inline edit with ESC
  - Enter edit mode, modify value, press ESC
  - Verify value reverts to original

---

## Phase 4: User Story 3 - Partial Payment Display (P2)

**Goal**: Show partial payment status directly in expense table row

### Backend Verification

- [ ] **T022** [US3] Verify backend includes paid_amount and pending_amount in list endpoint
  - Check GET `/api/v1/expenses` response
  - Ensure each expense object has `paid_amount` and `pending_amount` fields
  - If missing, add calculation in expense service layer
  - Calculate from existing `partial_payments` table (SUM of amounts WHERE expense_id = ?)

### Frontend (React/TypeScript)

- [ ] **T023** [P] [US3] Create PartialPaymentBadge component in `apps/frontend/src/components/expenses/PartialPaymentBadge.tsx`
  - Accept props: expense (with paidAmount, pendingAmount, amount fields)
  - Only render if paidAmount > 0
  - Display format: "Paid: $30 of $100, $70 pending"
  - Style as badge/chip (distinct visual indicator)

- [ ] **T024** [US3] Add PartialPaymentBadge to ExpenseTable rows
  - Import PartialPaymentBadge component
  - Add new column or inline badge in amount/status column
  - Pass expense data to badge component
  - Ensure badge is visible on desktop (mobile can be in modal still)

- [ ] **T025** [US3] Update pending expenses calculation in dashboard
  - Use expense.pendingAmount instead of expense.amount
  - Verify total pending expenses only counts unpaid portions

### Testing (US3)

- [ ] **T026** [US3] Manual test: View expense with partial payment
  - Find or create expense with partial payment
  - Verify badge appears in table row showing "Paid: $X of $Y, $Z pending"
  - Verify badge does NOT appear for fully unpaid expenses

- [ ] **T027** [US3] Manual test: Pending expenses calculation
  - Note total pending expenses
  - Verify it matches sum of all expense.pendingAmount values, not expense.amount

---

## Phase 5: Polish & Edge Cases

- [ ] **T028** [ALL] Add client-side validation for inline editing
  - Amount must be positive number
  - Required fields cannot be empty
  - Show inline error message, prevent save

- [ ] **T029** [ALL] Add error handling for inline edit failures
  - On save error, rollback to original value
  - Show user-friendly error toast/message
  - Log error for debugging

- [ ] **T030** [ALL] Responsive design for inline editing on mobile
  - Test inline edit on small screens
  - Consider fallback to modal on very small screens if needed

- [ ] **T031** [ALL] Update API contract validation
  - Verify OpenAPI spec matches implementation
  - Test all endpoints with Postman/curl

---

## Summary

**Total Tasks**: 31
- **Phase 1 (Setup)**: 2 tasks
- **Phase 2 (US1 - Transaction Dates)**: 10 tasks (6 backend, 4 frontend)
- **Phase 3 (US2 - Inline Editing)**: 9 tasks (7 frontend, 1 backend, 1 service)
- **Phase 4 (US3 - Partial Payments)**: 6 tasks (1 backend verify, 3 frontend, 2 tests)
- **Phase 5 (Polish)**: 4 tasks

**Parallelizable**: Tasks marked [P] can run in parallel

**Testing Strategy**: Manual testing after each user story phase
