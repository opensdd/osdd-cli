package main

import (
	"os"

	"github.com/opensdd/osdd-cli/internal/cmd/osdd/recipe"
	"github.com/opensdd/osdd-cli/internal/cmd/osdd/version"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "osdd",
		Short: "OpenSDD CLI for accessing OpenSDD flows",
		Long:  "osdd is a command-line interface for OpenSDD.\nFor more information, visit: https://opensdd.ai",
		Run: func(cmd *cobra.Command, args []string) {
			// Print help if no arguments provided
			check(cmd.Help())
		},
	}

	rootCmd.AddCommand(version.VersionCmd())
	rootCmd.AddCommand(recipe.Cmd)

	// Custom handling for unknown commands
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErr("Error: " + err.Error() + "\n")
		check(cmd.Help())
		os.Exit(1)
		return nil
	})
	check(rootCmd.Execute())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
