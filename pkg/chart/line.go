package chart

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/sensor"
)

type LineChart struct {
	Labels           []string           `json:"labels"`
	LineChartDataset []lineChartDataset `json:"lineChartdataset"`
}
type lineChartDataset struct {
	Data  []float64 `json:"data"`
	Label string    `json:"label"`
}

func LineDataSeries(db *sqlx.DB, ref string) (*LineChart, error) {
	data, err := sensor.GetRawDataWithLimitByRef(db, 1000, ref)
	if err != nil {
		return nil, fmt.Errorf("LineDataSeries() unable to load data: %v", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSONtoSlice(data[0].Data)
	if err != nil {
		return nil, err
	}

	// fetching fields and showcharts list.
	// showcharts used to determine if chart should be added or not.
	fields, showcharts, err := dataset.DatasetFieldAndShowCartList(ref, db)
	if err != nil {
		return nil, err
	}
	// if we encounter more data points than columns we'll just add a unknown one
	if len(raw) > len(fields) {
		fields = append(fields, "unknow column")
		showcharts = append(showcharts, true)
	}
	var series LineChart
	var out []lineChartDataset
	for i := 0; i < len(raw); i++ {
		if !showcharts[i] {
			continue
		}
		var ps lineChartDataset
		ps.Label = fields[i]
		for _, set := range data {
			decoded, err := helper.DecodeRawJSONtoSlice(set.Data)
			if err != nil {
				return nil, err
			}
			toFloatValue, err := strconv.ParseFloat(decoded[i], 64)
			if err != nil {
				return nil, fmt.Errorf("LineChartDataSeries(): unable to parse data point '%v', check type in column '%s'", decoded[i], fields[i])
			}
			ps.Data = append(ps.Data, toFloatValue)
		}
		out = append(out, ps)
	}
	// labels
	for _, set := range data {
		series.Labels = append(series.Labels, set.Time.Local().String())
	}
	series.LineChartDataset = out
	return &series, nil
}
