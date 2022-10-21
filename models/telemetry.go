package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
)

// A telemetry struct holds the telemetry components in memory.
type Telemetry struct {
	done        chan bool // closes the time.Ticker go routine.
	db          *sqlx.DB
	datasets    []Dataset      // in memory datasets
	sensors     []SensorDevice // in memory sensors
	controllers []Controller   // in memory controllers
}

// newTelemetryTicker ...
func newTelemetryTicker(db *sqlx.DB) *Telemetry {
	return &Telemetry{
		done: make(chan bool),
		db:   db,
	}
}

// startTicker starts telemetry ticker, running a initial check on startup.
func (telemetry *Telemetry) startTelemetryTicker(cfg map[string]interface{}, debug bool) {
	// for some reason map[string]interface{} thinks value is float64
	checkTelemetryTimer := cfg["checkTelemetryTimer"].(float64)
	duration := 0

	ticker := time.NewTicker(time.Duration(1 * time.Second))
	log.Printf("[INFO] Telemetry check every %d seconds\n", int(checkTelemetryTimer))
	telemetry.init(true)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if duration > int(checkTelemetryTimer) {
					if len(telemetry.datasets) > 0 {
						log.Println("[INFO] Telemetry check")
						telemetry.CheckDatasetTelemetry()
						telemetry.CheckSensorsTelemetry()
					}
					duration = 0
				}
				// controllers require extra precision so check is done every second
				telemetry.CheckControllersTelemetry()
				duration++
			}
		}
	}()
}

// UpdateTelemetryLists updates sensors and dataset lists
func (t *Telemetry) UpdateTelemetryLists() {
	t.init(false)
}

// init loads sensors, dataset and controllers into memory, caller can initiate telemetry check
// by passing true.
func (t *Telemetry) init(runTelemetryCheck bool) {
	// loading sensor into memory
	t.sensors = GetSensorsList(t.db)
	log.Printf("[INFO] loading telemetry sensor list...")
	if len(t.sensors) == 0 {
		log.Printf("[WARNING] no sensors available in database")
	}

	// loading dataset into memory
	t.datasets = GetDatasetsList(t.db)
	log.Printf("[INFO] loading telemetry dataset list...")
	if len(t.datasets) == 0 {
		log.Printf("[WARNING] No datasets available in database")
	}

	// loading controllers into memory
	t.controllers, _ = GetControllersList(t.db)
	// TODO handle error
	log.Printf("[INFO] loading telemetry controllers list...")
	if len(t.controllers) == 0 {
		log.Printf("[WARNING] No active controllers")
	}

	for _, dset := range t.sensors {
		fmt.Printf("=> monitoring sensor telemetry for '%s'\n", dset.Title)
	}

	for _, dset := range t.datasets {
		fmt.Printf("=> monitoring dataset telemetry for '%s'\n", dset.Title)
	}

	for _, c := range t.controllers {
		fmt.Printf("=> monitoring controllers telemetry for '%v'\n", c.Title)
	}

	if runTelemetryCheck {
		t.CheckSensorsTelemetry()
		t.CheckDatasetTelemetry()
		t.CheckControllersTelemetry()
	}
}

func (t *Telemetry) CheckSensorsTelemetry() {
	log.Println("[INFO] UpdateSensorsTelemetry() NOT IMPLEMENTED")
	// TODO: "ping" device by ip, waiting for arduino sketch
}

// UpdateDatasetTelemetry updates dataset telemetry,
// but also setting dataset to offline if telemetry due are over 60 seconds.
func (t *Telemetry) CheckDatasetTelemetry() {
	for _, dset := range t.datasets {
		// running query to get last signal received
		//TODO: use loadSensorData
		var sd SensorData
		err := t.db.Get(&sd, `
		select
			id,
			sensor_id,
			dataset_id,
			data,
			time
		from sensordata 
		where dataset_id=$1 
		and sensor_id=$2 
		order by id desc limit 1`, dset.ID, dset.SensorID)
		if err != nil {
			if err == sql.ErrNoRows {
				msg := "no data points"
				log.Printf("[WARNING] no data points for dataset_id: %d with name '%s'\n", dset.ID, dset.SensorTitle)
				_, err := t.db.Exec("update iot.sensors set dataset_telemetry=$1 where id=$2", msg, dset.SensorID)
				if err != nil {
					log.Printf("[ERROR] unable to update dataset telemetry status: %v", err)
				}
				return
			}
			log.Printf("[ERROR] problems with selecting sensordata: %v", err)
		}
		// updating sensor dataset telemetry with last recorded time
		_, err = t.db.Exec("update iot.sensors set dataset_telemetry=$1 where id=$2", sd.RecordingTime, sd.SensorID)
		if err != nil {
			log.Printf("[ERROR] unable to update dataset telemetry status: %v", err)
		}
		// determine if dataset telemetry is "offline"

		timeFuture := sd.RecordingTime.Unix() + int64(dset.IntervalSec)

		// fmt.Println("current", time.Now().Unix())
		// fmt.Println("future", timeFuture)
		// fmt.Println("ti-tf", time.Now().Unix()-timeFuture)

		// if (current time - next interval time) are more than 60 seconds
		// we can consider the telemetry to be offline
		if (time.Now().Unix() - timeFuture) > 60 {
			SetDatasetIDOffline(t.db, dset.ID)
		}
	}
}

// CheckControllersTelemetry ...
func (t *Telemetry) CheckControllersTelemetry() {
	for _, c := range t.controllers {
		switch c.Category {
		case "thresholdswitch":
			var ts []Thresholdswitch
			err := json.Unmarshal(*c.Items, &ts)
			if err != nil {
				log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
			}
			for _, item := range ts {
				if item.Datasource == "" {
					return
				}
				dset, column := getSpecificSensorDataPoint(item.Datasource)
				var data *json.RawMessage
				err := t.db.Get(&data, `
			select data from sensordata
			where dataset_id=$1
			order by id desc limit 1`, dset)
				if err != nil {
					log.Println("[ERROR] problem fetching sensor data point values: ", err)
					return
				}
				//fmt.Printf("checking datasource '%s'\n", item.Datasource)
				//fmt.Printf("dataset: %v and column: %v \n", dset, column)
				slice, err := helper.DecodeRawJSONtoSlice(data)
				if err != nil {
					fmt.Printf("something went wrong... %v\n", err)
				}
				if len(slice) == 0 {
					fmt.Println("controller skipped, empty dataset")
				}
				// fmt.Printf("Data: %v\n", slice)
				datapoint, err := strconv.ParseFloat(slice[column], 64)
				if err != nil {
					fmt.Printf("something went wrong... %v\n", err)
					datapoint = 0
				}

				c.CheckThresholdSwitchEntries(datapoint, t.db)
			}
		case "timeswitch":
			c.ChecktimeSwitchEntries(t.db)
		case "switch":
			// do we need to track switches other than sensor telemetry?
		default:
			log.Println("[ERROR] unsupported controller category:", c.Category)
		}
	}
}

// getSpecificSensorDataPoint takes a datasource string and returns dataset id and column.
//
// e.g. d0c1 will return dataset_id=0, column=1
func getSpecificSensorDataPoint(datasource string) (dataset_id, column int64) {
	re := regexp.MustCompile("[0-9]+")
	slice := re.FindAllString(datasource, -1)
	s0, err := strconv.ParseInt(slice[0], 10, 64)
	if err != nil {
		return 0, 0
	}
	s1, err := strconv.ParseInt(slice[1], 10, 64)
	if err != nil {
		return 0, 0
	}
	return s0, s1
}
