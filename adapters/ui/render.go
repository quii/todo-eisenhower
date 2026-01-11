package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

var (
	quadrantStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1).
			Width(40).
			Height(10)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("6"))

	todoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))
)

// RenderMatrix renders the Eisenhower matrix as a string
func RenderMatrix(m matrix.Matrix) string {
	doFirst := renderQuadrant("DO FIRST", m.DoFirst())
	schedule := renderQuadrant("SCHEDULE", m.Schedule())
	delegate := renderQuadrant("DELEGATE", m.Delegate())
	eliminate := renderQuadrant("ELIMINATE", m.Eliminate())

	// Arrange quadrants in 2x2 grid
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, doFirst, schedule)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, delegate, eliminate)

	return lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)
}

func renderQuadrant(title string, todos []todo.Todo) string {
	var content strings.Builder

	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n\n")

	if len(todos) == 0 {
		content.WriteString(todoStyle.Render("(empty)"))
	} else {
		for _, t := range todos {
			prefix := "• "
			if t.IsCompleted() {
				prefix = "✓ "
			}
			content.WriteString(todoStyle.Render(prefix + t.Description()))
			content.WriteString("\n")
		}
	}

	return quadrantStyle.Render(content.String())
}
