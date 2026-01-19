package matrix_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestMatrix(t *testing.T) {
	t.Run("categorizes priority A todos into DoFirst quadrant", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Review security audit", todo.PriorityA),
		}

		m := matrix.New(todos)

		doFirst := m.DoFirst()
		is.Equal(len(doFirst), 2) // expected 2 todos in DoFirst

		assertContainsTodo(is, doFirst, "Fix critical bug")
		assertContainsTodo(is, doFirst, "Review security audit")
	})

	t.Run("categorizes priority B todos into Schedule quadrant", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Plan Q2 roadmap", todo.PriorityB),
		}

		m := matrix.New(todos)

		schedule := m.Schedule()
		is.Equal(len(schedule), 1) // expected 1 todo in Schedule

		assertContainsTodo(is, schedule, "Plan Q2 roadmap")
	})

	t.Run("categorizes priority C todos into Delegate quadrant", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Respond to routine emails", todo.PriorityC),
		}

		m := matrix.New(todos)

		delegate := m.Delegate()
		is.Equal(len(delegate), 1) // expected 1 todo in Delegate

		assertContainsTodo(is, delegate, "Respond to routine emails")
	})

	t.Run("categorizes priority D and no priority todos into Eliminate quadrant", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Organize desk", todo.PriorityD),
			todo.New("Review old docs", todo.PriorityNone),
		}

		m := matrix.New(todos)

		eliminate := m.Eliminate()
		is.Equal(len(eliminate), 2) // expected 2 todos in Eliminate

		assertContainsTodo(is, eliminate, "Organize desk")
		assertContainsTodo(is, eliminate, "Review old docs")
	})

	t.Run("distributes todos across all quadrants", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Fix critical bug", todo.PriorityA),
			todo.New("Plan quarterly goals", todo.PriorityB),
			todo.New("Reply to emails", todo.PriorityC),
			todo.New("Clean workspace", todo.PriorityD),
		}

		m := matrix.New(todos)

		is.Equal(len(m.DoFirst()), 1)  // expected 1 todo in DoFirst
		is.Equal(len(m.Schedule()), 1) // expected 1 todo in Schedule
		is.Equal(len(m.Delegate()), 1) // expected 1 todo in Delegate
		is.Equal(len(m.Eliminate()), 1) // expected 1 todo in Eliminate
	})
}

func assertContainsTodo(is *is.I, todos []todo.Todo, description string) {
	is.Helper()
	for _, td := range todos {
		if td.Description() == description {
			return
		}
	}
	is.Fail() // expected to find todo with description
}

func TestMatrix_RemoveTodo(t *testing.T) {
	is := is.New(t)

	// Create todos
	todo1 := todo.New("First task", todo.PriorityA)
	todo2 := todo.New("Second task", todo.PriorityA)
	todo3 := todo.New("Third task", todo.PriorityB)

	m := matrix.New([]todo.Todo{todo1, todo2, todo3})

	// Remove a todo
	updated := m.RemoveTodo(todo2)

	// Should have 2 todos now
	is.Equal(len(updated.AllTodos()), 2)

	// Should still have todo1 and todo3
	allTodos := updated.AllTodos()
	is.Equal(allTodos[0].Description(), "First task")
	is.Equal(allTodos[1].Description(), "Third task")
}

func TestMatrix_FilterByTag(t *testing.T) {
	t.Run("filters by project tag", func(t *testing.T) {
		is := is.New(t)

		todos := []todo.Todo{
			todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
			todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
			todo.NewWithTags("Task 3", todo.PriorityA, []string{"WebApp"}, []string{}),
		}

		m := matrix.New(todos)
		filtered := m.FilterByTag("+WebApp")

		// Should only have tasks with +WebApp
		is.Equal(len(filtered.DoFirst()), 2)  // Both WebApp tasks are Priority A
		is.Equal(len(filtered.Schedule()), 0) // No Schedule tasks with WebApp
		is.Equal(len(filtered.AllTodos()), 2) // Total filtered
	})

	t.Run("filters by context tag", func(t *testing.T) {
		is := is.New(t)

		todos := []todo.Todo{
			todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
			todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{"phone"}),
			todo.NewWithTags("Task 3", todo.PriorityC, []string{"Backend"}, []string{"computer"}),
		}

		m := matrix.New(todos)
		filtered := m.FilterByTag("@computer")

		// Should only have tasks with @computer
		is.Equal(len(filtered.DoFirst()), 1)   // Task 1
		is.Equal(len(filtered.Schedule()), 0)  // No Schedule tasks
		is.Equal(len(filtered.Delegate()), 1)  // Task 3
		is.Equal(len(filtered.AllTodos()), 2)  // Total filtered
	})

	t.Run("returns empty matrix when no matches", func(t *testing.T) {
		is := is.New(t)

		todos := []todo.Todo{
			todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{"computer"}),
		}

		m := matrix.New(todos)
		filtered := m.FilterByTag("+NonExistent")

		is.Equal(len(filtered.AllTodos()), 0) // No matches
	})

	t.Run("returns original matrix when filter is empty", func(t *testing.T) {
		is := is.New(t)

		todos := []todo.Todo{
			todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{}),
			todo.NewWithTags("Task 2", todo.PriorityB, []string{"Mobile"}, []string{}),
		}

		m := matrix.New(todos)
		filtered := m.FilterByTag("")

		is.Equal(len(filtered.AllTodos()), 2) // All todos present
	})

	t.Run("filter is case insensitive", func(t *testing.T) {
		is := is.New(t)

		todos := []todo.Todo{
			todo.NewWithTags("Task 1", todo.PriorityA, []string{"WebApp"}, []string{}),
		}

		m := matrix.New(todos)
		filtered := m.FilterByTag("+webapp") // lowercase

		is.Equal(len(filtered.AllTodos()), 1) // Should match case-insensitively
	})
}
