package parser

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

// Parse reads todo.txt format from an io.Reader and returns a slice of Todos
func Parse(r io.Reader) ([]todo.Todo, error) {
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

func parseLine(line string) todo.Todo {
	completed := false
	var completionDate *time.Time
	priority := todo.PriorityNone
	description := line

	// Check for completion marker
	if completedPrefix.MatchString(description) {
		completed = true
		description = completedPrefix.ReplaceAllString(description, "")
	}

	// Extract and remove completion date if present at the beginning (format: x DATE (A) Description)
	if completed && datePattern.MatchString(description) {
		// Extract the date string
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			completionDate = &parsedDate
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

	// Extract and remove completion date if present after priority (format: x (A) DATE Description - backward compat)
	if completed && completionDate == nil && datePattern.MatchString(description) {
		dateStr := strings.TrimSpace(datePattern.FindString(description))
		if parsedDate, err := time.Parse("2006-01-02", dateStr); err == nil {
			completionDate = &parsedDate
		}
		description = datePattern.ReplaceAllString(description, "")
	}

	description = strings.TrimSpace(description)

	// Extract projects and contexts
	projects := extractTags(description, projectPattern)
	contexts := extractTags(description, contextPattern)

	// Create todo with tags (completed or not)
	if completed {
		if len(projects) > 0 || len(contexts) > 0 {
			return todo.NewCompletedWithTags(description, priority, completionDate, projects, contexts)
		}
		return todo.NewCompleted(description, priority, completionDate)
	}

	if len(projects) > 0 || len(contexts) > 0 {
		return todo.NewWithTags(description, priority, projects, contexts)
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
