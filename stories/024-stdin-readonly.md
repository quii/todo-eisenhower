# Story 024: Stdin Read-Only Mode

As a user
I want to pipe todo.txt content into eisenhower
So that I can visualize combined/filtered todos using standard Unix tools

## Background

Following the Unix philosophy and todo.txt ethos of plain text composability, eisenhower
should accept input from stdin for read-only visualization. This allows users to:
- Combine multiple todo files with `cat`
- Filter todos with `grep`, `awk`, `sed`
- View archives without risking changes
- Compose with any Unix tool

When reading from stdin, the tool enters read-only mode (no edits possible).

## Acceptance Criteria

### Scenario: Auto-detect piped input
Given I have two todo files "work.txt" and "personal.txt"
When I run `cat work.txt personal.txt | eisenhower`
Then eisenhower should start in read-only mode
And display the combined matrix of both files

### Scenario: Read-only mode indicator
Given I am viewing todos from stdin
When the UI renders
Then the header should show "(read-only)" indicator
And it should be clear that editing is disabled

### Scenario: Viewing operations still work
Given I am in read-only mode from stdin
When I press navigation keys (1, 2, 3, 4, ESC)
And I press 'i' for inventory
And I press 'f' for filtering
Then all viewing/navigation features should work normally

### Scenario: Editing operations are disabled
Given I am in read-only mode from stdin
And I am in focus mode on a quadrant
When I try to press editing keys:
  | Key | Operation |
  | a   | Add todo  |
  | e   | Edit todo |
  | d   | Archive   |
  | x   | Delete    |
  | m   | Move      |
  | Space | Toggle complete |
Then nothing should happen
Or a message should indicate "Read-only mode (viewing stdin)"

### Scenario: Quit works normally
Given I am in read-only mode from stdin
When I press 'q'
Then the program should exit cleanly

### Scenario: Normal file mode when file argument provided
Given I have a file "todo.txt"
When I run `eisenhower todo.txt`
Then eisenhower should start in normal edit mode
And editing operations should work

### Scenario: Default file mode when no arguments or pipe
Given I run `eisenhower` with no arguments and no pipe
When the program starts
Then it should use the default file ~/todo.txt
And start in normal edit mode

### Scenario: Combine with grep for filtering
Given I have a todo file with various projects
When I run `grep '+WebApp' todo.txt | eisenhower`
Then I should see only todos with the +WebApp project
And the display should be read-only

### Scenario: View archive file
Given I have completed todos in done.txt
When I run `eisenhower < done.txt`
Then I should see the archived todos in the matrix
And the display should be read-only

### Scenario: Explicit stdin with dash (optional)
Given I run `cat file.txt | eisenhower -`
Then it should behave the same as `cat file.txt | eisenhower`
And enter read-only mode

## Technical Notes

**Stdin Detection:**
```go
stat, _ := os.Stdin.Stat()
if (stat.Mode() & os.ModeCharDevice) == 0 {
    // Stdin is being piped - read from it
    readOnlyMode = true
}
```

**Mode Determination:**
1. If stdin is piped: read-only mode, read from stdin
2. Else if args[1] == "-": read-only mode, read from stdin (explicit)
3. Else if len(args) > 1: normal mode, use args[1] as file path
4. Else: normal mode, use default ~/todo.txt

**UI Changes:**
- Add readOnly bool to Model
- Update header to show "(read-only)" when readOnly is true
- In Update(), ignore edit key presses when readOnly is true
- Optionally show helpful message on edit attempt

**Repository Implications:**
- In read-only mode, use memory repository (no file writes)
- Parse stdin content into memory
- No SaveAll() calls

## Example Usage

```bash
# Combine multiple files
cat work.txt personal.txt | eisenhower

# Filter by project
grep '+WebApp' todo.txt | eisenhower

# View only high priority
grep '(A)' todo.txt | eisenhower

# View archive
eisenhower < done.txt

# Filter and combine
grep -v '^x' work.txt personal.txt | eisenhower

# Complex filtering
awk '/+Project/ && !/(D)/' todo.txt | eisenhower
```

## Future Considerations (NOT in this story)

- Output mode: pipe OUT filtered results (eisenhower as a filter)
- Watch mode for piped input (refresh on changes)
- Color output when piped to another tool
- JSON/structured output for scripting

## Non-Goals

- Editing piped content (read-only only)
- Writing back to multiple files
- Managing file splits/merges
- Concurrent file access
