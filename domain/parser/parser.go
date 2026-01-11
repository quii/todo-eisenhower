package parser

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	"github.com/quii/todo-eisenhower/domain/todo"
)

var (
	completedPrefix = regexp.MustCompile(`^x\s+`)
	priorityPattern = regexp.MustCompile(`^\(([A-D])\)\s+`)
	datePattern     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s+`)
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
	priority := todo.PriorityNone
	description := line

	// Check for completion marker
	if completedPrefix.MatchString(description) {
		completed = true
		description = completedPrefix.ReplaceAllString(description, "")
	}

	// Check for priority
	if priorityPattern.MatchString(description) {
		matches := priorityPattern.FindStringSubmatch(description)
		if len(matches) > 1 {
			priority = parsePriority(matches[1])
		}
		description = priorityPattern.ReplaceAllString(description, "")
	}

	// Remove completion date if present (after priority)
	if completed && datePattern.MatchString(description) {
		description = datePattern.ReplaceAllString(description, "")
	}

	description = strings.TrimSpace(description)

	if completed {
		return todo.NewCompleted(description, priority)
	}
	return todo.New(description, priority)
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
