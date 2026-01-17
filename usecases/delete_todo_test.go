package usecases_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestDeleteTodo(t *testing.T) {
	//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
	is := is.New(t)

	todo1 := todo.New("Task to keep", todo.PriorityA)
	todo2 := todo.New("Task to delete", todo.PriorityA)
	m := matrix.New([]todo.Todo{todo1, todo2})

	// Create a spy writer to capture what gets written
	writer := &SpyTodoWriter{}

	// Delete todo2
	updatedMatrix, err := usecases.DeleteTodo(writer, m, todo2)

	is.NoErr(err)
	is.Equal(len(updatedMatrix.AllTodos()), 1)
	is.Equal(updatedMatrix.AllTodos()[0].Description(), "Task to keep")

	// Writer should have been called to persist all todos
	is.True(writer.replaceAllCalled)
}

// SpyTodoWriter captures calls for testing
type SpyTodoWriter struct {
	replaceAllCalled bool
	lastContent      string
}

func (s *SpyTodoWriter) SaveTodo(line string) error {
	return nil
}

func (s *SpyTodoWriter) ReplaceAll(content string) error {
	s.replaceAllCalled = true
	s.lastContent = content
	return nil
}
