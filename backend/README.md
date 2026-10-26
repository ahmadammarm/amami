# Amami — Backend

REST API for amami, written in Go with Gin, GORM and PostgreSQL. It enforces role-based access on every endpoint and keeps the financial ledger append-only.

## Stack

| Concern | Library |
|---|---|
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL via [GORM](https://gorm.io) (`TranslateError` on, schema by `AutoMigrate`) |
| Auth | JWT (HS256, `golang-jwt/jwt`), bcrypt password hashing |
| Validation | `go-playground/validator` binding tags on request DTOs |
| Logging | Uber Zap |
| Rate limiting | In-memory: 100 requests/min on `/api/v1`, 5 requests/min on `/api/v1/auth` |

## Getting started

**Requirements:** Go 1.25+ and a running PostgreSQL server.

```bash
cp .env.example .env     # replace the example values with your own
go mod download
go run cmd/main.go -seed # once, on an empty database
go run cmd/main.go       # http://localhost:8080
```

### Configuration

Read from `.env` (or the process environment) by `pkg/config/config.go`:

| Variable | Default | Purpose |
|---|---|---|
| `DB_HOST`, `DB_PORT` | `localhost`, `5432` | PostgreSQL address |
| `DB_USER`, `DB_PASSWORD` | `postgres`, `postgres` | PostgreSQL credentials |
| `DB_NAME` | `amami_db` | Database name |
| `DB_SSLMODE` | `disable` | PostgreSQL `sslmode` |
| `DB_TIMEZONE` | `Asia/Jakarta` | Session time zone |
| `JWT_SECRET` | `super-secret-key-change-me` | Token signing key. **Always override.** |
| `APP_PORT` | `8080` | HTTP port |
| `UPLOAD_DIR` | `./uploads` | Where agenda documentation files are stored |
| `MAX_UPLOAD_MB` | `10` | Per-file upload limit |

CORS allows `http://localhost:5173` and `http://127.0.0.1:5173` (the Vite dev server); change `cmd/main.go` for other origins.

### What happens on start

Every `go run cmd/main.go`:

1. runs `AutoMigrate` for all models, plus idempotent pre/post steps in `pkg/database/postgres.go` (partial unique indexes for soft-deleted rows, data backfills);
2. syncs roles, permissions and grants to exactly match `internal/domain/permissions.go` (`seed.SyncRBAC`).

So changing who can do what means editing `domain.RolePermissions` and restarting; no re-seed is needed.

### Seeding

`go run cmd/main.go -seed` creates the super administrator (`admin@amami.org` / `admin123`), default mosque settings, three funds, and sample jamaah, transactions, qurban data and assets, then exits.

The sample data is **not idempotent**: run it once, on an empty database.

## Project layout

```
cmd/main.go            Entry point and the full route table
database/seed/         SyncRBAC and the sample-data seeder
internal/domain/       GORM models, enums, permission grants, domain errors
internal/dto/<module>/ Request/response shapes and mappers
internal/handler/      HTTP handlers (bind, call service, write response)
internal/service/      Business rules
internal/repository/   Database access and transactions
internal/middleware/   Auth, permission, rate-limit and body-limit middleware
internal/di/           Wiring: repository -> service -> handler per module
pkg/                   apperror, config, database, logger, storage, utils
```

Modules: `agenda`, `auth`, `dashboard`, `finance`, `inventory`, `jamaah`, `qurban`, `settings`, `user`, `zakat`.

## API conventions

All routes live under `/api/v1`. Every response uses one envelope:

```json
{ "success": true, "message": "Transactions fetched successfully", "data": { } }
{ "success": false, "message": "Validation failed", "errors": { "amount": "amount must be gt 0" } }
```

- Lists are paginated with `?page=` and `?limit=` (max 100) and return `{ items, total, page, limit, total_pages }`.
- Services return typed errors (`pkg/apperror`) that map to status codes: validation `400`, unauthorized `401`, forbidden `403`, not found `404`, conflict `409`. Unexpected errors are logged and returned as a generic `500`.
- Send the token as `Authorization: Bearer <token>`. Accounts awaiting their first password change get a restricted token that only works for `/auth/me` and `/auth/change-password`.

### Endpoints

| Area | Endpoints | Permission |
|---|---|---|
| Auth | `POST /auth/login`, `GET /auth/me`, `POST /auth/change-password` | signed in (login is public) |
| Users | `GET /users`, `GET /users/roles`, `POST /users/invite`, `PATCH /users/:id/status`, `DELETE /users/:id` | `user:manage` |
| Settings | `GET /settings/profile` | public |
| | `PUT /settings/profile`, `GET/PUT /settings/smtp`, `POST /settings/smtp/test`, `GET /settings/health`, `GET /settings/audit` | `settings:manage` |
| Finance | `GET /finance/funds`, `GET /finance/transactions` | `ledger:read` |
| | `POST /finance/funds`, `POST /finance/transactions`, `POST /finance/transactions/:id/reverse` | `ledger:write` |
| Dashboard | `GET /dashboard/metrics` | signed in |
| Jamaah | `GET/POST /jamaah`, `GET/PUT /jamaah/:id`, `GET /jamaah/:id/mustahik-survey` | `jamaah:manage` |
| | `PUT /jamaah/:id/mustahik-survey` | `mustahik:survey` |
| Zakat | `GET /zakat/donations`, `GET /zakat/distributions` | `zakat:read` |
| | `POST /zakat/donations`, `POST /zakat/distributions` | `zakat:manage` |
| Qurban | `GET /qurban/packages`, `GET /qurban/bookings`, `GET /qurban/animals` | `qurban:read` |
| | create/update/delete packages and animals, `POST /qurban/bookings` | `qurban:manage` |
| Inventory | `GET /inventory/assets` | `inventory:read` |
| | create/update/delete assets, `POST /inventory/assets/:id/loans`, `PUT /inventory/assets/loans/:loanId/return` | `inventory:manage` |
| Agenda | `GET /agenda`, `GET /agenda/:id/documentation` | `agenda:read` |
| | create/update/delete agendas, `POST /agenda/:id/documentation` (multipart `files`), `DELETE /agenda/:id/documentation/:docId` | `agenda:manage` |

Uploaded files are served from `/uploads/...`.

## Business rules worth knowing

- **Append-only ledger.** Transactions are never updated or deleted. Every posting goes through `repository/finance.PostTransaction`, which locks the fund row, rejects overdrafts and updates the balance in the same database transaction. A correction is `POST /finance/transactions/:id/reverse`: an opposite entry linked by `reversal_of_id`. A transaction can be reversed once, and a reversal cannot be reversed.
- **Zakat and the ledger.** Only whole-rupiah `IDR` amounts can be posted to a fund; rice (`KG`) and other goods are recorded without a fund. Distributions only go to jamaah marked as mustahik.
- **Mustahik score.** `service/jamaah.CalculateMustahikScore` turns survey answers into a 0-100 score: income per household member (50 points), dependents (20), housing (20), disability (10).
- **Soft deletes.** Users, jamaah, qurban packages and animals, assets and agendas are soft-deleted; their unique fields are only unique among live rows. Ledger and zakat records are never deleted.
- **Uploads.** Documentation files are checked by content (JPEG, PNG, WebP, PDF only) and size, stored under random names, and removed again if the database write fails.

## Tests

Unit tests use hand-written repository fakes and need no database:

```bash
go test ./...
go test ./internal/service/finance -run TestReverseTransaction
```

`internal/domain/permissions_test.go` fails if the grant table drifts from `docs/RBAC_VIEW_MATRIX.md`.
