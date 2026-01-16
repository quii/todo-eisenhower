package todo_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestNewCompleted(t *testing.T) {
	t.Run("creates completed todo with description and priority", func(t *testing.T) {
		is := is.New(t)

		description := "Deploy hotfix"
		priority := todo.PriorityA

		item := todo.NewCompleted(description, priority, nil)

		is.Equal(item.Description(), description)
		is.Equal(item.Priority(), priority)
		is.True(item.IsCompleted()) // expected todo to be completed
		is.True(item.CompletionDate() == nil) // expected no completion date when passing nil
	})
}
