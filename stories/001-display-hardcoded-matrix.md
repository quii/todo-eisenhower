# Story 001: Display Hardcoded Matrix

## Goal
As a developer, I want to see a working Eisenhower matrix TUI with hard-coded todos, so I can validate the application stack and UI rendering.

## Background
- Eisenhower matrix has 4 quadrants: Do First (Q1), Schedule (Q2), Delegate (Q3), Eliminate (Q4)
- This is iteration 0: establish project structure, TUI rendering, and domain architecture
- Use hard-coded in-memory todos to keep this slice small
- Future stories will add file loading and command-line arguments

## Acceptance Criteria

### Scenario: Display matrix with hard-coded todos
```gherkin
Given the application has hard-coded todos in memory
When I run "eisenhower"
Then I see a 2x2 matrix with labeled quadrants
And the "DO FIRST" quadrant contains at least one todo
And the "SCHEDULE" quadrant contains at least one todo
And the "DELEGATE" quadrant contains at least one todo
And the "ELIMINATE" quadrant contains at least one todo
```

### Scenario: Matrix displays quadrant labels
```gherkin
Given the application is running
Then the top-left quadrant shows "DO FIRST" or "Q1"
And the top-right quadrant shows "SCHEDULE" or "Q2"
And the bottom-left quadrant shows "DELEGATE" or "Q3"
And the bottom-right quadrant shows "ELIMINATE" or "Q4"
```

### Scenario: Todos appear in correct quadrants
```gherkin
Given hard-coded todos with priorities (A), (B), (C), (D)
When the matrix is rendered
Then priority (A) todos appear in "DO FIRST" quadrant
And priority (B) todos appear in "SCHEDULE" quadrant
And priority (C) todos appear in "DELEGATE" quadrant
And priority (D) todos appear in "ELIMINATE" quadrant
```

## Technical Notes
- Set up full project structure (go.mod, directories, linting)
- Implement basic Todo domain with Priority field
- Implement Matrix domain to categorize todos by quadrant
- Create Bubble Tea UI adapter to render the matrix
- Hard-code sample todos in main.go for now
- No file I/O, no command-line arguments yet

## Hard-coded Sample Data
Suggested todos for initial implementation:
```
(A) Fix critical production bug
(A) Review security audit findings
(B) Plan Q2 roadmap
(B) Research new framework options
(C) Respond to routine emails
(C) Attend weekly status meeting
(D) Organize old project files
(D) Update personal wiki
```

## Out of Scope
- Loading todos from files
- Command-line arguments
- Editing todos
- Completed todo handling (can add `x` prefix later)
- Overflow/truncation (keep it simple for now)

## Success Criteria
- Application compiles and runs
- TUI displays without errors
- Matrix layout is clear and readable
- Todos appear in correct quadrants based on priority
- All tests pass
- Linter passes
- Code follows DDD architecture (domain/usecases/adapters)
