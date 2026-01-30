package matrix_test

import (
	"testing"
	"time"

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

func TestMatrix_ArchiveTodoAt(t *testing.T) {
	t.Run("archives completed todo and removes from matrix", func(t *testing.T) {
		is := is.New(t)

		activeTodo := todo.New("Active task", todo.PriorityA)
		completedTodo := todo.New("Completed task", todo.PriorityA).ToggleCompletion(time.Now())

		m := matrix.New([]todo.Todo{activeTodo, completedTodo})

		// Archive the completed todo (index 1 in DoFirst)
		archived, updated, success := m.ArchiveTodoAt(matrix.DoFirstQuadrant, 1)

		is.True(success)                                        // should succeed
		is.Equal(archived.Description(), "Completed task")      // archived todo should match
		is.True(archived.IsCompleted())                         // archived todo should be completed
		is.Equal(len(updated.DoFirst()), 1)                     // should have one less todo
		is.Equal(updated.DoFirst()[0].Description(), "Active task") // remaining todo
		is.True(!updated.DoFirst()[0].IsCompleted())            // remaining todo should not be completed
	})

	t.Run("cannot archive uncompleted todo", func(t *testing.T) {
		is := is.New(t)

		activeTodo := todo.New("Active task", todo.PriorityA)
		m := matrix.New([]todo.Todo{activeTodo})

		_, updated, success := m.ArchiveTodoAt(matrix.DoFirstQuadrant, 0)

		is.True(!success)                      // should fail
		is.Equal(len(updated.DoFirst()), 1)    // matrix unchanged
	})

	t.Run("returns false for invalid index", func(t *testing.T) {
		is := is.New(t)

		activeTodo := todo.New("Active task", todo.PriorityA)
		m := matrix.New([]todo.Todo{activeTodo})

		_, updated, success := m.ArchiveTodoAt(matrix.DoFirstQuadrant, 99)

		is.True(!success)                   // should fail
		is.Equal(len(updated.DoFirst()), 1) // matrix unchanged
	})

	t.Run("archives from different quadrants", func(t *testing.T) {
		is := is.New(t)

		doFirstTodo := todo.New("Do First", todo.PriorityA).ToggleCompletion(time.Now())
		scheduleTodo := todo.New("Schedule", todo.PriorityB).ToggleCompletion(time.Now())
		delegateTodo := todo.New("Delegate", todo.PriorityC).ToggleCompletion(time.Now())
		eliminateTodo := todo.New("Eliminate", todo.PriorityD).ToggleCompletion(time.Now())

		m := matrix.New([]todo.Todo{doFirstTodo, scheduleTodo, delegateTodo, eliminateTodo})

		// Archive from each quadrant
		_, m, success := m.ArchiveTodoAt(matrix.DoFirstQuadrant, 0)
		is.True(success)
		is.Equal(len(m.DoFirst()), 0)

		_, m, success = m.ArchiveTodoAt(matrix.ScheduleQuadrant, 0)
		is.True(success)
		is.Equal(len(m.Schedule()), 0)

		_, m, success = m.ArchiveTodoAt(matrix.DelegateQuadrant, 0)
		is.True(success)
		is.Equal(len(m.Delegate()), 0)

		_, m, success = m.ArchiveTodoAt(matrix.EliminateQuadrant, 0)
		is.True(success)
		is.Equal(len(m.Eliminate()), 0)
	})
}

func TestMatrix_ArchiveCompletedInQuadrant(t *testing.T) {
	t.Run("archives all completed todos in quadrant", func(t *testing.T) {
		is := is.New(t)

		active1 := todo.New("Active 1", todo.PriorityA)
		active2 := todo.New("Active 2", todo.PriorityA)
		completed1 := todo.New("Completed 1", todo.PriorityA).ToggleCompletion(time.Now())
		completed2 := todo.New("Completed 2", todo.PriorityA).ToggleCompletion(time.Now())
		completed3 := todo.New("Completed 3", todo.PriorityA).ToggleCompletion(time.Now())

		m := matrix.New([]todo.Todo{active1, active2, completed1, completed2, completed3})

		archived, updated := m.ArchiveCompletedInQuadrant(matrix.DoFirstQuadrant)

		is.Equal(len(archived), 3)           // should archive 3 completed todos
		is.Equal(len(updated.DoFirst()), 2)  // should have 2 active todos remaining
		is.True(!updated.DoFirst()[0].IsCompleted()) // remaining todos are active
		is.True(!updated.DoFirst()[1].IsCompleted())
	})

	t.Run("returns empty slice when no completed todos in quadrant", func(t *testing.T) {
		is := is.New(t)

		active1 := todo.New("Active 1", todo.PriorityA)
		active2 := todo.New("Active 2", todo.PriorityA)

		m := matrix.New([]todo.Todo{active1, active2})

		archived, updated := m.ArchiveCompletedInQuadrant(matrix.DoFirstQuadrant)

		is.Equal(len(archived), 0)           // nothing to archive
		is.Equal(len(updated.DoFirst()), 2)  // matrix unchanged
	})

	t.Run("returns empty slice for empty quadrant", func(t *testing.T) {
		is := is.New(t)

		m := matrix.New([]todo.Todo{})

		archived, updated := m.ArchiveCompletedInQuadrant(matrix.DoFirstQuadrant)

		is.Equal(len(archived), 0)          // nothing to archive
		is.Equal(len(updated.DoFirst()), 0) // still empty
	})

	t.Run("only archives from specified quadrant", func(t *testing.T) {
		is := is.New(t)

		doFirstCompleted := todo.New("DoFirst Completed", todo.PriorityA).ToggleCompletion(time.Now())
		scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())

		m := matrix.New([]todo.Todo{doFirstCompleted, scheduleCompleted})

		archived, updated := m.ArchiveCompletedInQuadrant(matrix.DoFirstQuadrant)

		is.Equal(len(archived), 1)                              // only DoFirst completed
		is.Equal(archived[0].Description(), "DoFirst Completed")
		is.Equal(len(updated.DoFirst()), 0)                     // DoFirst is empty
		is.Equal(len(updated.Schedule()), 1)                    // Schedule unchanged
		is.True(updated.Schedule()[0].IsCompleted())            // Schedule todo still there
	})
}

func TestMatrix_ArchiveAllCompleted(t *testing.T) {
	t.Run("archives all completed todos across all quadrants", func(t *testing.T) {
		is := is.New(t)

		doFirstActive := todo.New("DoFirst Active", todo.PriorityA)
		doFirstCompleted1 := todo.New("DoFirst Completed 1", todo.PriorityA).ToggleCompletion(time.Now())
		doFirstCompleted2 := todo.New("DoFirst Completed 2", todo.PriorityA).ToggleCompletion(time.Now())
		scheduleCompleted := todo.New("Schedule Completed", todo.PriorityB).ToggleCompletion(time.Now())
		delegateActive := todo.New("Delegate Active", todo.PriorityC)
		delegateCompleted1 := todo.New("Delegate Completed 1", todo.PriorityC).ToggleCompletion(time.Now())
		delegateCompleted2 := todo.New("Delegate Completed 2", todo.PriorityC).ToggleCompletion(time.Now())
		delegateCompleted3 := todo.New("Delegate Completed 3", todo.PriorityC).ToggleCompletion(time.Now())
		// Eliminate has no completed todos

		m := matrix.New([]todo.Todo{
			doFirstActive, doFirstCompleted1, doFirstCompleted2,
			scheduleCompleted,
			delegateActive, delegateCompleted1, delegateCompleted2, delegateCompleted3,
		})

		archived, updated := m.ArchiveAllCompleted()

		is.Equal(len(archived), 6)            // 2 + 1 + 3 = 6 completed todos
		is.Equal(len(updated.DoFirst()), 1)   // 1 active remains
		is.Equal(len(updated.Schedule()), 0)  // all were completed
		is.Equal(len(updated.Delegate()), 1)  // 1 active remains
		is.Equal(len(updated.Eliminate()), 0) // was empty
	})

	t.Run("returns empty slice when no completed todos anywhere", func(t *testing.T) {
		is := is.New(t)

		active1 := todo.New("Active 1", todo.PriorityA)
		active2 := todo.New("Active 2", todo.PriorityB)
		active3 := todo.New("Active 3", todo.PriorityC)

		m := matrix.New([]todo.Todo{active1, active2, active3})

		archived, updated := m.ArchiveAllCompleted()

		is.Equal(len(archived), 0)           // nothing to archive
		is.Equal(len(updated.AllTodos()), 3) // all todos still present
	})

	t.Run("returns empty slice for empty matrix", func(t *testing.T) {
		is := is.New(t)

		m := matrix.New([]todo.Todo{})

		archived, updated := m.ArchiveAllCompleted()

		is.Equal(len(archived), 0)           // nothing to archive
		is.Equal(len(updated.AllTodos()), 0) // still empty
	})
}

func TestMatrix_Backlog(t *testing.T) {
	t.Run("categorizes priority E todos into Backlog", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Idea for later", todo.PriorityE),
			todo.New("Another idea", todo.PriorityE),
		}

		m := matrix.New(todos)

		backlog := m.Backlog()
		is.Equal(len(backlog), 2)
		assertContainsTodo(is, backlog, "Idea for later")
		assertContainsTodo(is, backlog, "Another idea")
	})

	t.Run("backlog is separate from Eisenhower quadrants", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Do first task", todo.PriorityA),
			todo.New("Backlog idea", todo.PriorityE),
		}

		m := matrix.New(todos)

		is.Equal(len(m.DoFirst()), 1)
		is.Equal(len(m.Backlog()), 1)
		is.Equal(len(m.Eliminate()), 0) // Backlog should NOT go to Eliminate
	})

	t.Run("AllTodos excludes backlog", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Do first", todo.PriorityA),
			todo.New("Schedule", todo.PriorityB),
			todo.New("Backlog idea", todo.PriorityE),
		}

		m := matrix.New(todos)

		allTodos := m.AllTodos()
		is.Equal(len(allTodos), 2) // Only Eisenhower quadrants
		for _, t := range allTodos {
			is.True(t.Priority() != todo.PriorityE) // No backlog items
		}
	})

	t.Run("AllTodosIncludingBacklog includes everything", func(t *testing.T) {
		is := is.New(t)
		todos := []todo.Todo{
			todo.New("Do first", todo.PriorityA),
			todo.New("Schedule", todo.PriorityB),
			todo.New("Backlog idea", todo.PriorityE),
		}

		m := matrix.New(todos)

		allTodos := m.AllTodosIncludingBacklog()
		is.Equal(len(allTodos), 3) // All todos including backlog
	})

	t.Run("AddTodo with priority E goes to backlog", func(t *testing.T) {
		is := is.New(t)
		m := matrix.New([]todo.Todo{})

		newTodo := todo.New("New idea", todo.PriorityE)
		m = m.AddTodo(newTodo)

		is.Equal(len(m.Backlog()), 1)
		is.Equal(len(m.Eliminate()), 0)
	})
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
