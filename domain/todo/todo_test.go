package todo_test

import (
	"testing"

	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestTodo(t *testing.T) {
	t.Run("create todo with description and priority A", func(t *testing.T) {
		description := "Fix critical production bug"
		priority := todo.PriorityA

		item := todo.New(description, priority)

		if item.Description() != description {
			t.Errorf("got description %q, want %q", item.Description(), description)
		}

		if item.Priority() != priority {
			t.Errorf("got priority %v, want %v", item.Priority(), priority)
		}
	})

	t.Run("create todo with no priority", func(t *testing.T) {
		description := "No priority task"

		item := todo.New(description, todo.PriorityNone)

		if item.Priority() != todo.PriorityNone {
			t.Errorf("got priority %v, want PriorityNone", item.Priority())
		}
	})

	t.Run("newly created todo is not completed", func(t *testing.T) {
		item := todo.New("Some task", todo.PriorityA)

		if item.IsCompleted() {
			t.Error("expected new todo to not be completed")
		}
	})
}
