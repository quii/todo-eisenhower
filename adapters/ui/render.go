package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

const (
	quadrantWidth  = 40
	quadrantHeight = 10
)

var (
	// Color palette - Eisenhower matrix themed
	urgentImportantColor = lipgloss.Color("#FF6B6B") // Red - Do First
	importantColor       = lipgloss.Color("#4ECDC4") // Teal - Schedule
	urgentColor          = lipgloss.Color("#FFE66D") // Yellow - Delegate
	neitherColor         = lipgloss.Color("#95E1D3") // Light green - Eliminate

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5F5F87")).
			Padding(0, 2).
			MarginBottom(1)

	matrixBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Underline(true)

	activeTodoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	completedTodoStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#808080")).
				Strikethrough(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))
)

// RenderMatrix renders the Eisenhower matrix as a string with optional file path header
func RenderMatrix(m matrix.Matrix, filePath string) string {
	var output strings.Builder

	// Render header if file path provided
	if filePath != "" {
		header := headerStyle.Render("📄 File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Add axis label
	urgentLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		Render("← URGENT →")

	output.WriteString(lipgloss.Place(quadrantWidth*2+3, 1, lipgloss.Center, lipgloss.Top, urgentLabel))
	output.WriteString("\n")

	// Render quadrant contents
	doFirst := renderQuadrantContent("🔥 DO FIRST", urgentImportantColor, m.DoFirst())
	schedule := renderQuadrantContent("📅 SCHEDULE", importantColor, m.Schedule())
	delegate := renderQuadrantContent("👥 DELEGATE", urgentColor, m.Delegate())
	eliminate := renderQuadrantContent("🗑️  ELIMINATE", neitherColor, m.Eliminate())

	// Create vertical divider that spans quadrant height
	verticalDivider := createVerticalDivider(quadrantHeight)

	// Build top row with divider
	topRow := lipgloss.JoinHorizontal(lipgloss.Top,
		doFirst,
		verticalDivider,
		schedule,
	)

	// Build horizontal divider line
	horizontalLine := dividerStyle.Render(strings.Repeat("─", quadrantWidth*2+1))

	// Build bottom row with divider
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top,
		delegate,
		verticalDivider,
		eliminate,
	)

	// Combine all parts
	matrixContent := lipgloss.JoinVertical(lipgloss.Left,
		topRow,
		horizontalLine,
		bottomRow,
	)

	// Wrap entire matrix in border
	matrix := matrixBorder.Render(matrixContent)
	output.WriteString(matrix)

	return output.String()
}

// createVerticalDivider creates a vertical divider that spans the given height
func createVerticalDivider(height int) string {
	var divider strings.Builder
	for i := 0; i < height; i++ {
		divider.WriteString("│\n")
	}
	return dividerStyle.Render(strings.TrimSuffix(divider.String(), "\n"))
}

// renderQuadrantContent renders just the content of a quadrant (no border)
func renderQuadrantContent(title string, color lipgloss.Color, todos []todo.Todo) string {
	var lines []string

	// Title
	quadrantTitle := titleStyle.
		Copy().
		Foreground(color).
		Render(title)
	lines = append(lines, quadrantTitle)
	lines = append(lines, "") // spacing

	if len(todos) == 0 {
		lines = append(lines, emptyStyle.Render("(no tasks)"))
	} else {
		// Limit display to prevent overflow
		displayLimit := 7
		for i, t := range todos {
			if i >= displayLimit {
				remaining := len(todos) - displayLimit
				lines = append(lines, emptyStyle.Render(fmt.Sprintf("... and %d more", remaining)))
				break
			}

			var todoLine string
			if t.IsCompleted() {
				todoLine = completedTodoStyle.Render("✓ " + t.Description())
			} else {
				todoLine = activeTodoStyle.Render("• " + t.Description())
			}
			lines = append(lines, todoLine)
		}
	}

	// Join lines and place in exact dimensions
	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Width(quadrantWidth).
		Height(quadrantHeight).
		Padding(1, 2).
		AlignVertical(lipgloss.Top).
		Render(content)
}
