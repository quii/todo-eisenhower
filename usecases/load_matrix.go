package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
)

// LoadMatrix loads todos from the repository and returns an Eisenhower matrix
func LoadMatrix(repo TodoRepository) (matrix.Matrix, error) {
	todos, err := repo.LoadAll()
	if err != nil {
		return matrix.Matrix{}, err
	}

	return matrix.New(todos), nil
}
