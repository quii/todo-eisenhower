package usecases

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestChangePriority_ManagesPrioritisedDate(t *testing.T) {
	is := is.New(t)

	t.Run("moving to Priority A adds prioritised date", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		taskB := todo.NewFull("Task B", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

		err := repo.SaveAll([]todo.Todo{taskB})
		is.NoErr(err)

		m, err := LoadMatrix(repo)
		is.NoErr(err)

		// Move from Priority B to Priority A
		updatedMatrix, err := ChangePriority(repo, m, matrix.ScheduleQuadrant, 0, todo.PriorityA)
		is.NoErr(err)

		// Verify it moved to Do First
		doFirst := updatedMatrix.DoFirst()
		is.Equal(len(doFirst), 1)
		is.Equal(doFirst[0].Description(), "Task B")
		is.Equal(doFirst[0].Priority(), todo.PriorityA)

		// Verify prioritised date was added
		is.True(doFirst[0].PrioritisedDate() != nil)
	})

	t.Run("moving from Priority A removes prioritised date", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		taskA := todo.NewFull("Task A", todo.PriorityA, false, nil, &creationDate, nil, &prioritisedDate, nil, nil)

		err := repo.SaveAll([]todo.Todo{taskA})
		is.NoErr(err)

		m, err := LoadMatrix(repo)
		is.NoErr(err)

		// Move from Priority A to Priority B
		updatedMatrix, err := ChangePriority(repo, m, matrix.DoFirstQuadrant, 0, todo.PriorityB)
		is.NoErr(err)

		// Verify it moved to Schedule
		schedule := updatedMatrix.Schedule()
		is.Equal(len(schedule), 1)
		is.Equal(schedule[0].Description(), "Task A")
		is.Equal(schedule[0].Priority(), todo.PriorityB)

		// Verify prioritised date was removed
		is.True(schedule[0].PrioritisedDate() == nil)
	})

	t.Run("moving from B to C does not affect prioritised date", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		taskB := todo.NewFull("Task B", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

		err := repo.SaveAll([]todo.Todo{taskB})
		is.NoErr(err)

		m, err := LoadMatrix(repo)
		is.NoErr(err)

		// Move from Priority B to Priority C
		updatedMatrix, err := ChangePriority(repo, m, matrix.ScheduleQuadrant, 0, todo.PriorityC)
		is.NoErr(err)

		// Verify it moved to Delegate
		delegate := updatedMatrix.Delegate()
		is.Equal(len(delegate), 1)
		is.Equal(delegate[0].Description(), "Task B")
		is.Equal(delegate[0].Priority(), todo.PriorityC)

		// Verify no prioritised date (was nil, stays nil)
		is.True(delegate[0].PrioritisedDate() == nil)
	})

	t.Run("preserves other fields when changing priority", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC)
		taskB := todo.NewFull("Task B", todo.PriorityB, false, nil, &creationDate, &dueDate, nil, []string{"project"}, []string{"context"})

		err := repo.SaveAll([]todo.Todo{taskB})
		is.NoErr(err)

		m, err := LoadMatrix(repo)
		is.NoErr(err)

		// Move from Priority B to Priority A
		updatedMatrix, err := ChangePriority(repo, m, matrix.ScheduleQuadrant, 0, todo.PriorityA)
		is.NoErr(err)

		doFirst := updatedMatrix.DoFirst()
		is.Equal(len(doFirst), 1)

		// Verify all other fields preserved
		is.Equal(doFirst[0].Description(), "Task B")
		is.True(doFirst[0].CreationDate() != nil)
		is.True(doFirst[0].CreationDate().Equal(creationDate))
		is.True(doFirst[0].DueDate() != nil)
		is.True(doFirst[0].DueDate().Equal(dueDate))
		is.Equal(len(doFirst[0].Projects()), 1)
		is.Equal(doFirst[0].Projects()[0], "project")
		is.Equal(len(doFirst[0].Contexts()), 1)
		is.Equal(doFirst[0].Contexts()[0], "context")

		// And has new prioritised date
		is.True(doFirst[0].PrioritisedDate() != nil)
	})
}
