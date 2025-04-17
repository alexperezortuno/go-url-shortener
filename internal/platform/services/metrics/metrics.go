package metrics

import (
	"log"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

type MetricsService struct {
	client *statsd.Client
}

type MetricTag struct {
	Key   string
	Value string
}

type InstanceMetrics struct {
	Service *MetricsService
	Once    sync.Once
	Enabled bool
}

var (
	instance InstanceMetrics
)

// ConvertTags converts a slice of MetricTag to a slice of strings
func ConvertTags(tags []MetricTag) []string {
	var result []string
	for _, tag := range tags {
		result = append(result, tag.Key+":"+tag.Value)
	}
	return result
}

// InitializeMetricsService initializes the metrics service singleton
func InitializeMetricsService(address string, tags []string) *InstanceMetrics {
	instance.Once.Do(func() {
		client, err := statsd.New(address)
		if err != nil {
			log.Fatalf("Failed to initialize metrics client: %v", err)
		}
		client.Tags = tags
		instance.Service = &MetricsService{client: client}
		instance.Enabled = true
	})
	return &instance
}

// GetInstance returns the singleton instance of MetricsService
func GetInstance() *InstanceMetrics {
	if instance.Service == nil {
		log.Printf("MetricsService is not initialized. Call InitializeMetricsService first.")
		instance.Service = &MetricsService{}
		instance.Enabled = false
	}
	return &instance
}

// IncrementCounter increments a counter metric
func (m *InstanceMetrics) IncrementCounter(name string, tags []MetricTag) {
	if !m.Enabled {
		log.Printf("MetricsService is not enabled. Skipping counter increment.")
		return
	}
	err := m.Service.client.Incr(name, ConvertTags(tags), 1)
	if err != nil {
		log.Printf("Failed to send counter metric: %v", err)
	}
}

// RecordHistogram records a histogram metric
func (m *InstanceMetrics) RecordHistogram(name string, value float64, tags []MetricTag) {
	if !m.Enabled {
		log.Printf("MetricsService is not enabled. Skipping counter increment.")
		return
	}
	err := m.Service.client.Histogram(name, value, ConvertTags(tags), 1)
	if err != nil {
		log.Printf("Failed to send histogram metric: %v", err)
	}
}

// RecordTiming records a timing metric
func (m *InstanceMetrics) RecordTiming(name string, duration float64, tags []MetricTag) {
	if !m.Enabled {
		log.Printf("MetricsService is not enabled. Skipping counter increment.")
		return
	}
	err := m.Service.client.Timing(name, time.Duration(duration), ConvertTags(tags), 1)
	if err != nil {
		log.Printf("Failed to send timing metric: %v", err)
	}
}

func (m *InstanceMetrics) MonitorPerformance() {
	// Increment a counter for monitoring
	if !m.Enabled {
		log.Printf("MetricsService is not enabled. Skipping counter increment.")
		return
	}
	err := m.Service.client.Incr("metrics_service.metrics_sent", nil, 1)
	if err != nil {
		log.Printf("Failed to send self-monitoring metric: %v", err)
	}

	// Measure latency of sending metrics
	start := time.Now()
	// Simulate sending a metric
	time.Sleep(10 * time.Millisecond) // Replace with actual metric sending logic
	err = m.Service.client.Timing("metrics_service.latency", time.Since(start), nil, 1)
	if err != nil {
		log.Printf("Failed to send latency metric: %v", err)
	}
}

func (m *InstanceMetrics) HealthCheck() bool {
	if !m.Enabled {
		log.Printf("MetricsService is not enabled. Skipping counter increment.")
		return false
	}
	err := m.Service.client.Incr("metrics_service.health_check", nil, 1)
	return err == nil
}
