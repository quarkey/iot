package models

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/quarkey/iot/pkg/chart"
	"github.com/quarkey/iot/pkg/helper"
)

// LineChartDataSeries will generate a data structure that is fitted to ng2-charts.
func (s *Server) LineChartDataSeries(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		helper.RespondErrf(w, r, http.StatusBadRequest, "unable to parse string to int: %v", err)
		return
	}
	series, err := chart.LineDataSeries(s.DB, chi.URLParam(r, "reference"), limit)
	if err != nil {
		if err.Error() == "no data" {
			helper.RespondErr(w, r, http.StatusBadRequest, err)
			return
		}
		helper.RespondErr(w, r, 500, "Problems with loading line chart series: ", err)
		return
	}
	helper.Respond(w, r, 200, series)
}

// AreaChartDataSeries will generate a data structure that is fitted to ngx-charts.
func (s *Server) AreaChartDataSeries(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(chi.URLParam(r, "limit"))
	if err != nil {
		helper.RespondErrf(w, r, http.StatusBadRequest, "unable to parse string to int: %v", err)
		return
	}
	series, err := chart.AreaChartDataSeries(s.DB, chi.URLParam(r, "reference"), limit)
	if err != nil {
		if err.Error() == "no data" {
			helper.RespondErr(w, r, http.StatusBadRequest, err)
			return
		}
		helper.RespondErrf(w, r, 500, "Problems with loading area chart series: ", err)
		return
	}
	helper.Respond(w, r, 200, series)
}
