package models

import (
	"net/http"

	helper "github.com/quarkey/iot/json"
)

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	helper.Respond(w, r, 200, "it's alive!")
}
