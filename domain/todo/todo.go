// Package todo provides the Todo domain model following the todo.txt specification.
package todo

import (
	"time"
)

// Priority represents the priority level of a todo item
type Priority int

const (
	PriorityNone Priority = iota
	PriorityA
	PriorityB
	PriorityC
	PriorityD
	PriorityE
)

// Todo represents a single todo item
type Todo struct {
	description     string
	priority        Priority
	completed       bool
	completionDate  *time.Time // nil if not completed or no date recorded
	creationDate    *time.Time // nil if no creation date recorded
	dueDate         *time.Time // nil if no due date recorded
	prioritisedDate *time.Time // nil if not Priority A or no prioritised date recorded
	projects        []string
	contexts        []string
}

// NewFull is a comprehensive constructor used by the todotxt parser to create todos with all fields
// This allows the parser to set all fields without needing combinatorial constructors
func NewFull(description string, priority Priority, completed bool, completionDate, creationDate, dueDate, prioritisedDate *time.Time, projects, contexts []string) Todo {
	if projects == nil {
		projects = []string{}
	}
	if contexts == nil {
		contexts = []string{}
	}
	return Todo{
		description:     description,
		priority:        priority,
		completed:       completed,
		completionDate:  completionDate,
		creationDate:    creationDate,
		dueDate:         dueDate,
		prioritisedDate: prioritisedDate,
		projects:        projects,
		contexts:        contexts,
	}
}

// New creates a new Todo with the given description and priority
func New(description string, priority Priority) Todo {
	return Todo{
		description:     description,
		priority:        priority,
		completed:       false,
		completionDate:  nil,
		creationDate:    nil,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        []string{},
		contexts:        []string{},
	}
}

// NewWithCreationDate creates a new Todo with creation date
func NewWithCreationDate(description string, priority Priority, creationDate *time.Time) Todo {
	return Todo{
		description:     description,
		priority:        priority,
		completed:       false,
		completionDate:  nil,
		creationDate:    creationDate,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        []string{},
		contexts:        []string{},
	}
}

// NewCompleted creates a new completed Todo with the given description, priority, and optional completion date
func NewCompleted(description string, priority Priority, completionDate *time.Time) Todo {
	return Todo{
		description:     description,
		priority:        priority,
		completed:       true,
		completionDate:  completionDate,
		creationDate:    nil,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        []string{},
		contexts:        []string{},
	}
}

// NewCompletedWithDates creates a new completed Todo with both completion and creation dates
func NewCompletedWithDates(description string, priority Priority, completionDate, creationDate *time.Time) Todo {
	return Todo{
		description:     description,
		priority:        priority,
		completed:       true,
		completionDate:  completionDate,
		creationDate:    creationDate,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        []string{},
		contexts:        []string{},
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
		description:     description,
		priority:        priority,
		completed:       false,
		completionDate:  nil,
		creationDate:    nil,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        projects,
		contexts:        contexts,
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
		description:     description,
		priority:        priority,
		completed:       false,
		completionDate:  nil,
		creationDate:    creationDate,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        projects,
		contexts:        contexts,
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
		description:     description,
		priority:        priority,
		completed:       true,
		completionDate:  completionDate,
		creationDate:    nil,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        projects,
		contexts:        contexts,
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
		description:     description,
		priority:        priority,
		completed:       true,
		completionDate:  completionDate,
		creationDate:    creationDate,
		dueDate:         nil,
		prioritisedDate: nil,
		projects:        projects,
		contexts:        contexts,
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

// DueDate returns the todo's due date (nil if no due date recorded)
func (t Todo) DueDate() *time.Time {
	return t.dueDate
}

// PrioritisedDate returns the todo's prioritised date (nil if not Priority A or no date recorded)
func (t Todo) PrioritisedDate() *time.Time {
	return t.prioritisedDate
}

// Projects returns the todo's project tags
func (t Todo) Projects() []string {
	return t.projects
}

// Contexts returns the todo's context tags
func (t Todo) Contexts() []string {
	return t.contexts
}

// IsStale returns true if the todo has been sitting in its quadrant for too long
// Completed tasks are never stale
// Priority A: stale after 2 business days from prioritisedDate
// Priority B/C/D: stale after 5 business days from creationDate
func (t Todo) IsStale(now time.Time) bool {
	// Completed tasks are never stale
	if t.completed {
		return false
	}

	switch t.priority {
	case PriorityA:
		// Do First: check prioritised date
		if t.prioritisedDate == nil {
			return false // No prioritised date, can't be stale
		}
		businessDays := businessDaysBetween(*t.prioritisedDate, now)
		return businessDays > 2

	case PriorityB, PriorityC, PriorityD:
		// Schedule/Delegate/Eliminate: check creation date
		if t.creationDate == nil {
			return false // No creation date, can't be stale
		}
		businessDays := businessDaysBetween(*t.creationDate, now)
		return businessDays > 5

	default:
		// No priority or invalid priority
		return false
	}
}

// businessDaysBetween calculates the number of business days (Monday-Friday) between two dates
// Excludes weekends (Saturday and Sunday)
// Returns 0 if from and to are the same day
func businessDaysBetween(from, to time.Time) int {
	// Normalize to start of day for accurate day counting
	from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())

	// If same day or to is before from, return 0
	if !to.After(from) {
		return 0
	}

	businessDays := 0
	current := from

	for current.Before(to) {
		current = current.AddDate(0, 0, 1)
		// Count only weekdays (Monday = 1, Friday = 5)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			businessDays++
		}
	}

	return businessDays
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
		description:     t.description,
		priority:        t.priority,
		completed:       newCompleted,
		completionDate:  newCompletionDate,
		creationDate:    t.creationDate,
		dueDate:         t.dueDate,
		prioritisedDate: t.prioritisedDate,
		projects:        t.projects,
		contexts:        t.contexts,
	}
}

// ChangePriority returns a new Todo with the specified priority
func (t Todo) ChangePriority(newPriority Priority) Todo {
	return Todo{
		description:     t.description,
		priority:        newPriority,
		completed:       t.completed,
		completionDate:  t.completionDate,
		creationDate:    t.creationDate,
		dueDate:         t.dueDate,
		prioritisedDate: t.prioritisedDate,
		projects:        t.projects,
		contexts:        t.contexts,
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
	case PriorityE:
		return "E"
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

	// Add due date
	if t.dueDate != nil {
		result += " due:" + t.dueDate.Format("2006-01-02")
	}

	// Add prioritised date
	if t.prioritisedDate != nil {
		result += " prioritised:" + t.prioritisedDate.Format("2006-01-02")
	}

	result += "\n"
	return result
}
