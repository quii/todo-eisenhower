package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory017_ShowSummaryStatsForEachQuadrant(t *testing.T) {
	// Scenario: Show summary stats for each quadrant

	input := `(A) Task one
(A) Task two
x 2026-01-15 (A) Completed task
(B) Schedule task
(C) Delegate task
(D) Eliminate task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Should be in overview mode by default
	view := model.View()
	t.Logf("View output:\n%s", view)

	// Check for summary stats in each quadrant
	if !strings.Contains(view, "DO FIRST (3 tasks, 1 completed)") {
		t.Errorf("expected DO FIRST to show task count and completion stats, got view:\n%s", view)
	}
	if !strings.Contains(view, "SCHEDULE (1 task") || !strings.Contains(view, "0 completed") {
		t.Errorf("expected SCHEDULE to show task count and completion stats, got view:\n%s", view)
	}
	if !strings.Contains(view, "DELEGATE (1 task, 0 completed)") {
		t.Error("expected DELEGATE to show task count and completion stats")
	}
	if !strings.Contains(view, "ELIMINATE (1 task, 0 completed)") {
		t.Error("expected ELIMINATE to show task count and completion stats")
	}
}

func TestStory017_ShowTopNTodosAsSimpleList(t *testing.T) {
	// Scenario: Show top N todos as simple list

	input := `(A) First task
(A) Second task
(A) Third task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()
	t.Logf("View:\n%s", view)

	// Should show todos as simple bullet list
	if !strings.Contains(view, "• First task") {
		t.Error("expected first task to be shown with bullet point")
	}
	if !strings.Contains(view, "• Second task") {
		t.Error("expected second task to be shown with bullet point")
	}
	if !strings.Contains(view, "• Third task") {
		t.Error("expected third task to be shown with bullet point")
	}

	// Should NOT show table headers (this is overview, not focus mode)
	// Check for "Description" column header which would indicate table mode
	if strings.Contains(view, "Description") {
		t.Errorf("expected overview to use simple list, not table format (found 'Description' header)")
	}
}

func TestStory017_IndicateWhenThereAreMoreTodos(t *testing.T) {
	// Scenario: Indicate when there are more todos

	// Create more than 5 tasks to trigger the "... and N more" message
	input := `(A) Task 1
(A) Task 2
(A) Task 3
(A) Task 4
(A) Task 5
(A) Task 6
(A) Task 7
(A) Task 8`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()
	t.Logf("View:\n%s", view)

	// Should show a message indicating more todos
	if !strings.Contains(view, "and") || !strings.Contains(view, "more") {
		t.Errorf("expected message indicating there are more todos not shown")
	}

	// Should mention pressing the quadrant number to view all
	if !strings.Contains(view, "press 1 to view") {
		t.Errorf("expected hint to press 1 to view all todos in DO FIRST")
	}
}

func TestStory017_EmptyQuadrantShowsHelpfulMessage(t *testing.T) {
	// Scenario: Empty quadrant shows helpful message

	input := `(A) Only in DO FIRST`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Empty quadrants should show stats
	if !strings.Contains(view, "SCHEDULE (0 tasks") {
		t.Error("expected SCHEDULE to show 0 tasks in stats")
	}
	if !strings.Contains(view, "DELEGATE (0 tasks") {
		t.Error("expected DELEGATE to show 0 tasks in stats")
	}
	if !strings.Contains(view, "ELIMINATE (0 tasks") {
		t.Error("expected ELIMINATE to show 0 tasks in stats")
	}

	// Empty quadrants should show "(no tasks)"
	// Note: There will be multiple instances since we have multiple empty quadrants
	if !strings.Contains(view, "(no tasks)") {
		t.Error("expected empty quadrants to show '(no tasks)' message")
	}
}

func TestStory017_AllCompletedTodosShowsInStats(t *testing.T) {
	// Scenario: All completed todos shows in stats

	input := `x 2026-01-15 (A) Completed one
x 2026-01-15 (A) Completed two
x 2026-01-15 (A) Completed three`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show all tasks completed
	if !strings.Contains(view, "DO FIRST (3 tasks, 3 completed)") {
		t.Error("expected DO FIRST to show 3 tasks, 3 completed")
	}
}

func TestStory017_CompletedTodosShownWithVisualIndicator(t *testing.T) {
	// Scenario: Completed todos shown with visual indicator

	input := `x 2026-01-15 (A) Completed task
(A) Active task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Completed todos should have checkmark indicator
	if !strings.Contains(view, "✓") {
		t.Error("expected completed todos to show ✓ indicator")
	}

	// Active todos should have bullet point
	if !strings.Contains(view, "• Active task") {
		t.Error("expected active todos to show • bullet point")
	}
}

func TestStory017_QuadrantLayoutPreserved(t *testing.T) {
	// Scenario: Quadrant layout preserved

	input := `(A) DO FIRST task
(B) SCHEDULE task
(C) DELEGATE task
(D) ELIMINATE task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// All four quadrant titles should be present
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected DO FIRST quadrant to be shown")
	}
	if !strings.Contains(view, "SCHEDULE") {
		t.Error("expected SCHEDULE quadrant to be shown")
	}
	if !strings.Contains(view, "DELEGATE") {
		t.Error("expected DELEGATE quadrant to be shown")
	}
	if !strings.Contains(view, "ELIMINATE") {
		t.Error("expected ELIMINATE quadrant to be shown")
	}

	// Should have visual separation (borders)
	if !strings.Contains(view, "─") && !strings.Contains(view, "│") {
		t.Error("expected quadrants to have visual separation with borders")
	}
}

func TestStory017_NoTagsOrDatesInOverview(t *testing.T) {
	// Scenario: Overview shows simple descriptions without tags or dates

	input := `(A) 2026-01-10 Review code +WebApp @computer
x 2026-01-15 (A) 2026-01-12 Deploy feature +WebApp @terminal`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the task descriptions
	if !strings.Contains(view, "Review code") {
		t.Error("expected to see task description")
	}
	if !strings.Contains(view, "Deploy feature") {
		t.Error("expected to see completed task description")
	}

	// Tags should still be visible in the description (colorized)
	// but not in separate columns like in focus mode
	if !strings.Contains(view, "+WebApp") {
		t.Error("expected tags to still appear in description")
	}

	// Should NOT have "Description" column header which indicates table mode
	// Note: "Projects:" and "Contexts:" tag inventory is OK at the bottom
	if strings.Contains(view, "Description") {
		t.Error("expected overview to not show table column headers like 'Description'")
	}
}
