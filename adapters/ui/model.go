package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/quii/todo-eisenhower/domain/matrix"
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
	matrix   matrix.Matrix
	filePath string
	width    int
	height   int
	viewMode ViewMode
}

// NewModel creates a new UI model with the given matrix and file path
func NewModel(m matrix.Matrix, filePath string) Model {
	return Model{
		matrix:   m,
		filePath: filePath,
		viewMode: Overview,
	}
}

// Init initializes the model (required by tea.Model interface)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages (required by tea.Model interface)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the model (required by tea.Model interface)
func (m Model) View() string {
	var content string

	// Render based on current view mode
	switch m.viewMode {
	case FocusDoFirst:
		content = RenderFocusedQuadrant(
			m.matrix.DoFirst(),
			"DO FIRST",
			lipgloss.Color("#FF6B6B"),
			m.filePath,
			m.width,
			m.height,
		)
	case FocusSchedule:
		content = RenderFocusedQuadrant(
			m.matrix.Schedule(),
			"SCHEDULE",
			lipgloss.Color("#4ECDC4"),
			m.filePath,
			m.width,
			m.height,
		)
	case FocusDelegate:
		content = RenderFocusedQuadrant(
			m.matrix.Delegate(),
			"DELEGATE",
			lipgloss.Color("#FFE66D"),
			m.filePath,
			m.width,
			m.height,
		)
	case FocusEliminate:
		content = RenderFocusedQuadrant(
			m.matrix.Eliminate(),
			"ELIMINATE",
			lipgloss.Color("#95E1D3"),
			m.filePath,
			m.width,
			m.height,
		)
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
