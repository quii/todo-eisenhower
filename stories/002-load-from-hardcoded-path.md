# Story 002: Load from Hardcoded Path

## Goal
As a user, I want the application to load my todos from a todo.txt file at a fixed location, so I can see my actual tasks in the matrix.

## Background
- Story 001 established the TUI and domain architecture with hard-coded data
- Now we add file loading capability
- Use a hard-coded path (e.g., `~/todo.txt`) for simplicity
- This requires implementing the Parser domain and File adapter

## Acceptance Criteria

### Scenario: Load todos from hardcoded file path
```gherkin
Given a todo.txt file exists at "~/todo.txt" containing:
  """
  (A) Fix critical bug
  (B) Plan quarterly goals
  (C) Reply to emails
  (D) Clean workspace
  """
When I run "eisenhower"
Then the matrix displays todos from "~/todo.txt"
And the "DO FIRST" quadrant contains "Fix critical bug"
And the "SCHEDULE" quadrant contains "Plan quarterly goals"
And the "DELEGATE" quadrant contains "Reply to emails"
And the "ELIMINATE" quadrant contains "Clean workspace"
```

### Scenario: Handle missing file gracefully
```gherkin
Given no file exists at "~/todo.txt"
When I run "eisenhower"
Then the application displays an error message
And exits gracefully without crashing
```

### Scenario: Handle empty file
```gherkin
Given an empty file exists at "~/todo.txt"
When I run "eisenhower"
Then the matrix displays with all quadrants empty
```

### Scenario: Parse completed todos
```gherkin
Given a todo.txt file containing:
  """
  (A) Active task
  x (A) 2026-01-11 Completed task
  (B) Another active task
  """
When I run "eisenhower"
Then the "DO FIRST" quadrant shows both todos
And completed todos are visually distinct from active todos
```

### Scenario: Handle todos without priority
```gherkin
Given a todo.txt file containing:
  """
  (A) High priority task
  No priority task
  """
When I run "eisenhower"
Then the "DO FIRST" quadrant contains "High priority task"
And the "ELIMINATE" quadrant contains "No priority task"
```

## Technical Notes
- Implement Parser domain to parse todo.txt format
- Implement File adapter for reading files
- Create LoadTodoMatrix use case to orchestrate file loading and matrix creation
- Hard-code path `~/todo.txt` in main.go
- Handle file errors gracefully (missing, unreadable, etc.)
- Expand Todo domain to include completion status

## todo.txt Format Reference
Basic format:
```
(A) Task description
x (A) 2026-01-11 Completed task description
Task without priority
```

For this story, we only need to parse:
- Priority: `(A)`, `(B)`, `(C)`, `(D)` or none
- Completion: `x` prefix
- Description: remaining text

## Out of Scope
- Custom file paths (next story)
- Dates, contexts, projects from todo.txt spec
- File watching or auto-reload
- Header showing file path (can add later)

## Success Criteria
- Reads todos from `~/todo.txt`
- Correctly parses priority and completion status
- Handles errors gracefully
- All tests pass (including parser tests)
- Linter passes
