package acceptance_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory017_ShowSummaryStatsForEachQuadrant(t *testing.T) {
	// Scenario: Show summary stats for each quadrant
	is := is.New(t)

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
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Should be in overview mode by default
	view := model.View()
	t.Logf("View output:\n%s", view)

	// Check for summary stats in each quadrant
	is.True(strings.Contains(stripANSI(view), "Do First   3 tasks · 1 completed")) // expected DO FIRST to show task count and completion stats
	is.True(strings.Contains(stripANSI(view), "Schedule   1 task · 0 completed")) // expected SCHEDULE to show task count and completion stats
	is.True(strings.Contains(stripANSI(view), "Delegate   1 task · 0 completed")) // expected DELEGATE to show task count and completion stats
	is.True(strings.Contains(stripANSI(view), "Eliminate   1 task · 0 completed")) // expected ELIMINATE to show task count and completion stats
}

func TestStory017_ShowTopNTodosAsSimpleList(t *testing.T) {
	// Scenario: Show top N todos as simple list
	is := is.New(t)

	input := `(A) First task
(A) Second task
(A) Third task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()
	t.Logf("View:\n%s", view)

	// Should show todos as simple bullet list
	is.True(strings.Contains(stripANSI(view), "• First task")) // expected first task to be shown with bullet point
	is.True(strings.Contains(stripANSI(view), "• Second task")) // expected second task to be shown with bullet point
	is.True(strings.Contains(stripANSI(view), "• Third task")) // expected third task to be shown with bullet point

	// Should NOT show table headers (this is overview, not focus mode)
	// Check for "Description" column header which would indicate table mode
	is.True(!strings.Contains(stripANSI(view), "Description")) // expected overview to use simple list, not table format
}

func TestStory017_IndicateWhenThereAreMoreTodos(t *testing.T) {
	// Scenario: Indicate when there are more todos
	is := is.New(t)

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
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()
	t.Logf("View:\n%s", view)

	// Should show a message indicating more todos
	is.True(strings.Contains(stripANSI(view), "and") && strings.Contains(stripANSI(view), "more")) // expected message indicating there are more todos not shown

	// Should mention pressing the quadrant number to view all
	is.True(strings.Contains(stripANSI(view), "press 1 to view")) // expected hint to press 1 to view all todos in DO FIRST
}

func TestStory017_EmptyQuadrantShowsHelpfulMessage(t *testing.T) {
	// Scenario: Empty quadrant shows helpful message
	is := is.New(t)

	input := `(A) Only in DO FIRST`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Empty quadrants should show stats
	is.True(strings.Contains(stripANSI(view), "Schedule   0 tasks")) // expected SCHEDULE to show 0 tasks in stats
	is.True(strings.Contains(stripANSI(view), "Delegate   0 tasks")) // expected DELEGATE to show 0 tasks in stats
	is.True(strings.Contains(stripANSI(view), "Eliminate   0 tasks")) // expected ELIMINATE to show 0 tasks in stats

	// Empty quadrants should show "(no tasks)"
	// Note: There will be multiple instances since we have multiple empty quadrants
	is.True(strings.Contains(stripANSI(view), "(no tasks)")) // expected empty quadrants to show '(no tasks)' message
}

func TestStory017_AllCompletedTodosShowsInStats(t *testing.T) {
	// Scenario: All completed todos shows in stats
	is := is.New(t)

	input := `x 2026-01-15 (A) Completed one
x 2026-01-15 (A) Completed two
x 2026-01-15 (A) Completed three`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show all tasks completed
	is.True(strings.Contains(stripANSI(view), "Do First   3 tasks · 3 completed")) // expected DO FIRST to show 3 tasks, 3 completed
}

func TestStory017_CompletedTodosShownWithVisualIndicator(t *testing.T) {
	// Scenario: Completed todos shown with visual indicator
	is := is.New(t)

	input := `x 2026-01-15 (A) Completed task
(A) Active task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Completed todos should have checkmark indicator
	is.True(strings.Contains(stripANSI(view), "✓")) // expected completed todos to show ✓ indicator

	// Active todos should have bullet point
	is.True(strings.Contains(stripANSI(view), "• Active task")) // expected active todos to show • bullet point
}

func TestStory017_QuadrantLayoutPreserved(t *testing.T) {
	// Scenario: Quadrant layout preserved
	is := is.New(t)

	input := `(A) DO FIRST task
(B) SCHEDULE task
(C) DELEGATE task
(D) ELIMINATE task`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// All four quadrant titles should be present
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected DO FIRST quadrant to be shown
	is.True(strings.Contains(stripANSI(view), "Schedule")) // expected SCHEDULE quadrant to be shown
	is.True(strings.Contains(stripANSI(view), "Delegate")) // expected DELEGATE quadrant to be shown
	is.True(strings.Contains(stripANSI(view), "Eliminate")) // expected ELIMINATE quadrant to be shown

	// Should have visual separation (borders)
	is.True(strings.Contains(stripANSI(view), "─") && strings.Contains(stripANSI(view), "│")) // expected quadrants to have visual separation with borders
}

func TestStory017_NoTagsOrDatesInOverview(t *testing.T) {
	// Scenario: Overview shows simple descriptions without tags or dates
	is := is.New(t)

	input := `(A) 2026-01-10 Review code +WebApp @computer
x 2026-01-15 (A) 2026-01-12 Deploy feature +WebApp @terminal`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetSource(source).SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the task descriptions
	is.True(strings.Contains(stripANSI(view), "Review code")) // expected to see task description
	is.True(strings.Contains(stripANSI(view), "Deploy feature")) // expected to see completed task description

	// Tags should still be visible in the description (colorized)
	// but not in separate columns like in focus mode
	is.True(strings.Contains(stripANSI(view), "+WebApp")) // expected tags to still appear in description

	// Should NOT have "Description" column header which indicates table mode
	// Note: "Projects:" and "Contexts:" tag inventory is OK at the bottom
	is.True(!strings.Contains(stripANSI(view), "Description")) // expected overview to not show table column headers like 'Description'
}
