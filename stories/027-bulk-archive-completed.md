# Story 027: Bulk Archive Completed Tasks

As a user
I want to archive all completed tasks at once with Shift+D
So that I can quickly clear completed work without archiving tasks one by one

## Background

Users can already archive individual completed tasks with the 'd' key, which moves them to the done.txt file. When many tasks are completed (e.g., after a productive session), archiving them one at a time is tedious. Shift+D provides a bulk action that respects the current view context - archiving completed tasks in the focused quadrant, or all quadrants when in overview mode.

## Acceptance Criteria

```gherkin
Feature: Bulk Archive Completed Tasks

  Scenario: Bulk archive completed tasks in focused quadrant
    Given I am focused on the "Do" quadrant
    And the quadrant contains 3 completed tasks and 2 incomplete tasks
    When I press Shift+D
    Then all 3 completed tasks are moved to done.txt
    And the 2 incomplete tasks remain in the quadrant

  Scenario: Bulk archive completed tasks in overview mode
    Given I am in overview mode
    And the "Do" quadrant has 2 completed tasks
    And the "Schedule" quadrant has 1 completed task
    And the "Delegate" quadrant has 3 completed tasks
    And the "Delete" quadrant has 0 completed tasks
    When I press Shift+D
    Then all 6 completed tasks across all quadrants are moved to done.txt

  Scenario: No completed tasks to archive in focused quadrant
    Given I am focused on the "Do" quadrant
    And the quadrant contains only incomplete tasks
    When I press Shift+D
    Then nothing happens
    And the view remains unchanged

  Scenario: No completed tasks to archive in overview mode
    Given I am in overview mode
    And no quadrant contains completed tasks
    When I press Shift+D
    Then nothing happens
    And the view remains unchanged

  Scenario: Empty quadrant in focus mode
    Given I am focused on the "Schedule" quadrant
    And the quadrant is empty
    When I press Shift+D
    Then nothing happens

  Scenario: Mixed quadrants with some empty in overview mode
    Given I am in overview mode
    And the "Do" quadrant has 2 completed tasks
    And the "Schedule" quadrant is empty
    And the "Delegate" quadrant has only incomplete tasks
    And the "Delete" quadrant has 1 completed task
    When I press Shift+D
    Then the 3 completed tasks from "Do" and "Delete" are moved to done.txt
    And the incomplete tasks in "Delegate" remain unchanged
```

## Technical Notes

### Domain Layer

The matrix domain should expose two methods for bulk archiving:

- `ArchiveQuadrant(quadrant Quadrant) []Todo` - Archives all completed tasks in a specific quadrant, returns the archived todos
- `ArchiveAllCompleted() []Todo` - Archives all completed tasks across all quadrants, returns the archived todos

These methods should follow the same pattern as the existing single-task archive functionality.

### UI Layer

The UI adapter needs to:
- Handle Shift+D key binding
- Determine current context (focused quadrant vs overview mode)
- Call the appropriate use case based on context
- Update the view to reflect removed tasks

### Use Cases

Add a new use case (or extend existing archive use case):
- `ArchiveCompletedInQuadrant(quadrant)` - for focused mode
- `ArchiveAllCompleted()` - for overview mode

Both should coordinate with the existing done.txt persistence logic.

## Future Considerations (NOT in this story)

- Undo bulk archive action
- Confirmation dialog as an optional user preference
- Visual animation showing tasks being archived
- Count indicator briefly showing "Archived N tasks"
