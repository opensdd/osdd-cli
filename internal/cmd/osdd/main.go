package main

import (
	"os"

	"github.com/opensdd/osdd-cli/internal/ui"
	"github.com/opensdd/osdd-cli/internal/version"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "osdd",
		Short: "OpenSDD CLI for accessing OpenSDD flows",
		Long:  "osdd is a command-line interface for OpenSDD.\nFor more information, visit: https://opensdd.ai",
		Run: func(cmd *cobra.Command, args []string) {
			// Print help if no arguments provided
			cmd.Help()
		},
	}

	// Add version subcommand
	rootCmd.AddCommand(version.VersionCmd(ui.PrintVersion))

	// Custom handling for unknown commands
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErr("Error: " + err.Error() + "\n")
		cmd.Help()
		os.Exit(1)
		return nil
	})

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
