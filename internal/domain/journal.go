package domain

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// JournalEntry represents a double-entry journal transaction.
type JournalEntry struct {
	ID          int64         `json:"id"`
	Date        string        `json:"date"`
	Description string        `json:"description"`
	Lines       []JournalLine `json:"lines"`
}

// JournalLine represents a single posting.
type JournalLine struct {
	Account string  `json:"account"`
	Debit   float64 `json:"debit"`
	Credit  float64 `json:"credit"`
}

const balanceTolerance = 1e-9

// ValidateEntry checks the business rules for a journal entry.
func ValidateEntry(entry *JournalEntry) error {
	if entry == nil {
		return fmt.Errorf("journal entry is nil")
	}

	entry.Date = strings.TrimSpace(entry.Date)
	if entry.Date == "" {
		return fmt.Errorf("date is required")
	}
	if _, err := time.Parse("2006-01-02", entry.Date); err != nil {
		return fmt.Errorf("date must be YYYY-MM-DD: %w", err)
	}

	entry.Description = strings.TrimSpace(entry.Description)
	if entry.Description == "" {
		return fmt.Errorf("description is required")
	}

	if len(entry.Lines) < 2 {
		return fmt.Errorf("entry must have at least two lines")
	}

	var totalDebit, totalCredit float64
	for idx, line := range entry.Lines {
		if err := validateLine(line, idx); err != nil {
			return err
		}

		totalDebit += line.Debit
		totalCredit += line.Credit
	}

	if math.Abs(totalDebit-totalCredit) > balanceTolerance {
		return fmt.Errorf("entry debits (%.2f) do not equal credits (%.2f)", totalDebit, totalCredit)
	}

	return nil
}

func validateLine(line JournalLine, idx int) error {
	if strings.TrimSpace(line.Account) == "" {
		return fmt.Errorf("line %d account is required", idx+1)
	}

	if line.Debit < 0 || line.Credit < 0 {
		return fmt.Errorf("line %d has negative values", idx+1)
	}

	return nil
}

