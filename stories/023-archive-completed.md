# Story 023: Archive Completed Todos

As a user
I want to archive completed todos to a done.txt file
So that I can keep my main todo.txt focused on active work while preserving history

## Background

Following the todo.txt convention, completed todos should be movable to a `done.txt` file
in the same directory as `todo.txt`. From a UX perspective, archiving is like "deleting"
but the data is preserved in a separate file rather than lost.

Archiving is manual and per-todo (not bulk) - press 'd' on a completed todo to archive it.

## Acceptance Criteria

### Scenario: Archive a completed todo from focused quadrant
Given I have a completed todo "x 2026-01-20 Buy milk" in the Do First quadrant
And I am in focus mode on the Do First quadrant
When I select the completed todo
And I press 'd' (shift+a)
Then the todo should be removed from todo.txt
And the todo should be appended to done.txt in the same directory
And the matrix should update to reflect the removal

### Scenario: Cannot archive uncompleted todos
Given I have an uncompleted todo "(A) Active task" in the Do First quadrant
And I am in focus mode on the Do First quadrant
When I select the uncompleted todo
And I press 'd'
Then nothing should happen
Or a message should indicate "Only completed todos can be archived"
And the todo should remain in todo.txt

### Scenario: Archive from any quadrant
Given I have completed todos in different quadrants:
  | Quadrant  | Todo                               |
  | Do First  | x 2026-01-20 Task A               |
  | Schedule  | x 2026-01-19 Task B               |
  | Delegate  | x 2026-01-18 Task C               |
  | Eliminate | x 2026-01-17 Task D               |
When I focus on each quadrant and archive the completed todos
Then all archived todos should be moved to done.txt
And all should be removed from their respective quadrants

### Scenario: done.txt is created if it doesn't exist
Given done.txt does not exist
And I have a completed todo in the matrix
When I archive the completed todo
Then done.txt should be created in the same directory as todo.txt
And the archived todo should be written to it

### Scenario: done.txt is appended to if it exists
Given done.txt already contains:
  """
  x 2026-01-10 Previous task
  x 2026-01-11 Another old task
  """
And I archive a completed todo "x 2026-01-20 New task"
Then done.txt should contain:
  """
  x 2026-01-10 Previous task
  x 2026-01-11 Another old task
  x 2026-01-20 New task
  """
And the original content should be preserved

### Scenario: Archived todo preserves all metadata
Given I have a completed todo:
  """
  x 2026-01-20 2026-01-15 (A) Complex task +project @context due:2026-01-19
  """
When I archive it
Then done.txt should contain the exact same line
Including completion date, creation date, priority, projects, contexts, and due date

### Scenario: Cannot archive from overview mode
Given I am in overview mode
When I press 'd'
Then nothing should happen
Because archiving requires selecting a specific todo (focus mode only)

### Scenario: Multiple archives in one session
Given I have 3 completed todos in the Do First quadrant
When I archive the first todo
And I archive the second todo
And I archive the third todo
Then all 3 todos should be appended to done.txt in order
And the Do First quadrant should have 3 fewer todos

### Scenario: File errors are handled gracefully
Given done.txt exists but is read-only
When I attempt to archive a completed todo
Then an error message should be displayed
And the todo should remain in todo.txt (not lost)

### Scenario: Viewing archived todos (out of scope but documented)
Given I have archived todos in done.txt
When I want to view them
Then I can run `eisenhower ~/done.txt` in a separate terminal
Because done.txt is just another todo.txt file
But there's no in-app archive viewer (by design)

## Technical Notes

- Archive is triggered by pressing 'd' in focus mode
- Only works on completed todos (those with `x` prefix)
- The selected todo is identified by the current table row selection
- done.txt location: same directory as the input todo.txt file
- done.txt format: standard todo.txt format (append with newline)
- File operations should use the existing file adapter for consistency
- After archive, the matrix should be reloaded from todo.txt
- Cursor/selection should move to the next todo after archive (or stay if it was last)

## Implementation Considerations

- Add 'A' key handler in focus mode that checks if todo is completed
- Add an `Archive` method to the file repository that:
  1. Reads current done.txt (if exists)
  2. Appends the todo line
  3. Writes done.txt
  4. Removes the todo from todo.txt
- Ensure atomic operation (if archive succeeds but removal fails, rollback)
- Add confirmation dialog? (NO - keep it quick, it's not destructive)

## Future Considerations (NOT in this story)

- Bulk archive (archive all completed todos in quadrant)
- Auto-archive completed todos older than N days
- Un-archive functionality
- Archive statistics in inventory
- Configurable archive file name/location
