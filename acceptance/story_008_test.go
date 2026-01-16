package acceptance_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 008: Add Todo with Tag Reference

func TestStory008_PressAToEnterInputMode(t *testing.T) {
	// Scenario: Press 'a' to enter input mode in focused quadrant
	is := is.New(t)

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Add todo:")) // expected view to show input prompt

	// Should show help text for input mode
	is.True(strings.Contains(stripANSI(view), "Enter to save")) // expected view to show 'Enter to save' help text
	is.True(strings.Contains(stripANSI(view), "ESC to cancel")) // expected view to show 'ESC to cancel' help text

	// Should show tag reference headers
	is.True(strings.Contains(stripANSI(view), "Projects:")) // expected view to show 'Projects:' label
	is.True(strings.Contains(stripANSI(view), "Contexts:")) // expected view to show 'Contexts:' label
}

func TestStory008_AddSimpleTodoWithoutTags(t *testing.T) {
	// Scenario: Add a simple todo without tags
	is := is.New(t)

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Fix critical bug")) // expected view to show new todo

	// Should exit input mode
	is.True(!strings.Contains(stripANSI(view), "Add todo:")) // expected view to exit input mode after saving

	// Check that the todo was written to the file with priority (A) and creation date
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A)")) // expected todo to have priority (A)
	is.True(strings.Contains(written, "Fix critical bug")) // expected todo to contain description

	// Should have today's creation date
	today := time.Now().Format("2006-01-02")
	is.True(strings.Contains(written, today)) // expected todo to have today's creation date
}

func TestStory008_AddTodoWithProjectTags(t *testing.T) {
	// Scenario: Add todo with project tags
	is := is.New(t)

	input := "(B) Existing task +WebApp"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Plan sprint")) // expected view to show new todo
	is.True(strings.Contains(stripANSI(view), "WebApp")) // expected view to show +WebApp tag
	is.True(strings.Contains(stripANSI(view), "Mobile")) // expected view to show +Mobile tag

	// Check that the todo was written with priority (B) and creation date
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(B)")) // expected todo to have priority (B)
	is.True(strings.Contains(written, "Plan sprint +WebApp +Mobile")) // expected todo to contain description with tags

	today := time.Now().Format("2006-01-02")
	is.True(strings.Contains(written, today)) // expected todo to have today's creation date
}

func TestStory008_AddTodoWithContextTags(t *testing.T) {
	// Scenario: Add todo with context tags
	is := is.New(t)

	input := "(C) Existing task @phone"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Reply to emails")) // expected view to show new todo
	is.True(strings.Contains(stripANSI(view), "phone")) // expected view to show @phone tag
	is.True(strings.Contains(stripANSI(view), "office")) // expected view to show @office tag

	// Check that the todo was written with priority (C) and creation date
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(C)")) // expected todo to have priority (C)
	is.True(strings.Contains(written, "Reply to emails @phone @office")) // expected todo to contain description with tags

	today := time.Now().Format("2006-01-02")
	is.True(strings.Contains(written, today)) // expected todo to have today's creation date
}

func TestStory008_AddTodoWithMixedTags(t *testing.T) {
	// Scenario: Add todo with mixed tags
	is := is.New(t)

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Deploy to production")) // expected view to show new todo
	is.True(strings.Contains(stripANSI(view), "WebApp")) // expected view to show +WebApp tag
	is.True(strings.Contains(stripANSI(view), "computer")) // expected view to show @computer tag
	is.True(strings.Contains(stripANSI(view), "work")) // expected view to show @work tag

	// Check written content
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A)")) // expected todo to have priority (A)
	is.True(strings.Contains(written, "Deploy to production +WebApp @computer @work")) // expected todo to contain description with all tags

	// Should have today's creation date
	today := time.Now().Format("2006-01-02")
	is.True(strings.Contains(written, today)) // expected todo to have today's creation date
}

func TestStory008_CancelInputWithESC(t *testing.T) {
	// Scenario: Cancel input with ESC
	is := is.New(t)

	input := "(A) Existing task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(!strings.Contains(stripANSI(view), "Add todo:")) // expected view to exit input mode after ESC

	// Should not show the typed text
	is.True(!strings.Contains(stripANSI(view), "This should be discarded")) // expected typed text to be discarded

	// Should not have written anything to file
	written := source.writer.(*strings.Builder).String()
	is.Equal(written, "") // expected no todo to be written after cancel
}

func TestStory008_TagReferenceShowsExistingTags(t *testing.T) {
	// Scenario: Tag reference display shows existing tags
	is := is.New(t)

	input := "(A) Task 1 +WebApp @computer\n(B) Task 2 +Mobile @phone\n(C) Task 3 @office"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Projects:")) // expected view to show 'Projects:' label
	is.True(strings.Contains(stripANSI(view), "WebApp")) // expected view to show +WebApp in tag reference
	is.True(strings.Contains(stripANSI(view), "Mobile")) // expected view to show +Mobile in tag reference

	// Should show context tags
	is.True(strings.Contains(stripANSI(view), "Contexts:")) // expected view to show 'Contexts:' label
	is.True(strings.Contains(stripANSI(view), "computer")) // expected view to show @computer in tag reference
	is.True(strings.Contains(stripANSI(view), "phone")) // expected view to show @phone in tag reference
	is.True(strings.Contains(stripANSI(view), "office")) // expected view to show @office in tag reference
}

func TestStory008_EmptyTagReferenceWhenNoTags(t *testing.T) {
	// Scenario: Empty tag reference when no tags exist
	is := is.New(t)

	input := "(A) Task without tags\n(B) Another task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Projects: (none)")) // expected view to show 'Projects: (none)'
	is.True(strings.Contains(stripANSI(view), "Contexts: (none)")) // expected view to show 'Contexts: (none)'
}

func TestStory008_InputOnlyAvailableInFocusMode(t *testing.T) {
	// Scenario: Input only available in focus mode
	is := is.New(t)

	input := "(A) Task"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt").SetWriter(source)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode, press 'a'
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Should NOT show input mode
	is.True(!strings.Contains(stripANSI(view), "Add todo:")) // expected 'a' key to be ignored in overview mode

	// Should still show all quadrants (overview mode)
	is.True(strings.Contains(stripANSI(view), "Do First") && strings.Contains(stripANSI(view), "Schedule")) // expected view to remain in overview mode
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
			is := is.New(t)

			input := "(A) Existing task"
			source := &StubTodoSource{
				reader: strings.NewReader(input),
				writer: &strings.Builder{},
			}

			m, err := usecases.LoadMatrix(source)
			is.NoErr(err)

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
			is.True(strings.Contains(written, tt.expectedPrio)) // expected todo to have correct priority
			is.True(strings.Contains(written, "Test todo")) // expected todo to contain description

			// Should have today's creation date
			today := time.Now().Format("2006-01-02")
			is.True(strings.Contains(written, today)) // expected todo to have today's creation date
		})
	}
}

func TestStory008_NewTagsAreAccepted(t *testing.T) {
	// Scenario: New tags are accepted
	is := is.New(t)

	input := "(A) Task +WebApp +Mobile"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
		writer: &strings.Builder{},
	}

	m, err := usecases.LoadMatrix(source)
	is.NoErr(err)

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
	is.True(strings.Contains(stripANSI(view), "Build API")) // expected view to show new todo
	is.True(strings.Contains(stripANSI(view), "Backend")) // expected view to show +Backend tag

	// Check written content
	written := source.writer.(*strings.Builder).String()
	is.True(strings.Contains(written, "(A)")) // expected todo to have priority (A)
	is.True(strings.Contains(written, "Build API +Backend")) // expected todo to contain description with +Backend tag

	// Should have today's creation date
	today := time.Now().Format("2006-01-02")
	is.True(strings.Contains(written, today)) // expected todo to have today's creation date
}
