package models

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestAreaplotDataSeries(t *testing.T) {
	time.Sleep(3 * time.Second)

	tests := []struct {
		descr          string
		url            string
		method         string
		param          string
		got            []byte
		expextedBody   []byte
		expectedStatus int
	}{
		{
			"get sensor data by reference for plotting",
			"/chart/area/8a1bbddba98a8d8512787d311352d951",
			"GET",
			"",
			nil,
			nil,
			200,
		},
		{
			"failing: get sensor data by reference for plotting",
			"/chart/area/donotmatch",
			"GET",
			"",
			nil,
			nil,
			400,
		},
	}

	for _, test := range tests {
		log.Printf("%s ...", test.descr)
		url := fmt.Sprintf("%s%s%s", server.API_URL(), test.url, test.param)
		fmt.Println("testing", url)
		req, err := http.NewRequest(test.method, url, bytes.NewBuffer(test.got))
		if err != nil {
			t.Errorf("Failed creating request: %v", err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("Failed performing request: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != test.expectedStatus {
			t.Errorf("test failed: %s - expected: %d, got: %d", test.descr, test.expectedStatus, res.StatusCode)
		}
	}
	server.Stop(ctx)
}
