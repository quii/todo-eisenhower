# Story 008: Add Todo with Tag Reference

## User Story

**As a** user viewing a focused quadrant
**I want to** add a new todo with project and context tags
**So that** I can quickly capture tasks with proper categorization

## Acceptance Criteria

### Scenario 1: Press 'a' to enter input mode in focused quadrant
**Given** I am viewing the DO FIRST quadrant in focus mode
**When** I press 'a'
**Then** an input field appears at the bottom
**And** the input field is focused and ready for text entry
**And** existing project and context tags are displayed below the input
**And** help text shows "Enter to save • ESC to cancel"

### Scenario 2: Add a simple todo without tags
**Given** I am in input mode in the DO FIRST quadrant
**When** I type "Fix critical bug"
**And** I press Enter
**Then** the new todo appears in the DO FIRST quadrant
**And** the todo has priority (A)
**And** the input mode is exited
**And** the input field is cleared

### Scenario 3: Add todo with project tags
**Given** I am in input mode in the SCHEDULE quadrant
**When** I type "Plan sprint +WebApp +Mobile"
**And** I press Enter
**Then** the new todo appears in the SCHEDULE quadrant
**And** the todo has priority (B)
**And** the project tags +WebApp and +Mobile are displayed in color
**And** the tags match the colors of existing tags with those names

### Scenario 4: Add todo with context tags
**Given** I am in input mode in the DELEGATE quadrant
**When** I type "Reply to emails @phone @office"
**And** I press Enter
**Then** the new todo appears in the DELEGATE quadrant
**And** the todo has priority (C)
**And** the context tags @phone and @office are displayed in color

### Scenario 5: Add todo with mixed tags
**Given** I am in input mode in the DO FIRST quadrant
**When** I type "Deploy to production +WebApp @computer @work"
**And** I press Enter
**Then** the new todo has both project and context tags
**And** all tags are displayed with consistent colors

### Scenario 6: Cancel input with ESC
**Given** I am in input mode with text typed
**When** I press ESC
**Then** input mode is exited
**And** no new todo is created
**And** the typed text is discarded

### Scenario 7: Tag reference display shows existing tags
**Given** the matrix contains todos with tags +WebApp +Mobile @computer @phone @office
**When** I enter input mode in any focused quadrant
**Then** I see "Projects: +WebApp +Mobile" displayed below the input
**And** I see "Contexts: @computer @office @phone" displayed below that
**And** the tags are displayed in their consistent colors

### Scenario 8: Empty tag reference when no tags exist
**Given** the matrix contains todos with no tags
**When** I enter input mode
**Then** I see "Projects: (none)" displayed
**And** I see "Contexts: (none)" displayed

### Scenario 9: Input only available in focus mode
**Given** I am viewing the overview (all quadrants)
**When** I press 'a'
**Then** nothing happens
**And** input mode is not activated
**And** help text suggests pressing 1/2/3/4 to focus first

### Scenario 10: Auto-assign priority from quadrant
**Given** I am in focus mode on different quadrants
**When** I add a todo in DO FIRST, it gets priority (A)
**And** I add a todo in SCHEDULE, it gets priority (B)
**And** I add a todo in DELEGATE, it gets priority (C)
**And** I add a todo in ELIMINATE, it gets priority (D)

### Scenario 11: New tags are accepted
**Given** I am in input mode
**And** existing projects are +WebApp +Mobile
**When** I type "Build API +Backend"
**And** I press Enter
**Then** the todo is created with the new +Backend tag
**And** +Backend is displayed with a consistent color

## Technical Notes

### Implementation Approach

1. **Add Bubble Tea textinput component:**
   ```go
   import "github.com/charmbracelet/bubbles/textinput"

   type Model struct {
       // ... existing fields
       inputMode    bool
       input        textinput.Model
       allProjects  []string
       allContexts  []string
   }
   ```

2. **Extract all tags from matrix on init:**
   - Iterate through all quadrants
   - Use existing `extractTags()` logic from parser package
   - Store unique tag lists in model

3. **Update keyboard handling:**
   - Press 'a' in focus mode → enter input mode
   - Delegate input events to textinput.Model
   - Enter → parse input, create todo, save file, exit input mode
   - ESC → exit input mode without saving

4. **Rendering input mode:**
   - Show focused quadrant title
   - Show todos (limited by display height)
   - Show separator line
   - Show input field with prompt
   - Show tag reference lists (colorized)
   - Show help text: "Enter to save • ESC to cancel"

5. **Priority assignment:**
   ```go
   func (m Model) currentQuadrantPriority() todo.Priority {
       switch m.viewMode {
       case FocusDoFirst:
           return todo.PriorityA
       case FocusSchedule:
           return todo.PriorityB
       case FocusDelegate:
           return todo.PriorityC
       case FocusEliminate:
           return todo.PriorityD
       default:
           return todo.PriorityNone
       }
   }
   ```

6. **File persistence:**
   - Use existing todo.txt format
   - Append new todo to file
   - Reload matrix after save

7. **Tag extraction helper:**
   - Reuse `projectTagPattern` and `contextTagPattern` from render.go
   - Extract all unique tags from all todos across all quadrants

### Edge Cases

- Empty input (just Enter) → no todo created
- Whitespace only → no todo created
- Duplicate tags in input → deduplicated
- Very long input → truncate or wrap display
- Tag reference list too long → truncate with "... and X more"

## Definition of Done

- [ ] All acceptance test scenarios pass
- [ ] Input mode activates in focus mode only
- [ ] New todos are saved to file in todo.txt format
- [ ] Tags are colorized consistently with existing tags
- [ ] Priority is auto-assigned based on quadrant
- [ ] ESC properly cancels without creating todo
- [ ] Tag reference lists are displayed with colors
- [ ] Manual testing with real todo.txt file works
