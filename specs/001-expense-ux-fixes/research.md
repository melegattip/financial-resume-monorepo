# Research: Expense Management UX Improvements

**Feature**: 001-expense-ux-fixes
**Date**: 2026-02-05

## Technical Decisions

### 1. Database Schema for Transaction Dates

**Decision**: Add `transaction_date` column to existing `expenses` table

**Rationale**:
- Current system likely only has `created_at` and `due_date`
- Need separate field to track when expense actually occurred
- Allows retroactive and future expense creation
- Keeps expense in correct month regardless of when it was entered

**Alternatives Considered**:
- Using `created_at` as transaction date: Rejected because users need to backdate expenses
- Separate `historical_expenses` table: Rejected - adds unnecessary complexity

**Implementation**:
```sql
-- Add nullable transaction_date column
-- Default to created_at for existing records
-- Non-nullable for new records
ALTER TABLE expenses ADD COLUMN transaction_date TIMESTAMPTZ;
```

---

### 2. Partial Payments Display in Table Row

**Decision**: Display partial payment summary directly in the expense table row

**Current State**:
- Partial payments system already exists (`partial_payments` table)
- Backend already calculates `paid_amount` and `pending_amount`
- Data is available in modal but NOT shown in table row

**Problem (BUG-014)**:
- Users can't see partial payment status without opening modal
- Pending expenses calculation doesn't reflect partial payments
- No visual indicator in the row

**Solution**:
- Add partial payment badge/indicator to the row component
- Show format: "Paid: $30 of $100, $70 pending"
- Only show when expense has partial payments (paid_amount > 0)
- Ensure backend includes these calculated fields in list endpoint response

**No Backend Changes Needed**: Backend likely already has this logic, just need to consume it in frontend row display.

---

###3. Inline Editing Implementation

**Decision**: Use React state management with optimistic updates + server validation

**Rationale**:
- Provides instant feedback to user (feels fast)
- Falls back to server value if validation fails
- Prevents race conditions with single source of truth
- Follows React best practices for form handling

**Alternatives Considered**:
- Direct DOM manipulation: Rejected -not React-friendly, harder to maintain
- Full form modal: Rejected - current problem we're solving
- Uncontrolled inputs: Rejected - harder to validate and sync state

**Pattern**:
```typescript
// Optimistic update pattern
const handleCellEdit = async (field, value) => {
  // 1. Update UI immediately (optimistic)
  setLocalExpense(prev => ({ ...prev, [field]: value }));
  
  // 2. Validate and save to server
  try {
    const updated = await expenseService.update(expense.id, { [field]: value });
    // 3. Confirm with server response
    setExpense(updated);
  } catch (error) {
    // 4. Rollback on error
    setLocalExpense(expense);
    showError(error.message);
  }
};
```

---

### 4. Date Picker Library

**Decision**: Use native HTML5 date input or lightweight library like `react-datepicker`

**Rationale**:
- Need to support past, present, and future dates
- Must handle edge cases (leap years, month-end dates)
- Should be accessible and mobile-friendly
- Existing project likely already has a date picker

**Alternatives Considered**:
- Build custom date picker: Rejected - reinventing the wheel, accessibility concerns
- Heavy UI library (Material-UI DatePicker): Depends on existing dependencies

**Research Needed**:
- Check what date picker is currently used in the project
- Verify it supports the required date range (any past/future date)

---

### 5. Real-time Dashboard Updates

**Decision**: Use React state updates + recalculate totals client-side after each save

**Rationale**:
- Dashboard already has all expense data loaded
- Simple recalculation on client avoids extra API calls
- Provides instant feedback
- No need for WebSockets/polling for this use case

**Alternatives Considered**:
- Refetch entire dashboard after each edit: Rejected - slower, wasteful
- WebSocket for real-time updates: Rejected - overkill for single-user edits

**Implementation**: After successful expense update, trigger `recalculateTotals()` function that sums amounts client-side.

---

### 6. Validation Strategy

**Decision**: Client-side validation for UX + Server-side validation for security

**Rationale**:
- Client-side: Immediate feedback, better UX
- Server-side: Authoritative, prevents malicious requests
- Both layers are necessary for robust system

**Validation Rules**:
- Amount: Must be positive number
- Description: Non-empty (unless optional)
- Category: Must be from predefined list
- Transaction date: Any valid date (past/present/future)

---

## Best Practices Research

### Go Backend Patterns

**Handler Pattern**:
```go
// Use standard handler pattern for expense endpoints
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
    // 1. Parse and validate input
    // 2. Call service layer
    // 3. Return JSON response
}
```

**Service Layer**:
- Keep business logic in service layer (ExpenseService)
- Handlers are thin wrappers that parse HTTP and call services
- Services return domain errors, handlers map to HTTP status codes

### React/TypeScript Frontend Patterns

**Custom Hook for Inline Editing**:
```typescript
const useInlineEdit = (initialValue, onSave) => {
  const [value, setValue] = useState(initialValue);
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  
  const save = async () => {
    setIsSaving(true);
    try {
      await onSave(value);
      setIsEditing(false);
    } catch (error) {
      setValue(initialValue); // Rollback
    } finally {
      setIsSaving(false);
    }
  };
  
  return { value, setValue, isEditing, setIsEditing, save, isSaving };
};
```

**Table Component Pattern**:
- If using TanStack Table: Use column definitions with custom cell renderers
- Each editable cell is a separate component for reusability
- Handle keyboard events (Enter to save, Escape to cancel)

**Partial Payment Badge Component**:
```typescript
const PartialPaymentBadge = ({ expense }: { expense: Expense }) => {
  if (expense.paidAmount === 0) return null;
  
  return (
    <div className="partial-payment-badge">
      Paid: ${expense.paidAmount} of ${expense.amount}, ${expense.pendingAmount} pending
    </div>
  );
};
```

---

## Dependencies to Add/Update

### Backend (Go)
- No new dependencies needed (using stdlib + existing PostgreSQL driver)
- Partial payments system already exists

### Frontend (React/TypeScript)
- Possibly `react-datepicker` if not already present
- No other new dependencies required

---

## Risks & Mitigations

**Risk 1**: Migrating existing expenses without transaction_date
- **Mitigation**: Set `transaction_date = created_at` for all existing records in migration

**Risk 2**: Concurrent edits by same user in multiple tabs
- **Mitigation**: Optimistic locking or last-write-wins strategy (acceptable for single-user scenario)

**Risk 3**: Inline editing breaks on mobile/small screens
- **Mitigation**: Responsive design, possibly fallback to modal on very small screens

**Risk 4**: Backend doesn't include paid_amount/pending_amount in list endpoint
- **Mitigation**: Verify backend response, add calculated fields if missing (minimal backend change)
