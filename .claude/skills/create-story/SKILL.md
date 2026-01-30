# Create User Story

A skill for creating user stories with Gherkin scenarios in the stories/ directory.

## Instructions

You are helping the user create a new user story for this project. Follow this workflow:

### 1. Gather Initial Requirements

If the user provided a description with the skill invocation, use that. Otherwise, ask:
- "What feature or improvement would you like to add? Please describe the high-level problem or user need."

### 2. Ask Clarifying Questions

Ask ONE question at a time. Wait for the user's answer before asking the next question.

**Question Categories (in order of priority):**

1. **User Flow** - What triggers this feature? What should happen?
2. **Edge Cases** - What happens when [error/boundary condition]?
3. **Integration** - How does this interact with existing features?
4. **UI/UX** - Where does this appear? What feedback is shown?
5. **Scope** - Is this a small vertical slice or should it be broken down?

**Process:**
- Ask the most important question based on what you don't know yet
- Wait for the user's answer
- Based on their answer, decide the next most important question
- Continue until you have enough information to write clear Gherkin scenarios (typically 3-6 questions)
- If the user's answer reveals the feature is too large, ask if they want to break it into smaller stories

**When to stop asking questions:**
- You understand the main user flow
- You know the key edge cases
- You can write specific, testable scenarios
- The scope is clear and appropriately sized

### 3. Generate Story Number

Check the stories/ directory to find the next available story number:
- List files in stories/ directory
- Find the highest numbered story (e.g., 026-stale-tasks.md)
- Use the next number (e.g., 027)

### 4. Create Story File

Generate a story file with this structure:

```markdown
# Story NNN: [Feature Name]

As a [user type]
I want [capability]
So that [benefit/value]

## Background

[Provide context about why this feature is needed. Explain the current situation, pain points, or opportunities. This should help future readers understand the motivation.]

## Acceptance Criteria

\`\`\`gherkin
Feature: [Feature Name]

  Scenario: [First scenario name]
    Given [initial context]
    And [additional context if needed]
    When [action taken]
    Then [expected outcome]
    And [additional outcome if needed]

  Scenario: [Second scenario name]
    Given [initial context]
    When [action taken]
    Then [expected outcome]

  [Additional scenarios as needed]
\`\`\`

## Technical Notes

[Optional section for implementation guidance, architectural considerations, or technical constraints]

### [Subsection if needed]

- Implementation detail 1
- Implementation detail 2

## Future Considerations (NOT in this story)

[Optional section for ideas that came up but are out of scope]

- Future enhancement 1
- Future enhancement 2
```

**Gherkin Formatting Rules:**
- Use proper indentation (2 spaces for Given/When/Then/And)
- Each scenario should be clearly named and focused
- Given = initial state/context
- When = action/trigger
- Then = expected outcome
- And = additional conditions or outcomes
- Keep scenarios focused and readable
- Use realistic examples with specific values
- Include edge cases and error conditions

**Naming Convention:**
- File name: `NNN-feature-name.md` (lowercase, hyphens, descriptive)
- Example: `027-bulk-delete-todos.md`

### 5. Save to stories/ Directory

Write the story file to `stories/NNN-feature-name.md`

### 6. Confirm Completion

Tell the user:
- The story number and filename
- A brief summary of what was created
- Suggest next steps: "Ready to implement this story? I can help you start with TDD by writing the first test."

## Example Interaction

**User:** `/create-story Add ability to bulk delete multiple todos`

**Agent:** I'll help you create a user story for bulk delete functionality.

First question: How does the user select multiple todos for deletion? (For example: checkboxes, keyboard multi-select, select mode, etc.)

**User:** Let's use checkboxes

**Agent:** Got it. Should there be a confirmation dialog before deletion, or should it delete immediately?

**User:** Yes, confirmation dialog with count of items

**Agent:** Makes sense. Can users bulk delete across different quadrants, or only within the currently focused quadrant?

**User:** Only within the focused quadrant

**Agent:** Understood. Can both completed and uncompleted todos be bulk deleted, or should completed items be excluded?

**User:** Both can be deleted

**Agent:** Perfect, I have enough to create the story. Great! I've created **Story 027: Bulk Delete Todos** in `stories/027-bulk-delete-todos.md`.

The story includes:
- 8 Gherkin scenarios covering selection, confirmation, cross-quadrant deletion, and edge cases
- Technical notes about implementation considerations
- Future enhancements (undo, keyboard shortcuts) marked as out of scope

Ready to implement this story? I can help you start with TDD by writing the first test.

## Important Notes

- **Follow project conventions**: Check claude.md for architecture patterns
- **Small vertical slices**: If the feature is too large, suggest breaking it into multiple stories
- **TDD mindset**: Write scenarios that can become acceptance tests
- **Be specific**: Use concrete examples in scenarios, not abstract descriptions
- **Stay focused**: Each story should deliver one clear piece of value
- **Reference existing stories**: Look at stories/001-019.md for format examples

## Anti-Patterns to Avoid

❌ Don't create scenarios that are too abstract ("Then the system should work correctly")
❌ Don't include implementation details in scenarios (no "call API endpoint X")
❌ Don't make stories too large (if it needs >15 scenarios, probably too big)
❌ Don't forget the ```gherkin code block formatting
❌ Don't skip the clarifying questions phase
