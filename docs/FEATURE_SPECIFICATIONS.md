# Amami — Feature Specifications

Functional specification of amami, the mosque management system. Each feature is marked with its status:

| Status | Meaning |
|---|---|
| **Available** | Implemented in the API and the web app |
| **API only** | Implemented in the API; no screen in the web app yet |
| **Planned** | Specified here, not implemented |

Who may use each feature is defined in [`RBAC_VIEW_MATRIX.md`](RBAC_VIEW_MATRIX.md); the tables behind it in [`DATABASE_DOCUMENTATION.md`](DATABASE_DOCUMENTATION.md).

---

## 1. Public page and access

### 1.1 Public mosque page — Available
The landing page shows the mosque's name, address, phone, bank account for zakat and infaq (with a copy button), the zakat fitrah amount per person, and a summary of what the committee manages. All of it comes from the mosque profile (§2.3); nothing is placeholder content. The page is open to everyone; signed-in users see a button to the dashboard instead of the sign-in button.

### 1.2 Sign-in and first password — Available
- Committee members sign in with email and password. Passwords are hashed with bcrypt.
- Accounts created by an administrator start as `PENDING_PASSWORD_CHANGE`. Until the user sets a new password, their session can only reach the change-password screen.
- Suspended accounts cannot sign in.
- Sign-in is rate-limited to 5 attempts per minute.

### 1.3 Forgot password — Planned
Request a reset by email, receive a single-use link valid for 15 minutes, set a new password. Requires working SMTP (§2.4).

### 1.4 Account settings — Planned
Users update their own display name and phone number, and change their password after entering the current one.

---

## 2. Administration

### 2.1 Roles and permissions — Available
- Five roles: `SUPER_ADMIN`, `BENDAHARA`, `TAKMIR`, `SEKRETARIS`, `JAMAAH`.
- Each role holds a fixed set of permission codes (for example `ledger:read`, `ledger:write`, `mustahik:survey`). Reading and managing are separate codes.
- The API checks the permission on every request; the web app hides pages and actions the user's role cannot use, based on the codes returned at sign-in.
- Grants are defined in code and synced to the database on every server start. Changing them is a code change, not an admin screen.

### 2.2 User management — Available
- The super administrator invites committee members with name, username, email, role and an initial password. The password is shown once, to be passed on to the new member.
- Users can be suspended and reactivated.
- Users can be deleted (soft delete: the account can no longer sign in, history is kept).
- Nobody can delete their own account or a super administrator.

### 2.3 Mosque profile — Available
Name, address, phone, legal foundation number, bank name, account number and holder, zakat fitrah amount per person, and the active Hijri year. Shown on the public page. Closing the edit form with unsaved changes asks for confirmation.

### 2.4 Email (SMTP) settings — API only
Host, port, username, password and sender address, with a connection test. The password is never returned by the API but is stored unencrypted. No screen yet; nothing in the app sends email yet.

### 2.5 System health and audit log — API only
- A health check reports database and SMTP reachability.
- An audit log records sign-ins and password changes.
- Both have API endpoints but no screens yet.

### 2.6 Internal memos — Planned
Committee bulletins, with pinned memos shown on the dashboard. The `memos` table exists; there is no API or screen.

---

## 3. Finance

### 3.1 Funds and ledger — Available
- The mosque keeps separate funds (for example operational cash, building fund, orphan fund), each with a running balance.
- The treasurer records income (`CREDIT`) or expenses (`DEBIT`) against a fund, with a category and an optional description. Expenses larger than the fund balance are rejected.
- The ledger is append-only: entries cannot be edited or deleted.
- Each entry records who made it and when (set by the server).
- Every posting updates the fund balance in the same database transaction, so a balance always equals the sum of its entries, even under concurrent use.

### 3.2 Reversal entries — Available
- A mistaken entry is corrected by reversing it. This posts an opposite entry of the same amount to the same fund, with the reason as its description, linked to the original.
- A transaction can be reversed once, and a reversal cannot itself be reversed. The database enforces this, not only the app.
- The ledger shows reversed entries struck through and labels reversal entries.

### 3.3 Financial reports — Planned
Printable period reports per fund, and receipts for donors.

---

## 4. Zakat and social welfare

### 4.1 Jamaah registry — Available
- The community registry holds name, phone and address, and whether the person is a mustahik (eligible to receive zakat).
- Search by name is fuzzy (trigram).
- The list can be filtered to mustahik only and sorted by newest, highest score, or name.
- Committee accounts each have a linked jamaah profile.

### 4.2 Mustahik survey and score — Available
- A surveyor records monthly household income, number of dependents, housing (own, family, rent, none), whether someone in the household has a disability, and notes.
- Saving the survey recalculates a poverty score from 0 to 100:

| Factor | Points |
|---|---|
| Income per household member (income ÷ (dependents + 1)); full points at zero, none at Rp 1,200,000 or more | 50 |
| Dependents, capped at five | 20 |
| Housing: none 20, rent 15, family 10, own 0 | 20 |
| Disability in the household | 10 |

The score is shown in the registry and used to order mustahik when distributing zakat.

### 4.3 Zakat collection and distribution — Available
- **Collection:** record zakat fitrah (rice, in kg), zakat maal, or infaq from a muzakki. Cash must go to a fund and must be a whole rupiah amount; it is posted to the ledger in the same step. Rice is recorded without touching the ledger.
- **Distribution:** give cash (from a fund, posted as an expense) or rice to a mustahik. Only jamaah marked as mustahik can receive. The picker lists them by highest score first.
- Zakat records are never deleted.

### 4.4 Zakat calculator — Planned
Zakat maal (2.5% of assets above the nisab), zakat profesi from monthly income, with a configurable gold price for the nisab threshold.

### 4.5 Asnaf categories and NIK — Planned
Recording the asnaf category and national ID (NIK) of mustahik, with NIK masked in any shared view.

---

## 5. Qurban — Available

- **Packages:** name, type (goat, sheep, cow, 1/7 cow share), price, Hijri year and quota. The quota cannot be lowered below the number of bookings already made.
- **Bookings:** register a shohibul (jamaah) for a package with payment status pending or paid. The amount is taken from the package price, and the package quota is reduced in the same database transaction, so a package cannot be overbooked.
- **Animals:** tag number (unique), type, weight, vendor, and status (available, assigned, slaughtered).
- Packages and animals can be soft-deleted; existing bookings keep their history.

Planned: the intention name (*niat*) per booking, linking bookings to ledger payments, and assigning animals to bookings.

---

## 6. Inventory

### 6.1 Asset registry — Available
- Each asset has a name, SKU (unique), location, purchase date and price, and condition: good, on loan, needs repair, under repair, or broken.
- "On loan" is set and cleared only by the loan workflow, never by editing the asset.
- An asset on loan cannot be deleted.
- Deleting an asset is a soft delete; its loan history is kept.

### 6.2 Loans — API only
- Lend an asset in good condition to a jamaah, with a loan date and optional due date. Lending marks the asset as on loan.
- Returning records the date and condition notes, and marks the asset good again.
- A loan cannot be returned twice, and dates cannot run backwards.
- The web app has no loan screen yet.

Planned: dashboard alerts for loans past their due date.

---

## 7. Agenda — Available

- **Schedule:** recurring weekly activities with a title, description, day, time, location and status (scheduled, ongoing, completed, cancelled). They are listed in weekly order.
- **Documentation:**
  - Upload up to 10 photos or PDFs at a time per activity. File types are checked from the file content (JPEG, PNG, WebP, PDF); the size limit is configurable (10 MB by default).
  - Files can be viewed and deleted.
  - Deleting an activity keeps its documentation for the archive.

Planned: one-off dated events, speakers (*mubaligh*), attendance, and an activity report for the annual accountability report (LPJ).

---

## 8. Dashboard — Available

- Shows the Gregorian and Hijri date.
- Shows summary figures: fund balance and this month's income, zakat collected and number of mustahik, jamaah count and new members, scheduled activities, assets on loan, and assets needing repair.
- Charts:
  - six-month income and expenses
  - share of zakat distributed
  - qurban animals by type
  - assets by condition
  - new jamaah per month
- Each figure and chart appears only if the user's role can read that module.

---

## 9. Reliability and feedback — Available

- **Status pages:** not-found (404), no-access (403), and an error page for unexpected failures.
- **Outages:** when the server reports maintenance (503), cannot be reached, or the device is offline, the app shows a maintenance page.
  - The maintenance page retries every 15 seconds and when the connection returns, then goes back to the page the user was on.
  - Outages do not sign the user out.
- **Notifications:** results and errors appear as dialogs, and destructive actions ask for confirmation first.

Planned: database backup from the app.

---

## 10. Business rules

| Area | Rule | Why |
|---|---|---|
| Finance | Ledger entries are never edited or deleted; corrections are reversals, at most one per entry | Every rupiah stays traceable |
| Finance | Postings, fund balance updates and their records are one database transaction | Balances always match the ledger |
| Zakat | Only whole-rupiah IDR amounts are posted to funds; only mustahik receive distributions | Ledger stays in rupiah; aid reaches eligible recipients |
| Qurban | Booking reserves quota atomically; price comes from the package | No overbooking, no manipulated totals |
| Inventory | Only assets in good condition can be lent; "on loan" is managed by loans | Asset status stays truthful |
| Users | No self-deletion; super administrators cannot be deleted | Prevents locking the mosque out of its own system |
| General | People, assets, packages, animals and activities are soft-deleted | History stays available for reports |
