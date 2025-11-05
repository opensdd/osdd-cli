package recipe

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Cmd = create()
)

func create() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "recipe",
		Short:   "Commands handling OpenSDD recipes",
		Aliases: []string{"specs"},
	}
	cmd.AddCommand(executeCmd)
	return cmd
}

func check(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
