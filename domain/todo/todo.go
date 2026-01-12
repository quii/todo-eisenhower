package todo

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
	description string
	priority    Priority
	completed   bool
	projects    []string
	contexts    []string
}

// New creates a new Todo with the given description and priority
func New(description string, priority Priority) Todo {
	return Todo{
		description: description,
		priority:    priority,
		completed:   false,
		projects:    []string{},
		contexts:    []string{},
	}
}

// NewCompleted creates a new completed Todo with the given description and priority
func NewCompleted(description string, priority Priority) Todo {
	return Todo{
		description: description,
		priority:    priority,
		completed:   true,
		projects:    []string{},
		contexts:    []string{},
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
		description: description,
		priority:    priority,
		completed:   false,
		projects:    projects,
		contexts:    contexts,
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

// Projects returns the todo's project tags
func (t Todo) Projects() []string {
	return t.projects
}

// Contexts returns the todo's context tags
func (t Todo) Contexts() []string {
	return t.contexts
}
