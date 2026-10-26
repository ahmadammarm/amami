# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

amami (Aplikasi Manajemen Masjid) is a mosque ERP: a Go/Gin/GORM/PostgreSQL API in `backend/` and a Vue 3 + TypeScript SPA in `frontend/`. `docs/` is authoritative for intended behavior: `FEATURE_SPECIFICATIONS.md` (with each feature marked Available / API only / Planned), `RBAC_VIEW_MATRIX.md`, `DATABASE_DOCUMENTATION.md`. Update them when behavior changes.

## Commands

Backend (run from `backend/`):
```bash
go run cmd/main.go          # start API (default :8080); AutoMigrate + RBAC sync on every start
go run cmd/main.go -seed    # insert super admin (admin@amami.org / admin123), settings and sample data, then exit
go build ./...
go vet ./...
go test ./...
go test ./internal/service/finance -run TestReverseTransaction   # single test
```
Config is read from `backend/.env` via godotenv (`pkg/config/config.go`): `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `DB_SSLMODE`, `DB_TIMEZONE` (default `Asia/Jakarta`), `JWT_SECRET`, `APP_PORT`, `UPLOAD_DIR` (default `./uploads`), `MAX_UPLOAD_MB` (default 10). All have defaults. The sample-data part of `-seed` is not idempotent; run it once on an empty database.

Frontend (run from `frontend/`, uses pnpm):
```bash
pnpm install
pnpm dev        # Vite on :5173 (the only origins allowed by backend CORS)
pnpm build      # vue-tsc -b && vite build -- this is the type check
```
There is no lint or test script. API base URL comes from `VITE_API_BASE_URL`, defaulting to `http://localhost:8080/api/v1`.

## Backend architecture

- **Routes live in one place:** `cmd/main.go` defines every route under `/api/v1`, attaches middleware per group/route, and calls handlers from the `di.Config` struct. Adding an endpoint means editing `main.go`, not the module.
- **Per-module layering:** each module (`agenda`, `auth`, `dashboard`, `finance`, `inventory`, `jamaah`, `qurban`, `settings`, `user`, `zakat`) has a folder in each of `internal/handler`, `internal/service`, `internal/repository`, `internal/dto`, and `internal/di`. Layers are interfaces with `NewXxx` constructors. `internal/di/<module>/*_module.go` wires repo -> service -> handler; `internal/di/init.go` builds all modules into `di.Config`. DTOs own the JSON shape and a `ToXxxResponse` mapper; handlers never return `domain` structs.
- **Domain package:** GORM models are in `internal/domain/models.go`; string enums (`TransactionType`, `AssetStatus`, `AgendaStatus`, ...) in `enums.go`; permission codes, role names and the role -> permission grant table in `permissions.go`; ledger/booking rule errors in `errors.go`. Request DTOs repeat enum values in `oneof` binding tags; keep both in sync.
- **Schema changes:** applied by `db.AutoMigrate(...)` in `pkg/database/postgres.go` (no migration files); a new model must be added to that list. That file also holds idempotent pre/post-migration steps (dropping non-partial unique indexes, backfills). AutoMigrate never alters an existing index.
- **Errors and responses:** services return `pkg/apperror` errors (`Validation`, `NotFound`, `Conflict`, `Unauthorized`, `Forbidden`, `Internal`); handlers call `utils.HandleError(c, err)`, which maps the kind to a status and masks unknown errors as 500. Every endpoint responds `{success, message, data, errors}` via `utils.SuccessResponse/CreatedResponse/ErrorResponse/ValidationErrorResponse`. Lists use `utils.Page[T]` (`{items, total, page, limit, total_pages}`) with `utils.ParsePagination`. Read context values with `utils.CurrentUserID/CurrentRoleID`, path params with `utils.ParseUUIDParam/ParseUintParam`. The DB is opened with `TranslateError: true`, so unique/FK violations surface as `gorm.ErrDuplicatedKey`/`gorm.ErrForeignKeyViolated`.
- **Auth and RBAC:**
  - `middleware.RequireAuth()` validates the JWT and puts `user_id`, `role_id` and the token scope in the Gin context.
  - JWT has a scope claim (`utils.ScopeFullAccess` / `utils.ScopePasswordReset`). Users with status `PENDING_PASSWORD_CHANGE` get a password-reset-scoped token; route groups add `middleware.RejectScope(utils.ScopePasswordReset)` until `/auth/change-password` issues a full token.
  - `middleware.RequirePermission(roleRepo, domain.PermXxx)` checks the role's permission codes. Reads use `*:read`, writes `*:manage`/`ledger:write`; "manage" does not imply "read".
  - `domain.RolePermissions` is the single source of grants (mirrors `docs/RBAC_VIEW_MATRIX.md`). `seed.SyncRBAC` makes the DB match it exactly on every server start, and `internal/domain/permissions_test.go` fails if the map drifts from the documented matrix. Permissions are cached per process, so a restart applies changes.
  - `/auth/me` returns the user plus `permissions`; the frontend decides visibility from those codes.
- **Ledger:** `Transaction` rows are append-only. Every posting (manual entries, reversals, zakat) goes through `repository/finance.PostTransaction`, which locks the `Fund` row, rejects overdrafts and updates `CurrentBalance` inside the caller's DB transaction. Corrections use `POST /finance/transactions/:id/reverse`: an opposite entry linked by `ReversalOfID` (unique index, so at most one reversal per transaction; reversals cannot be reversed). Zakat only posts whole-rupiah `IDR` amounts to a fund.
- **Soft deletes:** `User`, `Jamaah`, `QurbanPackage`, `QurbanAnimal`, `Asset` and `Agenda` have `gorm.DeletedAt`; their unique indexes are partial (`WHERE deleted_at IS NULL`). Historical lists preload related rows with `Unscoped()` so names of deleted records still show. Ledger and zakat tables are never deleted.
- **Uploads:** agenda documentation goes through `pkg/storage.Storage` (local disk under `UPLOAD_DIR`, served at `/uploads`). The service sniffs content type (jpeg/png/webp/pdf only), enforces `MAX_UPLOAD_MB`, and removes stored files if the DB write fails.
- **Mustahik score:** `service/jamaah.CalculateMustahikScore` is a pure, weighted 0-100 function (income per capita 50, dependents 20, housing 20, disability 10). `PUT /jamaah/:id/mustahik-survey` stores the survey and the score in one transaction.
- **Tests:** table-driven unit tests with hand-written repository fakes (embed the interface, override the methods under test). No database is needed.
- Rate limiting: 100 req/min on `/api/v1`, 5 req/min on `/auth` (in-memory, `middleware/rate_limiter.go`).

## Frontend architecture

- **Atomic design under `src/components/`:** `atoms/` holds shadcn-vue (new-york, Reka UI) primitives plus the `arch/` mihrab-arch atoms (`ArchBadge`, `ArchFrame`). `components.json` aliases `ui` to `@/components/ui`, but that folder does not exist; generated components were moved to `atoms/` and are imported as `@/components/atoms/<name>`. `pages/` are route views; `templates/DashboardLayout.vue` (authenticated shell) and `templates/AuthLayout.vue` (login, first password); `organisms/AppSidebar.vue` is the nav.
- **Design system:** colour tokens named after mosque materials (`kapur`, `batu`, `jati`, `sajadah`, `kuningan`, `bata`, with `-tua`/`-muda` shades) and fonts (`font-display` Reem Kufi, `font-sans` Plus Jakarta Sans, `font-mono` IBM Plex Mono) live in `src/style.css`; use these classes, not Tailwind's default palette. Pages are built from molecules: `PageHeader`, `StatCard`, `DataPanel` (+ `PaginationBar`), `EmptyState`, `StatusBadge`, `FormDialog` + `FormField`. Inputs use the `.field-control` class, tables the `.tabel` class. All UI copy is Indonesian; formatters are in `src/lib/format.ts`.
- **Notifications:** there are no toasts. `feedback.success/error/confirm` from `src/composables/useFeedback.ts` queue dialogs rendered by `organisms/FeedbackDialog.vue` in `App.vue`; `confirm` resolves to a boolean and replaces `window.confirm`. Pass API errors straight to `feedback.error(title, error)`.
- **Data access:** `src/api/client.ts` exports `request<T>(config)`, which unwraps the response envelope, plus `getErrorMessage(error)` and `assetUrl(path)` for `/uploads/...` links. The Axios instance injects `Bearer` from `localStorage.token` and hard-redirects to `/login` on 401 (except for the login call). One service object per module lives in `src/api/services/<module>.ts`; response types are in `src/types/<module>.ts` (`ApiResponse<T>`, `Page<T>` in `types/api.ts`).
- **Forms:** Zod schemas in `src/schemas/<module>.ts` are the single definition of each request payload (`z.infer`). Pages bind them with VeeValidate (`useForm` + `toTypedSchema` + `defineField`) and show `errors.<field>`.
- **Query keys:** always use `src/api/queryKeys.ts`. Each module has an `all` prefix key; invalidate it after mutations (plus `queryKeys.dashboard.metrics` when totals change).
- **Auth state and RBAC:** the Pinia store `src/stores/auth.ts` holds the token, `requiresPasswordChange` (persisted in localStorage) and the user from `/auth/me`, and exposes `can(permission)` / `canAny(...)`. Permission codes are typed in `src/lib/permissions.ts`. Routes declare `meta.permission` (typed via `RouteMeta` augmentation in `src/router/index.ts`), the sidebar items carry the same `permission`, and pages gate write buttons with `authStore.can(...)`. Never check role names for access; `userRole` is only for display.
- Charts use ApexCharts (`vue3-apexcharts`, options typed as `ApexOptions`, colours taken from the design tokens).

## Conventions

- Enforce RBAC on both sides: a backend `RequirePermission` on every protected endpoint and the matching `meta.permission` / `can(...)` in the UI, per `docs/RBAC_VIEW_MATRIX.md`. To change grants, edit `domain.RolePermissions` and the expected table in `permissions_test.go`.
- Financial records are never updated or deleted; corrections use reversal entries.
- Use GORM soft delete for entities needed in historical reports; make their unique indexes partial.
- Validate with Zod on the frontend and `binding` tags on backend DTOs. No `any` or `@ts-ignore` in frontend code.
- Files use LF line endings.
