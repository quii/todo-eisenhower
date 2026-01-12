# Story 004: Parse and Render Projects and Contexts

## Goal
As a user, I want to see project tags (+tag) and context tags (@tag) parsed from my todos and rendered with consistent, distinct colors, so I can visually identify related tasks.

## Background
- Todo.txt format supports **projects** (`+ProjectName`) and **contexts** (`@ContextName`)
- A todo can have multiple projects and multiple contexts
- Tags can appear anywhere in the description
- Examples:
  - `(A) Call Mom +Family @phone`
  - `(B) Write report +Work +Q1Goals @office @computer`
  - `Review code +OpenSource @github`

### Color Consistency
- Same project/context should always have same color
- Use color hashing: generate color from tag name deterministically
- Colors should be bright and bold for visibility

## Acceptance Criteria

### Scenario: Parse single project tag
```gherkin
Given a todo.txt file containing:
  """
  (A) Deploy new feature +WebApp
  """
When I run the application
Then the todo is parsed with project "WebApp"
And the project tag is rendered in a consistent color
```

### Scenario: Parse single context tag
```gherkin
Given a todo.txt file containing:
  """
(B) Call client @phone
  """
When I run the application
Then the todo is parsed with context "phone"
And the context tag is rendered in a consistent color
```

### Scenario: Parse multiple projects and contexts
```gherkin
Given a todo.txt file containing:
  """
  (A) Write quarterly report +Work +Q1Goals @office @computer
  """
When I run the application
Then the todo is parsed with projects ["Work", "Q1Goals"]
And the todo is parsed with contexts ["office", "computer"]
And each tag is rendered with its consistent color
```

### Scenario: Parse tags anywhere in description
```gherkin
Given a todo.txt file containing:
  """
  (A) Review +OpenSource code for @github issues
  """
When I run the application
Then the todo is parsed with project "OpenSource"
And the todo is parsed with context "github"
And tags are highlighted in the rendered output
```

### Scenario: Consistent colors across multiple todos
```gherkin
Given a todo.txt file containing:
  """
  (A) Deploy feature +WebApp @production
  (B) Fix bug +WebApp @testing
  (C) Write docs +WebApp @office
  """
When I run the application
Then all three todos show "+WebApp" in the same color
And each context shows in its own consistent color
```

### Scenario: Todos without tags still render normally
```gherkin
Given a todo.txt file containing:
  """
  (A) Task with no tags
  (B) Another plain task
  """
When I run the application
Then both todos render normally without tag coloring
```

### Scenario: Extract clean description without inline tags for display
```gherkin
Given a todo "Review +OpenSource code for @github issues"
When rendering the todo
Then the description shows: "Review code for issues"
And tags are displayed separately with colors
```

## Technical Notes
- Extend Todo domain to include `Projects() []string` and `Contexts() []string`
- Parser extracts tags using regex: `\+\w+` for projects, `@\w+` for contexts
- UI generates consistent colors: hash tag name → color
- Tags should be rendered bold and bright
- Consider showing tags inline or separately (ask user preference)

## Out of Scope
- Filtering by project/context (future story)
- Editing tags (future story)
- Tag autocomplete (future story)
- Custom tag colors (future story)

## Clarifications
- **Tag Display**: Inline with original text - tags stay in description and are colorized
- **Color Palette**: Bright/bold colors for visibility (#FF6B6B, #4ECDC4, #FFE66D style)
- **Visual Distinction**: Projects (+) rendered bold, Contexts (@) rendered with different style/color range

## Resources
- [Todo.txt Format Specification](https://github.com/todotxt/todo.txt)
- [Todo.txt Syntax Overview](https://swiftodoapp.com/todotxt-syntax/syntax-overview/)
- [Plaintext Productivity Guide](https://plaintext-productivity.net/1-03-how-i-organize-my-todo-txt-file.html)
