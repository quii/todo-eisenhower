# Story 010: Tag Inventory Display

## User Story

**As a** user viewing the matrix overview
**I want** to see counts of my incomplete todos by tag
**So that** I can understand my current digital inventory and workload distribution

## Acceptance Criteria

### Scenario 1: Display project tag inventory
**Given** I have todos with tags:
- 3 incomplete todos with +strategy
- 2 incomplete todos with +hiring
- 1 incomplete todo with +architecture
- 1 completed todo with +strategy (should not count)
**When** I view the overview matrix
**Then** I see at the bottom:
```
Projects: +strategy (3)  +hiring (2)  +architecture (1)
```
**And** tags are ordered by count (descending)
**And** tags use their consistent hash colors

### Scenario 2: Display context tag inventory
**Given** I have todos with contexts:
- 5 incomplete todos with @computer
- 2 incomplete todos with @phone
- 1 incomplete todo with @office
**When** I view the overview matrix
**Then** I see at the bottom:
```
Projects: (none)
Contexts: @computer (5)  @phone (2)  @office (1)
```

### Scenario 3: Display both project and context inventory
**Given** I have both project and context tags in use
**When** I view the overview matrix
**Then** I see both lines:
```
Projects: +strategy (3)  +hiring (2)  +architecture (1)
Contexts: @computer (5)  @phone (2)  @office (1)
```
**And** they appear at the bottom of the matrix, above the help text

### Scenario 4: No tags in use
**Given** I have no todos or no todos with tags
**When** I view the overview matrix
**Then** I see:
```
Projects: (none)
Contexts: (none)
```

### Scenario 5: Inventory not shown in focus mode
**Given** I press 1/2/3/4 to focus on a quadrant
**When** I view the focused quadrant
**Then** the tag inventory is NOT displayed (focus mode only shows that quadrant's todos)

### Scenario 6: Counts update when adding todos
**Given** I add a new todo with +strategy tag
**When** I return to overview mode
**Then** the +strategy count increments by 1

## Technical Notes

### Implementation Approach

1. **Count incomplete todos by tag:**
   ```go
   func countTagInventory(m matrix.Matrix) (projectCounts, contextCounts map[string]int) {
       projectCounts = make(map[string]int)
       contextCounts = make(map[string]int)

       // Count from all quadrants
       for _, todo := range m.DoFirst() {
           if !todo.IsCompleted() {
               for _, p := range todo.Projects() {
                   projectCounts[p]++
               }
               for _, c := range todo.Contexts() {
                   contextCounts[c]++
               }
           }
       }
       // ... repeat for other quadrants

       return projectCounts, contextCounts
   }
   ```

2. **Sort by count (descending):**
   ```go
   func sortTagsByCount(tagCounts map[string]int) []string {
       type tagCount struct {
           tag   string
           count int
       }

       pairs := make([]tagCount, 0, len(tagCounts))
       for tag, count := range tagCounts {
           pairs = append(pairs, tagCount{tag, count})
       }

       sort.Slice(pairs, func(i, j int) bool {
           if pairs[i].count == pairs[j].count {
               return pairs[i].tag < pairs[j].tag // alphabetical tiebreaker
           }
           return pairs[i].count > pairs[j].count // descending by count
       })

       tags := make([]string, len(pairs))
       for i, p := range pairs {
           tags[i] = p.tag
       }
       return tags
   }
   ```

3. **Render inventory in overview:**
   - Add to `RenderMatrix()` function
   - Place above help text at bottom
   - Use dimmed/italic style to distinguish from main content
   - Format: `+tag (count)` with hash colors

4. **Update RenderMatrix signature:**
   - Pass tag counts or compute them inside render function
   - Probably cleaner to compute inside since we already have the matrix

### Styling

- Use gray/dimmed text for labels "Projects:" and "Contexts:"
- Use hash colors for tags themselves
- Use gray for count numbers in parentheses
- Italic style to show it's metadata/summary info
- Ensure proper spacing from matrix content above and help text below

### Edge Cases

- Empty matrix: show "(none)" for both
- Only projects or only contexts: show "(none)" for the unused type
- Very long tag list: might need wrapping or truncation (future enhancement)
- Tags with same count: sort alphabetically as tiebreaker

## Definition of Done

- [ ] Tag inventory displays at bottom of overview matrix
- [ ] Counts only incomplete todos
- [ ] Tags sorted by count (descending)
- [ ] Both project and context tags shown
- [ ] Tags rendered in consistent hash colors
- [ ] Inventory NOT shown in focus mode
- [ ] Counts update when adding new todos
- [ ] "(none)" shown when no tags of that type
- [ ] Acceptance tests pass
- [ ] Manual testing with real todo.txt file
