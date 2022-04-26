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
	data, err := loadData(s.DB, ref)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
	}
	if len(data) == 0 {
		helper.RespondErr(w, r, 500, "no data avalable with reference: ", ref)
		return
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSON(data[0].Data)
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

	var series LineChart
	var out []LineChartDataset
	for i := 0; i < len(raw); i++ {
		if !showcharts[i] {
			continue
		}
		var ps LineChartDataset
		ps.Label = fields[i]
		for _, set := range data {
			decoded, err := helper.DecodeRawJSON(set.Data)
			if err != nil {
				helper.RespondErr(w, r, 500, err)
				return
			}
			ps.Data = append(ps.Data, toFloatValueFn(decoded[i]))
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
	data, err := loadData(s.DB, ref)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
	}
	// decoding jsonRawMessage data column
	raw, err := helper.DecodeRawJSON(data[0].Data)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	fields, showcharts, err := DatasetFieldAndShowCartList(ref, s.DB)
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
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
			decoded, err := helper.DecodeRawJSON(set.Data)
			if err != nil {
				helper.RespondErr(w, r, 500, err)
				return
			}
			// converting data point from string to float
			toFloatValue, err := strconv.ParseFloat(decoded[i], 64)
			if err != nil {
				helper.RespondErr(w, r, 500, err)
				return
			}
			ps = append(ps, Point{Name: set.Time.String(), Value: toFloatValue})
		}
		out = append(out, AreaChart{Name: fields[i], Point: ps})
	}
	helper.Respond(w, r, 200, out)
}

func loadData(db *sqlx.DB, ref string) ([]Data, error) {
	var data []Data
	err := db.Select(&data, `
	select 
		a.id,
		a.data,
		a.time 
	from sensordata a, datasets b 
	where b.reference=$1 
	and b.id = a.dataset_id
	`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get dataset from db: %v", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("no data")
	}
	return data, nil
}

func toFloatValueFn(str string) float64 {
	toFloatValue, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return toFloatValue
}
