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
	prefinedMetrics []preDfinedMetric
	db              *sqlx.DB
}

func NewCollection(list []preDfinedMetric, db *sqlx.DB) *iotCollector {
	var md []*prometheus.Desc
	for _, metric := range list {
		newDesc := prometheus.NewDesc(metric.Name,
			metric.Help,
			nil, nil,
		)
		md = append(md, newDesc)
	}

	return &iotCollector{
		metrics:         md,
		prefinedMetrics: list,
		db:              db,
	}
}

// Describe ...
func (c *iotCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
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

	for _, metric := range c.prefinedMetrics {
		fmt.Println("=>", metric.Name)

		newMetric := prometheus.MustNewConstMetric(prometheus.NewDesc(metric.Name,
			metric.Help,
			nil, nil,
		), prometheus.GaugeValue, rand.Float64()) // inject real data

		newMetric = prometheus.NewMetricWithTimestamp(time.Now(), newMetric)

		ch <- newMetric
	}
}
