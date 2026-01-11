package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/file"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

func main() {
	filePath, err := getFilePath()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	source := file.NewSource(filePath)

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(m, filePath)

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}

// getFilePath returns the file path from CLI args or default ~/todo.txt
func getFilePath() (string, error) {
	var path string

	// Get path from args or use default
	if len(os.Args) > 1 {
		path = os.Args[1]
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting home directory: %w", err)
		}
		path = filepath.Join(homeDir, "todo.txt")
	}

	// Expand tilde
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expanding tilde: %w", err)
		}
		path = filepath.Join(homeDir, path[2:])
	}

	// Convert to absolute path if relative
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("converting to absolute path: %w", err)
	}

	return absPath, nil
}
