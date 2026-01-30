package usecases_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestArchiveCompletedInQuadrant(t *testing.T) {
	t.Run("archives all completed todos in a quadrant", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		active1 := todo.New("Active 1", todo.PriorityA)
		active2 := todo.New("Active 2", todo.PriorityA)
		completed1 := todo.New("Completed 1", todo.PriorityA).ToggleCompletion(time.Now())
		completed2 := todo.New("Completed 2", todo.PriorityA).ToggleCompletion(time.Now())

		err := repo.SaveAll([]todo.Todo{active1, active2, completed1, completed2})
		is.NoErr(err)

		m, err := usecases.LoadMatrix(repo)
		is.NoErr(err)

		updatedMatrix, err := usecases.ArchiveCompletedInQuadrant(repo, m, matrix.DoFirstQuadrant)
		is.NoErr(err)

		// Matrix should only have active todos
		is.Equal(len(updatedMatrix.DoFirst()), 2)

		// Archive should contain both completed todos
		archiveContents := repo.ArchiveString()
		is.True(strings.Contains(archiveContents, "Completed 1"))
		is.True(strings.Contains(archiveContents, "Completed 2"))

		// Verify persistence - reload and check
		reloaded, err := usecases.LoadMatrix(repo)
		is.NoErr(err)
		is.Equal(len(reloaded.DoFirst()), 2)
	})

	t.Run("no-op when no completed todos in quadrant", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		active := todo.New("Active", todo.PriorityA)
		err := repo.SaveAll([]todo.Todo{active})
		is.NoErr(err)

		m, err := usecases.LoadMatrix(repo)
		is.NoErr(err)

		updatedMatrix, err := usecases.ArchiveCompletedInQuadrant(repo, m, matrix.DoFirstQuadrant)
		is.NoErr(err)

		is.Equal(len(updatedMatrix.DoFirst()), 1)
		is.Equal(repo.ArchiveString(), "") // Nothing archived
	})

	t.Run("only archives from specified quadrant", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		doFirstCompleted := todo.New("DoFirst Completed", todo.PriorityA).ToggleCompletion(time.Now())
		scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())

		err := repo.SaveAll([]todo.Todo{doFirstCompleted, scheduleCompleted})
		is.NoErr(err)

		m, err := usecases.LoadMatrix(repo)
		is.NoErr(err)

		updatedMatrix, err := usecases.ArchiveCompletedInQuadrant(repo, m, matrix.DoFirstQuadrant)
		is.NoErr(err)

		// DoFirst should be empty
		is.Equal(len(updatedMatrix.DoFirst()), 0)
		// Schedule should still have its completed todo
		is.Equal(len(updatedMatrix.Schedule()), 1)

		// Only DoFirst todo should be archived
		archiveContents := repo.ArchiveString()
		is.True(strings.Contains(archiveContents, "DoFirst Completed"))
		is.True(!strings.Contains(archiveContents, "Schedule Completed"))
	})
}

func TestArchiveAllCompleted(t *testing.T) {
	t.Run("archives all completed todos across all quadrants", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		doFirstActive := todo.New("DoFirst Active", todo.PriorityA)
		doFirstCompleted := todo.New("DoFirst Completed", todo.PriorityA).ToggleCompletion(time.Now())
		scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())
		delegateActive := todo.New("Delegate Active", todo.PriorityC)
		delegateCompleted := todo.New("Delegate Completed", todo.PriorityC).ToggleCompletion(time.Now())

		err := repo.SaveAll([]todo.Todo{
			doFirstActive, doFirstCompleted,
			scheduleCompleted,
			delegateActive, delegateCompleted,
		})
		is.NoErr(err)

		m, err := usecases.LoadMatrix(repo)
		is.NoErr(err)

		updatedMatrix, err := usecases.ArchiveAllCompleted(repo, m)
		is.NoErr(err)

		// Only active todos should remain
		is.Equal(len(updatedMatrix.DoFirst()), 1)
		is.Equal(len(updatedMatrix.Schedule()), 0)
		is.Equal(len(updatedMatrix.Delegate()), 1)

		// All completed todos should be archived
		archiveContents := repo.ArchiveString()
		is.True(strings.Contains(archiveContents, "DoFirst Completed"))
		is.True(strings.Contains(archiveContents, "Schedule Completed"))
		is.True(strings.Contains(archiveContents, "Delegate Completed"))
	})

	t.Run("no-op when no completed todos anywhere", func(t *testing.T) {
		is := is.New(t)
		repo := memory.NewRepository()

		active1 := todo.New("Active 1", todo.PriorityA)
		active2 := todo.New("Active 2", todo.PriorityB)

		err := repo.SaveAll([]todo.Todo{active1, active2})
		is.NoErr(err)

		m, err := usecases.LoadMatrix(repo)
		is.NoErr(err)

		updatedMatrix, err := usecases.ArchiveAllCompleted(repo, m)
		is.NoErr(err)

		is.Equal(len(updatedMatrix.AllTodos()), 2)
		is.Equal(repo.ArchiveString(), "") // Nothing archived
	})
}
