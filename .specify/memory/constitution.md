<!--
================================================================================
SYNC IMPACT REPORT - Constitution Update
================================================================================

Version Change: [INITIAL] → 1.0.0 → 1.0.1

Modified Principles:
- [v1.0.0 NEW] I. Microservices Architecture & Modularity
- [v1.0.1 UPDATED] I. Modular Monolith Architecture (corrected from Microservices)
- [NEW] II. Security & Data Privacy (NON-NEGOTIABLE)
- [NEW] III. Type Safety & Code Quality
- [NEW] IV. Testing Standards (minimum 80% coverage)
- [NEW] V. Performance & Scalability
- [NEW] VI. Observability & Debugging
- [NEW] VII. API Versioning & Breaking Changes

Added Sections:
- Technology Constraints (Backend: Go, Frontend: React/TS, DB: PostgreSQL, Cache: Redis)
- Development Workflow (Code Review, Quality Gates, Database Migrations)
- Monorepo Organization (Directory Structure, Dependency Management)
- Governance (Constitution Authority, Amendments, Exceptions)

Removed Sections:
- None (initial version)

Template Alignment Status:
✅ plan-template.md - Aligned (Constitution Check section present)
✅ spec-template.md - Aligned (requirements and success criteria match principles)
✅ tasks-template.md - Aligned (testing, observability, and deployment tasks supported)
✅ checklist-template.md - Ready for constitution-based quality gates
✅ Commands (.claude/commands/*.md) - Compatible with current constitution

Follow-up TODOs:
- None - Constitution is complete and ready for use

Next Steps:
1. Use /speckit.specify to create first feature specification
2. Specifications will automatically enforce Security & Data Privacy principles
3. Plans will verify against constitution check gates
4. Tasks will include required testing (80% coverage) and observability

Ratification Notes:
- Initial constitution version for Financial Resume Monorepo
- v1.0.1: Corrected to modular monolith architecture (all modules deploy together)
- Deployment: Unified Docker container with Nginx reverse proxy + Supervisord
- Modules: users, api-gateway, ai-service, gamification (independent code, unified deployment)
- Security standards appropriate for financial application
- Type safety enforced via Go compilation + TypeScript strict mode
- Performance targets aligned with web application standards
================================================================================
-->

# Financial Resume Monorepo Constitution

## Core Principles

### I. Modular Monolith Architecture

**Non-Negotiable Rules:**
- Code MUST be organized into independent modules (users, api-gateway, ai-service, gamification)
- Modules communicate via well-defined internal APIs (in-process or HTTP)
- All modules deploy together as a unified backend in a single Docker container
- Nginx reverse proxy routes external requests to appropriate modules
- Supervisord manages all module processes within the container
- Shared code lives in `packages/go-shared` and MUST be versioned
- Each module SHOULD own its data, but can share databases for simplicity
- API contracts MUST be documented with OpenAPI/Swagger

**Rationale:** Modular monolith provides clear separation of concerns and maintainability while keeping deployment simple. All modules deploy together, reducing operational complexity while maintaining code organization. This is ideal for early-stage applications where the overhead of true microservices isn't justified, but clean boundaries are still essential for long-term maintainability.

### II. Security & Data Privacy (NON-NEGOTIABLE)

**Non-Negotiable Rules:**
- JWT tokens for authentication with proper expiration
- 2FA MUST be enforced for sensitive financial operations
- All financial data MUST be encrypted at rest and in transit
- Password hashing with bcrypt (minimum cost factor 12)
- No sensitive data in logs or error messages
- Environment variables for all secrets - never hardcoded
- HTTPS only in production

**Rationale:** Financial data is highly sensitive. Security breaches can cause financial loss and destroy user trust. Privacy and security are non-negotiable in fintech applications.

### III. Type Safety & Code Quality

**Non-Negotiable Rules:**
- Go: Strict compilation with `go vet` and `golangci-lint` passing
- TypeScript: Strict mode enabled (`strict: true` in tsconfig)
- No `any` types in TypeScript without explicit justification
- All public APIs MUST have comprehensive comments
- Code reviews MUST verify type safety before merge

**Rationale:** Type safety catches bugs at compile time, improves developer experience with autocomplete, and serves as living documentation. Critical for financial calculations where errors have real monetary consequences.

### IV. Testing Standards

**Non-Negotiable Rules:**
- Minimum 80% code coverage for business logic
- Unit tests for all service handlers and business logic
- Integration tests for API endpoints
- Frontend: Component tests for critical UI flows (checkout, auth, transactions)
- Financial calculations MUST have comprehensive test cases including edge cases
- Tests MUST run in CI/CD pipeline - failing tests block deployment

**Rationale:** Financial applications require high reliability. Comprehensive testing catches bugs before they affect users' money. Test coverage ensures confidence when refactoring.

### V. Performance & Scalability

**Performance Standards:**
- API response time: p95 < 200ms for read operations
- API response time: p95 < 500ms for write operations
- Frontend: First Contentful Paint < 1.5s
- Frontend: Time to Interactive < 3s
- Database queries MUST use indexes - no full table scans
- Cache frequently accessed data (Redis)

**Rationale:** User experience depends on performance. Slow financial apps frustrate users. Performance standards ensure the app scales as user base grows.

### VI. Observability & Debugging

**Requirements:**
- Structured logging with correlation IDs across services
- Log levels: DEBUG, INFO, WARN, ERROR consistently used
- Metrics collection for: request count, latency, error rate per service
- Health check endpoints (`/health`) for all services
- Request tracing across microservices
- Error tracking with stack traces (no sensitive data)

**Rationale:** Distributed microservices are complex to debug. Observability tools are essential for diagnosing issues in production, especially financial transactions.

### VII. API Versioning & Breaking Changes

**Versioning Rules:**
- API versioning: `/api/v1/`, `/api/v2/` etc.
- Semantic versioning for shared packages: MAJOR.MINOR.PATCH
- Breaking changes require new major version
- Deprecation warnings MUST be issued at least 1 month before removal
- Backward compatibility MUST be maintained within same major version

**Rationale:** Multiple services depend on APIs. Uncoordinated breaking changes cause system-wide failures. Versioning enables safe evolution of the platform.

## Technology Constraints

### Backend Stack

- **Language**: Go 1.21+
- **Framework**: stdlib-based with minimal dependencies
- **Database**: PostgreSQL 15+ for persistence
- **Cache**: Redis 7+ for session and data caching
- **Authentication**: JWT with refresh tokens
- **API Documentation**: Swagger/OpenAPI 3.0

### Frontend Stack

- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite
- **State Management**: React Query + Context API
- **Styling**: TailwindCSS or CSS Modules
- **Testing**: Vitest + React Testing Library

### Infrastructure

- **Containerization**: Docker with multi-stage builds
- **Orchestration**: Docker Compose (local), Render.com (production)
- **Reverse Proxy**: Nginx
- **Process Management**: Supervisord
- **CI/CD**: GitHub Actions (when configured)

## Development Workflow

### Code Review Requirements

- All changes MUST go through pull request review
- At least 1 approval required before merge
- Reviews MUST verify:
  - Code follows constitution principles
  - Tests are included and passing
  - No security vulnerabilities introduced
  - API contracts are maintained or properly versioned
  - Documentation is updated

### Quality Gates

Before deployment, MUST verify:
1. All tests passing (unit + integration)
2. Linter checks passing (golangci-lint, eslint)
3. Type checking passing (Go compilation, TypeScript)
4. Security scan passing (no critical vulnerabilities)
5. Code coverage ≥ 80% for new/modified code

### Database Migrations

- Migrations MUST be versioned and sequential
- Migrations MUST be reversible (up/down scripts)
- Test migrations on staging before production
- Never alter production data directly - always use migrations

## Monorepo Organization

### Directory Structure Rules

- `apps/` - Deployable applications (services, frontend)
- `packages/` - Shared libraries and utilities
- `infrastructure/` - Docker, configs, deployment scripts
- `specs/` - Feature specifications (spec-kit managed)
- `docs/` - Architecture docs, runbooks, ADRs

### Dependency Management

- Go modules per service with shared workspace (`go.work`)
- Frontend: pnpm workspace for package management
- Shared packages MUST declare explicit dependencies
- Circular dependencies are prohibited

## Governance

### Constitution Authority

- This constitution supersedes all other development practices
- All code reviews MUST verify compliance with constitution
- Violations MUST be documented and resolved before merge

### Amendments

- Constitution changes require:
  1. Documented proposal with rationale
  2. Team discussion and consensus
  3. Version bump with migration plan if needed
  4. Update to dependent templates and workflows

### Exceptions

- Exceptions to constitution require explicit justification
- Document exceptions as Architectural Decision Records (ADRs)
- Time-bound exceptions MUST have remediation plan

**Version**: 1.0.1 | **Ratified**: 2026-02-05 | **Last Amended**: 2026-02-05
