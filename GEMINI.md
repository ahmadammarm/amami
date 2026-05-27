# amami (Amanah Management & Mosque Intelligence)

**amami** is a comprehensive Enterprise Resource Planning (ERP) platform designed to professionalize mosque operations. It transforms traditional, manual processes into a digital, auditable, and data-driven management framework.

## 🏗️ Project Architecture

The project is structured as a monorepo with a Go backend and a Vue 3 frontend.

### Backend (`/backend`)
- **Language:** Go 1.24+
- **Framework:** Gin (HTTP Web Framework)
- **ORM:** GORM (PostgreSQL)
- **Logging:** Zap
- **DI:** Manual dependency injection orchestrated in `internal/di/init.go`.
- **Pattern:** Layered architecture:
  - `handler/`: HTTP request handling and response formatting.
  - `service/`: Core business logic.
  - `repository/`: Data access layer (GORM).
  - `domain/`: Shared models and entities.
  - `dto/`: Data Transfer Objects for API requests/responses.

### Frontend (`/frontend`)
- **Framework:** Vue 3 (Composition API with `<script setup>`)
- **Build Tool:** Vite
- **Styling:** Tailwind CSS 4
- **Components:** Headless UI via `reka-ui` and custom atomic design.
- **Form Handling:** VeeValidate + Zod
- **Structure:** Atomic Design (Atoms, Molecules, Organisms, Pages, Templates).

## 🚀 Getting Started

### Prerequisites
- Go 1.24 or higher
- Node.js (v18+) and `pnpm`
- PostgreSQL

### Backend Setup
1. Navigate to the backend directory: `cd backend`
2. Create a `.env` file (refer to `pkg/config/config.go` for expected variables).
3. Install dependencies: `go mod download`
4. Run migrations and seed the database:
   ```bash
   go run cmd/main.go -seed
   ```
5. Start the server:
   ```bash
   go run cmd/main.go
   ```

### Frontend Setup
1. Navigate to the frontend directory: `cd frontend`
2. Install dependencies: `pnpm install`
3. Start the development server:
   ```bash
   pnpm dev
   ```

## 🛠️ Development Conventions

### Backend
- **RBAC:** All protected endpoints must use `middleware.RequirePermission`.
- **Immutability:** Financial transactions follow a "No Update/No Delete" rule (see `docs/FEATURE_SPECIFICATIONS.md`).
- **Soft Deletes:** Use GORM's soft delete for entities requiring historical preservation.
- **DI:** New modules should be integrated into `internal/di/init.go`.

### Frontend
- **Atomic Design:** Follow the existing component hierarchy in `src/components`.
- **Type Safety:** Maintain strict TypeScript typing for all components and composables.
- **Styling:** Use Tailwind 4 utility classes; prefer Vanilla CSS for complex custom components.

## 📂 Key Files & Directories
- `docs/`: Comprehensive project documentation (Roadmap, Specs, RBAC Matrix).
- `backend/cmd/main.go`: Application entry point and route definitions.
- `backend/internal/domain/models.go`: Central source of truth for database schemas.
- `frontend/src/App.vue`: Main frontend entry point.
- `frontend/src/components/atoms/`: Reusable primitive UI components.

## 🗺️ Roadmap & Specifications
Refer to the `docs/` directory for detailed functional specifications and the project roadmap:
- `docs/FEATURE_SPECIFICATIONS.md`: Authoritative functional specs for all modules.
- `docs/ROADMAP.md`: Project vision and development milestones.
