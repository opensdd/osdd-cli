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
	var launchInstructionsFile string
	var idePaths []string
	cmd := &cobra.Command{
		Use:   "execute [recipe-id]",
		Short: "Executes a recipe by its ID. Optionally use --recipe-file to run from a local file (testing/debugging)",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			ui.PrintLogo()
			fmt.Println()

			var recipe *recipes.ExecutableRecipe
			var err error
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
			idePathsMap := map[string]string{}
			for _, p := range idePaths {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) != 2 {
					check(fmt.Errorf("invalid ide path %v: expected format -p ide_type=<path>", p))
				}
				idePathsMap[kv[0]] = kv[1]
			}
			genCtx := &core.GenerationContext{
				UserInput:     in,
				ExecRecipe:    recipe,
				IDEPaths:      idePathsMap,
				OutputCMDOnly: launchInstructionsFile != "",
			}
			materialized, err := execRecipe.Materialize(ctx, genCtx)
			check(err)
			basePath := ""
			if wsPath := materialized.GetWorkspacePath(); wsPath != "" {
				basePath = wsPath
			} else {
				basePath, err = filepath.Abs(".")
				check(err)
			}

			fmt.Printf("Storing into %v\n", basePath)
			res, err := execRecipe.Execute(ctx, genCtx)
			check(err)
			check(writeLaunchOutput(launchInstructionsFile, res))
		},
	}
	cmd.Flags().StringVarP(&ideType, "ide", "i", "", "Name of the IDE")
	cmd.Flags().StringVarP(&recipeFile, "recipe-file", "f", "", "Path to a local recipe file (YAML or JSON). When set, the recipe ID argument is ignored.")
	cmd.Flags().StringVarP(&launchInstructionsFile, "launch-instructions", "l", "", "Path to a file where launch instructions should be written to. If this parameter is provided, then the CLI will not execute the recipe, but will output the launch details instead into that file.")
	cmd.Flags().StringSliceVarP(&idePaths, "ide-paths", "p", nil, "Paths to ide executables in format `-p codex=<codex_path> -p claude=<claude_path>`")
	_ = cmd.MarkFlagRequired("ide")
	return cmd
}
func writeLaunchOutput(outFile string, res executable.RecipeExecutionResult) error {
	if outFile == "" || res.LaunchResult.LaunchDetails == nil {
		return nil
	}
	m := protojson.MarshalOptions{Indent: "  "}
	b, err := m.Marshal(res.LaunchResult.LaunchDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal launch details: %w", err)
	}
	return os.WriteFile(outFile, b, 0644)
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
