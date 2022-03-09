package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Server ....
type Server struct {
	DB     *sqlx.DB
	Config map[string]interface{}
}

// NewDB opens connection to database
func NewDB(path string, automigrate bool) *Server {
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
	if automigrate {

		// go-migrate
		sqlFiles, err := (&file.File{}).Open("file://database/migrations")
		if err != nil {
			log.Fatalf("[ERROR] open migration files error: %v", err)
		}
		instancedriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			log.Fatalf("[ERROR] withInstance error: %v", err)
		}
		m, err := migrate.NewWithInstance("file", sqlFiles, "postgres", instancedriver)
		if err != nil {
			log.Fatal("[ERROR] NewWithInstance error:", err)
		}

		versionBefore, dirty, err := m.Version()
		if err != nil {
			log.Printf("[ERROR] unable to get database version: %v\n", err)
		}
		log.Printf("[INFO] Database version: %v, dirty: %v\n", versionBefore, dirty)

		// this will upgrade the database to latest version.
		if err := m.Up(); err != nil {
			if strings.Contains(err.Error(), `no change`) {
				log.Printf("[INFO] Database auto migration: %v\n", err)
			} else {
				log.Printf("[ERROR] migration up error: %v\n", err)
			}
		}
		versionAfter, _, err := m.Version()
		if err != nil {
			log.Fatalf("[ERROR] unable to get database version: %v", err)
		}
		if versionAfter > versionBefore {
			log.Printf("[INFO] Database auto migrated from db version '%v' to '%v'", versionBefore, versionAfter)
		}
	}
	srv.DB = db
	log.Printf("[INFO] Connected to: %s (%s)", connectionstr, db.DriverName())
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
