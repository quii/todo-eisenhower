# Story 028: Backlog Quadrant

As a user
I want a separate Backlog area for ideas and tasks not yet ready for prioritization
So that I can capture thoughts without cluttering my Eisenhower matrix

## Background

The Eisenhower matrix is great for prioritizing actionable tasks, but sometimes you have ideas or tasks that aren't ready to be prioritized yet. Currently, these would need to go into one of the four quadrants (likely Eliminate), which muddies the purpose of that quadrant. A dedicated Backlog gives users a parking lot for ideas they want to capture but aren't ready to commit to.

The Backlog is intentionally hidden from the overview to keep focus on the prioritized work, but easily accessible when needed.

## Acceptance Criteria

```gherkin
Feature: Backlog Quadrant

  Scenario: Access backlog by pressing 5
    Given I am in overview mode
    When I press "5"
    Then I see the Backlog quadrant in focus mode
    And the title shows "Backlog"

  Scenario: Backlog uses priority E in todo.txt
    Given I have a todo.txt file with "(E) Idea for later"
    When I load the application
    And I press "5" to view the Backlog
    Then I see "Idea for later" in the Backlog

  Scenario: Add task directly to Backlog
    Given I am focused on the Backlog (pressed "5")
    When I press "a" to add a task
    And I enter "Research new framework"
    And I press Enter
    Then the task is saved with priority (E)
    And it appears in the Backlog

  Scenario: Backlog not shown in overview
    Given I have tasks in all quadrants including Backlog
    When I am in overview mode
    Then I see only the four Eisenhower quadrants
    And Backlog tasks are not displayed

  Scenario: Backlog count shown in help text
    Given I have 3 tasks in the Backlog
    When I am in overview mode
    Then the help text area shows "5: Backlog (3)"

  Scenario: Move task from Backlog to Do First
    Given I am focused on the Backlog
    And I have selected a task
    When I press "m" to enter move mode
    And I press "1"
    Then the task is moved to Do First quadrant
    And its priority changes from (E) to (A)

  Scenario: Move task from Do First to Backlog
    Given I am focused on Do First quadrant
    And I have selected a task
    When I press "m" to enter move mode
    And I press "5"
    Then the task is moved to the Backlog
    And its priority changes to (E)

  Scenario: Move mode shows Backlog option
    Given I am focused on any quadrant
    When I press "m" to enter move mode
    Then the move dialog shows options 1-5
    And option 5 is labeled "Backlog"

  Scenario: Empty Backlog
    Given I have no tasks with priority (E)
    When I press "5" to view the Backlog
    Then I see an empty Backlog view
    And I can still add tasks with "a"

  Scenario: Return to overview from Backlog
    Given I am focused on the Backlog
    When I press "0" or "Esc"
    Then I return to the overview mode
```

## Technical Notes

### Domain Layer

- Add `BacklogQuadrant` to the `QuadrantType` enum
- Map priority E to Backlog in matrix categorization
- Backlog should be excluded from `AllTodos()` used by overview, or a new method `EisenhowerTodos()` that excludes Backlog

### UI Layer

- Add `FocusBacklog` to `ViewMode` enum
- Handle "5" key in overview to focus Backlog
- Update move mode dialog to include option 5
- Update help text rendering to show Backlog count
- Backlog quadrant styling - suggest a muted/neutral color (gray?) to visually distinguish from the urgent/important matrix

### Considerations

- The overview matrix rendering should explicitly exclude Backlog items
- Stats in overview (if any) should not include Backlog
- Inventory dashboard excludes Backlog (per user preference)

## Future Considerations (NOT in this story)

- Backlog review mode: guided flow to process backlog items
- Auto-suggest moving stale Backlog items to Eliminate
- Search/filter within Backlog only
