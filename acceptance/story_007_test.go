package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 007: Quadrant Focus Mode + Remove Emojis

func TestStory007_NoEmojisInQuadrantTitles(t *testing.T) {
	// Scenario: No emojis in quadrant titles

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	view := model.View()

	// Check that emojis are NOT present
	emojis := []string{"🔥", "📅", "👥", "🗑️", "📄"}
	for _, emoji := range emojis {
		if strings.Contains(view, emoji) {
			t.Errorf("view should not contain emoji %q", emoji)
		}
	}

	// Verify titles are present without emojis
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain 'DO FIRST'")
	}
	if !strings.Contains(view, "SCHEDULE") {
		t.Error("expected view to contain 'SCHEDULE'")
	}
	if !strings.Contains(view, "DELEGATE") {
		t.Error("expected view to contain 'DELEGATE'")
	}
	if !strings.Contains(view, "ELIMINATE") {
		t.Error("expected view to contain 'ELIMINATE'")
	}
}

func TestStory007_FocusOnDoFirst(t *testing.T) {
	// Scenario: Focus on DO FIRST quadrant

	input := generateManyTodos(20)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "1" to focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show DO FIRST title prominently
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected focused view to contain 'DO FIRST' title")
	}

	// Should show file path header
	if !strings.Contains(view, "File: test.txt") {
		t.Error("expected focused view to contain file path header")
	}

	// Should show help text
	if !strings.Contains(view, "Press a to add a task") {
		t.Error("expected focused view to contain help text about adding tasks")
	}
	if !strings.Contains(view, "Press 1/2/3/4 to focus on a quadrant") {
		t.Error("expected focused view to contain help text about focusing")
	}
	if !strings.Contains(view, "Press ESC to return") {
		t.Error("expected focused view to contain help text about ESC")
	}

	// Should NOT show other quadrant titles
	if strings.Contains(view, "SCHEDULE") || strings.Contains(view, "DELEGATE") || strings.Contains(view, "ELIMINATE") {
		t.Error("focused view should not contain other quadrant titles")
	}
}

func TestStory007_FocusOnSchedule(t *testing.T) {
	// Scenario: Focus on SCHEDULE quadrant

	input := "(B) Schedule task\n(B) Another schedule task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "2" to focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	if !strings.Contains(view, "SCHEDULE") {
		t.Error("expected focused view to contain 'SCHEDULE' title")
	}

	if !strings.Contains(view, "Schedule task") {
		t.Error("expected focused view to contain schedule tasks")
	}
}

func TestStory007_FocusOnDelegate(t *testing.T) {
	// Scenario: Focus on DELEGATE quadrant

	input := "(C) Delegate task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "3" to focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	if !strings.Contains(view, "DELEGATE") {
		t.Error("expected focused view to contain 'DELEGATE' title")
	}
}

func TestStory007_FocusOnEliminate(t *testing.T) {
	// Scenario: Focus on ELIMINATE quadrant

	input := "(D) Eliminate task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press "4" to focus on ELIMINATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	if !strings.Contains(view, "ELIMINATE") {
		t.Error("expected focused view to contain 'ELIMINATE' title")
	}
}

func TestStory007_ReturnToOverviewWithESC(t *testing.T) {
	// Scenario: Return to overview with ESC

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	focusView := model.View()

	// Should be in focus mode (only DO FIRST visible)
	if strings.Contains(focusView, "SCHEDULE") {
		t.Error("focus mode should not show other quadrants")
	}

	// Press ESC to return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)
	overviewView := model.View()

	// Should now show all quadrants
	if !strings.Contains(overviewView, "DO FIRST") {
		t.Error("overview should contain DO FIRST")
	}
	if !strings.Contains(overviewView, "SCHEDULE") {
		t.Error("overview should contain SCHEDULE")
	}
	if !strings.Contains(overviewView, "DELEGATE") {
		t.Error("overview should contain DELEGATE")
	}
	if !strings.Contains(overviewView, "ELIMINATE") {
		t.Error("overview should contain ELIMINATE")
	}

	// Should show overview help text (without ESC)
	if !strings.Contains(overviewView, "Press 1/2/3/4 to focus on a quadrant") {
		t.Error("overview should show help text")
	}
	// Should NOT show ESC instruction in overview
	if strings.Contains(overviewView, "Press ESC to return") {
		t.Error("overview should not show ESC instruction")
	}
}

func TestStory007_JumpBetweenQuadrantsInFocusMode(t *testing.T) {
	// Scenario: Jump between quadrants in focus mode

	input := "(A) Task A\n(B) Task B\n(C) Task C\n(D) Task D"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST (1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view1 := model.View()
	if !strings.Contains(view1, "DO FIRST") {
		t.Error("should show DO FIRST")
	}

	// Jump to SCHEDULE (2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	view2 := model.View()
	if !strings.Contains(view2, "SCHEDULE") || strings.Contains(view2, "DO FIRST") {
		t.Error("should show SCHEDULE, not DO FIRST")
	}

	// Jump to ELIMINATE (4)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)
	view4 := model.View()
	if !strings.Contains(view4, "ELIMINATE") {
		t.Error("should show ELIMINATE")
	}

	// Jump back to DO FIRST (1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view1Again := model.View()
	if !strings.Contains(view1Again, "DO FIRST") {
		t.Error("should show DO FIRST again")
	}
}

func TestStory007_EmptyQuadrantInFocusMode(t *testing.T) {
	// Scenario: Empty quadrant in focus mode

	// No Priority D or untagged tasks, so ELIMINATE will be empty
	input := "(A) Task A\n(B) Task B\n(C) Task C"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on ELIMINATE (which is empty)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	if !strings.Contains(view, "ELIMINATE") {
		t.Error("should show ELIMINATE title")
	}

	if !strings.Contains(view, "(no tasks)") {
		t.Error("should show '(no tasks)' for empty quadrant")
	}

	if !strings.Contains(view, "Press 1/2/3/4 to focus on a quadrant") {
		t.Error("should show help text even for empty quadrant")
	}
}

func TestStory007_DisplayLimitScalesInFocusMode(t *testing.T) {
	// Scenario: Display limit scales in focus mode

	input := generateManyTodos(50)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Large terminal
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 50})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	largView := model.View()

	// Should show many todos (50 - 9 reserved = 41 displayable)
	// With 50 todos, should show "... and 9 more"
	if !strings.Contains(largView, "... and 9 more") {
		t.Error("large terminal should show '... and 9 more'")
	}

	// Smaller terminal
	updatedModel, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	smallView := model.View()

	// Should show fewer todos (30 - 9 = 21 displayable)
	// With 50 todos, should show "... and 29 more"
	if !strings.Contains(smallView, "... and 29 more") {
		t.Error("small terminal should show '... and 29 more'")
	}
}
