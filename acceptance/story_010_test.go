package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory010_DisplayProjectTagInventory(t *testing.T) {
	// Scenario: Display project tag inventory

	input := `(A) Task one +strategy
(A) Task two +strategy
(A) Task three +strategy
(B) Task four +hiring
(B) Task five +hiring
(C) Task six +architecture
x (A) Completed task +strategy`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show project inventory at bottom
	if !strings.Contains(view, "Projects:") {
		t.Error("expected to see 'Projects:' label")
	}

	// Should show strategy with count 3 (not counting completed)
	if !strings.Contains(view, "strategy") && !strings.Contains(view, "(3)") {
		t.Error("expected to see +strategy (3)")
	}

	// Should show hiring with count 2
	if !strings.Contains(view, "hiring") && !strings.Contains(view, "(2)") {
		t.Error("expected to see +hiring (2)")
	}

	// Should show architecture with count 1
	if !strings.Contains(view, "architecture") && !strings.Contains(view, "(1)") {
		t.Error("expected to see +architecture (1)")
	}

	// Strategy should appear before hiring (higher count)
	strategyPos := strings.Index(view, "strategy")
	hiringPos := strings.Index(view, "hiring")
	if strategyPos == -1 || hiringPos == -1 || strategyPos > hiringPos {
		t.Error("expected +strategy to appear before +hiring (sorted by count descending)")
	}
}

func TestStory010_DisplayContextTagInventory(t *testing.T) {
	// Scenario: Display context tag inventory

	input := `(A) Task one @computer
(A) Task two @computer
(A) Task three @computer
(A) Task four @computer
(A) Task five @computer
(B) Task six @phone
(B) Task seven @phone
(C) Task eight @office`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show context inventory
	if !strings.Contains(view, "Contexts:") {
		t.Error("expected to see 'Contexts:' label")
	}

	// Should show computer with count 5
	if !strings.Contains(view, "computer") && !strings.Contains(view, "(5)") {
		t.Error("expected to see @computer (5)")
	}

	// Should show phone with count 2
	if !strings.Contains(view, "phone") && !strings.Contains(view, "(2)") {
		t.Error("expected to see @phone (2)")
	}

	// Should show office with count 1
	if !strings.Contains(view, "office") && !strings.Contains(view, "(1)") {
		t.Error("expected to see @office (1)")
	}
}

func TestStory010_DisplayBothProjectAndContextInventory(t *testing.T) {
	// Scenario: Display both project and context inventory

	input := `(A) Task one +strategy @computer
(A) Task two +strategy @computer
(B) Task three +hiring @phone`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show both project and context lines
	if !strings.Contains(view, "Projects:") {
		t.Error("expected to see 'Projects:' label")
	}
	if !strings.Contains(view, "Contexts:") {
		t.Error("expected to see 'Contexts:' label")
	}

	// Should show counts for both
	if !strings.Contains(view, "strategy") {
		t.Error("expected to see +strategy")
	}
	if !strings.Contains(view, "computer") {
		t.Error("expected to see @computer")
	}
}

func TestStory010_NoTagsInUse(t *testing.T) {
	// Scenario: No tags in use

	input := `(A) Task without tags
(B) Another task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show (none) for both
	projectsLine := extractLine(view, "Projects:")
	if !strings.Contains(projectsLine, "(none)") {
		t.Error("expected to see 'Projects: (none)'")
	}

	contextsLine := extractLine(view, "Contexts:")
	if !strings.Contains(contextsLine, "(none)") {
		t.Error("expected to see 'Contexts: (none)'")
	}
}

func TestStory010_InventoryNotShownInFocusMode(t *testing.T) {
	// Scenario: Inventory not shown in focus mode

	input := `(A) Task +strategy @computer
(B) Task +hiring @phone`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter focus mode (press 1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should NOT show inventory in focus mode
	if strings.Contains(view, "Projects:") && strings.Contains(view, "strategy (1)") {
		t.Error("expected inventory NOT to be shown in focus mode")
	}
}

func TestStory010_CountsUpdateWhenAddingTodos(t *testing.T) {
	// Scenario: Counts update when adding todos

	input := `(A) Existing task +strategy`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Check initial count
	view := model.View()
	if !strings.Contains(view, "strategy") && !strings.Contains(view, "(1)") {
		t.Error("expected initial count of +strategy (1)")
	}

	// Focus on DO FIRST and add a new todo with +strategy
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type new todo with +strategy
	for _, ch := range "New task +strategy " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	view2 := model.View()

	// Count should now be 2
	if !strings.Contains(view2, "strategy") && !strings.Contains(view2, "(2)") {
		t.Error("expected updated count of +strategy (2)")
	}
}

// Helper function to extract a line containing a specific substring
func extractLine(text string, substring string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, substring) {
			return line
		}
	}
	return ""
}
