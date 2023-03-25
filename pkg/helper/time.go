package helper

import (
	"fmt"
	"time"
)

var TimeFormat = "2006-01-02 15:04:05"

// InTimeSpanString checks if the current time falls within the time range specified
// by the start and end times given as strings in "2006-01-02 15:04:05" format.
// Returns true if the current time is between start and end times, false otherwise.
// Returns false if there is a problem parsing the time strings.
func InTimeSpanString(start string, end string) bool {
	now := time.Now()
	t1, t2, err := ParseStrToLocalTime(start, end)
	if err != nil {
		return false
	}
	return now.After(*t1) && now.Before(*t2)
}

// InTimeSpan checks if the current time falls within the time range specified
// by the two time.Time objects t1 and t2. Returns true if the current time is between
// t1 and t2, false otherwise.
func InTimeSpan(t1 time.Time, t2 time.Time) bool {
	now := time.Now()
	return now.After(t1) && now.Before(t2)
}

// ParseStrToLocalTime parses two time strings in "2006-01-02 15:04:05" format into local time,
// and returns their respective time.Time pointers.
// Returns an error if the input strings cannot be parsed into time.Time objects.
func ParseStrToLocalTime(t1 string, t2 string) (*time.Time, *time.Time, error) {
	pt1, err := time.ParseInLocation(TimeFormat, t1, time.Now().Location())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse time string: %v", err)
	}
	pt2, err := time.ParseInLocation(TimeFormat, t2, time.Now().Location())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse time string: %v", err)
	}
	return &pt1, &pt2, nil
}
