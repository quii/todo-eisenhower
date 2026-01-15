package todo_test

import (
	"testing"

	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestNewCompleted(t *testing.T) {
	t.Run("creates completed todo with description and priority", func(t *testing.T) {
		description := "Deploy hotfix"
		priority := todo.PriorityA

		item := todo.NewCompleted(description, priority, nil)

		if item.Description() != description {
			t.Errorf("got description %q, want %q", item.Description(), description)
		}

		if item.Priority() != priority {
			t.Errorf("got priority %v, want %v", item.Priority(), priority)
		}

		if !item.IsCompleted() {
			t.Error("expected todo to be completed")
		}

		if item.CompletionDate() != nil {
			t.Error("expected no completion date when passing nil")
		}
	})
}
