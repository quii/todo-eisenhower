# Story 017: Compact Overview with Stats

## Status
Completed

## Context
Currently, the overview mode shows all todos in a simple list format, which can be overwhelming and doesn't scale well. The focus mode uses a nice table layout, but using tables for all 4 quadrants in overview would be too cramped.

We need an overview that gives users the "big picture" without overwhelming them with details.

## Goal
Update the overview mode to show summary statistics and a preview of top todos for each quadrant, making it easy to understand workload at a glance.

## Acceptance Criteria

### Scenario 1: Show summary stats for each quadrant
**Given** I have todos distributed across quadrants
**When** I view the overview
**Then** each quadrant should show a summary line like:
```
DO FIRST (5 tasks, 2 completed)
```

### Scenario 2: Show top N todos as simple list
**Given** I have multiple todos in a quadrant
**When** I view the overview
**Then** I should see the first 3-5 todos as a simple bulleted list
**And** I should NOT see table headers or columns

### Scenario 3: Indicate when there are more todos
**Given** I have more than 5 todos in a quadrant
**When** I view the overview
**Then** I should see a line like "... and 3 more (press 1 to view)"

### Scenario 4: Empty quadrant shows helpful message
**Given** I have an empty quadrant
**When** I view the overview
**Then** I should see:
```
DELEGATE (0 tasks)
  (no tasks)
```

### Scenario 5: All completed todos shows in stats
**Given** I have 3 todos all marked as completed in DO FIRST
**When** I view the overview
**Then** I should see "DO FIRST (3 tasks, 3 completed)"

### Scenario 6: Quadrant layout preserved
**Given** I am viewing the overview
**Then** the 4 quadrants should still be arranged in a 2x2 grid
**And** quadrants should have clear visual separation

### Scenario 7: Completed todos shown with visual indicator
**Given** I have completed todos in a quadrant
**When** I view the overview
**Then** completed todos should have a visual indicator (e.g., "✓" or strikethrough)

## Technical Notes
- Update `RenderOverview()` function in `adapters/ui/render.go`
- Show top 5 todos per quadrant (configurable constant)
- Summary format: `QUADRANT_NAME (N tasks, M completed)`
- Simple bullet list: `• Todo description` (no tags, dates, or other columns)
- "... and N more" message when todos exceed display limit
- Keep 2x2 grid layout using lipgloss
- Use subtle styling to differentiate completed todos

## Design Decisions
- **5 todos per quadrant** - Balances information with cleanliness
- **Simple bullets** - No table structure, just descriptions
- **Stats always shown** - Even for empty quadrants, provides context
- **Completed indicator** - Show ✓ prefix for completed todos
- **Press hint** - Remind users they can press 1-4 to see more

## Dependencies
- None - standalone visual update

## Out of Scope
- Showing tags/contexts in overview (focus mode is for details)
- Showing completion/creation dates in overview
- Pagination within overview mode
- Filtering or sorting in overview

## Example Layout
```
┌─ DO FIRST (5 tasks, 2 completed) ──┬─ SCHEDULE (3 tasks, 0 completed) ──┐
│ ✓ Review quarterly goals           │ • Plan sprint                      │
│ ✓ Fix critical bug                 │ • Update documentation             │
│ • Deploy to production             │ • Schedule team meeting            │
│ • Write technical docs             │                                    │
│ • Code review                      │                                    │
├────────────────────────────────────┼────────────────────────────────────┤
│ DELEGATE (1 task, 0 completed)     │ ELIMINATE (8 tasks, 1 completed)   │
│ • Reply to emails                  │ • Old project idea                 │
│                                    │ • Outdated task                    │
│                                    │ • Random note                      │
│                                    │ • Another old idea                 │
│                                    │ • Yet another task                 │
│                                    │   ... and 3 more (press 4 to view) │
└────────────────────────────────────┴────────────────────────────────────┘
```
