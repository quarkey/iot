package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/quarkey/iot/models"
)

func main() {
	srv := models.NewDB("sqlite3", "database.db")
	r := mux.NewRouter()

	r.HandleFunc("/api/sensors", srv.NewSensorReading).Methods("POST")

	// localhost:6001/health-check
	r.HandleFunc("/health-check", srv.HealthCheckHandler).Methods("GET")

	http.Handle("/", r)
	log.SetOutput(os.Stdout) // setting log output to the filehandler
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
	)(r)
	log.Fatal(http.ListenAndServe("192.168.10.119:6001", logRequest(corsHandler)))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s body: %v\n", r.RemoteAddr, r.Method, r.URL, r.Body)
		handler.ServeHTTP(w, r)
	})
}
