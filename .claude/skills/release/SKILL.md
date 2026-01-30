---
allowed-tools:
  - Bash(git status*)
  - Bash(git fetch*)
  - Bash(git tag*)
  - Bash(git rev-list*)
  - Bash(git push origin main --tags)
  - Bash(./check.sh*)
---

# Release

A skill for creating and pushing releases with semantic versioning.

## Instructions

You are helping the user create a new release. Follow this workflow:

### 1. Check Working Directory is Clean

Run `git status` to verify there are no uncommitted changes.

If there are uncommitted changes:
- List the uncommitted files
- Ask the user: "There are uncommitted changes. Would you like to commit them first, or abort the release?"
- Do NOT proceed with the release until the working directory is clean

### 2. Run Release Checks

Run `./check.sh` to verify the code is releasable.

If checks fail:
- Show the user the failure output
- Do NOT proceed with the release until the checks pass

### 3. Ensure Local is Up to Date

Run `git fetch` to get latest remote state, then check if local branch is behind remote.

If behind:
- Inform the user and ask if they want to pull first
- Do NOT proceed until resolved

### 4. Determine Next Version

Get the latest tag using:
```bash
git tag --sort=-v:refname | head -1
```

Parse the current version (format: `vMAJOR.MINOR.PATCH`) and calculate the next version:

**Default (minor bump):**
- v0.14.1 → v0.15.0
- v0.15.0 → v0.16.0

**If user specifies "patch":**
- v0.14.1 → v0.14.2
- v0.15.0 → v0.15.1

### 5. Confirm with User

Show the user:
- Current version
- Next version
- Number of commits since last tag (use `git rev-list <last-tag>..HEAD --count`)

Ask for confirmation before proceeding.

### 6. Create and Push Release

Execute these commands:
```bash
git tag <new-version>
git push origin main --tags
```

### 7. Confirm Success

Tell the user:
- The new tag that was created
- That the push was successful
- Remind them that the GitHub release will be created automatically by the pipeline

## Arguments

The skill accepts an optional argument to specify bump type:

- `/release` - Creates a minor version bump (default)
- `/release patch` - Creates a patch version bump

## Example Interaction

**User:** `/release`

**Agent:** Checking working directory...
✓ Working directory is clean

Running release checks...
✓ All checks passed

Fetching latest from remote...
✓ Local branch is up to date

Current version: v0.14.1
Next version: v0.15.0
Commits since last release: 3

Proceed with release? (y/n)

**User:** y

**Agent:**
✓ Created tag v0.15.0
✓ Pushed to origin/main with tags

Release v0.15.0 is complete! The GitHub release will be created automatically by the pipeline.

## Error Handling

- If `./check.sh` fails, show the output and abort
- If `git push` fails, inform the user and suggest checking their permissions or network connection
- If there are no existing tags, start from v0.1.0
- If the tag already exists, inform the user and abort
