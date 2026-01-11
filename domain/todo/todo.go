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
}

// New creates a new Todo with the given description and priority
func New(description string, priority Priority) Todo {
	return Todo{
		description: description,
		priority:    priority,
		completed:   false,
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
