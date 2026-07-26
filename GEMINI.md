# amami (Aplikasi Manajemen Masjid)

## Project Overview
**amami** is an Enterprise Resource Planning (ERP) platform designed to professionalize mosque operations. It transforms traditional, manual processes into a digital, auditable, and data-driven management framework. The system implements a strict Role-Based Access Control (RBAC) model to govern user permissions across various modules.

The project is structured as a monorepo, divided into a Go backend and a Vue 3 frontend.

## Architecture & Technology Stack

### Backend (`/backend`)
- **Language:** Go 1.25+
- **Framework:** Gin (HTTP Web Framework)
- **Database / ORM:** PostgreSQL with GORM
- **Authentication:** JWT, Argon2 (Crypto)
- **Logging:** Uber Zap
- **Validation:** go-playground/validator

### Frontend (`/frontend`)
- **Framework:** Vue 3 (Composition API, `<script setup>`)
- **Language:** TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS 4
- **State Management:** Pinia
- **Routing:** Vue Router
- **Data Fetching:** Vue Query (`@tanstack/vue-query`), Axios
- **Form Handling:** VeeValidate + Zod
- **UI Components:** Reka UI (Headless UI), Unovis (Charts)

### Documentation (`/docs`)
Contains authoritative project documentation, including:
- `FEATURE_SPECIFICATIONS.md`: Detailed functional specs for all modules (Auth, Ledger, Zakat, Qurban, Inventory, etc.).
- `RBAC_VIEW_MATRIX.md`: Defines role permissions (SUPER_ADMIN, BENDAHARA, TAKMIR, SEKRETARIS, JAMAAH) and their corresponding frontend visibility.
- `DATABASE_DOCUMENTATION.md`: Database schema and constraints.
- `ROADMAP.md`: Project vision and development milestones.

## Building and Running

### Prerequisites
- Go 1.25+
- Node.js (v18+) and `pnpm`
- PostgreSQL

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Create and configure your `.env` file based on configuration expectations (e.g., DB credentials, JWT secret).
3. Download dependencies:
   ```bash
   go mod download
   ```
4. Run migrations and seed data (if a seed script is available):
   ```bash
   go run cmd/main.go -seed
   ```
5. Start the server:
   ```bash
   go run cmd/main.go
   ```

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies using pnpm:
   ```bash
   pnpm install
   ```
3. Start the development server:
   ```bash
   pnpm dev
   ```
4. Build for production:
   ```bash
   pnpm build
   ```

## Development Conventions
- **Strict RBAC:** Both backend API endpoints and frontend routes/UI elements must enforce Role-Based Access Control as defined in the `RBAC_VIEW_MATRIX.md`. Frontend implements hard router guards for protected pages.
- **Immutable Ledger:** Financial records follow a strict "No Update / No Delete" rule. Corrections require explicit "Reversal Entries."
- **Data Preservation:** Use GORM's soft delete functionality for entities where historical reporting is necessary.
- **Component Architecture:** The frontend follows an Atomic Design pattern (Atoms, Molecules, Organisms, Pages, Templates) located in `src/components`.
- **Validation:** Always use Zod for frontend form validation and GORM/Gin validation tags on the backend to ensure data integrity.
