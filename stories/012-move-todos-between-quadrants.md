# Story 012: Move todos between quadrants

**As a** user managing my task priorities
**I want to** move todos between quadrants by changing their priority
**So that** I can quickly reorganize tasks as circumstances change

## Context

In the Eisenhower matrix, moving a todo between quadrants means changing its priority:
- Priority A → DO FIRST (Quadrant 1)
- Priority B → SCHEDULE (Quadrant 2)
- Priority C → DELEGATE (Quadrant 3)
- Priority D → ELIMINATE (Quadrant 4)

The number keys (1-4) should be context-aware:
- In Overview mode: Focus on that quadrant (current behavior)
- In Focus mode: Move selected todo to that quadrant (changes priority)

## Acceptance Criteria

### Scenario: Move todo from DO FIRST to SCHEDULE
```gherkin
Given I have a todo "(A) Review quarterly goals"
And I am in focus mode viewing DO FIRST
When I select the todo
And I press "2"
Then the todo priority changes to B
And the todo moves to SCHEDULE quadrant
And the file is updated with "(B) Review quarterly goals"
```

### Scenario: Move todo from DELEGATE to DO FIRST
```gherkin
Given I have a todo "(C) Update documentation"
And I am in focus mode viewing DELEGATE
When I select the todo
And I press "1"
Then the todo priority changes to A
And the todo moves to DO FIRST quadrant
And the file is updated with "(A) Update documentation"
```

### Scenario: Move todo to ELIMINATE (priority D)
```gherkin
Given I have a todo "(B) Optional feature idea"
And I am in focus mode viewing SCHEDULE
When I select the todo
And I press "4"
Then the todo priority changes to D
And the todo moves to ELIMINATE quadrant
And the file is updated with "(D) Optional feature idea"
```

### Scenario: Moving todo adjusts selection
```gherkin
Given I have 3 todos in DO FIRST
And the second todo is selected
When I press "2" to move it to SCHEDULE
Then the selection moves to the next remaining todo in DO FIRST
And the moved todo is no longer visible in DO FIRST
```

### Scenario: Moving last todo in quadrant returns to overview
```gherkin
Given I have only 1 todo in DELEGATE
And I am in focus mode viewing DELEGATE
When I press "1" to move it to DO FIRST
Then I automatically return to overview mode
And the moved todo is visible in DO FIRST quadrant
```

### Scenario: Moving todo preserves tags and completion status
```gherkin
Given I have a completed todo "x 2025-01-10 (A) Fix bug +WebApp @computer"
And I am in focus mode viewing DO FIRST
When I press "3" to move it to DELEGATE
Then the file contains "x 2025-01-10 (C) Fix bug +WebApp @computer"
And all tags are preserved
And completion status is preserved
```

### Scenario: Pressing current quadrant number does nothing
```gherkin
Given I have a todo "(B) Plan sprint"
And I am in focus mode viewing SCHEDULE (quadrant 2)
When I press "2"
Then nothing happens
And the todo remains in SCHEDULE
```

### Scenario: Number keys still focus quadrants in overview mode
```gherkin
Given I am in overview mode
When I press "1"
Then I enter focus mode on DO FIRST
And I can move todos using number keys
```

## Implementation Notes

- Update keyboard handler to check viewMode when handling 1-4 keys
- Add `ChangePriority()` method to Todo domain
- Add `MoveTodo()` or similar to Matrix domain
- Update help text in focus mode to indicate "Press 1/2/3/4 to move here"
- When moving a todo:
  - Update priority
  - Save entire matrix to file
  - Adjust selection index (move to next todo, or return to overview if none left)
- If pressing the number for current quadrant, do nothing (no-op)
- Preserve all other todo properties (description, tags, completion status)

## UI/UX Considerations

- Visual feedback: Brief flash or message showing "Moved to SCHEDULE"?
- Help text should indicate the dual purpose of number keys based on mode
- Moving todos should feel instant and responsive
- Consider adding undo functionality in a future story

## Out of Scope

- Moving multiple todos at once
- Drag-and-drop interface
- Undo/redo functionality
- Visual animations or transitions
