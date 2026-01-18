package todo_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// Property: ∀ todo: Parse(todo.String()) should reconstruct an equivalent todo
// This demonstrates the symmetry between serialization (String) and deserialization (Parse)
func TestPropertyParseFormatSymmetry(t *testing.T) {
	is := is.New(t)
	//nolint:gosec // G404: Using weak random for property testing, not security
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Run 100 random tests
	for range 100 {
		// Generate random todo
		original := generateRandomTodo(rng)

		// Format the todo to string
		formatted := original.String()

		// Parse it back
		parsed, err := todotxt.Unmarshal(strings.NewReader(formatted))
		is.NoErr(err)
		is.Equal(len(parsed), 1) // should parse exactly one todo

		// Round-trip: format again
		reformatted := parsed[0].String()

		// The formatted strings should be identical (idempotent)
		is.Equal(reformatted, formatted)
	}
}

// Property: ∀ todo: Parse(todo.String()) preserves all serializable properties
// Note: We check that a second round-trip yields the same todo, not that the parsed
// todo matches the original in-memory representation (dates may lose precision)
func TestPropertyParsePreservesProperties(t *testing.T) {
	is := is.New(t)
	//nolint:gosec // G404: Using weak random for property testing, not security
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Run 100 random tests
	for range 100 {
		original := generateRandomTodo(rng)

		// First round-trip
		formatted := original.String()
		parsed, err := todotxt.Unmarshal(strings.NewReader(formatted))
		is.NoErr(err)
		is.Equal(len(parsed), 1)

		first := parsed[0]

		// Second round-trip to verify stability
		reformatted := first.String()
		reparsed, err := todotxt.Unmarshal(strings.NewReader(reformatted))
		is.NoErr(err)
		is.Equal(len(reparsed), 1)

		second := reparsed[0]

		// After stabilization, all properties should match
		is.Equal(second.Description(), first.Description())
		is.Equal(second.Priority(), first.Priority())
		is.Equal(second.IsCompleted(), first.IsCompleted())

		// Compare projects
		firstProjects, secondProjects := first.Projects(), second.Projects()
		is.Equal(len(secondProjects), len(firstProjects))
		for j := range firstProjects {
			is.Equal(secondProjects[j], firstProjects[j])
		}

		// Compare contexts
		firstContexts, secondContexts := first.Contexts(), second.Contexts()
		is.Equal(len(secondContexts), len(firstContexts))
		for j := range firstContexts {
			is.Equal(secondContexts[j], firstContexts[j])
		}

		// Compare dates
		is.True(datesEqual(second.CreationDate(), first.CreationDate()))
		is.True(datesEqual(second.CompletionDate(), first.CompletionDate()))
	}
}

// generateRandomTodo creates a random todo for property testing
func generateRandomTodo(rng *rand.Rand) todo.Todo {
	// Random description (no newlines or special chars that would break parsing)
	description := randomAlphanumeric(rng, 5+rng.Intn(20))

	// Random priority
	priorities := []todo.Priority{
		todo.PriorityNone,
		todo.PriorityA,
		todo.PriorityB,
		todo.PriorityC,
		todo.PriorityD,
	}
	priority := priorities[rng.Intn(len(priorities))]

	// Random completion status
	completed := rng.Intn(2) == 1

	// Random dates (or nil)
	var creationDate *time.Time
	var completionDate *time.Time

	if rng.Intn(2) == 1 {
		d := time.Now().AddDate(0, 0, -rng.Intn(365))
		// Truncate to day precision to match todo.txt format
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		creationDate = &d
	}

	if completed && rng.Intn(2) == 1 {
		d := time.Now().AddDate(0, 0, -rng.Intn(30))
		// Truncate to day precision
		d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
		completionDate = &d
	}

	// Random projects (0-3 tags)
	var projects []string
	numProjects := rng.Intn(4)
	for range numProjects {
		projects = append(projects, "Project"+randomAlphanumeric(rng, 3))
	}

	// Random contexts (0-3 tags)
	var contexts []string
	numContexts := rng.Intn(4)
	for range numContexts {
		contexts = append(contexts, "context"+randomAlphanumeric(rng, 3))
	}

	// Create todo based on what we have
	if completed {
		if len(projects) > 0 || len(contexts) > 0 {
			return todo.NewCompletedWithTagsAndDates(description, priority, completionDate, creationDate, projects, contexts)
		}
		return todo.NewCompletedWithDates(description, priority, completionDate, creationDate)
	}

	if len(projects) > 0 || len(contexts) > 0 {
		return todo.NewWithTagsAndDates(description, priority, creationDate, projects, contexts)
	}
	if creationDate != nil {
		return todo.NewWithCreationDate(description, priority, creationDate)
	}
	return todo.New(description, priority)
}

// randomAlphanumeric generates a random alphanumeric string
func randomAlphanumeric(rng *rand.Rand, length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rng.Intn(len(chars))]
	}
	return string(result)
}

// datesEqual compares two dates at day precision (ignoring time)
func datesEqual(a, b *time.Time) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	// Compare as strings in todo.txt format (day precision)
	return a.Format("2006-01-02") == b.Format("2006-01-02")
}
