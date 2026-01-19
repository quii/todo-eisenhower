package acceptance_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 020: Edit Todo

func TestStory020_PressEToEnterEditMode(t *testing.T) {
	// Scenario: Press 'e' to enter edit mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Buy milk", todo.PriorityA, []string{"shopping"}, []string{"store"}),
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

	// Press 'e' to enter edit mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should be in edit mode
	is.True(strings.Contains(stripANSI(view), "Edit Todo"))

	// Should show the tag reference panel
	is.True(strings.Contains(stripANSI(view), "Projects") || strings.Contains(stripANSI(view), "Contexts"))
}

func TestStory020_EditTodoDescription(t *testing.T) {
	// Scenario: Edit todo description
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Original task", todo.PriorityA, []string{"shopping"}, []string{"store"}),
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

	// Press 'e' to enter edit mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	// Clear the pre-filled input (Ctrl+U clears the line in textinput)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	model = updatedModel.(ui.Model)

	// Type new description
	for _, ch := range "Updated task +shopping @store " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = updatedModel.(ui.Model)

	// Check the todo was updated
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Updated task")
	is.Equal(savedTodos[0].Projects(), []string{"shopping"})
	is.Equal(savedTodos[0].Contexts(), []string{"store"})
}

func TestStory020_EditAndRemoveTags(t *testing.T) {
	// Scenario: Edit and remove tags
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Review code", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	// Clear the pre-filled input (Ctrl+U clears the line in textinput)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	model = updatedModel.(ui.Model)

	// Type description without tags
	for _, ch := range "Review code" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = updatedModel.(ui.Model)

	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Review code")
	is.Equal(len(savedTodos[0].Projects()), 0)
	is.Equal(len(savedTodos[0].Contexts()), 0)
}

func TestStory020_PreserveCreationDate(t *testing.T) {
	// Scenario: Preserve creation date
	is := is.New(t)

	creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Original task", todo.PriorityA, &creationDate),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	// Clear the pre-filled input (Ctrl+U clears the line in textinput)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	model = updatedModel.(ui.Model)

	for _, ch := range "Updated description" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = updatedModel.(ui.Model)

	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Updated description")
	is.True(savedTodos[0].CreationDate() != nil)
	is.Equal(savedTodos[0].CreationDate().Format("2006-01-02"), "2026-01-15")
}

func TestStory020_PreserveCompletionDate(t *testing.T) {
	// Scenario: Preserve completion date for completed todos
	is := is.New(t)

	completionDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Task done", todo.PriorityA, &completionDate),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	for _, ch := range "Task done with more details" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = updatedModel.(ui.Model)

	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())
	is.True(savedTodos[0].CompletionDate() != nil)
	is.Equal(savedTodos[0].CompletionDate().Format("2006-01-02"), "2026-01-18")
}

func TestStory020_CancelEditWithESC(t *testing.T) {
	// Scenario: Cancel edit with ESC
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Original task", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	// Type something
	for _, ch := range "Changed" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = updatedModel.(ui.Model)

	// Todo should remain unchanged
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "Original task")
}

func TestStory020_EditModeOnlyInFocusMode(t *testing.T) {
	// Scenario: Edit mode only available in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press 'e'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should NOT be in edit mode
	is.True(!strings.Contains(stripANSI(view), "Edit Todo"))
}
