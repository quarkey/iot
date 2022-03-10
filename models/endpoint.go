package models

import "github.com/gorilla/mux"

func (s *Server) SetupEndpoints() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/api/sensors", s.GetSensorsList).Methods("GET")                   // many sensors
	s.Router.HandleFunc("/api/sensors/{reference}", s.GetSensorByReference).Methods("GET") // one sensor by reference

	s.Router.HandleFunc("/api/sensordata/{reference}", s.GetSensorDataByReference).Methods("GET") // sensordata by reference (listing all)
	s.Router.HandleFunc("/api/sensordata", s.SaveSensorReading).Methods("POST")                   // insert new reading

	s.Router.HandleFunc("/api/datasets", s.GetDatasetsList).Methods("GET")                   // many datasets
	s.Router.HandleFunc("/api/datasets/{reference}", s.GetDatasetByReference).Methods("GET") // one dataset by reference
	s.Router.HandleFunc("/api/datasets", s.NewDataset).Methods("POST")                       // insert new dataset

	s.Router.HandleFunc("/health", s.HealthCheckHandler).Methods("GET")
}
