package inputs

import (
	"context"
	"testing"

	osdd "github.com/opensdd/osdd-api/clients/go/osdd"
	"github.com/opensdd/osdd-api/clients/go/osdd/recipes"
)

// buildTestRecipe constructs a recipe with:
//   - entry1: simple user_input with params a(required), b(optional)
//   - entry2: combined with a user_input containing param a(required)
func buildTestRecipe() *recipes.Recipe {
	// entry 1: simple user_input
	pA := (&osdd.UserInputParameter_builder{Name: "a", Optional: false}).Build()
	pB := (&osdd.UserInputParameter_builder{Name: "b", Optional: true}).Build()
	ui1 := (&recipes.UserInputContextSource_builder{Entries: []*osdd.UserInputParameter{pA, pB}}).Build()
	from1 := (&recipes.ContextFrom_builder{UserInput: ui1}).Build()
	entry1 := (&recipes.ContextEntry_builder{Path: "path1", From: from1}).Build()

	// entry 2: combined with a user_input for "a"
	pA2 := (&osdd.UserInputParameter_builder{Name: "a", Optional: false}).Build()
	ui2 := (&recipes.UserInputContextSource_builder{Entries: []*osdd.UserInputParameter{pA2}}).Build()
	itemUI := (&recipes.CombinedContextSource_Item_builder{UserInput: ui2}).Build()
	combined := (&recipes.CombinedContextSource_builder{Items: []*recipes.CombinedContextSource_Item{itemUI}}).Build()
	from2 := (&recipes.ContextFrom_builder{Combined: combined}).Build()
	entry2 := (&recipes.ContextEntry_builder{Path: "path2", From: from2}).Build()

	ctx := (&recipes.Context_builder{Entries: []*recipes.ContextEntry{entry1, entry2}}).Build()
	rec := (&recipes.Recipe_builder{}).Build()
	rec.SetContext(ctx)
	return rec
}

func TestUserRequest_DoesNotMutateRecipeAndReturnsAnswers(t *testing.T) {
	recipe := buildTestRecipe()

	// Pre-conditions: ensure inputs are present
	entries := recipe.GetContext().GetEntries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].GetFrom().HasUserInput() || entries[0].GetFrom().HasText() {
		t.Fatalf("entry1 should initially have user_input and not text")
	}
	items := entries[1].GetFrom().GetCombined().GetItems()
	if len(items) != 1 || !items[0].HasUserInput() || items[0].HasText() {
		t.Fatalf("combined entry should initially contain user_input item only")
	}

	// Prepare deterministic answers: a -> "va", b -> "" (optional skip), a -> "vb" (last-write-wins)
	answersSeq := []string{"va", "", "vb"}
	idx := 0
	ask := func(_ string, validate func(string) error) (string, error) {
		if idx >= len(answersSeq) {
			t.Fatalf("ask called more times than expected")
		}
		val := answersSeq[idx]
		idx++
		if err := validate(val); err != nil {
			return "", err
		}
		return val, nil
	}

	u := &User{askFn: ask}
	m, err := u.Request(context.Background(), recipe)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Map expectations
	if got, ok := m["a"]; !ok || got != "vb" {
		t.Fatalf("expected a=vb in map (last write wins), got %v, present=%v", got, ok)
	}
	if _, ok := m["b"]; ok {
		t.Fatalf("expected b to be omitted due to empty optional, but it is present: %v", m["b"])
	}

	// Post-conditions: recipe must be unchanged
	entries = recipe.GetContext().GetEntries()
	if !entries[0].GetFrom().HasUserInput() || entries[0].GetFrom().HasText() {
		t.Fatalf("entry1 mutated: user_input removed or text set")
	}
	items = entries[1].GetFrom().GetCombined().GetItems()
	if !items[0].HasUserInput() || items[0].HasText() {
		t.Fatalf("combined item mutated: user_input removed or text set")
	}
}

func TestUserRequest_RespectsContextCancellation(t *testing.T) {
	recipe := buildTestRecipe()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	u := &User{askFn: func(_ string, _ func(string) error) (string, error) {
		t.Fatalf("askFn should not be called after context cancellation")
		return "", nil
	}}

	m, err := u.Request(ctx, recipe)
	if err == nil {
		t.Fatalf("expected an error due to context cancellation, got nil")
	}
	if len(m) != 0 {
		t.Fatalf("expected empty map on immediate cancellation, got %v", m)
	}

	// Ensure recipe unchanged
	entries := recipe.GetContext().GetEntries()
	if !entries[0].GetFrom().HasUserInput() || entries[0].GetFrom().HasText() {
		t.Fatalf("entry1 mutated after cancellation")
	}
}
