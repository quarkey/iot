package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	err = s.RegisterMetricsDBValue(dat)
	if err != nil {
		log.Printf("[ERROR] unable to register metrics to db: %v", err)
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

func (s *Server) RegisterMetricsDBValue(data SensorData) error {
	values, err := helper.DecodeRawJSONtoSlice(data.Data)
	if err != nil {
		return fmt.Errorf("[ERROR] unable to decode raw json to slice when registering metrics to db: %v", err)
	}
	for i, point := range values {
		num, err := strconv.ParseFloat(point, 64)
		if err != nil {
			return fmt.Errorf("[ERROR] unable to parse float when registering metrics to db: %v", err)
		}
		name := fmt.Sprintf("iot_dataset_%d_col%d", data.DatasetID, i)
		_, err = s.DB.Exec("update iot.metrics set value = $1 where name = $2", num, name)
		if err != nil {
			return fmt.Errorf("[ERROR] unable to update metrics in the database: %v", err)
		}
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
	ref := chi.URLParam(r, "reference")
	limitReq := chi.URLParam(r, "limit")
	limit := 1000
	if limitReq != "" {
		n, err := strconv.Atoi(limitReq)
		if err != nil {
			helper.RespondErr(w, r, 500, "unable parse limit parameter:", err)
			return
		}
		limit = n
	}

	data, err := sensor.GetRawDataWithLimitByRef(s.DB, limit, ref)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to get dataset from db:", err)
		return
	}
	helper.Respond(w, r, 200, data)
}

// ExportSensorDataToCSVEndpoint for exporting datasets used to create CSV report
func (s *Server) ExportSensorDataToCSVEndpoint(w http.ResponseWriter, r *http.Request) {
	result, err := sensor.ExportSensorDataToCSV(chi.URLParam(r, "reference"), s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to export dataset to csv:", err)
		return
	}
	helper.Respond(w, r, 200, result)
}
