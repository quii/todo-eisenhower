# Story 005: Fullscreen TUI

## User Story
As a user, I want the application to take over the full terminal screen (like Claude Code) so that I have an immersive viewing experience for my Eisenhower matrix.

## Background
Currently the app renders the matrix once and exits. We want it to:
- Take over the entire terminal (alt-screen mode)
- Stay running until the user quits
- Allow quitting with 'q' or Ctrl+C
- Center the matrix in the terminal

## Acceptance Criteria

```gherkin
Feature: Fullscreen TUI Experience

  Scenario: Application runs in alt-screen mode
    Given I have a todo.txt file
    When I run the application
    Then the terminal switches to alt-screen mode
    And the matrix is displayed
    And the application waits for input
    And my previous terminal content is preserved

  Scenario: User can quit the application
    Given the application is running in fullscreen mode
    When I press 'q'
    Then the application exits
    And the terminal returns to normal mode
    And my previous terminal content is restored

  Scenario: User can quit with Ctrl+C
    Given the application is running in fullscreen mode
    When I press Ctrl+C
    Then the application exits gracefully
    And the terminal returns to normal mode

  Scenario: Matrix is centered in terminal
    Given the application is running in fullscreen mode
    When the matrix is displayed
    Then it is centered horizontally and vertically
    And the file path header is visible at the top
```

## Technical Notes
- Use Bubble Tea's `tea.NewProgram()` with alt-screen mode
- Wire up the existing `ui.Model` that was created in Story 001
- Add keyboard handling for 'q' and Ctrl+C
- Use `tea.WithAltScreen()` option
- Current rendering logic in `ui.RenderMatrix` can be used in the View

## Out of Scope
- Window resizing responsiveness (future story)
- Interactive navigation (future story)
- Editing capabilities (future story)
