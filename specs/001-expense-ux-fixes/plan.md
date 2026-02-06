# Implementation Plan: Expense Management UX Improvements

**Branch**: `001-expense-ux-fixes` | **Date**: 2026-02-05 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-expense-ux-fixes/spec.md`

## Summary

This feature improves the UX of the expense management module by implementing three critical enhancements: (1) allowing users to create expenses with custom dates (past and future), (2) implementing Excel-style inline editing directly in the expense table, and (3) displaying partial payment status in the table row (partial payments system already exists, just need to show the data).

## Technical Context

**Language/Version**: Go 1.21+ (backend), TypeScript/React 18+ (frontend)  
**Primary Dependencies**: 
- Backend: Go stdlib, PostgreSQL driver (`lib/pq` or `pgx`), existing api-gateway framework
- Frontend: React, TanStack Table (or current table library), react-hook-form, date-fns or similar

**Storage**: PostgreSQL (existing database, likely `postgres-main`)  
**Testing**: Go testing package (backend), Vitest + React Testing Library (frontend)  
**Target Platform**: Web application (browser-based)  
**Project Type**: Web - Modular monolith with separate frontend/backend  
**Performance Goals**: 
- Inline edits save in <500ms (p95)
- Dashboard updates without full page reload
- Date picker renders instantly

**Constraints**: 
- Must not break existing expense functionality
- Must maintain backward compatibility with existing expense data
- Inline editing must work on all modern browsers

**Scale/Scope**: Moderate - affects core expense module used daily by all users

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ I. Modular Monolith Architecture
- **Status**: PASS
- Changes are contained within existing modules (api-gateway backend, frontend)
- No new modules required
- All deploys together in unified container

### ✅ II. Security & Data Privacy
- **Status**: PASS
- No new sensitive data introduced
- Existing JWT authentication applies to all endpoints
- Partial payment amounts are financial data - already encrypted at rest

### ✅ III. Type Safety & Code Quality
- **Status**: PASS
- Backend: Go strict compilation enforced
- Frontend: TypeScript strict mode enforced
- All new types will be properly defined

### ✅ IV. Testing Standards
- **Status**: PASS
- Will include unit tests for backend handlers
- Will include frontend component tests for inline editing
- Target: >80% coverage for new code

### ✅ V. Performance & Scalability
- **Status**: PASS
- Inline editing uses optimistic updates (< 500ms target)
- Dashboard calculations remain client-side (React state)
- Database impact minimal (same queries, just updated data)

### ✅ VI. Observability & Debugging
- **Status**: PASS
- Will add structured logging for partial payments
- Error handling for validation failures
- Client-side errors logged appropriately

### ✅ VII. API Versioning
- **Status**: PASS
- Changes are backward-compatible updates to existing `/api/v1/expenses/*` endpoints
- Adding optional fields (transaction_date, partial_payments)
- No breaking changes to API contracts

**Overall**: ✅ ALL GATES PASS - Proceed with implementation

## Project Structure

### Documentation (this feature)

```text
specs/001-expense-ux-fixes/
├── plan.md              # This file
├── research.md          # Phase 0 output (technical research)
├── data-model.md        # Phase 1 output (database schema changes)
├── quickstart.md        # Phase 1 output (how to test locally)
├── contracts/           # Phase 1 output (API contracts)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# Backend (Go - api-gateway module)
apps/api-gateway/
├── internal/
│   ├── models/
│   │   └── expense.go           # Update: add TransactionDate, PartialPayments fields
│   ├── handlers/
│   │   └── expenses.go          # Update: handle transaction_date, partial payments
│   ├── services/
│   │   ├── expense_service.go   # Update: calculation logic for pending amounts
│   │   └── partial_payment_service.go  # New: handle partial payment operations
│   └── validators/
│       └── expense_validator.go # Update: validate dates, partial payments
└── migrations/
    └── 00X_add_expense_transaction_date_and_partial_payments.sql  # New migration

# Frontend (React/TypeScript)
apps/frontend/
├── src/
│   ├── components/
│   │   ├── expenses/
│   │   │   ├── ExpenseTable.tsx         # Update: inline editing logic
│   │   │   ├── InlineEditCell.tsx       # New: reusable inline edit component
│   │   │   ├── PartialPaymentBadge.tsx  # New: visual indicator for partial payments
│   │   │   └── DatePicker.tsx           # Update: allow any date selection
│   │   └── forms/
│   │       └── ExpenseForm.tsx          # Update: add transaction_date field
│   ├── hooks/
│   │   └── useInlineEdit.ts             # New: hook for inline edit logic
│   ├── services/
│   │   └── expenseService.ts            # Update: API calls for partial payments
│   └── types/
│       └── expense.ts                   # Update: add transactionDate, partialPayments types
└── tests/
    └── components/
        └── expenses/
            ├── ExpenseTable.test.tsx
            └── InlineEditCell.test.tsx

# Shared (if applicable)
packages/go-shared/
└── (no changes expected for this feature)
```

**Structure Decision**: Web application architecture (Option 2). Changes span both backend (api-gateway module in Go) and frontend (React app). Following existing modular monolith pattern where frontend and backend deploy together in unified Docker container with Nginx routing.

## Complexity Tracking

> No constitution violations detected. This section left empty as all gates passed.
