package json

import (
	"encoding/json"
	"fmt"
)

// DecodeRawJSON takes a *json.RawMessage and unmarshal content to []string
func DecodeRawJSONtoSlice(raw *json.RawMessage) ([]string, error) {
	var values []string
	err := json.Unmarshal(*raw, &values)
	if err != nil {
		return nil, fmt.Errorf("unable to decode json rawmessage to slice: %v", err)
	}
	return values, nil
}

// JsonRawToString converts json.RawMessage to string, but returning emptu string if marshal fails.
func DecodeRawJSONtoString(raw *json.RawMessage) string {
	j, err := json.Marshal(&raw)
	if err != nil {
		return ""
	}
	return string(j)
}
