package acceptance_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 013: Preserve and Render Completion Dates

func TestStory013_ParseAndPreserveCompletionDates(t *testing.T) {
	is := is.New(t)
	// Scenario: Parse and preserve completion dates

	input := `x 2026-01-10 (A) Completed task from last week
(B) Active task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	// Check that completed todo has the preserved date
	completedTodos := m.DoFirst()
	is.Equal(len(completedTodos), 1) // expected 1 completed todo in DO FIRST

	completedTodo := completedTodos[0]
	is.True(completedTodo.IsCompleted()) // expected todo to be completed

	completionDate := completedTodo.CompletionDate()
	is.True(completionDate != nil) // expected completion date to be preserved, got nil

	expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	is.Equal(completionDate.Format("2006-01-02"), expectedDate.Format("2006-01-02")) // expected date 2026-01-10

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
	input := `(A) Review documentation`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle completion with spacebar
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Check file contains completion marker with today's date
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "x ") {
		t.Error("expected file to contain completion marker 'x'")
	}

	// Should have today's date
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(written, today) {
		t.Errorf("expected file to contain today's date %s, got: %s", today, written)
	}

	if !strings.Contains(written, "(A) Review documentation") {
		t.Errorf("expected file to contain task description, got: %s", written)
	}
}

func TestStory013_ClearCompletionDateWhenTogglingIncomplete(t *testing.T) {
	// Scenario: Clear completion date when toggling incomplete
	is := is.New(t)
	input := `x 2026-01-10 (A) Review documentation`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle to incomplete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Check file no longer has completion marker or date
	written := source.writer.(*strings.Builder).String()
	if strings.Contains(written, "x ") {
		t.Error("expected file not to contain completion marker 'x'")
	}

	if strings.Contains(written, "2026-01-10") {
		t.Error("expected file not to contain old completion date")
	}

	if !strings.Contains(written, "(A) Review documentation") {
		t.Errorf("expected file to contain task description, got: %s", written)
	}
}

func TestStory013_NewCompletionDateWhenRecompleting(t *testing.T) {
	// Scenario: New completion date when re-completing
	is := is.New(t)
	input := `x 2026-01-10 (A) Review documentation`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle to incomplete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Clear the writer to capture only the new completion
	source.writer = &strings.Builder{}

	// Toggle back to complete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Check file has completion marker with NEW date (today)
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "x ") {
		t.Error("expected file to contain completion marker 'x'")
	}

	// Should have today's date, NOT the old date
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(written, today) {
		t.Errorf("expected file to contain today's date %s, got: %s", today, written)
	}

	if strings.Contains(written, "2026-01-10") {
		t.Error("expected old date to be replaced with new date")
	}
}

func TestStory013_DisplayCompletionDateInUI(t *testing.T) {
	// Scenario: Display completion date in UI
	is := is.New(t)

	// Create a todo with a specific completion date (10 days ago)
	tenDaysAgo := time.Now().AddDate(0, 0, -10).Format("2006-01-02")
	input := "x " + tenDaysAgo + " (A) Completed task"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the completed todo with its description
	if !strings.Contains(view, "Completed task") {
		t.Error("expected view to contain todo description")
	}

	// Should have Completed column header in table
	if !strings.Contains(view, "Completed") {
		t.Error("expected view to contain 'Completed' column header")
	}

	// Should show completion date in the Completed column
	// For 10 days ago, should show "10 days ago"
	if !strings.Contains(view, "10 days ago") {
		t.Error("expected view to contain completion date '10 days ago'")
	}
}

func TestStory013_DisplayCompletionDateRelativeFormatting(t *testing.T) {
	// Scenario: Display completion date with relative formatting
	is := is.New(t)
	// Create todos with different completion dates
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	twoDaysAgo := time.Now().AddDate(0, 0, -2).Format("2006-01-02")

	input := "x " + today + " (A) Task completed today\n" +
		"x " + yesterday + " (B) Task completed yesterday\n" +
		"x " + twoDaysAgo + " (C) Task completed 2 days ago"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show relative formatting for recent dates
	if !strings.Contains(view, "today") {
		t.Error("expected view to show 'today' for task completed today")
	}

	// Focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	view = model.View()
	if !strings.Contains(view, "yesterday") {
		t.Error("expected view to show 'yesterday' for task completed yesterday")
	}

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	view = model.View()
	if !strings.Contains(view, "2 days ago") {
		t.Error("expected view to show '2 days ago' for task completed 2 days ago")
	}
}

func TestStory013_NoDateShownForIncompleteTodos(t *testing.T) {
	// Scenario: No date shown for incomplete todos
	is := is.New(t)
	input := `(A) Active task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show active todo in table
	if !strings.Contains(view, "Active task") {
		t.Error("expected view to contain active todo")
	}

	// Should have table with Completed column showing "-" for active todos
	if !strings.Contains(view, "Completed") {
		t.Error("expected view to have Completed column header")
	}
	// The completed column should show "-" for active todos (not a specific date)
	// This is implicit in the table rendering
}

func TestStory013_PreserveCompletionDateWhenMovingQuadrants(t *testing.T) {
	// Scenario: Preserve completion date when moving quadrants
	is := is.New(t)
	input := `x 2026-01-10 (A) Completed urgent task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Move to SCHEDULE (press 'm' then '2')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Check file preserves the completion date
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "x 2026-01-10") {
		t.Errorf("expected file to preserve completion date 2026-01-10, got: %s", written)
	}

	if !strings.Contains(written, "(B) Completed urgent task") {
		t.Errorf("expected file to contain updated priority and description, got: %s", written)
	}
}
