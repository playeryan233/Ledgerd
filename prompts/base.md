You are a senior backend engineer.

## Goal
Build a minimal, production-quality foundation for a personal accounting system using double-entry bookkeeping.

This is Phase 1 (MVP). Focus on correctness and simplicity.

---

## Requirements

### 1. Tech Stack
- Language: Go
- CLI: cobra
- Storage: file-based (JSON), no database
- Project structure: clean architecture (modular, maintainable)

---

### 2. Core Domain (VERY IMPORTANT)

Implement a double-entry accounting system:

#### Entities:

JournalEntry:
- ID
- Date
- Description
- Lines []JournalLine

JournalLine:
- Account (string, e.g. "Assets:Cash")
- Debit (float64)
- Credit (float64)

---

### 3. Business Rules

1. Every JournalEntry must satisfy:
   sum(debit) == sum(credit)

2. No negative values allowed

3. At least 2 lines per entry

4. Accounts are strings (no need for full COA yet)

---

### 4. Storage

- Persist entries as JSON files
- File structure:

data/
  journal.json

- Append-only (do not modify existing entries)

---

### 5. Ledger Engine

Implement:

- Validate(entry)
- AddEntry(entry)
- LoadEntries()
- ComputeBalance(account)

Balance rules:
- Assets / Expenses → debit - credit
- Liabilities / Income / Equity → credit - debit

(You can infer type by prefix: Assets:, Liabilities:, etc.)

---

### 6. CLI Commands

Implement using cobra:

1. add
   - input: JSON file or flags
   - adds a journal entry

2. list
   - print all entries

3. balance
   - input: account name
   - output: current balance

---

### 7. Code Structure

Use this layout:

/cmd
/internal
  /domain
  /service
  /storage
  /cli

---

### 8. Output Requirements

- Generate full working Go project
- Include:
  - go.mod
  - main.go
  - all packages
- Code must compile

---

### 9. Constraints

- No frameworks except cobra
- Keep code simple and readable
- No over-engineering
- Do not implement reports yet

---

### 10. Example Entry

{
  "date": "2026-04-05",
  "description": "Lunch",
  "lines": [
    {"account": "Expenses:Food", "debit": 50, "credit": 0},
    {"account": "Assets:Cash", "debit": 0, "credit": 50}
  ]
}

---

## Deliverable

Return a complete Go project with all files.
Do not explain, only output code.
