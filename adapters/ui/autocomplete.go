// Package ui provides the terminal user interface using Bubble Tea.
package ui

import (
	"strings"
)

// detectTrigger detects if the cursor is after a tag trigger (+ or @)
// Returns the trigger type, partial tag text, and whether a trigger was found
func detectTrigger(inputValue string) (triggerChar, partialTag string, found bool) {
	// Find the last occurrence of + or @
	lastPlus := strings.LastIndex(inputValue, "+")
	lastAt := strings.LastIndex(inputValue, "@")

	triggerPos := -1
	trigger := ""

	if lastPlus > lastAt {
		triggerPos = lastPlus
		trigger = "+"
	} else if lastAt >= 0 {
		triggerPos = lastAt
		trigger = "@"
	}

	if triggerPos == -1 {
		return "", "", false
	}

	// Extract text after the trigger
	afterTrigger := inputValue[triggerPos+1:]

	// If there's a space after the trigger, no autocomplete
	if strings.Contains(afterTrigger, " ") {
		return "", "", false
	}

	return trigger, afterTrigger, true
}

// filterTags filters tags by prefix (case-insensitive)
func filterTags(tags []string, prefix string) []string {
	if prefix == "" {
		return tags
	}

	lowerPrefix := strings.ToLower(prefix)
	var matches []string
	for _, tag := range tags {
		if strings.HasPrefix(strings.ToLower(tag), lowerPrefix) {
			matches = append(matches, tag)
		}
	}
	return matches
}

// completeTag replaces the partial tag with the completed tag in the input
func completeTag(inputValue, completedTag string) string {
	// Find the last trigger position
	lastPlus := strings.LastIndex(inputValue, "+")
	lastAt := strings.LastIndex(inputValue, "@")

	triggerPos := max(lastPlus, lastAt)

	if triggerPos == -1 {
		return inputValue
	}

	// Replace from trigger to end with completed tag + space
	prefix := inputValue[:triggerPos+1]
	return prefix + completedTag + " "
}

// detectDueTrigger detects if the cursor is after "due:"
// Returns the partial shortcut text and whether the trigger was found
func detectDueTrigger(inputValue string) (partialShortcut string, found bool) {
	// Find the last occurrence of "due:" (case-insensitive)
	lowerInput := strings.ToLower(inputValue)
	duePos := strings.LastIndex(lowerInput, "due:")

	if duePos == -1 {
		return "", false
	}

	// Extract text after "due:"
	afterDue := inputValue[duePos+4:]

	// If there's a space after the trigger, no autocomplete
	if strings.Contains(afterDue, " ") {
		return "", false
	}

	return afterDue, true
}

// getDateShortcuts returns the list of available date shortcuts
func getDateShortcuts() []string {
	return []string{
		"today",
		"tomorrow",
		"endofweek",
		"endofmonth",
		"endofquarter",
		"endofnextquarter",
		"endofyear",
		"+1d",
		"+3d",
		"+7d",
		"+1w",
		"+2w",
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
	}
}

// filterDateShortcuts filters date shortcuts by prefix (case-insensitive)
func filterDateShortcuts(prefix string) []string {
	shortcuts := getDateShortcuts()

	if prefix == "" {
		return shortcuts
	}

	lowerPrefix := strings.ToLower(prefix)
	var matches []string
	for _, shortcut := range shortcuts {
		if strings.HasPrefix(strings.ToLower(shortcut), lowerPrefix) {
			matches = append(matches, shortcut)
		}
	}
	return matches
}

// completeDateShortcut replaces the partial shortcut with the completed one
func completeDateShortcut(inputValue, completedShortcut string) string {
	// Find the last "due:" position (case-insensitive)
	lowerInput := strings.ToLower(inputValue)
	duePos := strings.LastIndex(lowerInput, "due:")

	if duePos == -1 {
		return inputValue
	}

	// Replace from "due:" to end with completed shortcut + space
	prefix := inputValue[:duePos+4]
	return prefix + completedShortcut + " "
}
