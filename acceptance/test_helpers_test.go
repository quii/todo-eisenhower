package acceptance_test

import (
	"regexp"
)

// stripANSI removes ANSI escape codes from a string for easier testing
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}
