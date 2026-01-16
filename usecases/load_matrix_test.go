package usecases_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
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
		is := is.New(t)
		source := StubTodoSource{
			data: `(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace
`,
		}

		m, err := usecases.LoadMatrix(source)

		is.NoErr(err)

		// Verify todos were loaded and categorized
		is.Equal(len(m.DoFirst()), 1)  // expected 1 todo in DoFirst
		is.Equal(len(m.Schedule()), 1) // expected 1 todo in Schedule
		is.Equal(len(m.Delegate()), 1) // expected 1 todo in Delegate
		is.Equal(len(m.Eliminate()), 1) // expected 1 todo in Eliminate
	})

	t.Run("returns error when source fails", func(t *testing.T) {
		is := is.New(t)
		source := StubTodoSource{
			err: errors.New("source error"),
		}

		_, err := usecases.LoadMatrix(source)

		is.True(err != nil) // expected error when source fails
	})

	t.Run("handles empty source", func(t *testing.T) {
		is := is.New(t)
		source := StubTodoSource{
			data: "",
		}

		m, err := usecases.LoadMatrix(source)

		is.NoErr(err)

		is.Equal(len(m.DoFirst()), 0)  // expected empty DoFirst
		is.Equal(len(m.Schedule()), 0) // expected empty Schedule
		is.Equal(len(m.Delegate()), 0) // expected empty Delegate
		is.Equal(len(m.Eliminate()), 0) // expected empty Eliminate
	})
}
