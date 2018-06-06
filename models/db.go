package models

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Server struct {
	DB *sqlx.DB
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
	return &Server{DB: db}
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// e.g. http.HandleFunc("/health-check", HealthCheckHandler)
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	var dbstatus string
	if err := s.DB.Ping(); err != nil {
		dbstatus = fmt.Sprintf("ping error: %v", err)
		log.Printf(dbstatus)
	} else {
		dbstatus = "alive"
	}
	fmt.Fprintf(w, `{"api": "alive", "db": "%s"}`, dbstatus)
}
