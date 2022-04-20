package main

import (
	"context"
	"flag"
	"log"

	_ "github.com/lib/pq" // postgres driver
	"github.com/quarkey/iot/models"
)

func main() {
	confPath := flag.String("conf", "", "path to configuration file")
	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()
	if *confPath == "" {
		log.Fatalf("ERROR: missing configuration jsonfile")
	}
	server := models.New(*confPath, true, *debug) //automigration=true
	server.SetupEndpoints()
	server.Run(context.Background())
}
