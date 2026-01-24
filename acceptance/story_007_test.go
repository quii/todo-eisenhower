package acceptance_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 007: Quadrant Focus Mode + Remove Emojis

func TestStory007_NoEmojisInQuadrantTitles(t *testing.T) {
	// Scenario: No emojis in quadrant titles
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	view := model.View()

	// Check that emojis are NOT present
	emojis := []string{"🔥", "📅", "👥", "🗑️", "📄"}
	for _, emoji := range emojis {
		is.True(!strings.Contains(stripANSI(view), emoji)) // view should not contain emoji
	}

	// Verify titles are present without emojis
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain 'DO FIRST'
	is.True(strings.Contains(stripANSI(view), "Schedule")) // expected view to contain 'SCHEDULE'
	is.True(strings.Contains(stripANSI(view), "Delegate")) // expected view to contain 'DELEGATE'
	is.True(strings.Contains(stripANSI(view), "Eliminate")) // expected view to contain 'ELIMINATE'
}

func TestStory007_FocusOnDoFirst(t *testing.T) {
	// Scenario: Focus on DO FIRST quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll(generateManyTodoObjects(20))
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "1" to focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show DO FIRST title prominently
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected focused view to contain 'DO FIRST' title

	// Should show file path header
	is.True(strings.Contains(stripANSI(view), "File: test.txt")) // expected focused view to contain file path header

	// Should show help text
	is.True(strings.Contains(stripANSI(view), "a to add")) // expected focused view to contain help text about adding tasks
	is.True(strings.Contains(stripANSI(view), "1-4 to jump")) // expected focused view to contain help text about jumping quadrants
	is.True(strings.Contains(stripANSI(view), "m to move")) // expected focused view to contain help text about moving todos
	is.True(strings.Contains(stripANSI(view), "ESC to return")) // expected focused view to contain help text about ESC

	// Should NOT show other quadrant titles
	is.True(!strings.Contains(stripANSI(view), "Schedule") && !strings.Contains(stripANSI(view), "Delegate") && !strings.Contains(stripANSI(view), "Eliminate")) // focused view should not contain other quadrant titles
}

func TestStory007_FocusOnSchedule(t *testing.T) {
	// Scenario: Focus on SCHEDULE quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Schedule task", todo.PriorityB),
		todo.New("Another schedule task", todo.PriorityB),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "2" to focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	is.True(strings.Contains(stripANSI(view), "Schedule")) // expected focused view to contain 'SCHEDULE' title
	is.True(strings.Contains(stripANSI(view), "Schedule task")) // expected focused view to contain schedule tasks
}

func TestStory007_FocusOnDelegate(t *testing.T) {
	// Scenario: Focus on DELEGATE quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Delegate task", todo.PriorityC),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "3" to focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	is.True(strings.Contains(stripANSI(view), "Delegate")) // expected focused view to contain 'DELEGATE' title
}

func TestStory007_FocusOnEliminate(t *testing.T) {
	// Scenario: Focus on ELIMINATE quadrant
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Eliminate task", todo.PriorityD),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "4" to focus on ELIMINATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	is.True(strings.Contains(stripANSI(view), "Eliminate")) // expected focused view to contain 'ELIMINATE' title
}

func TestStory007_ReturnToOverviewWithESC(t *testing.T) {
	// Scenario: Return to overview with ESC
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
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
	focusView := model.View()

	// Should be in focus mode (only DO FIRST visible)
	is.True(!strings.Contains(focusView, "Schedule")) // focus mode should not show other quadrants

	// Press ESC to return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)
	overviewView := model.View()

	// Should now show all quadrants
	is.True(strings.Contains(overviewView, "Do First")) // overview should contain DO FIRST
	is.True(strings.Contains(overviewView, "Schedule")) // overview should contain SCHEDULE
	is.True(strings.Contains(overviewView, "Delegate")) // overview should contain DELEGATE
	is.True(strings.Contains(overviewView, "Eliminate")) // overview should contain ELIMINATE

	// Should show overview help text (without ESC)
	is.True(strings.Contains(overviewView, "Press 1/2/3/4 to focus on a quadrant")) // overview should show help text

	// Should NOT show ESC instruction in overview
	is.True(!strings.Contains(overviewView, "Press ESC to return")) // overview should not show ESC instruction
}

func TestStory007_JumpBetweenQuadrantsInFocusMode(t *testing.T) {
	// Scenario: Jump between quadrants in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task A", todo.PriorityA),
		todo.New("Task B", todo.PriorityB),
		todo.New("Task C", todo.PriorityC),
		todo.New("Task D", todo.PriorityD),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST (1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view1 := model.View()
	is.True(strings.Contains(stripANSI(view1), "Do First")) // should show DO FIRST

	// Jump to SCHEDULE (2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view2 := model.View()
	is.True(strings.Contains(view2, "Schedule")) // should show SCHEDULE

	// Jump to ELIMINATE (4)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)
	view4 := model.View()
	is.True(strings.Contains(view4, "Eliminate")) // should show ELIMINATE

	// Jump back to DO FIRST (1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view1Again := model.View()
	is.True(strings.Contains(view1Again, "Do First")) // should show DO FIRST again
}

func TestStory007_EmptyQuadrantInFocusMode(t *testing.T) {
	// Scenario: Empty quadrant in focus mode
	is := is.New(t)

	// No Priority D or untagged tasks, so ELIMINATE will be empty
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Task A", todo.PriorityA),
		todo.New("Task B", todo.PriorityB),
		todo.New("Task C", todo.PriorityC),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on ELIMINATE (which is empty)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	is.True(strings.Contains(stripANSI(view), "Eliminate")) // should show ELIMINATE title
	is.True(strings.Contains(stripANSI(view), "(no tasks)")) // should show '(no tasks)' for empty quadrant
	is.True(strings.Contains(stripANSI(view), "m to move")) // should show help text even for empty quadrant
}

func TestStory007_DisplayLimitScalesInFocusMode(t *testing.T) {
	// Scenario: Display limit scales in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll(generateManyTodoObjects(50))
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Large terminal
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	largView := model.View()

	// Table should be rendered and handle many todos
	// With table-based rendering, all todos are in the table but scrollable
	is.True(strings.Contains(largView, "Task")) // large terminal should show task column header
	is.True(strings.Contains(stripANSI(largView), "Do First")) // large terminal should show quadrant title

	// Smaller terminal
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	smallView := model.View()

	// Table should still render in smaller terminal
	is.True(strings.Contains(smallView, "Task")) // small terminal should show task column header
	is.True(strings.Contains(smallView, "Do First")) // small terminal should show quadrant title
}
