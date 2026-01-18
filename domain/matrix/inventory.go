package matrix

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
)

// Inventory represents system health and WIP metrics for the matrix
type Inventory struct {
	// Per-quadrant metrics
	DoFirstActive      int
	DoFirstOldestDays  int
	ScheduleActive     int
	ScheduleOldestDays int
	DelegateActive     int
	DelegateOldestDays int
	EliminateActive    int
	EliminateOldestDays int

	TotalActive int

	// Throughput metrics
	CompletedLast7Days int
	AddedLast7Days     int

	// Tag breakdowns
	ContextBreakdown map[string]TagMetrics
	ProjectBreakdown map[string]TagMetrics

	// Construction state (private)
	matrix Matrix
	now    time.Time
}

// TagMetrics represents metrics for a specific tag (context or project)
type TagMetrics struct {
	Tag        string
	Count      int
	AvgAgeDays int
}

// NewInventory constructs an Inventory by analyzing the given matrix.
// The now parameter allows deterministic testing and follows dependency inversion.
func NewInventory(m Matrix, now time.Time) Inventory {
	inv := Inventory{
		ContextBreakdown: make(map[string]TagMetrics),
		ProjectBreakdown: make(map[string]TagMetrics),
		matrix:           m,
		now:              now,
	}

	// Process each quadrant's todos
	inv.processDoFirst()
	inv.processSchedule()
	inv.processDelegate()
	inv.processEliminate()

	// Count throughput
	inv.countThroughput()

	return inv
}

// processDoFirst calculates metrics for the DoFirst quadrant
func (inv *Inventory) processDoFirst() {
	inv.processQuadrant(inv.matrix.doFirst, &inv.DoFirstActive, &inv.DoFirstOldestDays)
}

// processSchedule calculates metrics for the Schedule quadrant
func (inv *Inventory) processSchedule() {
	inv.processQuadrant(inv.matrix.schedule, &inv.ScheduleActive, &inv.ScheduleOldestDays)
}

// processDelegate calculates metrics for the Delegate quadrant
func (inv *Inventory) processDelegate() {
	inv.processQuadrant(inv.matrix.delegate, &inv.DelegateActive, &inv.DelegateOldestDays)
}

// processEliminate calculates metrics for the Eliminate quadrant
func (inv *Inventory) processEliminate() {
	inv.processQuadrant(inv.matrix.eliminate, &inv.EliminateActive, &inv.EliminateOldestDays)
}

// countThroughput counts completed and added todos in the last 7 days
func (inv *Inventory) countThroughput() {
	sevenDaysAgo := inv.now.AddDate(0, 0, -7)

	for _, t := range inv.matrix.AllTodos() {
		if wasCompletedRecently(t, sevenDaysAgo) {
			inv.CompletedLast7Days++
		}
		if wasCreatedRecently(t, sevenDaysAgo) {
			inv.AddedLast7Days++
		}
	}
}

// processQuadrant handles metrics calculation for a single quadrant's todos
func (inv *Inventory) processQuadrant(todos []todo.Todo, activeCount, oldestDays *int) {
	for _, t := range todos {
		if !t.IsCompleted() {
			inv.TotalActive++
			*activeCount++

			if age, hasAge := ageInDays(t.CreationDate(), inv.now); hasAge {
				if age > *oldestDays {
					*oldestDays = age
				}
				inv.trackContexts(t.Contexts(), age)
				inv.trackProjects(t.Projects(), age)
			}
		}
	}
}

// trackContexts records ages for non-empty contexts in the inventory's context breakdown
func (inv *Inventory) trackContexts(contexts []string, age int) {
	for _, context := range contexts {
		if context == "" {
			continue
		}
		inv.trackInBreakdown(inv.ContextBreakdown, context, age)
	}
}

// trackProjects records ages for non-empty projects in the inventory's project breakdown
func (inv *Inventory) trackProjects(projects []string, age int) {
	for _, project := range projects {
		if project == "" {
			continue
		}
		inv.trackInBreakdown(inv.ProjectBreakdown, project, age)
	}
}

// trackInBreakdown updates a tag breakdown with a new age data point using running average
func (inv *Inventory) trackInBreakdown(breakdown map[string]TagMetrics, tag string, age int) {
	metrics := breakdown[tag]
	metrics.Tag = tag

	// Running average: new_avg = (old_avg * old_count + new_value) / new_count
	oldCount := metrics.Count
	metrics.Count++
	metrics.AvgAgeDays = ((metrics.AvgAgeDays * oldCount) + age) / metrics.Count

	breakdown[tag] = metrics
}

// ageInDays calculates the age in days from a creation date, returns (age, hasAge)
func ageInDays(creationDate *time.Time, now time.Time) (int, bool) {
	if creationDate == nil {
		return 0, false
	}
	return int(now.Sub(*creationDate).Hours() / 24), true
}

// wasCompletedRecently checks if a todo was completed after the threshold
func wasCompletedRecently(t todo.Todo, threshold time.Time) bool {
	if !t.IsCompleted() {
		return false
	}
	completionDate := t.CompletionDate()
	return completionDate != nil && completionDate.After(threshold)
}

// wasCreatedRecently checks if a todo was created after the threshold
func wasCreatedRecently(t todo.Todo, threshold time.Time) bool {
	creationDate := t.CreationDate()
	return creationDate != nil && creationDate.After(threshold)
}
