package inputs

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/opensdd/osdd-api/clients/go/osdd/recipes"
)

// askDefault uses promptui to ask the question with validation.
func askDefault(label string, validate func(string) error) (string, error) {
	p := promptui.Prompt{Label: label, Validate: validate}
	return p.Run()
}

type User struct {
	// askFn is a dependency-injected function used to ask the user for input.
	// It should return the entered value or an error. When nil, a default
	// implementation backed by promptui is used.
	askFn func(label string, validate func(string) error) (string, error)
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

	ask := u.askFn
	if ask == nil {
		ask = askDefault
	}

	answers := make(map[string]string)
	for _, p := range src.GetEntries() {
		if p == nil {
			continue
		}

		name := p.GetName()
		desc := p.GetDescription()
		optional := p.GetOptional()

		label := name
		if desc != "" {
			label = fmt.Sprintf("%s (%s)", name, desc)
		}

		validate := func(input string) error {
			if !optional && strings.TrimSpace(input) == "" {
				return errors.New("required")
			}
			return nil
		}

		// respect cancellation if provided
		select {
		case <-ctx.Done():
			return answers, ctx.Err()
		default:
		}

		result, err := ask(label, validate)
		if err != nil {
			return answers, err
		}
		if strings.TrimSpace(result) == "" {
			// optional and skipped
			continue
		}
		answers[name] = result
	}

	return answers, nil
}
