package models

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	helper "github.com/quarkey/iot/json"
)

type LineChart struct {
	Labels           []string           `json:"labels"`
	LineChartDataset []LineChartDataset `json:"lineChartdataset"`
}
type LineChartDataset struct {
	Data  []float64 `json:"data"`
	Label string    `json:"label"`
}

// LineChartDataSeries will generate a data structure that is fitted to ng2-charts.
func (s *Server) LineChartDataSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ref := vars["reference"]
	data, err := loadSensorData(s.DB, ref, 1000)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to load data: ", err)
		return
	}
	if len(data) == 0 {
		helper.RespondErr(w, r, 400, "no data for line chart")
		return
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSONtoSlice(data[0].Data)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}

	// fetching fields and showcharts list.
	// showcharts used to determine if chart should be added or not.
	fields, showcharts, err := DatasetFieldAndShowCartList(ref, s.DB)
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	// if we encounter more data points than columns we'll just add a unknown one
	if len(raw) > len(fields) {
		fields = append(fields, "unknow column")
		showcharts = append(showcharts, true)
	}
	var series LineChart
	var out []LineChartDataset
	for i := 0; i < len(raw); i++ {
		if !showcharts[i] {
			continue
		}
		var ps LineChartDataset
		ps.Label = fields[i]
		for _, set := range data {
			decoded, err := helper.DecodeRawJSONtoSlice(set.Data)
			if err != nil {
				helper.RespondErr(w, r, 500, err)
				return
			}
			toFloatValue, err := strconv.ParseFloat(decoded[i], 64)
			if err != nil {
				msg := fmt.Sprintf("LineChartDataSeries(): unable to parse data point '%v', check type for column %d", decoded[i], i)
				helper.RespondErr(w, r, 500, msg)
				return
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
	helper.Respond(w, r, 200, series)
}

type AreaChart struct {
	Name  string  `json:"name"`
	Point []Point `json:"series"`
}
type Point struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// AreaChartDataSeries will generate a data structure that is fitted to ngx-charts.
func (s *Server) AreaChartDataSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ref := vars["reference"]
	data, err := loadSensorData(s.DB, ref, 1000)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	if len(data) == 0 {
		helper.RespondErr(w, r, 400, "no data for area chart")
		return
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSONtoSlice(data[0].Data)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	fields, showcharts, err := DatasetFieldAndShowCartList(ref, s.DB)
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	// if we encounter more data points than columns we'll just add a unknown one
	if len(raw) > len(fields) {
		fields = append(fields, "unknow field")
		showcharts = append(showcharts, true)
	}
	// populate data structure fitted for ngx-charts
	var out []AreaChart
	for i := 0; i < len(raw); i++ {
		// skipping disabled charts
		if !showcharts[i] {
			continue
		}
		var ps []Point
		for _, set := range data {
			decoded, err := helper.DecodeRawJSONtoSlice(set.Data)
			if err != nil {
				helper.RespondErr(w, r, 500, err)
				return
			}
			// converting data point from string to float
			toFloatValue, err := strconv.ParseFloat(decoded[i], 64)
			if err != nil {
				msg := fmt.Sprintf("AreaChartDataSeries(): unable to parse data point '%v', check type for column %d", decoded[i], i)
				helper.RespondErr(w, r, 500, msg)
				return
			}
			ps = append(ps, Point{Name: set.Time.String(), Value: toFloatValue})
		}
		out = append(out, AreaChart{Name: fields[i], Point: ps})
	}
	helper.Respond(w, r, 200, out)
}

// loadSensorData fetches the last n records available for given
// dataset id and reference.
func loadSensorData(db *sqlx.DB, ref string, limit int) ([]Data, error) {
	var data []Data
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
	`, ref, limit)
	if err != nil {
		return nil, fmt.Errorf("unable to get dataset from db: %v", err)
	}
	return data, nil
}
