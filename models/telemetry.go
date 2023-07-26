package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
)

// A telemetry struct holds the telemetry components in memory.
type Telemetry struct {
	done            chan bool // closes the time.Ticker go routine.
	db              *sqlx.DB
	datasets        []Dataset      // in memory datasets
	sensors         []SensorDevice // in memory sensors
	controllers     ControllerList // in memory controllers
	storageLocation string
}

// newTelemetryTicker ...
func newTelemetryTicker(db *sqlx.DB, path string) *Telemetry {
	return &Telemetry{
		done:            make(chan bool),
		db:              db,
		storageLocation: path,
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
						// telemetry.CheckSensorsTelemetry()
					}
					duration = 0
				}
				// controllers require extra precision so a check is done every second
				telemetry.CheckControllersTelemetry()
				duration++
			}
		}
	}()
}

// UpdateTelemetryLists updates sensors and dataset lists
// TODO: remove t.init method
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
	t.controllers, _ = GetControllersListFromDB(t.db)
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
		// t.CheckSensorsTelemetry()
		t.CheckDatasetTelemetry()
		t.CheckControllersTelemetry()
	}
}

// func (t *Telemetry) CheckSensorsTelemetry() {
// 	log.Println("[INFO] UpdateSensorsTelemetry() NOT IMPLEMENTED")
// 	// TODO: "ping" device by ip, waiting for arduino sketch
// }

// UpdateDatasetTelemetry updates the telemetry data for a given dataset.
// It also checks if the telemetry data is overdue by more than 60 seconds, and if so,
// marks the dataset as offline.
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

		// if (current time - next interval time) is more than 60 seconds
		// we can consider the telemetry to be offline
		if (time.Now().Unix() - timeFuture) > 60 {
			dataset.SetOfflineByID(t.db, dset.ID)
		}
	}
}

// var mutex sync.Mutex

// CheckControllersTelemetry ...
func (t *Telemetry) CheckControllersTelemetry() {

	for _, c := range t.controllers {
		if !c.Active {
			continue
		}
		switch c.Category {
		case "thresholdswitch":
			var ts []Thresholdswitch
			err := json.Unmarshal(*c.Items, &ts)
			if err != nil {
				log.Printf("[ERROR] unable to unmarshal thresholdswitch json: %v", err)
			}
			// to be able to check thresold dataset is required
			for _, item := range ts {
				if item.Datasource == "" {
					return
				}
				jsonRaw, _, column, err := dataset.GetLastDataSourcePointRaw(t.db, item.Datasource)
				if err != nil {
					log.Printf("[ERROR] %v", err)
					return
				}
				//fmt.Printf("checking datasource '%s'\n", item.Datasource)
				//fmt.Printf("dataset: %v and column: %v \n", dset, column)
				slice, err := helper.DecodeRawJSONtoSlice(jsonRaw)
				if err != nil {
					fmt.Printf("something went wrong... %v\n", err)
				}
				if len(slice) == 0 {
					log.Printf("[WARNING] controller '%s' skipped, empty dataset", c.Category)
				}
				// fmt.Printf("Data: %v\n", slice)
				datapoint, err := strconv.ParseFloat(slice[column], 64)
				if err != nil {
					fmt.Printf("something went wrong... %v\n", err)
					datapoint = 0
				}

				c.CheckThresholdSwitchEntries(datapoint, t.db)
			}
		case "timeswitch", "timeswitchrepeat":
			c.CheckTimeSwitchEntries(t.db)
		case "switch":
			// do we need to track switches other than sensor telemetry?
		case "webcamstreamtimelapse":
			// do timelapse capture
			c.CheckWebCamStreamEntries(t.db, t.storageLocation)
		default:
			log.Println("[ERROR] unsupported controller category:", c.Category)
		}
	}
}
