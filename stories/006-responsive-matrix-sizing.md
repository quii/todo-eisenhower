# Story 006: Responsive Matrix Sizing

## User Story
As a user, I want the matrix to intelligently fill the available terminal space so that I can see more of my todos and make better use of my screen real estate.

## Background
Currently the matrix has hardcoded dimensions (quadrantWidth = 40, quadrantHeight = 10). Now that we have fullscreen mode and window size tracking, we should make the matrix responsive to terminal dimensions.

The matrix should:
- Calculate quadrant dimensions based on available terminal space
- Leave reasonable margins and spacing
- Adjust the number of visible todos based on quadrant height
- Handle window resize events gracefully (already captured via WindowSizeMsg)

## Acceptance Criteria

```gherkin
Feature: Responsive Matrix Sizing

  Scenario: Matrix fills available terminal space
    Given I have a todo.txt file with many todos
    When I run the application in a 120x40 terminal
    Then the matrix quadrants are larger than the default 40x10
    And more todos are visible per quadrant
    And the matrix is still centered

  Scenario: Matrix adjusts to different terminal sizes
    Given I have a todo.txt file
    When I run the application in a 200x60 terminal
    Then the quadrants are larger than in a 120x40 terminal
    And more todos are displayed

  Scenario: Matrix respects minimum dimensions
    Given I have a todo.txt file
    When I run the application in a very small terminal (80x24)
    Then the matrix uses minimum viable dimensions
    And does not break the layout
    And still displays at least some todos

  Scenario: Todo display limit scales with height
    Given I have 20 todos in the Do First quadrant
    When the quadrant height is 20 lines
    Then I can see approximately 15-17 todos (accounting for title and spacing)
    When the quadrant height is 10 lines
    Then I see approximately 7 todos (current behavior)

  Scenario: Window resize is handled gracefully
    Given the application is running
    When I resize my terminal window
    Then the matrix re-renders with new dimensions
    And the layout remains correct
```

## Technical Notes

### Calculation approach:
1. Get terminal width and height from WindowSizeMsg
2. Reserve space for:
   - Top/bottom margins
   - Matrix border
   - Axis label ("← URGENT →")
   - File path header
3. Calculate quadrant dimensions:
   - Width: (availableWidth - border - divider) / 2
   - Height: (availableHeight - border - divider - header) / 2
4. Adjust todo display limit based on quadrant height
   - Account for title (2 lines with spacing)
   - Each todo takes 1 line
   - Reserve 1-2 lines for "... and X more" if needed

### Changes needed:
- Make RenderMatrix accept optional width/height parameters
- Calculate dynamic quadrantWidth and quadrantHeight
- Scale displayLimit based on available height
- Add minimum dimension constraints (e.g., 30x8 per quadrant)
- Handle edge case where no window size received yet (fallback to defaults)

### Example calculation:
```
Terminal: 120x40
Header: 3 lines
Urgent label: 2 lines
Matrix border: 2 lines (top + bottom)
Horizontal divider: 1 line
Margins: 2 lines

Available for quadrants: 40 - 3 - 2 - 2 - 1 - 2 = 30 lines
Height per quadrant: 30 / 2 = 15 lines

Available width: 120 - 4 (border) - 1 (divider) = 115
Width per quadrant: 115 / 2 = 57 characters
```

## Out of Scope
- Responsive reflow during resize (just re-render with new dimensions)
- Horizontal scrolling for overflow
- Collapsible quadrants
