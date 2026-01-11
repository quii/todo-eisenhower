package acceptance_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory003_AcceptCustomFilePath(t *testing.T) {
	t.Run("Scenario: Display file path in header", func(t *testing.T) {
		// Given a todo.txt file at "/Users/chris/projects/todo.txt"
		source := todoSource{
			data: `(A) Test task
`,
		}
		filePath := "/Users/chris/projects/todo.txt"

		// When I run "eisenhower /Users/chris/projects/todo.txt"
		m, err := usecases.LoadMatrix(source)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Then the header displays "File: /Users/chris/projects/todo.txt"
		output := ui.RenderMatrix(m, filePath)

		if !strings.Contains(output, "File:") {
			t.Error("expected header to contain 'File:'")
		}
		if !strings.Contains(output, filePath) {
			t.Errorf("expected header to contain %q", filePath)
		}

		// And the matrix is displayed below the header
		if !strings.Contains(output, "DO FIRST") {
			t.Error("expected matrix to be displayed")
		}
	})

	t.Run("Scenario: Load from custom file path", func(t *testing.T) {
		// Given a todo.txt file at a custom path
		source := todoSource{
			data: `(A) Custom path task
(B) Another task
`,
		}

		// When I run "eisenhower /custom/path/todo.txt"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from that file
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doFirst := m.DoFirst()
		if len(doFirst) != 1 {
			t.Fatalf("expected 1 todo in DO FIRST, got %d", len(doFirst))
		}
		if doFirst[0].Description() != "Custom path task" {
			t.Errorf("expected 'Custom path task', got %q", doFirst[0].Description())
		}
	})

	t.Run("Scenario: Use default path when no argument provided", func(t *testing.T) {
		// Given a todo.txt file at "~/todo.txt"
		source := todoSource{
			data: `(A) Default path task
`,
		}

		// When I run "eisenhower" without arguments
		// (CLI parsing tested in main, we verify use case still works)
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from "~/todo.txt"
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		doFirst := m.DoFirst()
		if len(doFirst) != 1 {
			t.Fatalf("expected 1 todo in DO FIRST, got %d", len(doFirst))
		}
		if doFirst[0].Description() != "Default path task" {
			t.Errorf("expected 'Default path task', got %q", doFirst[0].Description())
		}
	})

	t.Run("Scenario: Handle relative paths", func(t *testing.T) {
		// Given a todo.txt file at "./todo.txt"
		// (File path expansion tested in main, we verify use case works)
		source := todoSource{
			data: `(C) Relative path task
`,
		}

		// When I run "eisenhower ./todo.txt"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from the current directory's todo.txt
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		delegate := m.Delegate()
		if len(delegate) != 1 {
			t.Fatalf("expected 1 todo in DELEGATE, got %d", len(delegate))
		}
		if delegate[0].Description() != "Relative path task" {
			t.Errorf("expected 'Relative path task', got %q", delegate[0].Description())
		}
	})

	t.Run("Scenario: Handle non-existent custom path", func(t *testing.T) {
		// Given no file exists at a custom path
		source := todoSource{
			err: errors.New("file not found"),
		}

		// When I run "eisenhower /path/to/missing.txt"
		_, err := usecases.LoadMatrix(source)

		// Then the application displays an error message
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})
}
