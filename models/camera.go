package models

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/quarkey/iot/pkg/event"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/webcam"
)

// InsertCameraEndpoint inserts a new camera into the database and returns the new camera.
func (s *Server) InsertCameraEndpoint(w http.ResponseWriter, r *http.Request) {
	c := webcam.NewCameraWithDB(s.DB)
	err := helper.DecodeBody(r, &c)
	if err != nil {
		helper.RespondErrf(w, r, 500, "unable to decode body: %v", err)
		return
	}
	newCamera, err := c.InsertNewCamera()
	if err != nil {
		helper.RespondErrf(w, r, 500, "unable insert camera: %v", err)
		return
	}
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.LogEvent(WebCamEvent, "webcam '%s' created", newCamera.Title)
	helper.Respond(w, r, 200, newCamera)
}

func (s *Server) UpdateCameraEndpoint(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c := webcam.NewCameraWithDB(s.DB)
	err := helper.DecodeBody(r, &c)
	if err != nil {
		helper.RespondErrf(w, r, 500, "unable to decode body: %v", err)
		return
	}
	updatedCamera, err := c.UpdateCamera()
	if err != nil {
		helper.RespondErrf(w, r, 500, "unable update camera: %v", err)
		return
	}
	s.Telemetry.UpdateTelemetryLists()
	e := event.New(s.DB)
	e.LogEvent(WebCamEvent, "webcam '%s' updated", updatedCamera.Title)
	helper.Respond(w, r, 200, updatedCamera)
}

func (s *Server) DeleteCameraEndpoint(w http.ResponseWriter, r *http.Request) {
	c := webcam.NewCameraWithDB(s.DB)
	err := helper.DecodeBody(r, &c)
	if err != nil {
		helper.RespondErrf(w, r, 500, "DeleteCameraEndpoint() unable to decode body: %v", err)
		return
	}
	nrows, err := c.DeleteCamera()
	if err != nil {
		helper.RespondErrf(w, r, 500, "DeleteCameraEndpoint() problems with deleting camera: %v", err)
		return
	}
	if nrows == 0 {
		helper.RespondErrf(w, r, 500, "DeleteCameraEndpoint() problems with deleting camera with id '%d', no rows affected", c.ID)
		return
	}
	helper.Respond(w, r, 200, "camera deleted")
}

func (s *Server) GetCameraListEndpoint(w http.ResponseWriter, r *http.Request) {
	cs, err := webcam.GetCameraList(s.DB)
	if err != nil {
		helper.RespondErrf(w, r, 500, "GetCameraListEndpoint(): problems with fetching list of cameras: %v", err)
		return
	}
	helper.Respond(w, r, 200, cs)
}

func (s *Server) GetCameraByIDEndpoint(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		helper.RespondErrf(w, r, 500, "GetCameraByIDEndpoint() problems %v", err)
		return
	}

	c, err := webcam.GetCameraByID(id, s.DB)
	if err != nil {
		helper.RespondErrf(w, r, 500, "GetCameraByIDEndpoint() problems with fetching camera with id '%d': %v", id, err)
		return
	}
	helper.Respond(w, r, 200, c)
}
