# Story 007: Quadrant Focus Mode + Remove Emojis

## User Story
As a user, I want to focus on one quadrant at a time so that I can see all my todos in that priority level and work through them systematically. I also want a cleaner interface without emojis.

## Background
Currently the app shows all 4 quadrants simultaneously, which gives a good overview but limits how many todos are visible per quadrant. When actually working, users need to focus on one priority level at a time.

The app should have two viewing modes:
- **Overview mode**: All 4 quadrants visible (current behavior)
- **Focus mode**: Single quadrant fills the screen, showing many more todos

Additionally, the current UI uses emojis (🔥📅👥🗑️📄) which should be removed for a cleaner, more professional appearance.

## Acceptance Criteria

```gherkin
Feature: Quadrant Focus Mode

  Scenario: No emojis in quadrant titles
    Given I am viewing the application
    Then the DO FIRST title does not contain "🔥"
    And the SCHEDULE title does not contain "📅"
    And the DELEGATE title does not contain "👥"
    And the ELIMINATE title does not contain "🗑️"
    And the file header does not contain "📄"

  Scenario: Focus on DO FIRST quadrant
    Given I have 20 todos in the DO FIRST quadrant
    And I am in overview mode
    When I press "1"
    Then the DO FIRST quadrant fills the entire screen
    And I can see significantly more todos than in overview mode
    And the quadrant title "DO FIRST" is prominently displayed
    And the file path header is shown at the top
    And I see "Press ESC for overview • Press 2/3/4 for other quadrants" at the bottom

  Scenario: Focus on SCHEDULE quadrant
    Given I am in overview mode
    When I press "2"
    Then the SCHEDULE quadrant fills the entire screen
    And the quadrant title "SCHEDULE" is displayed
    And I see the help text at the bottom

  Scenario: Focus on DELEGATE quadrant
    Given I am in overview mode
    When I press "3"
    Then the DELEGATE quadrant fills the entire screen
    And the quadrant title "DELEGATE" is displayed

  Scenario: Focus on ELIMINATE quadrant
    Given I am in overview mode
    When I press "4"
    Then the ELIMINATE quadrant fills the entire screen
    And the quadrant title "ELIMINATE" is displayed

  Scenario: Return to overview with ESC
    Given I am in focus mode viewing DO FIRST
    When I press "ESC"
    Then I return to the 4-quadrant overview mode
    And all quadrants are visible again

  Scenario: Return to overview with 'q' does NOT work
    Given I am in focus mode viewing DO FIRST
    When I press "q"
    Then the application quits
    And I do not return to overview

  Scenario: Jump between quadrants in focus mode
    Given I am in focus mode viewing DO FIRST (quadrant 1)
    When I press "2"
    Then I switch to focus mode viewing SCHEDULE
    And the transition is immediate
    When I press "4"
    Then I switch to focus mode viewing ELIMINATE
    When I press "1"
    Then I switch back to focus mode viewing DO FIRST

  Scenario: Display limit scales in focus mode
    Given I have 50 todos in the SCHEDULE quadrant
    And my terminal is 120x40
    When I press "2" to focus on SCHEDULE
    Then I can see approximately 30-35 todos (much more than the ~11 in overview)
    And the remaining todos show "... and X more"

  Scenario: Focus mode adapts to terminal size
    Given I am in focus mode viewing DO FIRST
    And my terminal is 200x60
    Then I can see approximately 50+ todos
    When I resize my terminal to 100x30
    Then the view adjusts to show fewer todos accordingly

  Scenario: Empty quadrant in focus mode
    Given the ELIMINATE quadrant is empty
    When I press "4" to focus on ELIMINATE
    Then I see "ELIMINATE" as the title
    And I see "(no tasks)" styled appropriately
    And I see the help text at the bottom

  Scenario: Completed todos visible in focus mode
    Given I have completed todos in the DO FIRST quadrant
    When I press "1" to focus
    Then completed todos are shown with strikethrough
    And they are visually distinct from active todos
```

## Technical Notes

### Emoji removal:
Emojis to remove from `adapters/ui/render.go`:
- Quadrant titles: `"🔥 DO FIRST"` → `"DO FIRST"`
- Quadrant titles: `"📅 SCHEDULE"` → `"SCHEDULE"`
- Quadrant titles: `"👥 DELEGATE"` → `"DELEGATE"`
- Quadrant titles: `"🗑️  ELIMINATE"` → `"ELIMINATE"`
- File header: `"📄 File: "` → `"File: "`

These are the only emojis in the codebase. Unicode symbols (•, ✓, ←, →) are kept as they're functional typography.

### State management:
- Add `ViewMode` enum: `Overview`, `FocusDoFirst`, `FocusSchedule`, `FocusDelegate`, `FocusEliminate`
- Add `viewMode` field to Model
- Default to `Overview` on startup

### Keyboard handling in Update:
```go
case tea.KeyMsg:
    switch msg.String() {
    case "1":
        m.viewMode = FocusDoFirst
    case "2":
        m.viewMode = FocusSchedule
    case "3":
        m.viewMode = FocusDelegate
    case "4":
        m.viewMode = FocusEliminate
    case "esc":
        m.viewMode = Overview
    case "q", "ctrl+c":
        return m, tea.Quit
    }
```

### View rendering:
- Add new function: `RenderFocusedQuadrant(todos []todo.Todo, title, filePath, helpText string, width, height int)`
- Calculate display limit based on full screen height (reserve space for header, title, help footer)
- Use lipgloss.Place to center the focused quadrant content
- Render help text at bottom: `"Press ESC for overview • Press 2/3/4 for other quadrants"`

### Display limit calculation in focus mode:
- Full screen height minus:
  - File header: 3 lines
  - Quadrant title: 2 lines
  - Help text footer: 2 lines
  - Margins: 2 lines
- Example: 40 line terminal = ~31 available lines for todos

### Styling:
- Keep existing color scheme and styling
- Help text should be subtle (gray, italic, centered at bottom)
- Quadrant title should be prominent (larger, centered, with color)

## Out of Scope
- Navigation within quadrant (selecting specific todos) - Story 008
- Editing todos - Future story
- Showing todo details - Future story
- Scrolling within focus mode - Future story (for now, "... and X more" is fine)

## Design Decisions
- 'q' still quits from focus mode (doesn't return to overview) - consistent behavior
- ESC is the explicit "go back" key
- Number keys work in both modes (jump to focus from overview, switch focus in focus mode)
- Focus mode centers content, doesn't left-align like overview quadrants
