package acceptance_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 019: Digital Inventory Dashboard

func TestStory019_ViewInventoryDashboard(t *testing.T) {
	// Scenario: View inventory dashboard with active todos
	is := is.New(t)

	tenDaysAgo := time.Now().AddDate(0, 0, -10)
	fiveDaysAgo := time.Now().AddDate(0, 0, -5)

	input := "(A) " + tenDaysAgo.Format("2006-01-02") + " Urgent task +project1 @people\n" +
		"(A) " + fiveDaysAgo.Format("2006-01-02") + " Another urgent +project1 @architecture\n" +
		"(B) " + fiveDaysAgo.Format("2006-01-02") + " Important task +project2 @people"

	repository := memory.NewRepository(input)

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

	input := "(A) " + veryOld.Format("2006-01-02") + " Very old task\n" +
		"(B) " + stale.Format("2006-01-02") + " Stale task"

	repository := memory.NewRepository(input)

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
	input := "x " + threeDaysAgo.Format("2006-01-02") + " " + tenDaysAgo.Format("2006-01-02") + " (A) Completed 1\n" +
		"x " + threeDaysAgo.Format("2006-01-02") + " " + tenDaysAgo.Format("2006-01-02") + " (A) Completed 2\n" +
		"x " + twoDaysAgo.Format("2006-01-02") + " " + tenDaysAgo.Format("2006-01-02") + " (A) Completed 3\n" +
		"(A) " + threeDaysAgo.Format("2006-01-02") + " New task 1\n" +
		"(A) " + threeDaysAgo.Format("2006-01-02") + " New task 2\n" +
		"(A) " + threeDaysAgo.Format("2006-01-02") + " New task 3\n" +
		"(A) " + twoDaysAgo.Format("2006-01-02") + " New task 4\n" +
		"(B) " + twoDaysAgo.Format("2006-01-02") + " New task 5\n" +
		"(B) " + twoDaysAgo.Format("2006-01-02") + " New task 6\n" +
		"(B) " + twoDaysAgo.Format("2006-01-02") + " New task 7\n" +
		"(C) " + threeDaysAgo.Format("2006-01-02") + " New task 8\n" +
		"(C) " + threeDaysAgo.Format("2006-01-02") + " New task 9\n" +
		"(C) " + threeDaysAgo.Format("2006-01-02") + " New task 10\n" +
		"(C) " + threeDaysAgo.Format("2006-01-02") + " New task 11"

	repository := memory.NewRepository(input)

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
	input := "(A) " + eighteenDaysAgo.Format("2006-01-02") + " Task 1 @people +WebApp\n" +
		"(A) " + eighteenDaysAgo.Format("2006-01-02") + " Task 2 @people +WebApp\n" +
		"(A) " + eighteenDaysAgo.Format("2006-01-02") + " Task 3 @people +WebApp\n" +
		"(A) " + eighteenDaysAgo.Format("2006-01-02") + " Task 4 @people +WebApp\n" +
		"(B) " + eighteenDaysAgo.Format("2006-01-02") + " Task 5 @people\n" +
		"(B) " + eighteenDaysAgo.Format("2006-01-02") + " Task 6 @people\n" +
		"(C) " + tenDaysAgo.Format("2006-01-02") + " Task 7 @architecture +Mobile"

	repository := memory.NewRepository(input)

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

	input := "(A) Task 1"
	repository := memory.NewRepository(input)

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

	input := "(A) Task 1"
	repository := memory.NewRepository(input)

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
	input := "(A) " + tenDaysAgo.Format("2006-01-02") + " Task with date 1 +project1 @context1\n" +
		"(A) " + tenDaysAgo.Format("2006-01-02") + " Task with date 2 +project1 @context1\n" +
		"(A) " + tenDaysAgo.Format("2006-01-02") + " Task with date 3 +project1 @context1\n" +
		"(A) Task without date 1 +project2 @context2\n" +
		"(A) Task without date 2 +project2 @context2"

	repository := memory.NewRepository(input)

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

	input := "(A) Task 1\n(A) Task 2"
	repository := memory.NewRepository(input)

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
