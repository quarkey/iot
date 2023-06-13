package sensor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RawSensorData struct {
	ID   int              `json:"id"`
	Data *json.RawMessage `json:"data"`
	Time time.Time        `json:"time"`
}

// GetRawDataWithLimitByRef fetches the last n records available for given dataset id and reference.
func GetRawDataWithLimitByRef(db *sqlx.DB, limit int, reference string) ([]RawSensorData, error) {
	var data []RawSensorData
	err := db.Select(&data, `
	select * from (select 
		a.id,
		a.data,
		a.time 
		from sensordata a, datasets b 
		where b.reference=$1 
		and b.id = a.dataset_id
		order by id desc limit $2
	) as sortedItems -- required for pg sub selects.
	order by id asc;
	`, reference, limit)
	if err != nil {
		return nil, fmt.Errorf("unable to get dataset from db: %v", err)
	}
	return data, nil
}
func GetRawDataByDateAndRef(db *sqlx.DB, DateFrom string, DateTo string, reference string) ([]RawSensorData, error) {
	var data []RawSensorData
	err := db.Select(&data, `
		select 
			a.id,
			a.data,
			a.time 
		from sensordata a, datasets b
		where b.reference=$1
		and b.id = a.dataset_id
		and a.time between $2 and $3;`, reference, DateFrom, DateTo)
	if err != nil {
		return nil, fmt.Errorf("unable to get dataset from db: %v", err)
	}
	return data, nil
}
