package ui_test

import (
	"testing"
	"time"

	"github.com/quii/todo-eisenhower/adapters/ui"
)

// A date can arrive in a different location to the one the user is in (the
// todo.txt parser used to produce UTC). The count of days between two calendar
// dates must not depend on that.
func TestFormatDate_DateInDifferentLocationToNow(t *testing.T) {
	tests := []struct {
		daysAgo int
		want    string
	}{
		{0, "today"},
		{1, "yesterday"},
		{5, "5 days ago"},
		{10, "10 days ago"},
	}

	// Both sides of UTC: an offset either way used to move the count by a day.
	for _, zone := range []string{"Europe/London", "America/New_York"} {
		loc, err := time.LoadLocation(zone)
		if err != nil {
			t.Fatal(err)
		}

		original := time.Local
		time.Local = loc

		for _, tt := range tests {
			day := time.Now().In(loc).AddDate(0, 0, -tt.daysAgo)
			utcMidnight := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)

			if got := ui.FormatDate(&utcMidnight); got != tt.want {
				t.Errorf("%s: formatDate(%s) = %q, want %q", zone, utcMidnight, got, tt.want)
			}
		}

		time.Local = original
	}
}
