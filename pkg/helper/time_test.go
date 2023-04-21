package helper_test

import (
	"testing"
	"time"

	"github.com/quarkey/iot/pkg/helper"
)

var TimeFormat = "2006-01-02 15:04:05"

// func TestInTimeSpanIgnoreDate(t *testing.T) {
// 	t1, _ := time.Parse(TimeFormat, "12:03:32")
// 	t2, _ := time.Parse(TimeFormat, "01:34:31")
// 	// Test case 1: current time is before t1
// 	currentTime := time.Date(2023, time.March, 23, 12, 0, 0, 0, time.UTC)
// 	if helper.InTimeSpanIgnoreDate(t1, t2) {
// 		t.Errorf("Test case 1 failed: current time %v is not between %v and %v", currentTime, t1, t2)
// 	}
// 	// Test case 2: current time is after t2
// 	currentTime = time.Date(2023, time.August, 26, 2, 0, 0, 0, time.UTC)
// 	if helper.InTimeSpanIgnoreDate(t1, t2) {
// 		t.Errorf("Test case 2 failed: current time %v is not between %v and %v", currentTime, t1, t2)
// 	}
// }

func TestInTimeSpanIgnoreDatex(t *testing.T) {
	t1, _ := time.Parse(TimeFormat, "2023-03-25 18:00:11")
	t2, _ := time.Parse(TimeFormat, "2023-03-26 07:00:11")

	now, _ := time.Parse(TimeFormat, "2023-03-25 23:00:11")
	fail, _ := time.Parse(TimeFormat, "2023-03-25 14:00:11")
	// Test case 1: current time is before t1
	if !helper.InTimeSpanIgnoreDate(t1, t2, now) {
		t.Errorf("Test case 1 failed: current time falls not in timeframe t1=%v, t2=%v", t1, t2)
	}
	if helper.InTimeSpanIgnoreDate(t1, t2, fail) {
		t.Errorf("Test case 2 failed: current time falls not in timeframe t1=%v, t2=%v", t1, t2)
	}
	// // Test case 2: current time is after t2
	// if helper.InTimeSpanIgnoreDate(t1, t2) {
	// 	t.Errorf("Test case 2 failed: current time falls not in timeframe t1=%v, t2=%v", t1, t2)
	// }
}
