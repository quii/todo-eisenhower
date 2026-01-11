package ui_test

import (
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestRenderMatrix(t *testing.T) {
	t.Run("renders matrix with todos in all quadrants", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Plan quarterly goals", todo.PriorityB),
			todo.New("Reply to emails", todo.PriorityC),
			todo.New("Clean workspace", todo.PriorityD),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "")

		// Check that quadrant labels are present
		assertContains(t, output, "DO FIRST")
		assertContains(t, output, "SCHEDULE")
		assertContains(t, output, "DELEGATE")
		assertContains(t, output, "ELIMINATE")

		// Check that todos appear in the output
		assertContains(t, output, "Fix critical bug")
		assertContains(t, output, "Plan quarterly goals")
		assertContains(t, output, "Reply to emails")
		assertContains(t, output, "Clean workspace")
	})

	t.Run("renders matrix with empty quadrants", func(t *testing.T) {
		m := matrix.New([]todo.Todo{})

		output := ui.RenderMatrix(m, "")

		// Labels should still be present
		assertContains(t, output, "DO FIRST")
		assertContains(t, output, "SCHEDULE")
		assertContains(t, output, "DELEGATE")
		assertContains(t, output, "ELIMINATE")
	})

	t.Run("renders multiple todos in same quadrant", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("First urgent task", todo.PriorityA),
			todo.New("Second urgent task", todo.PriorityA),
			todo.New("Third urgent task", todo.PriorityA),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "")

		assertContains(t, output, "First urgent task")
		assertContains(t, output, "Second urgent task")
		assertContains(t, output, "Third urgent task")
	})

	t.Run("renders completed todos with visual distinction", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Active task", todo.PriorityA),
			todo.NewCompleted("Completed task", todo.PriorityA),
		}
		m := matrix.New(todos)

		output := ui.RenderMatrix(m, "")

		// Both should appear in output
		assertContains(t, output, "Active task")
		assertContains(t, output, "Completed task")

		// Completed should have strikethrough or different marker
		// We'll use "✓" prefix for completed tasks
		assertContains(t, output, "✓")
	})

	t.Run("renders file path header", func(t *testing.T) {
		m := matrix.New([]todo.Todo{})
		filePath := "/Users/chris/projects/todo.txt"

		output := ui.RenderMatrix(m, filePath)

		assertContains(t, output, "File:")
		assertContains(t, output, filePath)
	})
}

func assertContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("expected output to contain %q, but it didn't.\nOutput:\n%s", expected, output)
	}
}
