package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/quarkey/iot/models"
)

func main() {
	confPath := flag.String("conf", "", "path to your config")
	automigrate := flag.Bool("automigrate", false, "allow program to run postgres automigration")

	flag.Parse()

	if *confPath == "" {
		log.Fatalf("ERROR: missing configuration jsonfile")
	}
	server := models.New(*confPath, *automigrate)
	datasets := models.GetDatasetsList(server.DB)

	_ = models.New(*confPath, *automigrate)
	var jobs []chan string
	for i, ds := range datasets {
		log.Printf("[INFO] %d/%d simulating sensor for dataset '%s, interval duration: %d'\n", i+1, len(datasets), ds.Title, ds.IntervalSec)
		job := make(chan string)
		jobs = append(jobs, job)
		go runSim(&ds)
	}
	for _, result := range jobs {
		fmt.Println("kanal:", <-result)
	}
}
func runSim(ds *models.Dataset) {
	tick := 0
	url := "http://localhost:6001/api/sensordata"
	data := []byte(fmt.Sprintf(`{"sensor_id": %d,"dataset_id": %d,"data": [123.00,12.00]}`, ds.SensorID, ds.ID))
	for {
		tick++
		log.Printf("[RESULT] tick: %d %s", tick, ds.Title)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Failed creating request: %v", err)
			break
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Failed performing request: %v", err)
			break
		}
		res.Body.Close()
		time.Sleep(time.Duration(ds.IntervalSec) * time.Second)
	}
	fmt.Println("failed")
}

// takes an api url and changing :port number to n+10001
func makeURLwPort(url string, index int) string {
	r := regexp.MustCompile(`\d+`)
	port := r.FindString(url)
	nport, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("unable to convert to int: %v", err)
	}
	newport := strconv.Itoa(nport + 10000 + index)
	newurl := strings.Replace(url, port, newport, 1)
	return newurl
}
