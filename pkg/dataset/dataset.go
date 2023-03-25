package dataset

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

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

// GetSpecificSensorDataPoint takes a datasource string and extracts dataset id and column.
// e.g. d0c1 will return dataset_id=0, column=1
func GetSpecificSensorDataPoint(datasource string) (dataset_id, column int64) {
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

func GetLastDataSourcePointRaw(db *sqlx.DB, dataset string) (*json.RawMessage, int64, int64, error) {
	var data *json.RawMessage
	dset, column := GetSpecificSensorDataPoint(dataset)
	err := db.Get(&data, `
			select data from sensordata
			where dataset_id=$1
			order by id desc limit 1`, dset)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("unable to get last datapoint: %v", err)
	}
	return data, dset, column, nil
}
