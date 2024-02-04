package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type iotCollector struct {
	metrics           []*prometheus.Desc
	predefinedMetrics []preDfinedMetric
	db                *sqlx.DB
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
		metrics:           md,
		predefinedMetrics: list,
		db:                db,
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

	// get metrics from database
	metrics, err := GetRegisteredMetricsFromDB(c.db)
	if err != nil {
		log.Printf("[ERROR] unable to get metrics from db: %v", err)
	}

	for _, metric := range c.predefinedMetrics {
		fmt.Println("=>", metric.Name)
		var metricValue float64
		existingMetric, ok := metrics[metric.Name]
		if ok {
			metricValue = existingMetric.Value
		}
		newMetric := prometheus.MustNewConstMetric(prometheus.NewDesc(metric.Name,
			metric.Help,
			nil, nil,
		), prometheus.GaugeValue, metricValue) // inject real data

		newMetric = prometheus.NewMetricWithTimestamp(time.Now(), newMetric)

		ch <- newMetric
	}
}

// RegisterPrefinedMetrics registers predefined metrics into the database.
// Before registering, it clears all the existing metrics. If it fails to clear before registering,
// it returns an error.
// TODO: add insert to a transaction
func RegisterPrefinedMetrics(db *sqlx.DB, list []preDfinedMetric) error {
	err := ClearRegisteredMetrics(db)
	if err != nil {
		return err
	}
	for _, metric := range list {
		fmt.Println("regging", metric.Name)
		_, err := db.Exec(`insert into iot.metrics (metric_id, name, help, value) values ($1, $2, $3, $4);`,
			metric.ID, metric.Name, metric.Help, metric.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
func ClearRegisteredMetrics(db *sqlx.DB) error {
	_, err := db.Exec(`delete from iot.metrics;`)
	if err != nil {
		return err
	}
	return nil
}

func GetRegisteredMetricsFromDB(db *sqlx.DB) (map[string]preDfinedMetric, error) {
	metricsMap := make(map[string]preDfinedMetric)
	rows, err := db.Queryx(`select metric_id, name, help, value from iot.metrics;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m preDfinedMetric
		err = rows.Scan(&m.ID, &m.Name, &m.Help, &m.Value)
		if err != nil {
			return nil, err
		}
		metricsMap[m.Name] = m
	}
	return metricsMap, nil
}
