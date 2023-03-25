package controller

import (
	"fmt"
	"time"
)

// isCurrentTimeBetween checks if the current time is between two given times.
// t1 and t2 should be in the format "15:04" (HH:MM).
func isCurrentTimeBetween(t1, t2 string) bool {
	layout := "15:04"
	fmt.Println("current time:", time.Now().Format(layout))

	fmt.Println("T1: ", t1)
	fmt.Println("T2: ", t2)
	fmt.Println("Now:", time.Now().Format(layout))
	now := time.Now().Format(layout)

	if now >= t1 && now <= t2 {
		fmt.Println("inside time frame")
	}
	return now >= t1 && now <= t2
}

var TimeFormat = "2006-01-02 15:04:05"

func InTimeSpan(start string, end string) (bool, error) {
	now := time.Now()
	t1, err := time.ParseInLocation(TimeFormat, start, time.Now().Location())
	if err != nil {
		return false, fmt.Errorf("unable to parse time: %v", err)
	}
	t2, err := time.ParseInLocation(TimeFormat, end, time.Now().Location())
	if err != nil {
		return false, fmt.Errorf("unable to parse time: %v", err)
	}
	return now.After(t1) && now.Before(t2), nil
}
