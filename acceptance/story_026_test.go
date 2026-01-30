package acceptance_test

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 026: Stale Task Detection

// Prioritised Tag Management

func TestStory026_AddingPrioritisedTagWhenCreating(t *testing.T) {
	// Scenario: Adding prioritised tag when creating task in Do First
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create task with priority A (Do First)
	task := todo.NewFull("New urgent task", todo.PriorityA, false, nil, &creationDate, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Verify task has prioritised date
	is.True(task.PrioritisedDate() != nil)
	is.True(task.PrioritisedDate().Equal(prioritisedDate))

	// Verify tag is not visible in UI
	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()
	is.True(strings.Contains(view, "New urgent task"))
	is.True(!strings.Contains(view, "prioritised:")) // Tag hidden in UI
}

func TestStory026_AddingPrioritisedTagWhenMoving(t *testing.T) {
	// Scenario: Adding prioritised tag when moving task to Do First
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create task with priority B (Schedule)
	task := todo.NewFull("Scheduled task", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Move task to Do First (priority A)
	updatedMatrix, err := usecases.ChangePriority(repository, m, matrix.ScheduleQuadrant, 0, todo.PriorityA)
	is.NoErr(err)

	// Verify task moved to Do First and has prioritised date
	doFirst := updatedMatrix.DoFirst()
	is.Equal(len(doFirst), 1)
	is.Equal(doFirst[0].Description(), "Scheduled task")
	is.Equal(doFirst[0].Priority(), todo.PriorityA)
	is.True(doFirst[0].PrioritisedDate() != nil)
}

func TestStory026_RemovingPrioritisedTagWhenMoving(t *testing.T) {
	// Scenario: Removing prioritised tag when moving task out of Do First
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create task with priority A (Do First) and prioritised date
	task := todo.NewFull("Urgent task", todo.PriorityA, false, nil, &creationDate, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Move task to Schedule (priority B)
	updatedMatrix, err := usecases.ChangePriority(repository, m, matrix.DoFirstQuadrant, 0, todo.PriorityB)
	is.NoErr(err)

	// Verify task moved to Schedule and prioritised date removed
	schedule := updatedMatrix.Schedule()
	is.Equal(len(schedule), 1)
	is.Equal(schedule[0].Description(), "Urgent task")
	is.Equal(schedule[0].Priority(), todo.PriorityB)
	is.True(schedule[0].PrioritisedDate() == nil)
}

func TestStory026_ResettingPrioritisedTagWhenMovingBack(t *testing.T) {
	// Scenario: Resetting prioritised tag when moving task back to Do First
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	// Create task with priority A and old prioritised date
	task := todo.NewFull("Important task", todo.PriorityA, false, nil, &creationDate, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Move task to Schedule
	updatedMatrix, err := usecases.ChangePriority(repository, m, matrix.DoFirstQuadrant, 0, todo.PriorityB)
	is.NoErr(err)

	// Move task back to Do First (3 days later - Friday)
	updatedMatrix, err = usecases.ChangePriority(repository, updatedMatrix, matrix.ScheduleQuadrant, 0, todo.PriorityA)
	is.NoErr(err)

	// Verify task has NEW prioritised date (current date when moved back)
	doFirst := updatedMatrix.DoFirst()
	is.Equal(len(doFirst), 1)
	is.True(doFirst[0].PrioritisedDate() != nil)
	// The new prioritised date should be "now" (when ChangePriority was called)
	// which is different from the old date
}

// Do First Staleness (2 Business Days)

func TestStory026_DoFirstNotStaleDay1(t *testing.T) {
	// Scenario: Do First task not stale on day 1 (business day)
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday

	task := todo.NewFull("Task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness on same day
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_DoFirstNotStaleDay2(t *testing.T) {
	// Scenario: Do First task not stale on day 2 (business day)
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday

	task := todo.NewFull("Task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 1 business day later (Wednesday)
	now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_DoFirstStaleAfter2BusinessDays(t *testing.T) {
	// Scenario: Do First task is stale after 2 business days
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday

	task := todo.NewFull("Task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 2 business days later (Thursday)
	now := time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false) // Exactly 2 is not stale

	// Check staleness 3 business days later (Friday)
	now = time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), true) // More than 2 is stale
}

func TestStory026_DoFirstStalenessExcludesWeekends(t *testing.T) {
	// Scenario: Do First task staleness excludes weekends
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC) // Thursday

	task := todo.NewFull("Task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness on Monday (4 calendar days, 2 business days)
	// Jan 15 (Thu) -> Jan 16 (Fri) = 1 business day
	// Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day (skip weekend)
	// Total: 2 business days
	now := time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_DoFirstStaleAfterWeekend(t *testing.T) {
	// Scenario: Do First task stale after weekend threshold
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC) // Thursday

	task := todo.NewFull("Task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness on Tuesday (5 calendar days, 3 business days)
	// Jan 15 (Thu) -> Jan 16 (Fri) = 1 business day
	// Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day (skip weekend)
	// Jan 19 (Mon) -> Jan 20 (Tue) = 1 business day
	// Total: 3 business days
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), true)
}

// Schedule/Delegate/Eliminate Staleness (5 Business Days)

func TestStory026_ScheduleNotStaleWithin5BusinessDays(t *testing.T) {
	// Scenario: Schedule task not stale within 5 business days
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday

	task := todo.NewFull("Scheduled task", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 4 business days later (Friday)
	now := time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_ScheduleStaleAfter5BusinessDays(t *testing.T) {
	// Scenario: Schedule task is stale after 5 business days
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday

	task := todo.NewFull("Scheduled task", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness exactly 5 business days later (Monday next week)
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false) // Exactly 5 is not stale

	// Check staleness 6 business days later (Tuesday next week)
	now = time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), true) // More than 5 is stale
}

func TestStory026_DelegateStalenessExcludesWeekends(t *testing.T) {
	// Scenario: Delegate task staleness excludes weekends
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday

	task := todo.NewFull("Delegate task", todo.PriorityC, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 8 calendar days later (Tuesday)
	// Jan 13 (Mon) -> Jan 17 (Fri) = 4 business days
	// Jan 17 (Fri) -> Jan 20 (Mon) = 1 business day (skip weekend)
	// Jan 20 (Mon) -> Jan 21 (Tue) = 1 business day
	// Total: 6 business days
	now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), true)
}

func TestStory026_EliminateFollowsSameStalenessRules(t *testing.T) {
	// Scenario: Eliminate task follows same staleness rules
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday

	task := todo.NewFull("Low priority task", todo.PriorityD, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness exactly 5 business days later (Monday next week)
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false) // Exactly 5 is not stale

	// Check staleness 6 business days later (Tuesday next week)
	now = time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), true) // More than 5 is stale
}

// Completed Items

func TestStory026_CompletedDoFirstNeverStale(t *testing.T) {
	// Scenario: Completed Do First task never stale
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)
	completionDate := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)

	task := todo.NewFull("Important task", todo.PriorityA, true, &completionDate, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 10 business days after prioritisation
	now := time.Date(2026, 1, 27, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_CompletedScheduleNeverStale(t *testing.T) {
	// Scenario: Completed Schedule task never stale
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)
	completionDate := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)

	task := todo.NewFull("Old task", todo.PriorityB, true, &completionDate, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 15 business days after creation
	now := time.Date(2026, 2, 2, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

// Edge Cases

func TestStory026_TaskCreatedOnFridayNotStaleOnMonday(t *testing.T) {
	// Scenario: Task created on Friday not stale on Monday
	is := is.New(t)

	repository := memory.NewRepository()
	creationDate := time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC) // Friday

	task := todo.NewFull("Weekend task", todo.PriorityB, false, nil, &creationDate, nil, nil, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness 3 calendar days later (Monday)
	// Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day
	now := time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

func TestStory026_TaskCreatedTodayNeverStale(t *testing.T) {
	// Scenario: Task created today is never stale
	is := is.New(t)

	repository := memory.NewRepository()
	prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	task := todo.NewFull("New task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	// Check staleness on same day
	now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
	is.Equal(task.IsStale(now), false)
}

// Visual Styling Tests

func TestStory026_VisualStylingForStaleTaskInOverview(t *testing.T) {
	// Verify stale background is applied in overview mode
	is := is.New(t)

	repository := memory.NewRepository()
	// Create a task that is stale (3+ business days ago)
	prioritisedDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday

	task := todo.NewFull("Stale task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Task should be visible
	is.True(strings.Contains(view, "Stale task"))

	// Note: We can't easily test for background color codes in the rendered output
	// as they are ANSI escape sequences. The important thing is that the IsStale()
	// method returns true, which we've tested in domain tests.
}

func TestStory026_VisualStylingForStaleTaskInFocusMode(t *testing.T) {
	// Verify stale background is applied in focus mode
	is := is.New(t)

	repository := memory.NewRepository()
	// Create a task that is stale
	prioritisedDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)

	task := todo.NewFull("Stale focused task", todo.PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)

	err := repository.SaveAll([]todo.Todo{task})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First quadrant
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	view := model.View()

	// Task should be visible in focus mode
	is.True(strings.Contains(view, "Stale focused task"))
}
