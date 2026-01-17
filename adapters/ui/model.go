package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// ViewMode represents the current viewing mode
type ViewMode int

const (
	Overview ViewMode = iota
	FocusDoFirst
	FocusSchedule
	FocusDelegate
	FocusEliminate
	Inventory
)

// Model represents the Bubble Tea model for the Eisenhower matrix UI
type Model struct {
	matrix             matrix.Matrix
	filePath           string
	width              int
	height             int
	viewMode           ViewMode
	inputMode          bool
	moveMode           bool       // true when in move mode (selecting quadrant to move to)
	deleteMode         bool       // true when in delete confirmation mode
	input              textinput.Model
	allProjects        []string
	allContexts        []string
	source             usecases.TodoSource
	writer             usecases.TodoWriter
	showSuggestions    bool
	suggestions        []string
	selectedSuggestion int
	selectedTodoIndex  int // index of selected todo in current quadrant
	todoTable          table.Model // table for displaying todos
	inventoryViewport  viewport.Model // viewport for scrollable inventory dashboard
}

// NewModel creates a new UI model with the given matrix and file path
func NewModel(m matrix.Matrix, filePath string) Model {
	// Extract all tags from the matrix
	projects, contexts := extractAllTags(m)

	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Enter todo description..."
	ti.CharLimit = 200
	ti.Width = 80

	return Model{
		matrix:      m,
		filePath:    filePath,
		viewMode:    Overview,
		inputMode:   false,
		input:       ti,
		allProjects: projects,
		allContexts: contexts,
	}
}

// SetSource sets the source for reloading todos
func (m Model) SetSource(s usecases.TodoSource) Model {
	m.source = s
	return m
}

// SetWriter sets the writer for saving todos
func (m Model) SetWriter(w usecases.TodoWriter) Model {
	m.writer = w
	return m
}

// Init initializes the model (required by tea.Model interface)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages (required by tea.Model interface)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle delete mode separately
		if m.deleteMode {
			switch msg.String() {
			case "y":
				// Confirm deletion
				m = m.deleteTodo()
				m.deleteMode = false
				return m, nil
			case "n", "esc":
				// Cancel deletion
				m.deleteMode = false
				return m, nil
			}
			// Ignore all other keys in delete mode
			return m, nil
		}

		// Handle move mode separately
		if m.moveMode {
			switch msg.String() {
			case "1":
				m = m.changeTodoPriority(todo.PriorityA)
				m.moveMode = false
				return m, nil
			case "2":
				m = m.changeTodoPriority(todo.PriorityB)
				m.moveMode = false
				return m, nil
			case "3":
				m = m.changeTodoPriority(todo.PriorityC)
				m.moveMode = false
				return m, nil
			case "4":
				m = m.changeTodoPriority(todo.PriorityD)
				m.moveMode = false
				return m, nil
			case "esc":
				m.moveMode = false
				return m, nil
			}
			// Ignore all other keys in move mode
			return m, nil
		}

		// Handle input mode separately
		if m.inputMode {
			// Handle autocomplete-specific keys when suggestions are visible
			if m.showSuggestions && len(m.suggestions) > 0 {
				switch msg.String() {
				case "down":
					m.selectedSuggestion = (m.selectedSuggestion + 1) % len(m.suggestions)
					return m, nil
				case "up":
					m.selectedSuggestion = (m.selectedSuggestion - 1 + len(m.suggestions)) % len(m.suggestions)
					return m, nil
				case "tab", "enter":
					// Complete the suggestion
					m = m.completeSuggestion()
					return m, nil
				case "esc":
					// Dismiss suggestions but stay in input mode
					m.showSuggestions = false
					return m, nil
				}
			}

			// Handle regular input mode keys
			switch msg.String() {
			case "enter":
				// Only save if suggestions are not visible
				if !m.showSuggestions {
					m = m.saveTodo()
					return m, nil
				}
			case "esc":
				// Cancel input mode entirely
				m.inputMode = false
				m.input.SetValue("")
				m.showSuggestions = false
				return m, nil
			}

			// Delegate to textinput for character input
			m.input, cmd = m.input.Update(msg)

			// Update autocomplete suggestions after input changes
			m = m.updateSuggestions()

			return m, cmd
		}

		// Normal mode key handling
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			if m.viewMode == Overview {
				// Overview mode: focus on quadrant
				m.viewMode = FocusDoFirst
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			} else {
				// Focus mode: jump to DO FIRST quadrant
				m.viewMode = FocusDoFirst
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			}
		case "2":
			if m.viewMode == Overview {
				// Overview mode: focus on quadrant
				m.viewMode = FocusSchedule
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			} else {
				// Focus mode: jump to SCHEDULE quadrant
				m.viewMode = FocusSchedule
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			}
		case "3":
			if m.viewMode == Overview {
				// Overview mode: focus on quadrant
				m.viewMode = FocusDelegate
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			} else {
				// Focus mode: jump to DELEGATE quadrant
				m.viewMode = FocusDelegate
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			}
		case "4":
			if m.viewMode == Overview {
				// Overview mode: focus on quadrant
				m.viewMode = FocusEliminate
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			} else {
				// Focus mode: jump to ELIMINATE quadrant
				m.viewMode = FocusEliminate
				m.selectedTodoIndex = 0
				m = m.rebuildTable()
			}
		case "m":
			// Enter move mode (only in focus mode with todos)
			if m.viewMode != Overview && len(m.currentQuadrantTodos()) > 0 {
				m.moveMode = true
			}
		case "backspace":
			// Enter delete mode (only in focus mode with todos)
			if m.viewMode != Overview && len(m.currentQuadrantTodos()) > 0 {
				m.deleteMode = true
			}
		case "i":
			// Toggle inventory mode (only in overview)
			if m.viewMode == Overview {
				m.viewMode = Inventory
				// Initialize viewport for inventory dashboard
				m.inventoryViewport = viewport.New(m.width, m.height-2) // Reserve space for help text
				content := RenderInventoryDashboard(m.matrix, m.width, 0) // Pass width for centering
				m.inventoryViewport.SetContent(content)
			} else if m.viewMode == Inventory {
				m.viewMode = Overview
			}
		case "esc":
			// Return to overview from any other mode
			if m.viewMode == Inventory {
				m.viewMode = Overview
			} else if m.viewMode != Overview {
				m.viewMode = Overview
			}
		case "a":
			// Enter input mode only if in focus mode
			if m.viewMode != Overview {
				m.inputMode = true
				m.input.Focus()
			}
		case "down", "s", "j":
			// Scroll down in inventory mode
			if m.viewMode == Inventory {
				var cmd tea.Cmd
				m.inventoryViewport, cmd = m.inventoryViewport.Update(msg)
				return m, cmd
			}
			// Navigate down in focus mode using table
			if m.viewMode != Overview {
				var cmd tea.Cmd
				m.todoTable, cmd = m.todoTable.Update(msg)
				m.selectedTodoIndex = m.todoTable.Cursor()
				return m, cmd
			}
		case "up", "w", "k":
			// Scroll up in inventory mode
			if m.viewMode == Inventory {
				var cmd tea.Cmd
				m.inventoryViewport, cmd = m.inventoryViewport.Update(msg)
				return m, cmd
			}
			// Navigate up in focus mode using table
			if m.viewMode != Overview {
				var cmd tea.Cmd
				m.todoTable, cmd = m.todoTable.Update(msg)
				m.selectedTodoIndex = m.todoTable.Cursor()
				return m, cmd
			}
		case " ":
			// Toggle completion in focus mode (space bar)
			if m.viewMode != Overview {
				m = m.toggleCompletion()
			}
		}

		// Handle pgup/pgdown in inventory mode
		if m.viewMode == Inventory {
			switch msg.String() {
			case "pgup", "pgdown", "home", "end":
				var cmd tea.Cmd
				m.inventoryViewport, cmd = m.inventoryViewport.Update(msg)
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update viewport size if in inventory mode
		if m.viewMode == Inventory {
			m.inventoryViewport.Width = msg.Width
			m.inventoryViewport.Height = msg.Height - 2
		}
	}
	return m, nil
}

// saveTodo saves the current input as a new todo
func (m Model) saveTodo() Model {
	description := m.input.Value()

	// Don't create empty todos
	if description == "" || len(description) == 0 {
		m.inputMode = false
		m.input.SetValue("")
		return m
	}

	if m.writer == nil {
		return m // No-op if no writer configured
	}

	// Determine priority from current quadrant
	priority := m.currentQuadrantPriority()

	// Use the AddTodo usecase
	updatedMatrix, err := usecases.AddTodo(m.writer, m.matrix, description, priority)
	if err != nil {
		// TODO: Show error to user in future story
		return m
	}

	m.matrix = updatedMatrix

	// Refresh tag lists
	m.allProjects, m.allContexts = extractAllTags(m.matrix)

	// Reset selection to first todo
	m.selectedTodoIndex = 0

	// Exit input mode
	m.inputMode = false
	m.input.SetValue("")

	// Rebuild table with new todo
	m = m.rebuildTable()

	return m
}

// updateSuggestions updates autocomplete suggestions based on current input
func (m Model) updateSuggestions() Model {
	inputValue := m.input.Value()

	// Detect if we're at a tag trigger
	trigger, partialTag, found := detectTrigger(inputValue)
	if !found {
		m.showSuggestions = false
		return m
	}

	// Get the appropriate tag list
	var tagList []string
	if trigger == "+" {
		tagList = m.allProjects
	} else if trigger == "@" {
		tagList = m.allContexts
	}

	// Filter tags by partial input
	m.suggestions = filterTags(tagList, partialTag)
	m.showSuggestions = len(m.suggestions) > 0 || partialTag != ""
	m.selectedSuggestion = 0 // Reset selection to first item

	return m
}

// completeSuggestion completes the currently selected suggestion
func (m Model) completeSuggestion() Model {
	if len(m.suggestions) == 0 {
		return m
	}

	selectedTag := m.suggestions[m.selectedSuggestion]
	completedValue := completeTag(m.input.Value(), selectedTag)

	m.input.SetValue(completedValue)
	m.showSuggestions = false
	m.suggestions = nil
	m.selectedSuggestion = 0

	// Move cursor to end
	m.input.SetCursor(len(completedValue))

	return m
}

// currentQuadrantPriority returns the priority for the current focused quadrant
func (m Model) currentQuadrantPriority() todo.Priority {
	switch m.viewMode {
	case FocusDoFirst:
		return todo.PriorityA
	case FocusSchedule:
		return todo.PriorityB
	case FocusDelegate:
		return todo.PriorityC
	case FocusEliminate:
		return todo.PriorityD
	default:
		return todo.PriorityNone
	}
}

// currentQuadrantTodos returns the todos for the current focused quadrant
func (m Model) currentQuadrantTodos() []todo.Todo {
	switch m.viewMode {
	case FocusDoFirst:
		return m.matrix.DoFirst()
	case FocusSchedule:
		return m.matrix.Schedule()
	case FocusDelegate:
		return m.matrix.Delegate()
	case FocusEliminate:
		return m.matrix.Eliminate()
	default:
		return []todo.Todo{}
	}
}

// currentQuadrantType returns the quadrant type for the current view mode
func (m Model) currentQuadrantType() matrix.QuadrantType {
	switch m.viewMode {
	case FocusDoFirst:
		return matrix.DoFirstQuadrant
	case FocusSchedule:
		return matrix.ScheduleQuadrant
	case FocusDelegate:
		return matrix.DelegateQuadrant
	case FocusEliminate:
		return matrix.EliminateQuadrant
	default:
		return matrix.DoFirstQuadrant
	}
}

// moveSelectionDown moves the selection to the next todo (with wrap-around)
func (m Model) moveSelectionDown() Model {
	todos := m.currentQuadrantTodos()
	if len(todos) == 0 {
		return m
	}
	m.selectedTodoIndex = (m.selectedTodoIndex + 1) % len(todos)
	return m
}

// moveSelectionUp moves the selection to the previous todo (with wrap-around)
func (m Model) moveSelectionUp() Model {
	todos := m.currentQuadrantTodos()
	if len(todos) == 0 {
		return m
	}
	m.selectedTodoIndex = (m.selectedTodoIndex - 1 + len(todos)) % len(todos)
	return m
}

// toggleCompletion toggles the completion status of the selected todo
func (m Model) toggleCompletion() Model {
	if m.writer == nil {
		return m // No-op if no writer configured
	}

	// Use the ToggleCompletion usecase
	quadrant := m.currentQuadrantType()
	updatedMatrix, err := usecases.ToggleCompletion(m.writer, m.matrix, quadrant, m.selectedTodoIndex)
	if err != nil {
		// TODO: Show error to user in future story
		return m
	}

	m.matrix = updatedMatrix

	// Rebuild table to reflect the change
	m = m.rebuildTable()

	return m
}

// changeTodoPriority changes the priority of the selected todo
func (m Model) changeTodoPriority(newPriority todo.Priority) Model {
	if m.writer == nil || m.source == nil {
		return m // No-op if no writer or source configured
	}

	// Use the ChangePriority usecase
	quadrant := m.currentQuadrantType()
	updatedMatrix, err := usecases.ChangePriority(m.source, m.writer, m.matrix, quadrant, m.selectedTodoIndex, newPriority)
	if err != nil {
		// TODO: Show error to user in future story
		return m
	}

	m.matrix = updatedMatrix

	// After moving a todo, adjust the view:
	// - If the current quadrant is now empty, return to overview
	// - Otherwise, adjust selection index if needed
	todos := m.currentQuadrantTodos()
	if len(todos) == 0 {
		m.viewMode = Overview
	} else {
		if m.selectedTodoIndex >= len(todos) {
			// If selected index is now out of bounds, select the last todo
			m.selectedTodoIndex = len(todos) - 1
		}
		// Rebuild table to reflect the change
		m = m.rebuildTable()
	}

	return m
}

// deleteTodo deletes the currently selected todo
func (m Model) deleteTodo() Model {
	if m.writer == nil {
		return m // No-op if no writer configured
	}

	// Get the todo to delete
	todos := m.currentQuadrantTodos()
	if m.selectedTodoIndex < 0 || m.selectedTodoIndex >= len(todos) {
		return m // Invalid index
	}

	todoToDelete := todos[m.selectedTodoIndex]

	// Use the DeleteTodo usecase
	updatedMatrix, err := usecases.DeleteTodo(m.writer, m.matrix, todoToDelete)
	if err != nil {
		// TODO: Show error to user in future story
		return m
	}

	m.matrix = updatedMatrix

	// After deleting a todo:
	// - If the current quadrant is now empty, return to overview
	// - Otherwise, adjust selection index if needed
	todos = m.currentQuadrantTodos()
	if len(todos) == 0 {
		m.viewMode = Overview
	} else {
		if m.selectedTodoIndex >= len(todos) {
			// If selected index is now out of bounds, select the last todo
			m.selectedTodoIndex = len(todos) - 1
		}
		// Rebuild table to reflect the change
		m = m.rebuildTable()
	}

	return m
}

// rebuildTable rebuilds the todo table based on current quadrant
func (m Model) rebuildTable() Model {
	if m.viewMode == Overview {
		return m // No table in overview mode
	}

	todos := m.currentQuadrantTodos()
	m.todoTable = buildTodoTable(todos, m.width, m.height, m.selectedTodoIndex)
	return m
}

// View renders the model (required by tea.Model interface)
func (m Model) View() string {
	var content string

	// Render based on current view mode
	switch m.viewMode {
	case FocusDoFirst:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.DoFirst(),
				"Do First",
				lipgloss.Color("#FF6B6B"),
				m.filePath,
				m.input,
				m.allProjects,
				m.allContexts,
				m.showSuggestions,
				m.suggestions,
				m.selectedSuggestion,
				m.width,
				m.height,
			)
		} else {
			content = RenderFocusedQuadrantWithTable(
				m.matrix.DoFirst(),
				"Do First",
				lipgloss.Color("#FF6B6B"),
				m.filePath,
				m.todoTable,
				m.width,
				m.height,
			)
		}
	case FocusSchedule:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Schedule(),
				"Schedule",
				lipgloss.Color("#4ECDC4"),
				m.filePath,
				m.input,
				m.allProjects,
				m.allContexts,
				m.showSuggestions,
				m.suggestions,
				m.selectedSuggestion,
				m.width,
				m.height,
			)
		} else {
			content = RenderFocusedQuadrantWithTable(
				m.matrix.Schedule(),
				"Schedule",
				lipgloss.Color("#4ECDC4"),
				m.filePath,
				m.todoTable,
				m.width,
				m.height,
			)
		}
	case FocusDelegate:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Delegate(),
				"Delegate",
				lipgloss.Color("#FFE66D"),
				m.filePath,
				m.input,
				m.allProjects,
				m.allContexts,
				m.showSuggestions,
				m.suggestions,
				m.selectedSuggestion,
				m.width,
				m.height,
			)
		} else {
			content = RenderFocusedQuadrantWithTable(
				m.matrix.Delegate(),
				"Delegate",
				lipgloss.Color("#FFE66D"),
				m.filePath,
				m.todoTable,
				m.width,
				m.height,
			)
		}
	case FocusEliminate:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Eliminate(),
				"Eliminate",
				lipgloss.Color("#95E1D3"),
				m.filePath,
				m.input,
				m.allProjects,
				m.allContexts,
				m.showSuggestions,
				m.suggestions,
				m.selectedSuggestion,
				m.width,
				m.height,
			)
		} else {
			content = RenderFocusedQuadrantWithTable(
				m.matrix.Eliminate(),
				"Eliminate",
				lipgloss.Color("#95E1D3"),
				m.filePath,
				m.todoTable,
				m.width,
				m.height,
			)
		}
	case Inventory:
		content = m.inventoryViewport.View()
	default: // Overview
		// Pass terminal dimensions to RenderMatrix for responsive sizing
		content = RenderMatrix(m.matrix, m.filePath, m.width, m.height)

		// Center the content in the terminal if we have dimensions
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}

		return content
	}

	// If in move mode, overlay the move dialog
	if m.moveMode {
		return RenderMoveOverlay(m.width, m.height)
	}

	// If in delete mode, overlay the delete confirmation dialog
	if m.deleteMode {
		return RenderDeleteOverlay(m.width, m.height)
	}

	// Focus mode content is already full-width and properly aligned
	return content
}

// NewModelWithWriter creates a model with explicit writer (for testing)
func NewModelWithWriter(m matrix.Matrix, filePath string, source usecases.TodoSource, writer usecases.TodoWriter) Model {
	model := NewModel(m, filePath)
	model.source = source
	model.writer = writer
	return model
}

// GetMatrix returns the current matrix (for testing)
func (m Model) GetMatrix() matrix.Matrix {
	return m.matrix
}
