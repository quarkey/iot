package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DecodeBody is decoding directly from the request reader
func DecodeBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}
func encodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

// Respond ....
func Respond(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) //setting statuscode
	if data != nil {
		encodeBody(w, r, data)
	}
}

// RespondErr ....
func RespondErr(w http.ResponseWriter, r *http.Request, status int, args ...interface{}) {
	Respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}
func RespondErrf(w http.ResponseWriter, r *http.Request, status int, format string, a ...any) {
	Respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": fmt.Sprintf(format, a...),
		},
	})
}

// RespondSuccess ....
func RespondSuccess(w http.ResponseWriter, r *http.Request, args ...interface{}) {
	Respond(w, r, 200, map[string]interface{}{
		"success": map[string]interface{}{
			"message": fmt.Sprint(args...),
		},
	})
}

// RespondHTTPErr takes a http status code and write it to client as text
func RespondHTTPErr(w http.ResponseWriter, r *http.Request, status int) {
	Respond(w, r, status, map[string]interface{}{
		"error": map[string]interface{}{
			"message": http.StatusText(status),
		},
	})
}
