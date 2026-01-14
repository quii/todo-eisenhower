package usecases

import (
	"regexp"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// AddTodo creates a new todo and adds it to the matrix
func AddTodo(writer TodoWriter, m matrix.Matrix, description string, priority todo.Priority) (matrix.Matrix, error) {
	// Extract tags from description
	projects := extractTagsFromDescription(description, `\+(\w+)`)
	contexts := extractTagsFromDescription(description, `@(\w+)`)

	// Create the todo using rich domain model
	var newTodo todo.Todo
	if len(projects) > 0 || len(contexts) > 0 {
		newTodo = todo.NewWithTags(description, priority, projects, contexts)
	} else {
		newTodo = todo.New(description, priority)
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
