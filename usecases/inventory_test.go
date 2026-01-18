package usecases_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestAnalyzeInventory_ActiveCounts(t *testing.T) {
	//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
	is := is.New(t)

	// Create todos in different quadrants
	todoA1 := todo.New("Urgent task", todo.PriorityA)
	todoA2 := todo.New("Another urgent", todo.PriorityA)
	todoB := todo.New("Important task", todo.PriorityB)

	m := matrix.New([]todo.Todo{todoA1, todoA2, todoB})

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	metrics := matrix.NewInventory(m, now)

	// Check active counts
	is.Equal(metrics.DoFirstActive, 2)
	is.Equal(metrics.ScheduleActive, 1)
	is.Equal(metrics.TotalActive, 3)
}

func TestAnalyzeInventory_OldestAge(t *testing.T) {
	//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
	is := is.New(t)

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Create todo with creation date 21 days before "now"
	oldDate := now.AddDate(0, 0, -21)
	oldTodo := todo.NewWithCreationDate("Old task", todo.PriorityA, &oldDate)

	// Recent todo 3 days before "now"
	recentDate := now.AddDate(0, 0, -3)
	recentTodo := todo.NewWithCreationDate("Recent task", todo.PriorityA, &recentDate)

	m := matrix.New([]todo.Todo{oldTodo, recentTodo})

	metrics := matrix.NewInventory(m, now)

	// Oldest should be exactly 21 days (deterministic now)
	is.Equal(metrics.DoFirstOldestDays, 21)
}

func TestAnalyzeInventory_ContextBreakdown(t *testing.T) {
	//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
	is := is.New(t)

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Create todos with contexts 18 and 2 days before "now"
	old := now.AddDate(0, 0, -18)
	recent := now.AddDate(0, 0, -2)

	todo1 := todo.NewWithTagsAndDates("Task 1", todo.PriorityA, &old, nil, []string{"people"})
	todo2 := todo.NewWithTagsAndDates("Task 2", todo.PriorityA, &old, nil, []string{"people"})
	todo3 := todo.NewWithTagsAndDates("Task 3", todo.PriorityB, &recent, nil, []string{"architecture"})

	m := matrix.New([]todo.Todo{todo1, todo2, todo3})

	metrics := matrix.NewInventory(m, now)

	// Should have context breakdown
	is.True(len(metrics.ContextBreakdown) > 0)

	// People context should have 2 items
	peopleMetrics, found := metrics.ContextBreakdown["people"]
	is.True(found)
	is.Equal(peopleMetrics.Count, 2)

	// Average age should be exactly 18 days (deterministic now)
	is.Equal(peopleMetrics.AvgAgeDays, 18)
}

func TestAnalyzeInventory_ProjectBreakdown(t *testing.T) {
	is := is.New(t)

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Create todos with projects 20 and 5 days before "now"
	old := now.AddDate(0, 0, -20)
	recent := now.AddDate(0, 0, -5)

	todo1 := todo.NewWithTagsAndDates("Task 1", todo.PriorityA, &old, []string{"WebApp"}, nil)
	todo2 := todo.NewWithTagsAndDates("Task 2", todo.PriorityA, &old, []string{"WebApp"}, nil)
	todo3 := todo.NewWithTagsAndDates("Task 3", todo.PriorityB, &recent, []string{"Mobile"}, nil)

	m := matrix.New([]todo.Todo{todo1, todo2, todo3})

	metrics := matrix.NewInventory(m, now)

	// Should have project breakdown
	is.True(len(metrics.ProjectBreakdown) > 0)

	// WebApp project should have 2 items
	webappMetrics, found := metrics.ProjectBreakdown["WebApp"]
	is.True(found)
	is.Equal(webappMetrics.Count, 2)

	// Average age should be exactly 20 days (deterministic now)
	is.Equal(webappMetrics.AvgAgeDays, 20)
}

func TestAnalyzeInventory_ExcludesTodosWithoutCreationDates(t *testing.T) {
	is := is.New(t)

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Create todos: one with date, one without
	withDate := now.AddDate(0, 0, -10)
	todoWithDate := todo.NewWithTagsAndDates("Has date", todo.PriorityA, &withDate, []string{"project1"}, []string{"context1"})
	todoNoDate := todo.NewWithTagsAndDates("No date", todo.PriorityA, nil, []string{"project2"}, []string{"context2"})

	m := matrix.New([]todo.Todo{todoWithDate, todoNoDate})

	metrics := matrix.NewInventory(m, now)

	// Should count both as active
	is.Equal(metrics.DoFirstActive, 2)
	is.Equal(metrics.TotalActive, 2)

	// But only the one with date should appear in breakdowns
	is.Equal(len(metrics.ProjectBreakdown), 1)
	is.Equal(len(metrics.ContextBreakdown), 1)

	_, hasProject1 := metrics.ProjectBreakdown["project1"]
	is.True(hasProject1)

	_, hasProject2 := metrics.ProjectBreakdown["project2"]
	is.True(!hasProject2) // Should not include project from todo without date
}

func TestAnalyzeInventory_Throughput(t *testing.T) {
	is := is.New(t)

	// Use fixed "now" for deterministic testing
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// Completed recently (within 7 days of "now")
	recentCompletion := now.AddDate(0, 0, -3)
	completed1 := todo.NewCompletedWithDates("Done 1", todo.PriorityA, &recentCompletion, nil)
	completed2 := todo.NewCompletedWithDates("Done 2", todo.PriorityA, &recentCompletion, nil)

	// Completed long ago (>7 days before "now")
	oldCompletion := now.AddDate(0, 0, -10)
	completedOld := todo.NewCompletedWithDates("Done old", todo.PriorityA, &oldCompletion, nil)

	// Created recently (within 7 days of "now")
	recentCreation := now.AddDate(0, 0, -2)
	newTodo := todo.NewWithCreationDate("New", todo.PriorityA, &recentCreation)

	m := matrix.New([]todo.Todo{completed1, completed2, completedOld, newTodo})

	metrics := matrix.NewInventory(m, now)

	// Should show 2 completed in last 7 days
	is.Equal(metrics.CompletedLast7Days, 2)

	// Should show 1 added in last 7 days
	is.Equal(metrics.AddedLast7Days, 1)
}

// Boundary tests for 7-day threshold (catches CONDITIONALS_BOUNDARY/NEGATION mutations)
func TestAnalyzeInventory_ThroughputBoundaries(t *testing.T) {
	t.Run("exactly 7 days ago should not be counted", func(t *testing.T) {
		is := is.New(t)

		// Use fixed "now" for deterministic testing
		now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		// Exactly 7 days before "now"
		exactly7Days := now.AddDate(0, 0, -7)
		completed := todo.NewCompletedWithDates("Done exactly 7", todo.PriorityA, &exactly7Days, nil)
		created := todo.NewWithCreationDate("Created exactly 7", todo.PriorityA, &exactly7Days)

		m := matrix.New([]todo.Todo{completed, created})
		metrics := matrix.NewInventory(m, now)

		// After(sevenDaysAgo) should exclude exactly 7 days ago
		is.Equal(metrics.CompletedLast7Days, 0) // not counted (too old)
		is.Equal(metrics.AddedLast7Days, 0)     // not counted (too old)
	})

	t.Run("6 days ago should be counted", func(t *testing.T) {
		is := is.New(t)

		// Use fixed "now" for deterministic testing
		now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		// 6 days before "now" (within threshold)
		within7Days := now.AddDate(0, 0, -6)
		completed := todo.NewCompletedWithDates("Done 6 days", todo.PriorityA, &within7Days, nil)
		created := todo.NewWithCreationDate("Created 6 days", todo.PriorityA, &within7Days)

		m := matrix.New([]todo.Todo{completed, created})
		metrics := matrix.NewInventory(m, now)

		// Should be counted (recent enough)
		is.Equal(metrics.CompletedLast7Days, 1)
		is.Equal(metrics.AddedLast7Days, 1)
	})

	t.Run("8 days ago should not be counted", func(t *testing.T) {
		is := is.New(t)

		// Use fixed "now" for deterministic testing
		now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		// 8 days before "now" (too old)
		moreThan7Days := now.AddDate(0, 0, -8)
		completed := todo.NewCompletedWithDates("Done 8 days", todo.PriorityA, &moreThan7Days, nil)
		created := todo.NewWithCreationDate("Created 8 days", todo.PriorityA, &moreThan7Days)

		m := matrix.New([]todo.Todo{completed, created})
		metrics := matrix.NewInventory(m, now)

		// Should not be counted (too old)
		is.Equal(metrics.CompletedLast7Days, 0)
		is.Equal(metrics.AddedLast7Days, 0)
	})
}

// Boundary tests for oldest age comparisons (catches CONDITIONALS_BOUNDARY mutations)
func TestAnalyzeInventory_OldestAgeBoundaries(t *testing.T) {
	t.Run("when two todos have same age, tracks the age correctly", func(t *testing.T) {
		is := is.New(t)

		// Use fixed "now" for deterministic testing
		now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		// Two todos with identical ages (15 days before "now")
		sameAge := now.AddDate(0, 0, -15)
		todo1 := todo.NewWithCreationDate("First", todo.PriorityA, &sameAge)
		todo2 := todo.NewWithCreationDate("Second", todo.PriorityA, &sameAge)

		m := matrix.New([]todo.Todo{todo1, todo2})
		metrics := matrix.NewInventory(m, now)

		// Should track the age (both are equally old) - exactly 15 days
		is.Equal(metrics.DoFirstOldestDays, 15)
	})

	t.Run("tracks oldest when one is older", func(t *testing.T) {
		is := is.New(t)

		// Use fixed "now" for deterministic testing
		now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		// Different ages (20 and 10 days before "now")
		older := now.AddDate(0, 0, -20)
		newer := now.AddDate(0, 0, -10)
		todo1 := todo.NewWithCreationDate("Older", todo.PriorityA, &older)
		todo2 := todo.NewWithCreationDate("Newer", todo.PriorityA, &newer)

		m := matrix.New([]todo.Todo{todo1, todo2})
		metrics := matrix.NewInventory(m, now)

		// Should track the older one (20 days, not 10)
		is.Equal(metrics.DoFirstOldestDays, 20)
	})
}
