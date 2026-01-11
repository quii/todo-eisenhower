package acceptance_test

import (
	"io"
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/usecases"
)

type todoSource struct {
	data string
	err  error
}

func (s todoSource) GetTodos() (io.ReadCloser, error) {
	if s.err != nil {
		return nil, s.err
	}
	return io.NopCloser(strings.NewReader(s.data)), nil
}

func TestStory002_LoadFromHardcodedPath(t *testing.T) {
	t.Run("Scenario: Load todos from hardcoded file path", func(t *testing.T) {
		// Given a todo.txt file exists at "~/todo.txt" containing:
		source := todoSource{
			data: `(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace
`,
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from "~/todo.txt"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// And the "DO FIRST" quadrant contains "Fix critical bug"
		doFirst := m.DoFirst()
		if len(doFirst) != 1 {
			t.Fatalf("expected 1 todo in DO FIRST, got %d", len(doFirst))
		}
		if doFirst[0].Description() != "Fix critical bug" {
			t.Errorf("expected 'Fix critical bug', got %q", doFirst[0].Description())
		}

		// And the "SCHEDULE" quadrant contains "Plan quarterly goals"
		schedule := m.Schedule()
		if len(schedule) != 1 {
			t.Fatalf("expected 1 todo in SCHEDULE, got %d", len(schedule))
		}
		if schedule[0].Description() != "Plan quarterly goals" {
			t.Errorf("expected 'Plan quarterly goals', got %q", schedule[0].Description())
		}

		// And the "DELEGATE" quadrant contains "Reply to emails"
		delegate := m.Delegate()
		if len(delegate) != 1 {
			t.Fatalf("expected 1 todo in DELEGATE, got %d", len(delegate))
		}
		if delegate[0].Description() != "Reply to emails" {
			t.Errorf("expected 'Reply to emails', got %q", delegate[0].Description())
		}

		// And the "ELIMINATE" quadrant contains "Clean workspace"
		eliminate := m.Eliminate()
		if len(eliminate) != 1 {
			t.Fatalf("expected 1 todo in ELIMINATE, got %d", len(eliminate))
		}
		if eliminate[0].Description() != "Clean workspace" {
			t.Errorf("expected 'Clean workspace', got %q", eliminate[0].Description())
		}
	})

	t.Run("Scenario: Handle missing file gracefully", func(t *testing.T) {
		// Given no file exists at "~/todo.txt"
		source := todoSource{
			err: io.ErrUnexpectedEOF, // simulating file not found
		}

		// When I run "eisenhower"
		_, err := usecases.LoadMatrix(source)

		// Then the application displays an error message
		// And exits gracefully without crashing
		if err == nil {
			t.Error("expected error when file doesn't exist, got nil")
		}
	})

	t.Run("Scenario: Handle empty file", func(t *testing.T) {
		// Given an empty file exists at "~/todo.txt"
		source := todoSource{
			data: "",
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays with all quadrants empty
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(m.DoFirst()) != 0 {
			t.Errorf("expected empty DO FIRST quadrant, got %d todos", len(m.DoFirst()))
		}
		if len(m.Schedule()) != 0 {
			t.Errorf("expected empty SCHEDULE quadrant, got %d todos", len(m.Schedule()))
		}
		if len(m.Delegate()) != 0 {
			t.Errorf("expected empty DELEGATE quadrant, got %d todos", len(m.Delegate()))
		}
		if len(m.Eliminate()) != 0 {
			t.Errorf("expected empty ELIMINATE quadrant, got %d todos", len(m.Eliminate()))
		}
	})

	t.Run("Scenario: Parse completed todos", func(t *testing.T) {
		// Given a todo.txt file containing:
		source := todoSource{
			data: `(A) Active task
x (A) 2026-01-11 Completed task
(B) Another active task
`,
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Then the "DO FIRST" quadrant shows both todos
		doFirst := m.DoFirst()
		if len(doFirst) != 2 {
			t.Fatalf("expected 2 todos in DO FIRST, got %d", len(doFirst))
		}

		// And completed todos are visually distinct from active todos
		hasActiveTodo := false
		hasCompletedTodo := false

		for _, td := range doFirst {
			if td.Description() == "Active task" && !td.IsCompleted() {
				hasActiveTodo = true
			}
			if td.Description() == "Completed task" && td.IsCompleted() {
				hasCompletedTodo = true
			}
		}

		if !hasActiveTodo {
			t.Error("expected to find active 'Active task'")
		}
		if !hasCompletedTodo {
			t.Error("expected to find completed 'Completed task'")
		}
	})

	t.Run("Scenario: Handle todos without priority", func(t *testing.T) {
		// Given a todo.txt file containing:
		source := todoSource{
			data: `(A) High priority task
No priority task
`,
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Then the "DO FIRST" quadrant contains "High priority task"
		doFirst := m.DoFirst()
		if len(doFirst) != 1 {
			t.Fatalf("expected 1 todo in DO FIRST, got %d", len(doFirst))
		}
		if doFirst[0].Description() != "High priority task" {
			t.Errorf("expected 'High priority task', got %q", doFirst[0].Description())
		}

		// And the "ELIMINATE" quadrant contains "No priority task"
		eliminate := m.Eliminate()
		if len(eliminate) != 1 {
			t.Fatalf("expected 1 todo in ELIMINATE, got %d", len(eliminate))
		}
		if eliminate[0].Description() != "No priority task" {
			t.Errorf("expected 'No priority task', got %q", eliminate[0].Description())
		}
	})
}
