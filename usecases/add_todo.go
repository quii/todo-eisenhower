// Package usecases contains application use cases and business logic.
package usecases

import (
	"regexp"
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// AddTodo creates a new todo and adds it to the matrix
func AddTodo(writer TodoWriter, m matrix.Matrix, description string, priority todo.Priority) (matrix.Matrix, error) {
	// Extract tags from description
	projects := extractTagsFromDescription(description, `\+(\w+)`)
	contexts := extractTagsFromDescription(description, `@(\w+)`)

	// Strip tags from description (they're stored separately)
	cleanDescription := stripTagsFromDescription(description)

	// Set creation date to now
	now := time.Now()
	creationDate := &now

	// Create the todo using rich domain model with creation date
	var newTodo todo.Todo
	if len(projects) > 0 || len(contexts) > 0 {
		newTodo = todo.NewWithTagsAndDates(cleanDescription, priority, creationDate, projects, contexts)
	} else {
		newTodo = todo.NewWithCreationDate(cleanDescription, priority, creationDate)
	}

	// Add todo to matrix
	updatedMatrix := m.AddTodo(newTodo)

	// Persist changes (implementation detail)
	err := saveTodo(writer, newTodo)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	return updatedMatrix, nil
}

// extractTagsFromDescription extracts tags matching the given regex pattern
func extractTagsFromDescription(description string, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(description, -1)
	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}
	return tags
}

// stripTagsFromDescription removes project and context tags from description
func stripTagsFromDescription(description string) string {
	// Remove project tags (+tag)
	projectPattern := regexp.MustCompile(`\+\w+`)
	description = projectPattern.ReplaceAllString(description, "")

	// Remove context tags (@tag)
	contextPattern := regexp.MustCompile(`@\w+`)
	description = contextPattern.ReplaceAllString(description, "")

	// Clean up extra whitespace
	re := regexp.MustCompile(`\s+`)
	description = re.ReplaceAllString(description, " ")

	return regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(description, "")
}
