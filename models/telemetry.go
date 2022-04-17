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
}

func newTelemetryTicker(db *sqlx.DB) *Telemetry {
	return &Telemetry{
		done: make(chan bool),
		db:   db,
	}
}

// startTicker starts telemetry ticker, running a initial check on startup
func (telemetry *Telemetry) startTelemetryTicker(cfg map[string]interface{}) {
	// for some reason map[string]interface{} thinks value is float64
	durfloat64 := cfg["checkTelemetryTimer"].(float64)
	ticker := time.NewTicker(time.Duration(int(durfloat64)) * time.Second)
	log.Printf("[INFO] Telemetry check every %d seconds\n", int(durfloat64))
	telemetry.init()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				if len(telemetry.datasets) == 0 {
					return
				}
				log.Println("[INFO] Telemetry check", t)
				telemetry.UpdateDatasetTelemetry()
			}
		}
	}()
}

// init loads dataset into memory
func (t *Telemetry) init() {
	t.datasets = GetDatasetsList(t.db)
	log.Printf("[INFO] Loading datasets ...")
	if len(t.datasets) == 0 {
		log.Printf("[WARNING] No datasets available in database")
	}
	for _, dset := range t.datasets {
		fmt.Printf("=> monitoring telemetry for '%s'\n", dset.Title)
	}
	t.UpdateDatasetTelemetry()
}

// UpdateDatasetTelemetry checks for duration between collected signal
// and tries to determine connectivity.
func (t *Telemetry) UpdateDatasetTelemetry() {
	for _, dset := range t.datasets {
		// running query to get last signal received
		var sd SensorData
		err := t.db.Get(&sd, "select id, sensor_id, dataset_id, data, time from sensordata where dataset_id=$1 and sensor_id=$2 order by id desc limit 1", dset.ID, dset.SensorID)
		if err != nil {
			if err == sql.ErrNoRows {
				return
			}
			log.Printf("[ERROR] problems with selecting sensordata: %v", err)
		}
		now := time.Now()
		tx := unixdiff(now, sd.RecordingTime)
		// checking the difference between now and last recording time,
		// if value is over intervalSec we can say that the sensor is running late.
		// This could indicate issues with the sensor device itself or server performance.
		if tx.diff() > int64(dset.IntervalSec) {
			log.Printf("[WARNING] no telemetry for %v: (sensor_id: %d - %s)\n", tx.duration(), sd.SensorID, dset.SensorTitle)
			msg := fmt.Sprintf("no telemetry for %s", tx.duration())
			_, err := t.db.Exec("update iot.sensors set dataset_telemetry=$1 where id=$2", msg, sd.SensorID)
			if err != nil {
				log.Printf("[ERROR] unable to update dataset_telemetry status: %v", err)
			}
		} else {
			msg := fmt.Sprintf("telemetry ok %s", tx.duration())
			_, err := t.db.Exec("update iot.sensors set dataset_telemetry=$1 where id=$2", msg, sd.SensorID)
			if err != nil {
				log.Printf("[ERROR] unable to update dataset_telemetry status: %v", err)
			}
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
