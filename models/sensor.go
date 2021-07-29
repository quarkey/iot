package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	helper "github.com/quarkey/iot/json"
)

// SensorRead ....
type SensorRead struct {
	SensorID  int              `json:"sensor_id"`
	DatasetID int              `json:"dataset_id"`
	Data      *json.RawMessage `json:"data"`
}

// SaveSensorReading is registering sensor readings (json) to database.
func (s *Server) SaveSensorReading(w http.ResponseWriter, r *http.Request) {
	// TODO: must verify arduino_key on input
	dat := SensorRead{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	_, err = s.DB.Exec("insert into sensordata(sensor_id, dataset_id, data) values($1,$2,$3)", dat.SensorID, dat.DatasetID, dat.Data)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
}

// Sensor meta information
type Sensor struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	ArduinoKey  string    `db:"arduino_key" json:"arduino_key"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// GetSensorsList fetches a list of all available sensors in the database
func (s *Server) GetSensorsList(w http.ResponseWriter, r *http.Request) {
	var sensors []Sensor
	err := s.DB.Select(&sensors, "select id, title, description, arduino_key, created_at from sensors")
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, sensors)
}

// GetSensorByReference is looking up a particular sensor based on a arduino_key
func (s *Server) GetSensorByReference(w http.ResponseWriter, r *http.Request) {
	var sensor Sensor
	vars := mux.Vars(r)
	err := s.DB.Get(&sensor, "select id, title, description, arduino_key, created_at from sensors where arduino_key=$1", vars["reference"])
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, sensor)
}

// Data JSON payload
type Data struct {
	ID   int              `json:"id"`
	Data *json.RawMessage `json:"data"`
	Time time.Time        `json:"time"`
}

// GetSensorDataByReference fetches a dataset by a reference
func (s *Server) GetSensorDataByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data []Data
	var datasetID int
	err := s.DB.Get(&datasetID, "select id from dataset where reference=$1", vars["reference"])
	if err != nil {
		log.Printf("unable to run query referece: %v", err)
	}
	err = s.DB.Select(&data, "select id, data, time from sensordata where dataset_id=$1;", datasetID)
	if err != nil {
		log.Printf("unable to run query two: %v", err)
	}
	helper.Respond(w, r, 200, data)
}
