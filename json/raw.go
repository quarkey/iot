package json

import (
	"encoding/json"
	"fmt"
)

// DecodeRawJSON takes a *json.RawMessage and unmarshal content to []string
func DecodeRawJSON(raw *json.RawMessage) ([]string, error) {
	var values []string
	err := json.Unmarshal(*raw, &values)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %v", err)
	}
	return values, nil
}
