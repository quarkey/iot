package models

import "github.com/gorilla/mux"

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

	s.Router.HandleFunc("/api/health", s.HealthCheckHandler).Methods("GET")

}
