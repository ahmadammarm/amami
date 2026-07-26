# Amami (Amanah Management & Mosque Intelligence)

Amami is a comprehensive Enterprise Resource Planning (ERP) platform designed to professionalize mosque operations. It transforms traditional, manual processes into a digital, auditable, and data-driven management framework. The system implements a strict Role-Based Access Control (RBAC) model to govern user permissions across various modules.

## Project Structure

The project is structured as a monorepo, divided into a Go backend and a Vue 3 frontend:

- `/backend` - The REST API server built with Go, Gin, and PostgreSQL.
- `/frontend` - The web client built with Vue 3, Vite, and Tailwind CSS.
- `/docs` - Contains authoritative project documentation, including feature specifications, RBAC view matrix, and roadmap.

## Key Features

- **Strict Role-Based Access Control (RBAC):** Permissions are strictly enforced based on roles (SUPER_ADMIN, BENDAHARA, TAKMIR, SEKRETARIS, JAMAAH) at both the API and UI levels.
- **Immutable Financial Ledger:** Financial records follow a strict "No Update / No Delete" rule. Corrections require explicit "Reversal Entries" to maintain absolute financial integrity and transparency.
- **Zakat Management:** Full tracking of Zakat Fitrah and Zakat Maal collection and distribution, with atomic ledger integration for cash donations.
- **Jamaah Database:** Community management system with dynamic tagging (e.g., Mustahik, Donatur).
- **Interactive Dashboards:** Role-specific dynamic dashboards that display aggregated metrics for finance, zakat, jamaah, and more.

## Development Setup

### Prerequisites

- Go 1.25+
- Node.js (v18+) and pnpm
- PostgreSQL

### Running the Application

To run the application locally, you will need to start both the backend and frontend servers in separate terminal instances. Please refer to the specific README files in the `/backend` and `/frontend` directories for detailed setup instructions.

## License

This project is proprietary.
