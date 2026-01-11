# Story 003: Accept Custom File Path

## Goal
As a user, I want to specify which todo.txt file to load, so I can manage multiple todo lists or use non-default locations.

## Background
- Story 001 established the TUI with hard-coded data
- Story 002 added file loading from `~/todo.txt`
- Now we add CLI argument parsing for custom file paths
- Default to `~/todo.txt` when no argument provided

## Acceptance Criteria

### Scenario: Load from custom file path
```gherkin
Given a todo.txt file at "/Users/chris/projects/work/todo.txt"
When I run "eisenhower /Users/chris/projects/work/todo.txt"
Then the matrix displays todos from that file
```

### Scenario: Use default path when no argument provided
```gherkin
Given a todo.txt file at "~/todo.txt"
When I run "eisenhower" without arguments
Then the matrix displays todos from "~/todo.txt"
```

### Scenario: Handle relative paths
```gherkin
Given a todo.txt file at "./todo.txt" in the current directory
When I run "eisenhower ./todo.txt"
Then the matrix displays todos from the current directory's todo.txt
```

### Scenario: Display file path in header
```gherkin
Given a todo.txt file at "/Users/chris/projects/todo.txt"
When I run "eisenhower /Users/chris/projects/todo.txt"
Then the header displays "File: /Users/chris/projects/todo.txt"
And the matrix is displayed below the header
```

### Scenario: Handle non-existent custom path
```gherkin
Given no file exists at "/path/to/missing.txt"
When I run "eisenhower /path/to/missing.txt"
Then the application displays "Error: file not found at /path/to/missing.txt"
And exits gracefully
```

### Scenario: Expand tilde in file paths
```gherkin
Given a todo.txt file in my home directory
When I run "eisenhower ~/projects/todo.txt"
Then the application correctly expands "~" to my home directory
And loads the file from the expanded path
```

## Technical Notes
- Add CLI argument parsing (use standard library or minimal flag package)
- Update main.go to accept file path argument
- Default to `~/todo.txt` when no argument provided
- Handle path expansion (tilde, relative paths)
- Add header component to TUI showing loaded file path
- Update LoadTodoMatrix use case to accept file path parameter

## Out of Scope
- Multiple file arguments
- Flag-based arguments (e.g., `--file` or `-f`)
- Help text or version flags (can add later)
- Config file for default path
- Watching file for changes

## Success Criteria
- Accepts custom file path as first argument
- Defaults to `~/todo.txt` when no argument provided
- Displays loaded file path in header
- Handles path expansion correctly
- All tests pass
- Linter passes
