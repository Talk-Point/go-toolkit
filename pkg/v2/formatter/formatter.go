package formatter

import (
	"fmt"
	"time"
)

// TimePeriodHumanReadable converts a time period in seconds to a human readable format.
// The function takes an int32 representing the number of seconds and returns a string.
// If the number of seconds is 0, it returns "0s".
// If the number of seconds is less than 60, it returns the number of seconds followed by "s".
// If the number of seconds is less than 3600 (1 hour), it returns the number of minutes and remaining seconds in the format "Xm Ys".
// If the number of seconds is less than 86400 (1 day), it returns the number of hours, minutes, and remaining seconds in the format "Xh Ym Zs".
// If the number of seconds is 86400 or more, it returns the number of days in the format "Xd".
func TimePeriodHumanReadable(seconds int32) string {
	if seconds == 0 {
		return "0s"
	} else if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	} else if seconds < 3600 {
		minutes := seconds / 60
		remainingSeconds := seconds % 60
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	} else if seconds < 86400 {
		hours := seconds / 3600
		remainingMinutes := (seconds % 3600) / 60
		remainingSeconds := (seconds % 3600) % 60
		return fmt.Sprintf("%dh %dm %ds", hours, remainingMinutes, remainingSeconds)
	} else {
		return fmt.Sprintf("%dd", seconds/86400)
	}
}

// TimeAbsoluteFormatter converts a time.Time to a human readable format relative to a reference time.Time.
// The function takes two time.Time arguments, date and referenceDate, and returns a string.
// If the date is before the referenceDate, it returns the date in the format "X days ago".
// If the date is after the referenceDate, it returns the date in the format "X days from now".
// If the date is the same as the referenceDate, it returns "today".
func TimeAbsoluteFormatter(date time.Time, referenceDate time.Time) string {
	duration := referenceDate.Sub(date)
	switch {
	case duration < 0:
		duration = -duration
		switch {
		case duration < time.Minute:
			return fmt.Sprintf("%d seconds from now", int(duration.Seconds()))
		case duration < time.Hour:
			return fmt.Sprintf("%d minutes from now", int(duration.Minutes()))
		case duration < 24*time.Hour:
			return fmt.Sprintf("%d hours from now", int(duration.Hours()))
		case duration < 7*24*time.Hour:
			return fmt.Sprintf("%d days from now", int(duration.Hours()/24))
		case duration < 30*24*time.Hour:
			return fmt.Sprintf("%d weeks from now", int(duration.Hours()/24/7))
		case duration < 12*30*24*time.Hour:
			return fmt.Sprintf("%d months from now", int(duration.Hours()/24/30))
		default:
			return fmt.Sprintf("%d years from now", int(duration.Hours()/24/365))
		}
	case duration > 0:
		switch {
		case duration < time.Minute:
			return fmt.Sprintf("%d seconds ago", int(duration.Seconds()))
		case duration < time.Hour:
			return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
		case duration < 24*time.Hour:
			return fmt.Sprintf("%d hours ago", int(duration.Hours()))
		case duration < 7*24*time.Hour:
			return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
		case duration < 30*24*time.Hour:
			return fmt.Sprintf("%d weeks ago", int(duration.Hours()/24/7))
		case duration < 12*30*24*time.Hour:
			return fmt.Sprintf("%d months ago", int(duration.Hours()/24/30))
		default:
			return fmt.Sprintf("%d years ago", int(duration.Hours()/24/365))
		}
	default:
		return "now"
	}
}
