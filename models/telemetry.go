package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Telemetry struct {
	done chan bool
	db   *sqlx.DB
}

func newTelemetryTicker(db *sqlx.DB) *Telemetry {
	return &Telemetry{
		done: make(chan bool),
		db:   db,
	}
}

// startTicker starts telemetry ticker
func (telemetry *Telemetry) startTelemetryTicker(cfg map[string]interface{}) {
	// for some reason map[string]interface{} thinks value is float64
	durfloat64 := cfg["checkTelemetryTimer"].(float64)
	ticker := time.NewTicker(time.Duration(int(durfloat64)) * time.Second)
	log.Printf("[INFO] Telemetry check every %d seconds\n", int(durfloat64))
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("[INFO] Telemetry check", t)
				telemetry.UpdateDatasetTelemetry()
			}
		}
	}()
}
func (t *Telemetry) UpdateDatasetTelemetry() {
	fmt.Println("WWW")
}
