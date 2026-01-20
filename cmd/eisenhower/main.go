package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/file"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/usecases"
)

func main() {
	// Check if stdin is being piped
	stat, _ := os.Stdin.Stat()
	isStdinPiped := (stat.Mode() & os.ModeCharDevice) == 0

	var repo usecases.TodoRepository
	var filePath string
	var readOnly bool

	if isStdinPiped {
		// Read from stdin in read-only mode
		todos, err := todotxt.Unmarshal(os.Stdin)
		if err != nil {
			fmt.Printf("Error parsing stdin: %v\n", err)
			os.Exit(1)
		}

		repo = memory.NewRepository()
		if err := repo.SaveAll(todos); err != nil {
			fmt.Printf("Error loading todos from stdin: %v\n", err)
			os.Exit(1)
		}

		filePath = "(stdin)"
		readOnly = true
	} else {
		// Normal file mode
		var err error
		filePath, err = getFilePath()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		repo = file.NewRepository(filePath)
		readOnly = false
	}

	m, err := usecases.LoadMatrix(repo)
	if err != nil {
		fmt.Printf("Error loading todos: %v\n", err)
		os.Exit(1)
	}

	model := ui.NewModel(m, filePath).SetRepository(repo)
	if readOnly {
		model = model.SetReadOnly(true)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
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
