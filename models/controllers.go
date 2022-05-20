package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
)

type Controllers struct {
	ID          int              `db:"id" json:"id"`
	SensorID    int              `db:"sensor_id" json:"sensor_id"`
	Category    string           `db:"category" json:"category"`
	Title       string           `db:"title" json:"title"`
	Description string           `db:"description" json:"description"`
	Items       *json.RawMessage `db:"items" json:"items"`
	Alert       bool             `db:"alert" json:"alert"`
	Active      bool             `db:"active" json:"active"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
}

// GetControllersList fetches a list of controllers. All types of errors will return empty a slice.
func GetControllersList(db *sqlx.DB) ([]Controllers, error) {
	var cs []Controllers
	err := db.Select(&cs, `
	select
		id,
		sensor_id,
		category,
		title,
		description,
		items,
		alert,
		active,
		created_at 
	from controllers`)
	if err != nil {
		return nil, fmt.Errorf("unable to get list of controllers: %v", err)
	}
	return cs, nil
}
func (s *Server) GetControllersListEndpoint(w http.ResponseWriter, r *http.Request) {
	cs, err := GetControllersList(s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
	}
	helper.Respond(w, r, 200, cs)
}
