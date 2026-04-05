package storage

import "ledgerd/internal/domain"

// Store abstracts persistence for journal entries.
type Store interface {
	LoadEntries() ([]domain.JournalEntry, error)
	AppendEntry(entry domain.JournalEntry) error
}

