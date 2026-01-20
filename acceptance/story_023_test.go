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

// Story 023: Archive Completed Todos

func TestStory023_ArchiveCompletedTodo(t *testing.T) {
	// Scenario: Archive a completed todo from focused quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	activeTodo := todo.New("Active task", todo.PriorityA)
	completedTodo := todo.New("Completed task", todo.PriorityA).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{activeTodo, completedTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Select the completed todo (should be at index 1, after the active todo)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // Move down
	model = updatedModel.(ui.Model)

	// Archive it with 'd'
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	// Verify the completed todo was removed from todo.txt
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1) // Should only have the active todo left

	// Verify the todo was added to archive
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "x"))                // Should have completion marker
	is.True(strings.Contains(archive, "Completed task"))   // Should have the description
}

func TestStory023_CannotArchiveUncompletedTodo(t *testing.T) {
	// Scenario: Cannot archive uncompleted todos
	is := is.New(t)

	repository := memory.NewRepository()
	activeTodo := todo.New("Active task", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{activeTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to archive with 'd'
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	// Verify nothing was archived
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1) // Todo should still be there

	archive := repository.ArchiveString()
	is.Equal(archive, "") // Archive should be empty
}

func TestStory023_ArchiveFromDifferentQuadrants(t *testing.T) {
	// Scenario: Archive from any quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	completedA := todo.New("Task A", todo.PriorityA).ToggleCompletion(time.Now())
	completedB := todo.New("Task B", todo.PriorityB).ToggleCompletion(time.Now())
	completedC := todo.New("Task C", todo.PriorityC).ToggleCompletion(time.Now())
	completedD := todo.New("Task D", todo.PriorityD).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{completedA, completedB, completedC, completedD})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Archive from Do First (quadrant 1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model = updatedModel.(ui.Model)

	// Archive from Schedule (quadrant 2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model = updatedModel.(ui.Model)

	// Archive from Delegate (quadrant 3)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	model = updatedModel.(ui.Model)

	// Archive from Eliminate (quadrant 4)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	// Verify all todos were archived
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 0) // All should be archived

	// Verify all are in archive
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "Task A"))
	is.True(strings.Contains(archive, "Task B"))
	is.True(strings.Contains(archive, "Task C"))
	is.True(strings.Contains(archive, "Task D"))
}

func TestStory023_ArchivePreservesMetadata(t *testing.T) {
	// Scenario: Archived todo preserves all metadata
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	completionDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	complexTodo := todo.NewCompletedWithTagsAndDates(
		"Complex task",
		todo.PriorityA,
		&completionDate,
		&creationDate,
		[]string{"project"},
		[]string{"context"},
	)

	err := repository.SaveAll([]todo.Todo{complexTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus and archive
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	// Verify archive contains all metadata
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "x"))             // Completion marker
	is.True(strings.Contains(archive, "2026-01-20"))    // Completion date
	is.True(strings.Contains(archive, "2026-01-15"))    // Creation date
	is.True(strings.Contains(archive, "(A)"))           // Priority
	is.True(strings.Contains(archive, "Complex task"))  // Description
	is.True(strings.Contains(archive, "+project"))      // Project tag
	is.True(strings.Contains(archive, "@context"))      // Context tag
}

func TestStory023_CannotArchiveFromOverview(t *testing.T) {
	// Scenario: Cannot archive from overview mode
	is := is.New(t)

	repository := memory.NewRepository()
	completedTodo := todo.New("Completed task", todo.PriorityA).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{completedTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Try to archive from overview (should do nothing)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	// Verify nothing was archived
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1) // Todo should still be there

	archive := repository.ArchiveString()
	is.Equal(archive, "") // Archive should be empty
}
