package dataset

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// ResetConnectivity sets all dataset telemetry status to offline.
func ResetConnectivity(db *sqlx.DB) error {
	_, err := db.DB.Exec(`update datasets set telemetry='offline'`)
	if err != nil {
		return fmt.Errorf("unable to set dataset telemetry to 'offline': %v", err)
	}
	return nil
}

// SetOnlineByID updating dataset and setting telemetry 'online'
func SetOnlineByID(db *sqlx.DB, id int) {
	_, err := db.Exec(`update datasets set telemetry='online' where id=$1`, id)
	if err != nil {
		log.Printf("[ERROR] unable to set dataset telemetry to 'online': %v", err)
	}
}

// SetOfflineByID updating dataset and setting telemetry 'offline'
func SetOfflineByID(db *sqlx.DB, id int) {
	_, err := db.Exec(`update datasets set telemetry='offline' where id=$1`, id)
	if err != nil {
		log.Printf("[ERROR] unable to set dataset telemetry to 'offline': %v", err)
	}
}
