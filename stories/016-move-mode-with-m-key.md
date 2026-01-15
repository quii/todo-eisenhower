# Story 016: Move Mode with 'm' Key

## Status
Completed

## Context
Currently, moving todos between quadrants uses Shift+1-4 shortcuts. However, Shift+number produces different symbols on different keyboard layouts (e.g., Shift+3 = £ on UK keyboards, # on US keyboards), making the feature broken for non-US layouts.

We need a keyboard-layout-independent way to move todos between quadrants.

## Goal
Replace Shift+1-4 shortcuts with an interactive move mode triggered by the 'm' key, which works consistently across all keyboard layouts.

## Acceptance Criteria

### Scenario 1: Enter move mode with 'm' key
**Given** I am viewing a quadrant with a todo selected
**When** I press 'm'
**Then** I should enter move mode
**And** I should see an overlay showing:
```
Move to quadrant:
  1. DO FIRST
  2. SCHEDULE
  3. DELEGATE
  4. ELIMINATE
Press ESC to cancel
```

### Scenario 2: Select destination quadrant
**Given** I am in move mode
**When** I press '2' (for SCHEDULE)
**Then** the selected todo should move to the SCHEDULE quadrant
**And** the priority should change to B
**And** move mode should exit
**And** I should return to normal focus mode

### Scenario 3: Cancel move mode with ESC
**Given** I am in move mode
**When** I press ESC
**Then** move mode should exit
**And** the todo should remain in its original quadrant
**And** I should return to normal focus mode

### Scenario 4: Move to each quadrant
**Given** I am in move mode
**When** I press 1, 2, 3, or 4
**Then** the todo should move to:
- 1 → DO FIRST (priority A)
- 2 → SCHEDULE (priority B)
- 3 → DELEGATE (priority C)
- 4 → ELIMINATE (priority D)

### Scenario 5: Moving to current quadrant is no-op
**Given** I am viewing DO FIRST with a todo selected
**And** I am in move mode
**When** I press '1'
**Then** the todo should remain in DO FIRST
**And** move mode should exit

### Scenario 6: Move mode only available in focus mode
**Given** I am in overview mode
**When** I press 'm'
**Then** nothing should happen (no move mode)

### Scenario 7: Move mode only available when todo selected
**Given** I am viewing an empty quadrant
**When** I press 'm'
**Then** nothing should happen (no move mode)

### Scenario 8: Help text shows 'm' for move
**Given** I am viewing a quadrant in focus mode
**Then** the help text should show "m to move"
**And** should NOT show "Shift+1-4 to move"

## Technical Notes
- Add `moveMode bool` field to Model
- Update keyboard handling in Update() to handle 'm' key
- Create RenderMoveOverlay() function to render the move mode UI
- Overlay should be centered and use lipgloss styling
- In move mode, only handle 1-4 and ESC keys, ignore other input
- Remove old Shift+1-4 handling code

## Design Decisions
- Use 'm' for move (mnemonic and available key)
- Show all 4 quadrants even if moving within same quadrant (consistency)
- ESC cancels (consistent with other cancel behavior)
- Move mode is modal - blocks other input while active
- Overlay uses lipgloss border and centered positioning

## Dependencies
- Story 012 (moving todos) - provides the underlying move functionality

## Out of Scope
- Batch moving multiple todos (future enhancement)
- Keyboard shortcuts for move destinations beyond 1-4
- Visual preview of where todo will move
