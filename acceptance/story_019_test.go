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

// Story 019: Digital Inventory Dashboard

func TestStory019_ViewInventoryDashboard(t *testing.T) {
	// Scenario: View inventory dashboard with active todos
	is := is.New(t)

	tenDaysAgo := time.Now().AddDate(0, 0, -10)
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTagsAndDates("Urgent task", todo.PriorityA, &tenDaysAgo, []string{"project1"}, []string{"people"}),
		todo.NewWithTagsAndDates("Another urgent", todo.PriorityA, &fiveDaysAgo, []string{"project1"}, []string{"architecture"}),
		todo.NewWithTagsAndDates("Important task", todo.PriorityB, &fiveDaysAgo, []string{"project2"}, []string{"people"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 'i' to enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show dashboard title
	is.True(strings.Contains(stripANSI(view), "Digital Inventory Dashboard"))

	// Should show quadrant metrics
	is.True(strings.Contains(stripANSI(view), "Quadrant Metrics"))
	is.True(strings.Contains(stripANSI(view), "Do First"))
	is.True(strings.Contains(stripANSI(view), "Schedule"))

	// Should show total WIP count
	is.True(strings.Contains(stripANSI(view), "Total WIP"))

	// Should show throughput section
	is.True(strings.Contains(stripANSI(view), "Throughput"))

	// Should show project breakdown
	is.True(strings.Contains(stripANSI(view), "Project Breakdown"))
	is.True(strings.Contains(stripANSI(view), "project1"))

	// Should show context breakdown
	is.True(strings.Contains(stripANSI(view), "Context Breakdown"))
	is.True(strings.Contains(stripANSI(view), "people"))
}

func TestStory019_StaleItemsWarning(t *testing.T) {
	// Scenario: Dashboard shows stale items warning
	is := is.New(t)

	veryOld := time.Now().AddDate(0, 0, -25) // 25 days ago - VERY STALE
	stale := time.Now().AddDate(0, 0, -16)   // 16 days ago - STALE

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithCreationDate("Very old task", todo.PriorityA, &veryOld),
		todo.NewWithCreationDate("Stale task", todo.PriorityB, &stale),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show VERY STALE warning for >21 days
	is.True(strings.Contains(stripANSI(view), "VERY STALE"))

	// Should show STALE warning for >14 days
	is.True(strings.Contains(stripANSI(view), "STALE"))
}

func TestStory019_ThroughputCalculation(t *testing.T) {
	// Scenario: Dashboard calculates throughput
	is := is.New(t)

	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	twoDaysAgo := time.Now().AddDate(0, 0, -2)
	tenDaysAgo := time.Now().AddDate(0, 0, -10)

	// 3 todos completed in last 7 days
	// 11 todos added in last 7 days (we'll create several)
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewCompletedWithDates("Completed 1", todo.PriorityA, &threeDaysAgo, &tenDaysAgo),
		todo.NewCompletedWithDates("Completed 2", todo.PriorityA, &threeDaysAgo, &tenDaysAgo),
		todo.NewCompletedWithDates("Completed 3", todo.PriorityA, &twoDaysAgo, &tenDaysAgo),
		todo.NewWithCreationDate("New task 1", todo.PriorityA, &threeDaysAgo),
		todo.NewWithCreationDate("New task 2", todo.PriorityA, &threeDaysAgo),
		todo.NewWithCreationDate("New task 3", todo.PriorityA, &threeDaysAgo),
		todo.NewWithCreationDate("New task 4", todo.PriorityA, &twoDaysAgo),
		todo.NewWithCreationDate("New task 5", todo.PriorityB, &twoDaysAgo),
		todo.NewWithCreationDate("New task 6", todo.PriorityB, &twoDaysAgo),
		todo.NewWithCreationDate("New task 7", todo.PriorityB, &twoDaysAgo),
		todo.NewWithCreationDate("New task 8", todo.PriorityC, &threeDaysAgo),
		todo.NewWithCreationDate("New task 9", todo.PriorityC, &threeDaysAgo),
		todo.NewWithCreationDate("New task 10", todo.PriorityC, &threeDaysAgo),
		todo.NewWithCreationDate("New task 11", todo.PriorityC, &threeDaysAgo),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show completed count
	is.True(strings.Contains(stripANSI(view), "Completed: 3 items"))

	// Should show added count
	is.True(strings.Contains(stripANSI(view), "Added: 11 items"))

	// Should show warning about adding faster than completing
	is.True(strings.Contains(stripANSI(view), "Adding faster than completing"))
}

func TestStory019_TagBreakdowns(t *testing.T) {
	// Scenario: Dashboard groups by context and project tags
	is := is.New(t)

	eighteenDaysAgo := time.Now().AddDate(0, 0, -18)
	tenDaysAgo := time.Now().AddDate(0, 0, -10)

	// 6 todos with @people context
	// 4 todos with +WebApp project
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTagsAndDates("Task 1", todo.PriorityA, &eighteenDaysAgo, []string{"WebApp"}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 2", todo.PriorityA, &eighteenDaysAgo, []string{"WebApp"}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 3", todo.PriorityA, &eighteenDaysAgo, []string{"WebApp"}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 4", todo.PriorityA, &eighteenDaysAgo, []string{"WebApp"}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 5", todo.PriorityB, &eighteenDaysAgo, []string{}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 6", todo.PriorityB, &eighteenDaysAgo, []string{}, []string{"people"}),
		todo.NewWithTagsAndDates("Task 7", todo.PriorityC, &tenDaysAgo, []string{"Mobile"}, []string{"architecture"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show context breakdown with @people
	is.True(strings.Contains(stripANSI(view), "people"))
	is.True(strings.Contains(stripANSI(view), "6")) // count in table

	// Should show project breakdown with +WebApp
	is.True(strings.Contains(stripANSI(view), "WebApp"))
	is.True(strings.Contains(stripANSI(view), "4")) // count in table

	// Should show average age column header
	is.True(strings.Contains(stripANSI(view), "Avg Age"))
}

func TestStory019_ExitInventoryMode(t *testing.T) {
	// Scenario: Press 'i' again to exit dashboard
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task 1", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should be in inventory mode
	is.True(strings.Contains(stripANSI(view), "Digital Inventory Dashboard"))

	// Press 'i' again to exit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view = model.View()

	// Should return to overview mode (no longer showing dashboard)
	is.True(!strings.Contains(stripANSI(view), "Digital Inventory Dashboard"))
	is.True(strings.Contains(stripANSI(view), "Do First") || strings.Contains(stripANSI(view), "DO_FIRST"))
}

func TestStory019_ExitInventoryModeWithESC(t *testing.T) {
	// Scenario: Press ESC to exit dashboard
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task 1", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)

	// Press ESC to exit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should return to overview mode
	is.True(!strings.Contains(stripANSI(view), "Digital Inventory Dashboard"))
}

func TestStory019_ExcludeTodosWithoutCreationDates(t *testing.T) {
	// Scenario: Todos without creation dates are excluded from age metrics
	is := is.New(t)

	tenDaysAgo := time.Now().AddDate(0, 0, -10)

	// 5 active todos in Do First: 3 with dates, 2 without
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTagsAndDates("Task with date 1", todo.PriorityA, &tenDaysAgo, []string{"project1"}, []string{"context1"}),
		todo.NewWithTagsAndDates("Task with date 2", todo.PriorityA, &tenDaysAgo, []string{"project1"}, []string{"context1"}),
		todo.NewWithTagsAndDates("Task with date 3", todo.PriorityA, &tenDaysAgo, []string{"project1"}, []string{"context1"}),
		todo.NewWithTags("Task without date 1", todo.PriorityA, []string{"project2"}, []string{"context2"}),
		todo.NewWithTags("Task without date 2", todo.PriorityA, []string{"project2"}, []string{"context2"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter inventory mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should show 5 active todos in Do First quadrant
	is.True(strings.Contains(stripANSI(view), "Do First"))
	is.True(strings.Contains(stripANSI(view), "5")) // active count in table

	// Should show project1 (from todos with dates) in breakdown
	is.True(strings.Contains(stripANSI(view), "project1"))

	// Should NOT show project2 (from todos without dates) in breakdown
	is.True(!strings.Contains(stripANSI(view), "project2"))

	// Should show context1 in breakdown
	is.True(strings.Contains(stripANSI(view), "context1"))

	// Should NOT show context2 in breakdown
	is.True(!strings.Contains(stripANSI(view), "context2"))
}

func TestStory019_InventoryOnlyInOverviewMode(t *testing.T) {
	// Scenario: 'i' key only works in overview mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task 1", todo.PriorityA),
		todo.New("Task 2", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Enter focus mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Try to press 'i' in focus mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Should NOT enter inventory mode (should stay in focus mode)
	is.True(!strings.Contains(stripANSI(view), "Digital Inventory Dashboard"))
}
