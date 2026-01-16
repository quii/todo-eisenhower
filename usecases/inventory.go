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
	
	// Context breakdown
	ContextBreakdown map[string]ContextMetrics
}

// ContextMetrics represents metrics for a specific context tag
type ContextMetrics struct {
	Context    string
	Count      int
	AvgAgeDays int
}

// AnalyzeInventory analyzes the matrix and returns WIP/inventory metrics
func AnalyzeInventory(m matrix.Matrix) InventoryMetrics {
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)
	
	metrics := InventoryMetrics{
		ContextBreakdown: make(map[string]ContextMetrics),
	}
	
	allTodos := m.AllTodos()
	
	// Track context ages for averaging
	contextAges := make(map[string][]int)
	
	for _, t := range allTodos {
		// Calculate age if creation date exists
		var ageDays int
		if creationDate := t.CreationDate(); creationDate != nil {
			ageDays = int(now.Sub(*creationDate).Hours() / 24)
		}
		
		// Count active todos and track oldest per quadrant
		if !t.IsCompleted() {
			metrics.TotalActive++
			
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
			
			// Track contexts for active todos
			for _, context := range t.Contexts() {
				if _, exists := contextAges[context]; !exists {
					contextAges[context] = []int{}
				}
				contextAges[context] = append(contextAges[context], ageDays)
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
		if creationDate := t.CreationDate(); creationDate != nil {
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
		
		metrics.ContextBreakdown[context] = ContextMetrics{
			Context:    context,
			Count:      count,
			AvgAgeDays: avgAge,
		}
	}
	
	return metrics
}
