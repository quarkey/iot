package models

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/sensor"
)

type TemperatureReportRequest struct {
	DateFrom    string `json:"date_from"`
	DateTo      string `json:"date_to"`
	DatasetID   int    `json:"dataset_id"`
	DatasetRef  string `json:"dataset_ref"`
	DataColumn  string `json:"data_column"`
	IncludeData bool   `json:"include_data"`
}

type TemperatureReport struct {
	Description string                 `json:"description"`
	DateFrom    string                 `json:"date_from"`
	DateTo      string                 `json:"date_to"`
	HighDate    string                 `json:"high_date"`
	LowDate     string                 `json:"low_date"`
	High        float64                `json:"high"`
	Low         float64                `json:"low"`
	Average     float64                `json:"average"`
	Datapoints  int                    `json:"datapoints"`
	Data        []sensor.RawSensorData `json:"data,omitempty"`
}

func (s *Server) GetTemperatureReport(w http.ResponseWriter, r *http.Request) {
	req, err := readTempRequest(r)
	if err != nil {
		helper.RespondErr(w, r, http.StatusBadRequest, "unable to read request:", err)
		return
	}
	// we must consider the given time is local timezone, CEST or CET
	format := "2006-01-02 15:04"
	df, err := helper.ParseToLocalTime(req.DateFrom, format)
	if err != nil {
		helper.RespondErr(w, r, http.StatusInternalServerError, "unable to convert to local time: ", err)
		return
	}
	dt, err := helper.ParseToLocalTime(req.DateTo, format)
	if err != nil {
		helper.RespondErr(w, r, http.StatusInternalServerError, "unable to convert to local time: ", err)
		return
	}
	// postgres uses UTC time
	data, err := sensor.GetRawDataByDateAndRef(s.DB, df.UTC().Format(format), dt.UTC().Format(format), req.DatasetRef)
	if err != nil {
		helper.RespondErr(w, r, http.StatusInternalServerError, "unable to get data from db:", err)
		return
	}
	_, col := dataset.GetSpecificSensorDataPoint(req.DataColumn)
	var max float64
	var min float64
	var datapoints int
	var highDate string
	var lowDate string
	for i, v := range data {
		slice, err := helper.DecodeRawJSONtoSlice(v.Data)
		if err != nil {
			helper.RespondErrf(w, r, http.StatusBadRequest, "unable to parse json: %v", err)
			return
		}
		datapoint, err := strconv.ParseFloat(slice[col], 64)
		if err != nil {
			helper.RespondErrf(w, r, http.StatusBadRequest, "unable to parse datapoint: %v\n", err)
			datapoint = 0
		}
		if i == 0 {
			min = datapoint
		}
		if datapoint > max {
			max = datapoint
			highDate = v.Time.Local().String()
		}
		if datapoint < min {
			min = datapoint
			lowDate = v.Time.Local().String()
		}
		datapoints++
	}
	avg := (min + max) / 2
	out := TemperatureReport{
		Description: "High, low and average",
		DateFrom:    req.DateFrom,
		DateTo:      req.DateTo,
		High:        max,
		Low:         min,
		HighDate:    highDate,
		LowDate:     lowDate,
		Average:     avg,
		Datapoints:  datapoints,
	}
	if req.IncludeData {
		out.Data = data
	}
	helper.Respond(w, r, 200, out)
}
func readTempRequest(r *http.Request) (TemperatureReportRequest, error) {
	var dat TemperatureReportRequest
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		return TemperatureReportRequest{}, fmt.Errorf("unable to decode body: %v", err)
	}
	return dat, nil
}
