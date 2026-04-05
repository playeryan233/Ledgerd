You are a senior Go backend engineer.

## Goal

Extend an existing double-entry accounting system (file-based, JSON storage, CLI) by implementing a Report Engine.

This is Phase 2.

You MUST build:
1. Balance Sheet (资产负债表)
2. Income Statement (利润表)

Do NOT modify existing core logic unless necessary.
Do NOT break backward compatibility.

---

## Existing System Context

The system already supports:
- JournalEntry with multiple lines (debit/credit)
- File-based storage (journal.json)
- Ledger logic (validation + balance calculation)
- CLI commands: add, list, balance

Accounts are strings like:
- Assets:Cash
- Liabilities:Huabei
- Income:Salary
- Expenses:Food

---

## Requirements

### 1. Report Engine Module

Create a new package:

/internal/report

---

### 2. Balance Sheet

Implement:

func GenerateBalanceSheet(entries []JournalEntry, asOf time.Time) BalanceSheet

---

### Balance Sheet Rules

1. Only include entries with date <= asOf

2. Account classification by prefix:

- "Assets:" → Assets
- "Liabilities:" → Liabilities
- "Equity:" → Equity

3. Income and Expenses MUST NOT appear directly

Instead:

Retained Earnings = total Income - total Expenses

Add this to Equity.

---

### Balance Calculation

- Assets / Expenses → debit - credit
- Liabilities / Income / Equity → credit - debit

---

### Output Structure

type BalanceSheet struct {
    Date time.Time

    Assets      []LineItem
    Liabilities []LineItem
    Equity      []LineItem

    TotalAssets      float64
    TotalLiabilities float64
    TotalEquity      float64
}

type LineItem struct {
    Name   string
    Amount float64
}

---

### 4. Income Statement

Implement:

func GenerateIncomeStatement(entries []JournalEntry, start, end time.Time) IncomeStatement

---

### Rules

1. Only include entries within [start, end]

2. Classification:

- "Income:" → Revenue
- "Expenses:" → Expense

3. Net Income:

NetIncome = TotalIncome - TotalExpense

---

### Output Structure

type IncomeStatement struct {
    StartDate time.Time
    EndDate   time.Time

    Income   []LineItem
    Expenses []LineItem

    TotalIncome  float64
    TotalExpense float64
    NetIncome    float64
}

---

### 5. Aggregation Logic

- Group by full account name
- Sum all balances per account

---

### 6. CLI Integration

Add commands:

1. balance-sheet
   --date YYYY-MM-DD

2. income-statement
   --start YYYY-MM-DD
   --end YYYY-MM-DD

---

### CLI Output (text format)

Example:

BALANCE SHEET (2026-04-30)

ASSETS
  Assets:Cash        1000
  Assets:Bank        5000

LIABILITIES
  Liabilities:Huabei 2000

EQUITY
  Retained Earnings  4000

TOTAL ASSETS:        6000
TOTAL L+E:           6000

---

INCOME STATEMENT (2026-04-01 ~ 2026-04-30)

INCOME
  Income:Salary      10000

EXPENSES
  Expenses:Food      2000

NET INCOME:          8000

---

### 7. Constraints

- No database
- Use existing storage layer
- Do not over-engineer
- Keep functions pure (no side effects)

---

### 8. Edge Cases

- Ignore zero-balance accounts
- Handle empty data gracefully
- Ensure Assets == Liabilities + Equity

If not equal → return error

---

## Deliverable

- All new code under /internal/report
- Updated CLI commands
- Fully compilable project

Do not explain anything.
Only output code.
