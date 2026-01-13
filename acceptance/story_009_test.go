package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 009: Tag Autocomplete

func TestStory009_TriggerAutocompleteWithPlus(t *testing.T) {
	// Scenario: Trigger autocomplete with +

	input := "(A) Task +WebApp\n(A) Task +Mobile\n(A) Task +Backend"
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

	// Type "Deploy feature +"
	for _, ch := range "Deploy feature +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show suggestions
	if !strings.Contains(view, "WebApp") {
		t.Error("expected autocomplete to show WebApp")
	}
	if !strings.Contains(view, "Mobile") {
		t.Error("expected autocomplete to show Mobile")
	}
	if !strings.Contains(view, "Backend") {
		t.Error("expected autocomplete to show Backend")
	}
}

func TestStory009_FilterSuggestionsAsIType(t *testing.T) {
	// Scenario: Filter suggestions as I type

	input := "(A) Existing task +WebApp\n(B) Other task +Mobile\n(C) Another task +Backend"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy feature +Web"
	for _, ch := range "Deploy feature +Web" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should only show WebApp
	if !strings.Contains(view, "WebApp") {
		t.Error("expected autocomplete to show WebApp")
	}
	if strings.Contains(view, "Mobile") {
		t.Error("expected autocomplete to NOT show Mobile")
	}
	if strings.Contains(view, "Backend") {
		t.Error("expected autocomplete to NOT show Backend")
	}
}

func TestStory009_NavigateSuggestionsWithArrows(t *testing.T) {
	// Scenario: Navigate suggestions with arrow keys

	input := "(A) Task +WebApp\n(A) Task +Mobile\n(A) Task +Backend"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +"
	for _, ch := range "Deploy +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Down arrow - should highlight second suggestion
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Down arrow again - should highlight third suggestion
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Down arrow again - should wrap to first suggestion
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Up arrow - should wrap to last suggestion
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})
	model = updatedModel.(ui.Model)

	// Just verify we can navigate without errors
	view := model.View()
	if !strings.Contains(view, "WebApp") {
		t.Error("expected suggestions to still be visible after navigation")
	}
}

func TestStory009_CompleteTagWithTab(t *testing.T) {
	// Scenario: Complete tag with Tab

	input := "(A) Task +API\n(A) Task +WebApp"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +"
	for _, ch := range "Deploy +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Tab to complete
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(ui.Model)

	// Suggestions should be dismissed
	// Input should show completed tag (we can't easily check the input value in tests,
	// but we can verify by continuing to type and then saving)

	// Type more text to verify we can continue
	for _, ch := range "more text" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	written := source.writer.(*strings.Builder).String()
	if !strings.Contains(written, "+API") {
		t.Errorf("expected todo to contain +API tag, got: %s", written)
	}
	if !strings.Contains(written, "more text") {
		t.Errorf("expected todo to contain 'more text', got: %s", written)
	}
}

func TestStory009_CompleteTagWithEnter(t *testing.T) {
	// Scenario: Complete tag with Enter (when suggestions visible)

	input := "(A) Task +API\n(A) Task +WebApp"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +"
	for _, ch := range "Deploy +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Navigate to second tag (WebApp)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updatedModel.(ui.Model)

	// Press Enter to complete (should NOT save todo)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Verify we're still in input mode
	view := model.View()
	if !strings.Contains(view, "Enter to save") {
		t.Error("expected to still be in input mode after completing with Enter")
	}

	// Verify no todo was saved yet
	written := source.writer.(*strings.Builder).String()
	if written != "" {
		t.Error("expected no todo to be saved when using Enter to complete suggestion")
	}
}

func TestStory009_DismissSuggestionsWithESC(t *testing.T) {
	// Scenario: Dismiss suggestions with ESC

	input := "(A) Task +WebApp"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +"
	for _, ch := range "Deploy +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Verify suggestions are showing
	view1 := model.View()
	if !strings.Contains(view1, "WebApp") {
		t.Error("expected suggestions to be visible before ESC")
	}

	// Press ESC to dismiss suggestions
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)

	// Should still be in input mode
	view2 := model.View()
	if !strings.Contains(view2, "Enter to save") {
		t.Error("expected to still be in input mode after ESC")
	}

	// Verify no todo was saved
	written := source.writer.(*strings.Builder).String()
	if written != "" {
		t.Error("expected no todo to be saved when dismissing suggestions")
	}
}

func TestStory009_AutocompleteContextTags(t *testing.T) {
	// Scenario: Autocomplete context tags with @

	input := "(B) Existing task @computer\n(C) Other task @phone\n(D) Another task @office"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Reply to emails @"
	for _, ch := range "Reply to emails @" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show all context tags
	if !strings.Contains(view, "computer") {
		t.Error("expected autocomplete to show computer")
	}
	if !strings.Contains(view, "phone") {
		t.Error("expected autocomplete to show phone")
	}
	if !strings.Contains(view, "office") {
		t.Error("expected autocomplete to show office")
	}

	// Type "p" to filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(ui.Model)

	view2 := model.View()

	// Should only show phone
	if !strings.Contains(view2, "phone") {
		t.Error("expected autocomplete to show phone")
	}
	if strings.Contains(view2, "computer") {
		t.Error("expected autocomplete to NOT show computer")
	}
	if strings.Contains(view2, "office") {
		t.Error("expected autocomplete to NOT show office")
	}
}

func TestStory009_MultipleTagsInOneInput(t *testing.T) {
	// Scenario: Multiple tags in one todo

	input := "(B) Existing task +WebApp @computer"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +"
	for _, ch := range "Deploy +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Complete with Tab
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model = updatedModel.(ui.Model)

	// Type " @"
	for _, ch := range " @" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show context tags
	if !strings.Contains(view, "computer") {
		t.Error("expected autocomplete to show context tags after @")
	}
	// Note: We can't check that "WebApp" doesn't appear because it will appear
	// in the input field text itself (from the completed tag), which is expected.
	// What matters is that the autocomplete suggestions show context tags.
}

func TestStory009_NoMatchesMessage(t *testing.T) {
	// Scenario: No suggestions available

	input := "(A) Task +WebApp"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +xyz" (no matches)
	for _, ch := range "Deploy +xyz" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show no matches message
	if !strings.Contains(view, "no matches") || !strings.Contains(view, "Space") {
		t.Error("expected autocomplete to show 'no matches' message")
	}
}

func TestStory009_CaseInsensitiveMatching(t *testing.T) {
	// Scenario: Case-insensitive matching

	input := "(A) Task +WebApp"
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

	// Enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Deploy +web" (lowercase)
	for _, ch := range "Deploy +web" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show WebApp (original case)
	if !strings.Contains(view, "WebApp") {
		t.Error("expected case-insensitive match to find WebApp")
	}
}
