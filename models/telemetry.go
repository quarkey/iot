package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Telemetry struct {
	done     chan bool
	db       *sqlx.DB
	datasets []Dataset
	sensors  []Sensor
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
	durfloat64 := cfg["checkTelemetryTimer"].(float64)
	ticker := time.NewTicker(time.Duration(int(durfloat64)) * time.Second)
	log.Printf("[INFO] Telemetry check every %d seconds\n", int(durfloat64))
	telemetry.init(true)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if len(telemetry.datasets) > 0 {
					log.Println("[INFO] Telemetry check", t)
					telemetry.CheckDatasetTelemetry()
					telemetry.CheckSensorsTelemetry()
				}
			}
		}
	}()
}

// UpdateTelemetryLists updates sensors and dataset lists
func (t *Telemetry) UpdateTelemetryLists() {
	t.init(false)
}

// init loads sensors and dataset into memory, caller can initiate telemetry check.
func (t *Telemetry) init(updateTelemetry bool) {
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

	for _, dset := range t.sensors {
		fmt.Printf("=> monitoring sensor telemetry for '%s'\n", dset.Title)
	}

	for _, dset := range t.datasets {
		fmt.Printf("=> monitoring dataset telemetry for '%s'\n", dset.Title)
	}
	if updateTelemetry {
		t.CheckSensorsTelemetry()
		t.CheckDatasetTelemetry()
	}
}

func (t *Telemetry) CheckSensorsTelemetry() {
	log.Println("[INFO] UpdateSensorsTelemetry() NOT IMPLEMENTED")
	// TODO: "ping" device by ip, waiting for arduino sketch
}

// UpdateDatasetTelemetry checks for duration between collected signals
// and tries to determine connectivity.
func (t *Telemetry) CheckDatasetTelemetry() {
	for _, dset := range t.datasets {
		// running query to get last signal received
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
		now := time.Now().UTC()
		tx := unixdiff(now, sd.RecordingTime)
		msg := fmt.Sprintf("since last telemetry %s", tx.duration())
		_, err = t.db.Exec("update iot.sensors set dataset_telemetry=$1 where id=$2", msg, sd.SensorID)
		if err != nil {
			log.Printf("[ERROR] unable to update dataset telemetry status: %v", err)
		}
	}
}

type diffunix struct {
	t1 int64
	t2 int64
}

// unixdiff ...
func unixdiff(t1, t2 time.Time) *diffunix {
	return &diffunix{
		t1: t1.Unix(),
		t2: t2.Unix(),
	}
}

// diff ...
func (d *diffunix) diff() int64 {
	return d.t1 - d.t2
}

// duration ...
func (d *diffunix) duration() time.Duration {
	return time.Duration(d.diff() * int64(time.Second))
}
