package models

import (
	"net/http"

	"github.com/quarkey/iot/pkg/event"
	"github.com/quarkey/iot/pkg/helper"
)

// EventLogEndpoint fetches a list of events with a limite from GET param count
func (s *Server) EventLogEndpoint(w http.ResponseWriter, r *http.Request) {
	e := event.New(s.DB)
	events, err := e.GetEventLogWithLimitRequest(r)
	if err != nil {
		helper.RespondErr(w, r, 500, "unable to fetch events from database: ", err)
		return
	}
	helper.Respond(w, r, 200, events)
}
