package ui

import (
	"sort"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// extractAllTags extracts all unique project and context tags from the matrix
func extractAllTags(m matrix.Matrix) (projects []string, contexts []string) {
	projectSet := make(map[string]bool)
	contextSet := make(map[string]bool)

	// Helper to extract tags from a quadrant
	extractFromQuadrant := func(todos []todo.Todo) {
		for _, t := range todos {
			for _, p := range t.Projects() {
				projectSet[p] = true
			}
			for _, c := range t.Contexts() {
				contextSet[c] = true
			}
		}
	}

	// Extract from all quadrants
	extractFromQuadrant(m.DoFirst())
	extractFromQuadrant(m.Schedule())
	extractFromQuadrant(m.Delegate())
	extractFromQuadrant(m.Eliminate())

	// Convert sets to sorted slices
	projects = make([]string, 0, len(projectSet))
	for p := range projectSet {
		projects = append(projects, p)
	}
	sort.Strings(projects)

	contexts = make([]string, 0, len(contextSet))
	for c := range contextSet {
		contexts = append(contexts, c)
	}
	sort.Strings(contexts)

	return projects, contexts
}
