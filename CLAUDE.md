# financial-resume-monorepo Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-02-10

## Active Technologies
- Go 1.23+ (workspace go.work at 1.24.4) + Gin 1.10 (HTTP), GORM 1.25 (ORM), zerolog (logging), golang-jwt/v5 (JWT), google/uuid (IDs), pquerna/otp (TOTP/2FA), x/crypto/bcrypt (password hashing), stretchr/testify (testing) (001-migrate-auth-module)
- PostgreSQL 15 (single database, shared with other modules via GORM) (001-migrate-auth-module)
- Go 1.23+ (workspace go.work at 1.24.4) + GORM 1.25 (ORM + AutoMigrate), zerolog (logging), google/uuid (ID generation), lib/pq (raw SQL where needed) (002-db-consolidation)
- PostgreSQL 15 — consolidating from 3 databases to 1 (`financial_resume`) (002-db-consolidation)

- Go 1.23+ (workspace go.work at 1.24.4) + Gin 1.10 (HTTP), GORM 1.25 (ORM), zerolog (structured logging), godotenv (env loading), golang-jwt/v5 (JWT), google/uuid, stretchr/testify (testing) (001-monolith-foundation)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.23+ (workspace go.work at 1.24.4)

## Code Style

Go 1.23+ (workspace go.work at 1.24.4): Follow standard conventions

## Recent Changes
- 002-db-consolidation: Added Go 1.23+ (workspace go.work at 1.24.4) + GORM 1.25 (ORM + AutoMigrate), zerolog (logging), google/uuid (ID generation), lib/pq (raw SQL where needed)
- 001-migrate-auth-module: Added Go 1.23+ (workspace go.work at 1.24.4) + Gin 1.10 (HTTP), GORM 1.25 (ORM), zerolog (logging), golang-jwt/v5 (JWT), google/uuid (IDs), pquerna/otp (TOTP/2FA), x/crypto/bcrypt (password hashing), stretchr/testify (testing)

- 001-monolith-foundation: Added Go 1.23+ (workspace go.work at 1.24.4) + Gin 1.10 (HTTP), GORM 1.25 (ORM), zerolog (structured logging), godotenv (env loading), golang-jwt/v5 (JWT), google/uuid, stretchr/testify (testing)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
