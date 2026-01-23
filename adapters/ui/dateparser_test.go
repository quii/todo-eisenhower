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
