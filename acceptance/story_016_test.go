package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory016_EnterMoveModeWithMKey(t *testing.T) {
	// Scenario: Enter move mode with 'm' key

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

	// Press 'm' to enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should see move mode overlay
	view := model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to see move mode overlay title")
	}
	if !strings.Contains(view, "1. DO FIRST") {
		t.Error("expected to see DO FIRST option")
	}
	if !strings.Contains(view, "2. SCHEDULE") {
		t.Error("expected to see SCHEDULE option")
	}
	if !strings.Contains(view, "3. DELEGATE") {
		t.Error("expected to see DELEGATE option")
	}
	if !strings.Contains(view, "4. ELIMINATE") {
		t.Error("expected to see ELIMINATE option")
	}
	if !strings.Contains(view, "Press ESC to cancel") {
		t.Error("expected to see cancel instruction")
	}
}

func TestStory016_SelectDestinationQuadrant(t *testing.T) {
	// Scenario: Select destination quadrant

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

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '2' for SCHEDULE
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view := model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to exit move mode after selection")
	}

	// Should have moved todo to priority B
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(B) Review quarterly goals") {
		t.Errorf("expected todo to be moved to priority B, got: %s", written)
	}
}

func TestStory016_CancelMoveModeWithESC(t *testing.T) {
	// Scenario: Cancel move mode with ESC

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

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should be in move mode
	view := model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to be in move mode")
	}

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view = model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to exit move mode after ESC")
	}

	// Todo should remain in DO FIRST (priority A)
	written := source.writer.(*strings.Builder).String()
	// No write should have happened - the builder should be empty
	if written != "" {
		t.Errorf("expected no changes to file after canceling move, got: %s", written)
	}
}

func TestStory016_MoveToEachQuadrant(t *testing.T) {
	// Scenario: Move to each quadrant

	tests := []struct {
		name            string
		initialPriority string
		destinationKey  string
		expectedPriority string
	}{
		{"Move to DO FIRST", "(B)", "1", "(A)"},
		{"Move to SCHEDULE", "(A)", "2", "(B)"},
		{"Move to DELEGATE", "(A)", "3", "(C)"},
		{"Move to ELIMINATE", "(A)", "4", "(D)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.initialPriority + " Test task"

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

			// Focus on the appropriate quadrant based on initial priority
			var focusKey rune
			switch tt.initialPriority {
			case "(A)":
				focusKey = '1'
			case "(B)":
				focusKey = '2'
			case "(C)":
				focusKey = '3'
			case "(D)":
				focusKey = '4'
			}

			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{focusKey}})
			model = updatedModel.(ui.Model)

			// Enter move mode
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
			model = updatedModel.(ui.Model)

			// Press destination key
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.destinationKey)})
			model = updatedModel.(ui.Model)

			// Check file was updated with expected priority
			written := source.writer.(*strings.Builder).String()
			expectedText := tt.expectedPriority + " Test task"
			if !strings.Contains(written, expectedText) {
				t.Errorf("expected file to contain '%s', got: %s", expectedText, written)
			}
		})
	}
}

func TestStory016_MovingToCurrentQuadrantIsNoOp(t *testing.T) {
	// Scenario: Moving to current quadrant is no-op

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

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '1' to move to DO FIRST (same quadrant)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view := model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to exit move mode after selection")
	}

	// Todo should remain in DO FIRST (no write should occur since priority unchanged)
	written := source.writer.(*strings.Builder).String()
	if written != "" {
		t.Errorf("expected no write to occur when moving to same quadrant, got: %s", written)
	}

	// Verify still viewing DO FIRST with the todo
	view = model.View()
	if !strings.Contains(view, "Review quarterly goals") {
		t.Error("expected todo to still be visible in DO FIRST")
	}
}

func TestStory016_MoveModeOnlyAvailableInFocusMode(t *testing.T) {
	// Scenario: Move mode only available in focus mode

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

	// In overview mode, press 'm'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should NOT show move mode overlay
	view := model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected 'm' to do nothing in overview mode")
	}

	// Should still be in overview mode
	if !strings.Contains(view, "DO FIRST") && !strings.Contains(view, "SCHEDULE") {
		t.Error("expected to remain in overview mode")
	}
}

func TestStory016_MoveModeOnlyAvailableWhenTodoSelected(t *testing.T) {
	// Scenario: Move mode only available when todo selected

	input := `(A) Task in DO FIRST`

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

	// Focus on SCHEDULE (empty quadrant)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)

	// Press 'm'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should NOT show move mode overlay
	view := model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected 'm' to do nothing in empty quadrant")
	}

	// Should show "(no tasks)"
	if !strings.Contains(view, "(no tasks)") {
		t.Error("expected to see '(no tasks)' in empty quadrant")
	}
}

func TestStory016_OtherKeysIgnoredInMoveMode(t *testing.T) {
	// Scenario: Other keys should be ignored while in move mode

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

	// Enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Should be in move mode
	view := model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to be in move mode")
	}

	// Press various other keys (should be ignored)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)
	view = model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to remain in move mode after pressing 'a'")
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	model = updatedModel.(ui.Model)
	view = model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to remain in move mode after pressing space")
	}

	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)
	view = model.View()
	if !strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to remain in move mode after pressing down arrow")
	}

	// Should still have no changes to the file
	written := source.writer.(*strings.Builder).String()
	if written != "" {
		t.Errorf("expected no changes to file while in move mode, got: %s", written)
	}

	// Press ESC to exit
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model = updatedModel.(ui.Model)

	// Should exit move mode
	view = model.View()
	if strings.Contains(view, "Move to quadrant:") {
		t.Error("expected to exit move mode after ESC")
	}
}
