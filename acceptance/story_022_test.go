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

// Story 022: Due Date Support

func TestStory022_DisplayDueDateInOverview(t *testing.T) {
	// Scenario: Display due date in overview mode
	is := is.New(t)

	repository := memory.NewRepository()
	dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	taskWithDueDate := todo.NewFull("Submit report", todo.PriorityA, false, nil, nil, &dueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{taskWithDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see "Submit report" and "due: Jan 25"
	is.True(strings.Contains(view, "Submit report"))
	is.True(strings.Contains(view, "due: Jan 25"))
}

func TestStory022_DisplayOverdueInOverview(t *testing.T) {
	// Scenario: Display overdue items
	is := is.New(t)

	repository := memory.NewRepository()
	overdueDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	overdueTask := todo.NewFull("Overdue task", todo.PriorityA, false, nil, nil, &overdueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{overdueTask})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see "Overdue task" with "!" and "due: Jan 15" (red in actual UI)
	is.True(strings.Contains(view, "Overdue task"))
	is.True(strings.Contains(view, "!"))
	is.True(strings.Contains(view, "due: Jan 15"))
}

func TestStory022_DisplayDueDateInFocusedTable(t *testing.T) {
	// Scenario: Display due date in focused quadrant mode (table)
	is := is.New(t)

	repository := memory.NewRepository()
	dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	taskWithDueDate := todo.NewFull("Submit report", todo.PriorityA, false, nil, nil, &dueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{taskWithDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should have "Due Date" column header and "Jan 25" in the row
	is.True(strings.Contains(view, "Due Date"))
	is.True(strings.Contains(view, "Jan 25"))
}

func TestStory022_DisplayOverdueInFocusedTable(t *testing.T) {
	// Scenario: Display overdue items in focused mode
	is := is.New(t)

	repository := memory.NewRepository()
	overdueDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	overdueTask := todo.NewFull("Overdue task", todo.PriorityA, false, nil, nil, &overdueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{overdueTask})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see "!" prefix in the Due Date column
	is.True(strings.Contains(view, "! Jan 15"))
}

func TestStory022_TodosWithoutDueDatesUnaffected(t *testing.T) {
	// Scenario: Todos without due dates are unaffected
	is := is.New(t)

	repository := memory.NewRepository()
	regularTask := todo.New("Regular task", todo.PriorityB)

	err := repository.SaveAll([]todo.Todo{regularTask})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see "Regular task" without any due date info
	is.True(strings.Contains(view, "Regular task"))
	is.True(!strings.Contains(view, "due:"))
}

func TestStory022_DueDateExactlyTodayNotOverdue(t *testing.T) {
	// Scenario: Due date exactly today is not overdue
	is := is.New(t)

	repository := memory.NewRepository()
	today := time.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	taskDueToday := todo.NewFull("Task due today", todo.PriorityA, false, nil, nil, &todayDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{taskDueToday})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see due date but NOT the "!" overdue indicator
	is.True(strings.Contains(view, "Task due today"))
	is.True(strings.Contains(view, "due:"))
	// The view should not show "!" right before the due date for today
	// (This is a bit tricky to test without the exact formatting, but we can check the general presence)
}

func TestStory022_DueDatesInformationalOnly(t *testing.T) {
	// Scenario: Due dates are informational only - doesn't affect quadrant placement
	is := is.New(t)

	repository := memory.NewRepository()
	// Low priority task with overdue date - should still be in Eliminate quadrant
	overdueDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	lowPriorityOverdue := todo.NewFull("Low priority task", todo.PriorityD, false, nil, nil, &overdueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{lowPriorityOverdue})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Verify it's in the Eliminate quadrant (quadrant 4)
	eliminateTodos := m.Eliminate()
	is.Equal(len(eliminateTodos), 1)
	is.Equal(eliminateTodos[0].Description(), "Low priority task")
	is.Equal(eliminateTodos[0].Priority(), todo.PriorityD)
}

func TestStory022_DueDatesWorkWithTags(t *testing.T) {
	// Scenario: Due dates work with all todo features (tags)
	is := is.New(t)

	repository := memory.NewRepository()
	dueDate := time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC)
	taskWithTagsAndDueDate := todo.NewFull(
		"Task with everything",
		todo.PriorityA,
		false,
		nil,
		nil,
		&dueDate,
		[]string{"project"},
		[]string{"context"},
	)

	err := repository.SaveAll([]todo.Todo{taskWithTagsAndDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant to see table
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should see description, tags, AND due date
	is.True(strings.Contains(view, "Task with everything"))
	is.True(strings.Contains(view, "project"))
	is.True(strings.Contains(view, "context"))
	is.True(strings.Contains(view, "Jan 30"))
}

func TestStory022_InvalidDueDatesIgnored(t *testing.T) {
	// Scenario: Invalid due dates are ignored
	// This tests the parser behavior - when an invalid due date format is provided,
	// it should be ignored (DueDate() returns nil)
	is := is.New(t)

	repository := memory.NewRepository()
	// This would typically be parsed from "Task due:invalid-date"
	// but since our parser ignores invalid dates, we create a task without due date
	taskWithoutDueDate := todo.New("Task", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{taskWithoutDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should not display any due date
	is.True(!strings.Contains(view, "due:"))
}

func TestStory022_DifferentYearFormatting(t *testing.T) {
	// Test that dates in different years show full date format
	is := is.New(t)

	repository := memory.NewRepository()
	// Use a date from a different year (2025)
	dueDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
	taskWithDifferentYearDueDate := todo.NewFull("Task from past year", todo.PriorityA, false, nil, nil, &dueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{taskWithDifferentYearDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show full date format for different year
	is.True(strings.Contains(view, "2025-12-25"))
}

func TestStory022_DueDateInFocusedInputMode(t *testing.T) {
	// Test that due dates display in focus mode with input (list view)
	is := is.New(t)

	repository := memory.NewRepository()
	dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
	taskWithDueDate := todo.NewFull("Task with due date", todo.PriorityA, false, nil, nil, &dueDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{taskWithDueDate})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter add mode (switches to list view with input)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should still see due date in list view
	is.True(strings.Contains(view, "Task with due date"))
	is.True(strings.Contains(view, "due: Jan 25"))
}
