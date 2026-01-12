package acceptance_test

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/ui"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 006: Responsive Matrix Sizing

func TestStory006_MatrixFillsAvailableSpace(t *testing.T) {
	// Scenario: Matrix fills available terminal space
	// Given I have a todo.txt file with many todos
	// When I run the application in a 120x40 terminal
	// Then the matrix quadrants are larger than the default 40x10
	// And more todos are visible per quadrant

	// Create todos (15 in Priority A to test display limit)
	input := generateManyTodos(15)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate large terminal (120x40)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// View should use larger dimensions
	view := updatedModel.View()

	// With a 120x40 terminal, we should be able to display more than the default 7 todos
	// Calculate expected: 40 - 10 (reserved) = 30, / 2 = 15 lines per quadrant
	// 15 - 3 (title, spacing, footer) = 12 todos displayable
	// We have 15 todos, so should see "... and 3 more"
	if !strings.Contains(view, "... and 3 more") {
		t.Errorf("expected to see '... and 3 more' with larger terminal")
	}

	// Verify matrix content is present
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain matrix content")
	}
}

func TestStory006_MatrixAdjustsToDifferentSizes(t *testing.T) {
	// Scenario: Matrix adjusts to different terminal sizes
	// Given I have a todo.txt file
	// When I run the application in a 200x60 terminal
	// Then the quadrants are larger than in a 120x40 terminal

	input := generateManyTodos(30)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate very large terminal (200x60)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 200, Height: 60})

	view := updatedModel.View()

	// With 200x60: 60 - 10 = 50, / 2 = 25 lines per quadrant
	// 25 - 3 = 22 todos displayable
	// We have 30, so should see "... and 8 more"
	if !strings.Contains(view, "... and 8 more") {
		t.Errorf("expected to see '... and 8 more' with very large terminal")
	}
}

func TestStory006_MatrixRespectsMinimumDimensions(t *testing.T) {
	// Scenario: Matrix respects minimum dimensions
	// Given I have a todo.txt file
	// When I run the application in a very small terminal (80x24)
	// Then the matrix uses minimum viable dimensions
	// And does not break the layout

	input := "(A) Test todo"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Simulate small terminal (80x24)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := updatedModel.View()

	// Matrix should still render without errors
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain matrix content even in small terminal")
	}

	// Should show the test todo
	if !strings.Contains(view, "Test todo") {
		t.Error("expected view to contain test todo")
	}
}

func TestStory006_DisplayLimitScalesWithHeight(t *testing.T) {
	// Scenario: Todo display limit scales with height
	// Given I have 20 todos in the Do First quadrant
	// When the quadrant height varies
	// Then the number of visible todos scales accordingly

	tests := []struct {
		name            string
		terminalHeight  int
		expectedTodos   int
		expectedMessage string
	}{
		{
			name:            "small terminal shows fewer todos",
			terminalHeight:  30,
			expectedTodos:   7, // (30-10)/2 = 10, 10-3 = 7 todos
			expectedMessage: "... and 13 more",
		},
		{
			name:            "large terminal shows more todos",
			terminalHeight:  50,
			expectedTodos:   17, // (50-10)/2 = 20, 20-3 = 17 todos
			expectedMessage: "... and 3 more",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := generateManyTodos(20)
			source := &StubTodoSource{
				reader: strings.NewReader(input),
			}

			m, err := usecases.LoadMatrix(source)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			model := ui.NewModel(m, "test.txt")
			updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: tt.terminalHeight})

			view := updatedModel.View()

			if !strings.Contains(view, tt.expectedMessage) {
				t.Errorf("expected to see %q, but didn't see it in output", tt.expectedMessage)
			}
		})
	}
}

func TestStory006_WindowResizeHandledGracefully(t *testing.T) {
	// Scenario: Window resize is handled gracefully
	// Given the application is running
	// When I resize my terminal window
	// Then the matrix re-renders with new dimensions

	input := generateManyTodos(15)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Start with medium terminal
	model1, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view1 := model1.View()

	// Resize to larger terminal
	model2, _ := model1.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
	view2 := model2.View()

	// Views should be different (different dimensions)
	if view1 == view2 {
		t.Error("expected views to be different after resize")
	}

	// Both should contain matrix content
	if !strings.Contains(view1, "DO FIRST") || !strings.Contains(view2, "DO FIRST") {
		t.Error("expected both views to contain matrix content")
	}
}

func TestStory006_DefaultDimensionsWhenNoWindowSize(t *testing.T) {
	// Verify that default dimensions are used when no window size is received
	// (edge case for initial render before WindowSizeMsg)

	input := generateManyTodos(10)
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	model := ui.NewModel(m, "test.txt")

	// Call View without setting window size (width and height are 0)
	view := model.View()

	// Should use default displayLimit of 7, showing "... and 3 more"
	if !strings.Contains(view, "... and 3 more") {
		t.Error("expected default display limit when no window size")
	}

	// Should still contain matrix content
	if !strings.Contains(view, "DO FIRST") {
		t.Error("expected view to contain matrix content")
	}
}

// Helper function to generate many Priority A todos for testing
func generateManyTodos(count int) string {
	var builder strings.Builder
	for i := 1; i <= count; i++ {
		builder.WriteString("(A) Task number ")
		builder.WriteString(string(rune('0' + (i/10)%10)))
		builder.WriteString(string(rune('0' + i%10)))
		builder.WriteString("\n")
	}
	return builder.String()
}
