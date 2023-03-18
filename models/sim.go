package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	helper "github.com/quarkey/iot/json"
)

type Sim struct {
	allowSim   bool
	simRunning bool
	done       chan bool
	// map of dset ids and interval second remaining
	items map[int]int
}

func NewSim() *Sim {
	return &Sim{
		allowSim:   true,
		simRunning: false,
		done:       make(chan bool),
		items:      make(map[int]int),
	}
}

// Start starts a ticker with 1 second duration
func (s *Server) StartSim(sim *Sim) {
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	log.Printf("[INFO] Simulator starting in 4 seconds...")
	s.simulator.simRunning = true
	go func() {
		for i := 1; i < 3; i++ {
			log.Printf("[INFO] %d", i)
			time.Sleep(1 * time.Second)
		}
		log.Printf("[INFO] Simulator started")

		for {
			select {
			case <-s.simulator.done:
				return
			case <-ticker.C:
				for _, dset := range s.Telemetry.datasets {
					// skipping if fields (data columns) are empty
					if dset.Fields == nil {
						break
					}
					item, ok := sim.items[dset.ID]
					if ok {
						if item == 0 {
							// fmt.Printf("signal '%s' fired, '%ds' to next!\n", dset.Title, dset.IntervalSec)
							sim.items[dset.ID] = dset.IntervalSec
							doSim(dset)
						}
						seconds := sim.items[dset.ID] - 1
						sim.items[dset.ID] = seconds
					}
					if !ok {
						// fmt.Printf("new sim item '%s' active with interval '%ds'\n", dset.Title, dset.IntervalSec)
						sim.items[dset.ID] = dset.IntervalSec
					}
				}
			}
		}
	}()
}

func (s *Server) StopSim() {
	s.simulator.done <- true
}

var sensorURL = "http://localhost:6001/api/sensordata"

func doSim(ds Dataset) {
	var dataPoints []string
	types, err := helper.DecodeRawJSONtoSlice(ds.Types)
	if err != nil {
		log.Printf("[ERROR] sim error: %v", err)
		return
	}
	if len(types) > 0 {
		for _, field := range types {
			switch field {
			case "int":
				// fmt.Println("int")
				rand.Seed(time.Now().UnixNano())
				num := 20 + rand.Int()*(60-20)
				dataPoints = append(dataPoints, fmt.Sprintf("%d", num))
			case "float", "float32", "float64":
				// fmt.Println("float")
				rand.Seed(time.Now().UnixNano())
				num := 20 + rand.Float64()*(60-20)
				dataPoints = append(dataPoints, fmt.Sprintf("%f", num))
			case "string":
				// fmt.Println("string")
				dataPoints = append(dataPoints, "string")
			default:
				// fmt.Println("type not supported:", field)
			}
		}
	}
	dataOut := struct {
		Sensor_id  int      `json:"sensor_id"`
		Dataset_id int      `json:"dataset_id"`
		Data       []string `json:"data"`
	}{
		Sensor_id:  ds.SensorID,
		Dataset_id: ds.ID,
		Data:       dataPoints,
	}
	payload, err := json.Marshal(dataOut)
	if err != nil {
		log.Printf("unable to create simulator data object: %v", err)
	}
	req, err := http.NewRequest("POST", sensorURL, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Failed creating request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Failed performing request: %v", err)
	}
	res.Body.Close()
}
