package models

import (
	"log"
	"net/http"

	helper "github.com/quarkey/iot/json"
)

// TempReading ....
type TempReading struct {
	SensorDescription string `json:"sensor_description"`
	Serial            string `json:"serial"`
	Temp              string `json:"temp"`
	Time              string `json:"time"`
}

// NewSensorReading ....
func (s *Server) NewSensorReading(w http.ResponseWriter, r *http.Request) {
	dat := TempReading{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	res, err := s.DB.Exec("insert into sensordata (sensor_description, serial, temp, time) values($1,$2,$3,$4)", dat.SensorDescription, dat.Serial, dat.Temp, dat.Time)

	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	if count, err := res.RowsAffected(); err == nil {
		helper.RespondSuccess(w, r, "Rows affected:", count)
	}
}

// Sensors ....
func (s *Server) Sensors(w http.ResponseWriter, r *http.Request) {
	helper.RespondErr(w, r, 500, "not implemented yet")
}
