package file_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/quii/todo-eisenhower/adapters/file"
)

func TestSource(t *testing.T) {
	t.Run("reads todos from file", func(t *testing.T) {
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "todo.txt")
		content := "(A) Test todo"

		err := os.WriteFile(tmpFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		source := file.NewSource(tmpFile)
		reader, err := source.GetTodos()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		defer reader.Close()

		got, err := io.ReadAll(reader)
		if err != nil {
			t.Fatalf("failed to read: %v", err)
		}

		if string(got) != content {
			t.Errorf("got %q, want %q", string(got), content)
		}
	})

	t.Run("returns error for non-existent file", func(t *testing.T) {
		source := file.NewSource("/non/existent/file.txt")

		_, err := source.GetTodos()

		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})
}
