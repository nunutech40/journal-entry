# 📖 Cara Pakai .agent Workflows

Template ini bisa di-copy ke project mana pun. Tinggal copy folder `.agent/` ke root project baru.

```bash
cp -r .agent/ ~/Documents/ProjectBaru/.agent/
```

---

## 📁 Struktur Folder

```
.agent/
├── HOW_TO_USE.md              ← File ini
├── workflows/                 ← Instruksi workflow untuk AI
│   ├── context-engineering.md ← Research & gather context
│   ├── solve-hard-problems.md ← Plan & implement step-by-step
│   └── error-recovery.md     ← Handle error terstruktur
└── outputs/                   ← Output dari workflow (auto-generated)
    ├── README.md              ← Panduan naming & versioning
    ├── research/              ← Hasil research
    │   └── [topik].md
    ├── plans/                 ← Implementation plan
    │   └── [topik].md
    └── errors/                ← Error log
        └── [topik].md
```

---

## 🚀 Quick Start

### 1. Research Dulu Sebelum Coding

Ketik:
```
/context-engineering — [Deskripsi apa yang mau di-research]
```

Contoh:
```
/context-engineering — Gue mau bikin fitur auth pakai Firebase di Flutter. 
Research dulu: package yang cocok, arsitektur yang bagus, 
dan flow login/register.
```

**Apa yang terjadi:**
- AI scan codebase, baca docs, cek Knowledge Items
- AI rangkum temuan dan tanya konfirmasi
- Output disimpan ke `.agent/outputs/research/[topik].md`

---

### 2. Bikin Plan & Implement

Ketik:
```
/solve-hard-problems — [Deskripsi task yang mau dikerjakan]
```

Contoh:
```
/solve-hard-problems — Implement payment gateway Midtrans di Go backend.
Butuh: create invoice, handle webhook, update payment status.
```

**Apa yang terjadi:**
- AI cek existing research di `.agent/outputs/`
- AI bikin implementation plan per-phase
- Plan disimpan ke `.agent/outputs/plans/[topik].md`
- AI implement per-phase, update progress di plan
- Setiap phase di-verify (build/test) sebelum lanjut

---

### 3. Lanjutin Task dari Sesi Sebelumnya

Ketik:
```
/solve-hard-problems — Lanjutin dari plan di .agent/outputs/plans/[topik].md
```

Contoh:
```
/solve-hard-problems — Lanjutin dari plan di .agent/outputs/plans/midtrans.md.
Phase 1 sudah done, lanjut Phase 2.
```

**Apa yang terjadi:**
- AI baca plan file → langsung tahu progress & konteks
- AI lanjut dari phase yang belum selesai
- Nggak perlu jelaskan ulang dari awal

---

### 4. Handle Error yang Stuck

Ketik:
```
/error-recovery
```

Atau langsung jelaskan error-nya:
```
/error-recovery — Build error setelah tambah package X. 
Error message: [paste error]
```

**Apa yang terjadi:**
- AI coba fix max 3 attempt dengan strategi berbeda
- Kalau 3x gagal → AI stop dan tanya kamu
- Error kompleks di-log ke `.agent/outputs/errors/[topik].md`

---

## 📋 Cheat Sheet

| Situasi | Prompt |
|---|---|
| **Project baru, mau research dulu** | `/context-engineering — Research untuk bikin [app]` |
| **Mau langsung bikin plan** | `/solve-hard-problems — Bikin plan untuk [fitur]` |
| **Research lalu plan** | Sesi 1: `/context-engineering`, Sesi 2: `/solve-hard-problems` |
| **Lanjutin kemarin** | `/solve-hard-problems — Lanjutin dari .agent/outputs/plans/[topik].md` |
| **Error stuck** | `/error-recovery — [jelaskan error]` |
| **Task simpel** | Langsung chat biasa, nggak perlu slash command |

---

## 💡 Tips

### Kapan Pakai Slash Command vs Chat Biasa?

| Task | Pakai |
|---|---|
| Hapus file, rename, task kecil | Chat biasa |
| Fix bug simpel (1 file) | Chat biasa |
| Tanya-tanya | Chat biasa |
| Fitur baru (multi-file) | `/solve-hard-problems` |
| Onboarding codebase baru | `/context-engineering` |
| Refactor besar | `/context-engineering` → `/solve-hard-problems` |
| Debug yang bikin frustasi | `/error-recovery` |

### Best Practices

1. **Pecah conversation per task** — Jangan 1 conversation super panjang. 
   Konteks AI akan shrink di conversation panjang. 
   Mulai conversation baru per fitur/task.

2. **1 sesi = 1 phase** — Kalau plan punya 3 phase, 
   kerjakan masing-masing di conversation terpisah. 
   Plan file akan jaga konteks antar sesi.

3. **Bilang "simpan ke .agent/outputs/"** — Kalau AI lupa simpan output
   ke subfolder yang sesuai (research/, plans/, errors/), ingatkan aja.

4. **Review output files** — Sesekali baca file di `.agent/outputs/` 
   untuk pastikan isinya akurat. Koreksi kalau ada yang salah.

5. **Tulis keputusan penting di kode** — Comment di kode lebih reliable 
   daripada file manapun:
   ```
   // DECISION: Pakai DI bukan Singleton karena rencana support 
   // multi-player lokal. Jangan ubah tanpa diskusi.
   ```

---

## 🔄 Workflow Pipeline (Full)

```
Sesi 1: Research
  /context-engineering — [deskripsi]
  └── Output: .agent/outputs/research/[topik].md

Sesi 2: Plan  
  /solve-hard-problems — Buat plan berdasarkan research di .agent/outputs/research/[topik].md
  └── Output: .agent/outputs/plans/[topik].md

Sesi 3: Implement Phase 1
  /solve-hard-problems — Implement Phase 1 dari .agent/outputs/plans/[topik].md
  └── Update: plans/[topik].md (Phase 1 ✅)

Sesi 4: Implement Phase 2
  /solve-hard-problems — Lanjut Phase 2 dari .agent/outputs/plans/[topik].md
  └── Update: plans/[topik].md (Phase 2 ✅)

...dan seterusnya sampai semua phase selesai.
```

---

## ⚠️ Versioning

Kalau research/plan perlu di-update (kesimpulan berubah):
- **Update langsung** konten yang berubah (jangan duplicate/copy-paste)
- **Tambah entry di Changelog** di bawah file (1-2 baris: apa + kenapa)
- **Update version** di frontmatter

Detail lengkap: lihat `.agent/outputs/README.md`

---

*Template ini diadaptasi dari metodologi HumanLayer (Context Engineering + 12-Factor Agents).*
*Kompatibel dengan Antigravity dan AI coding assistant manapun yang bisa baca file.*
