# Mutation Testing Results - Gremlins

**Date:** 2026-01-18 (Updated after test improvements)
**Tool:** Gremlins v0.6.0
**Command:** `gremlins unleash --timeout-coefficient 3 --workers 2 .`

---

## Overall Results

```
Killed: 52 (+10 from baseline)
Lived: 37 (-1 from baseline)
Not covered: 298 (-9 from baseline)
Timed out: 0
Not viable: 0
Skipped: 0

Test efficacy: 58.43% (+5.93% improvement)
Mutator coverage: 23.00% (+2.33% improvement)
```

### What These Numbers Mean

**Test Efficacy: 58.43%** ⚠️ (IMPROVED from 52.50%)
- 58.43% of mutations that were tested were caught by tests
- **37 mutations survived** (down from 38)
- Quality of test assertions has improved

**Mutator Coverage: 23.00%** ⚠️ (IMPROVED from 20.67%)
- 23% of possible mutations are covered by tests
- **298 mutations are not covered** (down from 307)
- Most uncovered mutations are in UI adapters (acceptable)

---

## Improvements Completed

### ✅ Domain Layer Tests Added (domain/parser/parser_test.go)
- **Tag boundary tests**: Tests for todos with only projects, only contexts, and empty tag lists
- **Creation date tests**: Tests for parsing todos with creation dates
- **Date error handling tests**: Tests for malformed dates in various positions
- **Impact**: Parser package now has 100% mutator coverage (up from 95.45%)

### ✅ Formatting Tests Added (usecases/save_todo_test.go)
- **FormatTodo unit tests**: Comprehensive tests for all todo formatting variations
- **Round-trip tests**: Parse → format → parse fidelity tests
- **Edge case tests**: Empty projects/contexts slices
- **Impact**: All save_todo.go mutations now covered and killed

### ✅ Inventory Boundary Tests Added (usecases/inventory_test.go)
- **7-day threshold tests**: Exact boundary tests for throughput metrics (6, 7, 8 days)
- **Age comparison tests**: Tests for same-age todos and oldest tracking
- **Impact**: Caught CONDITIONALS_BOUNDARY mutations in age calculations

### Test Efficacy Breakdown by Layer
- **Domain (parser, todo)**: 50-58% efficacy - acceptable for domain logic
- **Use cases**: ~60% efficacy - good improvement from boundary tests
- **Adapters (UI)**: Mostly not covered - acceptable trade-off

---

## Critical Findings

### Domain Layer (Business Logic) - 11 Lived Mutations

**domain/parser/parser.go**
```
LIVED CONDITIONALS_BOUNDARY at parser.go:80:19   (if len(text) > 0)
LIVED CONDITIONALS_BOUNDARY at parser.go:111:20  (if len(projects) > 0)
LIVED CONDITIONALS_BOUNDARY at parser.go:111:41  (|| len(contexts) > 0)
LIVED CONDITIONALS_BOUNDARY at parser.go:117:19  (if len(projects) > 0)
LIVED CONDITIONALS_BOUNDARY at parser.go:117:40  (|| len(contexts) > 0)
```

**What this means:**
- Changing `>` to `>=` doesn't break tests
- Tests don't verify edge cases: empty vs non-empty collections
- Missing test cases for boundary conditions

---

### Use Cases Layer - 9 Lived Mutations

**usecases/inventory.go**
```
LIVED CONDITIONALS_BOUNDARY at inventory.go:74:17   (age thresholds)
LIVED CONDITIONALS_BOUNDARY at inventory.go:79:17   (age comparisons)
LIVED CONDITIONALS_BOUNDARY at inventory.go:155:12  (metrics thresholds)
LIVED CONDITIONALS_BOUNDARY at inventory.go:174:12  (age calculations)
```

**usecases/load_matrix.go**
```
LIVED CONDITIONALS_NEGATION at load_matrix.go:22:43  (error handling)
LIVED CONDITIONALS_NEGATION at load_matrix.go:22:57  (error conditions)
```

**What this means:**
- Age threshold comparisons aren't tested at boundaries
- Error handling paths may not be properly tested

---

### Adapter Layer - 18 Lived Mutations

Most lived mutations in adapters (UI rendering logic):
- `adapters/ui/render.go` - 6 lived mutations
- `adapters/ui/help.go` - 4 lived mutations
- `adapters/ui/inventory.go` - 2 lived mutations

**Acceptable:** Adapters are harder to test and less critical than domain logic.

---

## Untested Code (307 NOT COVERED mutations)

Major gaps in test coverage:

### usecases/save_todo.go
```
NOT COVERED CONDITIONALS_NEGATION at save_todo.go:32:59
NOT COVERED CONDITIONALS_NEGATION at save_todo.go:36:53
NOT COVERED CONDITIONALS_NEGATION at save_todo.go:42:18
NOT COVERED ARITHMETIC_BASE at save_todo.go:50:48
NOT COVERED ARITHMETIC_BASE at save_todo.go:59:18
NOT COVERED ARITHMETIC_BASE at save_todo.go:64:18
```

**Issue:** `FormatTodo()` and formatting helpers are private and untested directly.

### usecases/toggle_completion.go
```
NOT COVERED CONDITIONALS_BOUNDARY at toggle_completion.go:24:11
NOT COVERED CONDITIONALS_NEGATION at toggle_completion.go:24:24
NOT COVERED CONDITIONALS_NEGATION at toggle_completion.go:37:9
```

**Issue:** Edge cases in toggle completion logic aren't tested.

---

## Recommendations

### Priority 1: Fix Domain Layer Gaps (High Impact)

**domain/parser/parser.go - Add boundary tests:**

```go
func TestParse_EmptyProjects(t *testing.T) {
    // Test: todos with no projects but with contexts
    input := "@context Task"
    todos, err := Parse(strings.NewReader(input))
    // Assert correct handling
}

func TestParse_EmptyContexts(t *testing.T) {
    // Test: todos with projects but no contexts
    input := "+project Task"
    todos, err := Parse(strings.NewReader(input))
    // Assert correct handling
}

func TestParse_BothEmpty(t *testing.T) {
    // Test: todos with neither projects nor contexts
    input := "Task"
    todos, err := Parse(strings.NewReader(input))
    // Assert correct handling
}
```

### Priority 2: Test Private Formatting Functions

**usecases/save_todo.go - Export or test via public API:**

```go
// Option 1: Test via public API (preferred)
func TestAddTodo_FormatsCorrectly(t *testing.T) {
    // Add todo, read back file, verify format
}

// Option 2: Export for testing (if needed)
// Create usecases/export_test.go:
package usecases

var FormatTodoForTesting = formatTodo  // expose for tests
```

### Priority 3: Add Edge Case Tests for Use Cases

**usecases/inventory.go - Test age boundaries:**

```go
func TestAnalyzeInventory_AgeBoundaries(t *testing.T) {
    // Test exact boundary conditions: 14 days, 21 days, etc.
    oldestDate14 := time.Now().AddDate(0, 0, -14)  // exactly 14 days
    oldestDate15 := time.Now().AddDate(0, 0, -15)  // just over
    // ... assert correct categorization
}
```

### Priority 4: Increase Overall Test Coverage

**Current gaps:**
- `toggle_completion.go` - not fully tested
- UI adapter edge cases (lower priority)
- Error handling paths in `load_matrix.go`

---

## Next Steps

1. **Add boundary tests for parser** (domain/parser/parser_test.go)
   - Empty collections edge cases
   - Date parsing boundaries

2. **Test FormatTodo indirectly** via acceptance tests
   - Add round-trip tests (parse → format → parse)

3. **Test inventory age calculations** at exact boundaries
   - 14, 21 day thresholds

4. **Re-run Gremlins** to measure improvement
   ```bash
   ~/go/bin/gremlins unleash --timeout-coefficient 3 .
   ```

5. **Set quality gates** in CI
   ```yaml
   - name: Mutation Testing
     run: |
       gremlins unleash --threshold-efficacy 70 --threshold-mcover 50
   ```

---

## Mutation Types Found

- **CONDITIONALS_BOUNDARY**: Changed `>` to `>=`, `<` to `<=`
- **CONDITIONALS_NEGATION**: Changed `==` to `!=`, `>` to `<=`
- **ARITHMETIC_BASE**: Changed `+` to `-`, `*` to `/`
- **INCREMENT_DECREMENT**: Changed `i++` to `i--`
- **INVERT_NEGATIVES**: Changed `-x` to `+x`

---

## References

- Full output: `gremlins-output.txt`
- JSON report: `gremlins-report.json`
- Gremlins docs: https://gremlins.dev/
