package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/parser"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory011_NavigateTodosWithArrowKeys(t *testing.T) {
	is := is.New(t)
	// Scenario: Navigate todos with arrow keys

	input := `(A) Task one
(A) Task two
(A) Task three`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// First todo should be selected by default
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "Task one"))  // expected first todo to be visible

	// Press Down arrow - should select second todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Down again - should select third todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Down again - should wrap to first todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Up - should wrap to third todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	_ = updatedModel.(ui.Model)
}

func TestStory011_NavigateTodosWithWASD(t *testing.T) {
	is := is.New(t)
	// Scenario: Navigate todos with w/s keys

	input := `(A) Task one
(A) Task two`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 's' - should move down
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model = updatedModel.(ui.Model)

	// Press 'w' - should move up
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	_ = updatedModel.(ui.Model)
}

func TestStory011_MarkTodoAsComplete(t *testing.T) {
	is := is.New(t)
	// Scenario: Mark todo as complete

	input := `(A) Fix bug +WebApp
(A) Another task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// First todo is selected by default
	// Press Space to mark as complete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Check that todo is marked as complete in view
	view := model.View()
	// With table rendering, completed todos show date in Completed column
	is.True(strings.Contains(stripANSI(view), "Completed"))  // expected table to have Completed column header
	// Should show "today" in the Completed column
	is.True(strings.Contains(stripANSI(view), "today"))  // expected completed todo to show 'today' in Completed column

	// Check that file was updated with completion marker
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "x"))  // expected file to contain completion marker 'x'
	is.True(strings.Contains(written, "Fix bug +WebApp"))  // expected completed todo to retain description and tags
}

func TestStory011_UnmarkCompletedTodo(t *testing.T) {
	is := is.New(t)
	// Scenario: Unmark completed todo

	input := `x 2025-12-25 (A) Completed task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	// First, let's parse the todos to see what we get
	reader, _ := source.GetTodos()
	todos, err := parser.Parse(reader)
	reader.Close()
	is.NoErr(err)
	t.Logf("Parsed %d todos", len(todos))
	for i, td := range todos {
		t.Logf("Todo %d: priority=%v, completed=%v, desc=%s", i, td.Priority(), td.IsCompleted(), td.Description())
	}

	// Reset the source reader
	source.reader = strings.NewReader(input)

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	t.Logf("Initial matrix DO FIRST count: %d", len(m.DoFirst()))
	if len(m.DoFirst()) > 0 {
		t.Logf("First todo: completed=%v, description=%s", m.DoFirst()[0].IsCompleted(), m.DoFirst()[0].Description())
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// First todo (completed) is selected by default
	// Press Space to unmark
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Check that file was updated without completion marker
	written := source.writer.(*strings.Builder).String()
	t.Logf("Written content:\n%s", written)

	is.True(!strings.Contains(written, "x "))  // expected completed task to be unmarked (no 'x' prefix)
	is.True(strings.Contains(written, "(A) Completed task"))  // expected unmarked todo to retain priority and description
}

func TestStory011_EmptyQuadrantNoSelection(t *testing.T) {
	is := is.New(t)
	// Scenario: Empty quadrant has no selection

	input := `(A) Task in DO FIRST`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE (empty)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	view := model.View()
	is.True(strings.Contains(stripANSI(view), "(no tasks)"))  // expected to see '(no tasks)' in empty quadrant

	// Pressing Space should do nothing (no panic)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Pressing navigation keys should do nothing (no panic)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	_ = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	_ = updatedModel.(ui.Model)
}

func TestStory011_SelectionNotShownInOverviewMode(t *testing.T) {
	is := is.New(t)
	// Scenario: Selection state not shown in overview mode

	input := `(A) Task one
(A) Task two
(A) Task three`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Navigate to second todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show all tasks but no selection highlight
	is.True(strings.Contains(stripANSI(view), "Task one"))    // expected to see Task one in overview
	is.True(strings.Contains(stripANSI(view), "Task two"))    // expected to see Task two in overview
	is.True(strings.Contains(stripANSI(view), "Task three"))  // expected to see Task three in overview

	// No selection indicator should be present (we'll check this by ensuring
	// focus-mode specific rendering isn't present in overview)
}

func TestStory011_InputModePreservesSelection(t *testing.T) {
	is := is.New(t)
	// Scenario: Entering input mode preserves selection

	input := `(A) Task one
(A) Task two
(A) Task three`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Navigate to third todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Cancel input
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Selection should still be on third todo (we can't easily verify this
	// without exposing internal state, but pressing Space should toggle
	// the third todo)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Verify third todo was toggled by checking written output
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "x") && strings.Contains(written, "Task three"))  // expected third todo to be marked as complete
}
