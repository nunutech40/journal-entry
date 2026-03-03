# .agent/outputs

Folder ini berisi output dari workflow, dibagi ke subfolder masing-masing:

```
outputs/
├── README.md          ← File ini
├── research/          ← Hasil research codebase (dari /context-engineering)
├── plans/             ← Implementation plan (dari /solve-hard-problems)
└── errors/            ← Error log untuk error kompleks (dari /error-recovery)
```

---

## Cara Pakai

1. **Sesi baru?** → Cek folder `research/` dan `plans/` dulu
2. **Mulai research?** → Simpan output ke `research/[topik].md`
3. **Mulai plan?** → Simpan output ke `plans/[topik].md`
4. **Error kompleks?** → Log ke `errors/[topik].md`

---

## Naming Convention

- Gunakan huruf kecil, pisahkan dengan dash: `research/auth-module.md`
- Topik harus deskriptif tapi singkat: `plans/websocket-multiplayer.md`
- Untuk error: `errors/build-crash-ios.md`

---

## ⚠️ Versioning Rules (PENTING)

### Prinsip: 1 File = 1 Source of Truth

Konten file SELALU versi terbaru. Perubahan dicatat di section Changelog di bawah file.

### Cara Update:

```
KALAU research/plan perlu di-update (kesimpulan berubah):
  1. UPDATE langsung konten yang berubah (jangan duplicate)
  2. Update field `version` dan `date` di frontmatter
  3. Tambahkan entry baru di section "Changelog" di BAWAH file
  4. Changelog cukup 1-2 baris: apa yang berubah + kenapa

JANGAN:
  → Copy-paste seluruh konten lama sebagai "archived version"
  → Bikin file jadi bloated dengan repeated content
  → Detail history = tugas git, bukan file markdown
```

### Contoh Changelog:

```markdown
## 📋 Changelog
| Versi | Tanggal    | Perubahan                          |
|-------|------------|-------------------------------------|
| v3    | 2026-03-03 | Tambah: architecture, ERD, shared code strategy |
| v2    | 2026-03-03 | Database: SQLite → PostgreSQL       |
| v1    | 2026-03-03 | Initial research                    |
```

### Rules untuk AI:

1. **Selalu baca file dari atas** — konten sudah versi terbaru
2. **Cek Changelog di bawah** — untuk pahami evolusi keputusan
3. **Jangan suggest keputusan yang sudah ditolak** tanpa alasan kuat (cek changelog kenapa berubah)
4. **Update, jangan duplicate** — ubah konten langsung, tambah 1 baris di changelog

---

## Lifecycle File

```
1. BARU    → Bikin file research/plan baru di subfolder yang sesuai
2. UPDATE  → Update konten langsung + tambah entry di Changelog
3. DONE    → Set status: completed di frontmatter
4. CLEANUP → File yang sudah >3 bulan dan irrelevant boleh dihapus (opsional)
```

