package models

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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
	Types       *json.RawMessage `db:"types" json:"types"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	SensorTitle string           `db:"sensor_title" json:"sensor_title"`
}

// GetDatasetsList fetches a list of all datasets
func (s *Server) GetDatasetsListEndpoint(w http.ResponseWriter, r *http.Request) {
	var datasets []Dataset
	err := s.DB.Select(&datasets,
		`select a.id, a.sensor_id, a.title, a.description,a.reference, 
				a.intervalsec, a.fields, a.types, a.created_at,
				b.title as sensor_title
		from datasets a, sensors b
		where a.sensor_id = b.id`)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, datasets)
}
func GetDatasetsList(db *sqlx.DB) []Dataset {
	var datasets []Dataset
	err := db.Select(&datasets,
		`select a.id, a.sensor_id, a.title, a.description,a.reference, 
				a.intervalsec, a.fields, a.types, a.created_at,
				b.title as sensor_title
		from datasets a, sensors b
		where a.sensor_id = b.id`)
	if err != nil {
		return nil
	}
	return datasets
}

// GetDatasetByReference fetches a dataset based on a reference
func (s *Server) GetDatasetByReference(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var dataset Dataset
	err := s.DB.Get(&dataset, `
	select a.id, a.sensor_id, a.title, a.description, 
	a.reference, a.intervalsec, a.fields, a.types, 
	a.created_at, b.title as sensor_title
	 from datasets a, sensors b
		where reference=$1
		and a.sensor_id = b.id
		`, vars["reference"])
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
	//TODO must be unique reference!
	_, err = s.DB.Exec(`insert into datasets(sensor_id, title, description, reference, intervalsec, fields, types, created_at) 
	values($1,$2,$3,$4,$5,$6,$7,$8)`, dat.SensorID, dat.Title, dat.Description, dat.Reference, dat.IntervalSec, dat.Fields, dat.Types, dat.CreatedAt)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.RespondSuccess(w, r)
}
