package utils

import (
	"fmt"
	"time"
)

func RelativeDate(date time.Time) string {
	now := time.Now()
	diff := now.Sub(date)

	days := int(diff.Hours() / 24)
	weeks := days / 7
	months := int(now.Month()) - int(date.Month()) + 12*(now.Year()-date.Year())

	switch {
	case days == 0:
		return "today"
	case days == 1:
		return "yesterday"
	case days < 7:
		return fmt.Sprintf("%d days ago", days)
	case weeks < 2:
		return "a week ago"
	case weeks < 4:
		return fmt.Sprintf("%d weeks ago", weeks)
	case months < 2:
		return "a month ago"
	case months < 12:
		return fmt.Sprintf("%d months ago", months)
	default:
		years := now.Year() - date.Year()
		if now.YearDay() < date.YearDay() {
			years--
		}
		if years < 2 {
			return "a year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}
