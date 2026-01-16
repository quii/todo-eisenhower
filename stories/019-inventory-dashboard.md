# Story 019: Digital Inventory Dashboard

**As a** user managing many todos  
**I want to** see my digital inventory health metrics  
**So that** I can understand my WIP and identify bottlenecks

## Acceptance Criteria

```gherkin
Feature: Digital Inventory Dashboard
  Users should be able to view inventory health metrics to manage WIP

  Scenario: View inventory dashboard with active todos
    Given I have multiple incomplete todos with various contexts
    And some todos have creation dates
    When I press 'i' in overview mode
    Then I see total active work count per quadrant
    And I see a health indicator (OK, HIGH, OVERLOADED)
    And I see age of oldest item per quadrant
    And I see a breakdown by context tags (@architecture, @people, etc)
    And I see completion metrics for last 7 days

  Scenario: Dashboard shows stale items warning
    Given I have a todo that is 21 days old
    When I view the inventory dashboard
    Then I see a warning indicator for that quadrant
    And the item is marked as "VERY STALE"

  Scenario: Dashboard calculates throughput
    Given I have completed 3 todos in the last 7 days
    And I have added 11 todos in the last 7 days
    When I view the inventory dashboard
    Then I see "Completed: 3 items"
    And I see "Added: 11 items"
    And I see a warning "Adding faster than completing"

  Scenario: Dashboard groups by context tags
    Given I have 6 todos with @people context
    And the average age is 18 days
    When I view the inventory dashboard
    Then I see "people: 6 items (avg 18 days old)"
    And I see a health indicator based on count and age

  Scenario: Press 'i' again to exit dashboard
    Given I am viewing the inventory dashboard
    When I press 'i' or ESC
    Then I return to the overview mode
```

## Technical Notes

**Domain Layer:**
- No new domain logic needed - use existing Matrix and Todo data
- Calculate metrics from existing creation/completion dates

**Use Case Layer:**
- `AnalyzeInventory(m matrix.Matrix) InventoryMetrics`
- Returns struct with:
  - Per-quadrant active counts and oldest item age
  - Context breakdown with counts and avg ages
  - 7-day throughput (completed vs added)
  - Health indicators based on thresholds

**UI Layer:**
- Add `inventoryMode` flag to Model
- Press 'i' in overview toggles inventory view
- New `RenderInventoryDashboard()` function
- Colorize health indicators (green/yellow/red)

**Thresholds (research-based):**
- Do First: >5 = HIGH, >8 = OVERLOADED
- Schedule: >3 = HIGH, >6 = OVERLOADED  
- Age: >14 days = STALE, >21 days = VERY STALE
- Context groups: >4 items = HIGH

**Metrics Calculations:**
- Active work: Count todos where IsCompleted() == false
- Age: time.Since(CreationDate) if date exists
- Throughput: Count todos with CompletionDate in last 7 days
- Add rate: Count todos with CreationDate in last 7 days
- Context groups: Extract contexts, group todos, calc avg age
