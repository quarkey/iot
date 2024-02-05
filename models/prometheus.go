package models

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

// iotCollector is a type used for prometheus collector implementation.
type iotCollector struct {
	metrics []*prometheus.Desc
	pdm     []PreDefinedMetric
	db      *sqlx.DB
}

// PreDefinedMetric.
type PreDefinedMetric struct {
	ID        int
	Reference string // dataset reference
	Name      string
	Help      string
	Value     float64
}

// NewCollector creates a new instance of iotCollector with a database connection.
func NewCollection(list []PreDefinedMetric, db *sqlx.DB) *iotCollector {
	var md []*prometheus.Desc
	for _, metric := range list {
		newDesc := prometheus.NewDesc(metric.Name,
			metric.Help,
			nil, nil,
		)
		md = append(md, newDesc)
	}

	return &iotCollector{
		metrics: md,
		pdm:     list,
		db:      db,
	}
}

// Describe implements the prometheus.Collector interface.
func (c *iotCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric
	}
}

// Collect is a custom implementation of the prometheus.Collector interface.
// This methods runs when prometheus scrapes /metrics endpoint. This enables us to inject real data  into prometheus.
// A db.Ping() is done before pulling data to verify the db connection.
func (c *iotCollector) Collect(ch chan<- prometheus.Metric) {
	err := c.db.Ping()
	if err != nil {
		log.Printf("[ERROR] Collect() is unable to ping database: %v\n", err)
	}
	metrics, err := RetrieveDBMetricsToMap(c.db)
	if err != nil {
		log.Printf("[ERROR] Collect() is unable to get metrics from db: %v", err)
	}

	// Loop through predefined metrics map and inject real data into prometheus
	for _, metric := range c.pdm {

		var metricValue float64
		metricsFromDB, ok := metrics[metric.Name]
		if ok {
			metricValue = metricsFromDB.Value
		}
		// replace the value with the real data on the fly
		newDesc := prometheus.MustNewConstMetric(prometheus.NewDesc(metric.Name,
			metric.Help,
			nil, nil,
		), prometheus.GaugeValue, metricValue)

		newDesc = prometheus.NewMetricWithTimestamp(time.Now(), newDesc)

		ch <- newDesc
	}
}

// RegisterPrefinedMetrics registers predefined metrics into the database.
// Before registering, it clears all the existing metrics. If it fails to clear before registering,
// it returns an error.
// TODO: add insert to a transaction
func RegisterPrefinedMetrics(db *sqlx.DB, list []PreDefinedMetric) error {
	err := ClearRegisteredMetrics(db)
	if err != nil {
		return err
	}
	for _, metric := range list {
		_, err := db.Exec(`insert into iot.metrics (metric_id, name, help, value) values ($1, $2, $3, $4);`,
			metric.ID, metric.Name, metric.Help, metric.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClearRegisteredMetrics clears all the registered metrics in the database.
func ClearRegisteredMetrics(db *sqlx.DB) error {
	_, err := db.Exec(`delete from iot.metrics;`)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveDBMetricsToMap retrives stored metrics in the database into a map. The key is the metric name
func RetrieveDBMetricsToMap(db *sqlx.DB) (map[string]PreDefinedMetric, error) {
	metricsMap := make(map[string]PreDefinedMetric)
	rows, err := db.Queryx(`select metric_id, name, help, value from iot.metrics;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m PreDefinedMetric
		err = rows.Scan(&m.ID, &m.Name, &m.Help, &m.Value)
		if err != nil {
			return nil, err
		}
		metricsMap[m.Name] = m
	}
	return metricsMap, nil
}
