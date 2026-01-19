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

// Story 021: Filter by Tag

func TestStory021_PressFToEnterFilterMode(t *testing.T) {
	// Scenario: Press 'f' to enter filter mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 'f' to enter filter mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show filter input prompt
	is.True(strings.Contains(stripANSI(view), "Filter by:"))
}

func TestStory021_FilterByProject(t *testing.T) {
	// Scenario: Filter by project
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
		todo.NewWithTags("Task 3", todo.PriorityA, []string{"WebApp"}, []string{"phone"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 'f' to enter filter mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	// Type "+WebApp"
	for _, ch := range "+WebApp " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to apply filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show filtered todos (only those with +WebApp)
	is.True(strings.Contains(stripANSI(view), "Task 1"))
	is.True(strings.Contains(stripANSI(view), "Task 3"))
	is.True(!strings.Contains(stripANSI(view), "Task 2"))
}

func TestStory021_FilterByContext(t *testing.T) {
	// Scenario: Filter by context
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
		todo.NewWithTags("Task 3", todo.PriorityC, []string{"Backend"}, []string{"computer"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 'f' to enter filter mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	// Type "@computer"
	for _, ch := range "@computer " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to apply filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show filtered todos (only those with @computer)
	is.True(strings.Contains(stripANSI(view), "Task 1"))
	is.True(strings.Contains(stripANSI(view), "Task 3"))
	is.True(!strings.Contains(stripANSI(view), "Task 2"))
}

func TestStory021_ClearFilter(t *testing.T) {
	// Scenario: Clear filter and return to overview
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Apply filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	for _, ch := range "+WebApp " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Press 'c' to clear filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should return to overview showing all todos (filter cleared)
	is.True(strings.Contains(stripANSI(view), "Task 1"))
	is.True(strings.Contains(stripANSI(view), "Task 2"))
}

func TestStory021_CancelFilterInput(t *testing.T) {
	// Scenario: Cancel filter input
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 'f' to enter filter mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	// Type something
	for _, ch := range "+Web" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should return to overview without filter applied
	is.True(!strings.Contains(stripANSI(view), "Filter by:"))
}

func TestStory021_FilterPersistsWhenNavigating(t *testing.T) {
	// Scenario: Filter persists when navigating quadrants
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		todo.NewWithTags("Task 2", todo.PriorityB, []string{"WebApp"}, []string{"phone"}),
		todo.NewWithTags("Task 3", todo.PriorityA, []string{"Mobile"}, []string{"computer"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Apply filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	model = updatedModel.(ui.Model)

	for _, ch := range "+WebApp " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Press '1' to focus on Do First
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should only see filtered todos
	is.True(strings.Contains(stripANSI(view), "Task 1"))
	is.True(!strings.Contains(stripANSI(view), "Task 3"))

	// Press '0' to return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}})
	model = updatedModel.(ui.Model)

	view = model.View()

	// Filter should still be active (only shows Task 1 and 2, not Task 3)
	is.True(strings.Contains(stripANSI(view), "Task 1"))
	is.True(strings.Contains(stripANSI(view), "Task 2"))
	is.True(!strings.Contains(stripANSI(view), "Task 3"))
}
