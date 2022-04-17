package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	helper "github.com/quarkey/iot/json"
)

// AreaChartDataSeries will generate a data structure that is fitted to ngx-charts.
func (s *Server) AreaChartDataSeries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data []Data
	err := s.DB.Select(&data, `select a.id, a.data, a.time from 
	sensordata a, datasets b where b.reference=$1 and b.id = a.dataset_id`, vars["reference"])
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to get dataset from db:", err)
		return
	}
	if len(data) == 0 {
		helper.RespondErr(w, r, 400, "no data")
		return
	}
	// decoding jsonRawMessage data column
	raw, err := decodeRaw(data[0].Data)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	// fetching field list
	fields, err := s.DatasetFieldsList(vars["reference"])
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	// showcharts used to determine if chart should added or not.
	showcharts, err := s.DatasetShowChartBools(vars["reference"])
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	// populate data structure fitted for ngx-charts
	var out []series
	for i := 0; i < len(raw); i++ {
		// skipping disabled charts
		if !showcharts[i] {
			continue
		}
		var ps []point
		for _, set := range data {
			decoded, err := decodeRaw(set.Data)
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
			ps = append(ps, point{Name: set.Time.String(), Value: toFloatValue})
		}
		out = append(out, series{Name: fields[i], Point: ps})
	}
	helper.Respond(w, r, 200, out)
}
func decodeRaw(dat *json.RawMessage) ([]string, error) {
	var sets []string
	err := json.Unmarshal(*dat, &sets)
	if err != nil {
		return nil, fmt.Errorf("unable to decode: %v", err)
	}
	return sets, nil
}

type series struct {
	Name  string  `json:"name"`
	Point []point `json:"series"`
}
type point struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
