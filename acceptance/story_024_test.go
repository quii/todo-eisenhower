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

// Story 024: Stdin Read-Only Mode

func TestStory024_ReadOnlyModeIndicator(t *testing.T) {
	// Scenario: Read-only mode indicator
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Create model in read-only mode
	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should indicate read-only mode
	is.True(strings.Contains(view, "(read-only)") || strings.Contains(view, "read-only"))
}

func TestStory024_ViewingOperationsWork(t *testing.T) {
	// Scenario: Viewing operations still work
	is := is.New(t)

	repository := memory.NewRepository()
	taskA := todo.New("Task A", todo.PriorityA)
	taskB := todo.New("Task B", todo.PriorityB)
	err := repository.SaveAll([]todo.Todo{taskA, taskB})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Navigation should work - focus on Do First
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view := model.View()
	is.True(strings.Contains(view, "Task A")) // Should show focused content

	// Navigate to Schedule
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view = model.View()
	is.True(strings.Contains(view, "Task B"))

	// Return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)
	view = model.View()
	is.True(strings.Contains(view, "Do First"))   // Back to overview
	is.True(strings.Contains(view, "Schedule"))
}

func TestStory024_AddTodoDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Add
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'a' to add
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should not enter add mode - no input field visible
	is.True(!strings.Contains(view, "Add todo:"))
}

func TestStory024_EditTodoDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Edit
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'e' to edit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should not enter edit mode
	is.True(!strings.Contains(view, "Edit Todo:"))
}

func TestStory024_ToggleCompletionDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Toggle completion
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to toggle completion with space
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeySpace})
	_ = updatedModel.(ui.Model)

	// Verify the task is still not completed (repository wasn't modified)
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.True(!todos[0].IsCompleted()) // Should still be incomplete
}

func TestStory024_DeleteDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Delete
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'x' to delete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should not show delete confirmation
	is.True(!strings.Contains(view, "Delete"))
	is.True(!strings.Contains(view, "confirm"))
}

func TestStory024_MoveDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Move
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'm' to move
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should not show move overlay
	is.True(!strings.Contains(view, "Move to"))
}

func TestStory024_ArchiveDisabledInReadOnly(t *testing.T) {
	// Scenario: Editing operations are disabled - Archive
	is := is.New(t)

	repository := memory.NewRepository()
	completedTask := todo.NewCompleted("Completed task", todo.PriorityA, nil)
	err := repository.SaveAll([]todo.Todo{completedTask})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'd' to archive
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	_ = updatedModel.(ui.Model)

	// Verify task wasn't archived
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1) // Should still be there
}

func TestStory024_FilteringWorksInReadOnly(t *testing.T) {
	// Scenario: Filtering still works in read-only
	is := is.New(t)

	repository := memory.NewRepository()
	task1 := todo.NewWithTags("Task 1", todo.PriorityA, []string{"project1"}, nil)
	task2 := todo.NewWithTags("Task 2", todo.PriorityA, []string{"project2"}, nil)
	err := repository.SaveAll([]todo.Todo{task1, task2})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Try to enter filter mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Filtering should still work in read-only mode
	is.True(strings.Contains(view, "Filter") || strings.Contains(view, "filter"))
}

func TestStory024_InventoryWorksInReadOnly(t *testing.T) {
	// Scenario: Inventory still works in read-only
	is := is.New(t)

	repository := memory.NewRepository()
	task := todo.New("Task", todo.PriorityA)
	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "(stdin)", repository).SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Try to enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Inventory should work
	is.True(strings.Contains(view, "Inventory") || strings.Contains(view, "Dashboard"))
}
