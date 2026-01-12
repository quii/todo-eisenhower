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
			} else {
				todoLine = activeTodoStyle.Render("• ") + description
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
	output.WriteString("\n\n")

	// Render tag reference
	projectLine := renderTagReference("Projects", projects, terminalWidth)
	contextLine := renderTagReference("Contexts", contexts, terminalWidth)
	output.WriteString(projectLine)
	output.WriteString("\n")
	output.WriteString(contextLine)
	output.WriteString("\n\n")

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
