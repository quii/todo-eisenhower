# Story 026: Stale Task Detection

As a user
I want to see visual indicators on tasks that have been sitting too long
So that I can identify tasks that need attention or should be reprioritized

## Background

Tasks that sit in quadrants for extended periods may indicate poor prioritization, procrastination, or tasks that should be eliminated. This story adds visual background highlighting to "stale" tasks based on:

- **Do First (Priority A)**: Stale after 2 business days from when prioritised
- **Schedule/Delegate/Eliminate (Priority B/C/D)**: Stale after 5 business days from creation

To accurately track Do First tasks, we introduce a hidden `prioritised:YYYY-MM-DD` tag that:
- Gets added when a task enters Do First (creation or move)
- Gets removed when a task leaves Do First
- Resets when a task re-enters Do First
- Is never visible in the UI or edit mode

## Acceptance Criteria

```gherkin
Feature: Stale Task Detection

  Scenario: Adding prioritised tag when creating task in Do First
    Given today is Tuesday, January 20, 2026
    When I create a new task "(A) New urgent task"
    Then the task should have tag "prioritised:2026-01-20"
    And the tag should not be visible in the UI
    And the tag should not appear in edit mode

  Scenario: Adding prioritised tag when moving task to Do First
    Given today is Tuesday, January 20, 2026
    And I have a task "(B) Scheduled task" in Schedule quadrant
    When I move the task to Do First quadrant
    Then the task priority should be "(A)"
    And the task should have tag "prioritised:2026-01-20"

  Scenario: Removing prioritised tag when moving task out of Do First
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Urgent task prioritised:2026-01-20" in Do First quadrant
    When I move the task to Schedule quadrant
    Then the task priority should be "(B)"
    And the task should not have a "prioritised:" tag

  Scenario: Resetting prioritised tag when moving task back to Do First
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Important task prioritised:2026-01-20" in Do First quadrant
    When I move the task to Schedule quadrant
    And 3 business days pass
    And today is Friday, January 23, 2026
    And I move the task back to Do First quadrant
    Then the task should have tag "prioritised:2026-01-23"
    And the task should not have tag "prioritised:2026-01-20"

  Scenario: Editing task in Do First preserves prioritised tag
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Original task prioritised:2026-01-15" in Do First quadrant
    When I edit the task description to "(A) Updated task +project"
    Then the task should still have tag "prioritised:2026-01-15"
    And the edit input should show "Updated task +project"
    And the edit input should not show "prioritised:2026-01-15"

  Scenario: Do First task not stale on day 1 (business day)
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Task prioritised:2026-01-20" in Do First
    When I view the matrix
    Then the task should not be marked as stale

  Scenario: Do First task not stale on day 2 (business day)
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Task prioritised:2026-01-20" in Do First
    When 1 business day passes
    And today is Wednesday, January 21, 2026
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Do First task is stale after 2 business days
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) Task prioritised:2026-01-20" in Do First
    When 2 business days pass
    And today is Thursday, January 22, 2026
    And I view the matrix
    Then the task should be marked as stale

  Scenario: Do First task staleness excludes weekends
    Given today is Thursday, January 15, 2026
    And I have a task "(A) Task prioritised:2026-01-15" in Do First
    When 4 calendar days pass including a weekend
    And today is Monday, January 19, 2026
    # Jan 15 (Thu) -> Jan 16 (Fri) = 1 business day
    # Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day (skip weekend)
    # Total: 2 business days
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Do First task stale after weekend threshold
    Given today is Thursday, January 15, 2026
    And I have a task "(A) Task prioritised:2026-01-15" in Do First
    When 5 calendar days pass including a weekend
    And today is Tuesday, January 20, 2026
    # Jan 15 (Thu) -> Jan 16 (Fri) = 1 business day
    # Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day (skip weekend)
    # Jan 19 (Mon) -> Jan 20 (Tue) = 1 business day
    # Total: 3 business days
    And I view the matrix
    Then the task should be marked as stale

  Scenario: Schedule task not stale within 5 business days
    Given today is Monday, January 13, 2026
    And I have a task "(B) Scheduled task" with creation date "2026-01-13" in Schedule
    When 4 business days pass
    And today is Friday, January 17, 2026
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Schedule task is stale after 5 business days
    Given today is Monday, January 13, 2026
    And I have a task "(B) Scheduled task" with creation date "2026-01-13" in Schedule
    When 5 business days pass
    And today is Monday, January 20, 2026
    And I view the matrix
    Then the task should be marked as stale

  Scenario: Delegate task staleness excludes weekends
    Given today is Monday, January 13, 2026
    And I have a task "(C) Delegate task" with creation date "2026-01-13" in Delegate
    When 8 calendar days pass including a weekend
    And today is Tuesday, January 21, 2026
    # Jan 13 (Mon) -> Jan 17 (Fri) = 4 business days
    # Jan 17 (Fri) -> Jan 20 (Mon) = 1 business day (skip weekend)
    # Jan 20 (Mon) -> Jan 21 (Tue) = 1 business day
    # Total: 6 business days
    And I view the matrix
    Then the task should be marked as stale

  Scenario: Eliminate task follows same staleness rules
    Given today is Monday, January 13, 2026
    And I have a task "(D) Low priority task" with creation date "2026-01-13" in Eliminate
    When 5 business days pass
    And today is Monday, January 20, 2026
    And I view the matrix
    Then the task should be marked as stale

  Scenario: Completed Do First task never stale
    Given today is Tuesday, January 13, 2026
    And I have a task "(A) Important task prioritised:2026-01-13" in Do First
    When 10 business days pass
    And today is Tuesday, January 27, 2026
    And I complete the task
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Completed Schedule task never stale
    Given today is Monday, January 13, 2026
    And I have a task "(B) Old task" with creation date "2026-01-13" in Schedule
    When 15 business days pass
    And today is Monday, February 2, 2026
    And I complete the task
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Task created on Friday not stale on Monday
    Given today is Friday, January 16, 2026
    And I have a task "(B) Weekend task" with creation date "2026-01-16" in Schedule
    When 3 calendar days pass including a weekend
    And today is Monday, January 19, 2026
    # Jan 16 (Fri) -> Jan 19 (Mon) = 1 business day
    And I view the matrix
    Then the task should not be marked as stale

  Scenario: Task created today is never stale
    Given today is Tuesday, January 20, 2026
    And I have a task "(A) New task prioritised:2026-01-20" in Do First
    When I view the matrix
    Then the task should not be marked as stale
```

## Technical Notes

### Domain Logic

- Add `IsStale() bool` method to `Todo` domain object
- Staleness calculation happens in the domain, not the UI
- Business days exclude Saturday and Sunday
- Day counting: Tuesday to Wednesday = 1 business day
- Completed tasks always return `false` for `IsStale()`
- Staleness thresholds:
  - Priority A: > 2 business days since `prioritised:` date
  - Priority B/C/D: > 5 business days since creation date

### Prioritised Tag Management

- Hidden tag format: `prioritised:YYYY-MM-DD`
- Added by `AddTodo` usecase when priority is A
- Added by `ChangePriority` usecase when moving to priority A
- Removed by `ChangePriority` usecase when moving from priority A
- Preserved by `EditTodo` usecase (description editing)
- Never displayed in UI rendering
- Filtered out in `FormatForInput()` for edit mode
- Included in `String()` method for persistence

### Business Day Calculation

```go
// Example helper function
func businessDaysBetween(from, to time.Time) int {
    // Count only Monday-Friday between dates
    // Skip weekends
}
```

### Visual Design

- Background color for stale tasks (to be determined during implementation)
- Applied in both overview and focus modes
- Only applies to uncompleted tasks
- Visual indicator is subtle but noticeable

## Future Considerations (NOT in this story)

- Configurable staleness thresholds
- Different staleness levels (warning vs critical)
- Staleness analytics in inventory view
- Option to dismiss staleness warning
- Staleness notifications or reminders
