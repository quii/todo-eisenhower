package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 008: Add Todo with Tag Reference

func TestStory008_PressAToEnterInputMode(t *testing.T) {
	// Scenario: Press 'a' to enter input mode in focused quadrant

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'a' to enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show input field
	if !strings.Contains(view, "Add todo:") {
		t.Error("expected view to show input prompt")
	}

	// Should show help text for input mode
	if !strings.Contains(view, "Enter to save") {
		t.Error("expected view to show 'Enter to save' help text")
	}
	if !strings.Contains(view, "ESC to cancel") {
		t.Error("expected view to show 'ESC to cancel' help text")
	}

	// Should show tag reference headers
	if !strings.Contains(view, "Projects:") {
		t.Error("expected view to show 'Projects:' label")
	}
	if !strings.Contains(view, "Contexts:") {
		t.Error("expected view to show 'Contexts:' label")
	}
}

func TestStory008_AddSimpleTodoWithoutTags(t *testing.T) {
	// Scenario: Add a simple todo without tags

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type the todo description
	for _, ch := range "Fix critical bug" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the new todo
	if !strings.Contains(view, "Fix critical bug") {
		t.Error("expected view to show new todo")
	}

	// Should exit input mode
	if strings.Contains(view, "Add todo:") {
		t.Error("expected view to exit input mode after saving")
	}

	// Check that the todo was written to the file with priority (A)
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Fix critical bug") {
		t.Errorf("expected todo to be written to file with priority (A), got: %s", written)
	}
}

func TestStory008_AddTodoWithProjectTags(t *testing.T) {
	// Scenario: Add todo with project tags

	input := "(B) Existing task +WebApp"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on SCHEDULE and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type the todo with project tags
	for _, ch := range "Plan sprint +WebApp +Mobile " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the new todo with tags
	if !strings.Contains(view, "Plan sprint") {
		t.Error("expected view to show new todo")
	}
	if !strings.Contains(view, "+WebApp") {
		t.Error("expected view to show +WebApp tag")
	}
	if !strings.Contains(view, "+Mobile") {
		t.Error("expected view to show +Mobile tag")
	}

	// Check that the todo was written with priority (B)
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(B) Plan sprint +WebApp +Mobile") {
		t.Errorf("expected todo to be written with priority (B) and tags, got: %s", written)
	}
}

func TestStory008_AddTodoWithContextTags(t *testing.T) {
	// Scenario: Add todo with context tags

	input := "(C) Existing task @phone"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DELEGATE and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type the todo with context tags
	for _, ch := range "Reply to emails @phone @office " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the new todo with context tags
	if !strings.Contains(view, "Reply to emails") {
		t.Error("expected view to show new todo")
	}
	if !strings.Contains(view, "@phone") {
		t.Error("expected view to show @phone tag")
	}
	if !strings.Contains(view, "@office") {
		t.Error("expected view to show @office tag")
	}

	// Check that the todo was written with priority (C)
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(C) Reply to emails @phone @office") {
		t.Errorf("expected todo to be written with priority (C) and tags, got: %s", written)
	}
}

func TestStory008_AddTodoWithMixedTags(t *testing.T) {
	// Scenario: Add todo with mixed tags

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type the todo with mixed tags
	for _, ch := range "Deploy to production +WebApp @computer @work " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show all tags
	if !strings.Contains(view, "Deploy to production") {
		t.Error("expected view to show new todo")
	}
	if !strings.Contains(view, "+WebApp") {
		t.Error("expected view to show +WebApp tag")
	}
	if !strings.Contains(view, "@computer") {
		t.Error("expected view to show @computer tag")
	}
	if !strings.Contains(view, "@work") {
		t.Error("expected view to show @work tag")
	}

	// Check written content
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Deploy to production +WebApp @computer @work") {
		t.Errorf("expected todo with all tags, got: %s", written)
	}
}

func TestStory008_CancelInputWithESC(t *testing.T) {
	// Scenario: Cancel input with ESC

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type some text
	for _, ch := range "This should be discarded" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press ESC to cancel
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should exit input mode
	if strings.Contains(view, "Add todo:") {
		t.Error("expected view to exit input mode after ESC")
	}

	// Should not show the typed text
	if strings.Contains(view, "This should be discarded") {
		t.Error("expected typed text to be discarded")
	}

	// Should not have written anything to file
	written := source.writer.(*strings.Builder).String()
	if written != "" {
		t.Errorf("expected no todo to be written after cancel, got: %s", written)
	}
}

func TestStory008_TagReferenceShowsExistingTags(t *testing.T) {
	// Scenario: Tag reference display shows existing tags

	input := "(A) Task 1 +WebApp @computer\n(B) Task 2 +Mobile @phone\n(C) Task 3 @office"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show project tags
	if !strings.Contains(view, "Projects:") {
		t.Error("expected view to show 'Projects:' label")
	}
	if !strings.Contains(view, "+WebApp") {
		t.Error("expected view to show +WebApp in tag reference")
	}
	if !strings.Contains(view, "+Mobile") {
		t.Error("expected view to show +Mobile in tag reference")
	}

	// Should show context tags
	if !strings.Contains(view, "Contexts:") {
		t.Error("expected view to show 'Contexts:' label")
	}
	if !strings.Contains(view, "@computer") {
		t.Error("expected view to show @computer in tag reference")
	}
	if !strings.Contains(view, "@phone") {
		t.Error("expected view to show @phone in tag reference")
	}
	if !strings.Contains(view, "@office") {
		t.Error("expected view to show @office in tag reference")
	}
}

func TestStory008_EmptyTagReferenceWhenNoTags(t *testing.T) {
	// Scenario: Empty tag reference when no tags exist

	input := "(A) Task without tags\n(B) Another task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show "(none)" for projects and contexts
	if !strings.Contains(view, "Projects: (none)") {
		t.Error("expected view to show 'Projects: (none)'")
	}
	if !strings.Contains(view, "Contexts: (none)") {
		t.Error("expected view to show 'Contexts: (none)'")
	}
}

func TestStory008_InputOnlyAvailableInFocusMode(t *testing.T) {
	// Scenario: Input only available in focus mode

	input := "(A) Task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press 'a'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should NOT show input mode
	if strings.Contains(view, "Add todo:") {
		t.Error("expected 'a' key to be ignored in overview mode")
	}

	// Should still show all quadrants (overview mode)
	if !strings.Contains(view, "DO FIRST") && !strings.Contains(view, "SCHEDULE") {
		t.Error("expected view to remain in overview mode")
	}
}

func TestStory008_AutoAssignPriorityFromQuadrant(t *testing.T) {
	// Scenario: Auto-assign priority from quadrant

	tests := []struct {
		name         string
		quadrantKey  rune
		expectedPrio string
	}{
		{"DO FIRST assigns priority A", '1', "(A)"},
		{"SCHEDULE assigns priority B", '2', "(B)"},
		{"DELEGATE assigns priority C", '3', "(C)"},
		{"ELIMINATE assigns priority D", '4', "(D)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := "(A) Existing task"
			source := &StubTodoSource{
				reader: strings.NewReader(input),
				writer: &strings.Builder{},
			}

			m, err := usecases.LoadMatrix(source)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			model := ui.NewModel(m, "test.txt").SetWriter(source)
			updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			model = updatedModel.(ui.Model)

			// Focus on the quadrant
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.quadrantKey}})
			model = updatedModel.(ui.Model)

			// Enter input mode
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			model = updatedModel.(ui.Model)

			// Type a todo
			for _, ch := range "Test todo" {
				updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
				model = updatedModel.(ui.Model)
			}

			// Save
			updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = updatedModel.(ui.Model)

			// Check written content
			written := source.writer.(*strings.Builder).String()
			if !strings.Contains(written, tt.expectedPrio+" Test todo") {
				t.Errorf("expected todo with priority %s, got: %s", tt.expectedPrio, written)
			}
		})
	}
}

func TestStory008_NewTagsAreAccepted(t *testing.T) {
	// Scenario: New tags are accepted

	input := "(A) Task +WebApp +Mobile"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type todo with new tag +Backend
	for _, ch := range "Build API +Backend " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should show the new todo with +Backend tag
	if !strings.Contains(view, "Build API") {
		t.Error("expected view to show new todo")
	}
	if !strings.Contains(view, "+Backend") {
		t.Error("expected view to show +Backend tag")
	}

	// Check written content
	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "(A) Build API +Backend") {
		t.Errorf("expected todo with +Backend tag, got: %s", written)
	}
}
