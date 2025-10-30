package inputs

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/opensdd/osdd-api/clients/go/osdd/recipes"
)

type User struct {
}

// Request requests user input for user inputs configured in the recipe's context.
// Each entry is requested one at a time. Required input cannot be skipped while optional
// can be skipped.
func (u *User) Request(ctx context.Context, recipe *recipes.Recipe) (map[string]string, error) {
	collected := make(map[string]string)
	if recipe == nil || !recipe.HasContext() {
		return collected, nil
	}

	ctxObj := recipe.GetContext()
	entries := ctxObj.GetEntries()
	if len(entries) == 0 {
		return collected, nil
	}

	for _, entry := range entries {
		if entry == nil {
			continue
		}
		from := entry.GetFrom()
		if from == nil {
			continue
		}

		// Handle simple user input source (collect only, do not mutate recipe)
		if ui := from.GetUserInput(); ui != nil {
			answers, err := u.promptForUserInput(ctx, ui)
			for k, v := range answers {
				collected[k] = v
			}
			if err != nil {
				return collected, err
			}
			continue
		}

		// Handle combined source that may include user input among items (collect only)
		if combined := from.GetCombined(); combined != nil {
			items := combined.GetItems()
			for _, it := range items {
				if it == nil {
					continue
				}
				uisrc := it.GetUserInput()
				if uisrc == nil {
					continue
				}
				answers, err := u.promptForUserInput(ctx, uisrc)
				for k, v := range answers {
					collected[k] = v
				}
				if err != nil {
					return collected, err
				}
			}
		}
	}

	return collected, nil
}

// promptForUserInput iterates over the user input parameters and prompts the user
// using an injected prompter. Required inputs cannot be skipped; optional inputs can be skipped
// by submitting an empty value. Returns a map[string]string with collected answers keyed by parameter name.
func (u *User) promptForUserInput(ctx context.Context, src *recipes.UserInputContextSource) (map[string]string, error) {
	if src == nil {
		return map[string]string{}, nil
	}

	var inputs []huh.Field
	results := make(map[string]*string)

	answers := make(map[string]string)
	for _, p := range src.GetEntries() {
		if p == nil {
			continue
		}

		validate := func(input string) error {
			if p.GetOptional() || strings.TrimSpace(input) != "" {
				return nil
			}
			return errors.New("required")
		}

		v := ""
		in := huh.NewInput().
			Title(p.GetName()).
			Description(p.GetDescription()).
			Prompt("?").
			Validate(validate).
			Value(&v)
		results[p.GetName()] = &v
		inputs = append(inputs, in)
	}

	err := huh.NewForm(huh.NewGroup(inputs...)).RunWithContext(ctx)
	if err != nil {
		return nil, err
	}
	for k, v := range results {
		val := *v
		if val != "" {
			answers[k] = val
		}
	}

	return answers, nil
}
