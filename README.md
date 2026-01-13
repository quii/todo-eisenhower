# Eisenhower Matrix TUI

A beautiful, fullscreen terminal-based Eisenhower matrix viewer for [todo.txt](http://todotxt.org/) files, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## What's an Eisenhower Matrix?

The Eisenhower Matrix helps prioritize tasks by organizing them into four quadrants:
- **DO FIRST** (Q1): Urgent & Important - Priority A tasks
- **SCHEDULE** (Q2): Important, Not Urgent - Priority B tasks
- **DELEGATE** (Q3): Urgent, Not Important - Priority C tasks
- **ELIMINATE** (Q4): Neither Urgent nor Important - Priority D and untagged tasks

## Features

**Fullscreen TUI** - Alt-screen mode that takes over your terminal, just like modern CLI tools

**Responsive Layout** - Automatically adapts to your terminal size, showing more todos in larger windows

**Colorized Tags** - Project (+tag) and context (@tag) tags are colorized consistently throughout

**todo.txt Format** - Full support for the todo.txt specification including:
- Priority markers: `(A)`, `(B)`, `(C)`, `(D)`
- Completion markers: `x` with date
- Project tags: `+ProjectName`
- Context tags: `@context`

**Beautiful Styling** - Color-coded quadrants and polished UI with [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Installation

### Download Pre-built Binary (Recommended)

Pre-built binaries are automatically built on every push to main. Download the latest build:

1. Go to the [Actions tab](https://github.com/yourusername/todo-eisenhower/actions)
2. Click on the most recent "CI" workflow run (green checkmark)
3. Scroll down to "Artifacts" section
4. Download the binary for your platform:
   - **Linux (x64)**: `eisenhower-linux-amd64`
   - **Linux (ARM64)**: `eisenhower-linux-arm64`
   - **macOS (Intel)**: `eisenhower-darwin-amd64`
   - **macOS (Apple Silicon)**: `eisenhower-darwin-arm64`
   - **Windows (x64)**: `eisenhower-windows-amd64.exe`

```bash
# After downloading, make it executable and move to PATH (macOS/Linux)
unzip eisenhower-darwin-arm64.zip  # GitHub wraps artifacts in zip
chmod +x eisenhower-darwin-arm64
sudo mv eisenhower-darwin-arm64 /usr/local/bin/eisenhower
```

Artifacts are kept for 90 days. For tagged releases, see the [Releases page](https://github.com/yourusername/todo-eisenhower/releases).

### Build from Source

Requires [Go 1.21+](https://golang.org/dl/)

```bash
# Clone the repository
git clone https://github.com/yourusername/todo-eisenhower.git
cd todo-eisenhower

# Build
go build -o eisenhower ./cmd/eisenhower/

# Optional: Install to your PATH
go install ./cmd/eisenhower
```

## Usage

```bash
# View your ~/todo.txt file
./eisenhower

# View a specific todo.txt file
./eisenhower /path/to/todo.txt

# Press 'q' or Ctrl+C to quit
```

### Example todo.txt

```
(A) Fix critical bug in authentication +WebApp @computer
(B) Plan quarterly goals +Work @office
(C) Reply to client emails @phone
(D) Organize workspace
x (A) 2026-01-11 Completed task +Project
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run linter (requires golangci-lint)
golangci-lint run
```

## Project Structure

This project follows Domain-Driven Design (DDD) with a hexagonal (ports & adapters) architecture:

```
.
├── domain/           # Domain layer (business logic)
│   ├── todo/        # Todo entities and value objects
│   ├── matrix/      # Matrix categorization logic
│   └── parser/      # todo.txt format parser
├── usecases/        # Application use cases (orchestration)
├── adapters/        # Adapters layer (external interfaces)
│   ├── ui/          # Bubble Tea TUI with Lipgloss styling
│   └── file/        # File system operations
├── cmd/
│   └── eisenhower/  # Main entry point (wiring)
├── acceptance/      # Acceptance tests (black-box)
└── stories/         # User stories (Gherkin format)
```

## Development

This project is developed using:
- **Test-Driven Development (TDD)** with black-box testing
- **Small vertical slices** - each story is a complete, releasable feature
- **Trunk-based development** - all commits go directly to main
- **Gherkin scenarios** for acceptance criteria

See [CLAUDE.md](./CLAUDE.md) for detailed development conventions and architecture guidelines.

## Roadmap

### Completed ✅
- [x] **Story 001**: Display matrix with hard-coded todos
- [x] **Story 002**: Load todos from hardcoded file path (~/todo.txt)
- [x] **Story 003**: Accept custom file path as CLI argument
- [x] **Story 004**: Parse and render todo.txt project (+tag) and context (@tag) tags with colors
- [x] **Story 005**: Fullscreen TUI with alt-screen mode
- [x] **Story 006**: Responsive matrix sizing that adapts to terminal dimensions
- [x] **Story 007**: Quadrant focus mode (press 1/2/3/4 to focus, ESC to return)
- [x] **Story 008**: Add todos with tag reference (press 'a' in focus mode)
- [x] **Story 009**: Tag autocomplete with + and @ triggers (arrow keys to navigate, Tab/Enter to complete)
- [x] **Story 010**: Tag inventory display showing counts by tag (sorted by count, overview mode only)

### Future Ideas 🚀
- Todo editing capabilities
- Mark todos as complete
- Delete todos
- Filtering by project or context
- Search functionality
- Due dates support
- Recurring tasks
- Multiple file support
