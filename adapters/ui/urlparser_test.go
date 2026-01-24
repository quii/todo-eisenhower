package ui

import (
	"testing"

	"github.com/matryer/is"
)

func TestExtractURLs(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "single HTTPS URL",
			text:     "Review doc https://docs.google.com/document/d/abc123",
			expected: []string{"https://docs.google.com/document/d/abc123"},
		},
		{
			name:     "single HTTP URL",
			text:     "Check http://example.com for details",
			expected: []string{"http://example.com"},
		},
		{
			name: "multiple URLs",
			text: "Review PR https://github.com/user/repo/pull/42 and staging https://staging.example.com",
			expected: []string{
				"https://github.com/user/repo/pull/42",
				"https://staging.example.com",
			},
		},
		{
			name:     "URL with query parameters",
			text:     "Search https://google.com/search?q=test&page=1",
			expected: []string{"https://google.com/search?q=test&page=1"},
		},
		{
			name:     "URL with anchor",
			text:     "Docs at https://example.com/page#section-1",
			expected: []string{"https://example.com/page#section-1"},
		},
		{
			name:     "URL at end of string",
			text:     "Check out https://example.com",
			expected: []string{"https://example.com"},
		},
		{
			name:     "URL at start of string",
			text:     "https://example.com is the site",
			expected: []string{"https://example.com"},
		},
		{
			name:     "no URLs",
			text:     "Regular task with no links +project",
			expected: nil,
		},
		{
			name:     "empty string",
			text:     "",
			expected: nil,
		},
		{
			name:     "URL with port",
			text:     "Local server https://localhost:8080/api",
			expected: []string{"https://localhost:8080/api"},
		},
		{
			name:     "URL with path and extension",
			text:     "Download https://example.com/files/report.pdf",
			expected: []string{"https://example.com/files/report.pdf"},
		},
		{
			name:     "mixed with tags",
			text:     "Review https://github.com/user/repo +project @work",
			expected: []string{"https://github.com/user/repo"},
		},
		{
			name:     "URL in parentheses",
			text:     "See (https://example.com) for more",
			expected: []string{"https://example.com)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			result := extractURLs(tt.text)
			is.Equal(result, tt.expected)
		})
	}
}

func TestRenderClickableURL(t *testing.T) {
	is := is.New(t)

	url := "https://example.com"
	result := renderClickableURL(url)

	// Should contain OSC 8 escape sequences
	is.True(result != url) // Should be different from plain URL
	is.True(len(result) > len(url)) // Should be longer due to escape codes

	// Should contain the URL twice (once in escape sequence, once as display text)
	expected := "\x1b]8;;https://example.com\x1b\\https://example.com\x1b]8;;\x1b\\"
	is.Equal(result, expected)
}
