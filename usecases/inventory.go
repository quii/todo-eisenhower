package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// InventoryMetrics represents system health and WIP metrics
type InventoryMetrics struct {
	// Per-quadrant metrics
	DoFirstActive     int
	DoFirstOldestDays int
	ScheduleActive    int
	ScheduleOldestDays int
	DelegateActive    int
	DelegateOldestDays int
	EliminateActive   int
	EliminateOldestDays int

	TotalActive       int

	// Throughput metrics
	CompletedLast7Days int
	AddedLast7Days     int

	// Tag breakdowns
	ContextBreakdown map[string]TagMetrics
	ProjectBreakdown map[string]TagMetrics
}

// TagMetrics represents metrics for a specific tag (context or project)
type TagMetrics struct {
	Tag        string
	Count      int
	AvgAgeDays int
}

// AnalyzeInventory analyzes the matrix and returns WIP/inventory metrics
func AnalyzeInventory(m matrix.Matrix) InventoryMetrics {
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	metrics := InventoryMetrics{
		ContextBreakdown: make(map[string]TagMetrics),
		ProjectBreakdown: make(map[string]TagMetrics),
	}

	allTodos := m.AllTodos()

	// Track tag ages for averaging (only for todos with creation dates)
	contextAges := make(map[string][]int)
	projectAges := make(map[string][]int)

	for _, t := range allTodos {
		// Only process age if creation date exists
		creationDate := t.CreationDate()
		var ageDays int
		hasCreationDate := false
		if creationDate != nil {
			ageDays = int(now.Sub(*creationDate).Hours() / 24)
			hasCreationDate = true
		}

		// Count active todos and track oldest per quadrant (only with creation dates)
		if !t.IsCompleted() {
			metrics.TotalActive++

			if hasCreationDate {
				switch t.Priority() {
				case todo.PriorityA:
					metrics.DoFirstActive++
					if ageDays > metrics.DoFirstOldestDays {
						metrics.DoFirstOldestDays = ageDays
					}
				case todo.PriorityB:
					metrics.ScheduleActive++
					if ageDays > metrics.ScheduleOldestDays {
						metrics.ScheduleOldestDays = ageDays
					}
				case todo.PriorityC:
					metrics.DelegateActive++
					if ageDays > metrics.DelegateOldestDays {
						metrics.DelegateOldestDays = ageDays
					}
				case todo.PriorityD, todo.PriorityNone:
					metrics.EliminateActive++
					if ageDays > metrics.EliminateOldestDays {
						metrics.EliminateOldestDays = ageDays
					}
				}
			} else {
				// Still count active todos without dates, just don't track age
				switch t.Priority() {
				case todo.PriorityA:
					metrics.DoFirstActive++
				case todo.PriorityB:
					metrics.ScheduleActive++
				case todo.PriorityC:
					metrics.DelegateActive++
				case todo.PriorityD, todo.PriorityNone:
					metrics.EliminateActive++
				}
			}

			// Track contexts for active todos (only with creation dates and non-empty)
			if hasCreationDate {
				for _, context := range t.Contexts() {
					if context != "" {
						if _, exists := contextAges[context]; !exists {
							contextAges[context] = []int{}
						}
						contextAges[context] = append(contextAges[context], ageDays)
					}
				}

				// Track projects for active todos (only with creation dates and non-empty)
				for _, project := range t.Projects() {
					if project != "" {
						if _, exists := projectAges[project]; !exists {
							projectAges[project] = []int{}
						}
						projectAges[project] = append(projectAges[project], ageDays)
					}
				}
			}
		}

		// Count completed in last 7 days
		if t.IsCompleted() {
			if completionDate := t.CompletionDate(); completionDate != nil {
				if completionDate.After(sevenDaysAgo) {
					metrics.CompletedLast7Days++
				}
			}
		}

		// Count added in last 7 days
		if creationDate != nil {
			if creationDate.After(sevenDaysAgo) {
				metrics.AddedLast7Days++
			}
		}
	}

	// Calculate context breakdown with averages
	for context, ages := range contextAges {
		count := len(ages)
		sum := 0
		for _, age := range ages {
			sum += age
		}
		avgAge := 0
		if count > 0 {
			avgAge = sum / count
		}

		metrics.ContextBreakdown[context] = TagMetrics{
			Tag:        context,
			Count:      count,
			AvgAgeDays: avgAge,
		}
	}

	// Calculate project breakdown with averages
	for project, ages := range projectAges {
		count := len(ages)
		sum := 0
		for _, age := range ages {
			sum += age
		}
		avgAge := 0
		if count > 0 {
			avgAge = sum / count
		}

		metrics.ProjectBreakdown[project] = TagMetrics{
			Tag:        project,
			Count:      count,
			AvgAgeDays: avgAge,
		}
	}

	return metrics
}
