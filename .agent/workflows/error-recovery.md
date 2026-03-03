---
description: Structured error recovery — no infinite loops. Adapted from HumanLayer's Factor 9 (Compact Errors) + Factor 8 (Own Control Flow).
---

# 🔧 Error Recovery Workflow

Workflow ini memastikan error handling yang terstruktur, tanpa spin-out atau infinite loop.
Diadaptasi dari HumanLayer's "Compact Errors into Context Window" + "Own Your Control Flow".

**Output: Error log di `.agent/outputs/errors/[topik].md` (kalau error kompleks)**

---

## Prinsip Utama

1. **Error = Data** — Setiap error memberikan informasi untuk perbaikan
2. **Max 3 Attempts** — Jangan loop tanpa batas, eskalasi ke user setelah 3x
3. **Analyze, Don't Repeat** — Setiap attempt harus beda strategi
4. **Escalate Early** — Lebih baik tanya daripada buang waktu
5. **Log Complex Errors** — Error yang butuh >1 attempt, simpan ke file

---

## Error Recovery Protocol

### Attempt 1: Analyze & Fix
```
1. Baca error message dengan TELITI (jangan skim)
2. Identifikasi root cause:
   - Syntax error? → Fix typo/syntax
   - Import error? → Cek dependency/path
   - Type error? → Cek model/entity definition
   - Runtime error? → Trace logic flow
3. Terapkan fix yang targeted
4. Verify → build/test lagi
```

### Attempt 2: Reassess Strategy
```
Kalau Attempt 1 gagal:
1. Jangan ulangi strategi yang sama!
2. Pertanyakan pendekatan:
   - Apakah pemahaman gue terhadap masalah sudah benar?
   - Apakah ada file/context yang belum gue baca?
   - Apakah ada pattern lain di codebase yang bisa diikuti?
3. Cari contoh serupa di codebase (grep_search)
4. Baca documentation kalau perlu
5. Cek .agent/outputs/research/ dan .agent/outputs/plans/ untuk konteks yang relevan
6. Terapkan pendekatan yang BEDA
7. Verify → build/test lagi
```

### Attempt 3: Last Try with Fresh Eyes
```
Kalau Attempt 2 masih gagal:
1. Step back — lihat masalah dari sudut pandang berbeda
2. Pertimbangkan:
   - Apakah ada bug di library/framework?
   - Apakah ada version mismatch?
   - Apakah environment/config bermasalah?
3. Coba pendekatan yang completely different
4. Verify → build/test lagi
```

### Escalation: Contact Human
```
Kalau Attempt 3 masih gagal — STOP dan tanya user:

"Gue sudah coba 3 pendekatan berbeda untuk masalah ini:

1. [Attempt 1]: [apa yang dicoba] → [kenapa gagal]
2. [Attempt 2]: [apa yang dicoba] → [kenapa gagal]  
3. [Attempt 3]: [apa yang dicoba] → [kenapa gagal]

Root cause yang gue curigai: [analisis]

Opsi yang bisa kita coba:
- A: [opsi A]
- B: [opsi B]
- C: [opsi lain]

Mau lanjut yang mana, atau ada ide lain?"
```

---

## 💾 Log Error Kompleks (Opsional)

Untuk error yang butuh >1 attempt atau error yang mungkin muncul lagi, simpan ke file:

```
Buat file: .agent/outputs/errors/[topik].md

Format:
---
topic: [Topik/fitur yang bermasalah]
date: [Tanggal]
status: resolved | unresolved | escalated
resolution_attempts: [jumlah attempt]
---

# Error Log: [Topik]

## Error Description
[Deskripsi error yang muncul, termasuk error message lengkap]

## Root Cause
[Analisis root cause yang teridentifikasi]

## Attempts
### Attempt 1
- **Strategy:** [apa yang dicoba]
- **Result:** ❌ Gagal / ✅ Berhasil
- **Learning:** [apa yang dipelajari]

### Attempt 2
- **Strategy:** [apa yang dicoba]
- **Result:** ❌ Gagal / ✅ Berhasil
- **Learning:** [apa yang dipelajari]

## Resolution
[Solusi yang berhasil, atau status eskalasi]

## Prevention
[Bagaimana mencegah error ini di masa depan]
```

---

## Error Categories & Quick Reference

### Build/Compile Errors
| Error Type | First Check | Common Fix |
|---|---|---|
| Import not found | File path, package name | Fix import path |
| Type mismatch | Model definition | Update types |
| Syntax error | Recent edits | Fix syntax |
| Missing dependency | pubspec.yaml / go.mod / Package.swift | Add dependency |

### Runtime Errors
| Error Type | First Check | Common Fix |
|---|---|---|
| Null pointer | Data flow, nullable fields | Add null checks |
| State error | BLoC/Cubit state | Fix state transitions |
| API error | Endpoint, payload | Fix request/response |
| Permission error | Config, credentials | Update permissions |

### Logic Errors (No Error Message)
| Symptom | First Check | Common Fix |
|---|---|---|
| Wrong output | Data transformation | Fix logic/mapping |
| UI not updating | State management | Fix emit/notify |
| Data not saving | Repository/API call | Fix persistence |
| Feature not triggering | Condition/route | Fix control flow |

---

## Anti-Patterns

❌ **Error Spinning** — Mencoba hal yang sama berulang kali  
❌ **Panic Fix** — Ubah banyak hal sekaligus tanpa analisis  
❌ **Stack Overflow Driven** — Copy-paste "fix" dari internet tanpa paham  
❌ **Silent Ignore** — Skip error dan lanjut ke task lain  
❌ **Over-Engineering** — Bikin solusi yang terlalu kompleks untuk error sederhana  
❌ **No Documentation** — Error kompleks yang nggak di-log, bisa muncul lagi  
