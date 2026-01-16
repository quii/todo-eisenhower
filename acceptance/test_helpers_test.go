package acceptance_test

import "regexp"

// stripANSI removes ANSI escape codes from a string for easier testing
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// StubTodoWriter captures write calls for testing
type StubTodoWriter struct {
	replaceAllCalled bool
	lastContent      string
	saveError        error
}

func (s *StubTodoWriter) SaveTodo(line string) error {
	return s.saveError
}

func (s *StubTodoWriter) ReplaceAll(content string) error {
	s.replaceAllCalled = true
	s.lastContent = content
	return s.saveError
}
