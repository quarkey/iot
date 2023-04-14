package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/event"
	"github.com/quarkey/iot/pkg/helper"
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
	Showcharts  *json.RawMessage `db:"showcharts" json:"showcharts"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	SensorTitle string           `db:"sensor_title" json:"sensor_title"`
	Telemetry   string           `db:"telemetry" json:"telemetry"`
}

// GetDatasetsList fetches a list of all datasets
func (s *Server) GetDatasetsListEndpoint(w http.ResponseWriter, r *http.Request) {
	var datasets []Dataset
	err := s.DB.Select(&datasets,
		`select 
			a.id,
			a.sensor_id,
			a.title,
			a.description,
			a.reference, 
			a.intervalsec,
			a.fields,
			a.types,
			a.showcharts,
			a.created_at,
			b.title as sensor_title,
			a.telemetry
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
	// TODO: add error handling
	var datasets []Dataset
	err := db.Select(&datasets,
		`select 
			a.id,
			a.sensor_id,
			a.title,
			a.description,
			a.reference,
			a.intervalsec,
			a.fields,
			a.types,
			a.showcharts,
			a.created_at,
			b.title as sensor_title,
			a.telemetry
		from datasets a, sensors b
		where a.sensor_id = b.id`)
	if err != nil {
		return nil
	}
	return datasets
}

// GetDatasetByReference fetches a dataset based on a reference
func (s *Server) GetDatasetByReference(w http.ResponseWriter, r *http.Request) {
	dataset, err := s.getDsetByRef(chi.URLParam(r, "reference"))
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, dataset)
}

func (s Server) getDsetByRef(ref string) (Dataset, error) {
	var dataset Dataset
	err := s.DB.Get(&dataset, `
	select 
		a.id, 
		a.sensor_id,
		a.title,
		a.description, 
		a.reference,
		a.intervalsec,
		a.fields,
		a.types,
		a.showcharts,
		a.created_at,
		b.title as sensor_title,
		a.telemetry
	from datasets a, sensors b
	where reference=$1
	and a.sensor_id = b.id
		`, ref)
	if err != nil {
		return dataset, err
	}
	return dataset, err
}

// NewDataset adds a new dataset to db
func (s *Server) NewDataset(w http.ResponseWriter, r *http.Request) {
	dat := Dataset{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}

	if dat.Reference == "" {
		helper.RespondErr(w, r, http.StatusBadRequest, "missing device reference")
		return
	}
	fmt.Println(dat)
	_, err = s.DB.Exec(`insert into datasets(sensor_id, title, description, reference, intervalsec, fields, types) 
	values($1,$2,$3,$4,$5,$6,$7)`,
		dat.SensorID,
		dat.Title,
		dat.Description,
		dat.Reference,
		dat.IntervalSec,
		dat.Fields,
		dat.Types)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	// also update telemetry dataset list
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.LogEvent(DatasetEvent, "dataset '%s' added", dat.Title)
	helper.RespondSuccess(w, r)
}

func (s *Server) UpdateDataset(w http.ResponseWriter, r *http.Request) {
	dat := Dataset{}
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	_, err = s.DB.Exec(`update iot.datasets set 
		title=$1, 
		description=$2,
		fields=$3,
		types=$4,
		intervalsec=$5,
		showcharts=$6
	where reference=$7
	and id=$8
	`, dat.Title, dat.Description, dat.Fields, dat.Types, dat.IntervalSec, dat.Showcharts, dat.Reference, dat.ID)
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondErr(w, r, 500, err)
		return
	}
	dataset, err := s.getDsetByRef(dat.Reference)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	// also update telemetry dataset list
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.LogEvent(DatasetEvent, "dataset '%s' updated", dataset.Title)
	helper.Respond(w, r, 200, dataset)
}

// DeleteDatasetByIDEndpoint deletes dataset and any sensor data associated with it.
// WARNING!!! THIS IS IRREVERSIBLE!
func (s *Server) DeleteDatasetByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	var dat Dataset
	err := helper.DecodeBody(r, &dat)
	if err != nil {
		log.Printf("unable to decode body: %v", err)
		helper.RespondErr(w, r, 500, err)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to start transactio: ", err)
	}
	// deleting sensor data first so forein key constraint wont get triggered
	_, err = tx.Exec(`delete from sensordata where dataset_id=$1`, dat.ID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			helper.RespondErr(w, r, 500, "unable to delete rolling back: ", err)
			return
		}
		helper.RespondErr(w, r, 500, "unable to delete dataset: ", err)
		return
	}
	// deleting dataset
	_, err = tx.Exec(`delete from datasets where id=$1`, dat.ID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			helper.RespondErr(w, r, 500, "unable to delete rolling back: ", err)
			return
		}
		helper.RespondErr(w, r, 500, "unable to delete dataset: ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to commit: ", err)
		return
	}
	// also update telemetry
	// also update telemetry dataset list
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.LogEvent(DatasetEvent, "dataset '%s' deleted", dat.Title)

	helper.RespondSuccess(w, r, 200)
}
