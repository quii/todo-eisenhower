package acceptance_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory003_AcceptCustomFilePath(t *testing.T) {
	t.Run("Scenario: Display file path in header", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file at "/Users/chris/projects/todo.txt"
		source := todoSource{
			data: `(A) Test task
`,
		}
		filePath := "/Users/chris/projects/todo.txt"

		// When I run "eisenhower /Users/chris/projects/todo.txt"
		m, err := usecases.LoadMatrix(source)
		is.NoErr(err)

		// Then the header displays "File: /Users/chris/projects/todo.txt"
		output := ui.RenderMatrix(m, filePath, 0, 0)

		is.True(strings.Contains(stripANSI(output), "File:")) // expected header to contain 'File:'
		is.True(strings.Contains(stripANSI(output), filePath)) // expected header to contain file path

		// And the matrix is displayed below the header
		is.True(strings.Contains(stripANSI(output), "Do First")) // expected matrix to be displayed
	})

	t.Run("Scenario: Load from custom file path", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file at a custom path
		source := todoSource{
			data: `(A) Custom path task
(B) Another task
`,
		}

		// When I run "eisenhower /custom/path/todo.txt"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from that file
		is.NoErr(err)

		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "Custom path task")
	})

	t.Run("Scenario: Use default path when no argument provided", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file at "~/todo.txt"
		source := todoSource{
			data: `(A) Default path task
`,
		}

		// When I run "eisenhower" without arguments
		// (CLI parsing tested in main, we verify use case still works)
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from "~/todo.txt"
		is.NoErr(err)

		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "Default path task")
	})

	t.Run("Scenario: Handle relative paths", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file at "./todo.txt"
		// (File path expansion tested in main, we verify use case works)
		source := todoSource{
			data: `(C) Relative path task
`,
		}

		// When I run "eisenhower ./todo.txt"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays todos from the current directory's todo.txt
		is.NoErr(err)

		delegate := m.Delegate()
		is.Equal(len(delegate), 1) // expected 1 todo in DELEGATE
		is.Equal(delegate[0].Description(), "Relative path task")
	})

	t.Run("Scenario: Handle non-existent custom path", func(t *testing.T) {
		is := is.New(t)
		// Given no file exists at a custom path
		source := todoSource{
			err: errors.New("file not found"),
		}

		// When I run "eisenhower /path/to/missing.txt"
		_, err := usecases.LoadMatrix(source)

		// Then the application displays an error message
		is.True(err != nil) // expected error for non-existent file
	})
}
