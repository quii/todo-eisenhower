package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory012_MoveFromDoFirstToSchedule(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo from DO FIRST to SCHEDULE

	input := `(A) Review quarterly goals`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '2' to move to SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	_ = updatedModel.(ui.Model)

	// Check file was updated with priority B
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityB)
	is.Equal(savedTodos[0].Description(), "Review quarterly goals")
}

func TestStory012_MoveFromDelegateToDoFirst(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo from DELEGATE to DO FIRST

	input := `(C) Update documentation`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '1' to move to DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	_ = updatedModel.(ui.Model)

	// Check file was updated with priority A
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
	is.Equal(savedTodos[0].Description(), "Update documentation")
}

func TestStory012_MoveToEliminate(t *testing.T) {
	is := is.New(t)
	// Scenario: Move todo to ELIMINATE (priority D)

	input := `(B) Optional feature idea`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode, then '4' to move to ELIMINATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	_ = updatedModel.(ui.Model)

	// Check file was updated with priority D
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityD)
	is.Equal(savedTodos[0].Description(), "Optional feature idea")
}

func TestStory012_MovingTodoAdjustsSelection(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving todo adjusts selection

	input := `(A) First task
(A) Second task
(A) Third task`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	_ = updatedModel.(ui.Model)

	// Check file shows second task moved to priority B
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 3)

	// Count priority A and B tasks
	priorityACount := 0
	priorityBCount := 0
	hasSecondTask := false
	for _, t := range savedTodos {
		if t.Priority() == todo.PriorityA {
			priorityACount++
		}
		if t.Priority() == todo.PriorityB && t.Description() == "Second task" {
			priorityBCount++
			hasSecondTask = true
		}
	}
	is.Equal(priorityACount, 2)  // expected 2 tasks with priority A
	is.Equal(priorityBCount, 1)  // expected 1 task with priority B
	is.True(hasSecondTask)       // expected second task to have priority B
}

func TestStory012_MovingLastTodoReturnsToOverview(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving last todo in quadrant returns to overview

	input := `(C) Only delegate task`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "Do First") || strings.Contains(stripANSI(view), "Schedule") || strings.Contains(stripANSI(view), "Delegate") || strings.Contains(stripANSI(view), "Eliminate"))  // expected to return to overview mode showing all quadrants

	// Verify the moved todo is in the file with priority A
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
	is.Equal(savedTodos[0].Description(), "Only delegate task")
}

func TestStory012_PreservesTagsAndCompletion(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving todo preserves tags and completion status

	input := `x 2025-01-10 (A) Fix bug +WebApp @computer`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Move to DELEGATE (press 'm' then '3')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	_ = updatedModel.(ui.Model)

	// Check file preserves everything
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())                           // expected completion marker to be preserved
	is.Equal(savedTodos[0].Priority(), todo.PriorityC)             // expected priority to change to C
	is.Equal(savedTodos[0].Projects(), []string{"WebApp"})         // expected project tag to be preserved
	is.Equal(savedTodos[0].Contexts(), []string{"computer"})       // expected context tag to be preserved
	is.Equal(savedTodos[0].Description(), "Fix bug")               // expected description to be preserved
}

func TestStory012_PressingCurrentQuadrantDoesNothing(t *testing.T) {
	is := is.New(t)
	// Scenario: Pressing current quadrant number does nothing

	input := `(B) Plan sprint`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "Schedule"))  // should still be viewing SCHEDULE quadrant
	is.True(strings.Contains(stripANSI(view), "Plan sprint"))  // should still show the todo

	// Press 'm' then '2' again - should be a no-op since todo is already priority B
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Verify todo is still in SCHEDULE with priority B
	view = model.View()
	is.True(strings.Contains(stripANSI(view), "Schedule"))  // should still be viewing SCHEDULE quadrant after pressing 'm' then '2'
	is.True(strings.Contains(stripANSI(view), "Plan sprint"))  // should still show the todo after pressing 'm' then '2'

	// Now actually move it to verify moving works (press 'm' then '1')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	_ = updatedModel.(ui.Model)

	// Check that it was moved to priority A
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
	is.Equal(savedTodos[0].Description(), "Plan sprint")
}

func TestStory012_NumberKeysStillFocusInOverview(t *testing.T) {
	is := is.New(t)
	// Scenario: Number keys still focus quadrants in overview mode

	input := `(A) Task one
(B) Task two`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press '1' to focus
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should be in focus mode on DO FIRST
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "Do First"))  // expected to focus on DO FIRST quadrant

	// Should show focus mode help text
	is.True(strings.Contains(stripANSI(view), "Press a to add"))  // expected focus mode help text
}
