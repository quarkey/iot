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
func InTimeSpanIgnoreDate(t1, t2, now time.Time) bool {
	diff := t2.Sub(t1)
	timeStart := t1.Format("15:04:05")
	// fmt.Println("timeStart:", timeStart)

	combined, err := time.Parse(TimeFormat, fmt.Sprintf("%s %s", t1.Format("2006-01-02"), timeStart))
	if err != nil {
		fmt.Printf("Error parsing time string: %v\n", err)
	}
	due := combined.Add(diff)

	// combined, _ := time.Parse(TimeFormat, fmt.Sprintf("%s %s", xdate, xtime))
	// fmt.Println("now:\t\t", now)
	// fmt.Println("combined:\t", combined)
	// fmt.Println("due:\t\t", due)
	// fmt.Println("diff:", diff)

	// fmt.Println("until:", time.Until(due))
	// fmt.Println("inspan:", InTimeSpan(combined, due, now))
	return InTimeSpan(combined, due, now)
}
