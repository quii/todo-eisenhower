package todotxt_test

import (
	"strings"
	"testing"
	"time"

	"github.com/quii/todo-eisenhower/domain/todotxt"
)

// Dates in todo.txt are calendar dates with no time or zone, so they mean the
// same day in the user's own location - not UTC.
func TestUnmarshal_DatesAreInLocalLocation(t *testing.T) {
	input := "x 2026-07-14 2026-07-10 Task due:2026-07-20\n"

	todos, err := todotxt.Unmarshal(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	want := map[string]time.Time{
		"creation":   time.Date(2026, 7, 10, 0, 0, 0, 0, time.Local),
		"completion": time.Date(2026, 7, 14, 0, 0, 0, 0, time.Local),
		"due":        time.Date(2026, 7, 20, 0, 0, 0, 0, time.Local),
	}
	got := map[string]*time.Time{
		"creation":   todos[0].CreationDate(),
		"completion": todos[0].CompletionDate(),
		"due":        todos[0].DueDate(),
	}

	for name, w := range want {
		g := got[name]
		if g == nil {
			t.Fatalf("%s date not parsed", name)
		}
		if !g.Equal(w) {
			t.Errorf("%s date = %s, want %s", name, g, w)
		}
	}
}
