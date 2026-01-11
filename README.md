# Eisenhower Matrix TUI

A terminal-based Eisenhower matrix viewer for todo.txt files, built with Go and Bubble Tea.

## What's an Eisenhower Matrix?

The Eisenhower Matrix helps prioritize tasks by organizing them into four quadrants:
- **DO FIRST** (Q1): Urgent & Important
- **SCHEDULE** (Q2): Important, Not Urgent
- **DELEGATE** (Q3): Urgent, Not Important
- **ELIMINATE** (Q4): Neither Urgent nor Important

## Current Status

**Story 001 Complete**: Display matrix with hard-coded todos

The application currently displays a working Eisenhower matrix with hard-coded sample data.

## Building and Running

```bash
# Build
go build -o eisenhower ./cmd/eisenhower/

# Run
./eisenhower

# Press 'q' or Ctrl+C to quit
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

```
.
├── domain/           # Domain layer (business logic)
│   ├── todo/        # Todo entities
│   ├── matrix/      # Matrix categorization logic
│   └── parser/      # todo.txt parsing (coming in Story 002)
├── usecases/        # Application use cases
├── adapters/        # Adapters layer
│   ├── ui/          # Bubble Tea TUI
│   └── file/        # File operations (coming in Story 002)
├── cmd/
│   └── eisenhower/  # Main entry point
└── stories/         # User stories (Gherkin format)
```

## Development

See [CLAUDE.md](./CLAUDE.md) for development conventions and architecture guidelines.

## Roadmap

- [x] Story 001: Display matrix with hard-coded todos
- [ ] Story 002: Load todos from hard-coded file path
- [ ] Story 003: Accept custom file path as argument
- [ ] Future: Editing, filtering, and more!
