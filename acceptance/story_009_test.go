package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 009: Tag Autocomplete

func TestStory009_TriggerAutocompleteWithPlus(t *testing.T) {
	is := is.New(t)
	// Scenario: Trigger autocomplete with +

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"Mobile"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"Backend"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "WebApp"))   // expected autocomplete to show WebApp
	is.True(strings.Contains(stripANSI(view), "Mobile"))   // expected autocomplete to show Mobile
	is.True(strings.Contains(stripANSI(view), "Backend"))  // expected autocomplete to show Backend
}

func TestStory009_FilterSuggestionsAsIType(t *testing.T) {
	is := is.New(t)
	// Scenario: Filter suggestions as I type

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Existing task", todo.PriorityA, []string{"WebApp"}, []string{}),
		todo.NewWithTags("Other task", todo.PriorityB, []string{"Mobile"}, []string{}),
		todo.NewWithTags("Another task", todo.PriorityC, []string{"Backend"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "WebApp"))    // expected autocomplete to show WebApp
	is.True(!strings.Contains(stripANSI(view), "Mobile"))   // expected autocomplete to NOT show Mobile
	is.True(!strings.Contains(stripANSI(view), "Backend"))  // expected autocomplete to NOT show Backend
}

func TestStory009_NavigateSuggestionsWithArrows(t *testing.T) {
	is := is.New(t)
	// Scenario: Navigate suggestions with arrow keys

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"Mobile"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"Backend"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "WebApp"))  // expected suggestions to still be visible after navigation
}

func TestStory009_CompleteTagWithTab(t *testing.T) {
	is := is.New(t)
	// Scenario: Complete tag with Tab

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"API"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	_ = updatedModel.(ui.Model)

	// Verify the todo was persisted with completed tag
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 3) // expected 3 todos (2 existing + 1 new)

	newTodo := savedTodos[2]
	is.Equal(len(newTodo.Projects()), 1)
	is.Equal(newTodo.Projects()[0], "API")
	is.True(strings.Contains(newTodo.Description(), "more text"))
}

func TestStory009_CompleteTagWithEnter(t *testing.T) {
	is := is.New(t)
	// Scenario: Complete tag with Enter (when suggestions visible)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"API"}, []string{}),
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "Enter to save"))  // expected to still be in input mode after completing with Enter

	// Verify no todo was saved yet
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 2) // expected only original 2 todos, no new one
}

func TestStory009_DismissSuggestionsWithESC(t *testing.T) {
	is := is.New(t)
	// Scenario: Dismiss suggestions with ESC

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(view1, "WebApp"))  // expected suggestions to be visible before ESC

	// Press ESC to dismiss suggestions
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)

	// Should still be in input mode
	view2 := model.View()
	is.True(strings.Contains(view2, "Enter to save"))  // expected to still be in input mode after ESC

	// Verify no todo was saved
	savedTodos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(savedTodos), 1) // expected only original todo
}

func TestStory009_AutocompleteContextTags(t *testing.T) {
	is := is.New(t)
	// Scenario: Autocomplete context tags with @

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Existing task", todo.PriorityB, []string{}, []string{"computer"}),
		todo.NewWithTags("Other task", todo.PriorityC, []string{}, []string{"phone"}),
		todo.NewWithTags("Another task", todo.PriorityD, []string{}, []string{"office"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "computer"))  // expected autocomplete to show computer
	is.True(strings.Contains(stripANSI(view), "phone"))     // expected autocomplete to show phone
	is.True(strings.Contains(stripANSI(view), "office"))    // expected autocomplete to show office

	// Type "p" to filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	model = updatedModel.(ui.Model)

	view2 := model.View()

	// Should only show phone
	is.True(strings.Contains(view2, "phone"))      // expected autocomplete to show phone
	is.True(!strings.Contains(view2, "computer"))  // expected autocomplete to NOT show computer
	is.True(!strings.Contains(view2, "office"))    // expected autocomplete to NOT show office
}

func TestStory009_MultipleTagsInOneInput(t *testing.T) {
	is := is.New(t)
	// Scenario: Multiple tags in one todo

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Existing task", todo.PriorityB, []string{"WebApp"}, []string{"computer"}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "computer"))  // expected autocomplete to show context tags after @
	// Note: We can't check that "WebApp" doesn't appear because it will appear
	// in the input field text itself (from the completed tag), which is expected.
	// What matters is that the autocomplete suggestions show context tags.
}

func TestStory009_NoMatchesMessage(t *testing.T) {
	is := is.New(t)
	// Scenario: No suggestions available

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "no matches") && strings.Contains(stripANSI(view), "Space"))  // expected autocomplete to show 'no matches' message
}

func TestStory009_CaseInsensitiveMatching(t *testing.T) {
	is := is.New(t)
	// Scenario: Case-insensitive matching

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityA, []string{"WebApp"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
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
	is.True(strings.Contains(stripANSI(view), "WebApp"))  // expected case-insensitive match to find WebApp
}

func TestStory009_NewTagsAvailableInAutocomplete(t *testing.T) {
	is := is.New(t)
	// Scenario: Newly created tags should appear in autocomplete

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.NewWithTags("Task", todo.PriorityB, []string{"OldTag"}, []string{}),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on DO FIRST and enter input mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Add a new todo with a new tag +NewTag
	for _, ch := range "First task +NewTag " {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	// Save the todo (priority A - goes to DO FIRST)
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updatedModel.(ui.Model)

	// Enter input mode again
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type "Second task +"
	for _, ch := range "Second task +" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
		model = updatedModel.(ui.Model)
	}

	view := model.View()

	// Should show both OldTag and NewTag in autocomplete
	// Note: OldTag is from SCHEDULE quadrant (priority B), so won't appear in DO FIRST's todo list
	is.True(strings.Contains(stripANSI(view), "OldTag"))  // expected autocomplete to show +OldTag tag
	is.True(strings.Contains(stripANSI(view), "NewTag"))  // expected autocomplete to show newly created +NewTag tag

	// Type "N" to filter
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
	model = updatedModel.(ui.Model)

	view2 := model.View()

	// Should only show NewTag
	is.True(strings.Contains(view2, "NewTag"))  // expected autocomplete to filter to +NewTag
	// OldTag should not appear in autocomplete suggestions (filtered out)
	// However, it might still appear in the visible todo list, so we can't check this reliably
}
