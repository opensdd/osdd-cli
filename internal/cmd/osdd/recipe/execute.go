package recipe

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensdd/osdd-api/clients/go/osdd/recipes"
	"github.com/opensdd/osdd-cli/internal/inputs"
	"github.com/opensdd/osdd-cli/internal/ui"
	"github.com/opensdd/osdd-core/core"
	"github.com/opensdd/osdd-core/core/executable"
	"github.com/opensdd/osdd-core/core/fetcher"
	"github.com/opensdd/osdd-core/core/utils"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

var (
	executeCmd = createExecuteCmd()
)

func createExecuteCmd() *cobra.Command {
	var ideType string
	var recipeFile string
	cmd := &cobra.Command{
		Use:   "execute [recipe-id]",
		Short: "Executes a recipe by its ID. Optionally use --recipe-file to run from a local file (testing/debugging)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ui.PrintLogo()
			fmt.Println()

			var (
				recipe *recipes.ExecutableRecipe
				err    error
			)
			if len(args) == 0 {
				check(fmt.Errorf("missing recipe ID: provide an ID or use --recipe-file"))
				return
			}

			if recipeFile != "" {
				info, statErr := os.Stat(recipeFile)
				if statErr != nil || info.IsDir() {
					check(fmt.Errorf("failed to access recipe file %s: %v", recipeFile, statErr))
					return
				}
				recipe, err = loadRecipeFromFile(recipeFile)
			} else {
				source := args[0]
				gh := &fetcher.GitHub{}
				recipe, err = gh.FetchRecipe(source)
			}
			check(err)
			// TODO: need to make it more flexible.
			utils.AllowedCommands = []string{"git", "devplan"}

			userIn := &inputs.User{}
			ctx := context.Background()
			in, err := userIn.Request(ctx, recipe.GetRecipe())
			check(err)
			execRecipe := executable.ForRecipe(recipe)
			if recipe.GetEntryPoint() == nil {
				recipe.SetEntryPoint(&recipes.EntryPoint{})
			}
			if recipe.GetEntryPoint().GetIdeType() == "" {
				recipe.GetEntryPoint().SetIdeType(ideType)
			}
			genCtx := &core.GenerationContext{
				UserInput:  in,
				ExecRecipe: recipe,
			}
			materialized, err := execRecipe.Materialize(ctx, genCtx)
			check(err)
			basePath := "."
			if wsPath := materialized.GetWorkspacePath(); wsPath != "" {
				basePath = wsPath
			}
			fmt.Printf("Storing into %v\n", basePath)
			check(execRecipe.Execute(ctx, genCtx))
		},
	}
	cmd.Flags().StringVarP(&ideType, "ide", "i", "", "Name of the IDE")
	cmd.Flags().StringVarP(&recipeFile, "recipe-file", "f", "", "Path to a local recipe file (YAML or JSON). When set, the recipe ID argument is ignored.")
	_ = cmd.MarkFlagRequired("ide")
	return cmd
}

func loadRecipeFromFile(path string) (*recipes.ExecutableRecipe, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipe file: %w", err)
	}
	exec := &recipes.ExecutableRecipe{}
	um := protojson.UnmarshalOptions{DiscardUnknown: true}

	// Parse strictly based on file extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		var y any
		if err := yaml.Unmarshal(b, &y); err != nil {
			return nil, fmt.Errorf("failed to parse YAML recipe: %w", err)
		}
		jb, err := json.Marshal(y)
		if err != nil {
			return nil, fmt.Errorf("failed to convert YAML to JSON: %w", err)
		}
		return exec, um.Unmarshal(jb, exec)
	case ".json":
		return exec, um.Unmarshal(b, exec)
	default:
		return nil, fmt.Errorf("unsupported recipe file extension %q: expected .yaml, .yml, or .json", ext)
	}
}
