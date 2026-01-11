package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// LoadMatrix loads todos and returns an Eisenhower matrix
// Story 001: Returns hard-coded todos
// Story 002: Will read from hard-coded file path
// Story 003: Will accept file path parameter
func LoadMatrix() matrix.Matrix {
	// Hard-coded sample todos for Story 001
	todos := []todo.Todo{
		todo.New("Fix critical production bug", todo.PriorityA),
		todo.New("Review security audit findings", todo.PriorityA),
		todo.New("Plan Q2 roadmap", todo.PriorityB),
		todo.New("Research new framework options", todo.PriorityB),
		todo.New("Respond to routine emails", todo.PriorityC),
		todo.New("Attend weekly status meeting", todo.PriorityC),
		todo.New("Organize old project files", todo.PriorityD),
		todo.New("Update personal wiki", todo.PriorityD),
	}

	return matrix.New(todos)
}
