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

	t.Run("create todo with projects", func(t *testing.T) {
		description := "Deploy feature"
		projects := []string{"WebApp", "Q1Goals"}

		item := todo.NewWithTags(description, todo.PriorityA, projects, nil)

		if len(item.Projects()) != 2 {
			t.Errorf("expected 2 projects, got %d", len(item.Projects()))
		}
		if item.Projects()[0] != "WebApp" {
			t.Errorf("expected first project to be 'WebApp', got %q", item.Projects()[0])
		}
	})

	t.Run("create todo with contexts", func(t *testing.T) {
		description := "Call client"
		contexts := []string{"phone", "morning"}

		item := todo.NewWithTags(description, todo.PriorityA, nil, contexts)

		if len(item.Contexts()) != 2 {
			t.Errorf("expected 2 contexts, got %d", len(item.Contexts()))
		}
		if item.Contexts()[0] != "phone" {
			t.Errorf("expected first context to be 'phone', got %q", item.Contexts()[0])
		}
	})

	t.Run("todo without tags has empty slices", func(t *testing.T) {
		item := todo.New("Simple task", todo.PriorityA)

		if len(item.Projects()) != 0 {
			t.Errorf("expected no projects, got %d", len(item.Projects()))
		}
		if len(item.Contexts()) != 0 {
			t.Errorf("expected no contexts, got %d", len(item.Contexts()))
		}
	})
}
