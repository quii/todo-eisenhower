# Stories

This directory contains user stories and requirements written in Gherkin format.

## Structure

Each story file follows the naming convention: `NNN-short-description.md`

Stories include:
- **Goal**: User-focused outcome
- **Background**: Context and domain knowledge
- **Acceptance Criteria**: Gherkin scenarios (Given/When/Then)
- **Technical Notes**: Implementation hints
- **Out of Scope**: What this story explicitly doesn't cover
- **Questions to Answer**: Clarifications needed before implementation

## Working with Stories

1. Review the story and acceptance criteria
2. Ask clarifying questions if anything is unclear
3. Implement using TDD, treating Gherkin scenarios as test cases
4. Verify all scenarios pass before considering the story complete
5. Work in small vertical slices - each story should be shippable

## Story Breakdown Philosophy

We work in very small vertical slices. Each story should be the smallest shippable increment that adds value.

**Example progression:**
1. Story 001: Display matrix with hard-coded in-memory todos (establishes architecture & TUI)
2. Story 002: Load todos from hard-coded file path (adds file I/O)
3. Story 003: Accept custom file path as argument (adds CLI flexibility)

Each story builds on the previous one, maintaining a working, releasable application at every step.

## Story Status

- `001-display-hardcoded-matrix.md` - Ready for implementation
- `002-load-from-hardcoded-path.md` - Ready (depends on 001)
- `003-load-from-custom-path.md` - Ready (depends on 002)
