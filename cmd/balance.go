package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:   "balance ACCOUNT",
	Short: "Show the current balance for an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		account := args[0]
		balance, err := app.ComputeBalance(account)
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s balance: %.2f\n", account, balance)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}

