package ui

import (
	"strings"
)

// detectTrigger detects if the cursor is after a tag trigger (+ or @)
// Returns the trigger type, partial tag text, and whether a trigger was found
func detectTrigger(inputValue string) (triggerChar string, partialTag string, found bool) {
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
func completeTag(inputValue string, completedTag string) string {
	// Find the last trigger position
	lastPlus := strings.LastIndex(inputValue, "+")
	lastAt := strings.LastIndex(inputValue, "@")

	triggerPos := lastPlus
	if lastAt > lastPlus {
		triggerPos = lastAt
	}

	if triggerPos == -1 {
		return inputValue
	}

	// Replace from trigger to end with completed tag + space
	prefix := inputValue[:triggerPos+1]
	return prefix + completedTag + " "
}
