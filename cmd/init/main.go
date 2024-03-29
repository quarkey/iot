package main

import (
	"flag"
	"log"

	_ "github.com/lib/pq" // postgres driver
	models "github.com/quarkey/iot/models"
)

func main() {
	confPath := flag.String("conf", "", "path to your config")
	automigrate := flag.Bool("automigrate", false, "allow program to run postgres automigration")
	debug := flag.Bool("debug", false, "enable debug mode")

	flag.Parse()

	if *confPath == "" {
		log.Fatalf("ERROR: missing configuration jsonfile")
	}
	server := models.New(*confPath, *automigrate, *debug)

	err := server.InsertTestdata()
	if err != nil {
		log.Fatalf("unable to insert test data: %v", err)
	}
}
