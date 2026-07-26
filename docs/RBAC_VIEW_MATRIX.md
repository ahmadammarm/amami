# amami ERP - RBAC & Frontend Access Matrix

This document defines the relationship between backend permissions and frontend visibility. The application uses a **Unified Dashboard** approach where the UI adapts dynamically to the authenticated user's role.

---

## 1. Role Definitions

### 1.1 SUPER_ADMIN
*   **Authority**: Total system control.
*   **Primary Objective**: System maintenance, user governance, and final auditing.
*   **Key Traits**: Only role capable of managing other users and global settings.

### 1.2 BENDAHARA (Treasurer)
*   **Authority**: Financial and Zakat operations.
*   **Primary Objective**: Accuracy of the ledger and fair distribution of charity.
*   **Key Traits**: Access to sensitive financial reports and current balances.

### 1.3 TAKMIR (Operational Leader)
*   **Authority**: Broad operational management.
*   **Primary Objective**: Overall mosque operational success (Zakat, Qurban, Inventory, etc.).
*   **Key Traits**: Access to almost all operational data, excluding User/System management.

### 1.4 SEKRETARIS (Secretary)
*   **Authority**: Administration and Documentation.
*   **Primary Objective**: Maintaining the community database, agendas, and historical archives.
*   **Key Traits**: Focus on Jamaah data, Agenda documentation, and Memos.

### 1.5 JAMAAH (Community Member)
*   **Authority**: Personal data and Public information.
*   **Primary Objective**: Self-service transparency and community participation.
*   **Key Traits**: Can view their own contribution history and public agendas.

---

## 2. Page-Level Access Matrix

| Sidebar Module | Sub-Page | Super Admin | Bendahara | Takmir | Sekretaris | Jamaah |
| :--- | :--- | :---: | :---: | :---: | :---: | :---: |
| **Dashboard** | Executive Summary | Read | Read | Read | Read | Hidden |
| | Operational Metrics | Read | Read | Read | Read | Read |
| **User Management**| User List / Invite | Full | Hidden | Hidden | Hidden | Hidden |
| | Role/Permissions View | Full | Hidden | Hidden | Hidden | Hidden |
| **Finance** | Immutable Ledger | Read | Full | Read | Read | Hidden |
| | Fund Management | Read | Full | Hidden | Hidden | Hidden |
| | Financial Reports | Read | Full | Read | Read | Hidden |
| **Zakat** | Zakat Calculator | Full | Full | Full | Hidden | Read |
| | Muzakki (Donors) | Full | Full | Full | Read | Personal |
| | Mustahik Distribution | Full | Full | Full | Read | Hidden |
| **Jamaah** | Community Registry | Full | Read | Full | Full | Hidden |
| | Mustahik Scoring | Read | Read | Full | Full | Hidden |
| **Qurban** | Bookings (Shohibul) | Read | Read | Full | Full | Personal |
| | Animal Inventory | Read | Hidden | Full | Read | Read |
| **Inventory** | Asset Registry | Read | Hidden | Full | Read | Read |
| | Loan Management | Read | Hidden | Full | Full | Personal |
| **Agenda** | Agenda / Calendar | Full | Read | Full | Full | Read |
| | Activity Documentation | Full | Read | Full | Full | Read |
| **Settings** | Mosque Profile | Full | Hidden | Hidden | Full | Read |
| | SMTP/System Config | Full | Hidden | Hidden | Hidden | Hidden |

---

## 3. UI Component-Level Constraints (In-Page)

### 3.1 Financial Ledger Page
*   **Bendahara**: Sees "Add Transaction" and "Reversal Entry" buttons.
*   **Super Admin / Takmir / Sekretaris**: Sees the ledger list but cannot add or reverse entries.

### 3.2 User Management
*   **Super Admin**: Only role that can see the "Invite User" button.
*   **Everyone Else**: Page is entirely hidden.

---

## 4. Summary Table of Permissions
| Permission Code | Super Admin | Bendahara | Takmir | Sekretaris | Jamaah |
| :--- | :---: | :---: | :---: | :---: | :---: |
| `user:manage` | ✅ | ❌ | ❌ | ❌ | ❌ |
| `ledger:read` | ✅ | ✅ | ✅ | ✅ | ❌ |
| `ledger:write` | ❌ | ✅ | ❌ | ❌ | ❌ |
| `zakat:manage` | ✅ | ✅ | ✅ | ❌ | ❌ |
| `jamaah:manage` | ✅ | ✅ | ✅ | ✅ | ❌ |
| `inventory:manage`| ❌ | ❌ | ✅ | ✅ | ❌ |
| `qurban:manage` | ❌ | ❌ | ✅ | ✅ | ❌ |
| `agenda:manage` | ✅ | ❌ | ✅ | ✅ | ❌ |
| `settings:manage` | ✅ | ❌ | ❌ | ✅ | ❌ |
