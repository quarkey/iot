package helper

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullString wrap sql.NullString to be able to json it pretty
type NullTime struct {
	sql.NullTime
}

// MarshalJSON NullString interface redefinition
func (r NullTime) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Time)
	}
	return json.Marshal("")
}

// UnmarshalJSON NullTime interface redefinition
func (r *NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		r.Valid = true
		r.Time = *x
	} else {
		r.Valid = false
	}
	return nil
}
