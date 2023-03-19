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
	"github.com/quarkey/iot/pkg/event"
)

// Controller data structure
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

// SwitchDefaultValues initial state
var SwitchDefaultValues = `[{
	"item_description": "",
	"on": false
}]`

// ThresholdswitchDefaultValues initial state
var ThresholdswitchDefaultValues = `
[{
	"item_description": "",
	"datasource": "",
	"operation": "",
	"threshold_limit": null,
	"on": false
}]`

// TimesSwitchDefaultValues initial state
var TimesSwitchDefaultValues = `
[{
	"item_description": "",
	"time_on": "",
	"time_off": "",
	"duration": null,
	"repeat": false,
	"on": false
}]
`
var SWITCH_ON = 1
var SWITCH_OFF = 0

type Thresholdswitch struct {
	Description    string  `json:"item_description"`
	Datasource     string  `json:"datasource"`
	Operation      string  `json:"operation"`
	ThresholdLimit float64 `json:"threshold_limit"`
}
type Timeswitch struct {
	Description string `json:"item_description"`
	TimeOn      string `json:"time_on"`
	TimeOff     string `json:"time_off"`
	Duration    int    `json:"duration"`
	Repeat      bool   `json:"repeat"`
	On          bool   `json:"on"`
}

// GetControllersList fetches a list of available controllers including not active ones.
func GetControllersList(db *sqlx.DB) ([]Controller, error) {
	var list []Controller
	err := db.Select(&list, `
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
	from controllers
	order by id`)
	if err != nil {
		return nil, fmt.Errorf("unable to get list of controllers from database: %v", err)
	}
	return list, nil
}

// GetControllerByID loads a specific controller from database by given ID
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

// GetControllersListEndpoint fetches a list of all available controllers.
func (s *Server) GetControllersListEndpoint(w http.ResponseWriter, r *http.Request) {
	list, err := GetControllersList(s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
	}
	helper.Respond(w, r, 200, list)
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

// AddNewControllerEndpoint adds a new controller to the database with an initial switch state and alerting to false.
func (s *Server) AddNewControllerEndpoint(w http.ResponseWriter, r *http.Request) {
	var dat Controller
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	var itemJSON string
	if dat.Category == "thresholdswitch" {
		itemJSON = ThresholdswitchDefaultValues
	}
	if dat.Category == "switch" {
		itemJSON = SwitchDefaultValues
	}
	if dat.Category == "timeswitch" {
		itemJSON = TimesSwitchDefaultValues
	}
	if dat.Category == "timeswitchrepeat" {
		itemJSON = TimesSwitchDefaultValues
	}

	var returning_id int
	err = s.DB.QueryRow(`insert into controllers(sensor_id, category, title, description, items, alert, active)
	values($1, $2, $3, $4, $5, $6, $7) returning id;
	`,
		dat.SensorID,
		dat.Category,
		dat.Title,
		dat.Description,
		itemJSON,
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

	e := event.New(s.DB)
	e.NewEvent(ControllerEvent, "controller '%s' created", dat.Title)
	// s.NewEvent(DatasetEvent, "dataset '%s' updated", dat.Title)

	helper.Respond(w, r, 200, returning_id)
}

// UpdatControllerByIDEndpoint updates database and memory values for a  controller item by an ID. NB! The method will not validate *json.RawMessage.
func (s *Server) UpdateControllerByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	dat := Controller{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondErr(w, r, 500, err)
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
	e := event.New(s.DB)
	e.NewEvent(DatasetEvent, "controller '%s' updated", dat.Title)
	helper.Respond(w, r, 200, "updated")
}

// ResetControllerSwitchValueEndpoint resets the active configuration (raw json) for a given controller to its default state.
// The method also sets the controller switch state to OFF and inactive.
func (s *Server) ResetControllerSwitchValueEndpoint(w http.ResponseWriter, r *http.Request) {
	var dat Controller
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondErr(w, r, 500, err)
		return
	}
	var defaultValues string
	switch dat.Category {
	case "switch":
		defaultValues = SwitchDefaultValues
	case "thresholdswitch":
		defaultValues = ThresholdswitchDefaultValues
	case "timeswitch":
		defaultValues = TimesSwitchDefaultValues
	}
	_, err = s.DB.Exec(`update controllers set
		items=$1,
		switch=0,
		active='f'
		where id=$2
	`, defaultValues, dat.ID)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to update controller item: ", err)
		return
	}
	helper.RespondSuccess(w, r, 200)
}

// DeleteControllerByIDEndpoint deletes a controller record by a given ID without any confirmation.
// Caution is advised when deleting controllers, as this may affect the IoT system.
func (s *Server) DeleteControllerByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	var dat Controller
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondErr(w, r, 500, err)
		return
	}
	_, err = s.DB.Exec(`delete from controllers where id=$1`, dat.ID)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to delete controller: ", err)
		return
	}
	helper.RespondSuccess(w, r, 200)
}

type switchState struct {
	ID          int  `db:"id" json:"id"`
	Active      bool `db:"active" json:"active"`
	SwitchState int  `db:"switch" json:"switch"`
}

// SetControllerSwitchState allows the caller to change the state to either on or off, but for only active controllers.
func (s *Server) SetControllerSwitchStateEndpoint(w http.ResponseWriter, r *http.Request) {
	// check current status
	vars := mux.Vars(r)
	tmp := vars["id"]
	reqID, _ := strconv.Atoi(tmp)
	var sw switchState

	err := s.DB.Get(&sw, `select id, active, switch from controllers where id=$1`, reqID)
	if err != nil {
		helper.RespondErr(w, r, 200, "unable to get data from db: ", err)
		return
	}
	if !sw.Active {
		helper.RespondErr(w, r, 500, "unable to change switch state: controller is currently set to inactive")
		return
	}
	switch vars["state"] {
	case "on":
		if sw.SwitchState == SWITCH_ON {
			helper.RespondErr(w, r, 500, "switch state already on!")
			return
		}
		// turn the switch on
		updateControllerSwitchStatebyID(s.DB, reqID, SWITCH_ON)
		sw.SwitchState = SWITCH_ON
	case "off":
		if sw.SwitchState == SWITCH_OFF {
			helper.RespondErr(w, r, 500, "switch state already off!")
			return
		}
		// turn the switch off
		updateControllerSwitchStatebyID(s.DB, reqID, SWITCH_OFF)
		sw.SwitchState = SWITCH_OFF
	default:
		helper.RespondErr(w, r, 500, "unknown state")
		return
	}
	helper.Respond(w, r, 200, sw)
}

// CheckThresholdEntries checks if a list of threshold switches is within the boundaries of a given condition and turns them ON or OFF.
// Supported operations: greather than, less than, equal and not equal.
func (c *Controller) CheckThresholdSwitchEntries(dataPoint float64, db *sqlx.DB) {
	if !c.Active {
		return
	}
	var ts []Thresholdswitch
	err := json.Unmarshal(*c.Items, &ts)
	if err != nil {
		log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
	}
	for _, item := range ts {
		// threshold switch operation
		switch item.Operation {
		case "greather than":
			// switching state based on threshold
			if dataPoint > item.ThresholdLimit {
				err := c.UpdateControllerSwitchState(db, SWITCH_ON)
				if err != nil {
					fmt.Println(err)
				}
				// fmt.Printf("gt switch on: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)
				return
			}
			err := c.UpdateControllerSwitchState(db, SWITCH_OFF)
			if err != nil {
				fmt.Println(err)
			}
			c.Switch = SWITCH_OFF
			// fmt.Printf("gt switch off %v - state: %d\n", dataPoint, c.Switch)
			// fmt.Printf("gt switch off: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)

		case "less than":
			// switching state based on threshold
			// fmt.Println("datapoint", dataPoint, "threshold", item.ThresholdLimit)
			if dataPoint < item.ThresholdLimit {
				err := c.UpdateControllerSwitchState(db, SWITCH_ON)
				if err != nil {
					fmt.Println(err)
				}
				// c.Switch = 1
				// fmt.Printf("lt switch on: %s -> %v - state: %d\n", c.Title, dataPoint, c.Switch)
				// fmt.Printf("lt switch on: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)
				return
			}
			err := c.UpdateControllerSwitchState(db, SWITCH_OFF)
			if err != nil {
				fmt.Println(err)
			}
			c.Switch = SWITCH_OFF
			// fmt.Printf("lt switch off: %s -> t:%v -> d:%v - state: %d\n", c.Title, item.ThresholdLimit, dataPoint, c.Switch)
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
}

func (c *Controller) ChecktimeSwitchRepeatEntries(db *sqlx.DB) {
	if !c.Active {
		return
	}
	var ts []Timeswitch
	err := json.Unmarshal(*c.Items, &ts)
	if err != nil {
		log.Printf("[ERROR] unable to unmarshal timeswitchrepeat json: %v", err)
	}

	for _, item := range ts {
		if isCurrentTimeBetween(item.TimeOn, item.TimeOff) {
			err := c.UpdateControllerSwitchState(db, 1)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = c.UpdateControllerSwitchState(db, 0)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// CheckTimeSwitchEtries checks if a list of time switches is within a timeframe and turns them ON or OFF.
func (c *Controller) ChecktimeSwitchEntries(db *sqlx.DB) {
	if !c.Active {
		return
	}
	var ts []Timeswitch
	err := json.Unmarshal(*c.Items, &ts)
	if err != nil {
		log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
	}
	for _, item := range ts {
		on, err := time.ParseInLocation(TimeFormat, item.TimeOn, time.Now().Location())
		if err != nil {
			fmt.Printf("unable to parse time: %v\n", err)
		}
		off, err := time.ParseInLocation(TimeFormat, item.TimeOff, time.Now().Location())
		if err != nil {
			fmt.Printf("unable to parse time: %v\n", err)
		}
		start := on.In(time.Now().Location())
		end := off.In(time.Now().Location())

		if inTimeSpan(start, end, time.Now()) {
			fmt.Printf("timeswitch: %s status: on \n", item.Description)
			err := c.UpdateControllerSwitchState(db, 1)
			if err != nil {
				fmt.Println(err)
			}
			return
		}
		err = c.UpdateControllerSwitchState(db, 0)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// inTimeSpan checks if "current time" is in between start and end time.
func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

// isCurrentTimeBetween checks if the current time is between two given times.
// t1 and t2 should be in the format "15:04" (HH:MM).
func isCurrentTimeBetween(t1, t2 string) bool {
	layout := "15:04:05"
	now := time.Now().Format(layout)
	return now >= t1 && now <= t2
}

// UpdateControllerSwitchState changes the controller switch state.
func (c *Controller) UpdateControllerSwitchState(db *sqlx.DB, switchState int) error {
	err := updateControllerSwitchStatebyID(db, c.ID, switchState)
	if err != nil {
		return err
	}
	return nil
}

// updateControllerSwitchStatebyID changes the state for a given controller id.
func updateControllerSwitchStatebyID(db *sqlx.DB, id, switchState int) error {
	_, err := db.Exec(`update controllers set
		switch=$1
		where id=$2
	`, switchState, id)
	if err != nil {
		return fmt.Errorf("unable to update controller: %v", err)
	}
	return nil
}
