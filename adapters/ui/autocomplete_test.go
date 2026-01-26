package ui

import (
	"testing"

	"github.com/matryer/is"
)

func TestDetectDueTrigger(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name             string
		input            string
		expectedPartial  string
		expectedFound    bool
	}{
		{
			name:            "detects due: with no text after",
			input:           "Task description due:",
			expectedPartial: "",
			expectedFound:   true,
		},
		{
			name:            "detects due: with partial text",
			input:           "Task description due:tom",
			expectedPartial: "tom",
			expectedFound:   true,
		},
		{
			name:            "detects due: case insensitive",
			input:           "Task description DUE:tomorrow",
			expectedPartial: "tomorrow",
			expectedFound:   true,
		},
		{
			name:            "no due: found",
			input:           "Task description +project",
			expectedPartial: "",
			expectedFound:   false,
		},
		{
			name:            "due: with space after stops trigger",
			input:           "Task description due:tomorrow another task",
			expectedPartial: "",
			expectedFound:   false,
		},
		{
			name:            "multiple due: uses last one",
			input:           "due:today task due:tom",
			expectedPartial: "tom",
			expectedFound:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			partial, found := detectDueTrigger(tt.input)
			is.Equal(found, tt.expectedFound)
			is.Equal(partial, tt.expectedPartial)
		})
	}
}

func TestFilterDateShortcuts(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name           string
		prefix         string
		expectedCount  int
		expectedFirst  string
	}{
		{
			name:          "empty prefix returns all shortcuts",
			prefix:        "",
			expectedCount: 18, // today, tomorrow, endofmonth, endofquarter, endofnextquarter, endofyear, +1d, +3d, +7d, +1w, +2w, 7 weekdays
			expectedFirst: "today",
		},
		{
			name:          "prefix 't' filters to today, tomorrow, tuesday, thursday",
			prefix:        "t",
			expectedCount: 4,
			expectedFirst: "today",
		},
		{
			name:          "prefix 'tod' filters to today",
			prefix:        "tod",
			expectedCount: 1,
			expectedFirst: "today",
		},
		{
			name:          "prefix '+1' filters to +1d and +1w",
			prefix:        "+1",
			expectedCount: 2,
			expectedFirst: "+1d",
		},
		{
			name:          "prefix 'fri' filters to friday",
			prefix:        "fri",
			expectedCount: 1,
			expectedFirst: "friday",
		},
		{
			name:          "prefix 'end' filters to all end shortcuts",
			prefix:        "end",
			expectedCount: 4,
			expectedFirst: "endofmonth",
		},
		{
			name:          "prefix 'endof' filters to all end shortcuts",
			prefix:        "endof",
			expectedCount: 4,
			expectedFirst: "endofmonth",
		},
		{
			name:          "prefix 'endofn' filters to endofnextquarter",
			prefix:        "endofn",
			expectedCount: 1,
			expectedFirst: "endofnextquarter",
		},
		{
			name:          "no matches returns empty",
			prefix:        "xyz",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			results := filterDateShortcuts(tt.prefix)
			is.Equal(len(results), tt.expectedCount)
			if tt.expectedCount > 0 {
				is.Equal(results[0], tt.expectedFirst)
			}
		})
	}
}

func TestCompleteDateShortcut(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name              string
		input             string
		completedShortcut string
		expected          string
	}{
		{
			name:              "completes partial shortcut",
			input:             "Task due:tom",
			completedShortcut: "tomorrow",
			expected:          "Task due:tomorrow ",
		},
		{
			name:              "completes with empty partial",
			input:             "Task due:",
			completedShortcut: "today",
			expected:          "Task due:today ",
		},
		{
			name:              "handles case insensitive due:",
			input:             "Task DUE:fri",
			completedShortcut: "friday",
			expected:          "Task DUE:friday ",
		},
		{
			name:              "no due: found returns original",
			input:             "Task +project",
			completedShortcut: "today",
			expected:          "Task +project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result := completeDateShortcut(tt.input, tt.completedShortcut)
			is.Equal(result, tt.expected)
		})
	}
}
