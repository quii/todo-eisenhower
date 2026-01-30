package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// RenderDeleteOverlay renders a confirmation overlay for delete mode
func RenderDeleteOverlay(terminalWidth, terminalHeight int) string {
	content := lipgloss.NewStyle().
		Bold(true).
		Render("Delete this todo?") + "\n\n"

	content += "  y - Yes, delete it\n"
	content += "  n - No, keep it\n\n"

	content += lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Italic(true).
		Render("Press ESC to cancel")

	return renderCenteredOverlay(content, 30, lipgloss.Color("#7B68EE"), terminalWidth, terminalHeight)
}
