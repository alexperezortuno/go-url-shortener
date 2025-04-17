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

var (
	instance *MetricsService
	once     sync.Once
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
func InitializeMetricsService(address string, globalTags []string) *MetricsService {
	once.Do(func() {
		client, err := statsd.New(address)
		if err != nil {
			log.Fatalf("Failed to initialize metrics client: %v", err)
		}
		client.Tags = globalTags
		instance = &MetricsService{client: client}
	})
	return instance
}

// GetInstance returns the singleton instance of MetricsService
func GetInstance() *MetricsService {
	if instance == nil {
		log.Fatal("MetricsService is not initialized. Call InitializeMetricsService first.")
	}
	return instance
}

// IncrementCounter increments a counter metric
func (m *MetricsService) IncrementCounter(name string, tags []string) {
	err := m.client.Incr(name, tags, 1)
	if err != nil {
		log.Printf("Failed to send counter metric: %v", err)
	}
}

// RecordHistogram records a histogram metric
func (m *MetricsService) RecordHistogram(name string, value float64, tags []string) {
	err := m.client.Histogram(name, value, tags, 1)
	if err != nil {
		log.Printf("Failed to send histogram metric: %v", err)
	}
}

// RecordTiming records a timing metric
func (m *MetricsService) RecordTiming(name string, duration float64, tags []string) {
	err := m.client.Timing(name, time.Duration(duration), tags, 1)
	if err != nil {
		log.Printf("Failed to send timing metric: %v", err)
	}
}

func (m *MetricsService) MonitorPerformance() {
	// Increment a counter for monitoring
	err := m.client.Incr("metrics_service.metrics_sent", nil, 1)
	if err != nil {
		log.Printf("Failed to send self-monitoring metric: %v", err)
	}

	// Measure latency of sending metrics
	start := time.Now()
	// Simulate sending a metric
	time.Sleep(10 * time.Millisecond) // Replace with actual metric sending logic
	err = m.client.Timing("metrics_service.latency", time.Since(start), nil, 1)
	if err != nil {
		log.Printf("Failed to send latency metric: %v", err)
	}
}

func (m *MetricsService) HealthCheck() bool {
	err := m.client.Incr("metrics_service.health_check", nil, 1)
	return err == nil
}
