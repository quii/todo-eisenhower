# Story 011: Mark Todos as Complete

## User Story

**As a** user viewing a focused quadrant
**I want** to navigate through todos and mark them as complete
**So that** I can track my progress and update my todo list without leaving the app

## Acceptance Criteria

### Scenario 1: Navigate todos with arrow keys
**Given** I am in focus mode on DO FIRST with 3 todos
**When** I view the quadrant
**Then** the first todo is highlighted
**When** I press Down arrow (or 's')
**Then** the second todo is highlighted
**When** I press Down arrow again
**Then** the third todo is highlighted
**When** I press Down arrow again
**Then** the first todo is highlighted (wrap around)
**When** I press Up arrow (or 'w')
**Then** the third todo is highlighted (wrap backwards)

### Scenario 2: Mark todo as complete
**Given** I am in focus mode with todo "Fix bug +WebApp" highlighted
**When** I press Enter
**Then** the todo is marked as complete (shows ✓ instead of •)
**And** the todo.txt file is updated with "x YYYY-MM-DD (A) Fix bug +WebApp"
**And** the todo remains in the same quadrant (completion doesn't move it)
**And** the selection moves to the next todo

### Scenario 3: Unmark completed todo
**Given** I am in focus mode with a completed todo highlighted
**When** I press Enter
**Then** the todo is marked as incomplete (shows • instead of ✓)
**And** the todo.txt file is updated to remove the "x YYYY-MM-DD" prefix
**And** the selection remains on the same todo

### Scenario 4: Selection persists when switching quadrants
**Given** I am in DO FIRST with the second todo selected
**When** I press ESC to return to overview
**And** I press 2 to focus on SCHEDULE
**Then** the first todo in SCHEDULE is selected (reset per quadrant)

### Scenario 5: Selection state not shown in overview mode
**Given** I am in focus mode with a todo selected
**When** I press ESC to return to overview
**Then** no todos are highlighted in overview mode
**And** todos still show completion status (✓ or •)

### Scenario 6: Entering input mode preserves selection
**Given** I am in focus mode with the third todo selected
**When** I press 'a' to enter input mode
**And** I press ESC to cancel input
**Then** the third todo is still selected

### Scenario 7: After adding todo, selection resets to first
**Given** I am in focus mode with the third todo selected
**When** I press 'a' and add a new todo
**Then** selection resets to the first todo in the list

### Scenario 8: Empty quadrant has no selection
**Given** I focus on a quadrant with no todos
**Then** no todo is highlighted
**And** pressing Enter does nothing
**And** navigation keys do nothing

## Technical Notes

### Implementation Approach

1. **Add selection state to Model:**
   ```go
   type Model struct {
       // ... existing fields
       selectedIndex  int  // index of selected todo in current quadrant
   }
   ```

2. **Keyboard handling in focus mode (not input mode):**
   ```go
   case tea.KeyMsg:
       if m.viewMode != Overview && !m.inputMode {
           switch msg.String() {
           case "down", "s":
               m = m.moveSelectionDown()
           case "up", "w":
               m = m.moveSelectionUp()
           case "enter":
               m = m.toggleCompletion()
           }
       }
   ```

3. **Todo selection navigation:**
   ```go
   func (m Model) moveSelectionDown() Model {
       todos := m.currentQuadrantTodos()
       if len(todos) == 0 {
           return m
       }
       m.selectedIndex = (m.selectedIndex + 1) % len(todos)
       return m
   }
   ```

4. **Toggle completion:**
   ```go
   func (m Model) toggleCompletion() Model {
       todos := m.currentQuadrantTodos()
       if len(todos) == 0 || m.selectedIndex >= len(todos) {
           return m
       }

       selectedTodo := todos[m.selectedIndex]
       updatedTodo := selectedTodo.ToggleCompletion()

       // Update in matrix
       m.matrix = m.matrix.UpdateTodo(selectedTodo, updatedTodo)

       // Write entire matrix back to file
       if m.writer != nil {
           _ = usecases.SaveAllTodos(m.writer, m.matrix)
       }

       // Move selection to next todo (or stay if unmarking)
       if updatedTodo.IsCompleted() {
           m.selectedIndex = (m.selectedIndex + 1) % len(todos)
       }

       return m
   }
   ```

5. **Rendering selected todo:**
   - Add `selectedIndex int` parameter to RenderFocusedQuadrant
   - When rendering todo at index i, check if i == selectedIndex
   - Apply highlight style (background color, bold, or indicator)
   - Example: `"> • Task"` or background highlight

6. **Domain changes needed:**
   - Add `ToggleCompletion()` method to Todo
   - Add `UpdateTodo(old, new Todo)` method to Matrix
   - Add `SaveAllTodos(writer, matrix)` use case to write all todos

### Styling for Selection

Option 1 - Selection indicator:
```
  • Task one
> • Task two (selected)
  • Task three
```

Option 2 - Background highlight:
```
• Task one
• Task two (with background color)
• Task three
```

Prefer Option 2 for cleaner look.

### Edge Cases

- Empty quadrant: no selection, navigation does nothing
- Single todo: up/down wraps to same todo
- Switching to input mode: preserve selection
- After adding todo: reset to first
- Switching quadrants: each has independent selection state
- Completed todos stay in same quadrant (priority doesn't change)
- File write failures: handle gracefully, show error

### File Format

Completion marker follows todo.txt spec:
- Incomplete: `(A) Task description +tag @context`
- Complete: `x 2026-01-13 (A) Task description +tag @context`

### Performance Considerations

- Writing entire file on each toggle is acceptable for small files
- For larger files (>1000 todos), consider incremental updates (future)
- Matrix rebuild after write ensures consistency

## Definition of Done

- [ ] Can navigate todos with arrow keys and w/s in focus mode
- [ ] Can mark todos as complete with Enter
- [ ] Can unmark completed todos with Enter
- [ ] Visual highlight shows selected todo
- [ ] Completion persists to file with correct format
- [ ] Selection wraps around at boundaries
- [ ] Selection state independent per quadrant
- [ ] No selection shown in overview mode
- [ ] All acceptance tests pass
- [ ] Manual testing with real todo.txt file works smoothly
