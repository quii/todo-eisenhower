# Story 018: Delete Todo

**As a** user viewing my todos in focus mode  
**I want to** delete a todo with confirmation  
**So that** I can remove tasks I no longer need

## Acceptance Criteria

```gherkin
Feature: Delete Todo
  Users should be able to delete todos they no longer need

  Scenario: Delete a todo with confirmation
    Given I have a todo file with multiple todos
    And I am in focus mode on a quadrant
    When I navigate to a todo
    And I press Backspace
    Then I see a confirmation dialog asking "Delete this todo?"
    When I press "y" to confirm
    Then the todo is removed from the matrix
    And the todo is removed from the file
    And I am still in focus mode

  Scenario: Cancel deletion with ESC
    Given I have a todo in focus mode
    When I press Backspace to enter delete mode
    And I press ESC in the confirmation dialog
    Then the todo is not deleted
    And I return to focus mode

  Scenario: Cancel deletion with 'n'
    Given I have a todo in focus mode
    When I press Backspace to enter delete mode
    And I press 'n' in the confirmation dialog
    Then the todo is not deleted
    And I return to focus mode

  Scenario: Delete when only one todo exists
    Given I have exactly one todo in a quadrant
    When I delete it
    Then the quadrant becomes empty
    And I am returned to overview mode
```

## Technical Notes

- Use case: `DeleteTodo(writer TodoWriter, m matrix.Matrix, todoToDelete todo.Todo) (Matrix, error)`
- Domain: `Matrix.RemoveTodo(todo.Todo) Matrix` - returns new matrix without the todo
- UI: New delete mode (similar to move mode) with confirmation dialog
- Key binding: Backspace in focus mode
- Confirmation keys: 'y' (yes), 'n' (no), ESC (cancel)
