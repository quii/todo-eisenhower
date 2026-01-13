# Story 009: Tag Autocomplete

## User Story

**As a** user adding a new todo
**I want** autocomplete suggestions when typing project and context tags
**So that** I can quickly reuse existing tags without typos and discover available tags

## Acceptance Criteria

### Scenario 1: Trigger autocomplete with +
**Given** I am in input mode with existing project tags +WebApp, +Mobile, +Backend
**When** I type "Deploy feature +"
**Then** a suggestion box appears below the input
**And** all 3 project tags are shown as suggestions
**And** the first suggestion is highlighted
**And** tags are shown in their consistent colors

### Scenario 2: Filter suggestions as I type
**Given** autocomplete is showing +WebApp, +Mobile, +Backend
**When** I continue typing "Deploy feature +Web"
**Then** only +WebApp is shown as a suggestion
**And** it remains highlighted

### Scenario 3: Navigate suggestions with arrow keys
**Given** autocomplete is showing +WebApp, +Mobile, +Backend
**When** I press Down arrow
**Then** the second suggestion (+Mobile) is highlighted
**When** I press Down arrow again
**Then** the third suggestion (+Backend) is highlighted
**When** I press Down arrow again
**Then** the first suggestion (+WebApp) is highlighted (wrap around)
**When** I press Up arrow
**Then** the third suggestion (+Backend) is highlighted (wrap backwards)

### Scenario 4: Complete tag with Tab
**Given** autocomplete is showing +WebApp, +Mobile with +WebApp highlighted
**When** I press Tab
**Then** the input value becomes "Deploy feature +WebApp "
**And** autocomplete suggestions are dismissed
**And** cursor is positioned after the completed tag with a space

### Scenario 5: Complete tag with Enter
**Given** autocomplete is showing +WebApp, +Mobile with +Mobile highlighted
**When** I press Enter while a suggestion is selected
**Then** the input value becomes "Deploy feature +Mobile "
**And** autocomplete suggestions are dismissed
**And** the cursor is positioned for continued typing
**And** the todo is NOT saved (Enter only saves when suggestions are not shown)

### Scenario 6: Dismiss suggestions with ESC
**Given** autocomplete is showing suggestions
**When** I press ESC
**Then** autocomplete suggestions are dismissed
**And** I remain in input mode
**And** my typed text is preserved

### Scenario 7: Autocomplete context tags with @
**Given** I am in input mode with existing context tags @computer, @phone, @office
**When** I type "Reply to emails @"
**Then** autocomplete shows all 3 context tags
**When** I type "Reply to emails @p"
**Then** autocomplete shows @phone only
**When** I press Tab
**Then** the input value becomes "Reply to emails @phone "

### Scenario 8: Multiple tags in one todo
**Given** I am in input mode
**When** I type "Deploy +WebApp @"
**Then** autocomplete shows context tags (not project tags)
**When** I complete with @computer
**And** I type " +"
**Then** autocomplete shows project tags again
**And** I can add multiple tags to the same todo

### Scenario 9: No suggestions available
**Given** I am in input mode
**When** I type "Deploy +xyz"
**And** no existing tags match "xyz"
**Then** autocomplete shows "(no matches - press Space to create new tag)"
**When** I press Space
**Then** autocomplete is dismissed
**And** "+xyz " is in the input

### Scenario 10: Continue typing dismisses autocomplete
**Given** autocomplete is showing suggestions for "+Web"
**When** I press Space
**Then** autocomplete suggestions are dismissed
**And** "+Web " remains in the input
**And** I can continue typing normally

### Scenario 11: Autocomplete styling matches UI
**Given** autocomplete is visible
**Then** the suggestion box has a subtle border
**And** the selected suggestion has a highlighted background
**And** tags are shown in their consistent hash colors
**And** the box appears directly below the input field

### Scenario 12: Case-insensitive matching
**Given** existing tag +WebApp
**When** I type "+web"
**Then** +WebApp appears as a suggestion
**When** I complete it
**Then** it inserts +WebApp (preserves original case)

## Technical Notes

### Implementation Approach

1. **Autocomplete State in Model:**
   ```go
   type Model struct {
       // ... existing fields
       showSuggestions  bool
       suggestions      []string
       selectedSuggestion int
       currentPrefix    string  // e.g., "+Web"
   }
   ```

2. **Detect Trigger Characters:**
   - Monitor textinput value changes
   - Find last `+` or `@` before cursor
   - Extract partial tag (e.g., "+Web")
   - Filter tags by prefix match (case-insensitive)

3. **Keyboard Handling in Input Mode:**
   ```go
   case tea.KeyMsg:
       if m.showSuggestions {
           switch msg.String() {
           case "up":
               m.selectedSuggestion = (m.selectedSuggestion - 1 + len(m.suggestions)) % len(m.suggestions)
           case "down":
               m.selectedSuggestion = (m.selectedSuggestion + 1) % len(m.suggestions)
           case "tab":
               m = m.completeSuggestion()
           case "enter":
               if m.showSuggestions {
                   m = m.completeSuggestion()
                   return m, nil // Don't save todo
               }
               // Otherwise save todo
           case "esc":
               m.showSuggestions = false
           }
       }
   ```

4. **Filtering Logic:**
   ```go
   func filterTags(tags []string, prefix string) []string {
       prefix = strings.ToLower(prefix)
       var matches []string
       for _, tag := range tags {
           if strings.HasPrefix(strings.ToLower(tag), prefix) {
               matches = append(matches, tag)
           }
       }
       return matches
   }
   ```

5. **Tag Completion:**
   - Replace partial tag with selected suggestion
   - Add space after completed tag
   - Update cursor position
   - Dismiss suggestions

6. **Rendering Suggestions:**
   - Create new `renderAutocomplete()` function
   - Show up to 5 suggestions
   - Highlight selected suggestion with background color
   - Position below input field
   - Show "(no matches)" message when empty

7. **Update RenderFocusedQuadrantWithInput:**
   - Check if `showSuggestions` is true
   - Render autocomplete box between input and tag reference
   - Adjust display limit to account for autocomplete height

### Edge Cases

- Cursor in middle of word: don't trigger autocomplete
- Multiple spaces after tag: don't trigger
- Tag at start of input: "+WebApp" should trigger
- Backspace removes trigger character: dismiss autocomplete
- Very long tag list: show scrollable list (future enhancement)
- Duplicate tags: already handled by hash set in tag extraction

### Performance Considerations

- Filtering happens on every keystroke (acceptable for <100 tags)
- Case-insensitive comparison is simple and fast
- No need for fuzzy matching yet (can add later)

### Future Enhancements (not in this story)

- Fuzzy matching (e.g., "wap" matches "WebApp")
- Show tag usage count
- Recently used tags first
- Mouse click to select suggestion

## Definition of Done

- [ ] All 12 acceptance test scenarios pass
- [ ] Autocomplete triggers on `+` and `@`
- [ ] Arrow keys navigate suggestions
- [ ] Tab and Enter complete suggestions
- [ ] ESC dismisses autocomplete
- [ ] Case-insensitive filtering works
- [ ] Multiple tags can be autocompleted in one input
- [ ] No matches message shows when appropriate
- [ ] Autocomplete integrates cleanly with existing input mode
- [ ] Manual testing with real todo.txt file works smoothly
