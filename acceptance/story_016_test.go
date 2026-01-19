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

func TestStory016_EnterMoveModeWithMKey(t *testing.T) {
	is := is.New(t)
	// Scenario: Enter move mode with 'm' key

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should see move mode overlay
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:"))   // expected to see move mode overlay title
	is.True(strings.Contains(stripANSI(view), "1. Do First"))         // expected to see DO FIRST option
	is.True(strings.Contains(stripANSI(view), "2. Schedule"))         // expected to see SCHEDULE option
	is.True(strings.Contains(stripANSI(view), "3. Delegate"))         // expected to see DELEGATE option
	is.True(strings.Contains(stripANSI(view), "4. Eliminate"))        // expected to see ELIMINATE option
	is.True(strings.Contains(stripANSI(view), "Press ESC to cancel")) // expected to see cancel instruction
}

func TestStory016_SelectDestinationQuadrant(t *testing.T) {
	is := is.New(t)
	// Scenario: Select destination quadrant

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '2' for SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view := model.View()
	is.True(!strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to exit move mode after selection

	// Should have moved todo to priority B
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityB)
	is.Equal(savedTodos[0].Description(), "Review quarterly goals")
}

func TestStory016_CancelMoveModeWithESC(t *testing.T) {
	is := is.New(t)
	// Scenario: Cancel move mode with ESC

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should be in move mode
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to be in move mode

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view = model.View()
	is.True(!strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to exit move mode after ESC

	// Todo should remain in DO FIRST (priority A)
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
	is.Equal(savedTodos[0].Description(), "Review quarterly goals")
}

func TestStory016_MoveToEachQuadrant(t *testing.T) {
	// Scenario: Move to each quadrant

	tests := []struct {
		name             string
		initialPriority  string
		destinationKey   string
		expectedPriority string
	}{
		{"Move to DO FIRST", "(B)", "1", "(A)"},
		{"Move to SCHEDULE", "(A)", "2", "(B)"},
		{"Move to DELEGATE", "(A)", "3", "(C)"},
		{"Move to ELIMINATE", "(A)", "4", "(D)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			repository := memory.NewRepository()
			var priority todo.Priority
			switch tt.initialPriority {
			case "(A)":
				priority = todo.PriorityA
			case "(B)":
				priority = todo.PriorityB
			case "(C)":
				priority = todo.PriorityC
			case "(D)":
				priority = todo.PriorityD
			default:
				priority = todo.PriorityNone
			}
			err := repository.SaveAll([]todo.Todo{
				todo.New("Test task", priority),
			})
			is.NoErr(err)

			m, err := usecases.LoadMatrix(repository)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			model := ui.NewModelWithRepository(m, "test.txt", repository)
			updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			model = updatedModel.(ui.Model)

			// Focus on the appropriate quadrant based on initial priority
			var focusKey rune
			switch tt.initialPriority {
			case "(A)":
				focusKey = '1'
			case "(B)":
				focusKey = '2'
			case "(C)":
				focusKey = '3'
			case "(D)":
				focusKey = '4'
			}

			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{focusKey}})
			model = updatedModel.(ui.Model)

			// Enter move mode
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
			model = updatedModel.(ui.Model)

			// Press destination key
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.destinationKey)})
			_ = updatedModel.(ui.Model)

			// Check file was updated with expected priority
			savedTodos, err := repository.LoadAll()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(savedTodos) != 1 {
				t.Fatalf("expected 1 todo, got %d", len(savedTodos))
			}

			// Map expected priority string to priority constant
			var expectedPriority todo.Priority
			switch tt.expectedPriority {
			case "(A)":
				expectedPriority = todo.PriorityA
			case "(B)":
				expectedPriority = todo.PriorityB
			case "(C)":
				expectedPriority = todo.PriorityC
			case "(D)":
				expectedPriority = todo.PriorityD
			}

			if savedTodos[0].Priority() != expectedPriority {
				t.Errorf("expected priority %v, got %v", expectedPriority, savedTodos[0].Priority())
			}
			if savedTodos[0].Description() != "Test task" {
				t.Errorf("expected description 'Test task', got '%s'", savedTodos[0].Description())
			}
		})
	}
}

func TestStory016_MovingToCurrentQuadrantIsNoOp(t *testing.T) {
	// Scenario: Moving to current quadrant is no-op
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '1' to move to DO FIRST (same quadrant)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view := model.View()
	if strings.Contains(stripANSI(view), "Move to quadrant:") {
		t.Error("expected to exit move mode after selection")
	}

	// Todo should remain in DO FIRST (no write should occur since priority unchanged)
	savedTodos, err := repository.LoadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(savedTodos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(savedTodos))
	}
	if savedTodos[0].Priority() != todo.PriorityA {
		t.Errorf("expected priority A, got %v", savedTodos[0].Priority())
	}

	// Verify still viewing DO FIRST with the todo
	view = model.View()
	if !strings.Contains(stripANSI(view), "Review quarterly goals") {
		t.Error("expected todo to still be visible in DO FIRST")
	}
}

func TestStory016_MoveModeOnlyAvailableInFocusMode(t *testing.T) {
	// Scenario: Move mode only available in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press 'm'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should NOT show move mode overlay
	view := model.View()
	if strings.Contains(stripANSI(view), "Move to quadrant:") {
		t.Error("expected 'm' to do nothing in overview mode")
	}

	// Should still be in overview mode
	if !strings.Contains(stripANSI(view), "Do First") && !strings.Contains(stripANSI(view), "Schedule") {
		t.Error("expected to remain in overview mode")
	}
}

func TestStory016_MoveModeOnlyAvailableWhenTodoSelected(t *testing.T) {
	// Scenario: Move mode only available when todo selected

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task in DO FIRST", todo.PriorityA),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, err := usecases.LoadMatrix(repository)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE (empty quadrant)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press 'm'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should NOT show move mode overlay
	view := model.View()
	if strings.Contains(stripANSI(view), "Move to quadrant:") {
		t.Error("expected 'm' to do nothing in empty quadrant")
	}

	// Should show "(no tasks)"
	if !strings.Contains(stripANSI(view), "(no tasks)") {
		t.Error("expected to see '(no tasks)' in empty quadrant")
	}
}

func TestStory016_OtherKeysIgnoredInMoveMode(t *testing.T) {
	is := is.New(t)
	// Scenario: Other keys should be ignored while in move mode

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review quarterly goals", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should be in move mode
	view := model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to be in move mode

	// Press various other keys (should be ignored)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)
	view = model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to remain in move mode after pressing 'a'

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)
	view = model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to remain in move mode after pressing space

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)
	view = model.View()
	is.True(strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to remain in move mode after pressing down arrow

	// Should still have no changes to the file
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
	is.Equal(savedTodos[0].Description(), "Review quarterly goals")

	// Press ESC to exit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view = model.View()
	is.True(!strings.Contains(stripANSI(view), "Move to quadrant:")) // expected to exit move mode after ESC
}
