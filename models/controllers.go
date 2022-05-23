package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
)

type Controller struct {
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

var ThresholdswitchDefaultValues = `
[{
	"item_description": "",
	"datasource": "",
	"operation": "",
	"threshold_limit": null,
	"on": false
}]`

type Thresholdswitch struct {
	Description    string  `json:"item_description"`
	Datasource     string  `json:"datasource"`
	Operation      string  `json:"operation"`
	ThresholdLimit float64 `json:"threshold_limit"`
}

// GetControllersList fetches a list of controllers. All types of errors will return empty a slice.
func GetControllersList(db *sqlx.DB) ([]Controller, error) {
	var cs []Controller
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

// GetControllerByID loads a specific controller from database by given id
func GetControllerByID(db *sqlx.DB, cid int) (Controller, error) {
	var c Controller
	err := db.Get(&c, `
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
	from controllers where id=$1`, cid)
	if err != nil {
		return c, fmt.Errorf("unable to get '%d' controllers: %v", cid, err)
	}
	return c, nil
}

// GetControllersListEndpoint loads all available controllers
func (s *Server) GetControllersListEndpoint(w http.ResponseWriter, r *http.Request) {
	cs, err := GetControllersList(s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
	}
	helper.Respond(w, r, 200, cs)
}

// GetControllerByIDEndpoint loads a specific controller from database by given id
func (s *Server) GetControllerByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, _ := strconv.Atoi(vars["cid"])
	c, err := GetControllerByID(s.DB, n)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	helper.Respond(w, r, 200, c)
}

// AddControllerEndpoint adds a new controller to the database with initial switch state set to false
func (s *Server) AddNewControllerEndpoint(w http.ResponseWriter, r *http.Request) {
	var dat Controller
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	var returning_id int
	err = s.DB.QueryRow(`insert into controllers(sensor_id, category, title, description, items, alert, active)
	values($1, $2, $3, $4, $5, $6, $7) returning id;
	`,
		dat.SensorID,
		dat.Category,
		dat.Title,
		dat.Description,
		//dat.Switch,
		//dat.Items,
		//dat.Alert,
		//dat.Active,
		ThresholdswitchDefaultValues,
		false,
		false).Scan(&returning_id)

	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}

	// out := fmt.Sprintf(`{"returning_id": "%d"}`, returning_id)
	// also update telemetry dataset list
	s.Telemetry.UpdateTelemetryLists()
	s.NewEvent(DatasetEvent, "dataset '%s' updated", dat.Title)

	helper.Respond(w, r, 200, returning_id)

}

// UpdateControllerByIDEndpoint updates the entire row including json items
func (s *Server) UpdateControllerByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	dat := Controller{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	_, err = s.DB.Exec(`
	update iot.controllers set 
		category=$1,
		title=$2,
		description=$3,
		items=$4,
		active=$5
	where id=$6`,
		dat.Category,
		dat.Title,
		dat.Description,
		dat.Items,
		dat.Active,
		dat.ID,
	)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondErr(w, r, 500, err)
		return
	}
	// also update telemetry dataset list
	s.Telemetry.UpdateTelemetryLists()
	s.NewEvent(DatasetEvent, "controller '%s' updated", dat.Title)
	helper.Respond(w, r, 200, "updated")
}

func (c Controller) Check(dataPoint float64, db *sqlx.DB) {
	if !c.Active {
		return
	}
	switch c.Category {
	case "switch":
	case "thresholdswitch":
		var ts []Thresholdswitch
		err := json.Unmarshal(*c.Items, &ts)
		if err != nil {
			log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
		}
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
					fmt.Printf("gt switch on: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)
					return
				}
				err := c.UpdateControllerSwitchState(db, 0)
				if err != nil {
					fmt.Println(err)
				}
				c.Switch = 0
				// fmt.Printf("gt switch off %v - state: %d\n", dataPoint, c.Switch)
				fmt.Printf("gt switch off: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)

			case "less than":
				// switching state based on threshold
				// fmt.Println("datapoint", dataPoint, "threshold", item.ThresholdLimit)
				if dataPoint < item.ThresholdLimit {
					err := c.UpdateControllerSwitchState(db, 1)
					if err != nil {
						fmt.Println(err)
					}
					c.Switch = 1
					// fmt.Printf("lt switch on: %s -> %v - state: %d\n", c.Title, dataPoint, c.Switch)
					fmt.Printf("lt switch on: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)

					return
				}
				err := c.UpdateControllerSwitchState(db, 0)
				if err != nil {
					fmt.Println(err)
				}
				c.Switch = 0
				fmt.Printf("lt switch off: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)
				// fmt.Printf("ls switch off %v - state: %d\n", dataPoint, c.Switch)
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

func (c Controller) UpdateControllerSwitchState(db *sqlx.DB, switchState int) error {
	_, err := db.Exec(`update controllers set
		switch=$1
		where id=$2
	`, switchState, c.ID)
	if err != nil {
		return fmt.Errorf("unable to update controller: %v", err)
	}
	return nil
}
