package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
	"github.com/quarkey/iot/pkg/event"
)

// Sensor meta information
type SensorDevice struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	ArduinoKey  string    `db:"arduino_key" json:"arduino_key"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	// nullString because selecting a device without reference
	// will produce records with empty values
	DatasetTelem *sql.NullString `db:"dataset_telemetry" json:"dataset_telemetry"`
	SensorIP     string          `db:"sensor_ip" json:"sensor_ip"`
}

// GetSensorsList fetches a list of all available sensors in the database
func (s *Server) GetSensorsListEndpoint(w http.ResponseWriter, r *http.Request) {
	var sensors []SensorDevice
	err := s.DB.Select(&sensors, `
	select
		id,
		title,
		description,
		arduino_key,
		dataset_telemetry,
		sensor_ip,
		created_at 
	from sensors order by id`)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to select sensorlist: ", err)
		return
	}
	helper.Respond(w, r, 200, sensors)
}

func GetSensorsList(db *sqlx.DB) []SensorDevice {
	// TODO: add error handling
	var sensors []SensorDevice
	err := db.Select(&sensors, `
	select
		id,
		title,
		description,
		arduino_key,
		dataset_telemetry,
		sensor_ip,
		created_at 
	from sensors`)
	if err != nil {
		return nil
	}
	return sensors
}

// GetSensorByReference is looking up a particular sensor based on a arduino_key
func (s *Server) GetSensorByReference(w http.ResponseWriter, r *http.Request) {
	var sensor SensorDevice
	vars := mux.Vars(r)
	err := s.DB.Get(&sensor, `
	select
		id,
		title,
		description,
		arduino_key,
		created_at,
		sensor_ip,
		dataset_telemetry
	from sensors 
	where arduino_key=$1`,
		vars["reference"])
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondErr(w, r, 500, "unable to get sensor by reference:", err)
		return
	}
	helper.Respond(w, r, 200, sensor)
}

type NewDevice struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// AddNewDevice adds a fresh device to the database
func (s *Server) AddNewDevice(w http.ResponseWriter, r *http.Request) {
	dat := NewDevice{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to read sensordata:", err)
		return
	}
	id := uuid.New()
	_, err = s.DB.Exec("insert into iot.sensors(title, description, arduino_key) values($1, $2, $3) returning arduino_key", dat.Title, dat.Description, id.String())
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to add new sensor device:", err)
		return
	}
	//TODO only return uuid, full device not needed
	var device SensorDevice
	err = s.DB.Get(&device, `
	select
		id,
		title,
		description,
		arduino_key,
		created_at,
		dataset_telemetry,
		sensor_ip
	from sensors 
	where arduino_key=$1`, id)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondErr(w, r, 500, "unable to get sensor by reference:", err)
		return
	}
	s.Telemetry.UpdateTelemetryLists()

	e := event.New(s.DB)
	e.NewEvent(SernsorEvent, "sensor '%s' added", device.Title)

	helper.Respond(w, r, 200, device)
}

// UpdateDevice updates sensor metadata fields
func (s *Server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	var device = SensorDevice{}
	err := helper.DecodeBody(r, &device)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondErr(w, r, 500, "unable to decode body: ", err)
		return
	}
	_, err = s.DB.Exec(`update iot.sensors set title=$1, description=$2, sensor_ip=$3 where arduino_key=$4`,
		device.Title,
		device.Description,
		device.SensorIP,
		device.ArduinoKey,
	)
	fmt.Println(device)
	if err != nil {
		log.Println(err)
		helper.RespondErr(w, r, 500, "unable to update device: ", err)
		return
	}
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.NewEvent(SernsorEvent, "sensor '%s' updated", device.Title)
	helper.RespondSuccess(w, r)
}
