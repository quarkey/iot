package models

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type iotCollector struct {
	metrics []*prometheus.Desc
}

// NewCollection creates a new iotCollector
func NewCollection(list []string) *iotCollector {
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

// Describe ...
func (c *iotCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		fmt.Println("Setting up metric:", metric.String())
		ch <- metric
	}
}

// Collect...
func (c *iotCollector) Collect(ch chan<- prometheus.Metric) {
	fmt.Println("Custom Collect")
	metricValue := rand.Float64()
	for _, metric := range c.metrics {
		newMetric := prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, metricValue)
		newMetric = prometheus.NewMetricWithTimestamp(time.Now(), newMetric)
		ch <- newMetric
	}
}
