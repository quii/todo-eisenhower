package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/file"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	source := file.NewSource(filepath.Join(homeDir, "todo.txt"))

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(m)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
