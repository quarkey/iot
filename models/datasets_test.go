package models

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

// var server *Server
// var ctx context.Context

func init() {
	path := "../config/exampleconfig.json"
	server = New(path, false)
	server.SetupEndpoints()
	ctx = context.Background()
	go server.Run(ctx)
}
func TestDatasets(t *testing.T) {
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
		// datasets
		{
			"GetDatasetsListEndpoint()",
			"/datasets",
			"GET",
			"",
			nil,
			nil,
			200,
		},
		{
			"GetDatasetByReference()",
			"/datasets/",
			"GET",
			"8a1bbddba98a8d8512787d311352d951",
			nil,
			nil,
			200,
		},
		{
			"NewDataset()",
			"/datasets",
			"POST",
			"",
			[]byte(`{
				"sensor_id": 1,
				"title": "test 007",
				"description": "this is a test",
				"reference": "balle",
				"intervalsec": 32,
				"fields": "['kinetic energy']"
			  }`),
			nil,
			200,
		},
	}

	for _, test := range tests {
		log.Printf("%s ...", test.descr)
		url := fmt.Sprintf("%s%s%s", server.API_URL(), test.url, test.param)
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
