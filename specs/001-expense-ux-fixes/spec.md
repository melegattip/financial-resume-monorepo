# Feature Specification: Expense Management UX Improvements

**Feature Branch**: `001-expense-ux-fixes`  
**Created**: 2026-02-05  
**Status**: Draft  
**Input**: User description: "Mejorar UX del módulo de gestión de gastos para resolver tres bugs críticos: 1) Permitir crear gastos con fechas retroactivas seleccionando el mes deseado (BUG-015), 2) Implementar edición inline tipo Excel directamente en la tabla sin modales lentos, con confirmación simple mediante check (BUG-013), 3) Corregir visualización de pagos parciales mostrando indicador visual claro y monto pagado, asegurando que se descuente correctamente de gastos pendientes (BUG-014)."

## User Scenarios & Testing

### User Story 1 - Create Expenses with Custom Dates (Priority: P1)

Users frequently need to record expenses that occurred in previous months (e.g., discovering an old receipt, correcting historical data, or categorizing past transactions). Currently, the system doesn't provide a way to select a custom date when creating an expense, forcing all new expenses to be created in the current month.

**Why this priority**: This is blocking basic functionality. Users cannot maintain accurate historical financial records without the ability to backdate expenses.

**Independent Test**: Can be fully tested by creating an expense and verifying that a date picker allows selecting any past month, and that the expense appears in the correct month's dashboard.

**Acceptance Scenarios**:

1. **Given** a user is creating a new expense, **When** they open the creation form, **Then** they see a date field that defaults to today's date
2. **Given** a user is on the expense creation form, **When** they click the date field, **Then** a date picker appears allowing them to select any past, current or future date
3. **Given** a user selects a date in a previous month, **When** they complete the expense creation, **Then** the expense appears in that month's dashboard, not the current month
4. **Given** a user creates an expense for January while viewing March, **When** they navigate to January's dashboard, **Then** they see the newly created expense

---

### User Story 2 - Inline Excel-style Expense Editing (Priority: P1)

Users find the current modal-based editing workflow slow and cumbersome for making quick updates to expenses. They often need to edit multiple fields across several expenses (correcting amounts, updating categories, adjusting descriptions) and the current flow requires opening a modal, waiting for it to load, making changes, and closing it for each edit.

**Why this priority**: This is the most common daily interaction with expenses. Slow editing directly impacts productivity and user satisfaction.

**Independent Test**: Can be fully tested by clicking on any cell in the expense table, editing the value inline, confirming with a checkmark, and verifying the change persists without a page refresh.

**Acceptance Scenarios**:

1. **Given** a user hovers over an expense row, **When** they click on any editable cell (amount, description, category, due date), **Then** that cell becomes editable in place
2. **Given** a cell is in edit mode, **When** the user modifies the value, **Then** a checkmark icon appears next to the cell
3. **Given** a user has modified a cell value, **When** they click the checkmark, **Then** the change is saved immediately without a full page refresh
4. **Given** a user is editing a cell, **When** they press Escape or click outside, **Then** the changes are discarded and the cell returns to its original value
5. **Given** a user edits multiple cells in succession, **When** each is confirmed, **Then** the dashboard totals update automatically without requiring a page reload

---

### User Story 3 - Clear Partial Payment Tracking (Priority: P2)

Users who make partial payments on expenses (e.g., paying off a credit card balance in installments, or splitting a large bill) currently cannot see the payment status clearly. When a partial payment is recorded, it doesn't update the "pending expenses" calculation, and there's no visual indicator on the expense row showing that it's been partially paid.

**Why this priority**: Accurate debt tracking is critical for financial planning. Hidden partial payments lead to double-counting of pending debts.

**Independent Test**: Can be fully tested by creating an expense, making a partial payment, and verifying that both the visual indicator and the pending expenses calculation reflect the partial payment correctly.

**Acceptance Scenarios**:

1. **Given** an expense with a $100 total amount, **When** a user records a $30 partial payment, **Then** the expense row displays a visual indicator (e.g., "Paid: $30 of $100")
2. **Given** an expense has been partially paid, **When** viewing the dashboard, **Then** only the remaining unpaid amount ($70 in this case) counts toward "pending expenses"
3. **Given** multiple partial payments have been made, **When** viewing the expense, **Then** the total paid amount is displayed (e.g., "Paid: $80 of $100" after two payments)
4. **Given** an expense is fully paid through multiple partial payments, **When** all payments sum to the total, **Then** the expense is marked as "Paid in Full" and removed from pending expenses
5. **Given** a user views an expense with partial payments, **When** they hover over or click the payment indicator, **Then** they see a breakdown of all partial payments made (dates and amounts)

---

### Edge Cases

- What happens when a user tries to create an expense with a future date? Should be allowed
- How does the system handle editing multiple cells simultaneously if network latency causes save conflicts? the confirm checkbox must appear including all modified cells in the same row, save the changes and update the dashboard totals without a page refresh.
- What happens if a partial payment amount exceeds the total expense amount? the system should prevent the user from making a partial payment that exceeds the total expense amount.
- How does the system handle inline editing when validation fails (e.g., negative amount, empty required field)? the system should prevent the user from making a partial payment that exceeds the total expense amount.
- What happens to pending expenses calculation when an expense is deleted that had partial payments? the system should remove the expense from the pending expenses calculation.
- How does the date picker handle leap years and month-end dates (e.g., selecting Feb 30)? the system should show user a calendar with the correct number of days for each month.

## Requirements

### Functional Requirements

- **FR-001**: System MUST provide a date selector when creating a new expense
- **FR-002**: Date selector MUST allow selecting any date
- **FR-003**: System MUST allow selecting future dates for expenses
- **FR-004**: Expenses MUST be displayed in the dashboard corresponding to their transaction date, not creation date
- **FR-005**: System MUST support inline editing of expense fields (amount, description, category, due date) directly in the table
- **FR-006**: Inline editing MUST provide visual confirmation (checkmark icon) before saving changes of the row
- **FR-007**: System MUST allow canceling inline edits without saving (Escape key or click outside)
- **FR-008**: System MUST update dashboard totals and calculations immediately after confirming inline edits
- **FR-009**: System MUST track partial payments separately from total expense amount
- **FR-010**: System MUST display a visual indicator on expenses that have partial payments
- **FR-011**: Partial payment indicator MUST show all paid amount, total amount and pending amount (e.g., "$30 of $100 paid, $70 pending")
- **FR-012**: Pending expenses calculation MUST only include unpaid portions of expenses with partial payments
- **FR-013**: System MUST maintain a history of all partial payments for each expense
- **FR-014**: System MUST automatically mark an expense as "Paid in Full" when partial payments sum to total amount
- **FR-015**: Inline editing MUST validate field values while editing and before allowing save (positive amounts, required fields, valid dates)
- **FR-016**: Inline editing MUST prevent saving invalid values (e.g., negative amounts, empty required fields, invalid dates)
- **FR-017**: Inline editing MUST prevent saving partial payments that exceed the total expense amount

### Key Entities

- **Expense**: Represents a financial obligation with amount, description, category, transaction date, due date, and payment status
- **Partial Payment**: Represents a payment made toward an expense, with amount and payment date, linked to parent expense
- **Dashboard View**: Aggregated view of expenses filtered by month, showing totals and pending amounts

## Success Criteria

### Measurable Outcomes

- **SC-001**: Users can create expenses with custom dates in under 10 seconds (vs. currently impossible)
- **SC-002**: Users can edit an expense field inline in under 5 seconds (vs. 15+ seconds with modal)
- **SC-003**: 100% of partial payments are correctly reflected in pending expenses calculation
- **SC-004**: Users can identify expenses with partial payments at a glance without clicking or hovering
- **SC-005**: Zero reported cases of date-related bugs after implementation (expenses appearing in wrong month)
- **SC-006**: Editing multiple expenses in succession takes 50% less time than current modal workflow
- **SC-007**: Pending expenses calculation matches manual calculation 100% of the time, including partial payments
