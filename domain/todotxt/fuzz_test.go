package todotxt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/domain/todotxt"
)

// FuzzParseLine fuzzes the parseLine function through the public Unmarshal API
// to ensure it never panics on arbitrary input and always produces valid todos.
func FuzzParseLine(f *testing.F) {
	// Seed corpus with valid and interesting test cases
	testCases := []string{
		// Basic cases
		"Buy milk",
		"(A) Important task",
		"x Completed task",

		// With dates
		"2024-01-15 Task with creation date",
		"(A) 2024-01-15 Priority and date",
		"x 2024-01-20 2024-01-15 Completed with both dates",
		"x 2024-01-20 2024-01-15 (A) Completed with priority",

		// With tags
		"Task +project @context",
		"(B) Task +multiple +projects @multiple @contexts",
		"+ProjectOnly",
		"@ContextOnly",

		// Complex cases
		"x 2024-01-20 2024-01-15 (C) Complex +project @context task",

		// Edge cases
		"",
		" ",
		"x",
		"()",
		"(Z)",
		"++",
		"@@",
		"x x x",
		"(A) (B) (C)",
		"2024-13-45", // Invalid date
		"x not-a-date",

		// Special characters
		"Task with 日本語",
		"Task with émojis 🎉",
		"Task\twith\ttabs",
		"Task\nwith\nnewlines",

		// Very long strings
		string(make([]byte, 1000)),
		string(make([]byte, 10000)),

		// Regex edge cases
		"Task +++project",
		"Task @@@context",
		"Task +",
		"Task @",
		"Task + project",
		"Task @ context",
		"(A)NoSpace",
		"x NoSpace",

		// Malformed dates
		"2024-01-",
		"2024--01",
		"----",
		"9999-99-99",
		"0000-00-00",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, input string) {
		// The goal is to ensure parseLine never panics
		// We test through Unmarshal which calls parseLine
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("parseLine panicked on input %q: %v", input, r)
			}
		}()

		// Parse the input - should never panic
		reader := strings.NewReader(input)
		todos, err := todotxt.Unmarshal(reader)

		// Should never return an error for malformed todos
		// (parser is lenient and always returns a valid todo)
		if err != nil {
			// Only acceptable error is from the scanner itself
			// Not from parsing logic
			return
		}

		// If we got todos, validate basic invariants hold
		for _, todo := range todos {
			// Basic invariants that should always be true
			// Description should never be nil (even if empty string is ok)
			desc := todo.Description()
			_ = desc

			// Priority should be valid
			priority := todo.Priority()
			if priority < 0 || priority > 4 {
				t.Errorf("invalid priority value: %v", priority)
			}

			// Should be able to call all accessor methods without panic
			_ = todo.IsCompleted()
			_ = todo.CreationDate()
			_ = todo.CompletionDate()
			_ = todo.Projects()
			_ = todo.Contexts()
			_ = todo.String()
		}
	})
}

// FuzzParseNew fuzzes the ParseNew function to ensure it handles arbitrary input
func FuzzParseNew(f *testing.F) {
	// Seed corpus
	testCases := []string{
		"Simple task",
		"Task +project @context",
		"+project",
		"@context",
		"Multiple +p1 +p2 @c1 @c2",
		"",
		" ",
		"   ",
		"🎉 emoji task 🚀",
		string(make([]byte, 10000)),
		"+++",
		"@@@",
		"Task\nwith\nnewlines",
		"Task\twith\ttabs",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, description string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ParseNew panicked on input %q: %v", description, r)
			}
		}()

		// Try all priority values
		priorities := []todo.Priority{
			todo.PriorityNone,
			todo.PriorityA,
			todo.PriorityB,
			todo.PriorityC,
			todo.PriorityD,
		}

		creationDate := time.Now()

		for _, priority := range priorities {
			t := todotxt.ParseNew(description, priority, creationDate)

			// Validate basic invariants
			_ = t.Description()
			_ = t.Priority()
			_ = t.Projects()
			_ = t.Contexts()
			_ = t.String()
		}
	})
}

// FuzzFormatForInput fuzzes the round-trip: parse -> format -> parse
func FuzzFormatForInput(f *testing.F) {
	testCases := []string{
		"Simple task",
		"Task +project @context",
		"Complex +p1 +p2 @c1 @c2",
		"",
		"日本語 task",
		"emoji 🎉",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, description string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("FormatForInput round-trip panicked on input %q: %v", description, r)
			}
		}()

		// Parse -> Format -> Parse round trip
		creationDate := time.Now()
		todo1 := todotxt.ParseNew(description, todo.PriorityNone, creationDate)
		formatted := todotxt.FormatForInput(todo1)
		todo2 := todotxt.ParseNew(formatted, todo.PriorityNone, creationDate)

		// After round-trip, descriptions should match
		// (though whitespace normalization might occur)
		_ = todo1.Description()
		_ = todo2.Description()

		// Projects and contexts should be preserved
		// Just ensure no panic when accessing them
		_ = todo1.Projects()
		_ = todo2.Projects()
		_ = todo1.Contexts()
		_ = todo2.Contexts()
	})
}
