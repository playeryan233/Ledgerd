package service

import (
	"fmt"
	"strings"

	"ledgerd/internal/domain"
	"ledgerd/internal/storage"
)

// LedgerService orchestrates validation and persistence logic.
type LedgerService struct {
	store storage.Store
}

// NewLedgerService constructs a ledger service.
func NewLedgerService(store storage.Store) *LedgerService {
	return &LedgerService{store: store}
}

// Validate enforces double-entry rules on a journal entry.
func (s *LedgerService) Validate(entry *domain.JournalEntry) error {
	return domain.ValidateEntry(entry)
}

// AddEntry validates and persists a journal entry.
func (s *LedgerService) AddEntry(entry *domain.JournalEntry) error {
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}

	if err := s.Validate(entry); err != nil {
		return err
	}

	entries, err := s.store.LoadEntries()
	if err != nil {
		return err
	}

	if entry.ID == 0 {
		entry.ID = nextEntryID(entries)
	} else {
		for _, existing := range entries {
			if existing.ID == entry.ID {
				return fmt.Errorf("entry with id %d already exists", entry.ID)
			}
		}
	}

	return s.store.AppendEntry(*entry)
}

// LoadEntries returns all persisted journal entries.
func (s *LedgerService) LoadEntries() ([]domain.JournalEntry, error) {
	return s.store.LoadEntries()
}

// ComputeBalance calculates the account balance using ledger rules.
func (s *LedgerService) ComputeBalance(account string) (float64, error) {
	account = strings.TrimSpace(account)
	if account == "" {
		return 0, fmt.Errorf("account is required")
	}

	entries, err := s.store.LoadEntries()
	if err != nil {
		return 0, err
	}

	var debitTotal, creditTotal float64
	for _, entry := range entries {
		for _, line := range entry.Lines {
			if line.Account == account {
				debitTotal += line.Debit
				creditTotal += line.Credit
			}
		}
	}

	if inferAccountType(account) == accountTypeDebitNormal {
		return debitTotal - creditTotal, nil
	}

	return creditTotal - debitTotal, nil
}

func nextEntryID(entries []domain.JournalEntry) int64 {
	var maxID int64
	for _, entry := range entries {
		if entry.ID > maxID {
			maxID = entry.ID
		}
	}
	return maxID + 1
}

type accountType int

const (
	accountTypeDebitNormal accountType = iota
	accountTypeCreditNormal
)

func inferAccountType(account string) accountType {
	account = strings.TrimSpace(account)
	if account == "" {
		return accountTypeDebitNormal
	}

	prefix := strings.ToLower(account)
	if idx := strings.Index(prefix, ":"); idx >= 0 {
		prefix = prefix[:idx]
	}

	switch prefix {
	case "assets", "asset", "expenses", "expense":
		return accountTypeDebitNormal
	case "liabilities", "liability", "income", "revenue", "equity":
		return accountTypeCreditNormal
	default:
		return accountTypeDebitNormal
	}
}
