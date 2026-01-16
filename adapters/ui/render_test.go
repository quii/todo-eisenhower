package ui_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestRenderMatrix(t *testing.T) {
	t.Run("renders matrix with todos in all quadrants", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Plan quarterly goals", todo.PriorityB),
			todo.New("Reply to emails", todo.PriorityC),
			todo.New("Clean workspace", todo.PriorityD),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "", 0, 0)

		// Check that quadrant labels are present
		is.True(strings.Contains(output, "DO FIRST"))  // expected view to contain 'DO FIRST'
		is.True(strings.Contains(output, "SCHEDULE"))  // expected view to contain 'SCHEDULE'
		is.True(strings.Contains(output, "DELEGATE"))  // expected view to contain 'DELEGATE'
		is.True(strings.Contains(output, "ELIMINATE")) // expected view to contain 'ELIMINATE'

		// Check that todos appear in the output
		is.True(strings.Contains(output, "Fix critical bug"))     // expected view to contain todo
		is.True(strings.Contains(output, "Plan quarterly goals")) // expected view to contain todo
		is.True(strings.Contains(output, "Reply to emails"))      // expected view to contain todo
		is.True(strings.Contains(output, "Clean workspace"))      // expected view to contain todo
	})

	t.Run("renders matrix with empty quadrants", func(t *testing.T) {
		is := is.New(t)
		m := matrix.New([]todo.Todo{})

		output := ui.RenderMatrix(m, "", 0, 0)

		// Labels should still be present
		is.True(strings.Contains(output, "DO FIRST"))  // expected view to contain 'DO FIRST'
		is.True(strings.Contains(output, "SCHEDULE"))  // expected view to contain 'SCHEDULE'
		is.True(strings.Contains(output, "DELEGATE"))  // expected view to contain 'DELEGATE'
		is.True(strings.Contains(output, "ELIMINATE")) // expected view to contain 'ELIMINATE'
	})

	t.Run("renders multiple todos in same quadrant", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("First urgent task", todo.PriorityA),
			todo.New("Second urgent task", todo.PriorityA),
			todo.New("Third urgent task", todo.PriorityA),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "", 0, 0)

		is.True(strings.Contains(output, "First urgent task"))  // expected view to contain todo
		is.True(strings.Contains(output, "Second urgent task")) // expected view to contain todo
		is.True(strings.Contains(output, "Third urgent task"))  // expected view to contain todo
	})

	t.Run("renders completed todos with visual distinction", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Active task", todo.PriorityA),
			todo.NewCompleted("Completed task", todo.PriorityA, nil),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "", 0, 0)

		// Both should appear in output
		is.True(strings.Contains(output, "Active task"))    // expected view to contain active task
		is.True(strings.Contains(output, "Completed task")) // expected view to contain completed task

		// Completed should have strikethrough or different marker
		// We'll use "✓" prefix for completed tasks
		is.True(strings.Contains(output, "✓")) // expected view to show checkmark for completed task
	})

	t.Run("renders file path header", func(t *testing.T) {
		is := is.New(t)
		m := matrix.New([]todo.Todo{})
		filePath := "/Users/chris/projects/todo.txt"

		output := ui.RenderMatrix(m, filePath, 0, 0)

		is.True(strings.Contains(output, "File:"))  // expected view to contain 'File:' label
		is.True(strings.Contains(output, filePath)) // expected view to contain file path
	})
}
