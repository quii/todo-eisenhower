# Story 020: Edit Todo

As a user
I want to edit an existing todo
So that I can update its description and tags without recreating it

## Acceptance Criteria

### Scenario: Press 'e' to enter edit mode
Given I am in focus mode with a selected todo
When I press 'e'
Then I should enter edit mode
And the input field should be pre-filled with the current description and tags
And I should see the tag reference panel
And the header should show "Edit Todo"

### Scenario: Edit todo description
Given I have a todo "Buy milk +shopping @store"
And I am in edit mode for that todo
When I clear the input and type "Buy eggs +shopping @store"
And I press Enter
Then the todo description should be updated to "Buy eggs"
And the todo should still have project tag "shopping"
And the todo should still have context tag "store"

### Scenario: Edit and remove tags
Given I have a todo "Review code +WebApp @computer"
And I am in edit mode for that todo
When I clear the input and type "Review code"
And I press Enter
Then the todo description should be "Review code"
And the todo should have no project tags
And the todo should have no context tags

### Scenario: Edit and add tags
Given I have a todo "Fix bug"
And I am in edit mode for that todo
When I clear the input and type "Fix bug +WebApp @computer"
And I press Enter
Then the todo description should be "Fix bug"
And the todo should have project tag "WebApp"
And the todo should have context tag "computer"

### Scenario: Preserve creation date
Given I have a todo created on 2026-01-15
And I am in edit mode for that todo
When I type "Updated description"
And I press Enter
Then the todo's creation date should still be 2026-01-15

### Scenario: Preserve completion date for completed todos
Given I have a completed todo "Task done" completed on 2026-01-18
And I am in edit mode for that todo
When I type "Task done with more details"
And I press Enter
Then the todo should still be completed
And the completion date should still be 2026-01-18

### Scenario: Cancel edit with ESC
Given I am in edit mode
When I press ESC
Then I should exit edit mode
And the todo should remain unchanged
And I should return to focus mode

### Scenario: Edit mode only available in focus mode
Given I am in overview mode
When I press 'e'
Then I should remain in overview mode
And edit mode should not activate

### Scenario: Edit mode only available when todo selected
Given I am in focus mode on an empty quadrant
When I press 'e'
Then edit mode should not activate

### Scenario: Auto-complete shows existing tags
Given there are existing todos with tags "+WebApp", "+Mobile", "@computer"
And I am in edit mode
Then the tag reference should show these existing tags
And I can use them for auto-completion

### Scenario: Edit maintains todo priority
Given I have a priority A todo "Important task +project"
And I am in edit mode for that todo
When I type "Updated important task +project"
And I press Enter
Then the todo should still have priority A

## Technical Notes

- Edit reuses the same input component as Add
- The input field is pre-populated with: `description +projects @contexts`
- Completion status and dates are preserved
- Creation date is preserved
- Priority is preserved (determined by which quadrant it's in)
- Empty input should probably be rejected (or cancel the edit)
