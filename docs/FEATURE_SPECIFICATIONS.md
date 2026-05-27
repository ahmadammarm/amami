# amami - Comprehensive ERP Feature Specifications

## 1. Introduction
The **amami** platform is an Enterprise Resource Planning (ERP) system designed to professionalize mosque operations. It transforms traditional, manual processes into a digital, auditable, and data-driven management framework. This document serves as the authoritative functional specification for all system features.

---

## 2. Administrative & Identity Module

### 2.1 Management Akses & Role (RBAC)
**Objective**: To enforce the Principle of Least Privilege (PoLP), ensuring users only access data relevant to their specific duties.
- **Detailed Workflow**:
    1.  **Define Permissions**: System initializes with atomic codes like `ledger:read`, `zakat:calculate`, `inventory:loan`.
    2.  **Role Template**: Super Admin creates grouping roles:
        - **Super Admin**: Full system governance.
        - **Bendahara**: Financial and Zakat management.
        - **Takmir**: Physical operations and logistics.
        - **Sekretaris**: Administration and documentation.
        - **Jamaah**: Personal data and public information access.
    3.  **Mapping**: Permissions are assigned to roles in the `role_permissions` join table.
- **Key Data**: Role Name, Permission Code, Description.
- **Senior Constraints**:
    - Roles cannot be deleted if they are currently assigned to active users.
    - Permission checks happen at the API Gateway level (Go Middleware).

### 2.2 Manajemen User & Account Lifecycle
**Objective**: To manage the internal workforce (Takmir/Volunteers) and their credentials.
- **Detailed Workflow**:
    1.  **Creation**: Admin adds name, email, and assigns a Role.
    2.  **Activation**: User receives an initial password or an activation link.
    3.  **Audit**: System tracks `last_login_at` and account status (ACTIVE/SUSPENDED).
- **Key Data**: Username, Email (Unique), Password Hash (Argon2), RoleID.

### 2.3 Pengaturan Masjid & SMTP
**Objective**: To establish the system's organizational identity and communication capabilities.
- **Mosque Profile**:
    - **Data**: Name, Address, Phone, Logo, Legal Yayasan ID.
    - **Usage**: Injected into all PDF receipts, financial reports, and the public dashboard.
- **SMTP Configuration**:
    - **Data**: SMTP Host, Port, Username, Password (Encrypted), From Address.
    - **Connectivity**: System performs a "Ping/Handshake" test before saving.

### 2.4 Lupa Password & Pengaturan Akun
**Objective**: Self-service security management for users.
- **Workflow (Forgot Password)**:
    1.  User submits email.
    2.  System validates email existence.
    3.  Backend generates a cryptographically secure token (UUID) with a 15-minute expiration.
    4.  SMTP service sends a secure reset link.
- **Account Settings**: Users can update their display name, phone number, and change passwords (requiring old password validation).

---

## 3. Financial Management Module (The Ledger)

### 3.1 Manajemen Keuangan (Immutable Ledger)
**Objective**: To provide an auditable and transparent record of all mosque cash flows.
- **Detailed Workflow**:
    1.  **Fund Creation**: Establish pots like "Kas Masjid", "Renovasi Gedung", "Anak Yatim".
    2.  **Transaction Entry**: 
        - **Source**: Manual input (e.g., Friday Box).
        - **Type**: CREDIT (Inflow) or DEBIT (Outflow).
        - **Validation**: System checks if fund balance is sufficient for DEBIT.
    3.  **Finalization**: Atomic update to `funds.current_balance`.
- **Senior Constraints**: 
    - **NO UPDATE/DELETE**: Every change must be a new transaction.
    - **Audit Log**: Records `created_by` and `created_at` (server-side only).

---

## 4. Zakat & Social Welfare Module

### 4.1 Kalkulator Zakat (Internal Tool)
**Objective**: To assist Takmir in accurately assessing the religious obligations of the Jamaah.
- **Functionality**: 
    - **Zakat Maal**: Calculates 2.5% of assets if they meet the *Nisab* (threshold).
    - **Zakat Profesi**: Calculates based on monthly salary (with/without deduction options).
    - **Gold Price Integration**: Provision to input or fetch current gold price to update Nisab thresholds.

### 4.2 Database Mustahik & Warga
**Objective**: To maintain a verified registry of aid recipients.
- **Detailed Workflow**:
    1.  **Intake**: Register Citizen Name, NIK, Family Size, and Address.
    2.  **Surveying**: Capture socio-economic indicators (Income, House ownership, Disability).
    3.  **Scoring**: System applies weights to survey data to generate a "Poverty Score" (0-100).
- **Key Data**: Demographic data, Score, Survey History, Last Distribution Date.

### 4.3 Manajemen Zakat Fitrah & Mal (Distribution)
**Objective**: To ensure fair and targeted distribution of collected charity.
- **Workflow**: 
    1.  **Collection**: Record receipt of Zakat (Rice/Cash) from Muzakki.
    2.  **Targeting**: Filter Mustahik by Score and Asnaf category.
    3.  **Execution**: Record the distribution event linked to a ledger transaction.

---

## 5. Qurban Operations Module

### 5.1 Qurban Lifecycle Management
**Objective**: To manage the logistical complexity of animal sacrifice during Eid al-Adha.
- **Features**:
    - **Package Settings**: Define Cow (1/7 share), Whole Cow, or Goat with prices.
    - **Shohibul Management**: Track donor bookings, names for the "Niat," and payment status.
    - **Animal Registry**: Assign tag numbers to physical animals, record weight, and vendor info.
    - **Execution Tracking**: Real-time update from "ALIVE" to "SLAUGHTERED" on D-Day.
- **Business Rule**: System prevents over-booking if animal stock or handling capacity is reached.

---

## 6. Logistics & Communication Module

### 6.1 Manajemen Agenda & Kegiatan
**Objective**: To centralize the mosque’s calendar of spiritual and social events.
- **Workflow**: 
    1.  Create event (e.g., "Kajian Kitab").
    2.  Set Time, Location, and Mubaligh (Speaker).
    3.  Track attendance (optional) and documentation.

### 6.2 Dokumentasi Agenda & Dokumentasi
**Objective**: To build a historical archive of mosque activities for accountability.
- **Function**: Multi-file upload (Images/PDFs) per Agenda item.
- **Reporting**: Generates visual "Activity Reports" for annual LPJ (Laporan Pertanggungjawaban).

### 6.3 Manajemen Inventaris (Assets & Loans)
**Objective**: To track and preserve physical mosque property.
- **Registry**: Comprehensive SKU tracking, purchase history, and condition monitoring (GOOD/REPAIR/DISPOSED).
- **Loan Workflow**: 
    1.  Jamaah requests item (e.g., Funeral Tents).
    2.  Admin creates `AssetLoan` record.
    3.  Dashboard alerts when loan is past `due_date`.

### 6.4 Memo (Internal Communication)
**Objective**: To provide a collaborative environment for the management team.
- **Feature**: 
    - Internal digital bulletins. 
    - **Pinned Memos**: Critical alerts that stay fixed to the Admin Dashboard (e.g., "Water pump broken - Priority Fix").

---

## 7. Dashboard & Utilities

### 7.1 Executive Dashboard
**Objective**: High-level situational awareness.
- **Visuals**:
    - **Finance**: Bar charts for monthly Infaq/Expense trends.
    - **Operations**: Gauges for Qurban stock and Zakat collection progress.
    - **Alerts**: Notifications for Pinned Memos, Overdue Loans, and pending Mustahik surveys.

### 7.2 Utilitas (System Tools)
**Objective**: Maintenance and data safety.
- **Functionality**:
    - System Health Check (Database/SMTP connection status).
    - Database Backup Trigger (Manual/Scheduled).
    - Audit Log Viewer (Track who changed what and when).

---

## 8. Summary of Senior Business Rules

| Module | Constraint | Rationale |
| :--- | :--- | :--- |
| **Auth** | No Self-Deletion | Prevents administrative lockouts. |
| **Finance** | Atomic Writes | Ensures fund balance always matches transaction history. |
| **Zakat** | Data Privacy | Masking of Mustahik NIK/Personal details in public views. |
| **Inventory** | State Validation | Assets in 'REPAIR' status cannot be loaned out. |
| **General** | Soft Deletes | Preservation of data for historical reporting. |
