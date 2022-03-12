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

var server *Server
var ctx context.Context

func init() {
	path := "../config/exampleconfig.json"
	server = New(path, false)
	server.SetupEndpoints()
	ctx = context.Background()
	go server.Run(ctx)
}
func TestSensordata(t *testing.T) {
	time.Sleep(3 * time.Second)

	tests := []struct {
		descr          string
		url            string
		method         string
		got            []byte
		expextedBody   []byte
		expectedStatus int
	}{
		{
			"insert sensordata",
			"/sensordata",
			"POST",
			[]byte(`{"sensor_id": 1,"dataset_id": 1,"data": [123.00,12.00]}`),
			nil,
			200,
		},
		{
			"get sensor list",
			"/sensors",
			"GET",
			nil,
			nil,
			200,
		},
		{
			"testing a non existent endpoint",
			"/sensors_uggabugga",
			"GET",
			nil,
			nil,
			404,
		},
	}

	for _, test := range tests {
		log.Printf("%s ...", test.descr)
		url := fmt.Sprintf("%s%s", server.API_URL(), test.url)
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
			t.Errorf("unable to insert data to database")
		}
	}

	server.Stop(ctx)
}
