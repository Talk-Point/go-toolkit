package formatter

import (
	"testing"
	"time"
)

func TestTimePeriodHumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int32
		expected string
	}{
		{
			name:     "0 seconds",
			seconds:  0,
			expected: "0s",
		},
		{
			name:     "59 seconds",
			seconds:  59,
			expected: "59s",
		},
		{
			name:     "1 minute 1 second",
			seconds:  61,
			expected: "1m 1s",
		},
		{
			name:     "1 hour 1 minute 1 second",
			seconds:  3661,
			expected: "1h 1m 1s",
		},
		{
			name:     "1 day",
			seconds:  86400,
			expected: "1d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := TimePeriodHumanReadable(tt.seconds)
			if actual != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestTimeAbsoluteFormatter(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{"now", now, "now"},
		{"one second ago", now.Add(-1 * time.Second), "1 seconds ago"},
		{"one second from now", now.Add(1 * time.Second), "1 seconds from now"},
		{"one minute ago", now.Add(-1 * time.Minute), "1 minutes ago"},
		{"one minute from now", now.Add(1 * time.Minute), "1 minutes from now"},
		{"one hour ago", now.Add(-1 * time.Hour), "1 hours ago"},
		{"one hour from now", now.Add(1 * time.Hour), "1 hours from now"},
		{"one day ago", now.AddDate(0, 0, -1), "1 days ago"},
		{"one day from now", now.AddDate(0, 0, 1), "1 days from now"},
		{"one week ago", now.AddDate(0, 0, -7), "1 weeks ago"},
		{"one week from now", now.AddDate(0, 0, 7), "1 weeks from now"},
		{"one month ago", now.AddDate(0, -1, 0), "1 months ago"},
		{"one month from now", now.AddDate(0, 1, 0), "1 months from now"},
		{"one year ago", now.AddDate(-1, 0, 0), "1 years ago"},
		{"one year from now", now.AddDate(1, 0, 0), "1 years from now"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeAbsoluteFormatter(tt.date, now); got != tt.expected {
				t.Errorf("TimeAbsoluteFormatter() = %v, want %v", got, tt.expected)
			}
		})
	}
}
