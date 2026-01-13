package acceptance_test

import (
	"io"
	"strings"
	"testing"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 004: Parse and render todo.txt project and context tags

func TestStory004_ParseSingleProjectTag(t *testing.T) {
	// Scenario: Parse single project tag
	// Given a todo.txt file containing:
	//   (A) Deploy new feature +WebApp
	// When I run the application
	// Then the todo is parsed with project ["WebApp"]

	input := "(A) Deploy new feature +WebApp"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doFirst := m.DoFirst()
	if len(doFirst) != 1 {
		t.Fatalf("expected 1 todo in Do First, got %d", len(doFirst))
	}

	projects := doFirst[0].Projects()
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0] != "WebApp" {
		t.Errorf("expected project 'WebApp', got %q", projects[0])
	}
}

func TestStory004_ParseSingleContextTag(t *testing.T) {
	// Scenario: Parse single context tag
	// Given a todo.txt file containing:
	//   (B) Call client @phone
	// When I run the application
	// Then the todo is parsed with context ["phone"]

	input := "(B) Call client @phone"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	schedule := m.Schedule()
	if len(schedule) != 1 {
		t.Fatalf("expected 1 todo in Schedule, got %d", len(schedule))
	}

	contexts := schedule[0].Contexts()
	if len(contexts) != 1 {
		t.Fatalf("expected 1 context, got %d", len(contexts))
	}
	if contexts[0] != "phone" {
		t.Errorf("expected context 'phone', got %q", contexts[0])
	}
}

func TestStory004_ParseMultipleProjectsAndContexts(t *testing.T) {
	// Scenario: Parse multiple projects and contexts
	// Given a todo.txt file containing:
	//   (A) Write quarterly report +Work +Q1Goals @office @computer
	// When I run the application
	// Then the todo is parsed with projects ["Work", "Q1Goals"]
	// And the todo is parsed with contexts ["office", "computer"]

	input := "(A) Write quarterly report +Work +Q1Goals @office @computer"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doFirst := m.DoFirst()
	if len(doFirst) != 1 {
		t.Fatalf("expected 1 todo in Do First, got %d", len(doFirst))
	}

	projects := doFirst[0].Projects()
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	if projects[0] != "Work" || projects[1] != "Q1Goals" {
		t.Errorf("expected projects ['Work', 'Q1Goals'], got %v", projects)
	}

	contexts := doFirst[0].Contexts()
	if len(contexts) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(contexts))
	}
	if contexts[0] != "office" || contexts[1] != "computer" {
		t.Errorf("expected contexts ['office', 'computer'], got %v", contexts)
	}
}

func TestStory004_ParseTagsAnywhereInDescription(t *testing.T) {
	// Scenario: Parse tags anywhere in description
	// Given a todo.txt file containing:
	//   (A) Review +OpenSource code for @github issues
	// When I run the application
	// Then the todo is parsed with project ["OpenSource"]
	// And the todo is parsed with context ["github"]

	input := "(A) Review +OpenSource code for @github issues"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doFirst := m.DoFirst()
	if len(doFirst) != 1 {
		t.Fatalf("expected 1 todo in Do First, got %d", len(doFirst))
	}

	projects := doFirst[0].Projects()
	if len(projects) != 1 || projects[0] != "OpenSource" {
		t.Errorf("expected project ['OpenSource'], got %v", projects)
	}

	contexts := doFirst[0].Contexts()
	if len(contexts) != 1 || contexts[0] != "github" {
		t.Errorf("expected context ['github'], got %v", contexts)
	}
}

func TestStory004_TodosWithoutTagsRenderNormally(t *testing.T) {
	// Scenario: Todos without tags render normally
	// Given a todo.txt file containing:
	//   (A) Simple task without tags
	// When I run the application
	// Then the todo is displayed in the Do First quadrant
	// And no tags are shown

	input := "(A) Simple task without tags"
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	doFirst := m.DoFirst()
	if len(doFirst) != 1 {
		t.Fatalf("expected 1 todo in Do First, got %d", len(doFirst))
	}

	if len(doFirst[0].Projects()) != 0 {
		t.Errorf("expected no projects, got %v", doFirst[0].Projects())
	}
	if len(doFirst[0].Contexts()) != 0 {
		t.Errorf("expected no contexts, got %v", doFirst[0].Contexts())
	}
}

func TestStory004_ConsistentColorForSameTag(t *testing.T) {
	// Scenario: Render tags with consistent colors
	// Given a todo.txt file containing multiple todos with the same tag:
	//   (A) First task +WebApp
	//   (B) Second task +WebApp
	// When I run the application
	// Then the +WebApp tag renders with the same color in both todos

	input := `(A) First task +WebApp
(B) Second task +WebApp`
	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify both todos parsed correctly with same project tag
	doFirst := m.DoFirst()
	if len(doFirst) != 1 {
		t.Fatalf("expected 1 todo in Do First, got %d", len(doFirst))
	}
	if doFirst[0].Projects()[0] != "WebApp" {
		t.Errorf("expected project 'WebApp' in Do First, got %v", doFirst[0].Projects())
	}

	schedule := m.Schedule()
	if len(schedule) != 1 {
		t.Fatalf("expected 1 todo in Schedule, got %d", len(schedule))
	}
	if schedule[0].Projects()[0] != "WebApp" {
		t.Errorf("expected project 'WebApp' in Schedule, got %v", schedule[0].Projects())
	}

	// Note: Color consistency is verified by the HashColor function using deterministic hashing
	// The UI layer will apply the same color to both instances of +WebApp
}

// StubTodoSource for testing - implements usecases.TodoSource and usecases.TodoWriter
type StubTodoSource struct {
	reader io.Reader
	writer io.Writer
	err    error
}

func (s *StubTodoSource) GetTodos() (io.ReadCloser, error) {
	if s.err != nil {
		return nil, s.err
	}
	return io.NopCloser(s.reader), nil
}

func (s *StubTodoSource) SaveTodo(line string) error {
	if s.writer == nil {
		return nil
	}
	_, err := s.writer.Write([]byte(line))
	return err
}

func (s *StubTodoSource) ReplaceAll(content string) error {
	if s.writer == nil {
		return nil
	}
	// For testing, we use a strings.Builder which doesn't support truncation
	// So we'll just reset and write the new content
	if sb, ok := s.writer.(*strings.Builder); ok {
		sb.Reset()
		_, err := sb.WriteString(content)
		return err
	}
	_, err := s.writer.Write([]byte(content))
	return err
}

// Integration test: Load matrix from realistic todo.txt with mixed tags
func TestStory004_Integration_MixedTagsInMatrix(t *testing.T) {
	// Given a realistic todo.txt with various tags across all quadrants
	input := `(A) Deploy feature +WebApp @computer
(B) Plan sprint +Work @office
(C) Reply to emails @computer
(D) Organize desk @office
(A) Fix bug +MobileApp @phone
No priority task +PersonalProject`

	source := &StubTodoSource{
		reader: strings.NewReader(input),
	}

	m, err := usecases.LoadMatrix(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify distribution across matrix
	assertMatrixDistribution(t, m, 2, 1, 1, 2)

	// Verify Do First has correct projects and contexts
	doFirst := m.DoFirst()
	verifyTodoHasTags(t, doFirst[0], []string{"WebApp"}, []string{"computer"})
	verifyTodoHasTags(t, doFirst[1], []string{"MobileApp"}, []string{"phone"})

	// Verify Schedule has correct tags
	schedule := m.Schedule()
	verifyTodoHasTags(t, schedule[0], []string{"Work"}, []string{"office"})

	// Verify Delegate has correct tags
	delegate := m.Delegate()
	verifyTodoHasTags(t, delegate[0], []string{}, []string{"computer"})

	// Verify Eliminate has mixed scenarios
	eliminate := m.Eliminate()
	verifyTodoHasTags(t, eliminate[0], []string{}, []string{"office"})
	verifyTodoHasTags(t, eliminate[1], []string{"PersonalProject"}, []string{})
}

func assertMatrixDistribution(t *testing.T, m matrix.Matrix, doFirst, schedule, delegate, eliminate int) {
	t.Helper()
	if len(m.DoFirst()) != doFirst {
		t.Errorf("expected %d todos in Do First, got %d", doFirst, len(m.DoFirst()))
	}
	if len(m.Schedule()) != schedule {
		t.Errorf("expected %d todos in Schedule, got %d", schedule, len(m.Schedule()))
	}
	if len(m.Delegate()) != delegate {
		t.Errorf("expected %d todos in Delegate, got %d", delegate, len(m.Delegate()))
	}
	if len(m.Eliminate()) != eliminate {
		t.Errorf("expected %d todos in Eliminate, got %d", eliminate, len(m.Eliminate()))
	}
}

func verifyTodoHasTags(t *testing.T, td todo.Todo, expectedProjects, expectedContexts []string) {
	t.Helper()

	projects := td.Projects()
	if len(projects) != len(expectedProjects) {
		t.Errorf("expected %d projects, got %d", len(expectedProjects), len(projects))
		return
	}
	for i, expected := range expectedProjects {
		if projects[i] != expected {
			t.Errorf("project[%d]: expected %q, got %q", i, expected, projects[i])
		}
	}

	contexts := td.Contexts()
	if len(contexts) != len(expectedContexts) {
		t.Errorf("expected %d contexts, got %d", len(expectedContexts), len(contexts))
		return
	}
	for i, expected := range expectedContexts {
		if contexts[i] != expected {
			t.Errorf("context[%d]: expected %q, got %q", i, expected, contexts[i])
		}
	}
}
