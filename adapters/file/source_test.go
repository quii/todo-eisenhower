package file_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/file"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestRepository(t *testing.T) {
	t.Run("reads todos from file", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")
		content := "(A) Test todo\n"

		//nolint:gosec // G306: test file permissions intentionally match production (0o644)
		err := os.WriteFile(tmpFile, []byte(content), 0o644)
		is.NoErr(err) // failed to create test file

		repo := file.NewRepository(tmpFile)
		todos, err := repo.LoadAll()

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Test todo")
		is.Equal(todos[0].Priority(), todo.PriorityA)
	})

	t.Run("creates empty file if it doesn't exist", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "new-todo.txt")

		repo := file.NewRepository(tmpFile)
		todos, err := repo.LoadAll()

		is.NoErr(err)
		is.Equal(len(todos), 0)

		// Verify file was actually created
		_, err = os.Stat(tmpFile)
		is.NoErr(err) // file should exist
	})

	t.Run("saves todos to file", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")

		repo := file.NewRepository(tmpFile)

		// Save some todos
		todosToSave := []todo.Todo{
			todo.New("First task", todo.PriorityA),
			todo.New("Second task", todo.PriorityB),
		}
		err := repo.SaveAll(todosToSave)
		is.NoErr(err)

		// Load them back
		loaded, err := repo.LoadAll()
		is.NoErr(err)
		is.Equal(len(loaded), 2)
		is.Equal(loaded[0].Description(), "First task")
		is.Equal(loaded[1].Description(), "Second task")
	})
}
