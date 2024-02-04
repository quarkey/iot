package models

import (
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

type iotCollector struct {
	metrics         []*prometheus.Desc
	prefinedMetrics map[string][]dsetMetric
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

func NewCollectionDataset(dataMap map[string][]dsetMetric) *iotCollector {
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
// It is called when prometheus scrapes the metrics endpoint.
// This allows us to inject our own data into the prometheus metrics upon creation.
func (c *iotCollector) Collect(ch chan<- prometheus.Metric) {
	for k, v := range c.prefinedMetrics {
		for _, metric := range v {
			newMetric := prometheus.MustNewConstMetric(prometheus.NewDesc(k,
				metric.Help,
				nil, nil,
			), prometheus.GaugeValue, rand.Float64())
			newMetric = prometheus.NewMetricWithTimestamp(time.Now(), newMetric)
			ch <- newMetric
		}
	}
}
