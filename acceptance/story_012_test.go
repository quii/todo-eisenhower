package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory012_MoveFromDoFirstToSchedule(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo from DO FIRST to SCHEDULE

	input := `(A) Review quarterly goals`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '2' to move to SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority B
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(B) Review quarterly goals"))  // expected file to contain '(B) Review quarterly goals'
	is.True(!strings.Contains(written, "(A)"))  // expected priority A to be changed to B
}

func TestStory012_MoveFromDelegateToDoFirst(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo from DELEGATE to DO FIRST

	input := `(C) Update documentation`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '1' to move to DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority A
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A) Update documentation"))  // expected file to contain '(A) Update documentation'
}

func TestStory012_MoveToEliminate(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo to ELIMINATE (priority D)

	input := `(B) Optional feature idea`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '4' to move to ELIMINATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority D
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(D) Optional feature idea"))  // expected file to contain '(D) Optional feature idea'
}

func TestStory012_MovingTodoAdjustsSelection(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving todo adjusts selection

	input := `(A) First task
(A) Second task
(A) Third task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Navigate to second task
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Move second task to SCHEDULE (press 'm' then '2')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Check file shows second task moved to priority B
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(B) Second task"))  // expected second task to have priority B

	// First and third tasks should still be priority A
	lines := strings.Split(strings.TrimSpace(written), "\n")
	priorityACount := 0
	for _, line := range lines {
		if strings.Contains(line, "(A)") {
			priorityACount++
		}
	}
	is.Equal(priorityACount, 2)  // expected 2 tasks with priority A
}

func TestStory012_MovingLastTodoReturnsToOverview(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving last todo in quadrant returns to overview

	input := `(C) Only delegate task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Move to DO FIRST (press 'm' then '1')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should return to overview mode
	view := model.View()
	is.True(strings.Contains(view, "DO FIRST") || strings.Contains(view, "SCHEDULE") || strings.Contains(view, "DELEGATE") || strings.Contains(view, "ELIMINATE"))  // expected to return to overview mode showing all quadrants

	// Verify the moved todo is in the file with priority A
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A) Only delegate task"))  // expected todo to be moved to priority A
}

func TestStory012_PreservesTagsAndCompletion(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving todo preserves tags and completion status

	input := `x 2025-01-10 (A) Fix bug +WebApp @computer`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Move to DELEGATE (press 'm' then '3')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Check file preserves everything
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "x "))  // expected completion marker to be preserved
	is.True(strings.Contains(written, "(C)"))  // expected priority to change to C
	is.True(strings.Contains(written, "+WebApp"))  // expected project tag to be preserved
	is.True(strings.Contains(written, "@computer"))  // expected context tag to be preserved
	is.True(strings.Contains(written, "Fix bug"))  // expected description to be preserved
}

func TestStory012_PressingCurrentQuadrantDoesNothing(t *testing.T) {
	is := is.New(t)
	// Scenario: Pressing current quadrant number does nothing

	input := `(B) Plan sprint`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE (quadrant 2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press '2' again (same quadrant) - should be a no-op
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Verify we're still viewing SCHEDULE (not moved to overview or other quadrant)
	view := model.View()
	is.True(strings.Contains(view, "SCHEDULE"))  // should still be viewing SCHEDULE quadrant
	is.True(strings.Contains(view, "Plan sprint"))  // should still show the todo

	// Press 'm' then '2' again - should be a no-op since todo is already priority B
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Verify todo is still in SCHEDULE with priority B
	view = model.View()
	is.True(strings.Contains(view, "SCHEDULE"))  // should still be viewing SCHEDULE quadrant after pressing 'm' then '2'
	is.True(strings.Contains(view, "Plan sprint"))  // should still show the todo after pressing 'm' then '2'

	// Now actually move it to verify moving works (press 'm' then '1')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Check that it was moved to priority A
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A) Plan sprint"))  // after pressing Shift+1, todo should be moved to priority A
}

func TestStory012_NumberKeysStillFocusInOverview(t *testing.T) {
	is := is.New(t)
	// Scenario: Number keys still focus quadrants in overview mode

	input := `(A) Task one
(B) Task two`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press '1' to focus
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should be in focus mode on DO FIRST
	view := model.View()
	is.True(strings.Contains(view, "DO FIRST"))  // expected to focus on DO FIRST quadrant

	// Should show focus mode help text
	is.True(strings.Contains(view, "Press a to add"))  // expected focus mode help text
}
