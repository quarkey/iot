package event

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Event struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Event {
	return Event{
		db: db,
	}
}

// NewEvent stores an event message to the database
func (e *Event) LogEvent(category string, message string, v ...interface{}) {
	_, err := e.db.Exec(`insert into events (category, message) values($1, $2)`, category, fmt.Sprintf(message, v...))
	if err != nil {
		log.Printf("[ERROR] unable to log event %v", err)
	}
}

type EventLogMessage struct {
	ID        int        `db:"id" json:"id"`
	Category  string     `db:"category" json:"category"`
	Message   string     `db:"message" json:"message"`
	EventTime *time.Time `db:"event_time" json:"event_time"`
}

// GetEventLogWithLimit will fetch a list of event messages with a count limit
func (e *Event) GetEventLogWithLimit(limit int) ([]EventLogMessage, error) {
	var list []EventLogMessage
	err := e.db.Select(&list, "select id, category, message, event_time from events order by id desc limit $1", limit)
	if err != nil {
		return list, err
	}
	return list, nil
}

// GetEventLogWithLimitRequest parses a GET request that contains the variable "count" and
// fetching a list of event messages with that given count.
func (e *Event) GetEventLogWithLimitRequest(r *http.Request) ([]EventLogMessage, error) {
	n, err := strconv.Atoi(chi.URLParam(r, "count"))
	if err != nil {
		return nil, err
	}
	list, err := e.GetEventLogWithLimit(n)
	if err != nil {
		return list, fmt.Errorf("unable to fetch events from database '%v'", err)
	}
	return list, nil
}
