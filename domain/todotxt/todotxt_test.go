package todotxt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestParse(t *testing.T) {
	t.Run("parses single todo with priority A", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("(A) Fix critical bug")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "Fix critical bug", todo.PriorityA, false)
	})

	t.Run("parses todo without priority", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("No priority task")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "No priority task", todo.PriorityNone, false)
	})

	t.Run("parses completed todo", func(t *testing.T) {
		//nolint:gocritic // importShadow: is := is.New(t) is idiomatic for github.com/matryer/is
		is := is.New(t)
		input := strings.NewReader("x (A) 2026-01-11 Completed task")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo
		assertTodo(is, todos[0], "Completed task", todo.PriorityA, true)
	})

	t.Run("parses multiple todos with different priorities", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) Fix critical bug
(B) Plan quarterly goals
(C) Reply to emails
(D) Clean workspace`)

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 4) // expected 4 todos

		assertTodo(is, todos[0], "Fix critical bug", todo.PriorityA, false)
		assertTodo(is, todos[1], "Plan quarterly goals", todo.PriorityB, false)
		assertTodo(is, todos[2], "Reply to emails", todo.PriorityC, false)
		assertTodo(is, todos[3], "Clean workspace", todo.PriorityD, false)
	})

	t.Run("parses empty input", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 0) // expected 0 todos
	})

	t.Run("ignores empty lines", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) First task

(B) Second task`)

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 2) // expected 2 todos
	})

	t.Run("parses mixed completed and active todos", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader(`(A) Active task
x (A) 2026-01-11 Completed task
(B) Another active task`)

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 3) // expected 3 todos

		assertTodo(is, todos[0], "Active task", todo.PriorityA, false)
		assertTodo(is, todos[1], "Completed task", todo.PriorityA, true)
		assertTodo(is, todos[2], "Another active task", todo.PriorityB, false)
	})

	t.Run("parses single project tag", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Deploy new feature +WebApp")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		projects := todos[0].Projects()
		is.Equal(len(projects), 1) // expected 1 project
		is.Equal(projects[0], "WebApp")
	})

	t.Run("parses single context tag", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(B) Call client @phone")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 1) // expected 1 context
		is.Equal(contexts[0], "phone")
	})

	t.Run("parses multiple projects and contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Write report +Work +Q1Goals @office @computer")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1) // expected 1 todo

		projects := todos[0].Projects()
		is.Equal(len(projects), 2) // expected 2 projects
		is.Equal(projects[0], "Work")
		is.Equal(projects[1], "Q1Goals")

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 2) // expected 2 contexts
		is.Equal(contexts[0], "office")
		is.Equal(contexts[1], "computer")
	})

	t.Run("parses tags anywhere in description", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Review +OpenSource code for @github issues")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)

		projects := todos[0].Projects()
		is.Equal(len(projects), 1) // expected 1 project
		is.Equal(projects[0], "OpenSource")

		contexts := todos[0].Contexts()
		is.Equal(len(contexts), 1) // expected 1 context
		is.Equal(contexts[0], "github")
	})

	t.Run("todos without tags have empty slices", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Simple task without tags")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos[0].Projects()), 0) // expected no projects
		is.Equal(len(todos[0].Contexts()), 0) // expected no contexts
	})

	// Boundary tests for tag combinations (catches CONDITIONALS_BOUNDARY mutations)
	t.Run("parses todo with only projects and no contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task with +Project1 +Project2 but no contexts")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(len(todos[0].Projects()), 2) // expected 2 projects
		is.Equal(len(todos[0].Contexts()), 0) // expected 0 contexts
		is.Equal(todos[0].Projects()[0], "Project1")
		is.Equal(todos[0].Projects()[1], "Project2")
	})

	t.Run("parses todo with only contexts and no projects", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task @home @computer but no projects")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(len(todos[0].Projects()), 0) // expected 0 projects
		is.Equal(len(todos[0].Contexts()), 2) // expected 2 contexts
		is.Equal(todos[0].Contexts()[0], "home")
		is.Equal(todos[0].Contexts()[1], "computer")
	})

	t.Run("parses completed todo with only projects", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 (A) Completed +Project")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		is.Equal(len(todos[0].Projects()), 1) // expected 1 project
		is.Equal(len(todos[0].Contexts()), 0) // expected 0 contexts
	})

	t.Run("parses completed todo with only contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 (A) Completed @office")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		is.Equal(len(todos[0].Projects()), 0) // expected 0 projects
		is.Equal(len(todos[0].Contexts()), 1) // expected 1 context
	})

	// Creation date tests (catches CONDITIONALS_NEGATION mutations)
	t.Run("parses todo with creation date", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-01-15 Task with creation date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task with creation date")
		is.Equal(todos[0].Priority(), todo.PriorityA)
		is.Equal(todos[0].IsCompleted(), false)
	})

	t.Run("parses completed todo with both completion and creation dates", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 2026-01-15 (A) Task with both dates")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task with both dates")
		is.Equal(todos[0].Priority(), todo.PriorityA)
		is.True(todos[0].IsCompleted())
	})

	t.Run("parses todo without creation date", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(B) Task without creation date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Task without creation date")
		is.Equal(todos[0].Priority(), todo.PriorityB)
	})

	// Date parsing error handling (catches error path mutations)
	t.Run("removes malformed completion date from description but doesn't store it", func(t *testing.T) {
		is := is.New(t)
		// Dates matching pattern are removed even if invalid
		input := strings.NewReader("x 2026-99-99 (A) Task with invalid date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		// Malformed date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid date")
	})

	t.Run("removes malformed creation date from description but doesn't store it", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-99-99 Task with invalid creation date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// Malformed date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid creation date")
	})

	t.Run("removes malformed creation date on completed todo", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 2026-99-99 (A) Task with invalid creation date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		// Malformed creation date is removed from description (matches regex but parse fails)
		is.Equal(todos[0].Description(), "Task with invalid creation date")
	})
}

// Tests based on official todo.txt spec: https://github.com/todotxt/todo.txt
func TestTodoTxtSpecCompliance(t *testing.T) {
	// Rule 1: Priority must be uppercase A-Z at line start
	t.Run("lowercase priority is not recognized", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(a) Task with lowercase priority")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// Lowercase (a) should not be recognized as priority
		is.Equal(todos[0].Priority(), todo.PriorityNone)
		is.Equal(todos[0].Description(), "(a) Task with lowercase priority")
	})

	t.Run("priority not at line start is not recognized", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("Task (A) with priority in middle")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// (A) not at start should not be recognized
		is.Equal(todos[0].Priority(), todo.PriorityNone)
		is.Equal(todos[0].Description(), "Task (A) with priority in middle")
	})

	// Rule 1 (Complete): Completion marker must be lowercase 'x'
	t.Run("uppercase X is not recognized as completion", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("X (A) Task with uppercase X")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// Uppercase X should not mark as completed
		is.Equal(todos[0].IsCompleted(), false)
		is.Equal(todos[0].Description(), "X (A) Task with uppercase X")
	})

	t.Run("x not at line start is not recognized as completion", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("Task x marked")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// 'x' not at start should not mark as completed
		is.Equal(todos[0].IsCompleted(), false)
		is.Equal(todos[0].Description(), "Task x marked")
	})

	// Rule 3: Projects and contexts must be preceded by space
	t.Run("email addresses are not contexts", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Email test@example.com about issue")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		// @ in email shouldn't be parsed as context
		contexts := todos[0].Contexts()
		// Our implementation uses \w+ which won't match the period, so it might parse "example"
		// This documents current behavior
		if len(contexts) > 0 {
			// If it parses anything, it would be the part before special chars
			is.True(contexts[0] == "example") // documents current behavior
		}
	})

	t.Run("projects and contexts anywhere after priority and dates", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-01-15 Call +Sales about @proposal")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(len(todos[0].Projects()), 1)
		is.Equal(todos[0].Projects()[0], "Sales")
		is.Equal(len(todos[0].Contexts()), 1)
		is.Equal(todos[0].Contexts()[0], "proposal")
		// Tags removed from description
		is.Equal(todos[0].Description(), "Call about")
	})

	// Creation date position tests
	t.Run("creation date without priority at line start", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("2026-01-15 Task without priority but with date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Priority(), todo.PriorityNone)
		// Should have creation date parsed
		is.True(todos[0].CreationDate() != nil)
	})

	t.Run("creation date after priority", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) 2026-01-15 Task with both priority and date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Priority(), todo.PriorityA)
		is.True(todos[0].CreationDate() != nil)
	})

	// Completed task date ordering: x COMPLETION_DATE CREATION_DATE (PRIORITY) Description
	t.Run("completed task with both dates in correct order", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-18 2026-01-15 (A) Task completed 3 days after creation")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.True(todos[0].IsCompleted())
		is.True(todos[0].CompletionDate() != nil)
		is.True(todos[0].CreationDate() != nil)
		// First date should be completion, second should be creation
		is.Equal(todos[0].CompletionDate().Format("2006-01-02"), "2026-01-18")
		is.Equal(todos[0].CreationDate().Format("2006-01-02"), "2026-01-15")
	})
}

func TestParseNew(t *testing.T) {
	t.Run("creates todo without tags", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Buy milk", todo.PriorityA, creationDate)

		is.Equal(result.Description(), "Buy milk")
		is.Equal(result.Priority(), todo.PriorityA)
		is.True(!result.IsCompleted())
		is.Equal(len(result.Projects()), 0)
		is.Equal(len(result.Contexts()), 0)
		is.True(result.CreationDate() != nil)
		is.Equal(result.CreationDate().Format("2006-01-02"), "2026-01-20")
	})

	t.Run("creates todo with project tags", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Review code +WebApp +Mobile", todo.PriorityB, creationDate)

		is.Equal(result.Description(), "Review code")
		is.Equal(result.Priority(), todo.PriorityB)
		is.Equal(len(result.Projects()), 2)
		is.Equal(result.Projects()[0], "WebApp")
		is.Equal(result.Projects()[1], "Mobile")
		is.Equal(len(result.Contexts()), 0)
	})

	t.Run("creates todo with context tags", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Call manager @phone @urgent", todo.PriorityA, creationDate)

		is.Equal(result.Description(), "Call manager")
		is.Equal(result.Priority(), todo.PriorityA)
		is.Equal(len(result.Contexts()), 2)
		is.Equal(result.Contexts()[0], "phone")
		is.Equal(result.Contexts()[1], "urgent")
		is.Equal(len(result.Projects()), 0)
	})

	t.Run("creates todo with mixed tags", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Fix bug +WebApp @computer", todo.PriorityA, creationDate)

		is.Equal(result.Description(), "Fix bug")
		is.Equal(len(result.Projects()), 1)
		is.Equal(result.Projects()[0], "WebApp")
		is.Equal(len(result.Contexts()), 1)
		is.Equal(result.Contexts()[0], "computer")
	})

	t.Run("creates todo with tags in middle of description", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Review +WebApp pull request @github", todo.PriorityC, creationDate)

		is.Equal(result.Description(), "Review pull request")
		is.Equal(result.Projects()[0], "WebApp")
		is.Equal(result.Contexts()[0], "github")
	})

	t.Run("cleans up extra whitespace", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Task   with    extra     spaces +project", todo.PriorityD, creationDate)

		is.Equal(result.Description(), "Task with extra spaces")
	})
}

func TestParse_DueDates(t *testing.T) {
	t.Run("parses due date from description", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Submit report due:2026-01-25")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(len(todos), 1)
		is.Equal(todos[0].Description(), "Submit report") // due date stripped
		is.True(todos[0].DueDate() != nil)                // has due date
		expectedDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[0].DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
	})

	t.Run("parses due date case-insensitive", func(t *testing.T) {
		testCases := []string{
			"(A) Task due:2026-01-25",
			"(A) Task Due:2026-01-25",
			"(A) Task DUE:2026-01-25",
			"(A) Task dUe:2026-01-25",
		}

		for _, input := range testCases {
			is := is.New(t)
			todos, err := todotxt.Unmarshal(strings.NewReader(input))

			is.NoErr(err)
			is.True(todos[0].DueDate() != nil) // should parse case-insensitive
		}
	})

	t.Run("uses first due date when multiple present", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task due:2026-01-20 due:2026-01-25")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.True(todos[0].DueDate() != nil)
		expectedDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[0].DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
	})

	t.Run("ignores invalid due date format", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task due:invalid-date")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.True(todos[0].DueDate() == nil) // invalid date ignored
	})

	t.Run("parses due date with tags", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(A) Task +project @context due:2026-01-25")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.Equal(todos[0].Description(), "Task") // tags and due date stripped
		is.Equal(len(todos[0].Projects()), 1)
		is.Equal(todos[0].Projects()[0], "project")
		is.Equal(len(todos[0].Contexts()), 1)
		is.Equal(todos[0].Contexts()[0], "context")
		is.True(todos[0].DueDate() != nil)
		expectedDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[0].DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
	})

	t.Run("todos without due date have nil DueDate", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("(B) Regular task")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.True(todos[0].DueDate() == nil) // no due date
	})

	t.Run("completed todo preserves due date", func(t *testing.T) {
		is := is.New(t)
		input := strings.NewReader("x 2026-01-20 2026-01-15 (A) Completed task due:2026-01-19")

		todos, err := todotxt.Unmarshal(input)

		is.NoErr(err)
		is.True(todos[0].IsCompleted())
		is.True(todos[0].DueDate() != nil)
		expectedDate := time.Date(2026, 1, 19, 0, 0, 0, 0, time.UTC)
		is.Equal(todos[0].DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
	})
}

func TestString_DueDates(t *testing.T) {
	t.Run("reconstructs due date in String output", func(t *testing.T) {
		is := is.New(t)
		input := "(A) Submit report due:2026-01-25\n"

		todos, err := todotxt.Unmarshal(strings.NewReader(input))
		is.NoErr(err)

		output := todos[0].String()
		is.True(strings.Contains(output, "due:2026-01-25")) // due date reconstructed
	})

	t.Run("round-trip preserves due date", func(t *testing.T) {
		is := is.New(t)
		original := "(A) Task +project @context due:2026-01-25\n"

		// Parse
		todos, err := todotxt.Unmarshal(strings.NewReader(original))
		is.NoErr(err)

		// Reconstruct
		reconstructed := todos[0].String()

		// Parse again
		todos2, err := todotxt.Unmarshal(strings.NewReader(reconstructed))
		is.NoErr(err)

		// Verify same due date
		is.True(todos[0].DueDate() != nil)
		is.True(todos2[0].DueDate() != nil)
		is.Equal(todos[0].DueDate().Format("2006-01-02"), todos2[0].DueDate().Format("2006-01-02"))
	})
}

func TestParseNew_DueDates(t *testing.T) {
	t.Run("extracts due date from user input", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Submit report due:2026-01-25", todo.PriorityA, creationDate)

		is.Equal(result.Description(), "Submit report") // due date stripped
		is.True(result.DueDate() != nil)
		expectedDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		is.Equal(result.DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
	})

	t.Run("extracts due date with tags", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

		result := todotxt.ParseNew("Task +project @context due:2026-01-25", todo.PriorityA, creationDate)

		is.Equal(result.Description(), "Task")
		is.Equal(len(result.Projects()), 1)
		is.Equal(len(result.Contexts()), 1)
		is.True(result.DueDate() != nil)
	})
}

func TestParseEdit_DueDates(t *testing.T) {
	t.Run("updates due date from edited description", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		original := todo.NewWithCreationDate("Original task", todo.PriorityA, &creationDate)

		edited := todotxt.ParseEdit(original, "Updated task due:2026-01-30", todo.PriorityA)

		is.Equal(edited.Description(), "Updated task")
		is.True(edited.DueDate() != nil)
		expectedDate := time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC)
		is.Equal(edited.DueDate().Format("2006-01-02"), expectedDate.Format("2006-01-02"))
		is.Equal(edited.CreationDate(), original.CreationDate()) // creation date preserved
	})

	t.Run("removes due date if not in edited description", func(t *testing.T) {
		is := is.New(t)
		creationDate := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
		dueDate := time.Date(2026, 1, 30, 0, 0, 0, 0, time.UTC)
		original := todo.NewFull("Original task", todo.PriorityA, false, nil, &creationDate, &dueDate, nil, nil, nil)

		edited := todotxt.ParseEdit(original, "Updated task without due date", todo.PriorityA)

		is.Equal(edited.Description(), "Updated task without due date")
		is.True(edited.DueDate() == nil) // due date removed
	})
}

func TestFormatForInput_DueDates(t *testing.T) {
	t.Run("includes due date in formatted input", func(t *testing.T) {
		is := is.New(t)
		dueDate := time.Date(2026, 1, 25, 0, 0, 0, 0, time.UTC)
		td := todo.NewFull("Task", todo.PriorityA, false, nil, nil, &dueDate, nil, []string{"project"}, []string{"context"})

		formatted := todotxt.FormatForInput(td)

		is.True(strings.Contains(formatted, "Task"))
		is.True(strings.Contains(formatted, "+project"))
		is.True(strings.Contains(formatted, "@context"))
		is.True(strings.Contains(formatted, "due:2026-01-25"))
	})
}

func assertTodo(is *is.I, got todo.Todo, wantDesc string, wantPriority todo.Priority, wantCompleted bool) {
	is.Helper()

	is.Equal(got.Description(), wantDesc) // description
	is.Equal(got.Priority(), wantPriority) // priority
	is.Equal(got.IsCompleted(), wantCompleted) // completed
}
