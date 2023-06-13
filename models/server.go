package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver
	"github.com/quarkey/iot/pkg/dataset"
	"github.com/quarkey/iot/pkg/event"
	"github.com/quarkey/iot/pkg/helper"
	"github.com/quarkey/iot/pkg/hub"
	"github.com/quarkey/iot/pkg/webhooks"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// system event such as start, stop etc
var SystemEvent = "system"
var SernsorEvent = "sensor"
var DatasetEvent = "dataset"
var ControllerEvent = "controller"

var TimeFormat = "2006-01-02 15:04:05"

// Server ....
type Server struct {
	DB *sqlx.DB
	// Server JSON config file
	Config     map[string]interface{}
	Router     *chi.Mux
	httpServer *http.Server
	//Hub to keep track of socket connection
	// for live monitoring of datasets.
	Hub       *hub.Hub
	Telemetry *Telemetry
	Debug     bool
	simulator *Sim
	startTime time.Time
}

var GLOBALCONFIG map[string]interface{}

// //go:embed ../database/migrations/sql/*.sql
// var SQLfs embed.FS

// New initialize server and opens a database connection.
func New(path string, automigrate bool, debug bool) *Server {
	//TODO move timezone to config
	SetTimeZone("Europe/Oslo")
	srv := &Server{}
	srv.startTime = time.Now()
	log.Printf("[INFO] Loading config: %v", path)
	err := srv.loadcfg(path)
	if err != nil {
		log.Fatalf("unable to load config : %v", err)
	}
	if debug {
		srv.Debug = true
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
		sqlFiles, err := (&file.File{}).Open(srv.Config["migration"].(string))
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
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
	)(srv.Router)
	srv.httpServer = &http.Server{
		Addr:    srv.Config["api_addr"].(string),
		Handler: logRequest(corsHandler),
	}

	log.Printf("[INFO] Starting to listen on %s", srv.Config["api_addr"].(string))

	// resetting dataset connectivity on start to offline
	err := dataset.ResetConnectivity(srv.DB)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	e := event.New(srv.DB)
	e.LogEvent(SystemEvent, "Server started")
	// socket hub for live monitoring
	hub := hub.NewHub()
	srv.Hub = hub
	go srv.Hub.Run()
	// server ticker timer for scheduled tasks
	srv.Telemetry = newTelemetryTicker(srv.DB)
	srv.Telemetry.startTelemetryTicker(srv.Config, srv.Debug)

	// only when allowSim is set we start the simulator.
	state, exist := srv.Config["allowSim"]
	if exist {
		log.Printf("[INFO] allowSim is set to %v\n", state)
		if state == true {
			sim := NewSim()
			srv.simulator = sim
			srv.StartSim(sim)
		}
	}

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
				return
			}
		}
	}(ctx)
	err = srv.httpServer.ListenAndServe()
	if err != nil {
		log.Printf("[INFO] Service stopped")
		e.LogEvent(SystemEvent, "Server stopped")
	}
}

// Stop stops the webserver by shutting down context.Background
func (s *Server) Stop(ctx context.Context) {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Printf("[ERROR]: failed shutting down server after cancel request: %v", err)
		panic("PANIC!!!")
	}
}

// loadcfg reads the contents of a jsonfile
func (s *Server) loadcfg(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}
	if err := json.Unmarshal(data, &s.Config); err != nil {
		return fmt.Errorf("unable to unmarshal: %v", err)
	}
	GLOBALCONFIG = s.Config
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
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	out := map[string]interface{}{
		"start_time":         s.startTime.Local().Format(TimeFormat),
		"uptime":             time.Since(s.startTime).String(),
		"memory_alloc":       helper.BytesToHuman(int64(m.Alloc)),
		"memory_tot_alloc":   helper.BytesToHuman(int64(m.TotalAlloc)),
		"system_mem":         helper.BytesToHuman(int64(m.Sys)),
		"garbage_collection": m.NumGC,
		"pg_table_stats":     s.pgTableStats(),
	}
	helper.Respond(w, r, 200, out)
}
func (s *Server) TestCheckWebhooks(w http.ResponseWriter, r *http.Request) {
	wh, err := webhooks.ParseConfig(GLOBALCONFIG["discordConfig"].(string))
	if err != nil {
		log.Printf("[ERROR] unable to parse discord webhook configuration: %v", err)
	}
	wh.Discord.Sendf("testing connection ...")
}
func (s *Server) API_URL() string {
	return fmt.Sprintf("http://%s/api", s.Config["api_addr"].(string))
}
func (s *Server) SERVER_URL() string {
	return fmt.Sprintf("http://%s", s.Config["api_addr"].(string))
}

type pgStats struct {
	Relname    string `db:"relname" json:"relname"`
	Full_size  string `db:"full_size" json:"full_size"`
	Table_size string `db:"table_size" json:"table_size"`
	Index_size string `db:"index_size" json:"index_size"`
}

func (s *Server) pgTableStats() []pgStats {
	var pgStats []pgStats
	err := s.DB.Select(&pgStats, `select relname, pg_size_pretty(pg_total_relation_size(relname::regclass)) as full_size, pg_size_pretty(pg_relation_size(relname::regclass)) as table_size, pg_size_pretty(pg_total_relation_size(relname::regclass) - pg_relation_size(relname::regclass)) as index_size from pg_stat_user_tables order by pg_total_relation_size(relname::regclass) desc limit 10;`)
	if err != nil {
		log.Printf("[ERROR] unable to get pg stats: %v", err)
		return nil
	}
	return pgStats
}

func SetTimeZone(tz string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("[ERROR] unable to set timezone to '%s': %v", tz, err)
		return
	}
	time.Local = loc
	log.Printf("[INFO] Timezone set to '%s'", loc.String())
}
