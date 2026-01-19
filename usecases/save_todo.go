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

// saveTodo is a helper that loads all, adds one, and saves all
func saveTodo(repo TodoRepository, t todo.Todo) error {
	todos, err := repo.LoadAll()
	if err != nil {
		return err
	}

	m := matrix.New(todos)
	updated := m.AddTodo(t)

	return repo.SaveAll(updated.AllTodos())
}

// saveAllTodos is a helper that saves all todos from a matrix
func saveAllTodos(repo TodoRepository, m matrix.Matrix) error {
	return repo.SaveAll(m.AllTodos())
}
