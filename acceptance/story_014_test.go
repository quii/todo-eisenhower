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

func TestStory014_ParseAndPreserveCreationDates(t *testing.T) {
	is := is.New(t)
	// Scenario: Parse and preserve existing creation dates from file

	input := `(A) 2026-01-10 Task created on Jan 10
(B) 2026-01-12 Task created on Jan 12
(C) Task without creation date`

	todos, err := todotxt.Unmarshal(strings.NewReader(input))
	is.NoErr(err)

	is.Equal(len(todos), 3) // expected 3 todos

	// First todo should have creation date Jan 10
	is.True(todos[0].CreationDate() != nil) // expected first todo to have a creation date

	expectedDate := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(todos[0].CreationDate().Year(), todos[0].CreationDate().Month(), todos[0].CreationDate().Day(), 0, 0, 0, 0, time.UTC)
	is.True(actualDate.Equal(expectedDate)) // expected creation date 2026-01-10

	// Second todo should have creation date Jan 12
	is.True(todos[1].CreationDate() != nil) // expected second todo to have a creation date

	expectedDate = time.Date(2026, 1, 12, 0, 0, 0, 0, time.UTC)
	actualDate = time.Date(todos[1].CreationDate().Year(), todos[1].CreationDate().Month(), todos[1].CreationDate().Day(), 0, 0, 0, 0, time.UTC)
	is.True(actualDate.Equal(expectedDate)) // expected creation date 2026-01-12

	// Third todo should not have creation date
	is.True(todos[2].CreationDate() == nil) // expected third todo to not have a creation date
}

func TestStory014_NewTodosGetCreationDateSet(t *testing.T) {
	is := is.New(t)
	// Scenario: Set creation date to today when adding new todos

	repository := memory.NewRepository()

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	_ = updatedModel.(ui.Model)

	// Check that creation date was set to today
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Description(), "New task")
	is.Equal(savedTodos[0].Priority(), todo.PriorityA)

	today := time.Now().Format("2006-01-02")
	is.True(savedTodos[0].CreationDate() != nil)
	is.Equal(savedTodos[0].CreationDate().Format("2006-01-02"), today)
}

func TestStory014_DisplayCreationDatesInUI(t *testing.T) {
	is := is.New(t)
	// Scenario: Display creation dates consistently in the UI

	// Create a todo from 5 days ago
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Task from five days ago", todo.PriorityA, &fiveDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode to see detailed view
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should display "5 days ago" in the Created column
	is.True(strings.Contains(stripANSI(view), "5 days ago")) // expected view to show '5 days ago' in Created column

	// Should have Created column header
	is.True(strings.Contains(stripANSI(view), "Created")) // expected view to show Created column header
}

func TestStory014_PreserveCreationDateOnToggle(t *testing.T) {
	is := is.New(t)
	// Scenario: Toggling completion preserves creation date

	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Task to toggle", todo.PriorityA, &threeDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Toggle completion with spacebar
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	_ = updatedModel.(ui.Model)

	// Verify creation date is preserved in written output
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.True(savedTodos[0].IsCompleted())

	// Should contain the creation date (3 days ago)
	creationDateStr := threeDaysAgo.Format("2006-01-02")
	is.True(savedTodos[0].CreationDate() != nil)
	is.Equal(savedTodos[0].CreationDate().Format("2006-01-02"), creationDateStr)
}

func TestStory014_PreserveCreationDateOnMove(t *testing.T) {
	is := is.New(t)
	// Scenario: Moving between quadrants preserves creation date

	twoDaysAgo := time.Now().AddDate(0, 0, -2)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Task to move", todo.PriorityA, &twoDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	_ = updatedModel.(ui.Model)

	// Verify creation date is preserved in written output
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1)
	is.Equal(savedTodos[0].Priority(), todo.PriorityB)

	// Should contain the creation date (2 days ago)
	creationDateStr := twoDaysAgo.Format("2006-01-02")
	is.True(savedTodos[0].CreationDate() != nil)
	is.Equal(savedTodos[0].CreationDate().Format("2006-01-02"), creationDateStr)
}

func TestStory014_FriendlyDateFormatting(t *testing.T) {
	is := is.New(t)
	// Scenario: Display dates in friendly format (today, yesterday, N days ago)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	sevenDaysAgo := today.AddDate(0, 0, -7)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Task created today", todo.PriorityA, &today),
		todo.NewWithCreationDate("Task created yesterday", todo.PriorityB, &yesterday),
		todo.NewWithCreationDate("Task from a week ago", todo.PriorityC, &sevenDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode on DO_FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show "today" in Created column
	is.True(strings.Contains(stripANSI(view), "today")) // expected view to show 'today' in Created column

	// Switch to SCHEDULE quadrant to see "yesterday"
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	is.True(strings.Contains(stripANSI(view), "yesterday")) // expected view to show 'yesterday' in Created column

	// Switch to DELEGATE quadrant to see "7 days ago"
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	is.True(strings.Contains(stripANSI(view), "7 days ago")) // expected view to show '7 days ago' in Created column
}

func TestStory014_HandleTodosWithoutCreationDate(t *testing.T) {
	is := is.New(t)
	// Scenario: Application gracefully handles todos without creation dates

	date1 := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Task with date", todo.PriorityA, &date1),
		todo.New("Task without date", todo.PriorityB),
		todo.NewWithCreationDate("Another task with date", todo.PriorityC, &date2),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	var updatedModel tea.Model
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = updatedModel.(ui.Model)

	// Enter focus mode to view todos
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should display the task with date
	is.True(strings.Contains(stripANSI(view), "Task with date")) // expected view to show task with date

	// Switch to SCHEDULE to see task without date
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	// Should display the task without date (no date info shown)
	is.True(strings.Contains(stripANSI(view), "Task without date")) // expected view to show task without date
}

func TestStory014_ParseCompletedTodoWithCreationDate(t *testing.T) {
	is := is.New(t)
	// Scenario: Parse completed todos with both completion and creation dates

	input := "x 2026-01-15 2026-01-10 (A) Completed task"

	todos, err := todotxt.Unmarshal(strings.NewReader(input))
	is.NoErr(err)

	is.Equal(len(todos), 1) // expected 1 todo

	todo := todos[0]

	// Should be completed
	is.True(todo.IsCompleted()) // expected todo to be completed

	// Should have completion date Jan 15
	is.True(todo.CompletionDate() != nil) // expected todo to have completion date

	expectedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	actualDate := time.Date(todo.CompletionDate().Year(), todo.CompletionDate().Month(), todo.CompletionDate().Day(), 0, 0, 0, 0, time.UTC)
	is.True(actualDate.Equal(expectedDate)) // expected completion date 2026-01-15

	// Should have creation date Jan 10
	is.True(todo.CreationDate() != nil) // expected todo to have creation date

	expectedDate = time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	actualDate = time.Date(todo.CreationDate().Year(), todo.CreationDate().Month(), todo.CreationDate().Day(), 0, 0, 0, 0, time.UTC)
	is.True(actualDate.Equal(expectedDate)) // expected creation date 2026-01-10

	// Should have correct description
	is.Equal(todo.Description(), "Completed task")
}
