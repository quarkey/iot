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
	Description string  `json:"description"`
	DateFrom    string  `json:"date_from"`
	DateTo      string  `json:"date_to"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	Average     float64 `json:"average"`
	Datapoints  int     `json:"datapoints"`
}

func (s *Server) GetTemperatureReport(w http.ResponseWriter, r *http.Request) {
	req, err := readTempRequest(r)
	if err != nil {
		helper.RespondErr(w, r, http.StatusBadRequest, "unable to read request:", err)
		return
	}
	data, err := sensor.GetRawDataByDateAndRef(s.DB, req.DateFrom, req.DateTo, req.DatasetRef)
	if err != nil {
		helper.RespondErr(w, r, http.StatusInternalServerError, "unable to get data from db:", err)
		return
	}
	_, col := dataset.GetSpecificSensorDataPoint(req.DataColumn)
	var max float64
	var min float64
	var datapoints int
	for i, v := range data {
		slice, err := helper.DecodeRawJSONtoSlice(v.Data)
		if err != nil {
			fmt.Printf("something went wrong... %v\n", err)
		}
		datapoint, err := strconv.ParseFloat(slice[col], 64)
		if err != nil {
			fmt.Printf("something went wrong... %v\n", err)
			datapoint = 0
		}
		if i == 0 {
			min = datapoint
		}
		if datapoint > max {
			max = datapoint
		}
		if datapoint < min {
			min = datapoint
		}
		datapoints++
	}
	avg := (min + max) / 2
	out := TemperatureReport{
		Description: "Temperature Report",
		DateFrom:    req.DateFrom,
		DateTo:      req.DateTo,
		High:        max,
		Low:         min,
		Average:     avg,
		Datapoints:  datapoints,
	}
	fmt.Println("OUT:", out)
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
