package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 018: Delete Todo

func TestStory018_DeleteTodoWithConfirmation(t *testing.T) {
	// Scenario: Delete a todo with confirmation
	is := is.New(t)

	input := "(A) Task to keep\n(A) Task to delete\n"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Navigate to second todo
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Backspace to enter delete mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show delete confirmation dialog
	is.True(strings.Contains(stripANSI(view), "Delete this todo?")) // expected delete confirmation

	// Press 'y' to confirm
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
	model = updatedModel.(ui.Model)

	// Todo should be deleted from matrix
	updatedMatrix := model.GetMatrix()
	is.Equal(len(updatedMatrix.DoFirst()), 1)
	is.Equal(updatedMatrix.DoFirst()[0].Description(), "Task to keep")

	// Verify the todo was deleted by checking the repository
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Task to keep")
}

func TestStory018_CancelDeletionWithESC(t *testing.T) {
	// Scenario: Cancel deletion with ESC
	is := is.New(t)

	input := "(A) Task to keep\n"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus and enter delete mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updatedModel.(ui.Model)

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Todo should NOT be deleted
	updatedMatrix := model.GetMatrix()
	is.Equal(len(updatedMatrix.DoFirst()), 1)

	// Verify no changes were made to the repository
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Task to keep")
}

func TestStory018_CancelDeletionWithN(t *testing.T) {
	// Scenario: Cancel deletion with 'n'
	is := is.New(t)

	input := "(A) Task to keep\n"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus and enter delete mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model = updatedModel.(ui.Model)

	// Press 'n' to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	model = updatedModel.(ui.Model)

	// Todo should NOT be deleted
	updatedMatrix := model.GetMatrix()
	is.Equal(len(updatedMatrix.DoFirst()), 1)

	// Verify no changes were made to the repository
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Task to keep")
}
