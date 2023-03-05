package models

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServerHealth(t *testing.T) {
	path := "../config/exampleconfig.json"
	s := New(path, false, false)
	s.SetupEndpoints()
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s.Run(ctx)
		wg.Done()
	}()
	// log.Print("[INFO] Waiting 3 seconds for services to start")
	time.Sleep(3 * time.Second)
	// log.Print("[TEST] checking api health status ...")

	req, err := http.NewRequest(http.MethodGet, "http://localhost:6001/api/health", nil)
	if err != nil {
		t.Errorf("Failed creating request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("Failed performing request: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("unable to check api health")
	}
	cancel()
	wg.Wait()
}

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
