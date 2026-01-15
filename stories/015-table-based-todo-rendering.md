# Story 015: Table-based Todo Rendering

## Status
Planned

## Context
Currently, todos are rendered as run-on lines with all information (description, tags, dates) concatenated together:
```
• Task description +project @context (added 2 days ago)
✓ Completed task +project (added 5 days ago, completed today)
```

This format becomes difficult to read as we add more metadata (creation dates, completion dates, priority changes, etc.). The bubbles library provides a table component that would allow us to display todo information in a structured, columnar format.

## Goal
Replace the current line-based todo rendering with a table component from the bubbles library, displaying todos in a structured columnar format that makes metadata easier to scan and read.

## Acceptance Criteria

### Scenario 1: Display todos in table format
**Given** I have todos with various metadata (priority, tags, dates)
**When** I view a quadrant in focus mode
**Then** I should see todos displayed in a table with columns:
- Status (✓ or • indicator)
- Description (with colorized tags inline)
- Projects (comma-separated list)
- Contexts (comma-separated list)
- Created (friendly date format)
- Completed (friendly date format, or empty for active todos)

### Scenario 2: Navigate table with keyboard
**Given** I am viewing a quadrant with multiple todos
**When** I press up/down arrow keys or w/s
**Then** the table selection should move between rows
**And** the selected row should be highlighted

### Scenario 3: Table columns adjust to terminal width
**Given** I have a terminal of varying width
**When** I resize the terminal
**Then** the table should adjust column widths appropriately
**And** description column should use available space
**And** metadata columns should have fixed, reasonable widths

### Scenario 4: Empty quadrant shows appropriate message
**Given** I am viewing an empty quadrant
**When** I enter focus mode
**Then** I should see "(no tasks)" instead of an empty table

### Scenario 5: Table respects display limit
**Given** I have more todos than can fit on screen
**When** I view the quadrant
**Then** the table should show only the visible portion
**And** I should be able to scroll through all todos

### Scenario 6: Completed todos are visually distinct
**Given** I have both active and completed todos
**When** I view them in the table
**Then** completed todos should have:
- ✓ in the status column
- Dimmed/grayed out styling
- Completion date populated

### Scenario 7: Table headers are clear and concise
**Given** I am viewing a table of todos
**When** I look at the column headers
**Then** they should be:
- Status (or blank)
- Task
- Projects
- Contexts
- Created
- Completed

### Scenario 8: Existing keyboard shortcuts still work
**Given** I am viewing a table of todos
**When** I press space, ESC, Shift+number, etc.
**Then** existing functionality (toggle complete, exit focus, move quadrants) should still work
**And** the table navigation should not interfere

## Technical Notes
- Use `github.com/charmbracelet/bubbles/table` component
- The table should integrate with the existing Model structure
- Selected row should match `selectedTodoIndex` in the model
- Table styling should match the existing lipgloss styles
- Consider how input mode interacts with the table
- Tag colorization should still work within the description column

## Design Decisions
- **Column order**: Status, Task, Projects, Contexts, Created, Completed
- **Date columns**: Use the existing friendly format ("today", "yesterday", "N days ago")
- **Tag display**: Keep tags inline with description but also show in dedicated columns for scanning
- **Navigation**: Reuse existing up/down/w/s handling, just update to work with table rows
- **Overflow**: Description column can truncate with "..." if needed, or wrap depending on width

## Dependencies
- Story 014 (creation dates) - provides data for Created column
- Story 013 (completion dates) - provides data for Completed column
- Story 004 (projects/contexts) - provides data for tag columns

## Out of Scope
- Sortable columns (future enhancement)
- Column reordering (future enhancement)
- Column width customization (future enhancement)
- Filtering/searching within table (future enhancement)
