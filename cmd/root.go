package cmd

import (
	"github.com/spf13/cobra"
	"ledgerd/internal/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "ledgerd",
		Short: "A minimal double-entry accounting CLI",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if app != nil {
				return nil
			}

			var err error
			app, err = cli.NewApp(dataPath)
			return err
		},
	}

	dataPath string
	app      *cli.App
)

func init() {
	rootCmd.PersistentFlags().StringVar(&dataPath, "data", "data/journal.json", "Path to journal JSON file")
}

// Execute runs the CLI.
func Execute() error {
	return rootCmd.Execute()
}
