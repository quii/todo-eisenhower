package ui

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderHelp renders help text with colored key bindings
func renderHelp(parts ...string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00D9FF")). // Bright cyan (works in both modes)
		Bold(true)

	textStyle := lipgloss.NewStyle().
		Foreground(TextSecondary)

	separatorStyle := lipgloss.NewStyle().
		Foreground(TextMuted)

	var result strings.Builder
	for i, part := range parts {
		// Split by spaces to identify keys vs text
		words := strings.Fields(part)
		expectKey := false
		for j, word := range words {
			lower := strings.ToLower(word)
			switch {
			case lower == "press":
				result.WriteString(textStyle.Render(word))
				expectKey = true
			case expectKey && isKeyBinding(word):
				result.WriteString(keyStyle.Render(word))
				expectKey = false
			default:
				result.WriteString(textStyle.Render(word))
				expectKey = false
			}

			if j < len(words)-1 {
				result.WriteString(" ")
			}
		}

		// Add separator between parts
		if i < len(parts)-1 {
			result.WriteString(separatorStyle.Render(" • "))
		}
	}

	return result.String()
}

// isKeyBinding checks if a word looks like a key binding
func isKeyBinding(word string) bool {
	// Check for single keys or key combinations
	keyBindings := []string{
		"a", "q", "esc", "enter", "space",
		"1", "2", "3", "4",
		"1/2/3/4", "1-4",
		"↑↓/w/s", "ctrl+c", "m",
	}

	lower := strings.ToLower(word)
	return slices.Contains(keyBindings, lower)
}
