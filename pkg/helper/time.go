package helper

import "time"

var TimeFormat = "2006-01-02 15:04:05"

func InTimeSpan(start string, end string) bool {
	now := time.Now()
	t1, err := time.ParseInLocation(TimeFormat, start, time.Now().Location())
	if err != nil {
		return false
	}
	t2, err := time.ParseInLocation(TimeFormat, end, time.Now().Location())
	if err != nil {
		return false
	}
	return now.After(t1) && now.Before(t2)
}
