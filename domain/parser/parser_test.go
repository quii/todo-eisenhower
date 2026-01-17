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
}

func assertTodo(is *is.I, got todo.Todo, wantDesc string, wantPriority todo.Priority, wantCompleted bool) {
	is.Helper()

	is.Equal(got.Description(), wantDesc) // description
	is.Equal(got.Priority(), wantPriority) // priority
	is.Equal(got.IsCompleted(), wantCompleted) // completed
}
