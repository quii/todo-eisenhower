package acceptance_test

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
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
		is := is.New(t)
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
		is.NoErr(err)

		// And the "DO FIRST" quadrant contains "Fix critical bug"
		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "Fix critical bug")

		// And the "SCHEDULE" quadrant contains "Plan quarterly goals"
		schedule := m.Schedule()
		is.Equal(len(schedule), 1) // expected 1 todo in SCHEDULE
		is.Equal(schedule[0].Description(), "Plan quarterly goals")

		// And the "DELEGATE" quadrant contains "Reply to emails"
		delegate := m.Delegate()
		is.Equal(len(delegate), 1) // expected 1 todo in DELEGATE
		is.Equal(delegate[0].Description(), "Reply to emails")

		// And the "ELIMINATE" quadrant contains "Clean workspace"
		eliminate := m.Eliminate()
		is.Equal(len(eliminate), 1) // expected 1 todo in ELIMINATE
		is.Equal(eliminate[0].Description(), "Clean workspace")
	})

	t.Run("Scenario: Handle missing file gracefully", func(t *testing.T) {
		is := is.New(t)
		// Given no file exists at "~/todo.txt"
		source := todoSource{
			err: io.ErrUnexpectedEOF, // simulating file not found
		}

		// When I run "eisenhower"
		_, err := usecases.LoadMatrix(source)

		// Then the application displays an error message
		// And exits gracefully without crashing
		is.True(err != nil) // expected error when file doesn't exist
	})

	t.Run("Scenario: Handle empty file", func(t *testing.T) {
		is := is.New(t)
		// Given an empty file exists at "~/todo.txt"
		source := todoSource{
			data: "",
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		// Then the matrix displays with all quadrants empty
		is.NoErr(err)

		is.Equal(len(m.DoFirst()), 0)   // expected empty DO FIRST quadrant
		is.Equal(len(m.Schedule()), 0)  // expected empty SCHEDULE quadrant
		is.Equal(len(m.Delegate()), 0)  // expected empty DELEGATE quadrant
		is.Equal(len(m.Eliminate()), 0) // expected empty ELIMINATE quadrant
	})

	t.Run("Scenario: Parse completed todos", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file containing:
		source := todoSource{
			data: `(A) Active task
x (A) 2026-01-11 Completed task
(B) Another active task
`,
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		is.NoErr(err)

		// Then the "DO FIRST" quadrant shows both todos
		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 2) // expected 2 todos in DO FIRST

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

		is.True(hasActiveTodo)      // expected to find active 'Active task'
		is.True(hasCompletedTodo)   // expected to find completed 'Completed task'
	})

	t.Run("Scenario: Handle todos without priority", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file containing:
		source := todoSource{
			data: `(A) High priority task
No priority task
`,
		}

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(source)

		is.NoErr(err)

		// Then the "DO FIRST" quadrant contains "High priority task"
		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "High priority task")

		// And the "ELIMINATE" quadrant contains "No priority task"
		eliminate := m.Eliminate()
		is.Equal(len(eliminate), 1) // expected 1 todo in ELIMINATE
		is.Equal(eliminate[0].Description(), "No priority task")
	})
}
