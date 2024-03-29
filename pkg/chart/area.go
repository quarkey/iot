package chart

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/sensor"
)

type AreaChart struct {
	Name  string  `json:"name"`
	Point []point `json:"series"`
}
type point struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

func AreaChartDataSeries(db *sqlx.DB, ref string, limit int) (*[]AreaChart, error) {
	data, err := sensor.GetRawDataWithLimitByRef(db, limit, ref)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSONtoSlice(data[0].Data)
	if err != nil {
		return nil, err
	}
	fields, showcharts, err := dataset.DatasetFieldAndShowCartList(ref, db)
	if err != nil {
		return nil, err
	}
	// if we encounter more data points than columns we'll just add a unknown one
	if len(raw) > len(fields) {
		fields = append(fields, "unknow field")
		showcharts = append(showcharts, true)
	}
	// FIXME:  when adding a column to an existing dataset this function will produce
	// an error because the sensordata array is not matching the number of fields.
	// the code below checks if more fields than data is present, this is not a fix.
	isEmpty := false
	count := len(raw)
	if len(raw) < len(fields) {
		fmt.Print("hmmm")
		isEmpty = true
		count++
	}

	// populate data structure fitted for ngx-charts
	var out []AreaChart
	for i := 0; i < len(raw); i++ {
		// skipping disabled charts
		if !showcharts[i] {
			continue
		}
		if isEmpty {
			if (len(raw) - 1) == i {
				fmt.Println("should add dummy point here")
			}
		}
		var ps []point
		for _, set := range data {
			decoded, err := helper.DecodeRawJSONtoSlice(set.Data)
			if err != nil {
				return nil, err
			}
			// converting data point from string to float
			toFloatValue, err := strconv.ParseFloat(decoded[i], 64)
			if err != nil {
				return nil, fmt.Errorf("AreaChartDataSeries(): unable to parse data point '%v', check type in column '%s'", decoded[i], fields[i])

			}
			ps = append(ps, point{Name: set.Time.String(), Value: toFloatValue})
		}
		out = append(out, AreaChart{Name: fields[i], Point: ps})
	}
	return &out, nil
}
