package helper

import (
	"database/sql"
	"encoding/json"
)

// NullString wrap sql.NullString to be able to json it pretty
type NullString struct {
	sql.NullString
}

// MarshalJSON NullString interface redefinition
func (r NullString) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.String)
	}
	return json.Marshal("")
}

// UnmarshalJSON NullString interface redefinition
func (r *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		r.Valid = true
		r.String = *x
	} else {
		r.Valid = false
	}
	return nil
}
