package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
)

// AnalyzeInventory analyzes the matrix and returns WIP/inventory metrics
func AnalyzeInventory(m matrix.Matrix) matrix.Inventory {
	// Tell the domain to calculate inventory metrics
	// Use case provides "now" (application concern)
	return m.CalculateInventory(time.Now())
}
