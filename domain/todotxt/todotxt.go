// Package todotxt provides encoding and decoding for the todo.txt format.
// It follows the convention of encoding packages like encoding/json.
package todotxt

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
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
	recPattern             = regexp.MustCompile(`(?i)rec:(\+?\d+[dwmy])`)
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

	// Extract due date (case-insensitive)
	var dueDate *time.Time
	if dueDatePattern.MatchString(description) {
		matches := dueDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			dueDate = parseDate(matches[1])
		}
	}

	// Extract recurrence tag (only valid when due date is present)
	// Must be extracted/removed BEFORE project extraction because rec:+2w contains +2w which matches project pattern
	var recurrence *todo.Recurrence
	if recPattern.MatchString(description) {
		matches := recPattern.FindStringSubmatch(description)
		if len(matches) > 1 && dueDate != nil {
			recurrence = parseRecurrence(matches[1])
		}
		description = recPattern.ReplaceAllString(description, "")
	}

	// Extract prioritised date (case-insensitive)
	var prioritisedDate *time.Time
	if prioritisedDatePattern.MatchString(description) {
		matches := prioritisedDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			prioritisedDate = parseDate(matches[1])
		}
	}

	// Extract projects and contexts (after rec removal to avoid +2w matching as project)
	projects := extractTags(description, projectPattern)
	contexts := extractTags(description, contextPattern)

	// Remove tags, due date, and prioritised date from description now that they're extracted
	description = projectPattern.ReplaceAllString(description, "")
	description = contextPattern.ReplaceAllString(description, "")
	description = dueDatePattern.ReplaceAllString(description, "")
	description = prioritisedDatePattern.ReplaceAllString(description, "")
	// Clean up extra whitespace
	description = strings.Join(strings.Fields(description), " ")
	description = strings.TrimSpace(description)

	// Use the comprehensive constructor with all extracted fields
	return todo.NewFull(description, priority, completed, completionDate, creationDate, dueDate, prioritisedDate, recurrence, projects, contexts)
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

// parseRecurrence parses a rec: tag value like "+2w" or "1m" into a Recurrence
func parseRecurrence(s string) *todo.Recurrence {
	if s == "" {
		return nil
	}

	relative := false
	if s[0] == '+' {
		relative = true
		s = s[1:]
	}

	// Last character is the unit
	if len(s) < 2 {
		return nil
	}

	unitChar := s[len(s)-1]
	intervalStr := s[:len(s)-1]

	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		return nil
	}

	var unit todo.RecurrenceUnit
	switch unitChar {
	case 'd', 'D':
		unit = todo.Day
	case 'w', 'W':
		unit = todo.Week
	case 'm', 'M':
		unit = todo.Month
	case 'y', 'Y':
		unit = todo.Year
	default:
		return nil
	}

	rec := todo.NewRecurrence(relative, interval, unit)
	return &rec
}

// parsedDescription holds the result of parsing a user input description
type parsedDescription struct {
	description     string
	projects        []string
	contexts        []string
	dueDate         *time.Time
	prioritisedDate *time.Time
	recurrence      *todo.Recurrence
}

// parseDescriptionAndTags extracts clean description, projects, contexts, due date, and recurrence from user input
func parseDescriptionAndTags(description string) parsedDescription {
	var result parsedDescription

	// Extract due date (case-insensitive)
	if dueDatePattern.MatchString(description) {
		matches := dueDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			result.dueDate = parseDate(matches[1])
		}
	}

	// Extract recurrence (only valid when due date is present)
	// Must be extracted/removed BEFORE project extraction because rec:+2w contains +2w which matches project pattern
	if recPattern.MatchString(description) {
		matches := recPattern.FindStringSubmatch(description)
		if len(matches) > 1 && result.dueDate != nil {
			result.recurrence = parseRecurrence(matches[1])
		}
		description = recPattern.ReplaceAllString(description, "")
	}

	// Extract projects and contexts (after rec removal to avoid +2w matching as project)
	result.projects = extractTags(description, projectPattern)
	result.contexts = extractTags(description, contextPattern)

	// Extract prioritised date (case-insensitive) - but this is ignored in ParseNew/ParseEdit
	// We extract it here for consistency with file parsing, but ParseNew/ParseEdit should manage it themselves
	if prioritisedDatePattern.MatchString(description) {
		matches := prioritisedDatePattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			result.prioritisedDate = parseDate(matches[1])
		}
	}

	// Remove tags, due date, and prioritised date from description to get clean text
	cleanDesc := projectPattern.ReplaceAllString(description, "")
	cleanDesc = contextPattern.ReplaceAllString(cleanDesc, "")
	cleanDesc = dueDatePattern.ReplaceAllString(cleanDesc, "")
	cleanDesc = prioritisedDatePattern.ReplaceAllString(cleanDesc, "")

	// Clean up extra whitespace
	cleanDesc = strings.Join(strings.Fields(cleanDesc), " ")
	result.description = strings.TrimSpace(cleanDesc)

	return result
}

// ParseNew creates a new todo from user input.
// It parses the description to extract tags and due date, and creates the appropriate todo with creation date.
// If priority is A (Do First), automatically sets prioritised date to creation date.
// This is the primary way to create todos from user input in the application.
func ParseNew(description string, priority todo.Priority, creationDate time.Time) todo.Todo {
	parsed := parseDescriptionAndTags(description)

	// Set prioritised date for Priority A (Do First quadrant)
	var prioritisedDate *time.Time
	if priority == todo.PriorityA {
		prioritisedDate = &creationDate
	}

	// Use comprehensive constructor with all extracted fields
	return todo.NewFull(parsed.description, priority, false, nil, &creationDate, parsed.dueDate, prioritisedDate, parsed.recurrence, parsed.projects, parsed.contexts)
}

// ParseEdit updates a todo from user input while preserving dates and completion status.
// It parses the new description to extract tags and due date, and creates an updated todo that preserves
// the original creation date, completion date, completion status, and prioritised date.
// Note: The due date from the new description REPLACES the original due date (it's not preserved from original).
// Note: The prioritised date is PRESERVED from the original (never changed by editing).
func ParseEdit(original todo.Todo, newDescription string, priority todo.Priority) todo.Todo {
	parsed := parseDescriptionAndTags(newDescription)

	// Preserve original dates and completion status
	creationDate := original.CreationDate()
	completionDate := original.CompletionDate()
	prioritisedDate := original.PrioritisedDate()
	isCompleted := original.IsCompleted()

	// Use comprehensive constructor with all fields
	return todo.NewFull(parsed.description, priority, isCompleted, completionDate, creationDate, parsed.dueDate, prioritisedDate, parsed.recurrence, parsed.projects, parsed.contexts)
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

	if rec := t.Recurrence(); rec != nil {
		result.WriteString(" rec:")
		result.WriteString(rec.String())
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
