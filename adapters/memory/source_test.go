package memory_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

func TestMemoryRepository_RoundTrip(t *testing.T) {
	is := is.New(t)

	// Start with some initial content
	initialContent := `(A) First task
(B) Second task
`
	repo := memory.NewRepository(initialContent)

	// Load using repository
	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	// Verify we got the todos
	is.Equal(len(m.DoFirst()), 1)
	is.Equal(len(m.Schedule()), 1)
	is.Equal(m.DoFirst()[0].Description(), "First task")
	is.Equal(m.Schedule()[0].Description(), "Second task")

	// Add a new todo by loading, modifying, saving
	todos, _ := repo.LoadAll()
	todos = append(todos, todo.New("Third task", todo.PriorityC))
	err = repo.SaveAll(todos)
	is.NoErr(err)

	// Read back and verify
	m, err = usecases.LoadMatrix(repo)
	is.NoErr(err)
	is.Equal(len(m.AllTodos()), 3)

	// Replace all todos
	replacementTodos := []todo.Todo{
		todo.New("Only task", todo.PriorityA),
	}
	err = repo.SaveAll(replacementTodos)
	is.NoErr(err)

	// Verify replacement worked
	m, err = usecases.LoadMatrix(repo)
	is.NoErr(err)
	is.Equal(len(m.AllTodos()), 1)
	is.Equal(m.AllTodos()[0].Description(), "Only task")

	// Verify the buffer content is valid todo.txt format
	expected := "(A) Only task\n"
	is.Equal(repo.String(), expected)
}

func TestMemoryRepository_EmptySource(t *testing.T) {
	is := is.New(t)

	repo := memory.NewRepository("")

	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)
	is.Equal(len(m.AllTodos()), 0)
}

func TestMemoryRepository_UsesRealMarshal(t *testing.T) {
	// This test proves we're using the real todotxt.Marshal,
	// not duplicating serialization logic
	is := is.New(t)

	repo := memory.NewRepository("")

	// Create a todo with projects and contexts
	todoWithTags := todo.New("Deploy feature +WebApp @computer", todo.PriorityA)
	err := repo.SaveAll([]todo.Todo{todoWithTags})
	is.NoErr(err)

	// The real Marshal should preserve the full todo.txt format
	content := repo.String()
	is.Equal(content, "(A) Deploy feature +WebApp @computer\n")

	// And Unmarshal should parse it back correctly
	m, err := usecases.LoadMatrix(repo)
	is.NoErr(err)

	todos := m.DoFirst()
	is.Equal(len(todos), 1)
	is.Equal(len(todos[0].Projects()), 1)
	is.Equal(todos[0].Projects()[0], "WebApp")
	is.Equal(len(todos[0].Contexts()), 1)
	is.Equal(todos[0].Contexts()[0], "computer")
}
