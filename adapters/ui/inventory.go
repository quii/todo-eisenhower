package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
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
