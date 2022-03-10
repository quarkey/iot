package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	helper "github.com/quarkey/iot/json"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Server ....
type Server struct {
	DB         *sqlx.DB
	Config     map[string]interface{}
	Router     *mux.Router
	httpServer *http.Server
}

// New initialize server and opens a database connection.
func New(path string, automigrate bool) *Server {
	srv := &Server{}
	err := srv.loadcfg(path)
	if err != nil {
		log.Fatalf("unable to load config : %v", err)
	}
	driver := srv.Config["driver"].(string)
	connectionstr := srv.Config["connectString"].(string)
	db, err := sqlx.Open(driver, connectionstr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("unable to ping db: %v", err)
	}
	if automigrate {

		// go-migrate
		sqlFiles, err := (&file.File{}).Open("file://database/migrations")
		if err != nil {
			log.Fatalf("[ERROR] open migration files error: %v", err)
		}
		instancedriver, err := postgres.WithInstance(db.DB, &postgres.Config{})
		if err != nil {
			log.Fatalf("[ERROR] withInstance error: %v", err)
		}
		m, err := migrate.NewWithInstance("file", sqlFiles, "postgres", instancedriver)
		if err != nil {
			log.Fatal("[ERROR] NewWithInstance error:", err)
		}

		versionBefore, dirty, err := m.Version()
		if err != nil {
			log.Printf("[ERROR] unable to get database version: %v\n", err)
		}
		log.Printf("[INFO] Database version: %v, dirty: %v\n", versionBefore, dirty)

		// this will upgrade the database to latest version.
		if err := m.Up(); err != nil {
			if strings.Contains(err.Error(), `no change`) {
				log.Printf("[INFO] Database auto migration: %v\n", err)
			} else {
				log.Printf("[ERROR] migration up error: %v\n", err)
			}
		}
		versionAfter, _, err := m.Version()
		if err != nil {
			log.Fatalf("[ERROR] unable to get database version: %v", err)
		}
		if versionAfter > versionBefore {
			log.Printf("[INFO] Database auto migrated from db version '%v' to '%v'", versionBefore, versionAfter)
		}
	}
	srv.DB = db
	log.Printf("[INFO] Connected to: %s (%s)", connectionstr, db.DriverName())
	return srv
}

// Run starts the webserver
func (srv *Server) Run(ctx context.Context) {
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
	)(srv.Router)
	srv.httpServer = &http.Server{
		Addr:    srv.Config["api_addr"].(string),
		Handler: logRequest(corsHandler),
	}

	log.Printf("[INFO] Starting to listen on %s", srv.Config["api_addr"].(string))

	go func(ctx context.Context) {
		signalCh := make(chan os.Signal, 1024)
		signal.Notify(signalCh, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, os.Interrupt)
		for {
			select {
			case <-ctx.Done():
				log.Print("[INFO] Received cancel request, shutting down...")
				srv.Stop(ctx)

				return
			case sig := <-signalCh:
				log.Printf("[INFO] Received signal %v, shutting down...\n", sig)
				srv.Stop(ctx)
			}
		}
	}(ctx)
	err := srv.httpServer.ListenAndServe()
	if err != nil {
		log.Printf("[INFO] Service stopped")
	}
}

// Stop stops the webserver by shutting down context.Background
func (s *Server) Stop(ctx context.Context) {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Printf("ERROR: failed shutting down server after cancel request: %v", err)
	}
}

// loadcfg reads the contents of a jsonfile
func (s *Server) loadcfg(path string) error {
	// TODO use io.reader
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}
	if err := json.Unmarshal(data, &s.Config); err != nil {
		return fmt.Errorf("unable to unmarshal: %v", err)
	}
	return nil
}

// logRequest is a middleware that prints out all incoming requests in a nice way
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s body: %v\n", r.RemoteAddr, r.Method, r.URL, r.Body)
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	helper.Respond(w, r, 200, "it's alive!")
}
