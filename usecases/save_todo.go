package usecases

import (
	"fmt"
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// TodoWriter is the interface for writing todos
type TodoWriter interface {
	SaveTodo(line string) error
	ReplaceAll(content string) error
}

// saveTodo formats and appends a single todo to the file (private helper)
func saveTodo(writer TodoWriter, t todo.Todo) error {
	line := FormatTodo(t)
	return writer.SaveTodo(line)
}

// FormatTodo converts a Todo to todo.txt format
func FormatTodo(t todo.Todo) string {
	// Format: x YYYY-MM-DD (PRIORITY) Description +project @context
	var result string

	// Add completion marker if completed
	if t.IsCompleted() {
		completionDate := time.Now().Format("2006-01-02")
		result = fmt.Sprintf("x %s ", completionDate)
	}

	// Add priority if present
	if t.Priority() != todo.PriorityNone {
		priorityLetter := priorityToString(t.Priority())
		result += fmt.Sprintf("(%s) ", priorityLetter)
	}

	// Add description
	result += t.Description()

	// Note: tags are already included in the description
	// The todo.txt format has tags inline with the description

	result += "\n"
	return result
}

// saveAllTodos writes all todos from the matrix back to the file (private helper)
func saveAllTodos(writer TodoWriter, m matrix.Matrix) error {
	var content string

	// Format all todos from all quadrants
	for _, t := range m.AllTodos() {
		content += FormatTodo(t)
	}

	return writer.ReplaceAll(content)
}

// priorityToString converts a Priority to its string representation
func priorityToString(p todo.Priority) string {
	switch p {
	case todo.PriorityA:
		return "A"
	case todo.PriorityB:
		return "B"
	case todo.PriorityC:
		return "C"
	case todo.PriorityD:
		return "D"
	default:
		return ""
	}
}
