---
topic: Journal Entry Web App — Go + HTMX + Alpine.js
date: 2026-03-03
version: 5
status: in-progress
research: ../research/journal-entry.md
phases_total: 7
phases_completed: 4
---

# Plan: Journal Entry Web App

## Summary

Implementasi web app journal entry (pembukuan/akuntansi) menggunakan Go + HTMX + Alpine.js + PostgreSQL. App ini untuk mencatat transaksi keuangan dengan double-entry bookkeeping. Dikerjakan dalam 7 phase, bottom-up: scaffold → database → fitur per fitur → polish.

## Research Reference

Baca: `.agent/outputs/research/journal-entry.md` (v3, 15 sections, semua keputusan arsitektur sudah final)

## Desired End State

Web app yang bisa:
- Manage Chart of Accounts (COA) standar Indonesia
- CRUD journal entries dengan validasi double-entry (debit == kredit)
- Lihat General Ledger per akun
- Lihat Trial Balance (neraca saldo)
- Dashboard ringkasan keuangan
- Server-rendered HTML via Go templates, interaktif via HTMX + Alpine.js
- Data tersimpan di PostgreSQL

## What We're NOT Doing (MVP)

- ❌ Auth / multi-user (Phase 2+ nanti)
- ❌ Laporan Neraca & Laba Rugi (Phase 2+ nanti)
- ❌ Export PDF/CSV/Excel
- ❌ Recurring transactions
- ❌ Bank reconciliation
- ❌ File attachment
- ❌ Financial period closing
- ❌ Deployment / Docker setup

---

## Phase 1: Project Scaffold & Base Layout

**Status:** ✅ Completed

**Goal:** Server jalan, halaman kosong tampil di browser dengan layout lengkap (nav, sidebar, content).

**Files:**
- [x] `go.mod` — Init Go module
- [x] `cmd/web/main.go` — Entry point: load .env, init DB pool, init server, start
- [x] `internal/server/server.go` — HTTP server config, graceful shutdown
- [x] `internal/server/routes.go` — Chi router setup, static file serving, middleware
- [x] `internal/shared/middleware/logging.go` — Request logging
- [x] `internal/shared/middleware/recovery.go` — Panic recovery
- [x] `internal/shared/response/htmx.go` — RenderPage, RenderPartial, RenderError, SetTrigger
- [x] `templates/layout/base.html` — Base layout: head (CDN links), nav, sidebar, {{ block "content" }}, footer
- [x] `templates/dashboard/index.html` — Placeholder dashboard (extends base)
- [x] `templates/components/_toast.html` — Toast notification component
- [x] `static/css/style.css` — CSS: variables, reset, layout, shared components (c-btn, c-table, c-form-field, c-modal, c-toast)
- [x] `static/js/app.js` — Alpine.js init, toast event listener, HTMX config
- [x] `.env.example` — DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME, APP_PORT
- [x] `.air.toml` — Hot reload config (watch .go + .html + .css + .js)
- [x] `.gitignore` — Go binary, .env, air tmp/

**Steps:**
1. `go mod init` + install deps (chi, pgx, pgxpool, godotenv)
2. Buat seluruh folder structure sesuai research section 8
3. Buat `.env.example`, `.air.toml`, `.gitignore`
4. Buat `main.go`: load env → init pgxpool → init server → listen
5. Buat `server.go`: graceful shutdown dengan context
6. Buat `routes.go`: chi router + static files + middleware stack
7. Buat middleware: logging (method, path, duration) + recovery (catch panic)
8. Buat `response/htmx.go`: helper functions untuk render template + HTMX headers
9. Buat `base.html`: full layout dengan HTMX + Alpine.js CDN, nav, sidebar, content block
10. Buat `style.css`: modern dark/light design system, BEM naming
11. Buat `app.js`: Alpine init + toast handler
12. Buat placeholder dashboard page
13. Verify: `go build` + manual test

**Success Criteria:**

#### Automated Verification:
- [x] `go build ./cmd/web/` — berhasil tanpa error
- [ ] Server start: `go run cmd/web/main.go` — tidak crash (dengan PostgreSQL running)

#### Manual Verification:
- [ ] `http://localhost:8080` → base layout tampil (nav, sidebar, content area)
- [ ] CSS loaded (styling terlihat, bukan plain HTML)
- [ ] Console: HTMX & Alpine.js loaded (no 404, no JS errors)
- [ ] Hot reload: edit file → server auto-restart (air)

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 2: Database Migrations & Seed Data

**Status:** ✅ Completed

**Goal:** Database schema ready (3 tabel + indexes), COA standar Indonesia ter-seed.

**Files:**
- [x] `db/migrations/001_create_accounts.sql` — UP: create accounts table + indexes + CHECK constraint; DOWN: drop
- [x] `db/migrations/002_create_journal_entries.sql` — UP: create journal_entries + indexes; DOWN: drop
- [x] `db/migrations/003_create_entry_lines.sql` — UP: create entry_lines + indexes + debit XOR credit constraint; DOWN: drop
- [x] `db/seeds/chart_of_accounts.sql` — Default COA Indonesia (PSAK EMKM, 27 akun, ON CONFLICT DO NOTHING)

**Steps:**
1. Install goose CLI
2. Buat 3 migration files sesuai schema di research section 5
3. Buat seed file: COA standar Indonesia dari research section 2
4. Jalankan `goose up` → verify tabel terbuat
5. Jalankan seed → verify data masuk
6. Test `goose down` → verify reversible

**Success Criteria:**

#### Automated Verification:
- [ ] `goose -dir db/migrations postgres "$DSN" up` — berhasil
- [ ] `goose -dir db/migrations postgres "$DSN" down` — berhasil (reversible)
- [x] `go build ./cmd/web/` — tetap berhasil

#### Manual Verification:
- [ ] Database: 3 tabel ada (accounts, journal_entries, entry_lines)
- [ ] Accounts: ~25 akun COA ter-seed (cek via psql: `SELECT code, name, type FROM accounts`)
- [ ] Constraints: FK, UNIQUE, indexes terdaftar

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 3: Chart of Accounts (Full CRUD)

**Status:** ✅ Completed

**Goal:** User bisa lihat daftar akun, tambah, edit, hapus. Semua interaktif via HTMX.

**Files:**
- [x] `internal/account/model.go` — Account struct, CreateRequest, UpdateRequest, TypeLabel
- [x] `internal/account/repository.go` — Repository interface + pgx implementation (GetAll, GetByID, GetByCode, Create, Update, Delete)
- [x] `internal/account/service.go` — Service interface + implementation (validasi: code unique, name required, type valid, whitespace trimming)
- [x] `internal/account/handler.go` — HandleList, HandleCreateForm, HandleCreate, HandleEditForm, HandleUpdate, HandleDelete
- [x] `internal/account/service_test.go` — 11 unit tests (mock repo, no DB dependency)
- [x] `templates/account/list.html` — Full page: account table + badges + HTMX delete
- [x] `templates/account/form.html` — Create/edit form with validation error display
- [ ] `templates/account/_row.html` — Skipped (using full page reload via HX-Redirect instead of row swap)
- [ ] `templates/components/_modal.html` — Skipped (using hx-confirm instead)
- [ ] `templates/components/_empty_state.html` — Inline in templates
- [x] Update `internal/server/routes.go` — Register /accounts routes + account templates
- [x] Update `cmd/web/main.go` — Wire: account.Repo → account.Service → account.Handler
- [x] Update `templates/layout/base.html` — Already had sidebar links
- [x] Update `static/css/style.css` — Added .account-* styles, code style

**Steps:**
1. Buat model.go: Account struct (matches DB schema)
2. Buat repository.go: interface + pgx queries
3. Buat service.go: interface + business logic (validate code unique, name required)
4. Buat service_test.go: test validation logic (mock repo)
5. Buat handler.go: HTTP handlers, render templates
6. Buat templates: list, form, _row, components
7. Wire DI di main.go, register routes
8. Test: automated + manual CRUD

**Success Criteria:**

#### Automated Verification:
- [x] `go build ./cmd/web/` — berhasil
- [x] `go test ./internal/account/...` — 11 tests pass
- [x] Test: CreateAccount dengan code duplikat → error
- [x] Test: CreateAccount dengan name kosong → error
- [x] Test: CreateAccount valid → success

#### Manual Verification:
- [x] Buka /accounts → tabel COA tampil (27 akun seeded)
- [x] Klik "Tambah Akun" → form tampil
- [ ] Submit form → row baru muncul di tabel (HTMX, tanpa full reload)
- [ ] Edit akun → form pre-filled → update berhasil
- [ ] Hapus akun → modal konfirmasi → row hilang
- [ ] Toast notification muncul setelah setiap aksi
- [x] Empty state tampil kalau tabel kosong
- [x] Sidebar: "Daftar Akun" ter-highlight saat di halaman akun

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 4: Journal Entry (CRUD + Double-Entry Validation)

**Status:** ✅ Completed

**Goal:** User bisa buat, lihat, edit, hapus journal entry. Validasi debit == kredit. Dynamic lines via HTMX + Alpine.js.

> ⚠️ **Ini phase paling kompleks.** Bisa dipecah jadi 2 sub-phase kalau terlalu besar:
> - 4a: Model + Repository + Service (data layer)
> - 4b: Handler + Templates (presentation layer)

**Files:**
- [x] `internal/journal/model.go` — JournalEntry, EntryLine, CreateRequest, UpdateRequest, FormatDate
- [x] `internal/journal/repository.go` — pgx with DB transactions (Create/Update atomic), GetAll with SUM JOIN, GetByID with account JOIN
- [x] `internal/journal/service.go` — Validation: date, description, min 2 lines, debit XOR credit per line, account exists, SUM balanced (float tolerance)
- [x] `internal/journal/handler.go` — CRUD handlers + parseFormToRequest (array fields), friendlyError, account dropdown data
- [x] `internal/journal/service_test.go` — 10 unit tests (mock repo for both journal + account)
- [x] `templates/journal/list.html` — Journal list with date format, totals, HTMX delete
- [x] `templates/journal/form.html` — Alpine.js dynamic form: add/remove lines, live totals, balance check, disabled submit
- [ ] `templates/journal/_row.html` — Skipped (using full page reload via HX-Redirect)
- [ ] `templates/journal/_entry_line.html` — Skipped (using Alpine x-for instead of HTMX partial)
- [x] Update `internal/server/routes.go` — Register /journals routes + journal templates
- [x] Update `cmd/web/main.go` — Wire: journalRepo → journalSvc(+accountRepo) → journalHandler(+accountSvc)
- [x] Update `templates/layout/base.html` — Already had sidebar links
- [x] Update `static/css/style.css` — Added .journal-form-* styles
- [ ] Update `static/js/app.js` — Not needed, Alpine.js inline handles calculation

**Steps:**
1. Buat model.go: JournalEntry (header) + EntryLine (lines) + request types
2. Buat repository.go:
   - GetAll (with pagination info), GetByID (JOIN with lines + account names)
   - Create: pgx transaction → insert journal_entry → insert entry_lines
   - Update: transaction → delete old lines → insert new lines
   - Delete: soft delete atau cascade
3. Buat service.go:
   - Validate: SUM(debit) == SUM(credit)
   - Validate: setiap line → debit > 0 XOR credit > 0
   - Validate: minimal 2 lines
   - Validate: account_id exists (via account.Repository interface)
4. Buat service_test.go:
   - Test: debit != credit → error
   - Test: line with both debit and credit > 0 → error
   - Test: < 2 lines → error
   - Test: valid entry → success
5. Buat handler.go: CRUD handlers + HandleAddLine (return partial _entry_line.html)
6. Buat templates:
   - form.html: Alpine.js x-data track lines, computed totalDebit/totalCredit/diff
   - _entry_line.html: account dropdown + debit/credit inputs
   - list.html: table with date, description, total, status
7. Wire DI, register routes
8. Test automated + manual

**Success Criteria:**

#### Automated Verification:
- [x] `go build ./cmd/web/` — berhasil
- [x] `go test ./internal/journal/...` — 10 tests pass
- [x] `go test ./...` — all 21 tests pass (11 account + 10 journal)
- [x] Test: debit ≠ credit → error returned
- [x] Test: single line → error returned
- [x] Test: line with debit AND credit → error returned
- [x] Test: valid entry (2+ lines, balanced) → success

#### Manual Verification:
- [x] Buka /journals → list jurnal tampil (empty state)
- [x] Buat jurnal baru → form dengan dynamic lines tampil
- [x] Tambah baris: klik "Tambah Baris" → line baru muncul (Alpine.js, no reload)
- [x] Account dropdown: pilih akun → akun terpilih
- [x] Live balance: Alpine.js hitung total debit, total kredit, selisih secara realtime
- [ ] Submit dengan selisih ≠ 0 → error message (button disabled)
- [ ] Submit balanced → sukses, redirect ke list, toast tampil
- [ ] Edit jurnal → form pre-filled dengan existing lines
- [ ] Hapus jurnal → confirm dialog → jurnal hilang dari list

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 5: General Ledger & Trial Balance

**Status:** ⬜ Not started

**Goal:** User bisa lihat buku besar per akun dan neraca saldo seluruh akun.

**Files:**
- [ ] `internal/report/model.go` — LedgerEntry, TrialBalanceRow structs
- [ ] `internal/report/repository.go` — Aggregate queries (SUM, GROUP BY)
- [ ] `internal/report/service.go` — Filter, format, running balance calculation
- [ ] `internal/report/handler.go` — HandleLedger, HandleTrialBalance
- [ ] `templates/report/ledger.html` — Full page: buku besar (filter account + date range)
- [ ] `templates/report/trial_balance.html` — Full page: neraca saldo (semua akun)
- [ ] Update `internal/server/routes.go` — Register /reports routes
- [ ] Update `cmd/web/main.go` — Wire report DI
- [ ] Update `templates/layout/base.html` — Sidebar: links "Buku Besar", "Neraca Saldo"
- [ ] Update `static/css/style.css` — .report-* styles

**Steps:**
1. Buat model.go: LedgerEntry (date, ref, desc, debit, credit, balance), TrialBalanceRow (account code, name, type, debit_sum, credit_sum)
2. Buat repository.go:
   - GetLedger(accountID, dateFrom, dateTo): query entry_lines JOIN journal_entries, ordered by date
   - GetTrialBalance(): SUM(debit), SUM(credit) GROUP BY account_id, dari posted entries saja
3. Buat service.go: calculate running balance untuk ledger, format angka
4. Buat handler.go: render reports, handle filter params (HTMX partial update kalau filter berubah)
5. Buat templates: tabel dengan totals row di bawah
6. Wire DI, register routes

**Success Criteria:**

#### Automated Verification:
- [ ] `go build ./cmd/web/` — berhasil
- [ ] `go test ./internal/report/...` — pass (kalau ada)

#### Manual Verification:
- [ ] Buku Besar: pilih akun → semua transaksi tampil dengan running balance
- [ ] Buku Besar: filter tanggal → data ter-filter
- [ ] Neraca Saldo: semua akun tampil dengan total debit & kredit
- [ ] Neraca Saldo: baris total di bawah → total debit == total kredit
- [ ] Angka: format ribuan (1.000.000) benar, desimal (,00) benar
- [ ] Empty state: kalau belum ada transaksi, tampilkan pesan yang jelas

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 6: Dashboard

**Status:** ⬜ Not started

**Goal:** Halaman home menampilkan ringkasan keuangan: summary cards + recent entries.

**Files:**
- [ ] `internal/dashboard/handler.go` — Aggregate data via account.Service + journal.Service
- [ ] Update `templates/dashboard/index.html` — Summary cards + recent journal entries table
- [ ] Update `internal/server/routes.go` — Set "/" → dashboard
- [ ] Update `cmd/web/main.go` — Wire dashboard handler (inject services)
- [ ] Update `static/css/style.css` — .dashboard-* styles (cards layout)

**Steps:**
1. Buat handler.go: query total aset, kewajiban, pendapatan, beban + 10 jurnal terakhir
2. Update template: cards grid + recent entries table
3. Set sebagai halaman home (/)

**Success Criteria:**

#### Automated Verification:
- [ ] `go build ./cmd/web/` — berhasil

#### Manual Verification:
- [ ] Dashboard: 4 summary cards tampil (Aset, Kewajiban, Pendapatan, Beban) dengan angka benar
- [ ] Dashboard: 10 jurnal terakhir tampil
- [ ] Dashboard: angka format Rp (Rp 1.000.000)
- [ ] Navigation: klik logo/home → kembali ke dashboard
- [ ] Dashboard jadi default halaman saat buka root URL

**⏸️ Pause untuk manual verification sebelum lanjut.**

---

## Phase 7: Polish & Edge Cases

**Status:** ⬜ Not started

**Goal:** App terasa polished dan production-ready (secara UI/UX).

**Files:**
- [ ] Update semua templates — Responsive (mobile-friendly)
- [ ] Update `static/css/style.css` — Media queries, hover effects, transitions, animations
- [ ] Update `templates/components/_toast.html` — Animate in/out
- [ ] Buat error templates — 404.html, 500.html
- [ ] Update `templates/layout/base.html` — Active menu state berdasarkan current path
- [ ] Update `internal/server/routes.go` — Custom 404/500 handlers
- [ ] Update semua forms — Client-side validation (Alpine.js)
- [ ] Update HTMX elements — Loading indicators, hx-indicator

**Steps:**
1. Responsive: sidebar toggle di mobile, tabel horizontal scroll
2. Error pages: 404 dan 500 dengan design yang consistent
3. Toast: slide-in animation, auto-dismiss 3s
4. Active menu: chi middleware atau template logic, highlight current page di sidebar
5. Loading states: hx-indicator spinner saat request in-flight
6. Form UX: disable submit button saat loading, clear form setelah sukses
7. HTMX + Alpine morph: verify state preserved saat swap
8. Accessibility check: labels, keyboard nav, semantic HTML
9. Cross-check semua halaman, fix visual inconsistencies

**Success Criteria:**

#### Automated Verification:
- [ ] `go build ./cmd/web/` — berhasil
- [ ] `go test ./...` — semua existing test tetap pass

#### Manual Verification:
- [ ] Responsive: semua halaman usable di viewport 375px
- [ ] 404: akses URL invalid → custom 404 page
- [ ] Toast: animasi smooth, auto-dismiss
- [ ] Menu: item aktif ter-highlight sesuai halaman
- [ ] Loading: spinner terlihat saat HTMX request
- [ ] Forms: validasi client-side berjalan (required fields, numeric)
- [ ] Overall: app terasa "finished", bukan prototype

**⏸️ Final review dengan user.**

---

## Testing Strategy

### Unit Tests (per phase):
- **Phase 3:** `account/service_test.go` — validate code unique, name required, type valid
- **Phase 4:** `journal/service_test.go` — debit==credit, single line, both debit+credit, valid entry
- **Phase 5:** `report/service_test.go` (opsional) — running balance calculation

### Mocking Strategy:
- Mock **Repository interface** untuk test Service
- Gunakan Go standard `testing` package (no external test framework)
- Table-driven tests dimana applicable

### Integration Tests (opsional, setelah MVP):
- Full CRUD flow: create account → create journal → verify ledger → verify trial balance

### Manual Testing Checklist (per phase):
- Dicantumkan di "Manual Verification" section masing-masing phase

---

## Risks & Mitigations

| Risk | Mitigation |
|---|---|
| PostgreSQL belum terinstall | Docker: `docker run -e POSTGRES_PASSWORD=pass -p 5432:5432 postgres:15` |
| Template error sulit di-debug | Parse all templates saat startup (fail fast), bukan lazy |
| HTMX swap hilangkan Alpine state | Alpine Morph plugin + `hx-ext="morph"` pada elemen yang punya x-data |
| Phase 4 terlalu besar (1 sesi) | Pecah jadi 4a (data layer) dan 4b (presentation layer) |
| Float precision untuk uang | NUMERIC(15,2) di PostgreSQL, format string di Go |
| CSS collision antar fitur | BEM naming: `.account-*` vs `.journal-*` vs `.c-*` (shared) |

## Decisions Log

- [2026-03-03] Feature-based modular structure — sesuai research v3
- [2026-03-03] PostgreSQL bukan SQLite — future-proof, concurrent writes
- [2026-03-03] MVP tanpa auth — fokus fitur akuntansi dulu
- [2026-03-03] Testing embedded di setiap phase — filosofi HumanLayer
- [2026-03-03] 7 phases, each in separate session — prevent context shrinking

## Progress Notes

- 2026-03-03 — Plan v1 dibuat. Research v3 sudah lengkap. Siap implementasi.
- 2026-03-03 — Phase 1 completed. 15 files dibuat, build berhasil. PostgreSQL belum tersedia (perlu start OrbStack/Docker). Manual verification pending.
- 2026-03-03 — Phase 2 completed. 3 migration files + 1 seed file. Tambah DB-level constraint (debit XOR credit). Manual verification done.
- 2026-03-03 — Phase 3 completed. Account CRUD: 5 Go files + 2 templates + CSS. 11 unit tests pass. List + form pages render. Manual testing untuk HTMX CRUD flow masih pending user test.
- 2026-03-03 — Phase 4 completed. Journal Entry CRUD: 5 Go files + 2 templates + CSS. 10 unit tests (21 total). DB transactions for atomic create/update. Alpine.js dynamic form with live balance calculation. Cross-module validation (account exists). Manual CRUD testing pending.

---

## 📋 Changelog

| Versi | Tanggal    | Perubahan |
|-------|------------|-----------|
| v5    | 2026-03-03 | Phase 4 completed, journal CRUD, 21 total tests |
| v4    | 2026-03-03 | Phase 3 completed, account CRUD, 11 tests pass |
| v3    | 2026-03-03 | Phase 2 completed, migrations + seed created |
| v2    | 2026-03-03 | Phase 1 completed, all files created, build pass |
| v1    | 2026-03-03 | Initial plan: 7 phases, testing strategy, risks |
