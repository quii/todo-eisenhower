package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
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
)

// Model represents the Bubble Tea model for the Eisenhower matrix UI
type Model struct {
	matrix             matrix.Matrix
	filePath           string
	width              int
	height             int
	viewMode           ViewMode
	inputMode          bool
	input              textinput.Model
	allProjects        []string
	allContexts        []string
	source             usecases.TodoSource
	writer             usecases.TodoWriter
	showSuggestions    bool
	suggestions        []string
	selectedSuggestion int
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
			m.viewMode = FocusDoFirst
		case "2":
			m.viewMode = FocusSchedule
		case "3":
			m.viewMode = FocusDelegate
		case "4":
			m.viewMode = FocusEliminate
		case "esc":
			m.viewMode = Overview
		case "a":
			// Enter input mode only if in focus mode
			if m.viewMode != Overview {
				m.inputMode = true
				m.input.Focus()
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
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

	// Determine priority from current quadrant
	priority := m.currentQuadrantPriority()

	// Parse tags from description (they're already in the description)
	// Create the todo
	t := todo.New(description, priority)

	// Save to file if writer is set
	if m.writer != nil {
		_ = usecases.SaveTodo(m.writer, t)
	}

	// Add the todo to the matrix in memory
	m.matrix = m.matrix.AddTodo(t)

	// Refresh tag lists
	m.allProjects, m.allContexts = extractAllTags(m.matrix)

	// Exit input mode
	m.inputMode = false
	m.input.SetValue("")

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

// View renders the model (required by tea.Model interface)
func (m Model) View() string {
	var content string

	// Render based on current view mode
	switch m.viewMode {
	case FocusDoFirst:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.DoFirst(),
				"DO FIRST",
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
			content = RenderFocusedQuadrant(
				m.matrix.DoFirst(),
				"DO FIRST",
				lipgloss.Color("#FF6B6B"),
				m.filePath,
				m.width,
				m.height,
			)
		}
	case FocusSchedule:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Schedule(),
				"SCHEDULE",
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
			content = RenderFocusedQuadrant(
				m.matrix.Schedule(),
				"SCHEDULE",
				lipgloss.Color("#4ECDC4"),
				m.filePath,
				m.width,
				m.height,
			)
		}
	case FocusDelegate:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Delegate(),
				"DELEGATE",
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
			content = RenderFocusedQuadrant(
				m.matrix.Delegate(),
				"DELEGATE",
				lipgloss.Color("#FFE66D"),
				m.filePath,
				m.width,
				m.height,
			)
		}
	case FocusEliminate:
		if m.inputMode {
			content = RenderFocusedQuadrantWithInput(
				m.matrix.Eliminate(),
				"ELIMINATE",
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
			content = RenderFocusedQuadrant(
				m.matrix.Eliminate(),
				"ELIMINATE",
				lipgloss.Color("#95E1D3"),
				m.filePath,
				m.width,
				m.height,
			)
		}
	default: // Overview
		// Pass terminal dimensions to RenderMatrix for responsive sizing
		content = RenderMatrix(m.matrix, m.filePath, m.width, m.height)

		// Center the content in the terminal if we have dimensions
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
		}

		return content
	}

	// Focus mode content is already full-width and properly aligned
	return content
}
