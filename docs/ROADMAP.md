# amami ERP - Comprehensive Project Roadmap

## 1. Project Vision
**amami** (Amanah Management & Mosque Intelligence) is designed to professionalize mosque operations through data integrity, financial transparency, and community-centric logistics.

---

## 2. Development Milestones

### Milestone 1: User Onboarding & Identity (Current Phase)
*   **Goal**: Enable the Super Admin to build and manage the Takmir workforce.
*   **Key Features**:
    *   **User Invitation**: Admin-led creation of staff accounts with specific RBAC roles.
    *   **Account Lifecycle**: Activation, suspension, and role updates.
    *   **Profile Linkage**: Automatic 1:1 mapping between `User` accounts and physical `Jamaah` records.
    *   **Mosque Profile**: Centralized branding (Logo, Name, Yayasan ID) for report generation.
    *   **SMTP Integration**: System-wide email capabilities for alerts and invitations.

### Milestone 2: The Immutable Ledger (Core Finance)
*   **Goal**: Establish a bulletproof "Single Source of Truth" for all cash flows.
*   **Key Features**:
    *   **Fund Management**: CRUD for restricted pots (e.g., Orphan Fund vs. General Fund).
    *   **Transaction Engine**: Multi-table atomic operations for `CREDIT` and `DEBIT`.
    *   **Strict Immutability**: Enforcement of "No Update/No Delete" rules at the service layer.
    *   **Reversal Logic**: Standardized workflow for correcting data entry errors via offsetting entries.

### Milestone 3: Community Registry & Mustahik Scoring
*   **Goal**: Build a data-driven foundation for social welfare.
*   **Key Features**:
    *   **Fuzzy-Searchable Jamaah Database**: High-performance registry with GIN/Trigram indexing.
    *   **Mustahik Profiling**: Collection of socio-economic indicators (Income, Dependents, Housing).
    *   **Scoring Algorithm**: Automated calculation of "Poverty Scores" to prioritize aid objectively.

### Milestone 4: Zakat Lifecycle (Collection & Distribution)
*   **Goal**: Automate and audit the religious obligation of charity.
*   **Key Features**:
    *   **Zakat Calculators**: Maal, Profession, and Fitrah assessment tools.
    *   **Collection Tracking**: Digital receipts and automatic ledger entries for received Zakat.
    *   **Targeted Distribution**: Aid distribution workflows that filter recipients by poverty score and Asnaf category.

### Milestone 5: Qurban Operational Management
*   **Goal**: Manage the high-intensity logistics of Eid al-Adha.
*   **Key Features**:
    *   **Package Bookings**: Shohibul (donor) management with shares and intentions (Niat).
    *   **Animal Inventory**: Tracking physical stock (Alive -> Slaughtered -> Distributed).
    *   **D-Day Dashboard**: Real-time slaughtering progress and distribution log.

### Milestone 6: Logistics, Assets & Events
*   **Goal**: Preserve mosque property and manage the spiritual calendar.
*   **Key Features**:
    *   **Inventory SKU Tracking**: Condition monitoring (Good/Repair/Disposed) for mosque property.
    *   **Asset Loan Workflow**: Request-to-return lifecycle for borrowing mosque equipment.
    *   **Agenda & LPJ**: Event scheduling with Mubaligh (speaker) management and activity photo documentation.

### Milestone 7: Intelligence, Analytics & Utilities
*   **Goal**: High-level situational awareness and system maintenance.
*   **Key Features**:
    *   **Executive Dashboard**: Financial trend visualizations (Inflow vs. Outflow charts).
    *   **Pinned Memos**: Internal management bulletin board.
    *   **System Health**: Real-time monitoring of DB/SMTP connectivity and automated backup triggers.

---

## 3. Engineering Standards (Senior Level)

### 3.1 Architectural Principles
*   **Layered Separation**: Repository (Data) -> Service (Logic) -> Handler (HTTP).
*   **Dependency Injection**: Modularized via `internal/di/` to ensure clean testing and decoupling.
*   **DRY Utilities**: Centralized logging (Zap), validation (Validator v10), and response handling.

### 3.2 Security & Integrity
*   **RBAC Middleware**: Permission-level checks on every sensitive endpoint.
*   **Rate Limiting**: Brute-force protection on authentication and financial submissions.
*   **Audit Logging**: Every write operation must record the `UserID` and timestamp.
*   **Database Atomicity**: Use of SQL Transactions for all multi-step business flows.

### 3.3 Reliability
*   **Environment-Driven**: 12-Factor App methodology using `.env` for all configurations.
*   **Soft Deletes**: Use of GORM's soft-delete where historical reporting is required.
*   **Automated Migrations**: Schema updates managed via Go code to ensure environment parity.
