package json

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestJSONrawToString(t *testing.T) {
	tests := []struct {
		name     string
		args     []byte
		wantType string
	}{
		{
			"DecodeRawJSONtoString",
			json.RawMessage(`["123.00","12.00"]`),
			`string`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DecodeRawJSONtoString((*json.RawMessage)(&tt.args))
			if reflect.TypeOf(got).String() != tt.wantType {
				t.Errorf("DecodeRawJSONtoString() = %v, want %v", got, tt.wantType)
			}
		})
	}
}

func TestDecodeRawJSONtoSlice(t *testing.T) {
	tests := []struct {
		name     string
		args     []byte
		want     string
		wantType string
	}{
		{
			"number array",
			json.RawMessage(`["123.00","12.00"]`),
			`["123.00","12.00"]`,
			`[]string`,
		},
		{
			"complete object",
			json.RawMessage(`'{"name":"John", "age":30, "car":"null"}'`),
			`'{"name":"John", "age":30, "car":"null"}'`,
			`[]string`,
		},
		{
			"complete object",
			json.RawMessage(`'{"name":"John", "age":30, "car":"null"}'`),
			`'{"name":"John", "age":30, "car":"null"}'`,
			`[]string`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := DecodeRawJSONtoSlice((*json.RawMessage)(&tt.args))
			if reflect.TypeOf(got).String() != tt.wantType {
				t.Errorf("DecodeRawJSONtoSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
