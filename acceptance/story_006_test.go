package acceptance_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/quii/todo-eisenhower/adapters/memory"
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
	is := is.New(t)

	// Create todos (15 in Priority A to test display limit)
	input := generateManyTodos(15)
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate large terminal (120x40)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// View should use larger dimensions
	view := updatedModel.View()

	// Overview mode shows top 5 todos per quadrant (Story 017)
	// We have 15 todos, so should see "... and 10 more"
	is.True(strings.Contains(stripANSI(view), "... and 10 more")) // expected to see '... and 10 more' in overview mode

	// Verify matrix content is present
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain matrix content
}

func TestStory006_MatrixAdjustsToDifferentSizes(t *testing.T) {
	// Scenario: Matrix adjusts to different terminal sizes
	// Given I have a todo.txt file
	// When I run the application in a 200x60 terminal
	// Then the quadrants are larger than in a 120x40 terminal
	is := is.New(t)

	input := generateManyTodos(30)
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate very large terminal (200x60)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 200, Height: 60})

	view := updatedModel.View()

	// Overview mode shows top 5 todos per quadrant (Story 017)
	// We have 30 todos, so should see "... and 25 more"
	is.True(strings.Contains(stripANSI(view), "... and 25 more")) // expected to see '... and 25 more' in overview mode
}

func TestStory006_MatrixRespectsMinimumDimensions(t *testing.T) {
	// Scenario: Matrix respects minimum dimensions
	// Given I have a todo.txt file
	// When I run the application in a very small terminal (80x24)
	// Then the matrix uses minimum viable dimensions
	// And does not break the layout
	is := is.New(t)

	input := "(A) Test todo"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Simulate small terminal (80x24)
	updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	view := updatedModel.View()

	// Matrix should still render without errors
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain matrix content even in small terminal

	// Should show the test todo
	is.True(strings.Contains(stripANSI(view), "Test todo")) // expected view to contain test todo
}

func TestStory006_DisplayLimitScalesWithHeight(t *testing.T) {
	// Scenario: Todo display limit scales with height
	// Given I have 20 todos in the Do First quadrant
	// When the quadrant height varies
	// Then the number of visible todos scales accordingly

	tests := []struct {
		name            string
		terminalHeight  int
		expectedMessage string
	}{
		{
			name:            "small terminal shows top 5 todos",
			terminalHeight:  30,
			expectedMessage: "... and 15 more", // Overview always shows top 5
		},
		{
			name:            "large terminal shows top 5 todos",
			terminalHeight:  50,
			expectedMessage: "... and 15 more", // Overview always shows top 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			input := generateManyTodos(20)
			repository := memory.NewRepository(input)

			m, err := usecases.LoadMatrix(repository)
			is.NoErr(err)

			model := ui.NewModel(m, "test.txt")
			updatedModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: tt.terminalHeight})

			view := updatedModel.View()

			is.True(strings.Contains(stripANSI(view), tt.expectedMessage)) // expected to see message
		})
	}
}

func TestStory006_WindowResizeHandledGracefully(t *testing.T) {
	// Scenario: Window resize is handled gracefully
	// Given the application is running
	// When I resize my terminal window
	// Then the matrix re-renders with new dimensions
	is := is.New(t)

	input := generateManyTodos(15)
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Start with medium terminal
	model1, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	view1 := model1.View()

	// Resize to larger terminal
	model2, _ := model1.Update(tea.WindowSizeMsg{Width: 150, Height: 50})
	view2 := model2.View()

	// Views should be different (different dimensions)
	is.True(view1 != view2) // expected views to be different after resize

	// Both should contain matrix content
	is.True(strings.Contains(view1, "Do First") && strings.Contains(view2, "Do First")) // expected both views to contain matrix content
}

func TestStory006_DefaultDimensionsWhenNoWindowSize(t *testing.T) {
	// Verify that default dimensions are used when no window size is received
	// (edge case for initial render before WindowSizeMsg)
	is := is.New(t)

	input := generateManyTodos(10)
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	model := ui.NewModel(m, "test.txt")

	// Call View without setting window size (width and height are 0)
	view := model.View()

	// Overview mode shows top 5, so with 10 todos should see "... and 5 more"
	is.True(strings.Contains(stripANSI(view), "... and 5 more")) // expected overview to show top 5 todos (... and 5 more)

	// Should still contain matrix content
	is.True(strings.Contains(stripANSI(view), "Do First")) // expected view to contain matrix content
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
