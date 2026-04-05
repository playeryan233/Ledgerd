package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"ledgerd/internal/domain"
)

var (
	addFilePath     string
	addDate         string
	addDescription  string
	addLineSegments []string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a journal entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		var entry *domain.JournalEntry
		var err error

		if addFilePath != "" {
			if addDate != "" || addDescription != "" || len(addLineSegments) > 0 {
				return fmt.Errorf("--file cannot be combined with other entry flags")
			}

			entry, err = loadEntryFromFile(addFilePath)
		} else {
			entry, err = buildEntryFromFlags()
		}

		if err != nil {
			return err
		}

		if err := app.AddEntry(entry); err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Entry recorded with id %d\n", entry.ID)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addFilePath, "file", "", "Path to JSON file containing a journal entry")
	addCmd.Flags().StringVar(&addDate, "date", "", "Entry date (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Entry description")
	addCmd.Flags().StringArrayVar(&addLineSegments, "line", nil, "Journal line as account,debit,credit (repeatable)")

	rootCmd.AddCommand(addCmd)
}

func loadEntryFromFile(path string) (*domain.JournalEntry, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entry domain.JournalEntry
	if err := json.Unmarshal(bytes, &entry); err != nil {
		return nil, err
	}

	return &entry, nil
}

func buildEntryFromFlags() (*domain.JournalEntry, error) {
	if strings.TrimSpace(addDate) == "" {
		return nil, fmt.Errorf("--date is required when --file is not provided")
	}

	if strings.TrimSpace(addDescription) == "" {
		return nil, fmt.Errorf("--description is required when --file is not provided")
	}

	if len(addLineSegments) < 2 {
		return nil, fmt.Errorf("at least two --line values are required")
	}

	entry := &domain.JournalEntry{
		Date:        addDate,
		Description: addDescription,
	}

	for _, segment := range addLineSegments {
		line, err := parseLine(segment)
		if err != nil {
			return nil, err
		}
		entry.Lines = append(entry.Lines, line)
	}

	return entry, nil
}

func parseLine(value string) (domain.JournalLine, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 3 {
		return domain.JournalLine{}, fmt.Errorf("line %q must follow account,debit,credit", value)
	}

	account := strings.TrimSpace(parts[0])
	debit, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return domain.JournalLine{}, fmt.Errorf("invalid debit in %q: %w", value, err)
	}

	credit, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return domain.JournalLine{}, fmt.Errorf("invalid credit in %q: %w", value, err)
	}

	return domain.JournalLine{
		Account: account,
		Debit:   debit,
		Credit:  credit,
	}, nil
}
