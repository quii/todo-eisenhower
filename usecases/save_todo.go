package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// TodoRepository is the port for loading and saving todos
// Implementations handle all the persistence details (files, marshaling, etc)
type TodoRepository interface {
	LoadAll() ([]todo.Todo, error)
	SaveAll(todos []todo.Todo) error
}

func saveAllTodos(repo TodoRepository, m matrix.Matrix) error {
	return repo.SaveAll(m.AllTodos())
}
