package todo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestTodo(t *testing.T) {
	t.Run("create todo with description and priority A", func(t *testing.T) {
		is := is.New(t)
		description := "Fix critical production bug"
		priority := todo.PriorityA

		item := todo.New(description, priority)

		is.Equal(item.Description(), description)
		is.Equal(item.Priority(), priority)
	})

	t.Run("create todo with no priority", func(t *testing.T) {
		is := is.New(t)
		description := "No priority task"

		item := todo.New(description, todo.PriorityNone)

		is.Equal(item.Priority(), todo.PriorityNone)
	})

	t.Run("newly created todo is not completed", func(t *testing.T) {
		is := is.New(t)
		item := todo.New("Some task", todo.PriorityA)

		is.True(!item.IsCompleted()) // expected new todo to not be completed
	})

	t.Run("create todo with projects", func(t *testing.T) {
		is := is.New(t)
		description := "Deploy feature"
		projects := []string{"WebApp", "Q1Goals"}

		item := todo.NewWithTags(description, todo.PriorityA, projects, nil)

		is.Equal(len(item.Projects()), 2) // expected 2 projects
		is.Equal(item.Projects()[0], "WebApp")
	})

	t.Run("create todo with contexts", func(t *testing.T) {
		is := is.New(t)
		description := "Call client"
		contexts := []string{"phone", "morning"}

		item := todo.NewWithTags(description, todo.PriorityA, nil, contexts)

		is.Equal(len(item.Contexts()), 2) // expected 2 contexts
		is.Equal(item.Contexts()[0], "phone")
	})

	t.Run("todo without tags has empty slices", func(t *testing.T) {
		is := is.New(t)
		item := todo.New("Simple task", todo.PriorityA)

		is.Equal(len(item.Projects()), 0) // expected no projects
		is.Equal(len(item.Contexts()), 0) // expected no contexts
	})
}
