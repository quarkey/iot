package models

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/quarkey/iot/hub"
)

func (s *Server) SetupEndpoints() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/api/sensors", s.AddNewDevice).Methods("POST")
	s.Router.HandleFunc("/api/sensors", s.UpdateDevice).Methods("PUT")
	s.Router.HandleFunc("/api/sensors", s.GetSensorsListEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/sensors/{reference}", s.GetSensorByReference).Methods("GET")

	s.Router.HandleFunc("/api/sensordata/{reference}", s.GetSensorDataByReference).Methods("GET")
	s.Router.HandleFunc("/api/sensordata", s.SaveSensorReading).Methods("POST")
	s.Router.HandleFunc("/api/syncdata", s.SyncSensorData).Methods("POST")

	s.Router.HandleFunc("/api/datasets", s.GetDatasetsListEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/datasets", s.UpdateDataset).Methods("PUT")
	s.Router.HandleFunc("/api/datasets/{reference}", s.GetDatasetByReference).Methods("GET")
	s.Router.HandleFunc("/api/datasets", s.NewDataset).Methods("POST")

	// charts
	s.Router.HandleFunc("/api/chart/area/{reference}", s.AreaChartDataSeries).Methods("GET")
	s.Router.HandleFunc("/api/chart/line/{reference}", s.LineChartDataSeries).Methods("GET")

	// events
	s.Router.HandleFunc("/api/events/{count}", s.EventLogEndpoint).Methods("GET")

	s.Router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(s.Hub, w, r)
	})
	// export
	s.Router.HandleFunc("/api/exportdataset/{reference}", s.ExportSensorDataToCSVEndpoint).Methods("GET")

	s.Router.HandleFunc("/api/dashboard", s.DashboardInfoEndpoint).Methods("GET")
	s.Router.HandleFunc("/api/health", s.HealthCheckHandler).Methods("GET")

}
