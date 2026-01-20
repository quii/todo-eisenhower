package ui

import "time"

// formatDueDate formats a due date for display.
// Returns "Jan 25" for current year, "2025-12-25" for other years, "" for nil.
func formatDueDate(dueDate *time.Time, now time.Time) string {
	if dueDate == nil {
		return ""
	}

	// If same year, show month abbreviation and day
	if dueDate.Year() == now.Year() {
		return dueDate.Format("Jan 02")
	}

	// Different year, show full date
	return dueDate.Format("2006-01-02")
}

// isDueOverdue checks if a due date is in the past (before today).
// Today is NOT considered overdue.
// Dates are normalized to start of day for comparison.
func isDueOverdue(dueDate *time.Time, now time.Time) bool {
	if dueDate == nil {
		return false
	}

	// Normalize both dates to start of day for comparison
	dueDay := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, dueDate.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Overdue if due date is before today
	return dueDay.Before(today)
}

// formatDueDateWithOverdue formats a due date and returns whether it's overdue.
// This is the main function used by rendering code.
// Returns: (formattedDate string, isOverdue bool)
func formatDueDateWithOverdue(dueDate *time.Time, now time.Time) (string, bool) {
	formatted := formatDueDate(dueDate, now)
	overdue := isDueOverdue(dueDate, now)
	return formatted, overdue
}
