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

// Story 027: Bulk Archive Completed Tasks

func TestStory027_BulkArchiveInFocusedQuadrant(t *testing.T) {
	// Scenario: Bulk archive completed tasks in focused quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	active1 := todo.New("Active 1", todo.PriorityA)
	active2 := todo.New("Active 2", todo.PriorityA)
	completed1 := todo.New("Completed 1", todo.PriorityA).ToggleCompletion(time.Now())
	completed2 := todo.New("Completed 2", todo.PriorityA).ToggleCompletion(time.Now())
	completed3 := todo.New("Completed 3", todo.PriorityA).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{active1, active2, completed1, completed2, completed3})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press Shift+D to bulk archive
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify only active todos remain
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 2)

	// Verify all completed todos were archived
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "Completed 1"))
	is.True(strings.Contains(archive, "Completed 2"))
	is.True(strings.Contains(archive, "Completed 3"))
	is.True(!strings.Contains(archive, "Active 1"))
	is.True(!strings.Contains(archive, "Active 2"))
}

func TestStory027_BulkArchiveInOverviewMode(t *testing.T) {
	// Scenario: Bulk archive completed tasks in overview mode (all quadrants)
	is := is.New(t)

	repository := memory.NewRepository()
	doFirstActive := todo.New("DoFirst Active", todo.PriorityA)
	doFirstCompleted1 := todo.New("DoFirst Completed 1", todo.PriorityA).ToggleCompletion(time.Now())
	doFirstCompleted2 := todo.New("DoFirst Completed 2", todo.PriorityA).ToggleCompletion(time.Now())
	scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())
	delegateActive := todo.New("Delegate Active", todo.PriorityC)
	delegateCompleted1 := todo.New("Delegate Completed 1", todo.PriorityC).ToggleCompletion(time.Now())
	delegateCompleted2 := todo.New("Delegate Completed 2", todo.PriorityC).ToggleCompletion(time.Now())
	delegateCompleted3 := todo.New("Delegate Completed 3", todo.PriorityC).ToggleCompletion(time.Now())
	// Eliminate has no completed todos
	eliminateActive := todo.New("Eliminate Active", todo.PriorityD)

	err := repository.SaveAll([]todo.Todo{
		doFirstActive, doFirstCompleted1, doFirstCompleted2,
		scheduleCompleted,
		delegateActive, delegateCompleted1, delegateCompleted2, delegateCompleted3,
		eliminateActive,
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press Shift+D to bulk archive all
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify only active todos remain (3 active across all quadrants)
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 3)

	// Verify all 6 completed todos were archived
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "DoFirst Completed 1"))
	is.True(strings.Contains(archive, "DoFirst Completed 2"))
	is.True(strings.Contains(archive, "Schedule Completed"))
	is.True(strings.Contains(archive, "Delegate Completed 1"))
	is.True(strings.Contains(archive, "Delegate Completed 2"))
	is.True(strings.Contains(archive, "Delegate Completed 3"))
}

func TestStory027_NoCompletedTasksInFocusedQuadrant(t *testing.T) {
	// Scenario: No completed tasks to archive in focused quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	active1 := todo.New("Active 1", todo.PriorityA)
	active2 := todo.New("Active 2", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{active1, active2})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press Shift+D (should be no-op)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify nothing changed
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 2)

	archive := repository.ArchiveString()
	is.Equal(archive, "")
}

func TestStory027_NoCompletedTasksInOverviewMode(t *testing.T) {
	// Scenario: No completed tasks to archive in overview mode
	is := is.New(t)

	repository := memory.NewRepository()
	active1 := todo.New("Active 1", todo.PriorityA)
	active2 := todo.New("Active 2", todo.PriorityB)
	active3 := todo.New("Active 3", todo.PriorityC)

	err := repository.SaveAll([]todo.Todo{active1, active2, active3})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press Shift+D (should be no-op)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify nothing changed
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 3)

	archive := repository.ArchiveString()
	is.Equal(archive, "")
}

func TestStory027_EmptyQuadrantInFocusMode(t *testing.T) {
	// Scenario: Empty quadrant in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	// Only add a todo in DoFirst, leave Schedule empty
	active := todo.New("Active", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{active})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Schedule quadrant (which is empty)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press Shift+D (should be no-op)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify nothing changed
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)

	archive := repository.ArchiveString()
	is.Equal(archive, "")
}

func TestStory027_OnlyArchivesFromSpecifiedQuadrant(t *testing.T) {
	// Scenario: In focus mode, only archives from the focused quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	doFirstCompleted := todo.New("DoFirst Completed", todo.PriorityA).ToggleCompletion(time.Now())
	scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{doFirstCompleted, scheduleCompleted})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant only
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press Shift+D to bulk archive
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify only DoFirst todo was archived, Schedule todo remains
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.Equal(todos[0].Description(), "Schedule Completed")

	// Verify only DoFirst was archived
	archive := repository.ArchiveString()
	is.True(strings.Contains(archive, "DoFirst Completed"))
	is.True(!strings.Contains(archive, "Schedule Completed"))
}

func TestStory027_CannotBulkArchiveInReadOnlyMode(t *testing.T) {
	// Scenario: Cannot bulk archive when in read-only mode
	is := is.New(t)

	repository := memory.NewRepository()
	completed := todo.New("Completed", todo.PriorityA).ToggleCompletion(time.Now())

	err := repository.SaveAll([]todo.Todo{completed})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	model = model.SetReadOnly(true)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press Shift+D (should be no-op in read-only mode)
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'D'}})

	// Verify nothing was archived
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)

	archive := repository.ArchiveString()
	is.Equal(archive, "")
}
