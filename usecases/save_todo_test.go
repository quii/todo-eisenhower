package usecases_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/parser"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestFormatTodo(t *testing.T) {
	t.Run("formats simple todo with priority", func(t *testing.T) {
		is := is.New(t)
		td := todo.New("Fix critical bug", todo.PriorityA)

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) Fix critical bug\n")
	})

	t.Run("formats todo without priority", func(t *testing.T) {
		is := is.New(t)
		td := todo.New("No priority task", todo.PriorityNone)

		result := usecases.FormatTodo(td)

		is.Equal(result, "No priority task\n")
	})

	t.Run("formats completed todo with priority", func(t *testing.T) {
		is := is.New(t)
		completionDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
		td := todo.NewCompleted("Completed task", todo.PriorityA, &completionDate)

		result := usecases.FormatTodo(td)

		is.Equal(result, "x 2026-01-18 (A) Completed task\n")
	})

	t.Run("formats todo with projects", func(t *testing.T) {
		is := is.New(t)
		td := todo.NewWithTags("Deploy feature", todo.PriorityA, []string{"WebApp"}, nil)

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) Deploy feature +WebApp\n")
	})

	t.Run("formats todo with contexts", func(t *testing.T) {
		is := is.New(t)
		td := todo.NewWithTags("Call client", todo.PriorityB, nil, []string{"phone"})

		result := usecases.FormatTodo(td)

		is.Equal(result, "(B) Call client @phone\n")
	})

	t.Run("formats todo with multiple projects and contexts", func(t *testing.T) {
		is := is.New(t)
		td := todo.NewWithTags(
			"Write report",
			todo.PriorityA,
			[]string{"Work", "Q1Goals"},
			[]string{"office", "computer"},
		)

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) Write report +Work +Q1Goals @office @computer\n")
	})

	t.Run("formats completed todo with tags", func(t *testing.T) {
		is := is.New(t)
		completionDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
		td := todo.NewCompletedWithTags("Finished task", todo.PriorityA, &completionDate, []string{"Project"}, []string{"office"})

		result := usecases.FormatTodo(td)

		is.Equal(result, "x 2026-01-18 (A) Finished task +Project @office\n")
	})

	t.Run("formats todo with creation date", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		td := todo.NewWithCreationDate("Task with date", todo.PriorityA, &creationDate)

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) 2026-01-15 Task with date\n")
	})

	t.Run("formats completed todo with both completion and creation dates", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		completionDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
		td := todo.NewCompletedWithDates("Task with both dates", todo.PriorityA, &completionDate, &creationDate)

		result := usecases.FormatTodo(td)

		is.Equal(result, "x 2026-01-18 2026-01-15 (A) Task with both dates\n")
	})

	// Round-trip tests: format -> parse -> format should be identical
	t.Run("round-trip: simple todo", func(t *testing.T) {
		is := is.New(t)
		original := todo.New("Test task", todo.PriorityA)

		formatted := usecases.FormatTodo(original)
		parsed, err := parser.Parse(strings.NewReader(formatted))

		is.NoErr(err)
		is.Equal(len(parsed), 1)

		reformatted := usecases.FormatTodo(parsed[0])
		is.Equal(reformatted, formatted) // should be identical
	})

	t.Run("round-trip: todo with tags", func(t *testing.T) {
		is := is.New(t)
		original := todo.NewWithTags(
			"Complex task",
			todo.PriorityB,
			[]string{"Project1", "Project2"},
			[]string{"office", "computer"},
		)

		formatted := usecases.FormatTodo(original)
		parsed, err := parser.Parse(strings.NewReader(formatted))

		is.NoErr(err)
		is.Equal(len(parsed), 1)

		reformatted := usecases.FormatTodo(parsed[0])
		is.Equal(reformatted, formatted) // should be identical
	})

	t.Run("round-trip: completed todo", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		completionDate := time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC)
		original := todo.NewCompletedWithTagsAndDates(
			"Done task",
			todo.PriorityA,
			&completionDate,
			&creationDate,
			[]string{"Work"},
			[]string{"done"},
		)

		formatted := usecases.FormatTodo(original)
		parsed, err := parser.Parse(strings.NewReader(formatted))

		is.NoErr(err)
		is.Equal(len(parsed), 1)

		reformatted := usecases.FormatTodo(parsed[0])
		is.Equal(reformatted, formatted) // should be identical
	})

	// Edge cases for CONDITIONALS mutations
	t.Run("formats todo with empty projects slice", func(t *testing.T) {
		is := is.New(t)
		td := todo.NewWithTags("Task", todo.PriorityA, []string{}, nil)

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) Task\n") // no trailing space for empty tags
	})

	t.Run("formats todo with empty contexts slice", func(t *testing.T) {
		is := is.New(t)
		td := todo.NewWithTags("Task", todo.PriorityA, nil, []string{})

		result := usecases.FormatTodo(td)

		is.Equal(result, "(A) Task\n") // no trailing space for empty tags
	})
}
