# Story 014: Support Creation Dates

## User Story
**As a** todo.txt user
**I want** creation dates to be preserved and displayed consistently
**So that** I can see when tasks were created and maintain full todo.txt format compliance

## Background
The todo.txt format supports creation dates:
- Active todos: `(A) 2026-01-15 Task description`
- Completed todos: `x 2026-01-10 2026-01-05 (A) Task description` (completion date, then creation date)

Currently we:
- ✅ Support completion dates
- ❌ Don't preserve or display creation dates
- ❌ Don't set creation dates when adding new todos

Also, date formatting should be consistently friendly:
- "today"
- "yesterday"
- "2 days ago", "3 days ago", etc. (for all older dates)

## Acceptance Criteria

### Scenario: Parse and preserve creation dates
```gherkin
Given a todo.txt file containing:
  """
  (A) 2026-01-10 Task created on Jan 10
  (B) Task without creation date
  """
When I load the file
Then the first todo should have creation date "2026-01-10"
And the second todo should have no creation date
```

### Scenario: Set creation date when adding new todo
```gherkin
Given I'm in focus mode on DO FIRST quadrant
When I press 'a' and add a new todo "Buy groceries"
Then the file should contain "(A) 2026-01-15 Buy groceries"
And the creation date should be today's date
```

### Scenario: Display creation date in UI
```gherkin
Given a todo created 3 days ago "(A) 2026-01-12 Buy groceries"
When I view the quadrant
Then it should display "• Buy groceries (added 3 days ago)"
```

### Scenario: Display both creation and completion dates
```gherkin
Given a completed todo "x 2026-01-14 2026-01-10 (A) Finished task"
When I view the quadrant
Then it should display "✓ Finished task (added 5 days ago, completed yesterday)"
```

### Scenario: Friendly date formatting for all dates
```gherkin
Given todos with various creation dates
When I view the UI
Then dates should show as:
  | Date | Format |
  | Today | "today" |
  | Yesterday | "yesterday" |
  | 2 days ago | "2 days ago" |
  | 30 days ago | "30 days ago" |
  | 365 days ago | "365 days ago" |
```

### Scenario: No creation date shown for todos without dates
```gherkin
Given a todo without creation date "(A) Old todo"
When I view the quadrant
Then it should display "• Old todo"
And no creation date should be shown
```

### Scenario: Preserve creation date when toggling completion
```gherkin
Given a todo "(A) 2026-01-10 Important task"
When I toggle it to complete
Then the file should contain "x 2026-01-15 2026-01-10 (A) Important task"
And the creation date "2026-01-10" should be preserved
```

### Scenario: Preserve creation date when moving quadrants
```gherkin
Given a todo "(A) 2026-01-10 Task to move"
When I move it to SCHEDULE (priority B)
Then the file should contain "(B) 2026-01-10 Task to move"
And the creation date "2026-01-10" should be preserved
```

## Technical Notes

### Domain Model Changes
- Add `creationDate *time.Time` field to `Todo` struct
- Update all constructors to accept optional creation date
- Update `ToggleCompletion()` to preserve creation date
- Update `ChangePriority()` to preserve creation date (already does this pattern)

### Parser Changes
- Parse creation date from position after priority
- Handle both formats: with and without completion date
- Store in Todo domain object

### Formatter Changes
- Write creation date after priority for active todos
- Write creation date after completion date for completed todos
- Format: `(A) YYYY-MM-DD Description` or `x COMP_DATE CREATION_DATE (A) Description`

### AddTodo Usecase Changes
- Set creation date to current date when adding new todos
- Include in todo.txt output

### UI Rendering Changes
- Update `formatCompletionDate` to use "N days ago" for all dates beyond yesterday
- Create unified date formatting function for consistency
- Display creation dates with "added X ago" format
- Show both dates for completed todos: "(added X ago, completed Y ago)"

## Definition of Done
- [ ] Domain model stores creation date
- [ ] Parser preserves creation dates from file
- [ ] Formatter writes creation dates in correct position
- [ ] AddTodo sets creation date to now when creating new todos
- [ ] Date formatting uses "today", "yesterday", "N days ago" consistently
- [ ] UI displays creation dates for todos that have them
- [ ] UI displays both dates for completed todos
- [ ] Creation dates preserved when toggling completion
- [ ] Creation dates preserved when moving quadrants
- [ ] All acceptance tests pass
- [ ] Existing tests updated for new format
- [ ] Manual testing confirms dates display correctly
