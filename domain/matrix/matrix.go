// Package matrix provides the Eisenhower Matrix domain model for organizing todos into quadrants.
package matrix

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/domain/todotxt"
)

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
	todos := m.getTodosForQuadrant(quadrant)

	if index < 0 || index >= len(todos) {
		return m
	}

	todos[index] = newTodo
	return m.setTodosForQuadrant(quadrant, todos)
}

// EditTodo edits the todo at the given index in the specified quadrant with a new description
func (m Matrix) EditTodo(quadrant QuadrantType, index int, newDescription string) Matrix {
	todos := m.getTodosForQuadrant(quadrant)

	if index < 0 || index >= len(todos) {
		return m
	}

	originalTodo := todos[index]
	priority := originalTodo.Priority()
	updatedTodo := todotxt.ParseEdit(originalTodo, newDescription, priority)

	return m.UpdateTodoAtIndex(quadrant, index, updatedTodo)
}

// RemoveTodo removes a todo from the matrix by comparing descriptions
// Returns a new Matrix without the specified todo
func (m Matrix) RemoveTodo(todoToRemove todo.Todo) Matrix {
	quadrants := []QuadrantType{DoFirstQuadrant, ScheduleQuadrant, DelegateQuadrant, EliminateQuadrant}
	for _, q := range quadrants {
		todos := m.getTodosForQuadrant(q)
		m = m.setTodosForQuadrant(q, removeFromSlice(todos, todoToRemove))
	}
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
	quadrants := []QuadrantType{DoFirstQuadrant, ScheduleQuadrant, DelegateQuadrant, EliminateQuadrant}
	for _, q := range quadrants {
		all = append(all, m.getTodosForQuadrant(q)...)
	}
	return all
}

// FilterByTag returns a new Matrix containing only todos that match the given tag filter.
// Filter format: "+project" for projects, "@context" for contexts.
// Returns an empty matrix if the filter doesn't match any todos.
func (m Matrix) FilterByTag(filter string) Matrix {
	if filter == "" {
		return m
	}

	return Matrix{
		doFirst:   filterTodosByTag(m.doFirst, filter),
		schedule:  filterTodosByTag(m.schedule, filter),
		delegate:  filterTodosByTag(m.delegate, filter),
		eliminate: filterTodosByTag(m.eliminate, filter),
	}
}

// filterTodosByTag filters a slice of todos by project or context tag
func filterTodosByTag(todos []todo.Todo, filter string) []todo.Todo {
	if filter == "" {
		return todos
	}

	if len(filter) < 2 {
		return []todo.Todo{}
	}

	prefix := filter[0]
	tag := filter[1:]

	filtered := make([]todo.Todo, 0)
	for _, t := range todos {
		match := false

		if prefix == '+' {
			// Filter by project
			for _, p := range t.Projects() {
				if equalsFold(p, tag) {
					match = true
					break
				}
			}
		} else if prefix == '@' {
			// Filter by context
			for _, c := range t.Contexts() {
				if equalsFold(c, tag) {
					match = true
					break
				}
			}
		}

		if match {
			filtered = append(filtered, t)
		}
	}
	return filtered
}

// equalsFold is a simple case-insensitive string comparison
func equalsFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		// Convert to lowercase if uppercase letter
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}

// ToggleCompletionAt toggles the completion status of a todo at the specified position.
// Returns the updated matrix and true if successful, or the original matrix and false if invalid.
// The now parameter allows deterministic testing and follows dependency inversion.
func (m Matrix) ToggleCompletionAt(quadrant QuadrantType, index int, now time.Time) (Matrix, bool) {
	todos := m.getTodosForQuadrant(quadrant)

	// Validate index
	if index < 0 || index >= len(todos) {
		return m, false // No-op if invalid
	}

	// Toggle completion on the todo
	selectedTodo := todos[index]
	updatedTodo := selectedTodo.ToggleCompletion(now)

	// Update in place (completion doesn't change priority/quadrant)
	return m.UpdateTodoAtIndex(quadrant, index, updatedTodo), true
}

// ChangePriorityAt changes the priority of a todo at the specified position.
// Returns the updated matrix and true if successful, or the original matrix and false if invalid/unchanged.
// Note: Changing priority may move the todo to a different quadrant.
func (m Matrix) ChangePriorityAt(quadrant QuadrantType, index int, newPriority todo.Priority) (Matrix, bool) {
	todos := m.getTodosForQuadrant(quadrant)

	// Validate index
	if index < 0 || index >= len(todos) {
		return m, false // No-op if invalid
	}

	// Get current todo and check if priority is already the same
	selectedTodo := todos[index]
	if selectedTodo.Priority() == newPriority {
		return m, false // No-op if priority unchanged
	}

	// Change priority on the todo
	updatedTodo := selectedTodo.ChangePriority(newPriority)

	// Priority change moves todo between quadrants, so remove old and add new
	return m.RemoveTodo(selectedTodo).AddTodo(updatedTodo), true
}

// ArchiveTodoAt archives (removes) a completed todo at the specified position.
// Returns the archived todo, updated matrix, and true if successful.
// Returns the original matrix and false if the index is invalid or the todo is not completed.
func (m Matrix) ArchiveTodoAt(quadrant QuadrantType, index int) (todo.Todo, Matrix, bool) {
	todos := m.getTodosForQuadrant(quadrant)

	// Validate index
	if index < 0 || index >= len(todos) {
		return todo.Todo{}, m, false
	}

	selectedTodo := todos[index]

	// Only archive completed todos
	if !selectedTodo.IsCompleted() {
		return todo.Todo{}, m, false
	}

	// Remove the todo from the matrix
	return selectedTodo, m.RemoveTodo(selectedTodo), true
}

// getTodosForQuadrant is a helper that returns the todos for a given quadrant
func (m Matrix) getTodosForQuadrant(quadrant QuadrantType) []todo.Todo {
	switch quadrant {
	case DoFirstQuadrant:
		return m.doFirst
	case ScheduleQuadrant:
		return m.schedule
	case DelegateQuadrant:
		return m.delegate
	case EliminateQuadrant:
		return m.eliminate
	default:
		return []todo.Todo{}
	}
}

// setTodosForQuadrant is a helper that sets the todos for a given quadrant
func (m Matrix) setTodosForQuadrant(quadrant QuadrantType, todos []todo.Todo) Matrix {
	switch quadrant {
	case DoFirstQuadrant:
		m.doFirst = todos
	case ScheduleQuadrant:
		m.schedule = todos
	case DelegateQuadrant:
		m.delegate = todos
	case EliminateQuadrant:
		m.eliminate = todos
	}
	return m
}
