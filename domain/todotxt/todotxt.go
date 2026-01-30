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

// DateFormat is the standard date format used in todo.txt files (ISO 8601)
const DateFormat = "2006-01-02"

// parseDate parses a date string in DateFormat and returns a pointer to the time.
// Returns nil if the string cannot be parsed.
func parseDate(s string) *time.Time {
	if t, err := time.Parse(DateFormat, s); err == nil {
		return &t
	}
	return nil
}

var (
	completedPrefix       = regexp.MustCompile(`^x\s+`)
	priorityPattern       = regexp.MustCompile(`^\(([A-E])\)\s+`)
	datePattern           = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s+`)
	projectPattern        = regexp.MustCompile(`\+(\w+)`)
	contextPattern        = regexp.MustCompile(`@(\w+)`)
	dueDatePattern        = regexp.MustCompile(`(?i)due:(\d{4}-\d{2}-\d{2})`)
	prioritisedDatePattern = regexp.MustCompile(`(?i)prioritised:(\d{4}-\d{2}-\d{2})`)
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
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		completionDate = parseDate(dateStr)
		description = datePattern.ReplaceAllString(description, "")
	}

	// For completed todos, next date (before priority) is creation date
	// Format: x COMP_DATE CREATION_DATE (A) Description
	if completed && datePattern.MatchString(description) {
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		creationDate = parseDate(dateStr)
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
		creationDate = parseDate(dateStr)
		description = datePattern.ReplaceAllString(description, "")
	}

	description = strings.TrimSpace(description)

	// Extract projects and contexts
	projects := extractTags(description, projectPattern)
	contexts := extractTags(description, contextPattern)

	// Extract due date (case-insensitive)
	var dueDate *time.Time
	if dueDatePattern.MatchString(description) {
		matches := dueDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			dueDate = parseDate(matches[1])
		}
	}

	// Extract prioritised date (case-insensitive)
	var prioritisedDate *time.Time
	if prioritisedDatePattern.MatchString(description) {
		matches := prioritisedDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			prioritisedDate = parseDate(matches[1])
		}
	}

	// Remove tags, due date, and prioritised date from description now that they're extracted
	description = projectPattern.ReplaceAllString(description, "")
	description = contextPattern.ReplaceAllString(description, "")
	description = dueDatePattern.ReplaceAllString(description, "")
	description = prioritisedDatePattern.ReplaceAllString(description, "")
	// Clean up extra whitespace
	description = strings.Join(strings.Fields(description), " ")
	description = strings.TrimSpace(description)

	// Use the comprehensive constructor with all extracted fields
	return todo.NewFull(description, priority, completed, completionDate, creationDate, dueDate, prioritisedDate, projects, contexts)
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
	case "E":
		return todo.PriorityE
	default:
		return todo.PriorityNone
	}
}

// parseDescriptionAndTags extracts clean description, projects, contexts, and due date from user input
func parseDescriptionAndTags(description string) (cleanDesc string, projects, contexts []string, dueDate, prioritisedDate *time.Time) {
	// Extract projects and contexts from description
	projects = extractTags(description, projectPattern)
	contexts = extractTags(description, contextPattern)

	// Extract due date (case-insensitive)
	if dueDatePattern.MatchString(description) {
		matches := dueDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			dueDate = parseDate(matches[1])
		}
	}

	// Extract prioritised date (case-insensitive) - but this is ignored in ParseNew/ParseEdit
	// We extract it here for consistency with file parsing, but ParseNew/ParseEdit should manage it themselves
	if prioritisedDatePattern.MatchString(description) {
		matches := prioritisedDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			prioritisedDate = parseDate(matches[1])
		}
	}

	// Remove tags, due date, and prioritised date from description to get clean text
	cleanDesc = projectPattern.ReplaceAllString(description, "")
	cleanDesc = contextPattern.ReplaceAllString(cleanDesc, "")
	cleanDesc = dueDatePattern.ReplaceAllString(cleanDesc, "")
	cleanDesc = prioritisedDatePattern.ReplaceAllString(cleanDesc, "")

	// Clean up extra whitespace
	cleanDesc = strings.Join(strings.Fields(cleanDesc), " ")
	cleanDesc = strings.TrimSpace(cleanDesc)

	return cleanDesc, projects, contexts, dueDate, prioritisedDate
}

// ParseNew creates a new todo from user input.
// It parses the description to extract tags and due date, and creates the appropriate todo with creation date.
// If priority is A (Do First), automatically sets prioritised date to creation date.
// This is the primary way to create todos from user input in the application.
func ParseNew(description string, priority todo.Priority, creationDate time.Time) todo.Todo {
	cleanDesc, projects, contexts, dueDate, _ := parseDescriptionAndTags(description)

	// Set prioritised date for Priority A (Do First quadrant)
	var prioritisedDate *time.Time
	if priority == todo.PriorityA {
		prioritisedDate = &creationDate
	}

	// Use comprehensive constructor with all extracted fields
	return todo.NewFull(cleanDesc, priority, false, nil, &creationDate, dueDate, prioritisedDate, projects, contexts)
}

// ParseEdit updates a todo from user input while preserving dates and completion status.
// It parses the new description to extract tags and due date, and creates an updated todo that preserves
// the original creation date, completion date, completion status, and prioritised date.
// Note: The due date from the new description REPLACES the original due date (it's not preserved from original).
// Note: The prioritised date is PRESERVED from the original (never changed by editing).
func ParseEdit(original todo.Todo, newDescription string, priority todo.Priority) todo.Todo {
	cleanDesc, projects, contexts, dueDate, _ := parseDescriptionAndTags(newDescription)

	// Preserve original dates and completion status
	creationDate := original.CreationDate()
	completionDate := original.CompletionDate()
	prioritisedDate := original.PrioritisedDate()
	isCompleted := original.IsCompleted()

	// Use comprehensive constructor with all fields
	return todo.NewFull(cleanDesc, priority, isCompleted, completionDate, creationDate, dueDate, prioritisedDate, projects, contexts)
}

// FormatForInput formats a todo for user input by combining description, tags, and due date.
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

	if dueDate := t.DueDate(); dueDate != nil {
		result.WriteString(" due:")
		result.WriteString(dueDate.Format(DateFormat))
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
