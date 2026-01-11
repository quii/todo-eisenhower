package matrix_test

import (
	"testing"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestMatrix(t *testing.T) {
	t.Run("categorizes priority A todos into DoFirst quadrant", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Review security audit", todo.PriorityA),
		}

		m := matrix.New(todos)

		doFirst := m.DoFirst()
		if len(doFirst) != 2 {
			t.Fatalf("expected 2 todos in DoFirst, got %d", len(doFirst))
		}

		assertContainsTodo(t, doFirst, "Fix critical bug")
		assertContainsTodo(t, doFirst, "Review security audit")
	})

	t.Run("categorizes priority B todos into Schedule quadrant", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Plan Q2 roadmap", todo.PriorityB),
		}

		m := matrix.New(todos)

		schedule := m.Schedule()
		if len(schedule) != 1 {
			t.Fatalf("expected 1 todo in Schedule, got %d", len(schedule))
		}

		assertContainsTodo(t, schedule, "Plan Q2 roadmap")
	})

	t.Run("categorizes priority C todos into Delegate quadrant", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Respond to routine emails", todo.PriorityC),
		}

		m := matrix.New(todos)

		delegate := m.Delegate()
		if len(delegate) != 1 {
			t.Fatalf("expected 1 todo in Delegate, got %d", len(delegate))
		}

		assertContainsTodo(t, delegate, "Respond to routine emails")
	})

	t.Run("categorizes priority D and no priority todos into Eliminate quadrant", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Organize desk", todo.PriorityD),
			todo.New("Review old docs", todo.PriorityNone),
		}

		m := matrix.New(todos)

		eliminate := m.Eliminate()
		if len(eliminate) != 2 {
			t.Fatalf("expected 2 todos in Eliminate, got %d", len(eliminate))
		}

		assertContainsTodo(t, eliminate, "Organize desk")
		assertContainsTodo(t, eliminate, "Review old docs")
	})

	t.Run("distributes todos across all quadrants", func(t *testing.T) {
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Plan quarterly goals", todo.PriorityB),
			todo.New("Reply to emails", todo.PriorityC),
			todo.New("Clean workspace", todo.PriorityD),
		}

		m := matrix.New(todos)

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
}

func assertContainsTodo(t *testing.T, todos []todo.Todo, description string) {
	t.Helper()
	for _, td := range todos {
		if td.Description() == description {
			return
		}
	}
	t.Errorf("expected to find todo with description %q", description)
}
