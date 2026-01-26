package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// parseDateShortcut converts date shortcuts to YYYY-MM-DD format
// Supports: today, tomorrow, +3d, +2w, friday, jan25, 2026-01-20, etc.
func parseDateShortcut(shortcut string, now time.Time) (string, error) {
	if shortcut == "" {
		return "", fmt.Errorf("empty shortcut")
	}

	shortcut = strings.ToLower(strings.TrimSpace(shortcut))

	// If already in YYYY-MM-DD format, return as-is
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, shortcut); matched {
		return shortcut, nil
	}

	// Absolute shortcuts
	switch shortcut {
	case "today", "tod":
		return now.Format("2006-01-02"), nil
	case "tomorrow", "tom":
		return now.AddDate(0, 0, 1).Format("2006-01-02"), nil
	case "endofmonth":
		// Last day of current month
		year, month, _ := now.Date()
		lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, now.Location())
		return lastDay.Format("2006-01-02"), nil
	case "endofquarter":
		// Last day of current calendar quarter
		return endOfQuarter(now).Format("2006-01-02"), nil
	case "endofnextquarter":
		// Last day of next calendar quarter
		return endOfNextQuarter(now).Format("2006-01-02"), nil
	case "endofyear":
		// December 31st of current year
		return time.Date(now.Year(), 12, 31, 0, 0, 0, 0, now.Location()).Format("2006-01-02"), nil
	}

	// Relative days: +3d, +7d
	if matched, _ := regexp.MatchString(`^\+\d+d$`, shortcut); matched {
		days, _ := strconv.Atoi(shortcut[1 : len(shortcut)-1])
		return now.AddDate(0, 0, days).Format("2006-01-02"), nil
	}

	// Relative weeks: +2w, +1w
	if matched, _ := regexp.MatchString(`^\+\d+w$`, shortcut); matched {
		weeks, _ := strconv.Atoi(shortcut[1 : len(shortcut)-1])
		return now.AddDate(0, 0, weeks*7).Format("2006-01-02"), nil
	}

	// Weekday names (next occurrence)
	weekdays := map[string]time.Weekday{
		"monday":    time.Monday,
		"mon":       time.Monday,
		"tuesday":   time.Tuesday,
		"tue":       time.Tuesday,
		"wednesday": time.Wednesday,
		"wed":       time.Wednesday,
		"thursday":  time.Thursday,
		"thu":       time.Thursday,
		"friday":    time.Friday,
		"fri":       time.Friday,
		"saturday":  time.Saturday,
		"sat":       time.Saturday,
		"sunday":    time.Sunday,
		"sun":       time.Sunday,
	}

	if targetDay, ok := weekdays[shortcut]; ok {
		return nextWeekday(now, targetDay).Format("2006-01-02"), nil
	}

	// Month abbreviations: jan25, feb14, dec31
	monthPattern := regexp.MustCompile(`^(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)(\d{1,2})$`)
	if matches := monthPattern.FindStringSubmatch(shortcut); matches != nil {
		monthAbbr := matches[1]
		day, _ := strconv.Atoi(matches[2])

		months := map[string]time.Month{
			"jan": time.January, "feb": time.February, "mar": time.March,
			"apr": time.April, "may": time.May, "jun": time.June,
			"jul": time.July, "aug": time.August, "sep": time.September,
			"oct": time.October, "nov": time.November, "dec": time.December,
		}

		month := months[monthAbbr]
		year := now.Year()

		// Use current year
		date := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
		return date.Format("2006-01-02"), nil
	}

	return "", fmt.Errorf("unrecognized date shortcut: %s", shortcut)
}

// nextWeekday returns the next occurrence of the given weekday
// If today is the target weekday, it returns next week's occurrence
func nextWeekday(from time.Time, targetDay time.Weekday) time.Time {
	daysUntil := int(targetDay - from.Weekday())
	if daysUntil <= 0 {
		daysUntil += 7 // Next week
	}
	return from.AddDate(0, 0, daysUntil)
}

// endOfQuarter returns the last day of the current calendar quarter
// Q1: Jan-Mar (ends March 31), Q2: Apr-Jun (ends June 30)
// Q3: Jul-Sep (ends September 30), Q4: Oct-Dec (ends December 31)
func endOfQuarter(from time.Time) time.Time {
	year := from.Year()
	month := from.Month()

	var endMonth time.Month
	switch {
	case month <= 3: // Q1
		endMonth = 3
	case month <= 6: // Q2
		endMonth = 6
	case month <= 9: // Q3
		endMonth = 9
	default: // Q4
		endMonth = 12
	}

	// Last day of the end month
	return time.Date(year, endMonth+1, 0, 0, 0, 0, 0, from.Location())
}

// endOfNextQuarter returns the last day of the next calendar quarter
// Handles year rollover (Q4 → Q1 of next year)
func endOfNextQuarter(from time.Time) time.Time {
	year := from.Year()
	month := from.Month()

	var endMonth time.Month
	var endYear int

	switch {
	case month <= 3: // Q1 → Q2
		endMonth = 6
		endYear = year
	case month <= 6: // Q2 → Q3
		endMonth = 9
		endYear = year
	case month <= 9: // Q3 → Q4
		endMonth = 12
		endYear = year
	default: // Q4 → Q1 of next year
		endMonth = 3
		endYear = year + 1
	}

	// Last day of the end month
	return time.Date(endYear, endMonth+1, 0, 0, 0, 0, 0, from.Location())
}

// expandDateShortcuts finds due:shortcut patterns and expands them to due:YYYY-MM-DD
func expandDateShortcuts(input string, now time.Time) string {
	// Pattern to match due:shortcut (stops at space or end of string)
	pattern := regexp.MustCompile(`due:(\S+)`)

	// Find first match only (in case of multiple due: tags)
	match := pattern.FindStringSubmatch(input)
	if match == nil {
		return input // No due: found
	}

	shortcut := match[1]
	expanded, err := parseDateShortcut(shortcut, now)
	if err != nil {
		// If parsing fails, leave as-is
		return input
	}

	// Replace only the first occurrence
	return strings.Replace(input, "due:"+shortcut, "due:"+expanded, 1)
}
