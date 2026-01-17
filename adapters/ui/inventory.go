package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// TagInventory represents the count of incomplete todos for a tag
type TagInventory struct {
	Tag   string
	Count int
}

// countTagInventory counts incomplete todos by project and context tags
func countTagInventory(m matrix.Matrix) (projectCounts, contextCounts map[string]int) {
	projectCounts = make(map[string]int)
	contextCounts = make(map[string]int)

	// Helper to count tags from a quadrant
	countFromQuadrant := func(todos []todo.Todo) {
		for _, t := range todos {
			if !t.IsCompleted() {
				for _, p := range t.Projects() {
					projectCounts[p]++
				}
				for _, c := range t.Contexts() {
					contextCounts[c]++
				}
			}
		}
	}

	// Count from all quadrants
	countFromQuadrant(m.DoFirst())
	countFromQuadrant(m.Schedule())
	countFromQuadrant(m.Delegate())
	countFromQuadrant(m.Eliminate())

	return projectCounts, contextCounts
}

// sortTagsByCount sorts tags by count (descending), with alphabetical tiebreaker
func sortTagsByCount(tagCounts map[string]int) []TagInventory {
	inventory := make([]TagInventory, 0, len(tagCounts))
	for tag, count := range tagCounts {
		inventory = append(inventory, TagInventory{Tag: tag, Count: count})
	}

	sort.Slice(inventory, func(i, j int) bool {
		if inventory[i].Count == inventory[j].Count {
			return inventory[i].Tag < inventory[j].Tag // alphabetical tiebreaker
		}
		return inventory[i].Count > inventory[j].Count // descending by count
	})

	return inventory
}

// renderTagInventory renders the tag inventory display for overview mode
func renderTagInventory(m matrix.Matrix, width int) string {
	projectCounts, contextCounts := countTagInventory(m)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	var output strings.Builder

	// Render projects line
	output.WriteString(labelStyle.Render("Projects (+): "))
	if len(projectCounts) == 0 {
		output.WriteString(labelStyle.Render("(none)"))
	} else {
		projectInventory := sortTagsByCount(projectCounts)
		for i, item := range projectInventory {
			tagWithPrefix := "+" + item.Tag
			color := HashColor(item.Tag)
			tagStyle := lipgloss.NewStyle().Foreground(color)

			output.WriteString(tagStyle.Render(tagWithPrefix))
			output.WriteString(countStyle.Render(fmt.Sprintf(" (%d)", item.Count)))

			if i < len(projectInventory)-1 {
				output.WriteString("  ")
			}
		}
	}
	output.WriteString("\n")

	// Render contexts line
	output.WriteString(labelStyle.Render("Contexts (@): "))
	if len(contextCounts) == 0 {
		output.WriteString(labelStyle.Render("(none)"))
	} else {
		contextInventory := sortTagsByCount(contextCounts)
		for i, item := range contextInventory {
			tagWithPrefix := "@" + item.Tag
			color := HashColor(item.Tag)
			tagStyle := lipgloss.NewStyle().Foreground(color)

			output.WriteString(tagStyle.Render(tagWithPrefix))
			output.WriteString(countStyle.Render(fmt.Sprintf(" (%d)", item.Count)))

			if i < len(contextInventory)-1 {
				output.WriteString("  ")
			}
		}
	}

	return output.String()
}

// RenderInventoryDashboard renders the digital inventory dashboard view
func RenderInventoryDashboard(m matrix.Matrix, width, height int) string {
	metrics := usecases.AnalyzeInventory(m)

	// Color palette
	okColor := lipgloss.Color("#10B981")
	highColor := lipgloss.Color("#F59E0B")
	overloadedColor := lipgloss.Color("#EF4444")
	staleColor := lipgloss.Color("#F59E0B")
	veryStaleColor := lipgloss.Color("#EF4444")

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#6366F1")).
		Padding(0, 2).
		MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#A78BFA")).
		MarginBottom(1)

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF"))

	var output strings.Builder

	// Title
	output.WriteString(titleStyle.Render("📊 Digital Inventory Dashboard"))
	output.WriteString("\n\n")

	// LEFT COLUMN: Quadrant Metrics Table
	quadrantRows := [][]string{}

	formatQuadrantRow := func(name string, active, oldest int, highThreshold, overloadThreshold int) []string {
		// Health indicator
		var health string
		var healthColor lipgloss.Color
		if active <= highThreshold {
			health = "OK"
			healthColor = okColor
		} else if active <= overloadThreshold {
			health = "HIGH"
			healthColor = highColor
		} else {
			health = "OVERLOADED"
			healthColor = overloadedColor
		}

		healthText := lipgloss.NewStyle().Foreground(healthColor).Bold(true).Render(health)

		// Age indicator
		var ageText string
		if oldest > 0 {
			ageStr := fmt.Sprintf("%dd", oldest)
			if oldest > 21 {
				ageText = lipgloss.NewStyle().Foreground(veryStaleColor).Bold(true).Render(ageStr + " (VERY STALE)")
			} else if oldest > 14 {
				ageText = lipgloss.NewStyle().Foreground(staleColor).Render(ageStr + " (STALE)")
			} else {
				ageText = ageStr
			}
		} else {
			ageText = "-"
		}

		return []string{name, fmt.Sprintf("%d", active), healthText, ageText}
	}

	quadrantRows = append(quadrantRows, formatQuadrantRow("Do First", metrics.DoFirstActive, metrics.DoFirstOldestDays, 5, 8))
	quadrantRows = append(quadrantRows, formatQuadrantRow("Schedule", metrics.ScheduleActive, metrics.ScheduleOldestDays, 3, 6))
	quadrantRows = append(quadrantRows, formatQuadrantRow("Delegate", metrics.DelegateActive, metrics.DelegateOldestDays, 3, 6))
	quadrantRows = append(quadrantRows, formatQuadrantRow("Eliminate", metrics.EliminateActive, metrics.EliminateOldestDays, 3, 6))

	quadrantTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))).
		Headers("Quadrant", "WIP", "Health", "Oldest").
		Rows(quadrantRows...)

	// TOP SECTION: Centered Quadrant Metrics
	// Use the passed width for centering, or default to 100 if not provided
	contentWidth := width
	if contentWidth <= 0 {
		contentWidth = 100
	}

	quadrantSection := lipgloss.JoinVertical(
		lipgloss.Left,
		sectionStyle.Render("Quadrant Metrics"),
		quadrantTable.String(),
		"",
		labelStyle.Render(fmt.Sprintf("Total WIP: %d", metrics.TotalActive)),
	)

	// Center the quadrant section horizontally
	topSection := lipgloss.NewStyle().
		Width(contentWidth).
		AlignHorizontal(lipgloss.Center).
		Render(quadrantSection) + "\n"

	// BOTTOM SECTION: Project and Context Tables side-by-side
	var leftColumn string  // Projects
	var rightColumn string // Contexts

	// Project Breakdown Table
	if len(metrics.ProjectBreakdown) > 0 {
		projects := make([]usecases.TagMetrics, 0, len(metrics.ProjectBreakdown))
		for _, pm := range metrics.ProjectBreakdown {
			projects = append(projects, pm)
		}
		sort.Slice(projects, func(i, j int) bool {
			if projects[i].Count == projects[j].Count {
				return projects[i].Tag < projects[j].Tag
			}
			return projects[i].Count > projects[j].Count
		})

		// Limit to top 10 projects
		maxDisplay := 10
		displayProjects := projects
		if len(projects) > maxDisplay {
			displayProjects = projects[:maxDisplay]
		}

		projectRows := [][]string{}
		for _, pm := range displayProjects {
			tagColor := HashColor(pm.Tag)
			tagStyle := lipgloss.NewStyle().Foreground(tagColor).Bold(true)

			var status string
			if pm.Count > 4 {
				status = lipgloss.NewStyle().Foreground(highColor).Bold(true).Render("HIGH")
			} else {
				status = lipgloss.NewStyle().Foreground(okColor).Render("OK")
			}

			projectRows = append(projectRows, []string{
				tagStyle.Render("+" + pm.Tag),
				fmt.Sprintf("%d", pm.Count),
				fmt.Sprintf("%dd", pm.AvgAgeDays),
				status,
			})
		}

		projectTable := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))).
			Headers("Project", "WIP", "Avg Age", "Status").
			Rows(projectRows...)

		leftColumn = sectionStyle.Render("Project Breakdown") + "\n" + projectTable.String()

		// Add "...and N more" footer if there are additional projects
		if len(projects) > maxDisplay {
			remaining := len(projects) - maxDisplay
			footerText := fmt.Sprintf("...and %d more", remaining)
			leftColumn += "\n" + labelStyle.Render(footerText)
		}
	}

	// Context Breakdown Table
	if len(metrics.ContextBreakdown) > 0 {
		contexts := make([]usecases.TagMetrics, 0, len(metrics.ContextBreakdown))
		for _, cm := range metrics.ContextBreakdown {
			contexts = append(contexts, cm)
		}
		sort.Slice(contexts, func(i, j int) bool {
			if contexts[i].Count == contexts[j].Count {
				return contexts[i].Tag < contexts[j].Tag
			}
			return contexts[i].Count > contexts[j].Count
		})

		// Limit to top 10 contexts
		maxDisplay := 10
		displayContexts := contexts
		if len(contexts) > maxDisplay {
			displayContexts = contexts[:maxDisplay]
		}

		contextRows := [][]string{}
		for _, cm := range displayContexts {
			tagColor := HashColor(cm.Tag)
			tagStyle := lipgloss.NewStyle().Foreground(tagColor).Bold(true)

			var status string
			if cm.Count > 4 {
				status = lipgloss.NewStyle().Foreground(highColor).Bold(true).Render("HIGH")
			} else {
				status = lipgloss.NewStyle().Foreground(okColor).Render("OK")
			}

			contextRows = append(contextRows, []string{
				tagStyle.Render("@" + cm.Tag),
				fmt.Sprintf("%d", cm.Count),
				fmt.Sprintf("%dd", cm.AvgAgeDays),
				status,
			})
		}

		contextTable := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#6366F1"))).
			Headers("Context", "WIP", "Avg Age", "Status").
			Rows(contextRows...)

		rightColumn = sectionStyle.Render("Context Breakdown") + "\n" + contextTable.String()

		// Add "...and N more" footer if there are additional contexts
		if len(contexts) > maxDisplay {
			remaining := len(contexts) - maxDisplay
			footerText := fmt.Sprintf("...and %d more", remaining)
			rightColumn += "\n" + labelStyle.Render(footerText)
		}
	}

	// Assemble the full layout
	output.WriteString(topSection)
	output.WriteString("\n")

	// Two-column layout for projects and contexts
	if leftColumn != "" || rightColumn != "" {
		twoColumnLayout := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftColumn,
			strings.Repeat(" ", 4), // spacer
			rightColumn,
		)

		// Center the two-column layout
		centeredTwoColumn := lipgloss.NewStyle().
			Width(contentWidth).
			AlignHorizontal(lipgloss.Center).
			Render(twoColumnLayout)

		output.WriteString(centeredTwoColumn)
		output.WriteString("\n\n")
	}

	// Throughput Section
	throughputContent := lipgloss.JoinVertical(
		lipgloss.Left,
		sectionStyle.Render("Throughput (Last 7 Days)"),
		labelStyle.Render(fmt.Sprintf("  Completed: %d items  |  Added: %d items",
			metrics.CompletedLast7Days, metrics.AddedLast7Days)),
	)

	// Warning if adding faster than completing
	if metrics.AddedLast7Days > metrics.CompletedLast7Days {
		warningStyle := lipgloss.NewStyle().
			Foreground(highColor).
			Bold(true)
		throughputContent += labelStyle.Render("  |  ") + warningStyle.Render("⚠ Adding faster than completing")
	}

	// Center throughput section
	centeredThroughput := lipgloss.NewStyle().
		Width(contentWidth).
		AlignHorizontal(lipgloss.Center).
		Render(throughputContent)

	output.WriteString(centeredThroughput)
	output.WriteString("\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6B7280")).
		Italic(true).
		MarginTop(2)

	helpText := lipgloss.NewStyle().
		Width(contentWidth).
		AlignHorizontal(lipgloss.Center).
		Render(helpStyle.Render("↑/↓/PgUp/PgDn: Scroll  |  'i' or ESC: Return to overview"))

	output.WriteString(helpText)

	// The viewport will handle the final content, so we don't need to center with lipgloss.Place
	return output.String()
}
