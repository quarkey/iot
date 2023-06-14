package helper_test

import (
	"log"
	"testing"

	"github.com/quarkey/iot/pkg/helper"
)

var TimeFormat = "2006-01-02 15:04:05"

func TestParseToLocalTime(t *testing.T) {
	t1 := "2021-03-25 15:00:00" // +0100 CET
	w1 := "2021-03-25 15:00:00 +0100 CET"

	t2 := "2021-07-25 18:00:00" // +0200 CEST
	w2 := "2021-07-25 18:00:00 +0200 CEST"

	format := TimeFormat
	t1Out, err := helper.ParseToLocalTime(t1, format)
	if err != nil {
		log.Fatalf("unable to parse time string: %v", err)
	}
	t2Out, err := helper.ParseToLocalTime(t2, format)
	if err != nil {
		log.Fatalf("unable to parse time string: %v", err)
	}
	if t1Out.String() != w1 {
		t.Errorf("expected %v, got %v", w1, t1Out)
	}
	if t2Out.String() != w2 {
		t.Errorf("expected %v, got %v", w2, t2Out)
	}
	// fmt.Println(t1Out)
	// fmt.Println(t2Out)
}
