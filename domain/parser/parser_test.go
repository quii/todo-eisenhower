package parser_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/parser"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestParse(t *testing.T) {
	t.Run("parses single todo with priority A", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("(A) Fix critical bug")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "Fix critical bug", todo.PriorityA, false)
	})

	t.Run("parses todo without priority", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("No priority task")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "No priority task", todo.PriorityNone, false)
	})

	t.Run("parses completed todo", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("x (A) 2026-01-11 Completed task")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "Completed task", todo.PriorityA, true)
	})

	t.Run("parses multiple todos with different priorities", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace`)

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 4) // expected 4 todos

		assertTodo(is, todos[0], "Fix critical bug", todo.PriorityA, false)
		assertTodo(is, todos[1], "Plan quarterly goals", todo.PriorityB, false)
		assertTodo(is, todos[2], "Reply to emails", todo.PriorityC, false)
		assertTodo(is, todos[3], "Clean workspace", todo.PriorityD, false)
	})

	t.Run("parses empty input", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 0) // expected 0 todos
	})

	t.Run("ignores empty lines", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) First task

(B) Second task`)

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 2) // expected 2 todos
	})

	t.Run("parses mixed completed and active todos", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) Active task
x (A) 2026-01-11 Completed task
(B) Another active task`)

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 3) // expected 3 todos

		assertTodo(is, todos[0], "Active task", todo.PriorityA, false)
		assertTodo(is, todos[1], "Completed task", todo.PriorityA, true)
		assertTodo(is, todos[2], "Another active task", todo.PriorityB, false)
	})

	t.Run("parses single project tag", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Deploy new feature +WebApp")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		projects := todos[0].Projects()
		is.Equal(len(projects), 1) // expected 1 project
		is.Equal(projects[0], "WebApp")
	})

	t.Run("parses single context tag", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(B) Call client @phone")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 1) // expected 1 context
		is.Equal(contexts[0], "phone")
	})

	t.Run("parses multiple projects and contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Write report +Work +Q1Goals @office @computer")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		projects := todos[0].Projects()
		is.Equal(len(projects), 2) // expected 2 projects
		is.Equal(projects[0], "Work")
		is.Equal(projects[1], "Q1Goals")

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 2) // expected 2 contexts
		is.Equal(contexts[0], "office")
		is.Equal(contexts[1], "computer")
	})

	t.Run("parses tags anywhere in description", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Review +OpenSource code for @github issues")

		todos, err := parser.Parse(input)

		is.NoErr(err)

		projects := todos[0].Projects()
		is.Equal(len(projects), 1) // expected 1 project
		is.Equal(projects[0], "OpenSource")

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 1) // expected 1 context
		is.Equal(contexts[0], "github")
	})

	t.Run("todos without tags have empty slices", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Simple task without tags")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos[0].Projects()), 0) // expected no projects
		is.Equal(len(todos[0].Contexts()), 0) // expected no contexts
	})

	// Boundary tests for tag combinations (catches CONDITIONALS_BOUNDARY mutations)
	t.Run("parses todo with only projects and no contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task with +Project1 +Project2 but no contexts")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(len(todos[0].Projects()), 2) // expected 2 projects
		is.Equal(len(todos[0].Contexts()), 0) // expected 0 contexts
		is.Equal(todos[0].Projects()[0], "Project1")
		is.Equal(todos[0].Projects()[1], "Project2")
	})

	t.Run("parses todo with only contexts and no projects", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task @home @computer but no projects")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(len(todos[0].Projects()), 0) // expected 0 projects
		is.Equal(len(todos[0].Contexts()), 2) // expected 2 contexts
		is.Equal(todos[0].Contexts()[0], "home")
		is.Equal(todos[0].Contexts()[1], "computer")
	})

	t.Run("parses completed todo with only projects", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 (A) Completed +Project")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		is.Equal(len(todos[0].Projects()), 1) // expected 1 project
		is.Equal(len(todos[0].Contexts()), 0) // expected 0 contexts
	})

	t.Run("parses completed todo with only contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 (A) Completed @office")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		is.Equal(len(todos[0].Projects()), 0) // expected 0 projects
		is.Equal(len(todos[0].Contexts()), 1) // expected 1 context
	})

	// Creation date tests (catches CONDITIONALS_NEGATION mutations)
	t.Run("parses todo with creation date", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-01-15 Task with creation date")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task with creation date")
		is.Equal(todos[0].Priority(), todo.PriorityA)
		is.Equal(todos[0].IsCompleted(), false)
	})

	t.Run("parses completed todo with both completion and creation dates", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 2026-01-15 (A) Task with both dates")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task with both dates")
		is.Equal(todos[0].Priority(), todo.PriorityA)
		is.True(todos[0].IsCompleted())
	})

	t.Run("parses todo without creation date", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(B) Task without creation date")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task without creation date")
		is.Equal(todos[0].Priority(), todo.PriorityB)
	})

	// Date parsing error handling (catches error path mutations)
	t.Run("removes malformed completion date from description but doesn't store it", func(t *testing.T) {
		is := is.New(t)
		// Dates matching pattern are removed even if invalid
		input := strings.NewReader("x 2026-99-99 (A) Task with invalid date")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		// Malformed date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid date")
	})

	t.Run("removes malformed creation date from description but doesn't store it", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-99-99 Task with invalid creation date")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// Malformed date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid creation date")
	})

	t.Run("removes malformed creation date on completed todo", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 2026-99-99 (A) Task with invalid creation date")

		todos, err := parser.Parse(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		// Malformed creation date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid creation date")
	})
}

func assertTodo(is *is.I, got todo.Todo, wantDesc string, wantPriority todo.Priority, wantCompleted bool) {
	is.Helper()

	is.Equal(got.Description(), wantDesc) // description
	is.Equal(got.Priority(), wantPriority) // priority
	is.Equal(got.IsCompleted(), wantCompleted) // completed
}
