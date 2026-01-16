package file_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/file"
)

func TestSource(t *testing.T) {
	t.Run("reads todos from file", func(t *testing.T) {
		is := is.New(t)
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")
		content := "(A) Test todo"

		err := os.WriteFile(tmpFile, []byte(content), 0644)
		is.NoErr(err) // failed to create test file

		source := file.NewSource(tmpFile)
		reader, err := source.GetTodos()

		is.NoErr(err)
		defer reader.Close()

		got, err := io.ReadAll(reader)
		is.NoErr(err) // failed to read

		is.Equal(string(got), content)
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		is := is.New(t)
		source := file.NewSource("/non/existent/file.txt")

		_, err := source.GetTodos()

		is.True(err != nil) // expected error for non-existent file
	})
}
