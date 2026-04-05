package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"ledgerd/internal/report"
)

var (
	balanceSheetDate string
)

var balanceSheetCmd = &cobra.Command{
	Use:   "balance-sheet",
	Short: "Generate a balance sheet as of a date",
	RunE: func(cmd *cobra.Command, args []string) error {
		dateStr := strings.TrimSpace(balanceSheetDate)
		if dateStr == "" {
			return fmt.Errorf("--date is required")
		}

		asOf, err := time.Parse(reportDateLayout, dateStr)
		if err != nil {
			return fmt.Errorf("invalid date: %w", err)
		}

		entries, err := app.ListEntries()
		if err != nil {
			return err
		}

		sheet, err := report.GenerateBalanceSheet(entries, asOf)
		if err != nil {
			return err
		}

		printBalanceSheet(cmd.OutOrStdout(), sheet)
		return nil
	},
}

func init() {
	balanceSheetCmd.Flags().StringVar(&balanceSheetDate, "date", "", "As-of date (YYYY-MM-DD)")
	rootCmd.AddCommand(balanceSheetCmd)
}

func printBalanceSheet(w io.Writer, sheet report.BalanceSheet) {
	fmt.Fprintf(w, "BALANCE SHEET (%s)\n\n", sheet.Date.Format(reportDateLayout))

	printReportSection(w, "ASSETS", sheet.Assets)
	printReportSection(w, "LIABILITIES", sheet.Liabilities)
	printReportSection(w, "EQUITY", sheet.Equity)

	fmt.Fprintf(w, "TOTAL ASSETS:        %10.2f\n", sheet.TotalAssets)
	fmt.Fprintf(w, "TOTAL L+E:           %10.2f\n", sheet.TotalLiabilities+sheet.TotalEquity)
}
