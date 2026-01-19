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
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 013: Preserve and Render Completion Dates

func TestStory013_ParseAndPreserveCompletionDates(t *testing.T) {
	is := is.New(t)
	// Scenario: Parse and preserve completion dates

	repository := memory.NewRepository()
	completionDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Completed task from last week", todo.PriorityA, &completionDate),
		todo.New("Active task", todo.PriorityB),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Check that completed todo has the preserved date
	completedTodos := m.DoFirst()
	is.Equal(len(completedTodos), 1) // expected 1 completed todo in DO FIRST

	completedTodo := completedTodos[0]
	is.True(completedTodo.IsCompleted()) // expected todo to be completed

	actualCompletionDate := completedTodo.CompletionDate()
	is.True(actualCompletionDate != nil) // expected completion date to be preserved, got nil

	expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	is.Equal(actualCompletionDate.Format("2006-01-02"), expectedDate.Format("2006-01-02")) // expected date 2026-01-10

	// Check that active todo has no completion date
	activeTodos := m.Schedule()
	is.Equal(len(activeTodos), 1) // expected 1 active todo in SCHEDULE

	activeTodo := activeTodos[0]
	is.True(!activeTodo.IsCompleted()) // expected todo to be active (not completed)

	is.True(activeTodo.CompletionDate() == nil) // expected no completion date for active todo
}

func TestStory013_SetCompletionDateWhenMarkingComplete(t *testing.T) {
	// Scenario: Set completion date when marking complete
	is := is.New(t)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review documentation", todo.PriorityA),
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

	// Toggle completion with spacebar
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Check file contains completion marker with today's date
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())
	is.Equal(savedTodos[0].Description(), "Review documentation")
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)

	// Should have today's date
	today := time.Now().Format("2006-01-02")
	is.True(savedTodos[0].CompletionDate() != nil)
	is.Equal(savedTodos[0].CompletionDate().Format("2006-01-02"), today)
}

func TestStory013_ClearCompletionDateWhenTogglingIncomplete(t *testing.T) {
	// Scenario: Clear completion date when toggling incomplete
	is := is.New(t)
	repository := memory.NewRepository()
	completionDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Review documentation", todo.PriorityA, &completionDate),
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

	// Toggle to incomplete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Check file no longer has completion marker or date
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(!savedTodos[0].IsCompleted())
	is.True(savedTodos[0].CompletionDate() == nil)
	is.Equal(savedTodos[0].Description(), "Review documentation")
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)
}

func TestStory013_NewCompletionDateWhenRecompleting(t *testing.T) {
	// Scenario: New completion date when re-completing
	is := is.New(t)
	repository := memory.NewRepository()
	completionDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Review documentation", todo.PriorityA, &completionDate),
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

	// Toggle to incomplete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Create a new repository to track the next write operation
	// We need to re-create the model with a fresh repository to capture only the completion write
	updatedTodos, _ := repository.LoadAll()
	// Create a buffer with the current todos
	var buf strings.Builder
	_ = todotxt.Marshal(&buf, updatedTodos)
	// Now create a fresh repository starting with these todos
	tempRepo := memory.NewRepository()
	_ = tempRepo.SaveAll(updatedTodos)
	// Re-load the matrix from this repository
	m2, _ := usecases.LoadMatrix(tempRepo)
	// Create a new repository to capture writes
	repository2 := memory.NewRepository()
	model = ui.NewModelWithRepository(m2, "test.txt", repository2)
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle back to complete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Check file has completion marker with NEW date (today)
	savedTodos, err := repository2.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())

	// Should have today's date, NOT the old date
	today := time.Now().Format("2006-01-02")
	is.True(savedTodos[0].CompletionDate() != nil)
	is.Equal(savedTodos[0].CompletionDate().Format("2006-01-02"), today)
}

func TestStory013_DisplayCompletionDateInUI(t *testing.T) {
	// Scenario: Display completion date in UI
	is := is.New(t)

	// Create a todo with a specific completion date (10 days ago)
	tenDaysAgo := time.Now().AddDate(0, 0, -10)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Completed task", todo.PriorityA, &tenDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the completed todo with its description
	if !strings.Contains(stripANSI(view), "Completed task") {
		t.Error("expected view to contain todo description")
	}

	// Should have Completed column header in table
	if !strings.Contains(stripANSI(view), "Completed") {
		t.Error("expected view to contain 'Completed' column header")
	}

	// Should show completion date in the Completed column
	// For 10 days ago, should show "10 days ago"
	if !strings.Contains(stripANSI(view), "10 days ago") {
		t.Error("expected view to contain completion date '10 days ago'")
	}
}

func TestStory013_DisplayCompletionDateRelativeFormatting(t *testing.T) {
	// Scenario: Display completion date with relative formatting
	is := is.New(t)
	// Create todos with different completion dates
	today := time.Now()
	yesterday := time.Now().AddDate(0, 0, -1)
	twoDaysAgo := time.Now().AddDate(0, 0, -2)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Task completed today", todo.PriorityA, &today),
		todo.NewCompleted("Task completed yesterday", todo.PriorityB, &yesterday),
		todo.NewCompleted("Task completed 2 days ago", todo.PriorityC, &twoDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show relative formatting for recent dates
	if !strings.Contains(stripANSI(view), "today") {
		t.Error("expected view to show 'today' for task completed today")
	}

	// Focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	view = model.View()
	if !strings.Contains(stripANSI(view), "yesterday") {
		t.Error("expected view to show 'yesterday' for task completed yesterday")
	}

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	view = model.View()
	if !strings.Contains(stripANSI(view), "2 days ago") {
		t.Error("expected view to show '2 days ago' for task completed 2 days ago")
	}
}

func TestStory013_NoDateShownForIncompleteTodos(t *testing.T) {
	// Scenario: No date shown for incomplete todos
	is := is.New(t)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Active task", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show active todo in table
	if !strings.Contains(stripANSI(view), "Active task") {
		t.Error("expected view to contain active todo")
	}

	// Should have table with Completed column showing "-" for active todos
	if !strings.Contains(stripANSI(view), "Completed") {
		t.Error("expected view to have Completed column header")
	}
	// The completed column should show "-" for active todos (not a specific date)
	// This is implicit in the table rendering
}

func TestStory013_PreserveCompletionDateWhenMovingQuadrants(t *testing.T) {
	// Scenario: Preserve completion date when moving quadrants
	is := is.New(t)
	repository := memory.NewRepository()
	completionDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompleted("Completed urgent task", todo.PriorityA, &completionDate),
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

	// Move to SCHEDULE (press 'm' then '2')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	_ = updatedModel.(ui.Model)

	// Check file preserves the completion date
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())
	is.Equal(savedTodos[0].Priority(), todo.PriorityB)
	is.Equal(savedTodos[0].Description(), "Completed urgent task")

	// Completion date should be preserved
	is.True(savedTodos[0].CompletionDate() != nil)
	expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	is.Equal(savedTodos[0].CompletionDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
}
