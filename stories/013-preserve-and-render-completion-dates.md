# Story 013: Preserve and Render Completion Dates

## User Story
**As a** todo.txt user
**I want** completion dates to be preserved and displayed
**So that** I can see when I completed each task and maintain todo.txt format compliance

## Background
The todo.txt format specifies that completed todos should include a completion date:
```
x 2026-01-15 (A) Completed task description
```

Currently, we:
- ✅ Write completion dates when marking todos complete
- ❌ Don't preserve existing completion dates from the file
- ❌ Overwrite dates if you toggle complete→incomplete→complete
- ❌ Don't display completion dates in the UI

## Acceptance Criteria

### Scenario: Parse and preserve completion dates
```gherkin
Given a todo.txt file containing:
  """
  x 2026-01-10 (A) Completed task from last week
  (B) Active task
  """
When I load the file
Then the completed todo should preserve the date "2026-01-10"
And the active todo should have no completion date
```

### Scenario: Set completion date when marking complete
```gherkin
Given an active todo "(A) Review documentation"
When I toggle it to complete
Then it should have today's date as the completion date
And the file should contain "x 2026-01-15 (A) Review documentation"
```

### Scenario: Clear completion date when toggling incomplete
```gherkin
Given a completed todo "x 2026-01-10 (A) Review documentation"
When I toggle it to incomplete
Then it should have no completion date
And the file should contain "(A) Review documentation" (no x marker or date)
```

### Scenario: New completion date when re-completing
```gherkin
Given a completed todo "x 2026-01-10 (A) Review documentation"
When I toggle it to incomplete (Jan 14)
And I toggle it back to complete (Jan 15)
Then it should have the date "2026-01-15" (today)
And the file should contain "x 2026-01-15 (A) Review documentation"
```

### Scenario: Display completion date in UI
```gherkin
Given a completed todo "x 2026-01-10 (A) Completed task"
When I view the quadrant
Then the todo should display "✓ Completed task (completed 2026-01-10)"
```

### Scenario: Display completion date with relative formatting
```gherkin
Given a completed todo from 2 days ago "x 2026-01-13 (A) Recent task"
When I view the quadrant
Then the todo should display "✓ Recent task (completed 2 days ago)"
```

### Scenario: No date shown for incomplete todos
```gherkin
Given an active todo "(A) Active task"
When I view the quadrant
Then the todo should display "• Active task"
And no completion date should be shown
```

### Scenario: Preserve completion date when moving quadrants
```gherkin
Given a completed todo in DO FIRST "x 2026-01-10 (A) Completed urgent task"
When I move it to SCHEDULE (priority B)
Then it should preserve the date "2026-01-10"
And the file should contain "x 2026-01-10 (B) Completed urgent task"
```

## Technical Notes

### Domain Model Changes
- Add `completionDate *time.Time` field to `Todo` struct (pointer = nullable)
- Update `ToggleCompletion()` to:
  - Set current date when marking complete (always overwrites any existing date)
  - Clear date (set to nil) when marking incomplete
- Update all constructors to accept optional completion date

### Parser Changes
- Parse completion date from `x YYYY-MM-DD` format
- Store in Todo domain object
- Handle both formats: `x DATE (PRIORITY) Description` (standard)

### Formatter Changes
- Use stored completion date instead of `time.Now()`
- Only write date if todo is completed and has a date
- Format as `x YYYY-MM-DD (PRIORITY) Description`

### UI Rendering
- Show completion date next to completed todos
- Use relative formatting for recent dates (today, yesterday, X days ago)
- Use absolute date for older dates (e.g., "Jan 10, 2026")
- Consider format: `✓ Description (completed Jan 10)` or `✓ Description [Jan 10]`

## Definition of Done
- [ ] Domain model stores completion date as nullable field
- [ ] Parser preserves completion dates from file
- [ ] Formatter writes stored completion date (not current date)
- [ ] ToggleCompletion sets date to today when marking complete
- [ ] ToggleCompletion clears date when marking incomplete
- [ ] UI displays completion dates for completed todos
- [ ] Completion dates get new date each time marked complete
- [ ] Completion dates preserved when moving quadrants
- [ ] All acceptance tests pass
- [ ] Existing tests updated for new format
- [ ] Manual testing confirms dates work correctly
