package sensor

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
)

var TimeFormat = "2006-01-02 15:04:05"

// ExportSensorDataToCSV generates a csv dataset with corresponding columns
// Exported data will include id, and time columns and then data points
func ExportSensorDataToCSV(ref string, db *sqlx.DB) (interface{}, error) {
	datalabel, _, err := dataset.DatasetFieldAndShowCartList(ref, db)
	if err != nil {
		return nil, fmt.Errorf("unable to get datasetfields: %v", err)

	}
	dat, err := GetRawDataWithLimitByRef(db, 1000, ref)
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
		slice, _ := helper.DecodeRawJSONtoSlice(x.Data)
		row := []string{strconv.Itoa(x.ID), x.Time.Format(TimeFormat)}
		row = append(row, slice...)
		csv = append(csv, row)
	}
	return csv, nil
}
