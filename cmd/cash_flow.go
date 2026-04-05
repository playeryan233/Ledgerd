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
	cashFlowStart string
	cashFlowEnd   string
)

var cashFlowCmd = &cobra.Command{
	Use:   "cash-flow",
	Short: "Generate a cash flow statement for a period",
	RunE: func(cmd *cobra.Command, args []string) error {
		startStr := strings.TrimSpace(cashFlowStart)
		endStr := strings.TrimSpace(cashFlowEnd)
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

		statement, err := report.GenerateCashFlowStatement(entries, startDate, endDate)
		if err != nil {
			return err
		}

		printCashFlowStatement(cmd.OutOrStdout(), statement)
		return nil
	},
}

func init() {
	cashFlowCmd.Flags().StringVar(&cashFlowStart, "start", "", "Start date inclusive (YYYY-MM-DD)")
	cashFlowCmd.Flags().StringVar(&cashFlowEnd, "end", "", "End date inclusive (YYYY-MM-DD)")
	rootCmd.AddCommand(cashFlowCmd)
}

func printCashFlowStatement(w io.Writer, statement report.CashFlowStatement) {
	fmt.Fprintf(w, "CASH FLOW STATEMENT (%s ~ %s)\n\n", statement.StartDate.Format(reportDateLayout), statement.EndDate.Format(reportDateLayout))

	printSignedSection(w, "OPERATING", statement.Operating)
	printSignedSection(w, "INVESTING", statement.Investing)
	printSignedSection(w, "FINANCING", statement.Financing)

	printSignedTotal(w, "NET OPERATING:", statement.NetOperating)
	printSignedTotal(w, "NET INVESTING:", statement.NetInvesting)
	printSignedTotal(w, "NET FINANCING:", statement.NetFinancing)
	printSignedTotal(w, "NET INCREASE:", statement.NetIncrease)
	printUnsignedTotal(w, "BEGINNING CASH:", statement.BeginningCash)
	printUnsignedTotal(w, "ENDING CASH:", statement.EndingCash)
}
