package models

import (
	"net/http"

	helper "github.com/quarkey/iot/json"
)

type dashinfo struct {
	Sensor      int `json:"sensors"`
	Datasets    int `json:"datasets"`
	Controllers int `json:"controllers"`
}

func (s *Server) DashboardInfoEndpoint(w http.ResponseWriter, r *http.Request) {
	var sensors, datasets int
	err := s.DB.Get(&sensors, `select count(*) as "sensors" from sensors`)
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	err = s.DB.Get(&datasets, `select count(*) as "datasets" from datasets`)
	if err != nil {
		helper.RespondErr(w, r, 400, err)
		return
	}
	info := dashinfo{
		Sensor:      sensors,
		Datasets:    datasets,
		Controllers: 0,
	}
	helper.Respond(w, r, 200, info)
}
