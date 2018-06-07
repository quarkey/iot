package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	helper "github.com/quarkey/iot/json"
)

// TempReading ....
type SensorRead struct {
	SensorID  int              `json:"sensor_id"`
	DatasetID int              `json:"dataset_id"`
	Data      *json.RawMessage `json:"data"`
}

// NewSensorReading is registering sensor readings (json) to database.
func (s *Server) NewSensorReading(w http.ResponseWriter, r *http.Request) {
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

type Sensor struct {
	ID          int       `db:"id" json:"id"`
	Description string    `db:"description" json:"description"`
	ArduinoKey  string    `db:"arduino_key" json:"arduino_key"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// Sensors lists all available sensors in the database
func (s *Server) Sensors(w http.ResponseWriter, r *http.Request) {
	var sensors []Sensor
	err := s.DB.Select(&sensors, "select id, description, arduino_key, created_at from sensor")
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, sensors)
}
