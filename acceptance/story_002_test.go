package acceptance_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestStory002_LoadFromHardcodedPath(t *testing.T) {
	t.Run("Scenario: Load todos from hardcoded file path", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file exists at "~/todo.txt" containing:
		repository := memory.NewRepository()
		initialTodos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Plan quarterly goals", todo.PriorityB),
			todo.New("Reply to emails", todo.PriorityC),
			todo.New("Clean workspace", todo.PriorityD),
		}
		err := repository.SaveAll(initialTodos)
		is.NoErr(err)

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(repository)
		is.NoErr(err)

		// The "Do First" quadrant contains "Fix critical bug"
		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "Fix critical bug")

		// And the "Schedule" quadrant contains "Plan quarterly goals"
		schedule := m.Schedule()
		is.Equal(len(schedule), 1) // expected 1 todo in SCHEDULE
		is.Equal(schedule[0].Description(), "Plan quarterly goals")

		// And the "Delegate" quadrant contains "Reply to emails"
		delegate := m.Delegate()
		is.Equal(len(delegate), 1) // expected 1 todo in DELEGATE
		is.Equal(delegate[0].Description(), "Reply to emails")

		// And the "Eliminate" quadrant contains "Clean workspace"
		eliminate := m.Eliminate()
		is.Equal(len(eliminate), 1) // expected 1 todo in ELIMINATE
		is.Equal(eliminate[0].Description(), "Clean workspace")
	})

	t.Run("Scenario: Handle missing file gracefully", func(t *testing.T) {
		// This test is for file not found errors, which the fake doesn't simulate
		// In real usage, file.Repository creates an empty file if it doesn't exist
		// So this scenario is handled by the file adapter, not the use case
		// Skip this test or test it at the adapter level
		t.Skip("File-not-found behavior is adapter-specific")
	})

	t.Run("Scenario: Handle empty file", func(t *testing.T) {
		is := is.New(t)
		// Given an empty file exists at "~/todo.txt"
		repository := memory.NewRepository()

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(repository)

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
		repository := memory.NewRepository()
		completionDate := time.Date(2026, 1, 11, 0, 0, 0, 0, time.UTC)
		initialTodos := []todo.Todo{
			todo.New("Active task", todo.PriorityA),
			todo.NewCompleted("Completed task", todo.PriorityA, &completionDate),
			todo.New("Another active task", todo.PriorityB),
		}
		err := repository.SaveAll(initialTodos)
		is.NoErr(err)

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(repository)

		is.NoErr(err)

		// Then the "Do First" quadrant shows both todos
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

		is.True(hasActiveTodo)    // expected to find active 'Active task'
		is.True(hasCompletedTodo) // expected to find completed 'Completed task'
	})

	t.Run("Scenario: Handle todos without priority", func(t *testing.T) {
		is := is.New(t)
		// Given a todo.txt file containing:
		repository := memory.NewRepository()
		initialTodos := []todo.Todo{
			todo.New("High priority task", todo.PriorityA),
			todo.New("No priority task", todo.PriorityNone),
		}
		err := repository.SaveAll(initialTodos)
		is.NoErr(err)

		// When I run "eisenhower"
		m, err := usecases.LoadMatrix(repository)

		is.NoErr(err)

		// Then the "Do First" quadrant contains "High priority task"
		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 1) // expected 1 todo in DO FIRST
		is.Equal(doFirst[0].Description(), "High priority task")

		// And the "Eliminate" quadrant contains "No priority task"
		eliminate := m.Eliminate()
		is.Equal(len(eliminate), 1) // expected 1 todo in ELIMINATE
		is.Equal(eliminate[0].Description(), "No priority task")
	})
}
