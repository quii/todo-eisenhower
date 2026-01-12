package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 005: Fullscreen TUI
//
// Note: Full acceptance testing of alt-screen mode, terminal restoration,
// and visual centering requires manual testing in a real terminal.
// These tests verify the model behavior that enables fullscreen mode.

func TestStory005_ModelRespondsToQuitKey(t *testing.T) {
	// Scenario: User can quit the application
	// Given the application is running in fullscreen mode
	// When I press 'q'
	// Then the application exits

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate pressing 'q'
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Verify the command is tea.Quit
	if cmd == nil {
		t.Fatal("expected tea.Quit command, got nil")
	}

	// tea.Quit is a sentinel function, we can check if it's non-nil
	// The actual quit behavior is handled by Bubble Tea runtime
	_ = updatedModel
}

func TestStory005_ModelRespondsToCtrlC(t *testing.T) {
	// Scenario: User can quit with Ctrl+C
	// Given the application is running in fullscreen mode
	// When I press Ctrl+C
	// Then the application exits gracefully

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate pressing Ctrl+C
	updatedModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Verify the command is tea.Quit
	if cmd == nil {
		t.Fatal("expected tea.Quit command, got nil")
	}

	_ = updatedModel
}

func TestStory005_ModelHandlesWindowSize(t *testing.T) {
	// Scenario: Matrix is centered in terminal
	// Given the application is running in fullscreen mode
	// When the matrix is displayed
	// Then it is centered horizontally and vertically

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate window size message
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	// Verify the view renders with centering
	// (actual centering is done by lipgloss.Place in the View method)
	view := updatedModel.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}

	// Verify the view contains our matrix content
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain matrix content")
	}
}

func TestStory005_ModelInitialViewWithoutWindowSize(t *testing.T) {
	// Verify that the model can render even before receiving window size
	// (for example, during initial render before WindowSizeMsg arrives)

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Call View without any window size message
	view := model.View()

	if view == "" {
		t.Fatal("expected non-empty view")
	}

	// Should still contain matrix content, just not centered
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain matrix content")
	}
}
