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
		projects:       projects,
		contexts:       contexts,
	}
}

// NewCompletedWithTags creates a new completed Todo with tags and optional completion date
func NewCompletedWithTags(description string, priority Priority, completionDate *time.Time, projects, contexts []string) Todo {
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

// Projects returns the todo's project tags
func (t Todo) Projects() []string {
	return t.projects
}

// Contexts returns the todo's context tags
func (t Todo) Contexts() []string {
	return t.contexts
}

// ToggleCompletion returns a new Todo with the completion status toggled
// When marking complete: sets completion date to now
// When marking incomplete: clears completion date
func (t Todo) ToggleCompletion() Todo {
	newCompleted := !t.completed
	var newCompletionDate *time.Time

	if newCompleted {
		// Marking as complete: set date to now
		now := time.Now()
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
		projects:       t.projects,
		contexts:       t.contexts,
	}
}
