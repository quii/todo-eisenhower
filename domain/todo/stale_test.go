package todo

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestBusinessDaysBetween(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected int
	}{
		{
			name:     "same day returns 0",
			from:     time.Date(2026, 1, 20, 10, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 1, 20, 15, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "Monday to Tuesday is 1 business day",
			from:     time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // Monday
			to:       time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC), // Tuesday
			expected: 1,
		},
		{
			name:     "Tuesday to Thursday is 2 business days",
			from:     time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC), // Tuesday
			to:       time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC), // Thursday
			expected: 2,
		},
		{
			name:     "Friday to Monday is 1 business day (skip weekend)",
			from:     time.Date(2026, 1, 16, 0, 0, 0, 0, time.UTC), // Friday
			to:       time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // Monday
			expected: 1,
		},
		{
			name:     "Thursday to Tuesday is 3 business days (skip weekend)",
			from:     time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC), // Thursday
			to:       time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC), // Tuesday
			expected: 3,
		},
		{
			name:     "Monday to Friday is 4 business days",
			from:     time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // Monday
			to:       time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC), // Friday
			expected: 4,
		},
		{
			name:     "Monday to next Monday is 5 business days",
			from:     time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // Monday
			to:       time.Date(2026, 1, 26, 0, 0, 0, 0, time.UTC), // Next Monday
			expected: 5,
		},
		{
			name:     "Saturday to Sunday is 0 business days",
			from:     time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC), // Saturday
			to:       time.Date(2026, 1, 18, 0, 0, 0, 0, time.UTC), // Sunday
			expected: 0,
		},
		{
			name:     "Saturday to Monday is 1 business day",
			from:     time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC), // Saturday
			to:       time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC), // Monday
			expected: 1,
		},
		{
			name:     "to before from returns 0",
			from:     time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result := businessDaysBetween(tt.from, tt.to)
			is.Equal(result, tt.expected)
		})
	}
}

func TestIsStale_PriorityA(t *testing.T) {
	is := is.New(t)

	t.Run("not stale on same day", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 1, 20, 15, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), false)
	})

	t.Run("not stale after 1 business day", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday
		now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)             // Wednesday

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), false)
	})

	t.Run("not stale after 2 business days", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday
		now := time.Date(2026, 1, 22, 0, 0, 0, 0, time.UTC)             // Thursday

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), false)
	})

	t.Run("stale after 3 business days", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC) // Tuesday
		now := time.Date(2026, 1, 23, 0, 0, 0, 0, time.UTC)             // Friday

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), true)
	})

	t.Run("excludes weekends when calculating staleness", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC) // Thursday
		now := time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC)             // Monday (4 calendar days, 2 business days)

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), false) // Only 2 business days
	})

	t.Run("stale after weekend when threshold exceeded", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC) // Thursday
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)             // Tuesday (5 calendar days, 3 business days)

		td := NewFull("Task", PriorityA, false, nil, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), true) // 3 business days > 2
	})

	t.Run("not stale if no prioritised date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityA, false, nil, nil, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false)
	})
}

func TestIsStale_OtherPriorities(t *testing.T) {
	is := is.New(t)

	t.Run("Priority B not stale within 5 business days", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday
		now := time.Date(2026, 1, 17, 0, 0, 0, 0, time.UTC)          // Friday (4 business days)

		td := NewFull("Task", PriorityB, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false)
	})

	t.Run("Priority B stale after 5 business days", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)          // Next Monday (5 business days)

		td := NewFull("Task", PriorityB, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false) // Exactly 5 is not stale
	})

	t.Run("Priority B stale after more than 5 business days", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC) // Monday
		now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)          // Tuesday (6 business days)

		td := NewFull("Task", PriorityB, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), true)
	})

	t.Run("Priority C follows same rules as B", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityC, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), true)
	})

	t.Run("Priority D follows same rules as B", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 13, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityD, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), true)
	})

	t.Run("not stale if no creation date", func(t *testing.T) {
		is := is.New(t)
		now := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityB, false, nil, nil, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false)
	})
}

func TestIsStale_CompletedTasks(t *testing.T) {
	is := is.New(t)

	t.Run("completed Priority A task never stale", func(t *testing.T) {
		is := is.New(t)
		prioritisedDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		completionDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC) // Way in the future

		td := NewFull("Task", PriorityA, true, &completionDate, nil, nil, &prioritisedDate, nil, nil)
		is.Equal(td.IsStale(now), false)
	})

	t.Run("completed Priority B task never stale", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		completionDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityB, true, &completionDate, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false)
	})
}

func TestIsStale_NoPriority(t *testing.T) {
	is := is.New(t)

	t.Run("no priority task never stale", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		now := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

		td := NewFull("Task", PriorityNone, false, nil, &creationDate, nil, nil, nil, nil)
		is.Equal(td.IsStale(now), false)
	})
}
