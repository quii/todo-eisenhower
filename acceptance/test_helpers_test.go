package acceptance_test

import (
	"io"
	"regexp"
)

// stripANSI removes ANSI escape codes from a string for easier testing
func stripANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// StubTodoWriter captures write calls for testing
type StubTodoWriter struct {
	saveAllTodosCalled bool
	lastContent        string
	saveError          error
}

func (s *StubTodoWriter) GetAppendWriter() (io.WriteCloser, error) {
	if s.saveError != nil {
		return nil, s.saveError
	}
	return &nopWriteCloser{writer: s}, nil
}

func (s *StubTodoWriter) GetReplaceWriter() (io.WriteCloser, error) {
	s.saveAllTodosCalled = true
	if s.saveError != nil {
		return nil, s.saveError
	}
	return &nopWriteCloser{writer: s}, nil
}

func (s *StubTodoWriter) Write(p []byte) (n int, err error) {
	s.lastContent += string(p)
	return len(p), nil
}

type nopWriteCloser struct {
	writer io.Writer
}

func (n *nopWriteCloser) Write(p []byte) (int, error) {
	return n.writer.Write(p)
}

func (n *nopWriteCloser) Close() error {
	return nil
}
