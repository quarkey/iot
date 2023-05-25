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
	now1, _ := time.Parse(TimeFormat, "2023-03-25 23:00:11")
	t1, _ := time.Parse(TimeFormat, "2023-03-25 18:00:11")
	t2, _ := time.Parse(TimeFormat, "2023-03-26 07:00:11")

	now2, _ := time.Parse(TimeFormat, "2023-03-26 06:00:09")
	t3, _ := time.Parse(TimeFormat, "2023-03-25 18:00:11")
	t4, _ := time.Parse(TimeFormat, "2023-03-26 06:00:11")

	if !helper.InTimeSpanIgnoreDate(t1, t2, now1) {
		t.Errorf("current time falls not in timeframe \nt1=%v \nt2=%v \nc1=%s", t1, t2, now1)
	}
	if !helper.InTimeSpanIgnoreDate(t3, t4, now2) {
		t.Errorf("current time falls not in timeframe \n t3=%v \n t4=%v \nnow=%s", t3, t4, now2)
	}
}
