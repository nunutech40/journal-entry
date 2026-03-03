---
description: Build rich context before coding. Adapted from HumanLayer's Context Engineering methodology (Factor 3 + Factor 13).
---

# 🧠 Context Engineering Workflow

Workflow ini memastikan AI punya konteks yang cukup sebelum mulai ngoding.
Diadaptasi dari HumanLayer's "Own Your Context Window" + "Pre-fetch All Context".

**Inti: AI engineering = Context Engineering. Output bagus = input bagus.**

**Output: File research di `.agent/outputs/research/[topik].md`**

---

## Kapan Pakai Workflow Ini?

Gunakan workflow ini ketika:
- Mengerjakan fitur baru di codebase yang sudah ada
- Debugging masalah yang kompleks
- Refactoring kode yang saling terkait
- Mengerjakan task yang melibatkan banyak file/modul
- Onboarding ke codebase baru

---

## Step 1: Scan Landscape

### 1.1 Identifikasi Area yang Terdampak
// turbo
```
Gunakan find_by_name dan grep_search untuk:
- Cari file yang terkait dengan fitur/modul yang akan diubah
- Cari keyword/function/class yang relevan
- Identifikasi test files yang perlu di-update
```

### 1.2 Baca Outline Dulu, Detail Kemudian
// turbo
```
Untuk setiap file yang relevan:
1. view_file_outline → pahami struktur (class, function, dll)
2. view_code_item → baca detail function/class yang spesifik
3. view_file → baca section tertentu kalau perlu konteks lebih
```

### 1.3 Trace Data Flow
// turbo
```
Ikuti alur data dari awal sampai akhir:
- UI/Screen → BLoC/Cubit → Repository → API/Database
- Atau sebaliknya untuk flow yang dimulai dari backend
- Catat setiap transformasi data yang terjadi
```

---

## Step 2: Gather External Context

### 2.1 Knowledge Items
// turbo
```
- Cek KI summaries untuk topik terkait
- Baca artifact dari KI yang relevan
- Perhatikan troubleshooting notes dari KI sebelumnya
```

### 2.2 Previous Outputs
// turbo
```
Cek apakah sudah ada research/plan sebelumnya:
- Baca file di .agent/outputs/research/ dan .agent/outputs/plans/ yang relevan
- Gunakan sebagai starting point, jangan mulai dari nol
```

### 2.3 Documentation
```
Kalau perlu:
- Baca API documentation (read_url_content)
- Cek library/package documentation
- Review changelog kalau ada upgrade
```

### 2.4 Previous Conversations
```
Kalau task ini kelanjutan dari sebelumnya:
- Cari conversation yang relevan dari summaries
- Baca conversation logs/artifacts kalau perlu detail
```

---

## Step 3: Pre-fetch (Ambil Semua Sekaligus)

### Prinsip: Jangan Tunggu Error Untuk Cari Konteks

```
SEBELUM mulai coding, kumpulkan:

□ Semua file yang akan diubah (sudah dibaca)
□ File-file yang terkait/depend pada file yang akan diubah
□ Test files yang relevan
□ Config files kalau ada
□ Model/Entity definitions
□ API endpoint definitions
□ Existing patterns yang serupa (cari contoh implementasi yang mirip)
```

### Contoh Pre-fetch Checklist untuk Flutter:
```
□ Screen/Page file
□ BLoC/Cubit file
□ State file
□ Repository file
□ Data Source file
□ Entity/Model file
□ Router/Navigation config
□ Dependency injection setup
□ Existing similar feature (sebagai referensi pattern)
```

### Contoh Pre-fetch Checklist untuk Go Backend:
```
□ Handler file
□ Service/UseCase file
□ Repository file
□ Model/Struct definitions
□ Route definitions
□ Middleware config
□ Database migration files
□ Existing similar endpoint (sebagai referensi pattern)
```

### Contoh Pre-fetch Checklist untuk iOS/Swift:
```
□ ViewController file
□ View files (UIKit/SwiftUI)
□ Model/Entity file
□ Network/API client file
□ Storyboard/XIB atau SwiftUI preview
□ Info.plist / config files
□ Dependency setup (CocoaPods/SPM)
□ Existing similar feature (sebagai referensi pattern)
```

---

## Step 4: Synthesize & Plan

### 4.1 Rangkum Temuan
```
Sebelum mulai coding, sampaikan ke user:
- "Berdasarkan analisis gue, ini yang perlu diubah: ..."
- "Pattern yang sudah ada di codebase adalah: ..."
- "Potential risks/edge cases: ..."
```

### 4.2 Confirm Understanding
```
Tanya ke user:
- "Apakah pemahaman gue sudah benar?"
- "Ada hal lain yang perlu gue pertimbangkan?"
- "Mau gue mulai dari mana?"
```

---

## Step 5: 💾 Simpan Research Output (WAJIB)

### Simpan hasil research ke file agar bisa dibaca di sesi berikutnya.

### ⚠️ Versioning Rule:
```
PRINSIP: 1 file = 1 source of truth. Konten SELALU versi terbaru.

KALAU file research untuk topik ini SUDAH ADA:
  → UPDATE langsung konten yang berubah (jangan duplicate)
  → Update field `version` dan `date` di frontmatter
  → Tambahkan entry baru di section "Changelog" di BAWAH file
  → Changelog cukup 1-2 baris: apa yang berubah + kenapa

KALAU belum ada:
  → Buat file baru dengan version: 1

JANGAN:
  → Copy-paste seluruh konten lama sebagai "archived version"
  → Bikin file jadi bloated dengan repeated content
  → Detail history = tugas git, bukan file markdown
```

### Format file:
```
Buat file: .agent/outputs/research/[topik-singkat].md

Isi file:
---
topic: [Nama topik/fitur]
date: [Tanggal terakhir update]
version: [nomor versi, mulai dari 1]
status: completed
related_files:
  - path/to/file1
  - path/to/file2
---

# Research: [Topik]

## Summary
[Ringkasan 2-3 kalimat tentang apa yang ditemukan]

## File yang Relevan
- `path/to/file1` — [deskripsi singkat perannya]
- `path/to/file2` — [deskripsi singkat perannya]

## Architecture / Data Flow
[Diagram atau penjelasan alur data]

## Existing Patterns
[Pattern yang sudah ada di codebase dan harus diikuti]

## Keputusan & Alasan (Decisions)
- [Keputusan 1]: [Alasan kenapa]
- [Keputusan 2]: [Alasan kenapa]

## Risks & Edge Cases
- [Risk 1]
- [Edge case 1]

## Notes
[Catatan tambahan yang penting untuk konteks]

---

## 📋 Changelog
| Versi | Tanggal    | Perubahan                          |
|-------|------------|------------------------------------|
| v1    | YYYY-MM-DD | Initial research                   |
```

---

## Tips Context Engineering

### DO ✅
- **Baca dulu, coding belakangan** — Minimal 30% effort di phase context
- **Pre-fetch aggressively** — Lebih baik kebanyakan context daripada kekurangan
- **Follow existing patterns** — Cari contoh serupa di codebase, jangan bikin dari nol
- **Explain your understanding** — Sampaikan pemahaman ke user sebelum mulai
- **Simpan research** — Selalu save output ke `.agent/outputs/research/`

### DON'T ❌
- **Jangan langsung edit** — Tanpa baca file terkait dulu
- **Jangan asumsi** — Struktur/pattern yang belum diverifikasi
- **Jangan abaikan test** — Selalu cari dan baca test files
- **Jangan skip KI check** — Knowledge Items ada untuk dipakai
- **Jangan buang research** — Selalu simpan ke `.agent/outputs/research/` untuk sesi berikutnya
