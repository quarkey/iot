package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	Showcharts  *json.RawMessage `db:"showcharts" json:"showcharts"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	SensorTitle string           `db:"sensor_title" json:"sensor_title"`
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
	dataset, err := s.getDsetByRef(vars["reference"])
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
		b.title as sensor_title
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
	helper.Respond(w, r, 200, dataset)
}

func DatasetFieldAndShowCartList(ref string, db *sqlx.DB) ([]string, []bool, error) {
	fields, err := datasetFieldsList(ref, db)
	if err != nil {
		return nil, nil, err
	}
	showCharts, err := datasetShowChartBools(ref, db)
	if err != nil {
		return nil, nil, err
	}
	return fields, showCharts, nil
}

// DatasetFieldsList returns a list of column labels for a dataset by given reference
func datasetFieldsList(ref string, db *sqlx.DB) ([]string, error) {
	var raw *json.RawMessage
	err := db.Get(&raw, `select fields from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get labels from datasets: %v", err)
	}
	var fields []string
	err = json.Unmarshal(*raw, &fields)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v", err)
	}
	return fields, nil
}

// datasetShowChartBools returns a list of boolean values to indicate if chart should show or not.
// Method takes in account that showcharts fields from pg may be empty.
// Also i think there is a bug in this function, second marshal fails because
// pg res can be validated as bool [true, "true"]
func datasetShowChartBools(ref string, db *sqlx.DB) ([]bool, error) {
	// fetching fields to determine field count
	var raw *json.RawMessage
	err := db.Get(&raw, `select fields from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get fields from datasets: %v", err)
	}
	var fields []string
	err = json.Unmarshal(*raw, &fields)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %v", err)
	}
	var nextRaw *json.RawMessage
	err = db.Get(&nextRaw, `select showcharts from datasets where reference=$1`, ref)
	if err != nil {
		return nil, fmt.Errorf("unable to get showchart bools from dataset: %v", err)
	}
	var nextOut []bool
	// defaulting to false for all datasets
	// if boolean values is missing in the database
	if nextRaw == nil {
		if len(nextOut) == 0 {
			var tmp []bool
			for i := 0; i < len(fields); i++ {
				tmp = append(tmp, false)
			}
			nextOut = tmp
		}
		return nextOut, nil
	}
	err = json.Unmarshal(*nextRaw, &nextOut)
	if err != nil {
		if err.Error() == "json: cannot unmarshal string into Go value of type bool" {
			// we can handle them as bool if they are stored as strings in the database
			var str []string
			var tmp []bool
			err := json.Unmarshal(*nextRaw, &str)
			if err != nil {
				return nil, fmt.Errorf("DatasetShowChartBools: %v ", err)
			}
			for _, v := range str {
				boolValue, err := strconv.ParseBool(v)
				if err != nil {
					return nil, fmt.Errorf("DatasetShowChartBools: %v", err)
				}
				tmp = append(tmp, boolValue)
			}
			nextOut = tmp
		}
	}
	return nextOut, nil
}
