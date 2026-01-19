package usecases_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/file"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// TestRepositoryContract_Memory runs the contract tests against the memory repository
func TestRepositoryContract_Memory(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		repo := memory.NewRepository()
		repositoryContract(t, repo)
	})

	t.Run("repository with initial data", func(t *testing.T) {
		repo := memory.NewRepository()
		initialTodo := todo.New("Existing task", todo.PriorityA)
		err := repo.SaveAll([]todo.Todo{initialTodo})
		if err != nil {
			t.Fatalf("failed to create initial data: %v", err)
		}
		repositoryContractWithInitialData(t, repo)
	})
}

// TestRepositoryContract_File runs the contract tests against the file repository
func TestRepositoryContract_File(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")
		repo := file.NewRepository(tmpFile)
		repositoryContract(t, repo)
	})

	t.Run("repository with initial data", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")

		// Create file with initial data
		repo := file.NewRepository(tmpFile)
		initialTodo := todo.New("Existing task", todo.PriorityA)
		err := repo.SaveAll([]todo.Todo{initialTodo})
		if err != nil {
			t.Fatalf("failed to create initial data: %v", err)
		}

		repositoryContractWithInitialData(t, repo)
	})
}

// repositoryContract tests the contract for an empty repository
func repositoryContract(t *testing.T, repo usecases.TodoRepository) {
	t.Helper()
	is := is.New(t)

	t.Run("LoadAll returns empty list for new repository", func(t *testing.T) {
		todos, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(todos), 0)
	})

	t.Run("SaveAll and LoadAll round-trip", func(t *testing.T) {
		todosToSave := []todo.Todo{
			todo.New("First task", todo.PriorityA),
			todo.New("Second task", todo.PriorityB),
		}

		err := repo.SaveAll(todosToSave)
		is.NoErr(err)

		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), 2)
		is.Equal(loaded[0].Description(), "First task")
		is.Equal(loaded[0].Priority(), todo.PriorityA)
		is.Equal(loaded[1].Description(), "Second task")
		is.Equal(loaded[1].Priority(), todo.PriorityB)
	})

	t.Run("SaveAll preserves all todo properties", func(t *testing.T) {
		// Create a todo with all properties
		now := time.Now()
		creationDate := now.AddDate(0, 0, -5) // 5 days ago
		completionDate := now

		testTodo := todo.NewCompletedWithTagsAndDates(
			"Task with everything",
			todo.PriorityA,
			&completionDate,
			&creationDate,
			[]string{"Project"},
			[]string{"context"},
		)

		err := repo.SaveAll([]todo.Todo{testTodo})
		is.NoErr(err)

		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), 1)

		// Verify all properties preserved
		retrieved := loaded[0]
		is.True(retrieved.Description() == "Task with everything" ||
			retrieved.Description() == "Task with everything +Project @context") // Description format varies
		is.Equal(retrieved.Priority(), todo.PriorityA)
		is.Equal(len(retrieved.Projects()), 1)
		is.Equal(retrieved.Projects()[0], "Project")
		is.Equal(len(retrieved.Contexts()), 1)
		is.Equal(retrieved.Contexts()[0], "context")
		is.True(retrieved.IsCompleted())
		is.True(retrieved.CreationDate() != nil)
		is.Equal(retrieved.CreationDate().Format("2006-01-02"), creationDate.Format("2006-01-02"))
		is.True(retrieved.CompletionDate() != nil)
		is.Equal(retrieved.CompletionDate().Format("2006-01-02"), completionDate.Format("2006-01-02"))
	})

	t.Run("SaveAll with empty list clears repository", func(t *testing.T) {
		// First save some todos
		err := repo.SaveAll([]todo.Todo{
			todo.New("Task to be cleared", todo.PriorityA),
		})
		is.NoErr(err)

		// Then save empty list
		err = repo.SaveAll([]todo.Todo{})
		is.NoErr(err)

		// Should be empty now
		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), 0)
	})

	t.Run("SaveAll replaces all previous content", func(t *testing.T) {
		// Save initial todos
		err := repo.SaveAll([]todo.Todo{
			todo.New("First batch 1", todo.PriorityA),
			todo.New("First batch 2", todo.PriorityB),
		})
		is.NoErr(err)

		// Save different todos (full replacement)
		err = repo.SaveAll([]todo.Todo{
			todo.New("Second batch 1", todo.PriorityC),
		})
		is.NoErr(err)

		// Should only have second batch
		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), 1)
		is.Equal(loaded[0].Description(), "Second batch 1")
		is.Equal(loaded[0].Priority(), todo.PriorityC)
	})
}

// repositoryContractWithInitialData tests the contract when repository has initial data
func repositoryContractWithInitialData(t *testing.T, repo usecases.TodoRepository) {
	t.Helper()
	is := is.New(t)

	t.Run("LoadAll returns existing data", func(t *testing.T) {
		todos, err := repo.LoadAll()
		is.NoErr(err)
		is.True(len(todos) > 0) // Should have at least the initial data
	})

	t.Run("SaveAll can add to existing data", func(t *testing.T) {
		// Load existing
		existing, err := repo.LoadAll()
		is.NoErr(err)
		initialCount := len(existing)

		// Add a new todo
		existing = append(existing, todo.New("New task", todo.PriorityD))
		err = repo.SaveAll(existing)
		is.NoErr(err)

		// Verify count increased
		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), initialCount+1)
	})
}
