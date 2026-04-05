package cmd

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recorded journal entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := app.ListEntries()
		if err != nil {
			return err
		}

		encoder := json.NewEncoder(cmd.OutOrStdout())
		encoder.SetIndent("", "  ")
		return encoder.Encode(entries)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

