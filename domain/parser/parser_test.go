package parser_test

import (
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/domain/parser"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestParse(t *testing.T) {
	t.Run("parses single todo with priority A", func(t *testing.T) {
		input := strings.NewReader("(A) Fix critical bug")

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("expected 1 todo, got %d", len(todos))
		}
		assertTodo(t, todos[0], "Fix critical bug", todo.PriorityA, false)
	})

	t.Run("parses todo without priority", func(t *testing.T) {
		input := strings.NewReader("No priority task")

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("expected 1 todo, got %d", len(todos))
		}
		assertTodo(t, todos[0], "No priority task", todo.PriorityNone, false)
	})

	t.Run("parses completed todo", func(t *testing.T) {
		input := strings.NewReader("x (A) 2026-01-11 Completed task")

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 1 {
			t.Fatalf("expected 1 todo, got %d", len(todos))
		}
		assertTodo(t, todos[0], "Completed task", todo.PriorityA, true)
	})

	t.Run("parses multiple todos with different priorities", func(t *testing.T) {
		input := strings.NewReader(`(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace`)

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 4 {
			t.Fatalf("expected 4 todos, got %d", len(todos))
		}

		assertTodo(t, todos[0], "Fix critical bug", todo.PriorityA, false)
		assertTodo(t, todos[1], "Plan quarterly goals", todo.PriorityB, false)
		assertTodo(t, todos[2], "Reply to emails", todo.PriorityC, false)
		assertTodo(t, todos[3], "Clean workspace", todo.PriorityD, false)
	})

	t.Run("parses empty input", func(t *testing.T) {
		input := strings.NewReader("")

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 0 {
			t.Fatalf("expected 0 todos, got %d", len(todos))
		}
	})

	t.Run("ignores empty lines", func(t *testing.T) {
		input := strings.NewReader(`(A) First task

(B) Second task`)

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 2 {
			t.Fatalf("expected 2 todos, got %d", len(todos))
		}
	})

	t.Run("parses mixed completed and active todos", func(t *testing.T) {
		input := strings.NewReader(`(A) Active task
x (A) 2026-01-11 Completed task
(B) Another active task`)

		todos, err := parser.Parse(input)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(todos) != 3 {
			t.Fatalf("expected 3 todos, got %d", len(todos))
		}

		assertTodo(t, todos[0], "Active task", todo.PriorityA, false)
		assertTodo(t, todos[1], "Completed task", todo.PriorityA, true)
		assertTodo(t, todos[2], "Another active task", todo.PriorityB, false)
	})
}

func assertTodo(t *testing.T, got todo.Todo, wantDesc string, wantPriority todo.Priority, wantCompleted bool) {
	t.Helper()

	if got.Description() != wantDesc {
		t.Errorf("description: got %q, want %q", got.Description(), wantDesc)
	}
	if got.Priority() != wantPriority {
		t.Errorf("priority: got %v, want %v", got.Priority(), wantPriority)
	}
	if got.IsCompleted() != wantCompleted {
		t.Errorf("completed: got %v, want %v", got.IsCompleted(), wantCompleted)
	}
}
