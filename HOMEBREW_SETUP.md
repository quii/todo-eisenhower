# Homebrew Setup Guide

This project uses GoReleaser to automatically publish to Homebrew.

## Setup Steps

### 1. Create Homebrew Tap Repository

Create a new repository on GitHub named `homebrew-tap`:

```bash
gh repo create quii/homebrew-tap --public --description "Homebrew tap for quii's projects"
```

Or create it manually at: https://github.com/new
- Name: `homebrew-tap`
- Description: "Homebrew tap for quii's projects"
- Public repository

### 2. Grant GitHub Token Permissions

The `GITHUB_TOKEN` in GitHub Actions needs permission to push to your tap repository.

**Option A: Default GITHUB_TOKEN (Recommended)**

The workflow uses the default `GITHUB_TOKEN`. GoReleaser will automatically push to your tap if it's in the same GitHub account.

**Option B: Personal Access Token (if needed)**

If the default token doesn't work, create a Personal Access Token:

1. Go to: https://github.com/settings/tokens/new
2. Name: "GoReleaser Homebrew Tap"
3. Scopes needed:
   - `repo` (full control of private repositories)
   - `write:packages` (if you want to publish packages)
4. Add as repository secret named `GORELEASER_TOKEN`
5. Update `.github/workflows/release.yml` to use `${{ secrets.GORELEASER_TOKEN }}`

### 3. Create a Release

Tag and push a version:

```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

GoReleaser will:
- Run tests
- Build binaries for all platforms
- Create GitHub release with binaries
- Create/update Homebrew formula in `quii/homebrew-tap`

### 4. Users Install Your App

```bash
# Add the tap
brew tap quii/tap

# Install todo-eisenhower
brew install todo-eisenhower
```

Or in one command:
```bash
brew install quii/tap/todo-eisenhower
```

## Testing Locally

Test GoReleaser locally without publishing:

```bash
# Install GoReleaser
brew install goreleaser

# Test the build (creates dist/ directory but doesn't publish)
goreleaser release --snapshot --clean --skip=publish

# Check what would be released
goreleaser release --skip=publish --skip=validate
```

## Updating the Formula

When you push a new tag, GoReleaser automatically:
1. Updates the formula with new version
2. Updates SHA256 checksums
3. Commits and pushes to your tap

Users update with:
```bash
brew update
brew upgrade todo-eisenhower
```

## Troubleshooting

### "Permission denied" when pushing to tap
- Check that `homebrew-tap` repository exists and is public
- Verify GitHub token has `repo` scope
- Ensure token is not expired

### Formula not found
- Verify the tap repository name is exactly `homebrew-tap`
- Check that GoReleaser successfully pushed to the tap (check Actions logs)
- Try: `brew tap --repair`

### Build fails
- Run tests locally: `go test ./...`
- Check GoReleaser config: `goreleaser check`
- Review Actions logs for specific errors
