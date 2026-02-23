package acceptance_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 029: Recurring Tasks

func TestStory029_ParseRelativeRecurrence(t *testing.T) {
	// Scenario: Parse a relative recurring task from todo.txt
	is := is.New(t)

	input := "(B) Prep for catchup due:2026-03-09 rec:+2w\n"
	todos, err := todotxt.Unmarshal(strings.NewReader(input))

	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.Equal(todos[0].Description(), "Prep for catchup")
	is.True(todos[0].DueDate() != nil)
	is.True(todos[0].Recurrence() != nil)
	is.True(todos[0].Recurrence().IsRelative())
	is.Equal(todos[0].Recurrence().Interval(), 2)
	is.Equal(todos[0].Recurrence().Unit(), todo.Week)
}

func TestStory029_ParseStrictRecurrence(t *testing.T) {
	// Scenario: Parse a strict recurring task from todo.txt
	is := is.New(t)

	input := "(B) Pay rent due:2026-03-01 rec:1m\n"
	todos, err := todotxt.Unmarshal(strings.NewReader(input))

	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.True(todos[0].Recurrence() != nil)
	is.True(!todos[0].Recurrence().IsRelative())
	is.Equal(todos[0].Recurrence().Interval(), 1)
	is.Equal(todos[0].Recurrence().Unit(), todo.Month)
}

func TestStory029_InvalidRecIgnored(t *testing.T) {
	// Scenario: Invalid rec: tag is ignored
	is := is.New(t)

	input := "(B) Task due:2026-03-01 rec:invalid\n"
	todos, err := todotxt.Unmarshal(strings.NewReader(input))

	is.NoErr(err)
	is.True(todos[0].Recurrence() == nil)
}

func TestStory029_RecWithoutDueIgnored(t *testing.T) {
	// Scenario: rec: without due: is ignored
	is := is.New(t)

	input := "(B) Task rec:+2w\n"
	todos, err := todotxt.Unmarshal(strings.NewReader(input))

	is.NoErr(err)
	is.True(todos[0].Recurrence() == nil)
}

func TestStory029_CompleteRelativeRecurring(t *testing.T) {
	// Scenario: Completing a relative recurring task creates successor from completion date
	is := is.New(t)

	repo := memory.NewRepository()
	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
	rec := todo.NewRecurrence(true, 2, todo.Week)
	task := todo.NewFull("Prep for catchup", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, []string{"Work"}, []string{"office"})

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	// Complete the recurring task
	updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
	is.NoErr(err)

	// Verify: completed original + active successor in Schedule
	todos := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)
	is.Equal(len(todos), 2)

	// Original: completed, no rec tag
	is.True(todos[0].IsCompleted())
	is.True(todos[0].Recurrence() == nil)

	// Successor: active, with rec tag, updated due date
	successor := todos[1]
	is.True(!successor.IsCompleted())
	is.True(successor.Recurrence() != nil)
	is.Equal(successor.Recurrence().String(), "+2w")
	expectedDue := time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC) // now + 2 weeks
	is.Equal(successor.DueDate().Format("2006-01-02"), expectedDue.Format("2006-01-02"))

	// Successor preserves projects and contexts
	is.Equal(successor.Projects(), []string{"Work"})
	is.Equal(successor.Contexts(), []string{"office"})
	is.Equal(successor.Description(), "Prep for catchup")
	is.Equal(successor.Priority(), todo.PriorityB)
}

func TestStory029_CompleteStrictRecurring(t *testing.T) {
	// Scenario: Completing a strict recurring task creates successor from due date
	is := is.New(t)

	repo := memory.NewRepository()
	now := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	rec := todo.NewRecurrence(false, 1, todo.Month)
	task := todo.NewFull("Pay rent", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
	is.NoErr(err)

	successor := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)[1]
	expectedDue := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC) // due + 1 month
	is.Equal(successor.DueDate().Format("2006-01-02"), expectedDue.Format("2006-01-02"))
}

func TestStory029_StrictSkipPastStaleDates(t *testing.T) {
	// Scenario: Strict recurrence advances past stale dates
	is := is.New(t)

	repo := memory.NewRepository()
	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC) // way past
	rec := todo.NewRecurrence(false, 1, todo.Week)
	task := todo.NewFull("Weekly review", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	updatedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
	is.NoErr(err)

	successor := updatedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)[1]
	// Due date should be after today
	is.True(successor.DueDate().After(now))
}

func TestStory029_UncompletingDoesNotRemoveSuccessor(t *testing.T) {
	// Scenario: Uncompleting a task doesn't undo the successor
	is := is.New(t)

	repo := memory.NewRepository()
	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
	rec := todo.NewRecurrence(true, 2, todo.Week)
	task := todo.NewFull("Task", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	// Complete it
	completedMatrix, err := usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
	is.NoErr(err)
	is.Equal(len(completedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)), 2)

	// Uncomplete the first one
	later := now.Add(time.Hour)
	uncompletedMatrix, err := usecases.ToggleCompletion(repo, completedMatrix, matrix.ScheduleQuadrant, 0, later)
	is.NoErr(err)

	// Both todos still exist
	is.Equal(len(uncompletedMatrix.GetTodosForQuadrant(matrix.ScheduleQuadrant)), 2)
}

func TestStory029_RoundTripPersistence(t *testing.T) {
	// Scenario: Recurring task survives save/load round-trip
	is := is.New(t)

	repo := memory.NewRepository()
	dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
	rec := todo.NewRecurrence(true, 2, todo.Week)
	task := todo.NewFull("Catchup", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, []string{"Work"}, nil)

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Load and verify recurrence survives
	reloaded, err := repo.LoadAll()
	is.NoErr(err)
	is.Equal(len(reloaded), 1)
	is.True(reloaded[0].Recurrence() != nil)
	is.Equal(reloaded[0].Recurrence().String(), "+2w")
	is.Equal(reloaded[0].Description(), "Catchup")
	is.Equal(reloaded[0].Projects(), []string{"Work"})
}

func TestStory029_CompletedOriginalSavedWithoutRec(t *testing.T) {
	// Scenario: Completed original is saved without rec: tag
	is := is.New(t)

	repo := memory.NewRepository()
	now := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
	dueDate := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
	rec := todo.NewRecurrence(true, 2, todo.Week)
	task := todo.NewFull("Task", todo.PriorityB, false, nil, nil, &dueDate, nil, &rec, nil, nil)

	err := repo.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	_, err = usecases.ToggleCompletion(repo, m, matrix.ScheduleQuadrant, 0, now)
	is.NoErr(err)

	// Reload and verify
	reloaded, err := repo.LoadAll()
	is.NoErr(err)
	is.Equal(len(reloaded), 2)

	// Find completed one
	for _, td := range reloaded {
		if td.IsCompleted() {
			is.True(td.Recurrence() == nil)
		} else {
			is.True(td.Recurrence() != nil)
		}
	}
}
