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
	server = New(path, false, false)
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
		param          string
		got            []byte
		expextedBody   []byte
		expectedStatus int
	}{
		{
			"insert sensordata",
			"/sensordata",
			"POST",
			"",
			[]byte(`{"sensor_id": 1,"dataset_id": 1,"data": ["123.00","12.00"]}`),
			nil,
			200,
		},
		{
			"get sensor list",
			"/sensors",
			"GET",
			"",
			nil,
			nil,
			200,
		},
		{
			"testing a non existent endpoint",
			"/sensors_uggabugga",
			"GET",
			"",
			nil,
			nil,
			404,
		},
		{
			"get sensor data by reference",
			"/sensordata/8a1bbddba98a8d8512787d311352d951",
			"GET",
			"",
			nil,
			nil,
			200,
		},
		{
			"add new sensor endpoint",
			"/sensors",
			"POST",
			"",
			// sensor details
			[]byte(`{"title": "wrom","description": "jalla", "sensor_ip": "10.0.0.123"}`),
			nil,
			200,
		},
		{
			"update sensor",
			"/sensors",
			"PUT",
			"",
			// sensor details
			[]byte(`{"id":1,"title":"kongle","description":"suppe","arduino_key":"8a1bbddba98a8d8512787d311352d951", "sensor_ip": "10.0.0.99"}`),
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
