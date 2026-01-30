package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 028: Backlog Quadrant

func TestStory028_AccessBacklogByPressingFive(t *testing.T) {
	// Scenario: Access backlog by pressing 5
	is := is.New(t)

	repository := memory.NewRepository()
	backlogTodo := todo.New("Idea for later", todo.PriorityE)

	err := repository.SaveAll([]todo.Todo{backlogTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 5 to focus on Backlog
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	model = updatedModel.(ui.Model)

	// Verify we can see the backlog
	view := model.View()
	is.True(strings.Contains(view, "Backlog"))
	is.True(strings.Contains(view, "Idea for later"))
}

func TestStory028_BacklogUsesPriorityE(t *testing.T) {
	// Scenario: Backlog uses priority E in todo.txt
	is := is.New(t)

	repository := memory.NewRepository()
	backlogTodo := todo.New("Research new tech", todo.PriorityE)

	err := repository.SaveAll([]todo.Todo{backlogTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Backlog should contain the todo
	is.Equal(len(m.Backlog()), 1)
	is.Equal(m.Backlog()[0].Description(), "Research new tech")

	// It should NOT be in Eliminate
	is.Equal(len(m.Eliminate()), 0)
}

func TestStory028_AddTaskDirectlyToBacklog(t *testing.T) {
	// Scenario: Add task directly to Backlog
	is := is.New(t)

	repository := memory.NewRepository()
	err := repository.SaveAll([]todo.Todo{})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Press 5 to focus on Backlog
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	model = updatedModel.(ui.Model)

	// Press 'a' to add
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	model = updatedModel.(ui.Model)

	// Type description
	for _, r := range "New framework idea" {
		updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		model = updatedModel.(ui.Model)
	}

	// Press Enter to save
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Verify it was saved with priority E
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.Equal(todos[0].Priority(), todo.PriorityE)
	is.Equal(todos[0].Description(), "New framework idea")
}

func TestStory028_BacklogNotShownInOverview(t *testing.T) {
	// Scenario: Backlog not shown in overview
	is := is.New(t)

	repository := memory.NewRepository()
	doFirstTodo := todo.New("Urgent task", todo.PriorityA)
	backlogTodo := todo.New("Backlog idea", todo.PriorityE)

	err := repository.SaveAll([]todo.Todo{doFirstTodo, backlogTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// In overview mode
	view := model.View()

	// Should show Do First task
	is.True(strings.Contains(view, "Urgent task"))

	// Should NOT show backlog task in the matrix area
	// But SHOULD show backlog count in help text
	is.True(strings.Contains(view, "Backlog (1)"))
}

func TestStory028_MoveTaskFromBacklogToDoFirst(t *testing.T) {
	// Scenario: Move task from Backlog to Do First
	is := is.New(t)

	repository := memory.NewRepository()
	backlogTodo := todo.New("Promote this idea", todo.PriorityE)

	err := repository.SaveAll([]todo.Todo{backlogTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Backlog
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '1' to move to Do First
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	// Verify the todo moved to Do First with priority A
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.Equal(todos[0].Priority(), todo.PriorityA)
}

func TestStory028_MoveTaskFromDoFirstToBacklog(t *testing.T) {
	// Scenario: Move task from Do First to Backlog
	is := is.New(t)

	repository := memory.NewRepository()
	doFirstTodo := todo.New("Defer this task", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{doFirstTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Press '5' to move to Backlog
	_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})

	// Verify the todo moved to Backlog with priority E
	todos, err := repository.LoadAll()
	is.NoErr(err)
	is.Equal(len(todos), 1)
	is.Equal(todos[0].Priority(), todo.PriorityE)
}

func TestStory028_MoveOverlayShowsBacklogOption(t *testing.T) {
	// Scenario: Move mode shows Backlog option
	is := is.New(t)

	repository := memory.NewRepository()
	doFirstTodo := todo.New("Task", todo.PriorityA)

	err := repository.SaveAll([]todo.Todo{doFirstTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Do First
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	model = updatedModel.(ui.Model)

	// Press 'm' to enter move mode
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model = updatedModel.(ui.Model)

	// Verify move overlay shows Backlog option
	view := model.View()
	is.True(strings.Contains(view, "5. Backlog"))
}

func TestStory028_ReturnToOverviewFromBacklog(t *testing.T) {
	// Scenario: Return to overview from Backlog
	is := is.New(t)

	repository := memory.NewRepository()
	backlogTodo := todo.New("Idea", todo.PriorityE)

	err := repository.SaveAll([]todo.Todo{backlogTodo})
	is.NoErr(err)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModelWithRepository(m, "test.txt", repository)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = updatedModel.(ui.Model)

	// Focus on Backlog
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})
	model = updatedModel.(ui.Model)

	// Verify we're in backlog view
	view := model.View()
	is.True(strings.Contains(view, "Backlog"))

	// Press Esc to return to overview
	updatedModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
	model = updatedModel.(ui.Model)

	// Verify we're back in overview (shows the matrix with quadrant labels)
	view = model.View()
	is.True(strings.Contains(view, "Do First"))
	is.True(strings.Contains(view, "Schedule"))
}
