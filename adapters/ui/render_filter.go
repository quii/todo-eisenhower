package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// RenderFilterInput renders the filter input UI
func RenderFilterInput(
	filePath string,
	input textinput.Model,
	showSuggestions bool,
	suggestions []string,
	selectedSuggestion int,
	terminalWidth, terminalHeight int,
) string {
	var output strings.Builder

	// Render file path header
	if filePath != "" {
		header := headerStyle.
			Width(terminalWidth).
			Align(lipgloss.Center).
			Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(TextPrimary).
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render("Filter Todos")
	output.WriteString(title)
	output.WriteString("\n\n")

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(TextSecondary).
		Italic(true).
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render("Type a project (+) or context (@) to filter")
	output.WriteString(instructions)
	output.WriteString("\n\n")

	// Input prompt
	promptText := "Filter by: "
	inputPrompt := lipgloss.NewStyle().
		Foreground(TextPrimary).
		Render(promptText)

	// Center the input section
	inputLine := inputPrompt + input.View()
	centeredInput := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(inputLine)
	output.WriteString(centeredInput)
	output.WriteString("\n\n")

	// Render autocomplete suggestions if visible
	if showSuggestions && len(suggestions) > 0 {
		autocompleteBox := renderFilterAutocomplete(suggestions, selectedSuggestion, terminalWidth)
		output.WriteString(autocompleteBox)
		output.WriteString("\n\n")
	}

	// Help text
	helpText := renderHelp("Enter to apply", "ESC to cancel")
	centeredHelp := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(helpText)
	output.WriteString(centeredHelp)

	// Center everything vertically
	content := output.String()
	return lipgloss.Place(terminalWidth, terminalHeight, lipgloss.Center, lipgloss.Center, content)
}

// renderFilterAutocomplete renders autocomplete suggestions for filter mode
func renderFilterAutocomplete(suggestions []string, selectedIndex, width int) string {
	var lines []string
	maxSuggestions := 8 // Show more suggestions in filter mode

	for i, suggestion := range suggestions {
		if i >= maxSuggestions {
			break
		}

		// Determine color based on prefix
		var color lipgloss.TerminalColor
		var prefix byte
		if suggestion != "" {
			prefix = suggestion[0]
		}

		switch prefix {
		case '+', '@':
			color = HashColor(suggestion[1:])
		default:
			color = TextPrimary
		}

		if i == selectedIndex {
			// Highlighted suggestion
			suggestionStyle := lipgloss.NewStyle().
				Foreground(color).
				Background(SelectionBg).
				Bold(true).
				Padding(0, 1)
			lines = append(lines, suggestionStyle.Render(suggestion))
		} else {
			// Regular suggestion
			suggestionStyle := lipgloss.NewStyle().
				Foreground(color).
				Padding(0, 1)
			lines = append(lines, suggestionStyle.Render(suggestion))
		}
	}

	// Box style
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderAccent).
		Padding(0, 1).
		Align(lipgloss.Center)

	content := strings.Join(lines, "\n")
	box := boxStyle.Render(content)

	// Center the box
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(width).
		Render(box)
}
