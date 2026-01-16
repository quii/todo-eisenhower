package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderDeleteOverlay renders a confirmation overlay for delete mode
func RenderDeleteOverlay(terminalWidth, terminalHeight int) string {
	// Build overlay content
	content := lipgloss.NewStyle().
		Bold(true).
		Render("Delete this todo?") + "\n\n"

	content += "  y - Yes, delete it\n"
	content += "  n - No, keep it\n\n"

	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Render("Press ESC to cancel")

	// Create bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7B68EE")).
		Padding(1, 2).
		Width(30)

	box := boxStyle.Render(content)

	// Center the box in the terminal
	return lipgloss.Place(
		terminalWidth,
		terminalHeight,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}
