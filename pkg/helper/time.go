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

// InTimeSpanIgnoreDate checks if the current time falls within the time range specified by t1 and t2, ignoring the date.
// This function creates new time.Time objects using the current date and the hour, minute, and second values from t1 and t2.
// If t2 is before t1, it's assumed to be on the following day. Returns true if the current time is between t1 and t2, false otherwise.
func InTimeSpanIgnoreDate(t1 time.Time, t2 time.Time) bool {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), t1.Hour(),
		t1.Minute(), t1.Second(), t1.Nanosecond(), t1.Location())

	endTime := time.Date(now.Year(), now.Month(), now.Day(), t2.Hour(),
		t2.Minute(), t2.Second(), t2.Nanosecond(), t2.Location())

	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}
	return now.After(startTime) && now.Before(endTime)
}
