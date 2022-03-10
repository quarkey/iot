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
	flag.Parse()
	if *confPath == "" {
		log.Fatalf("ERROR: missing configuration jsonfile")
	}
	srv := models.NewDB(*confPath, true) //automigration=true
	srv.SetupEndpoints()
	srv.Run(context.Background())
}
