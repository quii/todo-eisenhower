package matrix

import "github.com/quii/todo-eisenhower/domain/todo"

// Matrix represents an Eisenhower matrix organizing todos by quadrant
type Matrix struct {
	doFirst   []todo.Todo
	schedule  []todo.Todo
	delegate  []todo.Todo
	eliminate []todo.Todo
}

// New creates a new Matrix and categorizes the given todos into quadrants
func New(todos []todo.Todo) Matrix {
	m := Matrix{
		doFirst:   make([]todo.Todo, 0),
		schedule:  make([]todo.Todo, 0),
		delegate:  make([]todo.Todo, 0),
		eliminate: make([]todo.Todo, 0),
	}

	for _, t := range todos {
		switch t.Priority() {
		case todo.PriorityA:
			m.doFirst = append(m.doFirst, t)
		case todo.PriorityB:
			m.schedule = append(m.schedule, t)
		case todo.PriorityC:
			m.delegate = append(m.delegate, t)
		case todo.PriorityD, todo.PriorityNone:
			m.eliminate = append(m.eliminate, t)
		}
	}

	return m
}

// DoFirst returns todos in the "Do First" quadrant (urgent and important)
func (m Matrix) DoFirst() []todo.Todo {
	return m.doFirst
}

// Schedule returns todos in the "Schedule" quadrant (important, not urgent)
func (m Matrix) Schedule() []todo.Todo {
	return m.schedule
}

// Delegate returns todos in the "Delegate" quadrant (urgent, not important)
func (m Matrix) Delegate() []todo.Todo {
	return m.delegate
}

// Eliminate returns todos in the "Eliminate" quadrant (neither urgent nor important)
func (m Matrix) Eliminate() []todo.Todo {
	return m.eliminate
}
