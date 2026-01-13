package ui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

const (
	defaultQuadrantWidth  = 40
	defaultQuadrantHeight = 10
	minQuadrantWidth      = 30
	minQuadrantHeight     = 8
)

var (
	projectTagPattern = regexp.MustCompile(`\+(\w+)`)
	contextTagPattern = regexp.MustCompile(`@(\w+)`)
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
// terminalWidth and terminalHeight are optional (0 = use defaults)
func RenderMatrix(m matrix.Matrix, filePath string, terminalWidth, terminalHeight int) string {
	// Calculate quadrant dimensions based on terminal size
	quadrantWidth, quadrantHeight := calculateQuadrantDimensions(terminalWidth, terminalHeight)
	displayLimit := calculateDisplayLimit(quadrantHeight)

	var output strings.Builder

	// Render header if file path provided
	if filePath != "" {
		header := headerStyle.Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Render quadrant contents
	doFirst := renderQuadrantContent("DO FIRST", urgentImportantColor, m.DoFirst(), quadrantWidth, quadrantHeight, displayLimit)
	schedule := renderQuadrantContent("SCHEDULE", importantColor, m.Schedule(), quadrantWidth, quadrantHeight, displayLimit)
	delegate := renderQuadrantContent("DELEGATE", urgentColor, m.Delegate(), quadrantWidth, quadrantHeight, displayLimit)
	eliminate := renderQuadrantContent("ELIMINATE", neitherColor, m.Eliminate(), quadrantWidth, quadrantHeight, displayLimit)

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
	output.WriteString("\n\n")

	// Add tag inventory
	inventory := renderTagInventory(m, terminalWidth)
	output.WriteString(inventory)
	output.WriteString("\n\n")

	// Add help text
	helpText := renderHelp("Press 1/2/3/4 to focus on a quadrant")
	output.WriteString(helpText)

	return output.String()
}

// RenderFocusedQuadrant renders a single quadrant in fullscreen focus mode
func RenderFocusedQuadrant(todos []todo.Todo, title string, color lipgloss.Color, filePath string, terminalWidth, terminalHeight int) string {
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

	// Calculate display limit for focus mode
	// Reserve: header (3), title (2), help text (2), margins (2) = 9 lines
	displayLimit := terminalHeight - 9
	if displayLimit < 5 {
		displayLimit = 5
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

	// Render help text at bottom
	helpText := renderHelp("Press a to add a task", "Press 1/2/3/4 to focus on a quadrant", "Press ESC to return")
	centeredHelp := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(helpText)
	output.WriteString(centeredHelp)

	return output.String()
}

// calculateQuadrantDimensions calculates optimal quadrant size based on terminal dimensions
func calculateQuadrantDimensions(terminalWidth, terminalHeight int) (width, height int) {
	// Use defaults if no terminal size provided
	if terminalWidth == 0 || terminalHeight == 0 {
		return defaultQuadrantWidth, defaultQuadrantHeight
	}

	// Reserve space for:
	// - File header: 3 lines
	// - Matrix border: 2 lines (top + bottom)
	// - Horizontal divider: 1 line
	// - Margins: 4 lines for spacing
	reservedHeight := 10

	availableHeight := terminalHeight - reservedHeight
	if availableHeight < minQuadrantHeight*2 {
		height = minQuadrantHeight
	} else {
		height = availableHeight / 2
	}

	// Reserve space for:
	// - Matrix border: 4 chars (left + right padding)
	// - Vertical divider: 1 char
	// - Margins: 6 chars
	reservedWidth := 11

	availableWidth := terminalWidth - reservedWidth
	if availableWidth < minQuadrantWidth*2 {
		width = minQuadrantWidth
	} else {
		width = availableWidth / 2
	}

	return width, height
}

// calculateDisplayLimit determines how many todos to show based on quadrant height
func calculateDisplayLimit(quadrantHeight int) int {
	// Reserve 2 lines for title + spacing
	// Reserve 1 line for potential "... and X more" message
	availableLines := quadrantHeight - 3

	// Each todo takes 1 line
	// Ensure at least 3 todos are shown
	if availableLines < 3 {
		return 3
	}

	return availableLines
}

// createVerticalDivider creates a vertical divider that spans the given height
func createVerticalDivider(height int) string {
	var divider strings.Builder
	for i := 0; i < height; i++ {
		divider.WriteString("│\n")
	}
	return dividerStyle.Render(strings.TrimSuffix(divider.String(), "\n"))
}

// colorizeDescription replaces project and context tags with colored versions
func colorizeDescription(description string) string {
	// Colorize project tags (+tag) with bold styling
	description = projectTagPattern.ReplaceAllStringFunc(description, func(match string) string {
		tag := match[1:] // Remove the + prefix
		color := HashColor(tag)
		style := lipgloss.NewStyle().
			Foreground(color).
			Bold(true)
		return style.Render(match)
	})

	// Colorize context tags (@tag) with normal styling but colored
	description = contextTagPattern.ReplaceAllStringFunc(description, func(match string) string {
		tag := match[1:] // Remove the @ prefix
		color := HashColor(tag)
		style := lipgloss.NewStyle().
			Foreground(color)
		return style.Render(match)
	})

	return description
}

// renderQuadrantContent renders just the content of a quadrant (no border)
func renderQuadrantContent(title string, color lipgloss.Color, todos []todo.Todo, width, height, displayLimit int) string {
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

	// Join lines and place in exact dimensions
	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Padding(1, 2).
		AlignVertical(lipgloss.Top).
		Render(content)
}
