package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type iotCollector struct {
	metrics         []*prometheus.Desc
	prefinedMetrics map[string][]dsetMetric
	db              *sqlx.DB
}

// NewCollection creates a new iotCollector
func NewCollection(list []string, db *sqlx.DB) *iotCollector {
	var md []*prometheus.Desc
	for _, metric := range list {
		metric := prometheus.NewDesc(metric,
			"dynamically created metric",
			nil, nil,
		)
		md = append(md, metric)
	}
	return &iotCollector{
		metrics: md,
	}
}

func NewCollectionDataset(dataMap map[string][]dsetMetric, db *sqlx.DB) *iotCollector {
	var md []*prometheus.Desc
	for k, v := range dataMap {
		for _, metric := range v {
			newDesc := prometheus.NewDesc(k,
				metric.Help,
				nil, nil,
			)
			md = append(md, newDesc)
		}
	}
	return &iotCollector{
		metrics:         md,
		prefinedMetrics: dataMap,
		db:              db,
	}
}

// Describe ...
func (c *iotCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		// fmt.Println("Setting up metric:", metric.String())
		ch <- metric
	}
}

// Collect is a custom implementation of the prometheus.Collector interface.
// This methods is invoked when prometheus scrapes metrics endpoint.
// By making a custom implementation allows us to inject db into the prometheus collector.
func (c *iotCollector) Collect(ch chan<- prometheus.Metric) {
	// pull metrics information from database
	err := c.db.Ping()
	if err != nil {
		fmt.Printf("[ERROR] unable to ping database: %v\n", err)
	} else {
		fmt.Println("[INFO] database connection is OK")
	}
	for k, v := range c.prefinedMetrics {
		for _, metric := range v {
			fmt.Println("=>", k)
			newMetric := prometheus.MustNewConstMetric(prometheus.NewDesc(k,
				metric.Help,
				nil, nil,
			), prometheus.GaugeValue, rand.Float64())
			newMetric = prometheus.NewMetricWithTimestamp(time.Now(), newMetric)
			ch <- newMetric
		}
	}
}
