package usecases_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/usecases"
)

type StubTodoSource struct {
	data string
	err  error
}

func (s StubTodoSource) GetTodos() (io.ReadCloser, error) {
	if s.err != nil {
		return nil, s.err
	}
	return io.NopCloser(strings.NewReader(s.data)), nil
}

func TestLoadMatrix(t *testing.T) {
	t.Run("loads and parses todos from source", func(t *testing.T) {
		source := StubTodoSource{
			data: `(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace
`,
		}

		m, err := usecases.LoadMatrix(source)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify todos were loaded and categorized
		if len(m.DoFirst()) != 1 {
			t.Errorf("expected 1 todo in DoFirst, got %d", len(m.DoFirst()))
		}
		if len(m.Schedule()) != 1 {
			t.Errorf("expected 1 todo in Schedule, got %d", len(m.Schedule()))
		}
		if len(m.Delegate()) != 1 {
			t.Errorf("expected 1 todo in Delegate, got %d", len(m.Delegate()))
		}
		if len(m.Eliminate()) != 1 {
			t.Errorf("expected 1 todo in Eliminate, got %d", len(m.Eliminate()))
		}
	})

	t.Run("returns error when source fails", func(t *testing.T) {
		source := StubTodoSource{
			err: errors.New("source error"),
		}

		_, err := usecases.LoadMatrix(source)

		if err == nil {
			t.Error("expected error when source fails, got nil")
		}
	})

	t.Run("handles empty source", func(t *testing.T) {
		source := StubTodoSource{
			data: "",
		}

		m, err := usecases.LoadMatrix(source)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(m.DoFirst()) != 0 || len(m.Schedule()) != 0 || len(m.Delegate()) != 0 || len(m.Eliminate()) != 0 {
			t.Error("expected all quadrants to be empty")
		}
	})
}
