package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/sensor"
)

// A structure that holds Sensor data
type SensorData struct {
	ID            int              `db:"id" json:"id"`
	SensorID      int              `db:"sensor_id" json:"sensor_id"`
	DatasetID     int              `db:"dataset_id" json:"dataset_id"`
	Data          *json.RawMessage `db:"data" json:"data"`
	RecordingTime time.Time        `db:"time" json:"time"`
}

// SaveSensorReading is registering sensor readings (json) to database.
func (s *Server) SaveSensorReading(w http.ResponseWriter, r *http.Request) {
	dat := SensorData{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to read sensordata:", err)
		return
	}
	err = saveReadings([]SensorData{dat}, s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to save reading:", err)
		return
	}
	// setting dataset to online and broadcasting only when clients are connected

	dataset.SetOnlineByID(s.DB, dat.DatasetID)
	if len(s.Hub.Clients) > 0 {
		b, err := json.Marshal(&dat)
		if err != nil {
			helper.RespondErr(w, r, 500, "unable to marshal sensor", err)
		}
		s.Hub.Broadcast <- b
	}
}
func saveReadings(datapoints []SensorData, db *sqlx.DB) error {
	for _, r := range datapoints {
		_, err := db.Exec("insert into sensordata(sensor_id, dataset_id, data) values($1,$2,$3)", r.SensorID, r.DatasetID, r.Data)
		if err != nil {
			return fmt.Errorf("unablet to save sensor reading to db: %v", err)
		}
	}
	return nil
}
func saveReadingsTx(datapoints []SensorData, db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("tx failed: %v", err)
	}
	err = saveReadings(datapoints, db)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("tx failed: %v", err)
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}
	return nil
}

// SyncDataset ...
// this is not in use.
func (s *Server) SyncSensorData(w http.ResponseWriter, r *http.Request) {
	dat := []SensorData{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to read sensordata:", err)
		return
	}
	err = saveReadingsTx(dat, s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, "synchronization failed", err)
		return
	}
}

// Data JSON payload
type SensorRawJSONData struct {
	ID   int              `json:"id"`
	Data *json.RawMessage `json:"data"`
	Time time.Time        `json:"time"`
}

// GetSensorDataByReferenceEndpoint fetches a dataset by a reference
func (s *Server) GetSensorDataByReferenceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	data, err := sensor.GetRawDataWithLimitByRef(s.DB, 1000, vars["reference"])
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to get dataset from db:", err)
		return
	}
	helper.Respond(w, r, 200, data)
}

// ExportSensorDataToCSVEndpoint for exporting datasets used to create CSV report
func (s *Server) ExportSensorDataToCSVEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result, err := sensor.ExportSensorDataToCSV(vars["reference"], s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to export dataset to csv:", err)
		return
	}
	helper.Respond(w, r, 200, result)
}
