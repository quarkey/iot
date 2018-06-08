package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	helper "github.com/quarkey/iot/json"
)

// Dataset ....
type Dataset struct {
	ID          int              `db:"id" json:"id"`
	SensorID    int              `db:"sensor_id" json:"sensor_id"`
	Title       string           `db:"title" json:"title"`
	Description string           `db:"description" json:"description"`
	Reference   string           `db:"reference" json:"reference"`
	IntervalSec int              `db:"intervalsec" json:"intervalsec"`
	Fields      *json.RawMessage `db:"fields" json:"fields"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	SensorTitle string           `json:"sensor_title,omitempty"`
}

// Datasets ....
func (s *Server) Datasets(w http.ResponseWriter, r *http.Request) {
	var datasets []Dataset
	var newdataset []Dataset

	err := s.DB.Select(&datasets, "select id, sensor_id, title, description, reference, intervalsec, fields, created_at from dataset")
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	// adding sensor title to payload, ugly i know :(
	for _, d := range datasets {
		var modified Dataset
		modified.ID = d.ID
		modified.SensorID = d.SensorID
		modified.Title = d.Title
		modified.Description = d.Description
		modified.Reference = d.Reference
		modified.IntervalSec = d.IntervalSec
		modified.Fields = d.Fields
		modified.CreatedAt = d.CreatedAt
		modified.SensorTitle = s.getArduinoTitle(d.Reference)
		fmt.Printf("looking up: %s, title: %s\n", d.Reference, d.SensorTitle)
		newdataset = append(newdataset, modified)
	}
	// fmt.Println(datasets)
	helper.Respond(w, r, 200, newdataset)
}
func (s Server) getArduinoTitle(arduinokey string) string {
	var sensortitle string
	err := s.DB.Get(&sensortitle, "select title from sensor where arduino_key=$1", arduinokey)
	if err != nil {
		return ""
	}
	return sensortitle
}

// DatasetByReference ....
func (s *Server) DatasetByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var dataset Dataset
	dataset.SensorTitle = s.getArduinoTitle(vars["reference"])
	err := s.DB.Get(&dataset, "select sensor_id, description, reference, intervalsec, fields, created_at from dataset where reference=$1", vars["reference"])
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, dataset)
}

// NewDataset ...
func (s *Server) NewDataset(w http.ResponseWriter, r *http.Request) {
	dat := Dataset{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	fmt.Println(dat.Fields)
	//TODO must be unique reference!
	_, err = s.DB.Exec("insert into dataset(sensor_id, description, reference, intervalsec, fields, created_at) values($1,$2,$3,$4,$5,$6)", dat.SensorID, dat.Description, dat.Reference, dat.IntervalSec, dat.Fields, dat.CreatedAt)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.RespondErr(w, r, http.StatusNotImplemented, "not implemented")
}
