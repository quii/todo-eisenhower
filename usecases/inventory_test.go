package usecases_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestAnalyzeInventory_ActiveCounts(t *testing.T) {
	is := is.New(t)

	// Create todos in different quadrants
	todoA1 := todo.New("Urgent task", todo.PriorityA)
	todoA2 := todo.New("Another urgent", todo.PriorityA)
	todoB := todo.New("Important task", todo.PriorityB)
	
	m := matrix.New([]todo.Todo{todoA1, todoA2, todoB})

	metrics := usecases.AnalyzeInventory(m)

	// Check active counts
	is.Equal(metrics.DoFirstActive, 2)
	is.Equal(metrics.ScheduleActive, 1)
	is.Equal(metrics.TotalActive, 3)
}

func TestAnalyzeInventory_OldestAge(t *testing.T) {
	is := is.New(t)

	// Create todo with creation date 21 days ago
	oldDate := time.Now().AddDate(0, 0, -21)
	oldTodo := todo.NewWithCreationDate("Old task", todo.PriorityA, &oldDate)
	
	// Recent todo
	recentDate := time.Now().AddDate(0, 0, -3)
	recentTodo := todo.NewWithCreationDate("Recent task", todo.PriorityA, &recentDate)

	m := matrix.New([]todo.Todo{oldTodo, recentTodo})

	metrics := usecases.AnalyzeInventory(m)

	// Oldest should be ~21 days
	is.True(metrics.DoFirstOldestDays >= 20)
	is.True(metrics.DoFirstOldestDays <= 22)
}

func TestAnalyzeInventory_ContextBreakdown(t *testing.T) {
	is := is.New(t)

	// Create todos with contexts
	old := time.Now().AddDate(0, 0, -18)
	recent := time.Now().AddDate(0, 0, -2)

	todo1 := todo.NewWithTagsAndDates("Task 1", todo.PriorityA, &old, nil, []string{"people"})
	todo2 := todo.NewWithTagsAndDates("Task 2", todo.PriorityA, &old, nil, []string{"people"})
	todo3 := todo.NewWithTagsAndDates("Task 3", todo.PriorityB, &recent, nil, []string{"architecture"})

	m := matrix.New([]todo.Todo{todo1, todo2, todo3})

	metrics := usecases.AnalyzeInventory(m)

	// Should have context breakdown
	is.True(len(metrics.ContextBreakdown) > 0)
	
	// People context should have 2 items
	peopleMetrics, found := metrics.ContextBreakdown["people"]
	is.True(found)
	is.Equal(peopleMetrics.Count, 2)
	
	// Average age should be ~18 days
	is.True(peopleMetrics.AvgAgeDays >= 17)
	is.True(peopleMetrics.AvgAgeDays <= 19)
}

func TestAnalyzeInventory_Throughput(t *testing.T) {
	is := is.New(t)

	// Completed recently (within 7 days)
	recentCompletion := time.Now().AddDate(0, 0, -3)
	completed1 := todo.NewCompletedWithDates("Done 1", todo.PriorityA, &recentCompletion, nil)
	completed2 := todo.NewCompletedWithDates("Done 2", todo.PriorityA, &recentCompletion, nil)

	// Completed long ago (>7 days)
	oldCompletion := time.Now().AddDate(0, 0, -10)
	completedOld := todo.NewCompletedWithDates("Done old", todo.PriorityA, &oldCompletion, nil)

	// Created recently
	recentCreation := time.Now().AddDate(0, 0, -2)
	newTodo := todo.NewWithCreationDate("New", todo.PriorityA, &recentCreation)

	m := matrix.New([]todo.Todo{completed1, completed2, completedOld, newTodo})

	metrics := usecases.AnalyzeInventory(m)

	// Should show 2 completed in last 7 days
	is.Equal(metrics.CompletedLast7Days, 2)
	
	// Should show 1 added in last 7 days
	is.Equal(metrics.AddedLast7Days, 1)
}
