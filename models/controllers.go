package models

import (
	"encoding/json"
	"fmt"
	"log"
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
	Switch      int              `db:"switch" json:"switch"`
	Items       *json.RawMessage `db:"items" json:"items"`
	Alert       bool             `db:"alert" json:"alert"`
	Active      bool             `db:"active" json:"active"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
}
type Thresholdswitch struct {
	Description    string  `json:"description"`
	Datasource     string  `json:"datasource"`
	Operation      string  `json:"operation"`
	ThresholdLimit float64 `json:"threshold_limit"`
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
		switch,
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

func (c Controllers) Check(dataPoint float64, db *sqlx.DB) {
	var ts []Thresholdswitch
	err := json.Unmarshal(*c.Items, &ts)
	if err != nil {
		log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
	}

	switch c.Category {
	case "switch":
	case "thresholdswitch":
		for _, item := range ts {
			// threshold switch operation
			switch item.Operation {
			// TODO: update controller active if threshold is met
			case "greather than":
				// switching state based on threshold
				if dataPoint > item.ThresholdLimit {
					err := c.UpdateControllerSwitchState(db, 1)
					if err != nil {
						fmt.Println(err)
					}
					c.Switch = 1
					fmt.Printf("gt switch on: %s -> %v - state: %d\n", c.Title, dataPoint, c.Switch)
					return
				}
				err := c.UpdateControllerSwitchState(db, 0)
				if err != nil {
					fmt.Println(err)
				}
				c.Switch = 0
				fmt.Printf("gt switch off %v - state: %d\n", dataPoint, c.Switch)
			case "less than":
				// switching state based on threshold
				if dataPoint < item.ThresholdLimit {
					err := c.UpdateControllerSwitchState(db, 1)
					if err != nil {
						fmt.Println(err)
					}
					c.Switch = 1
					fmt.Printf("lt switch on: %s -> %v - state: %d\n", c.Title, dataPoint, c.Switch)
					return
				}
				err := c.UpdateControllerSwitchState(db, 0)
				if err != nil {
					fmt.Println(err)
				}
				c.Switch = 0
				fmt.Printf("ls switch off %v - state: %d\n", dataPoint, c.Switch)
			case "equal":
				if dataPoint == item.ThresholdLimit {
					fmt.Println("equal not implemented")
				}
			case "not equal":
				if dataPoint == item.ThresholdLimit {
					fmt.Println("not equal not implemented")
				}
			}
		}
	case "timeswitch":
	default:
		return
	}
}

func (c Controllers) UpdateControllerSwitchState(db *sqlx.DB, switchState int) error {
	_, err := db.Exec(`update controllers set
		switch=$1
		where id=$2
	`, switchState, c.ID)
	if err != nil {
		return fmt.Errorf("unable to update controller: %v", err)
	}
	return nil
}
