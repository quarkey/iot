package models

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/quarkey/iot/pkg/hub"
)

func (s *Server) SetupEndpoints() {
	s.Router = chi.NewRouter()
	// sensors
	s.Router.Post("/api/sensors", s.AddNewDevice)
	s.Router.Put("/api/sensors", s.UpdateDevice)
	s.Router.Get("/api/sensors", s.GetSensorsListEndpoint)
	s.Router.Get("/api/sensors/{reference}", s.GetSensorByReference)

	// sensordata
	s.Router.Get("/api/sensordata/{reference}/{limit}", s.GetSensorDataByReferenceEndpoint)
	s.Router.Post("/api/sensordata", s.SaveSensorReading)
	s.Router.Post("/api/syncdata", s.SyncSensorData)

	// datasets
	s.Router.Get("/api/datasets", s.GetDatasetsListEndpoint)
	s.Router.Put("/api/datasets", s.UpdateDataset)
	s.Router.Get("/api/datasets/{reference}", s.GetDatasetByReference)
	s.Router.Post("/api/datasets", s.NewDataset)
	s.Router.Post("/api/datasets/delete", s.DeleteDatasetByIDEndpoint)

	// charts
	s.Router.Get("/api/chart/area/{reference}/{limit}", s.AreaChartDataSeries)
	s.Router.Get("/api/chart/line/{reference}/{limit}", s.LineChartDataSeries)

	// reports
	s.Router.Post("/api/report/temperature", s.GetTemperatureReport)

	// events
	s.Router.Get("/api/events/{count}", s.GetEventLogListEndpoint)

	// controllers
	s.Router.Get("/api/controllers", s.GetControllersListEndpoint)
	s.Router.Get("/api/controllers/{cid}", s.GetControllerByIDEndpoint)
	s.Router.Post("/api/controllers", s.AddNewControllerEndpoint)
	s.Router.Put("/api/controllers", s.UpdateControllerByIDEndpoint)

	s.Router.Get("/api/controller/{id}/state/{state}", s.SetControllerStateEndpoint)
	s.Router.Get("/api/controller/{id}/switch/{state}", s.SetControllerSwitchStateEndpoint)
	s.Router.Get("/api/controller/{id}/alert/{state}", s.SetControllerAlertStateEndpoint)

	s.Router.Post("/api/controller/reset", s.ResetControllerSwitchValueEndpoint)
	s.Router.Post("/api/controller/delete", s.DeleteControllerByIDEndpoint)

	// socket upgrader for live dataset monitoring
	s.Router.HandleFunc("/api/live", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(s.Hub, w, r)
	})
	// export
	s.Router.Get("/api/exportdataset/{reference}", s.ExportSensorDataToCSVEndpoint)

	s.Router.Get("/api/dashboard", s.DashboardInfoEndpoint)
	s.Router.Get("/api/health", s.HealthCheckHandler)
	s.Router.Get("/api/webhook/test", s.TestCheckWebhooks)
}
