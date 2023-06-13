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
func InDateTimeSpanString(start string, end string) bool {
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
func InTimeSpan(t1 time.Time, t2 time.Time, now time.Time) bool {
	return now.After(t1) && now.Before(t2)
}

// ParseStrToLocalTime parses two time strings in "2006-01-02 15:04:05" format into CEST +0200.
// Returns an error if the input strings cannot be parsed into time.Time objects.
func ParseStrToLocalTime(t1 string, t2 string) (*time.Time, *time.Time, error) {
	pt1, err := localTimeFixedZone(t1)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse time string: %v", err)
	}
	pt2, err := localTimeFixedZone(t2)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to parse time string: %v", err)
	}
	return pt1, pt2, nil
}

func localTimeFixedZone(t string) (*time.Time, error) {
	zone := time.FixedZone("CEST", 2*60*60)
	pt1, err := time.ParseInLocation(TimeFormat, t, zone)
	if err != nil {
		return nil, err
	}
	return &pt1, nil
}

// InTimeSpanIgnoreDate checks if the current time falls within the time range specified by t1 and t2, ignoring the date.
// This is useful, for example, when working with recurring events that occur at the same time every day.
// func ParseInTimeSpanString(st1, st2 string) bool {
// 	t1, err := time.Parse(TimeFormat, fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), st1))
// 	if err != nil {
// 		fmt.Printf("Error parsing time string: %v\n", err)
// 	}
// 	t2, err := time.Parse(TimeFormat, fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), st2))
// 	if err != nil {
// 		fmt.Printf("Error parsing time string: %v\n", err)
// 	}
// 	// fmt.Println("t1:", t1.Local())
// 	// fmt.Println("t2:", t2.Local())
// 	// fmt.Println("now:", time.Now())
// 	return InTimeSpan(t1.Local(), t2.Local(), time.Now())
// }

func ParseTimeString(ts string) (time.Time, error) {
	t, err := time.Parse(TimeFormat, fmt.Sprintf("%s %s", time.Now().Format("2006-01-02"), ts))
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse time string: %v", err)

	}
	return t, nil
}
