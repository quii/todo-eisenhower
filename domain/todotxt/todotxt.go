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

// parseDescriptionAndTags extracts clean description, projects, and contexts from user input
func parseDescriptionAndTags(description string) (cleanDesc string, projects, contexts []string) {
	// Extract projects and contexts from description
	projects = extractTags(description, projectPattern)
	contexts = extractTags(description, contextPattern)

	// Remove tags from description to get clean text
	cleanDesc = projectPattern.ReplaceAllString(description, "")
	cleanDesc = contextPattern.ReplaceAllString(cleanDesc, "")

	// Clean up extra whitespace
	cleanDesc = strings.Join(strings.Fields(cleanDesc), " ")
	cleanDesc = strings.TrimSpace(cleanDesc)

	return cleanDesc, projects, contexts
}

// ParseNew creates a new todo from user input.
// It parses the description to extract tags and creates the appropriate todo with creation date.
// This is the primary way to create todos from user input in the application.
func ParseNew(description string, priority todo.Priority, creationDate time.Time) todo.Todo {
	cleanDesc, projects, contexts := parseDescriptionAndTags(description)

	// Create todo with appropriate constructor based on what we have
	if len(projects) > 0 || len(contexts) > 0 {
		return todo.NewWithTagsAndDates(cleanDesc, priority, &creationDate, projects, contexts)
	}

	return todo.NewWithCreationDate(cleanDesc, priority, &creationDate)
}

// ParseEdit updates a todo from user input while preserving dates and completion status.
// It parses the new description to extract tags and creates an updated todo that preserves
// the original creation date, completion date, and completion status.
func ParseEdit(original todo.Todo, newDescription string, priority todo.Priority) todo.Todo {
	cleanDesc, projects, contexts := parseDescriptionAndTags(newDescription)

	// Preserve original dates and completion status
	creationDate := original.CreationDate()
	completionDate := original.CompletionDate()
	isCompleted := original.IsCompleted()

	// Create updated todo with appropriate constructor
	if isCompleted {
		if len(projects) > 0 || len(contexts) > 0 {
			return todo.NewCompletedWithTagsAndDates(cleanDesc, priority, completionDate, creationDate, projects, contexts)
		}
		return todo.NewCompletedWithDates(cleanDesc, priority, completionDate, creationDate)
	}

	if len(projects) > 0 || len(contexts) > 0 {
		if creationDate != nil {
			return todo.NewWithTagsAndDates(cleanDesc, priority, creationDate, projects, contexts)
		}
		return todo.NewWithTags(cleanDesc, priority, projects, contexts)
	}

	if creationDate != nil {
		return todo.NewWithCreationDate(cleanDesc, priority, creationDate)
	}

	return todo.New(cleanDesc, priority)
}

// FormatForInput formats a todo for user input by combining description and tags.
// This is the inverse of parsing - it converts a structured todo back to input format.
func FormatForInput(t todo.Todo) string {
	var result strings.Builder
	result.WriteString(t.Description())

	for _, project := range t.Projects() {
		result.WriteString(" +")
		result.WriteString(project)
	}

	for _, context := range t.Contexts() {
		result.WriteString(" @")
		result.WriteString(context)
	}

	return result.String()
}

// ParseDescription extracts the description text and tags from a todo description string.
// It returns the clean description (with tags removed), projects, and contexts.
//
// Deprecated: Use ParseNew instead for creating new todos.
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
