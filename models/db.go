package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
)

// Server ....
type Server struct {
	DB     *sqlx.DB
	Config map[string]interface{}
}

// NewDB opens connection to database
func NewDB(path string) *Server {
	srv := &Server{}
	err := srv.loadcfg(path)
	if err != nil {
		log.Fatalf("unable to load config : %v", err)
	}
	driver := srv.Config["driver"].(string)
	connectionstr := srv.Config["connectString"].(string)

	db, err := sqlx.Open(driver, connectionstr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to ping db: %v", err)
	}
	srv.DB = db
	log.Printf("Connected to: %s (%s)", connectionstr, db.DriverName())
	return srv
}

// loadcfg
func (s *Server) loadcfg(path string) error {
	// TODO use io.reader
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}
	if err := json.Unmarshal(data, &s.Config); err != nil {
		return fmt.Errorf("unable to unmarshal: %v", err)
	}
	return nil
}
