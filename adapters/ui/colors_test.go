package ui_test

import (
	"testing"

	"github.com/quii/todo-eisenhower/adapters/ui"
)

func TestContrastingForeground(t *testing.T) {
	tests := []struct {
		name       string
		background string
		want       string
	}{
		{name: "white background gets black text", background: "#FFFFFF", want: "#000000"},
		{name: "black background gets white text", background: "#000000", want: "#FFFFFF"},
		{name: "light pastel background gets black text", background: "#C7CEEA", want: "#000000"},
		{name: "saturated indigo background gets white text", background: "#6366F1", want: "#FFFFFF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(ui.ContrastingForeground(tt.background))
			if got != tt.want {
				t.Errorf("ContrastingForeground(%s) = %s, want %s", tt.background, got, tt.want)
			}
		})
	}
}
