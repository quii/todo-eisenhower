package usecases_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestToggleCompletion_RecurringTask(t *testing.T) {
	t.Run("completing a relative recurring task creates successor from completion date", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 2, todo.Week)
		task := todo.NewFull("Prep for catchup", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)

		// Should have 2 todos in Schedule: completed original + new successor
		todos := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(todos), 2)

		// First: completed original without rec tag
		is.True(todos[0].IsCompleted())
		is.True(todos[0].Recurrence() == nil)

		// Second: new successor with updated due date
		is.True(!todos[1].IsCompleted())
		is.True(todos[1].Recurrence() != nil)
		is.Equal(todos[1].Recurrence().String(), "+2w")
		expectedDue := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[1].DueDate().Format("2006-01-02"), expectedDue.Format("2006-01-02"))
	})

	t.Run("completing a strict recurring task creates successor from due date", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(false, 1, todo.Month)
		task := todo.NewFull("Pay rent", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)
		todos := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(todos), 2)

		// Successor due date: April 1
		expectedDue := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[1].DueDate().Format("2006-01-02"), expectedDue.Format("2006-01-02"))
	})

	t.Run("strict recurrence skips past stale dates", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC) // way past
		rec := todo.NewRecurrence(false, 1, todo.Week)
		task := todo.NewFull("Weekly review", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)
		todos := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		successor := todos[1]
		// Successor due date should be in the future
		is.True(successor.DueDate().After(now))
	})

	t.Run("successor preserves projects and contexts", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 2, todo.Week)
		task := todo.NewFull("Catchup", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, []string{"Work"}, []string{"office"})

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)
		successor := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)[1]
		is.Equal(successor.Description(), "Catchup")
		is.Equal(successor.Projects(), []string{"Work"})
		is.Equal(successor.Contexts(), []string{"office"})
	})

	t.Run("successor always created in Schedule quadrant", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 1, todo.Day)
		task := todo.NewFull("Daily standup", todo.PriorityA, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.DoFirstQuadrant, 0, now)

		is.NoErr(err)
		// Completed original stays in DoFirst
		doFirst := updatedMatrix.GetTodosForQuadrant(matrix.DoFirstQuadrant)
		is.Equal(len(doFirst), 1)
		is.True(doFirst[0].IsCompleted())

		// Successor created in Schedule (Priority B)
		schedule := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(schedule), 1)
		is.True(!schedule[0].IsCompleted())
		is.Equal(schedule[0].Priority(), todo.PriorityB)
	})

	t.Run("completed original is saved without rec tag", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 2, todo.Week)
		task := todo.NewFull("Task", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		_, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)
		// Verify persisted data doesn't have rec: on completed todo
		saved := repo.String()
		// Load and check
		allTodos, err := repo.LoadAll()
		is.NoErr(err)
		_ = saved

		completedTodo := allTodos[0]
		is.True(completedTodo.IsCompleted())
		is.True(completedTodo.Recurrence() == nil)
	})

	t.Run("uncompleting does not remove successor", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 2, todo.Week)
		task := todo.NewFull("Task", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		// Complete it first
		completedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
		is.NoErr(err)
		todos := completedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(todos), 2)

		// Uncomplete the first one (the completed original)
		later := now.Add(time.Hour)
		uncomletedMatrix, err := usecases.ToggleCompletion(repo, completedMatrix, matrix.ScheduleQuadrant, 0, later)
		is.NoErr(err)
		// Both todos should still exist
		todos = uncomletedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(todos), 2)
	})

	t.Run("non-recurring task is unaffected", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		task := todo.New("Regular task", todo.PriorityB)

		m := matrix.New([]todo.Todo{task})
		updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)

		is.NoErr(err)
		todos := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
		is.Equal(len(todos), 1) // no successor created
		is.True(todos[0].IsCompleted())
	})

	t.Run("persistence verified - successor saved to repo", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()
		now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		rec := todo.NewRecurrence(true, 2, todo.Week)
		task := todo.NewFull("Task", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

		m := matrix.New([]todo.Todo{task})
		_, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
		is.NoErr(err)

		// Reload from repo and verify
		reloaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(reloaded), 2) // original + successor

		// Find the active one (successor)
		var successor todo.Todo
		for _, t := range reloaded {
			if !t.IsCompleted() {
				successor = t
			}
		}
		is.True(successor.Recurrence() != nil)
		expectedDue := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC)
		is.Equal(successor.DueDate().Format("2006-01-02"), expectedDue.Format("2006-01-02"))
	})
}
