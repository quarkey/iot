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
func NewDB(driver, source string) *Server {
	db, err := sqlx.Open(driver, source)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to ping db: %v", err)
	}
	log.Printf("Connected to: %s (%s)", source, db.DriverName())
	s := &Server{DB: db}
	s.loadcfg()
	return s
}

// loadcfg
func (s *Server) loadcfg() error {
	// TODO use io.reader
	defaultPath := "./exampleconfig.json"
	data, err := ioutil.ReadFile(defaultPath)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}
	if err := json.Unmarshal(data, &s.Config); err != nil {
		return fmt.Errorf("unable to unmarshal: %v", err)
	}
	fmt.Println(s.Config)
	return nil
}
