package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory012_MoveFromDoFirstToSchedule(t *testing.T) {
	// Scenario: Move todo from DO FIRST to SCHEDULE

	input := `(A) Review quarterly goals`

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

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press Shift+2 (@) to move to SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority B
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(B) Review quarterly goals") {
		t.Errorf("expected file to contain '(B) Review quarterly goals', got: %s", written)
	}
	if strings.Contains(written, "(A)") {
		t.Errorf("expected priority A to be changed to B, got: %s", written)
	}
}

func TestStory012_MoveFromDelegateToDoFirst(t *testing.T) {
	// Scenario: Move todo from DELEGATE to DO FIRST

	input := `(C) Update documentation`

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

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Press Shift+1 (!) to move to DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority A
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Update documentation") {
		t.Errorf("expected file to contain '(A) Update documentation', got: %s", written)
	}
}

func TestStory012_MoveToEliminate(t *testing.T) {
	// Scenario: Move todo to ELIMINATE (priority D)

	input := `(B) Optional feature idea`

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

	// Focus on SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press Shift+4 ($) to move to ELIMINATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'$'}})
	model = updatedModel.(ui.Model)

	// Check file was updated with priority D
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(D) Optional feature idea") {
		t.Errorf("expected file to contain '(D) Optional feature idea', got: %s", written)
	}
}

func TestStory012_MovingTodoAdjustsSelection(t *testing.T) {
	// Scenario: Moving todo adjusts selection

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

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Navigate to second task
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Move second task to SCHEDULE (press Shift+2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	model = updatedModel.(ui.Model)

	// Check file shows second task moved to priority B
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(B) Second task") {
		t.Errorf("expected second task to have priority B, got: %s", written)
	}

	// First and third tasks should still be priority A
	lines := strings.Split(strings.TrimSpace(written), "\n")
	priorityACount := 0
	for _, line := range lines {
		if strings.Contains(line, "(A)") {
			priorityACount++
		}
	}
	if priorityACount != 2 {
		t.Errorf("expected 2 tasks with priority A, got %d", priorityACount)
	}
}

func TestStory012_MovingLastTodoReturnsToOverview(t *testing.T) {
	// Scenario: Moving last todo in quadrant returns to overview

	input := `(C) Only delegate task`

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

	// Focus on DELEGATE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)

	// Move to DO FIRST (press Shift+1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model = updatedModel.(ui.Model)

	// Should return to overview mode
	view := model.View()
	if !strings.Contains(view, "DO FIRST") && !strings.Contains(view, "SCHEDULE") &&
	   !strings.Contains(view, "DELEGATE") && !strings.Contains(view, "ELIMINATE") {
		t.Error("expected to return to overview mode showing all quadrants")
	}

	// Verify the moved todo is in the file with priority A
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Only delegate task") {
		t.Errorf("expected todo to be moved to priority A, got: %s", written)
	}
}

func TestStory012_PreservesTagsAndCompletion(t *testing.T) {
	// Scenario: Moving todo preserves tags and completion status

	input := `x 2025-01-10 (A) Fix bug +WebApp @computer`

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

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Move to DELEGATE (press Shift+3)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'#'}})
	model = updatedModel.(ui.Model)

	// Check file preserves everything
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "x ") {
		t.Error("expected completion marker to be preserved")
	}
	if !strings.Contains(written, "(C)") {
		t.Error("expected priority to change to C")
	}
	if !strings.Contains(written, "+WebApp") {
		t.Error("expected project tag to be preserved")
	}
	if !strings.Contains(written, "@computer") {
		t.Error("expected context tag to be preserved")
	}
	if !strings.Contains(written, "Fix bug") {
		t.Error("expected description to be preserved")
	}
}

func TestStory012_PressingCurrentQuadrantDoesNothing(t *testing.T) {
	// Scenario: Pressing current quadrant number does nothing

	input := `(B) Plan sprint`

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

	// Focus on SCHEDULE (quadrant 2)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press '2' again (same quadrant) - should be a no-op
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Verify we're still viewing SCHEDULE (not moved to overview or other quadrant)
	view := model.View()
	if !strings.Contains(view, "SCHEDULE") {
		t.Error("should still be viewing SCHEDULE quadrant")
	}
	if !strings.Contains(view, "Plan sprint") {
		t.Error("should still show the todo")
	}

	// Press Shift+2 (@) again - should be a no-op since todo is already priority B
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'@'}})
	model = updatedModel.(ui.Model)

	// Verify todo is still in SCHEDULE with priority B
	view = model.View()
	if !strings.Contains(view, "SCHEDULE") {
		t.Error("should still be viewing SCHEDULE quadrant after Shift+2")
	}
	if !strings.Contains(view, "Plan sprint") {
		t.Error("should still show the todo after Shift+2")
	}

	// Now actually move it to verify moving works (press Shift+1)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'!'}})
	model = updatedModel.(ui.Model)

	// Check that it was moved to priority A
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Plan sprint") {
		t.Errorf("after pressing Shift+1, todo should be moved to priority A, got: %s", written)
	}
}

func TestStory012_NumberKeysStillFocusInOverview(t *testing.T) {
	// Scenario: Number keys still focus quadrants in overview mode

	input := `(A) Task one
(B) Task two`

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

	// In overview mode, press '1' to focus
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should be in focus mode on DO FIRST
	view := model.View()
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected to focus on DO FIRST quadrant")
	}

	// Should show focus mode help text
	if !strings.Contains(view, "Press a to add") {
		t.Error("expected focus mode help text")
	}
}
