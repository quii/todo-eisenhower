# CLAUDE.md - Project Conventions & Architecture

## Overview
A Go TUI application for visualizing todo.txt files as an Eisenhower matrix using Charm's Bubble Tea framework.

## Architecture Philosophy

### Domain-Driven Design (DDD) with Ports & Adapters
We follow Hexagonal Architecture (Ports and Adapters), which is a key part of DDD. The domain is at the center, with ports (interfaces) defining how to interact with the outside world, and adapters implementing those ports.

**Ports**: Interfaces the domain/use cases need (e.g., `TodoSource.GetTodos()`)
**Adapters**: Concrete implementations (e.g., file system, HTTP, database, UI frameworks)

**When to create adapters:**
- ✅ When you have a port (interface) and need to implement it
  - Example: `TodoSource` port → `file.Source` adapter wraps `os.Open()`
- ❌ When you're just wrapping stdlib with no abstraction
  - Bad: Creating a wrapper around `os.Open()` that returns `io.ReadCloser` with no interface
- ✅ When stdlib already implements the port you need
  - Example: Parser works with `io.Reader` - `os.Open()` already provides this, no adapter needed

**The nuance**: If you have a port (interface), create an adapter to implement it, even if it's a thin wrapper. If you don't have a port, don't wrap unnecessarily.

We organize code around three bounded contexts:

1. **Todo Domain** (`domain/todo/`)
   - Todo entity and business rules
   - Priority, completion state, todo.txt format concerns
   - No knowledge of matrices or UI

2. **Matrix Domain** (`domain/matrix/`)
   - Eisenhower matrix logic (quadrant categorization)
   - Depends on Todo domain
   - No knowledge of UI or file formats

3. **Parser Domain** (`domain/parser/`)
   - todo.txt file format parsing/serialization
   - Converts between file format and Todo entities
   - Works with `io.Reader` - not coupled to files or any specific input source
   - No business logic, no I/O concerns

### Use Cases (`usecases/`)
- Orchestrate domain objects to fulfill application requirements
- Example: `LoadTodoMatrix` use case coordinates Parser and Matrix domains
- Keep use cases thin - domain objects do the heavy lifting
- **Use cases are the application's API** - this is the boundary main.go calls

**Use cases represent user actions, not implementation details:**
- ✅ `AddTodo()` - user adds a new todo
- ✅ `ToggleCompletion()` - user marks a todo complete/incomplete
- ✅ `ChangePriority()` - user moves a todo to different quadrant
- ❌ `SaveAllTodos()` - implementation detail (private helper)
- ❌ `FormatTodo()` - implementation detail (private helper)

**Rule of thumb**: If you can infer it from a user story ("As a user, I want to..."), it's probably a usecase. If it's about "how" the system persists or formats data, it's a private helper.

**Usecase pattern:**
```go
// Usecase shows the business intent clearly
func AddTodo(writer TodoWriter, m matrix.Matrix, description string, priority Priority) (Matrix, error) {
    // 1. Use rich domain model
    newTodo := todo.New(description, priority)
    updatedMatrix := m.AddTodo(newTodo)

    // 2. Persist changes (private helper - implementation detail)
    err := saveTodo(writer, newTodo)
    if err != nil {
        return m, err // Return original on failure
    }

    return updatedMatrix, nil
}

// Private helper handles persistence details
func saveTodo(writer TodoWriter, t todo.Todo) error {
    line := FormatTodo(t)
    return writer.SaveTodo(line)
}
```

### The main.go Stability Rule
- **main.go should be stable and rarely change**
- main.go is just wiring - it calls use cases and sets up infrastructure
- If you find yourself changing main.go for new features, it's a smell
- Business logic changes should manifest as changes to use cases, not main
- Main should only change when:
  - Adding fundamentally new commands or modes
  - Changing infrastructure concerns (CLI flags, logging setup)
  - Fundamental architectural shifts

**Example:**
```go
// Good: main.go stays stable across stories
func main() {
    filePath := getFilePath() // might change for CLI args
    matrix, err := usecases.LoadMatrix(filePath)
    // ... run UI
}

// Story 001: usecases.LoadMatrix() returns hard-coded todos
// Story 002: usecases.LoadMatrix() reads from file
// Story 003: usecases.LoadMatrix() accepts parameter
// main.go doesn't need to change - the use case evolves
```

### Adapters (`adapters/`)
- **UI Adapter** (`adapters/ui/`): Bubble Tea components and views
- **File Adapter** (`adapters/file/`): Implements `TodoSource` port using filesystem
- Dependencies point inward toward domain/use cases (Dependency Inversion)
- Adapters implement ports defined in use cases/domain

## Testing Conventions

### Test-Driven Development (TDD)
- Write tests first, always
- Red → Green → Refactor cycle
- Tests document expected behavior

### Black-Box Testing
- All tests use `_test` packages (e.g., `package todo_test`)
- Tests interact with public API only, as a consumer would
- Forces good API design
- Example: `domain/todo/todo_test.go` has `package todo_test`

### Test Organization
```
domain/todo/
  todo.go           // implementation
  todo_test.go      // package todo_test
  export_test.go    // package todo - exports internals for testing if needed
```

### Testing Bubble Tea Components
- Separate **presentation logic** from **view rendering**
- Test presentation logic in black-box tests
- Bubble Tea's `Update()` and `View()` should be thin adapters
- Example: Extract matrix layout logic into testable functions

### Acceptance Testing
- Create acceptance tests in `acceptance/` that map to Gherkin scenarios
- Test use cases with stub implementations (e.g., `StubTodoSource`)
- Use `strings.NewReader` or `bytes.Buffer` instead of real filesystem
- Each test should clearly reference its Gherkin scenario
- Acceptance tests verify stories work end-to-end at the business logic level
- We accept some risk with UI (Bubble Tea) for now, but use case tests give confidence

### CRITICAL: Use Tests to Verify Behavior, Not Manual Scripts

**The architecture is designed for testability - use it.**

❌ **NEVER** resort to:
- Writing temporary Go scripts in /tmp to test parsing
- Creating one-off test files that aren't part of the test suite
- Manual testing without writing automated tests first
- Running the full TUI app just to verify a domain or parser change

✅ **ALWAYS**:
- Write unit tests to verify behavior
- Run existing tests with `go test ./domain/parser -v`
- Add new tests to the existing test suite
- Use `strings.NewReader()` for parser tests
- Trust the architecture - it's designed to make testing easy

**Why this matters:**
1. **Token efficiency**: Writing a test is faster than writing temporary scripts
2. **Permanent artifact**: Tests document expected behavior forever
3. **Regression prevention**: Tests prevent future breakage
4. **TDD compliance**: This is what TDD means - test first, always
5. **Architecture justification**: The entire hexagonal design is for testability

**Example - Testing parser changes:**
```go
// ❌ DON'T: Create /tmp/test_parser.go
// ✅ DO: Write a unit test in domain/parser/parser_test.go

func TestParse_CreationDates(t *testing.T) {
    input := "(A) 2026-01-10 Task created on Jan 10"
    todos, err := Parse(strings.NewReader(input))

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if cd := todos[0].CreationDate(); cd == nil {
        t.Error("expected creation date to be parsed")
    } else if cd.Format("2006-01-02") != "2026-01-10" {
        t.Errorf("expected date 2026-01-10, got %s", cd.Format("2006-01-02"))
    }
}
```

**If you need to verify something works:**
1. Look for existing tests you can run
2. Write a new test if none exists
3. Run tests with `go test ./path/to/package -v`
4. Use test output to verify behavior

**Remember**: If you're reaching for a temporary script, you've abandoned TDD.

**Example:**
```go
// acceptance/story_002_test.go
func TestStory002_LoadFromHardcodedPath(t *testing.T) {
    t.Run("Scenario: Load todos from hardcoded file path", func(t *testing.T) {
        source := StubTodoSource{data: "(A) Task\n"}
        m, err := usecases.LoadMatrix(source)
        // assertions...
    })
}
```

### Test Coverage
- Focus on behavior, not implementation details
- Don't test framework code (e.g., Bubble Tea internals)
- Test edge cases and error conditions
- Use table-driven tests for multiple scenarios

## Code Organization

```
.
├── domain/
│   ├── todo/          # Todo entities and business rules
│   ├── matrix/        # Eisenhower matrix logic
│   └── parser/        # todo.txt format handling
├── usecases/          # Application use cases (define ports)
├── adapters/
│   ├── file/          # File system adapter (implements TodoSource port)
│   └── ui/            # Bubble Tea TUI components
├── acceptance/        # Acceptance tests mapping to Gherkin scenarios
├── cmd/
│   └── eisenhower/    # Main entry point (wires up adapters)
└── CLAUDE.md          # This file
```

## Linting & Code Quality

### golangci-lint Configuration
- Use `.golangci.yml` in project root
- Enable linters: `gofmt`, `goimports`, `govet`, `staticcheck`, `errcheck`, `gosimple`, `ineffassign`
- Consider: `revive`, `unparam`, `unused`, `gocritic`
- Run before committing: `golangci-lint run`

### Code Style
- Follow standard Go conventions
- Use `gofmt` and `goimports`
- Avoid naked returns
- Handle errors explicitly, no `_` discards without reason
- Prefer small, focused functions

### Architecture Enforcement

**Automated enforcement via `architecture_test.go`:**

The project uses architecture tests to enforce hexagonal architecture boundaries. These tests run automatically with `go test ./...` and will fail CI if violated.

**Enforced rules:**

1. **Domain Layer Purity** (`TestDomainLayerPurity`)
   - ❌ No imports of `adapters` packages
   - ❌ No imports of `usecases` packages
   - ❌ No infrastructure imports: `os`, `net/http`, `database/sql`
   - ✅ Can import: stdlib (except infrastructure), other domain packages

2. **Use Case Layer Purity** (`TestUseCaseLayerPurity`)
   - ❌ No imports of `adapters` packages
   - ✅ Can import: stdlib, domain packages
   - ✅ Must define ports (interfaces) for external needs

**What happens when violated:**

```bash
$ go test
--- FAIL: TestDomainLayerPurity (0.00s)
    architecture_test.go:29: Domain layer architecture violations found:

          ❌ domain/todo/bad_code.go imports forbidden package 'os'

        Domain packages must not import infrastructure or outer layers.
        Define a port (interface) instead and let adapters implement it.
```

**Why this matters:**
- Prevents accidental coupling to infrastructure
- Enforces testability (domain is pure, easy to test)
- Makes code portable (swap adapters without changing domain)
- Self-documenting architecture (tests explain the rules)

## Bubble Tea / Lipgloss Guidelines

### Learning Approach
- Bubble Tea follows Elm architecture: Model → Update → View
- Start with simple models, iterate
- Lipgloss handles styling (colors, borders, layout)

### Testing Strategy
- **Model**: Pure data structure (easy to test)
- **Update**: Pure function `Update(msg) (Model, Cmd)` (test with different messages)
- **View**: String rendering (test formatting functions separately)

### Separation of Concerns
- Keep Bubble Tea code in `adapters/ui/`
- Don't let `tea.Model` leak into domain or use cases
- Adapter translates between domain types and UI messages

## Development Workflow

### Planning Work: Small Vertical Slices
- Work in **very small vertical slices** - the smallest shippable increment
- Each slice should take one focused session to complete
- Before implementing, ask clarifying questions to understand requirements
- Use **Gherkin** (Given/When/Then) to document and communicate requirements
- Store stories in `stories/` directory

**Example of breaking down work:**
- ❌ Too large: "Load and display todos from file with custom path"
- ✅ Better approach:
  1. Story 001: Display matrix with hard-coded in-memory todos
  2. Story 002: Load todos from hard-coded file path
  3. Story 003: Accept file path as command-line argument

**Why start with hard-coded data?**
- Establishes the full application stack (TUI, domain, architecture)
- Proves the Bubble Tea rendering works
- Creates foundation for incremental enhancement
- Each step adds exactly one new capability
- Always have a working, demonstrable application

### Trunk-Based Development
- Work directly on `main` branch (no feature branches, no pull requests)
- Each commit must be releasable
- Commit frequently with working, tested code
- Never commit broken code or failing tests
- Use feature flags or incomplete-but-safe code if needed
- Each push should pass all tests and linting

**Commit discipline:**
- Small commits (one logical change)
- Each commit leaves the codebase in a working state
- Write clear commit messages explaining the "why"
- Run tests and linter before every commit

### TDD Cycle
1. **Review story** - understand acceptance criteria (Gherkin scenarios)
2. **Write failing test** (black-box, `_test` package)
3. **Run test** - confirm it fails for the right reason
4. **Write minimal code** to pass the test
5. **Run linter** - `golangci-lint run`
6. **Refactor** - improve design while tests pass
7. **Verify acceptance criteria** - ensure story scenarios pass
8. **Commit** - working, tested code

## Dependencies

- **Bubble Tea**: TUI framework (`github.com/charmbracelet/bubbletea`)
- **Lipgloss**: Styling library (`github.com/charmbracelet/lipgloss`)
- Consider: `github.com/charmbracelet/bubbles` for reusable components

## Incremental Development

- Build features iteratively
- Start with simplest version (e.g., view-only before editing)
- Each iteration should be shippable
- Refactor toward extensibility as patterns emerge

## Common Pitfalls to Avoid

- ❌ Testing implementation details instead of behavior
- ❌ Coupling domain logic to UI framework
- ❌ Fat use cases with business logic (push to domain)
- ❌ Skipping tests "just this once"
- ❌ Mocking excessively - prefer real objects in tests

## Naming Conventions

- **Domains**: Singular noun (`todo`, not `todos`)
- **Use Cases**: Verb phrase (`LoadTodoMatrix`, `ParseTodoFile`)
- **Tests**: `Test<FunctionName>` or `Test<Scenario>`
- **Interfaces**: `-er` suffix when appropriate (`Parser`, `Loader`)

## Questions & Clarifications

When uncertain:
- Default to simpler solution
- Add complexity only when needed
- Consult Go proverbs and standard library for idiomatic patterns
- Remember: "A little copying is better than a little dependency"
