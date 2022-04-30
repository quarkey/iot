package models

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	helper "github.com/quarkey/iot/json"
)

type EventLogMessage struct {
	ID        int        `db:"id" json:"id"`
	Category  string     `db:"category" json:"category"`
	Message   string     `db:"message" json:"message"`
	EventTime *time.Time `db:"event_time" json:"event_time"`
}

// NewEvent stores an event message to the database
func (s *Server) NewEvent(category string, message string, v ...interface{}) {
	_, err := s.DB.Exec(`insert into events (category, message) values($1, $2)`, category, fmt.Sprintf(message, v...))
	if err != nil {
		log.Printf("[ERROR] unable to log event %v", err)
	}
}

// EventLogEndpoint ...
func (s *Server) EventLogEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	count := vars["count"]
	n, err := strconv.Atoi(count)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	var events []EventLogMessage
	err = s.DB.Select(&events, "select id, category, message, event_time from events order by id desc limit $1", n)
	if err != nil {
		helper.RespondErr(w, r, 500, err)
		return
	}
	helper.Respond(w, r, 200, events)
}
