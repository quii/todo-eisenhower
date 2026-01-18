package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
)

// AnalyzeInventory analyzes the matrix and returns WIP/inventory metrics
func AnalyzeInventory(m matrix.Matrix) matrix.Inventory {
	// Construct inventory from matrix
	// Use case provides "now" (application concern)
	return matrix.NewInventory(m, time.Now())
}
