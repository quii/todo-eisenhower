package todo

import (
	"fmt"
	"time"
)

// RecurrenceUnit represents the time unit for recurrence intervals
type RecurrenceUnit int

const (
	Day RecurrenceUnit = iota
	Week
	Month
	Year
)

// Recurrence represents a recurring task schedule.
// It can be relative (from completion date) or strict (from original due date).
type Recurrence struct {
	relative bool
	interval int
	unit     RecurrenceUnit
}

// NewRecurrence creates a new Recurrence value object
func NewRecurrence(relative bool, interval int, unit RecurrenceUnit) Recurrence {
	return Recurrence{
		relative: relative,
		interval: interval,
		unit:     unit,
	}
}

// IsRelative returns true if the recurrence is relative to completion date
func (r Recurrence) IsRelative() bool {
	return r.relative
}

// Interval returns the recurrence interval
func (r Recurrence) Interval() int {
	return r.interval
}

// Unit returns the recurrence time unit
func (r Recurrence) Unit() RecurrenceUnit {
	return r.unit
}

// NextDueDate calculates the next due date for a recurring task.
// For relative recurrence: completionDate + interval
// For strict recurrence: currentDueDate + interval, advancing until result > today
func (r Recurrence) NextDueDate(currentDueDate, completionDate, today time.Time) time.Time {
	if r.relative {
		return addInterval(completionDate, r.interval, r.unit)
	}

	// Strict: advance from current due date until past today
	next := addInterval(currentDueDate, r.interval, r.unit)
	for !next.After(today) {
		next = addInterval(next, r.interval, r.unit)
	}
	return next
}

func addInterval(t time.Time, interval int, unit RecurrenceUnit) time.Time {
	switch unit {
	case Day:
		return t.AddDate(0, 0, interval)
	case Week:
		return t.AddDate(0, 0, interval*7)
	case Month:
		return t.AddDate(0, interval, 0)
	case Year:
		return t.AddDate(interval, 0, 0)
	default:
		return t
	}
}

// String returns the recurrence in todo.txt rec: tag format (e.g., "+2w" or "1m")
func (r Recurrence) String() string {
	prefix := ""
	if r.relative {
		prefix = "+"
	}

	var unitStr string
	switch r.unit {
	case Day:
		unitStr = "d"
	case Week:
		unitStr = "w"
	case Month:
		unitStr = "m"
	case Year:
		unitStr = "y"
	}

	return fmt.Sprintf("%s%d%s", prefix, r.interval, unitStr)
}
