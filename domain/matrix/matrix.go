// Package matrix provides the Eisenhower Matrix domain model for organizing todos into quadrants.
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

// AddTodo adds a todo to the appropriate quadrant based on its priority
func (m Matrix) AddTodo(t todo.Todo) Matrix {
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
	return m
}

// QuadrantType identifies which quadrant a todo belongs to
type QuadrantType int

const (
	DoFirstQuadrant QuadrantType = iota
	ScheduleQuadrant
	DelegateQuadrant
	EliminateQuadrant
)

// UpdateTodoAtIndex updates the todo at the given index in the specified quadrant
func (m Matrix) UpdateTodoAtIndex(quadrant QuadrantType, index int, newTodo todo.Todo) Matrix {
	switch quadrant {
	case DoFirstQuadrant:
		if index >= 0 && index < len(m.doFirst) {
			m.doFirst[index] = newTodo
		}
	case ScheduleQuadrant:
		if index >= 0 && index < len(m.schedule) {
			m.schedule[index] = newTodo
		}
	case DelegateQuadrant:
		if index >= 0 && index < len(m.delegate) {
			m.delegate[index] = newTodo
		}
	case EliminateQuadrant:
		if index >= 0 && index < len(m.eliminate) {
			m.eliminate[index] = newTodo
		}
	}
	return m
}

// RemoveTodo removes a todo from the matrix by comparing descriptions
// Returns a new Matrix without the specified todo
func (m Matrix) RemoveTodo(todoToRemove todo.Todo) Matrix {
	m.doFirst = removeFromSlice(m.doFirst, todoToRemove)
	m.schedule = removeFromSlice(m.schedule, todoToRemove)
	m.delegate = removeFromSlice(m.delegate, todoToRemove)
	m.eliminate = removeFromSlice(m.eliminate, todoToRemove)
	return m
}

// removeFromSlice removes todos matching the given todo from a slice
func removeFromSlice(todos []todo.Todo, todoToRemove todo.Todo) []todo.Todo {
	result := make([]todo.Todo, 0, len(todos))
	for _, t := range todos {
		// Compare by description since Todo doesn't have an ID
		if t.Description() != todoToRemove.Description() {
			result = append(result, t)
		}
	}
	return result
}

// AllTodos returns all todos from all quadrants
func (m Matrix) AllTodos() []todo.Todo {
	all := make([]todo.Todo, 0)
	all = append(all, m.doFirst...)
	all = append(all, m.schedule...)
	all = append(all, m.delegate...)
	all = append(all, m.eliminate...)
	return all
}
