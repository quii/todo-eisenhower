package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/domain/matrix"
)

// Model represents the Bubble Tea model for the Eisenhower matrix UI
type Model struct {
	matrix   matrix.Matrix
	filePath string
}

// NewModel creates a new UI model with the given matrix and file path
func NewModel(m matrix.Matrix, filePath string) Model {
	return Model{
		matrix:   m,
		filePath: filePath,
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
		// Quit on 'q' or Ctrl+C
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the model (required by tea.Model interface)
func (m Model) View() string {
	return RenderMatrix(m.matrix, m.filePath)
}
