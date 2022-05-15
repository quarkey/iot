package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
)

// SensorData ....
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
	SetDatasetIDOnline(s.DB, dat.DatasetID)
	s.Telemetry.AutomaticUpdateSensorIPAdress(r.RemoteAddr, dat.SensorID)
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

// Sensor meta information
type Sensor struct {
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
	var sensors []Sensor
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

func GetSensorsList(db *sqlx.DB) []Sensor {
	var sensors []Sensor
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
	var sensor Sensor
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
	data, err := getSensorDataByReference(s.DB, vars["reference"])
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to get dataset from db:", err)
		return
	}
	helper.Respond(w, r, 200, data)
}
func getSensorDataByReference(db *sqlx.DB, ref string) ([]Data, error) {
	var data []Data
	err := db.Select(&data, `
	select 
		a.id,
		a.data,
		a.time 
	from sensordata a, datasets b 
	where b.reference=$1 
	and b.id = a.dataset_id`,
		ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get sensor data")
	}
	return data, nil
}

// ExportSensorDataToCSV generates a csv dataset with corresponding columns
// Exported data will include id, and time columns and then data points
func ExportSensorDataToCSV(ref string, db *sqlx.DB) (interface{}, error) {
	datalabel, _, err := DatasetFieldAndShowCartList(ref, db)
	if err != nil {
		return nil, fmt.Errorf("unable to get datasetfields: %v", err)

	}
	dat, err := getSensorDataByReference(db, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get data: %v", err)
	}
	var csv [][]string
	var header []string
	// adding id and time columns
	header = append(header, "id", "time")
	header = append(header, datalabel...)
	csv = append(csv, header)
	for _, x := range dat {
		slice, _ := helper.DecodeRawJSON(x.Data)
		row := []string{strconv.Itoa(x.ID), x.Time.Format(TimeFormat)}
		row = append(row, slice...)
		csv = append(csv, row)
	}
	return csv, nil
}

// ExportSensorDataToCSVEndpoint for exporting datasets used to create CSV report
func (s *Server) ExportSensorDataToCSVEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	result, err := ExportSensorDataToCSV(vars["reference"], s.DB)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to export dataset to csv:", err)
		return
	}
	helper.Respond(w, r, 200, result)
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
	var device Sensor
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
	s.NewEvent(SernsorEvent, "sensor '%s' added", device.Title)
	helper.Respond(w, r, 200, device)
}

// UpdateDevice updates sensor metadata fields
func (s *Server) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	var device = Sensor{}
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
	s.NewEvent(SernsorEvent, "sensor '%s' updated", device.Title)
	helper.RespondSuccess(w, r)
}
