#!/bin/bash
# Install git hooks for this repository

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Installing git hooks..."

# Install pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash
# Pre-commit hook to run tests and linter

echo "Running pre-commit checks..."
./check.sh

if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Pre-commit checks failed. Commit aborted."
    echo "Fix the issues above and try again."
    exit 1
fi

exit 0
EOF

chmod +x "$HOOKS_DIR/pre-commit"

echo "✓ Pre-commit hook installed successfully!"
echo ""
echo "The hook will run ./check.sh before every commit."
echo "To bypass the hook temporarily, use: git commit --no-verify"
