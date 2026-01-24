package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
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

// buildDescriptionWithTags reconstructs the full description with tags for display
func buildDescriptionWithTags(t todo.Todo) string {
	var b strings.Builder
	b.WriteString(t.Description())

	// Add project tags
	for _, project := range t.Projects() {
		b.WriteString(" +")
		b.WriteString(project)
	}

	// Add context tags
	for _, context := range t.Contexts() {
		b.WriteString(" @")
		b.WriteString(context)
	}

	return b.String()
}

var (
	// Color palette - Eisenhower matrix themed
	urgentImportantColor = lipgloss.Color("#FF6B6B") // Red - Do First
	importantColor       = lipgloss.Color("#4ECDC4") // Teal - Schedule
	urgentColor          = lipgloss.Color("#FFE66D") // Yellow - Delegate
	neitherColor         = lipgloss.Color("#95E1D3") // Light green - Eliminate

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7B68EE")).
			Padding(0, 2).
			MarginBottom(1)

	matrixBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7B68EE")).
			Padding(0)

	activeTodoStyle = lipgloss.NewStyle().
			Foreground(TextPrimary)

	completedTodoStyle = lipgloss.NewStyle().
				Foreground(CompletedColor).
				Strikethrough(true)

	emptyStyle = lipgloss.NewStyle().
			Foreground(EmptyColor).
			Italic(true)

	dividerStyle = lipgloss.NewStyle().
			Foreground(BorderColor)
)

// RenderMatrix renders the Eisenhower matrix as a string with optional file path header
// terminalWidth and terminalHeight are optional (0 = use defaults)
// Note: Filtering should be applied at the domain level before calling this function
func RenderMatrix(m matrix.Matrix, filePath string, terminalWidth, terminalHeight int) string {
	return RenderMatrixWithFilterHint(m, filePath, terminalWidth, terminalHeight, "", false)
}

// RenderMatrixWithFilterHint renders the matrix with a filter hint for help text
// The activeFilter parameter is only used for display (help text), not for actual filtering
func RenderMatrixWithFilterHint(m matrix.Matrix, filePath string, terminalWidth, terminalHeight int, activeFilter string, readOnly bool) string {
	// Calculate quadrant dimensions based on terminal size
	quadrantWidth, quadrantHeight := calculateQuadrantDimensions(terminalWidth, terminalHeight)
	// For overview mode, always show top 5 todos per quadrant (cleaner, more consistent)
	displayLimit := 5

	var output strings.Builder

	// Render header if file path provided
	if filePath != "" {
		header := headerStyle.Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Get todos from matrix
	doFirstTodos := m.DoFirst()
	scheduleTodos := m.Schedule()
	delegateTodos := m.Delegate()
	eliminateTodos := m.Eliminate()

	// Render quadrant contents
	doFirst := renderQuadrantContent("Do First", urgentImportantColor, doFirstTodos, quadrantWidth, quadrantHeight, displayLimit, 1)
	schedule := renderQuadrantContent("Schedule", importantColor, scheduleTodos, quadrantWidth, quadrantHeight, displayLimit, 2)
	delegate := renderQuadrantContent("Delegate", urgentColor, delegateTodos, quadrantWidth, quadrantHeight, displayLimit, 3)
	eliminate := renderQuadrantContent("Eliminate", neitherColor, eliminateTodos, quadrantWidth, quadrantHeight, displayLimit, 4)

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
	matrixView := matrixBorder.Render(matrixContent)
	output.WriteString(matrixView)
	output.WriteString("\n\n")

	// Add tag inventory
	inventory := renderTagInventory(m, terminalWidth)
	output.WriteString(inventory)
	output.WriteString("\n\n")

	// Add help text with filter instructions
	// Note: Overview mode help text doesn't include editing commands, so it's the same for read-only and normal mode
	var helpText string
	if activeFilter != "" {
		helpText = renderHelp("Press 1/2/3/4 to focus on a quadrant", "Press c to clear filter", "Press i for inventory", "Press q to quit")
	} else {
		helpText = renderHelp("Press 1/2/3/4 to focus on a quadrant", "Press f to filter", "Press i for inventory", "Press q to quit")
	}
	output.WriteString(helpText)

	return output.String()
}

// RenderFocusedQuadrant renders a single quadrant in fullscreen focus mode
func RenderFocusedQuadrant(todos []todo.Todo, title string, color lipgloss.Color, filePath string, selectedIndex, terminalWidth, terminalHeight int) string {
	var output strings.Builder

	// Render file path header with full width and center alignment
	if filePath != "" {
		header := headerStyle.
			Width(terminalWidth).
			Align(lipgloss.Center).
			Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Calculate display limit for focus mode
	// Reserve: header (3), title (2), help text (2), margins (2) = 9 lines
	displayLimit := max(terminalHeight-9, 5)

	// Render prominent quadrant title with gradient background
	titleText := fmt.Sprintf(" %s ", title)
	gradientTitle := GradientBackground(titleText, color, lightenColor(color, 0.5))
	
	focusTitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Bold(true).
		Render(gradientTitle)
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

			// Build full description with tags for display
			description := buildDescriptionWithTags(t)
			
			// Colorize tags in description
			description = colorizeDescription(description)

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

			// Highlight selected todo
			if i == selectedIndex {
				selectedStyle := lipgloss.NewStyle().
					Background(SelectionBg).
					Bold(true)
				todoLine = selectedStyle.Render(todoLine)
			}

			lines = append(lines, todoLine)
		}
	}

	todosContent := strings.Join(lines, "\n")
	output.WriteString(todosContent)
	output.WriteString("\n\n")

	// Render help text at bottom
	helpText := renderHelp("↑↓/w/s navigate", "Space to toggle", "a to add", "o to open URL", "Press 1-4 to jump", "m to move", "ESC to return")
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

// createVerticalDivider creates a vertical divider that spans the given height
func createVerticalDivider(height int) string {
	var divider strings.Builder
	for range height {
		divider.WriteString("│\n")
	}
	return dividerStyle.Render(strings.TrimSuffix(divider.String(), "\n"))
}

// formatDate formats a date with consistent friendly formatting
// Returns "today", "yesterday", or "N days ago" for all dates
func formatDate(date *time.Time) string {
	if date == nil {
		return ""
	}

	now := time.Now()
	// Normalize both to start of day for comparison
	dateDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	daysDiff := int(today.Sub(dateDay).Hours() / 24)

	switch daysDiff {
	case 0:
		return "today"
	case 1:
		return "yesterday"
	default:
		// Always use "N days ago" for consistency
		return fmt.Sprintf("%d days ago", daysDiff)
	}
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
func renderQuadrantContent(title string, color lipgloss.Color, todos []todo.Todo, width, height, displayLimit, quadrantNumber int) string {
	var lines []string

	// Calculate stats
	totalTasks := len(todos)
	completedTasks := 0
	for _, t := range todos {
		if t.IsCompleted() {
			completedTasks++
		}
	}

	// Title with gradient-style background block
	taskWord := "tasks"
	if totalTasks == 1 {
		taskWord = "task"
	}
	
	// Create title with gradient background
	titleText := fmt.Sprintf(" %s ", title)
	gradientTitle := GradientBackground(titleText, color, lightenColor(color, 0.5))
	statsText := fmt.Sprintf(" %d %s · %d completed ", totalTasks, taskWord, completedTasks)
	
	statsBlock := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(statsText)
	
	quadrantTitle := lipgloss.JoinHorizontal(lipgloss.Top, gradientTitle, statsBlock)
	lines = append(lines, quadrantTitle, "") // spacing

	if len(todos) == 0 {
		lines = append(lines, emptyStyle.Render("  (no tasks)"))
	} else {
		for i, t := range todos {
			if i >= displayLimit {
				remaining := len(todos) - displayLimit
				hint := fmt.Sprintf("  ... and %d more (press %d to view)", remaining, quadrantNumber)
				lines = append(lines, emptyStyle.Render(hint))
				break
			}

			// Build full description with tags for display
			description := buildDescriptionWithTags(t)

			// Colorize tags in description
			description = colorizeDescription(description)

			var todoLine string
			if t.IsCompleted() {
				todoLine = completedTodoStyle.Render("✓ ") + description
			} else {
				todoLine = activeTodoStyle.Render("• ") + description
			}

			// Append due date if present
			dueDateText, isOverdue := formatDueDateWithOverdue(t.DueDate(), time.Now())
			if dueDateText != "" {
				if isOverdue {
					// Red with ! prefix for overdue
					overdueStyle := lipgloss.NewStyle().Foreground(overdueColor)
					todoLine += overdueStyle.Render(fmt.Sprintf(" ! due: %s", dueDateText))
				} else {
					// Cyan for upcoming due dates
					dueDateStyle := lipgloss.NewStyle().Foreground(dueDateColor)
					todoLine += dueDateStyle.Render(fmt.Sprintf(" due: %s", dueDateText))
				}
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

// RenderFocusedQuadrantWithTable renders a quadrant in focus mode using a table
func RenderFocusedQuadrantWithTable(
	todos []todo.Todo,
	title string,
	color lipgloss.Color,
	filePath string,
	todoTable table.Model,
	terminalWidth, terminalHeight int,
	readOnly bool,
) string {
	var output strings.Builder

	// Render file path header with full width and center alignment
	if filePath != "" {
		header := headerStyle.
			Width(terminalWidth).
			Align(lipgloss.Center).
			Render("File: " + filePath)
		output.WriteString(header)
		output.WriteString("\n\n")
	}

	// Render prominent quadrant title with gradient background
	titleText := fmt.Sprintf(" %s ", title)
	gradientTitle := GradientBackground(titleText, color, lightenColor(color, 0.5))
	
	focusTitle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Bold(true).
		Render(gradientTitle)
	output.WriteString(focusTitle)
	output.WriteString("\n\n")

	// Render table or empty message
	if len(todos) == 0 {
		emptyMsg := emptyStyle.Render("(no tasks)")
		centeredMsg := lipgloss.NewStyle().
			Width(terminalWidth).
			Align(lipgloss.Center).
			Render(emptyMsg)
		output.WriteString(centeredMsg)
	} else {
		// Render the table
		output.WriteString(todoTable.View())
		output.WriteString("\n\n")

		// Render detail pane for selected todo
		selectedIndex := todoTable.Cursor()
		if selectedIndex >= 0 && selectedIndex < len(todos) {
			detailPane := renderTodoDetailPane(todos[selectedIndex], terminalWidth)
			output.WriteString(detailPane)
		}
	}

	output.WriteString("\n\n")

	// Render help text at bottom
	var helpText string
	if readOnly {
		// Read-only mode: only show viewing/navigation commands
		helpText = renderHelp("Press 1-4 to jump", "Press ESC to return", "Press q to quit")
	} else {
		// Normal mode: show all commands including editing
		helpText = renderHelp("a to add", "o to open URL", "d to archive", "1-4 to jump", "m to move", "ESC to return")
	}
	centeredHelp := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(terminalWidth).
		Render(helpText)
	output.WriteString(centeredHelp)

	return output.String()
}

// buildTodoTable creates a table.Model from a list of todos
func buildTodoTable(todos []todo.Todo, terminalWidth, terminalHeight, selectedIndex int) table.Model {
	// Calculate column widths based on terminal width
	// Reserve some width for borders, padding, etc.
	availableWidth := max(terminalWidth-10, 80)

	// Define column widths (these are approximate ratios)
	// Task gets most of the space, other columns get fixed widths
	projectsWidth := 15
	contextsWidth := 15
	createdWidth := 12
	completedWidth := 12
	dueDateWidth := 12
	taskWidth := max(availableWidth-projectsWidth-contextsWidth-createdWidth-completedWidth-dueDateWidth, 30)

	columns := []table.Column{
		{Title: "Task", Width: taskWidth},
		{Title: "Projects", Width: projectsWidth},
		{Title: "Contexts", Width: contextsWidth},
		{Title: "Created", Width: createdWidth},
		{Title: "Completed", Width: completedWidth},
		{Title: "Due Date", Width: dueDateWidth},
	}

	// Build rows from todos
	rows := make([]table.Row, len(todos))
	for i, t := range todos {
		// Task: description is already clean (tags extracted by parser)
		// Table will truncate if too long; full description shown in detail pane
		taskDesc := t.Description()

		// Projects: comma-separated list
		projects := strings.Join(t.Projects(), ", ")
		if projects == "" {
			projects = "-"
		}

		// Contexts: comma-separated list
		contexts := strings.Join(t.Contexts(), ", ")
		if contexts == "" {
			contexts = "-"
		}

		// Created: friendly date format
		created := formatDate(t.CreationDate())
		if created == "" {
			created = "-"
		}

		// Completed: friendly date format (empty for active todos)
		completed := formatDate(t.CompletionDate())
		if completed == "" {
			completed = "-"
		}

		// Due Date: formatted with overdue indicator
		dueDateText, isOverdue := formatDueDateWithOverdue(t.DueDate(), time.Now())
		var dueDate string
		if dueDateText != "" {
			if isOverdue {
				// Add ! prefix for overdue (red styling limited in table)
				dueDate = "! " + dueDateText
			} else {
				dueDate = dueDateText
			}
		} else {
			dueDate = "-"
		}

		rows[i] = table.Row{taskDesc, projects, contexts, created, completed, dueDate}
	}

	// Track which rows are completed for styling
	completedRows := make(map[int]bool)
	for i, t := range todos {
		if t.IsCompleted() {
			completedRows[i] = true
		}
	}

	// Calculate table height - should fit within terminal
	// Reserve space for header, title, help text, etc.
	tableHeight := max(terminalHeight-15, 5)

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	// Style the table
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(BorderColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(TextPrimary).
		Background(SelectionBg).
		Bold(true)

	// Style completed rows with green foreground
	// Note: The bubbles table doesn't directly support per-row styling,
	// so we'll use a green foreground for the entire cell style when rendering
	// For now, we'll apply green to all cells if the row is completed
	// This is a limitation of the table component - better per-row styling would require custom rendering

	t.SetStyles(s)

	// Set cursor to the selected index
	if selectedIndex >= 0 && selectedIndex < len(rows) {
		t.SetCursor(selectedIndex)
	}

	return t
}

// renderTodoDetailPane renders a detail pane showing full information for a todo
func renderTodoDetailPane(t todo.Todo, terminalWidth int) string {
	var output strings.Builder

	// Style definitions with adaptive colors
	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#666666", Dark: "#888888"}).
		Bold(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"})

	// Border style for the detail pane
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#AAAAAA", Dark: "#555555"}).
		Padding(0, 1).
		Width(terminalWidth - 4)

	var details strings.Builder

	// Full description
	details.WriteString(labelStyle.Render("Description: "))
	details.WriteString(valueStyle.Render(t.Description()))
	details.WriteString("\n")

	// Projects
	if len(t.Projects()) > 0 {
		details.WriteString(labelStyle.Render("Projects: "))
		for i, p := range t.Projects() {
			color := HashColor(p)
			projectStyle := lipgloss.NewStyle().Foreground(color)
			details.WriteString(projectStyle.Render("+" + p))
			if i < len(t.Projects())-1 {
				details.WriteString(", ")
			}
		}
		details.WriteString("\n")
	}

	// Contexts
	if len(t.Contexts()) > 0 {
		details.WriteString(labelStyle.Render("Contexts: "))
		for i, c := range t.Contexts() {
			color := HashColor(c)
			contextStyle := lipgloss.NewStyle().Foreground(color)
			details.WriteString(contextStyle.Render("@" + c))
			if i < len(t.Contexts())-1 {
				details.WriteString(", ")
			}
		}
		details.WriteString("\n")
	}

	// Created date
	if created := formatDate(t.CreationDate()); created != "" {
		details.WriteString(labelStyle.Render("Created: "))
		details.WriteString(valueStyle.Render(created))
		details.WriteString("\n")
	}

	// Completion date
	if t.IsCompleted() {
		if completed := formatDate(t.CompletionDate()); completed != "" {
			details.WriteString(labelStyle.Render("Completed: "))
			details.WriteString(valueStyle.Render(completed))
			details.WriteString("\n")
		}
	}

	// Due date
	if t.DueDate() != nil {
		dueDateText, isOverdue := formatDueDateWithOverdue(t.DueDate(), time.Now())
		details.WriteString(labelStyle.Render("Due: "))
		if isOverdue {
			overdueStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true)
			details.WriteString(overdueStyle.Render(dueDateText + " (OVERDUE)"))
		} else {
			dueDateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00CED1"))
			details.WriteString(dueDateStyle.Render(dueDateText))
		}
		details.WriteString("\n")
	}

	output.WriteString(borderStyle.Render(details.String()))
	return output.String()
}

// RenderMoveOverlay renders an overlay for move mode
func RenderMoveOverlay(terminalWidth, terminalHeight int) string {
	// Build overlay content
	content := lipgloss.NewStyle().
		Bold(true).
		Render("Move to quadrant:") + "\n\n"

	content += "  1. Do First\n"
	content += "  2. Schedule\n"
	content += "  3. Delegate\n"
	content += "  4. Eliminate\n\n"

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

// RenderURLSelectionOverlay renders the URL selection dialog overlay
func RenderURLSelectionOverlay(urls []string, selectedIndex, terminalWidth, terminalHeight int) string {
	var output strings.Builder

	// Build overlay content
	titleStyle := lipgloss.NewStyle().Bold(true)
	output.WriteString(titleStyle.Render("Select URL to open:"))
	output.WriteString("\n\n")

	// Render each URL with selection indicator
	for i, url := range urls {
		var line string
		if i == selectedIndex {
			// Highlight selected URL
			line = lipgloss.NewStyle().
				Background(SelectionBg).
				Foreground(TextPrimary).
				Padding(0, 1).
				Render("> " + url)
		} else {
			line = "  " + url
		}
		output.WriteString(line)
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(lipgloss.NewStyle().
		Foreground(TextSecondary).
		Italic(true).
		Render("Press Enter to open • ESC to cancel"))

	content := output.String()

	// Create bordered box with adaptive width
	boxWidth := 60
	// Adjust width based on longest URL
	for _, url := range urls {
		if len(url)+4 > boxWidth {
			boxWidth = min(len(url)+4, terminalWidth-10)
		}
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BorderAccent).
		Padding(1, 2).
		Width(boxWidth)

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
