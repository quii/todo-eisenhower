package acceptance_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory010_DisplayProjectTagInventory(t *testing.T) {
	// Scenario: Display project tag inventory
	is := is.New(t)

	input := `(A) Task one +strategy
(A) Task two +strategy
(A) Task three +strategy
(B) Task four +hiring
(B) Task five +hiring
(C) Task six +architecture
x (A) Completed task +strategy`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show project inventory at bottom
	is.True(strings.Contains(stripANSI(view), "Projects (+):")) // expected to see 'Projects:' label

	// Should show strategy with count 3 (not counting completed)
	is.True(strings.Contains(stripANSI(view), "strategy") && strings.Contains(stripANSI(view), "(3)")) // expected to see +strategy (3)

	// Should show hiring with count 2
	is.True(strings.Contains(stripANSI(view), "hiring") && strings.Contains(stripANSI(view), "(2)")) // expected to see +hiring (2)

	// Should show architecture with count 1
	is.True(strings.Contains(stripANSI(view), "architecture") && strings.Contains(stripANSI(view), "(1)")) // expected to see +architecture (1)

	// Strategy should appear before hiring (higher count)
	strategyPos := strings.Index(view, "strategy")
	hiringPos := strings.Index(view, "hiring")
	is.True(strategyPos != -1 && hiringPos != -1 && strategyPos < hiringPos) // expected +strategy to appear before +hiring (sorted by count descending)
}

func TestStory010_DisplayContextTagInventory(t *testing.T) {
	// Scenario: Display context tag inventory
	is := is.New(t)

	input := `(A) Task one @computer
(A) Task two @computer
(A) Task three @computer
(A) Task four @computer
(A) Task five @computer
(B) Task six @phone
(B) Task seven @phone
(C) Task eight @office`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show context inventory
	is.True(strings.Contains(stripANSI(view), "Contexts (@):")) // expected to see 'Contexts:' label

	// Should show computer with count 5
	is.True(strings.Contains(stripANSI(view), "computer") && strings.Contains(stripANSI(view), "(5)")) // expected to see @computer (5)

	// Should show phone with count 2
	is.True(strings.Contains(stripANSI(view), "phone") && strings.Contains(stripANSI(view), "(2)")) // expected to see @phone (2)

	// Should show office with count 1
	is.True(strings.Contains(stripANSI(view), "office") && strings.Contains(stripANSI(view), "(1)")) // expected to see @office (1)
}

func TestStory010_DisplayBothProjectAndContextInventory(t *testing.T) {
	// Scenario: Display both project and context inventory
	is := is.New(t)

	input := `(A) Task one +strategy @computer
(A) Task two +strategy @computer
(B) Task three +hiring @phone`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show both project and context lines
	is.True(strings.Contains(stripANSI(view), "Projects (+):")) // expected to see 'Projects:' label
	is.True(strings.Contains(stripANSI(view), "Contexts (@):")) // expected to see 'Contexts:' label

	// Should show counts for both
	is.True(strings.Contains(stripANSI(view), "strategy")) // expected to see +strategy
	is.True(strings.Contains(stripANSI(view), "computer")) // expected to see @computer
}

func TestStory010_NoTagsInUse(t *testing.T) {
	// Scenario: No tags in use
	is := is.New(t)

	input := `(A) Task without tags
(B) Another task`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show (none) for both
	projectsLine := extractLine(view, "Projects (+):")
	is.True(strings.Contains(projectsLine, "(none)")) // expected to see 'Projects: (none)'

	contextsLine := extractLine(view, "Contexts (@):")
	is.True(strings.Contains(contextsLine, "(none)")) // expected to see 'Contexts: (none)'
}

func TestStory010_InventoryNotShownInFocusMode(t *testing.T) {
	// Scenario: Inventory not shown in focus mode
	is := is.New(t)

	input := `(A) Task +strategy @computer
(B) Task +hiring @phone`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter focus mode (press 1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should NOT show inventory in focus mode
	is.True(!(strings.Contains(stripANSI(view), "Projects (+):") && strings.Contains(stripANSI(view), "strategy (1)"))) // expected inventory NOT to be shown in focus mode
}

func TestStory010_CountsUpdateWhenAddingTodos(t *testing.T) {
	// Scenario: Counts update when adding todos
	is := is.New(t)

	input := `(A) Existing task +strategy`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Check initial count
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "strategy") && strings.Contains(stripANSI(view), "(1)")) // expected initial count of +strategy (1)

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
	is.True(strings.Contains(view2, "strategy") && strings.Contains(view2, "(2)")) // expected updated count of +strategy (2)

	// Verify the new todo was persisted to repository
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 2) // expected 2 todos (existing + new)

	newTodo := savedTodos[1]
	is.True(strings.Contains(newTodo.Description(), "New task")) // description should contain base text
	is.Equal(newTodo.Priority(), todo.PriorityA)
	is.Equal(len(newTodo.Projects()), 1)
	is.Equal(newTodo.Projects()[0], "strategy")
}

// Helper function to extract a line containing a specific substring
func extractLine(text, substring string) string {
	for line := range strings.SplitSeq(text, "\n") {
		if strings.Contains(line, substring) {
			return line
		}
	}
	return ""
}
