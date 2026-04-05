You are a senior Go backend engineer.

## Goal

Extend the accounting system by implementing a Cash Flow Statement (Direct Method).

---

## Requirements

Create a new function:

func GenerateCashFlowStatement(
    entries []JournalEntry,
    start time.Time,
    end time.Time,
) CashFlowStatement

---

## Cash Definition

Cash accounts are:

- Assets:Cash
- Assets:Bank
- Assets:Alipay
- Assets:WeChat

Use prefix match: "Assets:" AND in predefined cash list.

---

## Core Logic

For each JournalEntry within [start, end]:

1. Check if ANY line involves a cash account
   - If no → ignore this entry

2. For each cash line:

   If debit → cash inflow (+)
   If credit → cash outflow (-)

3. Determine category by the COUNTERPART account:

- If counterpart is "Income:" or "Expenses:" → Operating
- If counterpart is "Assets:" (non-cash) → Investing
- If counterpart is "Liabilities:" or "Equity:" → Financing

---

## Data Structures

type CashFlowStatement struct {
    StartDate time.Time
    EndDate   time.Time

    Operating  []LineItem
    Investing  []LineItem
    Financing  []LineItem

    NetOperating  float64
    NetInvesting  float64
    NetFinancing  float64

    NetIncrease float64

    BeginningCash float64
    EndingCash    float64
}

---

## Beginning / Ending Cash

- BeginningCash:
  sum of all cash account balances before start

- EndingCash:
  sum of all cash account balances at end

---

## Validation

Ensure:

EndingCash = BeginningCash + NetIncrease

If not → return error

---

## CLI

Add command:

cash-flow
  --start YYYY-MM-DD
  --end YYYY-MM-DD

---

## Output Example

CASH FLOW STATEMENT

OPERATING
  Income:Salary        +10000
  Expenses:Food        -2000

INVESTING
  (none)

FINANCING
  (none)

NET INCREASE:          +8000
BEGINNING CASH:        5000
ENDING CASH:           13000

---

## Constraints

- No database
- Pure functions
- Reuse existing balance logic if possible

---

## Edge Cases

- Ignore entries without cash accounts
- Multiple cash lines in one entry → handle correctly
- Skip zero amounts

---

## Deliverable

- New report module
- CLI integration
- Clean, testable code

Do not explain anything.
Only output code.
