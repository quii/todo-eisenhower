// Package todo provides the Todo domain model following the todo.txt specification.
package todo

import "time"

// Priority represents the priority level of a todo item
type Priority int

const (
	PriorityNone Priority = iota
	PriorityA
	PriorityB
	PriorityC
	PriorityD
)

// Todo represents a single todo item
type Todo struct {
	description    string
	priority       Priority
	completed      bool
	completionDate *time.Time // nil if not completed or no date recorded
	creationDate   *time.Time // nil if no creation date recorded
	projects       []string
	contexts       []string
}

// New creates a new Todo with the given description and priority
func New(description string, priority Priority) Todo {
	return Todo{
		description:    description,
		priority:       priority,
		completed:      false,
		completionDate: nil,
		creationDate:   nil,
		projects:       []string{},
		contexts:       []string{},
	}
}

// NewWithCreationDate creates a new Todo with creation date
func NewWithCreationDate(description string, priority Priority, creationDate *time.Time) Todo {
	return Todo{
		description:    description,
		priority:       priority,
		completed:      false,
		completionDate: nil,
		creationDate:   creationDate,
		projects:       []string{},
		contexts:       []string{},
	}
}

// NewCompleted creates a new completed Todo with the given description, priority, and optional completion date
func NewCompleted(description string, priority Priority, completionDate *time.Time) Todo {
	return Todo{
		description:    description,
		priority:       priority,
		completed:      true,
		completionDate: completionDate,
		creationDate:   nil,
		projects:       []string{},
		contexts:       []string{},
	}
}

// NewCompletedWithDates creates a new completed Todo with both completion and creation dates
func NewCompletedWithDates(description string, priority Priority, completionDate, creationDate *time.Time) Todo {
	return Todo{
		description:    description,
		priority:       priority,
		completed:      true,
		completionDate: completionDate,
		creationDate:   creationDate,
		projects:       []string{},
		contexts:       []string{},
	}
}

// NewWithTags creates a new Todo with projects and contexts
func NewWithTags(description string, priority Priority, projects, contexts []string) Todo {
	if projects == nil {
		projects = []string{}
	}
	if contexts == nil {
		contexts = []string{}
	}
	return Todo{
		description:    description,
		priority:       priority,
		completed:      false,
		completionDate: nil,
		creationDate:   nil,
		projects:       projects,
		contexts:       contexts,
	}
}

// NewWithTagsAndDates creates a new Todo with tags and dates
func NewWithTagsAndDates(description string, priority Priority, creationDate *time.Time, projects, contexts []string) Todo {
	if projects == nil {
		projects = []string{}
	}
	if contexts == nil {
		contexts = []string{}
	}
	return Todo{
		description:    description,
		priority:       priority,
		completed:      false,
		completionDate: nil,
		creationDate:   creationDate,
		projects:       projects,
		contexts:       contexts,
	}
}

// NewCompletedWithTags creates a new completed Todo with tags and optional completion date
func NewCompletedWithTags(description string, priority Priority, completionDate *time.Time, projects, contexts []string) Todo{
	if projects == nil {
		projects = []string{}
	}
	if contexts == nil {
		contexts = []string{}
	}
	return Todo{
		description:    description,
		priority:       priority,
		completed:      true,
		completionDate: completionDate,
		creationDate:   nil,
		projects:       projects,
		contexts:       contexts,
	}
}

// NewCompletedWithTagsAndDates creates a new completed Todo with tags and both dates
func NewCompletedWithTagsAndDates(description string, priority Priority, completionDate, creationDate *time.Time, projects, contexts []string) Todo {
	if projects == nil {
		projects = []string{}
	}
	if contexts == nil {
		contexts = []string{}
	}
	return Todo{
		description:    description,
		priority:       priority,
		completed:      true,
		completionDate: completionDate,
		creationDate:   creationDate,
		projects:       projects,
		contexts:       contexts,
	}
}

// Description returns the todo's description
func (t Todo) Description() string {
	return t.description
}

// Priority returns the todo's priority
func (t Todo) Priority() Priority {
	return t.priority
}

// IsCompleted returns whether the todo is completed
func (t Todo) IsCompleted() bool {
	return t.completed
}

// CompletionDate returns the todo's completion date (nil if not completed or no date)
func (t Todo) CompletionDate() *time.Time {
	return t.completionDate
}

// CreationDate returns the todo's creation date (nil if no date recorded)
func (t Todo) CreationDate() *time.Time {
	return t.creationDate
}

// Projects returns the todo's project tags
func (t Todo) Projects() []string {
	return t.projects
}

// Contexts returns the todo's context tags
func (t Todo) Contexts() []string {
	return t.contexts
}

// ToggleCompletion returns a new Todo with the completion status toggled
// When marking complete: sets completion date to the provided time
// When marking incomplete: clears completion date
// The now parameter allows deterministic testing and follows dependency inversion
func (t Todo) ToggleCompletion(now time.Time) Todo {
	newCompleted := !t.completed
	var newCompletionDate *time.Time

	if newCompleted {
		// Marking as complete: set date to provided time
		newCompletionDate = &now
	} else {
		// Marking as incomplete: clear date
		newCompletionDate = nil
	}

	return Todo{
		description:    t.description,
		priority:       t.priority,
		completed:      newCompleted,
		completionDate: newCompletionDate,
		creationDate:   t.creationDate,
		projects:       t.projects,
		contexts:       t.contexts,
	}
}

// ChangePriority returns a new Todo with the specified priority
func (t Todo) ChangePriority(newPriority Priority) Todo {
	return Todo{
		description:    t.description,
		priority:       newPriority,
		completed:      t.completed,
		completionDate: t.completionDate,
		creationDate:   t.creationDate,
		projects:       t.projects,
		contexts:       t.contexts,
	}
}

// String converts the Priority to its string representation
func (p Priority) String() string {
	switch p {
	case PriorityA:
		return "A"
	case PriorityB:
		return "B"
	case PriorityC:
		return "C"
	case PriorityD:
		return "D"
	default:
		return ""
	}
}

// String converts a Todo to todo.txt format
// This is the inverse operation of parser.Parse()
// Format: x COMP_DATE CREATION_DATE (PRIORITY) Description +project @context
// Or: (PRIORITY) CREATION_DATE Description +project @context
func (t Todo) String() string {
	var result string

	// Add completion marker and date if completed
	if t.completed {
		result = "x "
		if t.completionDate != nil {
			result += t.completionDate.Format("2006-01-02") + " "
		}
		// Add creation date after completion date for completed todos
		if t.creationDate != nil {
			result += t.creationDate.Format("2006-01-02") + " "
		}
	}

	// Add priority if present
	if t.priority != PriorityNone {
		result += "(" + t.priority.String() + ") "
	}

	// Add creation date after priority for active todos (if not already added)
	if !t.completed {
		if t.creationDate != nil {
			result += t.creationDate.Format("2006-01-02") + " "
		}
	}

	// Add description
	result += t.description

	// Add project tags
	for _, project := range t.projects {
		result += " +" + project
	}

	// Add context tags
	for _, context := range t.contexts {
		result += " @" + context
	}

	result += "\n"
	return result
}
