package ui

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
)

// extractURLs finds all HTTP/HTTPS URLs in a string
func extractURLs(text string) []string {
	// Pattern matches http:// or https:// followed by non-whitespace characters
	pattern := regexp.MustCompile(`https?://\S+`)
	matches := pattern.FindAllString(text, -1)
	return matches
}

// openURL opens a URL in the system's default browser
func openURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// renderClickableURL wraps a URL with OSC 8 escape sequences for terminal clickability
// Format: \x1b]8;;URL\x1b\\TEXT\x1b]8;;\x1b\\
func renderClickableURL(url string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, url)
}

// ExtractURLsForTest exposes extractURLs for testing
func ExtractURLsForTest(text string) []string {
	return extractURLs(text)
}
