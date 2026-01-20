package ui_test

import (
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/ui"
)

func TestFormatDueDate(t *testing.T) {
	t.Run("formats date in current year as month and day", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)

		result := ui.FormatDueDate(&dueDate, now)

		is.Equal(result, "Jan 25")
	})

	t.Run("formats date in different year with full date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)

		result := ui.FormatDueDate(&dueDate, now)

		is.Equal(result, "2025-12-25")
	})

	t.Run("returns empty string for nil date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

		result := ui.FormatDueDate(nil, now)

		is.Equal(result, "")
	})

	t.Run("formats various months correctly", func(t *testing.T) {
		testCases := []struct {
			month    time.Month
			expected string
		}{
			{time.January, "Jan 15"},
			{time.February, "Feb 15"},
			{time.March, "Mar 15"},
			{time.April, "Apr 15"},
			{time.May, "May 15"},
			{time.June, "Jun 15"},
			{time.July, "Jul 15"},
			{time.August, "Aug 15"},
			{time.September, "Sep 15"},
			{time.October, "Oct 15"},
			{time.November, "Nov 15"},
			{time.December, "Dec 15"},
		}

		for _, tc := range testCases {
			is := is.New(t)
			now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
			dueDate := time.Date(2026, tc.month, 15, 0, 0, 0, 0, time.UTC)

			result := ui.FormatDueDate(&dueDate, now)

			is.Equal(result, tc.expected)
		}
	})
}

func TestIsDueOverdue(t *testing.T) {
	t.Run("past date is overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 10, 30, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 15, 14, 0, 0, 0, time.UTC)

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(result) // past date should be overdue
	})

	t.Run("today is not overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 10, 30, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 20, 14, 0, 0, 0, time.UTC) // Same day, different time

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(!result) // today is NOT overdue
	})

	t.Run("future date is not overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 10, 30, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 25, 8, 0, 0, 0, time.UTC)

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(!result) // future date is not overdue
	})

	t.Run("nil date is not overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := ui.IsDueOverdue(nil, now)

		is.True(!result) // nil is not overdue
	})

	t.Run("normalizes dates to start of day", func(t *testing.T) {
		is := is.New(t)
		// Now is late in the day
		now := time.Date(2026, 1, 20, 23, 59, 59, 0, time.UTC)
		// Due date is early in the same day
		dueDate := time.Date(2026, 1, 20, 0, 0, 1, 0, time.UTC)

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(!result) // same day should not be overdue regardless of time
	})

	t.Run("yesterday is overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 19, 23, 59, 59, 0, time.UTC)

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(result) // yesterday is overdue
	})

	t.Run("tomorrow is not overdue", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 23, 59, 59, 0, time.UTC)
		dueDate := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)

		result := ui.IsDueOverdue(&dueDate, now)

		is.True(!result) // tomorrow is not overdue
	})
}

func TestFormatDueDateWithOverdue(t *testing.T) {
	t.Run("returns formatted date and overdue flag for past date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)

		formatted, isOverdue := ui.FormatDueDateWithOverdue(&dueDate, now)

		is.Equal(formatted, "Jan 15")
		is.True(isOverdue)
	})

	t.Run("returns formatted date and not overdue for future date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)

		formatted, isOverdue := ui.FormatDueDateWithOverdue(&dueDate, now)

		is.Equal(formatted, "Jan 25")
		is.True(!isOverdue)
	})

	t.Run("returns formatted date and not overdue for today", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		formatted, isOverdue := ui.FormatDueDateWithOverdue(&dueDate, now)

		is.Equal(formatted, "Jan 20")
		is.True(!isOverdue) // today is NOT overdue
	})

	t.Run("returns empty string and not overdue for nil date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		formatted, isOverdue := ui.FormatDueDateWithOverdue(nil, now)

		is.Equal(formatted, "")
		is.True(!isOverdue)
	})

	t.Run("formats different year with overdue flag", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)

		formatted, isOverdue := ui.FormatDueDateWithOverdue(&dueDate, now)

		is.Equal(formatted, "2025-12-25")
		is.True(isOverdue) // past year is overdue
	})
}
