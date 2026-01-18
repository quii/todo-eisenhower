package todo_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestToggleCompletion(t *testing.T) {
	t.Run("marking incomplete todo as complete sets completion date", func(t *testing.T) {
		is := is.New(t)

		// Given an incomplete todo
		incompleteTodo := todo.New("Write tests", todo.PriorityA)
		is.Equal(incompleteTodo.IsCompleted(), false)
		is.True(incompleteTodo.CompletionDate() == nil)

		// When we toggle completion with a specific time
		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		completedTodo := incompleteTodo.ToggleCompletion(now)

		// Then it should be marked as complete with the provided time
		is.Equal(completedTodo.IsCompleted(), true)
		is.True(completedTodo.CompletionDate() != nil)
		is.Equal(*completedTodo.CompletionDate(), now)
	})

	t.Run("marking complete todo as incomplete clears completion date", func(t *testing.T) {
		is := is.New(t)

		// Given a completed todo
		completionDate := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)
		completedTodo := todo.NewCompletedWithDates("Write tests", todo.PriorityA, &completionDate, nil)
		is.Equal(completedTodo.IsCompleted(), true)
		is.True(completedTodo.CompletionDate() != nil)

		// When we toggle completion
		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		incompleteTodo := completedTodo.ToggleCompletion(now)

		// Then it should be marked as incomplete with no completion date
		is.Equal(incompleteTodo.IsCompleted(), false)
		is.True(incompleteTodo.CompletionDate() == nil)
	})

	t.Run("toggling preserves description", func(t *testing.T) {
		is := is.New(t)

		original := todo.New("Important task", todo.PriorityB)
		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		toggled := original.ToggleCompletion(now)

		is.Equal(toggled.Description(), "Important task")
	})

	t.Run("toggling preserves priority", func(t *testing.T) {
		is := is.New(t)

		original := todo.New("Important task", todo.PriorityB)
		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		toggled := original.ToggleCompletion(now)

		is.Equal(toggled.Priority(), todo.PriorityB)
	})

	t.Run("toggling preserves creation date", func(t *testing.T) {
		is := is.New(t)

		creationDate := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
		original := todo.NewWithCreationDate("Task with date", todo.PriorityA, &creationDate)

		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		toggled := original.ToggleCompletion(now)

		is.True(toggled.CreationDate() != nil)
		is.Equal(*toggled.CreationDate(), creationDate)
	})

	t.Run("toggling preserves contexts", func(t *testing.T) {
		is := is.New(t)

		original := todo.NewWithTagsAndDates("Task", todo.PriorityA, nil, nil, []string{"home", "urgent"})

		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		toggled := original.ToggleCompletion(now)

		contexts := toggled.Contexts()
		is.Equal(len(contexts), 2)
		is.Equal(contexts[0], "home")
		is.Equal(contexts[1], "urgent")
	})

	t.Run("toggling preserves projects", func(t *testing.T) {
		is := is.New(t)

		original := todo.NewWithTagsAndDates("Task", todo.PriorityA, nil, []string{"ProjectX", "ProjectY"}, nil)

		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
		toggled := original.ToggleCompletion(now)

		projects := toggled.Projects()
		is.Equal(len(projects), 2)
		is.Equal(projects[0], "ProjectX")
		is.Equal(projects[1], "ProjectY")
	})

	t.Run("toggling twice returns to original state", func(t *testing.T) {
		is := is.New(t)

		original := todo.New("Task", todo.PriorityA)
		now := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)

		// Toggle to complete
		completed := original.ToggleCompletion(now)
		is.Equal(completed.IsCompleted(), true)

		// Toggle back to incomplete
		later := time.Date(2024, 1, 16, 14, 30, 0, 0, time.UTC)
		incomplete := completed.ToggleCompletion(later)
		is.Equal(incomplete.IsCompleted(), false)

		// Should be back to incomplete state (no completion date)
		is.True(incomplete.CompletionDate() == nil)
		is.Equal(incomplete.Description(), original.Description())
		is.Equal(incomplete.Priority(), original.Priority())
	})

	t.Run("uses deterministic time not system clock", func(t *testing.T) {
		is := is.New(t)

		todo1 := todo.New("Task", todo.PriorityA)

		// Complete at different times
		time1 := time.Date(2020, 1, 1, 10, 0, 0, 0, time.UTC)
		time2 := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

		completed1 := todo1.ToggleCompletion(time1)
		completed2 := todo1.ToggleCompletion(time2)

		// Should have different completion dates
		is.Equal(*completed1.CompletionDate(), time1)
		is.Equal(*completed2.CompletionDate(), time2)
	})
}
