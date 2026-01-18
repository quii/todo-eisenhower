package architecture_test

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDomainLayerPurity ensures domain packages follow hexagonal architecture rules:
// - No imports of adapters (hexagonal architecture violation)
// - No imports of usecases (dependency inversion violation)
// - No direct infrastructure access (os, net/http, database/sql)
func TestDomainLayerPurity(t *testing.T) {
	forbiddenImports := []string{
		"github.com/quii/todo-eisenhower/adapters",
		"github.com/quii/todo-eisenhower/usecases",
		"os",           // filesystem access - use ports instead
		"net/http",     // network access - use ports instead
		"database/sql", // database access - use ports instead
	}

	violations := findImportViolations("domain", forbiddenImports)

	if len(violations) > 0 {
		t.Errorf("Domain layer architecture violations found:\n\n%s\n\n"+
			"Domain packages must not import infrastructure or outer layers.\n"+
			"Define a port (interface) instead and let adapters implement it.",
			strings.Join(violations, "\n"))
	}
}

// TestUseCaseLayerPurity ensures use cases follow hexagonal architecture rules:
// - Can import domain packages (allowed)
// - Cannot import adapters (must define ports/interfaces instead)
func TestUseCaseLayerPurity(t *testing.T) {
	forbiddenImports := []string{
		"github.com/quii/todo-eisenhower/adapters",
	}

	violations := findImportViolations("usecases", forbiddenImports)

	if len(violations) > 0 {
		t.Errorf("Use case layer architecture violations found:\n\n%s\n\n"+
			"Use cases must not import adapters.\n"+
			"Define ports (interfaces) in use cases and let adapters implement them.",
			strings.Join(violations, "\n"))
	}
}

// findImportViolations scans all .go files in the given directory (recursively)
// and returns a list of violations where forbidden imports are used
func findImportViolations(layerPath string, forbiddenImports []string) []string {
	var violations []string

	err := filepath.Walk(layerPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only check .go files, skip test files
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file to extract imports
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Check each import against forbidden list
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)

			for _, forbidden := range forbiddenImports {
				if importPath == forbidden || strings.HasPrefix(importPath, forbidden+"/") {
					violations = append(violations,
						fmt.Sprintf("  ❌ %s imports forbidden package '%s'", path, importPath))
				}
			}
		}

		return nil
	})

	if err != nil {
		violations = append(violations, fmt.Sprintf("  ⚠️  Error scanning %s: %v", layerPath, err))
	}

	return violations
}
