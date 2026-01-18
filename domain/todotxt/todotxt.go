// Package todotxt provides encoding and decoding for the todo.txt format.
// It follows the convention of encoding packages like encoding/json.
package todotxt

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
)

var (
	completedPrefix = regexp.MustCompile(`^x\s+`)
	priorityPattern = regexp.MustCompile(`^\(([A-D])\)\s+`)
	datePattern     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s+`)
	projectPattern  = regexp.MustCompile(`\+(\w+)`)
	contextPattern  = regexp.MustCompile(`@(\w+)`)
)

// Unmarshal reads todo.txt format from an io.Reader and returns a slice of Todos.
// This is the inverse operation of Marshal.
func Unmarshal(r io.Reader) ([]todo.Todo, error) {
	var todos []todo.Todo

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		t := parseLine(line)
		todos = append(todos, t)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// Marshal writes todos in todo.txt format to an io.Writer.
// This is the inverse operation of Unmarshal.
func Marshal(w io.Writer, todos []todo.Todo) error {
	for _, t := range todos {
		if _, err := w.Write([]byte(t.String())); err != nil {
			return err
		}
	}
	return nil
}

func parseLine(line string) todo.Todo {
	completed := false
	var completionDate *time.Time
	var creationDate *time.Time
	priority := todo.PriorityNone
	description := line

	// Check for completion marker
	if completedPrefix.MatchString(description) {
		completed = true
		description = completedPrefix.ReplaceAllString(description, "")
	}

	// Extract and remove completion date if present at the beginning (format: x DATE ...)
	if completed && datePattern.MatchString(description) {
		// Extract the date string
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			completionDate = &parsedDate
		}
		description = datePattern.ReplaceAllString(description, "")
	}

	// For completed todos, next date (before priority) is creation date
	// Format: x COMP_DATE CREATION_DATE (A) Description
	if completed && datePattern.MatchString(description) {
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			creationDate = &parsedDate
		}
		description = datePattern.ReplaceAllString(description, "")
	}

	// Check for priority
	if priorityPattern.MatchString(description) {
		matches := priorityPattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			priority = parsePriority(matches[1])
		}
		description = priorityPattern.ReplaceAllString(description, "")
	}

	// After priority, any date is the creation date
	// Format: (A) CREATION_DATE Description (for active todos or if not parsed yet)
	if creationDate == nil && datePattern.MatchString(description) {
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			creationDate = &parsedDate
		}
		description = datePattern.ReplaceAllString(description, "")
	}

	description = strings.TrimSpace(description)

	// Extract projects and contexts
	projects := extractTags(description, projectPattern)
	contexts := extractTags(description, contextPattern)

	// Remove tags from description now that they're extracted
	description = projectPattern.ReplaceAllString(description, "")
	description = contextPattern.ReplaceAllString(description, "")
	// Clean up extra whitespace
	description = strings.Join(strings.Fields(description), " ")
	description = strings.TrimSpace(description)

	// Create todo with appropriate constructor based on what we have
	if completed {
		if len(projects) > 0 || len(contexts) > 0 {
			return todo.NewCompletedWithTagsAndDates(description, priority, completionDate, creationDate, projects, contexts)
		}
		return todo.NewCompletedWithDates(description, priority, completionDate, creationDate)
	}

	if len(projects) > 0 || len(contexts) > 0 {
		if creationDate != nil {
			return todo.NewWithTagsAndDates(description, priority, creationDate, projects, contexts)
		}
		return todo.NewWithTags(description, priority, projects, contexts)
	}

	if creationDate != nil {
		return todo.NewWithCreationDate(description, priority, creationDate)
	}

	return todo.New(description, priority)
}

// extractTags extracts all matching tags using the given pattern
func extractTags(text string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(text, -1)
	tags := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			tags = append(tags, match[1])
		}
	}
	return tags
}

func parsePriority(p string) todo.Priority {
	switch p {
	case "A":
		return todo.PriorityA
	case "B":
		return todo.PriorityB
	case "C":
		return todo.PriorityC
	case "D":
		return todo.PriorityD
	default:
		return todo.PriorityNone
	}
}

// ParseDescription extracts the description text and tags from a todo description string.
// It returns the clean description (with tags removed), projects, and contexts.
// This is useful when creating new todos from user input.
func ParseDescription(description string) (cleanDesc string, projects, contexts []string) {
	// Extract projects and contexts
	projects = extractTags(description, projectPattern)
	contexts = extractTags(description, contextPattern)

	// Remove tags from description
	cleanDesc = projectPattern.ReplaceAllString(description, "")
	cleanDesc = contextPattern.ReplaceAllString(cleanDesc, "")

	// Clean up extra whitespace
	cleanDesc = strings.Join(strings.Fields(cleanDesc), " ")
	cleanDesc = strings.TrimSpace(cleanDesc)

	return cleanDesc, projects, contexts
}
