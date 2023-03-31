package models

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quarkey/iot/pkg/hub"
)

func (s *Server) SetupEndpoints() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/api/sensors", s.AddNewDevice).Methods("POST")
	s.Router.HandleFunc("/api/sensors", s.UpdateDevice).Methods("PUT")
	s.Router.HandleFunc("/api/sensors", s.GetSensorsListEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/sensors/{reference}", s.GetSensorByReference).Methods("GET")

	s.Router.HandleFunc("/api/sensordata/{reference}", s.GetSensorDataByReferenceEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/sensordata", s.SaveSensorReading).Methods("POST")
	s.Router.HandleFunc("/api/syncdata", s.SyncSensorData).Methods("POST")

	// datasets
	s.Router.HandleFunc("/api/datasets", s.GetDatasetsListEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/datasets", s.UpdateDataset).Methods("PUT")
	s.Router.HandleFunc("/api/datasets/{reference}", s.GetDatasetByReference).Methods("GET")
	s.Router.HandleFunc("/api/datasets", s.NewDataset).Methods("POST")
	s.Router.HandleFunc("/api/datasets/delete", s.DeleteDatasetByIDEndpoint).Methods("POST")

	// charts
	s.Router.HandleFunc("/api/chart/area/{reference}", s.AreaChartDataSeries).Methods("GET")
	s.Router.HandleFunc("/api/chart/line/{reference}", s.LineChartDataSeries).Methods("GET")

	// events
	s.Router.HandleFunc("/api/events/{count}", s.EventLogEndpoint).Methods("GET")

	// controllers
	s.Router.HandleFunc("/api/controllers", s.GetControllersListEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/controllers/{cid}", s.GetControllerByIDEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/controllers", s.AddNewControllerEndpoint).Methods("POST")
	s.Router.HandleFunc("/api/controllers", s.UpdateControllerByIDEndpoint).Methods("PUT")

	s.Router.HandleFunc("/api/controller/{id}/state/{state}", s.SetControllerStateEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/controller/{id}/switch/{state}", s.SetControllerSwitchStateEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/controller/{id}/alert/{state}", s.SetControllerAlertStateEndpoint).Methods("GET")

	s.Router.HandleFunc("/api/controller/reset", s.ResetControllerSwitchValueEndpoint).Methods("POST")
	s.Router.HandleFunc("/api/controller/delete", s.DeleteControllerByIDEndpoint).Methods("POST")

	// socket upgrader for live dataset monitoring
	s.Router.HandleFunc("/api/live", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(s.Hub, w, r)
	})
	// export
	s.Router.HandleFunc("/api/exportdataset/{reference}", s.ExportSensorDataToCSVEndpoint).Methods("GET")

	s.Router.HandleFunc("/api/dashboard", s.DashboardInfoEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/health", s.HealthCheckHandler).Methods("GET")
}
