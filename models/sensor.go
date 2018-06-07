package models

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	helper "github.com/quarkey/iot/json"
)

// TempReading ....
type SensorRead struct {
	SensorID  int       `json:"sensor_id"`
	DatasetID int       `json:"dataset_id"`
	Data      []float64 `json:"data"`
}

// NewSensorReading ....
func (s *Server) NewSensorReading(w http.ResponseWriter, r *http.Request) {
	dat := SensorRead{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	_, err = s.DB.Exec("insert into sensordata(sensor_id, dataset_id, data) values($1,$2,$3)", dat.SensorID, dat.DatasetID, arrayToString(dat.Data, ","))
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
}

// Sensors ....
func (s *Server) Sensors(w http.ResponseWriter, r *http.Request) {
	helper.RespondErr(w, r, 500, "not implemented yet")
}

// arrayToString converts a slice of floats to json array
func arrayToString(a []float64, delim string) string {
	return fmt.Sprintf("[%s]", strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]"))
}
