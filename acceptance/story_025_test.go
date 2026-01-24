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

// Story 025: URL Support
// As a user, I want to attach URLs to tasks and open them from the TUI
// so that I can quickly access related documents, PRs, and resources.

func TestStory025_HelpTextShowsOpenURLCommand(t *testing.T) {
	is := is.New(t)

	// Given a matrix with a task
	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{
		todo.New("Review design doc", todo.PriorityA),
	})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// When I focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)
	view := model.View()

	// Then help text should show "o to open URL"
	is.True(strings.Contains(stripANSI(view), "o to open URL")) // expected help text to show open URL command
}

func TestStory025_ExtractsSingleURL(t *testing.T) {
	is := is.New(t)

	// This test verifies that extractURLs correctly finds a single URL
	url := "https://docs.google.com/document/d/abc123"
	description := "Review design doc " + url + " +design"

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find exactly one URL
	is.Equal(len(urls), 1)           // expected to find one URL
	is.Equal(urls[0], url) // expected URL to match
}

func TestStory025_ExtractsMultipleURLs(t *testing.T) {
	is := is.New(t)

	// Given a task description with multiple URLs
	url1 := "https://github.com/user/repo/pull/42"
	url2 := "https://staging.example.com"
	description := "Review PR " + url1 + " and staging " + url2 + " +review"

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find both URLs
	is.Equal(len(urls), 2)             // expected to find two URLs
	is.Equal(urls[0], url1) // expected first URL to match
	is.Equal(urls[1], url2) // expected second URL to match
}

func TestStory025_ExtractsNoURLs(t *testing.T) {
	is := is.New(t)

	// Given a task description with no URLs
	description := "Regular task without links +project @work"

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find no URLs
	is.Equal(len(urls), 0) // expected to find no URLs
}

func TestStory025_SupportsHTTPAndHTTPS(t *testing.T) {
	is := is.New(t)

	// Given a task with both HTTP and HTTPS URLs
	httpURL := "http://example.com"
	httpsURL := "https://example.com"
	description := "Check " + httpURL + " and " + httpsURL

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find both URLs
	is.Equal(len(urls), 2)                // expected to find two URLs
	is.Equal(urls[0], httpURL)  // expected HTTP URL
	is.Equal(urls[1], httpsURL) // expected HTTPS URL
}

func TestStory025_URLsWithQueryParams(t *testing.T) {
	is := is.New(t)

	// Given a URL with query parameters
	url := "https://google.com/search?q=test&page=1"
	description := "Search " + url

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find URL with query params intact
	is.Equal(len(urls), 1)          // expected to find one URL
	is.Equal(urls[0], url) // expected URL with query params
}

func TestStory025_URLsWithAnchor(t *testing.T) {
	is := is.New(t)

	// Given a URL with an anchor
	url := "https://example.com/page#section-1"
	description := "See " + url + " for details"

	// Extract URLs
	urls := ui.ExtractURLsForTest(description)

	// Should find URL with anchor intact
	is.Equal(len(urls), 1)          // expected to find one URL
	is.Equal(urls[0], url) // expected URL with anchor
}
