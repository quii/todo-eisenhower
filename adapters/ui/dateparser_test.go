package ui

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestParseDateShortcut(t *testing.T) {
	is := is.New(t)

	// Fixed reference time for testing: 2026-01-20 (Tuesday)
	refTime := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		shortcut string
		expected string
		wantErr  bool
	}{
		// Absolute shortcuts
		{
			name:     "today",
			shortcut: "today",
			expected: "2026-01-20",
			wantErr:  false,
		},
		{
			name:     "tod (abbreviated)",
			shortcut: "tod",
			expected: "2026-01-20",
			wantErr:  false,
		},
		{
			name:     "tomorrow",
			shortcut: "tomorrow",
			expected: "2026-01-21",
			wantErr:  false,
		},
		{
			name:     "tom (abbreviated)",
			shortcut: "tom",
			expected: "2026-01-21",
			wantErr:  false,
		},

		// Relative days
		{
			name:     "plus 3 days",
			shortcut: "+3d",
			expected: "2026-01-23",
			wantErr:  false,
		},
		{
			name:     "plus 1 day",
			shortcut: "+1d",
			expected: "2026-01-21",
			wantErr:  false,
		},
		{
			name:     "plus 7 days",
			shortcut: "+7d",
			expected: "2026-01-27",
			wantErr:  false,
		},

		// Relative weeks
		{
			name:     "plus 2 weeks",
			shortcut: "+2w",
			expected: "2026-02-03",
			wantErr:  false,
		},
		{
			name:     "plus 1 week",
			shortcut: "+1w",
			expected: "2026-01-27",
			wantErr:  false,
		},

		// Weekdays (next occurrence)
		{
			name:     "monday (next Monday)",
			shortcut: "monday",
			expected: "2026-01-26",
			wantErr:  false,
		},
		{
			name:     "mon (abbreviated)",
			shortcut: "mon",
			expected: "2026-01-26",
			wantErr:  false,
		},
		{
			name:     "friday",
			shortcut: "friday",
			expected: "2026-01-23",
			wantErr:  false,
		},
		{
			name:     "fri (abbreviated)",
			shortcut: "fri",
			expected: "2026-01-23",
			wantErr:  false,
		},
		{
			name:     "tuesday (today is Tuesday, should give next week)",
			shortcut: "tuesday",
			expected: "2026-01-27",
			wantErr:  false,
		},

		// Month abbreviations
		{
			name:     "jan25 (January 25 this year)",
			shortcut: "jan25",
			expected: "2026-01-25",
			wantErr:  false,
		},
		{
			name:     "feb14 (February 14 this year)",
			shortcut: "feb14",
			expected: "2026-02-14",
			wantErr:  false,
		},
		{
			name:     "dec31 (December 31 this year)",
			shortcut: "dec31",
			expected: "2026-12-31",
			wantErr:  false,
		},

		// End of period shortcuts (reference time: 2026-01-20, Tuesday in Q1)
		{
			name:     "endofmonth (end of January)",
			shortcut: "endofmonth",
			expected: "2026-01-31",
			wantErr:  false,
		},
		{
			name:     "endofquarter (end of Q1 - March 31)",
			shortcut: "endofquarter",
			expected: "2026-03-31",
			wantErr:  false,
		},
		{
			name:     "endofyear (December 31 of current year)",
			shortcut: "endofyear",
			expected: "2026-12-31",
			wantErr:  false,
		},

		// Already in correct format - pass through
		{
			name:     "ISO format already",
			shortcut: "2026-03-15",
			expected: "2026-03-15",
			wantErr:  false,
		},

		// Invalid inputs
		{
			name:     "invalid shortcut",
			shortcut: "invalidshortcut",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty string",
			shortcut: "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result, err := parseDateShortcut(tt.shortcut, refTime)

			if tt.wantErr {
				is.True(err != nil)
			} else {
				is.NoErr(err)
				is.Equal(result, tt.expected)
			}
		})
	}
}

func TestEndOfPeriodShortcuts(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name      string
		refDate   time.Time
		shortcut  string
		expected  string
	}{
		// End of Month tests - various months
		{
			name:     "endofmonth in January",
			refDate:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofmonth",
			expected: "2026-01-31",
		},
		{
			name:     "endofmonth in February (non-leap year)",
			refDate:  time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC),
			shortcut: "endofmonth",
			expected: "2026-02-28",
		},
		{
			name:     "endofmonth in February (leap year)",
			refDate:  time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			shortcut: "endofmonth",
			expected: "2024-02-29",
		},
		{
			name:     "endofmonth on last day of month",
			refDate:  time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
			shortcut: "endofmonth",
			expected: "2026-01-31",
		},

		// End of Quarter tests - all four quarters
		{
			name:     "endofquarter in Q1 (January)",
			refDate:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-03-31",
		},
		{
			name:     "endofquarter in Q1 (March)",
			refDate:  time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-03-31",
		},
		{
			name:     "endofquarter in Q2 (April)",
			refDate:  time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-06-30",
		},
		{
			name:     "endofquarter in Q3 (July)",
			refDate:  time.Date(2026, 7, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-09-30",
		},
		{
			name:     "endofquarter in Q4 (October)",
			refDate:  time.Date(2026, 10, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-12-31",
		},
		{
			name:     "endofquarter on last day of quarter",
			refDate:  time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
			shortcut: "endofquarter",
			expected: "2026-03-31",
		},

		// End of Year tests
		{
			name:     "endofyear in January",
			refDate:  time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofyear",
			expected: "2026-12-31",
		},
		{
			name:     "endofyear in December",
			refDate:  time.Date(2026, 12, 15, 0, 0, 0, 0, time.UTC),
			shortcut: "endofyear",
			expected: "2026-12-31",
		},
		{
			name:     "endofyear on last day of year",
			refDate:  time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
			shortcut: "endofyear",
			expected: "2026-12-31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result, err := parseDateShortcut(tt.shortcut, tt.refDate)
			is.NoErr(err)
			is.Equal(result, tt.expected)
		})
	}
}

func TestExpandDateShortcuts(t *testing.T) {
	is := is.New(t)

	// Fixed reference time for testing: 2026-01-20
	refTime := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tomorrow in description",
			input:    "Task description due:tomorrow +project",
			expected: "Task description due:2026-01-21 +project",
		},
		{
			name:     "expand +3d shortcut",
			input:    "Deploy app due:+3d @work",
			expected: "Deploy app due:2026-01-23 @work",
		},
		{
			name:     "expand friday shortcut",
			input:    "Meeting due:friday",
			expected: "Meeting due:2026-01-23",
		},
		{
			name:     "already ISO format - no change",
			input:    "Task due:2026-02-15 +project",
			expected: "Task due:2026-02-15 +project",
		},
		{
			name:     "no due date - no change",
			input:    "Task with no due date +project",
			expected: "Task with no due date +project",
		},
		{
			name:     "multiple due dates (should expand first one only)",
			input:    "Task due:tomorrow and due:friday",
			expected: "Task due:2026-01-21 and due:friday",
		},
		{
			name:     "invalid shortcut - leave as is",
			input:    "Task due:invalidshortcut",
			expected: "Task due:invalidshortcut",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result := expandDateShortcuts(tt.input, refTime)
			is.Equal(result, tt.expected)
		})
	}
}
