package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// RenderFocusedQuadrantWithInput renders a quadrant in focus mode with input UI
func RenderFocusedQuadrantWithInput(
	todos []todo.Todo,
	title string,
	color lipgloss.Color,
	filePath string,
	input textinput.Model,
	projects, contexts []string,
	showSuggestions bool,
	suggestions []string,
	selectedSuggestion int,
	terminalWidth, terminalHeight int,
) string {
	var output strings.Builder

	// Render file path header with full width and center alignment
	if filePath != "" {
		header := headerStyle.
			Copy().
			Width(terminalWidth).
			Align(lipgloss.Center).
			Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Calculate display limit for focus mode with input
	// Reserve: header (3), title (2), input section (6), help text (2), margins (2) = 15 lines
	displayLimit := terminalHeight - 15
	if displayLimit < 3 {
		displayLimit = 3
	}

	// Render prominent quadrant title
	focusTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Underline(true).
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(title)
	output.WriteString(focusTitle)
	output.WriteString("\n\n")

	// Render todos
	var lines []string
	if len(todos) == 0 {
		lines = append(lines, emptyStyle.Render("(no tasks)"))
	} else {
		for i, t := range todos {
			if i >= displayLimit {
				remaining := len(todos) - displayLimit
				lines = append(lines, emptyStyle.Render(fmt.Sprintf("... and %d more", remaining)))
				break
			}

			// Colorize tags in description
			description := colorizeDescription(t.Description())

			var todoLine string
			if t.IsCompleted() {
				todoLine = completedTodoStyle.Render("✓ ") + description
				// Add date information
				createdStr := formatDate(t.CreationDate())
				completedStr := formatDate(t.CompletionDate())
				if createdStr != "" && completedStr != "" {
					dateInfo := emptyStyle.Render(fmt.Sprintf(" (added %s, completed %s)", createdStr, completedStr))
					todoLine += dateInfo
				} else if completedStr != "" {
					dateInfo := emptyStyle.Render(fmt.Sprintf(" (completed %s)", completedStr))
					todoLine += dateInfo
				}
			} else {
				todoLine = activeTodoStyle.Render("• ") + description
				// Add creation date for active todos
				if createdStr := formatDate(t.CreationDate()); createdStr != "" {
					dateInfo := emptyStyle.Render(fmt.Sprintf(" (added %s)", createdStr))
					todoLine += dateInfo
				}
			}
			lines = append(lines, todoLine)
		}
	}

	todosContent := strings.Join(lines, "\n")
	output.WriteString(todosContent)
	output.WriteString("\n\n")

	// Render input section
	output.WriteString(dividerStyle.Render(strings.Repeat("─", terminalWidth)))
	output.WriteString("\n\n")

	// Input prompt and field
	inputPrompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render("Add todo: ")
	output.WriteString(inputPrompt)
	output.WriteString(input.View())
	output.WriteString("\n")

	// Render autocomplete suggestions if visible, otherwise show tag reference
	if showSuggestions {
		trigger, partialTag, _ := detectTrigger(input.Value())
		autocompleteBox := renderAutocomplete(suggestions, selectedSuggestion, trigger, partialTag, terminalWidth)
		output.WriteString("\n")
		output.WriteString(autocompleteBox)
		output.WriteString("\n\n")
	} else {
		output.WriteString("\n")
		// Render tag reference when not autocompleting
		projectLine := renderTagReference("Projects", projects, terminalWidth)
		contextLine := renderTagReference("Contexts", contexts, terminalWidth)
		output.WriteString(projectLine)
		output.WriteString("\n")
		output.WriteString(contextLine)
		output.WriteString("\n\n")
	}

	// Render help text at bottom
	helpText := renderHelp("Enter to save", "ESC to cancel")
	centeredHelp := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(helpText)
	output.WriteString(centeredHelp)

	return output.String()
}

// renderTagReference renders a line showing available tags
func renderTagReference(label string, tags []string, width int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	var content string
	if len(tags) == 0 {
		content = labelStyle.Render(label + ": (none)")
	} else {
		// Build tag list with colors
		var tagParts []string
		prefix := "+"
		if label == "Contexts" {
			prefix = "@"
		}
		for _, tag := range tags {
			tagWithPrefix := prefix + tag
			color := HashColor(tag)
			styled := lipgloss.NewStyle().
				Foreground(color).
				Render(tagWithPrefix)
			tagParts = append(tagParts, styled)
		}
		content = labelStyle.Render(label+": ") + strings.Join(tagParts, " ")
	}

	return content
}

// renderAutocomplete renders the autocomplete suggestion box
func renderAutocomplete(suggestions []string, selectedIndex int, trigger string, partialTag string, width int) string {
	if len(suggestions) == 0 {
		// Show "no matches" message
		noMatchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)
		return noMatchStyle.Render("  (no matches - press Space to create new tag)")
	}

	var lines []string
	maxSuggestions := 5
	for i, suggestion := range suggestions {
		if i >= maxSuggestions {
			break
		}

		// Use the trigger character passed in (+ or @)
		tagWithPrefix := trigger + suggestion

		// Color the tag
		color := HashColor(suggestion)

		if i == selectedIndex {
			// Highlighted suggestion
			suggestionStyle := lipgloss.NewStyle().
				Foreground(color).
				Background(lipgloss.Color("#444444")).
				Bold(true).
				Padding(0, 1)
			lines = append(lines, suggestionStyle.Render(tagWithPrefix))
		} else {
			// Regular suggestion
			suggestionStyle := lipgloss.NewStyle().
				Foreground(color).
				Padding(0, 1)
			lines = append(lines, suggestionStyle.Render(tagWithPrefix))
		}
	}

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666666")).
		Padding(0, 1)

	content := strings.Join(lines, "\n")
	return boxStyle.Render(content)
}
