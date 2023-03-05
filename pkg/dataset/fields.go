package dataset

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

func DatasetFieldAndShowCartList(ref string, db *sqlx.DB) ([]string, []bool, error) {
	fields, err := datasetFieldsList(ref, db)
	if err != nil {
		return nil, nil, err
	}
	showCharts, err := datasetShowChartBools(ref, db)
	if err != nil {
		return nil, nil, err
	}
	return fields, showCharts, nil
}

// DatasetFieldsList returns a list of column labels for a dataset by given reference
func datasetFieldsList(ref string, db *sqlx.DB) ([]string, error) {
	var raw *json.RawMessage
	err := db.Get(&raw, `select fields from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get labels from datasets: %v", err)
	}
	var fields []string
	err = json.Unmarshal(*raw, &fields)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v", err)
	}
	return fields, nil
}

// datasetShowChartBools returns a list of boolean values to indicate if chart should show or not.
// Method takes in account that showcharts fields from pg may be empty.
// Also i think there is a bug in this function, second marshal fails because
// pg res can be validated as bool [true, "true"]
func datasetShowChartBools(ref string, db *sqlx.DB) ([]bool, error) {
	// fetching fields to determine field count
	var raw *json.RawMessage
	err := db.Get(&raw, `select fields from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get fields from datasets: %v", err)
	}
	var fields []string
	err = json.Unmarshal(*raw, &fields)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v", err)
	}
	var nextRaw *json.RawMessage
	err = db.Get(&nextRaw, `select showcharts from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get showchart bools from dataset: %v", err)
	}
	var nextOut []bool
	// defaulting to false for all datasets
	// if boolean values is missing in the database
	if nextRaw == nil {
		if len(nextOut) == 0 {
			var tmp []bool
			for i := 0; i < len(fields); i++ {
				tmp = append(tmp, false)
			}
			nextOut = tmp
		}
		return nextOut, nil
	}
	err = json.Unmarshal(*nextRaw, &nextOut)
	if err != nil {
		if err.Error() == "json: cannot unmarshal string into Go value of type bool" {
			// we can handle them as bool if they are stored as strings in the database
			var str []string
			var tmp []bool
			err := json.Unmarshal(*nextRaw, &str)
			if err != nil {
				return nil, fmt.Errorf("DatasetShowChartBools: %v ", err)
			}
			for _, v := range str {
				boolValue, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("DatasetShowChartBools: %v", err)
				}
				tmp = append(tmp, boolValue)
			}
			nextOut = tmp
		}
	}
	return nextOut, nil
}
