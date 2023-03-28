package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Discord struct {
	URL    string `json:"url"`
	Prefix string `json:"prefix"`
}

func ParseConfig(path string) (*Webhooks, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var webhooks Webhooks
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %v", err)
	}
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return nil, fmt.Errorf("unable to parse discord webhook json: %v", err)
	}
	return &webhooks, nil
}

func (d *Discord) Sendf(format string, a ...any) {
	data := map[string]any{
		"Username": "iot-alerter",
		"content":  fmt.Sprintf(d.Prefix+format, a...),
	}

	// Marshal the payload to JSON
	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling payload:", err)
		return
	}

	// Make the HTTP POST request with the payload
	resp, err := http.Post(d.URL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	// Print the response status code
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("[ERROR] discord webhook integration failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		log.Fatalf("[ERROR] discord webhook integration message failed: %v", err)
		fmt.Println("json payload:", (string(payload)))
		fmt.Println("status code:", resp.StatusCode)
		fmt.Println("body from server:", string(dat))
	}
}
