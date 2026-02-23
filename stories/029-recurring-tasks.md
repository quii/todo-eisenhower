# Story 029: Recurring Tasks

As a user
I want to set tasks as recurring with a defined interval
So that when I complete them, a new instance is automatically created with an updated due date

## Background

Some tasks repeat on a regular cadence - weekly 1:1 prep, monthly rent, fortnightly catchups. Currently, users must manually re-create these tasks each time. This story adds a `rec:` tag following the established todo.txt recurring task convention.

**Format:** `rec:[+]<number><unit>`

- **Units:** `d` (days), `w` (weeks), `m` (months), `y` (years)
- **`+` prefix (relative):** Next due date is calculated from the **completion date**. Use when the gap between occurrences matters more than a fixed schedule (e.g., "Prepare for catchup with Dave every 2 weeks from whenever I last did it").
- **No `+` prefix (strict):** Next due date is calculated from the **current due date**. Use for fixed-schedule items that must happen on specific dates regardless of when completed (e.g., "Pay rent on the 1st of every month").

**Examples:**
```
(B) Prepare for catchup with Dave due:2026-03-09 rec:+2w
(C) Pay rent due:2026-03-01 rec:1m
(B) Weekly review +GTD due:2026-02-28 rec:+1w
(A) Quarterly OKR check-in due:2026-03-31 rec:3m
```

## Acceptance Criteria

```gherkin
Feature: Recurring Tasks

  # --- Parsing ---

  Scenario: Parse recurring tag (relative)
    Given I have a todo "(B) Prep for catchup due:2026-03-09 rec:+2w"
    When the todo is parsed
    Then the recurrence should be relative with interval 2 weeks
    And the due date should be 2026-03-09

  Scenario: Parse recurring tag (strict)
    Given I have a todo "(C) Pay rent due:2026-03-01 rec:1m"
    When the todo is parsed
    Then the recurrence should be strict with interval 1 month
    And the due date should be 2026-03-01

  Scenario: Parse recurring tag with days
    Given I have a todo "(A) Daily standup due:2026-03-02 rec:+1d"
    When the todo is parsed
    Then the recurrence should be relative with interval 1 day

  Scenario: Parse recurring tag with years
    Given I have a todo "(B) Annual review due:2026-06-15 rec:1y"
    When the todo is parsed
    Then the recurrence should be strict with interval 1 year

  Scenario: Invalid recurring tag is ignored
    Given I have a todo "(A) Task rec:invalid"
    When the todo is parsed
    Then the todo should have no recurrence
    And "rec:invalid" should remain in the description

  Scenario: Recurring tag without due date is ignored
    Given I have a todo "(B) Task with no due date rec:+1w"
    When the todo is parsed
    Then the todo should have no recurrence
    And "rec:+1w" should remain in the description
    Because recurrence requires a due date as an anchor

  # --- Completion behaviour ---

  Scenario: Completing a relative recurring task creates successor
    Given today is 2026-03-12
    And I have a todo "(B) Prep for catchup due:2026-03-09 rec:+2w"
    When I mark the todo as complete
    Then the original todo should be marked complete with completion date 2026-03-12
    And the original todo should not have a rec: tag
    And a new todo should be created:
      | Field       | Value                                    |
      | Description | Prep for catchup                         |
      | Priority    | B                                        |
      | Due Date    | 2026-03-26                               |
      | Recurrence  | +2w                                      |
      | Created     | 2026-03-12                               |
    Because relative recurrence calculates from completion date (2026-03-12 + 2w)

  Scenario: Completing a strict recurring task creates successor
    Given today is 2026-03-05
    And I have a todo "(C) Pay rent due:2026-03-01 rec:1m"
    When I mark the todo as complete
    Then the original todo should be marked complete with completion date 2026-03-05
    And the original todo should not have a rec: tag
    And a new todo should be created:
      | Field       | Value                                    |
      | Description | Pay rent                                 |
      | Priority    | C                                        |
      | Due Date    | 2026-04-01                               |
      | Recurrence  | 1m                                       |
      | Created     | 2026-03-05                               |
    Because strict recurrence calculates from original due date (2026-03-01 + 1m)

  Scenario: Successor preserves projects and contexts
    Given today is 2026-03-12
    And I have a todo "(B) Weekly review +GTD @computer due:2026-03-09 rec:+1w"
    When I mark the todo as complete
    Then the new todo should have project "+GTD"
    And the new todo should have context "@computer"

  Scenario: Successor appears in same quadrant
    Given today is 2026-03-12
    And I have a todo "(B) Task due:2026-03-09 rec:+1w" in the Schedule quadrant
    When I mark the todo as complete
    Then the new todo should appear in the Schedule quadrant
    And the completed todo should remain in the Schedule quadrant

  Scenario: Uncompleting a recurring task does not undo successor
    Given today is 2026-03-12
    And I have completed a recurring todo that spawned a successor
    When I unmark the completed todo
    Then the successor should remain in the list
    And the original should be marked incomplete
    And the original should not have a rec: tag
    Because the successor is an independent task once created

  Scenario: Strict recurrence skips past dates
    Given today is 2026-04-10
    And I have a todo "(C) Pay rent due:2026-03-01 rec:1m"
    When I mark the todo as complete
    Then the new todo should have due date 2026-05-01
    Because the next due date (2026-04-01) is already in the past,
    so it advances until the due date is in the future

  # --- Display ---

  Scenario: Recurring indicator in overview mode
    Given I have a todo "(B) Prep for catchup due:2026-03-09 rec:+2w"
    When I view the overview matrix
    Then the todo should show a recurrence indicator
    And the description should read something like "Prep for catchup"
    And the indicator should be visually distinct from the due date

  Scenario: Recurring indicator in focused mode
    Given I have a todo "(B) Prep for catchup due:2026-03-09 rec:+2w"
    When I press '2' to focus on Schedule
    Then the due date column should include the recurrence indicator
    And the indicator should show the interval (e.g., "Mar 09 [2w]")

  Scenario: Non-recurring tasks are unaffected
    Given I have a todo "(A) One-off task due:2026-03-15"
    When I view the matrix
    Then no recurrence indicator should be shown
    And the display should look identical to before this feature

  # --- Editing ---

  Scenario: Adding recurrence via add input
    Given I am in focus mode on Schedule quadrant
    When I press 'a' to add a todo
    And I type "Prep for catchup due:2026-03-23 rec:+2w"
    And I press Enter
    Then a new todo should be created with recurrence "+2w"
    And due date should be 2026-03-23

  Scenario: Adding recurrence via edit input
    Given I have a todo "(B) Prep for catchup due:2026-03-09" with no recurrence
    When I press 'e' to edit the todo
    And I append " rec:+2w" to the input
    And I press Enter
    Then the todo should now have recurrence "+2w"

  Scenario: Removing recurrence via edit input
    Given I have a todo "(B) Prep for catchup due:2026-03-09 rec:+2w"
    When I press 'e' to edit the todo
    And I remove "rec:+2w" from the input
    And I press Enter
    Then the todo should no longer have recurrence

  Scenario: rec: tag appears in edit input
    Given I have a todo "(B) Task due:2026-03-09 rec:+2w"
    When I press 'e' to edit the todo
    Then the input should show "Task due:2026-03-09 rec:+2w"
    Because unlike prioritised:, the rec: tag is user-managed

  # --- Persistence ---

  Scenario: Recurring tag round-trips through file
    Given I have a todo file containing "(B) Task due:2026-03-09 rec:+2w"
    When the file is loaded and saved
    Then the file should still contain "rec:+2w"

  Scenario: Completed recurring original is saved without rec: tag
    Given I complete a recurring todo
    When the file is saved
    Then the completed line should not contain "rec:"
    Because completed instances are historical records
```

## Technical Notes

### Domain Changes (`domain/todo/`)

- Add `Recurrence` value object:
  - `Relative bool` — whether `+` prefix was present
  - `Interval int` — numeric amount
  - `Unit string` — one of `d`, `w`, `m`, `y`
- Add `recurrence *Recurrence` field to `Todo`
- Add `Recurrence() *Recurrence` getter
- Add `WithRecurrence(r *Recurrence) Todo` method
- Add `NextDueDate(completionDate time.Time) time.Time` method on `Recurrence`:
  - If relative: `completionDate + interval`
  - If strict: `currentDueDate + interval` (advancing until future)
- Recurrence is only valid when a due date is also present

### Parser Changes (`domain/todotxt/`)

- Parse `rec:` tag similar to `due:` — extract from description, populate field
- Valid format: `rec:[+]<digits><d|w|m|y>` (case-insensitive)
- Invalid format: leave in description as literal text, set no recurrence
- `rec:` without a `due:` on the same todo: leave in description, set no recurrence
- `Marshal()`: include `rec:` tag in output (unlike `prioritised:` which is hidden)
- `FormatForInput()`: include `rec:` tag (it is user-managed, not hidden)
- Strip `rec:` from completed instances when marshalling

### Use Case Changes (`usecases/`)

- Modify `ToggleCompletion` use case:
  - If marking complete AND todo has recurrence:
    1. Create successor todo with `NextDueDate()`, same priority/projects/contexts/recurrence
    2. Strip `rec:` from the now-completed original
    3. Add successor to matrix
    4. Persist both changes
  - If unmarking (uncompleting): no special behaviour — successor already exists independently

### UI Changes (`adapters/ui/`)

- Display recurrence indicator next to due date
- Suggested format: `Mar 09 [2w]` or `Mar 09 [+2w]`
- No special input handling needed — `rec:` is parsed from free-text input like `due:`

### Edge Cases

- **`rec:` without `due:`**: Treated as plain text, not parsed as recurrence
- **Multiple `rec:` tags**: Use first valid one, ignore rest
- **Strict recurrence past due**: Advance due date by interval until it lands in the future
- **Completing then uncompleting**: Successor is independent; uncompleting does not remove it
- **Moving recurring task between quadrants**: Recurrence is preserved, priority changes as normal
- **Archiving**: Completed originals (without `rec:`) archive normally

## Definition of Done

- [ ] `Recurrence` value object in domain with `NextDueDate()` logic
- [ ] Parser handles `rec:` tag (parse, marshal, format for input)
- [ ] `rec:` without `due:` is ignored (left as description text)
- [ ] Invalid `rec:` values are ignored (left as description text)
- [ ] Completing a recurring task creates a successor with correct due date
- [ ] Successor preserves priority, projects, contexts, and recurrence
- [ ] Completed original has `rec:` stripped
- [ ] Strict recurrence advances past stale dates
- [ ] Recurrence indicator shown in overview and focus modes
- [ ] `rec:` tag appears in edit input (user-managed)
- [ ] Round-trip persistence works correctly
- [ ] Uncompleting does not remove successor
- [ ] All acceptance tests pass
- [ ] All existing tests still pass

## Future Considerations (NOT in this story)

- Recurrence without due date (e.g., "every 2 weeks from creation date")
- Named recurrence patterns (e.g., `rec:weekdays`, `rec:MWF`)
- "Skip this occurrence" option
- Recurrence history / streak tracking in inventory dashboard
- Business-day-aware recurrence (skip weekends)
- End date for recurrence (e.g., `rec:+1w until:2026-12-31`)
