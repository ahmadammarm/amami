# Amami - Backend

The backend of Amami is a high-performance RESTful API server built with Go, designed to handle complex business logic, enforce strict Role-Based Access Control (RBAC), and maintain an immutable financial ledger.

## Technology Stack

- **Language:** Go 1.25+
- **Framework:** Gin (HTTP Web Framework)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT, Argon2 (Cryptography)
- **Logging:** Uber Zap
- **Validation:** go-playground/validator

## Architecture

The backend follows a clean architecture pattern to separate concerns and ensure maintainability:

1. **Handler (Controllers):** Parses HTTP requests, validates inputs, and formats responses.
2. **Service (Business Logic):** Executes business rules, orchestrates operations, and enforces constraints.
3. **Repository (Data Access):** Handles database queries, ORM interactions, and atomic transactions.
4. **Domain (Models):** Defines data structures, GORM mapping, and interfaces.
5. **Dependency Injection (DI):** Wires the components together in the `internal/di` package.

### Immutable Ledger & Concurrency

Financial transactions (like Zakat collection) are executed using atomic database transactions. The system utilizes explicit database locking (`clause.Locking{Strength: "UPDATE"}`) to prevent race conditions when multiple concurrent requests attempt to update the same financial fund balance. 

Records in the financial ledger cannot be updated or deleted. Any adjustments require a compensatory transaction.

## Getting Started

### Configuration

1. Ensure PostgreSQL is running.
2. Create a `.env` file in the `backend` directory with the following variables:
   - `DB_HOST`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT`
   - `JWT_SECRET`
   - `APP_PORT`

### Database Seeding

To initialize the database with required permissions, roles, the super admin account, and default configurations, run the seeder:

```bash
go run cmd/main.go -seed
```

### Running the Server

Start the API server in development mode:

```bash
go run cmd/main.go
```

The server will typically start on port 8080 (or as defined in your environment configuration).
