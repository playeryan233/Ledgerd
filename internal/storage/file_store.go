package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"ledgerd/internal/domain"
)

// FileStore persists journal entries as a JSON array.
type FileStore struct {
	path string
}

// NewFileStore prepares a file-backed store at the provided path.
func NewFileStore(path string) (*FileStore, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("storage path is required")
	}

	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		if err := os.WriteFile(path, []byte("[]"), 0o644); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &FileStore{path: path}, nil
}

// LoadEntries returns all journal entries from disk.
func (s *FileStore) LoadEntries() ([]domain.JournalEntry, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []domain.JournalEntry{}, nil
	}

	var entries []domain.JournalEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

// AppendEntry writes the new entry without altering previous data.
func (s *FileStore) AppendEntry(entry domain.JournalEntry) error {
	entries, err := s.LoadEntries()
	if err != nil {
		return err
	}

	entries = append(entries, entry)
	payload, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, payload, 0o644)
}

var _ Store = (*FileStore)(nil)

