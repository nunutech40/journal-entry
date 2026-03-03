---
description: Solve hard coding problems in complex codebases without "vibe coding". Adapted from HumanLayer's 12-Factor Agents & Context Engineering methodology.
---

# 🚫 No Vibes Allowed — Solving Hard Problems

Workflow ini diadaptasi dari metodologi HumanLayer (12-Factor Agents + Advanced Context Engineering).
Gunakan workflow ini ketika menghadapi task yang kompleks di codebase yang besar.

**Prinsip utama: Jangan langsung ngoding. Bangun konteks dulu, baru eksekusi.**

**Output: File plan di `.agent/outputs/plans/[topik].md`**

---

## Phase 1: Context Engineering (WAJIB sebelum nulis kode)

### 1.1 Cek Existing Research & Plan
// turbo
```
PERTAMA, cek apakah sudah ada research/plan sebelumnya:
- Baca file di .agent/outputs/research/ dan .agent/outputs/plans/ yang relevan
- Kalau sudah ada research → skip ke Phase 1.3
- Kalau sudah ada plan → review & update, skip ke Phase 2
- Kalau belum ada → lanjut ke 1.2
```

### 1.2 Pahami Scope
// turbo
- Baca file/folder yang relevan dengan task menggunakan `view_file_outline` dan `view_file`
- Cari pattern yang sudah ada di codebase menggunakan `grep_search` dan `find_by_name`
- Identifikasi dependency dan side-effect dari perubahan yang akan dibuat

### 1.3 Cek Knowledge Items
// turbo
- Cek KI summaries yang sudah ada — mungkin ada KI yang relevan
- Baca artifact dari KI yang relevan
- Cek conversation history kalau ada diskusi terkait sebelumnya

### 1.4 Pre-fetch All Context
// turbo
- Kumpulkan SEMUA informasi yang mungkin dibutuhkan SEBELUM mulai coding
- Baca API docs / library docs kalau perlu
- Pahami data flow end-to-end (dari UI → logic → API → database)
- **Jangan tunggu error untuk baru cari konteks — cari duluan!**

### 1.5 💾 Simpan Research (kalau belum ada)
```
Kalau belum ada file research, simpan ke:
.agent/outputs/research/[topik].md

(Ikuti format di context-engineering.md Step 5)
```

### 1.6 Buat Implementation Plan
```
Tulis rencana implementasi ke file:
.agent/outputs/plans/[topik].md

Format:
---
topic: [Nama topik/fitur]
date: [Tanggal plan dibuat]
status: draft | in-progress | completed
research: ../research/[topik].md
phases_total: [jumlah phase]
phases_completed: 0
---

# Plan: [Topik]

## Summary
[Ringkasan 2-3 kalimat tentang apa yang akan dikerjakan]

## Research Reference
Baca: `.agent/outputs/research/[topik].md`

## Phases

### Phase 1: [Nama Phase]
**Status:** ⬜ Not started
**Files:**
- [ ] `path/to/file1` — [apa yang diubah]
- [ ] `path/to/file2` — [apa yang diubah]

**Steps:**
1. [Step detail]
2. [Step detail]

**Success Criteria:**

#### Automated Verification:
- [ ] Build/compile berhasil: `go build ./...`
- [ ] Unit tests pass: `go test ./...`
- [ ] Lint pass: `golangci-lint run` (kalau ada)
- [ ] [Command spesifik lain, e.g. migration: `goose up`]

#### Manual Verification:
- [ ] [Fitur bisa diakses/dipakai sesuai harapan]
- [ ] [Edge case yang perlu dicek manual]
- [ ] [Nggak ada regresi di fitur terkait]

**⏸️ Setelah semua automated verification pass, pause untuk manual verification dari user sebelum lanjut ke phase berikutnya.**

---

### Phase 2: [Nama Phase]
**Status:** ⬜ Not started
[Sama strukturnya: Files, Steps, Success Criteria (Automated + Manual)]

---

## Testing Strategy

### Unit Tests:
- [Apa yang di-test per layer: handler, service, repository]
- [Key edge cases yang perlu di-cover]
- [Mocking strategy: mock interfaces, bukan concrete]

### Integration Tests:
- [End-to-end scenario yang perlu diverifikasi]
- [Database tests: pakai test DB atau mock?]

### Manual Testing Steps:
1. [Step spesifik untuk verifikasi fitur]
2. [Edge case yang perlu dicek manual]
3. [Performance check kalau relevan]

## Risks & Mitigations
- **Risk:** [deskripsi] → **Mitigation:** [cara handle]

## Decisions Log
- [Keputusan]: [Alasan] (tanggal)

## Progress Notes
- [tanggal] — [catatan progress]
```

### 1.7 Confirm Plan dengan User
```
Tanya ke user:
- "Ini plan-nya, sudah sesuai?"
- "Ada yang perlu ditambah/diubah?"
- "Mau mulai dari Phase berapa?"
```

---

## Phase 2: Small, Focused Execution

### 2.1 Pecah Task Jadi Unit Kecil
- Jangan bikin perubahan besar sekaligus
- Setiap unit harus bisa diverifikasi secara independen
- Urutan: Data Layer → Logic Layer → UI Layer (bottom-up)

### 2.2 ⚠️ Shared Code Impact Check (WAJIB sebelum edit)
```
SEBELUM edit file APAPUN, tanya dulu:

"Apakah file ini SHARED (dipake fitur lain)?"

Cara cek:
  → File di shared/, components/, layout/ → PASTI shared
  → File di internal/[fitur]/ → kemungkinan besar isolated
  → CSS class dengan prefix c- (e.g. .c-btn) → shared
  → CSS class dengan prefix fitur (e.g. .journal-form) → isolated
  → Template yang di-include dari banyak page → shared

KALAU SHARED:
  1. GREP dulu siapa semua yang pake (grep_search di seluruh project)
  2. LIST semua file caller/consumer
  3. Cek: perubahan ini backward-compatible?
     → YA (tambah fungsi baru, tambah parameter optional): langsung ubah
     → NGGAK (ubah signature, ubah return type, ubah HTML structure):
       a. Bikin fungsi/variant BARU, jangan ubah yang lama
       b. ATAU kalau harus ubah yang lama, update SEMUA caller
       c. TEST semua fitur yang terpengaruh
  4. Setelah ubah, verify: jalankan/test fitur lain yang depend

KALAU ISOLATED (file di internal/[fitur]/ sendiri):
  → Langsung edit, nggak perlu cek fitur lain

CONTOH file shared yang sering kesenggol:
  - shared/response/htmx.go → semua handler pake
  - shared/validate/validate.go → semua service pake
  - templates/layout/base.html → semua page extend
  - templates/components/_toast.html → banyak page pake
  - static/css/style.css (.c-* classes) → semua page pake
```

### 2.3 Own Your Control Flow
- Untuk setiap perubahan:
  1. **Check** — cek apakah file shared? (step 2.2)
  2. **Edit** file yang diperlukan
  3. **Verify** — jalankan build/test untuk memastikan tidak break
  4. **Compact errors** — kalau ada error, baca error message dengan teliti, perbaiki, coba lagi (max 3x, lihat `/error-recovery`)
  5. **Escalate** — kalau sudah 3x gagal, TANYA USER sebelum lanjut

### 2.4 Update Plan Setelah Setiap Phase
```
Setelah menyelesaikan setiap phase:
1. Buka file .agent/outputs/plans/[topik].md
2. Update status phase: ⬜ → ✅
3. Update phases_completed di frontmatter
4. Tambah catatan di Progress Notes
5. Kalau ada keputusan baru, tambah di Decisions Log
```

### 2.5 Error Recovery Protocol
- Error pertama: Baca error message, analisis root cause, perbaiki
- Error kedua: Cek apakah pendekatan sudah benar, mungkin perlu ganti strategi
- Error ketiga: STOP. Jelaskan situasi ke user, minta guidance (lihat `/error-recovery`)
- **Jangan loop tanpa batas mencoba hal yang sama!**

---

## Phase 3: Verification & Human Contact

### 3.1 Automated Verification (WAJIB setiap phase)
// turbo
```
Jalankan SEMUA automated checks sesuai Success Criteria di plan:

1. Build/compile:
   go build ./...

2. Unit tests:
   go test ./...
   → Kalau ada test baru yang ditulis di phase ini, pastikan pass
   → Kalau ada test lama yang break, fix SEBELUM lanjut

3. Lint (kalau ada):
   golangci-lint run

4. Command spesifik lain (migration, seed, dll)

SEMUA harus pass sebelum lanjut ke manual verification.
Kalau ada yang gagal → fix dulu (max 3 attempt, lihat /error-recovery)
```

### 3.2 Manual Verification (minta User)
```
Setelah automated verification pass:
1. Jelaskan ke user apa yang perlu di-test manual
2. Tunggu konfirmasi dari user
3. JANGAN lanjut ke phase berikutnya tanpa konfirmasi

Contoh:
  "Automated checks semua pass ✅
   Tolong cek manual:
   - [ ] Buka halaman X, pastikan Y tampil
   - [ ] Coba input Z, pastikan validasi jalan
   Sudah oke? Lanjut ke Phase berikutnya?"
```

### 3.3 Contact Human (User) untuk:
- **Approval** — sebelum melakukan perubahan yang high-stakes (hapus file, ubah schema DB, dll)
- **Clarification** — kalau ada yang ambigu di requirement
- **Review** — setelah selesai, jelaskan apa yang sudah diubah dan kenapa

### 3.4 💾 Final Update Plan
```
Update file .agent/outputs/plans/[topik].md:
- Set status: completed
- Update semua phase status
- Tambah final notes di Progress Notes
- Pastikan Decisions Log lengkap
- Pastikan Testing Strategy section ter-update (test apa yang sudah ada)
```

### 3.5 Summary ke User
```
- Berikan ringkasan perubahan yang sudah dibuat
- List file yang diubah
- Jelaskan keputusan desain yang non-obvious
- Sebutkan limitasi atau hal yang belum di-cover
- Status testing: test apa saja yang sudah ditulis dan pass
- Reference: "Detail lengkap ada di .agent/outputs/plans/[topik].md"
```

---

## Anti-Patterns (JANGAN LAKUKAN)

❌ **Vibe Coding** — Langsung nulis kode tanpa paham konteks  
❌ **Boil the Ocean** — Coba ubah semuanya sekaligus  
❌ **Error Spinning** — Loop terus mencoba hal yang sama tanpa analisis  
❌ **Assumption-Driven** — Asumsi tanpa verifikasi  
❌ **Copy-Paste Blindly** — Copy kode dari tempat lain tanpa adaptasi  
❌ **Ignore Existing Patterns** — Bikin pattern baru padahal sudah ada yang serupa di codebase  
❌ **Throwaway Research** — Research tanpa simpan output ke file  

## Best Practices (LAKUKAN)

✅ **Context First** — Selalu bangun konteks sebelum coding  
✅ **Check Existing Outputs** — Baca research/plan sebelumnya di `.agent/outputs/research/` dan `.agent/outputs/plans/`  
✅ **Pre-fetch** — Ambil semua info yang mungkin diperlukan di awal  
✅ **Small Steps** — Pecah task besar jadi langkah kecil  
✅ **Verify Each Step** — Verifikasi setiap perubahan sebelum lanjut (automated lalu manual)  
✅ **Test as You Go** — Tulis unit test bersamaan dengan kode, bukan nanti  
✅ **Ask When Unsure** — Tanya user kalau ragu  
✅ **Follow Existing Patterns** — Ikuti convention yang sudah ada di codebase  
✅ **Explain Decisions** — Jelaskan kenapa, bukan cuma apa  
✅ **Save Everything** — Simpan research ke `.agent/outputs/research/`, plan ke `.agent/outputs/plans/`  
✅ **Update Plan** — Selalu update progress di file plan  
