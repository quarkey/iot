package models

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServerHealth(t *testing.T) {
	path := "../config/exampleconfig.json"
	s := New(path, false)
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
