#!/bin/bash
set -e

echo "Running tests..."
go test ./...

echo ""
echo "Running linter..."
~/go/bin/golangci-lint run

echo ""
echo "✓ All checks passed!"
