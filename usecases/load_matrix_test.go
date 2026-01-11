package usecases_test

import (
	"testing"

	"github.com/quii/todo-eisenhower/usecases"
)

func TestLoadMatrix(t *testing.T) {
	t.Run("returns matrix with hard-coded todos for Story 001", func(t *testing.T) {
		m := usecases.LoadMatrix()

		// Verify we have todos in each quadrant
		if len(m.DoFirst()) == 0 {
			t.Error("expected DoFirst quadrant to have todos")
		}
		if len(m.Schedule()) == 0 {
			t.Error("expected Schedule quadrant to have todos")
		}
		if len(m.Delegate()) == 0 {
			t.Error("expected Delegate quadrant to have todos")
		}
		if len(m.Eliminate()) == 0 {
			t.Error("expected Eliminate quadrant to have todos")
		}
	})

	t.Run("returns todos with expected descriptions", func(t *testing.T) {
		m := usecases.LoadMatrix()

		// Verify some expected hard-coded todos exist
		doFirst := m.DoFirst()
		hasExpectedTodo := false
		for _, td := range doFirst {
			if td.Description() == "Fix critical production bug" {
				hasExpectedTodo = true
				break
			}
		}
		if !hasExpectedTodo {
			t.Error("expected to find 'Fix critical production bug' in DoFirst quadrant")
		}
	})
}
