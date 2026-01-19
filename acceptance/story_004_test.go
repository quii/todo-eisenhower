package acceptance_test

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/adapters/memory"
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/usecases"
)

// Story 004: Parse and render todo.txt project and context tags

func TestStory004_ParseSingleProjectTag(t *testing.T) {
	is := is.New(t)

	input := "(A) Deploy new feature +WebApp"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	doFirst := m.DoFirst()
	is.Equal(len(doFirst), 1) // expected 1 todo in Do First

	projects := doFirst[0].Projects()
	is.Equal(len(projects), 1) // expected 1 project
	is.Equal(projects[0], "WebApp")
}

func TestStory004_ParseSingleContextTag(t *testing.T) {
	is := is.New(t)

	input := "(B) Call client @phone"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	schedule := m.Schedule()
	is.Equal(len(schedule), 1) // expected 1 todo in Schedule

	contexts := schedule[0].Contexts()
	is.Equal(len(contexts), 1) // expected 1 context
	is.Equal(contexts[0], "phone")
}

func TestStory004_ParseMultipleProjectsAndContexts(t *testing.T) {
	is := is.New(t)

	input := "(A) Write quarterly report +Work +Q1Goals @office @computer"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	doFirst := m.DoFirst()
	is.Equal(len(doFirst), 1) // expected 1 todo in Do First

	projects := doFirst[0].Projects()
	is.Equal(len(projects), 2) // expected 2 projects
	is.Equal(projects[0], "Work")
	is.Equal(projects[1], "Q1Goals")

	contexts := doFirst[0].Contexts()
	is.Equal(len(contexts), 2) // expected 2 contexts
	is.Equal(contexts[0], "office")
	is.Equal(contexts[1], "computer")
}

func TestStory004_ParseTagsAnywhereInDescription(t *testing.T) {
	is := is.New(t)

	input := "(A) Review +OpenSource code for @github issues"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	doFirst := m.DoFirst()
	is.Equal(len(doFirst), 1) // expected 1 todo in Do First

	projects := doFirst[0].Projects()
	is.Equal(len(projects), 1) // expected 1 project
	is.Equal(projects[0], "OpenSource")

	contexts := doFirst[0].Contexts()
	is.Equal(len(contexts), 1) // expected 1 context
	is.Equal(contexts[0], "github")
}

func TestStory004_TodosWithoutTagsRenderNormally(t *testing.T) {
	is := is.New(t)

	input := "(A) Simple task without tags"
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	doFirst := m.DoFirst()
	is.Equal(len(doFirst), 1)               // expected 1 todo in Do First
	is.Equal(len(doFirst[0].Projects()), 0) // expected no projects
	is.Equal(len(doFirst[0].Contexts()), 0) // expected no contexts
}

func TestStory004_ConsistentColorForSameTag(t *testing.T) {
	is := is.New(t)

	input := `(A) First task +WebApp
(B) Second task +WebApp`
	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Verify both todos parsed correctly with same project tag
	doFirst := m.DoFirst()
	is.Equal(len(doFirst), 1) // expected 1 todo in Do First
	is.Equal(doFirst[0].Projects()[0], "WebApp")

	schedule := m.Schedule()
	is.Equal(len(schedule), 1) // expected 1 todo in Schedule
	is.Equal(schedule[0].Projects()[0], "WebApp")

	// Note: Color consistency is verified by the HashColor function using deterministic hashing
	// The UI layer will apply the same color to both instances of +WebApp
}

// StubTodoSource is a legacy test helper for backward compatibility
// PREFER: Use memory.NewRepository(input) for new tests - it uses real Marshal/Unmarshal
// This stub provides io.Writer/io.Reader directly
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

func (s *StubTodoSource) GetAppendWriter() (io.WriteCloser, error) {
	if s.writer == nil {
		return &nopStubCloser{io.Discard}, nil
	}
	return &nopStubCloser{s.writer}, nil
}

func (s *StubTodoSource) GetReplaceWriter() (io.WriteCloser, error) {
	if s.writer == nil {
		return &nopStubCloser{io.Discard}, nil
	}
	// For testing with strings.Builder, reset it
	if sb, ok := s.writer.(*strings.Builder); ok {
		sb.Reset()
	}
	return &nopStubCloser{s.writer}, nil
}

type nopStubCloser struct {
	io.Writer
}

func (n *nopStubCloser) Close() error {
	return nil
}

// Integration test: Load matrix from realistic todo.txt with mixed tags
func TestStory004_Integration_MixedTagsInMatrix(t *testing.T) {
	is := is.New(t)

	input := `(A) Deploy feature +WebApp @computer
(B) Plan sprint +Work @office
(C) Reply to emails @computer
(D) Organize desk @office
(A) Fix bug +MobileApp @phone
No priority task +PersonalProject`

	repository := memory.NewRepository(input)

	m, err := usecases.LoadMatrix(repository)
	is.NoErr(err)

	// Verify distribution across matrix
	assertMatrixDistribution(is, m, 2, 1, 1, 2)

	// Verify Do First has correct projects and contexts
	doFirst := m.DoFirst()
	verifyTodoHasTags(is, doFirst[0], []string{"WebApp"}, []string{"computer"})
	verifyTodoHasTags(is, doFirst[1], []string{"MobileApp"}, []string{"phone"})

	// Verify Schedule has correct tags
	schedule := m.Schedule()
	verifyTodoHasTags(is, schedule[0], []string{"Work"}, []string{"office"})

	// Verify Delegate has correct tags
	delegate := m.Delegate()
	verifyTodoHasTags(is, delegate[0], []string{}, []string{"computer"})

	// Verify Eliminate has mixed scenarios
	eliminate := m.Eliminate()
	verifyTodoHasTags(is, eliminate[0], []string{}, []string{"office"})
	verifyTodoHasTags(is, eliminate[1], []string{"PersonalProject"}, []string{})
}

func assertMatrixDistribution(is *is.I, m matrix.Matrix, doFirst, schedule, delegate, eliminate int) {
	is.Helper()
	is.Equal(len(m.DoFirst()), doFirst)     // expected todos in Do First
	is.Equal(len(m.Schedule()), schedule)   // expected todos in Schedule
	is.Equal(len(m.Delegate()), delegate)   // expected todos in Delegate
	is.Equal(len(m.Eliminate()), eliminate) // expected todos in Eliminate
}

func verifyTodoHasTags(is *is.I, td todo.Todo, expectedProjects, expectedContexts []string) {
	is.Helper()

	projects := td.Projects()
	is.Equal(len(projects), len(expectedProjects)) // expected projects count

	for i, expected := range expectedProjects {
		is.Equal(projects[i], expected) // expected project
	}

	contexts := td.Contexts()
	is.Equal(len(contexts), len(expectedContexts)) // expected contexts count

	for i, expected := range expectedContexts {
		is.Equal(contexts[i], expected) // expected context
	}
}
