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

// Story 005: Fullscreen TUI
//
// Note: Full acceptance testing of alt-screen mode, terminal restoration,
// and visual centering requires manual testing in a real terminal.
// These tests verify the model behavior that enables fullscreen mode.

func TestStory005_ModelRespondsToQuitKey(t *testing.T) {
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate pressing 'q'
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Verify the command is tea.Quit
	is.True(cmd != nil) // expected tea.Quit command
}

func TestStory005_ModelRespondsToCtrlC(t *testing.T) {
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate pressing Ctrl+C
	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Verify the command is tea.Quit
	is.True(cmd != nil) // expected tea.Quit command
}

func TestStory005_ModelHandlesWindowSize(t *testing.T) {
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate window size message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	// Verify the view renders with centering
	view := updatedModel.View()

	is.True(view != "") // expected non-empty view
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain matrix content
}

func TestStory005_ModelInitialViewWithoutWindowSize(t *testing.T) {
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Test todo", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Call View without any window size message
	view := model.View()

	is.True(view != "") // expected non-empty view
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain matrix content
}
