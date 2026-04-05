package cli

import (
	"ledgerd/internal/domain"
	"ledgerd/internal/service"
	"ledgerd/internal/storage"
)

// App wires CLI commands to the ledger service layer.
type App struct {
	ledger *service.LedgerService
}

// NewApp creates a CLI application bound to a data store.
func NewApp(dataPath string) (*App, error) {
	store, err := storage.NewFileStore(dataPath)
	if err != nil {
		return nil, err
	}

	return &App{
		ledger: service.NewLedgerService(store),
	}, nil
}

// AddEntry validates and persists a journal entry.
func (a *App) AddEntry(entry *domain.JournalEntry) error {
	return a.ledger.AddEntry(entry)
}

// ListEntries fetches all journal entries.
func (a *App) ListEntries() ([]domain.JournalEntry, error) {
	return a.ledger.LoadEntries()
}

// ComputeBalance calculates a single account balance.
func (a *App) ComputeBalance(account string) (float64, error) {
	return a.ledger.ComputeBalance(account)
}

