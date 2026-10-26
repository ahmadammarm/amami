# Amami

**amami** (*Aplikasi Manajemen Masjid*) is a management system for mosques. It gives the mosque committee (pengurus) one auditable place to keep the cash ledger, zakat, qurban, the community registry, assets and the event calendar, and gives the congregation (jamaah) a public page with the mosque's profile and donation details.

The interface is in Bahasa Indonesia.

## What it does

| Module | What the committee can do |
|---|---|
| **Keuangan** (finance) | Keep separate funds and an append-only ledger. Entries are never edited or deleted; mistakes are corrected with a linked reversal entry. |
| **Zakat** | Record zakat and infaq received from muzakki, and distributions to mustahik. Cash posts to a fund automatically; rice (kg) is recorded without touching the ledger. |
| **Jamaah** | Maintain the community registry, mark mustahik, and run a needs survey that produces a 0-100 poverty score used to prioritise zakat. |
| **Qurban** | Define packages and quotas, register shohibul (the price comes from the package), and track animals from arrival to slaughter. |
| **Inventaris** | Register mosque assets with their condition; lend them to jamaah and record returns. |
| **Agenda** | Publish the weekly schedule of studies and events, and archive photo/PDF documentation per event. |
| **Pengaturan / Pengguna** | Edit the mosque profile and bank account shown publicly; invite committee members and assign roles. |

Access is role-based and enforced on both the API and the UI:

| Role | Typical access |
|---|---|
| `SUPER_ADMIN` | Users, settings, and read access across modules |
| `BENDAHARA` (treasurer) | Ledger writes, zakat, jamaah |
| `TAKMIR` | Zakat, jamaah and mustahik surveys, qurban, inventory, agenda |
| `SEKRETARIS` (secretary) | Jamaah and surveys, qurban, inventory, agenda, settings |
| `JAMAAH` | Read-only inventory and agenda |

The exact grants are in [`docs/RBAC_VIEW_MATRIX.md`](docs/RBAC_VIEW_MATRIX.md).

## Repository layout

```
backend/    Go REST API (Gin, GORM, PostgreSQL)
frontend/   Vue 3 + TypeScript single-page app (Vite, Tailwind CSS 4)
docs/       Feature specifications, RBAC matrix, database notes, roadmap
```

See [`backend/README.md`](backend/README.md) and [`frontend/README.md`](frontend/README.md) for details.

## Quick start

**Requirements:** Go 1.25+, PostgreSQL, Node.js 24 (see `frontend/.nvmrc`), pnpm.

```bash
# 1. API
cd backend
cp .env.example .env        # then replace the example values with your own
go run cmd/main.go -seed    # once, on an empty database: admin account + sample data
go run cmd/main.go          # serves http://localhost:8080

# 2. Web app (second terminal)
cd frontend
pnpm install
pnpm dev                    # serves http://localhost:5173
```

Sign in at <http://localhost:5173/login> with the seeded administrator:

| Email | Password |
|---|---|
| `admin@amami.org` | `admin123` |

Change this password before using the system with real data. Other committee accounts are created from **Pengguna → Undang pengurus**; new users must set their own password at first sign-in.

## Documentation

| Document | Contents |
|---|---|
| [`docs/FEATURE_SPECIFICATIONS.md`](docs/FEATURE_SPECIFICATIONS.md) | What each module does, and which features are available, API-only or planned |
| [`docs/RBAC_VIEW_MATRIX.md`](docs/RBAC_VIEW_MATRIX.md) | Roles, page visibility and permission grants |
| [`docs/DATABASE_DOCUMENTATION.md`](docs/DATABASE_DOCUMENTATION.md) | Tables, relationships and integrity rules |
| [`CLAUDE.md`](CLAUDE.md) | Architecture and conventions for contributors |

## License

Proprietary. All rights reserved.
