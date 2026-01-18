// Package matrix provides the Eisenhower Matrix domain model for organizing todos into quadrants.
package matrix

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
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

// ToggleCompletionAt toggles the completion status of a todo at the specified position.
// Returns the updated matrix and true if successful, or the original matrix and false if invalid.
func (m Matrix) ToggleCompletionAt(quadrant QuadrantType, index int) (Matrix, bool) {
	todos := m.getTodosForQuadrant(quadrant)

	// Validate index
	if index < 0 || index >= len(todos) {
		return m, false // No-op if invalid
	}

	// Toggle completion on the todo
	selectedTodo := todos[index]
	updatedTodo := selectedTodo.ToggleCompletion()

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

// Inventory represents system health and WIP metrics for the matrix
type Inventory struct {
	// Per-quadrant metrics
	DoFirstActive      int
	DoFirstOldestDays  int
	ScheduleActive     int
	ScheduleOldestDays int
	DelegateActive     int
	DelegateOldestDays int
	EliminateActive    int
	EliminateOldestDays int

	TotalActive int

	// Throughput metrics
	CompletedLast7Days int
	AddedLast7Days     int

	// Tag breakdowns
	ContextBreakdown map[string]TagMetrics
	ProjectBreakdown map[string]TagMetrics
}

// TagMetrics represents metrics for a specific tag (context or project)
type TagMetrics struct {
	Tag        string
	Count      int
	AvgAgeDays int
}

// CalculateInventory analyzes the matrix and returns WIP/inventory metrics.
// The now parameter allows deterministic testing and follows dependency inversion.
func (m Matrix) CalculateInventory(now time.Time) Inventory {
	sevenDaysAgo := now.AddDate(0, 0, -7)

	inventory := Inventory{
		ContextBreakdown: make(map[string]TagMetrics),
		ProjectBreakdown: make(map[string]TagMetrics),
	}

	// Track tag ages for averaging (only for todos with creation dates)
	contextAges := make(map[string][]int)
	projectAges := make(map[string][]int)

	// Process each quadrant
	quadrants := []struct {
		todos         []todo.Todo
		activeCount   *int
		oldestDays    *int
	}{
		{m.doFirst, &inventory.DoFirstActive, &inventory.DoFirstOldestDays},
		{m.schedule, &inventory.ScheduleActive, &inventory.ScheduleOldestDays},
		{m.delegate, &inventory.DelegateActive, &inventory.DelegateOldestDays},
		{m.eliminate, &inventory.EliminateActive, &inventory.EliminateOldestDays},
	}

	for _, q := range quadrants {
		for _, t := range q.todos {
			// Only process age if creation date exists
			creationDate := t.CreationDate()
			var ageDays int
			hasCreationDate := false
			if creationDate != nil {
				ageDays = int(now.Sub(*creationDate).Hours() / 24)
				hasCreationDate = true
			}

			// Count active todos and track oldest per quadrant (only with creation dates)
			if !t.IsCompleted() {
				inventory.TotalActive++
				*q.activeCount++

				if hasCreationDate && ageDays > *q.oldestDays {
					*q.oldestDays = ageDays
				}

				// Track contexts for active todos (only with creation dates and non-empty)
				if hasCreationDate {
					for _, context := range t.Contexts() {
						if context != "" {
							if _, exists := contextAges[context]; !exists {
								contextAges[context] = []int{}
							}
							contextAges[context] = append(contextAges[context], ageDays)
						}
					}

					// Track projects for active todos (only with creation dates and non-empty)
					for _, project := range t.Projects() {
						if project != "" {
							if _, exists := projectAges[project]; !exists {
								projectAges[project] = []int{}
							}
							projectAges[project] = append(projectAges[project], ageDays)
						}
					}
				}
			}

			// Count completed in last 7 days
			if t.IsCompleted() {
				if completionDate := t.CompletionDate(); completionDate != nil {
					if completionDate.After(sevenDaysAgo) {
						inventory.CompletedLast7Days++
					}
				}
			}

			// Count added in last 7 days
			if creationDate != nil {
				if creationDate.After(sevenDaysAgo) {
					inventory.AddedLast7Days++
				}
			}
		}
	}

	// Calculate context breakdown with averages
	for context, ages := range contextAges {
		count := len(ages)
		sum := 0
		for _, age := range ages {
			sum += age
		}
		avgAge := 0
		if count > 0 {
			avgAge = sum / count
		}

		inventory.ContextBreakdown[context] = TagMetrics{
			Tag:        context,
			Count:      count,
			AvgAgeDays: avgAge,
		}
	}

	// Calculate project breakdown with averages
	for project, ages := range projectAges {
		count := len(ages)
		sum := 0
		for _, age := range ages {
			sum += age
		}
		avgAge := 0
		if count > 0 {
			avgAge = sum / count
		}

		inventory.ProjectBreakdown[project] = TagMetrics{
			Tag:        project,
			Count:      count,
			AvgAgeDays: avgAge,
		}
	}

	return inventory
}
