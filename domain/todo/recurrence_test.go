package todo_test

import (
	"testing"
	"time"

	"github.com/quii/todo-eisenhower/domain/todo"
)

func TestRecurrence_NextDueDate_Relative(t *testing.T) {
	tests := []struct {
		name           string
		interval       int
		unit           todo.RecurrenceUnit
		currentDueDate time.Time
		completionDate time.Time
		today          time.Time
		wantDueDate    time.Time
	}{
		{
			name:           "relative 2 weeks from completion",
			interval:       2,
			unit:           todo.Week,
			currentDueDate: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "relative 1 day from completion",
			interval:       1,
			unit:           todo.Day,
			currentDueDate: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 11, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 3, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "relative 1 month from completion",
			interval:       1,
			unit:           todo.Month,
			currentDueDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "relative 1 year from completion",
			interval:       1,
			unit:           todo.Year,
			currentDueDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2027, 1, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := todo.NewRecurrence(true, tt.interval, tt.unit)
			got := rec.NextDueDate(tt.currentDueDate, tt.completionDate, tt.today)
			if !got.Equal(tt.wantDueDate) {
				t.Errorf("NextDueDate() = %v, want %v", got, tt.wantDueDate)
			}
		})
	}
}

func TestRecurrence_NextDueDate_Strict(t *testing.T) {
	tests := []struct {
		name           string
		interval       int
		unit           todo.RecurrenceUnit
		currentDueDate time.Time
		completionDate time.Time
		today          time.Time
		wantDueDate    time.Time
	}{
		{
			name:           "strict 2 weeks, completed on time",
			interval:       2,
			unit:           todo.Week,
			currentDueDate: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "strict 1 month, completed on time",
			interval:       1,
			unit:           todo.Month,
			currentDueDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "strict skip past stale - due date was weeks ago",
			interval:       1,
			unit:           todo.Week,
			currentDueDate: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "strict 1 day, completed late",
			interval:       1,
			unit:           todo.Day,
			currentDueDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "strict 1 year, skip past stale",
			interval:       1,
			unit:           todo.Year,
			currentDueDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			completionDate: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			today:          time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC),
			wantDueDate:    time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := todo.NewRecurrence(false, tt.interval, tt.unit)
			got := rec.NextDueDate(tt.currentDueDate, tt.completionDate, tt.today)
			if !got.Equal(tt.wantDueDate) {
				t.Errorf("NextDueDate() = %v, want %v", got, tt.wantDueDate)
			}
		})
	}
}

func TestRecurrence_String(t *testing.T) {
	tests := []struct {
		name     string
		relative bool
		interval int
		unit     todo.RecurrenceUnit
		want     string
	}{
		{"relative 2 weeks", true, 2, todo.Week, "+2w"},
		{"strict 1 month", false, 1, todo.Month, "1m"},
		{"relative 3 days", true, 3, todo.Day, "+3d"},
		{"strict 1 year", false, 1, todo.Year, "1y"},
		{"relative 1 day", true, 1, todo.Day, "+1d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := todo.NewRecurrence(tt.relative, tt.interval, tt.unit)
			got := rec.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRecurrence_Getters(t *testing.T) {
	rec := todo.NewRecurrence(true, 2, todo.Week)

	if !rec.IsRelative() {
		t.Error("expected IsRelative() to be true")
	}
	if rec.Interval() != 2 {
		t.Errorf("expected Interval() = 2, got %d", rec.Interval())
	}
	if rec.Unit() != todo.Week {
		t.Errorf("expected Unit() = Week, got %d", rec.Unit())
	}

	rec2 := todo.NewRecurrence(false, 1, todo.Month)
	if rec2.IsRelative() {
		t.Error("expected IsRelative() to be false")
	}
}
