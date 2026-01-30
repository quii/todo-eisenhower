# Story 022: Due Date Support

As a user
I want to see due dates on my todos
So that I can track deadlines without them affecting my urgency/importance prioritization

## Background

The todo.txt format supports due dates with the format `due:YYYY-MM-DD`. This story adds
visual indicators for due dates without changing the core Eisenhower matrix logic (priorities
remain the source of truth for urgency/importance).

## Acceptance Criteria

```gherkin
Feature: Due Date Support

  Scenario: Display due date in overview mode
    Given I have a todo "(A) Submit report due:2026-01-25"
    And today is 2026-01-20
    When I view the overview matrix
    Then I should see "Submit report due: Jan 25" in the Do First quadrant
    And the due date should be displayed in a distinct color (not red)

  Scenario: Display due date in focused quadrant mode
    Given I have multiple todos with due dates in the Do First quadrant
    When I press '1' to focus on Do First
    Then I should see a table with columns: [Description] [Due Date] [Created]
    And due dates should be formatted as "Jan 25" or "2026-01-25"
    And due dates should be displayed in a distinct color

  Scenario: Display overdue items
    Given I have a todo "(A) Overdue task due:2026-01-15"
    And today is 2026-01-20
    When I view the overview matrix
    Then I should see "Overdue task ! due: Jan 15"
    And the exclamation mark and due date should be in red
    And the text should clearly indicate this is overdue

  Scenario: Display overdue items in focused mode
    Given I have a todo "(A) Overdue task due:2026-01-15"
    And today is 2026-01-20
    When I press '1' to focus on Do First
    Then the due date column should show "! Jan 15" in red
    Or "! 2026-01-15" in red

  Scenario: Todos without due dates are unaffected
    Given I have a todo "(B) Regular task"
    When I view the matrix
    Then I should see "Regular task" without any due date information
    And it should look the same as before this feature

  Scenario: Due date exactly today
    Given I have a todo "(A) Task due:2026-01-20"
    And today is 2026-01-20
    Then the due date should be displayed in a distinct color (not red)
    Because it's not overdue yet

  Scenario: Due dates don't affect sorting
    Given I have three todos:
      | Description | Priority | Due Date   |
      | Task A      | A        | 2026-01-15 |
      | Task B      | A        | 2026-01-25 |
      | Task C      | A        | 2026-01-20 |
    When I view the Do First quadrant
    Then todos should appear in the order they were added
    And NOT sorted by due date

  Scenario: Due dates are informational only
    Given I have a todo "(D) Low priority due:2026-01-15"
    And the task is overdue
    When I view the matrix
    Then the task should still be in the Eliminate quadrant
    Because priority is the source of truth for quadrant placement

  Scenario: Due dates work with all todo features
    Given I have a todo "(A) Task +project @context due:2026-01-25"
    When I view the matrix
    Then I should see the description, colorized tags, AND the due date
    And all features should work together (edit, complete, move, filter)

  Scenario: Invalid due dates are ignored
    Given I have a todo "(A) Task due:invalid-date"
    When I view the matrix
    Then the invalid due date should not be displayed
    Or should be displayed as-is without special formatting

  Scenario: Multiple due date tags (edge case)
    Given I have a todo "(A) Task due:2026-01-20 due:2026-01-25"
    When I view the matrix
    Then only the first due date should be used (2026-01-20)
    Or both should be displayed as-is without special formatting
```

## Technical Notes

- Due date parsing should use the existing `todotxt` package
- Date comparison should use `time.Time` for accuracy
- Due dates should NOT be removed from the description text (they remain part of it)
- Visual indicators should be added in the rendering layer, not the domain
- Color choices:
  - Not yet due: A subtle distinct color (maybe cyan/blue tint)
  - Overdue: Red with exclamation mark prefix
- Date formatting: Short format like "Jan 25" for current year, "2025-12-25" for other years
- This feature is purely display - no new domain logic, just parsing and rendering

## Future Considerations (NOT in this story)

- Filter by due date (Story 024?)
- Sort by due date (Story 025?)
- Inventory analytics for due dates (Story 026?)
- "Due soon" warnings
- Recurring due dates
