package models

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	helper "github.com/quarkey/iot/json"
)

// Dataset ....
type Dataset struct {
	ID          int              `db:"id" json:"id"`
	SensorID    int              `db:"sensor_id" json:"sensor_id"`
	Description string           `db:"description" json:"description"`
	Reference   string           `db:"reference" json:"reference"`
	IntervalSec int              `db:"intervalsec" json:"intervalsec"`
	Fields      *json.RawMessage `db:"fields" json:"fields"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
}

// Datasets ....
func (s *Server) Datasets(w http.ResponseWriter, r *http.Request) {
	var datasets []Dataset
	err := s.DB.Select(&datasets, "select id, sensor_id, description, reference, intervalsec, fields, created_at from dataset")
	if err != nil {
		log.Printf("unable to run query: %v", err)
		helper.RespondHTTPErr(w, r, 500)
		return
	}
	helper.Respond(w, r, 200, datasets)

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
