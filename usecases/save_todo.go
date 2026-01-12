package usecases

import (
	"fmt"

	"github.com/quii/todo-eisenhower/domain/todo"
)

// TodoWriter is the interface for writing todos
type TodoWriter interface {
	SaveTodo(line string) error
}

// SaveTodo formats a todo according to todo.txt format and writes it to the sink
func SaveTodo(writer TodoWriter, t todo.Todo) error {
	line := FormatTodo(t)
	return writer.SaveTodo(line)
}

// FormatTodo converts a Todo to todo.txt format
func FormatTodo(t todo.Todo) string {
	// Format: (PRIORITY) Description +project @context
	var result string

	// Add priority if present
	if t.Priority() != todo.PriorityNone {
		priorityLetter := priorityToString(t.Priority())
		result = fmt.Sprintf("(%s) ", priorityLetter)
	}

	// Add description
	result += t.Description()

	// Note: tags are already included in the description
	// The todo.txt format has tags inline with the description

	result += "\n"
	return result
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
