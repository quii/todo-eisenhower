package acceptance_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/parser"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory014_ParseAndPreserveCreationDates(t *testing.T) {
	// Scenario: Parse and preserve existing creation dates from file

	input := `(A) 2026-01-10 Task created on Jan 10
(B) 2026-01-12 Task created on Jan 12
(C) Task without creation date`

	todos, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todos) != 3 {
		t.Fatalf("expected 3 todos, got %d", len(todos))
	}

	// First todo should have creation date Jan 10
	if todos[0].CreationDate() == nil {
		t.Error("expected first todo to have a creation date")
	} else {
		expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		actualDate := time.Date(todos[0].CreationDate().Year(), todos[0].CreationDate().Month(), todos[0].CreationDate().Day(), 0, 0, 0, 0, time.UTC)
		if !actualDate.Equal(expectedDate) {
			t.Errorf("expected creation date %s, got %s", expectedDate.Format("2006-01-02"), actualDate.Format("2006-01-02"))
		}
	}

	// Second todo should have creation date Jan 12
	if todos[1].CreationDate() == nil {
		t.Error("expected second todo to have a creation date")
	} else {
		expectedDate := time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)
		actualDate := time.Date(todos[1].CreationDate().Year(), todos[1].CreationDate().Month(), todos[1].CreationDate().Day(), 0, 0, 0, 0, time.UTC)
		if !actualDate.Equal(expectedDate) {
			t.Errorf("expected creation date %s, got %s", expectedDate.Format("2006-01-02"), actualDate.Format("2006-01-02"))
		}
	}

	// Third todo should not have creation date
	if todos[2].CreationDate() != nil {
		t.Error("expected third todo to not have a creation date")
	}
}

func TestStory014_NewTodosGetCreationDateSet(t *testing.T) {
	// Scenario: Set creation date to today when adding new todos

	input := ""
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode on DO_FIRST quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type a new todo
	for _, ch := range "New task" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Check that creation date was set to today
	written := source.writer.(*strings.Builder).String()
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(written, today) {
		t.Errorf("expected new todo to have today's creation date %s, got: %s", today, written)
	}
	if !strings.Contains(written, "(A) "+today+" New task") {
		t.Errorf("expected todo in format '(A) %s New task', got: %s", today, written)
	}
}

func TestStory014_DisplayCreationDatesInUI(t *testing.T) {
	// Scenario: Display creation dates consistently in the UI

	// Create a todo from 5 days ago
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)
	input := "(A) " + fiveDaysAgo.Format("2006-01-02") + " Task from five days ago"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode to see detailed view
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should display "5 days ago" in the Created column
	if !strings.Contains(view, "5 days ago") {
		t.Errorf("expected view to show '5 days ago' in Created column, got: %s", view)
	}
	// Should have Created column header
	if !strings.Contains(view, "Created") {
		t.Errorf("expected view to show Created column header, got: %s", view)
	}
}

func TestStory014_PreserveCreationDateOnToggle(t *testing.T) {
	// Scenario: Toggling completion preserves creation date

	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	input := "(A) " + threeDaysAgo.Format("2006-01-02") + " Task to toggle"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle completion with spacebar
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)

	// Verify creation date is preserved in written output
	written := source.writer.(*strings.Builder).String()

	// Should contain the creation date (3 days ago)
	creationDateStr := threeDaysAgo.Format("2006-01-02")
	if !strings.Contains(written, creationDateStr) {
		t.Errorf("expected creation date %s to be preserved, got: %s", creationDateStr, written)
	}

	// Should be in completed format: x COMPLETION_DATE CREATION_DATE (A) Description
	if !strings.Contains(written, "x") {
		t.Errorf("expected todo to be marked as completed, got: %s", written)
	}
}

func TestStory014_PreserveCreationDateOnMove(t *testing.T) {
	// Scenario: Moving between quadrants preserves creation date

	twoDaysAgo := time.Now().AddDate(0, 0, -2)
	input := "(A) " + twoDaysAgo.Format("2006-01-02") + " Task to move"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode on DO_FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Move to SCHEDULE quadrant (press 'm' then '2')
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Verify creation date is preserved in written output
	written := source.writer.(*strings.Builder).String()

	// Should contain the creation date (2 days ago)
	creationDateStr := twoDaysAgo.Format("2006-01-02")
	if !strings.Contains(written, creationDateStr) {
		t.Errorf("expected creation date %s to be preserved, got: %s", creationDateStr, written)
	}

	// Should have new priority B
	if !strings.Contains(written, "(B)") {
		t.Errorf("expected todo to have priority B, got: %s", written)
	}
}

func TestStory014_FriendlyDateFormatting(t *testing.T) {
	// Scenario: Display dates in friendly format (today, yesterday, N days ago)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	sevenDaysAgo := today.AddDate(0, 0, -7)

	input := "(A) " + today.Format("2006-01-02") + " Task created today\n" +
		"(B) " + yesterday.Format("2006-01-02") + " Task created yesterday\n" +
		"(C) " + sevenDaysAgo.Format("2006-01-02") + " Task from a week ago"

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode on DO_FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show "today" in Created column
	if !strings.Contains(view, "today") {
		t.Errorf("expected view to show 'today' in Created column, got: %s", view)
	}

	// Switch to SCHEDULE quadrant to see "yesterday"
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	if !strings.Contains(view, "yesterday") {
		t.Errorf("expected view to show 'yesterday' in Created column, got: %s", view)
	}

	// Switch to DELEGATE quadrant to see "7 days ago"
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	if !strings.Contains(view, "7 days ago") {
		t.Errorf("expected view to show '7 days ago' in Created column, got: %s", view)
	}
}

func TestStory014_HandleTodosWithoutCreationDate(t *testing.T) {
	// Scenario: Application gracefully handles todos without creation dates

	input := `(A) 2026-01-10 Task with date
(B) Task without date
(C) 2026-01-05 Another task with date`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode to view todos
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should display the task with date
	if !strings.Contains(view, "Task with date") {
		t.Error("expected view to show task with date")
	}

	// Switch to SCHEDULE to see task without date
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	// Should display the task without date (no date info shown)
	if !strings.Contains(view, "Task without date") {
		t.Error("expected view to show task without date")
	}
}

func TestStory014_ParseCompletedTodoWithCreationDate(t *testing.T) {
	// Scenario: Parse completed todos with both completion and creation dates

	input := "x 2026-01-15 2026-01-10 (A) Completed task"

	todos, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}

	todo := todos[0]

	// Should be completed
	if !todo.IsCompleted() {
		t.Error("expected todo to be completed")
	}

	// Should have completion date Jan 15
	if todo.CompletionDate() == nil {
		t.Error("expected todo to have completion date")
	} else {
		expectedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		actualDate := time.Date(todo.CompletionDate().Year(), todo.CompletionDate().Month(), todo.CompletionDate().Day(), 0, 0, 0, 0, time.UTC)
		if !actualDate.Equal(expectedDate) {
			t.Errorf("expected completion date %s, got %s", expectedDate.Format("2006-01-02"), actualDate.Format("2006-01-02"))
		}
	}

	// Should have creation date Jan 10
	if todo.CreationDate() == nil {
		t.Error("expected todo to have creation date")
	} else {
		expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
		actualDate := time.Date(todo.CreationDate().Year(), todo.CreationDate().Month(), todo.CreationDate().Day(), 0, 0, 0, 0, time.UTC)
		if !actualDate.Equal(expectedDate) {
			t.Errorf("expected creation date %s, got %s", expectedDate.Format("2006-01-02"), actualDate.Format("2006-01-02"))
		}
	}

	// Should have correct description
	if todo.Description() != "Completed task" {
		t.Errorf("expected description 'Completed task', got %s", todo.Description())
	}
}
