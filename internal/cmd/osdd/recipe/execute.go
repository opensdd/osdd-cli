package recipe

import (
	"context"
	"fmt"

	"github.com/opensdd/osdd-cli/internal/inputs"
	"github.com/opensdd/osdd-cli/internal/ui"
	"github.com/opensdd/osdd-core/core/fetcher"
	"github.com/spf13/cobra"
)

var (
	executeCmd = createExecuteCmd()
)

func createExecuteCmd() *cobra.Command {
	var recipeID string
	cmd := &cobra.Command{
		Use:   "execute",
		Short: "Executes a recipe by its ID",
		Run: func(_ *cobra.Command, _ []string) {
			ui.PrintLogo()
			fmt.Println()
			gh := &fetcher.GitHub{}
			recipe, err := gh.FetchRecipe(recipeID)
			check(err)
			userIn := &inputs.User{}
			ctx := context.Background()
			in, err := userIn.Request(ctx, recipe.GetRecipe())
			check(err)
			fmt.Println(in)
		},
	}
	cmd.Flags().StringVarP(&recipeID, "id", "i", "", "ID of the recipe to execute")
	_ = cmd.MarkFlagRequired("id")
	return cmd
}
