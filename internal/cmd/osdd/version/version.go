package version

import (
	"fmt"

	"github.com/opensdd/osdd-cli/internal/ui"
	"github.com/spf13/cobra"
)

// Version is set via ldflags at build time
// Usage: go build -ldflags "-X 'github.com/opensdd/osdd-cli/internal/version.Version=v1.0.0'"
var Version string

// GetVersion returns the injected version or "dev" if unset
func GetVersion() string {
	if Version == "" {
		return "dev"
	}
	return Version
}

// IsSet returns true if Version was injected at build time
func IsSet() bool {
	return Version != ""
}

// VersionCmd returns the Cobra command for the version subcommand.
//
// The printFunc parameter is called with the version string to display output.
// The function never returns an error - version display always succeeds.
func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the version of osdd",
		Long:  "Display the version of osdd with ASCII art",
		Run: func(cmd *cobra.Command, args []string) {
			ui.PrintLogo()
			version := GetVersion()
			// Print version
			fmt.Printf("\nOpenSDD CLI version %s\n", version)

			if version == "dev" {
				fmt.Println("WARNING: Version not set at build time")
			}
		},
	}
}
