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
	incomeStart string
	incomeEnd   string
)

var incomeStatementCmd = &cobra.Command{
	Use:   "income-statement",
	Short: "Generate an income statement for a period",
	RunE: func(cmd *cobra.Command, args []string) error {
		startStr := strings.TrimSpace(incomeStart)
		endStr := strings.TrimSpace(incomeEnd)

		if startStr == "" || endStr == "" {
			return fmt.Errorf("--start and --end are required")
		}

		startDate, err := time.Parse(reportDateLayout, startStr)
		if err != nil {
			return fmt.Errorf("invalid start date: %w", err)
		}

		endDate, err := time.Parse(reportDateLayout, endStr)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}

		if startDate.After(endDate) {
			return fmt.Errorf("--start must be on or before --end")
		}

		entries, err := app.ListEntries()
		if err != nil {
			return err
		}

		statement, err := report.GenerateIncomeStatement(entries, startDate, endDate)
		if err != nil {
			return err
		}

		printIncomeStatement(cmd.OutOrStdout(), statement)
		return nil
	},
}

func init() {
	incomeStatementCmd.Flags().StringVar(&incomeStart, "start", "", "Start date inclusive (YYYY-MM-DD)")
	incomeStatementCmd.Flags().StringVar(&incomeEnd, "end", "", "End date inclusive (YYYY-MM-DD)")
	rootCmd.AddCommand(incomeStatementCmd)
}

func printIncomeStatement(w io.Writer, statement report.IncomeStatement) {
	fmt.Fprintf(w, "INCOME STATEMENT (%s ~ %s)\n\n", statement.StartDate.Format(reportDateLayout), statement.EndDate.Format(reportDateLayout))

	printReportSection(w, "INCOME", statement.Income)
	printReportSection(w, "EXPENSES", statement.Expenses)

	fmt.Fprintf(w, "TOTAL INCOME:        %10.2f\n", statement.TotalIncome)
	fmt.Fprintf(w, "TOTAL EXPENSES:      %10.2f\n", statement.TotalExpense)
	fmt.Fprintf(w, "NET INCOME:          %10.2f\n", statement.NetIncome)
}
