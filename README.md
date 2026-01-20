# Eisenhower Matrix TUI

A beautiful, fullscreen terminal-based Eisenhower matrix viewer for [todo.txt](http://todotxt.org/) files, built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Quickstart

Install via Homebrew:

```bash
brew install quii/tap/todo-eisenhower
```

Run with your todo.txt file:

```bash
eisenhower ~/todo.txt
```

Or use the default location (`~/todo.txt`):

```bash
eisenhower
```

Press `1`, `2`, `3`, or `4` to focus on a quadrant. Press `q` to quit.

## What's an Eisenhower Matrix?

The Eisenhower Matrix helps prioritize tasks by organizing them into four quadrants:
- **Do First** (Q1): Urgent & Important - Priority A tasks
- **Schedule** (Q2): Important, Not Urgent - Priority B tasks
- **Delegate** (Q3): Urgent, Not Important - Priority C tasks
- **Eliminate** (Q4): Neither Urgent nor Important - Priority D and untagged tasks

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

### Homebrew (macOS/Linux) - Recommended

```bash
brew install quii/tap/todo-eisenhower
```

### Download Pre-built Binary

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

## Upgrading

### Homebrew

```bash
brew upgrade quii/tap/todo-eisenhower
```

### Pre-built Binary or Source Build

Download the latest binary from the [Releases page](https://github.com/quii/todo-eisenhower/releases) or rebuild from source using the same steps as installation.

## Usage

```bash
# Run with default ~/todo.txt
eisenhower

# Run with a specific file
eisenhower /path/to/your/todo.txt

# Read from stdin (read-only mode) - combine multiple files
cat work.txt personal.txt | eisenhower

# Filter by project before viewing
grep '+WebApp' todo.txt | eisenhower

# View only high priority tasks
grep '(A)' todo.txt | eisenhower

# View archive file
eisenhower < done.txt

# Complex filtering with awk
awk '/+Project/ && !/(D)/' todo.txt | eisenhower
```

When reading from stdin (piped input), eisenhower enters read-only mode. All viewing and navigation features work normally (1-4 keys, filtering, inventory), but editing operations are disabled.

### Keyboard Controls

**Overview Mode:**
- `1`, `2`, `3`, `4` - Focus on a quadrant
- `q` - Quit

**Focus Mode:**
- `↑`/`↓` or `w`/`s` - Navigate todos
- `Space` - Toggle completion
- `a` - Add new todo
- `m` - Move todo to another quadrant
- `1`, `2`, `3`, `4` - Jump to different quadrant
- `ESC` - Return to overview

### Example todo.txt

```
(A) Fix critical bug in authentication +WebApp @computer
(B) Plan quarterly goals +Work @office
(C) Reply to client emails @phone
(D) Organize workspace
x (A) 2026-01-11 Completed task +Project
```

## Development

### Running Tests and Linter

```bash
# Run both tests and linter (recommended)
./check.sh

# Or run individually:
go test ./...              # Run tests
golangci-lint run          # Run linter
```

### Git Hooks

Install the pre-commit hook to automatically run tests and linter before every commit:

```bash
./scripts/install-hooks.sh
```

This ensures code quality by running `./check.sh` before each commit. To bypass the hook temporarily:

```bash
git commit --no-verify
```

See [CLAUDE.md](./CLAUDE.md) for detailed development conventions and architecture guidelines.

## Roadmap

### Completed ✅
- [x] **Story 001**: Display matrix with hard-coded todos
- [x] **Story 002**: Load todos from hardcoded file path (~/todo.txt)
- [x] **Story 003**: Accept custom file path as CLI argument
- [x] **Story 004**: Parse and render todo.txt project (+tag) and context (@tag) tags with colors
- [x] **Story 005**: Fullscreen TUI with alt-screen mode
- [x] **Story 006**: Responsive matrix sizing that adapts to terminal dimensions
- [x] **Story 007**: Quadrant focus mode (press 1/2/3/4 to jump between quadrants, ESC to return)
- [x] **Story 008**: Add todos with tag reference (press 'a' in focus mode)
- [x] **Story 009**: Tag autocomplete with + and @ triggers (arrow keys to navigate, Tab/Enter to complete)
- [x] **Story 010**: Tag inventory display showing counts by tag (sorted by count, overview mode only)
- [x] **Story 011**: Mark todos complete/incomplete (arrow keys or w/s to navigate, Space to toggle)
- [x] **Story 012**: Move todos between quadrants (Shift+1/2/3/4 to move, 1/2/3/4 to jump)
- [x] **Story 013**: Preserve and render completion dates (shows relative dates like "yesterday", "2 days ago")
- [x] **Story 014**: Parse and preserve creation dates from existing files
- [x] **Story 016**: Move todos between quadrants with 'm' key and overlay selector
- [x] **Story 017**: Show summary stats for each quadrant (task counts, completed counts)
- [x] **Story 018**: Delete todos with 'x' key and confirmation
- [x] **Story 019**: Digital inventory dashboard with metrics and analytics (press 'i')
- [x] **Story 020**: Edit todo descriptions and tags (press 'e' in focus mode)
- [x] **Story 021**: Filter todos by project or context tag (press 'f')
- [x] **Story 022**: Due date support with visual indicators (due:YYYY-MM-DD format)
- [x] **Story 023**: Archive completed todos to done.txt (press 'd')
- [x] **Story 024**: Stdin read-only mode for Unix composability (pipe todos for viewing)

### Future Ideas 🚀
- Search functionality (fuzzy search across descriptions)
- Recurring tasks
- Undo/redo functionality
- Bulk operations (archive all completed, delete all in quadrant)
- Custom sorting options
- Export to different formats
