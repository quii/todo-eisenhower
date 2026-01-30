# Story 021: Filter by Tag

As a user
I want to filter todos by project or context
So that I can focus on specific areas of work or contexts

## Acceptance Criteria

```gherkin
Feature: Filter by Tag

  Scenario: Press 'f' to enter filter mode
    Given I am in overview mode
    When I press 'f'
    Then I should enter filter mode
    And I should see an input prompt "Filter by: "
    And I should see suggestions for all projects and contexts

  Scenario: Filter by project
    Given I am in filter mode
    When I type "+WebApp"
    Then I should see autocomplete suggestions for projects starting with "WebApp"
    When I press Enter
    Then I should see a filtered view showing only todos with project "WebApp"
    And todos should be organized by quadrant
    And the view should show "Filtered by: +WebApp" at the top

  Scenario: Filter by context
    Given I am in filter mode
    When I type "@computer"
    Then I should see autocomplete suggestions for contexts starting with "computer"
    When I press Enter
    Then I should see a filtered view showing only todos with context "computer"
    And todos should be organized by quadrant
    And the view should show "Filtered by: @computer" at the top

  Scenario: Autocomplete works without prefix
    Given I am in filter mode
    When I type "Web" (without +)
    Then I should see autocomplete suggestions for projects matching "Web"
    When I type "comp" (without @)
    Then I should see autocomplete suggestions for contexts matching "comp"

  Scenario: Autocomplete shows both projects and contexts
    Given I am in filter mode
    And there is a project "+work" and a context "@work"
    When I type "work"
    Then I should see both "+work" and "@work" in suggestions

  Scenario: Clear filter and return to overview
    Given I am viewing filtered todos
    When I press ESC or 'c' (clear)
    Then the filter should be cleared
    And I should return to overview mode

  Scenario: Cancel filter input
    Given I am in filter mode (typing)
    When I press ESC
    Then I should cancel filter mode
    And return to overview without applying a filter

  Scenario: Empty filter input shows all tags
    Given I am in filter mode
    And the input is empty
    Then I should see all projects and contexts as suggestions
    Organized with projects first, then contexts

  Scenario: No matches shows empty state
    Given I am in filter mode
    When I type a tag that doesn't exist
    And press Enter
    Then I should see "No todos match filter: +NonExistent"
    And pressing ESC should return to overview

  Scenario: Filter persists when navigating quadrants
    Given I have a filter "+WebApp" applied
    When I press '1' to focus on Do First quadrant
    Then I should only see Do First todos with "+WebApp"
    When I press '0' to return to overview
    Then I should still see the filtered overview

  Scenario: Can edit/complete/move filtered todos
    Given I have a filter applied
    When I focus on a quadrant and select a todo
    Then I should be able to edit (e), complete (space), or move (m) as normal
    And the filter should remain active after the operation
```

## Technical Notes

- Filter input reuses the same autocomplete component as add/edit
- Filter state is part of the model (activeFilter string, filterMode bool)
- Rendering shows filtered todos across all quadrants (like inventory but with quadrant sections)
- Filter should be case-insensitive
- Only one filter active at a time (simplicity)
- Filter by exact tag match (not substring of description)
