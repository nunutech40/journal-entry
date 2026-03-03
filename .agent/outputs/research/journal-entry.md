---
topic: Journal Entry (Accounting) Web App — Go + HTMX + Alpine.js
date: 2026-03-03
version: 3
status: completed
related_files: []
---

# Research: Journal Entry Web App

## Summary

Riset tentang standar fitur journal entry app (pembukuan/akuntansi), tech stack Go + HTMX + Alpine.js, database schema double-entry bookkeeping, arsitektur modular (feature-based), dan project structure yang optimal. App ini adalah web-based accounting tool untuk mencatat transaksi keuangan (debit/kredit, uang masuk/keluar).

---

## 1. Standar Fitur Journal Entry App

### Core Features (MVP)
1. **Double-Entry System** — Setiap transaksi punya minimal 2 entries: debit & kredit, harus balance
2. **Chart of Accounts (COA)** — Daftar akun terstruktur: Asset, Liability, Equity, Revenue, Expense
3. **Journal Entry CRUD** — Buat, lihat, edit (via correction entry), hapus (soft delete) transaksi
4. **General Ledger** — Buku besar: semua transaksi per akun
5. **Trial Balance** — Neraca saldo: ringkasan saldo semua akun
6. **Dashboard** — Ringkasan keuangan: total pemasukan, pengeluaran, saldo

### Extended Features (Phase 2+)
7. **Laporan Keuangan** — Neraca (Balance Sheet), Laba Rugi (Income Statement)
8. **Filter & Search** — By tanggal, akun, kategori, deskripsi
9. **Export** — PDF, CSV, Excel
10. **Kategori / Tags** — Untuk grouping transaksi
11. **Recurring Transactions** — Transaksi berulang (gaji, sewa, dll)
12. **Bank Reconciliation** — Cocokkan data dengan rekening bank
13. **Multi-user** — Login, roles
14. **Attachment** — Upload bukti transaksi (foto nota, invoice)
15. **Financial Period** — Periode akuntansi (bulanan, tahunan)

### Accounting Rules yang Harus Dipatuhi
- Total debit HARUS = total kredit untuk setiap journal entry
- Transaksi bersifat **immutable** — koreksi via reversing entry, bukan edit langsung
- Gunakan tipe data **NUMERIC** untuk amount, JANGAN float
- Audit trail: siapa buat, kapan, perubahan apa

---

## 2. Chart of Accounts (COA) Standard Indonesia

Berdasarkan PSAK EMKM (standar untuk UMKM Indonesia):

```
1xxx — ASET (Assets)
  1100 Kas (Cash)
  1200 Bank
  1300 Piutang Usaha (Accounts Receivable)
  1400 Persediaan (Inventory)
  1500 Beban Dibayar di Muka (Prepaid Expenses)
  1600 Aset Tetap (Fixed Assets)
  1700 Akumulasi Penyusutan (Accumulated Depreciation)

2xxx — KEWAJIBAN (Liabilities)
  2100 Utang Usaha (Accounts Payable)
  2200 Utang Gaji (Salaries Payable)
  2300 Utang Pajak (Taxes Payable)
  2400 Utang Bank (Bank Loans)

3xxx — EKUITAS (Equity)
  3100 Modal Pemilik (Owner's Capital)
  3200 Laba Ditahan (Retained Earnings)
  3300 Prive (Drawings)

4xxx — PENDAPATAN (Revenue)
  4100 Pendapatan Penjualan (Sales Revenue)
  4200 Pendapatan Jasa (Service Revenue)
  4900 Pendapatan Lain-lain (Other Income)

5xxx — HARGA POKOK (Cost of Goods Sold)
  5100 HPP (COGS)

6xxx — BEBAN (Expenses)
  6100 Beban Gaji (Salary Expense)
  6200 Beban Sewa (Rent Expense)
  6300 Beban Listrik/Air/Telepon (Utilities)
  6400 Beban Pemasaran (Marketing Expense)
  6500 Beban Perlengkapan (Supplies Expense)
  6600 Beban Penyusutan (Depreciation Expense)
  6700 Beban Transportasi (Transportation Expense)
  6800 Beban Administrasi (Admin Expense)
  6900 Beban Lain-lain (Other Expense)
```

---

## 3. Tech Stack

| Layer | Technology | Version | Alasan |
|---|---|---|---|
| **Backend** | Go (Golang) | 1.22+ | Performant, strongly typed, native HTTP server |
| **Router** | chi | v5.2.5 | Route grouping, middleware, subrouter support |
| **Server Interaction** | HTMX | 2.0.7 | Server-driven interactivity, partial page updates |
| **Client Interactivity** | Alpine.js | 3.15.8 | Lightweight client-side state (modals, dropdowns) |
| **Template Engine** | Go `html/template` | stdlib | Built-in, no dependency |
| **Database** | PostgreSQL | 15+ | Scalable, concurrent writes, production-ready |
| **DB Driver** | pgx | v5 | Modern, performant PostgreSQL driver for Go |
| **Migration** | goose | latest | Simple, Go-native, auto up/down migrations |
| **CSS** | Vanilla CSS | - | Flexible, no build step |
| **Hot Reload** | Air | latest | Live reload saat development |
| **Env Config** | godotenv | latest | Load `.env` file for local development |

### CDN untuk Frontend (no npm/build step)
- HTMX: `https://unpkg.com/htmx.org@2.0.7`
- Alpine.js: `https://unpkg.com/alpinejs@3.15.8`
- Alpine Morph: `https://unpkg.com/@alpinejs/morph@3.15.8`

---

## 4. ERD (Entity Relationship Diagram)

```
┌─────────────────────┐
│      accounts       │
├─────────────────────┤
│ id          SERIAL  │──┐ (self-ref: parent)
│ code        VARCHAR │  │
│ name        VARCHAR │  │
│ type        VARCHAR │  │
│ parent_id   INT FK  │──┘
│ description TEXT    │
│ is_active   BOOL   │
│ created_at  TSTZ   │
│ updated_at  TSTZ   │
└─────────┬───────────┘
          │ 1
          │
          │ N
┌─────────┴───────────┐         ┌─────────────────────┐
│    entry_lines      │         │   journal_entries    │
├─────────────────────┤         ├─────────────────────┤
│ id          SERIAL  │         │ id          SERIAL  │
│ journal_entry_id FK │────N:1──│ entry_date  DATE    │
│ account_id     FK   │         │ description TEXT    │
│ description TEXT    │         │ reference   VARCHAR │
│ debit     NUMERIC   │         │ is_posted   BOOL   │
│ credit    NUMERIC   │         │ created_at  TSTZ   │
│ created_at  TSTZ   │         │ updated_at  TSTZ   │
└─────────────────────┘         └─────────────────────┘

Relationships:
- journal_entries 1 ──→ N entry_lines  (1 jurnal punya banyak baris debit/kredit)
- accounts 1 ──→ N entry_lines          (1 akun bisa muncul di banyak transaksi)
- accounts 1 ──→ N accounts             (parent-child untuk sub-akun)

Constraints (enforce di app layer):
- Per entry_line: debit > 0 XOR credit > 0 (nggak boleh dua-duanya)
- Per journal_entry: SUM(debit) == SUM(credit) dari semua entry_lines-nya
```

---

## 5. Database Schema (PostgreSQL)

```sql
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL,
    parent_id INTEGER REFERENCES accounts(id),
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE journal_entries (
    id SERIAL PRIMARY KEY,
    entry_date DATE NOT NULL,
    description TEXT NOT NULL,
    reference_number VARCHAR(50),
    is_posted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE entry_lines (
    id SERIAL PRIMARY KEY,
    journal_entry_id INTEGER NOT NULL REFERENCES journal_entries(id) ON DELETE CASCADE,
    account_id INTEGER NOT NULL REFERENCES accounts(id),
    description TEXT,
    debit NUMERIC(15,2) DEFAULT 0,
    credit NUMERIC(15,2) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_entry_lines_journal ON entry_lines(journal_entry_id);
CREATE INDEX idx_entry_lines_account ON entry_lines(account_id);
CREATE INDEX idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX idx_accounts_code ON accounts(code);
CREATE INDEX idx_accounts_type ON accounts(type);
```

---

## 6. Go Backend Architecture (Feature-Based Modular)

### Prinsip Modularitas
- **Feature-based**: Setiap fitur punya folder sendiri (handler + service + repo + model)
- **Ubah 1 fitur → file lain nggak kesenggol**, kecuali di `shared/`
- **Depend on interfaces**: service terima interface, bukan struct konkret
- **Shared code = stable & additive only**: tambah fungsi baru boleh, ubah signature lama → cek semua caller dulu

### Dependency Flow (siapa depend siapa)

```
main.go
  │
  ├─ creates: pgxpool.Pool (DB connection)
  ├─ creates: shared/* (helpers)
  │
  ├─ creates: account.Repository  (terima *pgxpool.Pool)
  ├─ creates: account.Service     (terima account.Repository interface)
  ├─ creates: account.Handler     (terima account.Service interface)
  │
  ├─ creates: journal.Repository  (terima *pgxpool.Pool)
  ├─ creates: journal.Service     (terima journal.Repository + account.Repository interfaces)
  ├─ creates: journal.Handler     (terima journal.Service interface)
  │
  ├─ creates: dashboard.Handler   (terima account.Service + journal.Service interfaces)
  ├─ creates: report.Handler      (terima report.Service interface)
  │
  └─ registers: routes (chi.Router)
```

### Interface Pattern

```go
// internal/account/repository.go
type Repository interface {
    GetAll(ctx context.Context) ([]Model, error)
    GetByID(ctx context.Context, id int) (Model, error)
    Create(ctx context.Context, a Model) (Model, error)
    Update(ctx context.Context, a Model) error
    Delete(ctx context.Context, id int) error
}

// internal/account/service.go
type Service interface {
    ListAccounts(ctx context.Context) ([]Model, error)
    GetAccount(ctx context.Context, id int) (Model, error)
    CreateAccount(ctx context.Context, req CreateRequest) (Model, error)
    // ...
}

// Concrete implementations:
type repositoryImpl struct { pool *pgxpool.Pool }
type serviceImpl struct { repo Repository }  // ← depends on interface, not concrete
```

### Handler Pattern (HTMX response)

```go
// internal/account/handler.go
type Handler struct {
    svc    Service
    tmpl   *template.Template
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
    // 1. Parse form
    // 2. Validate (return error partial kalau gagal)
    // 3. Call service
    // 4. Render partial template (HTMX swap)
}
```

### Error Handling ke HTMX

```go
// Sukses: render partial HTML, HTMX swap ke target
// Error:  render error partial, HTMX swap ke error container
//
// Response headers yang penting:
// HX-Trigger: showToast  → trigger Alpine.js toast notification
// HX-Retarget: #errors   → redirect swap ke error container
// HX-Reswap: innerHTML   → override swap strategy
```

---

## 7. Frontend Architecture (Templates + HTMX + Alpine.js)

### Template Composition

```
base.html (layout)
  ├─ <head>: meta, CSS, HTMX, Alpine.js
  ├─ <nav>: navigation bar
  ├─ <aside>: sidebar (menu)
  ├─ <main>: {{ block "content" . }}  ← diisi oleh page template
  └─ <footer>

page template (e.g. account/list.html)
  ├─ {{ define "content" }}
  ├─ Heading, filters, actions
  ├─ Table/list → rows loaded via HTMX
  └─ Modal forms (Alpine.js x-data)

partial template (e.g. account/_row.html)
  ├─ Single <tr> atau component
  ├─ Di-swap oleh HTMX setelah CRUD
  └─ Nggak include layout/head/nav
```

### HTMX Patterns yang Dipakai

| Pattern | Contoh | Notes |
|---|---|---|
| **List + inline edit** | `hx-get="/accounts"` → render table rows | Full page load pertama, partial setelahnya |
| **Form submit** | `hx-post="/accounts" hx-target="#account-list"` | Server return updated list partial |
| **Delete with confirm** | Alpine.js modal → `hx-delete="/accounts/5"` | Confirm dulu, baru HTMX request |
| **Dynamic form rows** | `hx-get="/journal/add-line" hx-target="#lines" hx-swap="beforeend"` | Tambah baris debit/kredit |
| **Toast notification** | Server set `HX-Trigger: showToast` header | Alpine.js listen event, show toast |
| **Error display** | Server return `422` + error HTML, `HX-Retarget="#errors"` | Swap error message ke container |

### Alpine.js Patterns yang Dipakai

| Pattern | Contoh | Notes |
|---|---|---|
| **Modal** | `x-data="{ open: false }"` | Toggle confirmation dialog |
| **Form state** | `x-data="{ lines: [...] }"` | Track debit/credit lines client-side |
| **Live calculation** | `x-text="totalDebit - totalCredit"` | Show balance diff realtime |
| **Dropdown** | `x-data="{ open: false }" @click.away="open = false"` | Account type filter |

### ⚠️ HTMX + Alpine.js DOM Issue
Ketika HTMX swap DOM, Alpine.js state di elemen yang di-swap akan hilang.
**Solusi:** Include Alpine Morph plugin + gunakan `hx-ext="morph"` pada elemen yang punya Alpine state.

### CSS Architecture

```css
/* === Naming Convention: BEM-inspired === */

/* Shared components (prefix: c-) */
.c-btn { }
.c-btn--primary { }
.c-btn--danger { }
.c-modal { }
.c-toast { }
.c-table { }
.c-form-field { }

/* Feature-specific (prefix: feature name) */
.account-list { }
.account-form { }
.journal-form { }
.journal-lines { }
.dashboard-summary { }

/* Layout */
.l-sidebar { }
.l-content { }
.l-header { }
```

**Rule:** Ubah `.journal-form` → nggak affect `.account-form`. Ubah `.c-btn` → cek semua pemakai dulu.

---

## 8. Project Structure (Feature-Based)

```
journal-entry/
├── .agent/                          ← Workflow & research outputs
├── cmd/
│   └── web/
│       └── main.go                  ← Entry point: init DB, init services, register routes, start server
│
├── internal/
│   ├── account/                     ← 🟢 Feature: Chart of Accounts
│   │   ├── handler.go               ← HTTP handlers (CRUD endpoints)
│   │   ├── service.go               ← Business logic + validation
│   │   ├── repository.go            ← PostgreSQL queries
│   │   └── model.go                 ← Account struct + request/response types
│   │
│   ├── journal/                     ← 🟢 Feature: Journal Entry
│   │   ├── handler.go               ← HTTP handlers (CRUD + entry lines)
│   │   ├── service.go               ← Business logic (debit=credit validation)
│   │   ├── repository.go            ← PostgreSQL queries (transactions)
│   │   └── model.go                 ← JournalEntry + EntryLine structs
│   │
│   ├── report/                      ← 🟢 Feature: Reports (Ledger, Trial Balance)
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   │
│   ├── dashboard/                   ← 🟢 Feature: Dashboard
│   │   └── handler.go               ← Aggregates data from account + journal services
│   │
│   ├── shared/                      ← 🔒 Shared utilities (STABLE, additive only)
│   │   ├── response/
│   │   │   └── htmx.go              ← HTMX response helpers (render partial, set headers)
│   │   ├── validate/
│   │   │   └── validate.go          ← Common validation functions
│   │   ├── money/
│   │   │   └── money.go             ← Decimal/money formatting helpers
│   │   └── middleware/
│   │       ├── logging.go           ← Request logging
│   │       └── recovery.go          ← Panic recovery
│   │
│   └── server/                      ← 🔧 Server setup
│       ├── routes.go                ← All route registration (chi router)
│       └── server.go                ← HTTP server config, graceful shutdown
│
├── templates/                       ← Go HTML templates
│   ├── layout/
│   │   └── base.html                ← Base layout (head, nav, sidebar, footer)
│   ├── components/                  ← 🔒 Reusable UI components
│   │   ├── _toast.html              ← Toast notification
│   │   ├── _modal.html              ← Confirmation modal
│   │   ├── _form_field.html         ← Reusable form field wrapper
│   │   ├── _pagination.html         ← Pagination links
│   │   └── _empty_state.html        ← "No data" placeholder
│   ├── account/                     ← Feature-specific templates
│   │   ├── list.html                ← Full page: account list
│   │   ├── form.html                ← Full page: create/edit account
│   │   └── _row.html                ← Partial: single account row (HTMX swap)
│   ├── journal/
│   │   ├── list.html                ← Full page: journal entry list
│   │   ├── form.html                ← Full page: create/edit journal entry
│   │   ├── _row.html                ← Partial: single journal row
│   │   └── _entry_line.html         ← Partial: single debit/credit line
│   ├── report/
│   │   ├── ledger.html              ← Full page: general ledger
│   │   └── trial_balance.html       ← Full page: trial balance
│   └── dashboard/
│       └── index.html               ← Full page: dashboard
│
├── static/
│   ├── css/
│   │   └── style.css                ← Single CSS file (BEM naming)
│   ├── js/
│   │   └── app.js                   ← Minimal JS (Alpine.js init, toast handler)
│   └── img/
│
├── db/
│   ├── migrations/                  ← goose migrations
│   │   ├── 001_create_accounts.sql
│   │   ├── 002_create_journal_entries.sql
│   │   └── 003_create_entry_lines.sql
│   └── seeds/
│       └── chart_of_accounts.sql    ← Default COA Indonesia
│
├── .env.example                     ← Template environment variables
├── .air.toml                        ← Hot reload config
├── go.mod
├── go.sum
└── README.md
```

---

## 9. Shared Code Strategy (⚡ Penting)

### Prinsip: "Ubah 1 fitur, fitur lain nggak kesenggol"

#### Kapan kode masuk `shared/`?
| Kriteria | Masuk shared? | Contoh |
|---|---|---|
| Dipake 2+ fitur | ✅ Ya | `response/htmx.go`, `money/money.go` |
| Logic spesifik 1 fitur | ❌ Nggak | Account validation rules |
| Bisa berubah per fitur | ❌ Nggak | Bikin di masing-masing fitur |

#### Rules untuk shared code:

```
1. ADDITIVE ONLY
   ✅ Tambah fungsi baru           → aman, nggak affect siapapun
   ❌ Ubah signature fungsi lama   → HARUS cek semua caller dulu!

2. JANGAN bikin "God helper"
   ✅ shared/money/money.go        → focused: format, parse, calculate
   ✅ shared/validate/validate.go  → focused: email, required, numeric
   ❌ shared/utils/utils.go        → dump segala macam = anti-pattern

3. KALAU shared function butuh behavior beda per fitur:
   → Opsi A: Tambah parameter/option (preferred)
     func FormatMoney(amount decimal, opts ...FormatOption) string
   → Opsi B: Return interface, biar caller customize
   → Opsi C (last resort): Copy ke feature package, jadi independent
     Duplication > wrong abstraction

4. SEBELUM ubah shared code:
   → grep_search semua pemakai
   → pastikan perubahan backward-compatible
   → kalau nggak compatible, bikin fungsi BARU (v2), jangan ubah yang lama
```

#### Contoh konkret:

```go
// shared/response/htmx.go — dipakai oleh SEMUA handlers

// ✅ Fungsi yang STABLE (jangan ubah signature)
func RenderPartial(w http.ResponseWriter, tmpl *template.Template, name string, data any) error
func RenderError(w http.ResponseWriter, tmpl *template.Template, message string) error
func SetTrigger(w http.ResponseWriter, event string)

// ✅ Boleh TAMBAH fungsi baru
func RenderPartialWithStatus(w http.ResponseWriter, status int, tmpl *template.Template, name string, data any) error

// ❌ JANGAN ubah RenderPartial signature yang lama!
```

### Template Components: Sama prinsipnya

```
Ubah templates/components/_toast.html
  → Cek: siapa yang pake? Semua page.
  → Aman kalau: tambahin parameter optional (e.g. {{if .Icon}})
  → Berbahaya kalau: ubah structure HTML yang existing (break CSS/JS)

Ubah templates/journal/_row.html
  → Aman: cuma dipake di journal feature
  → Nggak affect account, report, dashboard
```

---

## 10. Data Flow

```
User Action (Browser)
  │
  ├─ HTMX Request (hx-get, hx-post, hx-put, hx-delete)
  │   │
  │   └─→ chi Router (internal/server/routes.go)
  │         │
  │         ├─→ Middleware (logging, recovery)
  │         │
  │         └─→ Feature Handler (e.g. internal/journal/handler.go)
  │               │
  │               ├─→ Feature Service (internal/journal/service.go)
  │               │     │
  │               │     ├─→ Validation (debit == credit?)
  │               │     │
  │               │     └─→ Feature Repository (internal/journal/repository.go)
  │               │           │
  │               │           └─→ PostgreSQL (via pgxpool)
  │               │
  │               └─→ Render Template
  │                     │
  │                     ├─ Full page: templates/journal/list.html (extends base.html)
  │                     └─ Partial: templates/journal/_row.html (HTMX swap)
  │
  └─ Alpine.js (client-side only, no server call)
      ├─ Toggle modal/dropdown
      ├─ Live debit/credit balance calculation
      └─ Client-side form validation
```

---

## 11. Keputusan & Alasan (Decisions)

| Keputusan | Alasan |
|---|---|
| **Feature-based structure** (bukan layer-based) | Ubah 1 fitur nggak kesenggol fitur lain. Semua file terkait 1 fitur ada dalam 1 folder |
| **PostgreSQL** (bukan SQLite) | Future-proof, concurrent writes, production-ready |
| **pgx** (bukan lib/pq) | Modern driver, better performance, native PostgreSQL types |
| **chi** (bukan stdlib mux) | Route grouping per feature, middleware chaining, subrouter |
| **goose** (untuk migration) | Simple, Go-native, auto up/down, version tracking |
| **godotenv** (untuk env config) | Simple `.env` file, nggak perlu tool external |
| **Interface-based DI** | Testable, mockable, service nggak depend pada concrete repo |
| **Shared = additive only** | Tambah fungsi baru aman, ubah yang lama harus cek caller |
| **Template prefix `_` = partial** | Jelas mana yang full page vs HTMX partial |
| **BEM CSS naming** | Feature-specific styles nggak saling tabrakan |
| **Go `html/template`** (bukan Templ) | Built-in, no extra dependency, cukup untuk project ini |
| **Separate debit/credit columns** | Sesuai standar akuntansi, mudah di-validate |
| **Server-side rendering** (bukan SPA) | HTMX philosophy: server renders HTML, client swap fragment |
| **No auth di MVP** | Fokus ke fitur akuntansi dulu |
| **Immutable transactions** | Best practice akuntansi. MVP: allow edit draft entries |

---

## 12. Risks & Edge Cases

| Risk | Mitigation |
|---|---|
| Float precision untuk uang | Gunakan NUMERIC(15,2), bukan float |
| Debit ≠ Credit saat submit | Validate di service layer + UI feedback via HTMX |
| HTMX swap hilangkan Alpine state | Alpine Morph plugin + `hx-ext="morph"` |
| Browser back button behavior | HTMX `hx-push-url` untuk URL management |
| PostgreSQL down | Connection pooling (pgxpool) + retry + error handling |
| Shared code change breaks features | Additive only rule + grep all callers before change |
| Template component change | Prefix partial with `_`, feature templates isolated |
| CSS collision antar fitur | BEM naming: `.account-*` vs `.journal-*` |

---

## 13. Go Libraries

| Library | Purpose |
|---|---|
| `net/http` (stdlib) | HTTP server |
| `html/template` (stdlib) | Template rendering |
| `github.com/go-chi/chi/v5` | Router + middleware |
| `github.com/jackc/pgx/v5` | PostgreSQL driver |
| `github.com/jackc/pgx/v5/pgxpool` | Connection pooling |
| `github.com/pressly/goose/v3` | Database migrations |
| `github.com/joho/godotenv` | Load `.env` file |
| `github.com/cosmtrek/air` (dev) | Hot reload |

---

## 14. Kapan Pakai HTMX vs Alpine.js

| Kebutuhan | Pakai |
|---|---|
| CRUD operations | **HTMX** — server renders HTML fragment |
| Load data dari server | **HTMX** — `hx-get`, `hx-post` |
| Form submission + server validation | **HTMX** — server returns error/success HTML |
| Toggle dropdown/modal | **Alpine.js** — `x-show`, `x-data` |
| Client-side form validation (instant) | **Alpine.js** — `x-model`, computed |
| Tab switching (no server call) | **Alpine.js** — `x-show` + `x-on:click` |
| Live balance calculation | **Alpine.js** — computed property |
| Delete confirmation dialog | **Alpine.js** — modal state |
| Toast notification | **Both** — HTMX trigger event → Alpine.js show toast |

---

## 15. Notes

- **Bahasa UI:** Indonesia (sesuai target user). Bisa bilingual nanti
- **Responsive:** Desktop-first, tapi harus usable di mobile
- **MVP Scope:** COA management + Journal Entry CRUD + General Ledger + Trial Balance + Dashboard
- **Phase 2:** Reports (Neraca, Laba Rugi), Export, Filter, Period closing
- **Phase 3:** Multi-user, Auth, Attachment, Recurring transactions

---

## 📋 Changelog

| Versi | Tanggal    | Perubahan |
|-------|------------|-----------|
| v3    | 2026-03-03 | Tambah: Go architecture (feature-based modular), frontend architecture, ERD, shared code strategy, updated versions, migration tool (goose), env config (godotenv) |
| v2    | 2026-03-03 | Database: SQLite → PostgreSQL |
| v1    | 2026-03-03 | Initial research |
